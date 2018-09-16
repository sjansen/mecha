package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sjansen/watchman"
)

func connect() *watchman.Client {
	os.Stdout.Write([]byte("Connecting to Watchman... "))
	os.Stdout.Sync()
	c, err := watchman.Connect()
	if err != nil {
		fmt.Println("FAILURE")
		die(err)
	}
	fmt.Println("SUCCESS")

	return c
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func loop(c *watchman.Client) {
	for n := range c.Notifications() {
		cn, ok := n.(*watchman.ChangeNotification)
		if !ok || cn.IsFreshInstance {
			continue
		}
		fmt.Printf(
			"Update: (clock=%q)\n",
			cn.Clock,
		)
		files := cn.Files
		for _, file := range files {
			switch file.Type {
			case "d":
				fmt.Printf("  %9s  %s/\n",
					file.Change, file.Name,
				)
			case "l":
				fmt.Printf("  %9s  %s -> %s\n",
					file.Change, file.Name, file.Target,
				)
			default:
				fmt.Printf("  %9s  %s\n",
					file.Change, file.Name,
				)
			}
		}
		fmt.Println()
	}
}

func mkdir() (dir string, err error) {
	dir, err = ioutil.TempDir("", "watchman-client-test")
	if err != nil {
		return
	}

	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return
	}

	path := filepath.Join(dir, ".watchmanconfig")
	err = ioutil.WriteFile(path, []byte(`{"idle_reap_age_seconds": 300}`+"\n"), os.ModePerm)
	return
}

func main() {
	c := connect()
	fmt.Printf("version: %s\n\n", c.Version())

	dir, err := mkdir()
	if err != nil {
		die(err)
	}
	defer os.RemoveAll(dir)

	fmt.Printf("Watching: %s\n\n", dir)
	watch, err := c.AddWatch(dir)
	if err != nil {
		die(err)
	}

	go func() {
		loop(c)
	}()

	_, err = watch.Subscribe("example", dir)
	if err != nil {
		die(err)
	}

	for i, basename := range []string{"foo", "bar", "baz"} {
		filename := filepath.Join(dir, basename)
		ioutil.WriteFile(filename, []byte{}, os.ModePerm)
		time.Sleep(time.Duration(i) * time.Second)
		os.Remove(filename)
		time.Sleep(time.Second)
	}
}
