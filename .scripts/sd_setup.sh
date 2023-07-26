#!/bin/bash

reset=${2}

for entry in `ls ${1}`; do
	if [ "${entry: -5}" == ".json" ]; then
		echo "Importing ${1}/${entry}..."
		HYPER_SETTINGS=settings/dev/roles/hd-1 hyper sd submit-records ${reset} ${1}/${entry}
		# we only call reset on the first file
		reset=''
	fi
done
