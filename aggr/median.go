package aggr

import "sort"

type median struct {
	v []float64
}

func newMedian(args []string) (Aggregator, error) {
	return &median{}, nil
}

func (a *median) Push(v interface{}) error {
	f, err := parseFloat(v)
	if err != nil {
		return err
	}
	a.v = append(a.v, f)
	return nil
}

func (a *median) Aggr() interface{} {
	var v float64
	if len(a.v) == 1 {
		v = a.v[0]
	} else if len(a.v) > 1 {
		sort.Float64s(a.v)
		l := len(a.v)
		if l%2 == 0 {
			v = (a.v[l/2-1] + a.v[l/2]) / 2
		} else {
			v = a.v[l/2]
		}
	}
	a.v = a.v[:0]
	return v
}
