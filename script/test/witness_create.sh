#!/bin/bash


readonly GRPC_URL="127.0.0.1:30002"
readonly CREATOR_ACCOUNT="EM2ZsSw7RWYC229Z1ib7ujKhken9GFR7dBkTTEbBWMKeLpVas"
readonly CREATOR_ACCOUNT_SECKEY="2yquS3ySrGWPEKywCPzX4RTJugqRh7kJSo5aehsLYPEWkUxBWA39oMrZ7ZxuM4fgyXYs2cPwh5n8aNNpH5x2VyK1"

readonly WITNESS_NUM=1

readonly WITNESS_NAME=(
EM2ZsSi4y3AYqvhbfyzHwDKShtpiNpCQK4WsgTgavup51N2UB
)

readonly WITNESS_PUBKEY=(
6sNQa7PV2SFzqCBtQUcQYJGGoU7XaB6R4xuCQVXNZe6b
)

iwallet wallet --import ${CREATOR_ACCOUNT} ${CREATOR_ACCOUNT_SECKEY}
for (( i = 0; i < ${WITNESS_NUM}; i++ ))
do
    iwallet -s ${GRPC_URL} --address ${CREATOR_ACCOUNT} --amount_limit "ram:1024|iost:100" account --create ${WITNESS_NAME[i]} --initial_balance 0 --initial_gas_pledge 100 --initial_ram 1024 --owner ${WITNESS_PUBKEY[i]} --active ${WITNESS_PUBKEY[i]}
done
