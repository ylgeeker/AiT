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
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
)

type slot struct {
	gen     int
	spec    Unit
	cancel  context.CancelFunc
	done    chan struct{}
	state   string
	running bool
	usage   Usage
}

type daemon struct {
	root   string
	parent context.Context
	sup    *Supervisor

	mu    sync.Mutex
	slots map[string]*slot
}

func (d *daemon) acceptAdmin(ctx context.Context, ln net.Listener) {
	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go serveConn(conn, d.handle)
	}
}

func (d *daemon) handle(req Request) Response {
	switch req.Op {
	case opStatus:
		return Response{OK: true, Status: d.statusText()}
	case opReload:
		return d.reload()
	case opStart:
		return d.startOne(req.Name)
	case opStop:
		return d.stopOne(req.Name)
	default:
		return Response{Error: "unknown op"}
	}
}

func (d *daemon) statusText() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	names := make([]string, 0, len(d.slots))
	for name := range d.slots {
		names = append(names, name)
	}
	sort.Strings(names)
	var b strings.Builder
	for _, name := range names {
		sl := d.slots[name]
		fmt.Fprint(&b, formatStatusLine(name, sl.state, sl.spec, sl.usage))
	}
	return b.String()
}

func (d *daemon) reload() Response {
	cfgPath, _, _ := Layout(d.root)
	doc, err := Load(cfgPath)
	if err != nil {
		return Response{Error: err.Error()}
	}
	d.reconcile(doc.Units)
	return Response{OK: true}
}

func (d *daemon) startOne(name string) Response {
	cfgPath, _, _ := Layout(d.root)
	doc, err := Load(cfgPath)
	if err != nil {
		return Response{Error: err.Error()}
	}
	for _, unit := range doc.Units {
		if unit.Name != name {
			continue
		}
		d.launch(unit)
		return Response{OK: true}
	}
	return Response{Error: ErrUnitMissing.Error()}
}

func (d *daemon) stopOne(name string) Response {
	d.mu.Lock()
	sl := d.slots[name]
	if sl == nil {
		d.mu.Unlock()
		return Response{Error: ErrUnitMissing.Error()}
	}
	sl.gen++
	cancel := sl.cancel
	sl.cancel = nil
	sl.running = false
	sl.state = stateStopped
	d.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	return Response{OK: true}
}

func (d *daemon) noteSample(name string, usage Usage) {
	d.mu.Lock()
	defer d.mu.Unlock()
	sl := d.slots[name]
	if sl == nil || !sl.running {
		return
	}
	sl.usage = usage
}

func (d *daemon) stopAll() {
	d.mu.Lock()
	names := make([]string, 0, len(d.slots))
	for name := range d.slots {
		names = append(names, name)
	}
	d.mu.Unlock()
	for _, name := range names {
		_ = d.stopOne(name)
	}
}

func (d *daemon) reconcile(units []Unit) {
	want := map[string]Unit{}
	for _, unit := range units {
		want[unit.Name] = unit
	}
	d.mu.Lock()
	var stop []string
	var start []Unit
	for name, sl := range d.slots {
		next, ok := want[name]
		if !ok && sl.running {
			stop = append(stop, name)
			continue
		}
		if !ok {
			delete(d.slots, name)
			continue
		}
		delete(want, name)
		if sameUnit(sl.spec, next) {
			continue
		}
		start = append(start, next)
	}
	for _, unit := range want {
		start = append(start, unit)
	}
	d.mu.Unlock()
	for _, name := range stop {
		_ = d.stopOne(name)
	}
	for _, unit := range start {
		d.launch(unit)
	}
}

func (d *daemon) launch(unit Unit) {
	d.mu.Lock()
	sl := d.slots[unit.Name]
	if sl == nil {
		sl = &slot{}
		d.slots[unit.Name] = sl
	}
	sl.gen++
	gen := sl.gen
	if sl.cancel != nil {
		sl.cancel()
	}
	prevDone := sl.done
	sl.cancel = nil
	sl.running = false
	d.mu.Unlock()
	if prevDone != nil {
		<-prevDone
	}
	if d.parent == nil || d.parent.Err() != nil {
		return
	}
	ctx, cancel := context.WithCancel(d.parent)
	done := make(chan struct{})
	d.mu.Lock()
	sl = d.slots[unit.Name]
	if sl == nil || sl.gen != gen {
		d.mu.Unlock()
		cancel()
		return
	}
	sl.spec = unit
	sl.cancel = cancel
	sl.done = done
	sl.state = stateRunning
	sl.running = true
	d.mu.Unlock()
	go func() {
		defer close(done)
		result, _ := d.sup.Watch(ctx, unit)
		d.mu.Lock()
		defer d.mu.Unlock()
		cur := d.slots[unit.Name]
		if cur == nil || cur.gen != gen {
			return
		}
		cur.running = false
		cur.cancel = nil
		if result.Failed {
			cur.state = stateFailed
			return
		}
		cur.state = stateStopped
	}()
}
