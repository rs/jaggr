package aggr

type count struct {
	c uint32
}

func (a *count) Push(v interface{}) error {
	a.c++
	return nil
}

func (a *count) Aggr() interface{} {
	v := a.c
	a.c = 0
	return v
}
