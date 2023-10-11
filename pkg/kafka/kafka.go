package kafka

import (
	"sync"

	"github.com/Shopify/sarama"
)

const TopicOperateLogs = "operate_logs"
const TopicApiLogs = "api_logs"
const TopicErrorLogs = "error_logs"

var (
	wg sync.WaitGroup
	Lp *LogProducerManager
)

func Init() {
	// cfg := config.Config.Log
	// Lp = new(LogProducerManager)
	// Lp.InitLogProducer()
	// go InitLogConsumer(TopicOperateLogs, cfg.KafkaConsumerNum)
	// go InitLogConsumer(TopicApiLogs, cfg.KafkaConsumerNum)
	// go InitLogConsumer(TopicErrorLogs, cfg.KafkaConsumerNum)
}

// type MConsumerGroup struct {
// 	sarama.ConsumerGroup
// 	groupID string
// 	topics  []string
// }

// type MConsumerGroupConfig struct {
// 	KafkaVersion   sarama.KafkaVersion
// 	OffsetsInitial int64
// 	IsReturnErr    bool
// }

// func InitKafka() {
// 	go InitKafkaGroup()
// 	// cfg := config.Config.Kafka
// 	// config := sarama.NewConfig()
// 	// config.Net.DialTimeout = 3 * time.Second
// 	// config.Consumer.Offsets.AutoCommit.Enable = false
// 	// consumer, err := sarama.NewConsumer([]string{cfg.KafkaAddress}, config)
// 	// if err != nil {
// 	// 	panic(fmt.Sprintf("kafka启动失败, err: %v", err))
// 	// }
// 	// logger.Sugar.Info("kafka 链接成功")
// 	// go ConsumerHistory(consumer)

// }

// func ConsumerHistory(consumer sarama.Consumer) error {
// 	// 根据topic获取所有的分区列表
// 	topic := config.Config.Kafka.KafkaHistoryTopic
// 	partitionList, err := consumer.Partitions(topic)
// 	if err != nil {
// 		logger.Sugar.Error(fmt.Sprintf("fail to get list of partition,err: %v", err))
// 		return err
// 	}
// 	logger.Sugar.Debug(fmt.Sprintf("获取分区成功 %s %v", topic, partitionList))
// 	// 遍历所有的分区
// 	for p := range partitionList {
// 		//针对每一个分区创建一个对应分区的消费者
// 		pc, err := consumer.ConsumePartition(topic, int32(p), sarama.OffsetNewest)
// 		if err != nil {
// 			logger.Sugar.Error(fmt.Sprintf("failed to start consumer for partition %d,err:%v\n", p, err))
// 		}
// 		logger.Sugar.Debug("已订阅消费者")
// 		defer pc.AsyncClose()
// 		wg.Add(1)
// 		//异步从每个分区消费信息
// 		go func(sarama.PartitionConsumer) {
// 			for msg := range pc.Messages() {
// 				logger.Sugar.Debug(fmt.Sprintf("partition:%d Offse:%d Key:%v Value:%s \n",
// 					msg.Partition, msg.Offset, msg.Key, msg.Value))
// 				err := useCase.MessageUseCase.RecvMqttMsg(msg.Value)
// 				if err == nil {
// 					// consumer.
// 					//标记已处理
// 				} else {
// 					//标记未处理
// 				}
// 			}
// 		}(pc)
// 	}
// 	wg.Wait()
// 	logger.Sugar.Debug("已关闭消费者")
// 	return nil
// }

// type consumerGroupHandler struct {
// 	name string
// }

// func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
// func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
// func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
// 	logger.Sugar.Debug(fmt.Sprintf("订阅等待消息"))
// 	msgChan := make(chan []byte, 1)
// 	isDone := make(chan bool, 1)
// 	defer func(isDone chan bool) {
// 		logger.Sugar.Debug(fmt.Sprintf("发送关闭信号"))
// 		isDone <- true
// 	}(isDone)
// 	go func(msgChan chan []byte) {
// 		msgs := [][]byte{}
// 		timeTicker := time.NewTicker(1 * time.Second)
// 		defer timeTicker.Stop()
// 		for {
// 			select {
// 			case <-isDone:
// 				logger.Sugar.Debug(fmt.Sprintf("收到关闭信号"))
// 				return
// 			case msg := <-msgChan:
// 				msgs = append(msgs, msg)
// 				// if len(msgs) >= 100 {
// 				// 	for {
// 				// 		err := useCase.MessageUseCase.RecvMqttMsgBatch(msgs)
// 				// 		if err == nil {
// 				// 			sess.Commit()
// 				// 			msgs = [][]byte{}
// 				// 			logger.Sugar.Debug(fmt.Sprintf("定时写入完成"))
// 				// 			break
// 				// 		} else {
// 				// 			logger.Sugar.Debug(fmt.Sprintf("定时写入异常%v", err))
// 				// 			time.Sleep(1 * time.Second)
// 				// 		}

