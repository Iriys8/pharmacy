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

up-all:
	docker compose up -d

down-all:
	docker compose down

build-all:
	docker compose build

reup-all:
	make down-all
	make build-all
	make up-all