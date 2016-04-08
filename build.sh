NAME=crate
PACKAGE=github.com/armada-io/crate
HARDWARE=$(uname -m)
VERSION=0.0.1
#GOOS=linux
#GOARCH=amd64

# setup GOPATH to be absolute path of script
export GOPATH=/vagrant


case $1 in
--int) TEST_OPTIONS="-tags int"
esac

# fmt & simplify
gofmt -w -s . command/ pid1/

# build and test
#go build -o bin/$NAME armada/$NAME/main && 
go test $TEST_OPTIONS $PACKAGE && go install $PACKAGE

