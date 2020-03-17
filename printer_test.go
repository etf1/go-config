package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/etf1/config"
)

type object struct {
	myPrivateParam   string
	MyObfuscateParam bool `print:"-"`

	MyBoolParam     bool
	MyIntParam      int
	MyStringParam   string
	MyDurationParam time.Duration

	SubObject subObject
}
type subObject struct {
	myPrivateParam   string
	MyObfuscateParam bool `print:"-"`

	MyBoolParam     bool
	MyIntParam      int
	MyStringParam   string
	MyDurationParam time.Duration
}

//-----------------------------------
//MyObfuscateParam| *** Hidden value ***|bool `print:"-"`
//     MyBoolParam|                 true|bool ``
//      MyIntParam|                  111|int ``
//   MyStringParam|         string value|string ``
// MyDurationParam|                222ns|time.Duration ``
//------ SubObject
//MyObfuscateParam| *** Hidden value ***|bool `print:"-"`
//     MyBoolParam|                false|bool ``
//      MyIntParam|                  333|int ``
//   MyStringParam|     SUB string value|string ``
// MyDurationParam|                444ns|time.Duration ``
//-----------------------------------
func TestTableString(t *testing.T) {
	obj := object{
		myPrivateParam:   "private_value",
		MyObfuscateParam: false,
		MyBoolParam:      true,
		MyIntParam:       111,
		MyStringParam:    "string value",
		MyDurationParam:  222,
		SubObject: subObject{
			myPrivateParam:  "SUB_private_value",
			MyIntParam:      333,
			MyStringParam:   "SUB string value",
			MyDurationParam: 444,
		},
	}

	assert.Equal(t,
		"\n-----------------------------------\n MyObfuscateParam| \x1b[0m*** Hidden value ***|\x1b[1;34mbool\x1b[0m \x1b[1;92m`print:\"-\"`\x1b[0m\n      MyBoolParam|                 \x1b[0mtrue|\x1b[1;34mbool\x1b[0m \x1b[1;92m``\x1b[0m\n       MyIntParam|                  \x1b[0m111|\x1b[1;34mint\x1b[0m \x1b[1;92m``\x1b[0m\n    MyStringParam|         \x1b[0mstring value|\x1b[1;34mstring\x1b[0m \x1b[1;92m``\x1b[0m\n  MyDurationParam|                \x1b[0m222ns|\x1b[1;34mtime.Duration\x1b[0m \x1b[1;92m``\x1b[0m\n##### SubObject #####\n MyObfuscateParam| \x1b[0m*** Hidden value ***|\x1b[1;34mbool\x1b[0m \x1b[1;92m`print:\"-\"`\x1b[0m\n      MyBoolParam|                \x1b[0mfalse|\x1b[1;34mbool\x1b[0m \x1b[1;92m``\x1b[0m\n       MyIntParam|                  \x1b[0m333|\x1b[1;34mint\x1b[0m \x1b[1;92m``\x1b[0m\n    MyStringParam|     \x1b[0mSUB string value|\x1b[1;34mstring\x1b[0m \x1b[1;92m``\x1b[0m\n  MyDurationParam|                \x1b[0m444ns|\x1b[1;34mtime.Duration\x1b[0m \x1b[1;92m``\x1b[0m\n-----------------------------------\n",
		config.TableString(obj),
	)
}
