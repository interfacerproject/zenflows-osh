FROM golang:1.19-bullseye AS builder
ADD . /app
WORKDIR /app
RUN go build loshImport.go

FROM dyne/devuan:chimaera
WORKDIR /root
ENV HOST=0.0.0.0
ENV PORT=80
EXPOSE 80
COPY --from=builder /app/loshImport /root/
CMD ["./loshImport.go"]
