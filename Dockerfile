FROM alpine:3.17.0
ENV DOCKER true
RUN apk add --no-cache aws-cli && \
    apk add --no-cache netcat-openbsd && \
    apk add --no-cache jq && \
    apk add --no-cache netcat-openbsd && \
    apk add --no-cache bash && \
    apk add --no-cache sed && \ 
    apk add --no-cache curl && \
    apk add --no-cache openssh-client-default 
WORKDIR /app
COPY ["openspot-ng.sh","id_ecdsa_crc", "common.sh","./"]
COPY templates templates
RUN chmod +x openspot-ng.sh

ENTRYPOINT [ "/app/openspot-ng.sh" ]

