default: test

clean:
	@for I in spikes/*/main.go; do \
	  pushd `dirname "$$I"` >/dev/null; \
	  rm `basename $$PWD` 2>/dev/null || true; \
	  popd >/dev/null; \
	done

spikes:
	@for I in spikes/*/main.go; do \
	  echo ; \
	  echo $$I; \
	  pushd `dirname "$$I"` >/dev/null; \
	  echo `basename $$PWD` > .gitignore; \
	  echo ----------; \
	  echo '1+2' | go run *.go; \
	  echo ==========; \
	  popd >/dev/null; \
	  echo ; \
	done

test:
	@scripts/run-all-tests
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true

test-coverage:
	mkdir -p dist
	go test -coverprofile=dist/coverage.out ./internal/...
	go tool cover -html=dist/coverage.out

test-docker:
	docker-compose --version
	docker-compose up --abort-on-container-exit --exit-code-from=go --force-recreate

.PHONY: clean default spikes test test-docker
