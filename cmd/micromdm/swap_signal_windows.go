// +build test windows

package main

import "os"

// TODO(issues/686) implement signals for windows
var defaultSwapSignal = os.Signal(nil)
