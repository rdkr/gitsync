readme:
	echo '# gitsync' > README.md
	echo '```' >> README.md
	go run main.go --help >> README.md
	echo '```' >> README.md
