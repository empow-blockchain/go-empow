const POST_PREFIX = "p_"
const POST_STATISTIC_PREFIX = "s_"
const LIKE_PREFIX = "l_"
const COMMENT_PREFIX = "c_"
const REPLY_COMMENT_PREFIX = "rc_"
const REPORT_PREFIX = "r_"
const LEVEL_PREFIX = "lv_"
const INTEREST_PREEFIX = "i_"
const USER_PREFIX = "u_"
const USER_FOLLOW_PREFIX = "f_"
const TOTAl_LIKE = "totalLike"
const LIKE_BY_LEVEL = "likeByLevel"
const REST_AMOUNT = "restAmount"
const REPORT_TAG_ARRAY = "reportTagArray"

const BLOCK_NUMBER_PER_DAY = 172800

class Social {

    init() {
        storage.put(TOTAl_LIKE, "0")
        storage.put(REST_AMOUNT, "0")
        storage.put(REPORT_TAG_ARRAY, "[]")
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

    _titleValidate(title) {
        if(title.length <= 0 || title.length > 255) {
            throw new Error("title must be length greater than 0 and less than 255 > " + title.length)
        }
    }

    _contentValidate(content) {
        if(typeof content !== "object") {
            throw new Error("content must be object")
        }
    }

    _tagValidate(tag) {
        if(!Array.isArray(tag)) throw new Error("tag must be array")
        if(tag.length > 50) throw new Error("tag must be length less than 50 > " + tag.length)
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

    _requireAuth(address, permission) {
        const ret = blockchain.requireAuth(address, permission);
        if (ret !== true) {
            throw new Error("require auth > " + address)
        }
    }

    _updateRestAmount(amount) {
        let restAmount = storage.get(REST_AMOUNT)
        storage.put(REST_AMOUNT, amount.plus(restAmount).toString())
    }

    post (address, title, content, tag) {
        
        this._requireAuth(address, "active")
        this._titleValidate(title)
        this._tagValidate(tag)

        let id = tx.time
        let postObj = {
            postId: id.toString(),
            time: id,
            title: title,
            content: content,
            tag: tag,
        }

        let postStatisticObj = {
            author: address,
            totalLike:  0,
            realLike: 0,
            totalComment: 0,
            totalCommentAndReply: 0,
            totalReport: 0,
            realLikeArray: [0],
            lastBlockWithdraw: block.number
        }

        storage.put(POST_PREFIX + id, JSON.stringify(postObj))
        storage.put(POST_STATISTIC_PREFIX + id, JSON.stringify(postStatisticObj))

        blockchain.receipt(JSON.stringify([id.toString()]))
    }

    like(address, postId) {
        this._requireAuth(address, "active")
        // check exist post
        let postStatisticObj = storage.get(POST_STATISTIC_PREFIX + postId)
        if(!postStatisticObj) {
            throw new Error("PostId not exist > " + postId)
        }
        postStatisticObj = JSON.parse(postStatisticObj)

        // check author like
        if(address === postStatisticObj.author) {
            throw new Error("You can't like own post > " + postStatisticObj.author)
        }

        // check is exist like on this post
        if(storage.mapHas(LIKE_PREFIX + postId, address)) {
            throw new Error("You have been like this postId > " + postId)
        }

        // check level
        let level = Math.floor(storage.get(LEVEL_PREFIX + address))

        if(!level) {
            storage.put(LEVEL_PREFIX + address, "1")
            level = 1
        }

        let likeByLevel = JSON.parse(storage.get(LIKE_BY_LEVEL))
        let amountLike = likeByLevel[level - 1]

        // update post statistic
        postStatisticObj.totalLike++
        let bn = block.number
        let totalDayLike = Math.floor((bn - postStatisticObj.lastBlockWithdraw) / BLOCK_NUMBER_PER_DAY)

        if(typeof postStatisticObj.realLikeArray[totalDayLike] !== "number") {
            postStatisticObj.realLikeArray[totalDayLike] = 0
        }

        postStatisticObj.realLikeArray[totalDayLike] += amountLike
        postStatisticObj.realLike += amountLike

        storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))

        // like
        storage.mapPut(LIKE_PREFIX + postId, address, JSON.stringify(amountLike))

        // update total like
        this._updateTotalLike(new Float64(amountLike))

        blockchain.receipt(JSON.stringify([address, postId]))
    }

