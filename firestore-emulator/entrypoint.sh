#!/bin/sh

exec firebase emulators:start --only firestore,ui --import=./data --export-on-exit --project="${CLOUDSDK_CORE_PROJECT}"
