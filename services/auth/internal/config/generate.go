package config

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../bin/minimock -i UserServiceConfig -o mocks -s "_minimock.go" -g
