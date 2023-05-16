#!/usr/bin/env bash

echo "=========== START LOCAL HTTP / HTTPS TUNNEL AND UPDATE TELEGRAM WEBHOOK TO USE IT ==========="
printf "\n"

# Set environment variables
source ./.env

# Start NGROK in background
echo "Starting ngrok"
ngrok http $PORT > /dev/null &

# Wait for ngrok to be available
while ! nc -z localhost 4040; do
  sleep 0.2
done

NGROK_REMOTE_URL="$(curl -s http://localhost:4040/api/tunnels | jq ".tunnels[0].public_url")"
if test -z "${NGROK_REMOTE_URL}"
then
  echo "ERROR: ngrok doesn't seem to return a valid URL (${NGROK_REMOTE_URL})."
  exit 1
fi

# Trim double quotes from variable
NGROK_REMOTE_URL=$(echo ${NGROK_REMOTE_URL} | tr -d '"')

WEBHOOK_URL=$NGROK_REMOTE_URL/telegram-update
STATUSCODE=$(curl --silent --output /dev/null --write-out "%{http_code}" https://api.telegram.org/bot$TELEGRAM_API_KEY/setWebhook?url=$WEBHOOK_URL)

if test $STATUSCODE -ne 200; then
  echo "ERROR: cannot set webhook: status code: $STATUSCODE"
  exit 1
fi

echo "Webhook set successfully to $WEBHOOK_URL"