package integration

import (
	"fmt"
	"testing"

	"github.com/empow-blockchain/go-empow/ilog"

	"github.com/empow-blockchain/go-empow/core/tx"
	. "github.com/empow-blockchain/go-empow/verifier"
	"github.com/empow-blockchain/go-empow/vm/database"
	. "github.com/smartystreets/goconvey/convey"
)

func prepareSocial(t *testing.T, s *Simulator, acc *TestAccount) {
	s.Head.Number = 0

	// deploy social.empow
	setNonNativeContract(s, "social.empow", "social.js", ContractPath)

	r, err := s.Call("social.empow", "init", `[]`, acc0.ID, acc0.KeyPair)
	if err != nil || r.Status.Code != tx.Success {
		t.Fatal(err, r)
	}

	r, err = s.Call("social.empow", "initAdmin", fmt.Sprintf(`["%s"]`, acc0.ID), acc0.ID, acc0.KeyPair)
	if err != nil || r.Status.Code != tx.Success {
		t.Fatal(err, r)
	}

	// deploy vote point
	setNonNativeContract(s, "vote_point.empow", "vote_point.js", ContractPath)

	r, err = s.Call("vote_point.empow", "initAdmin", fmt.Sprintf(`["%s"]`, acc0.ID), acc0.ID, acc0.KeyPair)
	if err != nil || r.Status.Code != tx.Success {
		t.Fatal(err, r)
	}
}

func Test_post(t *testing.T) {
	ilog.Stop()

	Convey("test post validate", t, func() {
		s := NewSimulator()
		defer s.Clear()

		createAccountsWithResource(s)
		prepareSocial(t, s, acc0)

		r, err := s.Call("social.empow", "post", fmt.Sprintf(`["%v","test title", "test content",
		  	["#test", "#flower"]
		  ]`, acc2.ID), acc2.ID, acc2.KeyPair)

		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldContainSubstring, "content must be array")

		r, err = s.Call("social.empow", "post", fmt.Sprintf(`["%v","test title", [
			{
			  "type": "photo",
			  "data": "https://cpmr-islands.org/wp-content/uploads/sites/4/2019/07/Happy-Test-Screen-01-825x510.png"
			}
		  ],
		  	{"test":"test","flower":"flower"}
		  ]`, acc2.ID), acc2.ID, acc2.KeyPair)

		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldContainSubstring, "tag must be array")

		r, err = s.Call("social.empow", "post", fmt.Sprintf(`["%v","AM1xXp3YremFaSqH5P580STfUQUjkehJUlXlIs2uYCjnaIAZcHk8VpUcgWjQPHnXoAqVhAqjwbaspjWKtDGjumb1sva9zv5lxN0So9AoVlB7FSkN40Fn6EisNsGLJ5ryxCLFiS02Ymcod9NHQ94oyISXxrpznqILKCpXx3yzSxFw677yc9MJ9j6HJhCHVHmj0lxDlOPfKLwlJ88OyK4Zy3SJvD5ctMiW4GrPa4r5xPjSYbQyNo3aNKqacWuooRLP", [
			{
			  "type": "photo",
			  "data": "https://cpmr-islands.org/wp-content/uploads/sites/4/2019/07/Happy-Test-Screen-01-825x510.png"
			}
		  ],
		  	["tag", "flower"]
		  ]`, acc2.ID), acc2.ID, acc2.KeyPair)

		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldContainSubstring, "title must be length greater than 0 and less than 255")

		r, err = s.Call("social.empow", "post", fmt.Sprintf(`["%v","test title", [
			{
			  "type": "photo",
			  "data": "https://cpmr-islands.org/wp-content/uploads/sites/4/2019/07/Happy-Test-Screen-01-825x510.png"
			}
		  ],
		  	["tag", "flower"]
		  ]`, acc2.ID), acc2.ID, acc2.KeyPair)

		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-p_%v", s.Head.Time))), ShouldEqual, `{"time":1541541540000000000,"title":"test title","content":[{"data":"https://cpmr-islands.org/wp-content/uploads/sites/4/2019/07/Happy-Test-Screen-01-825x510.png","type":"photo"}],"tag":["tag","flower"]}`)
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", s.Head.Time))), ShouldEqual, `{"author":"user_2","totalLike":0,"realLike":0,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[0],"lastBlockWithdraw":0}`)
		So(r.GasUsage, ShouldEqual, 3930800)
	})
}

