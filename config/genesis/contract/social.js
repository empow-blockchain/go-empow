const POST_PREFIX = "p_"
const POST_STATISTIC_PREFIX = "s_"
const LIKE_PREFIX = "l_"
const LIKE_COMMENT_PREFIX = "lc_"
const COMMENT_PREFIX = "c_"
const REPLY_COMMENT_PREFIX = "rc_"
const INTEREST_PREEFIX = "i_"
const USER_PREFIX = "u_"
const USER_LIKE_STATISTIC = "ul_"
const USER_FOLLOW_PREFIX = "f_"
const STAKE_USER_STATISTIC = "u_"
const TOTAl_LIKE = "totalLike"
const LIKE_RATIO = 0.00025          // Staking 1 EM -> 1 Like = 0.25 EM, 
const LIKE_MAX = 50                 // Staking 1,000,000 EM -> 1 Like = 250 EM > LIKE_MAX -> 1 Like = 50 EM
const LIKE_MAX_PER_DAY = 10
const REST_AMOUNT = "restAmount"
const BLOCK_TAG_PREFIX = "bt_"
const BLOCK_NUMBER_PER_DAY = 172800
const REPORT_TAG_ARRAY = "reportTagArray"
const REPORT_PENDING_ARRAY = "reportPendingArray"
const REPORT_VALIDATOR_PREFIX = "rp_"
const REPORT_VALIDATOR_PER_POST = 9
const REPORT_VALIDATOR_REWARD = 10
const REPORT_VALIDATOR_MINIMUM_STAKING = 10000

class Social {

    init() {
        storage.put(TOTAl_LIKE, "0")
        storage.put(REST_AMOUNT, "0")
        storage.put(REPORT_TAG_ARRAY, "[]")
    }

    initAdmin(adminAddress) {
        const bn = block.number;
        if (bn !== 0) {
            throw new Error("init out of genesis block")
        }
        storage.put("adminAddress", adminAddress);
    }

    can_update(data) {
        const admin = storage.get("adminAddress");
        this._requireAuth(admin, "active");
        return true;
    }

    _titleValidate(title) {
        if (title.length <= 0 || title.length > 255) {
            throw new Error("title must be length greater than 0 and less than 255 > " + title.length)
        }
    }

    _contentValidate(content) {
        if (typeof content !== "object") {
            throw new Error("content must be object")
        }
    }

    _tagValidate(tag) {
        if (!Array.isArray(tag)) throw new Error("tag must be array")
        if (tag.length > 50) throw new Error("tag must be length less than 50 > " + tag.length)
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

        if (!totalLike) {
            storage.put(TOTAl_LIKE, amount.toString())
        } else {
            storage.put(TOTAl_LIKE, amount.plus(totalLike).toString())
        }
    }

