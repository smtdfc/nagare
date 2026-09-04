package adapters

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/message"
)

type OpenAICompatibleAdapter struct {
	BaseURL string
	Client  *openai.Client
	Models  []string
	logger  *logger.BaseLogger
}

func (o *OpenAICompatibleAdapter) TransformToProviderMessage(msg message.Message) (responses.ResponseInputItemUnionParam, error) {
	switch t := msg.(type) {
	case *message.TextMessage:
		item := &responses.ResponseInputItemMessageParam{
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
		case message.USER:
			item.Role = "user"
			return responses.ResponseInputItemUnionParam{
				OfInputMessage: item,
			}, nil

		case message.SYSTEM:
			item.Role = "system"
			return responses.ResponseInputItemUnionParam{
				OfInputMessage: item,
			}, nil
		case message.DEVELOPER:
			item.Role = "developer"
			return responses.ResponseInputItemUnionParam{
				OfInputMessage: item,
			}, nil
		case message.AGENT:
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

	case *message.ToolResultMessage:
		return responses.ResponseInputItemUnionParam{
			OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
				CallID: t.CallID,
				Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
					OfString: openai.String(t.Result),
				},
			},
		}, nil

	case *message.ToolCallMessage:
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

func (o *OpenAICompatibleAdapter) TransformToolDeclarations(tools tool.ListTool) ([]responses.ToolUnionParam, error) {
	toolParams := make([]responses.ToolUnionParam, len(tools))

	for i, tool := range tools {
		params, err := helpers.UnmarshalJson[map[string]any](tool.GetArgsSchema())

		if err != nil {
			return nil, custom_errors.ErrInvalidToolSchema
		}

		toolParams[i] = responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        tool.GetName(),
				Description: openai.String(tool.GetDescription()),
				Parameters:  *params,
			},
		}
	}
	return toolParams, nil
}

func (o *OpenAICompatibleAdapter) Send(ctx context.Context, model string, listMessage message.ListMessage, tools tool.ListTool) (message.MessageReadOnlyChannel, error) {
	if !slices.Contains(o.Models, model) {
		return nil, custom_errors.ErrModelNotSupportedByProvider
	}

	inputs := responses.ResponseInputParam{}
	listTool, err := o.TransformToolDeclarations(tools)
	if err != nil {
		o.logger.Error("TransformToolDeclarations error", "error", err)
		return nil, err
	}

	for _, msg := range listMessage {
		input, err := o.TransformToProviderMessage(msg)
		if err != nil {
			o.logger.Error("TransformToProviderMessage error", "error", err)
			return nil, err
		}
		inputs = append(inputs, input)
	}

	outputChannel := make(chan message.Message)

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
				outputChannel <- message.NewResponseFailedMessage(fmt.Sprintf("%s", err.Code), err.Message)
			case "response.output_item.done":
				if event.AsResponseOutputItemAdded().Item.Type == "function_call" {
					item := event.AsResponseOutputItemAdded().Item
					outputChannel <- message.NewToolCallMessage(
						item.CallID,
						item.Name,
						item.Arguments.OfString,
					)
				}

			case "response.output_text.delta":
				if event.Delta != "" {
					outputChannel <- message.NewTextMessage(
						message.AGENT,
						event.Delta,
					)
				}

			case "response.reasoning_text.delta":
				if event.Delta != "" {
					outputChannel <- message.NewReasoningMessage(
						event.Delta,
					)
				}
			}
		}

		if err := stream.Err(); err != nil {
			o.logger.Error("stream error", "error", err)
			if strings.Contains(err.Error(), "404") {
				outputChannel <- message.NewResponseFailedMessage(
					"404",
					fmt.Sprintf("Model %s not found", model),
				)
			} else if strings.Contains(err.Error(), "429") {
				outputChannel <- message.NewResponseFailedMessage(
					"429",
					fmt.Sprintf("Quota exceed: %s", err.Error()),
				)
			} else {
				outputChannel <- message.NewResponseFailedMessage(
					"400",
					err.Error(),
				)
			}
		}
	})()

	return outputChannel, nil
}

func (o *OpenAICompatibleAdapter) GetModels(ctx context.Context) ([]string, error) {
	resp, err := o.Client.Models.List(ctx)
	if err != nil {
		o.logger.Error("list models error", "error", err)
		return nil, custom_errors.ErrGetListModelFailed
	}
	data := resp.Data

	result := make([]string, 0, len(data))
	for _, model := range data {
		result = append(result, model.ID)
	}
	return result, nil
}

func NewOpenAICompatibleAdapter(baseURL, APIKey string, Models []string, logger *logger.BaseLogger) *OpenAICompatibleAdapter {
	client := openai.NewClient(option.WithAPIKey(APIKey), option.WithBaseURL(baseURL))

	return &OpenAICompatibleAdapter{
		BaseURL: baseURL,
		Client:  &client,
		Models:  Models,
		logger:  logger.With("module", "openai-compatible-adapter"),
	}
}
