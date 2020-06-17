// Copyright (c) 2020 Shivaram Lingamneni
// released under the Apache 2.0 license

package ldap

import (
	"fmt"
	"os"
	"strings"
)

// logging stub that just dumps everything to stderr

type Logger struct{}

func (l *Logger) Debug(args ...string) {
	fmt.Fprintf(os.Stderr, "%s\n", strings.Join(args, " : "))
}

var log Logger
