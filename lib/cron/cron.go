package cron

import "context"

type selfCron struct {
}

type selfExe struct {
	schedule string
	exeFunc  func(ctx context.Context)
}

func init() {
	
}
