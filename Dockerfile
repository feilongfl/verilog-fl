From alpine

# enable testing repo to install iverilog
COPY scripts/alpine_setrepo.sh .
RUN  sh alpine_setrepo.sh && apk update

# install package
COPY scripts/requirements.txt .
RUN apk --no-cache add \
    bash \
    python3 \
    python3-dev \
    py3-pip \
    build-base \
    iverilog \
    && pip install -r requirements.txt

# config lib search for ldconfig
# ref: https://www.musl-libc.org/doc/1.0.0/manual.html
RUN echo /lib:/usr/local/lib:/usr/lib:`cocotb-config --lib-dir` > /etc/ld-musl-x86_64.path
