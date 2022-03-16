#
# this dockerfile is for tests purposes only
# see e2e target in Makefile
#

FROM debian:bookworm

RUN useradd -ms /bin/bash test

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && update-ca-certificates

COPY bin/binenv /home/test/binenv

USER test
WORKDIR /home/test

RUN ./binenv update && ./binenv install binenv && rm binenv
RUN echo -e '\nexport PATH=~/.binenv:$PATH' >> ~/.bashrc
RUN echo 'source <(binenv completion bash)' >> ~/.bashrc

# COPY distributions/distributions_test.yaml /home/test/.config/binenv/distributions.yaml
COPY distributions/distributions.yaml /home/test/.config/binenv/distributions.yaml
# COPY DISTRIBUTIONS_test.md /home/test/.config/binenv/DISTRIBUTIONS.md
COPY DISTRIBUTIONS.md /home/test/.config/binenv/DISTRIBUTIONS.md

COPY scripts/e2e_tests.sh /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/e2e_tests.sh"]