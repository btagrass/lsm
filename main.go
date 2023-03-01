package main

import (
	"embed"
	"fmt"
	"io/fs"
	"lsm/api"
	_ "lsm/docs"
	"lsm/job"
	"lsm/mgt"
	"net/http"

	"github.com/btagrass/go.core/cmd"
	"github.com/btagrass/go.core/htp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

//go:embed web/dist
var dist embed.FS

// 入口
func main() {
	cmd.Execute(
		"lsm",
		"直播流媒体系统",
		cmd.Install,
		cmd.Start,
		cmd.Status,
		cmd.Stop,
		cmd.Uninstall,
		&cobra.Command{
			Use:   "run",
			Short: "运行",
			Run: func(c *cobra.Command, args []string) {
				group := &errgroup.Group{}
				// 应用服务
				api := &http.Server{
					Addr:    fmt.Sprintf(":%d", htp.Port),
					Handler: api.Api(),
				}
				group.Go(api.ListenAndServe)
				// 管理服务
				mgt := &http.Server{
					Addr:    fmt.Sprintf(":%d", htp.Port+1),
					Handler: mgt.Mgt(),
				}
				group.Go(mgt.ListenAndServe)
				// 界面服务
				dist, _ := fs.Sub(dist, "web/dist")
				engine := gin.Default()
				engine.StaticFS("/", http.FS(dist))
				web := &http.Server{
					Addr:    fmt.Sprintf(":%d", htp.Port+2),
					Handler: engine,
				}
				group.Go(web.ListenAndServe)
				// 作业服务
				group.Go(job.Run)
				if err := group.Wait(); err != nil {
					logrus.Fatal(err)
				}
			},
		},
	)
}