    likeWithdraw(postId) {
        // check exist post
        let postStatisticObj = storage.get(POST_STATISTIC_PREFIX + postId)
        if(!postStatisticObj) {
             throw new Error("PostId not exist > " + postId)
        }
        postStatisticObj = JSON.parse(postStatisticObj)
        let address = postStatisticObj.author

        this._requireAuth(address, "active")

        // calc
        const bn = block.number
        let totalDayLike = Math.floor((bn - postStatisticObj.lastBlockWithdraw) / BLOCK_NUMBER_PER_DAY)

        if(totalDayLike === 0) {
            throw new Error("You can withdraw like after 1 day")
        }

        if(typeof postStatisticObj.realLikeArray[totalDayLike] !== "number") {
            postStatisticObj.realLikeArray[totalDayLike] = 0
        }

        let count = 2;
        let canWithdraw = new Float64("0")
        let maxiumWithdraw = new Float64("0")
        const currentDay = Math.floor(bn / BLOCK_NUMBER_PER_DAY)

        for(let i = currentDay; i >= 0; i--) {
            if(postStatisticObj.realLikeArray.length - count < 0) {
                break
            }

            let realLike = postStatisticObj.realLikeArray[postStatisticObj.realLikeArray.length - count]
            if(typeof realLike !== "number") {
                count++
                continue
            }
            realLike = new Float64(realLike)

            let amountEMPerLike = storage.get(INTEREST_PREEFIX + i)
            if(!amountEMPerLike) continue
            
            maxiumWithdraw = maxiumWithdraw.plus(realLike)
            canWithdraw = canWithdraw.plus(realLike.multi(amountEMPerLike))

            count++
        }

        if(canWithdraw.gt(maxiumWithdraw)) {
            // update rest amount
            this._updateRestAmount(canWithdraw.minus(maxiumWithdraw))
            canWithdraw = maxiumWithdraw
        }

        if(canWithdraw.eq(0)) {
            throw new Error("Amount EM can withdraw is zero")
        }

        blockchain.withdraw(address, canWithdraw.toFixed(8), "like withdraw")

        // save post statistic
        postStatisticObj.realLikeArray = postStatisticObj.realLikeArray.slice(totalDayLike, totalDayLike + 1)
        postStatisticObj.lastBlockWithdraw = block.number

        // reward vote point
        blockchain.callWithAuth("vote.empow", "issueVotePoint", [address, canWithdraw.toFixed(8)])

        storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))

        blockchain.receipt(JSON.stringify([address, postId, canWithdraw.toFixed(8)]))
    }

    topup(amount) {
        this._requireAuth("base.empow", "active")
        amount = this._fixAmount(amount);

        // calc interest per 1 like
        let totalLike = storage.get(TOTAl_LIKE)
        if(Math.floor(totalLike) === 0) {
            totalLike = "1"
        }
        let interest = amount
        if(totalLike !== "0") {
            interest = amount.div(totalLike)
        }

        // update interest
        const bn = block.number
        const currentDay = Math.floor(bn / BLOCK_NUMBER_PER_DAY)
        storage.put(INTEREST_PREEFIX + currentDay, interest.toFixed(8))

        // reset total like
        this._resetTotalLike()
    }

    comment(address, postId, type, parentId, content) {
        this._requireAuth(address, "active")
        // check exist post
        let postStatisticObj = storage.get(POST_STATISTIC_PREFIX + postId)
        if(!postStatisticObj) {
             throw new Error("PostId not exist > " + postId)
        }
        postStatisticObj = JSON.parse(postStatisticObj)

        if(type !== "comment" && type !== "reply") {
            throw new Error("Wrong comment type > " + type)
        }

        if(content.length <= 0) {
            throw new Error("Content not blank > " + content)
        }

        let commentId = postStatisticObj.totalComment
        let subCommentId = 0 
        
        if(type === "comment") {
            const commentObj = {
                address: address,
                type: "comment",
                postId: postId,
                commentId: postStatisticObj.totalComment,
                parentId: 0,
                totalReply: 0,
                content: content
            }
            storage.put(COMMENT_PREFIX + postId + "_" + postStatisticObj.totalComment, JSON.stringify(commentObj))
            postStatisticObj.totalComment++
            postStatisticObj.totalCommentAndReply++
            storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))
        } else {
            // check exist comment
            let commentObj = storage.get(COMMENT_PREFIX + postId + "_" + parentId)
            if(!commentObj) {
                throw new Error("Comment not exist > " + parentId)
            }
            
            commentObj = JSON.parse(commentObj)

            commentId = parentId
            subCommentId = commentObj.totalReply

            let subCommentObj = {
                address: address,
                type: "reply",
                postId: postId,
                commentId: subCommentId,
                parentId: parentId,
                content: content
            }
            storage.put(REPLY_COMMENT_PREFIX + postId + "_" + parentId + "_" + commentObj.totalReply, JSON.stringify(subCommentObj))
            commentObj.totalReply++
            storage.put(COMMENT_PREFIX + postId + "_" + parentId, JSON.stringify(commentObj))
            postStatisticObj.totalCommentAndReply++
            storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))
        }

        blockchain.receipt(JSON.stringify([type, postId, commentId, subCommentId]))
    }

    report(address, postId, tag) {
        this._requireAuth(address, "active")
        // check post exist
        let postStatisticObj = storage.get(POST_STATISTIC_PREFIX + postId)
        if(!postStatisticObj) {
             throw new Error("PostId not exist > " + postId)
        }
        postStatisticObj = JSON.parse(postStatisticObj)
        // check tag exist
        const reportTagArray = JSON.parse(storage.get(REPORT_TAG_ARRAY))
        
        if(reportTagArray.indexOf(tag) === -1) {
            throw new Error("report tag not exist > " + tag)
        }

        // check report 2 times
        const reported = storage.get(REPORT_PREFIX + postId + "_" + address)
        if(reported) {
            throw new Error("can report 2 times > " + address)
        }
        
        if(storage.mapHas(REPORT_PREFIX + postId, tag)) {
            let current = Math.floor(storage.mapGet(REPORT_PREFIX + postId, tag))
            current++
            storage.mapPut(REPORT_PREFIX + postId, tag, current.toString())
        } else {
            storage.mapPut(REPORT_PREFIX + postId, tag, "1")
        }

        postStatisticObj.totalReport++
        storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))

        storage.put(REPORT_PREFIX + postId + "_" + address, "true")

        blockchain.receipt(JSON.stringify([address, postId, tag]))
    }

    follow (address, target) {
        this._requireAuth(address, "active")
        // check target exist
        if(!storage.globalMapHas("auth.empow", "auth", target)) {
            throw new Error("Follow user not exist > " + target)
        }
        // check is following
        if(storage.mapHas(USER_FOLLOW_PREFIX + target, address)) {
            throw new Error("You are following this user > " + target)
        }

        // set follow
        storage.mapPut(USER_FOLLOW_PREFIX + target, address, "true")

        blockchain.receipt(JSON.stringify([address, target]))
    }

    unfollow(address, target) {
        this._requireAuth(address, "active")
        // check target exist
        if(!storage.globalMapHas("auth.empow", "auth", target)) {
            throw new Error("Follow user not exist > " + target)
        }
        // check is following
        if(!storage.mapHas(USER_FOLLOW_PREFIX + target, address)) {
            throw new Error("You have not followed this user > " + target)
        }

        // del follow
        storage.mapDel(USER_FOLLOW_PREFIX + target, address)

        blockchain.receipt(JSON.stringify([address, target]))
    }

    updateProfile(address, info) {
        this._requireAuth(address, "active")
        storage.put(USER_PREFIX + address, JSON.stringify(info))
        blockchain.receipt(JSON.stringify([address, info]))
    }

    // admin only
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
            throw new Error("permission denied");
        }

        let currentLevel = Math.floor(storage.get(LEVEL_PREFIX + address))

        if(!currentLevel) {
            currentLevel = 1
        }

        if(!currentLevel || level <= currentLevel) return true;

        this._updateLevel(address, level)
    }

    addReportTag(tag) {
        const admin = storage.get("adminAddress");
        this._requireAuth(admin, "active")

        let reportTagArray = JSON.parse(storage.get(REPORT_TAG_ARRAY))

        if(reportTagArray.indexOf(tag) !== -1) {
            throw new Error("tag is exist > " + tag)
        }

        reportTagArray.push(tag)
        storage.put(REPORT_TAG_ARRAY, JSON.stringify(reportTagArray))
    }
}

module.exports = Social
