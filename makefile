include .env
export 

export  PROJECT_ROOT=${shell pwd}
env-up:
	@docker compose up -d todoapp-postgres
env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Очистить все volume файлы Окружения? Опасность утери данных.[y/n]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down  todoapp-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo  "Файлы окружения очищены";\
	else \
		echo "Очиста окружения отменена "; \
	fi 

env-port-forward: 
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ];then \
		echo  "Отсутвует необоходимый параметр seq. Пример: migrate-create seq=<название миграции>"; \
		exit 1; \
		fi; \


	docker compose run --rm todoapp-postgres-migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq  "$(seq)"

migrate-up:
	@make migrate-action action=up


migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ];then \
	echo  "Отсутвует необоходимый параметр action. Пример: migrate-action action=<название версии миграции>"; \
	exit 1; \
	fi; \

	docker compose run --rm  todoapp-postgres-migrate \
	-path /migrations \
	-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable \
	"$(action)"

log-clean:
	@read -p "Очистить все лог файлы? Опасность утери логов.[y/n]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs/* && \
		echo  "Лог файлы очищены";\
	else \
		echo "Очиста лог файлов отменена "; \
	fi
todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go