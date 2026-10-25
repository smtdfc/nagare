package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	pkgs_messages "github.com/smtdfc/nagare/pkgs/messages"
)

type fakeApprovalBindings struct {
	approve bool
	calls   []string
}

func (f *fakeApprovalBindings) RefreshTask(*core_context.ExecuteContext) {}

func (f *fakeApprovalBindings) CreateTask(*core_context.ExecuteContext, string, string, string, string, bool, string, string, string) (string, error) {
	return "", nil
}

func (f *fakeApprovalBindings) FindToolsByCategories(*core_context.ExecuteContext, []string) ([]tool.Metadata, error) {
	return nil, nil
}

func (f *fakeApprovalBindings) CallTool(*core_context.ExecuteContext, string, string) *tool.Result {
	return nil
}

func (f *fakeApprovalBindings) RequestToolApproval(_ *core_context.ExecuteContext, title string) bool {
	f.calls = append(f.calls, title)
	return f.approve
}

func approvalTestContext() *core_context.ExecuteContext {
	return &core_context.ExecuteContext{
		Context:        context.Background(),
		MessageChannel: make(chan pkgs_messages.Message, 1),
	}
}

func withBindings(t *testing.T, tl tool.Tool, bindings tool.Bindings) tool.Tool {
	t.Helper()

	original := tl.GetBindings()
	t.Cleanup(func() {
		tl.WithBindings(original)
	})
	return tl.WithBindings(bindings)
}
