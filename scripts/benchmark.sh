#!/bin/sh
set -eu
url="${1:-http://localhost:8081}"
# Rotating client IDs: accepted-request throughput (includes new bucket creation).
wrk -t4 -c100 -d30s --latency -s scripts/traffic.lua "$url"
# Fixed client ID: rejection-heavy traffic after the initial burst.
wrk -t4 -c100 -d30s --latency -H 'X-Client-ID: burst-client' "$url"
