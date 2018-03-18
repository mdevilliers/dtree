package main

import (
	"log"
	"os"

	"github.com/mdevilliers/dtree"
	"github.com/mdevilliers/dtree/parsers"
	"github.com/mdevilliers/dtree/repo"
	"github.com/mdevilliers/dtree/store"
	"github.com/spf13/cobra"
)

var generateDiagramCommand = &cobra.Command{
	Use:   "gen",
	Short: "Generate a dependancy diagram.",
	RunE: func(cmd *cobra.Command, args []string) error {

		all, err := repo.FromCheckedOut(_config.Root).Paths()

		if err != nil {
			return err
		}

		allNodes := map[string]dtree.Node{}
		allEdges := []dtree.Edge{}

		for _, r := range all {

			nodes, edges, err := parsers.Glide().Parse(r.Path)

			if err != nil {
				log.Println("error parsing : ", r, err)
				continue
			}

			allEdges = append(allEdges, edges...)

			for _, n := range nodes {
				allNodes[n.Name] = n
			}

		}

		nodesArr := []dtree.Node{}

		for _, n := range allNodes {
			nodesArr = append(nodesArr, n)
		}

		s := store.InMemory(nodesArr, allEdges)

		var nn []dtree.Node
		var ee []dtree.Edge

		if !_config.Reverse {
			nn, ee = s.FromNode(_config.Focus)
		} else {
			nn, ee = s.ToNode(_config.Focus)
		}

		postProcessNodes(nn)

		return writeDot(nn, ee, os.Stdout)

	},
}
