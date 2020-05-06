package toolkit

//go:generate go get github.com/golang/mock/mockgen
//go:generate go install github.com/golang/mock/mockgen
//go:generate mockgen -source=./web/httpclient.go -destination=./mock/mock_httpclient.go -package=mock
