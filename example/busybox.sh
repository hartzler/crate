#!/bin/sh
./bin/crate new --cargo=https://s3.amazonaws.com/armada-crates/cargo/09f2059bcc00314b072d35b436b462bdc12d3cf7d3ca88e6467403764b1a8b9e.cargo --address=10.4.1.1/16 $1
./bin/crate shell $1
