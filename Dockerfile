FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY .env .
COPY out /bin/out

ENTRYPOINT [ "/bin/out" ]