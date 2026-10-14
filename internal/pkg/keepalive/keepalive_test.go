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
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPing_PongAndNotFound(t *testing.T) {
	h := PingHandler()
	ok := httptest.NewRecorder()
	servePing(h, ok, http.MethodGet, PingPath)
	if ok.Code != http.StatusOK || ok.Body.String() != PongBody {
		t.Fatalf("ping=%d %q", ok.Code, ok.Body.String())
	}
	miss := httptest.NewRecorder()
	servePing(h, miss, http.MethodGet, "/status")
	if miss.Code != http.StatusNotFound {
		t.Fatalf("other=%d", miss.Code)
	}
	post := httptest.NewRecorder()
	servePing(h, post, http.MethodPost, PingPath)
	if post.Code != http.StatusNotFound {
		t.Fatalf("post=%d", post.Code)
	}
}

func servePing(h http.Handler, w http.ResponseWriter, method, path string) {
	req := httptest.NewRequest(method, path, nil)
	h.ServeHTTP(w, req)
}

func TestDecide_ThrottleThenKill(t *testing.T) {
	limits := Limits{CPUPercent: 10}
	usage := Usage{CPUPercent: 80}
	if got := Decide(usage, limits, HardCaps{}, Hold{Duty: time.Second}, Sustain); got != ActionThrottle {
		t.Fatalf("early=%v", got)
	}
	if got := Decide(usage, limits, HardCaps{}, Hold{Duty: Sustain}, Sustain); got != ActionKill {
		t.Fatalf("sustain=%v", got)
	}
	if got := Decide(usage, limits, HardCaps{CPU: true}, Hold{Duty: Sustain}, Sustain); got != ActionNone {
		t.Fatalf("kernel cap=%v", got)
	}
	pause := DutyPause(80, 10, DutyWindow)
	if pause <= 0 || pause > DutyWindow {
		t.Fatalf("pause=%s", pause)
	}
}

func TestDecide_DiskFail(t *testing.T) {
	got := Decide(Usage{DiskSpaceBytes: 20}, Limits{DiskSpaceBytes: 10}, HardCaps{}, Hold{}, Sustain)
	if got != ActionDiskFail {
		t.Fatalf("disk=%v", got)
	}
}

func TestWatch_AbnormalRestartsCleanExitDoesNot(t *testing.T) {
	codes := []int{1, 0}
	var i int
	sup := &Supervisor{
		Sleep: func(time.Duration) {},
		Now:   time.Now,
		Run: func(ctx context.Context, unit Unit) (int, error) {
			code := codes[i]
			i++
			return code, nil
		},
	}
	res, err := sup.Watch(context.Background(), Unit{Name: "demo", Command: []string{"true"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.Starts != 2 || res.Failed || len(res.Sleeps) != 1 || res.Sleeps[0] != BackoffMin {
		t.Fatalf("result=%+v", res)
	}
}

func TestWatch_BudgetAndStableReset(t *testing.T) {
	var calls int
	now := time.Unix(1_700_000_000, 0)
	sup := &Supervisor{
		Sleep: func(time.Duration) {},
		Now:   func() time.Time { return now },
		Run: func(ctx context.Context, unit Unit) (int, error) {
			calls++
			now = now.Add(time.Millisecond)
			return 1, nil
		},
	}
	res, err := sup.Watch(context.Background(), Unit{Name: "demo", Command: []string{"false"}})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Failed || res.Starts != RestartBudget+1 {
		t.Fatalf("budget result=%+v", res)
	}
	if res.Sleeps[0] != BackoffMin || res.Sleeps[len(res.Sleeps)-1] > BackoffMax {
		t.Fatalf("sleeps=%v", res.Sleeps)
	}

	calls = 0
	now = time.Unix(0, 0)
	sup.Run = func(ctx context.Context, unit Unit) (int, error) {
		calls++
		now = now.Add(StableRun)
		if calls == 1 {
			return 1, nil
		}
		return 0, nil
	}
	res, err = sup.Watch(context.Background(), Unit{Name: "demo", Command: []string{"false"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.Failed || len(res.Sleeps) != 1 || res.Sleeps[0] != BackoffMin {
		t.Fatalf("reset result=%+v", res)
	}
}

func TestWatch_DiskStillOverDoesNotRestart(t *testing.T) {
	var starts int
	sup := &Supervisor{
		Sleep:    func(time.Duration) {},
		Now:      time.Now,
		DiskUsed: func(string) (uint64, error) { return 100, nil },
		Run: func(ctx context.Context, unit Unit) (int, error) {
			starts++
			return 1, nil
		},
	}
	res, err := sup.Watch(context.Background(), Unit{
		Name: "demo", Command: []string{"false"}, Cwd: "/tmp",
		Limits: Limits{DiskSpaceBytes: 10},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Failed || starts != 0 {
		t.Fatalf("result=%+v starts=%d", res, starts)
	}
}

func TestWatch_DiskPercentStillOverDoesNotRestart(t *testing.T) {
	var starts int
	sup := &Supervisor{
		Sleep:  func(time.Duration) {},
		Now:    time.Now,
		Volume: func(string) (Volume, error) { return Volume{UsedPercent: 95}, nil },
		Run: func(ctx context.Context, unit Unit) (int, error) {
			starts++
			return 1, nil
		},
	}
	res, err := sup.Watch(context.Background(), Unit{
		Name: "demo", Command: []string{"false"}, Cwd: t.TempDir(),
		Limits: Limits{DiskSpacePercent: 90},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Failed || starts != 0 {
		t.Fatalf("result=%+v starts=%d", res, starts)
	}
}

func TestWriteCPUMax_FakeSysfs(t *testing.T) {
	dir := t.TempDir()
	if err := WriteCPUMax(dir, 50); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "cpu.max"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "50000 100000\n" {
		t.Fatalf("cpu.max=%q", data)
	}
}

func TestAcquire_SecondFails(t *testing.T) {
	dir := t.TempDir()
	_, _, lock := Layout(dir)
	release, err := Acquire(lock)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if _, err := Acquire(lock); err != ErrAlreadyRunning {
		t.Fatalf("second=%v", err)
	}
}

func TestDirUsage_SkipsSymlink(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a"), []byte("abcd"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "a"), filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}
	got, err := DirUsage(dir, DiskWalkBudget)
	if err != nil {
		t.Fatal(err)
	}
	if got != 4 {
		t.Fatalf("usage=%d", got)
	}
}

func TestUpsert_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ConfigFileName)
	doc, err := Upsert(File{}, Unit{Name: "demo", Command: []string{"sleep", "1"}})
	if err != nil {
		t.Fatal(err)
	}
	if err := Save(path, doc); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(path)
	if err != nil || len(loaded.Units) != 1 || loaded.Units[0].Name != "demo" {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
	if _, err := Upsert(loaded, Unit{Name: "bad/name", Command: []string{"x"}}); err != ErrBadName {
		t.Fatalf("name=%v", err)
	}
}
