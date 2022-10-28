package furnish

import (
	"os"

	"github.com/tenderly/furnish/pkg/module/modules/shell"
	"github.com/tenderly/furnish/pkg/module/modules/ssh"
	"github.com/tenderly/furnish/pkg/module/modules/xcode"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/tenderly/furnish/pkg/module"
	"github.com/tenderly/furnish/pkg/pkgmanager"
)

var _ module.Stage = (*Stage)(nil)

type Stage struct {
	module.BaseDependable `yaml:",inline"`

	// Unknown module.Arbitrary `yaml:",inline" json:"modules"`
	XCodeSelect *xcode.XCodeSelect  `yaml:"xcode-select" json:"xcode-select"`
	SSH         *ssh.SSH            `yaml:"ssh" json:"ssh"`
	Packages    pkgmanager.Packages `yaml:"packages"     json:"packages"`
	Shell       shell.Shell         `yaml:"shell"        json:"shell"`
}

func (s *Stage) Modules() module.Modules {
	modules := make(module.Modules, 0)
	if s.XCodeSelect != nil && s.XCodeSelect.Enabled {
		modules = append(modules, s.XCodeSelect)
	}
	if s.SSH != nil && s.SSH.Enabled {
		modules = append(modules, s.SSH)
	}
	for _, s := range s.Shell {
		modules = append(modules, s)
	}
	for _, p := range s.Packages {
		modules = append(modules, p)
	}
	return modules
}

func (s *Stage) Initialize() error {
	modules := s.Modules()
	for _, m := range modules {
		m.SetParent(s.GetID())
		s.AddChildren(m.GetID())
	}
	return nil
}

type Stages map[string]*Stage

func (ss Stages) Modules() module.Modules {
	modules := make(module.Modules, 0)
	for _, p := range ss {
		modules = append(modules, p.Modules()...)
	}
	return modules
}

func (ss Stages) Stages() module.Stages {
	stages := make(module.Stages, 0, len(ss))
	for _, p := range ss {
		stages = append(stages, p)
	}
	return stages
}

func (ss Stages) Initialize() error {
	for id, phase := range ss {
		// This happens if someone creates an empty phase due to unmarshaling
		// It's better to allow it and skip than panic
		if phase == nil {
			ss[id] = &Stage{BaseDependable: module.BaseDependable{}}
			phase = ss[id]
		}

		phase.ID = module.ID(id)
		phase.Initialize()
		color.Blue("[init] initialized phase '%s'", phase.GetID())
	}
	return nil
}

type Global struct {
	PackageManagers pkgmanager.MultiManagerConfig `yaml:"package-managers" json:"package_managers"`
}

func (g *Global) Validate() error { return g.PackageManagers.Validate() }

func (g *Global) Initialize() error { return g.PackageManagers.Initialize() }

type Declaration struct {
	FileVersion module.Version `yaml:"version" json:"file_version"`
	Global      Global         `yaml:"global"  json:"global"`
	Phases      Stages         `yaml:",inline" json:"stages"`
}

func (d *Declaration) Modules() module.Modules { return d.Phases.Modules() }

func (d *Declaration) Validate() error { return d.Global.Validate() }

func (d *Declaration) Initialize() error {
	if err := d.Global.Initialize(); err != nil {
		return errors.Wrap(err, "initializing global")
	}
	if err := d.Phases.Initialize(); err != nil {
		return errors.Wrap(err, "initializing stages")
	}
	return nil
}

func (d *Declaration) Stages() module.Stages { return d.Phases.Stages() }

func Load(path string) (*Declaration, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "config file read")
	}

	decl := &Declaration{}
	if err := yaml.Unmarshal(file, decl); err != nil {
		return nil, errors.Wrap(err, "unmarshal config")
	}

	return decl, nil
}
