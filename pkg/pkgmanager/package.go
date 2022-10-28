package pkgmanager

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/tenderly/furnish/pkg/module"
)

type assertExists func(ctx context.Context) (bool, error)

type pkgApplier string

const (
	pkgApplierInstall = "install"
	pkgApplierUpdate  = "update"
	pkgApplierDelete  = "delete"
	pkgApplierEmpty   = ""
)

var _ module.Module = (*Package)(nil)

type Package struct {
	module.BaseDependable `yaml:",inline"`

	Name      module.ID   `yaml:"name"    json:"name"`
	Applier   pkgApplier  `yaml:"applier" json:"applier"`
	Source    string      `yaml:"source"  json:"source"`
	Manager   ManagerName `yaml:"manager" json:"manager"`
	Optional  bool        `yaml:"optional" json:"optional"`
	Mandatory bool        `yaml:"mandatory" json:"mandatory"`

	meta string `yaml:"-" json:"-"`
}

func (p *Package) GetID() module.ID { return p.Name }

func (p *Package) IsMandatory() bool { return p.Mandatory }

func (p *Package) IsOptional() bool { return p.Optional }

func (p *Package) Apply(ctx context.Context) (bool, string, error) {
	m, err := p.manager()
	if err != nil {
		return false, "", err
	}
	p.setMeta(m)

	switch p.Applier {
	case pkgApplierInstall, pkgApplierEmpty:
		return p.apply(ctx, m.Install, p.assertExists(m.Exists, false))
	case pkgApplierUpdate:
		return p.apply(ctx, m.Update, p.assertExists(m.Exists, true))
	case pkgApplierDelete:
		return p.apply(ctx, m.Delete, p.assertExists(m.Exists, true))
	default:
		return p.apply(ctx, nil, func(context.Context) (bool, error) { return false, nil })
	}
}

func (p *Package) apply(ctx context.Context, apply ApplyFunc, exists assertExists) (bool, string, error) {
	if ok, err := exists(ctx); err != nil || !ok {
		return ok, p.meta, err
	}
	if err := apply(ctx, p); err != nil {
		return false, p.meta, errors.Wrap(err, fmt.Sprintf("couldn't %s package", p.Applier))
	}
	return true, p.meta, nil
}

func (p *Package) assertExists(existsFn ExistsFunc, expected bool) assertExists {
	return func(ctx context.Context) (bool, error) {
		exists, err := existsFn(ctx, p)
		if err != nil {
			return false, errors.Wrap(err, "couldn't check if package exists")
		}
		if exists != expected {
			return false, nil
		}
		return true, nil
	}
}

func (p *Package) manager() (Manager, error) {
	manager, err := Default()
	if err != nil && p.Manager != "" {
		return nil, errors.Wrap(err, "couldn't provide default manager")
	}
	if p.Manager != "" {
		custom, err := ProvideManager(p.Manager)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't provide set manager")
		}
		manager = custom
	}
	return manager, nil
}

func (p *Package) setMeta(m Manager) {
	if p.Applier == pkgApplierEmpty {
		p.Applier = pkgApplierInstall
	}
	p.meta = fmt.Sprintf("applier: %s; manager: %s", p.Applier, m.Name())
}

type Packages []*Package
