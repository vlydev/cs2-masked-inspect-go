package cs2inspect

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"math"
	"regexp"
	"strings"
)

// ------------------------------------------------------------------
// Compiled regular expressions
// ------------------------------------------------------------------

var (
	hybridURLRe  = regexp.MustCompile(`(?i)S\d+A\d+D([0-9A-Fa-f]+)$`)
	inspectURLRe = regexp.MustCompile(`(?i)(?:%20|\s|\+)A([0-9A-Fa-f]+)`)
	maskedURLRe  = regexp.MustCompile(`(?i)csgo_econ_action_preview(?:%20|\s|\+)([0-9A-Fa-f]{10,})$`)
	classicURLRe = regexp.MustCompile(`(?i)csgo_econ_action_preview(?:%20|\s)[SM]\d+A\d+D\d+$`)
	hexLettersRe = regexp.MustCompile(`[A-Fa-f]`)
)

// ------------------------------------------------------------------
// Checksum
// ------------------------------------------------------------------

// computeChecksum computes the 4-byte big-endian checksum used in the binary format.
// buffer = [0x00] + proto_bytes
// crc = crc32(buffer) using IEEE polynomial
// xored = (crc & 0xFFFF) ^ (uint32(len(proto_bytes)) * crc)  [unsigned 32-bit]
// checksum = big-endian uint32 of xored
func computeChecksum(protoBytes []byte) [4]byte {
	buf := make([]byte, 1+len(protoBytes))
	buf[0] = 0x00
	copy(buf[1:], protoBytes)

	crcVal := crc32.ChecksumIEEE(buf)
	protoLen := uint32(len(protoBytes))
	xored := (crcVal & 0xFFFF) ^ (protoLen * crcVal)

	var result [4]byte
	binary.BigEndian.PutUint32(result[:], xored)
	return result
}

// ------------------------------------------------------------------
// URL extraction
// ------------------------------------------------------------------

// extractHex extracts the hex payload from any supported inspect link format.
func extractHex(input string) string {
	stripped := strings.TrimSpace(input)

	// 1. Hybrid format: S\d+A\d+D<hexproto> — if hex part contains a-f letters
	if m := hybridURLRe.FindStringSubmatch(stripped); m != nil {
		if hexLettersRe.MatchString(m[1]) {
			return m[1]
		}
	}

	// 2. Classic/market URL: A<hex> preceded by %20, space, or +
	// Only use if stripped A yields even-length hex.
	if m := inspectURLRe.FindStringSubmatch(stripped); m != nil {
		if len(m[1])%2 == 0 {
			return m[1]
		}
	}

	// 3. Pure masked format: csgo_econ_action_preview%20<hexblob>
	if m := maskedURLRe.FindStringSubmatch(stripped); m != nil {
		return m[1]
	}

	// 4. Bare hex — strip whitespace
	return strings.ReplaceAll(stripped, " ", "")
}

// ------------------------------------------------------------------
// Sticker encode / decode
// ------------------------------------------------------------------

func encodeSticker(s Sticker) []byte {
	w := &ProtoWriter{}
	w.WriteUint32(1, s.Slot)
	w.WriteUint32(2, s.StickerID)
	if s.Wear != nil {
		w.WriteFloat32Fixed(3, *s.Wear)
	}
	if s.Scale != nil {
		w.WriteFloat32Fixed(4, *s.Scale)
	}
	if s.Rotation != nil {
		w.WriteFloat32Fixed(5, *s.Rotation)
	}
	w.WriteUint32(6, s.TintID)
	if s.OffsetX != nil {
		w.WriteFloat32Fixed(7, *s.OffsetX)
	}
	if s.OffsetY != nil {
		w.WriteFloat32Fixed(8, *s.OffsetY)
	}
	if s.OffsetZ != nil {
		w.WriteFloat32Fixed(9, *s.OffsetZ)
	}
	w.WriteUint32(10, s.Pattern)
	if s.HighlightReel != nil {
		w.WriteUint32(11, *s.HighlightReel)
	}
	return w.Bytes()
}

