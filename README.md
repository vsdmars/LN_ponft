# ponft is a tool to query request to plat-telemetry /li/ponf endpoint

## usage:

ponft random pick post body from request.txt

request.txt is embedded inside the ponft binary, you don't have to prepare request.txt
for running ponft.


### -help for help

$ ponft -h

```help
Usage of ponft:
-minute int
test period (in minutes) (default 1)
-qps int
QPS (default 100)
```

### query with qps 642 run in 10 minutes period

$ ponft -qps 642 -minute 10
