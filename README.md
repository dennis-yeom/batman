# batman
file watcher

## creating module
```
go mod init projectName
```

## cobra & viper set up
make sure you have cobra and viper set up:
```
go get github.com/spf13/cobra
go get github.com/spf13/viper
```

## running test command line
```
go run main.go --help
```
![alt text](image.png)


```
go run main.go test
```
![alt text](image-1.png)


## reddis client in go
https://redis.io/docs/latest/develop/connect/clients/go/

go get command to install reddis:
```
go get github.com/redis/go-redis/v9
```

## starting redis server
```
redis-server --port 6380
```

## setting.getting values in redis
```
go run main.go set -k dennis -v 1995
```

```
go run main.go get -k dennis
```

![alt text](image-3.png)