FROM alpine:latest

RUN apk add --no-cache bash curl openjdk11-jre

COPY synopsys-detect-9.7.0.jar /opt/

COPY release/linux/amd64/blackduck-plugin /bin/

ENTRYPOINT ["/bin/blackduck-plugin"]