    _resetTotalLike() {
        storage.put(TOTAl_LIKE, "0")
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

    _checkPostExist(postId) {
        let postStatisticObj = storage.get(POST_STATISTIC_PREFIX + postId)
        if (!postStatisticObj) {
            throw new Error("PostId not exist > " + postId)
        }

        postStatisticObj = JSON.parse(postStatisticObj)
        if (postStatisticObj.deleted) {
            throw new Error("PostId not exist > " + postId)
        }

        return postStatisticObj
    }

    _checkCommentExist(postId, commentId) {
        let commentObj = storage.get(COMMENT_PREFIX + postId + "_" + commentId)

        if (!commentObj) {
            throw new Error("CommentId not exist > " + commentId)
        }

        return JSON.parse(commentObj)
    }

    _checkTagExistInReportTagArray(tag) {
        let reportTagArray = storage.get(REPORT_TAG_ARRAY)

        if (!reportTagArray) throw new Error("report tag array not exist")

        reportTagArray = JSON.parse(reportTagArray)

        if (!reportTagArray.includes(tag)) {
            throw new Error("block tag not exist > " + tag)
        }
    }

    _updateLikeReward(obj, amountLike) {
        let bn = block.number
        let totalDayLike = Math.floor((bn - obj.lastBlockWithdraw) / BLOCK_NUMBER_PER_DAY)

        if (typeof obj.realLikeArray[totalDayLike] !== "number") {
            obj.realLikeArray[totalDayLike] = 0
        }

        obj.realLikeArray[totalDayLike] += amountLike
        obj.realLike += amountLike

        // update total like today
        this._updateTotalLike(new Float64(amountLike))

        return obj
    }

    _calcCanWithdraw(obj) {
        const bn = block.number
        let totalDayLike = Math.floor((bn - obj.lastBlockWithdraw) / BLOCK_NUMBER_PER_DAY)

        if (totalDayLike === 0) {
            throw new Error("You can withdraw like after 1 day")
        }

        if (typeof obj.realLikeArray[totalDayLike] !== "number") {
            obj.realLikeArray[totalDayLike] = 0
        }

        let count = 2;
        let canWithdraw = new Float64("0")
        let maxiumWithdraw = new Float64("0")
        const currentDay = Math.floor(bn / BLOCK_NUMBER_PER_DAY)

        for (let i = currentDay; i >= 0; i--) {
            if (obj.realLikeArray.length - count < 0) {
                break
            }

            let realLike = obj.realLikeArray[obj.realLikeArray.length - count]
            if (typeof realLike !== "number") {
                count++
                continue
            }
            realLike = new Float64(realLike)

            let amountEMPerLike = storage.get(INTEREST_PREEFIX + i)
            if (!amountEMPerLike) continue

            maxiumWithdraw = maxiumWithdraw.plus(realLike)
            canWithdraw = canWithdraw.plus(realLike.multi(amountEMPerLike))

            count++
        }

        if (canWithdraw.gt(maxiumWithdraw)) {
            // update rest amount
            this._updateRestAmount(canWithdraw.minus(maxiumWithdraw))
            canWithdraw = maxiumWithdraw
        }

        if (canWithdraw.eq(0)) {
            throw new Error("Amount EM can withdraw is zero")
        }

        return { canWithdraw, totalDayLike }
    }

    _updateLikeStatisticForUser(address) {
        const bn = block.number
        let likeStatistic = storage.get(USER_LIKE_STATISTIC + address)

        if (!likeStatistic) {
            likeStatistic = {
                lastLikeBlock: bn,
                toDayLike: 1
            }

            storage.put(USER_LIKE_STATISTIC + address, JSON.stringify(likeStatistic))
        } else {
            likeStatistic = JSON.parse(likeStatistic)
            const toDay = Math.floor(bn / BLOCK_NUMBER_PER_DAY)
            const lastLikeOnDay = Math.floor(likeStatistic.lastLikeBlock / BLOCK_NUMBER_PER_DAY)

            if (toDay === lastLikeOnDay) {
                likeStatistic.lastLikeBlock = bn
                likeStatistic.toDayLike++
            } else {
                likeStatistic.lastLikeBlock = bn
                likeStatistic.toDayLike = 1
            }

            storage.put(USER_LIKE_STATISTIC + address, JSON.stringify(likeStatistic))
        }
    }

    _isMaxLikePerDay(address) {
        let likeStatistic = storage.get(USER_LIKE_STATISTIC + address)
        if (!likeStatistic) return false;

        likeStatistic = JSON.parse(likeStatistic)

        if(likeStatistic.toDayLike > LIKE_MAX_PER_DAY) {
            return true
        } else {
            return false
        }
    }

    _getStaking(address) {
        let userStatisticObj = storage.globalGet("stake.empow", STAKE_USER_STATISTIC + address)
        if(!userStatisticObj) return 0
        userStatisticObj = JSON.parse(userStatisticObj)
        return userStatisticObj.staking
    }

    _getAmountLike (address) {
        // get staking amount
        const staking = this._getStaking(address)
        const amountLike = staking * LIKE_RATIO
        if(amountLike > LIKE_MAX) return LIKE_MAX
        if(this._isMaxLikePerDay(address)) return 0
        return amountLike
    }

    post(address, title, content, tag) {

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
            totalLike: 0,
            realLike: 0,
            totalComment: 0,
            totalCommentAndReply: 0,
            totalShare: 0,
            realLikeArray: [0],
            lastBlockWithdraw: block.number
        }

        storage.put(POST_PREFIX + id, JSON.stringify(postObj))
        storage.put(POST_STATISTIC_PREFIX + id, JSON.stringify(postStatisticObj))

        blockchain.receipt(JSON.stringify([id.toString()]))
    }

