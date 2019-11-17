const oneYearNano = new Float64("31536000000000000");
const issueRate = new Float64("0.04").div(oneYearNano);
const activePermission = "active";

class IssueContract {
    init() {
        storage.put("FoundationAccount", "");
    }

    _initEM(config, witnessInfo) {
        blockchain.callWithAuth("token.empow", "create", [
            "em",
            "issue.empow",
            config.EMTotalSupply,
            {
                "can_transfer": true,
                "decimal": config.EMDecimal
            }
        ]);
        for (const info of witnessInfo) {
            if (info.Balance !== 0) {
                blockchain.callWithAuth("token.empow", "issue", [
                    "em",
                    info.Address,
                    new Float64(info.Balance).toFixed()
                ]);
            }
        }
        storage.put("EMDecimal", new Int64(config.EMDecimal).toFixed());
        storage.put("lastIssueTime", this._getBlockTime().toFixed());
    }

    /**
     * genesisConfig = {
     *      FoundationAccount string
     *      EMTotalSupply   int64
     *      EMDecimal       int64
     * }
     * witnessInfo = [{
     *      Address      string
     *      Owner   string
     *      Active  string
     *      Balance int64
     * }]
     */
    initGenesis(adminAddress, genesisConfig, witnessInfo) {
        const bn = block.number;
        if(bn !== 0) {
            throw new Error("init out of genesis block")
        }
        storage.put("adminAddress", adminAddress);
        storage.put("FoundationAccount", genesisConfig.FoundationAccount);

        this._initEM(genesisConfig, witnessInfo);
    }

    can_update(data) {
        const admin = storage.get("adminAddress");
        this._requireAuth(admin, activePermission);
        return true;
    }

    _requireAuth(account, permission) {
        const ret = blockchain.requireAuth(account, permission);
        if (ret !== true) {
            throw new Error("require auth failed. ret = " + ret);
        }
    }

    _getBlockTime() {
        return new Float64(block.time);
    }

    _mapGet(k, f) {
        const val = storage.mapGet(k, f);
        if (val === "") {
            return null;
        }
        return JSON.parse(val);
    }

    _mapPut(k, f, v, p) {
        storage.mapPut(k, f, JSON.stringify(v), p);
    }

    _mapDel(k, f) {
        storage.mapDel(k, f);
    }

    _issueEM(account, amount) {
        const amountStr = ((typeof amount === "string") ? amount : amount.toFixed(this._get("EMDecimal")));
        const args = ["em", account, amountStr];
        blockchain.callWithAuth("token.empow", "issue", args);
    }

    issueEMToSell(amount) {
        const admin = storage.get("adminAddress");
        if(!blockchain.requireAuth(admin, "active")) {
            throw new Error("issue permission denied");
        }
        this._issueEM(admin, amount);
    }

    // issueEM
    issueEM() {
        const admin = storage.get("adminAddress");
        const whitelist = ["base.empow", admin];
        let auth = false;
        for (const c of whitelist) {
            if (blockchain.requireAuth(c, "active")) {
                auth = true;
                break;
            }
        }
        if (!auth) {
            throw new Error("issue permission denied");
        }
        const lastIssueTime = storage.get("lastIssueTime");
        if (lastIssueTime === null || lastIssueTime === 0 || lastIssueTime === undefined) {
            throw new Error("lastIssueTime not set.");
        }
        const currentTime = this._getBlockTime();
        const gap = currentTime.minus(lastIssueTime);
        if (gap.lte(0)) {
            return;
        }

        const foundationAcc = storage.get("FoundationAccount");
        const decimal = JSON.parse(storage.get("EMDecimal"));
        if (!foundationAcc) {
            throw new Error("FoundationAccount not set.");
        }

        storage.put("lastIssueTime", currentTime.toFixed());

        const contractName = blockchain.contractName();
        const supply = new Float64(blockchain.callWithAuth("token.empow", "supply", ["em"])[0]);

        const issueAmount = supply.multi(issueRate).multi(gap);
        const onePercentAmount = issueAmount.multi("0.25");
        const halfOnePercentAmount = onePercentAmount.div("2");
        const onePointFivePercentAmount = onePercentAmount.plus(halfOnePercentAmount);
        // issue to foundation
        this._issueEM(foundationAcc, issueAmount.minus(onePointFivePercentAmount).minus(onePointFivePercentAmount).toFixed(decimal));
        // issue to producer with block reward
        this._issueEM("bonus.empow", onePointFivePercentAmount.toFixed(decimal));
        this._issueEM(contractName, onePointFivePercentAmount.toFixed(decimal));

        // issue to producer with vote percent
        const succ = blockchain.callWithAuth("vote_producer.empow", "topupCandidateBonus", [
            onePercentAmount.toFixed(decimal),
            contractName
        ])[0];
        if (!succ) {
            // transfer bonus to foundation if topup failed
            blockchain.transfer(contractName, foundationAcc, onePercentAmount.toFixed(decimal), "");
        }

        const balance = new Float64(blockchain.callWithAuth("token.empow", "balanceOf", ["em", contractName])[0])

        if(balance.gte(halfOnePercentAmount)) {
            // issue to stake
            const succ2 = blockchain.callWithAuth("stake.empow", "topup", [
                halfOnePercentAmount.toFixed(decimal)
            ])[0];
            if (!succ2) {
                // transfer bonus to foundation if topup failed
                blockchain.transfer(contractName, foundationAcc, halfOnePercentAmount.toFixed(decimal), "");
            }
        } else {
            throw new Error("not enough balance")
        }

        
    }
}

module.exports = IssueContract;
