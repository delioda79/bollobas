#!/bin/sh
# Usage: ./exec_migrate <.env> <(up)/down>

migration_dir=file://./script/sql/migrations

switch=up
if [[ $1 ]]; then
    source $1
fi
if [[ $2 ]]; then
    switch=$2
fi

migrate_cmd() {
    $1migrate -database mysql://${MYSQL_USERNAME}:${MYSQL_PASS}@tcp\(${MYSQL_WRITE}:${MYSQL_PORT}\)/${MYSQL_DB} -verbose -source $migration_dir $switch
}

migrate_db() {
    # Attempt to run migrations, retry if the database container is not yet ready.
    echo "Running migrations on ${MYSQL_DB}"
    i=0
    migrate_cmd ./
    while [ $? -ne 0 -a $i -lt 10 ]; do
        echo "Database not ready (attempt #$i), retrying.."
        sleep 2
        i=`expr $i + 1`
    migrate_cmd ./
    done
}

check_success() {
    # Exit if the last command failed.
    if [ $? -ne 0 ]; then
        echo "Last command failed, exiting.."
        exit 1
    fi
}

# Run default migration.
if [[ ${MYSQL_DB} =~ "test" ]]; then
    migrate_cmd
else
    migrate_db
fi
check_success

echo "Migration succeeded."
