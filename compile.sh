export GOPATH=/root/go/
export GOBIN=/root/go/src/pineapple/bin
go install pineapple/src/master
go install pineapple/src/server
go install pineapple/src/client
go install pineapple/src/clientnew
export GOBIN=