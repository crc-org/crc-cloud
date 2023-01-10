FROM quay.io/centos/centos:stream9
ENV CONTAINER true
ENV PATH $PATH:/root/.awscliv2/binaries
RUN rpm -Uvh https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm &&\
    dnf install -y          \
    which                   \
    jq                      \
    nc                      \
    pip                     \
    figlet                  \
    openssh-clients         \
    less &&                 \
    pip install awscliv2 && \
    awscliv2 -i 

WORKDIR /app
COPY . .
RUN chmod +x crc-cloud.sh

ENTRYPOINT [ "/app/crc-cloud.sh" ]

