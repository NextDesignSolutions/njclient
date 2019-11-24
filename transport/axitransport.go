package transport

import (
	"encoding/binary"
	"fmt"
	"github.com/NextDesignSolutions/axitransport"
	"github.com/NextDesignSolutions/bytehelpers"
	"github.com/NextDesignSolutions/njclient"
	"sync"
)

const (
	BYTES_PER_WORD = 0x4
	AXI_DATA_WIDTH = 0x4
	AXI_READ       = true
	AXI_WRITE      = false
)

type NextjtagAxiTransport struct {
	axi        *njclient.AxiHandle
	cache_attr *njclient.AxiCacheAttributes
	incr_mode  bool
	endianess  binary.ByteOrder
	mux        sync.Mutex
}

func NewClient(axi_handle *njclient.AxiHandle) (*NextjtagAxiTransport, error) {

	default_cache_attr := njclient.NewAxiCacheAttributes(axitransport.BUFFERABLE,
		axitransport.MODIFIABLE,
		axitransport.READ_ALLOC,
		axitransport.WRITE_ALLOC)
	if default_cache_attr == nil {
		return nil, fmt.Errorf("failed to generate default cache attributes")
	}

	default_incr_mode := true

	t := &NextjtagAxiTransport{
		axi:        axi_handle,
		cache_attr: default_cache_attr,
		incr_mode:  default_incr_mode,
		endianess:  binary.LittleEndian,
	}
	return t, nil
}

func (t *NextjtagAxiTransport) SetCacheAttributes(bufferable bool, modifiable bool, read_alloc bool, write_alloc bool) error {
	t.mux.Lock()
	new_cache_attr := njclient.NewAxiCacheAttributes(bufferable,
		modifiable,
		read_alloc,
		write_alloc)
	if new_cache_attr == nil {
		t.mux.Unlock()
		return fmt.Errorf("failed to generate new cache attributes")
	}
	t.cache_attr = new_cache_attr
	t.mux.Unlock()
	return nil
}

func (t *NextjtagAxiTransport) SetIncrMode(mode bool) error {
	t.mux.Lock()
	t.incr_mode = mode
	t.mux.Unlock()
	return nil
}
func (t *NextjtagAxiTransport) SetEndianess(endian bool) error {
	t.mux.Lock()
	if endian == axitransport.BIG_ENDIAN {
		t.endianess = binary.BigEndian
	} else {
		t.endianess = binary.LittleEndian
	}
	t.mux.Unlock()
	return nil
}

func (t *NextjtagAxiTransport) AccessSizeCheck(size int) (err error) {
	if size%AXI_DATA_WIDTH != 0 {
		return fmt.Errorf("access size(%d) must be a multiple of %d bytes\n", size, AXI_DATA_WIDTH)
	}
	return nil
}

func (t *NextjtagAxiTransport) ReadAxi(addr uint64, size int) (data []byte, err error) {
	err = t.AccessSizeCheck(size)
	if err != nil {
		return nil, err
	}
	count := size / AXI_DATA_WIDTH
	opts := njclient.NewAxiTransactionOptions(t.incr_mode, count)
	if opts == nil {
		return nil, fmt.Errorf("failed to create AxiTransactionOptions")
	}

	trans := njclient.NewAxiTransaction(addr, AXI_READ, opts, t.cache_attr, nil)
	if trans == nil {
		return nil, fmt.Errorf("failed to create axi read transaction addr:0x%08x size:0x%x", addr, size)
	}

	result, err := t.axi.IssueTransaction(trans)
	if err != nil {
		return nil, fmt.Errorf("AXI read failed addr:0x%08x size:0x%x", addr, size)
	}

	if result.Response != "OKAY" {
		return nil, fmt.Errorf("AXI read failed, received error from hardware, addr:0x%08x size:0x%x", addr, size)
	}
	if result.Value == nil || len(*result.Value) == 0 {
		return nil, fmt.Errorf("no data returned for AXI read addr:0x%08x size:0x%x", addr, size)
	}

	data, err = bytehelpers.Uint32sToBytes(*result.Value, t.endianess)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data to bytes for AXI read addr:0x%08x size:0x%x", addr, size)
	}
	return data, nil
}
func (t *NextjtagAxiTransport) WriteAxi(addr uint64, data []byte) (err error) {
	size := len(data)
	err = t.AccessSizeCheck(size)
	if err != nil {
		return err
	}
	count := size / AXI_DATA_WIDTH

	write_data, err := bytehelpers.BytesToUint32s(data, t.endianess)
	if err != nil {
		return fmt.Errorf("failed to encode data to words for AXI write addr:0x%08x size:0x%x", addr, size)
	}

	opts := njclient.NewAxiTransactionOptions(t.incr_mode, count)
	if opts == nil {
		return fmt.Errorf("failed to create AxiTransactionOptions")
	}

	trans := njclient.NewAxiTransaction(addr, AXI_WRITE, opts, t.cache_attr, &write_data)
	if trans == nil {
		return fmt.Errorf("failed to create axi write transaction addr:0x%08x size:0x%x", addr, size)
	}

	result, err := t.axi.IssueTransaction(trans)
	if err != nil {
		return fmt.Errorf("AXI write failed addr:0x%08x size:0x%x", addr, size)
	}

	if result.Response != "OKAY" {
		return fmt.Errorf("AXI read failed, received error from hardware, addr:0x%08x size:0x%x", addr, size)
	}
	return nil
}

func (t *NextjtagAxiTransport) ReadAxi32(addr uint64) (value uint32, err error) {

	data, err := t.ReadAxi(addr, 4)
	if err != nil {
		return 0, err
	}
	u32s, err := bytehelpers.BytesToUint32s(data, t.endianess)
	if err != nil {
		return 0, err
	}
	return u32s[0], err
}

func (t *NextjtagAxiTransport) WriteAxi32(addr uint64, value uint32) (err error) {
	d := []uint32{value}
	write_data, err := bytehelpers.Uint32sToBytes(d, t.endianess)
	if err != nil {
		return err
	}
	err = t.WriteAxi(addr, write_data)
	return err
}

func (t *NextjtagAxiTransport) ReadAxi64(addr uint64) (value uint64, err error) {

	data, err := t.ReadAxi(addr, 8)
	if err != nil {
		return 0, err
	}
	u64s, err := bytehelpers.BytesToUint64s(data, t.endianess)
	if err != nil {
		return 0, err
	}
	return u64s[0], err
}

func (t *NextjtagAxiTransport) WriteAxi64(addr uint64, value uint64) (err error) {
	d := []uint64{value}
	write_data, err := bytehelpers.Uint64sToBytes(d, t.endianess)
	if err != nil {
		return err
	}
	err = t.WriteAxi(addr, write_data)
	return err
}
