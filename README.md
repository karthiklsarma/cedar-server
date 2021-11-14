# cedar-engine

[Contributors |](https://github.com/karthiklsarma/cedar-engine/graphs/contributors)
[Forks |](https://github.com/karthiklsarma/cedar-engine/network/members)
[Issues |](https://github.com/karthiklsarma/cedar-engine/issues)
[MIT License |](https://github.com/karthiklsarma/cedar-engine/blob/main/LICENSE)

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

- Execute deploy cedar [script](https://github.com/karthiklsarma/cedar-deploy/blob/main/cedar-deploy.sh)
- Once the Kubernetes cluster and Container registry is deployed, Execute
  > az acr build --registry cedarcr --image cedar-server:v1 .
