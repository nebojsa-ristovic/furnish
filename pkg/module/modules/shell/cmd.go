package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

var shell = "/bin/bash"

func init() {
	if BinaryExists("zsh") {
		shell = "/bin/zsh"
		color.HiBlue("[init] zsh found: shell changed to zsh")
	}
}

func BinaryExists(pkg string) bool {
	if err := exec.Command(shell, "-c", fmt.Sprintf("which -s %s", pkg)).Run(); err != nil {
		return false
	}
	return true
}

func Exec(args ...string) error {
	allArgs := append(append(make([]string, 0, len(args)+1), "-c"), args...)
	cmd := exec.Command(shell, allArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func ExecSilent(args ...string) error {
	allArgs := append(append(make([]string, 0, len(args)+1), "-c"), args...)
	if err := exec.Command(shell, allArgs...).Run(); err != nil {
		return err
	}
	return nil
}

func ExecOutput(args ...string) (string, error) {
	allArgs := append(append(make([]string, 0, len(args)+1), "-c"), args...)
	out, err := exec.Command(shell, allArgs...).Output()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", out), nil
}

func Script(path string) error {
	cmd := exec.Command(shell, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func ScriptSilent(path string) error {
	if err := exec.Command(shell, path).Run(); err != nil {
		return err
	}
	return nil
}

func ScriptOutput(path string) (string, error) {
	out, err := exec.Command(shell, path).Output()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", out), nil
}
