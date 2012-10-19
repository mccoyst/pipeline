// © 2012 Steve McCoy. Available under the MIT license.

package pipeline

import (
	"os"
	"os/exec"
)

func ExampleNew() {
	c := exec.Command
	cmds, err := New(
		c("echo", "hi"),
		c("sed", "s/h/H/"),
		c("sed", "s/i/ello/"))

	if err != nil {
		os.Stderr.WriteString("Oops: "+err.Error()+"\n")
		return
	}

	cmds.Last().Stdout = os.Stdout

	err = cmds.Start()
	if err != nil {
		os.Stderr.WriteString("Oops: "+err.Error()+"\n")
		return
	}

	errs := cmds.Wait()
	if len(errs) > 0 {
		for _, err := range errs {
			os.Stderr.WriteString("Oops: "+err.Error()+"\n")
		}
		return
	}

	// Output: Hello
}
