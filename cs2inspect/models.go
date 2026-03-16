// Package cs2inspect provides encoding and decoding of CS2 masked inspect links.
//
// Binary format:
//
//	[key_byte] [proto_bytes XOR'd with key] [4-byte checksum XOR'd with key]
//
// For tool-generated links key_byte = 0x00 (no XOR needed).
// For native CS2 links key_byte != 0x00 — every byte is XOR'd before parsing.
package cs2inspect

// Sticker represents a sticker or keychain applied to a CS2 item.
//
// Maps to the Sticker protobuf message nested inside CEconItemPreviewDataBlock.
// The same message is used for both stickers (field 12) and keychains (field 20).
type Sticker struct {
	Slot          uint32
	StickerID     uint32
	Wear          *float32 // wire type 5 (fixed32 LE), nil = omitted
	Scale         *float32 // wire type 5 (fixed32 LE), nil = omitted
	Rotation      *float32 // wire type 5 (fixed32 LE), nil = omitted
	TintID        uint32
	OffsetX       *float32 // wire type 5 (fixed32 LE), nil = omitted
	OffsetY       *float32 // wire type 5 (fixed32 LE), nil = omitted
	OffsetZ       *float32 // wire type 5 (fixed32 LE), nil = omitted
	Pattern       uint32
	HighlightReel *uint32 // varint, nil = omitted
}

// ItemPreviewData represents a CS2 item as encoded in an inspect link.
//
// Fields map directly to the CEconItemPreviewDataBlock protobuf message
// used by the CS2 game coordinator.
//
// PaintWear is stored as float32 (IEEE 754). On the wire it is reinterpreted
// as a uint32 varint. A nil PaintWear means the field is absent.
type ItemPreviewData struct {
	AccountID          uint32
	ItemID             uint64
	DefIndex           uint32
	PaintIndex         uint32
	Rarity             uint32
	Quality            uint32
	PaintWear          *float32 // varint encoding of float32 bits, nil = omitted
	PaintSeed          uint32
	KillEaterScoreType uint32
	KillEaterValue     uint32
	CustomName         string
	Stickers           []Sticker
	Inventory          uint32
	Origin             uint32
	QuestID            uint32
	DropReason         uint32
	MusicIndex         uint32
	EntIndex           int32
	PetIndex           uint32
	Keychains          []Sticker
}
