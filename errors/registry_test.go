package errors

import (
	"errors"
	"math/bits"
	"testing"

	. "github.com/jucardi/go-testx/testx"
)

type someErrorImpl struct {
	msg string
}

func (s *someErrorImpl) Error() string { return s.msg }

type someErrWithCode struct {
	msg  string
	code ErrFlag
}

func (s *someErrWithCode) Error() string              { return s.msg }
func (s *someErrWithCode) Flags() ErrFlag             { return s.code }
func (s *someErrWithCode) HasFlags(code ErrFlag) bool { return HasFlags(s, code) }

func TestRegistry(t *testing.T) {
	Convey("Testing the error registry", t, func() {
		err1 := errors.New("error1")
		err2 := errors.New("error2")
		err3 := &someErrorImpl{msg: "error3"}
		err4 := &someErrorImpl{msg: "error4"}

		Convey("Ensuring all test errors are not registered", t, func() {
			ShouldBeFalse(Match(err1, FlagNotFound))
			ShouldBeFalse(Match(err2, FlagNotFound))
			ShouldBeFalse(Match(err3, FlagNotFound))
			ShouldBeFalse(Match(err4, FlagNotFound))
		})

		Convey("Registering an instance and attempting to match codes", t, func() {
			ShouldEqual(0, bits.OnesCount(uint(errRegistry[err1])))
			SetFlagsByInstance(err1, FlagNotFound)
			ShouldBeTrue(Match(err1, FlagNotFound))
			ShouldBeFalse(Match(err2, FlagNotFound))
			ShouldBeFalse(Match(err3, FlagNotFound))
			ShouldBeFalse(Match(err4, FlagNotFound))
			ShouldEqual(1, bits.OnesCount(uint(errRegistry[err1])))

			SetFlagsByInstance(err1, FlagUnauthorized)
			ShouldBeTrue(Match(err1, FlagNotFound))
			ShouldBeFalse(Match(err2, FlagNotFound))
			ShouldBeFalse(Match(err3, FlagNotFound))
			ShouldBeFalse(Match(err4, FlagNotFound))
			ShouldBeTrue(Match(err1, FlagUnauthorized))
			ShouldBeFalse(Match(err2, FlagUnauthorized))
			ShouldBeFalse(Match(err3, FlagUnauthorized))
			ShouldBeFalse(Match(err4, FlagUnauthorized))
			ShouldEqual(2, bits.OnesCount(uint(errRegistry[err1])))
		})

		Convey("Registering an error type and attempting to match codes", t, func() {
			ShouldNotError(
				SetFlagsByType((*someErrorImpl)(nil), FlagNotFound),
			)
			ShouldBeTrue(Match(err1, FlagNotFound))
			ShouldBeFalse(Match(err2, FlagNotFound))
			ShouldBeTrue(Match(err3, FlagNotFound))
			ShouldBeTrue(Match(err4, FlagNotFound))

			ShouldNotError(
				SetFlagsByType((*someErrorImpl)(nil), FlagUnauthorized),
			)
			ShouldBeTrue(Match(err1, FlagNotFound))
			ShouldBeFalse(Match(err2, FlagNotFound))
			ShouldBeTrue(Match(err3, FlagNotFound))
			ShouldBeTrue(Match(err4, FlagNotFound))
			ShouldBeTrue(Match(err1, FlagUnauthorized))
			ShouldBeFalse(Match(err2, FlagUnauthorized))
			ShouldBeTrue(Match(err3, FlagUnauthorized))
			ShouldBeTrue(Match(err4, FlagUnauthorized))

		})

		Convey("Testing unsetting error flags for error instances", t, func() {
			ShouldBeTrue(Match(err1, FlagNotFound))
			ShouldBeTrue(Match(err1, FlagUnauthorized))

			UnsetFlagsByInstance(err1, FlagNotFound)
			ShouldBeFalse(Match(err1, FlagNotFound))
			ShouldBeTrue(Match(err1, FlagUnauthorized))

			UnsetFlagsByInstance(err1, FlagUnauthorized)
			ShouldBeFalse(Match(err1, FlagNotFound))
			ShouldBeFalse(Match(err1, FlagUnauthorized))
		})

		Convey("Testing unsetting error flags for error types", t, func() {
			ShouldBeTrue(Match(err3, FlagNotFound))
			ShouldBeTrue(Match(err3, FlagUnauthorized))
			ShouldBeTrue(Match(err4, FlagNotFound))
			ShouldBeTrue(Match(err4, FlagUnauthorized))

			ShouldNotError(
				UnsetFlagsByType(err4, FlagNotFound),
			)
			ShouldBeFalse(Match(err3, FlagNotFound))
			ShouldBeTrue(Match(err3, FlagUnauthorized))
			ShouldBeFalse(Match(err4, FlagNotFound))
			ShouldBeTrue(Match(err4, FlagUnauthorized))

			ShouldNotError(
				UnsetFlagsByType(err4, FlagUnauthorized),
			)
			ShouldBeFalse(Match(err3, FlagNotFound))
			ShouldBeFalse(Match(err3, FlagUnauthorized))
			ShouldBeFalse(Match(err4, FlagNotFound))
			ShouldBeFalse(Match(err4, FlagUnauthorized))
		})

		Convey("Testing deregister for error instances", t, func() {
			ShouldBeFalse(Match(err1, FlagNotFound))
			ShouldBeFalse(Match(err1, FlagUnauthorized))
			ShouldBeFalse(Match(err1, FlagUnhandled))

			SetFlagsByInstance(err1, FlagNotFound)
			SetFlagsByInstance(err1, FlagUnauthorized)
			SetFlagsByInstance(err1, FlagUnhandled)

			ShouldBeTrue(Match(err1, FlagNotFound))
			ShouldBeTrue(Match(err1, FlagUnauthorized))
			ShouldBeTrue(Match(err1, FlagUnhandled))

			DeregisterInstance(err1)
			ShouldBeFalse(Match(err1, FlagNotFound))
			ShouldBeFalse(Match(err1, FlagUnauthorized))
			ShouldBeFalse(Match(err1, FlagUnhandled))
		})

		Convey("Testing deregister for error types", t, func() {
			ShouldBeFalse(Match(err3, FlagNotFound))
			ShouldBeFalse(Match(err3, FlagUnauthorized))
			ShouldBeFalse(Match(err3, FlagUnhandled))
			ShouldBeFalse(Match(err4, FlagNotFound))
			ShouldBeFalse(Match(err4, FlagUnauthorized))
			ShouldBeFalse(Match(err4, FlagUnhandled))

			ShouldNotError(
				SetFlagsByType(err3, FlagNotFound),
			)
			ShouldNotError(
				SetFlagsByType(err4, FlagUnauthorized),
			)
			ShouldNotError(
				SetFlagsByType(err3, FlagUnhandled),
			)

			ShouldBeTrue(Match(err3, FlagNotFound))
			ShouldBeTrue(Match(err3, FlagUnauthorized))
			ShouldBeTrue(Match(err3, FlagUnhandled))
			ShouldBeTrue(Match(err4, FlagNotFound))
			ShouldBeTrue(Match(err4, FlagUnauthorized))
			ShouldBeTrue(Match(err4, FlagUnhandled))

			ShouldNotError(
				DeregisterType(err4),
			)
			ShouldBeFalse(Match(err3, FlagNotFound))
			ShouldBeFalse(Match(err3, FlagUnauthorized))
			ShouldBeFalse(Match(err3, FlagUnhandled))
			ShouldBeFalse(Match(err4, FlagNotFound))
			ShouldBeFalse(Match(err4, FlagUnauthorized))
			ShouldBeFalse(Match(err4, FlagUnhandled))
		})
		Convey("Testing error which implements IErrorWithFlags", t, func() {
			err5 := &someErrWithCode{msg: "error5", code: FlagNotFound}
			ShouldBeTrue(Match(err5, FlagNotFound))
			ShouldBeFalse(Match(err5, FlagUnauthorized))

			SetFlagsByInstance(err5, FlagUnauthorized)
			ShouldBeTrue(Match(err5, FlagNotFound))
			ShouldBeTrue(Match(err5, FlagUnauthorized))

			UnsetFlagsByInstance(err5, FlagUnauthorized)
			ShouldBeTrue(Match(err5, FlagNotFound))
			ShouldBeFalse(Match(err5, FlagUnauthorized))

			UnsetFlagsByInstance(err5, FlagNotFound)
			ShouldBeTrue(Match(err5, FlagNotFound))
			ShouldBeFalse(Match(err5, FlagUnauthorized))
		})
	})
}
