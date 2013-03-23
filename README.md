Package pipeline provides a datatype for simple chaining of the standard os/exec package's *Cmd.

For example:

		c := exec.Command
		cmds, err := New(
			c("echo", "hi"),
			c("sed", "s/h/H/"),
			c("sed", "s/i/ello/"))

		…

		_ = cmds.Start()
		_ = cmds.Wait()

		// Output: Hello
