#!/bin/bash

set -e

for name in $@
do
  echo -e "\e[0;32mProcessing book ${name} ...\e[0m"
  ../go/bin/tex ${name}.txt
  for i in 1 2 3
  do
    xelatex ${name}.tex
  done
done
