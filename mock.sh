go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen

echo $GOPATH/bin/mockgen


go mod tidy

go mod download github.com/golang/mock

go generate ./...