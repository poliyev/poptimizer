VENV_NAME=.venv
PYTHON=${VENV_NAME}/bin/python3

new:
	brew install go
	brew install golangci-lint
	brew install python
	brew install mongodb-community
	brew install mongodb-database-tools
	make venv

venv:
	rm -rf $(VENV_NAME)
	python3 -m venv $(VENV_NAME)
	make update

update:
	brew upgrade python
	brew upgrade mongodb-community
	brew upgrade mongodb-database-tools
	${PYTHON} -m pip install -U pip
	${PYTHON} -m pip install -U -r requirements.txt

update_go:
	brew upgrade go
	brew upgrade golangci-lint
	cd data;go get -u -t -v ./...;go mod tidy

db_recover:
	@echo "Recover MongoDB"
	mongod --dbpath ./db --repair --directoryperdb

db_stop:
	@echo "Stop MongoDB"
	pkill -x mongod

test:
	pytest poptimizer -v --cov=poptimizer --cov-report=term-missing --cov-report=xml --setup-show

test_go:
	cd data;go test ./... -covermode=atomic -race

lint:
	mypy poptimizer
	flake8 poptimizer

lint_go:
	cd data;golangci-lint run