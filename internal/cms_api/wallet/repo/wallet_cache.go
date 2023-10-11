package repo

var WalletCache = new(walletCache)

type walletCache struct{}

const UserWalletKey = "wallet:%s"
const RedpackSingleKey = "redpack:single:%d"
const RedpackGroupKey = "redpack:group:%d"
const SignReward = "sign:%s"
