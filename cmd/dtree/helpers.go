package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/mdevilliers/dtree"
	"github.com/mdevilliers/dtree/output"
	"github.com/mdevilliers/dtree/parsers"
	"github.com/mdevilliers/dtree/repo"
	"github.com/mdevilliers/dtree/store"
)

type storer interface {
	FindNodes(str string) ([]dtree.Node, error)
	FromNode(focus dtree.Node) ([]dtree.Node, []dtree.Edge)
	ToNode(focus dtree.Node) ([]dtree.Node, []dtree.Edge)
}

func initStore(cfg *config) (storer, error) {

	all, err := repo.FromCheckedOut(cfg.Root)

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

func outputDot(_ string, n []dtree.Node, e []dtree.Edge) error {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err := output.ToDot(n, e, w)

	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	fmt.Println(b.String())

	return nil
}

func outputSvg(fragment string, n []dtree.Node, e []dtree.Edge) error {

	dotfile, err := ioutil.TempFile("", "dtree_")

	if err != nil {
		return err
	}

	defer dotfile.Close() // nolint  errcheck

	err = output.ToDot(n, e, dotfile)

	if err != nil {
		return err
	}

	data, err := executeCommand(fmt.Sprintf("dot -Tsvg %s", dotfile.Name()))

	if err != nil {
		return err
	}

	now := time.Now()

	fileName := fmt.Sprintf("output_%s_%v.svg", fragment, now.Unix())
	fileName = strings.Replace(fileName, "/", "_", -1)

	err = ioutil.WriteFile(fileName, data, 0644)

	if err != nil {
		return err
	}

	_, err = executeCommand(fmt.Sprintf("xdg-open %s", fileName))

	return err

}

func executeCommand(cmd string) ([]byte, error) {

	out, err := exec.Command("sh", "-c", cmd).Output() // nolint: gosec

	if err != nil {
		return []byte{}, err
	}

	return out, nil
}
