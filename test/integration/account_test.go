package integration

// import (
// 	. "github.com/smartystreets/goconvey/convey"
// )

// func TestAccountInfo(t *testing.T) {
// 	ilog.SetLevel(ilog.LevelInfo)
// 	s := verifier.NewSimulator()
// 	defer s.Clear()
// 	Convey("test of Auth", t, func() {
// 		ca, err := s.Compile("auth.empow", "../../config/genesis/contract/account", "../../config/genesis/contract/account.js")
// 		So(err, ShouldBeNil)
// 		s.Visitor.SetContract(ca)
// 		s.Visitor.SetContract(native.GasABI())
// 		s.Visitor.SetContract(native.TokenABI())

// 		acc := prepareAuth(t, s)
// 		s.SetGas(acc.ID, 1e8)
// 		s.SetGas("myidid", 1e8)
// 		s.SetRAM(acc.ID, 1000)
// 		s.SetRAM("myidid", 1000)
// 		err = createToken(t, s, acc)
// 		So(err, ShouldBeNil)

// 		ilog.Info(acc.ID, acc.KeyPair)
// 		r, err := s.Call("auth.empow", "signUp", array2json([]interface{}{"myidid", acc.KeyPair.ReadablePubkey(), acc.KeyPair.ReadablePubkey()}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(r.Status.Message, ShouldEqual, "")
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldStartWith, `{"id":"myidid",`)

// 		r, err = s.Call("auth.empow", "signUp", array2json([]interface{}{"invalid#id", acc.KeyPair.ReadablePubkey(), acc.KeyPair.ReadablePubkey()}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		ilog.Info(r.Status.Message)
// 		So(r.Status.Message, ShouldContainSubstring, "id contains invalid character")

// 		acc1, _ := host.ReadAuth(s.Visitor, "myidid")
// 		So(acc1.Referrer, ShouldEqual, acc.ID)

// 		r, err = s.Call("auth.empow", "addPermission", array2json([]interface{}{"myidid", "perm1", 1}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(r.Status.Message, ShouldEqual, "")
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldContainSubstring, `"perm1":{"name":"perm1","groups":[],"items":[],"threshold":1}`)

// 		r, err = s.Call("auth.empow", "addPermission", array2json([]interface{}{"myidid", "perm1", 1}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(r.Status.Message, ShouldContainSubstring, "permission already exist")

// 		r, err = s.Call("auth.empow", "dropPermission", array2json([]interface{}{"myidid", "perm1"}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(r.Status.Message, ShouldEqual, "")
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldNotContainSubstring, `"perm1":{"name":"perm1","groups":[],"items":[],"threshold":1}`)

// 		r, err = s.Call("auth.empow", "assignPermission", array2json([]interface{}{"myidid", "active", "@acc1", 1}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(r.Status.Message, ShouldContainSubstring, "unexpected item")

// 		r, err = s.Call("auth.empow", "assignPermission", array2json([]interface{}{"myidid", "active", "EMPOW1234", 1}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldContainSubstring, `{"id":"EMPOW1234","is_key_pair":true,"weight":1}`)

// 		r, err = s.Call("auth.empow", "assignPermission", array2json([]interface{}{"myidid", "active", "acc1@active", 1}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldContainSubstring, `{"id":"acc1","permission":"active","is_key_pair":false,"weight":1}`)

// 		r, err = s.Call("auth.empow", "revokePermission", array2json([]interface{}{"myidid", "active", "acc1@active"}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldNotContainSubstring, `{"id":"acc1","permission":"active","is_key_pair":false,"weight":1}`)

// 		r, err = s.Call("auth.empow", "revokePermission", array2json([]interface{}{"myidid", "active", "acc2"}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(r.Status.Message, ShouldContainSubstring, "item not found")

// 		r, err = s.Call("auth.empow", "revokePermission", array2json([]interface{}{"myidid", "active", "EMPOW1234"}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldNotContainSubstring, `{"id":"EMPOW1234","is_key_pair":true,"weight":1}`)

// 		r, err = s.Call("auth.empow", "addGroup", array2json([]interface{}{"myidid", "grp0"}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldContainSubstring, `"groups":{"grp0":{"name":"grp0","items":[]}}`)

// 		r, err = s.Call("auth.empow", "addGroup", array2json([]interface{}{"myidid", "grp0"}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(r.Status.Message, ShouldContainSubstring, "group already exist")

// 		r, err = s.Call("auth.empow", "assignGroup", array2json([]interface{}{"myidid", "grp0", "acc1@active", 1}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldContainSubstring, `"groups":{"grp0":{"name":"grp0","items":[{"id":"acc1","permission":"active","is_key_pair":false,"weight":1}]}`)

// 		r, err = s.Call("auth.empow", "revokeGroup", array2json([]interface{}{"myidid", "grp0", "acc1@active"}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldNotContainSubstring, `"groups":{"grp0":{"name":"grp0","items":[{"id":"acc1","permission":"active","is_key_pair":false,"weight":1}]}`)

// 		r, err = s.Call("auth.empow", "assignPermissionToGroup", array2json([]interface{}{"myidid", "active", "grp0"}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldContainSubstring, `"groups":["grp0"]`)

// 		r, err = s.Call("auth.empow", "revokePermissionInGroup", array2json([]interface{}{"myidid", "active", "grp0"}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldNotContainSubstring, `"groups":["grp0"]`)

// 		r, err = s.Call("auth.empow", "dropGroup", array2json([]interface{}{"myidid", "grp0"}), acc.ID, acc.KeyPair)
// 		So(err, ShouldBeNil)
// 		So(database.Unmarshal(s.Visitor.MGet("auth.empow-auth", "myidid")), ShouldNotContainSubstring, `"groups":{"grp0":{"name":"grp0","items":[]}}`)

// 	})
// }
