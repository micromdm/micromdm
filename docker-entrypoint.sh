#!/bin/sh

CONFIG_PATH=${CONFIG_PATH:-/data}
API_KEY=${API_KEY:-secret}
SERVER_URL=${SERVER_URL:-https://mdm.example.com/}
WEBHOOK_URL=${WEBHOOK_URL:-http://127.0.0.1:5000/webhook}
WEBHOOK_AUTH_USER=${WEBHOOK_AUTH_USER:-}
WEBHOOK_AUTH_PASSWORD=${WEBHOOK_AUTH_PASSWORD:-}
MYSQL_USER=${MYSQL_USER:-micromdm}
MYSQL_PASSWORD=${MYSQL_PASSWORD:-micromdm}
MYSQL_DATABASE=${MYSQL_DATABASE:-micromdm}
MYSQL_HOST=${MYSQL_HOST:-127.0.0.1}
MYSQL_PORT=${MYSQL_PORT:-3306}
SCEP_CLIENT_VALIDITY=${SCEP_CLIENT_VALIDITY:-10000}

micromdm serve \
  -config-path ${CONFIG_PATH} \
  -api-key ${API_KEY} \
  -server-url ${SERVER_URL} \
  -scep-client-validity ${SCEP_CLIENT_VALIDITY} \
  -command-webhook-url ${WEBHOOK_URL} \
  -command-webhook-auth-user ${WEBHOOK_AUTH_USER} \
  -command-webhook-auth-pass ${WEBHOOK_AUTH_PASSWORD} \
  -mysql-username ${MYSQL_USER} \
  -mysql-password ${MYSQL_PASSWORD} \
  -mysql-database ${MYSQL_DATABASE} \
  -mysql-host ${MYSQL_HOST} \
  -mysql-port ${MYSQL_PORT} \
  -tls=false
