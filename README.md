# cs2-masked-inspect (Go)

Pure Go library for encoding and decoding CS2 masked inspect links — no runtime dependencies, requires Go 1.21+.

[![Tests](https://github.com/vlydev/cs2-masked-inspect-go/actions/workflows/tests.yml/badge.svg)](https://github.com/vlydev/cs2-masked-inspect-go/actions/workflows/tests.yml)

## Installation

```
go get github.com/vlydev/cs2-masked-inspect-go
```

## Usage

### Deserialize a CS2 inspect link

```go
import "github.com/vlydev/cs2-masked-inspect-go/cs2inspect"

// Accepts raw hex, steam:// URLs, hybrid S/A/D URLs, and classic inspect URLs.
item, err := cs2inspect.Deserialize("steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20E3F3...")
if err != nil {
    log.Fatal(err)
}
fmt.Println(item.DefIndex)   // 7 (AK-47)
fmt.Println(item.PaintIndex) // 422
fmt.Println(item.PaintSeed)  // 922
if item.PaintWear != nil {
    fmt.Printf("%.5f\n", *item.PaintWear) // 0.04121
}
fmt.Println(len(item.Stickers)) // 5
```

### Serialize

```go
pw := float32(0.005411375779658556)
data := &cs2inspect.ItemPreviewData{
    DefIndex:   60,
    PaintIndex: 440,
    PaintSeed:  353,
    Rarity:     5,
    PaintWear:  &pw,
}
hexStr, err := cs2inspect.Serialize(data)
if err != nil {
    log.Fatal(err)
}
fmt.Println(hexStr)
// "00183C20B803280538E9A3C5DD0340E102C246A0D1"
```

### IsMasked / IsClassic

```go
// Returns true if the link contains a decodable protobuf payload.
cs2inspect.IsMasked("steam://rungame/730/.../+csgo_econ_action_preview%20A00183C...")
// true

// Returns true if the link is a classic S/A/D format with a decimal D value.
cs2inspect.IsClassic("steam://rungame/730/.../+csgo_econ_action_preview%20S123A456D789")
// true
```

### Working with stickers and keychains

```go
item, _ := cs2inspect.Deserialize(hexStr)

for _, s := range item.Stickers {
    fmt.Printf("Slot %d: StickerID=%d", s.Slot, s.StickerID)
    if s.Wear != nil {
        fmt.Printf(" wear=%.4f", *s.Wear)
    }
    fmt.Println()
}

for _, kc := range item.Keychains {
    fmt.Printf("Keychain StickerID=%d", kc.StickerID)
    if kc.HighlightReel != nil {
        fmt.Printf(" highlightReel=%d", *kc.HighlightReel)
    }
    fmt.Println()
}
```

## Gen codes

Generate a Steam inspect URL from item parameters (defindex, paintindex, paintseed, paintwear):

```go
import "github.com/vlydev/cs2-masked-inspect-go/cs2inspect"

// Generate a Steam inspect URL
url, err := cs2inspect.Generate(7, 474, 306, 0.22540508, nil)

// Convert to gen code
item := &cs2inspect.ItemPreviewData{DefIndex: 7, PaintIndex: 474, PaintSeed: 306}
pw := float32(0.22540508); item.PaintWear = &pw
code := cs2inspect.ToGenCode(item, "!gen") // "!gen 7 474 306 0.22540508"

// Parse a gen code
item2, err := cs2inspect.ParseGenCode("!gen 7 474 306 0.22540508")

// Convert an existing inspect link directly to a gen code
code, err := cs2inspect.GenCodeFromLink("steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20001A...", "!gen")
// "!gen 7 474 306 0.22540508"
```

## API Reference

### Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `Serialize` | `func Serialize(data *ItemPreviewData) (string, error)` | Encode item data to uppercase hex. Key byte is always `0x00`. |
| `Deserialize` | `func Deserialize(input string) (*ItemPreviewData, error)` | Decode hex payload or any inspect URL format. |
| `IsMasked` | `func IsMasked(link string) bool` | True if the link has an embedded protobuf payload. |
| `IsClassic` | `func IsClassic(link string) bool` | True if the link is a classic `S/A/D` format with decimal D. |

### ItemPreviewData fields

| Go Field | Proto Field | Wire Type | Description |
|----------|-------------|-----------|-------------|
| `AccountID` | 1 | varint | Steam account ID |
| `ItemID` | 2 | varint (uint64) | Asset ID |
| `DefIndex` | 3 | varint | Item definition index |
| `PaintIndex` | 4 | varint | Skin/paint index |
| `Rarity` | 5 | varint | Item rarity |
| `Quality` | 6 | varint | Item quality |
| `PaintWear` | 7 | varint (float32 bits) | Wear value 0.0–1.0, nil = omitted |
| `PaintSeed` | 8 | varint | Pattern seed (0–1000) |
| `KillEaterScoreType` | 9 | varint | StatTrak kill type |
| `KillEaterValue` | 10 | varint | StatTrak counter |
| `CustomName` | 11 | length-delimited | Name tag |
| `Stickers` | 12 | length-delimited (repeated) | Applied stickers |
| `Inventory` | 13 | varint | Inventory flags |
| `Origin` | 14 | varint | Item origin |
| `QuestID` | 15 | varint | Quest ID |
| `DropReason` | 16 | varint | Drop reason |
| `MusicIndex` | 17 | varint | Music kit index |
| `EntIndex` | 18 | varint (int32) | Entity index |
| `PetIndex` | 19 | varint | Pet/agent index |
| `Keychains` | 20 | length-delimited (repeated) | Keychains (same Sticker type) |

### Sticker fields

| Go Field | Proto Field | Wire Type | Description |
|----------|-------------|-----------|-------------|
| `Slot` | 1 | varint | Sticker slot (0–4) |
| `StickerID` | 2 | varint | Sticker definition ID |
| `Wear` | 3 | fixed32 LE | Sticker wear, nil = omitted |
| `Scale` | 4 | fixed32 LE | Scale, nil = omitted |
| `Rotation` | 5 | fixed32 LE | Rotation, nil = omitted |
| `TintID` | 6 | varint | Tint ID |
| `OffsetX` | 7 | fixed32 LE | X offset, nil = omitted |
| `OffsetY` | 8 | fixed32 LE | Y offset, nil = omitted |
| `OffsetZ` | 9 | fixed32 LE | Z offset, nil = omitted |
| `Pattern` | 10 | varint | Pattern index |
| `HighlightReel` | 11 | varint | Highlight reel, nil = omitted |

## Validation rules

**Serialize:**
- `PaintWear` must be in `[0.0, 1.0]` if set; returns error otherwise.
- `CustomName` must not exceed 100 characters; returns error otherwise.

**Deserialize:**
- Hex payload must not exceed 4096 characters; returns error otherwise.
- Decoded payload must be at least 6 bytes; returns error otherwise.

## Binary format

```
[key_byte] [proto_bytes XOR'd with key] [4-byte checksum XOR'd with key]
```

- `key_byte = 0x00` for tool-generated links (no XOR).
- `key_byte != 0x00` for native CS2 links (XOR every byte including key itself).

**Checksum:**
```
buffer  = [0x00] + proto_bytes
crc     = CRC32-IEEE(buffer)
xored   = (crc & 0xFFFF) ^ (uint32(len(proto_bytes)) * crc)  [unsigned 32-bit]
checksum = big-endian uint32(xored & 0xFFFFFFFF)
```

## Test vectors

| Vector | Key | DefIndex | PaintIndex | PaintSeed | PaintWear | Notes |
|--------|-----|----------|------------|-----------|-----------|-------|
| `00183C20B803280538E9A3C5DD0340E102C246A0D1` | 0x00 | 60 | 440 | 353 | ≈0.005411 | Tool-generated |
| `E3F33674…` (96 bytes) | 0xE3 | 7 (AK-47) | 422 | 922 | ≈0.04121 | Native CS2, 5 stickers |
| `00180720DA03…` | 0x00 | 7 | 474 | 306 | ≈0.6337 | CSFloat vector A |
| `00180720C80A…` | 0x00 | – | 1352 | – | ≈0.99 | CSFloat vector B, 4 stickers |
| `A2B2A2BA…` | 0xA2 | 1355 | – | – | nil | CSFloat vector C, 1 keychain |

## Running tests

```
go test ./...
```

## Contributing

Pull requests are welcome. Please ensure all tests pass and follow standard Go conventions.

## License

MIT — see [LICENSE](LICENSE).
