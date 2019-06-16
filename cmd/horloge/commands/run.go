/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/majordome-iot/horloge"
	"github.com/spf13/cobra"
)

func getServer() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/ping", horloge.HTTPHandlerPing())
	e.GET("/health_check", horloge.HTTPHandlerHealthCheck())
	e.GET("/version", horloge.HTTPHandlerVersion())
	e.POST("/jobs", horloge.HTTPHandlerRegisterJob(runner))
	e.GET("/jobs", horloge.HTTPHandlerListJobs(runner))
	e.GET("/jobs/:id", horloge.HTTPHandlerJobDetail(runner))
	e.DELETE("/jobs/:id", horloge.HTTPHandlerDeleteJob(runner))

	return e
}

func getSync() horloge.Sync {
	switch s := sync; s {
	case "redis":
		fmt.Printf("Syncing with redis %s with db %d \n", redisAddr, redisDb)

		return horloge.NewSyncRedis(runner, redisAddr, redisPasswd, redisDb)
	case "file":
		fmt.Printf("Syncing with file: %s\n", filePath)

		return horloge.NewSyncDisk(runner, filePath)
	default:
		fmt.Println("No sync")
		return horloge.NewSyncNone()
	}
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Creates an horloge runner with a web interface",
	Long:  `Something`,
	Run: func(cmd *cobra.Command, args []string) {
		runner = horloge.NewRunner()
		runner.Sync(getSync())
		e := getServer()

		go func() {
			addr := fmt.Sprintf("%s:%d", bind, port)
			fmt.Printf("🕒 Horloge v%s\n", horloge.Version)
			fmt.Printf("Http server powered by Echo v%s\n", echo.Version)
			e.Logger.Fatal(e.Start(addr))
		}()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan

		fmt.Println("Shutdown signal received, exiting...")
		e.Shutdown(context.Background())

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&bind, "bind", "b", "127.0.0.1", "Bind to")
	runCmd.Flags().IntVarP(&port, "port", "p", 6432, "Port to listen to")

	// sync options
	runCmd.Flags().StringVarP(&sync, "sync", "s", "none", "Sync method to use")

	// sync file options
	runCmd.Flags().StringVar(&filePath, "file-path", "./horloge.json", "File path to sync to (used with `file` sync)")
}
