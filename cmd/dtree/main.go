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
	Root    string `envconfig:"ROOT" default:"/home/mark/go/src/github.com/karhoo/"`
	Reverse bool   `envconfig:"REVERSE" default:"false"`
	ToDot   bool   `envconfig:"DOT" default:"false"`
	ToSvg   bool   `envconfig:"SVG" default:"true"`
}

func newConfig() *config {
	c := &config{}
	envconfig.Process("", c)
	return c
}

func (o *config) addFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Root, "root", "r", o.Root, "path to repositories")
	fs.BoolVarP(&o.Reverse, "reverse", "e", o.Reverse, "reverse the focus")
	fs.BoolVarP(&o.ToDot, "dot", "d", o.ToDot, "output to dot file")
	fs.BoolVarP(&o.ToSvg, "svg", "s", o.ToSvg, "output to svg (requires 'dot' installed)")
	fs.Parse(os.Args)
}

func init() {
	rootCmd.AddCommand(grepCommand)
	_config.addFlags(rootCmd.PersistentFlags())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
