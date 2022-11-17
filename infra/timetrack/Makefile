up:
	docker-compose up

down:
	- docker-compose down

clear: down
	- docker volume rm timetrack_grafana_data
	- docker volume rm timetrack_mysql_data
	- docker volume rm timetrack_prometheus_data