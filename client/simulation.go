package main

import (
	"fmt"
	nomad "github.com/hashicorp/nomad/api"
)

type NomadSimulation struct {
	client *nomad.Client
}

func NewNomadSimulation() *NomadSimulation {
	// First, let's establish the connection to the nomad server.
	var conf = nomad.DefaultConfig()
	var client, err = nomad.NewClient(conf)
	ExitOnError(err)

	// Ping the server and get back a response.
	resp, err := client.Status().Leader()
	ExitOnError(err)
	fmt.Println(resp)

	// Now, ping it for a list of nodes
	nodes, _, err := client.Nodes().List(&nomad.QueryOptions{})
	ExitOnError(err)
	for _, node := range nodes {
		fmt.Printf(nodeTmpl, node.Address, node.ID, node.Name, node.Status)
		fmt.Println()
	}

	return &NomadSimulation{client}
}
