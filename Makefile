test:
	go generate
	go test -race -v ./...

readme:
	echo '# gitsync [![Build Status](https://travis-ci.org/rdkr/gitsync.svg)](https://travis-ci.org/rdkr/gitsync) [![codecov.io](https://codecov.io/github/rdkr/gitsync/coverage.svg)](https://codecov.io/github/rdkr/gitsync)' > README.md
	echo '```' >> README.md
	go run main.go --help >> README.md
	echo '```' >> README.md
