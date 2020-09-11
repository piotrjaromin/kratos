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

```bash
go run main.go attack -test-file=./mock-server/mock-attack.js -duration=60 -ramp-up-time=20 -max-rps=50
```

## TODO

- targeter creation
- default options from command line
