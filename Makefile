all: build serve

serve:
	go run main.go

build:
	go run build/build.go .

.PHONY: build