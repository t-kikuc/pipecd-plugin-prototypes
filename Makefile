# Original: https://github.com/pipe-cd/pipecd/blob/plugin/master/Makefile

.PHONY: build/plugin
build/plugin: PLUGINS_BIN_DIR ?= ~/.piped/plugins
build/plugin: PLUGINS_SRC_DIR ?= ./pkg/app/pipedv1/plugin
build/plugin: PLUGINS_OUT_DIR ?= ./.artifacts/plugins
build/plugin: PLUGINS ?= $(shell find $(PLUGINS_SRC_DIR) -mindepth 1 -maxdepth 1 -type d | while read -r dir; do basename "$$dir"; done | paste -sd, -) # comma separated list of plugins. eg: PLUGINS=kubernetes,ecs,lambda
build/plugin:
	mkdir -p $(PLUGINS_BIN_DIR)
	@echo "Building plugins..."
	@for plugin in $(shell echo $(PLUGINS) | tr ',' ' '); do \
		echo "Building plugin: $$plugin"; \
		go build -o $(PLUGINS_OUT_DIR)/$$plugin $(PLUGINS_SRC_DIR)/$$plugin \
			&& cp $(PLUGINS_OUT_DIR)/$$plugin $(PLUGINS_BIN_DIR)/$$plugin; \
	done
	@echo "Plugins are built and copied to $(PLUGINS_BIN_DIR)"

.PHONY: test/go
test/go: COVERAGE ?= false
test/go: COVERAGE_OPTS ?= -covermode=atomic
test/go: COVERAGE_OUTPUT ?= coverage.out
test/go:
ifeq ($(COVERAGE), true)
	go test -failfast -race $(COVERAGE_OPTS) -coverprofile=$(COVERAGE_OUTPUT).tmp ./pkg/... ./cmd/...
	cat $(COVERAGE_OUTPUT).tmp | grep -v ".pb.go\|.pb.validate.go" > $(COVERAGE_OUTPUT)
	rm -rf $(COVERAGE_OUTPUT).tmp
else
	go test -failfast -race ./pkg/... ./cmd/...
endif

.PHONY: run/piped
run/piped: CONFIG_FILE ?=
run/piped: INSECURE ?= false
run/piped: LAUNCHER ?= false
run/piped: LOG_ENCODING ?= humanize
run/piped: EXPERIMENTAL ?= false
run/piped:
ifeq ($(EXPERIMENTAL), true)
	go run cmd/pipedv1/main.go piped --tools-dir=/tmp/piped-bin --config-file=$(CONFIG_FILE) --insecure=$(INSECURE) --log-encoding=$(LOG_ENCODING)
else ifeq ($(LAUNCHER),true)
	go run cmd/launcher/main.go launcher --config-file=$(CONFIG_FILE) --insecure=$(INSECURE) --log-encoding=$(LOG_ENCODING)
else
	go run cmd/piped/main.go piped --tools-dir=/tmp/piped-bin --config-file=$(CONFIG_FILE) --insecure=$(INSECURE) --log-encoding=$(LOG_ENCODING)
endif
