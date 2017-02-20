package database

import (
	"context"
	"errors"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
)

func initContainer() (err error) {
	// Create cli instance
	cli, err := client.NewEnvClient()
	if err != nil {
		glog.Fatal(err)
		return err
	}

	// Get running containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		glog.Fatal(err)
		return err
	}

	// Check if m2svis container exists
	for _, container := range containers {
		if strings.Contains(container.Names[0], dbConfig.Container) {
			return nil
		}
	}

	// TODO: create one
	// Ask users to create container if not found
	warning := "Cannot find m2svis container. Create with following command:\n\t"
	warning += "docker run --name m2svis -p 3306:3306 -e MYSQL_DATABASE=m2svis "
	warning += "-e MYSQL_USER=m2svis -e MYSQL_PASSWORD=m2svis -e MYSQL_ROOT_PASSWORD=mysqlroot -d mysql"
	glog.Warning(warning)
	return errors.New("no m2svis container found")
}
