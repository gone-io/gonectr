package priest

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func TimeStat(processName string) func() {
	if isStat {
		beginTime := time.Now()
		return func() {
			log.Infof("stat <%s> process use time:%v\n", processName, time.Since(beginTime))
		}
	}
	return func() {}
}
