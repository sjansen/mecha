default: spikes

spikes:
	@for I in spikes/*/main.go; do \
	  echo ; \
	  echo $$I; \
	  echo ----------; \
	  echo '1+2' | go run $$I; \
	  echo ==========; \
	  echo ; \
	done

.PHONY: default spikes
