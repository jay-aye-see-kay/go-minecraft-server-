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
	fmt.Printf("got / request\n")
	io.WriteString(w, "<div>Pages:</div>")
	io.WriteString(w, "<div><ul>")
	io.WriteString(w, "<li><a href='/containers'>list of containers</a></li>")
	io.WriteString(w, "<li><a href='/state-machine'>a graphviz representation of the state machine</a></li>")
	io.WriteString(w, "</ul></div>")
}

func getStateMachineDiagram(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/state-machine.html"))

	tmpl.Execute(w, mcss.MakeStateMachine().ToGraph())
}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/state-machine", getStateMachineDiagram)
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