// 				// 	}
// 				// }
// 			case <-timeTicker.C:
// 				for {
// 					if len(msgs) == 0 {
// 						break
// 					}
// 					err := useCase.MessageUseCase.RecvMqttMsgBatch(msgs)
// 					if err == nil {
// 						sess.Commit()
// 						msgs = [][]byte{}
// 						logger.Sugar.Debug(fmt.Sprintf("定时写入完成"))
// 						break
// 					} else {
// 						logger.Sugar.Debug(fmt.Sprintf("定时写入异常%v", err))
// 						time.Sleep(1 * time.Second)
// 					}

// 				}

// 			}
// 		}

// 	}(msgChan)
// 	for msg := range claim.Messages() {
// 		logger.Sugar.Debug(fmt.Sprintf("topic:%s Offse:%d \n", msg.Topic, msg.Offset))
// 		// useCase.MessageUseCase.RecvMqttMsg(msg.Value)
// 		// sess.Commit()
// 		msgChan <- msg.Value
// 		sess.MarkMessage(msg, "")

// 	}

// 	return nil
// }

// func InitKafkaGroup() {

// 	cfg := config.Config.Kafka
// 	config := sarama.NewConfig()
// 	config.Consumer.Return.Errors = true
// 	// config.Version = sarama.V0_11_0_2
// 	config.Consumer.Offsets.AutoCommit.Enable = false
// 	config.Consumer.Offsets.Initial = sarama.OffsetOldest
// 	config.Net.DialTimeout = 3 * time.Second

// 	group, err := sarama.NewConsumerGroup([]string{cfg.KafkaAddress}, cfg.KafkaHistoryTopic+"123", config)
// 	//sarama.NewConsumerGroupFromClient()
// 	if err != nil {
// 		panic(err)
// 	}
// 	logger.Sugar.Info("kafka 链接成功")
// 	defer group.Close()

// 	for {
// 		handler := consumerGroupHandler{name: cfg.KafkaHistoryTopic}
// 		logger.Sugar.Debug("重试链接 kafka", cfg.KafkaHistoryTopic)
// 		err := group.Consume(context.Background(), []string{cfg.KafkaHistoryTopic}, handler)
// 		if err != nil {
// 			logger.Sugar.Error(err.Error())
// 			time.Sleep(1 * time.Second)
// 		}
// 	}
// }

// // InitByTest 初始化数据库配置，仅用在单元测试
// func InitByTest() {
// 	fmt.Println("init kafka")
// 	logger.Target = logger.Console

// 	InitKafka()

// }
type PushLog struct {
	Topic string
	Data  string
}

type LogProducerManager struct {
	channelNum int
	chs        []chan *PushLog
	curr       int
}

type LogProducer struct {
	client sarama.SyncProducer
	ch     chan *PushLog
}

// func (m *LogProducerManager) InitLogProducer() {
// 	cfg := config.Config

// 	m.channelNum = cfg.Log.KafkaProducerNum
// 	if m.channelNum == 0 {
// 		m.channelNum = 1
// 	}

// 	for i := 0; i < m.channelNum; i++ {
// 		ch := make(chan *PushLog, 10000)
// 		m.chs = append(m.chs, ch)

// 		client := newLogProducer([]string{cfg.Kafka.KafkaAddress}, m.chs[i])
// 		go client.SendMessageToTopic()
// 	}
// }

// func (m *LogProducerManager) Push(topic, value string) {
// 	data := PushLog{Topic: topic, Data: value}
// 	if m.curr >= m.channelNum {
// 		m.curr = 0
// 	} else {
// 		m.curr++
// 	}
// 	index := m.curr % m.channelNum
// 	m.chs[index] <- &data
// }

// func newLogProducer(brokerList []string, ch chan *PushLog) LogProducer {
// 	config := sarama.NewConfig()
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Retry.Max = 10
// 	config.Producer.Return.Successes = true

// 	producer, err := sarama.NewSyncProducer(brokerList, config)
// 	if err != nil {
// 		logger.Sugar.Errorw("Failed to start Sarama producer:", "error", err)
// 	}
// 	instance := LogProducer{}
// 	instance.client = producer
// 	instance.ch = ch
// 	return instance
// }

// func (p *LogProducer) SendMessageToTopic() {
// 	for {
// 		data := <-p.ch
// 		partition, offset, err := p.client.SendMessage(&sarama.ProducerMessage{
// 			Topic: data.Topic,
// 			Value: sarama.StringEncoder(data.Data),
// 		})
// 		if err != nil {
// 			logger.Sugar.Errorw("Failed to SendMessage Sarama producer:", "error:", err, "topic:", data.Topic, "partition:", partition, "offset:", offset)
// 		}
// 	}

// }

