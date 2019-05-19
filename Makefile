BINARY_NAME=EventsStore
BINARY_UNIX=$(BINARY_NAME)_unix
SOURCE_DIRECTORY=src


all: clean deps test build

build: 
		cd $(SOURCE_DIRECTORY) && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o  $(BINARY_NAME) -v
test: 
		go test -v ./... 
clean: 
		cd $(SOURCE_DIRECTORY) && go clean && rm -f $(BINARY_NAME) && rm -f $(BINARY_UNIX)
deps:
		cd $(SOURCE_DIRECTORY) && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && dep ensure