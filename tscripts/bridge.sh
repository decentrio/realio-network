
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
realio-networkd tx bridge bridge-in 1000ario realio14ze2f3e55gy079jf0ca5qkamjmffwmuklh6xzj --from ${KEYS[0]} --keyring-backend $KEYRING --home "$HOMEDIR" -y --chain-id ${CHAINID} --fees 10922000030ario 