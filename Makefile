APPNAME = "blog"
# h - help
h help:
	@echo "h help 	- this help"
	@echo "build 	- build and the app"
	@echo "run 	- run the app"
	@echo "clean 	- clean app trash"
	@echo "swag 	- generate swag docs"
	@echo "dev 	- generate docs and run"
	@echo "test 	- run all tests"
.PHONY: h

# build - build the app
build:
	go build -o $(APPNAME)
.PHONY: build

# run - build and run the app
run: build
	./$(APPNAME)
.PHONY: run

clean:
	rm ./$(APPNAME)
.PHONY: clean

# swag - generate swagger docs
swag:
	swag init
.PHONY: swag

# dev - generate docs and run
dev: swag run
.PHONY: dev

# test - run all tests
test:
	go test ./infra
.PHONY: test
