package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	pull, err := cli.ImagePull(ctx, "hello-world", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, pull)
	pull.Close()

	fmt.Println("--")

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "hello-world",
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	defer func() {
		fmt.Println("--")
		fmt.Println("cleaning...")
		wait := 1
		_ = cli.ContainerStop(ctx, resp.ID, container.StopOptions{
			Timeout: &wait,
		})
		_ = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveLinks:   false,
			RemoveVolumes: true,
		})
	}()

	fmt.Println("ID: ", resp.ID)

	fmt.Println("--")

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	go func() {
		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}
		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}()

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	fmt.Println("--")

	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case status := <-statusCh:
		fmt.Println("exit =", status.StatusCode)
	}

}
