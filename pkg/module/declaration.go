package module

type Declaration interface {
	Initialize() error
	Modules() Modules
	Validate() error
	Stages() Stages
}

type Stage interface {
	Dependable

	Initialize() error
	Modules() Modules
}
