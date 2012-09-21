// © 2012 Steve McCoy. Available under the MIT license.

/*
Package cmdchain provides a datatype Chain for simple chaining of
the standard os/exec package's *Cmd.
*/
package cmdchain

import (
	"errors"
	"os/exec"
)

type Chain []*exec.Cmd

// New returns a new Chain. If the Cmds' Stdins and Stdouts cannot
// be piped together for any reason, nil and the error are returns.
//
// This function also returns an error if len(cmds) == 0.
func New(cmds ...*exec.Cmd) (Chain, error) {
	if len(cmds) == 0 {
		return nil, errors.New("cmdchain.New() requires at least one command")
	}

	for i := 1; i < len(cmds); i++ {
		out, err := cmds[i-1].StdoutPipe()
		if err != nil {
			return nil, err
		}
		cmds[i].Stdin = out
	}

	return Chain(cmds), nil
}

// First returns the Chain's initial Cmd.
func (c Chain) First() *exec.Cmd {
	return c[0]
}

// Last returns the last Cmd in the Chain.
func (c Chain) Last() *exec.Cmd {
	return c[len(c)-1]
}

// Start starts all of the process in the chain. If any fails to start,
// all previous processes are killed.
func (c Chain) Start() error {
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
func (c Chain) Wait() []error {
	errs := make([]error, 0)

	for _, c := range c {
		err := c.Wait()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
