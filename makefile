.PHONY:  clean  default  linux  macos  refresh  spikes  test  test-coverage  test-docker

default: test

clean:
	@for I in spikes/*/main.go; do \
	  pushd `dirname "$$I"` >/dev/null; \
	  rm `basename $$PWD` 2>/dev/null || true; \
	  popd >/dev/null; \
	done

linux:
	GOOS=linux GOARCH=amd64 go build -o mecha

macos:
	GOOS=darwin GOARCH=amd64 go build -o mecha

refresh:
	cookiecutter gh:sjansen/cookiecutter-golang --output-dir .. --config-file .cookiecutter.yaml --no-input --overwrite-if-exists
	git checkout go.mod go.sum

spikes:
	scripts/run-all-spikes

test:
	@scripts/run-all-tests
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true

test-coverage:
	mkdir -p dist
	go test -coverpkg .,./internal/... -coverprofile=dist/coverage.txt -tags integration ./...
	go tool cover -html=dist/coverage.txt

test-docker:
	docker-compose --version
	docker-compose build --pull go
	docker-compose up --abort-on-container-exit --exit-code-from=go --force-recreate
