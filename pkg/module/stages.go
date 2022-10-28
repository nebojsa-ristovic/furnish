package module

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/tenderly/furnish/pkg/util"

	"github.com/fatih/color"

	"github.com/tenderly/furnish/pkg/log"
)

var (
	fmtFailedMandatory = color.HiRedString("[fatal] failed module is mandatory. aborting\n")
	fmtFailedDep       = color.HiRedString("[fatal] the failed module is a dependency for %s. aborting\n")
	fmtErrorApply      = color.RedString(
		"  [%d/%d] '%s' error applying module.",
	) + color.WhiteString(
		"\n\tmeta: [%s]\n",
	) + color.RedString(
		"  err: %s\n",
	)
	fmtSkip    = color.YellowString("  [%d/%d] '%s' skipped - already applied") + color.WhiteString("\n\tmeta: [%s]\n")
	fmtSuccess = color.GreenString("  [%d/%d] '%s' applied.") + color.WhiteString("\n\tmeta: [%s]\n")
)

type Stages []Stage

func (s Stages) Apply(ctx context.Context) ([]bool, error) {
	return newStagesApplier(s).ApplyMany(ctx)
}

type stagesApplier struct {
	stages Stages
}

func newStagesApplier(s Stages) *stagesApplier {
	dependables := make(Dependables, 0, len(s))
	for _, s := range s {
		dependables = append(dependables, s.(Dependable))
	}
	dependables = dependables.Sort()

	sortedStages := make(Stages, 0, len(dependables))
	for _, r := range dependables {
		stage, ok := r.(Stage)
		if !ok {
			log.Info("relator not a stage")
			continue
		}
		sortedStages = append(sortedStages, stage)
	}

	return &stagesApplier{stages: sortedStages}
}

func (dma *stagesApplier) ApplyMany(ctx context.Context) ([]bool, error) {
	for _, s := range dma.stages {
		modules := s.Modules()
		if len(modules) == 0 {
			color.Yellow("stage '%s' empty, skipping", s.GetID())
			continue
		}

		dependables := make(Dependables, 0, len(modules))
		for _, m := range modules {
			dependables = append(dependables, m.(Dependable))
		}
		dependables = dependables.Sort()

		sorted := make(Modules, 0, len(dependables))
		for _, r := range dependables {
			module, ok := r.(Module)
			if !ok {
				log.Info("relator not a module")
				continue
			}
			sorted = append(sorted, module)
		}

		if _, err := dma.applyMany(ctx, s.GetID(), sorted); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (dma *stagesApplier) applyMany(ctx context.Context, stage ID, modules Modules) ([]bool, error) {
	if dma == nil {
		return nil, errors.New("no packages found")
	}

	color.Blue("[%s]", stage)
	failed, skipped, applied := make(map[ID]error), make(IDs, 0), make(IDs, 0)
	total := len(modules)

	for i, m := range modules {
		if m.IsOptional() {
			if !util.ReadConfirmation(
				fmt.Sprintf("Module %s is optional.\nIf you wish to install it press Y/y.", m.GetID()),
				"y",
			) {
				skipped = append(skipped, m.GetID())
				continue
			}
		}
		ok, meta, err := m.Apply(ctx)
		if err != nil {
			fmt.Printf(fmtErrorApply, i+1, total, m.GetID(), meta, err.Error())
			log.Debug("error applying module", "module", m, "err", err)

			if m.IsDependency() {
				fmt.Printf(fmtFailedDep, m.GetDependants())
				os.Exit(1)
			}
			if m.IsMandatory() {
				fmt.Printf(fmtFailedMandatory)
				os.Exit(1)
			}
			failed[m.GetID()] = err
			continue
		}
		if !ok {
			fmt.Printf(fmtSkip, i+1, total, m.GetID(), meta)
			skipped = append(skipped, m.GetID())
			continue
		}

		applied = append(applied, m.GetID())
		fmt.Printf(fmtSuccess, i+1, total, m.GetID(), meta)
	}

	dma.printResults(stage, len(applied), len(skipped), len(failed), total)

	return []bool{}, nil
}

func (dma *stagesApplier) printResults(stage ID, applied, skipped, failed, total int) {
	color.White("\n")
	// color.White("summary: ")
	if applied > 0 {
		color.Green("[✔] applied %d of %d modules", applied, total)
	}
	if skipped > 0 {
		color.Yellow("[-] skipped %d of %d modules", skipped, total)
	}
	if failed > 0 {
		color.Red("[✘] failed applying %d of %d modules", failed, total)
	}
	color.White("\n")
}
