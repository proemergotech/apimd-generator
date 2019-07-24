package generator

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	typeBody  = "body"
	typeParam = "param"
	typeQuery = "query"
)

var urlRegex = regexp.MustCompile(`:(\w+)`)

type Collector struct {
	values  map[int]*Value
	markers map[string]*Marker
}

func newCollector() *Collector {
	return &Collector{
		values:  make(map[int]*Value),
		markers: make(map[string]*Marker),
	}
}

type markedRouteKey struct {
	groupIndex int
	markerKey  string
	routeIndex int
}

func (c *Collector) collect(d Definitons) *Document {
	factory := newFactory(c)
	groups := d.Groups(factory)

	markedRoutes := make(map[markedRouteKey]*Route, len(c.markers))
	for markerK := range c.markers {
		factory.markerKey = markerK
		mg := d.Groups(factory)
		for groupI, group := range mg {
			for routeI, route := range group.GetRoutes() {
				markedRoutes[markedRouteKey{groupIndex: groupI, markerKey: markerK, routeIndex: routeI}] = route
			}
		}
	}

	result := &Document{
		Name:       d.Name(),
		Usage:      d.Usage(),
		Categories: make([]*DocCategory, 0),
	}

	for groupI, group := range groups {
		routes := group.GetRoutes()

		docRoutes := make([]*DocRoute, 0, len(routes))
		for routeI, route := range routes {
			path := normalizePath(route.Path)
			docRoute := &DocRoute{
				Name:        route.Name,
				Method:      route.Method,
				Description: route.Description,
				Params:      make(map[string]*DocValue),
				Query:       make(map[string]*DocValue),
				Path:        path,
			}

			if route.Request != nil {
				markedRequests := make(map[string]interface{})
				for mk, mr := range markedRoutes {
					if mk.groupIndex == groupI && mk.routeIndex == routeI {
						markedRequests[mk.markerKey] = mr.Request
					}
				}
				valueTree, err := c.createTree(d, route.Request, markedRequests)
				if err != nil {
					log.Fatalf("%+v", errors.Wrapf(err, "parsing request: [%v] %v", route.Method, route.Path))
				}

				params := c.toMap(c.docValues(valueTree, typeParam))
				for k, p := range params {
					pVal, ok := p.(*DocValue)
					if !ok {
						log.Fatalf("%+v", errors.New("nested object supplied as param for route: "+route.Name))
					}

					docRoute.Params[k] = pVal
				}

				query := c.toMap(c.docValues(valueTree, typeQuery))
				queryKeys := make([]string, 0)
				for k, q := range query {
					qVal, ok := q.(*DocValue)
					if !ok {
						if values, ok := q.([]interface{}); ok && len(values) != 0 {
							for _, value := range values {
								if v, ok := value.(*DocValue); ok {
									if qVal == nil {
										qVal = v
										qVal.APIMDType = "array"
									}	else {
										mergeQueryDocValues(qVal, v)
									}
								}
							}
						}
						if qVal == nil {
							log.Fatalf("%+v", errors.New("invalid nested object supplied as query for route: "+route.Name))
						}
					}

					docRoute.Params[k] = qVal
					queryKeys = append(queryKeys, k)
				}

				if len(queryKeys) > 0 {
					sort.Strings(queryKeys)
					path += "{?" + strings.Join(queryKeys, ",") + "}"
				}

				docRoute.Path = path

				docRoute.RequestBody = c.docValues(valueTree, typeBody)
			}

			docRoute.ResponseBodies = make(map[int]interface{})
			for statusCode, resp := range route.Responses {
				markedResponses := make(map[string]interface{})
				for mk, mr := range markedRoutes {
					if mk.groupIndex == groupI && mk.routeIndex == routeI {
						markedResponses[mk.markerKey] = mr.Responses[statusCode]
					}
				}
				valueTree, err := c.createTree(d, resp, markedResponses)
				if err != nil {
					log.Fatalf("%+v", errors.Wrapf(err, "parsing response: [%v] %v", route.Method, route.Path))
				}

				docRoute.ResponseBodies[statusCode] = c.docValues(valueTree, typeBody)
			}

			docRoutes = append(docRoutes, docRoute)
		}

		docGroup := &DocGroup{
			Name:   group.GetName(),
			Prefix: group.GetRoutePrefix(),
			Routes: docRoutes,
		}

		var category *DocCategory
		for _, cat := range result.Categories {
			if cat.Name == group.GetCategory() {
				category = cat
				break
			}
		}
		if category == nil {
			category = &DocCategory{
				Name:   group.GetCategory(),
				Groups: make([]*DocGroup, 0),
			}
			result.Categories = append(result.Categories, category)
		}
		category.Groups = append(category.Groups, docGroup)
	}

	return result
}

func mergeQueryDocValues(target *DocValue, source *DocValue) {
	target.Opt = target.Opt || source.Opt
	target.Desc = joinNonEmpty(", ", target.Desc, source.Desc)
	target.Value = joinNonEmpty(",", target.Value, source.Value)
}

