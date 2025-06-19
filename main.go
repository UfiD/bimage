package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/moby/go-archive"
)

func main() {
	// Создаем клиент Docker
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Создаем временную директорию с Dockerfile
	tempDir, err := os.MkdirTemp("", "docker-build")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем Dockerfile
	dockerfileContent := `FROM alpine:latest
RUN apk add --no-cache curl
CMD ["sh", "-c", "echo 'Контейнер успешно запущен!' && sleep 3600"]
`

	dockerfilePath := tempDir + "/Dockerfile"
	err = os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
	if err != nil {
		panic(err)
	}

	// Создаем tar-архив для сборки образа
	buildContext, err := archive.TarWithOptions(tempDir, &archive.TarOptions{})
	if err != nil {
		panic(err)
	}

	// Собираем образ
	imageName := "my-dynamic-image"
	buildResponse, err := cli.ImageBuild(ctx, buildContext, build.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Remove:     true,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer buildResponse.Body.Close()

	fmt.Println("Сборка образа:")
	_, err = io.Copy(os.Stdout, buildResponse.Body)
	if err != nil {
		panic(err)
	}

	// Создаём контейнер
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	// Запускаем контейнер
	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Контейнер успешно запущен с ID: %s\n", resp.ID)
}
