package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"go-minecraft-server/mcss"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func getStateMachineDiagram(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/state-machine.html"))
	tmpl.Execute(w, mcss.MakeStateMachine().ToGraph())
}

// example function playing with the docker api
func getDockerInfo(w http.ResponseWriter, r *http.Request) {
	cli, _ := client.NewClientWithOpts(client.FromEnv)

	// containers
	io.WriteString(w, "\nContainters:\n")
	containers, _ := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	for _, container := range containers {
		img_name := fmt.Sprintf("%s %s\n", container.ID[:10], container.Image)
		io.WriteString(w, img_name)
	}

	// networks
	io.WriteString(w, "\nNetworks:\n")
	networks, _ := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	for _, network := range networks {
		io.WriteString(w, fmt.Sprintf("%s %s\n", network.ID, network.Name))
	}

	// images
	io.WriteString(w, "\nImages:\n")
	images, _ := cli.ImageList(context.Background(), types.ImageListOptions{})
	for _, image := range images {
		io.WriteString(w, fmt.Sprintf("%s %d MB\n", image.ID, image.Size/1_000_000))
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/state-machine", getStateMachineDiagram)
	http.HandleFunc("/docker-info", getDockerInfo)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
