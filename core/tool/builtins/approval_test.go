package declarations

import (
	"context"
	"testing"

	core_context "github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/discovery/manager"
	"github.com/smtdfc/nagare/core/tool"
	pkgs_messages "github.com/smtdfc/nagare/pkgs/messages"
)

type mockCronJobBindings struct{}

func (m *mockCronJobBindings) RefreshCronJob(*core_context.ExecuteContext) {}

func (m *mockCronJobBindings) CreateCronJob(*core_context.ExecuteContext, string, string, string, string) (string, error) {
	return "", nil
}

func (m *mockCronJobBindings) FindToolsByCategories(*core_context.ExecuteContext, []string) ([]tool.Metadata, error) {
	return nil, nil
}

func (m *mockCronJobBindings) CallTool(*core_context.ExecuteContext, string, string) *tool.Result {
	return nil
}

func (m *mockCronJobBindings) RequestToolApproval(*core_context.ExecuteContext, string) bool {
	return true
}

type fakeApprovalBindings struct {
	mockCronJobBindings
	approve bool
	calls   []string
}

// GetDiscoveryMgr implements [tool.Bindings].
func (f *fakeApprovalBindings) GetDiscoveryMgr() *manager.DiscoveryManager {
	panic("unimplemented")
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
