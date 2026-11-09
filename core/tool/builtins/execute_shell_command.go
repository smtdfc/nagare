package declarations

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	coreContext "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

const defaultShellTimeoutSeconds = 60
const maxOutputBytes = 256 * 1024 // 256 KB max output to avoid flooding context

type ExecuteShellCommandInput struct {
	Command        string `json:"command" jsonschema_description:"The full shell command string to execute (e.g. 'ls -la', 'git status', 'cargo build'). Supports shell syntax like pipes, redirects, and environment variables."`
	WorkingDir     string `json:"working_dir,omitempty" jsonschema_description:"Optional working directory where the command should be run. Defaults to the current working directory."`
	TimeoutSeconds *int   `json:"timeout_seconds,omitempty" jsonschema_description:"Execution timeout in seconds. Default is 60 seconds."`
}

type ExecuteShellCommandOutput struct {
	ExitCode int    `json:"exit_code" jsonschema_description:"Process exit code (0 usually indicates success)."`
	Stdout   string `json:"stdout" jsonschema_description:"Standard output produced by the command."`
	Stderr   string `json:"stderr" jsonschema_description:"Standard error produced by the command."`
	Duration string `json:"duration" jsonschema_description:"Command execution duration."`
}

var ExecuteShellCommandTool = tool.DefineTool(
	"execute_shell_command_tool",
	"Execute a command in the system shell (bash/sh on Unix, cmd.exe on Windows) synchronously and return its stdout, stderr, and exit code. Requires user approval.",
	func(ctx *coreContext.ExecuteContext, args *ExecuteShellCommandInput, bindings tool.Bindings) (*ExecuteShellCommandOutput, error) {
		if args == nil || args.Command == "" {
			return nil, fmt.Errorf("command cannot be empty")
		}

		workingDir := args.WorkingDir
		if workingDir != "" {
			resolvedDir, err := system.ResolveFilePath(workingDir)
			if err != nil {
				return nil, fmt.Errorf("invalid working_dir: %w", err)
			}
			workingDir = resolvedDir
		}

		timeout := defaultShellTimeoutSeconds
		if args.TimeoutSeconds != nil && *args.TimeoutSeconds > 0 {
			timeout = *args.TimeoutSeconds
		}

		approvalTitle := fmt.Sprintf("Allow the agent to execute shell command: %s", args.Command)
		if workingDir != "" {
			approvalTitle = fmt.Sprintf("Allow the agent to execute shell command in %s: %s", workingDir, args.Command)
		}

		if err := requireApproval(ctx, bindings, approvalTitle); err != nil {
			return nil, err
		}

		execCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.CommandContext(execCtx, "cmd.exe", "/c", args.Command)
		} else {
			cmd = exec.CommandContext(execCtx, "sh", "-c", args.Command)
		}

		if workingDir != "" {
			cmd.Dir = workingDir
		}

		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		startTime := time.Now()
		runErr := cmd.Run()
		elapsed := time.Since(startTime)

		exitCode := 0
		if runErr != nil {
			if exitErr, ok := runErr.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else if execCtx.Err() == context.DeadlineExceeded {
				return nil, fmt.Errorf("command timed out after %d seconds", timeout)
			} else {
				return nil, fmt.Errorf("failed to execute command: %w", runErr)
			}
		}

		stdoutStr := stdoutBuf.String()
		stderrStr := stderrBuf.String()

		if len(stdoutStr) > maxOutputBytes {
			stdoutStr = stdoutStr[:maxOutputBytes] + "\n...[stdout truncated: exceeded 256KB limit]"
		}
		if len(stderrStr) > maxOutputBytes {
			stderrStr = stderrStr[:maxOutputBytes] + "\n...[stderr truncated: exceeded 256KB limit]"
		}

		return &ExecuteShellCommandOutput{
			ExitCode: exitCode,
			Stdout:   stdoutStr,
			Stderr:   stderrStr,
			Duration: elapsed.Round(time.Millisecond).String(),
		}, nil
	},
	[]string{tool.ProcessManagementCategory, tool.SystemCategory},
)
