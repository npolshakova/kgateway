#!/usr/bin/env bash

# Copyright Istio Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

# Copy credentials from mountpoints using su-exec
uid=$(id -u)
gid=$(id -g)

shopt -s dotglob

if [[ -d /config ]]; then
  # Make a copy of the host's config secrets. Do not copy docker sockets.
  su-exec 0:0 rsync -a --exclude=docker*.sock --exclude=/config/.config/gcloud/logs/* /config/ /config-copy/ || true

  # Set the ownership of the host's config secrets to that of the container
  su-exec 0:0 chown -R "${uid}":"${gid}" /config-copy || true

  # Permit only the UID:GID to read the copy of the host's config secrets
  chmod -R 700 /config-copy || true

  # If docker_for_mac plaintext-passwords.json exists, import it into config.json
  if [[ -f /config-copy/.docker/plaintext-passwords.json ]]; then
    auth_value=$(jq -r '.auths."https://index.docker.io/v1/".auth' /config-copy/.docker/plaintext-passwords.json)
    if [[ "${auth_value}" == "null" ]]; then
      echo "Missing docker credentials."
    fi
    encode_value=$(echo "${auth_value}" | base64 --decode | base64)
    jq --arg auth "${encode_value}" '.auths."https://index.docker.io/v1/".auth=$auth' /config-copy/.docker/config.json > /config-copy/.docker/config-tmp.json
    jq 'del(.credsStore)' /config-copy/.docker/config-tmp.json > /config-copy/.docker/config.json
  fi
fi

# Add user based upon passed UID. Skip if run as root.
if [[ "${uid}" -ne 0 ]]; then
  su-exec 0:0 useradd --uid "${uid}" --system user || true
fi

# Set ownership of /home to UID:GID
su-exec 0:0 chown "${uid}":"${gid}" /home || true

# Copy the config secrets without changing permissions nor ownership
if [[ -d /config-copy ]]; then
  cp -R /config-copy/* /home/ 2>/dev/null || true
fi

exec "$@"


