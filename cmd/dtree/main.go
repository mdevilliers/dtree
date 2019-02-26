package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use: "dtree",
}

var _config = newConfig()

type config struct {
	Root      string `envconfig:"ROOT" default:"/home/mark/go/src/github.com/karhoo/"`
	Reverse   bool   `envconfig:"REVERSE" default:"false"`
	ToDot     bool   `envconfig:"DOT" default:"false"`
	ToSvg     bool   `envconfig:"SVG" default:"true"`
	Predicate string `envconfig:"PREDICATE" default:""`
}

func newConfig() *config {
	c := &config{}
	err := envconfig.Process("", c)

	if err != nil {
		panic(err)
	}

	return c
}

func (o *config) addFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Root, "root", "r", o.Root, "path to repositories")
	fs.BoolVarP(&o.Reverse, "reverse", "e", o.Reverse, "reverse the focus")
	fs.BoolVarP(&o.ToDot, "dot", "d", o.ToDot, "output to dot file")
	fs.BoolVarP(&o.ToSvg, "svg", "s", o.ToSvg, "output to svg (requires 'dot' installed)")
	fs.StringVarP(&o.Predicate, "predicate", "p", o.Predicate, "trim responses to only include nodes and edges containing this value")
	err := fs.Parse(os.Args)

	if err != nil {
		panic(err)
	}
}

func main() {

	rootCmd.AddCommand(grepCommand)
	_config.addFlags(rootCmd.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
