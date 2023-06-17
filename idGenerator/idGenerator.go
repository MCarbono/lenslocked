package idGenerator

import "github.com/google/uuid"

type IDGenerator interface {
	New() string
}

type IDGeneratorImpl struct{}

func (i *IDGeneratorImpl) New() string {
	return uuid.New().String()
}
