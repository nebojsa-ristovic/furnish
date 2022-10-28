package module

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/stevenle/topsort"

	"github.com/tenderly/furnish/pkg/log"
)

var fmtNoDep = color.RedString("module is a dependency of") + color.HiMagentaString(" '%s' ") +
	color.RedString("but not found:") + color.HiMagentaString(" '%s' \n")
var fmtRemove = color.HiRedString("[fatal] please add the module or remove it as a dependency\n")

type Related interface {
	GetParent() ID
	SetParent(ID)
	GetChildren() IDs
	AddChildren(...ID)
}

type Dependable interface {
	Identifier

	GetDependencies() IDs
	GetDependants() IDs
	IsDependency() bool
	AddDependencies(...ID)
	AddDependants(...ID)

	mustEmbedBaseDependable()
}

var (
	_ Dependable = (*BaseDependable)(nil)
	_ Related    = (*BaseDependable)(nil)
)

type BaseDependable struct {
	ID           ID      `yaml:"id"           json:"id"`
	Version      Version `yaml:"version"      json:"version"`
	Description  string  `yaml:"description"  json:"description"`
	Dependencies IDs     `yaml:"dependencies" json:"dependencies"`
	Dependants   IDs     `yaml:"-"            json:"dependants"`

	Children IDs `yaml:"-" json:"children,omitempty"`
	Parent   ID  `yaml:"-" json:"parent"`
}

func (bd *BaseDependable) GetID() ID { return bd.ID }

func (bd *BaseDependable) GetDescription() string { return bd.Description }

func (bd *BaseDependable) GetVersion() Version { return bd.Version }

func (bd *BaseDependable) GetDependencies() IDs { return bd.Dependencies }

func (bd *BaseDependable) GetDependants() IDs { return bd.Dependants }

func (bd *BaseDependable) IsDependency() bool { return len(bd.Dependants) > 0 }

func (bd *BaseDependable) AddDependencies(ids ...ID) {
	bd.Dependencies = append(bd.Dependencies, ids...)
}

func (bd *BaseDependable) AddDependants(ids ...ID) {
	bd.Dependants = append(bd.Dependants, ids...)
}

func (bd *BaseDependable) GetParent() ID { return bd.Parent }

func (bd *BaseDependable) SetParent(id ID) { bd.Parent = id }

func (bd *BaseDependable) GetChildren() IDs { return bd.Children }

func (bd *BaseDependable) AddChildren(ids ...ID) { bd.Children = append(bd.Children, ids...) }

func (bd *BaseDependable) mustEmbedBaseDependable() {}

type Dependables []Dependable

func (dd Dependables) Map() map[ID]Dependable {
	dependerMap := make(map[ID]Dependable, len(dd))
	for _, r := range dd {
		dependerMap[r.GetID()] = r
	}
	return dependerMap
}

func (dd Dependables) Sort() Dependables {
	if len(dd) == 0 {
		return dd
	}

	sortedRelaters := make(Dependables, 0, len(dd))
	added := make(map[ID]struct{}, 0)
	for _, r := range dd {
		if len(r.GetDependencies()) == 0 {
			log.Debug("no deps", "r", r.GetID())
			sortedRelaters = append(sortedRelaters, r)
			added[r.GetID()] = struct{}{}
		}
	}

	relaterMap := dd.Map()
	graph := topsort.NewGraph()
	for _, r := range dd {
		graph.AddNode(r.GetID().String())
	}
	for _, r := range dd {
		for _, d := range r.GetDependencies() {
			graph.AddEdge(r.GetID().String(), d.String())
			dependency, ok := relaterMap[d]
			if !ok {
				fmt.Printf(fmtNoDep, r.GetID(), dependency)
				fmt.Printf(fmtRemove)
				os.Exit(1)
			}
			dependency.AddDependants(r.GetID())
		}
	}

	for _, r := range dd {
		if _, ok := added[r.GetID()]; ok {
			continue
		}

		sorted, edd := graph.TopSort(r.GetID().String())
		if edd != nil {
			log.Error("something went wrong when sorting", "edd", edd)
		}

		for _, id := range sorted {
			if _, ok := added[ID(id)]; ok {
				continue
			}
			relater, _ := relaterMap[ID(id)]
			sortedRelaters = append(sortedRelaters, relater)
			added[ID(id)] = struct{}{}
		}
	}

	return sortedRelaters
}
