include help.mk

### Manage Dependencies

## Install dependencies
deps.install:
	# install golanglint-ci into ./bin
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.41.1
	# Download packages in go.mod file
	go mod download
	# Install NodeJS dependences
	npm ci
	touch ./node_modules/go.mod
.PHONY: deps.install

## Update dependencies
deps.update:
	# Update go dependencies in go.mod and go.sum
	go get -u ./cmd/...
	go get -u ./internal/...
	go mod tidy
	# Update package.json dependencies and lockfile
	npm update;
	touch ./node_modules/go.mod
.PHONY: deps.update

### Code verification and static analysis

## Run code verification
verify:
	# Lint go files
	./bin/golangci-lint --version
	./bin/golangci-lint run ./...
	# Lint OpenAPI spec
	./node_modules/.bin/spectral --version
	./node_modules/.bin/spectral lint --fail-severity=warn ./openapi.yml
.PHONY: verify

## Run code verification and autofix issues where possible
verify.fix:
	echo 'TODO: add go code verification autofixing where possible'
.PHONY: verify.fix

### Testing

## Run tests
test: test.unit
.PHONY: test

## Run tests and output reports
test.report: test.unit.report
.PHONY: test.report

## Run unit tests
test.unit:
	go run github.com/joho/godotenv/cmd/godotenv@v1.4.0 -f .env.test go test -count=1 -v -p=1  ./internal/... ./cmd/...
.PHONY: test.unit

## Run unit tests and output reports
test.unit.report:
	mkdir -p reports
	go run github.com/joho/godotenv/cmd/godotenv@v1.4.0 -f .env.test go test -json -count=1 -coverprofile=reports/test-unit.out -v -p 5 ./internal/... ./cmd/... > reports/test-unit.json
.PHONY: test.unit.report

### Devstack

## Start the devstack
devstack.start:
	docker-compose up -d --remove-orphans devstack
.PHONY: devstack.start

## Stop the devstack
devstack.stop:
	docker-compose down --remove-orphans
.PHONY: devstack.stop

## Clean/reset the devstack
devstack.clean:
	docker-compose down --remove-orphans --volumes
	docker-compose rm --stop --force
.PHONY: devstack.clean

## Restart the devstack
devstack.restart: devstack.stop devstack.start
.PHONY: devstack.restart

## Clean/reset and restart the devstack
devstack.recreate: devstack.clean devstack.start
.PHONY: devstack.recreate

### Dev


## Start local development server
dev:
	./node_modules/.bin/cdk synth --no-staging 'Dev/*' > template.yaml
	sam local start-api
.PHONY: dev

### Deployment

## Deploy dev stage
deploy.dev:
	./node_modules/.bin/cdk deploy --app=cdk.out 'Dev/*'
.PHONY: deploy.dev
