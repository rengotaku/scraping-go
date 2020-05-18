# scraping-go
scraping website go

# environment
https://github.com/gin-gonic/gin/blob/master/mode.go

* debug
* release
* test

# build
## web
$ GOOS=linux GOARCH=amd64 go build -o ./scraping_web

## batch
$ cd batch
$ GOOS=linux GOARCH=amd64 go build -o ./batch