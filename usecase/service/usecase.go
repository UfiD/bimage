package service

import (
	consumer "bimage/infrastructure"
)

type Usecase struct {
	Consumer consumer.Consumer
}

func New(c consumer.Consumer) *Usecase {
	return &Usecase{
		Consumer: c,
	}
}

func (uc *Usecase) Do(task, language string) string {
	return uc.Consumer.Do(task, language)
}
