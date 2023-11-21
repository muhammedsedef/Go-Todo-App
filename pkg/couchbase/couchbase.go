package couchbase

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"time"
)

func ConnectCluster(host, username, password string) *gocb.Cluster {
	cluster, err := gocb.Connect(
		host,
		gocb.ClusterOptions{
			Username: username,
			Password: password,
			TimeoutsConfig: gocb.TimeoutsConfig{
				KVTimeout: 10 * time.Second,
			},
		})

	if err != nil {
		panic(fmt.Sprintf("%s, error : %s", "cannot connect couchbase cluster", err.Error()))
	}
	fmt.Println("couchbase connection success")
	return cluster
}