func decodeSticker(data []byte) (Sticker, error) {
	r := NewProtoReader(data)
	fields, err := r.ReadAllFields()
	if err != nil {
		return Sticker{}, fmt.Errorf("decoding sticker: %w", err)
	}

	var s Sticker
	for _, f := range fields {
		switch f.FieldNum {
		case 1:
			s.Slot = uint32(f.Varint)
		case 2:
			s.StickerID = uint32(f.Varint)
		case 3:
			v := f.Float32LE()
			s.Wear = &v
		case 4:
			v := f.Float32LE()
			s.Scale = &v
		case 5:
			v := f.Float32LE()
			s.Rotation = &v
		case 6:
			s.TintID = uint32(f.Varint)
		case 7:
			v := f.Float32LE()
			s.OffsetX = &v
		case 8:
			v := f.Float32LE()
			s.OffsetY = &v
		case 9:
			v := f.Float32LE()
			s.OffsetZ = &v
		case 10:
			s.Pattern = uint32(f.Varint)
		case 11:
			v := uint32(f.Varint)
			s.HighlightReel = &v
		}
	}

	return s, nil
}

// ------------------------------------------------------------------
// ItemPreviewData encode / decode
// ------------------------------------------------------------------

func encodeItem(item *ItemPreviewData) []byte {
	w := &ProtoWriter{}
	w.WriteUint32(1, item.AccountID)
	w.WriteUint64(2, item.ItemID)
	w.WriteUint32(3, item.DefIndex)
	w.WriteUint32(4, item.PaintIndex)
	w.WriteUint32(5, item.Rarity)
	w.WriteUint32(6, item.Quality)

	// PaintWear: float32 reinterpreted as uint32 varint
	if item.PaintWear != nil {
		w.WriteUint32(7, math.Float32bits(*item.PaintWear))
	}

	w.WriteUint32(8, item.PaintSeed)
	w.WriteUint32(9, item.KillEaterScoreType)
	w.WriteUint32(10, item.KillEaterValue)
	w.WriteString(11, item.CustomName)

	for _, sticker := range item.Stickers {
		w.WriteRawBytes(12, encodeSticker(sticker))
	}

	w.WriteUint32(13, item.Inventory)
	w.WriteUint32(14, item.Origin)
	w.WriteUint32(15, item.QuestID)
	w.WriteUint32(16, item.DropReason)
	w.WriteUint32(17, item.MusicIndex)
	w.WriteInt32(18, item.EntIndex)
	w.WriteUint32(19, item.PetIndex)

	for _, kc := range item.Keychains {
		w.WriteRawBytes(20, encodeSticker(kc))
	}

	return w.Bytes()
}

func decodeItem(data []byte) (*ItemPreviewData, error) {
	r := NewProtoReader(data)
	fields, err := r.ReadAllFields()
	if err != nil {
		return nil, fmt.Errorf("decoding item: %w", err)
	}

	item := &ItemPreviewData{}

	for _, f := range fields {
		switch f.FieldNum {
		case 1:
			item.AccountID = uint32(f.Varint)
		case 2:
			item.ItemID = f.Varint
		case 3:
			item.DefIndex = uint32(f.Varint)
		case 4:
			item.PaintIndex = uint32(f.Varint)
		case 5:
			item.Rarity = uint32(f.Varint)
		case 6:
			item.Quality = uint32(f.Varint)
		case 7:
			v := math.Float32frombits(uint32(f.Varint))
			item.PaintWear = &v
		case 8:
			item.PaintSeed = uint32(f.Varint)
		case 9:
			item.KillEaterScoreType = uint32(f.Varint)
		case 10:
			item.KillEaterValue = uint32(f.Varint)
		case 11:
			item.CustomName = string(f.Bytes)
		case 12:
			s, err := decodeSticker(f.Bytes)
			if err != nil {
				return nil, err
			}
			item.Stickers = append(item.Stickers, s)
		case 13:
			item.Inventory = uint32(f.Varint)
		case 14:
			item.Origin = uint32(f.Varint)
		case 15:
			item.QuestID = uint32(f.Varint)
		case 16:
			item.DropReason = uint32(f.Varint)
		case 17:
			item.MusicIndex = uint32(f.Varint)
		case 18:
			item.EntIndex = int32(f.Varint)
		case 19:
			item.PetIndex = uint32(f.Varint)
		case 20:
			kc, err := decodeSticker(f.Bytes)
			if err != nil {
				return nil, err
			}
			item.Keychains = append(item.Keychains, kc)
		}
	}

	return item, nil
}

