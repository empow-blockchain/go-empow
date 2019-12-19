const COUNT_PACKAGE_PREFIX = 'c_'                               // c_address
const PACKAGE_INFO_PREFIX = 'p_'                                // p_address_count
const INTEREST_PREFIX = 'i_'                                    // i_day
const TOTAL_STAKE_AMOUNT = 'totalStakeAmount'
const REST_AMOUNT = "restAmount"
const MAXIMUM_PERCENT_PER_DAY = new Float64("0.000833333333")   // 30%/year
const BLOCK_NUMBER_PER_DAY = 172800
const THREE_DAY_NANO = 259200*1e9

class Stake {
    init() {
        storage.put(TOTAL_STAKE_AMOUNT, "0")
        storage.put(REST_AMOUNT, "0")
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
        let packageId = 0

        if(countPackage) {
            packageId = countPackage
        }
        
        this._updatePackageInfo(address, packageId, data)
        this._updateAddressCountPackage(address, 1)
        
        return packageId
    }

    _fixAmount(amount) {
        amount = new Float64(new Float64(amount).toFixed(8));
        if (amount.lte("0")) {
            throw new Error("amount must be positive");
        }
        return amount;
    }

    _requireAuth(address, permission) {
        const ret = blockchain.requireAuth(address, permission);
        if (ret !== true) {
            throw new Error("require auth > " + address)
        }
    }

    _calcInterest(stakeAmount, totalDayStake) {
        const bn = block.number
        const currentDay = Math.floor(bn / BLOCK_NUMBER_PER_DAY)
        const stopWithdrawDay = currentDay - totalDayStake

        let amountCanWithdraw = new Float64("0")
        const packageAmount = new Float64(stakeAmount)
        
        for(let i = currentDay; i >= stopWithdrawDay; i--) {
            let interest = storage.get(INTEREST_PREFIX + i)
            if(!interest) continue;
            amountCanWithdraw = amountCanWithdraw.plus(packageAmount.multi(interest))
        }

        return amountCanWithdraw
    }

    _calcMaxInterest(stakeAmount, totalDayStake) {
        const packageAmount = new Float64(stakeAmount)
        return packageAmount.multi(MAXIMUM_PERCENT_PER_DAY).multi(totalDayStake)
    }

    stake(address, amount) {
        this._requireAuth(address, "active")
        // send EM to stake.empow
        blockchain.callWithAuth("token.empow", "transfer", ["em", address, blockchain.contractName(), amount, "stake EM"])
        // create package info
        let packageInfo = {
            lastBlockWithdraw: block.number,
            unstake: false,
            amount: amount
        }
        const packageId = this._addNewPackage(address, packageInfo)
        this._updateTotalStakeAmount(amount)

        blockchain.receipt(JSON.stringify([address, amount, packageId]))
    }

    topup(amount) {
        this._requireAuth("base.empow", "active")
        amount = this._fixAmount(amount);

        // calc interest per 1 EM stake
        const totalStakeAmount = storage.get(TOTAL_STAKE_AMOUNT)
        let interest = amount
        if(totalStakeAmount !== "0") {
            interest = amount.div(totalStakeAmount)
        }

        // insert interest to array
        const bn = block.number
        const currentDay = Math.floor(bn / BLOCK_NUMBER_PER_DAY)
        storage.put(INTEREST_PREFIX + currentDay, interest.toFixed(8))

        return true
    }

    withdraw(address, packageNumber) {
        this._requireAuth(address, "active")
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

        if(totalDayStake <= 0) {
            throw new Error("package withdraw less than 1 day > " + packageNumber)
        }

        let amountCanWithdraw = this._calcInterest(packageInfo.amount, totalDayStake)
        const maxAmountCanWithdraw = this._calcMaxInterest(packageInfo.amount, totalDayStake)

        if(amountCanWithdraw.gt(maxAmountCanWithdraw)) {
            this._updateRestAmount(amountCanWithdraw.minus(maxAmountCanWithdraw))
            amountCanWithdraw = maxAmountCanWithdraw
        }

        blockchain.withdraw(address, amountCanWithdraw.toFixed(8), "withdraw stake")
        packageInfo.lastBlockWithdraw = bn
        this._updatePackageInfo(address, packageNumber, packageInfo)
        blockchain.receipt(JSON.stringify([address, amountCanWithdraw.toFixed(8)]))
    }

