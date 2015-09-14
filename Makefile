TARGET_OSARCH="linux/386"

.PHONY: all clean deps build

all: clean deps build

deps:
	go get -d -v ./...
	go get github.com/mitchellh/gox
	go get github.com/garyburd/redigo/redis

build: deps
	gox -osarch=$(TARGET_OSARCH) -output monitor-sidekiq-dead-queue

clean:
	rm -f monitor-sidekiq-dead-queue

