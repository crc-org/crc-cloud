FROM quay.io/centos/centos:stream9
ENV CONTAINER true
ENV PATH $PATH:/root/.awscliv2/binaries
RUN rpm -Uvh https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm &&\
    dnf install -y which &&\
    dnf install -y jq &&\
    dnf install -y nc &&\
    dnf install -y pip &&\
    dnf install -y figlet &&\
    dnf install -y openssh-clients &&\
    dnf install -y less &&\
    pip install awscliv2 &&\
    awscliv2 -i 

WORKDIR /app
COPY . .
RUN chmod +x openspot-ng.sh

ENTRYPOINT [ "/app/openspot-ng.sh" ]

