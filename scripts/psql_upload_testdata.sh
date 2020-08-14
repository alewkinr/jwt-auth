#!/bin/sh
connection_string=$1
dir_sql_scripts=$2
for file in $dir_sql_scripts
do
  if [ -f "$file" ]
  then
    psql "$connection_string" -f "$file"
  fi
  echo "complete upload testdata"
done