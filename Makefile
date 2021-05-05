PHONY: run

run:
	go run main.go > .semaphore/semaphore.yml

test:
	go run main.go > .semaphore/semaphore.test.yml