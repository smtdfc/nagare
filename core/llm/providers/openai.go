package providers

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/shared/messages"
)

type OpenAICompatibleProviderAdapter struct {
	Client *openai.Client
	Models []string
}

func (o *OpenAICompatibleProviderAdapter) TransformToProviderMessage(msg messages.Message) (responses.ResponseInputItemUnionParam, error) {
	switch t := msg.(type) {
	case *messages.Text:
		message := &responses.ResponseInputItemMessageParam{
			Type: "message",
			Content: responses.ResponseInputMessageContentListParam{
				responses.ResponseInputContentUnionParam{
					OfInputText: &responses.ResponseInputTextParam{
						Text: t.Content,
					},
				},
			},
		}

		switch t.Role {
		case messages.USER:
			message.Role = "user"
			return responses.ResponseInputItemUnionParam{
				OfInputMessage: message,
			}, nil

		case messages.SYSTEM:
			message.Role = "system"
			return responses.ResponseInputItemUnionParam{
				OfInputMessage: message,
			}, nil
		case messages.DEVELOPER:
			message.Role = "developer"
			return responses.ResponseInputItemUnionParam{
				OfInputMessage: message,
			}, nil
		case messages.AGENT:
			return responses.ResponseInputItemUnionParam{
				OfOutputMessage: &responses.ResponseOutputMessageParam{
					Content: []responses.ResponseOutputMessageContentUnionParam{
						{
							OfOutputText: &responses.ResponseOutputTextParam{
								Text: t.Content,
							},
						},
					},
				},
			}, nil
		}

	case *messages.ToolCallResult:
		return responses.ResponseInputItemUnionParam{
			OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
				CallID: t.CallID,
				Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
					OfString: openai.String(t.Result),
				},
			},
		}, nil

	case *messages.ToolCall:
		return responses.ResponseInputItemUnionParam{
			OfFunctionCall: &responses.ResponseFunctionToolCallParam{
				CallID:    t.CallID,
				Name:      t.Name,
				Arguments: t.Args,
			},
		}, nil
	}

	return responses.ResponseInputItemUnionParam{}, nil
}

func (o *OpenAICompatibleProviderAdapter) TransformToolDeclarations(tools domains.ListTool) ([]responses.ToolUnionParam, error) {
	toolParams := make([]responses.ToolUnionParam, len(tools))

	for i, tool := range tools {
		var params map[string]any

		// println(tool.GetArgumentsSchema())
		err := json.Unmarshal([]byte(tool.GetArgs()), &params)
		if err != nil {
			return nil, custom_errors.NewLLMProviderError(fmt.Sprintf("JSON parse error %s: %v\n", tool.GetName(), err))
		}

		toolParams[i] = responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        tool.GetName(),
				Description: openai.String(tool.GetDescription()),
				Parameters:  params,
			},
		}
	}
	return toolParams, nil
}

func (o *OpenAICompatibleProviderAdapter) Chat(model string, ctx domains.Context, listMessage messages.ListMessage, tools domains.ListTool) (domains.MessageChannel, error) {
	if !slices.Contains(o.Models, model) {
		return nil, custom_errors.NewLLMProviderError(fmt.Sprintf("Model compatibility error: The current provider does not support the requested model '%s'.", model))
	}

	inputs := responses.ResponseInputParam{}
	listTool, err := o.TransformToolDeclarations(tools)
	if err != nil {
		return nil, err
	}

	for _, msg := range listMessage {
		input, err := o.TransformToProviderMessage(msg)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}

	outputChannel := make(chan messages.Message)

	go (func() {
		defer close(outputChannel)
		stream := o.Client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
			Model: model,
			Input: responses.ResponseNewParamsInputUnion{
				OfInputItemList: inputs,
			},
			Tools:       listTool,
			Temperature: param.NewOpt(0.1),
			TopP:        param.NewOpt(0.9),
		})

		for stream.Next() {
			event := stream.Current()
			switch event.Type {
			case "response.created":
				outputChannel <- &messages.ResponseStarted{}

			case "response.completed":
				// usage := event.AsResponseCompleted().Response.Usage
				// cb(&messages.ResponseStatsMessage{
				// 	InputTokens:     usage.InputTokens,
				// 	OutputTokens:    usage.OutputTokens,
				// 	ReasoningTokens: usage.OutputTokensDetails.ReasoningTokens,
				// 	TotalTokens:     usage.TotalTokens,
				// })
				outputChannel <- messages.NewResponseCompleted()

			case "response.failed":
				resp := event.AsResponseFailed().Response
				err := resp.Error
				// usage := resp.Usage
				// cb(&messages.ResponseStatsMessage{
				// 	InputTokens:     usage.InputTokens,
				// 	OutputTokens:    usage.OutputTokens,
				// 	ReasoningTokens: usage.OutputTokensDetails.ReasoningTokens,
				// 	TotalTokens:     usage.TotalTokens,
				// })
				outputChannel <- messages.NewResponseFailed(fmt.Sprintf("%s", err.Code), err.Message)

			case "response.output_item.done":
				if event.AsResponseOutputItemAdded().Item.Type == "function_call" {
					item := event.AsResponseOutputItemAdded().Item
					outputChannel <- messages.NewToolCall(
						item.Name,
						item.Arguments.OfString,
						item.CallID,
					)
				}

			case "response.output_text.delta":
				if event.Delta != "" {
					outputChannel <- messages.NewText(
						event.Delta,
						messages.AGENT,
					)
				}

			case "response.reasoning_text.delta":
				if event.Delta != "" {
					outputChannel <- messages.NewReasoning(
						event.Delta,
					)
				}

			}
		}

		if err := stream.Err(); err != nil {
			llm.LLMLogger.Error("stream error", "error", err)
			if strings.Contains(err.Error(), "404") {
				outputChannel <- messages.NewResponseFailed(
					"404",
					fmt.Sprintf("Model %s not found", model),
				)
			} else {
				outputChannel <- messages.NewResponseFailed(
					"400",
					err.Error(),
				)
			}
		}
	})()

	return outputChannel, nil
}

func NewOpenAICompatibleProviderAdapter(baseURL, APIKey string, Models []string) *OpenAICompatibleProviderAdapter {
	client := openai.NewClient(option.WithAPIKey(APIKey), option.WithBaseURL(baseURL))
	return &OpenAICompatibleProviderAdapter{
		Client: &client,
		Models: Models,
	}

}
