package pkgmanager

import (
	"errors"

	"github.com/fatih/color"

	"github.com/tenderly/furnish/pkg/module/modules/shell"
	"github.com/tenderly/furnish/pkg/util"
)

var supprotedPackageManagers = map[ManagerName]Defaults{
	TypeBrew: macOSBrewDefaults,
}

type Config struct {
	Name    ManagerName `yaml:"name"    json:"name"`
	Path    string      `yaml:"path"    json:"path"`
	Default bool        `yaml:"default" json:"default"`

	// Prefix is a specific macos field
	Prefix string `yaml:"prefix" json:"prefix"`
}

func (c *Config) Validate() error {
	if c.Name == "" {
		return errors.New("no name for package manager")
	}
	defaults, ok := supprotedPackageManagers[c.Name]
	if !ok {
		return errors.New("unsupported package manager")
	}
	if c.Path == "" {
		c.Path = defaults.Path()
	}

	if ok := shell.BinaryExists(c.Path); !ok {
		color.Yellow("[warn] passed path package manager not found")
		if ok := defaults.BinaryExists(); ok {
			color.HiBlue("[init] default path found. using default path %s", defaults.Path())
			c.Path = defaults.Path()
			return nil
		}
		if defaults.HowToInstall() != "" &&
			util.ReadConfirmation(
				"Brew not found but we can install it.\nIf you wish to install brew press Y/y.",
				"y",
			) {
			return shell.Exec(defaults.HowToInstall())
		}
		color.Red("package manager %s not found, or it's the wrong path", c.Name)
		return errors.New("package manager doesn't exist, or it's the wrong path")
	}

	return nil
}

type MultiManagerConfig []*Config

func (mmc MultiManagerConfig) Validate() error {
	if len(mmc) == 0 {
		return errors.New("no manager provided, need to provide atleast 1 manager as the default")
	}
	if len(mmc) == 1 {
		mmc[0].Default = true
	}
	for _, c := range mmc {
		if err := c.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (mmc MultiManagerConfig) Initialize() error { return ConfigureManagers(mmc) }
