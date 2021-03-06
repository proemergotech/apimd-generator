package generator

type Definitons interface {
	Name() string
	OutputPath() string
	Usage() []string
	Groups(f *Factory) []Group
	ParseIndex(index interface{}) (int, error)
}
