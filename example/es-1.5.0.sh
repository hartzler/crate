#!/bin/sh
./bin/crate new --cargo=https://s3.amazonaws.com/armada-crates/cargo/26cf93db776362bebc35cff379a65fa17aacf5d429785a080e3f3a73701caa71.cargo --cargo=https://s3.amazonaws.com/armada-crates/cargo/b9a3956b41431a9551da413a36a2322458fa097fdad11d98367710e52ae0b3bd.cargo --address=10.4.0.1/16 $1
./bin/crate run --pid=es --env=JAVA_HOME=/opt/jdk --env=ES_MIN_MEM=128m --env=ES_MAX_MEM=128m $1 /opt/es/bin/elasticsearch
