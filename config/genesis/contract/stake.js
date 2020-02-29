const PACKAGE_INFO_PREFIX = 'p_'                                // p_address_count
const INTEREST_PREFIX = 'i_'                                    // i_day
const USER_STATISTIC = 'u_'                                     // u_address
const TOTAL_STAKE_AMOUNT = 'totalStakeAmount'
const REST_AMOUNT = "restAmount"
const MAXIMUM_PERCENT_PER_DAY = new Float64("0.000833333333")   // 30%/year
const BLOCK_NUMBER_PER_DAY = 172800
const THREE_DAY_NANO = 259200*1e9
const MINIMUM_STAKE = 1                                          // 1 EM

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

    _updateRestAmount(amount) {
        let restAmount = storage.get(REST_AMOUNT)
        storage.put(REST_AMOUNT, amount.plus(restAmount).toString())
    }

    _updatePackageInfo(address, packageId, data) {
        const prefix = PACKAGE_INFO_PREFIX + address + "_" + packageId
        storage.put(prefix, JSON.stringify(data))
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

    _updateUserStatistic(address, type, amount) {

        let userStatisticObj = storage.get(USER_STATISTIC + address)

        if(!userStatisticObj) {
            userStatisticObj = {
                countPackage: 1,
                staking: amount
            }

            storage.put(USER_STATISTIC + address, JSON.stringify(userStatisticObj))

            return 0
        }

        userStatisticObj = JSON.parse(userStatisticObj)

        if(type === "stake") {
            userStatisticObj.countPackage++
            userStatisticObj.staking += amount
        }

        if(type === "unstake") {
            userStatisticObj.staking -= amount
        }

        storage.put(USER_STATISTIC + address, JSON.stringify(userStatisticObj))

        return userStatisticObj.countPackage - 1;
    }

    stake(address, amount) {
        this._requireAuth(address, "active")
        let amountObj = new Float64(amount)
        if(amountObj.lt(MINIMUM_STAKE)) {
            throw new Error("minimum stake " + MINIMUM_STAKE + " EM")
        }
        // send EM to stake.empow
        blockchain.callWithAuth("token.empow", "transfer", ["em", address, blockchain.contractName(), amount, "stake EM"])
        // create package info
        let packageInfo = {
            startBlock: block.number,
            lastBlockWithdraw: block.number,
            unstake: false,
            amount: amount,
            startTime: tx.time,
            lastWithdrawTime: tx.time
        }
        const packageId = this._updateUserStatistic(address, "stake", amount)
        this._updatePackageInfo(address, packageId, packageInfo)
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

    withdraw(address, packageId) {
        this._requireAuth(address, "active")
        // check package exist
        const packageInfoString = storage.get(PACKAGE_INFO_PREFIX + address + "_" + packageId)
        if(!packageInfoString) {
            throw new Error("package not exist > " + packageId)
        }
        let packageInfo = JSON.parse(packageInfoString)
        // check package unstake
        if(packageInfo.unstake) {
            throw new Error("package has been unstake > " + packageId)
        }
        // calc interest
        const bn = block.number
        const totalDayStake = Math.floor((bn - packageInfo.lastBlockWithdraw) / BLOCK_NUMBER_PER_DAY)

        if(totalDayStake <= 0) {
            throw new Error("package withdraw less than 1 day > " + packageId)
        }

        let amountCanWithdraw = this._calcInterest(packageInfo.amount, totalDayStake)
        const maxAmountCanWithdraw = this._calcMaxInterest(packageInfo.amount, totalDayStake)

        if(amountCanWithdraw.gt(maxAmountCanWithdraw)) {
            this._updateRestAmount(amountCanWithdraw.minus(maxAmountCanWithdraw))
            amountCanWithdraw = maxAmountCanWithdraw
        }

        blockchain.withdraw(address, amountCanWithdraw.toFixed(8), "withdraw stake")
        packageInfo.lastBlockWithdraw = bn
        packageInfo.lastWithdrawTime = tx.time
        this._updatePackageInfo(address, packageId, packageInfo)
        blockchain.receipt(JSON.stringify([address, amountCanWithdraw.toFixed(8), packageId]))
    }

    unstake (address, packageId) {
        this._requireAuth(address, "active")

        const packageInfoString = storage.get(PACKAGE_INFO_PREFIX + address + "_" + packageId)
        if(!packageInfoString) {
            throw new Error("package not exist > " + packageId)
        }
        let packageInfo = JSON.parse(packageInfoString)

        if(packageInfo.unstake) {
            throw new Error("package has been unstake > " + packageId)
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
            packageInfo.lastWithdrawTime = tx.time
            this._updatePackageInfo(address, packageId, packageInfo)
            blockchain.receipt(JSON.stringify([address, amountCanWithdraw.toFixed(8)]))
        }

        // transfer freeze
        const freezeTime = tx.time + THREE_DAY_NANO
        const stakeAmount = new Float64(packageInfo.amount)
        blockchain.callWithAuth("token.empow", "transferFreeze", ["em", "stake.empow", address, packageInfo.amount, freezeTime, "unstake"])

        this._updateTotalStakeAmount(stakeAmount.negated())
        packageInfo.unstake = true
        this._updatePackageInfo(address, packageId, packageInfo)
        this._updateUserStatistic(address, "unstake", packageInfo.amount)

        blockchain.receipt(JSON.stringify([address, packageInfo.amount, packageId]))
    }
}

module.exports = Stake
