#!/bin/bash
set -xeu

# always returns true so set -e doesn't exit if it is not running.
killall realio-networkd || true
rm -rf $HOME/.realio-network/

mkdir $HOME/.realio-network

# init all three validators
realio-networkd init --chain-id=testing_3301-1 validator --home=$HOME/.realio-network

# create keys for all three validators
mnemonic1="ozone unfold device pave lemon potato omit insect column wise cover hint narrow large provide kidney episode clay notable milk mention dizzy muffin crazy"
mnemonic2="resource slight describe risk lamp delay output lake cannon dismiss account amused"

## realio1ffdjc76wrqtuef5kjhadk9ujhp3xym8x4m3hha
echo $mnemonic1 | realio-networkd keys add validator --recover --keyring-backend=test --home=$HOME/.realio-network
## realio17eu6f96fk6wqpdajgaqwp70jdxj2jjaxp4regz
echo $mnemonic2 | realio-networkd keys add test --recover --keyring-backend=test --home=$HOME/.realio-network

VALIDATOR_CONFIG=$HOME/.realio-network/config
sed -i '' 's/enable = false/enable = true/' $VALIDATOR_CONFIG/app.toml
sed -i '' 's/"voting_period": "172800s"/"voting_period": "30s"/' $VALIDATOR_CONFIG/genesis.json
sed -i '' 's/"unbonding_time": "[^"]*"/"unbonding_time": "30s"/' $VALIDATOR_CONFIG/genesis.json
sed -i '' 's/"expedited_voting_period": "[^"]*"/"expedited_voting_period": "15s"/'  $VALIDATOR_CONFIG/genesis.json

# create validator node with tokens to transfer to the three other nodes
realio-networkd genesis add-genesis-account $(realio-networkd keys show test -a --keyring-backend=test --home=$HOME/.realio-network) 100000000000000000000000ario,100000000000000000000000stake --home=$HOME/.realio-network 
realio-networkd genesis add-genesis-account $(realio-networkd keys show validator -a --keyring-backend=test --home=$HOME/.realio-network) 1000000000000000000000000ario --home=$HOME/.realio-network 
realio-networkd genesis gentx validator 100000000000000000000000ario --keyring-backend=test --home=$HOME/.realio-network --chain-id=testing_3301-1

realio-networkd genesis collect-gentxs --home=$HOME/.realio-network

# # start all three validators/
# onomyd start --home=$HOME/.onomyd/validator1
realio-networkd start --home=$HOME/.realio-network

