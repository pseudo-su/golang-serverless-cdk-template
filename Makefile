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

.PHONY: codegen
codegen.diagram:
	npx cfn-dia draw.io --output-file docs/cdk-diagram.drawio --ci-mode
	/Applications/draw.io.app/Contents/MacOS/draw.io --export -o docs/cdk-diagram.png docs/cdk-diagram.drawio

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

.PHONY: dev
dev:
	sam-beta-cdk local start-api

.PHONY: package
package:
	npx cdk synth

.PHONY: deploy.dev
deploy.dev:
	npx cdk deploy --app=cdk.out 'Dev/*'

