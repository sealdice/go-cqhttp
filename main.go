// Package main
package main

import (
	"context"
	"net"

	"github.com/Mrs4s/go-cqhttp/cmd/gocq"
	"github.com/Mrs4s/go-cqhttp/global/terminal"

	_ "unsafe"

	_ "github.com/Mrs4s/go-cqhttp/db/leveldb"   // leveldb 数据库支持
	_ "github.com/Mrs4s/go-cqhttp/modules/silk" // silk编码模块
	// 其他模块
	// _ "github.com/Mrs4s/go-cqhttp/db/sqlite3"   // sqlite3 数据库支持
	// _ "github.com/Mrs4s/go-cqhttp/db/mongodb"    // mongodb 数据库支持
	// _ "github.com/Mrs4s/go-cqhttp/modules/pprof" // pprof 性能分析
)

func main() {
	terminal.SetTitle()
	dnsHack()
	gocq.InitBase()
	gocq.PrepareData()
	gocq.LoginInteract()
	_ = terminal.DisableQuickEdit()
	_ = terminal.EnableVT100()
	gocq.WaitSignal()
	_ = terminal.RestoreInputMode()
}

func dnsHack() {
	var (
		dnsResolverIP    = "114.114.114.114:53" // Google DNS resolver.
		dnsResolverProto = "udp"                // Protocol to use for the DNS resolver
	)
	var dialer net.Dialer
	net.DefaultResolver = &net.Resolver{
		PreferGo: false,
		Dial: func(context context.Context, _, _ string) (net.Conn, error) {
			conn, err := dialer.DialContext(context, dnsResolverProto, dnsResolverIP)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}
