
KEYS[0]="dev0"
KEYS[1]="dev1"
KEYS[2]="dev2"
CHAINID="realionetworklocal_7777-1"
MONIKER="realionetworklocal"

KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# Set dedicated home directory for the realio-networkd instance
HOMEDIR="$HOME/.realio-network"
# to trace evm
#TRACE="--trace"
TRACE=""
realio-networkd start --pruning=nothing "$TRACE" --log_level $LOGLEVEL --minimum-gas-prices=0.00001ario --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --home "$HOMEDIR"