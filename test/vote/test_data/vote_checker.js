const VPContract = "vote_producer.empow";
const BonusContract = "bonus.empow";
const IssueContract = "issue.empow";
const TokenContract = "token.empow";

class VoteChecker {
    init() {
    }

    _get(k) {
        const val = storage.get(k);
        if (val === "") {
            return null;
        }
        return JSON.parse(val);
    }

	_put(k, v, p) {
        storage.put(k, JSON.stringify(v), p);
    }

    _call(contract, api, args) {
        const ret = blockchain.callWithAuth(contract, api, JSON.stringify(args));
        if (ret && Array.isArray(ret) && ret.length >= 1) {
            return ret[0] === "" ? "" : JSON.parse(ret[0]);
        }
        return ret;
    }

    vote(from, to, amount) {
        let voteInfo1 =  this._call(VPContract, "getProducer", [to]);
        blockchain.callWithAuth(VPContract, "vote", [from, to, amount]);
        let voteInfo2 =  this._call(VPContract, "getProducer", [to]);
        let data = this._get("vote") || {};
        let f = data[from] || {};
        f[to] = {
            vote: new Float64(amount).plus(f[to] || 0).toFixed(8),
            VPContract: new Float64(voteInfo2.voteInfo.votes).minus(voteInfo1.voteInfo.votes).toFixed(8),
        };
        data[from] = f;
        this._put("vote", data);
        return data;
    }

    unvote(from, to, amount) {
        let voteInfo1 =  this._call(VPContract, "getProducer", [to]);
        blockchain.callWithAuth(VPContract, "unvote", [from, to, amount]);
        let voteInfo2 =  this._call(VPContract, "getProducer", [to]);
        let data = this._get("unvote") || {};
        let f = data[from] || {};
        f[to] = {
            vote: new Float64(amount).plus(f[to] || 0).toFixed(8),
            VPContract: new Float64(voteInfo2.voteInfo.votes).minus(voteInfo1.voteInfo.votes).toFixed(8),
        };
        data[from] = f;
        this._put("unvote", data);
        return data;
    }

    issueEM() {
        let total1 = blockchain.callWithAuth(TokenContract, "supply", ["em"])[0];
        blockchain.callWithAuth(IssueContract, "issueEM", []);
        let total2 = blockchain.callWithAuth(TokenContract, "supply", ["em"])[0];
        let data = this._get("issueEM") || [];
        data.push(new Float64(total2).minus(total1).toFixed(8));
        this._put("issueEM", data);
        return data;
    }

    exchangeEMPOW() {
        let publisher = blockchain.publisher();
        let balance10 = blockchain.callWithAuth(TokenContract, "balanceOf", ["em", BonusContract])[0];
        let balance11 = blockchain.callWithAuth(TokenContract, "balanceOf", ["em", publisher])[0];
        blockchain.callWithAuth(BonusContract, "exchangeEMPOW", [publisher, "0"]);
        let balance20 = blockchain.callWithAuth(TokenContract, "balanceOf", ["em", BonusContract])[0];
        let balance21 = blockchain.callWithAuth(TokenContract, "balanceOf", ["em", publisher])[0];
        let data = this._get("exchangeEMPOW") || {};
        data[publisher] = {
            BonusContract: new Float64(balance20).minus(balance10).toFixed(8),
            publisher: new Float64(balance21).minus(balance11).toFixed(8),
        }
        this._put("exchangeEMPOW", data);
        return data;
    }

    candidateWithdraw() {
        let publisher = blockchain.publisher();
        let balance10 = blockchain.callWithAuth(TokenContract, "balanceOf", ["em", VPContract])[0];
        let balance11 = blockchain.callWithAuth(TokenContract, "balanceOf", ["em", publisher])[0];
        let bonus = blockchain.callWithAuth(VPContract, "getCandidateBonus", [publisher])[0];
        blockchain.callWithAuth(VPContract, "candidateWithdraw", [publisher]);
        let balance20 = blockchain.callWithAuth(TokenContract, "balanceOf", ["em", VPContract])[0];
        let balance21 = blockchain.callWithAuth(TokenContract, "balanceOf", ["em", publisher])[0];
        let voteInfo =  this._call(VPContract, "getProducer", [publisher]);
        let vote = this._get("vote");
        let unvote = this._get("unvote");
        let votes = {};
        for (let a in vote) {
            let v = (vote[a][publisher] || {})["vote"] || "0";
            votes[a] = new Float64(v).plus(votes[a] || "0").toFixed(8);
        }
        for (let a in unvote) {
            let v = (unvote[a][publisher] || {})["vote"] || "0";
            votes[a] = new Float64(votes[a] || "0").minus(v).toFixed(8);
        }
        let data = this._get("candidateWithdraw") || {};
        data[publisher] = {
            bonus: bonus,
            VPContract: new Float64(balance20).minus(balance10).toFixed(8),
            publisher: new Float64(balance21).minus(balance11).toFixed(8),
            votes: votes,
            totalVotes: voteInfo.voteInfo.votes,
        }
        this._put("candidateWithdraw", data);
        return data;
    }
}

module.exports = VoteChecker;