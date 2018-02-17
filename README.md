# gin-commitment
[![CircleCI](https://circleci.com/gh/tyauvil/gin-commitment/tree/develop.svg?style=svg)](https://circleci.com/gh/tyauvil/gin-commitment/tree/develop)

### Go webapp using Gin framework
https://github.com/gin-gonic/gin

#### Endpoints
This app implements the following endpoints:
```
/ - returns html page with the commit message
/p/:sha - permanent link to a html commit message based on short sha256 hash
/json - returns message in json format { "message" : "wip" }
/commit.txt - returns plain text message
/healthz - returns 200 and "OK"
/robots.txt - returns generic robots.txt
```

#### Configuration
The following environment variables can be set:
```
GIN_MODE=[debug,release]
DOMAIN=localhost
TLS=true - enables TLS (setting DOMAIN required)
```

#### Running locally
```
docker-compose up
```

#### todo
Add more tests

Inspired by https://github.com/ngerakines/commitment
