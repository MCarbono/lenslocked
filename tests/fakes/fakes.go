package fakes

type IDGeneratorFake struct{}

func (i *IDGeneratorFake) Generate() string {
	return "fakeUUID"
}

func NewIDGeneratorFake() *IDGeneratorFake {
	return &IDGeneratorFake{}
}
