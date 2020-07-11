// +build !windows test

package main

import "syscall"

var defaultSwapSignal = syscall.SIGUSR2
