package controller

import "context"

type IController interface {
	Start(ctx context.Context, threads int)error
	Stop(stopCh <- chan struct{})error

	AddHook(hook IHook)error
	RemoveHook(hook IHook)error
}