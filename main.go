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

func getContainers(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		img_name := fmt.Sprintf("%s %s\n", container.ID[:10], container.Image)
		io.WriteString(w, img_name)
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/state-machine", getStateMachineDiagram)
	http.HandleFunc("/containers", getContainers)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
