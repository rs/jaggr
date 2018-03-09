package aggr

import "fmt"

type aggrsParser struct {
	exp string
	pos int
}

func (p *aggrsParser) parse() (map[string]Aggregator, error) {
	aggrs := map[string]Aggregator{}
	for p.more() {
		name := p.scanName()
		args := []string{}
		if len(name) == 0 {
			return nil, fmt.Errorf("expect name got `%c' at char %d", p.peekAt(p.pos), p.pos)
		}
		aggr := rAggrs[name]
		if aggr == nil {
			return nil, fmt.Errorf("unknown aggregation: %v", name)
		}
		if p.expect('[') {
			for {
				arg := p.scanArg()
				if len(arg) == 0 {
					return nil, fmt.Errorf("%s: zero length argument", name)
				}
				args = append(args, arg)
				if p.expect(']') {
					break
				}
				if !p.expect(',') {
					return nil, fmt.Errorf("%s: expect `,' got `%c' at char %d", name, p.peekAt(p.pos), p.pos)
				}
			}
		}
		var err error
		aggrs[name], err = aggr(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", name, err)
		}
		if !p.expect(',') {
			break
		}
	}
	if p.more() {
		return nil, fmt.Errorf("unexpected character `%c' at char %d", p.peek(), p.pos)
	}
	return aggrs, nil
}

func (p *aggrsParser) scanName() string {
	start := p.pos
	for p.more() {
		c := p.peek()
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			break
		}
		p.pos++
	}
	return p.exp[start:p.pos]
}

func (p *aggrsParser) scanArg() string {
	start := p.pos
	for p.more() {
		c := p.peek()
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '*') {
			break
		}
		p.pos++
	}
	return p.exp[start:p.pos]
}

// more returns true if there is more data to parse.
func (p *aggrsParser) more() bool {
	return p.pos < len(p.exp)
}

// expect advances the cursor if the current char is equal to c or return
// false otherwise.
func (p *aggrsParser) expect(c byte) bool {
	if p.peek() == c {
		p.pos++
		return true
	}
	return false
}

// peek returns the char at the current position without advancing the cursor.
func (p *aggrsParser) peek() byte {
	if p.more() {
		return p.exp[p.pos]
	}
	return 0
}

// peek returns the char at the given position without moving the cursor.
func (p *aggrsParser) peekAt(pos int) byte {
	if pos < len(p.exp) {
		return p.exp[pos]
	}
	return 0
}
