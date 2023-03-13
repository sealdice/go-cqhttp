// Package main
package main

import (
	"github.com/Mrs4s/go-cqhttp/cmd/gocq"
	"github.com/Mrs4s/go-cqhttp/global/terminal"
	"sync"
	"time"

	_ "unsafe"

	_ "github.com/Mrs4s/go-cqhttp/db/leveldb"   // leveldb 数据库支持
	_ "github.com/Mrs4s/go-cqhttp/modules/silk" // silk编码模块
	// 其他模块
	// _ "github.com/Mrs4s/go-cqhttp/db/sqlite3"   // sqlite3 数据库支持
	// _ "github.com/Mrs4s/go-cqhttp/db/mongodb"    // mongodb 数据库支持
	// _ "github.com/Mrs4s/go-cqhttp/modules/pprof" // pprof 性能分析
)

func main() {
	SetDefaultNS([]string{"114.114.114.114:53", "8.8.8.8:53"}, false)
	terminal.SetTitle()
	gocq.InitBase()
	gocq.PrepareData()
	gocq.LoginInteract()
	_ = terminal.DisableQuickEdit()
	_ = terminal.EnableVT100()
	gocq.WaitSignal()
	_ = terminal.RestoreInputMode()
}

//go:linkname defaultNS net.defaultNS
var defaultNS []string

// need to keep sync with go version
//
//go:linkname resolvConf net.resolvConf
var resolvConf resolverConfig

// copy from /src/net/dnsconfig_unix.go
type dnsConfig struct {
	servers    []string      // server addresses (in host:port form) to use
	search     []string      // rooted suffixes to append to local name
	ndots      int           // number of dots in name to trigger absolute lookup
	timeout    time.Duration // wait before giving up on a query, including retries
	attempts   int           // lost packets before giving up on server
	rotate     bool          // round robin among servers
	unknownOpt bool          // anything unknown was encountered
	lookup     []string      // OpenBSD top-level database "lookup" order
	err        error         // any error that occurs during open of resolv.conf
	mtime      time.Time     // time of resolv.conf modification
	soffset    uint32        // used by serverOffset
}

// copy from /src/net/dnsclient_unix.go
type resolverConfig struct {
	initOnce sync.Once // guards init of resolverConfig

	// ch is used as a semaphore that only allows one lookup at a
	// time to recheck resolv.conf.
	ch          chan struct{} // guards lastChecked and modTime
	lastChecked time.Time     // last time resolv.conf was checked

	mu        sync.RWMutex // protects dnsConfig
	dnsConfig *dnsConfig   // parsed resolv.conf structure used in lookups
}

//go:linkname (*resolverConfig).tryUpdate net.(*resolverConfig).tryUpdate
func (conf *resolverConfig) tryUpdate(name string)

// need to put a empty .s file

func SetDefaultNS(addrs []string, loadFromSystem bool) {
	if resolvConf.dnsConfig == nil {
		resolvConf.tryUpdate("")
	}

	if loadFromSystem {
		now := time.Now()
		resolvConf.lastChecked = now.Add(-6 * time.Second)
		resolvConf.dnsConfig.mtime = now
	}

	resolvConf.dnsConfig.servers = addrs
	defaultNS = addrs
}

// --------------
