package permissions

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerbToString(t *testing.T) {
	vs := Verbs(GET, DELETE)
	assert.Equal(t, "GET,DELETE", vs.String())
}

func TestRuleToJSON(t *testing.T) {
	r := Rule{
		Type:  "io.cozy.contacts",
		Verbs: Verbs(GET, POST),
	}

	b, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"io.cozy.contacts","verbs":["GET","POST"]}`, string(b))
}

func TestSetToJSON(t *testing.T) {
	s := Set{
		Rule{
			Title:       "contacts",
			Description: "Required for autocompletion on @name",
			Type:        "io.cozy.contacts",
			Verbs:       Verbs(GET),
		},
		Rule{
			Title:       "images",
			Description: "Required for the background",
			Type:        "io.cozy.files",
			Verbs:       Verbs(GET),
			Values:      []string{"io.cozy.files.music-dir"},
		},
		Rule{
			Title:       "mail",
			Description: "Required to send a congratulations email to your friends",
			Type:        "io.cozy.jobs",
			Selector:    "worker",
			Values:      []string{"sendmail"},
		},
	}

	b, err := json.Marshal(s)
	assert.NoError(t, err)
	assertEqualJSON(t, b, `{
    "contacts": {
      "type": "io.cozy.contacts",
      "description": "Required for autocompletion on @name",
      "verbs": ["GET"]
    },
    "images": {
      "type": "io.cozy.files",
      "description": "Required for the background",
      "verbs": ["GET"],
      "values": ["io.cozy.files.music-dir"]
    },
    "mail": {
      "type": "io.cozy.jobs",
      "description": "Required to send a congratulations email to your friends",
      "selector": "worker",
      "values": ["sendmail"]
    }
  }`)

}

func TestJSON2Set(t *testing.T) {
	jsonSet := []byte(`{
    "contacts": {
      "type": "io.cozy.contacts",
      "description": "Required for autocompletion on @name",
      "verbs": ["GET"]
    },
    "images": {
      "type": "io.cozy.files",
      "description": "Required for the background",
      "verbs": ["GET"],
      "values": ["io.cozy.files.music-dir"]
    },
    "mail": {
      "type": "io.cozy.jobs",
      "description": "Required to send a congratulations email to your friends",
      "selector": "worker",
      "values": ["sendmail"]
    }
  }`)
	var s Set
	err := json.Unmarshal(jsonSet, &s)
	assert.NoError(t, err)
}

func TestSetToString(t *testing.T) {
	s := Set{
		Rule{
			Title:       "contacts",
			Description: "Required for autocompletion on @name",
			Type:        "io.cozy.contacts",
		},
		Rule{
			Title:       "images",
			Description: "Required for the background",
			Type:        "io.cozy.files",
			Verbs:       Verbs(GET),
			Values:      []string{"io.cozy.files.music-dir"},
		},
	}

	out, err := s.MarshalScopeString()
	assert.NoError(t, err)
	assert.Equal(t, out, "io.cozy.contacts io.cozy.files:GET:io.cozy.files.music-dir")
}

func TestStringToSet(t *testing.T) {

	s, err := UnmarshalScopeString("io.cozy.contacts io.cozy.files:GET:io.cozy.files.music-dir")

	assert.NoError(t, err)
	assert.Len(t, s, 2)
	assert.Equal(t, "io.cozy.contacts", s[0].Type)
	assert.Equal(t, "io.cozy.files", s[1].Type)
	assert.Len(t, s[1].Verbs, 1)
	assert.Equal(t, GET, s[1].Verbs[0])
	assert.Len(t, s[1].Values, 1)
	assert.Equal(t, "io.cozy.files.music-dir", s[1].Values[0])

}

func assertEqualJSON(t *testing.T, value []byte, expected string) {
	expectedBytes := new(bytes.Buffer)
	err := json.Compact(expectedBytes, []byte(expected))
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes.String(), string(value))
}
