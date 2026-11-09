package computer_use

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/external_module"
)

const MODULE_ID = "nagare_computer_use"
const BUFFER_SIZE = 30

type ComputerUseOrchestrator struct {
	channels *external_module.ModuleChannelPair[any, any]
}

func (c *ComputerUseOrchestrator) GetAppInstance(ctx context.Context, appName string) (string, error) {
	reqID := uuid.New().String()
	reply, err := external_module.RequestReply[GetAppInstanceRequest, GetAppInstanceResponse](
		ctx,
		c.channels.InvokeChan,
		c.channels.ReplyChan,
		GET_APP_INSTANCE,
		GetAppInstanceRequest{
			Name: appName,
		},
		reqID,
	)
	if err != nil {
		return "", err
	}

	if reply.Data.IsError {
		return "", errors.New(reply.Data.Error)
	}

	return reply.Data.InstanceID, nil
}

func (c *ComputerUseOrchestrator) ListRunningApps(ctx context.Context) ([]AppInfo, error) {
	reqID := uuid.New().String()
	reply, err := external_module.RequestReply[ListRunningAppsRequest, ListRunningAppsResponse](
		ctx,
		c.channels.InvokeChan,
		c.channels.ReplyChan,
		LIST_RUNNING_APPS,
		ListRunningAppsRequest{},
		reqID,
	)
	if err != nil {
		return nil, err
	}

	if reply.Data.IsError {
		return nil, errors.New(reply.Data.Error)
	}

	return reply.Data.Apps, nil
}

func (c *ComputerUseOrchestrator) OpenApp(ctx context.Context, name string, args []string) (*uint32, error) {
	reqID := uuid.New().String()
	reply, err := external_module.RequestReply[OpenAppRequest, OpenAppResponse](
		ctx,
		c.channels.InvokeChan,
		c.channels.ReplyChan,
		OPEN_APP,
		OpenAppRequest{
			Name: name,
			Args: args,
		},
		reqID,
	)
	if err != nil {
		return nil, err
	}

	if reply.Data.IsError {
		return nil, errors.New(reply.Data.Error)
	}

	return reply.Data.PID, nil
}

func (c *ComputerUseOrchestrator) GetUITree(ctx context.Context, instanceID string, maxDepth *int) (*UINode, string, error) {
	reqID := uuid.New().String()
	reply, err := external_module.RequestReply[GetUITreeRequest, GetUITreeResponse](
		ctx,
		c.channels.InvokeChan,
		c.channels.ReplyChan,
		GET_UI_TREE,
		GetUITreeRequest{
			InstanceID: instanceID,
			MaxDepth:   maxDepth,
		},
		reqID,
	)
	if err != nil {
		return nil, "", err
	}
	if reply.Data.IsError {
		return nil, "", errors.New(reply.Data.Error)
	}
	return reply.Data.Tree, reply.Data.Summary, nil
}

func (c *ComputerUseOrchestrator) ClickElement(ctx context.Context, instanceID, refID string) error {
	reqID := uuid.New().String()
	reply, err := external_module.RequestReply[ClickElementRequest, ClickElementResponse](
		ctx,
		c.channels.InvokeChan,
		c.channels.ReplyChan,
		CLICK_ELEMENT,
		ClickElementRequest{
			InstanceID: instanceID,
			RefID:      refID,
		},
		reqID,
	)
	if err != nil {
		return err
	}
	if reply.Data.IsError {
		return errors.New(reply.Data.Error)
	}
	return nil
}

func (c *ComputerUseOrchestrator) InputText(ctx context.Context, instanceID, refID, text string) error {
	reqID := uuid.New().String()
	reply, err := external_module.RequestReply[InputTextRequest, InputTextResponse](
		ctx,
		c.channels.InvokeChan,
		c.channels.ReplyChan,
		INPUT_TEXT,
		InputTextRequest{
			InstanceID: instanceID,
			RefID:      refID,
			Text:       text,
		},
		reqID,
	)
	if err != nil {
		return err
	}
	if reply.Data.IsError {
		return errors.New(reply.Data.Error)
	}
	return nil
}

func (c *ComputerUseOrchestrator) MouseAction(ctx context.Context, req MouseActionRequest) error {
	reqID := uuid.New().String()
	reply, err := external_module.RequestReply[MouseActionRequest, MouseActionResponse](
		ctx,
		c.channels.InvokeChan,
		c.channels.ReplyChan,
		MOUSE_ACTION,
		req,
		reqID,
	)
	if err != nil {
		return err
	}
	if reply.Data.IsError {
		return errors.New(reply.Data.Error)
	}
	return nil
}

func (c *ComputerUseOrchestrator) KeyboardShortcut(ctx context.Context, key string, modifiers []string) error {
	reqID := uuid.New().String()
	reply, err := external_module.RequestReply[KeyboardShortcutRequest, KeyboardShortcutResponse](
		ctx,
		c.channels.InvokeChan,
		c.channels.ReplyChan,
		KEYBOARD_SHORTCUT,
		KeyboardShortcutRequest{
			Key:       key,
			Modifiers: modifiers,
		},
		reqID,
	)
	if err != nil {
		return err
	}
	if reply.Data.IsError {
		return errors.New(reply.Data.Error)
	}
	return nil
}

func (c *ComputerUseOrchestrator) WindowControl(ctx context.Context, instanceID, action string) error {
	reqID := uuid.New().String()
	reply, err := external_module.RequestReply[WindowControlRequest, WindowControlResponse](
		ctx,
		c.channels.InvokeChan,
		c.channels.ReplyChan,
		WINDOW_CONTROL,
		WindowControlRequest{
			InstanceID: instanceID,
			Action:     action,
		},
		reqID,
	)
	if err != nil {
		return err
	}
	if reply.Data.IsError {
		return errors.New(reply.Data.Error)
	}
	return nil
}

func (c *ComputerUseOrchestrator) QueryElementByText(ctx context.Context, instanceID, text string, exact *bool) ([]ElementMatch, error) {
	reqID := uuid.New().String()
	reply, err := external_module.RequestReply[QueryElementByTextRequest, QueryElementByTextResponse](
		ctx,
		c.channels.InvokeChan,
		c.channels.ReplyChan,
		QUERY_ELEMENT_BY_TEXT,
		QueryElementByTextRequest{
			InstanceID: instanceID,
			Text:       text,
			Exact:      exact,
		},
		reqID,
	)
	if err != nil {
		return nil, err
	}
	if reply.Data.IsError {
		return nil, errors.New(reply.Data.Error)
	}
	return reply.Data.Matches, nil
}

// @Injectable
func NewComputerUseOrchestrator(externalModuleHost *external_module.ExternalModuleHost) (*ComputerUseOrchestrator, error) {
	invokeChan, replyChan, err := externalModuleHost.RegisterChannels(MODULE_ID, BUFFER_SIZE)
	if err != nil {
		return nil, err
	}

	return &ComputerUseOrchestrator{
		channels: &external_module.ModuleChannelPair[any, any]{
			InvokeChan: invokeChan,
			ReplyChan:  replyChan,
		},
	}, nil
}
