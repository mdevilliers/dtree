package main

import (
	"strings"

	"github.com/mdevilliers/dtree"
	"github.com/spf13/cobra"
)

var grepCommand = &cobra.Command{
	Use:   "grep [TERM]",
	Short: "grep for a versions and dependancies for a package.",
	RunE: func(cmd *cobra.Command, args []string) error {

		store, err := initStore(_config)
		if err != nil {
			return err
		}

		fragment := args[0]
		nodes, _ := store.All()

		matches := []string{}

		for _, node := range nodes {
			if strings.Contains(node.Name, fragment) {
				matches = append(matches, node.Name)
			}
		}

		var nn []dtree.Node
		var ee []dtree.Edge

		for _, match := range matches {
			if !_config.Reverse {
				n, e := store.FromNode(match)
				nn = append(nn, n...)
				ee = append(ee, e...)

			} else {
				n, e := store.ToNode(match)
				nn = append(nn, n...)
				ee = append(ee, e...)
			}
		}
		err = outputImage(nn, ee)
		return err

	},
}
