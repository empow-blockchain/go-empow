class VotePoint {
    init() {
        this._initVotePoint()
    }

    initAdmin(adminID) {
        const bn = block.number;
        if(bn !== 0) {
            throw new Error("init out of genesis block");
        }
        storage.put("adminID", adminID);
    }

    _initVotePoint() {
        const bn = block.number;
        if(bn !== 0) {
            throw new Error("init out of genesis block");
        }
        blockchain.callWithAuth("token.empow", "create", [
            "vote",
            "vote.empow",
            90000000000,
            {
                "onlyIssuerCanTransfer": true,
                "decimal": 8
            }
        ]);
    }

    can_update(data) {
        const admin = storage.get("adminID");
        return blockchain.requireAuth(admin, "active");
    }
}

module.exports = VotePoint