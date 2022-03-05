all: build serve

serve:
	go run main.go

build:
	cd build; go run build.go

clean:
	rm -rf public

.PHONY: build