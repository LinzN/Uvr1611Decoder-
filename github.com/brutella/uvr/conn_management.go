/*
 * Copyright (C) 2020. Niklas Linz - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the LGPLv3 license, which unfortunately won't be
 * written for another century.
 *
 * You should have received a copy of the LGPLv3 license with
 * this file. If not, please write to: niklas.linz@enigmar.de
 *
 */

package uvr

import (
	"fmt"
	"github.com/brutella/can"
	"github.com/brutella/canopen"
	"time"
)

func Disconnect(serverID uint8, clientID uint8, bus *can.Bus) error {
	b := []byte{
		0x80 + byte(serverID),
		0x01, 0x1F,
		0x00,
		byte(serverID),
		byte(clientID),
		0x80,
		0x12,
	}

	return sendConnManagementData(b, serverID, clientID, bus)
}

func Connect(serverID uint8, clientID uint8, bus *can.Bus) error {
	b := []byte{
		0x80 + byte(serverID),
		0x00, 0x1F,
		0x00,
		byte(serverID),
		byte(clientID),
		0x80,
		0x12,
	}

	return sendConnManagementData(b, serverID, clientID, bus)
}

func sendConnManagementData(b []byte, serverID uint8, clientID uint8, bus *can.Bus) error {
	c := &canopen.Client{bus, time.Second * 2}
	frm := canopen.Frame{
		CobID: uint16(MPDOClientServerConnManagement) + uint16(clientID),
		Data:  b,
	}

	respCobID := uint32(MPDOClientServerConnManagement) + uint32(serverID)
	req := canopen.NewRequest(frm, respCobID)
	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	frm = resp.Frame

	if b0 := frm.Data[0]; b0 != 0x80+byte(clientID) {
		return fmt.Errorf("Invalid MPDO address %v\n", b0)
	}

	if b4, b5 := frm.Data[4], frm.Data[5]; b4 != 0x40+byte(clientID) || b5 != 0x06 {
		return fmt.Errorf("Invalid 0x640 + client id %X %X\n", b5, b4)
	}

	if b7 := frm.Data[7]; b7 != 0x00 {
		return fmt.Errorf("Invalid byte 7 %X", b7)
	}

	return nil
}
