# kratos performance platform

## TODO Plan

1. Run test using simple js script as source from command line -- currently being developed
2. Start server and run test from js source (either s3 file or string from http req)
3. Create master that will sync all nodes, gather metrics at the end of test
4. Gather metrics at interval from tests


## Cli

kratos has 3 basic functionalities:

- run standalone test from localhost
- start server in master mode which controls slaves to start tests and collect metrics
- start server in slave mode which connects to master and runs tests


## Development

Command to run
```
go run main.go attack -test-file=./example/rp.js -duration=60 -ramp-up-time=20 -max-rps=50
```

## TODO

- pacer vs users ramup
- targeter creation
- default options from command line


echo 'GET https://euc1-resume-points.playback.dazn-dev.com/health' | \
    vegeta attack -rate 50 -duration 30s | \
    vegeta encode | \
    jaggr @count=rps hist\[100,200,300,400,500\]:code p25,p50,p95:latency sum:bytes_in sum:bytes_out | \
    jplot rps+code.hist.100+code.hist.200+code.hist.300+code.hist.400+code.hist.500 latency.p95+latency.p50+latency.p25 bytes_in.sum+bytes_out.sum