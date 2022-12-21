FROM alpine:3.17.0
ENV CONTAINER true
RUN apk add --no-cache aws-cli && \
    apk add --no-cache netcat-openbsd && \
    apk add --no-cache jq && \
    apk add --no-cache netcat-openbsd && \
    apk add --no-cache bash && \
    apk add --no-cache sed && \ 
    apk add --no-cache curl && \
    apk add --no-cache figlet && \
    apk add --no-cache openssh-client-default 
WORKDIR /app
COPY . .
RUN chmod +x openspot-ng.sh

ENTRYPOINT [ "/app/openspot-ng.sh" ]

