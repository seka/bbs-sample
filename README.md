## bbs-sample

Simple bbs of golang.

## Installation

```sh
$ git clone https://github.com/seka/bbs-sample $GOPATH/src/github.com/seka/bbs-sample
$ cd $GOPATH/src/github.com/seka/bbs-sample
$ dep ensure
```

dep is the go vendoring tool in detail below.

[golang/dep](https://github.com/golang/dep)

## Run

### Docker

```sh
$ cd $GOPATH/src/github.com/seka/bbs-sample
$ docker-machine start docker
$ eval $(docker-machine env docker)
$ docker-machine ssh docker -L 8080:localhost:8080
$ docker-compose up
```

### Vagrant

```sh
# sample database
$ vagrant box add https://atlas.hashicorp.com/viniciusfs/boxes/centos7/
$ cd $GOPATH/src/github.com/seka/bbs-sample/script
$ vagrant up

# application
$ cd $GOPATH/src/github.com/seka/bbs-sample
$ go run cmd/bbsd/main.go
```

note: Vagrant version less then 1.4.1 if you using macOS