// ------------------------------------------------------------------
// Public API
// ------------------------------------------------------------------

// Serialize encodes an ItemPreviewData to an uppercase hex inspect-link payload.
//
// The returned string can be appended to a steam:// inspect URL or used standalone.
// The key_byte is always 0x00 (no XOR applied).
//
// Returns an error if:
//   - PaintWear is set and outside [0.0, 1.0]
//   - CustomName exceeds 100 characters
func Serialize(data *ItemPreviewData) (string, error) {
	if data.PaintWear != nil && (*data.PaintWear < 0.0 || *data.PaintWear > 1.0) {
		return "", fmt.Errorf("paintwear must be in [0.0, 1.0], got %v", *data.PaintWear)
	}
	if len([]rune(data.CustomName)) > 100 {
		return "", fmt.Errorf("customname must not exceed 100 characters, got %d", len([]rune(data.CustomName)))
	}

	protoBytes := encodeItem(data)

	// Build: [0x00] + proto_bytes
	buf := make([]byte, 1+len(protoBytes))
	buf[0] = 0x00
	copy(buf[1:], protoBytes)

	// Append checksum
	checksum := computeChecksum(protoBytes)
	result := append(buf, checksum[:]...)

	return strings.ToUpper(hex.EncodeToString(result)), nil
}

// Deserialize decodes an inspect-link hex payload (or full URL) into an ItemPreviewData.
//
// Accepts:
//   - A raw uppercase or lowercase hex string
//   - A full steam://rungame/... inspect URL
//   - A CS2-style csgo://rungame/... URL
//
// Handles the XOR obfuscation used in native CS2 links.
//
// Returns an error if:
//   - The payload exceeds 4096 hex chars
//   - The payload is too short (< 6 bytes) or invalid hex
func Deserialize(input string) (*ItemPreviewData, error) {
	hexStr := extractHex(input)

	if len(hexStr) > 4096 {
		preview := input
		if len(preview) > 64 {
			preview = preview[:64] + "..."
		}
		return nil, fmt.Errorf("payload too long (max 4096 hex chars): %q", preview)
	}

	raw, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("payload too short or invalid hex: %q", input)
	}

	if len(raw) < 6 {
		return nil, fmt.Errorf("payload too short or invalid hex: %q", input)
	}

	key := raw[0]
	decrypted := raw

	if key != 0 {
		decrypted = make([]byte, len(raw))
		for i, b := range raw {
			decrypted[i] = b ^ key
		}
	}

	// Layout: [key_byte] [proto_bytes] [4-byte checksum]
	protoBytes := decrypted[1 : len(decrypted)-4]

	return decodeItem(protoBytes)
}

// IsMasked returns true if the link contains a decodable protobuf payload
// that can be decoded offline (pure masked or hybrid with hex proto).
func IsMasked(link string) bool {
	s := strings.TrimSpace(link)
	if maskedURLRe.MatchString(s) {
		return true
	}
	if m := hybridURLRe.FindStringSubmatch(s); m != nil {
		return hexLettersRe.MatchString(m[1])
	}
	return false
}

// IsClassic returns true if the link is a classic S/A/D inspect URL
// with a decimal D value (not a hex payload).
func IsClassic(link string) bool {
	return classicURLRe.MatchString(strings.TrimSpace(link))
}
