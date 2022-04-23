#!/bin/bash

network_type=${1:online}

installer_target_path="dist"

if [ -d $installer_target_path ]; then
  rm -rf $installer_target_path
fi
mkdir -p $installer_target_path

# Offline artifacts
if [[ "${network_type}" == "offline" ]];
then
  pushd scripts/linux/basicdemo || exit 1
  ./export.sh || exit 1
  popd > /dev/null || exit 1
  
  pushd scripts/arm64/basicdemo || exit 1
  ./export.sh || exit 1
  popd > /dev/null || exit 1

fi

cp -r scripts/linux $installer_target_path
cp -r scripts/arm64 $installer_target_path