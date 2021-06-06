#!/usr/bin/env bash
set -u

for i in $(seq 1 10); do
  if [[ "$i" -gt "1" ]]; then
    echo "MySQL is unavailable - sleeping ..."
    sleep 5
  fi
  if mysql --host="${WAIT_DB_HOST}" --user="${WAIT_DB_USER}" --password="${WAIT_DB_PASSWORD}" --execute="select 1 from dual;" >/dev/null 2>&1; then
    echo "MySQL connection established."
    exit 0
  fi
done
echo "MySQL not found!"
exit 1
