const POST_PREFIX = "p_"
const POST_STATISTIC_PREFIX = "s_"
const LIKE_PREFIX = "l_"
const COMMENT_PREFIX = "c_"
const REPORT_PREFIX = "r_"
const LEVEL_PREFIX = "lv_"
const LIKE_ARRAY = "likeArray"
const TOTAl_LIKE = "totalLike"
const LIKE_BY_LEVEL = "likeByLevel"

const BLOCK_NUMBER_PER_DAY = 172800

class Social {

    init() {
        storage.put(LIKE_ARRAY, [])
        storage.put(TOTAl_LIKE, "0")
        this.initLevel()
    }

    initLevel() {
        let likeAmount = [0.01, 10, 15, 18, 20, 25, 30]

        storage.put(LIKE_BY_LEVEL, JSON.stringify(likeAmount))
    }

    initAdmin (adminAddress) {
        const bn = block.number;
        if(bn !== 0) {
            throw new Error("init out of genesis block")
        }
        storage.put("adminAddress", adminAddress);
    }

    can_update(data) {
        const admin = storage.get("adminAddress");
        this._requireAuth(admin, "active");
        return true;
    }

    _postValidate(title, content, tag) {
        if(tag.length === undefined || tag.length === 0) {
            throw new Error("Tag is not correct format")
        }

        if(title.length > 255) {
            throw new Error("Title is greater than 255 characters")
        }
    }

    _fixAmount(amount) {
        amount = new Float64(new Float64(amount).toFixed(8));
        if (amount.lte("0")) {
            throw new Error("amount must be positive");
        }
        return amount;
    }

    _updateTotalLike(amount) {
        let totalLike = storage.get(TOTAl_LIKE)

        if(!totalLike) {
            storage.put(TOTAl_LIKE, amount.toString())
        } else {
            storage.put(TOTAl_LIKE, amount.plus(totalLike).toString())
        }
    }

    _resetTotalLike() {
        storage.put(TOTAl_LIKE, "0")
    }

    _updateLevel(address, level) {
        storage.put(LEVEL_PREFIX + address, level)
    }

    post (address, title, content, tag) {
        blockchain.requireAuth(address)

        try {
            content = JSON.parse(content)
        } catch(e) {
            throw new Error("Content is not correct format")
        }

        try {
            tag = JSON.parse(tag)
        } catch(e) {
            throw new Error("Tag is not correct format")
        }

        this._postValidate(title,content,tag)

        let id = tx.time
        let postObj = {
            time: id,
            title: title,
            content: content,
            tag: tag,
        }

        let postStatisticObj = {
            author: address,
            totalLike:  0,
            totalComment: 0,
            realLikeArray: [0],
            lastBlockWithdraw: block.number
        }

        storage.put(POST_PREFIX + id, JSON.stringify(postObj))
        storage.put(POST_STATISTIC_PREFIX + id, JSON.stringify(postStatisticObj))
    }

    like(address, postId) {
        blockchain.requireAuth(address)

        // check exist post
        let postStatisticObj = storage.get(POST_PREFIX + postId)
        if(!postStatisticObj) {
            throw new Error("PostId not exist > " + postId)
        }
        postStatisticObj = JSON.parse(postStatisticObj)

        // check is exist like on this post
        if(storage.mapHas(LIKE_PREFIX + postId, address)) {
            throw new Error("You have been like this postId > " + postId)
        }

        // check level
        let level = storage.get(LEVEL_PREFIX + address)

        if(!level) {
            storage.put(LEVEL_PREFIX + address, 1)
            level = 1
        }

        let likeByLevel = JSON.parse(storage.get(LIKE_BY_LEVEL))
        let amountLike = likeByLevel[level - 1]

        // update post statistic
        postStatisticObj.totalLike += 1
        let bn = block.number
        let totalDayLike = Math.floor((bn - postStatisticObj.lastBlockWithdraw) / BLOCK_NUMBER_PER_DAY)
        postStatisticObj.realLikeArray[totalDayLike] += amountLike

        storage.put(POST_STATISTIC_PREFIX + address, JSON.stringify(postStatisticObj))

        // like
        storage.mapPut(LIKE_PREFIX + postId, address, amountLike)

        // update total like
        this._updateTotalLike(new Float64(amountLike))

        blockchain.receipt(JSON.stringify([address, postId]))
    }

    likeWithdraw(postId) {
        // check exist post
        let postStatisticObj = storage.get(POST_PREFIX + postId)
        if(!postStatisticObj) {
             throw new Error("PostId not exist > " + postId)
        }
        postStatisticObj = JSON.parse(postStatisticObj)
        let address = postStatisticObj.author

        blockchain.requireAuth(address, "active")

        // calc
        let totalDayLike = Math.floor((bn - postStatisticObj.lastBlockWithdraw) / BLOCK_NUMBER_PER_DAY)

        if(totalDayLike === 0) {
            throw new Error("You can withdraw like after 1 day")
        }

        if(typeof postStatisticObj.realLikeArray[totalDayLike] !== "number") {
            postStatisticObj.realLikeArray[totalDayLike] = 0
        }

        let likeArray = JSON.parse(storage.get(LIKE_ARRAY))
        let count = 2;
        let canWithdraw = 0;
        let maxiumWithdraw = 0;

        for(i = likeArray.length - 1; i >= 0; i--) {
            if(postStatisticObj.realLikeArray.length - count < 0) {
                break;
            }

            let realLike = postStatisticObj.realLikeArray[postStatisticObj.realLikeArray.length - count]
            let amountEMPerLike = likeArray[i]

            maxiumWithdraw += realLike
            canWithdraw += realLike * amountEMPerLike

            count++
        }

        if(canWithdraw > maxiumWithdraw) {
            canWithdraw = maxiumWithdraw
        }

        canWithdraw = new Float64(canWithdraw)

        if(canWithdraw.eq(0)) {
            throw new Error("Amount EM can withdraw is zero")
        }

        blockchain.withdraw(address, canWithdraw.toFixed(8), "like withdraw")

        // save post statistic
        postStatisticObj.realLikeArray = postStatisticObj.realLikeArray.slice(totalDayLike - 1, totalDayLike)
        postStatisticObj.lastBlockWithdraw = block.number

        storage.put(POST_STATISTIC_PREFIX + address, JSON.stringify(postStatisticObj))
    }

    topup(amount) {
        blockchain.requireAuth("issue.empow", "active")
        amount = this._fixAmount(amount);

        blockchain.deposit("issue.empow", amount.toFixed(), "");

        // calc interest per 1 like
        const totalLike = storage.get(TOTAl_LIKE)
        let interest = amount
        if(totalLike !== "0") {
            interest = amount.div(totalLike)
        }

        // update like array
        let likeArray = JSON.parse(storage.get(LIKE_ARRAY))
        likeArray.push(interest.toFixed(8))
        storage.put(JSON.stringify(likeArray))

        // reset total like
        this._resetTotalLike()
    }

    upLevel(address, level) {
        const admin = storage.get("adminAddress");
        const whitelist = ["auth.empow", admin];
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

        let currentLevel = storage.get(LEVEL_PREFIX + address)

        if(!currentLevel || level <= currentLevel) return true;

        this._updateLevel(address, level)
    }
}

module.exports = Social