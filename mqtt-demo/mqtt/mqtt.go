package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type Client struct {
	Client          MQTT.Client
	ClientOpts      MQTT.ClientOptions
	Host            string
	Username        string
	Password        string
	ClientID        string
	justPublish     bool
	lostFunction    []lostFunction
	allLostFunction func()
}

type lostFunction struct {
	topic     string
	payloadCh chan string
}

func NewMQTTClient(host, username, password, clientID string) Client {
	c := Client{}
	c.Host = host
	c.Username = username
	c.Password = password
	c.ClientID = clientID
	c.justPublish = true
	c.lostFunction = []lostFunction{}
	return c
}

func (c *Client) formatLostFunction() func() bool {
	return func() bool {
		if !c.justPublish {
			for _, v := range c.lostFunction {
				if !c.subscribe(v.topic, v.payloadCh) {
					return false
				}
			}
			return true
		} else {
			return true
		}
	}
}

func (c *Client) AddLostFunction(topic string, payloadCh chan string) {
	c.justPublish = false
	t := lostFunction{
		topic:     topic,
		payloadCh: payloadCh,
	}
	c.lostFunction = append(c.lostFunction, t)
}

func (c *Client) MqttConnect() bool {
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + c.Host).SetClientID(c.ClientID)
	opts.SetKeepAlive(time.Duration(30) * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(time.Duration(1) * time.Second)
	opts.SetUsername(c.Username)
	opts.SetPassword(c.Password)
	if !c.justPublish {
		var lostFunction MQTT.ConnectionLostHandler = func(c_ MQTT.Client, err_ error) {
			c.formatLostFunction()
		}
		opts.SetConnectionLostHandler(lostFunction)
	}
	c.Client = MQTT.NewClient(opts)
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		return false
	} else {
		return c.formatLostFunction()()
	}
}

func (c *Client) subscribe(topic string, payloadCh chan string) bool {
	var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		var content = msg.Payload()
		payload := string(content[:])
		var data = payload
		payloadCh <- data
	}
	if token := c.Client.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
		return false
	} else {
		return true
	}
}

func (c *Client) Subscribe(topic string, qos byte, callback MQTT.MessageHandler) bool {
	if token := c.Client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		return false
	} else {
		return true
	}
}

func (c *Client) LoopProcessData(processDataFunc func(d string), payloadCh chan string) {
	for {
		var data = <-payloadCh
		go processDataFunc(data)
	}
}

func (c *Client) Publish(topic string, qos byte, retained bool, payload interface{}) bool {
	if token := c.Client.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		return false
	} else {
		return true
	}
}

func (c *Client) Unsubscribe(topicName string, do func()) {
	token := c.Client.Unsubscribe(topicName)
	if token.Wait() {
		do()
	} else {
		time.Sleep(time.Second * 5)
		c.Unsubscribe(topicName, do)
	}
}
