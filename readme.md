# csp-report-collector

An http applicaton that collects reports generated from browsers for [Content Security Policy violations](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP#enabling_reporting).


## Installation
Copy the`config.example.json` file to `csp-report-collector.json` and fill in [datasource](https://github.com/go-sql-driver/mysql#dsn-data-source-name). The datasource must be a mysql database with a `csp_reports` table.

```bash
# download the prebuilt executable
wget https://github.com/ClearC2/csp-report-collector/releases/download/<release-tag>/csp-report-collector.linux-amd64 

# or clone the repo and build yourself
GOOS=linux GOARCH=amd64 go build -o csp-report-collector.linux-amd64 csp-report-collector.go

# deploy
scp ./csp-report-collector.linux-amd64 user@server:/srv/csp-report-collector/
scp ./csp-report-collector.json user@server:/srv/csp-report-collector/csp-report-collector.json
```

Create a service file on the target server to run the application:

```service
# /etc/systemd/system/csp-report-collector.service
[Unit]
Description=Go CSP report collector
After=network-online.target
[Service]
User=root
Restart=on-failure
ExecStart=/srv/csp-report-collector/csp-report-collector.linux-amd64 /srv/csp-report-collector/csp-report-collector.json
[Install]
WantedBy=multi-user.target
```
Enable and start the service:
```bash
systemctl enable csp-report-collector.service
service csp-report-collector start
```

The csp-report-collector will be running on port 3010.


#### Running locally
Create a local config file first.
```bash
# run locally
go run csp-report-collector.go ./csp-report-collector.json
```
