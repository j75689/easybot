## Build

```bash
cd server
rice embed-go

cd ..
go build
```

## Build Plugin

#### Plugins

```bash
go build -buildmode=plugin -o http.so http.go
go build -buildmode=plugin -o ./plugins/Request.so ./plugins/Request.go
```

#### goloader

```bash
cd plugins
go tool compile -I $GOPATH/pkg/darwin_amd64 ../lib/plugins/*.go
```

## HEROKU

#### deploy

```bash
git push heroku master
```
