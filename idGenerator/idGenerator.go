package idGenerator

import "github.com/google/uuid"

type IDGenerator interface {
	Generate() string
}

type IDGeneratorImpl struct{}

func (i *IDGeneratorImpl) Generate() string {
	return uuid.New().String()
}

func New() *IDGeneratorImpl {
	return &IDGeneratorImpl{}
}
