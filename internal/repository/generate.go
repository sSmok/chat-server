package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i ChatRepositoryI -o ./mocks/ -s "_minimock.go"
