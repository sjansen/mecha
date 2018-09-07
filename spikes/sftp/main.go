package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/pkg/sftp"
)

var (
	HOST = flag.String("host", "localhost", "ssh server hostname")
	PORT = flag.Int("port", 22, "ssh server port")
	USER = flag.String("user", os.Getenv("USER"), "ssh username")
	PASS = flag.String("pass", "", "ssh password")
)

func init() {
	flag.Parse()
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func connect() (conn *ssh.Client, err error) {
	var auths []ssh.AuthMethod

	// ssh-agent
	if sockname := os.Getenv("SSH_AUTH_SOCK"); sockname != "" {
		if conn, err := net.Dial("unix", sockname); err == nil {
			auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(conn).Signers))

		}
	}

	// password
	if *PASS != "" {
		auths = append(auths, ssh.Password(*PASS))
	}

	config := ssh.ClientConfig{
		User:            *USER,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", *HOST, *PORT)
	conn, err = ssh.Dial("tcp", addr, &config)
	return
}

type byName []os.FileInfo

func (a byName) Len() int      { return len(a) }
func (a byName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool {
	return a[i].Name() < a[j].Name()
}

func main() {
	conn, err := connect()
	if err != nil {
		die(err)
	}
	defer conn.Close()

	c, err := sftp.NewClient(conn)
	if err != nil {
		die(err)
	}
	defer c.Close()

	wd, err := c.Getwd()
	if err != nil {
		die(err)
	}

	files, err := c.ReadDir(wd)
	if err != nil {
		die(err)
	}

	sort.Sort(byName(files))

	fmt.Println(wd)
	for i, f := range files {
		if f.IsDir() {
			fmt.Printf("  % 24s/", f.Name())
		} else {
			fmt.Printf("  % 25s", f.Name())
		}
		if (i % 3) == 2 {
			fmt.Println()
		}
	}
	fmt.Println()
}
