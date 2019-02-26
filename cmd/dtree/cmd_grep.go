package main

import (
	"errors"

	"github.com/mdevilliers/dtree"
	"github.com/spf13/cobra"
)

var usageStr = "grep [TERM]"

var grepCommand = &cobra.Command{
	Use:   usageStr,
	Short: "grep for versions and dependancies of a package.",
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return errors.New("specify a [TERM]")
		}

		store, err := initStore(_config)
		if err != nil {
			return err
		}

		fragment := args[0]

		if fragment == "" {
			return errors.New(usageStr)
		}

		matches, err := store.FindNodes(fragment)

		if err != nil {
			return err
		}

		nn := dtree.Nodes{}
		ee := dtree.Edges{}

		for _, match := range matches {
			if !_config.Reverse {
				n, e := store.FromNode(match)
				nn = append(nn, n...)
				ee = append(ee, e...)
			} else {
				n, e := store.ToNode(match)
				nn = append(nn, n...)
				ee = append(ee, e...)
			}
		}

		nn = nn.Predicate(_config.Predicate)
		ee = ee.Predicate(_config.Predicate)

		if _config.ToDot {
			return outputDot(fragment, nn, ee)
		} else if _config.ToSvg {
			return outputSvg(fragment, nn, ee)
		}
		return errors.New("no output configured")

	},
}
