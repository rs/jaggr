package aggr

type max struct {
	c uint32
	v float64
}

func newMax(args []string) (Aggregator, error) {
	return &max{}, nil
}

func (a *max) Push(v interface{}) error {
	f, err := parseFloat(v)
	if err != nil {
		return err
	}
	if a.c == 0 || a.v < f {
		a.v = f
	}
	a.c++
	return nil
}

func (a *max) Aggr() interface{} {
	v := a.v
	*a = max{}
	return v
}
