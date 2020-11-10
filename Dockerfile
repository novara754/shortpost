FROM golang:alpine

WORKDIR /shortpost
COPY . .

RUN mkdir build
RUN go build -o ./build/shortpost ./shortpost

CMD ["./build/shortpost"]
