package contracts

const (
	InjectSystemPromptEvent        Event = "plugin:inject_system_prompt"
	InjectSystemPromptSuccessEvent Event = "plugin:inject_system_prompt:success"
	InjectSystemPromptFailedEvent  Event = "plugin:inject_system_prompt:failed"

	ClearSystemPromptEvent        Event = "plugin:clear_system_prompt"
	ClearSystemPromptSuccessEvent Event = "plugin:clear_system_prompt:success"
	ClearSystemPromptFailedEvent  Event = "plugin:clear_system_prompt:failed"

	// MaxSystemPromptBytes limits the size of a plugin-provided system prompt
	// fragment so a misbehaving plugin cannot exhaust gateway memory or build
	// unbounded prompts.
	MaxSystemPromptBytes = 64 * 1024
)

// InjectSystemPromptEventPayload sets the plugin's system prompt fragment.
// The fragment is appended to every chat invoke regardless of session or
// channel. An empty Content clears the fragment.
type InjectSystemPromptEventPayload struct {
	Content string `json:"content"`
}

type InjectSystemPromptSuccessEventPayload struct{}

type InjectSystemPromptFailedEventPayload struct {
	Cause string `json:"cause"`
}

type ClearSystemPromptEventPayload struct{}

type ClearSystemPromptSuccessEventPayload struct{}

type ClearSystemPromptFailedEventPayload struct {
	Cause string `json:"cause"`
}
