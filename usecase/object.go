package usecase

type Object interface {
	Do(task, language string) string
}
