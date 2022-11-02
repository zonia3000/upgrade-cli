package spawn

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// Spawns a process and optionally captures standard output and standard error
//   - gocmd             receiver of the "&exec.Cmd" created and used for the execution or nil if your don't care
//   - baseCmd           the program to execute
//   - args              the arguments in the spawn package special format [*1]
//   - env               a string->string map that specifies the spawn execution environment
//   - opts              activation of optional functions (see Options)
//
// Notes:
//   - [*1] Allows:
//   - simple strings
//   - nesting of string arrays to any level, that are simply flatten
//   - known objects like the SubOptionArg (check the spawn.SOA function).
//
// noinspection GoNameStartsWithPackageName
func Spawn(gocmd **exec.Cmd, baseCmd string, args []interface{}, env Environ, opts Options) (Res, error) {
	if gocmd == nil {
		tmpCmd := &exec.Cmd{}
		gocmd = &tmpCmd
	}

	// PREPARE THE COMMAND AND ITS ARGS
	if opts.WithSudo {
		args = PrependToArgs(args, baseCmd)
		*gocmd = exec.Command("sudo", MkRawSpawnArgs(args)...)
	} else {
		args := MkRawSpawnArgs(args)
		*gocmd = exec.Command(baseCmd, args...)
	}

	// SET THE SPAWN ENVIRONMENT
	if len(env) > 0 {
		(*gocmd).Env = os.Environ()
		(*gocmd).Env = addEnv((*gocmd).Env, env)
	}

	if opts.Interactive {
		(*gocmd).Stdin = os.Stdin
	}

	var stdoutPipe, stderrPipe io.ReadCloser
	var stdoutErr, stderrErr error
	if opts.CaptureStdout {
		stdoutPipe, stdoutErr = (*gocmd).StdoutPipe()
	} else {
		(*gocmd).Stdout = os.Stdout
		stdoutPipe = nil
	}
	if opts.CaptureStderr {
		stderrPipe, stderrErr = (*gocmd).StderrPipe()
	} else {
		(*gocmd).Stderr = os.Stdout
		stderrPipe = nil
	}

	// START THE PROCESS
	err := (*gocmd).Start()
	if err != nil {
		return Res{}, err
	}

	// COLLECT THE OUTPUT
	var capturedStdout []byte
	var capturedStderr []byte
	if stdoutErr == nil && opts.CaptureStdout {
		capturedStdout, _ = ioutil.ReadAll(stdoutPipe)
	}
	if stderrErr == nil && opts.CaptureStderr {
		capturedStderr, _ = ioutil.ReadAll(stderrPipe)
	}

	// WAIT FOR PROCESS TO COMPLETE
	_ = (*gocmd).Wait()

	return Res{string(capturedStdout), string(capturedStderr)}, err
}

// Optional functionalities activation flags
type Options struct {
	WithSudo      bool // the command is run using sudo (in the platforms that supports it)
	Interactive   bool // the command stdin is attached to the tty
	CaptureStdout bool // the command standard output is intercepted (see Res)
	CaptureStderr bool // the command standard error is intercepted (see Res)
}

// Composes a simple sub-option assignment argument:
//
// Example:
//   - CODE:   SOA("o", "key", "value")
//   - RESULT: -o key=value
func SOA(prefix, key, value string) SubOptionArg {
	return SubOptionArg{prefix, key, value}
}
