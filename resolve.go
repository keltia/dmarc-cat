package main

import (
	"net"
)

type Resolver interface {
	LookupAddr(addr string) ([]string, error)
}

type NullResolver struct{}

func (NullResolver) LookupAddr(addr string) ([]string, error) {
	return []string{addr}, nil
}

type RealResolver struct{}

func (r RealResolver) LookupAddr(addr string) ([]string, error) {
	return net.LookupAddr(addr)
}
