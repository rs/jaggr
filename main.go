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

var version = "master"

func main() {
	flag.Usage = func() {
		out := os.Stderr
		fmt.Fprintln(out, "Usage: jaggr [OPTIONS] FIELD_DEF [FIELD_DEF...]:")
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, "OPTIONS:")
		flag.PrintDefaults()
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, "FIELD_DEF: <aggr>[,<aggr>...]:path[=alias]")
		fmt.Fprintln(out, "  aggr:")
		fmt.Fprintln(out, "    - min, max, mean: Computes the min, max, mean of the field's values during the sample interval.")
		fmt.Fprintln(out, "    - median, p#: The p1 to p99 compute the percentile of the field's values during the sample interval.")
		fmt.Fprintln(out, "    - sum: Sum all values for the field.")
		fmt.Fprintln(out, "    - [bucket1,bucketN]hist: Count number of values between bucket and bucket+1.")
		fmt.Fprintln(out, "    - [bucket1,bucketN]cat: Count number of values equal to the define buckets (can be non-number values). The special `*` matches values that fit in none of the defined buckets.")
		fmt.Fprintln(out, "  path:")
		fmt.Fprintln(out, "    JSON field path (eg: field.sub-field).")
		fmt.Fprintln(out, "  alias:")
		fmt.Fprintln(out, "    Optional name to use instead of the field path on the output.")
	}
	showVersion := flag.Bool("version", false, "Show version")
	interval := flag.Duration("interval", time.Second, "Sampling interval")
	flag.Parse()

	if *showVersion {
		println(version)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	fields, err := aggr.NewFields(flag.Args())
	if err != nil {
		fatal("invalid argument: ", err)
	}

	go func() {
		t := time.NewTicker(*interval)
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
