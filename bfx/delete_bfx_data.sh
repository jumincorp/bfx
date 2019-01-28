#!/usr/bin/env bash
PROMETHEUS_URL="http://localhost:9090/prometheus"

curl -w "\nReturn Code: %{http_code}\n" -X POST -d "{
  \"matchers\": [{
  \"type\": \"EQ\",
  \"name\": \"namespace\",
  \"value\": \"bfx\"
  }]
}" -H "Content-Type: application/json" "${PROMETHEUS_URL}"/api/v2/admin/tsdb/delete_series
