all: build

build:
	go build -o bin/xserver src/main.go

example-tests:
	python3 -m pytest -x -v tests