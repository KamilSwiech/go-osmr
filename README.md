# go-osmr
Golang web app for parsing route request from OSMR project into json with most important data.

## Example Usage
```
$ curl 'http://localhost:8080/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219'
```
Where `src` is source and `dst` is destination. There can be only one source but multiple destiations.

## Installation
> Prerequsite: Installed docker

### From repo
1. Clone repostitory:
```
$ git clone https://github.com/KamilSwiech/go-osmr.git
```
2. Cd to go-osmr directory and run docker build with custom tag:
```
$ cd go-osmr
$ docker build . -t go-osmr
```
3. Run docker image in detached mode with exposed port 8080
```
$ docker run -d -p 8080:8080 go-osmr
```

### From dockerhub
```
$ docker run -d -p 8080:8080 swiechkamil/go-osmr
```

