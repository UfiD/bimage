package types

const GolangDockerfileContent = `FROM golang:1.21-alpine

WORKDIR /app

COPY task.go .

RUN go build -o task task.go

CMD ["./task"]
`

const PythonDockerfileContent = `FROM python:alpine

WORKDIR /app

COPY task.py .

CMD ["python", "task.py"]
`

const CPPDockerfileContent = `FROM alpine:latest

RUN apk add --no-cache g++ musl-dev

WORKDIR /app

COPY task.cpp .

RUN g++ -o task task.cpp

CMD ["./task"]
`

func GetDockerfileContent(language string) string {
	switch language {
	case "go":
		return GolangDockerfileContent
	case "python":
		return PythonDockerfileContent
	case "cpp":
		return CPPDockerfileContent
	default:
		return ""
	}
}
