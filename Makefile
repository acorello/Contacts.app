ARTEFACT_PATH?=contacts-app
IMAGE_NAME=contacts-app
IMAGE_TAGGED_NAME=contacts-app:latest
CONTAINER_REGISTRY=registry.digitalocean.com/acorello
IMAGE_REPOSITORY_TAGGED_NAME=$(CONTAINER_REGISTRY)/$(IMAGE_TAGGED_NAME)

.PHONY: executable
executable:
	@echo "Building" $(ARTEFACT_PATH)
	@go build -trimpath -o $(ARTEFACT_PATH)

.PHONY: container.image
container.image:
	@echo "Building image: " $(IMAGE_NAME)
	@docker build \
		--tag $(IMAGE_TAGGED_NAME) \
		--tag $(IMAGE_REPOSITORY_TAGGED_NAME) \
		--file _docker/Dockerfile .

.PHONY: deployed
deployed: container.image
	@docker push $(IMAGE_REPOSITORY_TAGGED_NAME)