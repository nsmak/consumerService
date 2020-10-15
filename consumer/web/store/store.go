package store

import "github.com/nsmak/consumerService/consumer/models"

//go:generate mockgen -destination=../domain/mock_db_test.go -package=domain_test -source=store.go DataStore

// DataStore - interface of data storage service
type DataStore interface {
	CreateConsumer(c models.Consumer) (models.Consumer, error)
	ConsumerIsExist(email string) (bool, error)
	Consumer(email string) (models.Consumer, error)
}
