package testing

import (
	"github.com/celer-network/sgn/testing/log"
)

func ChkErr(e error, msg string) {
	if e != nil {
		log.Fatalln(msg, e)
	}
}
