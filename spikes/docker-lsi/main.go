package main

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		panic(err)
	}
	for _, img := range imgs {
		fmt.Println("Labels: ", img.Labels)
		fmt.Println("  ID: ", img.ID)
		fmt.Println("  Created: ", img.Created)
		fmt.Println("  ParentId: ", img.ParentID)
		fmt.Println("  RepoDigests: ", img.RepoDigests)
		fmt.Println("  RepoTags: ", img.RepoTags)
		fmt.Println("  Size: ", img.Size)
		fmt.Println("  VirtualSize: ", img.VirtualSize)
	}
}
