#!/bin/bash

# Get the latest version number of protoc (removing the 'v' prefix)
PROTOC_VERSION=$(curl -s https://api.github.com/repos/protocolbuffers/protobuf/releases/latest | grep 'tag_name' | cut -d\" -f4 | cut -c 2-)

# Define the protoc zip file name
ZIP_FILE="protoc-$PROTOC_VERSION-linux-x86_64.zip"

# Download the zip file
curl -OL "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/$ZIP_FILE"

# Unzip the file, deleting the zip file, ignore existing files
unzip -o $ZIP_FILE -d ./protoc

# Move protoc to /usr/local/bin/
sudo mv ./protoc/bin/* /usr/local/bin/

# Move protoc include files /usr/local/include/
sudo mv ./protoc/include/* /usr/local/include/

# Cleanup
rm -rf ./protoc
rm $ZIP_FILE

# Verify installation
protoc --version
