package runner

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/fd/go-shellwords/shellwords"
)

var _ Execer = (*Docker)(nil)

type Docker struct {
}

type dockerProcess struct {
	name   string
	cmd    *exec.Cmd
	mtx    sync.Mutex
	output []*OutputLine
	w      io.Writer
	r      io.Reader
}

func (*Docker) Exec(wdir, name, impl, command string) (Process, error) {
	proc := &dockerProcess{name: name}

	cmdArgs, err := shellwords.Split(command)
	if err != nil {
		return nil, err
	}

	args := []string{
		"run",
		// --net container:th-c-test
		"--name", "th-interop-" + name,
		"--volume", wdir + ":/shared",
		"--rm",
		"-it",
		"telehash/test-" + impl,
	}

	args = append(args, cmdArgs...)

	cmd := exec.Command("docker", args...)
	cmd.Dir = os.TempDir()

	w, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if p, err := cmd.StderrPipe(); err == nil {
		go proc.readLines("stderr", p)
	} else {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	proc.cmd = cmd
	proc.w = w
	proc.r = r
	return proc, nil
}

func (p *dockerProcess) Output() []*OutputLine {
	return p.output
}

func (p *dockerProcess) OutputReader() io.Reader {
	return p.r
}

func (p *dockerProcess) Kill() error {
	p.cmd.Process.Kill()
	return exec.Command("docker", "rm", "-f", "th-interop-"+p.name).Run()
}

func (p *dockerProcess) Wait() error {
	return p.cmd.Wait()
}

func (p *dockerProcess) readLines(stream string, r io.Reader) {
	buf := bufio.NewReader(r)

	var (
		line *OutputLine
	)

	for {
		if line == nil {
			line = &OutputLine{Process: p.name, Stream: stream, Time: time.Now()}

			p.mtx.Lock()
			p.output = append(p.output, line)
			p.mtx.Unlock()
		}

		chunk, prefix, err := buf.ReadLine()
		line.Line += string(chunk)

		if !prefix {
			line = nil
		}

		if err != nil {
			return
		}
	}
}
