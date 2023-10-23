export GOFLAGS=-trimpath

OUT_DIR?=_tmp/built

$(OUT_DIR):
	mkdir -p $(OUT_DIR)

EXENAME := contacts

GOOS = _

.build: OUT=$(OUT_DIR)/$(GOOS)/$(EXENAME)
.build: $(OUT_DIR)
	go build -o $(OUT)

.PHONY: build.linux
build.linux: GOOS = linux
build.linux: .build

.PHONY: build.macos
build.macos: GOOS = darwin
build.macos: .build

.PHONY: deployable
deployable: build.linux

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)