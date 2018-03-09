# jaggr: JSON Aggregation CLI
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/jaggr/master/LICENSE) [![Build Status](https://travis-ci.org/rs/jaggr.svg?branch=master)](https://travis-ci.org/rs/jaggr)

Jaggr is a command line tool to aggregate in real time a series of JSON logs. The main goal of this tool is to prepare data for plotting with [jplot](https://github.com/rs/jplot).

## Install

```
go get -u github.com/rs/jaggr
```

## Usage

Given the input below, generate one line per second with mean, min, max:

```
{"code": 200, "latency": 4788000, "error": ""}
{"code": 200, "latency": 5785000, "error": ""}
{"code": 200, "latency": 4162000, "error": ""}
{"code": 502, "latency": 4461000, "error": "i/o error"}
{"code": 200, "latency": 5884000, "error": ""}
{"code": 200, "latency": 4702000, "error": ""}
...
```

```
tail -f log.json | jaggr @count=rps hist[200,502,*]:code min,max,mean:latency
```

Output will be on line per second as follow:

```
{"count":123, "code": {"hist": {"200": 100, 502: 13, "*": 10}}, "latency":{"min": 4461000, "max": 5884000, "mean": 4483000}}
```

So here we give a stream of real-time requests to jaggr standard input and request the aggregation of the `code` and `latency` fields. For the `code` we request an histogram with some known error codes with an "other" bucket defined by `*`. The `latency` field is aggregated using minimum, maximum and mean. In addition, `@count` adds an extra field indicating the total number of lines aggregated. The `=` sign can be used on any field to rename it, here we use it to say that the count is an `rps` as we are using the default aggregation time of 1 second.

Note that any field not specified in the argument list are removed from the output (i.e. `error` field).

## Recipes

### Vegeta

Jaggr can be used to integrate [vegeta](https://github.com/tsenart/vegeta) with [jplot](https://github.com/rs/jplot) as follow:

```
echo 'GET http://localhost' | \
    vegeta attack -rate 5000 -workers 100 -duration 10m | vegeta dump | \
    jaggr @count=rps \
          hist[100,200,300,400,500]:code \
          p25,p50,p95:latency \
          p25,p50,p95:bytes_in \
          p25,p50,p95:bytes_out | \
    jplot rps+code.hist.100+code.hist.200+code.hist.300+code.hist.400+code.hist.500 \
          latency.p95+latency.p50+latency.p25 \
          bytes_in.p95+bytes_in.p50+bytes_in.p25 \
          bytes_out.p95+bytes_out.p50+bytes_out.p25
```

![](doc/vegeta.gif)