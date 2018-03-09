package aggr

type mean struct {
	c float64
	v float64
}

func newMean(args []string) (Aggregator, error) {
	return &mean{}, nil
}

func (a *mean) Push(v interface{}) error {
	f, err := parseFloat(v)
	if err != nil {
		return err
	}
	a.v += f
	a.c++
	return nil
}

func (a *mean) Aggr() interface{} {
	v := a.v / a.c
	*a = mean{}
	return v
}
