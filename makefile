default: test

DIRS=`go list ./... | grep -v /spikes/`

spikes:
	@for I in spikes/*/main.go; do \
	  echo ; \
	  echo $$I; \
	  pushd `dirname "$$I"` >/dev/null; \
	  echo ----------; \
	  echo '1+2' | go run *.go; \
	  echo ==========; \
	  popd >/dev/null; \
	  echo ; \
	done

test:
	go test -tags integration ./...
	@echo ========================================
	go vet ./...
	golangci-lint run
	#golint -set_exit_status $(DIRS)
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true

.PHONY: default spikes test
