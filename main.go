package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/elgs/gojq"
	"github.com/rs/jaggr/aggr"
)

func main() {
	flag.Parse()
	fields, err := aggr.NewFields(flag.Args())
	if err != nil {
		fatal("invalid argument: ", err)
	}

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			b, err := json.Marshal(fields.Aggr())
			if err != nil {
				fatal("JSON marshal error: ", err)
			}
			fmt.Fprintln(os.Stdout, string(b))
		}
	}()

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		jq, err := gojq.NewStringQuery(s.Text())
		if err != nil {
			fatal("invalid input: ", err)
		}
		if err := fields.Push(jq); err != nil {
			fatal("aggregation error: ", err)
		}
	}
}

func fatal(a ...interface{}) {
	fmt.Println(append([]interface{}{"jaggr: "}, a...)...)
	os.Exit(1)
}
