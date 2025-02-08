#!/bin/bash

work_dir=$(cd $(dirname $0); pwd)

if [ "$1" = "online" ]; then
  sed -i "s/configName = \"test_config\"/configName = \"config\"/g" ${work_dir}/config/config.go
else
  sed -i "s/configName = \"config\"/configName = \"test_config\"/g" ${work_dir}/config/config.go
fi