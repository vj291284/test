#!/bin/bash

PACKAGE="main"

gen_version_file() {
    echo "package $PACKAGE
//Version: dynamically generated at build time for
//version reference
var Version =  struct{
     BuildTime string
     BuildUser string
     BuildHost string
     GitTag    string
     GitBranch string
}{
    BuildTime: \"$(date)\",
    BuildUser: \"$(whoami)\",
    BuildHost: \"$(hostname)\",
    GitTag: \"$(git describe --always --tag --dirty)\",
    GitBranch: \"$(git rev-parse --abbrev-ref HEAD)\",
}" | gofmt > version.go
    return $?
}

go_run() {
    echo "Running"
    go build $*
    local retval=$?
    if [[ $retval -ne 0 ]]; then
	return $retval
    fi
}

go_build() {
    echo "Building"
    go build $*
    local retval=$?
    if [[ $retval -ne 0 ]]; then
	return $retval
    fi
}

go_test() {
    echo "Testing"
    go test $*
    local retval=$?
    if [[ $retval -ne 0 ]]; then
	return $retval
    fi
}

echo "Gen version"
gen_version_file
RETVAL=$?
if [[ $RETVAL -ne 0 ]]; then
    echo "Gen version failed"
    exit $RETVAL
fi

CMD="$1"
shift
case "$CMD" in 
    run)
	go_run $*
	exit $?
	;;
    build)
	go_build $*
	exit $?
	;;
    test)
	go_test $*
	exit $?
	;;
    *)
	echo "Usage: $0 (run|build|test) <args>"
	exit 1
esac
