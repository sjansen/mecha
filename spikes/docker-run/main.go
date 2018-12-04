package main

import (
	"fmt"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	err = client.PullImage(
		docker.PullImageOptions{
			Repository:   "hello-world",
			Tag:          "latest",
			OutputStream: os.Stdout,
		},
		docker.AuthConfiguration{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("--")

	c, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: "",
		Config: &docker.Config{
			Image:        "hello-world",
			AttachStdin:  false,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          false,
		},
		HostConfig:       nil,
		NetworkingConfig: nil,
		Context:          nil,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("ID: ", c.ID)

	attached := make(chan struct{})
	go func() {
		_ = client.AttachToContainer(docker.AttachToContainerOptions{
			Container:    c.ID,
			OutputStream: os.Stdout,
			ErrorStream:  os.Stderr,
			Stdin:        false,
			Stdout:       true,
			Stderr:       true,
			Stream:       true,
			Success:      attached,
		})
	}()
	<-attached
	attached <- struct{}{}

	fmt.Println("--")

	err = client.StartContainer(c.ID, nil)
	if err != nil {
		panic(err)
	}

	status, err := client.WaitContainer(c.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println("--")

	fmt.Println("exit =", status)
}
