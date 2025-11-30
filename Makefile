.PHONY: run stop clean logs

run:
	docker-compose up -d --build

stop:
	docker-compose down

clean:
	docker-compose down -v

logs:
	docker-compose logs -f

tidy:
	go mod tidy
	go fmt ./...