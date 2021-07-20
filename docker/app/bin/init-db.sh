#!/usr/bin/env bash
echo 'Runing migrations...'
/go-backend/migrate up

echo 'Deleting mysql-client...'
apk del mysql-client

echo 'Start application...'
/go-backend/app