.PHONY:  migrateup

migrateup:
	migrate -path db/migrations -database "sqlite3://gossip.db" -verbose up

migratedown:
	migrate -path db/migrations -database "sqlite3://gossip.db" -verbose down
