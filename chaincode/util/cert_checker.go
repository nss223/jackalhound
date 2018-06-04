//
// cert_checker.go
// Copyright (C) 2018 jack <jack@HP-WorkStation>
//
// Distributed under terms of the MIT license.
//

package util

import (
	_ "log"
)

func IsAdmin(cert string) bool {
	if cert == "admin" {
		return true
	} else if len(cert) < 6 {
		return false
	} else {
		return cert[:6] == "Admin@"
	}
}
