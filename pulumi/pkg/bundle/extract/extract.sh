#!/bin/sh

# Function to downlad the component from 
# the url $1 in case of failure it retries
download () {
    local download_url="${1}"
    local download_result=1
    while [[ ${download_result} -ne 0 ]]
    do
        curl -s --insecure -LO -C - ${download_url}
        download_result=$?
    done
}

# Checks the shasum256 of an asset
# $1 asset name 
# $2 shasum filename
# Return 1 if not valid, 0 if valid
check_shasum() {
    cat ${2} | grep ${1} | sha256sum -c -
    return ${?}
} 

####################
###### MAIN ########
####################

# set -xuo 

# Get asset names
BUNDLE_NAME=${BUNDLE_DOWNLOAD_URL##*/}
SHASUMFILE_NAME=${SHASUMFILE_DOWNLOAD_URL##*/}

# Download assets and check 
curl -s --insecure -LO "${SHASUMFILE_DOWNLOAD_URL}"
download ${BUNDLE_DOWNLOAD_URL}
check_shasum ${BUNDLE_NAME} ${SHASUMFILE_NAME} 
if [[ ${?} -ne 0 ]]; then
    echo "Error with downloading ${BUNDLE_DOWNLOAD_URL}"
    exit 1
fi

# Uncompress
mv -v ${BUNDLE_NAME} "${BUNDLE_NAME%'crcbundle'}zst"
unzstd "${BUNDLE_NAME%'crcbundle'}zst" -o bundle.tar --quiet
rm "${BUNDLE_NAME%'crcbundle'}zst"
mkdir -p bundle
tar -vxf bundle.tar -C bundle --strip-components=1
rm bundle.tar

# Export image to raw
qemu-img convert bundle/crc.qcow2 disk.raw
# Export booting private key  
cp bundle/id_ecdsa_crc id_ecdsa
rm -rf bundle

# TODO keep it for checking best option to import images 
# tar --format=oldgnu -Sczf /tmp/crc.tar.gz disk.raw
