OUT_DIR?=_tmp/built

.PHONY: .build
.build: ARTEFACT_NAME?=contacts.$(GOOS)
.build: ARTEFACT_PATH=$(OUT_DIR)/$(ARTEFACT_NAME)
.build:
	@mkdir -p $(OUT_DIR)
	@echo "Building" $(ARTEFACT_PATH)
	@GOOS=$(GOOS) go build -trimpath -o $(ARTEFACT_PATH)

.PHONY: build.linux
build.linux: GOOS = linux
build.linux: .build

.PHONY: build.macos
build.macos: GOOS = darwin
build.macos: .build

.PHONY: clean
clean:
	@echo "Deleting $(OUT_DIR)"
	@rm -rf $(OUT_DIR)