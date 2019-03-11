package generator

import "net/http"

type Group interface {
	GetName() string
	GetRoutePrefix() string
	GetRoutes() []*Route
	GetCategory() string
}

type Route struct {
	Name        string
	Method      string
	Path        string
	Description []string
	Request     interface{}
	Responses   map[int]interface{}
}

type HTTPGroup struct {
	Name        string
	RoutePrefix string
	Routes      []*HTTPRoute
}

type HTTPRoute struct {
	Name        string
	Method      string
	Path        string
	Description []string
	Request     interface{}
	Responses   map[int]interface{}
}

type ConsumedMessagesGroup struct {
	Name        string
	RoutePrefix string
	Events      []*GEBEvent
}

type FiredEventsGroup struct {
	Name        string
	RoutePrefix string
	Events      []*GEBEvent
}

type GEBEvent struct {
	Name        string
	EventName   string
	Description []string
	Body        interface{}
}

type CentrifugeGroup struct {
	Name   string
	Events []*CentrifugeEvent
}

type CentrifugeEvent struct {
	Name        string
	Namespace   string
	Channel     string
	Description []string
	Params      interface{}
}

func (g *HTTPGroup) GetName() string {
	return g.Name
}

func (g *HTTPGroup) GetRoutePrefix() string {
	return g.RoutePrefix
}

func (g *HTTPGroup) GetRoutes() []*Route {
	result := make([]*Route, 0, len(g.Routes))
	for _, r := range g.Routes {
		result = append(result, &Route{
			Name:        r.Name,
			Method:      r.Method,
			Path:        r.Path,
			Description: r.Description,
			Request:     r.Request,
			Responses:   r.Responses,
		})
	}
	return result
}

func (g *HTTPGroup) GetCategory() string {
	return "Http"
}

func (g *ConsumedMessagesGroup) GetName() string {
	return g.Name
}

func (g *ConsumedMessagesGroup) GetRoutePrefix() string {
	return "/(geb-in)" + g.RoutePrefix
}

func (g *ConsumedMessagesGroup) GetCategory() string {
	return "Consumed GEB Messages"
}

func (g *ConsumedMessagesGroup) GetRoutes() []*Route {
	result := make([]*Route, 0, len(g.Events))
	for _, e := range g.Events {
		result = append(result, &Route{
			Name:        e.Name,
			Path:        e.EventName,
			Method:      http.MethodPost,
			Description: e.Description,
			Request:     e.Body,
			Responses:   map[int]interface{}{0: nil},
		})
	}

	return result
}

func (g *FiredEventsGroup) GetName() string {
	return g.Name
}

func (g *FiredEventsGroup) GetRoutePrefix() string {
	return "/(geb-out)" + g.RoutePrefix
}

func (g *FiredEventsGroup) GetCategory() string {
	return "Fired GEB Events"
}

func (g *FiredEventsGroup) GetRoutes() []*Route {
	result := make([]*Route, 0, len(g.Events))
	for _, e := range g.Events {
		result = append(result, &Route{
			Name:        e.Name,
			Path:        e.EventName,
			Method:      http.MethodGet,
			Description: e.Description,
			Responses:   map[int]interface{}{0: e.Body},
		})
	}

	return result
}

func (g *CentrifugeGroup) GetName() string {
	return g.Name
}

func (g *CentrifugeGroup) GetRoutePrefix() string {
	return "/(centrifuge)"
}

func (g *CentrifugeGroup) GetCategory() string {
	return "Fired Centrifuge Events"
}

func (g *CentrifugeGroup) GetRoutes() []*Route {
	result := make([]*Route, 0, len(g.Events))
	for _, e := range g.Events {
		result = append(result, &Route{
			Name:        e.Name,
			Path:        e.Namespace + "\\:" + e.Channel,
			Method:      http.MethodPost,
			Description: e.Description,
			Request:     e.Params,
			Responses:   map[int]interface{}{0: nil},
		})
	}

	return result
}
