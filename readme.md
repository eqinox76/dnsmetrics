# usage

```
go build
./dnsmetric -dnstargets "ip1, ip2, ip3"

```

Example output

```
curl -s localhost:8080/metrics
# HELP dnsquery dns query time in s
# TYPE dnsquery summary
dnsquery{success="true",target="8.8.8.8",quantile="0.5"} 0.011894335
dnsquery{success="true",target="8.8.8.8",quantile="0.9"} 0.019869637
dnsquery{success="true",target="8.8.8.8",quantile="0.99"} 0.020876582
dnsquery_sum{success="true",target="8.8.8.8"} 0.303105881
dnsquery_count{success="true",target="8.8.8.8"} 22
# HELP promhttp_metric_handler_errors_total Total number of internal errors encountered by the promhttp metric handler.
# TYPE promhttp_metric_handler_errors_total counter
promhttp_metric_handler_errors_total{cause="encoding"} 0
promhttp_metric_handler_errors_total{cause="gathering"} 0
```
