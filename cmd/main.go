package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	furnish "github.com/tenderly/furnish/pkg"
	"github.com/tenderly/furnish/pkg/log"
)

func main() {
	log := log.NewZapLogger(&log.Config{LogVerbosity: log.InfoVerbosity, ServiceName: "mac-wrap"})
	app := cli.NewApp()
	app.Name = "furnish"

	app.Commands = append(app.Commands, DebugPrintCmd(), RunCmd())
	if err := app.Run(os.Args); err != nil {
		log.Error("failed running app", "err", err)
	}
}

func RunCmd() *cli.Command {
	return &cli.Command{
		Name:        "apply",
		Description: "applies the modules",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "--config example.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			cfgPath := c.String("config")
			if cfgPath == "" {
				cfgPath = "furnish.yaml"
			}

			decl, err := furnish.Load(cfgPath)
			if err != nil {
				return errors.Wrap(err, "reading config")
			}

			if err = decl.Validate(); err != nil {
				return errors.Wrap(err, "validate")
			}

			if err = decl.Initialize(); err != nil {
				return errors.Wrap(err, "initialize")
			}

			color.White("\n\n")
			decl.Stages().Apply(c.Context)

			return nil
		},
	}
}

func DebugPrintCmd() *cli.Command {
	return &cli.Command{
		Name:        "debug",
		Description: "prints the configuration file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "--config example.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			cfgPath := c.String("config")
			if cfgPath == "" {
				cfgPath = "furnishl.yaml"
			}

			decl, err := furnish.Load(cfgPath)
			if err != nil {
				return errors.Wrap(err, "reading config")
			}

			log.Info("configuration", "cfg", decl)

			color.White("\n\n")
			_, err = decl.Stages().Apply(c.Context)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
