build-cli-plugin:
	go build -o ./dist/ -ldflags "-s -w -X github.com/koki-develop/docker-tags/cmd.cliPlugin=true"
