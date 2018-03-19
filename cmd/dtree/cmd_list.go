package main

import (
	"fmt"
	"strings"

	"github.com/mdevilliers/dtree"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "ls",
	Short: "List dependancies.",
	RunE: func(cmd *cobra.Command, args []string) error {

		store, err := InitStore(_config)
		if err != nil {
			return err
		}

		grouped := store.GroupAll()

		for k, g := range grouped {

			fmt.Println(k)

			x := normaliseEdgeArrByVersion(g)

			for version, edge := range x {
				fmt.Println("\t", version, " - ", strings.Join(edge, ","))
			}
		}

		return nil
	},
}

func normaliseEdgeArrByVersion(edges []dtree.Edge) map[string][]string {

	normalised := map[string][]string{}

	for _, e := range edges {

		key := e.Version
		_, contains := normalised[key]

		if !contains {
			normalised[key] = []string{e.Source.Name}
		} else {
			normalised[key] = append(normalised[key], e.Source.Name)
		}
	}

	return normalised
}
