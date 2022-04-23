#!/bin/bash

network_type=${1:?}
os_type=${2:?}
arch_type=${3:?}

build_offline(){
  # Offline artifacts
  if [[ "${arch_type}" == "amd64" ]]; then
      pushd "scripts/linux/basicdemo" || exit 1
      ./export.sh || exit 1
      popd > /dev/null || exit 1
  elif [[ "${arch_type}" == "arm64" ]]; then
      pushd "scripts/arm64/basicdemo" || exit 1
      ./export.sh || exit 1
      popd > /dev/null || exit 1
  fi
}


installer_target_path="dist"

if [ -d $installer_target_path ]; then
  rm -rf $installer_target_path
fi
mkdir -p $installer_target_path

# Offline artifacts
if [[ "${network_type}" == "offline" ]];
then
  build_offline ${arch_type}|| exit 1
fi

if [[ "${arch_type}" == "amd64" ]]; then
      cp -r scripts/linux $installer_target_path|| exit 1
elif [[ "${arch_type}" == "arm64" ]]; then
      cp -r scripts/arm64 $installer_target_path|| exit 1
fi

if [[ "${os_type}" != windows ]];
  then
    cp scripts/start.sh $installer_target_path || exit 1
    cp scripts/stop.sh $installer_target_path || exit 1
  fi

