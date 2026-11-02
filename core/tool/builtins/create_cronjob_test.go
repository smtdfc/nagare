package declarations

import (
	"testing"
)

func TestCreateCronJobToolMetadata(t *testing.T) {
	tool := CreateCronJobTool
	if tool.GetName() != "create_cronjob_tool" {
		t.Errorf("expected name 'create_cronjob_tool', got %q", tool.GetName())
	}
	categories := tool.GetCategories()
	if len(categories) != 1 || categories[0] != "cronjob_management" {
		t.Errorf("unexpected categories: %v", categories)
	}
}
