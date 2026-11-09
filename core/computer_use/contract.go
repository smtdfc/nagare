package computer_use

const (
	GET_APP_INSTANCE  = "get_app_instance"
	LIST_RUNNING_APPS = "list_running_apps"
	OPEN_APP          = "open_app"
	GET_UI_TREE       = "get_ui_tree"
	CLICK_ELEMENT     = "click_element"
	INPUT_TEXT        = "input_text"
	MOUSE_ACTION      = "mouse_action"
	KEYBOARD_SHORTCUT = "keyboard_shortcut"
	WINDOW_CONTROL    = "window_control"
	QUERY_ELEMENT_BY_TEXT = "query_element_by_text"
)

type GetAppInstanceRequest struct {
	Name string `json:"name"`
}

type GetAppInstanceResponse struct {
	IsError    bool   `json:"is_error"`
	InstanceID string `json:"instance_id"`
	Error      string `json:"error,omitempty"`
}

type AppInfo struct {
	Name string  `json:"name"`
	PID  *uint32 `json:"pid,omitempty"`
}

type ListRunningAppsRequest struct{}

type ListRunningAppsResponse struct {
	IsError bool      `json:"is_error"`
	Apps    []AppInfo `json:"apps"`
	Error   string    `json:"error,omitempty"`
}

type OpenAppRequest struct {
	Name string   `json:"name"`
	Args []string `json:"args,omitempty"`
}

type OpenAppResponse struct {
	IsError bool    `json:"is_error"`
	PID     *uint32 `json:"pid,omitempty"`
	Error   string  `json:"error,omitempty"`
}

type GetUITreeRequest struct {
	InstanceID string `json:"instance_id"`
	MaxDepth   *int   `json:"max_depth,omitempty"`
}

type UINode struct {
	RefID    string   `json:"ref_id"`
	Role     string   `json:"role"`
	Name     string   `json:"name"`
	Value    *string  `json:"value,omitempty"`
	X        int32    `json:"x"`
	Y        int32    `json:"y"`
	Width    uint32   `json:"width"`
	Height   uint32   `json:"height"`
	Children []UINode `json:"children"`
}

type GetUITreeResponse struct {
	IsError bool    `json:"is_error"`
	Tree    *UINode `json:"tree,omitempty"`
	Summary string  `json:"summary"`
	Error   string  `json:"error,omitempty"`
}

type ClickElementRequest struct {
	InstanceID string `json:"instance_id"`
	RefID      string `json:"ref_id"`
}

type ClickElementResponse struct {
	IsError bool   `json:"is_error"`
	Error   string `json:"error,omitempty"`
}

type InputTextRequest struct {
	InstanceID string `json:"instance_id,omitempty"`
	RefID      string `json:"ref_id,omitempty"`
	Text       string `json:"text"`
}

type InputTextResponse struct {
	IsError bool   `json:"is_error"`
	Error   string `json:"error,omitempty"`
}

type MouseActionRequest struct {
	Action       string   `json:"action"` // "move", "click", "right_click", "double_click", "scroll"
	X            *float64 `json:"x,omitempty"`
	Y            *float64 `json:"y,omitempty"`
	ScrollDeltaX *float64 `json:"scroll_delta_x,omitempty"`
	ScrollDeltaY *float64 `json:"scroll_delta_y,omitempty"`
}

type MouseActionResponse struct {
	IsError bool   `json:"is_error"`
	Error   string `json:"error,omitempty"`
}

type KeyboardShortcutRequest struct {
	Key       string   `json:"key"`
	Modifiers []string `json:"modifiers,omitempty"`
}

type KeyboardShortcutResponse struct {
	IsError bool   `json:"is_error"`
	Error   string `json:"error,omitempty"`
}

type WindowControlRequest struct {
	InstanceID string `json:"instance_id"`
	Action     string `json:"action"` // "activate", "minimize", "maximize", "restore", "close"
}

type WindowControlResponse struct {
	IsError bool   `json:"is_error"`
	Error   string `json:"error,omitempty"`
}

type QueryElementByTextRequest struct {
	InstanceID string `json:"instance_id"`
	Text       string `json:"text"`
	Exact      *bool  `json:"exact,omitempty"`
}

type ElementMatch struct {
	RefID  string  `json:"ref_id"`
	Role   string  `json:"role"`
	Name   string  `json:"name"`
	Value  *string `json:"value,omitempty"`
	X      int32   `json:"x"`
	Y      int32   `json:"y"`
	Width  uint32  `json:"width"`
	Height uint32  `json:"height"`
}

type QueryElementByTextResponse struct {
	IsError bool           `json:"is_error"`
	Matches []ElementMatch `json:"matches"`
	Error   string         `json:"error,omitempty"`
}
