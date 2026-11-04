package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/message"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/helpers"
	"github.com/smtdfc/nagare/pkgs/messages"
	"google.golang.org/genai"
	"google.golang.org/genai/interactions/models/apierrors"
	"google.golang.org/genai/interactions/models/components"
	"google.golang.org/genai/interactions/models/interactions"
	"google.golang.org/genai/interactions/models/operations"
)

type GoogleAdapter struct {
	client *genai.Client
	Models []string
	logger *logger.BaseLogger
}

// Embed implements [provider.ModelProviderAdapter].
func (g *GoogleAdapter) Embed(ctx context.Context, model string, inputs []string) ([][]float32, error) {
	contents := make([]*genai.Content, 0, len(inputs))
	for _, input := range inputs {
		contents = append(contents, genai.NewContentFromText(input, genai.RoleUser))
	}

	resp, err := g.client.Models.EmbedContent(ctx, model, contents, nil)
	if err != nil {
		g.logger.Error("embed error", "error", err)
		return nil, custom_errors.ErrModelProviderAdapter
	}

	result := make([][]float32, 0, len(resp.Embeddings))
	for _, embedding := range resp.Embeddings {
		result = append(result, embedding.Values)
	}

	return result, nil
}

// GetModels implements [provider.ModelProviderAdapter].
func (g *GoogleAdapter) GetModels(ctx context.Context) ([]string, error) {
	result := make([]string, 0)

	for model, err := range g.client.Models.All(ctx) {
		if err != nil {
			g.logger.Error("list models error", "error", err)
			return nil, custom_errors.ErrGetListModelFailed
		}

		result = append(result, strings.TrimPrefix(model.Name, "models/"))
	}

	return result, nil
}

// TransformToSystemInstruction joins system/developer messages into the single
// system instruction string expected by the Interactions API.
func (g *GoogleAdapter) TransformToSystemInstruction(inputs messages.ListMessage) *string {
	parts := make([]string, 0, len(inputs))
	for _, msg := range inputs {
		text, ok := msg.(*messages.TextMessage)
		if !ok || (text.Role != messages.SYSTEM && text.Role != messages.DEVELOPER) {
			continue
		}

		if content := strings.TrimSpace(text.Content); content != "" {
			parts = append(parts, content)
		}
	}

	if len(parts) == 0 {
		return nil
	}

	instruction := strings.Join(parts, "\n\n")
	return &instruction
}

func (g *GoogleAdapter) TransformToProviderMessage(msg messages.Message) (interactions.Step, error) {
	switch t := msg.(type) {
	case *messages.TextMessage:
		switch t.Role {
		case messages.USER:
			return interactions.NewStep(interactions.UserInputStep{
				Content: []interactions.Content{
					interactions.NewContent(interactions.TextContent{
						Text: t.Content,
					}),
				},
			}), nil

		case messages.AGENT:
			return interactions.NewStep(interactions.ModelOutputStep{
				Content: []interactions.Content{
					interactions.NewContent(interactions.TextContent{
						Text: t.Content,
					}),
				},
			}), nil
		}

	case *messages.ToolCallMessage:
		arguments := map[string]any{}
		if strings.TrimSpace(t.Args) != "" {
			if err := json.Unmarshal([]byte(t.Args), &arguments); err != nil {
				return interactions.Step{}, fmt.Errorf("invalid tool call arguments: %w", err)
			}
		}

		return interactions.NewStep(interactions.FunctionCallStep{
			ID:        t.CallID,
			Name:      t.Name,
			Arguments: arguments,
		}), nil

	case *messages.ToolResultMessage:
		step := interactions.FunctionResultStep{
			CallID: t.CallID,
			Result: interactions.NewFunctionResultStepResultUnion(t.Result),
		}
		if t.Name != "" {
			step.Name = genai.Ptr(t.Name)
		}

		return interactions.NewStep(step), nil
	}

	return interactions.Step{}, custom_errors.ErrMessageNotSupport
}

func (g *GoogleAdapter) TransformToProviderInput(inputs messages.ListMessage) ([]interactions.Step, error) {
	steps := make([]interactions.Step, 0, len(inputs))
	for _, msg := range inputs {
		if text, ok := msg.(*messages.TextMessage); ok &&
			(text.Role == messages.SYSTEM || text.Role == messages.DEVELOPER) {
			continue
		}

		step, err := g.TransformToProviderMessage(msg)
		if err != nil {
			return nil, err
		}

		steps = append(steps, step)
	}

	return steps, nil
}

func (g *GoogleAdapter) TransformToolDeclarations(tools tool.ListTool) ([]interactions.Tool, error) {
	declarations := make([]interactions.Tool, 0, len(tools))

	for _, t := range tools {
		parameters, err := helpers.UnmarshalJson[map[string]any](t.GetArgsSchema())
		if err != nil {
			g.logger.Error("TransformToolDeclarations error", "tool", t.GetName(), "error", err)
			return nil, custom_errors.ErrInvalidToolSchema
		}

		declarations = append(declarations, interactions.NewTool(interactions.Function{
			Name:        genai.Ptr(t.GetName()),
			Description: genai.Ptr(t.GetDescription()),
			Parameters:  *parameters,
		}))
	}

	return declarations, nil
}

