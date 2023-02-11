package job

import "github.com/robfig/cron/v3"

// 运行
func Run() error {
	cron := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	cron.Start()

	return nil
}
