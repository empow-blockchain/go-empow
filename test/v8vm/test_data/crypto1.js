'use strict';
class crypto1 {
    sha3(msg) {
        return EMPOWCrypto.sha3(msg);
    }
    verify(algo, msg, sig, pubkey) {
        return EMPOWCrypto.verify(algo, msg, sig, pubkey);
    }

}

module.exports = crypto1;