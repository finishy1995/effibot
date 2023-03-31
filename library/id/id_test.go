package id

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIDMarshalJSON(t *testing.T) {
	id := ID(10018820)
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(id)
	got := `"000000000098e004"`
	want := strings.TrimSpace(buf.String())

	assert.Equal(t, got, want)
}

func TestIDUnmarshalJSONHexString(t *testing.T) {
	j := []byte(`"000000000098e004"`)
	var got ID
	if err := json.Unmarshal(j, &got); err != nil {
		t.Fatal(err)
	}
	want := ID(10018820)

	assert.Equal(t, got, want)
}

func TestIDUnmarshalJSONInt(t *testing.T) {
	j := []byte(`10018820`)
	var got ID
	if err := json.Unmarshal(j, &got); err != nil {
		t.Fatal(err)
	}
	want := ID(10018820)

	assert.Equal(t, got, want)
}

func TestIDUnmarshalJSONNonInt(t *testing.T) {
	j := []byte(`[]`)
	var got ID
	err := json.Unmarshal(j, &got)

	assert.NotNil(t, err)
}

func TestIDUnmarshalJSONNonHexString(t *testing.T) {
	j := []byte(`"woo"`)
	var got ID
	err := json.Unmarshal(j, &got)

	assert.NotNil(t, err)
}

func TestIDGeneration(t *testing.T) {
	n := 10000
	ids := make(map[ID]bool, n)
	for i := 0; i < n; i++ {
		id := GenerateID()
		if ids[id] {
			assert.FailNow(t, "duplicate ID", id)
		}
		ids[id] = true
	}
}

func TestParseID(t *testing.T) {
	want := ID(10018181901)
	got, err := ParseID(want.String())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, got, want)
}

func TestParseIDError(t *testing.T) {
	_, err := ParseID("woo")
	assert.NotNil(t, err)
}

func BenchmarkIDGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateID()
	}
}
