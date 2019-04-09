package main

import (
	"go-learn/mqtt-demo/mqtt"
	"log"
	"time"
)

//MqttRequestClient 全局mqtt服务对象
var MqttRequestClient mqtt.Client

func main() {
	MqttRequestClient = mqtt.NewMQTTClient("123.207.34.70:1880", "easygo", "easygo123", "test")
	if !MqttRequestClient.MqttConnect() {
		log.Fatal("mqtt_err", "Mqtt connect failed!")
	}
	time.Sleep(time.Second * 3)
}
