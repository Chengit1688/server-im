package mqtt

import (
	"fmt"
	"im/config"
	"im/pkg/logger"
	"im/pkg/util"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var defaultEMQXClientManager *EMQXClientManager
var MqttClient *EMQXClient

type EMQXData struct {
	OperationID string
	Topic       string
	Payload     string
}

type EMQXClient struct {
	client   mqtt.Client
	address  string
	username string
	password string
	clientID string
	ch       chan *EMQXData
}

//func NewEMQXClient(address string, username string, password string, ch chan *EMQXData) *EMQXClient {
func NewEMQXClient(address string, username string, password string) *EMQXClient {
	if MqttClient != nil {
		return MqttClient
	}
	station := config.Config.Station
	//clientID := fmt.Sprintf("%s_system_%s%d", station, util.RandID(10), time.Now().Nanosecond())
	clientID := fmt.Sprintf("%s_system_%s%d", station, util.RandID(10), time.Now().Nanosecond())
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", address))
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)

	opts.OnConnect = func(client mqtt.Client) {
		logger.Sugar.Debugw(clientID, "func", "OnConnect", "info", fmt.Sprintf("emqx client connected"))
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		logger.Sugar.Debugw(clientID, "func", "OnConnectionLost", "info", fmt.Sprintf("emqx client disconnect, error: %v", err))
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	MqttClient = &EMQXClient{
		address:  address,
		username: username,
		password: password,
		clientID: clientID,
		//clientID: fmt.Sprintf("%s_system", station),
		//ch:       ch,
		client: client,
	}
	return MqttClient
}

func (c *EMQXClient) Publish() {
	//opts := mqtt.NewClientOptions()
	//opts.AddBroker(fmt.Sprintf("tcp://%s", c.address))
	//opts.SetClientID(c.clientID)
	//opts.SetUsername(c.username)
	//opts.SetPassword(c.password)
	//opts.SetKeepAlive(60 * time.Second)
	//opts.SetPingTimeout(10 * time.Second)
	//
	//opts.OnConnect = func(client mqtt.Client) {
	//	logger.Sugar.Debugw(c.clientID, "func", "OnConnect", "info", fmt.Sprintf("emqx client connected"))
	//}
	//opts.OnConnectionLost = func(client mqtt.Client, err error) {
	//	logger.Sugar.Debugw(c.clientID, "func", "OnConnectionLost", "info", fmt.Sprintf("emqx client disconnect, error: %v", err))
	//}
	//
	//client := mqtt.NewClient(opts)
	//if token := client.Connect(); token.Wait() && token.Error() != nil {
	//	panic(token.Error())
	//}

	for {
		// 向主题发送消息
		data := <-c.ch
		if token := c.client.Publish(data.Topic, 0, false, data.Payload); token.Wait() && token.Error() != nil {
			logger.Sugar.Errorw(data.OperationID, "func", "Publish", "error", fmt.Sprintf("publish error, client id: %s, topic: %s, error: %v", c.clientID, data.Topic, token.Error()))
		}
		logger.Sugar.Debugw(data.OperationID, "func", "Publish", "info", fmt.Sprintf("	publish success, client id: %s, topic: %s, data: %s", c.clientID, data.Topic, data.Payload))
	}
}

func (c *EMQXClient) Subscribe(callback mqtt.MessageHandler, topics []string) {
	//opts := mqtt.NewClientOptions()
	//opts.AddBroker(fmt.Sprintf("tcp://%s", c.address))
	//opts.SetClientID(c.clientID)
	//opts.SetUsername(c.username)
	//opts.SetPassword(c.password)
	//opts.SetKeepAlive(60 * time.Second)
	//opts.SetPingTimeout(10 * time.Second)
	//
	//opts.OnConnect = func(client mqtt.Client) {
	//	logger.Sugar.Debugw(c.clientID, "func", "OnConnect", "info", fmt.Sprintf("emqx client connected"))
	//}
	//opts.OnConnectionLost = func(client mqtt.Client, err error) {
	//	logger.Sugar.Debugw(c.clientID, "func", "OnConnectionLost", "info", fmt.Sprintf("emqx client disconnect, error: %v", err))
	//}
	//
	//client := mqtt.NewClient(opts)
	//if token := client.Connect(); token.Wait() && token.Error() != nil {
	//	panic(token.Error())
	//}
	//client := NewEMQXClient(cfg.MQTTAddress, cfg.MQTTUsername, cfg.MQTTPassword)
	for _, topic := range topics {
		if token := c.client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
			logger.Sugar.Debugw("", "func", "Subscribe", topic, "error", "token error", token.Error())
		}
		logger.Sugar.Debugw("", "func", "Subscribe", topic, "done")
	}
	select {}
}

type EMQXClientManager struct {
	channelNum int
	chs        []chan *EMQXData
}

func (m *EMQXClientManager) Start() {
	cfg := config.Config.EMQXServer

	m.channelNum = cfg.MQTTMaxConnection
	if m.channelNum == 0 {
		m.channelNum = 10
	}
	client := NewEMQXClient(cfg.MQTTAddress, cfg.MQTTUsername, cfg.MQTTPassword)
	for i := 0; i < m.channelNum; i++ {
		ch := make(chan *EMQXData, 10000)
		m.chs = append(m.chs, ch)

		logger.Sugar.Infow("", "info", "new emqx client")
		//client := NewEMQXClient(cfg.MQTTAddress, cfg.MQTTUsername, cfg.MQTTPassword, m.chs[i])
		client.ch = m.chs[i]
		go client.Publish()
		logger.Sugar.Infow("", "info", fmt.Sprintf("emqx client start success, client id: %s", client.clientID))
	}
	//订阅topic

}

func (m *EMQXClientManager) Publish(operationID string, topic string, payload string) {
	data := &EMQXData{}
	data.OperationID = operationID
	data.Topic = topic
	data.Payload = payload

	// 相同topic发送到一个channel中
	index := util.StringHash(topic) % uint32(m.channelNum)
	m.chs[index] <- data
}

func Init() {
	defaultEMQXClientManager = new(EMQXClientManager)
	defaultEMQXClientManager.Start()
}
