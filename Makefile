# Original: https://github.com/pipe-cd/pipecd/blob/plugin/master/Makefile

.PHONY: build/plugin
build/plugin: PLUGINS_BIN_DIR ?= ~/.piped/plugins
build/plugin: PLUGINS_SRC_DIR ?= ./plugins/
build/plugin: PLUGINS_OUT_DIR ?= $(shell pwd)/.artifacts
build/plugin: PLUGINS ?= $(shell find $(PLUGINS_SRC_DIR) -mindepth 1 -maxdepth 1 -type d | while read -r dir; do basename "$$dir"; done | paste -sd, -) # comma separated list of plugins. eg: PLUGINS=kubernetes,ecs,lambda
build/plugin:
	mkdir -p $(PLUGINS_BIN_DIR)
	mkdir -p $(PLUGINS_OUT_DIR)
	@echo "Building plugins..."
	@for plugin in $(shell echo $(PLUGINS) | tr ',' ' '); do \
		echo "\nBuilding plugin: $$plugin"; \
		pushd $(PLUGINS_SRC_DIR)/$$plugin > /dev/null; \
		go build -o $(PLUGINS_OUT_DIR)/$$plugin ./ \
			&& cp $(PLUGINS_OUT_DIR)/$$plugin $(PLUGINS_BIN_DIR)/$$plugin; \
		popd > /dev/null; \
	done
	@echo "Plugins are built and copied to $(PLUGINS_BIN_DIR)"

.PHONY: test/plugins
test/plugins: PLUGINS_SRC_DIR ?= ./plugins/
test/plugins: PLUGINS ?= $(shell find $(PLUGINS_SRC_DIR) -mindepth 1 -maxdepth 1 -type d | while read -r dir; do basename "$$dir"; done | paste -sd, -) # comma separated list of plugins. eg: PLUGINS=kubernetes,ecs,lambda
test/plugins:
	@for plugin in $(shell echo $(PLUGINS) | tr ',' ' '); do \
		echo "\nTest plugin: $$plugin"; \
		pushd $(PLUGINS_SRC_DIR)/$$plugin; \
		go test -failfast -race ./...; \
		popd; \
	done

