package kafka

import (
	"github.com/vkaushik/saga/trace"
	"sync"
	"time"

	"github.com/vkaushik/saga/storage/kafka/topic"

	"github.com/juju/errors"

	"github.com/Shopify/sarama"
)

// BrokerAddresses is the list of broker addresses e.g. ["192.168.1.19:9092", "192.168.1.21:9092"]
type BrokerAddresses []string

// PartitionCount is the number of partitions for any topic that'll be created by saga
type PartitionCount int32

// ReplicaCount is the number of replicas for any topic that'll be created by saga
type ReplicaCount int16

// ConsumerWaitDuration is the timeout duration for which consumer will wait for the message
type ConsumerWaitDuration time.Duration

// Version is the Kafka version e.g. "2.7.0"
type Version string

// MaxLogMessages is the maximum number of logs messages that will be published to Kafka in a single Transaction.
type MaxLogMessages int

// Kafka object provides storage functions to provide persistence medium for Saga. Saga is as reliable as it's persistence medium.
type Kafka struct {
	ver     sarama.KafkaVersion
	brokers BrokerAddresses

	topic    Topic
	producer sarama.SyncProducer
	consumer sarama.Consumer
	logger   trace.Logger

	pc      PartitionCount
	rc      ReplicaCount
	dur     ConsumerWaitDuration
	maxMsgs MaxLogMessages
}

func (k *Kafka) TxIDAlreadyExists(s string) (bool, error) {
	panic("implement me")
}

// Topic provides kafka-topic management functions.
type Topic interface {
	IsTopicAlreadyCreated(topicName string) (exists bool, err error)
	GetAllTopics() (topicNames []string, err error)
	CreateTopic(topicName string, partitionCount int32, replicaCount int16) (err error)
	DeleteTopic(topicName string) (err error)
}

// New to create a new Kafka object. It accepts functional options to support maximum customization.
func New(ver Version, brk BrokerAddresses, options ...func(*Kafka) error) (*Kafka, error) {
	k := &Kafka{brokers: brk}
	k.SetDefaults()
	var err error

	if k.topic, err = topic.New(string(ver), brk); err != nil {
		return nil, errors.Annotate(err, "could not create Topic object")
	}

	if k.producer, err = sarama.NewSyncProducer(k.brokers, nil); err != nil {
		return nil, errors.Annotate(err, "could not create new sync producer")
	}

	if k.consumer, err = sarama.NewConsumer(k.brokers, nil); err != nil {
		return nil, errors.Annotate(err, "could not create new consumer")
	}

	// Warning: Please be careful with options
	for _, setter := range options {
		if err = setter(k); err != nil {
			return nil, errors.Annotate(err, "issue while setting options")
		}
	}

	return k, nil
}

// SetDefaults sets sensible defaults to the Kafka config
func (k *Kafka) SetDefaults() {
	k.pc = 1
	k.rc = 1
	k.dur = ConsumerWaitDuration(5000)
}

// SetNumberOfPartitions is the functional option to set PartitionCount
func SetNumberOfPartitions(pc PartitionCount) func(*Kafka) error {
	return func(k *Kafka) error {
		if pc <= 1 {
			return errors.New("number of partitions must be greater than 0")
		}
		k.pc = pc
		return nil
	}
}

// SetNumberOfReplicas is the functional option to set ReplicaCount
func SetNumberOfReplicas(rc ReplicaCount) func(*Kafka) error {
	return func(k *Kafka) error {
		if rc <= 1 {
			return errors.New("number of replicas must be greater than 0")
		}
		k.rc = rc
		return nil
	}
}

// SetConsumeReturnDuration is the functional option to set ConsumerWaitDuration
func SetConsumeReturnDuration(dur ConsumerWaitDuration) func(*Kafka) error {
	return func(k *Kafka) error {
		k.dur = dur
		return nil
	}
}

// SetNumberMaximumLogMessages is the functional option to set MaxLogMessages
func SetNumberMaximumLogMessages(mc MaxLogMessages) func(*Kafka) error {
	return func(k *Kafka) error {
		if mc <= 1 {
			return errors.New("number of replicas must be greater than 0")
		}
		k.maxMsgs = mc
		return nil
	}
}

// AppendLog
func (k *Kafka) AppendLog(txID string, data string) error {
	topicExists, err := k.topic.IsTopicAlreadyCreated(txID)
	if err != nil {
		return errors.Annotatef(err, "could not check if topic: %v, is already created", txID)
	}
	if !topicExists {
		err = k.topic.CreateTopic(txID, int32(k.pc), int16(k.rc))
		if err != nil {
			return errors.Annotatef(err, "could not create new topic: %v", txID)
		}
	}

	msg := &sarama.ProducerMessage{Topic: txID, Value: sarama.StringEncoder(data)}
	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return errors.Annotatef(err, "could not publish kafka message with data: %v, to topic: %v", data, txID)
	}

	k.logger.Info("message data: ,published to partition %d at offset %d\n", data, partition, offset)

	return nil
}

// Lookup
func (k *Kafka) Lookup(txID string) ([]string, error) {
	partitionList, err := k.consumer.Partitions(txID)
	if err != nil {
		return nil, errors.Annotatef(err, "could not get partitions for topic: %v", txID)
	}

	var wg sync.WaitGroup
	// TODO: refactor channel sizes, wg usage
	msgs := make(chan *sarama.ConsumerMessage, int(k.maxMsgs))
	errs := make(chan error, int(k.maxMsgs))
	for _, partition := range partitionList {
		wg.Add(1)
		go consumePartition(k.consumer, txID, partition, msgs, errs, time.Duration(k.dur), &wg)
	}

	wg.Wait()

	data := make([]string, 0, int(k.maxMsgs))
	timer := time.NewTimer(time.Duration(k.dur))
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			return data, nil
		case msg := <-msgs:
			data = append(data, string(msg.Value))
			timer.Reset(time.Duration(k.dur))
		case err := <-errs:
			return data, errors.Annotate(err, "one of the partition-consumers failed")
		}
	}
}

// TODO: refactor parameter types and signature
func consumePartition(consumer sarama.Consumer, topic string, partition int32,
	msgs chan *sarama.ConsumerMessage, errs chan error, waitTime time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		errs <- err
		return
	}
	defer func() {
		_ = pc.Close()
	}()
	timer := time.NewTimer(waitTime)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			return
		case msg := <-pc.Messages():
			msgs <- msg
			timer.Reset(waitTime)
		}
	}
}
