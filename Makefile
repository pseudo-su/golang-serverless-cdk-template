.PHONY: deps.install
deps.install:
	# install golanglint-ci into ./bin
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.41.1
	# Download packages in go.mod file
	go mod download
	# Install NodeJS dependences
	npm ci

.PHONY: deps.update
deps.update:
	# Update go dependencies in go.mod and go.sum
	go get -u ./lambda/...
	go get -u ./internal/...
	go mod tidy
	# Update package.json dependencies and lockfile
	npm update;

.PHONY: verify
verify:
	# Lint go files
	./bin/golangci-lint run ./...
	# Lint OpenAPI spec
	./node_modules/.bin/spectral lint --fail-severity=warn ./openapi.yml

.PHONY: verify.fix
verify.fix:
	echo 'TODO: add go code verification autofixing where possible'

.PHONY: test
test:
	npx dotenv -c -- go test -v -p=1 ./lambda/...
	npx dotenv -c -- go test -v -p=1 ./internal/...

.PHONY: devstack.start
devstack.start:
	docker-compose up -d --remove-orphans devstack

.PHONY: devstack.stop
devstack.stop:
	docker-compose down --remove-orphans

.PHONY: devstack.clean
devstack.clean:
	rm -rf ./devstack/postgres/data/*
	docker-compose rm --stop --force

.PHONY: devstack.restart
devstack.restart: devstack.stop devstack.start

.PHONY: devstack.recreate
devstack.recreate: devstack.clean devstack.restart

.PHONY: dev.start
dev.start:
	if [ -z "${STAGE}" ]; then echo "STAGE environment variable not set"; exit 1; fi
	npx sst start --stage ${STAGE}

.PHONY: dev.stop
dev.stop:
	if [ -z "${STAGE}" ]; then echo "STAGE environment variable not set"; exit 1; fi
	npx sst remove --stage ${STAGE}

.PHONY: dev.clean
dev.clean:
	rm -rf .build/

.PHONY: dev.restart
dev.restart: dev.stop dev.start

.PHONY: dev.recreate
dev.recreate: dev.clean dev.restart
