// © 2012 Steve McCoy. Available under the MIT license.

/*
Package pipeline provides a datatype for simple chaining of
the standard os/exec package's *Cmd.
*/
package pipeline

import (
	"errors"
	"io"
	"os/exec"
)

type P []*exec.Cmd

// New returns a new Chain. If the Cmds' Stdins and Stdouts cannot
// be piped together for any reason, nil and the error are returned.
//
// This function also returns an error if len(cmds) == 0.
func New(cmds ...*exec.Cmd) (P, error) {
	if len(cmds) == 0 {
		return nil, errors.New("pipeline.New() requires at least one command")
	}

	outs := make([]io.ReadCloser, 0, len(cmds)-1)
	for i := 1; i < len(cmds); i++ {
		out, err := cmds[i-1].StdoutPipe()
		if err != nil {
			for _, o := range outs {
				o.Close()
			}
			return nil, err
		}
		outs = append(outs, out)
		cmds[i].Stdin = out
	}

	return P(cmds), nil
}

// First returns the Chain's initial Cmd.
func (c P) First() *exec.Cmd {
	return c[0]
}

// Last returns the last Cmd in the Chain.
func (c P) Last() *exec.Cmd {
	return c[len(c)-1]
}

// Start starts all of the process in the chain. If any fails to start,
// all previous processes are killed.
func (c P) Start() error {
	var err error
	i := 0
	for ; i < len(c); i++ {
		err = c[i].Start()
		if err != nil {
			break
		}
	}

	if err == nil {
		return nil
	}

	for j := 0; j < i; j++ {
		_ = c[j].Process.Kill() //TODO: What can be done if this fails?
	}

	return err
}

// Wait waits for all of the process in the chain to finish, and
// returns any non-nil errors that they return.
func (c P) Wait() []error {
	errs := make([]error, 0)

	for _, c := range c {
		err := c.Wait()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// String returns a string of the form cmd0 | cmd1 | cmd2.
// It does not try to produce a shell-compatible string or print
// the commands' arguments.
func (c P) String() string {
	s := ""

	if len(c) > 0 {
		s = c[0].Path
	}

	for _, cmd := range c[1:] {
		s += " | " + cmd.Path
	}

	return s
}
