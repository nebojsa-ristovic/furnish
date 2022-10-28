package module

import (
	"context"
)

type Arbitrary = map[string]interface{}

type Module interface {
	Dependable
	Related

	IsOptional() bool
	IsMandatory() bool
	Apply(context.Context) (bool, string, error)
}

type Modules []Module

func (mm Modules) Map() map[ID]Module {
	moduleMap := make(map[ID]Module, len(mm))
	for _, m := range mm {
		moduleMap[m.GetID()] = m
	}
	return moduleMap
}
