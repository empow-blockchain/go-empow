#!/bin/bash

readonly GRPC_URL="127.0.0.1:30002"
readonly APPROVER_ACCOUNT="EM2ZsSw7RWYC229Z1ib7ujKhken9GFR7dBkTTEbBWMKeLpVas"
readonly APPROVER_ACCOUNT_SECKEY="2yquS3ySrGWPEKywCPzX4RTJugqRh7kJSo5aehsLYPEWkUxBWA39oMrZ7ZxuM4fgyXYs2cPwh5n8aNNpH5x2VyK1"

readonly WITNESS_NUM=1

readonly WITNESS_NAME=(
EM2ZsSi4y3AYqvhbfyzHwDKShtpiNpCQK4WsgTgavup51N2UB
)

iwallet wallet --import ${APPROVER_ACCOUNT} ${APPROVER_ACCOUNT_SECKEY}
for (( i = 0; i < ${WITNESS_NUM}; i++ ))
do
    iwallet -s ${GRPC_URL} --address ${APPROVER_ACCOUNT} call 'vote_producer.empow' 'approveRegister' '["'${WITNESS_NAME[i]}'"]'
done
