package recipes

import (
	"testing"

	"github.com/curator-go/curator"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLockInternalsDriver(t *testing.T) {
	Convey("Given a StandardLockInternalsDriver", t, func() {
		driver := NewStandardLockInternalsDriver()

		Convey("Should implements LockInternalsSorter", func() {
			So(driver, ShouldImplement, (*LockInternalsSorter)(nil))

			Convey("When fix a string not contains lock name", func() {
				s := driver.FixForSorting("some", "lock")

				Convey("Return orignal string", func() {
					So(s, ShouldEqual, "some")
				})
			})

			Convey("When fix a string contains lock name", func() {
				s := driver.FixForSorting("/path/lock/tail", "lock")

				Convey("Return the tail", func() {
					So(s, ShouldEqual, "/tail")
				})
			})

			Convey("When fix a string without tail", func() {
				s := driver.FixForSorting("/path/lock", "lock")

				Convey("Return a empty string", func() {
					So(s, ShouldBeEmpty)
				})
			})
		})

		Convey("Should implements LockInternalsDriver", func() {
			So(driver, ShouldImplement, (*LockInternalsDriver)(nil))

			Convey("When gets the lock", func() {

				Convey("When the lock is nonexists", func() {
					ret, err := driver.GetsTheLock(nil, []string{}, "lock", 0)

					Convey("Return NoNode error", func() {
						So(ret, ShouldBeNil)
						So(err, ShouldEqual, curator.ErrNoNode)
					})
				})

				Convey("When max leases less than lock index", func() {
					ret, err := driver.GetsTheLock(nil, []string{"1st", "2nd", "lock", "4th"}, "lock", 3)

					Convey("Get the lock", func() {
						So(ret, ShouldNotBeNil)
						So(ret.GetsTheLock, ShouldBeTrue)
						So(ret.PathToWatch, ShouldBeEmpty)
						So(err, ShouldBeNil)
					})
				})

				Convey("When max leases large than lock index", func() {
					ret, err := driver.GetsTheLock(nil, []string{"1st", "2nd", "lock", "4th"}, "lock", 2)

					Convey("Return the path to watch", func() {
						So(ret, ShouldNotBeNil)
						So(ret.GetsTheLock, ShouldBeFalse)
						So(ret.PathToWatch, ShouldEqual, "1st")
						So(err, ShouldBeNil)
					})
				})
			})

			Convey("When creates the lock", func() {
				mocks := newMockBuilder(t)

				client := mocks.Build()

				So(client.Start(), ShouldBeNil)

				Convey("When lock with data", func() {
					mocks.conn.On("Create", "/lock", []byte("data"), int32(curator.EPHEMERAL_SEQUENTIAL), curator.OPEN_ACL_UNSAFE).Return("/lock", nil).Once()

					path, err := driver.CreatesTheLock(client, "/lock", []byte("data"))

					Convey("Create lock file with data", func() {
						So(path, ShouldEqual, "/lock")
						So(err, ShouldBeNil)

						mocks.Check(t)
					})
				})

				Convey("When lock without data", func() {
					mocks.conn.On("Create", "/lock", mocks.builder.DefaultData, int32(curator.EPHEMERAL_SEQUENTIAL), curator.OPEN_ACL_UNSAFE).Return("/lock", nil).Once()

					path, err := driver.CreatesTheLock(client, "/lock", nil)

					Convey("Create lock file without data", func() {
						So(path, ShouldEqual, "/lock")
						So(err, ShouldBeNil)

						mocks.Check(t)
					})
				})
			})
		})
	})
}

func TestLockInternals(t *testing.T) {
	Convey("Given lockInternals", t, func() {
		mocks := newMockBuilder(t)

		client := mocks.Build()

		So(client.Start(), ShouldBeNil)

		Convey("base on invalidated path", func() {
			internal, err := newLockInternals(client, mocks.driver, "invalid", LockPrefix, 3)

			So(internal, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("base on a validated path", func() {
			internal, err := newLockInternals(client, mocks.driver, "/path", LockPrefix, 3)

			So(internal, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

		mocks.Check(t)
	})
}

func TestInterProcessMutex(t *testing.T) {
	Convey("Given an InterProcessMutex base on a path", t, func() {

	})
}
