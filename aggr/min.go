package aggr

type min struct {
	c uint32
	v float64
}

func newMin(args []string) (Aggregator, error) {
	return &min{}, nil
}

func (a *min) Push(v interface{}) error {
	f, err := parseFloat(v)
	if err != nil {
		return err
	}
	if a.c == 0 || a.v > f {
		a.v = f
	}
	a.c++
	return nil
}

func (a *min) Aggr() interface{} {
	v := a.v
	*a = min{}
	return v
}
