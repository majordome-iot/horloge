/*
Copyright © 2019 Samori Gorse <samorigorse+github@gmail.com>

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
	"github.com/shinuza/horloge"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var websocketBridgeCmd = &cobra.Command{
	Use:   "websocket-bridge",
	Short: "A websocket server that forwards messages from redis",
	Long: `Starts a http server that subscribes to horloge/job on redis
and forwards each message it receives.
Use / with a Websocket client to receive messages.
Use /ping to receive a pong response.`,
	Run: func(cmd *cobra.Command, args []string) {
		bind := viper.GetString("websocket.bind")
		port := viper.GetInt("websocket.port")
		addr := fmt.Sprintf("%s:%d", bind, port)

		fmt.Printf("Connecting to redis server at %s on database %d...\n", redisAddr, redisDB)
		fmt.Printf("Websocket server listening on %s...\n", bind)

		httpServer := horloge.NewWebsocketServer()
		redisClient := horloge.NewRedisClient(redisAddr, redisPasswd, redisDB)
		signalChan := make(chan os.Signal, 1)

		redisClient.AddPublishHandler(func(msg *redis.Message) {
			httpServer.Publish(msg.Payload)
		})

		go redisClient.Wait()

		httpServer.Run(addr)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
	},
}

func init() {
	rootCmd.AddCommand(websocketBridgeCmd)

	websocketBridgeCmd.Flags().StringP("bind", "b", "127.0.0.1", "Addr to listen to")
	websocketBridgeCmd.Flags().IntP("port", "p", 5000, "Port to listen to")

	viper.BindPFlag("websocket.bind", websocketBridgeCmd.Flags().Lookup("bind"))
	viper.BindPFlag("websocket.port", websocketBridgeCmd.Flags().Lookup("port"))
}
