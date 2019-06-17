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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/majordome-iot/horloge"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bind string
var port int
var sync string

func getServer() *gin.Engine {
	r := gin.Default()

	// Routes
	r.GET("/ping", horloge.HTTPHandlerPing())
	r.GET("/health_check", horloge.HTTPHandlerHealthCheck())
	r.GET("/version", horloge.HTTPHandlerVersion())
	r.POST("/jobs", horloge.HTTPHandlerRegisterJob(runner))
	r.GET("/jobs", horloge.HTTPHandlerListJobs(runner))
	r.GET("/jobs/:id", horloge.HTTPHandlerJobDetail(runner))
	r.DELETE("/jobs/:id", horloge.HTTPHandlerDeleteJob(runner))

	return r
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs an horloge runner with a web interface",
	Long: `Starts a web API with the following routes:

GET /jobs returns a list of jobs
POST /jobs creates a new job
GET /jobs/{id} returns job with id {id}
DELETE /jobs/{id} delets job with id {id}`,
	Run: func(cmd *cobra.Command, args []string) {
		runner = horloge.NewRunner()

		if sync != "" {
			runner.Sync(horloge.NewSyncRedis(runner, redisAddr, redisPasswd, redisDB))
		}

		r := getServer()

		go func() {
			addr := fmt.Sprintf("%s:%d", bind, port)
			fmt.Printf("🕒 Horloge v%s\n", horloge.Version)
			fmt.Printf("Http server powered by Gin %s\n", gin.Version)
			r.Run(addr)
		}()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan

		fmt.Println("Shutdown signal received, exiting...")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&bind, "bind", "b", "127.0.0.1", "Addr to listen to")
	runCmd.Flags().IntVarP(&port, "port", "p", 6432, "Port to listen on")
	runCmd.Flags().StringVarP(&sync, "sync", "s", "", "Sync method to use")

	viper.BindPFlag("run.bind", mqttBridgeCmd.Flags().Lookup("bind"))
	viper.BindPFlag("run.port", mqttBridgeCmd.Flags().Lookup("port"))
	viper.BindPFlag("run.sync", mqttBridgeCmd.Flags().Lookup("sync"))
}
