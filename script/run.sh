#!/bin/sh
set -e

# migrate
echo "Running migrations on ${MYSQL_USERNAME}"
./migrate -database mysql://${MYSQL_USERNAME}:${MYSQL_PASS}@tcp\(${MYSQL_WRITE}:${MYSQL_PORT}\)/${MYSQL_DB} -verbose -source file://./migrations up

# run
./bollobas