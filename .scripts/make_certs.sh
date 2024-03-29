#!/bin/bash
# This script generates test & development certificates. Not for production use!

# You can add an entry to this list to generate a certificate for a given
# operator. The name of the operator will be added as the common name as well
# as the subject alternative name (SAN), which is required for some newer
# TLS libraries.
declare -a certs=("op-1" "op-2" "hd-1" "hd-2" "ls-1" "internal-server" "public-proxy-1.ga" "private-proxy-1.ga" "sd-1" "demo-app" "quic-1" "quic-2")
declare -A groups=(["hd-1"]="sd-admin")

O="IRIS"
ST="Berlin"
L="Berlin"
C="DE"
OU="IT"
CN="Testing-Development"
# using less than 1024 here will result in a TLS handshake failure in Go
# using less than 2048 will cause e.g. 'curl' to complain that the ciper is too weak
LEN="2048"

# Please note that we add a ".local" name to the wildcard subject alternative names as
# second-level wildcards (e.g. "*.internal-server") will not work. Probably a good security
# measure as one could otherwise register a wildcard like "*.com"

DIRECTORY_FILE=../directory/002_certificates.json

touch index.txt
echo 1000 > serial 

openssl genrsa -out root.key ${LEN}
openssl req -x509 -new -nodes -key root.key -sha256 -days 1024 -out root.crt -subj "/C=${C}/ST=${ST}/L=${L}/O=${O}/OU=${OU}/CN=${CN}"

openssl genrsa -out intermediate.key ${LEN}
openssl rsa -in "intermediate.key" -pubout -out "intermediate.pub";
openssl req -new -sha256 -key "intermediate.key" -subj "/C=${C}/ST=${ST}/L=${L}/O=${O}/OU=${OU}/CN=intermediate" -out "intermediate.csr";
# openssl x509 -req -in "intermediate.csr" -CA root.crt -CAkey root.key -CAcreateserial -out "intermediate.crt" -extensions SAN -days 500 -sha256;
openssl ca -extensions v3_intermediate_ca \
	-days 300 -batch -notext -md sha256 \
	-in intermediate.csr \
	-config ../../../.scripts/root.conf \
	-out intermediate.crt

echo -n "{\"records\": [" > $DIRECTORY_FILE

for cert in "${certs[@]}"
do

	if [[ -n $LAST ]]; then
		echo -n "," >> $DIRECTORY_FILE
	fi;

	openssl genrsa -out "${cert}.key" ${LEN};
	openssl rsa -in "${cert}.key" -pubout -out "${cert}.pub";
	openssl req -new -sha256 -key "${cert}.key" -subj "/C=${C}/ST=${ST}/L=${L}/O=${O}/OU=${OU}/CN=${cert}" -addext "subjectAltName = DNS:${cert},DNS:*.${cert}.local" -out "${cert}.csr";
	openssl x509 -req -in "${cert}.csr" -CA intermediate.crt -CAkey intermediate.key -CAcreateserial -out "${cert}.crt" -extensions SAN -extfile <(printf "[SAN]\nsubjectAltName = DNS:${cert},DNS:*.${cert}.local") -days 500 -sha256;
	# we add the intermediate certificate
	cat intermediate.crt >> "${cert}.crt"

	# we generate alternative certificates for testing the pinning (we won't add them to the directory)
	openssl genrsa -out "${cert}-alt.key" ${LEN};
	openssl rsa -in "${cert}-alt.key" -pubout -out "${cert}-alt.pub";
	openssl req -new -sha256 -key "${cert}-alt.key" -subj "/C=${C}/ST=${ST}/L=${L}/O=${O}/OU=${OU}/CN=${cert}" -addext "subjectAltName = DNS:${cert},DNS:*.${cert}.local" -out "${cert}-alt.csr";
	openssl x509 -req -in "${cert}-alt.csr" -CA intermediate.crt -CAkey intermediate.key -CAcreateserial -out "${cert}-alt.crt" -extensions SAN -extfile <(printf "[SAN]\nsubjectAltName = DNS:${cert},DNS:*.${cert}.local") -days 500 -sha256;
	# we add the intermediate certificate
	cat intermediate.crt >> "${cert}-alt.crt"

	# the signing certificates use ECDSA and are for signing service directory entries
	openssl ecparam -genkey -name prime256v1 -noout -out "${cert}-sign.key";
	openssl ec -in "${cert}-sign.key" -pubout -out "${cert}-sign.pub";
	openssl req -new -sha256 -key "${cert}-sign.key" -subj "/C=${C}/ST=${ST}/L=${L}/O=${O}/OU=${OU}/CN=${cert}" -addext "keyUsage=digitalSignature" -addext "subjectAltName = URI:iris-name://${cert},URI:iris-group://${groups[${cert}]},DNS:${cert}"  -out "${cert}-sign.csr";
	openssl x509 -req -in "${cert}-sign.csr" -CA intermediate.crt -CAkey intermediate.key -CAcreateserial -out "${cert}-sign.crt"  -extensions SANKey -extfile <(printf "[SANKey]\nsubjectAltName = URI:iris-name://${cert},URI:iris-group://${groups[${cert}]},DNS:${cert}\nkeyUsage = digitalSignature") -days 500 -sha256;
	# we add the intermediate certificate
	cat intermediate.crt >> "${cert}-sign.crt"

	# we generate alternative certificates for testing the pinning (we won't add them to the directory)
	openssl ecparam -genkey -name prime256v1 -noout -out "${cert}-sign-alt.key";
	openssl ec -in "${cert}-sign-alt.key" -pubout -out "${cert}-sign-alt.pub";
	openssl req -new -sha256 -key "${cert}-sign-alt.key" -subj "/C=${C}/ST=${ST}/L=${L}/O=${O}/OU=${OU}/CN=${cert}" -addext "keyUsage=digitalSignature" -addext "subjectAltName = URI:iris-name://${cert},URI:iris-group://${groups[${cert}]},DNS:${cert}"  -out "${cert}-sign-alt.csr";
	openssl x509 -req -in "${cert}-sign-alt.csr" -CA intermediate.crt -CAkey intermediate.key -CAcreateserial -out "${cert}-sign-alt.crt"  -extensions SANKey -extfile <(printf "[SANKey]\nsubjectAltName = URI:iris-name://${cert},URI:iris-group://${groups[${cert}]},DNS:${cert}\nkeyUsage = digitalSignature") -days 500 -sha256;
	# we add the intermediate certificate
	cat intermediate.crt >> "${cert}-sign-alt.crt"

	FINGERPRINT_SIGNING=`openssl x509 -noout -fingerprint -sha256 -inform pem -in "${cert}-sign.crt" | sed -e 's/://g' | sed -r 's/.*=(.*)$/\1/g' | awk '{print tolower($0)}'`
	FINGERPRINT_ENCRYPTION=`openssl x509 -noout -fingerprint -sha256 -inform pem -in "${cert}.crt" | sed -e 's/://g' | sed -r 's/.*=(.*)$/\1/g' | awk '{print tolower($0)}'`
	echo -n "{\"section\": \"certificates\", \"created_at\":\"`date --rfc-3339=seconds | sed 's/ /T/'`\", \"name\": \"${cert}\", \"data\": [{\"fingerprint\": \"${FINGERPRINT_SIGNING}\", \"key_usage\": \"signing\"},{\"fingerprint\": \"${FINGERPRINT_ENCRYPTION}\", \"key_usage\": \"encryption\"}]}" >> $DIRECTORY_FILE
	LAST=1
done

echo -n "]}" >> $DIRECTORY_FILE
