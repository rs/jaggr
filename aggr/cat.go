package aggr

import (
	"fmt"
)

type cat map[string]uint32

func newCat(args []string) (Aggregator, error) {
	h := cat{}
	for _, bucket := range args {
		h[bucket] = 0
	}
	return h, nil
}

func (a cat) Push(v interface{}) error {
	bucket := fmt.Sprint(v)
	if _, found := a[bucket]; found {
		a[bucket]++
	} else if _, found := a["*"]; found {
		a["*"]++
	}
	return nil
}

func (a cat) Aggr() interface{} {
	v := make(map[string]uint32, len(a))
	for b, c := range a {
		v[b] = c
		a[b] = 0
	}
	return v
}
