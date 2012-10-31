package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Command embeds a exec.Cmd but also includes buffers for stdin, stdout
// and stderr. These buffers are automatically attached when "New" is called.
type Command struct {
	*exec.Cmd
	BufStdin, BufStdout, BufStderr *bytes.Buffer
}

// New creates a new pointer to a Command. Byte buffers are created and
// attached to the command's Stdin, Stdout and Stderr.
func New(name string, arg ...string) *Command {
	var stdin, stdout, stderr *bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return &Command{
		Cmd:       cmd,
		BufStdin:  stdin,
		BufStdout: stdout,
		BufStderr: stderr,
	}
}

// Run calls (*exec.Cmd).Run on the embedded command. If (*exec.Cmd).Run returns
// an error, then Run will also return the error. But Run also checks the
// stderr buffer, and if it isn't empty, an error is returned with the contents
// of stderr.
func (cmd *Command) Run() error {
	fullCmd := strings.Join(cmd.Args, " ")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error starting '%s': %s.", fullCmd, err)
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// Wait calls (*exec.Cmd).Wait on the embedded command and handles errors
// as described in Run().
// Note that you may call (*Command).Start() since the Command type embeds a
// *exec.Cmd type.
func (cmd *Command) Wait() error {
	if err := cmd.Cmd.Wait(); err != nil {
		fullCmd := strings.Join(cmd.Args, " ")
		if cmd.BufStderr.Len() > 0 {
			return fmt.Errorf("Error running '%s': %s.\n\n%s",
				fullCmd, err, cmd.BufStderr.String())
		}
		return fmt.Errorf("Error running '%s': %s.", fullCmd, err)
	}
	return nil
}
