env GOOS=linux go build -o bin/import import.go
env GOOS=linux go build -o bin/search search.go
sls deploy -v
