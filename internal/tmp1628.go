package athenaLed

import (
	"fmt"
	"os"
)

const (
	tmp1628FrameLen   = 28
	tmp1628SlotMask   = 0b00011111
	tmp1628DevicePath = "/dev/tmp1628-led"
)

type tmp1628Device struct {
	file *os.File
}

var openTmp1628DeviceFile = func(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY, 0)
}

func openTmp1628Device(path string) (*tmp1628Device, error) {
	file, err := openTmp1628DeviceFile(path)
	if err != nil {
		return nil, err
	}
	return &tmp1628Device{file: file}, nil
}

func (dev *tmp1628Device) Close() error {
	if dev == nil || dev.file == nil {
		return nil
	}
	err := dev.file.Close()
	dev.file = nil
	return err
}

func (dev *tmp1628Device) WriteData(values []byte, status byte) error {
	if dev == nil || dev.file == nil {
		return fmt.Errorf("tmp1628 backend not initialized")
	}
	frame := encodeTMP1628Frame(values, status)
	_, err := dev.file.Write(frame)
	return err
}

func encodeTMP1628Frame(values []byte, status byte) []byte {
	var slots [tmp1628FrameLen]byte
	for i := 0; i < WIDTH && i < len(values); i++ {
		slots[i] = values[i] & tmp1628SlotMask
	}
	slots[tmp1628FrameLen-1] = status & tmp1628SlotMask

	frame := make([]byte, tmp1628FrameLen)
	for i := 0; i < tmp1628FrameLen/2; i++ {
		packed := uint16(slots[2*i]) | (uint16(slots[2*i+1]) << 5)
		frame[2*i] = byte(packed)
		frame[2*i+1] = byte(packed >> 8)
	}
	return frame
}
