#!/bin/bash

ls dirty_tmpl | while read line; do
    name=${line%%.*}
    cmd="./dirty-go -tmpl=dirty_tmpl/$line -out=dirty_out/$name.go"
    echo $cmd
    $cmd
done

