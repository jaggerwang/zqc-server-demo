package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/fsouza/go-dockerclient"
)

func startMongo(tries int, delay time.Duration) (addrs string, ctn *docker.Container, err error) {
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return addrs, ctn, err
	}

	ctn, err = client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: "daocloud.io/jaggerwang/mongodb",
		},
		HostConfig: &docker.HostConfig{
			PublishAllPorts: true,
		},
	})
	if err != nil {
		return addrs, ctn, err
	}
	err = client.StartContainer(ctn.ID, &docker.HostConfig{})
	if err != nil {
		return addrs, ctn, err
	}

	for i := 0; i < tries; i++ {
		ctn, err = client.InspectContainer(ctn.ID)
		if err != nil {
			return addrs, ctn, err
		}
		portBinding, ok := ctn.NetworkSettings.Ports["27017/tcp"]
		if !ok {
			time.Sleep(delay)
			continue
		}
		addrs = fmt.Sprintf("%v:%v", portBinding[0].HostIP, portBinding[0].HostPort)
		return addrs, ctn, nil
	}
	return addrs, ctn, errors.New("start mongodb failed")
}

func removeMongo(ctn *docker.Container) (err error) {
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return err
	}
	return client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            ctn.ID,
		RemoveVolumes: true,
		Force:         true,
	})
}
