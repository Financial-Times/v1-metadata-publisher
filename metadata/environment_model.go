package metadata

import (
	"strings"
)

type Cluster struct {
	address  string
	username string
	password string
}

func GetCluster(address string, credentials string) *Cluster {
	cluster := Cluster{address: address}
	if credentials != "" {
		auth := strings.Split(credentials, ":")
		cluster.username = auth[0]
		cluster.password = auth[1]
	}
	return &cluster
}

func (c *Cluster) GetAddress() string {
	return c.address
}

func (c *Cluster) GetUsername() string {
	return c.username
}

func (c *Cluster) GetPassword() string {
	return c.password
}
