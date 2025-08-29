GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

# 生成配置文件
conf:
	@buf generate

.PHONY: fmt
# 格式化代码
fmt:
	@gofmt -s -w .

.PHONY: vet
# 代码检查 vet
vet:
	@go vet ./...

.PHONY: ci-lint
# 代码检查 lint
lint:
	@golangci-lint run ./...

# git 记录清除
git-clean:
	#清除开始
	@git checkout --orphan latest_branch
	@git add -A
	@git commit -am "clean"
	@git branch -D ${gitBranch}
	@git branch -m ${gitBranch}
	@git push -f origin ${gitBranch}
	#清除结束

# 创建新的 tag
git-new-tag:
	# 获取当前最新的 tag
	@echo "当前最新的 tag: $$(git tag -l | sort -V | tail -n 1)"
	# 生成新的版本号
	$(eval NEW_VERSION := $(shell git tag -l | sort -V | tail -n 1 | awk -F. '{$$NF=$$NF+1; print $$1"."$$2"."$$NF}'))
	@echo "新的版本号: ${NEW_VERSION}"
	# 创建新的 tag
	@git tag -a ${NEW_VERSION} -m "release ${NEW_VERSION}"
	# 推送新的 tag
	@git push origin ${NEW_VERSION}

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
