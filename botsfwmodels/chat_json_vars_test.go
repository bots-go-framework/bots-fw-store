package botsfwmodels

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetJSONVar_GetJSONVar_int(t *testing.T) {
	chat := &ChatBaseData{}

	require.NoError(t, SetJSONVar(chat, "count", 42))

	var got int
	require.NoError(t, GetJSONVar(chat, "count", &got))
	assert.Equal(t, 42, got)
}

func TestSetJSONVar_GetJSONVar_stringSlice(t *testing.T) {
	chat := &ChatBaseData{}
	in := []string{"alpha", "beta", "gamma"}

	require.NoError(t, SetJSONVar(chat, "tags", in))

	var got []string
	require.NoError(t, GetJSONVar(chat, "tags", &got))
	assert.Equal(t, in, got)
}

func TestGetJSONVar_absentKey_returnsNil(t *testing.T) {
	chat := &ChatBaseData{}
	var got int
	err := GetJSONVar(chat, "missing", &got)
	assert.NoError(t, err)
	assert.Equal(t, 0, got) // untouched
}

func TestGetJSONVar_invalidJSON_returnsError(t *testing.T) {
	chat := &ChatBaseData{}
	chat.SetVar("bad", "not-json")

	var got int
	err := GetJSONVar(chat, "bad", &got)
	assert.Error(t, err)
}

func TestSetJSONVar_GetJSONVar_noFloat64Mangling(t *testing.T) {
	// Ensure that storing an int and retrieving it as int does NOT produce
	// the float64 mangling that hand-rolled JSON round-trips via interface{} cause.
	chat := &ChatBaseData{}
	require.NoError(t, SetJSONVar(chat, "n", 99))

	var n int
	require.NoError(t, GetJSONVar(chat, "n", &n))
	assert.Equal(t, 99, n)
}

func TestSetJSONVar_GetJSONVar_struct(t *testing.T) {
	type pendingAction struct {
		Verb    string
		Subject string
	}
	chat := &ChatBaseData{}
	in := pendingAction{Verb: "confirm", Subject: "spot-123"}

	require.NoError(t, SetJSONVar(chat, "pending", in))

	var got pendingAction
	require.NoError(t, GetJSONVar(chat, "pending", &got))
	assert.Equal(t, in, got)
}
