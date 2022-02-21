FROM fedora:34
ENV GOPATH=/go
ENV GOPROXY=https://proxy.golang.org
ENV PATH=$PATH:/go/bin
RUN dnf install -y git make gcc gcc-c++ which iproute iputils procps-ng vim-minimal tmux net-tools htop tar jq npm openssl-devel perl rust cargo golang
RUN curl https://get.starport.network/starport! | bash
# the actual source code for this repo, this **only** includes checked in files!
# this is a bit of a pain but it does speed things up a lot
ADD baseledger.tar.gz /
# build steps for all codebases in this repo, must be below the add statement
RUN pushd /baseledger/baseledger && PATH=$PATH:/usr/local/bin/ starport chain build
RUN pushd /baseledger/orchestrator && PATH=$PATH:$HOME/.cargo/bin cargo build --all
RUN pushd /baseledger/orchestrator && PATH=$PATH:$HOME/.cargo/bin cargo run -- init

# keep the container running
CMD ["tail", "-f", "/dev/null"] 
