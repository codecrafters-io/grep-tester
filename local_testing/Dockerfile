FROM golang:1.24

# Install required packages including Python
RUN apt-get update && apt-get install -y \
    # Required for make test
    make \
    # Required for testing
    git \
    # Python (for running scripts)
    python3 \
    && rm -rf /var/lib/apt/lists/*

# Create symlink for python command (optional)
RUN ln -s /usr/bin/python3 /usr/bin/python

# Set working directory
WORKDIR /app

# Starting from Go 1.20, the go standard library is no longer compiled.
# Setting GODEBUG to "installgoroot=all" restores the old behavior
RUN GODEBUG="installgoroot=all" go install std

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Default command
CMD ["/bin/bash"]