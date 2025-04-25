OUT_DIR := out
APP_NAME := ttedit

.PHONY: all
all: clean build

build: 
	@mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(APP_NAME) ./main.go

clean:
	rm -rf $(OUT_DIR)
