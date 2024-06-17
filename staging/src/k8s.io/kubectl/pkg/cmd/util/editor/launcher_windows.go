/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"k8s.io/klog/v2"

	"k8s.io/kubectl/pkg/util/term"
)

// Enclose argument in double-quotes. Double each double-quote character as
// an escape sequence.
func cmdQuoteArg(arg string) string {
	var result strings.Builder
	result.WriteString(`"`)
	result.WriteString(strings.ReplaceAll(arg, `"`, `""`))
	result.WriteString(`"`)
	return result.String()
}

func (e Editor) args(path string) []string {
	args := make([]string, len(e.Args))
	copy(args, e.Args)
	if e.Shell {
		last := args[len(args)-1]
		if args[0] == windowsShell {
			// Use double-quotation around whole command line string
			// See https://stackoverflow.com/a/6378038
			args[len(args)-1] = fmt.Sprintf(`"%s %s"`, last, cmdQuoteArg(path))
		} else {
			args[len(args)-1] = fmt.Sprintf("%s %q", last, path)
		}
	} else {
		args = append(args, path)
	}
	return args
}

// Launch opens the described or returns an error. The TTY will be protected, and
// SIGQUIT, SIGTERM, and SIGINT will all be trapped.
func (e Editor) Launch(path string) error {
	if len(e.Args) == 0 {
		return fmt.Errorf("no editor defined, can't open %s", path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	args := e.args(abs)
	var cmd *exec.Cmd
	if e.Shell && args[0] == windowsShell {
		// Pass all arguments to cmd.exe as one string
		// See https://pkg.go.dev/os/exec#Command
		cmd = exec.Command(args[0])
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.SysProcAttr.CmdLine = strings.Join(args, " ")
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	klog.V(5).Infof("Opening file with editor %v", args)
	if err := (term.TTY{In: os.Stdin, TryDev: true}).Safe(cmd.Run); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("unable to launch the editor %q", strings.Join(e.Args, " "))
		}
		return fmt.Errorf("there was a problem with the editor %q", strings.Join(e.Args, " "))
	}
	return nil
}
