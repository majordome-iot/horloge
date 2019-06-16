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

	"github.com/go-redis/redis"
	"github.com/majordome-iot/horloge"
	"github.com/spf13/cobra"
)

var websocketBridgeCmd = &cobra.Command{
	Use:   "websocket-bridge",
	Short: "A websocket server that forwards messages from redis",
	Long: `Starts a http server that subscribes to horloge/job on redis
and forwards each message it receives.
Use / with a Websocket client to receive messages.
Use /ping to receive a pong response.`,
	Run: func(cmd *cobra.Command, args []string) {
		redisAddr, _ := rootCmd.PersistentFlags().GetString("redis-addr")
		redisPasswd, _ := rootCmd.PersistentFlags().GetString("redis-passwd")
		redisDb, _ := rootCmd.PersistentFlags().GetInt("redis-db")
		bind, _ := cmd.Flags().GetString("bind")

		fmt.Printf("Connecting to redis server at %s on database %d...\n", redisAddr, redisDb)
		fmt.Printf("Websocket server listening on %s...\n", bind)

		httpServer := horloge.NewWebsocketServer()
		redisClient := horloge.NewRedisClient(redisAddr, redisPasswd, redisDb)
		signalChan := make(chan os.Signal, 1)

		redisClient.AddPublishHandler(func(msg *redis.Message) {
			httpServer.Publish(msg.Payload)
		})

		redisClient.AddPublishHandler(func(msg *redis.Message) {
			httpServer.Publish(msg.Payload)
		})

		go redisClient.Wait()

		httpServer.Run(bind)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
	},
}

func init() {
	rootCmd.AddCommand(websocketBridgeCmd)
	websocketBridgeCmd.Flags().StringP("bind", "b", ":5000", "address on which the websocket server should listen to")
}
