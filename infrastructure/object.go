package consumer

type Consumer interface {
	Do(task, language string) string
}
