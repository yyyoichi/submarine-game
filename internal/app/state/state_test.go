package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	type user struct {
		name string
		age  int
	}
	t.Run("Test UseState", func(t *testing.T) {
		u, setUser := UseState[user]()
		assert.Zero(t, *u)
		setUser(user{"hoge", 20})
		assert.Equal(t, "hoge", u.name)
		assert.Equal(t, 20, u.age)
	})

	t.Run("Test UseGlobalState", func(t *testing.T) {
		u, setUser := UseGlobalState[user]("test", WithInitialValue(user{"taro", 20}))
		assert.Equal(t, 20, u.age)
		done := make(chan struct{})
		go func() {
			defer close(done)
			u, setUser := UseGlobalState[user]("test", WithInitialValue(user{"taro", 100}))
			assert.Equal(t, 20, u.age)
			setUser(user{"taro", 21})
		}()
		<-done
		assert.Equal(t, 21, u.age)
		setUser(user{"taro", 22})
		assert.Equal(t, 22, u.age)

		// retry
		u2, setUser2 := UseGlobalState[user]("test")
		assert.Equal(t, 22, u2.age)
		setUser2(user{"taro", 23})
		assert.Equal(t, 23, u2.age)
		assert.Equal(t, 23, u.age)
		assert.Equal(t, u, u2)
	})

	t.Run("Test UseGlobalState with validate", func(t *testing.T) {
		u, setUser := UseGlobalState[user](
			"test-validate",
			WithSetValidateFunc(func(u user) error {
				if u.age < 20 {
					return assert.AnError
				}
				return nil
			}),
			WithInitialValue(user{"taro", 20}),
		)
		assert.Equal(t, 20, u.age)
		err := setUser(user{"taro", 19})
		assert.ErrorIs(t, err, assert.AnError)

		u, setUser = UseGlobalState[user](
			"test-validate",
			WithSetValidateFunc(func(u user) error {
				return assert.AnError
			}),
			WithInitialValue(user{"taro", 100}),
		)
		// set only once at first time
		assert.NoError(t, setUser(user{"taro", 21}))
		assert.Equal(t, 21, u.age)
	})
}
