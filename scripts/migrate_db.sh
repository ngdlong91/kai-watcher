#!/bin/sh
source .env
migrate -path=repo/migrations/ -database $DB_URL $1 $2