const PREMIUM_USERNAME_PRICE = "1000" // 1000 EM
const USERNAME_ARRAY_PREFIX = "u_"
const SELECTED_USERNAME_PREFIX = "s_"

class Account {
    init() {
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
    _saveAccount(account) {
        storage.mapPut("auth", account.address, JSON.stringify(account));
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

    _addUsernameToArray(address, username) {
        let arr = storage.get(USERNAME_ARRAY_PREFIX + address)

        if(!arr) {
            arr = [username]
        } else {
            arr = JSON.parse(arr)
            arr.push(username)
        }

        storage.put(USERNAME_ARRAY_PREFIX + address, JSON.stringify(arr))
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
            this._ra("token.empow", "active")
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
        this._saveAccount(account);
        if (block.number !== 0) {
            const defaultGasPledge = "15";
            // const defaultRamBuy = 200; // 200 bytes
            blockchain.callWithAuth("gas.empow", "pledge", [blockchain.contractName(), address, defaultGasPledge]);
            // blockchain.callWithAuth("ram.empow", "buy", [blockchain.contractName(), address, defaultRamBuy]);
        }

        blockchain.receipt(JSON.stringify([address, owner, active]));
    }

    addNormalUsername (address, username) {
        if (this._hasUsername("newbie." + username)) {
            throw new Error("username existed > " + "newbie." + username);
        }
        // check can buy free username one time
        let addressUsernameArray = storage.get(USERNAME_ARRAY_PREFIX + address)
        if(addressUsernameArray) {
            addressUsernameArray = JSON.parse(addressUsernameArray)
            for(let i = 0; i < addressUsernameArray.length; i++) {
                if(addressUsernameArray[i].indexOf("newbie.") !== -1) {
                    throw new Error("You can only buy free username 1 time")
                }
            }
        }
        // check isvailid username
        this._checkNormalUsername(username)
        // require auth address
        this._ra(address, "active")
        // save to storage with prefix "newbie."
        storage.mapPut("username", "newbie." + username, address);
        // add username to username array
        this._addUsernameToArray(address, "newbie." + username)
        blockchain.receipt(JSON.stringify([address, "newbie." + username]))
    }

    addPremiumUsername(address, username) {
        if (this._hasUsername(username)) {
            throw new Error("username existed > " + username);
        }
        // check isvailid username
        this._checkPremiumUsername(username)
        // require auth address
        this._ra(address, "active")
        blockchain.callWithAuth("token.empow", "transfer", ["em", address, "deadaddr", PREMIUM_USERNAME_PRICE, "pay premium username"]);
        // pledge more gas
        const gasReward = "10"
        blockchain.callWithAuth("gas.empow", "pledge", [blockchain.contractName(), address, gasReward]);
        // reward vote point
        blockchain.callWithAuth("vote.empow", "issueVotePoint", [address, PREMIUM_USERNAME_PRICE])
        // save to storage
        storage.mapPut("username",username, address);
        // add username to array
        this._addUsernameToArray(address,username)
        blockchain.receipt(JSON.stringify([address,username]))
    }

    selectUsername(username) {
    	if (!this._hasUsername(username)) {
            throw new Error("username not exist > " + username);
        }

        let address = storage.mapGet("username", username)
        this._ra(address, "active")

        storage.put(SELECTED_USERNAME_PREFIX + address, username)

        blockchain.receipt(JSON.stringify([address,username]))
    }
}

module.exports = Account;
