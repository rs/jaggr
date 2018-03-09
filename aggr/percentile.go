package aggr

import (
	"math"
	"sort"
)

type percentile struct {
	p float64
	v []float64
}

func (a *percentile) Push(v interface{}) error {
	f, err := parseFloat(v)
	if err != nil {
		return err
	}
	a.v = append(a.v, f)
	return nil
}

func (a *percentile) Aggr() interface{} {
	defer func() {
		a.v = a.v[:0]
	}()
	l := len(a.v)
	switch l {
	case 0:
		return 0
	case 1:
		return a.v[0]
	}
	var v float64
	sort.Float64s(a.v)
	switch a.p {
	case 0:
		return a.v[0]
	case 100:
		return a.v[l-1]
	}
	i := math.Ceil((a.p/100)*float64(len(a.v)) - 1)
	v = a.v[int(i)]
	a.v = a.v[:0]
	return v
}
