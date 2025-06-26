package types

const GolangDockerfileContent = `FROM golang:1.21-alpine

WORKDIR /app

COPY task.go .

RUN go build -o task task.go

CMD ["./task"]
`

func GetDockerfileContent(language string) string {
	switch language {
	case "go":
		return GolangDockerfileContent
	default:
		return ""
	}
}
