FROM quay.io/centos/centos:stream9

RUN dnf install -y ansible-core sudo && \
    ansible-galaxy collection install community.general \
    community.crypto ansible.posix \
    kubernetes.core


RUN groupadd --gid 1000 user && \
    useradd --uid 1000 --gid 1000 -d /home/user user
COPY entrypoint.sh /home/user/entrypoint.sh
RUN chown -R user:user /home/user && chmod +x /home/user/entrypoint.sh

USER user

WORKDIR /home/user
ENTRYPOINT ["/home/user/entrypoint.sh"]
