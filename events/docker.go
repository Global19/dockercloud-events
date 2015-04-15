package events

import (
	dc "github.com/fsouza/go-dockerclient"
	"log"
	"os/exec"
	"strings"
)

type DockerClient struct{ client *dc.Client }

var (
	DockerHost   string
	DockerBinary string
)

func NewDockerClient(host string) (DockerClient, error) {
	c, err := dc.NewClient(host)
	if err != nil {
		return DockerClient{}, err
	}
	return DockerClient{c}, nil
}

func (self DockerClient) addEventListener() (listener chan *dc.APIEvents, err error) {
	listener = make(chan *dc.APIEvents, 10)
	return listener, self.client.AddEventListener(listener)
}

func (self DockerClient) removeEventListener(listener chan *dc.APIEvents) error {
	return self.client.RemoveEventListener(listener)
}

func (self DockerClient) inspect(id string) string {
	out_bytes, err := exec.Command(DockerBinary, "inspect", id).CombinedOutput()
	out := string(out_bytes)
	if err != nil {
		if !strings.Contains(out, "No such image or container") {
			extra := map[string]interface{}{"docker_inspcet_output": out}
			SendError(err, "docker inspect failed", extra)
		}
		log.Println("Docker inspect error:", err, out)
		return ""
	}
	return string(out)
}

func (self DockerClient) ps(opts *dc.ListContainersOptions) ([]dc.APIContainers, error) {
	return self.client.ListContainers(*opts)
}
