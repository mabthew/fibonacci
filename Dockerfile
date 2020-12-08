# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
WORKDIR /src
COPY . .

RUN go get github.com/golang/groupcache/lru

RUN go get github.com/julienschmidt/httprouter

RUN go build -o fibonacci . 

CMD ./fibonacci