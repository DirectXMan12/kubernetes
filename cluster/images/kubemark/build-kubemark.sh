#!/bin/bash

# Copyright 2015 The Kubernetes Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# We don't need all intermediate steps in the image, especially because
# we clean after ourselves. Thus instead of doing all of this in the Dockerfile
# we use this script.
apt-get update
apt-get install -y wget vim rsync ca-certificates
update-ca-certificates

chmod a+x /kubemark.sh

tar xzf /tmp/kubemark.tar.gz
cp kubernetes/server/bin/kubemark /

rm -rf /tmp/*
apt-get clean -y
apt-get autoremove -y
