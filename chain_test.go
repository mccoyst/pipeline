// © 2012 Steve McCoy. Available under the MIT license.

package pipeline

import (
	"os"
	"os/exec"
	"testing"
)

func ExampleNew() {
	c := exec.Command
	cmds, err := New(
		c("echo", "hi"),
		c("sed", "s/h/H/"),
		c("sed", "s/i/ello/"))

	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		return
	}

	cmds.Last().Stdout = os.Stdout

	err = cmds.Start()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		return
	}

	errs := cmds.Wait()
	if len(errs) > 0 {
		for _, err := range errs {
			os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		}
		return
	}

	// Output: Hello
}

func TestString(t *testing.T) {
	c := exec.Command

	cmds, err := New(c("echo"))
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	s := cmds.String()
	exp := "/bin/echo"
	if s != exp {
		t.Errorf("Got %q, expected %q", s, exp)
	}

	cmds, err = New(
		c("echo", "hi"),
		c("sed", "s/h/H/"),
		c("sed", "s/i/ello/"))

	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	s = cmds.String()
	exp = "/bin/echo hi | /usr/bin/sed s/h/H/ | /usr/bin/sed s/i/ello/"
	if s != exp {
		t.Errorf("Got %q, expected %q", s, exp)
	}
}
