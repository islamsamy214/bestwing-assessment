FROM ubuntu:24.04

LABEL maintainer="Islam Samy"

# Arguments
ARG WWWGROUP=1000
ARG WWWUSER=1000
ARG NODE_VERSION=22
ARG GOLANG_VERSION=1.23.5
# ARG MYSQL_CLIENT="mysql-client"
ARG POSTGRES_VERSION=17

# Workdir
WORKDIR /var/www/html

# Environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC
ENV SUPERVISOR_GO_COMMAND="/usr/local/go/bin/go run main.go"
ENV SUPERVISOR_GO_USER="app"
ENV PGSSLCERT /tmp/postgresql.crt
ENV GOCACHE=/var/tmp/go-cache
ENV GOPATH=/var/www/html/go
ENV GOMODCACHE=/var/www/html/go/pkg/mod
ENV GOBIN=/var/www/html/go/bin

# Define the timezone
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Set apt-get noninteractive
RUN echo "Acquire::http::Pipeline-Depth 0;" > /etc/apt/apt.conf.d/99custom && \
    echo "Acquire::http::No-Cache true;" >> /etc/apt/apt.conf.d/99custom && \
    echo "Acquire::BrokenProxy    true;" >> /etc/apt/apt.conf.d/99custom

# Install dependencies
RUN apt-get update && apt-get upgrade -y \
    && mkdir -p /etc/apt/keyrings \
    && apt-get install -y gnupg gosu curl ca-certificates zip unzip git supervisor sqlite3 libcap2-bin libpng-dev python3 dnsutils librsvg2-bin fswatch ffmpeg nano vim librdkafka-dev

# Install golang
RUN apt-get update && apt-get install -y wget && \
    wget https://go.dev/dl/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
    rm go$GOLANG_VERSION.linux-amd64.tar.gz && \
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile

# # Install nodejs
# RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
#     && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_VERSION.x nodistro main" > /etc/apt/sources.list.d/nodesource.list \
#     && apt-get update \
#     && apt-get install -y nodejs \
#     # && npm install -g npm \
#     # && npm install -g pnpm \
#     && npm install -g bun

# Install database clients
RUN curl -sS https://www.postgresql.org/media/keys/ACCC4CF8.asc | gpg --dearmor | tee /etc/apt/keyrings/pgdg.gpg >/dev/null \
    && echo "deb [signed-by=/etc/apt/keyrings/pgdg.gpg] http://apt.postgresql.org/pub/repos/apt noble-pgdg main" > /etc/apt/sources.list.d/pgdg.list \
    && apt-get update \
    # && apt-get install -y $MYSQL_CLIENT \
    && apt-get install -y postgresql-client-$POSTGRES_VERSION

# Clean up
RUN apt-get -y autoremove \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Create app user
RUN userdel -r ubuntu
RUN groupadd --force -g $WWWGROUP app
RUN useradd -ms /bin/bash --no-user-group -g $WWWGROUP -u 1337 -G sudo app

# Copy files
COPY start-container.sh /usr/local/bin/start-container.sh
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY . .

# Build the Go app
RUN /usr/local/go/bin/go mod tidy
RUN mkdir -p /var/tmp/go-cache 
RUN mkdir -p /var/www/html/go/pkg/mod

# Set permissions
RUN chown -R app /var/www/html/storage
RUN chmod +x /usr/local/bin/start-container.sh
RUN chown -R app:app /var/www/html
RUN chown app:app /var/tmp/go-cache

EXPOSE 8000/tcp

ENTRYPOINT ["start-container.sh"]