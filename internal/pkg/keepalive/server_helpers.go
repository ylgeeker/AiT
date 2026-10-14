/**
 * Copyright 2026 AiT authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/

package keepalive

import (
	"context"
	"errors"
	"net"
	"net/http"
	"reflect"
)

// Serve runs the ping listener, the admin endpoint, and every unit in keepalive.yaml.
// root is AppRoot. listen overrides the file address when non-empty. ctx cancel stops children
// as an operator stop.
// Returns ErrAlreadyRunning when another serve holds the lock, a listen or admin bind error,
// a config load error, or a non-ErrServerClosed error from http.Server.Serve.
func Serve(ctx context.Context, root, listen string) error {
	cfgPath, _, lockPath := Layout(root)
	release, err := Acquire(lockPath)
	if err != nil {
		return err
	}
	defer release()

	doc, err := Load(cfgPath)
	if err != nil {
		return err
	}
	addr := listen
	if addr == "" {
		addr = doc.Listen
	}
	if addr == "" {
		addr = DefaultListen
	}

	admin, err := listenAdmin(root)
	if err != nil {
		return err
	}
	defer admin.Close()

	httpLn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer httpLn.Close()

	srv := &daemon{root: root, parent: ctx, sup: &Supervisor{}, slots: map[string]*slot{}}
	srv.sup.OnSample = srv.noteSample
	srv.reconcile(doc.Units)
	httpSrv := &http.Server{Handler: PingHandler()}
	serveErr := make(chan error, 1)
	go func() { serveErr <- httpSrv.Serve(httpLn) }()
	go srv.acceptAdmin(ctx, admin)

	var runErr error
	select {
	case <-ctx.Done():
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			runErr = err
		}
	}
	srv.stopAll()
	shutCtx, cancel := context.WithTimeout(context.Background(), StopGrace)
	defer cancel()
	_ = httpSrv.Shutdown(shutCtx)
	return runErr
}

func sameUnit(a, b Unit) bool {
	return reflect.DeepEqual(a, b)
}

// NotifyReload asks a running serve to reread keepalive.yaml.
// root is AppRoot. Returns ErrDaemonDown when serve is not running.
func NotifyReload(root string) error {
	resp, err := Call(root, Request{Op: opReload})
	if err != nil {
		return err
	}
	if !resp.OK {
		return errors.New(resp.Error)
	}
	return nil
}

// QueryStatus returns the live unit lines from serve.
// root is AppRoot. Returns ErrDaemonDown when serve is not running.
func QueryStatus(root string) (string, error) {
	resp, err := Call(root, Request{Op: opStatus})
	if err != nil {
		return "", err
	}
	if !resp.OK {
		return "", errors.New(resp.Error)
	}
	return resp.Status, nil
}

// NotifyStart asks serve to run one unit.
// root is AppRoot and name is the unit name. Returns ErrDaemonDown when serve is not running.
func NotifyStart(root, name string) error {
	return notifyName(root, opStart, name)
}

// NotifyStop asks serve to stop one unit without restarting it.
// root is AppRoot and name is the unit name. Returns ErrDaemonDown when serve is not running.
func NotifyStop(root, name string) error {
	return notifyName(root, opStop, name)
}

func notifyName(root, op, name string) error {
	resp, err := Call(root, Request{Op: op, Name: name})
	if err != nil {
		return err
	}
	if !resp.OK {
		return errors.New(resp.Error)
	}
	return nil
}
