#!/bin/bash

CHAIN_ID=$(jq -r '.chainId' ./realio-network/key_seed.json)

jq -r .user0 ./realio-network/key_seed.json | ./realio-networkd keys add user0 --recover --keyring-backend test
jq -r .user1 ./realio-network/key_seed.json | ./realio-networkd keys add user1 --recover --keyring-backend test
jq -r .user2 ./realio-network/key_seed.json | ./realio-networkd keys add user2 --recover --keyring-backend test
jq -r .user3 ./realio-network/key_seed.json | ./realio-networkd keys add user3 --recover --keyring-backend test

jq -r .val0 ./realio-network/key_seed.json | ./realio-networkd keys add val0 --recover --keyring-backend test
jq -r .val1 ./realio-network/key_seed.json | ./realio-networkd keys add val1 --recover --keyring-backend test
jq -r .val2 ./realio-network/key_seed.json | ./realio-networkd keys add val2 --recover --keyring-backend test
jq -r .val3 ./realio-network/key_seed.json | ./realio-networkd keys add val3 --recover --keyring-backend test

sleep 10s

./realio-networkd tx bank send user0 $(./realio-networkd keys show user1 --keyring-backend test -a) 1.23rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario
./realio-networkd tx bank send user1 $(./realio-networkd keys show user2 --keyring-backend test -a) 3.45000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario
./realio-networkd tx bank send user2 $(./realio-networkd keys show user3 --keyring-backend test -a) 3.23000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario
./realio-networkd tx bank send user3 $(./realio-networkd keys show user0 --keyring-backend test -a) 2.23rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario

sleep 10s
echo "staking node0"
./realio-networkd tx staking delegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 400rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 320rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 631rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 411rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2

sleep 10s
echo "staking node1"
./realio-networkd tx staking delegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 822000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 613000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 641000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 955000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2

sleep 10s
echo "staking node2"
./realio-networkd tx staking delegate $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 523rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 723rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 426rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 819rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2
sleep 10s
echo "staking node3"
./realio-networkd tx staking delegate $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 548000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 746000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 859000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 960000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2

sleep 10s
echo "unbond node0"
./realio-networkd tx staking unbond $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 26rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 63rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 82rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 37rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2

sleep 10s
echo "unbond node1"
./realio-networkd tx staking unbond $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 26000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 63000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 62000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 72000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2

sleep 10s
echo "unbond node2"
./realio-networkd tx staking unbond $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 63rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 26rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 32rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 47rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2

sleep 10s
echo "unbond node3"
./realio-networkd tx staking unbond $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 56000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 53000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 12000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 32000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2

sleep 10s
echo "redelegate node0"
./realio-networkd tx staking redelegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 23rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 46rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 56rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 22rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2

sleep 10s
echo "redelegate node1"
./realio-networkd tx staking redelegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 61000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 24000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 26000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 43000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2

sleep 10s
echo "unbond node0"
./realio-networkd tx staking unbond $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 21rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 23rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 32rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 57rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2
sleep 10s
echo "unbond node1"
./realio-networkd tx staking unbond $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 46000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 53000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 62000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 62000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2
sleep 10s
echo "unbond node2"
./realio-networkd tx staking unbond $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 73rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 46rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 72rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 37rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2
sleep 10s
echo "unbond node3"
./realio-networkd tx staking unbond $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 36000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 53000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 72000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking unbond $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 72000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2

sleep 10s
echo "delegate node0"
./realio-networkd tx staking delegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 400rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 320rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 631rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) 411rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2
sleep 10s
echo "delegate node1"
./realio-networkd tx staking delegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 822000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 613000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 641000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) 955000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2
sleep 10s
echo "delegate node2"
./realio-networkd tx staking delegate $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 523rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 723rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 426rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 819rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2
sleep 10s
echo "delegate node3"
./realio-networkd tx staking delegate $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 548000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 746000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 859000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking delegate $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 960000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2
sleep 10s

./realio-networkd tx staking redelegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 23rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 46rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 56rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val0 --keyring-backend test -a --bech val) $(./realio-networkd keys show val2 --keyring-backend test -a --bech val) 22rio --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2
sleep 10s

./realio-networkd tx staking redelegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 61000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user0 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 24000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user1 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 26000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user2 --gas auto --gas-adjustment 1.2
./realio-networkd tx staking redelegate $(./realio-networkd keys show val1 --keyring-backend test -a --bech val) $(./realio-networkd keys show val3 --keyring-backend test -a --bech val) 43000000000000arst --chain-id $CHAIN_ID --keyring-backend test -y --gas-prices 1000000000ario --from user3 --gas auto --gas-adjustment 1.2
