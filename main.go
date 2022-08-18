package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "<div>Pages:</div>")
	io.WriteString(w, "<div><ul>")
	io.WriteString(w, "<li><a href='/containers'>list of containers</a></li>")
	io.WriteString(w, "</ul></div>")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}
		for _, container := range containers {
			img_name := fmt.Sprintf("%s %s\n", container.ID[:10], container.Image)
			io.WriteString(w, img_name)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
