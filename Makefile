export GOFLAGS=-trimpath

OUT_DIR?=_tmp/built
ARTEFACT_NAME?=contacts

$(OUT_DIR):
	@mkdir -p $(OUT_DIR)

.build: ARTEFACT_PATH=$(OUT_DIR)/$(GOOS)/$(ARTEFACT_NAME)
.build: $(OUT_DIR)
	@echo "Building" $(ARTEFACT_PATH)
	@go build -o $(ARTEFACT_PATH)

.PHONY: build.linux
build.linux: GOOS = linux
build.linux: .build

.PHONY: build.macos
build.macos: GOOS = darwin
build.macos: .build

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)