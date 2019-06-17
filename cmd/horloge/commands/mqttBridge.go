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
	"github.com/spf13/viper"
)

var mqttBridgeCmd = &cobra.Command{
	Use:   "mqtt-bridge",
	Short: "A bridge between redis and a mqtt broker",
	Long: `Creates a bridge that subscribes to horloge/job on redis
and forwards messages to the configured mqtt broker`,
	Run: func(cmd *cobra.Command, args []string) {
		mqttAddr := viper.GetString("mqtt.addr")

		fmt.Printf("Connecting to redis server at %s on database %d...\n", redisAddr, redisDB)
		fmt.Printf("Connecting to mqtt broker at %s...\n", mqttAddr)

		mqttClient := horloge.NewMQTTClient(mqttAddr)
		redisClient := horloge.NewRedisClient(redisAddr, redisPasswd, redisDB)
		signalChan := make(chan os.Signal, 1)

		redisClient.AddPublishHandler(func(msg *redis.Message) {
			mqttClient.Publish(horloge.MQTT_CHANNEL, msg.Payload)
		})

		go redisClient.Wait()

		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
	},
}

func init() {
	rootCmd.AddCommand(mqttBridgeCmd)

	mqttBridgeCmd.Flags().String("mqtt-addr", "tcp://localhost:1883", "address of the mqtt broker")

	viper.BindPFlag("mqtt.addr", mqttBridgeCmd.Flags().Lookup("mqtt-addr"))
}
