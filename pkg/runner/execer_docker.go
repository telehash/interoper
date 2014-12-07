package runner

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/fd/go-shellwords/shellwords"

	"github.com/telehash/interoper/pkg/test"
)

var _ Execer = (*Docker)(nil)

type Docker struct {
}

type dockerProcess struct {
	wgOut  sync.WaitGroup
	wgExit sync.WaitGroup
	name   string
	id     string

	cOut  chan Event
	cKill chan struct{}
	cDone chan struct{}
}

func (d *Docker) Exec(wdir, name string, desc *test.Process) (out <-chan Event, kill chan struct{}, done <-chan struct{}) {
	proc := &dockerProcess{name: name, id: randomContainerName()}
	proc.cOut = make(chan Event)
	proc.cKill = make(chan struct{})
	proc.cDone = make(chan struct{})

	proc.wgOut.Add(1)
	proc.wgExit.Add(1)
	go proc.runContainer(wdir, name, desc)
	go proc.outCloser()
	go proc.doneCloser()

	return proc.cOut, proc.cKill, proc.cDone
}

func (p *dockerProcess) runContainer(wdir, name string, desc *test.Process) {
	defer p.wgExit.Done()
	defer p.wgOut.Done()

	defer func() {
		defer func() { recover() }()
		close(p.cKill)
	}()

	{ // start the cleaner
		p.wgOut.Add(1)
		go p.cleanContainer()
	}

	var cmd *exec.Cmd

	{ // make the command
		cmdArgs, err := shellwords.Split(desc.Command)
		if err != nil {
			p.traceError(err)
			p.traceExit(1)
			return
		}

		cmdArgs = append([]string{
			"run",
			"--net", "host",
			"--name", p.id,
			"--volume", wdir + ":/shared",
			desc.Image,
		}, cmdArgs...)

		cmd = exec.Command("docker", cmdArgs...)
		cmd.Dir = os.TempDir()
	}

	// setup I/O
	var (
		err     error
		rSTDOUT io.ReadCloser
		rSTDERR io.ReadCloser
	)

	{ // setup STDOUT
		rSTDOUT, err = cmd.StdoutPipe()
		if err != nil {
			p.traceError(err)
			p.traceExit(1)
			return
		}
		defer rSTDOUT.Close()
	}

	{ // setup STDERR
		rSTDERR, err = cmd.StderrPipe()
		if err != nil {
			p.traceError(err)
			p.traceExit(1)
			return
		}
		defer rSTDERR.Close()
	}

	p.wgOut.Add(1)
	go p.readOutCommands(rSTDOUT)
	p.wgOut.Add(1)
	go p.readOutCommands(rSTDERR)

	p.wgExit.Add(1)
	p.wgOut.Add(1)
	go p.waitForKill()

	if err = cmd.Run(); err != nil {
		p.traceError(err)
	}
	p.traceExit(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
}

func (p *dockerProcess) cleanContainer() {
	defer p.wgOut.Done()

	// wait for exit
	p.wgExit.Wait()

	{ // remove the container
		err := exec.Command("docker", "rm", "-f", p.id).Run()
		if err != nil {
			p.traceError(err)
		}
	}
}

func (p *dockerProcess) waitForKill() {
	defer p.wgExit.Done()
	defer p.wgOut.Done()

	<-p.cKill

	{ // kill the container
		err := exec.Command("docker", "kill", p.id).Run()
		if err != nil {
			p.traceError(err)
		}
	}
}

func (p *dockerProcess) readOutCommands(r io.Reader) {
	defer p.wgOut.Done()

	var (
		buf  = bufio.NewReader(r)
		line string
	)

	for {
		chunk, prefix, err := buf.ReadLine()
		if err == io.EOF || err == io.ErrClosedPipe {
			break
		}
		if err != nil {
			p.traceError(err)
			return
		}

		line += string(chunk)

		if !prefix {
			event := Event{}
			if strings.HasPrefix(line, "{\"") {
				err = json.Unmarshal([]byte(line), &event)
				if err != nil {
					panic(err)
				}
				event.Role = p.name
				line = ""

			} else {
				event.Type = "log"
				event.Role = p.name
				event.RecordTime()
				event.Info = Info{"line": line}
				line = ""

			}
			p.cOut <- event
		}
	}
}

func (p *dockerProcess) outCloser() {
	p.wgOut.Wait()
	close(p.cOut)
}

func (p *dockerProcess) doneCloser() {
	p.wgExit.Wait()
	close(p.cDone)
}

func (p *dockerProcess) traceError(err error) {
	if err != nil {
		e := Event{Type: "log", Info: Info{"line": "error: " + err.Error()}}
		e.RecordTime()
		e.Role = p.name
		p.cOut <- e
	}
}

func (p *dockerProcess) traceExit(code int) {
	e := Event{Type: "exited", Info: Info{"exit_code": code}}
	e.RecordTime()
	e.Role = p.name
	p.cOut <- e
}

func randomContainerName() string {
	var buf [5]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("th-interop-%x", buf[:])
}