// func InitLogConsumer(topic string, count int) {
// 	if count == 0 {
// 		count = 1
// 	}
// 	keepRunning := true
// 	cfg := config.Config.Kafka
// 	config := sarama.NewConfig()
// 	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
// 	config.Consumer.Offsets.Initial = sarama.OffsetOldest
// 	config.Consumer.Fetch.Default = 2048 // 每次从 Kafka 服务器读取 2048 条消息
// 	config.Consumer.Fetch.Min = 1        // 每次从 Kafka 服务器读取的最小消息数为 1 条
// 	consumer := Consumer{
// 		ready: make(chan bool),
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	wg := &sync.WaitGroup{}
// 	client, err := sarama.NewConsumerGroup([]string{cfg.KafkaAddress}, fmt.Sprintf("%s-group", topic), config)
// 	if err != nil {
// 		logger.Sugar.Errorw("Error creating consumer group client: ", "error", err)
// 		cancel()
// 		return
// 	}
// 	for i := 0; i < count; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			for {
// 				// `Consume` should be called inside an infinite loop, when a
// 				// server-side rebalance happens, the consumer session will need to be
// 				// recreated to get the new claims
// 				if err := client.Consume(ctx, []string{topic}, &consumer); err != nil {
// 					logger.Sugar.Errorw("Error from consumer: ", "error", err)
// 				}
// 				// check if context was cancelled, signaling that the consumer should stop
// 				if ctx.Err() != nil {
// 					return
// 				}
// 				consumer.ready = make(chan bool)
// 			}
// 		}()

// 		<-consumer.ready // Await till the consumer has been set up
// 		logger.Sugar.Info("Sarama consumer up and running!...")
// 	}

// 	sigterm := make(chan os.Signal, 1)
// 	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
// 	logger.Sugar.Info("Sarama consumer keepRunning!...")
// 	for keepRunning {
// 		select {
// 		case <-ctx.Done():
// 			logger.Sugar.Info("terminating: context cancelled")
// 			keepRunning = false
// 		case <-sigterm:
// 			logger.Sugar.Info("terminating: via signal")
// 			keepRunning = false
// 		}
// 	}
// 	logger.Sugar.Info("Sarama consumer cancel!...")
// 	cancel()
// 	logger.Sugar.Info("Sarama consumer Wait!...")
// 	wg.Wait()
// 	logger.Sugar.Info("Sarama consumer done!...")
// 	if err = client.Close(); err != nil {
// 		logger.Sugar.Errorf("Error closing client: %v", err)
// 	}
// }

// type Consumer struct {
// 	ready chan bool
// }

// // Setup is run at the beginning of a new session, before ConsumeClaim
// func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
// 	// Mark the consumer as ready
// 	close(consumer.ready)
// 	return nil
// }

// // Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
// 	return nil
// }

// // ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
// 	var topic string
// 	var record bool
// 	var messageArray []interface{}

// 	record = false
// 	for message := range claim.Messages() {
// 		if !record {
// 			topic = message.Topic
// 			record = true
// 		}
// 		//logger.Sugar.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
// 		session.MarkMessage(message, "")
// 		switch topic {
// 		case TopicOperateLogs:
// 			log := model.OperateLogs{}
// 			json.Unmarshal(message.Value, &log)
// 			messageArray = append(messageArray, log)
// 		case TopicApiLogs:
// 			log := model.ApiLogs{}
// 			json.Unmarshal(message.Value, &log)
// 			messageArray = append(messageArray, log)
// 		case TopicErrorLogs:
// 			log := model.ErrorLogs{}
// 			json.Unmarshal(message.Value, &log)
// 			messageArray = append(messageArray, log)
// 		}
// 	}
// 	switch topic {
// 	case TopicOperateLogs:
// 		var data []model.OperateLogs
// 		for _, item := range messageArray {
// 			structs, _ := item.(model.OperateLogs)
// 			data = append(data, structs)
// 		}
// 		err := usecase.LogUseCase.OperateLogBatchAdd(data)
// 		if err != nil {
// 			logger.Sugar.Errorw("OperateLogBatchAdd", util.GetSelfFuncName(), "OperateLogBatchAdd error:", err)
// 		}
// 	case TopicApiLogs:
// 		var data []model.ApiLogs
// 		for _, item := range messageArray {
// 			structs, _ := item.(model.ApiLogs)
// 			data = append(data, structs)
// 		}
// 		err := usecase.LogUseCase.ApiLogBatchAdd(data)
// 		if err != nil {
// 			logger.Sugar.Errorw("ApiLogBatchAdd", util.GetSelfFuncName(), "ApiLogBatchAdd error:", err)
// 		}
// 	case TopicErrorLogs:
// 		var data []model.ErrorLogs
// 		for _, item := range messageArray {
// 			structs, _ := item.(model.ErrorLogs)
// 			data = append(data, structs)
// 		}
// 		err := usecase.LogUseCase.ErrorLogBatchAdd(data)
// 		if err != nil {
// 			logger.Sugar.Errorw("ErrorLogBatchAdd", util.GetSelfFuncName(), "ErrorLogBatchAdd error:", err)
// 		}
// 	}
// 	if session.Context().Err() != nil {
// 		return nil
// 	}
// 	return nil
// }
