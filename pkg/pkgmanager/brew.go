package pkgmanager

import (
	"context"
	"fmt"

	"github.com/tenderly/furnish/pkg/module/modules/shell"

	"github.com/fatih/color"
	"github.com/klauspost/cpuid/v2"
	"github.com/pkg/errors"
)

const TypeBrew ManagerName = "brew"

var _ ManagerInfo = (*macOSBrewInfo)(nil)

type macOSBrewInfo struct {
	path           string
	prefix         string
	installCommand string
}

func (mbi *macOSBrewInfo) HowToInstall() string { return mbi.installCommand }

func (*macOSBrewInfo) Name() ManagerName { return TypeBrew }

func (mbi *macOSBrewInfo) Cmd() string { return fmt.Sprintf("%s %s", mbi.prefix, TypeBrew) }

func (mbi *macOSBrewInfo) Path() string { return mbi.path }

func (mbi *macOSBrewInfo) BinaryExists() bool { return shell.BinaryExists(mbi.Path()) }

var macOSBrewDefaults = newMacOSBrewDefaults()

func newMacOSBrewDefaults() *macOSBrewInfo {
	path := "/usr/local/bin/brew"
	prefix := "arch -x86_64"
	if cpuid.CPU.VendorID != cpuid.Intel && cpuid.CPU.VendorID != cpuid.AMD {
		path = "/opt/homebrew/bin/brew"
		prefix = "arch -arm64"
	}
	return &macOSBrewInfo{
		path:           path,
		prefix:         prefix,
		installCommand: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`,
	}
}

type BrewPackageManager struct {
	info ManagerInfo
}

func NewBrewPackageManager(cfg *Config) (Manager, error) {
	if cfg.Prefix == "" {
		cfg.Prefix = macOSBrewDefaults.prefix
	}
	info := &macOSBrewInfo{
		path:           cfg.Path,
		prefix:         cfg.Prefix,
		installCommand: macOSBrewDefaults.installCommand,
	}

	color.HiBlue("[init] brew initialized, using cmd: %s", info.Cmd())
	return &BrewPackageManager{info: info}, nil
}

func (b *BrewPackageManager) Cmd() string { return b.info.Cmd() }

func (b *BrewPackageManager) Name() ManagerName { return b.info.Name() }

func (b *BrewPackageManager) Path() string { return b.info.Path() }

func (b *BrewPackageManager) BinaryExists() bool { return b.info.BinaryExists() }

func (b *BrewPackageManager) Install(ctx context.Context, pkg *Package) error {
	if err := shell.ExecSilent(fmt.Sprintf("%s install %s", b.Cmd(), pkg.Name)); err != nil {
		return errors.New("package doesn't exist")
	}
	return nil
}

func (b *BrewPackageManager) Exists(ctx context.Context, pkg *Package) (bool, error) {
	if err := shell.ExecSilent(fmt.Sprintf("%s list %s", b.Cmd(), pkg.Name)); err != nil {
		return false, nil
	}
	return true, nil
}

func (b *BrewPackageManager) Update(ctx context.Context, pkg *Package) error {
	if err := shell.ExecSilent(fmt.Sprintf("%s upgrade %s", b.Cmd(), pkg.Name)); err != nil {
		return errors.New("couldn't update package")
	}
	return nil
}

func (b *BrewPackageManager) Delete(ctx context.Context, pkg *Package) error {
	if err := shell.Exec(fmt.Sprintf("%s uninstall %s", b.Cmd(), pkg.Name)); err != nil {
		return errors.New("couldn't uninstall package")
	}
	return nil
}

func (b *BrewPackageManager) HowToInstall() string { return b.info.HowToInstall() }
