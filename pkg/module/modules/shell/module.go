package shell

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/tenderly/furnish/pkg/module"
)

type Shell []*Execution

var _ module.Module = (*Execution)(nil)

type Execution struct {
	module.BaseDependable `yaml:",inline"`

	Name      module.ID `yaml:"name"   json:"name,omitempty"`
	Cmd       string    `yaml:"cmd"    json:"cmd,omitempty"`
	File      string    `yaml:"file"   json:"file,omitempty"`
	Silent    bool      `yaml:"silent" json:"silent,omitempty"`
	Mandatory bool      `yaml:"mandatory" json:"mandatory,omitempty"`

	meta string `yaml:"-" json:"-"`
}

func (x *Execution) GetID() module.ID { return x.Name }

func (x *Execution) Validate() error {
	if x.Name == "" {
		return errors.New("shell execution must have name")
	}
	if x.Cmd == "" || x.File == "" {
		return errors.New("shell execution must have a cmd or a script provided")
	}
	return nil
}

func (x *Execution) IsOptional() bool { return false }

func (x *Execution) IsMandatory() bool { return x.Mandatory }

func (x *Execution) Apply(ctx context.Context) (bool, string, error) {
	x.setMeta()
	if x.File != "" {
		return x.applyScript(ctx)
	}
	return x.applyCmd(ctx)
}

func (x *Execution) applyScript(ctx context.Context) (bool, string, error) {
	if x.Silent {
		if err := ScriptSilent(x.File); err != nil {
			return false, x.meta, errors.Wrap(err, "couldn't exec script")
		}
		return true, x.meta, nil
	}
	if err := Script(x.File); err != nil {
		return false, x.meta, errors.Wrap(err, "couldn't exec script")
	}
	return true, x.meta, nil
}
func (x *Execution) applyCmd(ctx context.Context) (bool, string, error) {
	if x.Silent {
		if err := ExecSilent(x.Cmd); err != nil {
			return false, x.meta, errors.Wrap(err, "couldn't exec command")
		}
		return true, x.meta, nil
	}
	if err := Exec(x.Cmd); err != nil {
		return false, x.meta, errors.Wrap(err, "couldn't exec command")
	}
	return true, x.meta, nil
}

func (x *Execution) setMeta() {
	if x.File != "" {
		x.meta = fmt.Sprintf("mode: file; silent: %t; path: '%s'", x.Silent, x.File)
		return
	}
	x.meta = fmt.Sprintf("mode: cmd; silent: %t; cmd: '%s'", x.Silent, x.Cmd)
}
