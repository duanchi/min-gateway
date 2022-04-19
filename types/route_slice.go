package types

import "github.com/duanchi/min-gateway/mapper"

type RouteSlice []mapper.Route

func (s RouteSlice) Len() int      { return len(s) }
func (s RouteSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s RouteSlice) Less(i, j int) bool { return s[i].Sort < s[j].Sort }
