package tcp_proxy

import (
	"net"
	"strings"

	"github.com/bettercap/bettercap/log"
	"github.com/bettercap/bettercap/session"

	"github.com/robertkrimen/otto"

	"github.com/evilsocket/islazy/plugin"
)

type TcpProxyScript struct {
	*plugin.Plugin
	doOnData bool
}

func LoadTcpProxyScript(path string, sess *session.Session) (err error, s *TcpProxyScript) {
	log.Info("loading tcp proxy script %s ...", path)

	plug, err := plugin.Load(path)
	if err != nil {
		return
	}

	// define session pointer
	if err = plug.Set("env", sess.Env.Data); err != nil {
		log.Error("Error while defining environment: %+v", err)
		return
	}

	// run onLoad if defined
	if plug.HasFunc("onLoad") {
		if _, err = plug.Call("onLoad"); err != nil {
			log.Error("Error while executing onLoad callback: %s", "\nTraceback:\n  "+err.(*otto.Error).String())
			return
		}
	}

	s = &TcpProxyScript{
		Plugin:   plug,
		doOnData: plug.HasFunc("onData"),
	}
	return
}

func (s *TcpProxyScript) OnData(from, to net.Addr, data []byte) []byte {
	if s.doOnData {
		addrFrom := strings.Split(from.String(), ":")[0]
		addrTo := strings.Split(to.String(), ":")[0]

		if ret, err := s.Call("onData", addrFrom, addrTo, data); err != nil {
			log.Error("Error while executing onData callback: %s", err)
			return nil
		} else if ret != nil {
			array, ok := ret.([]byte)
			if !ok {
				log.Error("Error while casting exported value to array of byte: value = %+v", ret)
			}
			return array
		}
	}
	return nil
}