func Test_Like(t *testing.T) {
	ilog.Stop()

	Convey("test like", t, func() {
		s := NewSimulator()
		defer s.Clear()

		createAccountsWithResource(s)
		prepareSocial(t, s, acc0)

		r, err := s.Call("social.empow", "post", fmt.Sprintf(`["%v","test title", [
			{
			  "type": "photo",
			  "data": "https://cpmr-islands.org/wp-content/uploads/sites/4/2019/07/Happy-Test-Screen-01-825x510.png"
			}
		  ],
		  	["tag", "flower"]
		  ]`, acc2.ID), acc1.ID, acc1.KeyPair)

		postID := s.Head.Time

		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-p_%v", postID))), ShouldEqual, `{"time":1541541540000000000,"title":"test title","content":[{"data":"https://cpmr-islands.org/wp-content/uploads/sites/4/2019/07/Happy-Test-Screen-01-825x510.png","type":"photo"}],"tag":["tag","flower"]}`)
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":0,"realLike":0,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[0],"lastBlockWithdraw":0}`)
		So(r.GasUsage, ShouldEqual, 3930800)

		s.Head.Time += 1 * 1e9
		s.Head.Number += 1 * 2

		Convey("test withdraw validate", func() {
			r, err = s.Call("social.empow", "likeWithdraw", fmt.Sprintf(`["%v"]`, postID), acc1.ID, acc1.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "require auth")

			r, err = s.Call("social.empow", "likeWithdraw", fmt.Sprintf(`["%v"]`, postID), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "You can withdraw like after 1 day")
		})

		Convey("test like validate", func() {

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "vxcvew"]`, acc3.ID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "PostId not exist")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc2.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "require auth")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc2.ID, postID), acc2.ID, acc2.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "You can't like own post")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldContainSubstring, "You have been like this postId")
		})

		Convey("test like amount lv 1", func() {
			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":1,"realLike":0.01,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[0.01],"lastBlockWithdraw":0}`)
			So(database.MustUnmarshal(s.Visitor.MGet(fmt.Sprintf("social.empow-l_%v", postID), acc3.ID)), ShouldEqual, "0.01")
		})

		Convey("test like amount lv 2", func() {

			r, err = s.Call("social.empow", "upLevel", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, 2), acc0.ID, acc0.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get("social.empow-lv_user_3")), ShouldEqual, "2")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":1,"realLike":10,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[10],"lastBlockWithdraw":0}`)
			So(database.MustUnmarshal(s.Visitor.MGet(fmt.Sprintf("social.empow-l_%v", postID), acc3.ID)), ShouldEqual, "10")
		})

		Convey("test like amount lv 3", func() {

			r, err = s.Call("social.empow", "upLevel", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, 3), acc0.ID, acc0.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get("social.empow-lv_user_3")), ShouldEqual, "3")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":1,"realLike":15,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[15],"lastBlockWithdraw":0}`)
			So(database.MustUnmarshal(s.Visitor.MGet(fmt.Sprintf("social.empow-l_%v", postID), acc3.ID)), ShouldEqual, "15")
		})

		Convey("test like amount lv 4", func() {

			r, err = s.Call("social.empow", "upLevel", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, 4), acc0.ID, acc0.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get("social.empow-lv_user_3")), ShouldEqual, "4")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":1,"realLike":18,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[18],"lastBlockWithdraw":0}`)
			So(database.MustUnmarshal(s.Visitor.MGet(fmt.Sprintf("social.empow-l_%v", postID), acc3.ID)), ShouldEqual, "18")
		})

		Convey("test like amount lv 5", func() {

			r, err = s.Call("social.empow", "upLevel", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, 5), acc0.ID, acc0.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get("social.empow-lv_user_3")), ShouldEqual, "5")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":1,"realLike":20,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[20],"lastBlockWithdraw":0}`)
			So(database.MustUnmarshal(s.Visitor.MGet(fmt.Sprintf("social.empow-l_%v", postID), acc3.ID)), ShouldEqual, "20")
		})

		Convey("test like amount lv 6", func() {

			r, err = s.Call("social.empow", "upLevel", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, 6), acc0.ID, acc0.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get("social.empow-lv_user_3")), ShouldEqual, "6")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":1,"realLike":25,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[25],"lastBlockWithdraw":0}`)
			So(database.MustUnmarshal(s.Visitor.MGet(fmt.Sprintf("social.empow-l_%v", postID), acc3.ID)), ShouldEqual, "25")
		})

		Convey("test like amount lv 7", func() {

			r, err = s.Call("social.empow", "upLevel", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, 7), acc0.ID, acc0.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get("social.empow-lv_user_3")), ShouldEqual, "7")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":1,"realLike":30,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[30],"lastBlockWithdraw":0}`)
			So(database.MustUnmarshal(s.Visitor.MGet(fmt.Sprintf("social.empow-l_%v", postID), acc3.ID)), ShouldEqual, "30")
		})

		Convey("test multi like and withdraw", func() {
			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":1,"realLike":0.01,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[0.01],"lastBlockWithdraw":0}`)
			So(database.MustUnmarshal(s.Visitor.MGet(fmt.Sprintf("social.empow-l_%v", postID), acc3.ID)), ShouldEqual, "0.01")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc4.ID, postID), acc4.ID, acc4.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":2,"realLike":0.02,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[0.02],"lastBlockWithdraw":0}`)

			s.Head.Time += 2 * 24 * 60 * 60 * 1e9
			s.Head.Number += 2 * 24 * 60 * 60 * 2

			r, err = s.Call("social.empow", "upLevel", fmt.Sprintf(`["%v", "%v"]`, acc5.ID, 7), acc0.ID, acc0.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get("social.empow-lv_user_5")), ShouldEqual, "7")

			r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc5.ID, postID), acc5.ID, acc5.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":3,"realLike":30.02,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[0.02,null,30],"lastBlockWithdraw":0}`)
		})

		Convey("test topup and withdraw", func() {
			prepareIssue(s, acc0)
			prepareFakeBase(t, s)
			prepareNewProducerVote(t, s, acc0)
			prepareStake(t, s, acc0)

			for _, acc := range testAccounts {

				r, err = s.Call("social.empow", "upLevel", fmt.Sprintf(`["%v", "%v"]`, acc.ID, 2), acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				if acc.ID == "user_2" {
					continue
				}

				r, err = s.Call("social.empow", "like", fmt.Sprintf(`["%v", "%v"]`, acc.ID, postID), acc.ID, acc.KeyPair)
				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
			}

			So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":9,"realLike":90,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[90],"lastBlockWithdraw":0}`)

			s.Head.Time += 1 * 24 * 60 * 60 * 1e9
			s.Head.Number += 1 * 24 * 60 * 60 * 2

			r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

			So(err, ShouldBeNil)
			So(r.Status.Message, ShouldEqual, "")
			So(s.Visitor.TokenBalance("em", "social.empow"), ShouldEqual, int64(13698632876712))
			So(database.MustUnmarshal(s.Visitor.Get("social.empow-i_1")), ShouldEqual, "1522.07031963")

			Convey("withdraw in day", func() {
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, 0)

				r, err = s.Call("social.empow", "likeWithdraw", fmt.Sprintf(`["%v"]`, postID), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, 90*1e8)
				So(s.Visitor.TokenBalance("em", "social.empow"), ShouldEqual, 13698632876712-90*1e8)
				So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":9,"realLike":90,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[0],"lastBlockWithdraw":172800}`)
				So(database.MustUnmarshal(s.Visitor.Get("social.empow-restAmount")), ShouldEqual, "136896.3287667")
			})

			Convey("withdraw a few day", func() {
				s.Head.Time += 1 * 24 * 60 * 60 * 1e9
				s.Head.Number += 1 * 24 * 60 * 60 * 2

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				s.Head.Time += 1 * 24 * 60 * 60 * 1e9
				s.Head.Number += 1 * 24 * 60 * 60 * 2

				r, err = s.Call("base.empow", "issueEM", `[]`, acc0.ID, acc0.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")

				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, 0)

				r, err = s.Call("social.empow", "likeWithdraw", fmt.Sprintf(`["%v"]`, postID), acc2.ID, acc2.KeyPair)

				So(err, ShouldBeNil)
				So(r.Status.Message, ShouldEqual, "")
				So(s.Visitor.TokenBalance("em", acc2.ID), ShouldEqual, 90*1e8)
				So(s.Visitor.TokenBalance("em", "social.empow"), ShouldEqual, 41091402454777)
				So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":9,"realLike":90,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[0],"lastBlockWithdraw":518400}`)
				So(s.Visitor.TokenBalance("vote", acc2.ID), ShouldEqual, 90*1e8)
			})
		})
	})
}

