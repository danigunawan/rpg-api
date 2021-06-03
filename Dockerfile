# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.16.9 as builder

RUN mkdir -p $GOPATH/https://github.com/danigunawan
WORKDIR $GOPATH/https://github.com/danigunawan

ADD . .
ENV GO111MODULE=on
RUN go get ./...
RUN buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .

# Uncomment to run the binary in "production" mode:
# ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

RUN mkdir -p /keys
ADD keys/rsapub.pem /keys
ADD keys/rsakey.pem /keys

ENV JWT_PUBLIC_KEY=/keys/rsapub.pem
ENV JWT_PRIVATE_KEY=/keys/rsakey.pem

# Uncomment to run the migrations before running the binary:
CMD /bin/app migrate; /bin/app
# CMD exec /bin/app
