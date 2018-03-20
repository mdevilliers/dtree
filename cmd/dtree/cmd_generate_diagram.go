package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/mdevilliers/dtree"
	"github.com/spf13/cobra"
)

var generateDiagramCommand = &cobra.Command{
	Use:   "gen",
	Short: "Generate a dependancy diagram.",
	RunE: func(cmd *cobra.Command, args []string) error {

		store, err := initStore(_config)
		if err != nil {
			return err
		}

		var nn []dtree.Node
		var ee []dtree.Edge

		if _config.Focus == "" {
			return errors.New("specify a focus - tree cannot be displayed")
		} else {
			if !_config.Reverse {
				nn, ee = store.FromNode(_config.Focus)
			} else {
				nn, ee = store.ToNode(_config.Focus)
			}
		}

		postProcessNodes(nn)

		dotfile, err := ioutil.TempFile("", "dtree_")

		if err != nil {
			return err
		}

		defer dotfile.Close()

		err = writeDot(nn, ee, dotfile)

		if err != nil {
			return err
		}

		data, err := executeCommand(fmt.Sprintf("dot -Tsvg %s", dotfile.Name()))

		if err != nil {
			return err
		}

		now := time.Now()

		fileName := fmt.Sprintf("output_%s_%v.svg", _config.Focus, now.Unix())
		fileName = strings.Replace(fileName, "/", "_", -1)

		err = ioutil.WriteFile(fileName, data, 0644)

		if err != nil {
			return err
		}

		_, err = executeCommand(fmt.Sprintf("xdg-open %s", fileName))

		return err
	},
}

func writeDot(nodes []dtree.Node, edges []dtree.Edge, writer io.Writer) error {

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("digraph G {\n")

	for _, edge := range edges {
		buf.WriteString(fmt.Sprintf(`"%s"->"%s" [ label="%s" ] ;`, edge.Source.Name, edge.Target.Name, edge.Version))
		buf.WriteByte('\n')
	}
	for _, node := range nodes {

		fillcolor := "white"

		typez := node.Labels["type"]

		switch typez {
		case KarhooAPI:
			fillcolor = "blue"
		case KarhooSvc:
			fillcolor = "red"
		case KarhooLibrary:
			fillcolor = "green"

		}

		buf.WriteString(fmt.Sprintf(`"%s" [fillcolor=%s style=filled];`, node.Name, fillcolor))
		buf.WriteByte('\n')
	}
	buf.WriteString("}\n")
	_, err := writer.Write(buf.Bytes())
	return err
}

type dependancyType string

var (
	Default       = dependancyType("DEFAULT")
	KarhooLibrary = dependancyType("LIB")
	KarhooSvc     = dependancyType("SVC")
	KarhooAPI     = dependancyType("API")
)

func isKarhooSvc(name string) bool {
	return strings.Contains(name, "svc") && strings.Contains(name, "karhoo")
}

func isKarhooLib(name string) bool {
	return strings.Contains(name, "lib") && strings.Contains(name, "karhoo")
}

func isKarhooAPI(name string) bool {
	return strings.Contains(name, "api-v1") && strings.Contains(name, "karhoo")
}

func postProcessNodes(nodes []dtree.Node) {

	typeStr := "type"
	for _, n := range nodes {
		if isKarhooSvc(n.Name) {
			n.Labels[typeStr] = KarhooSvc
		} else if isKarhooLib(n.Name) {
			n.Labels[typeStr] = KarhooLibrary
		} else if isKarhooAPI(n.Name) {
			n.Labels[typeStr] = KarhooAPI
		} else {
			n.Labels[typeStr] = Default
		}
	}
}

func executeCommand(cmd string) ([]byte, error) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}
