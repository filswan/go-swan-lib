module github.com/filswan/go-swan-lib

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/filecoin-project/boost v1.5.0
	github.com/filecoin-project/go-address v1.0.0
	github.com/filecoin-project/go-cbor-util v0.0.1
	github.com/filecoin-project/go-state-types v0.9.8
	github.com/filecoin-project/lotus v1.18.0
	github.com/google/uuid v1.3.0
	github.com/ipfs/go-cid v0.2.0
	github.com/libp2p/go-libp2p v0.22.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/shopspring/decimal v1.3.1
	github.com/sirupsen/logrus v1.8.1
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
