const premiumUsernamePriceUSD = 10 // $10/premium username

class Account {
    init() {
        storage.put('EMPrice', "0.001")
    }
    initAdmin(adminAddress) {
        const bn = block.number;
        if(bn !== 0) {
            throw new Error("init out of genesis block")
        }
        storage.put("adminAddress", adminAddress);
    }
    can_update(data) {
        const admin = storage.get("adminAddress");
        return blockchain.requireAuth(admin, "active");
    }
    _saveAccount(account, payer) {
        if (payer === undefined) {
            payer = account.address
        }
        storage.mapPut("auth", account.address, JSON.stringify(account), payer);
    }

    _loadAccount(id) {
        let a = storage.mapGet("auth", id);
        return JSON.parse(a);
    }

    static _find(items, name) {
        for (let i = 0; i < items.length; i++) {
            if (items[i].id === name) {
                return i
            }
        }
        return -1
    }

    static _findPermission(items, name) {
        const len = name.indexOf("@");
        if (len < 0) {
            for (let i = 0; i < items.length; i++) {
                if (items[i].id === name && items[i].permission === undefined) {
                    return i
                }
            }
        } else if (len > 0) {
            for (let i = 0; i < items.length; i++) {
                if (items[i].id === name.substring(0, len) && items[i].permission === name.substring(len+1, name.length)) {
                    return i
                }
            }
        } else {
            throw "unexpected item"
        }
        return -1
    }

    _hasAccount(id) {
        return storage.mapHas("auth", id);
    }

    _hasUsername(username) {
        return storage.mapHas("username", username);
    }

    _ra(id) {
        if (!blockchain.requireAuth(id, "owner")) {
            throw new Error("require auth failed");
        }
    }

    _checkAddressValid(address) {
        if (block.number === 0) {
            return
        }
        if(address.length != 49) {
            throw new Error("wrong address")
        }
        if (!address.startsWith("EM")) {
            throw new Error("address invalid. address must start with EM");
        }
        for (let i in address) {
            let ch = address[i];
            if (!(ch >= 'A' && ch <= 'z' || ch >= '0' && ch <= '9')) {
                throw new Error("address invalid. address contains invalid character > " + ch);
            }
        }
    }

    _checkPermValid(perm) {
        if (block.number === 0) {
            return
        }
        if (perm.length < 1 || perm.length > 32) {
            throw new Error("id invalid. id length should be between 1,32 > " + id)
        }
        for (let i in perm) {
            let ch = perm[i];
            if (!(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch === '_')) {
                throw new Error("id invalid. id contains invalid character > " + ch);
            }
        }
    }

    _checkWeight(weight) {
        if (weight <= 0) {
            throw "weight less than zero"
        }
    }

    _checkNormalUsername(username) {
        if (block.number === 0) {
            return
        }
        if(username.length < 6 || username.length > 32) {
            throw new Error("username length must 6-32 characters")
        }
        for (let i in username) {
            let ch = username[i];
            if (!(ch >= 'A' && ch <= 'z' || ch >= '0' && ch <= '9')) {
                throw new Error("username invalid. username contains invalid character > " + ch);
            }
        }
    }

    _checkPremiumUsername(username) {
        if (block.number === 0) {
            return
        }
        if(username.length < 6 || username.length > 32) {
            throw new Error("username length must 6-32 characters")
        }
        for (let i in username) {
            let ch = username[i];
            if (!(ch >= 'A' && ch <= 'z' || ch >= '0' && ch <= '9') || ch == '_') {
                throw new Error("username invalid. username contains invalid character > " + ch);
            }
        }
    }

    /**
     * @param  {string} address - this is a string
     *
     */
    signUp(address, owner, active) {
        if (this._hasAccount(address)) {
            throw new Error("address existed > " + address);
        }
        this._checkAddressValid(address);
        if(block.number != 0) {
            blockchain.requireAuth("token.empow", "active")
        }
        let account = {};
        account.address = address;
        account.permissions = {};
        account.permissions.active = {
            name: "active",
            groups: [],
            items: [{
                id: active,
                is_key_pair: true,
                weight: 100,
            }],
            threshold: 100,
        };
        account.permissions.owner = {
            name: "owner",
            groups: [],
            items: [{
                id: owner,
                is_key_pair: true,
                weight: 100,
            }],
            threshold: 100,
        };
        account.groups = {};
        this._saveAccount(account, blockchain.contractName());
        if (block.number !== 0) {
            const defaultGasPledge = "15";
            const defaultRamBuy = 200; // 200 bytes
            blockchain.callWithAuth("gas.empow", "pledge", [blockchain.contractName(), address, defaultGasPledge]);
            blockchain.callWithAuth("ram.empow", "buy", [blockchain.contractName(), address, defaultRamBuy]);
        }

        blockchain.receipt(JSON.stringify([address, owner, active]));
    }

    addNormalUsername (address, username) {
        if (this._hasUsername("newbie." + username)) {
            throw new Error("username existed > " + "newbie." + username);
        }
        // check isvailid username
        this._checkNormalUsername(username)
        // require auth address
        blockchain.requireAuth(address, 'active')
        // save to storage with prefix "newbie."
        storage.mapPut("username", "newbie." + username, address, address);
    }

