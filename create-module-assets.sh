#!/bin/bash
set -o xtrace

cd -P -- "$(dirname "${0}")" || exit
readonly SCRIPT_ROOT=$(pwd)

# test for required dependencies
dependencies=(docker kubectl minikube)
for dependency in "${dependencies[@]}"
do
  if ! command -v "${dependency}" >/dev/null
  then
    printf 'Required command "%s" not found on PATH.\n' "${dependency}" >&2
    exit 1
  fi
done

# use the docker host in minikube
eval "$(minikube docker-env --shell bash)"

# read the image tag specified in the .section-proxy.yaml file
read_pod_template_file_image_tag () {
  local resource_file=$1
  grep -Eo --no-filename -e '\bgcr\.io/section-?io/[^:]+:\w[A-Za-z0-9_.-]*' "${resource_file}" || true
}

namespaced_image_tag=$(read_pod_template_file_image_tag "$(find ./ -name .section-proxy.yaml -print0)")

# build the example-module image
docker build --tag="${namespaced_image_tag}" example

# processes .section-proxy.yaml files to generate yaml files used to run the module on the Section platform
readonly processor_version=latest

workdir=$(mktemp -d) || exit
mkdir "${workdir}/src" "${workdir}/out"

src_files=()
while IFS= read -r -d $'\0' proxy_file
do
  proxy_name=$(basename "$(dirname "${proxy_file}")")

  mkdir "${workdir}/src/${proxy_name}"
  cp "${proxy_file}" "${workdir}/src/${proxy_name}/"

  src_files+=("/w/src/${proxy_name}/.section-proxy.yaml")
done < <(find ./ -name .section-proxy.yaml -print0)

# create a volume to mount the processed files
vid=$(docker volume create) || exit

tar -C "${workdir}" --mode=a+rwX -cf- . |
  docker container run -i --rm -v "${vid}:/w" -w /w busybox:1.28 tar -xf-

exit_code=0
docker container run --rm \
  -v "${vid}:/w" \
  -w /w/out/ \
  "gcr.io/section-io/package-processor:${processor_version}" \
    "${src_files[@]}" || exit_code=$?

docker container run --rm -v "${vid}:/w" -w /w/out busybox:1.28 tar -cf- . |
  tar -xmf-

docker volume rm "${vid:?vid undefined}" > /dev/null
rm -rf "${workdir:?workdir undefined}"

docker build -t package-sync-override "${SCRIPT_ROOT}" || exit

run_image=package-sync-override

# run the image which copies the files within to the mounted volume
docker run --rm \
  --volume /data/section-io/package-sync/override/:/override/:rw \
  "${run_image}" || exit

# delete the package-sync pod to force a restart and re-export
# to avoid the 10 minute wait.
kubectl delete pod -n section-shared -l name=package-sync --now

exit "${exit_code}"