    share(address, postId, title) {
        this._requireAuth(address, "active")
        this._titleValidate(title)
        // check exist post
        let postStatisticObj = this._checkPostExist(postId)
        let postObj = JSON.parse(storage.get(POST_PREFIX + postId))
        // check postId is share post
        if (postObj.content.type === "share") {
            postId = postObj.content.data
        }
        // +1 share to postId
        postStatisticObj.totalShare++
        storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))
        // post
        let id = tx.time
        postObj = {
            postId: id.toString(),
            time: id,
            title: title,
            content: {
                type: "share",
                data: postId
            },
            tag: [],
        }

        postStatisticObj = {
            author: address,
            totalLike: 0,
            realLike: 0,
            totalComment: 0,
            totalCommentAndReply: 0,
            totalShare: 0,
            realLikeArray: [0],
            lastBlockWithdraw: block.number
        }

        storage.put(POST_PREFIX + id, JSON.stringify(postObj))
        storage.put(POST_STATISTIC_PREFIX + id, JSON.stringify(postStatisticObj))

        blockchain.receipt(JSON.stringify([id.toString(), postId]))
    }

    like(address, postId) {
        this._requireAuth(address, "active")
        // check exist post
        let postStatisticObj = this._checkPostExist(postId)
        // check author like
        if (address === postStatisticObj.author) {
            throw new Error("You can't like own post > " + postStatisticObj.author)
        }

        // check is exist like on this post
        if (storage.mapHas(LIKE_PREFIX + postId, address)) {
            throw new Error("You have been like this postId > " + postId)
        }

        // get amount like by staking in stake.empow
        const amountLike = this._getAmountLike(address)

        if(amountLike > 0) {
            postStatisticObj = this._updateLikeReward(postStatisticObj, amountLike)
            this._updateLikeStatisticForUser(address)
        }

        // update post statistic
        postStatisticObj.totalLike++
        storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))

        // update liked
        storage.mapPut(LIKE_PREFIX + postId, address, JSON.stringify(amountLike))

        blockchain.receipt(JSON.stringify([address, postId]))
    }

    likeWithdraw(postId) {
        // check exist post
        let postStatisticObj = this._checkPostExist(postId)
        let address = postStatisticObj.author

        this._requireAuth(address, "active")

        // calc
        const result = this._calcCanWithdraw(postStatisticObj)
        const canWithdraw = result.canWithdraw
        const totalDayLike = result.totalDayLike

        blockchain.withdraw(address, canWithdraw.toFixed(8), "like withdraw: " + postId)

        // save post statistic
        postStatisticObj.realLikeArray = postStatisticObj.realLikeArray.slice(totalDayLike, totalDayLike + 1)
        postStatisticObj.lastBlockWithdraw = block.number

        // reward vote point
        blockchain.callWithAuth("vote.empow", "issueVotePoint", [address, canWithdraw.toFixed(8)])

        storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))

        blockchain.receipt(JSON.stringify([address, postId, canWithdraw.toFixed(8)]))
    }

    comment(address, postId, type, parentId, content, attachment) {
        this._requireAuth(address, "active")
        // check exist post
        let postStatisticObj = this._checkPostExist(postId)

        if (type !== "comment" && type !== "reply") {
            throw new Error("Wrong comment type > " + type)
        }

        if (content.length <= 0) {
            throw new Error("Content not blank > " + content)
        }

        let commentId = postStatisticObj.totalComment
        let subCommentId = 0

        if (type === "comment") {
            const commentObj = {
                address: address,
                type: "comment",
                time: tx.time,
                postId: postId,
                commentId: postStatisticObj.totalComment,
                parentId: -1,
                totalReply: 0,
                realLike: 0,
                totalLike: 0,
                realLikeArray: [0],
                lastBlockWithdraw: block.number,
                content: content,
                attachment: attachment
            }
            storage.put(COMMENT_PREFIX + postId + "_" + postStatisticObj.totalComment, JSON.stringify(commentObj))
            postStatisticObj.totalComment++
            postStatisticObj.totalCommentAndReply++
            storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))
        } else {
            // check exist comment
            let commentObj = storage.get(COMMENT_PREFIX + postId + "_" + parentId)
            if (!commentObj) {
                throw new Error("Comment not exist > " + parentId)
            }

            commentObj = JSON.parse(commentObj)

            commentId = parentId
            subCommentId = commentObj.totalReply

            let subCommentObj = {
                address: address,
                time: tx.time,
                type: "reply",
                postId: postId,
                commentId: subCommentId,
                parentId: parentId,
                content: content,
                attachment: attachment
            }
            storage.put(REPLY_COMMENT_PREFIX + postId + "_" + parentId + "_" + commentObj.totalReply, JSON.stringify(subCommentObj))
            commentObj.totalReply++
            storage.put(COMMENT_PREFIX + postId + "_" + parentId, JSON.stringify(commentObj))
            postStatisticObj.totalCommentAndReply++
            storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))
        }

        blockchain.receipt(JSON.stringify([type, postId, commentId, subCommentId]))
    }

    likeComment(address, postId, commentId) {
        this._requireAuth(address, "active")
        this._checkPostExist(postId)
        let commentObj = this._checkCommentExist(postId, commentId)

        // check can't like own comment
        if (commentObj.address === address) {
            throw new Error("you can't like your comment > " + postId + "_" + commentId)
        }
        // check liked this comment
        if (storage.mapHas(LIKE_COMMENT_PREFIX + postId + "_" + commentId, address)) {
            throw new Error("you have been like this comment > " + postId + "_" + commentId)
        }

        // get amount like by staking in stake.empow
        const amountLike = this._getAmountLike(address)

        if(amountLike > 0) {
            commentObj = this._updateLikeReward(commentObj, amountLike)
            this._updateLikeStatisticForUser(address)
        }

        // update commentObj
        commentObj.totalLike++
        storage.put(COMMENT_PREFIX + postId + "_" + commentId, JSON.stringify(commentObj))

        // update liked
        storage.mapPut(LIKE_COMMENT_PREFIX + postId + "_" + commentId, address, JSON.stringify(amountLike))

        blockchain.receipt(JSON.stringify([address, postId, commentId]))
    }

    likeCommentWithdraw(postId, commentId) {
        this._checkPostExist(postId)
        let commentObj = this._checkCommentExist(postId, commentId)

        this._requireAuth(commentObj.address, "active")
        // calc

        const result = this._calcCanWithdraw(commentObj)
        const canWithdraw = result.canWithdraw
        const totalDayLike = result.totalDayLike

        blockchain.withdraw(address, canWithdraw.toFixed(8), "like comment withdraw: " + postId + "_" + commentId)

        // save post statistic
        commentId.realLikeArray = commentId.realLikeArray.slice(totalDayLike, totalDayLike + 1)
        commentId.lastBlockWithdraw = block.number

        // reward vote point
        blockchain.callWithAuth("vote.empow", "issueVotePoint", [address, canWithdraw.toFixed(8)])

        storage.put(COMMENT_PREFIX + postId + "_" + commentId, JSON.stringify(commentObj))

        blockchain.receipt(JSON.stringify([address, postId, commentId, canWithdraw.toFixed(8)]))
    }

    report(address, postId, tag) {
        this._requireAuth(address, "active")
        // check post exist
        let postStatisticObj = this._checkPostExist(postId)
        // check tag exist
        const reportTagArray = JSON.parse(storage.get(REPORT_TAG_ARRAY))

        if (reportTagArray.indexOf(tag) === -1) {
            throw new Error("report tag not exist > " + tag)
        }

        // check post is tagged
        if (postStatisticObj[tag]) {
            throw new Error("this post has been tagged > " + tag)
        }

        // check post is exist on reportPendingArray
        if (postStatisticObj.inReportPending) {
            throw new Error("this post in report pending > " + postId)
        }

        // clear old validators
        const isExistOldValidators = storage.get(REPORT_VALIDATOR_PREFIX + postId + "_" + tag)
        if(isExistOldValidators) {
            storage.del(REPORT_VALIDATOR_PREFIX + postId + "_" + tag)
        }

        // push to reportPendingArray
        let reportObj = {
            address: address,
            postId: postId,
            tag: tag,
            time: tx.time
        }

        let reportPendingArray = storage.get(REPORT_PENDING_ARRAY)

        if (!reportPendingArray) {
            reportPendingArray = []
        } else {
            reportPendingArray = JSON.parse(reportPendingArray)
        }

        reportPendingArray.push(reportObj)

        storage.put(REPORT_PENDING_ARRAY, JSON.stringify(reportPendingArray))

        // update postStatisticObj
        postStatisticObj.inReportPending = tag
        storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))

        blockchain.receipt(JSON.stringify([address, postId, tag]))
    }

    verifyPost(address, postId, status) {
        this._requireAuth(address, "active")
        // check post exist
        let postStatisticObj = this._checkPostExist(postId)
        if (!postStatisticObj.inReportPending) {
            throw new Error("this post not in report pending")
        }

        // check minimum staking
        const staking = this._getStaking(address)
        if(staking < REPORT_VALIDATOR_MINIMUM_STAKING) {
            throw new Error("you need staking 10000 EM to become verifier")
        }

        // get validator info
        let reportValidator = storage.get(REPORT_VALIDATOR_PREFIX + postId + "_" + postStatisticObj.inReportPending)

        if (!reportValidator) {
            reportValidator = []
        } else {
            reportValidator = JSON.parse(reportValidator)
        }

        if (reportValidator.length >= REPORT_VALIDATOR_PER_POST) {
            throw new Error("this post has enough verifier")
        }

        for (let i = 0; i < reportValidator.length; i++) {
            if (reportValidator[i].address === address) {
                throw new Error("you have already verify this post")
            }
        }

        // push validator
        reportValidator.push({
            address: address,
            status: status
        })

        storage.put(REPORT_VALIDATOR_PREFIX + postId + "_" + postStatisticObj.inReportPending, JSON.stringify(reportValidator))

        blockchain.receipt(JSON.stringify([address, postId, status]))
    }

    delete(postId) {
        // check exist post
        let postStatisticObj = storage.get(POST_STATISTIC_PREFIX + postId)
        if (!postStatisticObj) {
            throw new Error("PostId not exist > " + postId)
        }
        postStatisticObj = JSON.parse(postStatisticObj)
        let address = postStatisticObj.author

        const admin = storage.get("adminAddress");
        const whitelist = [admin, address];
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

        // delete
        postStatisticObj.deleted = true
        storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))

        blockchain.receipt(JSON.stringify([postId]))
    }

    follow(address, target) {
        this._requireAuth(address, "active")
        // check target exist
        if (!storage.globalMapHas("auth.empow", "auth", target)) {
            throw new Error("Follow user not exist > " + target)
        }
        //
        if (address === target) {
            throw new Error("Can't follow yoursefl > " + target)
        }
        // check is following
        if (storage.mapHas(USER_FOLLOW_PREFIX + target, address)) {
            throw new Error("You are following this user > " + target)
        }
        // set follow
        storage.mapPut(USER_FOLLOW_PREFIX + target, address, "true")

        blockchain.receipt(JSON.stringify([address, target]))
    }

    unfollow(address, target) {
        this._requireAuth(address, "active")
        // check target exist
        if (!storage.globalMapHas("auth.empow", "auth", target)) {
            throw new Error("Follow user not exist > " + target)
        }
        // check is following
        if (!storage.mapHas(USER_FOLLOW_PREFIX + target, address)) {
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

    blockContent(address, tagArray) {
        this._requireAuth(address, "active")

        if (typeof tagArray.length !== "number") {
            throw new Error("tagArray not correct > " + JSON.stringify(tagArray))
        }

        for (let i = 0; i < tagArray.length; i++) {
            const tag = tagArray[i]
            this._checkTagExistInReportTagArray(tag)
        }

        storage.put(BLOCK_TAG_PREFIX + address, JSON.stringify(tagArray))

        blockchain.receipt(JSON.stringify([address, tagArray]))
    }

    // admin only
    _assignReportTag(postId, tag) {
        let postStatisticObj = JSON.parse(storage.get(POST_STATISTIC_PREFIX + postId))
        postStatisticObj[tag] = true
        delete postStatisticObj.inReportPending
        storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))
    }

    _removeInReportPending(postId) {
        let postStatisticObj = JSON.parse(storage.get(POST_STATISTIC_PREFIX + postId))
        delete postStatisticObj.inReportPending
        storage.put(POST_STATISTIC_PREFIX + postId, JSON.stringify(postStatisticObj))
    }

    _sendRewardToValidator(postId, validators) {
        let balance = new Float64(blockchain.callWithAuth("token.empow", "balanceOf", ["em", blockchain.contractName()])[0])

        if (balance.lt(REPORT_VALIDATOR_REWARD * validators.length)) return

        for (let i = 0; i < validators.length; i++) {
            blockchain.callWithAuth("token.empow", "transfer", ["em", blockchain.contractName(), validators[i], REPORT_VALIDATOR_REWARD.toString(), "verify reward: " + postId])
        }
    }

    _checkValidate(postId, tag) {
        let validator = storage.get(REPORT_VALIDATOR_PREFIX + postId + "_" + tag)
        if (!validator) return { done: false, status: false }
        validator = JSON.parse(validator)
        if (validator.length < REPORT_VALIDATOR_PER_POST) return { done: false, status: false };

        let countTrue = 0
        let validatorTrue = []
        let validatorFalse = []

        for (let i = 0; i < validator.length; i++) {
            if (validator[i].status === true) {
                countTrue++
                validatorTrue.push(validator[i].address)
            } else {
                validatorFalse.push(validator[i].address)
            }
        }

        if (countTrue >= Math.round(REPORT_VALIDATOR_PER_POST / 2)) {
            // assign tag for post
            this._assignReportTag(postId, tag)
            // send reward for validator
            this._sendRewardToValidator(postId, validatorTrue)
            return { done: true, status: true }
        } else {
            this._removeInReportPending(postId)
            this._sendRewardToValidator(postId, validatorFalse)
            return { done: true, status: false }
        }
    }

    _updateReportPendingArray(reportPendingArray, validatedCount) {
        if (validatedCount === 0) return;

        for (let i = 0; i < validatedCount; i++) {
            reportPendingArray.shift()
        }

        storage.put(REPORT_PENDING_ARRAY, JSON.stringify(reportPendingArray))
    }

    checkReport() {
        const admin = storage.get("adminAddress");
        const whitelist = ["base.empow", admin];
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

        // get report pending array
        let reportPendingArray = storage.get(REPORT_PENDING_ARRAY)
        if (!reportPendingArray) return
        reportPendingArray = JSON.parse(reportPendingArray)

        let postValidated = []

        for (let i = 0; i < reportPendingArray.length; i++) {
            const postId = reportPendingArray[i].postId
            const postStatisticObj = JSON.parse(storage.get(POST_STATISTIC_PREFIX + postId))

            if (postStatisticObj.deleted) {
                postValidated.push(postId)
                continue
            }

            const tag = postStatisticObj.inReportPending
            const result = this._checkValidate(postId, tag)

            if (!result.done) {
                this._updateReportPendingArray(reportPendingArray, postValidated.length)
                if (postValidated.length > 0) blockchain.receipt(JSON.stringify(postValidated))
                return;
            }

            if (result.done) {
                postValidated.push(postId)
            }
        }

        this._updateReportPendingArray(reportPendingArray, postValidated.length)
        if (postValidated.length > 0) blockchain.receipt(JSON.stringify(postValidated))
    }

    topup(amount) {
        this._requireAuth("base.empow", "active")
        amount = this._fixAmount(amount);

        // calc interest per 1 like
        let totalLike = storage.get(TOTAl_LIKE)
        if (Math.floor(totalLike) === 0) {
            totalLike = "1"
        }
        let interest = amount
        if (totalLike !== "0") {
            interest = amount.div(totalLike)
        }

        // update interest
        const bn = block.number
        const currentDay = Math.floor(bn / BLOCK_NUMBER_PER_DAY)
        storage.put(INTEREST_PREEFIX + currentDay, interest.toFixed(8))

        // reset total like
        this._resetTotalLike()
    }

    addReportTag(tag) {
        const admin = storage.get("adminAddress");
        this._requireAuth(admin, "active")

        let reportTagArray = storage.get(REPORT_TAG_ARRAY)

        if (reportTagArray) {
            reportTagArray = JSON.parse(reportTagArray)

            if (reportTagArray.indexOf(tag) !== -1) {
                throw new Error("tag is exist > " + tag)
            }

            reportTagArray.push(tag)
        } else {
            reportTagArray = []
        }

        storage.put(REPORT_TAG_ARRAY, JSON.stringify(reportTagArray))
        blockchain.receipt(JSON.stringify([tag]))
    }
}

module.exports = Social
