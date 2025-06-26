package codeprocessor

import (
	"bimage/infrastructure/consumer/types"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/go-archive"
)

type Consumer struct{}

func New() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Do(task, language string) string {

	exeDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("не удалось получить текущую директорию: %v", err))
	}

	tempDir := filepath.Join(exeDir, "docker-build-temp")
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		panic(fmt.Errorf("не удалось создать временную директорию: %v", err))
	}
	defer os.RemoveAll(tempDir)

	fileName := types.SelectFileExtension(language)

	goFilePath := filepath.Join(tempDir, fileName)
	err = os.WriteFile(goFilePath, []byte(task), 0644)
	if err != nil {
		panic(fmt.Errorf("не удалось создать Go файл: %v", err))
	}

	dockerfileContent := types.GetDockerfileContent(language)

	dockerfilePath := filepath.Join(tempDir, "Dockerfile")
	err = os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
	if err != nil {
		panic(fmt.Errorf("не удалось создать Dockerfile: %v", err))
	}

	fmt.Printf("Временная директория: %s\n", tempDir)
	fmt.Printf("Go файл создан: %s\n", goFilePath)
	fmt.Printf("Dockerfile создан: %s\n", dockerfilePath)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(fmt.Errorf("не удалось создать Docker клиент: %v", err))
	}

	buildContext, err := archive.TarWithOptions(tempDir, &archive.TarOptions{})
	if err != nil {
		panic(fmt.Errorf("не удалось создать архив для сборки: %v", err))
	}

	imageName := "go-task"
	buildResponse, err := cli.ImageBuild(ctx, buildContext, build.ImageBuildOptions{
		Tags:       []string{imageName},
		Remove:     true,
		Dockerfile: "Dockerfile",
	})
	if err != nil {
		panic(fmt.Errorf("ошибка сборки образа: %v", err))
	}
	defer buildResponse.Body.Close()

	fmt.Println("\nПроцесс сборки образа:")
	_, err = io.Copy(os.Stdout, buildResponse.Body)
	if err != nil {
		panic(fmt.Errorf("ошибка чтения вывода сборки: %v", err))
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		panic(fmt.Errorf("ошибка создания контейнера: %v", err))
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		panic(fmt.Errorf("ошибка запуска контейнера: %v", err))
	}

	fmt.Printf("\nКонтейнер запущен с ID: %s\n", resp.ID)

	var containerOutput bytes.Buffer

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(fmt.Errorf("ошибка ожидания контейнера: %v", err))
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     false,
		Tail:       "all",
	})
	if err != nil {
		panic(fmt.Errorf("ошибка получения логов: %v", err))
	}
	defer out.Close()

	if _, err := stdcopy.StdCopy(&containerOutput, &containerOutput, out); err != nil {
		panic(fmt.Errorf("ошибка копирования логов: %v", err))
	}

	fmt.Println("=== Вывод контейнера ===")
	fmt.Println(containerOutput.String())
	res := containerOutput.String()

	if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{}); err != nil {
		fmt.Printf("Предупреждение: не удалось удалить контейнер: %v\n", err)
	}
	return res
}
