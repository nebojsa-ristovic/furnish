package xcode

import (
	"context"
	"strings"

	"github.com/tenderly/furnish/pkg/module/modules/shell"

	"github.com/pkg/errors"

	"github.com/tenderly/furnish/pkg/module"
)

type Declaration struct {
	XCodeSelect bool `yaml:"xcode-select"`
}

var _ module.Module = (*XCodeSelect)(nil)

type XCodeSelect struct {
	module.BaseDependable `yaml:",inline"`

	Enabled   bool `yaml:"enabled"`
	Update    bool `yaml:"update"`
	Mandatory bool `yaml:"mandatory"`
}

func (_ *XCodeSelect) IsOptional() bool { return false }

func (x *XCodeSelect) IsMandatory() bool { return x.Mandatory }

func (x *XCodeSelect) Apply(ctx context.Context) (bool, string, error) {
	output, err := shell.ExecOutput("arch -arm64 xcode-select -p")
	if err != nil {
		return false, "exists", errors.Wrap(err, "failed checking if xcode-select is installed")
	}
	if strings.Contains(strings.ToLower(output), "developer") {
		return false, "install", nil
	}
	if err := shell.Exec("arch -arm64 xcode-select --install"); err != nil {
		return false, "install", errors.Wrap(err, "failed installing xcode-select")
	}
	return true, "install", nil
}

func (x *XCodeSelect) GetVersion() module.Version {
	output, err := shell.ExecOutput("xcode-select --version | cut -f3 -d' '")
	if err != nil {
		return x.Version
	}
	output = strings.Replace(output, "\n", "", 1)
	x.Version = module.Version(strings.Replace(output, ".", "", 1))
	return x.Version
}

func (x *XCodeSelect) GetID() module.ID { return "xcode-select" }
