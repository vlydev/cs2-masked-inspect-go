package cs2inspect

import (
	"errors"
	"math"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Known test vectors
// ---------------------------------------------------------------------------

// A real CS2 item encoded with XOR key 0xE3
const nativeHex = "E3F3367440334DE2FBE4C345E0CBE0D3E7DB6943400AE0A379E481ECEBE2F36F" +
	"D9DE2BDB515EA6E30D74D981ECEBE3F37BCBDE640D475DA6E35EFCD881ECEBE3" +
	"F359D5DE37E9D75DA6436DD3DD81ECEBE3F366DCDE3F8F9BDDA69B43B6DE81EC" +
	"EBE3F33BC8DEBB1CA3DFA623F7DDDF8B71E293EBFD43382B"

// A tool-generated link with key 0x00
const toolHex = "00183C20B803280538E9A3C5DD0340E102C246A0D1"

const hybridURL = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20" +
	"S76561199323320483A50075495125D" +
	"1101C4C4FCD4AB10092D31B8143914211829A1FAE3FD125119591141117308191301" +
	"EA550C1111912E3C111151D12C413E6BAC54D1D29BAD731E191501B92C2C9B6BF92F5411C25B2A731E191501B92C2C" +
	"EA2B182E5411F7212A731E191501B92C2C4F89C12F549164592A799713611956F4339F"

const classicURL = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20" +
	"S76561199842063946A49749521570D2751293026650298712"

// CSFloat test vectors
const csfloatA = "00180720DA03280638FBEE88F90340B2026BC03C96"
const csfloatB = "00180720C80A280638A4E1F5FB03409A0562040800104C62040801104C62040802104C62040803104C6D4F5E30"
const csfloatC = "A2B2A2BA69A882A28AA192AECAA2D2B700A3A5AAA2B286FA7BA0D684BE72"

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

func roundtrip(t *testing.T, data *ItemPreviewData) *ItemPreviewData {
	t.Helper()
	hex, err := Serialize(data)
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}
	result, err := Deserialize(hex)
	if err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}
	return result
}

func float32Ptr(v float32) *float32 { return &v }
func uint32Ptr(v uint32) *uint32    { return &v }

// ---------------------------------------------------------------------------
// Deserialize tests — native XOR link (key 0xE3)
// ---------------------------------------------------------------------------

