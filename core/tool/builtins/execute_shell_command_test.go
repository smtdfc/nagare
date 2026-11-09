package declarations

import (
	"encoding/json"
	"testing"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/stretchr/testify/assert"
)

func TestExecuteShellCommandTool_Metadata(t *testing.T) {
	assert.Equal(t, "execute_shell_command_tool", ExecuteShellCommandTool.GetName())
	assert.Contains(t, ExecuteShellCommandTool.GetCategories(), tool.ProcessManagementCategory)
	assert.Contains(t, ExecuteShellCommandTool.GetCategories(), tool.SystemCategory)
}

func TestExecuteShellCommandTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, ExecuteShellCommandTool, approvals)

	_, err := ExecuteShellCommandTool.Execute(ctx, `{"command":""}`)
	assert.Error(t, err)
	assert.Len(t, approvals.calls, 0)
}

func TestExecuteShellCommandTool_DeniedByUser(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: false}
	withBindings(t, ExecuteShellCommandTool, approvals)

	_, err := ExecuteShellCommandTool.Execute(ctx, `{"command":"echo hello"}`)
	assert.Error(t, err)
	assert.Len(t, approvals.calls, 1)
}

func TestExecuteShellCommandTool_Success(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, ExecuteShellCommandTool, approvals)

	rawOutput, err := ExecuteShellCommandTool.Execute(ctx, `{"command":"echo hello"}`)
	assert.NoError(t, err)
	assert.Len(t, approvals.calls, 1)

	var output ExecuteShellCommandOutput
	err = json.Unmarshal([]byte(rawOutput), &output)
	assert.NoError(t, err)
	assert.Equal(t, 0, output.ExitCode)
	assert.Contains(t, output.Stdout, "hello")
}
