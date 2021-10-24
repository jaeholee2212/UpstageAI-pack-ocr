FROM golang:alpine as stage

COPY . /src
WORKDIR /src
RUN go build

FROM alpine
COPY --from=stage /src/clogs /app/clogs
RUN ln -s /app/clogs /usr/local/bin/clogs
WORKDIR /app
CMD ["clogs", "--port", "9977"]