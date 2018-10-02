package config

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockFile struct {
	mock.Mock
}

func (f *mockFile) GetKey(name string) string {
	args := f.Called(name)
	return args.String(0)
}
func (f *mockFile) RemoveKey(name string) {
	f.Called(name)
	return
}
func (f *mockFile) SetKey(name, value string) {
	f.Called(name, value)
	return
}
func (f *mockFile) Save() error {
	args := f.Called()
	return args.Error(0)
}

func TestGetPinned(t *testing.T) {
	require := require.New(t)

	f := new(mockFile)
	defer f.AssertExpectations(t)
	f.On("GetKey", "core.version").Return("0")

	c := Config{Project: f}
	actual := c.GetPinned()
	require.Equal("0", actual)

}

func TestSetPinned(t *testing.T) {
	for _, tc := range []struct {
		Current string
		Request string
		R1      string
		R2      string
	}{
		{"1", "1", "1", "no change"},
		{"2", "3", "2", "3"},
		{"4", "", "4", "not pinned"},
		{"", "3", "not pinned", "3"},
	} {
		tc := tc
		name := "[" + tc.Current + ":" + tc.Request + "]"
		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			f := new(mockFile)
			defer f.AssertExpectations(t)
			f.On("GetKey", "core.version").Return(tc.Current)
			switch {
			case tc.Current == tc.Request:
				// noop
			case tc.Request == "":
				f.On("RemoveKey", "core.version").Return()
				f.On("Save").Return(nil)
			default:
				f.On("SetKey", "core.version", tc.Request).Return()
				f.On("Save").Return(nil)
			}
			c := Config{Project: f}

			r1, r2 := c.SetPinned(tc.Request)
			require.Equal(tc.R1, r1)
			require.Equal(tc.R2, r2)

			err := c.Save()
			require.NoError(err)
		})
	}
}
