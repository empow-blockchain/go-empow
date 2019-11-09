const oneYearNano = new Float64("31536000000000000");
const iostIssueRate = new Float64("0.0296").div(oneYearNano);
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
        storage.put("IOSTLastIssueTime", this._getBlockTime().toFixed());
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

    // issueEM to bonus.empow and iost foundation
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
            throw new Error("issue iost permission denied");
        }
        const lastIssueTime = storage.get("IOSTLastIssueTime");
        if (lastIssueTime === null || lastIssueTime === 0 || lastIssueTime === undefined) {
            throw new Error("IOSTLastIssueTime not set.");
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

        storage.put("IOSTLastIssueTime", currentTime.toFixed());

        const contractName = blockchain.contractName();
        const supply = new Float64(blockchain.callWithAuth("token.empow", "supply", ["em"])[0]);
        const issueAmount = supply.multi(iostIssueRate).multi(gap);
        const bonus = issueAmount.multi("0.33333333");
        // issue to foundation
        this._issueEM(foundationAcc, issueAmount.minus(bonus).minus(bonus).toFixed(decimal));
        // issue to producer with block reward
        this._issueEM("bonus.empow", bonus.toFixed(decimal));
        // issue to producer with vote percent
        this._issueEM(contractName, bonus.toFixed(decimal));

        const succ = blockchain.callWithAuth("vote_producer.empow", "topupCandidateBonus", [
            bonus.toFixed(decimal),
            contractName
        ])[0];
        if (!succ) {
            // transfer bonus to foundation if topup failed
            blockchain.transfer(contractName, foundationAcc, bonus.toFixed(decimal), "");
        }
    }
}

module.exports = IssueContract;
