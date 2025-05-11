build-cli-plugin:
	go build -o ./dist/ -ldflags "-s -w -X github.com/koki-develop/docker-tags/cmd.cliPlugin=true"
	mkdir -p $(HOME)/.docker/cli-plugins/
	cp $(PWD)/dist/docker-tags $(HOME)/.docker/cli-plugins/
