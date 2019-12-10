const ADMIN_PERMISSION = "active";

class Exchange {
    init() {

    }

    can_update(data) {
        this._requireAuth("EM2ZsSw7RWYC229Z1ib7ujKhken9GFR7dBkTTEbBWMKeLpVas", ADMIN_PERMISSION);
        return true;
    }

    _requireAuth(account, permission) {
        const ret = blockchain.requireAuth(account, permission);
        if (ret !== true) {
            throw new Error("require auth failed. ret = " + ret);
        }
    }

    /**
     *
     * @param tokenSym  {string}  token symbol
     * @param to        {string}  to account, create new account if empty
     * @param amount    {string}  token amount
     * @param memo      {string}  command:args, e.g. create:userName:ownerKey:activeKey
     *
     * // 1. normal transfer
     * transfer("em", "user1", "100.1", "")
     * // 2. create an account, buy initialRAM and pledge initialGas, then transfer
     * transfer("em", "", "100.1", "create:newUser2:OWNERKEY:ACTIVEKEY")
     */
    transfer(tokenSym, to, amount, memo) {
        let from = blockchain.publisher();
        if (to !== "") {
            // transfer to an exist account
            blockchain.call("token.empow", "transfer", [tokenSym, from, to, amount, memo]);

        } else if (to == "") {
            const minAmount = 100;
            const initialRAM = 1000;
            const initialGasPledged = 10;
            let bamount = new BigNumber(amount);
            if (bamount.lt(minAmount)) {
                throw new Error("transfer amount should be greater or equal to " + minAmount);
            }

            if (memo.startsWith("create:")) {
                if (tokenSym !== "em") {
                    throw new Error("must transfer EM if you want to create a new account");
                }
                // create account and then transfer to account
                let args = memo.split(":").slice(1);
                if (args.length !== 3) {
                    throw new Error("memo of transferring to a new account should be of format create:name:ownerKey:activeKey");
                }
                blockchain.call("auth.empow", "signUp", args);
                let rets = blockchain.call("ram.empow", "buy", [from, args[0], initialRAM]);
                let price = rets[0];

                let paid = new BigNumber(price).plus(new BigNumber(initialGasPledged));
                if (bamount.lt(paid)) {
                    throw new Error("amount not enough to buy 1kB RAM and pledge 10 EM Gas. need " + bamount.toString())
                }

                blockchain.transfer(from, args[0], bamount.minus(paid), memo);
            } else {
                throw new Error("unsupported command : " + memo + ", use create:XX")
            }
        }
    }

}

module.exports = Exchange;
