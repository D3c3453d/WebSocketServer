.PHONY: db server down local

db: #run first
	docker compose up -d db
server: #./cfg/db.env: POSTGRES_HOST=db
	docker compose build server && docker compose up server
down:
	docker compose down
local: #./cfg/db.env: POSTGRES_HOST=localhost
	go run ./app/cmd/main/main.go