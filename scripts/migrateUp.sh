#!/bin/bash

if [ -f .env ]; then
    source .env
fi


cd repositories/sql/schema
goose turso $DATABASE_URL up