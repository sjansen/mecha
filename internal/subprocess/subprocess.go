package subprocess

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Config struct {
	err    error
	cmd    *exec.Cmd
	env    []string
	stdout <-chan string
	stderr <-chan string
}

type ExitStatus struct {
	Code  int
	Error error
}

type Subprocess struct {
	PID    int
	Status <-chan *ExitStatus
	Stdout <-chan string
	Stderr <-chan string
}

func makeLineChannel(r io.Reader) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		reader := bufio.NewReader(r)
		for {
			line, err := reader.ReadString('\n')
			n := len(line) - 1
			if n >= 0 && line[n] == '\n' {
				if n >= 1 && line[n-1] == '\r' {
					line = line[:n-1]
				} else {
					line = line[:n]
				}
			}
			switch {
			case err == nil:
				ch <- line
			case err == io.EOF:
				ch <- line
				return
			default:
				return
			}
		}
	}()
	return ch
}

func New(ctx context.Context, name string, args ...string) *Config {
	return &Config{
		cmd: exec.CommandContext(ctx, name, args...),
	}
}

func (cfg *Config) CaptureStderrLines() *Config {
	if cfg.err != nil {
		return cfg
	}

	stderr, err := cfg.cmd.StderrPipe()
	if err != nil {
		cfg.err = err
		return cfg
	}
	cfg.stderr = makeLineChannel(stderr)
	return cfg
}

func (cfg *Config) CaptureStdoutLines() *Config {
	if cfg.err != nil {
		return cfg
	}

	stdout, err := cfg.cmd.StdoutPipe()
	if err != nil {
		cfg.err = err
		return cfg
	}
	cfg.stdout = makeLineChannel(stdout)
	return cfg
}

func (cfg *Config) SetEnv(name, value string) *Config {
	if cfg.err != nil {
		return cfg
	}
	if cfg.env == nil {
		cfg.env = os.Environ()
	}

	cfg.env = append(cfg.env, name+"="+value)
	return cfg
}

func (cfg *Config) UnsetEnv(name string) *Config {
	if cfg.err != nil {
		return cfg
	}
	if cfg.env == nil {
		cfg.env = os.Environ()
	}

	prefix := name + "="
	n := 0
	for i, s := range cfg.env {
		if !strings.HasPrefix(s, prefix) {
			if i != n {
				cfg.env[n] = s
			}
			n++
		}
	}
	cfg.env = cfg.env[:n]
	return cfg
}

func (cfg *Config) Start() (*Subprocess, error) {
	if cfg.err != nil {
		return nil, cfg.err
	}

	cmd := cfg.cmd
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	es := make(chan *ExitStatus)
	go func() {
		defer close(es)
		err := cmd.Wait()
		if err == nil {
			es <- &ExitStatus{}
		} else {
			es <- &ExitStatus{
				Code:  cmd.ProcessState.ExitCode(),
				Error: err,
			}
		}
	}()

	p := &Subprocess{
		PID:    cmd.Process.Pid,
		Stdout: cfg.stdout,
		Stderr: cfg.stderr,
		Status: es,
	}

	return p, nil
}
