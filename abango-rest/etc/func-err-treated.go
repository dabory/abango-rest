// Author : Eric Kim
// Build Date : 6 Jul 2023  Last Update 02 Aug 2018
// All rights are reserved.

package etc

import (
	"encoding/json"
)

func DbrUnmarshal(bytes []byte, target interface{}) error {
	if err := json.Unmarshal(bytes, &target); err == nil {
		return nil
	} else {
		return LogErr("LOOUJGYTG", "Format mismatch:\n"+string(bytes), err)
	}
}

func DbrMarshal(target interface{}) []byte {
	if bytes, err := json.Marshal(&target); err == nil {
		return bytes
	} else {
		LogErr("LOOUJGYTG", "Format mismatch:\n"+string(bytes), err)
		return nil
	}
}
