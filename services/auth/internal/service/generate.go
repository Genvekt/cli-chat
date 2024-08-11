package service

// Generate mocks for service layer

// Run the following commands from service root:
// "make install-deps" to install minimock into auth/bin
// "go generate ./internal/service" to generate service mocks

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../bin/minimock -i UserService -o mocks -s "_minimock.go" -g
