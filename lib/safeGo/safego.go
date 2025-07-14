package safeGo

import (
	"context"
	"fmt"
	"runtime/debug"
)

func SafeGo(ctx context.Context, fn func(context.Context)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("handle safe go unexpected error:", r)
				debug.PrintStack()
			}
		}()
		fn(ctx)
	}()
}
