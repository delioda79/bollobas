help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@echo "  serve           to serve the app."
	@echo "  stop            to stop the app."
	@echo "  lint            to perform linting."
	@echo "  test            to perform testing."
	@echo "  coverage        to perform coverage report."
	@echo "  ci              to run the tests on ci pipeline."
	@echo "  ci-cleanup      to kill & remove all ci containers."

serve:
	docker-compose -p bollobas -f infra/deploy/local/docker-compose.yml up -d

stop:
	docker-compose -p bollobas -f infra/deploy/local/docker-compose.yml down

lint:
	go fmt ./...
	golint `go list ./...`

test:
	go test ./... -mod=vendor -cover -race -timeout 60s

coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

ci:
	docker-compose -p bollobas -f infra/deploy/local/docker-compose.yml down
	docker-compose -p bollobas -f infra/deploy/local/docker-compose.yml build bollobas_ci
	docker-compose -p bollobas -f infra/deploy/local/docker-compose.yml run bollobas_ci ./script/sql/exec_migrations.sh
	docker-compose -p bollobas -f infra/deploy/local/docker-compose.yml run bollobas_ci ./script/ci.sh
	docker-compose -p bollobas -f infra/deploy/local/docker-compose.yml down

ci-cleanup:
	docker-compose -p bollobas -f infra/deploy/local/docker-compose.yml down
