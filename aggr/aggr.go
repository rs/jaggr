package aggr

import (
	"errors"
	"fmt"
	"strconv"
)

var rAggrs = map[string]func(args []string) (Aggregator, error){}

func init() {
	Register("min", newMin)
	Register("max", newMax)
	Register("mean", newMean)
	Register("median", newMedian)
	Register("hist", newHist)
	Register("cat", newCat)
	for i := 1; i < 100; i++ {
		p := float64(i)
		Register(fmt.Sprintf("p%d", i), func(args []string) (Aggregator, error) {
			return &percentile{p: p}, nil
		})
	}
}

// Aggregator accumulates data using Push and aggregate (and reset) the
// accumulated data on Aggr. The result is expected to be used as a JSON value.
// Numbers are expressed as float64.
type Aggregator interface {
	Push(interface{}) error
	Aggr() interface{}
}

// Register a custom aggregator.
func Register(name string, aggr func(args []string) (Aggregator, error)) error {
	if _, found := rAggrs[name]; found {
		return errors.New("already exists")
	}
	rAggrs[name] = aggr
	return nil
}

func parseFloat(v interface{}) (float64, error) {
	return strconv.ParseFloat(fmt.Sprint(v), 64)
}
