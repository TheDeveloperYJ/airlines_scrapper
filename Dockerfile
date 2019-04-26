FROM golang:1.8.5-jessie
# create a working directory

RUN go get "github.com/PuerkitoBio/goquery"
RUN go get "github.com/gorilla/mux"

WORKDIR /go/src/app
# add source code
ADD src src
# run main.go
CMD ["go", "run", "src/scrapper.go"]