package aggr

type sum struct {
	v float64
}

func newSum(args []string) (Aggregator, error) {
	return &sum{}, nil
}

func (a *sum) Push(v interface{}) error {
	f, err := parseFloat(v)
	if err != nil {
		return err
	}
	a.v += f
	return nil
}

func (a *sum) Aggr() interface{} {
	v := a.v
	a.v = 0
	return v
}
