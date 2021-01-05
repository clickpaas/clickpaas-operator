package controller

import "fmt"

type BaseController struct {
	hooks []IHook
}

func(bc *BaseController)GetHooks()[]IHook{
	return bc.hooks
}

func(bc *BaseController)AddHook(hook IHook)error{
	for _,h := range bc.hooks{
		if h == hook{
			return fmt.Errorf("hook has installed, add failed")
		}
	}
	bc.hooks = append(bc.hooks, hook)
	return nil
}


func(bc *BaseController)RemoveHook(hook IHook)error{
	for i, h := range bc.hooks{
		if h == hook{
			bc.hooks = append(bc.hooks[:i], bc.hooks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("hook is not installed, remove failed")
}
