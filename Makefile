
install:
	go install ./cmd/gosper

test-env:
	docker build -t automata-env -f dev/setup/Dockerfile .
	#docker rmi automata-env
