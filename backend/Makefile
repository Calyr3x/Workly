.PHONY: start lint checks stop done task-stage

MAKEFLAGS += --no-print-directory
GIT_BRANCH := $(shell git branch --show-current)

start:
	docker-compose up -d
	go run $(CURDIR)/app/main.go

lint:
	@echo "Running linters..."
	@(golangci-lint run ./... > lint.log 2>&1 || (cat lint.log && exit 1))
	@rm -f lint.log
	@echo "No lint errors found."

checks: lint

stop:
	docker-compose down

done:
	@$(MAKE) checks
	@$(MAKE) task-stage

task-stage:
	git add .
	git commit -m "$(commit)"
	git push origin $(GIT_BRANCH)