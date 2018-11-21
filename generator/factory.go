package generator

import (
	"fmt"
	"strconv"
)

type Factory struct {
	counter   int
	markerKey string
	c         *Collector
}

type Marker struct {
	Value           *Value
	jsonPlaceholder interface{}
}

type Value struct {
	index     int
	value     string
	desc      string
	typ       string
	opt       bool
	apimdType string
	factory   *Factory
}

func (v *Value) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"index: %v value: %v type: %v\"", v.index, v.value, v.typ)), nil
}

func newFactory(c *Collector) *Factory {
	return &Factory{
		counter: 8103623,
		c:       c,
	}
}

func (f *Factory) Param(val string) *Value {
	return f.newValue(val, typeParam)
}

func (f *Factory) Query(val string) *Value {
	return f.newValue(val, typeQuery)
}

func (f *Factory) Body(val string) *Value {
	return f.newValue(val, typeBody)
}

func (f *Factory) Marker(v *Value, jsonPlaceholder interface{}) bool {
	key := fmt.Sprintf("%s%s%v", v.value, v.desc, v.opt)

	if f.markerKey != "" {
		if f.markerKey == key {
			return true
		}
	} else {
		f.c.markers[key] = &Marker{
			Value:           v,
			jsonPlaceholder: jsonPlaceholder,
		}
	}

	return false
}

func (f *Factory) newValue(val string, typ string) *Value {
	f.counter++
	v := &Value{
		index:   f.counter,
		value:   val,
		typ:     typ,
		factory: f,
	}
	f.c.values[f.counter] = v

	return v
}

func (v *Value) docValue() *DocValue {
	return &DocValue{
		Value:     v.value,
		Desc:      v.desc,
		Opt:       v.opt,
		APIMDType: v.apimdType,
	}
}

func (v *Value) Index() int {
	return v.index
}

func (v *Value) Marker(jsonPlaceholder interface{}) bool {
	return v.factory.Marker(v, jsonPlaceholder)
}

func (v *Value) Description(d string) {
	v.desc = d
}

func (v *Value) Optional() {
	v.opt = true
}

func (v *Value) String() string {
	return strconv.Itoa(v.Index())
}

func (v *Value) StringPtr() *string {
	str := strconv.Itoa(v.Index())
	return &str
}

func (v *Value) Int() int {
	return v.Index()
}

func (v *Value) IntPtr() *int {
	i := v.Int()
	return &i
}

func (v *Value) Uint() uint {
	return uint(v.Index())
}

func (v *Value) UintPtr() *uint {
	i := v.Uint()
	return &i
}

func (v *Value) Int32() int32 {
	return int32(v.Index())
}

func (v *Value) Int32Ptr() *int32 {
	i := v.Int32()
	return &i
}

func (v *Value) Int64() int64 {
	return int64(v.Index())
}

func (v *Value) Int64Ptr() *int64 {
	i := v.Int64()
	return &i
}

func (v *Value) Uint32() uint32 {
	return uint32(v.Index())
}

func (v *Value) Uint32Ptr() *uint32 {
	i := v.Uint32()
	return &i
}

func (v *Value) Uint64() uint64 {
	return uint64(v.Index())
}

func (v *Value) Uint64Ptr() *uint64 {
	i := v.Uint64()
	return &i
}

func (v *Value) Float32() float32 {
	return float32(v.Index())
}

func (v *Value) Float32Ptr() *float32 {
	f := v.Float32()
	return &f
}

func (v *Value) Float64() float64 {
	return float64(v.Index())
}

func (v *Value) Float64Ptr() *float64 {
	f := v.Float64()
	return &f
}

func (v *Value) Byte() byte {
	isMarker := v.Marker(float64(173))

	if isMarker {
		return 173
	}

	return 0
}

func (v *Value) BytePtr() *byte {
	b := v.Byte()
	return &b
}

func (v *Value) Uint16() uint16 {
	isMarker := v.Marker(float64(173))

	if isMarker {
		return 173
	}

	return 0
}

func (v *Value) Uint16Ptr() *uint16 {
	i := v.Uint16()
	return &i
}

func (v *Value) Bool() bool {
	return v.Marker(true)
}

func (v *Value) BoolPtr() *bool {
	isMarker := v.Marker(true)

	return &isMarker
}