func (g *GoogleAdapter) buildCreateRequest(model string, inputs messages.ListMessage, tools tool.ListTool) (operations.CreateInteractionRequestBody, error) {
	steps, err := g.TransformToProviderInput(inputs)
	if err != nil {
		g.logger.Error("TransformToProviderInput error", "error", err)
		return operations.CreateInteractionRequestBody{}, err
	}

	declarations, err := g.TransformToolDeclarations(tools)
	if err != nil {
		return operations.CreateInteractionRequestBody{}, err
	}

	return operations.NewCreateInteractionRequestBody(interactions.CreateModelInteraction{
		Model:             interactions.Model(model),
		Input:             genai.Ptr(interactions.NewInteractionsInput(steps)),
		Stream:            genai.Ptr(true),
		Store:             genai.Ptr(false),
		SystemInstruction: g.TransformToSystemInstruction(inputs),
		Tools:             declarations,
	}), nil
}

// Send implements [provider.ModelProviderAdapter].
func (g *GoogleAdapter) Send(ctx context.Context, model string, inputs messages.ListMessage, tools tool.ListTool) (message.ReadOnlyChannel, error) {
	body, err := g.buildCreateRequest(model, inputs, tools)
	if err != nil {
		return nil, err
	}

	outputChannel := make(chan messages.Message)

	go func() {
		defer close(outputChannel)

		emit := func(msg messages.Message) bool {
			select {
			case outputChannel <- msg:
				return true
			case <-ctx.Done():
				return false
			}
		}

		res, err := g.client.Interactions.Create(ctx, operations.CreateInteractionRequest{
			Body: body,
		})
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			g.logger.Error("create interaction error", "error", err)
			code, cause := g.mapCreateError(err, model)
			emit(messages.NewResponseFailedMessage(code, cause))
			return
		}

		if res == nil || res.InteractionSSEStreamEvent == nil {
			g.logger.Error("create interaction returned no stream")
			emit(messages.NewResponseFailedMessage("500", custom_errors.ErrMessageStream.Error()))
			return
		}

		stream := res.InteractionSSEStreamEvent
		defer stream.Close()

		toolCalls := map[int]*toolCallBuffer{}
		emitted := false

		flushToolCall := func(index int) bool {
			buffer, ok := toolCalls[index]
			if !ok {
				return true
			}

			delete(toolCalls, index)

			if !emit(messages.NewToolCallMessage(buffer.id, buffer.name, buffer.args())) {
				return false
			}

			emitted = true
			return true
		}

		for stream.Next() {
			event := stream.Value()
			if event == nil {
				continue
			}

			if eventError := event.GetDataError(); eventError != nil {
				code, cause := "500", "interaction stream error"
				if detail := eventError.GetError(); detail != nil {
					if detail.GetCode() != nil {
						code = *detail.GetCode()
					}
					if detail.GetMessage() != nil {
						cause = *detail.GetMessage()
					}
				}

				g.logger.Error("interaction stream error event", "code", code, "cause", cause)
				emit(messages.NewResponseFailedMessage(code, cause))
				return
			}

			if start := event.GetDataStepStart(); start != nil {
				if call := start.GetStepFunctionCall(); call != nil {
					toolCalls[start.GetIndex()] = newToolCallBuffer(
						call.GetID(),
						call.GetName(),
						call.GetArguments(),
					)
				}
				continue
			}

			if delta := event.GetDataStepDelta(); delta != nil {
				index := delta.GetIndex()

				if text := delta.GetDeltaText(); text != nil && text.GetText() != "" {
					if !emit(messages.NewTextMessage(messages.AGENT, text.GetText())) {
						return
					}
					emitted = true
				}

				if thought := delta.GetDeltaThoughtSummary(); thought != nil {

					if content := thought.GetContentText(); content != nil && content.GetText() != "" {
						fmt.Println(thought)
						if !emit(messages.NewReasoningMessage(content.GetText())) {
							return
						}
						emitted = true
					}
				}

				if arguments := delta.GetDeltaArgumentsDelta(); arguments != nil {
					buffer, ok := toolCalls[index]
					if !ok {
						buffer = newToolCallBuffer("", "", nil)
						toolCalls[index] = buffer
					}
					if fragment := arguments.GetArguments(); fragment != nil {
						buffer.deltas.WriteString(*fragment)
					}
				}

				if signature := delta.GetDeltaThoughtSignature(); signature != nil {
					g.logger.Debug("thought signature received", "index", index)
				}

				continue
			}

			if stop := event.GetDataStepStop(); stop != nil {
				if !flushToolCall(stop.GetIndex()) {
					return
				}
				continue
			}

			if completed := event.GetDataInteractionCompleted(); completed != nil && !emitted {
				completedSteps := completed.GetInteraction().Steps
				if !g.emitCompletedSteps(completedSteps, emit) {
					return
				}
				emitted = len(completedSteps) > 0
			}
		}

		indexes := make([]int, 0, len(toolCalls))
		for index := range toolCalls {
			indexes = append(indexes, index)
		}
		sort.Ints(indexes)
		for _, index := range indexes {
			if !flushToolCall(index) {
				return
			}
		}

		if err := stream.Err(); err != nil {
			if ctx.Err() != nil || errors.Is(err, context.Canceled) {
				return
			}

			g.logger.Error("stream error", "error", err)
			code, cause := g.mapStreamError(err, model)
			emit(messages.NewResponseFailedMessage(code, cause))
		}
	}()

	return outputChannel, nil
}

