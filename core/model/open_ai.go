package model

import (
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/exceptions"
	"github.com/smtdfc/nagare/core/messages"
	"github.com/smtdfc/nagare/core/tool"
)

type OpenAICompatibleChatModel struct {
	Config *ChatModelConfig
}

func (o *OpenAICompatibleChatModel) Transform(msg messages.Message) (responses.ResponseInputItemUnionParam, error) {
	switch t := msg.(type) {
	case *messages.TextMessage:
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

	case *messages.ToolResultMessage:
		return responses.ResponseInputItemUnionParam{
			OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
				CallID: t.CallID,
				Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
					OfString: openai.String(t.Result),
				},
			},
		}, nil

	case *messages.ToolCallMessage:
		return responses.ResponseInputItemUnionParam{
			OfFunctionCall: &responses.ResponseFunctionToolCallParam{
				CallID:    t.CallID,
				Name:      t.FunctionName,
				Arguments: t.Args,
			},
		}, nil
	}

	return responses.ResponseInputItemUnionParam{}, nil
}

func (o *OpenAICompatibleChatModel) TransformToolDeclarations(tools tool.ToolMap) ([]responses.ToolUnionParam, error) {
	toolParams := []responses.ToolUnionParam{}

	for _, tool := range tools {
		var params map[string]any

		// println(tool.GetArgumentsSchema())
		err := json.Unmarshal([]byte(tool.GetArgumentsSchema()), &params)
		if err != nil {
			return nil, exceptions.NewToolException(fmt.Sprintf("JSON parse error %s: %v\n", tool.GetName(), err), tool.GetName())
		}

		toolParams = append(toolParams, responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        tool.GetName(),
				Description: openai.String(tool.GetDesc()),
				Parameters:  params,
			},
		})
	}
	return toolParams, nil
}

// chat implements [ChatModel].
func (o *OpenAICompatibleChatModel) Chat(ctx context.ExecuteContext, history messages.ListMessage, cb MessageCallback) error {
	client := openai.NewClient(
		option.WithAPIKey(o.Config.APIKey),
		option.WithBaseURL(o.Config.BaseURL),
	)

	inputs := responses.ResponseInputParam{}
	listTool, err := o.TransformToolDeclarations(ctx.Tools)
	if err != nil {
		return err
	}
	for _, msg := range history {
		input, err := o.Transform(msg)
		if err != nil {
			return err
		}
		inputs = append(inputs, input)
	}

	stream := client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
		Model: o.Config.Model,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: inputs,
		},

		Tools: listTool,
	})

	for stream.Next() {
		event := stream.Current()
		switch event.Type {
		case "response.created":
			cb(&messages.ResponseCreatedMessage{})

		case "response.completed":
			cb(&messages.ResponseCompletedMessage{})

		case "response.failed":
			respErr := event.AsResponseFailed().Response.Error
			cb(&messages.ResponseFailedMessage{
				Code:    fmt.Sprintf("%s", respErr.Code),
				Message: respErr.Message,
			})

		case "response.output_item.done":
			if event.AsResponseOutputItemAdded().Item.Type == "function_call" {
				item := event.AsResponseOutputItemAdded().Item
				cb(&messages.ToolCallMessage{
					CallID:       item.CallID,
					FunctionName: item.Name,
					Args:         item.Arguments.OfString,
				})
			}

		case "response.output_text.delta":
			if event.Delta != "" {
				cb(&messages.TextMessage{
					Role:    messages.AGENT,
					Content: event.Delta,
				})
			}

		case "response.reasoning_text.delta":
			if event.Delta != "" {
				cb(&messages.ReasoningMessage{
					Content: event.Delta,
				})
			}

		}
	}

	if err := stream.Err(); err != nil {
		appLogger.Info(fmt.Sprintf("Error when streaming data from LLM provider: %s", err), "error", err, "compatible", "open-ai")
		return exceptions.NewStreamException(err.Error())
	}

	return nil

}

func NewOpenAICompatibleClient(config *ChatModelConfig) ChatModel {
	return &OpenAICompatibleChatModel{
		Config: config,
	}
}
