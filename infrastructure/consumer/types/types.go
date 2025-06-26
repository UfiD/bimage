package types

func SelectFileExtension(language string) string {
	switch language {
	case "go":
		return "task.go"
	default:
		return ""
	}
}
