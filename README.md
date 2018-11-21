# API.md generator for go projects

## Development:

- if you update API.md.go, run `go run pack.go` in order to regenerate `API.md.tmpl.go`

## Usage

- add to Gopkg.toml: 
```
[[constraint]]
name = "gitlab.com/proemergotech/apimd-generator-go"
source = "git@gitlab.com:proemergotech/apimd-generator-go.git"
version = "0.1.0"
```
- create `<project_root>/apimd/main.go` with the example content
- run `go run apimd/main.go` from project root

### main.go example
```go
package main

import (
	"net/http"
	"strconv"
	"time"

	"gitlab.com/proemergotech/apimd-generator-go/generator"
	"gitlab.com/proemergotech/uuid-go"
)

func main() {
	g := generator.NewGenerator()
	d := &definitions{}
	g.Generate(d)
}

type definitions struct {
	factory *generator.Factory
}

type value struct {
	*generator.Value
}

func (d *definitions) Name() string {
	return "Dliver Irc Room Manager Service"
}

func (d *definitions) OutputPath() string {
	return "./API.md"
}

func (d *definitions) Usage() []string {
	return []string{
		"modify [definitions file](apimd/main.go)",
		"run `go run apimd/main.go`",
	}
}

func (d *definitions) Groups(factory *generator.Factory) []generator.Group {
	d.factory = factory
	defer func() {
		d.factory = nil
	}()

	return []generator.Group{
		&generator.HTTPGroup{
			Name:        "MyHttpGroup",
			RoutePrefix: "",
			Routes: []*generator.HTTPRoute{
				{
					Name:   "MyHttpRoute",
					Path:   "/my/:url_param/route",
					Method: http.MethodPost,
					Request: mypackage.MyRequest{
						UrlParam: d.param("example value").desc("example description").String(),
						QueryParam: d.query("example value").opt().String(),
						BodyParam: d.body("example value").Int(),
					},
					Responses: map[int]interface{}{
						http.StatusOK: mypackage.MyResponse{
							Param: d.body("example value").String(),
						},
						http.StatusBadRequest: d.httpError(apierr.MyError("my error param")),
					},
				},
			},
		},
		&generator.ConsumedMessagesGroup{
			Name:        "Messages",
			RoutePrefix: "/msg/my-service",
			Events: []*generator.GEBEvent{
				{
					Name:      "ExampleMessage",
					EventName: "/my-resource/my-action/v1",
					Body: mypackage.MyActionBody{
						Param: d.body("example value").String(),
					},
				},
			},
		},
		&generator.FiredEventsGroup{
			Name:        "Events",
			RoutePrefix: "/event/my-service",
			Events: []*generator.GEBEvent{
				{
					Name:      "ExampleEvent",
					EventName: "/my-resource/my-event/v1",
					Body: mypackage.MyEventBody{
                        Param: d.body("example value").String(),
                    },
				},
			},
		},
	}
}

func (*definitions) ParseIndex(index interface{}) (int, error) {
	switch ind := index.(type) {
	case float64:
		return int(ind), nil

	case string:
		t := &microtime.Time{}
		err := t.UnmarshalJSON([]byte("\"" + ind + "\""))
		if err == nil {
			return int(t.Unix()), nil
		}

		indInt, err := strconv.Atoi(ind)
		if err != nil {
			// use as-is
			return 0, nil
		}
		return indInt, nil

	default:
		return 0, nil
	}
}

func (d *definitions) param(val string) *value {
	return &value{Value: d.factory.Param(val)}
}

func (d *definitions) query(val string) *value {
	return &value{Value: d.factory.Query(val)}
}

func (d *definitions) body(val string) *value {
	return &value{Value: d.factory.Body(val)}
}

func (v *value) desc(d string) *value {
	v.Description(d)
	return v
}

func (v *value) opt() *value {
	v.Optional()
	return v
}

func (v *value) microtime() microtime.Time {
	return microtime.Time{Time: time.Unix(v.Int64(), 0)}
}

func (v *value) uuid() uuid.UUID {
	return uuid.UUID(v.String())
}

func (d *definitions) httpError(err error, details ...map[string]interface{}) mypackage.HTTPError {
	return mypackage.HTTPError{
	    Error: mypackage.Error{
    		Code:    d.body(apierr.Code(err)).String(),
    		Message: d.body(err.Error()).String(),
    		Details: details,
    	},
	} 
}
```