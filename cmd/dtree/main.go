package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/mdevilliers/dtree"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use: "dtree",
}

/*
go build && ./dtree XXXXX > example.dot && dot -Tsvg example.dot > example.svg
*/

var _config = newConfig()

type config struct {
	Root    string `envconfig:"ROOT" default:"/home/mark/go/src/github.com/karhoo/"`
	Focus   string `envconfig:"FOCUS" default:"api-v1-estimates"`
	Reverse bool   `envconfig:"REVERSE" default:"false"`
}

func newConfig() *config {
	c := &config{}
	envconfig.Process("", c)
	return c
}

func (o *config) addFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Root, "root", "r", o.Root, "path to repositories")
	fs.StringVarP(&o.Focus, "focus", "f", o.Focus, "dependancy you are interested in exploring")
	fs.BoolVarP(&o.Reverse, "reverse", "e", o.Reverse, "reverse the focus")
	fs.Parse(os.Args)
}

func init() {
	rootCmd.AddCommand(generateDiagramCommand)
	_config.addFlags(rootCmd.PersistentFlags())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
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
