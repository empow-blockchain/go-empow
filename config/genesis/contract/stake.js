const COUNT_PACKAGE_PREFIX = 'c_'                               // c_address
const PACKAGE_INFO_PREFIX = 'p_'                                // p_address_count
const TOTAL_STAKE_AMOUNT = 'totalStakeAmount'
const INTEREST_ARRAY = 'interestArray'
const REST_AMOUNT = "restAmount"
const MAXIMUM_PERCENT_PER_DAY = new Float64("0.000833333333")   // 30%/year
const BLOCK_NUMBER_PER_DAY = 172800
const THREE_DAY_NANO = 259200

class Stake {
    init() {
        storage.put(TOTAL_STAKE_AMOUNT, "0")
        storage.put(REST_AMOUNT, "0")
        storage.put(INTEREST_ARRAY, JSON.stringify([]))
    }

    initAdmin (adminAddress) {
        const bn = block.number;
        if(bn !== 0) {
            throw new Error("init out of genesis block")
        }
        storage.put("adminAddress", adminAddress);
    }

    can_update(data) {
        const admin = storage.get("adminAddress");
        this._requireAuth(admin, "active");
        return true;
    }

    _updateTotalStakeAmount(amount) {
        let totalStakeAmount = storage.get(TOTAL_STAKE_AMOUNT)
        totalStakeAmount = new Float64(totalStakeAmount)
        storage.put(TOTAL_STAKE_AMOUNT, totalStakeAmount.plus(amount).toString())
    }

    _getAddressCountPackage(address) {
        let countPackage = storage.get(COUNT_PACKAGE_PREFIX + address)
        if(!countPackage) return new Int64("0")
        else return new Int64(countPackage)
    }

    _updateAddressCountPackage(address, amount) {
        let countPackage = storage.get(COUNT_PACKAGE_PREFIX + address)
        if(!countPackage) {
            storage.put(COUNT_PACKAGE_PREFIX + address, "1")
        } else {
            countPackage = new Int64(countPackage)
            storage.put(COUNT_PACKAGE_PREFIX + address, countPackage.plus(amount).toString())
        }
    }

    _updateRestAmount(amount) {
        let restAmount = storage.get(REST_AMOUNT)
        storage.put(REST_AMOUNT, amount.plus(restAmount).toString())
    }

    _updatePackageInfo(address, packageNumber, data) {
        const prefix = PACKAGE_INFO_PREFIX + address + "_" + packageNumber
        storage.put(prefix, JSON.stringify(data))
    }

    _addNewPackage(address, data) {
        let countPackage = storage.get(COUNT_PACKAGE_PREFIX + address)

        if(countPackage) {
            this._updatePackageInfo(address, countPackage, data)
        } else {
            this._updatePackageInfo(address, "0", data)
        }

        this._updateAddressCountPackage(address, 1)
    }

    _fixAmount(amount) {
        amount = new Float64(new Float64(amount).toFixed(8));
        if (amount.lte("0")) {
            throw new Error("amount must be positive");
        }
        return amount;
    }

    stake(address, amount) {
        blockchain.requireAuth(address, "active")
        // send EM to stake.empow
        blockchain.callWithAuth("token.empow", "transfer", ["em", address, blockchain.contractName(), amount, "stake EM"])
        // create package info
        let packageInfo = {
            lastBlockWithdraw: block.number,
            unstake: false,
            amount: amount
        }
        this._addNewPackage(address, packageInfo)
        this._updateTotalStakeAmount(amount)

        blockchain.receipt(JSON.stringify([address, amount]))
    }

    topup(amount) {
        blockchain.requireAuth("issue.empow", "active")
        amount = this._fixAmount(amount);

        blockchain.deposit("issue.empow", amount.toFixed(), "");

        // calc interest per 1 EM stake
        const totalStakeAmount = storage.get(TOTAL_STAKE_AMOUNT)
        let interest = amount
        if(totalStakeAmount !== "0") {
            interest = amount.div(totalStakeAmount)
        }

        // insert interest to array
        let interestArray = JSON.parse(storage.get(INTEREST_ARRAY))
        interestArray.push(interest.toFixed(8))
        if(interestArray.length > 30) {
            interestArray.shift()
        }
        storage.put(INTEREST_ARRAY, JSON.stringify(interestArray))

        return true
    }

    withdraw(address, packageNumber) {
        blockchain.requireAuth(address, "active")
        // check package exist
        const packageInfoString = storage.get(PACKAGE_INFO_PREFIX + address + "_" + packageNumber)
        if(!packageInfoString) {
            throw new Error("package not exist > " + packageNumber)
        }
        let packageInfo = JSON.parse(packageInfoString)
        // check package unstake
        if(packageInfo.unstake) {
            throw new Error("package has been unstake > " + packageNumber)
        }
        // calc interest
        const bn = block.number
        const totalDayStake = Math.floor((bn - packageInfo.lastBlockWithdraw) / BLOCK_NUMBER_PER_DAY)
        const interestArray = JSON.parse(storage.get(INTEREST_ARRAY))
        let amountCanWithdraw = new Float64("0")
        const packageAmount = new Float64(packageInfo.amount)
        
        if(totalDayStake > interestArray.length) {
            totalDayStake = interestArray.length
        }

        if(totalDayStake <= 0) {
            throw new Error("package stake less than 1 day > " + packageNumber)
        }

        const maxAmountCanWithdraw = packageAmount.multi(MAXIMUM_PERCENT_PER_DAY.multi(totalDayStake))

        for(i = interestArray.length - 1; i >= interestArray.length - totalDayStake; i--) {
            amountCanWithdraw.plus(packageAmount.multi(interestArray[i]))
        }

        if(amountCanWithdraw.gt(maxAmountCanWithdraw)) {
            this._updateRestAmount(amountCanWithdraw.minus(maxAmountCanWithdraw))
            amountCanWithdraw = maxAmountCanWithdraw
        }

        blockchain.withdraw(address, amountCanWithdraw.toFixed(8), "withdraw stake")

        packageInfo.lastBlockWithdraw = bn

        this._updatePackageInfo(address, packageNumber, packageInfo)

        blockchain.receipt(JSON.stringify([address, amountCanWithdraw.toFixed(8)]))
    }

    unstake (address, packageNumber) {
        blockchain.requireAuth(address, "active")

        const packageInfoString = storage.get(PACKAGE_INFO_PREFIX + address + "_" + packageNumber)
        if(!packageInfoString) {
            throw new Error("package not exist > " + packageNumber)
        }
        let packageInfo = JSON.parse(packageInfoString)

        if(packageInfo.unstake) {
            throw new Error("package has been unstake > " + packageNumber)
        }

        const freezeTime = tx.time + THREE_DAY_NANO*1e9
        const stakeAmount = new Float64(packageInfo.amount)
        blockchain.callWithAuth("token.empow", "transferFreeze", ["em", "stake.empow", address, packageInfo.amount, freezeTime, "unstake"])

        this._updateTotalStakeAmount(stakeAmount.negated())
        packageInfo.unstake = true
        this._updatePackageInfo(address, packageNumber, packageInfo)

        blockchain.receipt(JSON.stringify([address, packageInfo.amount]))
    }
}

module.exports = Stake