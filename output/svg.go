package output

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/mdevilliers/dtree"
)

func ToDot(nodes []dtree.Node, edges []dtree.Edge, writer io.Writer) error {

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("digraph G {\n")

	for _, edge := range edges {
		buf.WriteString(fmt.Sprintf(`"%s"->"%s" [ label="%s" ] ;`, edge.Source.Name, edge.Target.Name, edge.Version))
		buf.WriteByte('\n')
	}
	for _, node := range nodes {

		fillcolor := "white"

		if isAPI(node.Name) {
			fillcolor = "blue"
		} else if isSvc(node.Name) {
			fillcolor = "red"
		} else if isLib(node.Name) {
			fillcolor = "green"
		}

		buf.WriteString(fmt.Sprintf(`"%s" [fillcolor=%s style=filled];`, node.Name, fillcolor))
		buf.WriteByte('\n')
	}
	buf.WriteString("}\n")
	_, err := writer.Write(buf.Bytes())
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
