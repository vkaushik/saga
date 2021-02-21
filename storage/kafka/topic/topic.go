package topic

import (
	"github.com/Shopify/sarama"
	"github.com/juju/errors"
)

// Topic provides kafka-topic management functions.
type Topic struct {
	admin sarama.ClusterAdmin
}

// NewWithClusterAdmin to create new Topic with already created ClusterAdmin
func NewWithClusterAdmin(admin sarama.ClusterAdmin) *Topic {
	return &Topic{admin: admin}
}

// New to create a new Topic object
func New(v string, brk []string) (*Topic, error) {
	var ver sarama.KafkaVersion
	var err error
	if ver, err = sarama.ParseKafkaVersion(string(v)); err != nil {
		return nil, errors.Annotate(err, "invalid kafka version valid e.g. 2.7.0")
	}

	var ca sarama.ClusterAdmin
	sc := sarama.NewConfig()
	sc.Version = ver

	if ca, err = sarama.NewClusterAdmin(brk, sarama.NewConfig()); err != nil {
		return nil, errors.Annotate(err, "could not create new cluster admin")
	}

	return NewWithClusterAdmin(ca), nil
}

// IsTopicAlreadyCreated checks if topic is already created in Kafka
func (t *Topic) IsTopicAlreadyCreated(topic string) (exists bool, err error) {
	existingTopics, err := t.GetAllTopics()
	if err != nil {
		return false, errors.Annotate(err, "could not fetch existing topics")
	}

	for _, tp := range existingTopics {
		if tp == topic {
			return true, nil
		}
	}

	return false, nil
}

// GetAllTopics fetches all topic names that are created in Kafka
func (t *Topic) GetAllTopics() ([]string, error) {
	topicNameToDetailsMap, err := t.admin.ListTopics()
	size := len(topicNameToDetailsMap)
	topics := make([]string, 0, size)
	for topicName := range topicNameToDetailsMap {
		topics = append(topics, topicName)
	}

	return topics, err
}

// CreateTopic creates a new Topic in Kafka
func (t *Topic) CreateTopic(topic string, pc int32, rc int16) (err error) {
	td := &sarama.TopicDetail{
		NumPartitions:     pc,
		ReplicationFactor: rc,
		ConfigEntries:     make(map[string]*string),
	}

	return t.admin.CreateTopic(topic, td, false)
}

// DeleteTopic deletes a Topic in Kafka
func (t *Topic) DeleteTopic(topic string) (err error) {
	return t.admin.DeleteTopic(topic)
}
