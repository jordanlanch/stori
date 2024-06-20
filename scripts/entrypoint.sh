#!/bin/sh

# Ejecutar migraciones
goose -dir /go/src/github.com/jordanlanch/stori-test/storage/migrations postgres "host=$DB_HOST user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME port=$DB_PORT sslmode=disable" up

# Iniciar la aplicación
exec wgo run -buildvcs=false main.go
