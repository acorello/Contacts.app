ARTEFACT_PATH?=contacts-app
IMAGE_NAME=contacts-app
IMAGE_TAGGED_NAME=$(IMAGE_NAME):latest
CONTAINER_REGISTRY=registry.digitalocean.com/acorello
IMAGE_REPOSITORY_TAGGED_NAME=$(CONTAINER_REGISTRY)/$(IMAGE_TAGGED_NAME)

.PHONY: executable
executable:
	@echo "Building" $(ARTEFACT_PATH)
	@go build -trimpath -o $(ARTEFACT_PATH)

.PHONY: container.image
container.image:
	@echo "Building image: " $(IMAGE_TAGGED_NAME)
	@docker build \
		--tag $(IMAGE_TAGGED_NAME) \
		--tag $(IMAGE_REPOSITORY_TAGGED_NAME) \
		--file _docker/Dockerfile .

.PHONY: container
container: container.image
	@docker run --rm --env HOST=0.0.0.0 -p 8080:8080 $(IMAGE_TAGGED_NAME)

.PHONY: deployed
deployed: container.image
	@docker push $(IMAGE_REPOSITORY_TAGGED_NAME)