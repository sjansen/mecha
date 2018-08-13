default: test

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
	go test -tags integration ./internal/... ./spikes/...
	@echo ========================================
	go vet ./internal/... ./spikes/...
	golint -set_exit_status internal/
	golint -set_exit_status spikes/
	gocyclo -over 15 internal/ spikes/
	@echo ========================================
	@git grep TODO  internal/ spikes/ || true
	@git grep FIXME internal/ spikes/ || true

.PHONY: default spikes test
