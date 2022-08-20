package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"go-minecraft-server/mcss"
)

const PREFIX = "mc-panel-go__"

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

	http.HandleFunc("/api/v1/servers", listServers)     // get only
	http.HandleFunc("/api/v1/server/new", createServer) // post only

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listServers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
	}

	cli, _ := client.NewClientWithOpts(client.FromEnv)
	containers, _ := cli.ContainerList(context.Background(), types.ContainerListOptions{})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(containers)
}

func createServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
	}

	createdContainer := runContainerInBackground()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdContainer)
}

// example code, pulling and running a container, need to update to use the mc container
func runContainerInBackground() container.ContainerCreateCreatedBody {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	imageName := "bfirsh/reticulate-splines"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{Image: imageName}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

	return resp
}
