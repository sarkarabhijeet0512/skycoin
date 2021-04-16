package main

import (
	"log"
	"net"
	"sky_coin_messanger/server"
)

func main() {
	// message := 2048
	// priv, pub := cipher.GenerateKeyPair(message)
	// privbytes := cipher.PrivateKeyToBytes(priv)
	// pubbytes := cipher.PublicKeyToBytes(pub)
	// pubKey := cipher.BytesToPublicKey(pubbytes)
	// PrivKey := cipher.BytesToPrivateKey(privbytes)
	s := server.NewServer()
	go s.Run()
	listner, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("unable to start the server %s", err.Error())
	}
	defer listner.Close()
	log.Println("starting server on port :8000")
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Println("unable to Accept server conncetion", err.Error())
			continue
		}
		go s.NewClient(conn)
	}
}