func Test_Comment(t *testing.T) {
	ilog.Stop()
	Convey("test comment", t, func() {
		s := NewSimulator()
		defer s.Clear()

		createAccountsWithResource(s)
		prepareSocial(t, s, acc0)

		r, err := s.Call("social.empow", "post", fmt.Sprintf(`["%v","test title", [
			{
			  "type": "photo",
			  "data": "https://cpmr-islands.org/wp-content/uploads/sites/4/2019/07/Happy-Test-Screen-01-825x510.png"
			}
		  ],
		  	["tag", "flower"]
		  ]`, acc2.ID), acc1.ID, acc1.KeyPair)

		postID := s.Head.Time

		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-p_%v", postID))), ShouldEqual, `{"time":1541541540000000000,"title":"test title","content":[{"data":"https://cpmr-islands.org/wp-content/uploads/sites/4/2019/07/Happy-Test-Screen-01-825x510.png","type":"photo"}],"tag":["tag","flower"]}`)
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":0,"realLike":0,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[0],"lastBlockWithdraw":0}`)
		So(r.GasUsage, ShouldEqual, 3930800)

		r, err = s.Call("social.empow", "comment", fmt.Sprintf(`["%v", "%v", "comment", 0, "test comment"]`, acc2.ID, postID), acc2.ID, acc2.KeyPair)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":0,"realLike":0,"totalComment":1,"totalCommentAndReply":1,"totalReport":0,"realLikeArray":[0],"lastBlockWithdraw":0}`)
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-c_%v_0", postID))), ShouldEqual, `{"totalReply":0,"content":"test comment"}`)
		So(r.GasUsage, ShouldEqual, 3689000)

		r, err = s.Call("social.empow", "comment", fmt.Sprintf(`["%v", "%v", "reply", 0, "test reply comment"]`, acc2.ID, postID), acc2.ID, acc2.KeyPair)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":0,"realLike":0,"totalComment":1,"totalCommentAndReply":2,"totalReport":0,"realLikeArray":[0],"lastBlockWithdraw":0}`)
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-c_%v_0", postID))), ShouldEqual, `{"totalReply":1,"content":"test comment"}`)
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-rc_%v_0_0", postID))), ShouldEqual, "test reply comment")
		So(r.GasUsage, ShouldEqual, 3756900)
	})
}

