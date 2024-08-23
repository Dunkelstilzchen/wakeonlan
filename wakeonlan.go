package wakeonlan

import (
	"encoding/hex"
	"net"
)

// Sends the WakeonLan magic package for the given MacAddress.
//
// Returns nil or error
func Wake(macAddress, networkBroadcastAddress string) error {
	magicPaket, err := buildMagicPackage(macAddress)
	if err != nil {
		return err
	}
	return send(networkBroadcastAddress, magicPaket)
}

// buildMagicPackage generates the wakeonlan magic packet for the given MacAddress.
//
//	The address needs to be of one the following formats:
//
//	00:00:00:00:00:00
//	00-00-00-00-00-00
//
// Returns: ([]byte, nil) or (nil, error)
func buildMagicPackage(macAddress string) ([]byte, error) {
	macBytes := make([]byte, 0)
	for i := 0; i < len(macAddress); i += 3 {
		val, err := hex.DecodeString(macAddress[i : i+2])
		if err != nil {
			return nil, err
		}
		macBytes = append(macBytes, val[0])
	}

	magicPackage := make([]byte, 6)

	// first 6 bytes all 0xff
	for i := 0; i < 6; i++ {
		magicPackage[i] = 255
	}

	// now 16 times macaddress bytes
	for i := 0; i < 16; i++ {
		for j := 0; j < len(macBytes); j++ {
			magicPackage = append(magicPackage, macBytes[j])
		}
	}

	return magicPackage, nil
}

// Sends the data array to the given address using udp4.
//
// Returns an error if unsuccessfull
func send(networkAddress string, data []byte) error {
	pc, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		return err
	}
	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp4", networkAddress)
	if err != nil {
		return err
	}

	_, err = pc.WriteTo(data, addr)
	if err != nil {
		return err
	}

	return nil
}
