go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen

echo $GOPATH/bin/mockgen


go mod tidy

go generate ./...