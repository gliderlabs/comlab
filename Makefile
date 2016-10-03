NAME := gosper
INSTALL_DIR ?= /usr/local/bin

install:
	glide install
	go build -o $(NAME) ./cmd/gosper
	install -m 755 $(NAME) $(INSTALL_DIR)/$(NAME)

test-env:
	docker build -t automata-env -f dev/setup/Dockerfile .
	#docker rmi automata-env
