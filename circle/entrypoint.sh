#!/bin/sh
export GOOGLE_APPLICATION_CREDENTIALS=/root/google/service_account.json
echo "$1" > ${GOOGLE_APPLICATION_CREDENTIALS}
chmod 600 ${GOOGLE_APPLICATION_CREDENTIALS}
app