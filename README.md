# net-delay-time-exporter
Service that ping remote server and collect delay time metric for Prometheus consumption

## Install

```sh
$ go get github.com/lonord/net-delay-time-exporter
```

## Usage

Start exporter server:

```sh
$ net-delay-time-exporter your.domain1.com your.domain2.com 123.456.789.100
```

Check the metric data by `curl`:

```sh
$ curl http://127.0.0.1:8080/metrics
```

Output data:

```sh
# ...
# HELP server_delay_time_ms Round-trip delay time in milliseconds
# TYPE server_delay_time_ms gauge
server_delay_time_ms{server="your.domain1.com"} 11
server_delay_time_ms{server="your.domain2.com"} 53
server_delay_time_ms{server="123.456.789.100"} 210
# HELP server_package_lost_rate The package loss rate during ping
# TYPE server_package_lost_rate gauge
server_package_lost_rate{server="your.domain1.com"} 0
server_package_lost_rate{server="your.domain2.com"} 10
server_package_lost_rate{server="123.456.789.100"} 50
```

### License

MIT