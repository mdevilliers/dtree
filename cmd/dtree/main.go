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
	rootCmd.AddCommand(listCommand)
	_config.addFlags(rootCmd.PersistentFlags())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
