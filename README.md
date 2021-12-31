# cedar-server

[Contributors |](https://github.com/karthiklsarma/cedar-server/graphs/contributors)
[Forks |](https://github.com/karthiklsarma/cedar-server/network/members)
[Issues |](https://github.com/karthiklsarma/cedar-server/issues)
[MIT License |](https://github.com/karthiklsarma/cedar-server/blob/main/LICENSE)

## To build and run cedar-server project on localhost:

From ./cedar-server Directory:

### To build cedar-server

- Execute
  > go build -o ./bin/cedar-server

### To run

- Execute from ./cedar-server
  > /bin/cedar-server

## To build and run project on docker container on local:

From ./cedar-server Directory:

### To build

- Execute
  > docker build . -t cedar-server:latest

### To run

- Execute
  > docker run -d -p 8080:8080 `<IMAGE ID from previous step>`

## To deploy on azure
- Once the Kubernetes cluster and Container registry is deployed, Execute
  > az acr build --registry cedarcr --image cedar-server:v1 .
