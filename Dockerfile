FROM golang:1.16

WORKDIR /app

# Create the docker image to run the workers
ADD . .

RUN go build -trimpath -a -o map_reduce -ldflags="-w -s" *.go

CMD [ "./map_reduce" ]