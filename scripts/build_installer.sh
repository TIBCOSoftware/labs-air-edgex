#!/bin/bash

installer_target_path="dist"

if [ -d $installer_target_path ]; then
  rm -rf $installer_target_path
fi
mkdir -p $installer_target_path

cp -r scripts/linux $installer_target_path
cp -r scripts/arm64 $installer_target_path