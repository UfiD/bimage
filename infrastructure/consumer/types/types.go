package types

func SelectFileExtension(language string) string {
	switch language {
	case "go":
		return "task.go"
	case "python":
		return "task.py"
	case "cpp":
		return "task.cpp"
	default:
		return ""
	}
}
