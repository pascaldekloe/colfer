// automatically generated, do not modify

package bench

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type FlatBuffers struct {
	_tab flatbuffers.Table
}

func (rcv *FlatBuffers) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *FlatBuffers) Key() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FlatBuffers) Host() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FlatBuffers) Addr(j int) int8 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetInt8(a + flatbuffers.UOffsetT(j * 1))
	}
	return 0
}

func (rcv *FlatBuffers) AddrLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *FlatBuffers) Port() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FlatBuffers) Size() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FlatBuffers) Hash() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FlatBuffers) Ratio() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FlatBuffers) Route() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func FlatBuffersStart(builder *flatbuffers.Builder) { builder.StartObject(8) }
func FlatBuffersAddKey(builder *flatbuffers.Builder, key int64) { builder.PrependInt64Slot(0, key, 0) }
func FlatBuffersAddHost(builder *flatbuffers.Builder, host flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(host), 0) }
func FlatBuffersAddAddr(builder *flatbuffers.Builder, addr flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(addr), 0) }
func FlatBuffersStartAddrVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(1, numElems, 1)
}
func FlatBuffersAddPort(builder *flatbuffers.Builder, port int32) { builder.PrependInt32Slot(3, port, 0) }
func FlatBuffersAddSize(builder *flatbuffers.Builder, size int64) { builder.PrependInt64Slot(4, size, 0) }
func FlatBuffersAddHash(builder *flatbuffers.Builder, hash uint64) { builder.PrependUint64Slot(5, hash, 0) }
func FlatBuffersAddRatio(builder *flatbuffers.Builder, ratio float64) { builder.PrependFloat64Slot(6, ratio, 0) }
func FlatBuffersAddRoute(builder *flatbuffers.Builder, route byte) { builder.PrependByteSlot(7, route, 0) }
func FlatBuffersEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
