package pkgmanager

import (
	"context"

	"github.com/pkg/errors"
)

type ManagerName string

type Defaults ManagerInfo

type ManagerInfo interface {
	Name() ManagerName
	Path() string
	BinaryExists() bool
	HowToInstall() string
	Cmd() string
}

type ApplyFunc func(context.Context, *Package) error
type ExistsFunc func(context.Context, *Package) (bool, error)

type Manager interface {
	ManagerInfo

	Install(ctx context.Context, pkg *Package) error
	Exists(ctx context.Context, pkg *Package) (bool, error)
	Update(ctx context.Context, pkg *Package) error
	Delete(ctx context.Context, pkg *Package) error
}

type ManagerProvider interface {
	Provide(name ManagerName) (Manager, error)
	Default() (Manager, error)
}

type managerProvider struct {
	def      ManagerName
	managers map[ManagerName]Manager
}

func (mp *managerProvider) Provide(name ManagerName) (Manager, error) {
	if m, ok := mp.managers[name]; ok {
		return m, nil
	}
	return nil, errors.Errorf("couldn't find manager with name %s", name)
}

func (mp *managerProvider) Default() (Manager, error) {
	if mp.def == "" {
		return nil, errors.New("no default manager set")
	}
	return mp.Provide(mp.def)
}

var globalManagerProvider = managerProvider{
	managers: map[ManagerName]Manager{},
}

func (mp *managerProvider) register(m Manager, def bool) {
	if _, ok := mp.managers[m.Name()]; !ok {
		mp.managers[m.Name()] = m
		if def && mp.def == "" {
			mp.def = m.Name()
		}
	}
}

func ConfigureManagers(mmc MultiManagerConfig) error {
	if err := mmc.Validate(); err != nil {
		return err
	}
	for _, cfg := range mmc {
		m, err := configureManager(cfg)
		if err != nil {
			return errors.Wrap(err, "couldn't configure manager")
		}
		globalManagerProvider.register(m, cfg.Default)
	}

	return nil
}

func configureManager(cfg *Config) (Manager, error) {
	if cfg.Name == TypeBrew {
		return NewBrewPackageManager(cfg)
	}
	return nil, errors.New("package manager not supported")
}

func ProvideManager(name ManagerName) (Manager, error) { return globalManagerProvider.Provide(name) }
func Default() (Manager, error)                        { return globalManagerProvider.Default() }
