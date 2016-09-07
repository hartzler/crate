#!/usr/bin/env bash
GOTAR=go1.7.linux-amd64.tar.gz
if [ ! -d /usr/local/go ]; then
  if [ ! -f "$GOTAR" ]; then
    wget -q https://storage.googleapis.com/golang/$GOTAR
  fi
  tar -xzf $GOTAR &>/dev/null
  mv go /usr/local/
  which go || echo "PATH=$PATH:/usr/local/go/bin" >> /etc/profile
fi
