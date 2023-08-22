#!/bin/sh
echo "initializing service directory..."
export HYPER_SETTINGS=/app/settings
/app/hyper sd submit-records --reset /directory/001_certificates.json
/app/hyper sd submit-records /directory/quic/001_default.json
echo "Done!"