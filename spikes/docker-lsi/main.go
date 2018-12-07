package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, img := range images {
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