    addPremiumUsername(address, username) {
        if (this._hasUsername(username)) {
            throw new Error("username existed > " + username);
        }
        // check isvailid username
        this._checkPremiumUsername(username)
        // require auth address
        blockchain.requireAuth(address, 'active')
        const EMPrice = storage.get('EMPrice')
        const EMneedToPay = new Float64(premiumUsernamePriceUSD).div(EMPrice)
        blockchain.callWithAuth("token.empow", "transfer", ["em", address, "deadaddr", EMneedToPay, "pay premium username"]);
        // pledge gas and buy ram
        const halfEMAmount = EMneedToPay.div(2).toFixed(8)
        const ramPrice = blockchain.callWithAuth("ram.empow", "getPrice", [])
        const ramAmount = new Float64(halfEMAmount).div(ramPrice).toFixed(0)

        blockchain.callWithAuth("gas.empow", "pledge", [blockchain.contractName(), address, halfEMAmount]);
        blockchain.callWithAuth("ram.empow", "buy", [blockchain.contractName(), address, ramAmount * 1.00]);
        // save to storage
        storage.mapPut("username",username, address, address);
    }

    addPermission(id, perm, thres) {
        this._ra(id);
        this._checkPermValid(perm);
        let acc = this._loadAccount(id);
        if (acc.permissions[perm] !== undefined) {
            throw new Error("permission already exist");
        }
        acc.permissions[perm] = {
            name: perm,
            groups: [],
            items: [],
            threshold: thres,
        };
        this._saveAccount(acc);

        blockchain.receipt(JSON.stringify([id, perm, thres]));
    }

    dropPermission(id, perm) {
        this._ra(id);
        if (perm === "active" || perm === "owner") {
            throw "drop active or owner is forbidden"
        }
        let acc = this._loadAccount(id);
        acc.permissions[perm] = undefined;
        this._saveAccount(acc);

        blockchain.receipt(JSON.stringify([id, perm]));
    }

    assignPermission(id, perm, un, weight) {
        this._ra(id);
        this._checkWeight(weight);
        let acc = this._loadAccount(id);
        const index = Account._find(acc.permissions[perm].items, un);
        if (index < 0) {
            const len = un.indexOf("@");
            if (len < 0) {
                acc.permissions[perm].items.push({
                    id: un,
                    is_key_pair: true,
                    weight: weight
                });
            } else if (len > 0) {
                acc.permissions[perm].items.push({
                    id: un.substring(0, len),
                    permission: un.substring(len+1, un.length),
                    is_key_pair: false,
                    weight: weight
                });
            } else {
                throw "unexpected item"
            }
        } else {
            acc.permissions[perm].items[index].weight = weight
        }
        this._saveAccount(acc);

        blockchain.receipt(JSON.stringify([id, perm, un, weight]));
    }

    revokePermission(id, perm, un) {
        this._ra(id);
        let acc = this._loadAccount(id);
        const index = Account._findPermission(acc.permissions[perm].items, un);
        if (index < 0) {
            throw new Error("item not found");
        } else {
            acc.permissions[perm].items.splice(index, 1)
        }
        this._saveAccount(acc);

        blockchain.receipt(JSON.stringify([id, perm, un]));
    }

    addGroup(id, grp) {
        this._ra(id);
        this._checkPermValid(grp);
        let acc = this._loadAccount(id);
        if (acc.groups[grp] !== undefined) {
            throw new Error("group already exist");
        }
        acc.groups[grp] = {
            name: grp,
            items: [],
        };
        this._saveAccount(acc);

        blockchain.receipt(JSON.stringify([id, grp]));
    }

    dropGroup(id, group) {
        this._ra(id);
        let acc = this._loadAccount(id);
        acc.groups[group] = undefined;
        for (let i = 0; i < acc.permissions.length; i++) {
            for (let j = 0; j < acc.permissions[i].groups.length; j++) {
                if (acc.permissions[i].groups[j] === group) {
                    acc.permissions[i].groups.splice(j, 1)
                }
            }
        }
        this._saveAccount(acc);

        blockchain.receipt(JSON.stringify([id, group]));
    }

    assignGroup(id, group, un, weight) {
        this._ra(id);
        this._checkWeight(weight);
        let acc = this._loadAccount(id);
        const index = Account._find(acc.groups[group].items, un);
        if (index < 0) {
            let len = un.indexOf("@");
            if (len < 0) {
                acc.groups[group].items.push({
                    id: un,
                    is_key_pair: true,
                    weight: weight
                });
            } else {
                acc.groups[group].items.push({
                    id: un.substring(0, len),
                    permission: un.substring(len+1, un.length),
                    is_key_pair: false,
                    weight: weight
                });
            }
        } else {
            acc.groups[group].items[index].weight = weight
        }

        this._saveAccount(acc);

        blockchain.receipt(JSON.stringify([id, group, un, weight]));
    }

    revokeGroup(id, grp, un) {
        this._ra(id);
        let acc = this._loadAccount(id);
        const index = Account._findPermission(acc.groups[grp].items, un);
        if (index < 0) {
            throw new Error("item not found");
        } else {
            acc.groups[grp].items.splice(index, 1)
        }
        this._saveAccount(acc);

        blockchain.receipt(JSON.stringify([id, grp, un]));
    }

    assignPermissionToGroup(id, perm, group) {
        this._ra(id);
        let acc = this._loadAccount(id);
        if (acc.groups[group] === undefined) {
            throw new Error("group does not exist");
        }
        acc.permissions[perm].groups.push(group);
        this._saveAccount(acc);

        blockchain.receipt(JSON.stringify([id, perm, group]));
    }

    revokePermissionInGroup(id, perm, group) {
        this._ra(id);
        let acc = this._loadAccount(id);
        let index = acc.permissions[perm].groups.indexOf(group);
        if (index > -1) {
            acc.permissions[perm].groups.splice(index, 1);
        }
        this._saveAccount(acc);

        blockchain.receipt(JSON.stringify([id, perm, group]));
    }
}

module.exports = Account;
