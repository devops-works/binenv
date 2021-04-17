FROM debian:buster-slim

RUN apt-get update && \
    apt-get install -y \
        bash \
        ca-certificates \
        curl \
        wget \
    && \
    rm -rf /var/lib/apt

RUN useradd -ms /bin/bash binenv
USER binenv
WORKDIR /home/binenv

RUN wget -q https://github.com/devops-works/binenv/releases/latest/download/binenv_linux_amd64 -O binenv && \
    chmod +x binenv && \
    ./binenv update && \
    ./binenv install binenv && \
    rm binenv && \
    echo 'export PATH=~/.binenv:$PATH' >> ~/.bashrc && \
    echo "source <(binenv completion bash)" >> ~/.bashrc

ADD docker/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]