    withdrawAll(address) {
        this._requireAuth(address, "active")
        // get count package
        let countPackage = storage.get(COUNT_PACKAGE_PREFIX + address)
        if(!countPackage) {
            throw new Error("Not found any package")
        }

        countPackage = Math.floor(countPackage)
        if(countPackage <= 0) {
            throw new Error("Not found any package")
        }
        // calc interest
        let amountCanWithdraw = new Float64("0")
        let maxAmountCanWithdraw = new Float64("0")
        const bn = block.number
        for(let i = 0; i < countPackage; i++) {
            const packageInfo = JSON.parse(storage.get(PACKAGE_INFO_PREFIX + address + "_" + i))
            const totalDayStake = Math.floor((bn - packageInfo.lastBlockWithdraw) / BLOCK_NUMBER_PER_DAY)

            if(totalDayStake <= 0) continue
            if(packageInfo.unstake) continue

            amountCanWithdraw = amountCanWithdraw.plus(this._calcInterest(packageInfo.amount, totalDayStake))
            maxAmountCanWithdraw = maxAmountCanWithdraw.plus(this._calcMaxInterest(packageInfo.amount, totalDayStake))
            // update package info
            packageInfo.lastBlockWithdraw = bn
            this._updatePackageInfo(address, i, packageInfo)
        }

        if(amountCanWithdraw.eq(0)) {
            throw new Error("All package can't withdraw")
        }

        if(amountCanWithdraw.gt(maxAmountCanWithdraw)) {
            this._updateRestAmount(amountCanWithdraw.minus(maxAmountCanWithdraw))
            amountCanWithdraw = maxAmountCanWithdraw
        }

        blockchain.withdraw(address, amountCanWithdraw.toFixed(8), "withdraw all stake")
        blockchain.receipt(JSON.stringify([address, amountCanWithdraw.toFixed(8)]))
    }

    unstake (address, packageNumber) {
        this._requireAuth(address, "active")

        const packageInfoString = storage.get(PACKAGE_INFO_PREFIX + address + "_" + packageNumber)
        if(!packageInfoString) {
            throw new Error("package not exist > " + packageNumber)
        }
        let packageInfo = JSON.parse(packageInfoString)

        if(packageInfo.unstake) {
            throw new Error("package has been unstake > " + packageNumber)
        }

        // check remain interest
        const bn = block.number
        const totalDayStake = Math.floor((bn - packageInfo.lastBlockWithdraw) / BLOCK_NUMBER_PER_DAY)

        if(totalDayStake > 0) {
            let amountCanWithdraw = this._calcInterest(packageInfo.amount, totalDayStake)
            const maxAmountCanWithdraw = this._calcMaxInterest(packageInfo.amount, totalDayStake)

            if(amountCanWithdraw.gt(maxAmountCanWithdraw)) {
                this._updateRestAmount(amountCanWithdraw.minus(maxAmountCanWithdraw))
                amountCanWithdraw = maxAmountCanWithdraw
            }

            blockchain.withdraw(address, amountCanWithdraw.toFixed(8), "withdraw stake")
            packageInfo.lastBlockWithdraw = bn
            this._updatePackageInfo(address, packageNumber, packageInfo)
            blockchain.receipt(JSON.stringify([address, amountCanWithdraw.toFixed(8)]))
        }

        // transfer freeze
        const freezeTime = tx.time + THREE_DAY_NANO
        const stakeAmount = new Float64(packageInfo.amount)
        blockchain.callWithAuth("token.empow", "transferFreeze", ["em", "stake.empow", address, packageInfo.amount, freezeTime, "unstake"])

        this._updateTotalStakeAmount(stakeAmount.negated())
        packageInfo.unstake = true
        this._updatePackageInfo(address, packageNumber, packageInfo)

        blockchain.receipt(JSON.stringify([address, packageInfo.amount]))
    }
}

module.exports = Stake
