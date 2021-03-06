package blockchain

import (
	"errors"
	"fmt"
	"net"

	"Blockchain-Framework/src/core"
)

//StdBCServer : create a standard server
func StdBCServer(sc core.ServerConfig) core.Server {

	s := core.NewServer(sc)
	StdAddListeners(s)
	return s

}
func stdServerConfig() core.ServerConfig {
	sc := core.NewServerConfig(":8080", 1)
	return sc
}

var defListen map[string]chan net.Conn

//StdAddListeners :  add the standard listeners to the server
//NEW CONNECTION
//NEW DOSCONNECTION
func StdAddListeners(s core.Server) {
	defListen = make(map[string]chan net.Conn)
	s.StatusIn("adding listener")
	defListen["connected"] = s.AddListener("connected", stdOnConnected)
	//defListen["bcmex"] = s.AddListener("bcmex", HandleBCmessage)

}
func stdOnConnected(s core.Server, conn net.Conn) {

	fmt.Println()
	for {
		if err := StdReciveMessage(s, conn); err != nil {
			fmt.Println(err)
			break
		}

	}
}
func discF(s string) {
	//TODO : implement disconnection
	fmt.Println("disconnected:", s)
}

//FIXME : implement messages for future users

//StdReciveMessage : Standard function to recive messages
func StdReciveMessage(s core.Server, conn net.Conn) error {
	mType := make([]byte, 1)

	_, err := conn.Read(mType)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("new mess", mType, "from ", conn.RemoteAddr().String())
	switch (core.MessageType)(mType[0]) {
	case core.VoidMessage:
		{
		}
	case core.BChainMessage:
		{
			fmt.Println("recived bc message")
			HandleBCmessage(s, conn)

		}
	case core.NameString:
		{

		}
	case core.NewInRoom:
		{
		}
	case core.ChatAll:
		{
		}
	case core.ChatRoom:
		{
			//TODO : implement ChatRoom
		}
	case core.ChatTo:
		{
			//TODO : implement ChatTo
		}
	default:
		{
			return errors.New("Recived unknown message type")
		}
	}
	return nil
}

//HandleBCmessage : recive bc messages and answer them
func HandleBCmessage(s core.Server, conn net.Conn) {
	mType := make([]byte, 1)
	_, err := conn.Read(mType)
	fmt.Println("recived:", mType[0], core.FmtBcMex((core.MessageType)(mType[0])))
	if err != nil {
		fmt.Println(err)
	}
	switch (core.MessageType)(mType[0]) {
	case core.AmReady:
		{
			fmt.Println(conn.RemoteAddr().String(), "is ready")
			s.AssignRoom("ready", conn.RemoteAddr().String())
		}
	case core.Config:
		{

			lenb := make([]byte, 1)
			_, err := conn.Read(lenb)
			if err != nil {
				fmt.Println(err)
			}
			hash := make([]byte, lenb[0])
			_, err = conn.Read(hash)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("recived hash:", string(hash))
			myhash, err := s.GetVar("ConfigHash")
			if err != nil {
				fmt.Println(err)
				return
			}
			verified := true
			myh := ([]byte)(myhash.(string))
			if len(myh) != len(hash) {
				verified = false
			} else {
				for i, v := range hash {
					if v != myh[i] {
						verified = false
					}
				}
			}
			fmt.Println("hash verif:", verified)
			fmt.Println("sending confirm config")
			confMess := core.NewMessage(core.ConfirmConfig, myh)
			bcm := core.NewMessage(core.BChainMessage, confMess)
			bcm.Send(nil, conn)
			s.AddConnection(conn.RemoteAddr().String(), conn)
			err = s.AssignRoom("nodes", conn.RemoteAddr().String())

			if err != nil {
				fmt.Println(err)
			}
		}
	case core.ConfirmConfig:
		{
			lenb := make([]byte, 1)
			_, err := conn.Read(lenb)
			if err != nil {
				fmt.Println(err)
			}
			hash := make([]byte, lenb[0])
			_, err = conn.Read(hash)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("recived hash:", string(hash))
			myhash, err := s.GetVar("ConfigHash")
			if err != nil {
				fmt.Println(err)
				return
			}
			verified := true
			myh := ([]byte)(myhash.(string))
			if len(myh) != len(hash) {
				verified = false
			} else {
				for i, v := range hash {
					if v != myh[i] {
						verified = false
					}
				}
			}
			fmt.Println("hash verif:", verified)
			s.AddConnection(conn.RemoteAddr().String(), conn)
			err = s.AssignRoom("nodes", conn.RemoteAddr().String())

			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
