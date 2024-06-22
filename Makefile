migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path migrations -database $(url) up

migrate-down:
	migrate -path migrations -database $(url) down

compose:
	docker-compose up --remove-orphans --build