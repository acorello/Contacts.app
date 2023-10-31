OUT_DIR?=_tmp/built
ARTEFACT_NAME?=contacts-app

$(OUT_DIR):
	@mkdir -p $(OUT_DIR)

.PHONY: local.build
local.build: ARTEFACT_PATH=$(OUT_DIR)/$(ARTEFACT_NAME)
local.build: $(OUT_DIR)
	@echo "Building" $(ARTEFACT_PATH)
	@go build -trimpath -o $(ARTEFACT_PATH)

.PHONY: local.clean
local.clean:
	@echo "Deleting $(OUT_DIR)"
	@rm -rf $(OUT_DIR)

.PHONY: docker.image
docker.image:
	@echo "Building docker image: " $(ARTEFACT_PATH)
	@docker build --tag $(ARTEFACT_NAME) --file _docker/Dockerfile .