func (g *GoogleAdapter) emitCompletedSteps(steps []interactions.Step, emit func(messages.Message) bool) bool {
	for _, step := range steps {
		if output := step.ModelOutputStep; output != nil {
			for _, content := range output.Content {
				if content.TextContent == nil || content.TextContent.Text == "" {
					continue
				}

				if !emit(messages.NewTextMessage(messages.AGENT, content.TextContent.Text)) {
					return false
				}
			}
		}

		if call := step.FunctionCallStep; call != nil {
			arguments, err := json.Marshal(call.GetArguments())
			if err != nil {
				g.logger.Error("marshal tool call arguments error", "tool", call.GetName(), "error", err)
				continue
			}

			if !emit(messages.NewToolCallMessage(call.GetID(), call.GetName(), string(arguments))) {
				return false
			}
		}
	}

	return true
}

func (g *GoogleAdapter) mapCreateError(err error, model string) (string, string) {
	var apiErr *apierrors.APIError
	if errors.As(err, &apiErr) {
		return mapProviderStatus(apiErr.StatusCode, apiErr.Error(), model)
	}

	var validationErr *apierrors.ResponseValidationError
	if errors.As(err, &validationErr) {
		return "500", validationErr.Message
	}

	var interruptedErr *apierrors.ContinuationInterruptedError
	if errors.As(err, &interruptedErr) {
		return "500", interruptedErr.Message
	}

	var serverErr *apierrors.CreateInteractionServerError
	if errors.As(err, &serverErr) {
		status := httpStatus(serverErr.HTTPMeta)
		if status < 500 {
			status = 500
		}
		return mapProviderStatus(
			status,
			providerErrorText(serverErr.Error_.GetCode(), serverErr.Error_.GetMessage(), err),
			model,
		)
	}

	var clientErr *apierrors.CreateInteractionClientError
	if errors.As(err, &clientErr) {
		return mapProviderStatus(
			httpStatus(clientErr.HTTPMeta),
			providerErrorText(clientErr.Error_.GetCode(), clientErr.Error_.GetMessage(), err),
			model,
		)
	}

	return "400", err.Error()
}

func mapProviderStatus(status int, cause string, model string) (string, string) {
	switch {
	case status == 404:
		return "404", fmt.Sprintf("Model %s not found", model)
	case status == 402:
		return "402", cause
	case status == 429:
		return "429", cause
	case status >= 500:
		return "500", cause
	default:
		return "400", cause
	}
}

func httpStatus(meta components.HTTPMetadata) int {
	if meta.Response == nil {
		return 0
	}

	return meta.Response.StatusCode
}

// providerErrorText prefers the API's human readable message over the raw JSON
// body the SDK error stringifies to.
func providerErrorText(code, message *string, fallback error) string {
	if message != nil && strings.TrimSpace(*message) != "" {
		return *message
	}

	if code != nil && strings.TrimSpace(*code) != "" {
		return *code
	}

	return fallback.Error()
}

func (g *GoogleAdapter) mapStreamError(err error, model string) (string, string) {
	return g.mapCreateError(err, model)
}

// toolCallBuffer accumulates a streamed function call. Arguments can arrive
// either complete in `step.start` or as fragments in `step.delta`; deltas win
// when present.
type toolCallBuffer struct {
	id     string
	name   string
	seeded string
	deltas strings.Builder
}

func newToolCallBuffer(id, name string, arguments map[string]any) *toolCallBuffer {
	buffer := &toolCallBuffer{id: id, name: name}

	if len(arguments) > 0 {
		if encoded, err := json.Marshal(arguments); err == nil {
			buffer.seeded = string(encoded)
		}
	}

	return buffer
}

func (b *toolCallBuffer) args() string {
	if deltas := strings.TrimSpace(b.deltas.String()); deltas != "" {
		return deltas
	}

	if seeded := strings.TrimSpace(b.seeded); seeded != "" {
		return seeded
	}

	return "{}"
}

func NewGoogleAdapter(baseURL, apiKey string, models []string, logger *logger.BaseLogger) (*GoogleAdapter, error) {
	config := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}

	if baseURL != "" {
		config.HTTPOptions = genai.HTTPOptions{BaseURL: baseURL}
	}

	client, err := genai.NewClient(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &GoogleAdapter{
		client: client,
		Models: models,
		logger: logger.With("module", "google-adapter"),
	}, nil
}
