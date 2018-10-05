# samson-job
A simple tool written in Go to kick off, monitor, and report on a single [Samson](https://github.com/zendesk/samson)
deploy.

## Usage
Required Environment Variables:
* `SAMSON_PROJECT` (id or permalink)
* `REFERENCE` (e.g. 'master', 'v123', '9e44cb0fe')
* `SAMSON_STAGE` (id or permalink)
* `SAMSON_TOKEN`
* `SAMSON_URL` (e.g. 'samson.cooldomain.org')

Optional Environment Variables:
* `POLL_INTERVAL` (in seconds, default: 30)
* `DEPLOY_TIMEOUT` (in minutes, default: 120)

`go run main.go`