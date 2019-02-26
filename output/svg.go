package output

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/mdevilliers/dtree"
)

// ToDot will write the tree of Nodes and Edges to the supplied io.Writer or return an error
func ToDot(nodes []dtree.Node, edges []dtree.Edge, writer io.Writer) error {

	buf := bytes.NewBuffer([]byte{})
	_, err := buf.WriteString("digraph G {\n")

	if err != nil {
		return err
	}

	for _, edge := range edges {
		_, err = buf.WriteString(fmt.Sprintf(`"%s"->"%s" [ label="%s" ] ;`, edge.Source.Name, edge.Target.Name, edge.Version))

		if err != nil {
			return err
		}
		err = buf.WriteByte('\n')

		if err != nil {
			return err
		}
	}
	for _, node := range nodes {

		fillcolor := "white"

		switch {

		case isAPI(node.Name):
			fillcolor = "blue"
		case isSvc(node.Name):
			fillcolor = "red"
		case isLib(node.Name):
			fillcolor = "green"

		}
		_, err = buf.WriteString(fmt.Sprintf(`"%s" [fillcolor=%s style=filled];`, node.Name, fillcolor))

		if err != nil {
			return err
		}

		err = buf.WriteByte('\n')

		if err != nil {
			return err
		}
	}
	_, err = buf.WriteString("}\n")

	if err != nil {
		return err
	}

	_, err = writer.Write(buf.Bytes())
	return err
}

func isSvc(name string) bool {
	return strings.Contains(name, "svc")
}

func isLib(name string) bool {
	return strings.Contains(name, "lib")
}

func isAPI(name string) bool {
	return strings.Contains(name, "api")
}
