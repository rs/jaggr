package aggr

import (
	"sort"
)

type bucket struct {
	name string
	val  float64
}

type buckets []bucket

func (b buckets) Len() int           { return len(b) }
func (b buckets) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b buckets) Less(i, j int) bool { return b[i].val < b[j].val }

type hist struct {
	buckets buckets
	counts  []uint32
}

func newHist(args []string) (Aggregator, error) {
	h := hist{
		buckets: make([]bucket, 0, len(args)),
		counts:  make([]uint32, len(args)),
	}
	for _, arg := range args {
		val, err := parseFloat(arg)
		if err != nil {
			return nil, err
		}
		h.buckets = append(h.buckets, bucket{name: arg, val: val})
	}
	sort.Sort(h.buckets)
	return h, nil
}

func (a hist) Push(v interface{}) error {
	f, err := parseFloat(v)
	if err != nil {
		return err
	}
	var i int
	for ; i < len(a.buckets)-1; i++ {
		if f >= a.buckets[i].val && f < a.buckets[i+1].val {
			break
		}
	}
	a.counts[i]++
	return nil
}

func (a hist) Aggr() interface{} {
	v := make(map[string]uint32, len(a.buckets))
	for i, b := range a.buckets {
		v[b.name] = a.counts[i]
		a.counts[i] = 0
	}
	return v
}