func joinNonEmpty(sep string, values ...string) string {
	a := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" {
			a = append(a, value)
		}
	}
	return strings.Join(a, sep)
}

func (c *Collector) createTree(d Definitons, data interface{}, markedData map[string]interface{}) (interface{}, error) {
	data2, err := encDec(data)
	if err != nil {
		return nil, err
	}

	valueTree, err := c.createIndexTree("", data2, d)
	if err != nil {
		return nil, err
	}
	for mk, mData := range markedData {
		mData2, err := encDec(mData)
		if err != nil {
			return nil, err
		}

		valueTree, err = c.extendTreeWithMarker(valueTree, mData2, c.markers[mk])
		if err != nil {
			return nil, err
		}
	}

	return valueTree, nil
}

func (c *Collector) createIndexTree(key string, data interface{}, defs Definitons) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	switch d := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(d))
		for k, v := range d {
			val, err := c.createIndexTree(key+"."+k, v, defs)
			if err != nil {
				return nil, err
			}
			if val != nil {
				result[k] = val
			}
		}
		return result, nil

	case []interface{}:
		result := make([]interface{}, 0, len(d))
		for k, v := range d {
			val, err := c.createIndexTree(key+"."+strconv.Itoa(k), v, defs)
			if err != nil {
				return nil, err
			}
			if val != nil {
				result = append(result, val)
			}
		}
		return result, nil

	default:
		apimdType, err := c.apimdType(d)
		if err != nil {
			return nil, err
		}

		for _, m := range c.markers {
			if m.jsonPlaceholder == d {
				return nil, errors.Errorf("marker value %v found while parsing tree, if you want to use this value, use it in the full form: d.body(\"%v\").cast()", d, d)
			}
		}

		ind, err := defs.ParseIndex(d)
		if err != nil {
			return nil, errors.Wrapf(err, "ParseIndex failed for key: %v, value: %v", key, d)
		}

		val := c.values[ind]
		if val != nil {
			val.apimdType = apimdType
			return val, nil
		}

		if isZero(d) {
			// ignore default values, to be able to omit struct fields from documentation
			return nil, nil
		}

		// use as-is
		val = &Value{value: fmt.Sprint(d), typ: typeBody, apimdType: apimdType}

		return val, nil
	}
}

func (c *Collector) extendTreeWithMarker(valueTree interface{}, markedTree interface{}, marker *Marker) (interface{}, error) {
	var err error

	switch mt := markedTree.(type) {
	case map[string]interface{}:
		vt, ok := valueTree.(map[string]interface{})
		if !ok {
			vt = make(map[string]interface{})
		}
		for k, v := range mt {
			vt[k], err = c.extendTreeWithMarker(vt[k], v, marker)
			if err != nil {
				return nil, err
			}
		}

		return vt, nil

	case []interface{}:
		vt, ok := valueTree.([]interface{})
		if !ok {
			vt = make([]interface{}, 0)
		}

		for i, v := range mt {
			vt[i], err = c.extendTreeWithMarker(vt[i], v, marker)
			if err != nil {
				return nil, err
			}
		}

		return vt, nil

	default:
		if marker.jsonPlaceholder != mt {
			return valueTree, nil
		}

		val := marker.Value
		val.apimdType, err = c.apimdType(mt)
		if err != nil {
			return nil, err
		}

		return marker.Value, nil
	}
}

func (c *Collector) docValues(v interface{}, typ string) interface{} {
	if v == nil {
		return nil
	}

	switch vt := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(vt))
		for k, v := range vt {
			val := c.docValues(v, typ)
			if val != nil {
				result[k] = val
			}
		}
		if len(result) == 0 {
			return nil
		}

		return result

	case []interface{}:
		result := make([]interface{}, 0, len(vt))
		for _, v := range vt {
			val := c.docValues(v, typ)
			if val != nil {
				result = append(result, val)
			}
		}
		if len(result) == 0 {
			return nil
		}

		return result

	case *Value:
		if vt.typ != typ {
			return nil
		}

		return vt.docValue()

	default:
		log.Fatalf("%+v", errors.Errorf("invalid type %T in docValues", v))
		return nil
	}
}

func (*Collector) toMap(v interface{}) map[string]interface{} {
	if v == nil {
		return map[string]interface{}{}
	}
	res, ok := v.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return res
}

func (*Collector) apimdType(jsonVal interface{}) (string, error) {
	switch jsonVal.(type) {
	case string:
		return "string", nil
	case bool:
		return "boolean", nil
	case float64:
		return "number", nil
	default:
		return "", errors.Errorf("apimdType unhandled jsonType: %T", jsonVal)
	}
}

func isZero(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// normalizePath see: TestNormalizePath
func normalizePath(path string) string {
	const placeholder = "{colon-esc-placeholr}"
	path = strings.Replace(path, "\\:", placeholder, -1)
	path = urlRegex.ReplaceAllString(path, "{$1}")
	path = strings.Replace(path, placeholder, ":", -1)

	return path
}
