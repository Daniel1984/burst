# Start from golang v1.15.5 base image
FROM golang:1.15.5

# create a working directory
WORKDIR /app

# we will be expecting to get API_PORT as arguments
ARG API_PORT

# Fetch dependencies on separate layer as they are less likely to
# change on every build and will therefore be cached for speeding
# up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# copy source from the host to the working directory inside
# the container
COPY . .

# This container exposes API_PORT from .env to the outside world
EXPOSE ${API_PORT}
