. scripts/tools.sh

image_name=$1
image_tag=$2
image_url=$3
component_name=$4
target_name=$5

echo "in build air edgex component"
echo ${image_url}
echo ${component_name}


pushd ./${component_name} > /dev/null

build_image $image_name "local_image_tag" $image_url "Dockerfile" $target_name

tag_image $image_name "local_image_tag" $image_name $image_tag $image_url