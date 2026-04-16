reup-la:
	docker compose down local-api
	docker compose build local-api
	docker compose up local-api 

reup-pa:
	docker compose down public-api
	docker compose build public-api 
	docker compose up public-api

reup-pf:
	docker compose down public-frontend
	docker compose build public-frontend
	docker compose up public-frontend

reup-lf:
	docker compose down local-frontend
	docker compose build local-frontend
	docker compose up local-frontend

reup-se:
	docker compose --file docker-compose.services.yml down
	docker compose --file docker-compose.services.yml build
	docker compose --file docker-compose.services.yml up -d

up-all:
	docker compose --file ./docker-compose.databases.yml up -d
	docker compose up -d
	docker compose --file docker-compose.services.yml up -d

down-all:
	docker compose --file docker-compose.services.yml down
	docker compose down
	docker compose --file ./docker-compose.databases.yml down

build-all:
	docker compose --file docker-compose.databases.yml build
	docker compose --file docker-compose.services.yml build 
	docker compose build

reup-all:
	make down-all
	make build-all
	make up-all