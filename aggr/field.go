package aggr

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/elgs/gojq"
)

// Field describe an aggregation on a field.
type Field struct {
	Path  string
	Name  string
	Aggrs map[string]Aggregator
}

// Fields holds a set of aggregation fields
type Fields struct {
	f  []Field
	mu sync.Mutex
}

// NewFields parses defs and create aggregation fields.
func NewFields(defs []string) (*Fields, error) {
	fields := make([]Field, 0, len(defs))
	for _, def := range defs {
		f, err := NewField(def)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", def, err)
		}
		fields = append(fields, f)
	}
	return &Fields{f: fields}, nil
}

// Push pushes new pre-parsed JSON data to the aggregations.
func (fs *Fields) Push(jq *gojq.JQ) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	for _, f := range fs.f {
		if err := f.Push(jq); err != nil {
			return err
		}
	}
	return nil
}

// Aggr gets and flush aggregated data.
func (fs *Fields) Aggr() map[string]interface{} {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	v := map[string]interface{}{}
	for _, f := range fs.f {
		v[f.Name] = f.Aggr()
	}
	return v
}

// NewField parses a field definition.
func NewField(def string) (Field, error) {
	if strings.HasPrefix(def, "@count") {
		name := def
		if idx := strings.LastIndexByte(name, '='); idx != -1 {
			name = name[idx+1:]
		}
		return Field{
			Path:  ".",
			Name:  name,
			Aggrs: map[string]Aggregator{"": &count{}},
		}, nil
	}
	idx := strings.LastIndexByte(def, ':')
	if idx == -1 {
		return Field{}, errors.New("missing aggregation definition")
	}
	path := def[idx+1:]
	p := &aggrsParser{exp: def[:idx]}
	aggrs, err := p.parse()
	if err != nil {
		return Field{}, err
	}
	f := Field{
		Path:  path,
		Name:  path,
		Aggrs: aggrs,
	}
	if idx = strings.LastIndexByte(f.Path, '='); idx != -1 {
		f.Name = f.Path[idx+1:]
		f.Path = f.Path[:idx]
	}
	return f, nil
}

// Push pushes new pre-parsed JSON data to the aggregations.
func (f *Field) Push(jq *gojq.JQ) error {
	v, err := jq.Query(f.Path)
	if err != nil {
		return err
	}
	for name, aggr := range f.Aggrs {
		if err := aggr.Push(v); err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
	}
	return nil
}

// Aggr gets and flush aggregated data.
func (f *Field) Aggr() interface{} {
	if f.Path == "." && f.Aggrs[""] != nil {
		// Count special field
		return f.Aggrs[""].Aggr()
	}
	v := map[string]interface{}{}
	for name, aggr := range f.Aggrs {
		v[name] = aggr.Aggr()
	}
	return v
}
