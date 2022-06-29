#!/bin/bash

network_type=${1:?}
os_type=${2:?}
arch_type=${3:?}
release_version=${4:?}

build_offline(){
  # Offline artifacts
  pushd "./installers/community/${arch_type}/edge" || exit 1
  ./export.sh || exit 1
  popd > /dev/null || exit 1
}

replace_release_version(){
  # Replace release version
  sed -i  "s/LABS_AIR_VERSION=GENERATED_VERSION/LABS_AIR_VERSION=${release_version}/" "${installer_target_path}/.env"
}

installer_target_path="dist"

if [ -d $installer_target_path ]; then
  rm -rf $installer_target_path
fi
mkdir -p $installer_target_path

# Offline artifacts
if [[ "${network_type}" == "offline" ]];
then
  build_offline || exit 1
fi

cp -r "./installers/community/" ${installer_target_path} || exit 1

replace_release_version

if [[ "${os_type}" != windows ]];
  then
    cp scripts/start.sh $installer_target_path || exit 1
    cp scripts/stop.sh $installer_target_path || exit 1
  fi

