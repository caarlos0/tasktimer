package cmd

import (
	"log"
)

type badgerStdLoggerAdapter struct{}

func (b badgerStdLoggerAdapter) Errorf(s string, i ...interface{}) {
	log.Printf("[ERR] "+s, i...)
}

func (b badgerStdLoggerAdapter) Warningf(s string, i ...interface{}) {
	log.Printf("[WARN] "+s, i...)
}

func (b badgerStdLoggerAdapter) Infof(s string, i ...interface{}) {
	log.Printf("[INFO] "+s, i...)
}

func (b badgerStdLoggerAdapter) Debugf(s string, i ...interface{}) {
	log.Printf("[DEBUG] "+s, i...)
}
