package generator

type Document struct {
	Name       string
	Usage      []string
	Categories []*DocCategory
}

type DocCategory struct {
	Name   string
	Groups []*DocGroup
}

type DocGroup struct {
	Name   string
	Prefix string
	Routes []*DocRoute
}

type DocRoute struct {
	Name           string
	Method         string
	Path           string
	Description    []string
	Params         map[string]*DocValue
	Query          map[string]*DocValue
	RequestBody    interface{}
	ResponseBodies map[int]interface{}
}

type DocValue struct {
	Value     string
	Desc      string
	Opt       bool
	APIMDType string
}
