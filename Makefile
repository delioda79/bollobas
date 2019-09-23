help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@echo "  serve           to serve the app."
	@echo "  lint            to perform linting."
	@echo "  test            to perform testing."
	@echo "  coverage        to perform coverage report."
	@echo "  ci              to run the tests on ci pipeline."
	@echo "  ci-cleanup      to kill & remove all ci containers."

serve:
	docker-compose up -d

stop:
	docker-compose down

lint:
	go fmt ./...
	golint `go list ./...`
test:
	go test ./...
coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
ci:
	docker-compose down
	docker-compose up -d --build
	docker-compose run bollobas-dev sh ./scripts/lint.sh
	docker-compose run bollobas-dev go test ./...
	docker-compose down

ci-cleanup:
	docker-compose down