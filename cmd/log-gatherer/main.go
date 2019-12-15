package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"

	"github.com/gekoil/log-gatherer/internal"
	"github.com/gekoil/log-gatherer/pkg/storage"

	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	loggers    internal.Containers = make(map[string]chan int)
	logStorage storage.LogStorage
)

func init() {
	logStorage = &storage.FileLogStorage{}
}

func main() {
	log.Print("Running server")
	http.HandleFunc("/attach", attachLog)
	http.HandleFunc("/detach/", detachLog)
	http.HandleFunc("/logs/", getLogs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func attachLog(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print(err)
	}
	options := new(internal.RequestOptions)
	err = json.Unmarshal(body, options)
	log.Printf("Recieved options: %#v", options)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	filter := filters.NewArgs()

	if options.Filter.ContainerId != "" {
		filter.Add("id", options.Filter.ContainerId)
	}
	if options.Filter.ContainerName != "" {
		filter.Add("name", options.Filter.ContainerName)
	}
	if options.Filter.ContainerLabel != "" {
		filter.Add("label", options.Filter.ContainerLabel)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filter,
	})

	for _, container := range containers {
		log.Printf("Starting logging for %#v", container.ID)

		stopChannel := make(chan int)
		loggers[container.ID] = stopChannel
		logStorage.Create(container.ID)

		go func(containerID string, stopChan chan int) {
			attachOptions := types.ContainerAttachOptions{
				Stdin:  false,
				Stdout: options.StdOut,
				Stderr: options.StdErr,
				Logs:   true,
				Stream: true,
			}
			out, err := cli.ContainerAttach(ctx, containerID, attachOptions)
			defer out.Close()
			if err != nil {
				log.Print(err)
				return
			}

			streamChan := make(chan int)
			go func() {
				stdOut := new(strings.Builder)
				stdErr := new(strings.Builder)
				_, err := stdcopy.StdCopy(stdOut, stdErr, out.Reader)
				if err != nil {
					log.Print(err)
				} else {
					log.Printf("Collected %d in output and %d in error", stdOut.Len(), stdErr.Len())
					logStorage.Insert(containerID, stdOut, stdErr)
				}
				close(streamChan)
			}()

			resultChan, errChan := cli.ContainerWait(ctx, container.ID, "")
			for {
				select {
				case <-streamChan:
				case <-stopChan:
				case <-resultChan:
					return
				case err := <-errChan:
					log.Print(err)
					return
				}
			}
		}(container.ID, stopChannel)
	}
	writer.WriteHeader(http.StatusCreated)
}

func detachLog(writer http.ResponseWriter, request *http.Request) {
	id := strings.TrimPrefix(request.URL.Path, "/detach/")
	log.Printf("Recieved call to detach %s", id)
	stopChan, ok := loggers[id]
	if ok {
		close(stopChan)
		delete(loggers, id)
		writer.WriteHeader(http.StatusAccepted)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func getLogs(writer http.ResponseWriter, request *http.Request) {
	id := strings.TrimPrefix(request.URL.Path, "/logs/")
	stdLogs, errLogs := logStorage.Get(id)

	_, err := io.Copy(writer, stdLogs)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	_, err = io.Copy(writer, errLogs)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
