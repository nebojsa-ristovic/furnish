package ssh

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tenderly/furnish/pkg/module"
	"github.com/tenderly/furnish/pkg/module/modules/shell"
)

var _ module.Module = (*SSH)(nil)

type SSH struct {
	module.BaseDependable `yaml:",inline"`

	Enabled   bool `yaml:"enabled"`
	Optional  bool `yaml:"optional"`
	Mandatory bool `yaml:"mandatory"`

	Output     string `yaml:"output"`
	Type       string `yaml:"type"`
	Passphrase string `yaml:"passphrase"`
	Comment    string `yaml:"comment"`
}

func (s *SSH) IsOptional() bool {
	return s.Enabled && s.Optional
}

func (s *SSH) IsMandatory() bool {
	return s.Enabled && s.Mandatory
}

func (s *SSH) build() string {
	cmd := "ssh-keygen"
	if s.Output != "" {
		cmd += fmt.Sprintf(" -f %s", s.Output)
	}
	if s.Type != "" {
		cmd += fmt.Sprintf(" -t %s", s.Type)
	}
	if s.Passphrase != "" {
		cmd += fmt.Sprintf(" -p %s", s.Passphrase)
	}
	if s.Comment != "" {
		cmd += fmt.Sprintf(" -C %s", s.Comment)
	}

	return cmd
}

func (s *SSH) Apply(ctx context.Context) (bool, string, error) {
	err := shell.Exec(s.build())
	if err != nil {
		return false, "generate", errors.Wrap(err, "failed generating ssh key")
	}

	return true, "generate", nil
}

func (s *SSH) GetID() module.ID {
	return "ssh"
}
