#
# this dockerfile is for tests purposes only
# see e2e target in Makefile
#

FROM debian:bookworm

RUN useradd -ms /bin/bash test

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates

COPY bin/binenv /home/test/binenv
COPY scripts/entrypoint.sh /entrypoint.sh

USER test
WORKDIR /home/test

RUN ./binenv update && ./binenv install binenv && rm binenv
RUN echo -e '\nexport PATH=~/.binenv:$PATH' >> ~/.bashrc
RUN echo 'source <(binenv completion bash)' >> ~/.bashrc

ENTRYPOINT ["/entrypoint.sh"]