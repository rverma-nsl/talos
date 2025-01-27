// Code generated by "enumer -type=LinkFlag -linecomment -text"; DO NOT EDIT.

//
package nethelpers

import (
	"fmt"
)

const _LinkFlagName = "UPBROADCASTDEBUGLOOPBACKPOINTTOPOINTNOTRAILERSRUNNINGNOARPPROMISCALLMULTIMASTERSLAVEMULTICASTPORTSELAUTOMEDIADYNAMICLOWER_UPDORMANTECHO"

var _LinkFlagMap = map[LinkFlag]string{
	1:      _LinkFlagName[0:2],
	2:      _LinkFlagName[2:11],
	4:      _LinkFlagName[11:16],
	8:      _LinkFlagName[16:24],
	16:     _LinkFlagName[24:36],
	32:     _LinkFlagName[36:46],
	64:     _LinkFlagName[46:53],
	128:    _LinkFlagName[53:58],
	256:    _LinkFlagName[58:65],
	512:    _LinkFlagName[65:73],
	1024:   _LinkFlagName[73:79],
	2048:   _LinkFlagName[79:84],
	4096:   _LinkFlagName[84:93],
	8192:   _LinkFlagName[93:100],
	16384:  _LinkFlagName[100:109],
	32768:  _LinkFlagName[109:116],
	65536:  _LinkFlagName[116:124],
	131072: _LinkFlagName[124:131],
	262144: _LinkFlagName[131:135],
}

func (i LinkFlag) String() string {
	if str, ok := _LinkFlagMap[i]; ok {
		return str
	}
	return fmt.Sprintf("LinkFlag(%d)", i)
}

var _LinkFlagValues = []LinkFlag{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144}

var _LinkFlagNameToValueMap = map[string]LinkFlag{
	_LinkFlagName[0:2]:     1,
	_LinkFlagName[2:11]:    2,
	_LinkFlagName[11:16]:   4,
	_LinkFlagName[16:24]:   8,
	_LinkFlagName[24:36]:   16,
	_LinkFlagName[36:46]:   32,
	_LinkFlagName[46:53]:   64,
	_LinkFlagName[53:58]:   128,
	_LinkFlagName[58:65]:   256,
	_LinkFlagName[65:73]:   512,
	_LinkFlagName[73:79]:   1024,
	_LinkFlagName[79:84]:   2048,
	_LinkFlagName[84:93]:   4096,
	_LinkFlagName[93:100]:  8192,
	_LinkFlagName[100:109]: 16384,
	_LinkFlagName[109:116]: 32768,
	_LinkFlagName[116:124]: 65536,
	_LinkFlagName[124:131]: 131072,
	_LinkFlagName[131:135]: 262144,
}

// LinkFlagString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func LinkFlagString(s string) (LinkFlag, error) {
	if val, ok := _LinkFlagNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to LinkFlag values", s)
}

// LinkFlagValues returns all values of the enum
func LinkFlagValues() []LinkFlag {
	return _LinkFlagValues
}

// IsALinkFlag returns "true" if the value is listed in the enum definition. "false" otherwise
func (i LinkFlag) IsALinkFlag() bool {
	_, ok := _LinkFlagMap[i]
	return ok
}

// MarshalText implements the encoding.TextMarshaler interface for LinkFlag
func (i LinkFlag) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for LinkFlag
func (i *LinkFlag) UnmarshalText(text []byte) error {
	var err error
	*i, err = LinkFlagString(string(text))
	return err
}
