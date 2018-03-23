package main

import (
	"log"

	"github.com/mdevilliers/dtree"
	"github.com/mdevilliers/dtree/parsers"
	"github.com/mdevilliers/dtree/repo"
	"github.com/mdevilliers/dtree/store"
)

type storer interface {
	All() ([]dtree.Node, []dtree.Edge)
	FromNode(focus string) ([]dtree.Node, []dtree.Edge)
	ToNode(focus string) ([]dtree.Node, []dtree.Edge)
	GroupAll() map[string][]dtree.Edge
	GroupOn(focus string) map[string][]dtree.Edge
}

func initStore(cfg *config) (storer, error) {

	all, err := repo.FromCheckedOut(cfg.Root).Paths()

	if err != nil {
		return nil, err
	}

	allNodes := map[string]dtree.Node{}
	allEdges := []dtree.Edge{}

	for _, r := range all {

		parser, err := parsers.New(r.Path)
		if err != nil {
			log.Println("error :", r, err)
			continue
		}

		nodes, edges, err := parser.Parse(r.Path)

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

	return store.InMemory(nodesArr, allEdges), nil
}
