#!/bin/sh

CONFIG_PATH=${CONFIG_PATH:-/data}
API_KEY=${API_KEY:-secret}
SERVER_URL=${SERVER_URL:-https://mdm.example.com/}
WEBHOOK_URL=${WEBHOOK_URL:-http://127.0.0.1:5000/webhook}
WEBHOOK_AUTH_USER=${WEBHOOK_AUTH_USER:-}
WEBHOOK_AUTH_PASSWORD=${WEBHOOK_AUTH_PASSWORD:-}
RDBMS=${MICROMDM_RDBMS:-mysql}
RDBMS_USER=${MICROMDM_RDBMS_USER:-micromdm}
RDBMS_PASSWORD=${MICROMDM_RDBMS_PASSWORD:-micromdm}
RDBMS_DATABASE=${MICROMDM_RDBMS_DATABASE:-micromdm}
RDBMS_HOST=${MICROMDM_RDBMS_HOST:-127.0.0.1}
RDBMS_PORT=${MICROMDM_RDBMS_PORT:-3306}
SCEP_CLIENT_VALIDITY=${SCEP_CLIENT_VALIDITY:-10000}

micromdm serve \
  -config-path ${CONFIG_PATH} \
  -api-key ${API_KEY} \
  -server-url ${SERVER_URL} \
  -scep-client-validity ${SCEP_CLIENT_VALIDITY} \
  -command-webhook-url ${WEBHOOK_URL} \
  -command-webhook-auth-user ${WEBHOOK_AUTH_USER} \
  -command-webhook-auth-pass ${WEBHOOK_AUTH_PASSWORD} \
  -rdbms ${RDBMS} \
  -rdbms-username ${RDBMS_USER} \
  -rdbms-password ${RDBMS_PASSWORD} \
  -rdbms-database ${RDBMS_DATABASE} \
  -rdbms-host ${RDBMS_HOST} \
  -rdbms-port ${RDBMS_PORT} \
  -tls=false
