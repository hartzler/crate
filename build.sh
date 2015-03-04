NAME=crate
HARDWARE=$(uname -m)
VERSION=0.0.1

# setup GOPATH to be absolute path of script
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
GOPATH=$DIR
echo $GOPATH
export GOPATH

vendor() {
	go get armada/$NAME
	find src -name ".git" -exec rm -rf {} \;
	find src -name ".hg" -exec rm -rf {} \;
}

case $1 in
--vendor) vendor;;
--int) TEST_OPTIONS="-tags int"
esac

# setup
mkdir -p bin/
if [ -f bin/$NAME ]; then
  rm bin/$NAME
fi

# fmt & simplify
gofmt -w -s src/armada/$NAME/

# build and test
go build -o bin/$NAME armada/$NAME/main && go test $TEST_OPTIONS armada/$NAME/...
