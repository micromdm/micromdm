package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

type mainIOFunc func(t *testing.T, stdin, stdout, stderr *bytes.Buffer)
type signalFunc func(t *testing.T, pidfile string, signal os.Signal, d time.Duration)

func mainUsage(t *testing.T, stdin, stdout, stderr *bytes.Buffer) {
	output := stderr.String()
	words := []string{"USAGE", "pidfile"}
	for _, word := range words {
		if !strings.Contains(output, word) {
			t.Errorf("expected %q in output, got:\n%s", word, output)
		}
	}
}

func checkRestarted(t *testing.T, stdin, stdout, stderr *bytes.Buffer) {
	output := stderr.String()
	logsub := `level=info msg="restarting process" signal=hangup`
	if !strings.Contains(output, logsub) {
		t.Errorf("want %q in output, got:\n%s", logsub, output)
	}
}

func checkExitAfterSignal(t *testing.T, stdin, stdout, stderr *bytes.Buffer) {
	output := stderr.String()
	logsub := `level=info exit="received signal`
	if !strings.Contains(output, logsub) {
		t.Errorf("want %q in output, got:\n%s", logsub, output)
	}
}

func checkLogSwap(t *testing.T, stdin, stdout, stderr *bytes.Buffer) {
	output := stderr.String()
	debugsub := `level=info msg="swapping level" debug=true`
	if !strings.Contains(output, debugsub) {
		t.Errorf("want %q in output, got:\n%s", debugsub, output)
	}
}

func logswap() signalFunc {
	return func(t *testing.T, pidfile string, _ os.Signal, d time.Duration) {
		signalAfter(t, pidfile, defaultSwapSignal, 50*time.Millisecond)
		signalAfter(t, pidfile, syscall.SIGTERM, 100*time.Millisecond)
	}
}

func sighup() signalFunc {
	return func(t *testing.T, pidfile string, _ os.Signal, d time.Duration) {
		signalAfter(t, pidfile, syscall.SIGHUP, 50*time.Millisecond)
		signalAfter(t, pidfile, syscall.SIGTERM, 100*time.Millisecond)
	}
}

func exitWith(s os.Signal) signalFunc {
	return func(t *testing.T, pidfile string, _ os.Signal, d time.Duration) {
		signalAfter(t, pidfile, s, 50*time.Millisecond)
	}
}

func signalAfter(t *testing.T, pidfile string, s os.Signal, d time.Duration) {
	t.Helper()

	// sleep until pidfile has a chance to exist.
	time.Sleep(10 * time.Millisecond)

	pids, err := ioutil.ReadFile(pidfile)
	if err != nil {
		t.Fatal(err)
	}

	pid, err := strconv.Atoi(string(pids))
	if err != nil {
		t.Fatal(err)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-time.After(d):
		if err := proc.Signal(s); err != nil {
			t.Fatal(err)
		}
		return
	}
}

func TestMain(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "micromdm-test-main")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		stdin  *bytes.Buffer
		stdout *bytes.Buffer
		stderr *bytes.Buffer
		args   []string
		exit   int
		check  mainIOFunc
		signal signalFunc

		// run test synchronously
		// for cases where parallel test is not possible
		synchronous bool
	}{
		{
			name:   "short usage",
			stdin:  new(bytes.Buffer),
			stdout: new(bytes.Buffer),
			stderr: new(bytes.Buffer),
			args:   []string{"micromdm", "-h"},
			exit:   2,
			check:  mainUsage,
		},
		{
			name:   "long usage",
			stdin:  new(bytes.Buffer),
			stdout: new(bytes.Buffer),
			stderr: new(bytes.Buffer),
			args:   []string{"micromdm", "-help"},
			exit:   2,
			check:  mainUsage,
		},
		{
			name:   "subcommand help usage",
			stdin:  new(bytes.Buffer),
			stdout: new(bytes.Buffer),
			stderr: new(bytes.Buffer),
			args:   []string{"micromdm", "help"},
			exit:   2,
			check:  mainUsage,
		},
		{
			name:   "exit on signal",
			stdin:  new(bytes.Buffer),
			stdout: new(bytes.Buffer),
			stderr: new(bytes.Buffer),
			args: []string{"micromdm", "-pidfile", filepath.Join(tmpdir, "sigint.pid"),
				"-database_url", filepath.Join(tmpdir, "micromdm.db")},
			exit:   1,
			check:  checkExitAfterSignal,
			signal: exitWith(syscall.SIGTERM),
		},
		{
			name:   "swap logs",
			stdin:  new(bytes.Buffer),
			stdout: new(bytes.Buffer),
			stderr: new(bytes.Buffer),
			args: []string{"micromdm", "-pidfile", filepath.Join(tmpdir, "sigusr.pid"),
				"-database_url", filepath.Join(tmpdir, "micromdm.db")},
			exit:   1,
			check:  checkLogSwap,
			signal: logswap(),

			// run this test synchronously to avoid the data race
			// bytes.Buffer is not thread safe, or is there an actual bug here?
			synchronous: true,
		},
		{
			name:   "restart on hup",
			stdin:  new(bytes.Buffer),
			stdout: new(bytes.Buffer),
			stderr: new(bytes.Buffer),
			args: []string{"micromdm", "-pidfile", filepath.Join(tmpdir, "sighup.pid"),
				"-database_url", filepath.Join(tmpdir, "micromdm.db")},
			exit:        1,
			check:       checkRestarted,
			signal:      sighup(),
			synchronous: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if !tt.synchronous {
				t.Parallel()
			}

			if tt.signal != nil {
				if runtime.GOOS == "windows" {
					t.Skip("TODO(issues/686): signal handling on windows is not done.")
				}

				// the signal and timing arguments are defaults, specified by middleware instead.
				go tt.signal(t, tt.args[2], syscall.SIGTERM, 50*time.Millisecond)
			}

			if got, want := micromdm(tt.args, tt.stdin, tt.stdout, tt.stderr), tt.exit; got != want {
				t.Fatalf("exit code: got %d, want %d", got, want)
			}

			tt.check(t, tt.stdin, tt.stdout, tt.stderr)
		})
	}
}
