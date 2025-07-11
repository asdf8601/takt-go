NAME=takt
head=5
bin=./bin/takt
testfile=./data/test_functional.csv

.PHONY: bin
bin:
	rm -rf $(bin)
	mkdir -p ./bin
	go build -o ./bin/$(NAME) main.go

.PHONY: build
build: bin

install: build
	cp -f $(bin) ~/.local/$(bin)

.PHONY: all
all: build lint test functional-test
	cp -f $(bin) ~/.local/$(bin)

.PHONY: test
test:
	@echo "Running Go unit tests..."
	go test -v
	@echo "Running integration tests..."
	go test -v -run TestCLIIntegration

.PHONY: functional-test
functional-test: build setup-test-file
	@echo "Running functional tests with test file: $(testfile)"
	TAKT_FILE=$(testfile) $(bin) version
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) cat $(head)
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) day $(head)
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) d $(head)
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) week $(head)
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) w $(head)
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) month $(head)
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) m $(head)
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) year $(head)
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) y $(head)
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) check "deleteMe"
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) c "deleteMe"
	@sleep 1
	@echo
	TAKT_FILE=$(testfile) $(bin) cat $(head)
	@sleep 1
	@echo "Functional tests completed. Cleaning up..."
	@$(MAKE) clean-test-file

.PHONY: setup-test-file
setup-test-file:
	@echo "Setting up test file: $(testfile)"
	@mkdir -p $(dir $(testfile))
	@echo "timestamp,kind,notes" > $(testfile)
	@echo "2024-07-26T18:00:00+02:00,out," >> $(testfile)
	@echo "2024-07-26T09:00:00+02:00,in," >> $(testfile)
	@echo "2024-07-25T15:00:19+02:00,out," >> $(testfile)
	@echo "2024-07-25T14:00:16+02:00,in," >> $(testfile)
	@echo "2024-06-25T13:00:11+02:00,out," >> $(testfile)
	@echo "2024-06-25T12:00:01+02:00,in," >> $(testfile)

.PHONY: clean-test-file
clean-test-file:
	@if [ -f "$(testfile)" ]; then \
		echo "Removing test file: $(testfile)"; \
		rm -f "$(testfile)"; \
	fi
	@if [ -f "$(testfile).bak" ]; then \
		echo "Removing test backup file: $(testfile).bak"; \
		rm -f "$(testfile).bak"; \
	fi

.PHONY: gotest
gotest:
	go test -v


lint:
	go fmt ./...
	go vet ./...
