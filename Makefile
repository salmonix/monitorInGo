# Set variables
BINARY=./_build/gmon

RELEASE=`git describe --tags --abbrev=0`
BRANCH=`git rev-parse --abbrev-ref HEAD`
REVISION=`git rev-parse HEAD`
BUILT=`date --utc +%FT%TZ`
LDFLAGS=-ldflags "-s -w -X main.Release=${RELEASE} -X main.Revision=${REVISION} \
-X main.Built=${BUILT} -X main.ReleaseBuild=true -X main.Branch=${BRANCH}"

.DEFAULT_GOAL: ${BINARY}


# Build - release branch is master
${BINARY}:
	@if [ "`git rev-parse --abbrev-ref HEAD`" != "master" ]; then \
		echo "\nWARNING: compiling from dev branch\n"; \
	fi
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}
	upx -9 ${BINARY}
	chmod 0750 ${BINARY}


clean:
	if [ -f ${BINARY} ];then rm ${BINARY};fi

.PHONY: clean install
