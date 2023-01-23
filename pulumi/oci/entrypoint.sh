#!/bin/sh

VALID_CONFIG=true
# Check required ENVs
if [ -z "${PROVIDER+x}" ]; then 
  echo "OPERATION is required"
  VALID_CONFIG=false
fi

if [ -z "${PROJECT_NAME+x}" ]; then 
  echo "PROJECT_NAME ENV is required"
  VALID_CONFIG=false  
fi

if [ -z "${BACKED_URL+x}" ]; then 
  echo "${INTERNAL_OUTPUT} will be used as backed url it will be exported as volume"
  BACKED_URL="file://${INTERNAL_OUTPUT}"
fi

if [ -z "${OUTPUT_FOLDER+x}" ]; then 
  echo "${INTERNAL_OUTPUT} will be used as output folder for connecion resources"
  OUTPUT_FOLDER="${INTERNAL_OUTPUT}"
fi

if [[ "${PROVIDER}" == "aws" ]]; then
  if [ -z "${AWS_ACCESS_KEY_ID+x}" ] || [ -z "${AWS_SECRET_ACCESS_KEY+x}" ] || [ -z "${AWS_DEFAULT_REGION+x}" ]; then 
    echo "AWS ENV for credentials are required"
    VALID_CONFIG=false  
  fi
fi

if [ -z "${PULUMI_CONFIG_PASSPHRASE+x}" ]; then 
  # https://www.pulumi.com/docs/reference/cli/environment-variables/
  PULUMI_CONFIG_PASSPHRASE="passphrase"
fi

if [ "${VALID_CONFIG}" = false ]; then
  echo "Add the required ENVs"
  exit 1
fi

AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION} \
PULUMI_CONFIG_PASSPHRASE=${PULUMI_CONFIG_PASSPHRASE} \
      crc-cloud import \
        --project-name "${PROJECT_NAME}" \
        --backed-url "${BACKED_URL}" \
        --output "${OUTPUT_FOLDER}" \
        --provider "${PROVIDER}" \
        --bundle-url "${BUNDLE_URL}" \
        --bundle-shasumfile-url "${BUNDLE_SHASUMFILE_URL}"