func TestDeserializeNative_ItemId(t *testing.T) {
	item, err := Deserialize(nativeHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.ItemID != 46876117973 {
		t.Errorf("expected ItemID=46876117973, got %d", item.ItemID)
	}
}

func TestDeserializeNative_DefIndex(t *testing.T) {
	item, err := Deserialize(nativeHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.DefIndex != 7 {
		t.Errorf("expected DefIndex=7 (AK-47), got %d", item.DefIndex)
	}
}

func TestDeserializeNative_PaintIndex(t *testing.T) {
	item, err := Deserialize(nativeHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintIndex != 422 {
		t.Errorf("expected PaintIndex=422, got %d", item.PaintIndex)
	}
}

func TestDeserializeNative_PaintSeed(t *testing.T) {
	item, err := Deserialize(nativeHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintSeed != 922 {
		t.Errorf("expected PaintSeed=922, got %d", item.PaintSeed)
	}
}

func TestDeserializeNative_PaintWear(t *testing.T) {
	item, err := Deserialize(nativeHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintWear == nil {
		t.Fatal("expected PaintWear to be non-nil")
	}
	if math.Abs(float64(*item.PaintWear)-0.04121) >= 0.0001 {
		t.Errorf("expected PaintWear≈0.04121, got %v", *item.PaintWear)
	}
}

func TestDeserializeNative_Rarity(t *testing.T) {
	item, err := Deserialize(nativeHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.Rarity != 3 {
		t.Errorf("expected Rarity=3, got %d", item.Rarity)
	}
}

func TestDeserializeNative_Quality(t *testing.T) {
	item, err := Deserialize(nativeHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.Quality != 4 {
		t.Errorf("expected Quality=4, got %d", item.Quality)
	}
}

func TestDeserializeNative_StickerCount(t *testing.T) {
	item, err := Deserialize(nativeHex)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Stickers) != 5 {
		t.Errorf("expected 5 stickers, got %d", len(item.Stickers))
	}
}

func TestDeserializeNative_StickerIDs(t *testing.T) {
	item, err := Deserialize(nativeHex)
	if err != nil {
		t.Fatal(err)
	}
	expected := []uint32{7436, 5144, 6970, 8069, 5592}
	if len(item.Stickers) != len(expected) {
		t.Fatalf("expected %d stickers, got %d", len(expected), len(item.Stickers))
	}
	for i, s := range item.Stickers {
		if s.StickerID != expected[i] {
			t.Errorf("sticker[%d]: expected StickerID=%d, got %d", i, expected[i], s.StickerID)
		}
	}
}

// ---------------------------------------------------------------------------
// Deserialize tests — tool-generated link (key 0x00)
// ---------------------------------------------------------------------------

func TestDeserializeTool_DefIndex(t *testing.T) {
	item, err := Deserialize(toolHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.DefIndex != 60 {
		t.Errorf("expected DefIndex=60, got %d", item.DefIndex)
	}
}

func TestDeserializeTool_PaintIndex(t *testing.T) {
	item, err := Deserialize(toolHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintIndex != 440 {
		t.Errorf("expected PaintIndex=440, got %d", item.PaintIndex)
	}
}

func TestDeserializeTool_PaintSeed(t *testing.T) {
	item, err := Deserialize(toolHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintSeed != 353 {
		t.Errorf("expected PaintSeed=353, got %d", item.PaintSeed)
	}
}

func TestDeserializeTool_PaintWear(t *testing.T) {
	item, err := Deserialize(toolHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintWear == nil {
		t.Fatal("expected PaintWear to be non-nil")
	}
	if math.Abs(float64(*item.PaintWear)-0.005411375779658556) >= 1e-7 {
		t.Errorf("expected PaintWear≈0.005411375779658556, got %v", *item.PaintWear)
	}
}

func TestDeserializeTool_Rarity(t *testing.T) {
	item, err := Deserialize(toolHex)
	if err != nil {
		t.Fatal(err)
	}
	if item.Rarity != 5 {
		t.Errorf("expected Rarity=5, got %d", item.Rarity)
	}
}

func TestDeserializeTool_LowercaseHex(t *testing.T) {
	item, err := Deserialize(strings.ToLower(toolHex))
	if err != nil {
		t.Fatal(err)
	}
	if item.DefIndex != 60 {
		t.Errorf("expected DefIndex=60 with lowercase hex, got %d", item.DefIndex)
	}
}

func TestDeserializeTool_SteamURL(t *testing.T) {
	url := "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20A" + toolHex
	item, err := Deserialize(url)
	if err != nil {
		t.Fatal(err)
	}
	if item.DefIndex != 60 {
		t.Errorf("expected DefIndex=60 from steam:// URL, got %d", item.DefIndex)
	}
}

func TestDeserializeTool_CsgoURLWithLiteralSpace(t *testing.T) {
	url := "csgo://rungame/730/76561202255233023/+csgo_econ_action_preview A" + toolHex
	item, err := Deserialize(url)
	if err != nil {
		t.Fatal(err)
	}
	if item.DefIndex != 60 {
		t.Errorf("expected DefIndex=60 from csgo:// URL, got %d", item.DefIndex)
	}
}

// ---------------------------------------------------------------------------
// Serialize tests
// ---------------------------------------------------------------------------

func TestSerialize_KnownOutput(t *testing.T) {
	pw := float32(0.005411375779658556)
	data := &ItemPreviewData{
		DefIndex:   60,
		PaintIndex: 440,
		PaintSeed:  353,
		PaintWear:  &pw,
		Rarity:     5,
	}
	result, err := Serialize(data)
	if err != nil {
		t.Fatal(err)
	}
	if result != toolHex {
		t.Errorf("expected %q, got %q", toolHex, result)
	}
}

func TestSerialize_Uppercase(t *testing.T) {
	data := &ItemPreviewData{DefIndex: 1}
	result, err := Serialize(data)
	if err != nil {
		t.Fatal(err)
	}
	if result != strings.ToUpper(result) {
		t.Errorf("expected uppercase, got %q", result)
	}
}

func TestSerialize_StartsWithDoubleZero(t *testing.T) {
	data := &ItemPreviewData{DefIndex: 1}
	result, err := Serialize(data)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(result, "00") {
		t.Errorf("expected result to start with '00', got %q", result)
	}
}

func TestSerialize_MinLength(t *testing.T) {
	data := &ItemPreviewData{DefIndex: 1}
	result, err := Serialize(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) < 12 {
		t.Errorf("expected length >= 12, got %d", len(result))
	}
}

// ---------------------------------------------------------------------------
// Round-trip tests
// ---------------------------------------------------------------------------

func TestRoundtrip_DefIndex(t *testing.T) {
	result := roundtrip(t, &ItemPreviewData{DefIndex: 7})
	if result.DefIndex != 7 {
		t.Errorf("expected DefIndex=7, got %d", result.DefIndex)
	}
}

func TestRoundtrip_PaintIndex(t *testing.T) {
	result := roundtrip(t, &ItemPreviewData{PaintIndex: 422})
	if result.PaintIndex != 422 {
		t.Errorf("expected PaintIndex=422, got %d", result.PaintIndex)
	}
}

func TestRoundtrip_PaintSeed(t *testing.T) {
	result := roundtrip(t, &ItemPreviewData{PaintSeed: 999})
	if result.PaintSeed != 999 {
		t.Errorf("expected PaintSeed=999, got %d", result.PaintSeed)
	}
}

func TestRoundtrip_PaintWearFloat32Precision(t *testing.T) {
	original := float32(0.123456789)
	// Expected after float32 round-trip
	expected := math.Float32frombits(math.Float32bits(original))

	result := roundtrip(t, &ItemPreviewData{PaintWear: &original})
	if result.PaintWear == nil {
		t.Fatal("expected PaintWear to be non-nil")
	}
	if math.Abs(float64(*result.PaintWear)-float64(expected)) >= 1e-7 {
		t.Errorf("expected PaintWear≈%v, got %v", expected, *result.PaintWear)
	}
}

func TestRoundtrip_LargeItemId(t *testing.T) {
	result := roundtrip(t, &ItemPreviewData{ItemID: 46876117973})
	if result.ItemID != 46876117973 {
		t.Errorf("expected ItemID=46876117973, got %d", result.ItemID)
	}
}

func TestRoundtrip_Stickers(t *testing.T) {
	data := &ItemPreviewData{
		DefIndex: 7,
		Stickers: []Sticker{
			{Slot: 0, StickerID: 7436},
			{Slot: 1, StickerID: 5144},
		},
	}
	result := roundtrip(t, data)
	if len(result.Stickers) != 2 {
		t.Fatalf("expected 2 stickers, got %d", len(result.Stickers))
	}
	if result.Stickers[0].StickerID != 7436 {
		t.Errorf("sticker[0]: expected StickerID=7436, got %d", result.Stickers[0].StickerID)
	}
	if result.Stickers[1].StickerID != 5144 {
		t.Errorf("sticker[1]: expected StickerID=5144, got %d", result.Stickers[1].StickerID)
	}
}

func TestRoundtrip_StickerSlot(t *testing.T) {
	data := &ItemPreviewData{
		Stickers: []Sticker{{Slot: 3, StickerID: 123}},
	}
	result := roundtrip(t, data)
	if result.Stickers[0].Slot != 3 {
		t.Errorf("expected Slot=3, got %d", result.Stickers[0].Slot)
	}
}

func TestRoundtrip_StickerWear(t *testing.T) {
	wear := float32(0.5)
	data := &ItemPreviewData{
		Stickers: []Sticker{{StickerID: 1, Wear: &wear}},
	}
	result := roundtrip(t, data)
	if result.Stickers[0].Wear == nil {
		t.Fatal("expected sticker Wear to be non-nil")
	}
	if math.Abs(float64(*result.Stickers[0].Wear)-0.5) >= 1e-6 {
		t.Errorf("expected Wear≈0.5, got %v", *result.Stickers[0].Wear)
	}
}

func TestRoundtrip_Keychains(t *testing.T) {
	data := &ItemPreviewData{
		Keychains: []Sticker{{Slot: 0, StickerID: 999, Pattern: 42}},
	}
	result := roundtrip(t, data)
	if len(result.Keychains) != 1 {
		t.Fatalf("expected 1 keychain, got %d", len(result.Keychains))
	}
	if result.Keychains[0].StickerID != 999 {
		t.Errorf("expected StickerID=999, got %d", result.Keychains[0].StickerID)
	}
	if result.Keychains[0].Pattern != 42 {
		t.Errorf("expected Pattern=42, got %d", result.Keychains[0].Pattern)
	}
}

func TestRoundtrip_CustomName(t *testing.T) {
	data := &ItemPreviewData{DefIndex: 7, CustomName: "My Knife"}
	result := roundtrip(t, data)
	if result.CustomName != "My Knife" {
		t.Errorf("expected CustomName=%q, got %q", "My Knife", result.CustomName)
	}
}

func TestRoundtrip_RarityAndQuality(t *testing.T) {
	data := &ItemPreviewData{Rarity: 6, Quality: 9}
	result := roundtrip(t, data)
	if result.Rarity != 6 {
		t.Errorf("expected Rarity=6, got %d", result.Rarity)
	}
	if result.Quality != 9 {
		t.Errorf("expected Quality=9, got %d", result.Quality)
	}
}

func TestRoundtrip_FullItemWith5Stickers(t *testing.T) {
	pw := float32(0.04121)
	data := &ItemPreviewData{
		ItemID:     46876117973,
		DefIndex:   7,
		PaintIndex: 422,
		Rarity:     3,
		Quality:    4,
		PaintWear:  &pw,
		PaintSeed:  922,
		Stickers: []Sticker{
			{Slot: 0, StickerID: 7436},
			{Slot: 1, StickerID: 5144},
			{Slot: 2, StickerID: 6970},
			{Slot: 3, StickerID: 8069},
			{Slot: 4, StickerID: 5592},
		},
	}
	result := roundtrip(t, data)
	if result.DefIndex != 7 {
		t.Errorf("expected DefIndex=7, got %d", result.DefIndex)
	}
	if result.PaintIndex != 422 {
		t.Errorf("expected PaintIndex=422, got %d", result.PaintIndex)
	}
	if result.PaintSeed != 922 {
		t.Errorf("expected PaintSeed=922, got %d", result.PaintSeed)
	}
	if len(result.Stickers) != 5 {
		t.Fatalf("expected 5 stickers, got %d", len(result.Stickers))
	}
	expected := []uint32{7436, 5144, 6970, 8069, 5592}
	for i, s := range result.Stickers {
		if s.StickerID != expected[i] {
			t.Errorf("sticker[%d]: expected StickerID=%d, got %d", i, expected[i], s.StickerID)
		}
	}
}

func TestRoundtrip_EmptyStickers(t *testing.T) {
	data := &ItemPreviewData{DefIndex: 7, Stickers: []Sticker{}}
	result := roundtrip(t, data)
	if len(result.Stickers) != 0 {
		t.Errorf("expected empty stickers, got %d", len(result.Stickers))
	}
}

func TestRoundtrip_NilPaintWear(t *testing.T) {
	data := &ItemPreviewData{DefIndex: 7, PaintWear: nil}
	result := roundtrip(t, data)
	if result.PaintWear != nil {
		t.Errorf("expected nil PaintWear, got %v", *result.PaintWear)
	}
}

// ---------------------------------------------------------------------------
// IsMasked / IsClassic tests
// ---------------------------------------------------------------------------

func TestIsMasked_PureHexURL(t *testing.T) {
	url := "steam://run/730//+csgo_econ_action_preview%20" + toolHex
	if !IsMasked(url) {
		t.Error("expected IsMasked=true for pure hex URL")
	}
}

func TestIsMasked_NativeMaskedURL(t *testing.T) {
	url := "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20" + nativeHex
	if !IsMasked(url) {
		t.Error("expected IsMasked=true for native masked URL")
	}
}

func TestIsMasked_HybridURL(t *testing.T) {
	if !IsMasked(hybridURL) {
		t.Error("expected IsMasked=true for hybrid URL")
	}
}

func TestIsMasked_ClassicURL(t *testing.T) {
	if IsMasked(classicURL) {
		t.Error("expected IsMasked=false for classic URL")
	}
}

func TestIsClassic_ClassicURL(t *testing.T) {
	if !IsClassic(classicURL) {
		t.Error("expected IsClassic=true for classic URL")
	}
}

func TestIsClassic_MaskedURL(t *testing.T) {
	url := "steam://run/730//+csgo_econ_action_preview%20" + toolHex
	if IsClassic(url) {
		t.Error("expected IsClassic=false for masked URL")
	}
}

func TestIsClassic_HybridURL(t *testing.T) {
	if IsClassic(hybridURL) {
		t.Error("expected IsClassic=false for hybrid URL")
	}
}

// ---------------------------------------------------------------------------
// Hybrid URL deserialization
// ---------------------------------------------------------------------------

func TestDeserializeHybridURL_ItemId(t *testing.T) {
	item, err := Deserialize(hybridURL)
	if err != nil {
		t.Fatal(err)
	}
	if item.ItemID != 50075495125 {
		t.Errorf("expected ItemID=50075495125, got %d", item.ItemID)
	}
}

// ---------------------------------------------------------------------------
// Checksum test
// ---------------------------------------------------------------------------

func TestChecksum_KnownOutput(t *testing.T) {
	pw := float32(0.005411375779658556)
	data := &ItemPreviewData{
		DefIndex:   60,
		PaintIndex: 440,
		PaintSeed:  353,
		PaintWear:  &pw,
		Rarity:     5,
	}
	result, err := Serialize(data)
	if err != nil {
		t.Fatal(err)
	}
	if result != toolHex {
		t.Errorf("expected %q, got %q", toolHex, result)
	}
}

// ---------------------------------------------------------------------------
// Defensive / validation tests
// ---------------------------------------------------------------------------

func TestDeserialize_PayloadTooLong(t *testing.T) {
	// 4098 hex chars = 2049 bytes
	longHex := strings.Repeat("00", 2049)
	_, err := Deserialize(longHex)
	if err == nil {
		t.Error("expected error for payload exceeding 4096 hex chars")
	}
}

func TestDeserialize_PayloadTooShort(t *testing.T) {
	_, err := Deserialize("0000")
	if err == nil {
		t.Error("expected error for payload too short")
	}
}

func TestSerialize_PaintWearTooHigh(t *testing.T) {
	pw := float32(1.1)
	data := &ItemPreviewData{PaintWear: &pw}
	_, err := Serialize(data)
	if err == nil {
		t.Error("expected error for PaintWear > 1.0")
	}
}

func TestSerialize_PaintWearTooLow(t *testing.T) {
	pw := float32(-0.1)
	data := &ItemPreviewData{PaintWear: &pw}
	_, err := Serialize(data)
	if err == nil {
		t.Error("expected error for PaintWear < 0.0")
	}
}

func TestSerialize_PaintWearZeroBoundary(t *testing.T) {
	pw := float32(0.0)
	data := &ItemPreviewData{PaintWear: &pw}
	_, err := Serialize(data)
	if err != nil {
		t.Errorf("expected no error for PaintWear=0.0, got %v", err)
	}
}

func TestSerialize_PaintWearOneBoundary(t *testing.T) {
	pw := float32(1.0)
	data := &ItemPreviewData{PaintWear: &pw}
	_, err := Serialize(data)
	if err != nil {
		t.Errorf("expected no error for PaintWear=1.0, got %v", err)
	}
}

func TestSerialize_CustomNameTooLong(t *testing.T) {
	data := &ItemPreviewData{CustomName: strings.Repeat("a", 101)}
	_, err := Serialize(data)
	if err == nil {
		t.Error("expected error for CustomName > 100 chars")
	}
}

func TestSerialize_CustomNameExactly100(t *testing.T) {
	data := &ItemPreviewData{CustomName: strings.Repeat("a", 100)}
	_, err := Serialize(data)
	if err != nil {
		t.Errorf("expected no error for CustomName=100 chars, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// CSFloat test vectors
// ---------------------------------------------------------------------------

func TestCSFloatA_DefIndex(t *testing.T) {
	item, err := Deserialize(csfloatA)
	if err != nil {
		t.Fatal(err)
	}
	if item.DefIndex != 7 {
		t.Errorf("expected DefIndex=7, got %d", item.DefIndex)
	}
}

func TestCSFloatA_PaintIndex(t *testing.T) {
	item, err := Deserialize(csfloatA)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintIndex != 474 {
		t.Errorf("expected PaintIndex=474, got %d", item.PaintIndex)
	}
}

func TestCSFloatA_PaintSeed(t *testing.T) {
	item, err := Deserialize(csfloatA)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintSeed != 306 {
		t.Errorf("expected PaintSeed=306, got %d", item.PaintSeed)
	}
}

func TestCSFloatA_Rarity(t *testing.T) {
	item, err := Deserialize(csfloatA)
	if err != nil {
		t.Fatal(err)
	}
	if item.Rarity != 6 {
		t.Errorf("expected Rarity=6, got %d", item.Rarity)
	}
}

func TestCSFloatA_PaintWearNotNil(t *testing.T) {
	item, err := Deserialize(csfloatA)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintWear == nil {
		t.Error("expected PaintWear to be non-nil")
	}
}

func TestCSFloatA_PaintWearApprox(t *testing.T) {
	item, err := Deserialize(csfloatA)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintWear == nil {
		t.Fatal("expected PaintWear to be non-nil")
	}
	if math.Abs(float64(*item.PaintWear)-0.6337) >= 0.001 {
		t.Errorf("expected PaintWear≈0.6337, got %v", *item.PaintWear)
	}
}

func TestCSFloatB_StickerCount(t *testing.T) {
	item, err := Deserialize(csfloatB)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Stickers) != 4 {
		t.Errorf("expected 4 stickers, got %d", len(item.Stickers))
	}
}

func TestCSFloatB_StickerIDs(t *testing.T) {
	item, err := Deserialize(csfloatB)
	if err != nil {
		t.Fatal(err)
	}
	for i, s := range item.Stickers {
		if s.StickerID != 76 {
			t.Errorf("sticker[%d]: expected StickerID=76, got %d", i, s.StickerID)
		}
	}
}

func TestCSFloatB_PaintIndex(t *testing.T) {
	item, err := Deserialize(csfloatB)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintIndex != 1352 {
		t.Errorf("expected PaintIndex=1352, got %d", item.PaintIndex)
	}
}

func TestCSFloatB_PaintWear(t *testing.T) {
	item, err := Deserialize(csfloatB)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintWear == nil {
		t.Fatal("expected PaintWear to be non-nil")
	}
	if math.Abs(float64(*item.PaintWear)-0.99) >= 0.01 {
		t.Errorf("expected PaintWear≈0.99, got %v", *item.PaintWear)
	}
}

func TestCSFloatC_DefIndex(t *testing.T) {
	item, err := Deserialize(csfloatC)
	if err != nil {
		t.Fatal(err)
	}
	if item.DefIndex != 1355 {
		t.Errorf("expected DefIndex=1355, got %d", item.DefIndex)
	}
}

func TestCSFloatC_Quality(t *testing.T) {
	item, err := Deserialize(csfloatC)
	if err != nil {
		t.Fatal(err)
	}
	if item.Quality != 12 {
		t.Errorf("expected Quality=12, got %d", item.Quality)
	}
}

func TestCSFloatC_KeychainCount(t *testing.T) {
	item, err := Deserialize(csfloatC)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Keychains) != 1 {
		t.Errorf("expected 1 keychain, got %d", len(item.Keychains))
	}
}

func TestCSFloatC_KeychainHighlightReel(t *testing.T) {
	item, err := Deserialize(csfloatC)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Keychains) == 0 {
		t.Fatal("no keychains")
	}
	if item.Keychains[0].HighlightReel == nil {
		t.Fatal("expected HighlightReel to be non-nil")
	}
	if *item.Keychains[0].HighlightReel != 345 {
		t.Errorf("expected HighlightReel=345, got %d", *item.Keychains[0].HighlightReel)
	}
}

func TestCSFloatC_NoPaintWear(t *testing.T) {
	item, err := Deserialize(csfloatC)
	if err != nil {
		t.Fatal(err)
	}
	if item.PaintWear != nil {
		t.Errorf("expected nil PaintWear, got %v", *item.PaintWear)
	}
}

// ---------------------------------------------------------------------------
// HighlightReel and nullable PaintWear round-trips
// ---------------------------------------------------------------------------

func TestRoundtrip_HighlightReel(t *testing.T) {
	hr := uint32(345)
	data := &ItemPreviewData{
		DefIndex:  7,
		Keychains: []Sticker{{Slot: 0, StickerID: 36, HighlightReel: &hr}},
	}
	result := roundtrip(t, data)
	if len(result.Keychains) == 0 {
		t.Fatal("no keychains after roundtrip")
	}
	if result.Keychains[0].HighlightReel == nil {
		t.Fatal("expected HighlightReel to be non-nil")
	}
	if *result.Keychains[0].HighlightReel != 345 {
		t.Errorf("expected HighlightReel=345, got %d", *result.Keychains[0].HighlightReel)
	}
}

// ---------------------------------------------------------------------------
// Sticker Slab test vectors
//
// Sticker Slabs: defIndex=1355, quality=8, keychains[0].StickerID=37 (placeholder)
// keychains[0].PaintKit = actual slab variant ID
//
// URL A: rarity=5, paintKit=7256
// URL B: rarity=3, paintKit=275
// ---------------------------------------------------------------------------

const stickerSlabA = "steam://run/730//+csgo_econ_action_preview%20918191895A9BB191B994A199F991E191339096999181B4F149A98D5C0889"
const stickerSlabB = "steam://run/730//+csgo_econ_action_preview%20CBDBCBD300C1EBCBE3C8FBC3A3CBBBCB69CACCC3CBDBEEAB58C9B8B67C83"

func TestStickerSlabA_DefIndex(t *testing.T) {
	item, err := Deserialize(stickerSlabA)
	if err != nil {
		t.Fatal(err)
	}
	if item.DefIndex != 1355 {
		t.Errorf("expected DefIndex=1355, got %d", item.DefIndex)
	}
}

func TestStickerSlabA_Quality(t *testing.T) {
	item, err := Deserialize(stickerSlabA)
	if err != nil {
		t.Fatal(err)
	}
	if item.Quality != 8 {
		t.Errorf("expected Quality=8, got %d", item.Quality)
	}
}

func TestStickerSlabA_Rarity(t *testing.T) {
	item, err := Deserialize(stickerSlabA)
	if err != nil {
		t.Fatal(err)
	}
	if item.Rarity != 5 {
		t.Errorf("expected Rarity=5, got %d", item.Rarity)
	}
}

func TestStickerSlabA_KeychainCount(t *testing.T) {
	item, err := Deserialize(stickerSlabA)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Keychains) != 1 {
		t.Errorf("expected 1 keychain, got %d", len(item.Keychains))
	}
}

func TestStickerSlabA_KeychainStickerID(t *testing.T) {
	item, err := Deserialize(stickerSlabA)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Keychains) == 0 {
		t.Fatal("no keychains")
	}
	if item.Keychains[0].StickerID != 37 {
		t.Errorf("expected StickerID=37, got %d", item.Keychains[0].StickerID)
	}
}

func TestStickerSlabA_KeychainPaintKit(t *testing.T) {
	item, err := Deserialize(stickerSlabA)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Keychains) == 0 {
		t.Fatal("no keychains")
	}
	if item.Keychains[0].PaintKit == nil {
		t.Fatal("expected PaintKit to be non-nil")
	}
	if *item.Keychains[0].PaintKit != 7256 {
		t.Errorf("expected PaintKit=7256, got %d", *item.Keychains[0].PaintKit)
	}
}

func TestStickerSlabB_DefIndex(t *testing.T) {
	item, err := Deserialize(stickerSlabB)
	if err != nil {
		t.Fatal(err)
	}
	if item.DefIndex != 1355 {
		t.Errorf("expected DefIndex=1355, got %d", item.DefIndex)
	}
}

func TestStickerSlabB_Quality(t *testing.T) {
	item, err := Deserialize(stickerSlabB)
	if err != nil {
		t.Fatal(err)
	}
	if item.Quality != 8 {
		t.Errorf("expected Quality=8, got %d", item.Quality)
	}
}

func TestStickerSlabB_Rarity(t *testing.T) {
	item, err := Deserialize(stickerSlabB)
	if err != nil {
		t.Fatal(err)
	}
	if item.Rarity != 3 {
		t.Errorf("expected Rarity=3, got %d", item.Rarity)
	}
}

func TestStickerSlabB_KeychainCount(t *testing.T) {
	item, err := Deserialize(stickerSlabB)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Keychains) != 1 {
		t.Errorf("expected 1 keychain, got %d", len(item.Keychains))
	}
}

func TestStickerSlabB_KeychainStickerID(t *testing.T) {
	item, err := Deserialize(stickerSlabB)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Keychains) == 0 {
		t.Fatal("no keychains")
	}
	if item.Keychains[0].StickerID != 37 {
		t.Errorf("expected StickerID=37, got %d", item.Keychains[0].StickerID)
	}
}

func TestStickerSlabB_KeychainPaintKit(t *testing.T) {
	item, err := Deserialize(stickerSlabB)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.Keychains) == 0 {
		t.Fatal("no keychains")
	}
	if item.Keychains[0].PaintKit == nil {
		t.Fatal("expected PaintKit to be non-nil")
	}
	if *item.Keychains[0].PaintKit != 275 {
		t.Errorf("expected PaintKit=275, got %d", *item.Keychains[0].PaintKit)
	}
}

func TestRoundtrip_PaintKit(t *testing.T) {
	pk := uint32(7256)
	data := &ItemPreviewData{
		DefIndex:  1355,
		Quality:   8,
		Rarity:    5,
		Keychains: []Sticker{{Slot: 0, StickerID: 37, PaintKit: &pk}},
	}
	result := roundtrip(t, data)
	if len(result.Keychains) == 0 {
		t.Fatal("no keychains after roundtrip")
	}
	if result.Keychains[0].StickerID != 37 {
		t.Errorf("expected StickerID=37, got %d", result.Keychains[0].StickerID)
	}
	if result.Keychains[0].PaintKit == nil {
		t.Fatal("expected PaintKit to be non-nil")
	}
	if *result.Keychains[0].PaintKit != 7256 {
		t.Errorf("expected PaintKit=7256, got %d", *result.Keychains[0].PaintKit)
	}
}

// ---------------------------------------------------------------------------
// Malformed URLs (regression: must reject cleanly with ErrMalformedInspectLink)
// ---------------------------------------------------------------------------

func TestMalformedURLs(t *testing.T) {
	cases := []struct {
		name string
		url  string
	}{
		{"truncated mid-keychain (defindex=1, key=0xAD)", "steam://run/730//+csgo_econ_action_preview%20ADBD1050393912ACB5AC8D45AC85A99DA9956A116D5FAEED21ACCFB4A5AFBD348EB0ADAD2D9280ADADDD6F90EDA37510E84D8BEE11CFB4A5ACBD348EB0ADAD2D9280ADAD5D6F906F2B4C13E84D93D591CFB9A5ADBD419EB0ADAD2D9290ADF22010E8B72FB213CFB4A5ADBD549EB0ADAD2D9280ADADED6C90CFD43F10E892DFE513CFB4A5ADBD549EB0ADAD2D9280ADAD85EE902F952210E82EB8A613C52E2D2D2DA1DDA90FACBBA5ADBD89902923AAECE83"},
		{"truncated mid-keychain (defindex=9, key=0xEE)", "steam://run/730//+csgo_econ_action_preview%20EEFE3144332550EFF6E7CE28E8C6EADEEAD642323218EDAE4DEA8CFAE6ECFE35A7F302BFD6D1D3004FCB50ABAE5CF8528CF7E6EEFE0DA7F3394DDED1C3EEEE7EAFD39E25B3D3AB9EAE70D28CFAE6ECFE32A7F3EEEE6ED1D3595B17D3AB9E8E65538CFAE6ECFE64ADF3EEEE6ED1D3E597F3D3AB2EEA1AD58CF7E6EDFE0BCAF302BFD6D1C3EEEEEF2DD3AEF5F552ABEE31A855866D6E6E6EE29EE64CEFF8E6EEFED8D3B6CBCBACABFD70EED1A31F96E7AFBE5"},
		{"truncated mid-keychain (defindex=1, key=0x4A)", "steam://run/730//+csgo_econ_action_preview%204A5A8EFCB1B9F44B524B6AA24B624E7A4E72BACFF6B8490AD449285E42485AAB75576316457577CA2422760F4A413E712853424A5AA679574A4ACA75674A4A8A8A770A7F85760FD04246F4285342495AB279574A4ACA75674A4A8A0A7799714A750F0A5140F7285342495AAD7277EB4547F40F0A00EB7122C9CACACA463A4EE84B5D424A5A4C776C02A10A0F34A5C17407F0145C0A1AA"},
		{"truncated mid-keychain (AK-47 1035, key=0xFA)", "steam://run/730//+csgo_econ_action_preview%20FAEA5766387F45FBE2FDDA71F2D2FECAF3C2142C0C0EF9BA7CFFB2FAAAFA98EEF2F8EA3BFBD7FAFA3ABBC7EA2FB7C7BF9ACC47C698E3F2F9EA03C9E7FAFA7AC5D7FAFABA3AC780CA89C4BFFAAAD9C198E3F2F9EA03C9E7FAFA7AC5D7FAFACEB9C7B60177C4BFAAB82AC698EEF2F9EA03C9E7FAFA7AC5C7C11558C4BFFA43DBC198F5F2FBEA13DEC7759A1147BF9A16F64692797A7A7AF68AF258FBEFF2FAEADEC7DEAC32BBBF0BF9BAC4B760382CC5A2D24"},
		{"truncated mid-keychain (defindex=40, key=0x9F)", "steam://run/730//+csgo_econ_action_preview%209F8F4F504C7C219E87B7BF629EB79CAF9BA73F1D53419CDF0699FD8B979F8F5CCF82050686A0A2F4038821DA17F3FD22FD8B979F8F49D4825C6AB7A0A25F50B224DA7F0CF222FD8B979D8F5DCF82F9F9B9A0A2B7B25422DA6F6F1422FD8B979D8F5DCF822781DAA0A247B731A2DA7F92EF23FD86979F8F5DCF827EE5CBA0B29F9FDFDFA285AD4322DA9FFD8AA3F71C1F1F1F93EF873D9E88979F8FDDA243C05EDFDA4CFB44A0D202B77EDFCF3F339C3DF89"},
		{"truncated mid-keychain (M4A1-S 1130, key=0xFA)", "steam://run/730//+csgo_econ_action_preview%20FAEA5B24060844FBE2C6DA10F2D2FECAF3C2631A3308F9BA47F8B2FAAAFA98EEF2FEEA29B2E781EED4C5C7776B78C4BFFAFA21CD98EEF2FAEA09BCE7F02DD9C5C7FE6C57C7BF7A58ED4698EEF2FEEA6DBDE79C9CDCC5C7929F2F47BFFA461BC398E3F2FBEA15BDE7E57FD1C5D7FAFA8AB8C7C2CEDC46BF12973EC798EEF2FEEA24B8E781EED4C5C75696FCC5BF2AEBF2C792797A7A7AF68AF258FBEDF2FAEAD2C7B3683FBBBF065F8CC5B763A382BAAA0C5"},
		{"truncated mid-keychain (defindex=35, key=0x4D)", "steam://run/730//+csgo_econ_action_preview%204D5D9DF8C7D2F34C556E6DDC4C654E7D4975B2D2AEBB4E0DAB4B2F59454E5D9A70604D4DBD8C7002A356F308CD4603F62F5445495DEB745011C20F72604D4D798F70797F9FF3089D49ABF12F5945495DEB74502B2BAB737045B3FAF308FD4BE8F12F5945495DF868508081417270BFB9D2F308ADD283F12F5945495DF86850AC375972702FA3CBF3083D2EBFF125CECDCDCD413D5AEF4C5A454D5D567084E4F80C08254C547200CE3E1D0D1DF9D24C63938"},
		{"truncated mid-keychain (AK-47 1171, key=0xCF)", "steam://run/730//+csgo_econ_action_preview%20CFDF6258412F71CED7C8EF5CC6E7C9FFCBF7465B3B38CC8F7ACEADD6C7CEDF3BF2D2F2C5D8F0E2CFCFE40CF27F72E1728AF786F5F2ADC0C7CCDF3BF2F241D88EF18A0F03F6F2ADDBC7CCDF3CF2E2CFCF0F8DF27B3FB7F18A6F5ACDF2ADD6C7CCDF3CF2D2C518ECF0E2CFCF8F0EF2D5C29AF18ACFCFD8F5ADDBC7CFDF3CF2E2CFCF6D8DF24DB0DA718A2FA6DBF3A74C4F4F4FC3BFC76DCED8C7CFDF87F27B950C8E8A37D0A3F182A2B8A48F9F4332CD2C2B6"},
		{"truncated mid-keychain (defindex=1 1050, key=0xCE)", "steam://run/730//+csgo_econ_action_preview%20CEDE51082D1C70CFD6CFEE54C6E6CAFEC7F631274538CD8E2DC886CE9ECEACDAC6CDDE0C8DD3CECE4EF1F382996B738B4E5E8D75ACD7C6CEDE0C8DD3CECE4EF1E3CECE8E0FF3A56603708BECDBE170ACD7C6CEDE0C8DD3CECE4EF1E3CECEDE0FF37682A5708B0650FB70ACD7C6CEDE0C8DD3CECE4EF1E3CECEDE0FF34E3CEC738BDAC6F870ACD7C6CDDE798AD3CECE4EF1E3CECE0E0EF3333BC5F18BAE8E9FF2A64D4E4E4EC2BECA6CCFD9C6CEDECFF37C6"},
		{"odd-length bare hex", "ABC"},
		{"empty string", ""},
		{"non-hex characters", "ZZZZZZZZZZZZ"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Deserialize(tc.url)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, ErrMalformedInspectLink) {
				t.Errorf("expected errors.Is(err, ErrMalformedInspectLink) == true, got %v", err)
			}
		})
	}
}

func TestMalformedURLs_PayloadTooLong_WrapsErrMalformedInspectLink(t *testing.T) {
	longHex := strings.Repeat("00", 2049)
	_, err := Deserialize(longHex)
	if !errors.Is(err, ErrMalformedInspectLink) {
		t.Errorf("expected ErrMalformedInspectLink wrap, got %v", err)
	}
}
