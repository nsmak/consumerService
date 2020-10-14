package store

import "github.com/nsmak/consumerService/consumer/models"

// DataStore - interface of data storage service
type DataStore interface {
	CreateConsumer(c models.Consumer) (models.Consumer, error)
	ConsumerIsExist(email string) (bool, error)
	Consumer(email string) (models.Consumer, error)
}