func Test_Report(t *testing.T) {
	ilog.Stop()

	Convey("test report", t, func() {
		s := NewSimulator()
		defer s.Clear()

		createAccountsWithResource(s)
		prepareSocial(t, s, acc0)

		r, err := s.Call("social.empow", "post", fmt.Sprintf(`["%v","test title", [
			{
			  "type": "photo",
			  "data": "https://cpmr-islands.org/wp-content/uploads/sites/4/2019/07/Happy-Test-Screen-01-825x510.png"
			}
		  ],
		  	["tag", "flower"]
		  ]`, acc2.ID), acc1.ID, acc1.KeyPair)

		postID := s.Head.Time

		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-p_%v", postID))), ShouldEqual, `{"time":1541541540000000000,"title":"test title","content":[{"data":"https://cpmr-islands.org/wp-content/uploads/sites/4/2019/07/Happy-Test-Screen-01-825x510.png","type":"photo"}],"tag":["tag","flower"]}`)
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":0,"realLike":0,"totalComment":0,"totalCommentAndReply":0,"totalReport":0,"realLikeArray":[0],"lastBlockWithdraw":0}`)
		So(r.GasUsage, ShouldEqual, 3930800)

		r, err = s.Call("social.empow", "report", fmt.Sprintf(`["%v", "%v", "18plus"]`, acc2.ID, postID), acc2.ID, acc2.KeyPair)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldContainSubstring, "report tag not exist")

		r, err = s.Call("social.empow", "addReportTag", `["18plus"]`, acc2.ID, acc2.KeyPair)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get("social.empow-reportTagArray")), ShouldEqual, `["18plus"]`)

		r, err = s.Call("social.empow", "addReportTag", `["violence"]`, acc2.ID, acc2.KeyPair)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.Get("social.empow-reportTagArray")), ShouldEqual, `["18plus","violence"]`)

		r, err = s.Call("social.empow", "report", fmt.Sprintf(`["%v", "%v", "18plus"]`, acc2.ID, postID), acc2.ID, acc2.KeyPair)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.MGet(fmt.Sprintf("social.empow-r_%v", postID), "18plus")), ShouldEqual, "1")
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-r_%v_%v", postID, acc2.ID))), ShouldEqual, "true")

		r, err = s.Call("social.empow", "report", fmt.Sprintf(`["%v", "%v", "18plus"]`, acc2.ID, postID), acc2.ID, acc2.KeyPair)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldContainSubstring, "can report 2 times")

		r, err = s.Call("social.empow", "report", fmt.Sprintf(`["%v", "%v", "18plus"]`, acc3.ID, postID), acc3.ID, acc3.KeyPair)
		So(err, ShouldBeNil)
		So(r.Status.Message, ShouldEqual, "")
		So(database.MustUnmarshal(s.Visitor.MGet(fmt.Sprintf("social.empow-r_%v", postID), "18plus")), ShouldEqual, "2")
		So(database.MustUnmarshal(s.Visitor.Get(fmt.Sprintf("social.empow-s_%v", postID))), ShouldEqual, `{"author":"user_2","totalLike":0,"realLike":0,"totalComment":0,"totalCommentAndReply":0,"totalReport":2,"realLikeArray":[0],"lastBlockWithdraw":0}`)
		So(r.GasUsage, ShouldEqual, 3776700)
	})
}
