docker-setup:
	docker network ls | grep -q app-net || docker network create app-net