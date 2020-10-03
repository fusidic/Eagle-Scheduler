local:
	GOOS=linux GOARCH=amd64 go build -o eagle-scheduler ./cmd/scheduler 
build:
	docker build --no-cache . -t fusidic/eagle-scheduler:0.3
push:
	docker push fusidic/eagle-scheduler:0.3

format:
	sudo gofmt -l -w .
clean:
	sudo rm -f scheduler