//go:build !mutagencompose

package main

func init() {
	panic("executable built with without correct tag")
}
