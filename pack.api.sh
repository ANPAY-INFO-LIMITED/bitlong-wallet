#!/bin/bash

start_time=$(date +%s)
currentPath=$(pwd)
currentFolderName="${currentPath##*/}"
specificString="wallet"

if [ "$currentFolderName" == "$specificString" ]; then
    echo "gomobile is in progress, please wait..."
    # shellcheck disable=SC2164
    cd api
    gomobile bind -target android -tags "litd_no_ui litd autopilotrpc signrpc walletrpc chainrpc invoicesrpc watchtowerrpc neutrinorpc peersrpc btlapi"
    # shellcheck disable=SC2103
    cd ..
    end_time=$(date +%s)
    time_taken=$((end_time - start_time))
    echo "Time cost: $time_taken seconds."
else
    echo "Wrong current directory, please run script in wallet."
    # shellcheck disable=SC2162
    read -p "Press enter to continue"
fi
