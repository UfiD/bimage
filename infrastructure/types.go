package consumer

type Consumer interface {
	Do(task string) string
}
