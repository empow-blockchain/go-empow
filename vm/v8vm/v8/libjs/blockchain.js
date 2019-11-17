let BlockChain = (function () {
    let bc = new EMPOWBlockchain;
    let storage = new EMPOWStorage;
    // get contractName
    let contractName = function () {
        let ctxInfo = JSON.parse(bc.contextInfo());
        return ctxInfo["contract_name"];
    };
    // get publisher
    let publisher = function () {
        let ctxInfo = JSON.parse(bc.contextInfo());
        return ctxInfo["publisher"];
    };
    // transfer EMPOWoken
    let transfer = function (from, to, amount, memo) {
        if (!(amount instanceof Float64)) {
            amount = new Float64(amount);
        }
        const args = ["em", from, to, amount.toString(), memo.toString()];
        return JSON.parse(bc.callWithAuth("token.empow", "transfer", JSON.stringify(args)));
    };
    return {
        // transfer EMPOWoken
        transfer: transfer,
        // withdraw EMPOWoken
        withdraw: function (to, amount, memo) {
            return transfer(contractName(), to, amount, memo);
        },
        // deposit EMPOWoken
        deposit: function (from, amount, memo) {
            return transfer(from, contractName(), amount, memo);
        },
        // get blockInfo
        blockInfo: function () {
            return bc.blockInfo();
        },
        // get transactionInfo
        txInfo: function () {
            return bc.txInfo();
        },
        // get transactionInfo
        contextInfo: function () {
            return bc.contextInfo();
        },
        // get contractName
        contractName: contractName,
        // get publisher
        publisher: publisher,
        // get contractOwner
        contractOwner: function() {
            return storage.globalMapGet("system.empow", "contract_owner", contractName(), "")
        },
        // call contract's api using args
        call: function (contract, api, args) {
            if (typeof args == "object") {
                args = JSON.stringify(args);
            }
            return JSON.parse(bc.call(contract, api, args));
        },
        // call contract's api using args with auth
        callWithAuth: function (contract, api, args) {
            if (typeof args == "object") {
                args = JSON.stringify(args);
            }
            return JSON.parse(bc.callWithAuth(contract, api, args));
        },
        // check account's permission
        requireAuth: function (accountID, permission) {
            return bc.requireAuth(accountID, permission);
        },
        // generate receipt
        receipt: function (content) {
            return bc.receipt(content);
        },
        // post event
        event: function (content) {
            return bc.event(content);
        },
    }
})();

module.exports = BlockChain;
