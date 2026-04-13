install:
	go install github.com/vektra/mockery/v2@v2.53.6
	go install github.com/go-task/task/v3/cmd/task@latest

mock:
	go generate ./...

migrateaddinit:
	migrate create -ext sql -dir migrations -seq init

migrateaddname:
	migrate create -ext sql -dir migrations -seq $(name)

migrateup:
	migrate -source "file://migrations" -database "postgres://postgres:123456@localhost:5432/postgres?sslmode=disable" up $(int)

migratedown:
	migrate -source "file://migrations" -database "postgres://postgres:123456@localhost:5432/postgres?sslmode=disable" down $(int)