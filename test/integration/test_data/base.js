class Base {
    init() {
    }

    initWitness(lst) {
        const map = {};
        for (const witness of lst) {
            map[witness] = 1;
        }
        storage.put("witness_produced", JSON.stringify(map));
    }

    stat() {
        blockchain.callWithAuth("vote_producer.empow", "stat", '[]');
    }

    issueContribute(data) {
        blockchain.callWithAuth("bonus.empow", "issueContribute", JSON.stringify([data]));
    }
}

module.exports = Base;
