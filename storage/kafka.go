package storage

// NewKafka to create new kafka client
func NewKafka() *Kafka {
	return &Kafka{}
}

// Kafka to provide persistence for Saga
type Kafka struct{}
