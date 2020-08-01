package account

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHashSha1(t *testing.T) {
	Convey("hash", t, func() {

		Convey("sha1", func() {

			So(hashSha1("123456"), ShouldEqual, "7c4a8d09ca3762af61e59520943dc26494f8941b")
			So(hashSha1("hello world"), ShouldEqual, "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed")
		})

		Convey("md5", func() {
			So(hashMd5("123456"), ShouldEqual, "e10adc3949ba59abbe56e057f20f883e")
			So(hashMd5("hello world"), ShouldEqual, "5eb63bbbe01eeed093cb22bb8f5acdc3")
		})

		Convey("pwdFromMd5", func() {
			So(HashPwdFromMd5("e10adc3949ba59abbe56e057f20f883e", "salt"), ShouldEqual, "faeac6cf04aade89a64d7cabe5e921a47fdf21e9")
			So(HashPwdFromMd5("5eb63bbbe01eeed093cb22bb8f5acdc3", "salt"), ShouldEqual, "9c2414acf2d5c5e6472af9049679ec2f56ca34df")

			So(HashPwdFromMd5("e10adc3949ba59abbe56e057f20f883e", "NewSalt"), ShouldEqual, "31d13cbeeb64b910f81ff8a0c390526ae8b9d9b0")
			So(HashPwdFromMd5("5eb63bbbe01eeed093cb22bb8f5acdc3", "NewSalt"), ShouldEqual, "40460f97392f90fe16a35c2d6ff1483d0c56337e")

			// default admin password
			So(HashPwdFromMd5("e10adc3949ba59abbe56e057f20f883e", "HelloAdm"), ShouldEqual, "fe9ecdd819e15d679294f3f4c3a2e4c9bebe8787")
		})

		Convey("pwd", func() {
			So(HashPwd("123456", "salt"), ShouldEqual, "faeac6cf04aade89a64d7cabe5e921a47fdf21e9")
			So(HashPwd("hello world", "salt"), ShouldEqual, "9c2414acf2d5c5e6472af9049679ec2f56ca34df")

			So(HashPwd("123456", "NewSalt"), ShouldEqual, "31d13cbeeb64b910f81ff8a0c390526ae8b9d9b0")
			So(HashPwd("hello world", "NewSalt"), ShouldEqual, "40460f97392f90fe16a35c2d6ff1483d0c56337e")

			// default admin password
			So(HashPwd("123456", "HelloAdm"), ShouldEqual, "fe9ecdd819e15d679294f3f4c3a2e4c9bebe8787")
		})

	})
}
