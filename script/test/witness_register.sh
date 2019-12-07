#!/bin/bash

readonly GRPC_URL="127.0.0.1:30002"
readonly WITNESS_NUM=1

readonly WITNESS_NAME=(
EM2ZsSi4y3AYqvhbfyzHwDKShtpiNpCQK4WsgTgavup51N2UB
)

readonly WITNESS_PUBKEY=(
6sNQa7PV2SFzqCBtQUcQYJGGoU7XaB6R4xuCQVXNZe6b
)

readonly WITNESS_SECKEY=(
1rANSfcRzr4HkhbUFZ7L1Zp69JZZHiDDq5v7dNSbbEqeU4jxy3fszV4HGiaLQEyqVpS1dKT9g7zCVRxBVzuiUzB
)

for (( i = 0; i < ${WITNESS_NUM}; i++ ))
do
    iwallet -s ${GRPC_URL} wallet --import ${WITNESS_NAME[i]} ${WITNESS_SECKEY[i]}
    iwallet -s ${GRPC_URL} --address ${WITNESS_NAME[i]} call "vote_producer.empow" "applyRegister" '["'${WITNESS_NAME[i]}'","'${WITNESS_PUBKEY[i]}'","location","url","",true]'
done
