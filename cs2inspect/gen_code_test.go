package cs2inspect

import (
	"math"
	"strings"
	"testing"
)

func TestToGenCodeKeychainWithPaintKit(t *testing.T) {
	wear := float32(0.0)
	pk := uint32(929)
	item := &ItemPreviewData{
		DefIndex:   1355,
		PaintIndex: 0,
		PaintSeed:  0,
		PaintWear:  &wear,
		Keychains:  []Sticker{{Slot: 0, StickerID: 37, Wear: &wear, PaintKit: &pk}},
	}
	got := ToGenCode(item, "")
	tokens := strings.Fields(got)
	n := len(tokens)
	if n < 3 {
		t.Fatalf("expected at least 3 tokens, got %d: %q", n, got)
	}
	if tokens[n-3] != "37" || tokens[n-2] != "0" || tokens[n-1] != "929" {
		t.Errorf("expected last 3 tokens to be '37 0 929', got: %q %q %q", tokens[n-3], tokens[n-2], tokens[n-1])
	}
}

func TestToGenCodeKeychainWithoutPaintKitNoExtraToken(t *testing.T) {
	wear := float32(0.0)
	item := &ItemPreviewData{
		DefIndex:   7,
		PaintIndex: 0,
		PaintSeed:  0,
		PaintWear:  &wear,
		Keychains:  []Sticker{{Slot: 0, StickerID: 36, Wear: &wear}},
	}
	got := ToGenCode(item, "")
	tokens := strings.Fields(got)
	n := len(tokens)
	if n < 2 {
		t.Fatalf("expected at least 2 tokens, got %d: %q", n, got)
	}
	if tokens[n-2] != "36" || tokens[n-1] != "0" {
		t.Errorf("expected last 2 tokens to be '36 0', got: %q %q", tokens[n-2], tokens[n-1])
	}
}

func TestGenCodeFromLinkSlabUrl(t *testing.T) {
	slabUrl := "steam://run/730//+csgo_econ_action_preview%20819181994A8BA181A982B189E981F181238086898191A4E1208698F309C9"
	code, err := GenCodeFromLink(slabUrl, "")
	if err != nil {
		t.Fatalf("GenCodeFromLink() error: %v", err)
	}
	tokens := strings.Fields(code)
	n := len(tokens)
	if n < 3 {
		t.Fatalf("expected at least 3 tokens, got %d: %q", n, code)
	}
	if tokens[n-3] != "37" || tokens[n-2] != "0" || tokens[n-1] != "929" {
		t.Errorf("expected last 3 tokens to be '37 0 929', got: %q %q %q", tokens[n-3], tokens[n-2], tokens[n-1])
	}
}

func TestToGenCodeBasic(t *testing.T) {
	wear := float32(0.22540508)
	item := &ItemPreviewData{
		DefIndex:   7,
		PaintIndex: 474,
		PaintSeed:  306,
		PaintWear:  &wear,
	}
	got := ToGenCode(item, "!gen")
	want := "!gen 7 474 306 0.22540508"
	if got != want {
		t.Errorf("ToGenCode() = %q, want %q", got, want)
	}
}

func TestToGenCodeWithStickerAndKeychain(t *testing.T) {
	wear := float32(0.22540508)
	sWear := float32(0.0)
	kcWear := float32(0.0)
	item := &ItemPreviewData{
		DefIndex:   7,
		PaintIndex: 941,
		PaintSeed:  2,
		PaintWear:  &wear,
		Stickers:   []Sticker{{Slot: 2, StickerID: 7203, Wear: &sWear}},
		Keychains:  []Sticker{{Slot: 0, StickerID: 36, Wear: &kcWear}},
	}
	got := ToGenCode(item, "!g")
	want := "!g 7 941 2 0.22540508 0 0 0 0 7203 0 0 0 0 0 36 0"
	if got != want {
		t.Errorf("ToGenCode() = %q, want %q", got, want)
	}
}

func TestParseGenCodeBasic(t *testing.T) {
	item, err := ParseGenCode("!gen 7 474 306 0.22540508")
	if err != nil {
		t.Fatalf("ParseGenCode() error: %v", err)
	}
	if item.DefIndex != 7 {
		t.Errorf("DefIndex = %d, want 7", item.DefIndex)
	}
	if item.PaintIndex != 474 {
		t.Errorf("PaintIndex = %d, want 474", item.PaintIndex)
	}
	if item.PaintSeed != 306 {
		t.Errorf("PaintSeed = %d, want 306", item.PaintSeed)
	}
	if item.PaintWear == nil || math.Abs(float64(*item.PaintWear)-0.22540508) > 1e-5 {
		t.Errorf("PaintWear = %v, want ~0.22540508", item.PaintWear)
	}
}

func TestParseGenCodeWithStickerAndKeychain(t *testing.T) {
	item, err := ParseGenCode("!g 7 941 2 0.22540508 0 0 0 0 7203 0 0 0 0 0 36 0")
	if err != nil {
		t.Fatalf("ParseGenCode() error: %v", err)
	}
	if len(item.Stickers) != 1 {
		t.Fatalf("Stickers len = %d, want 1", len(item.Stickers))
	}
	if item.Stickers[0].StickerID != 7203 {
		t.Errorf("Sticker StickerID = %d, want 7203", item.Stickers[0].StickerID)
	}
	if len(item.Keychains) != 1 {
		t.Fatalf("Keychains len = %d, want 1", len(item.Keychains))
	}
	if item.Keychains[0].StickerID != 36 {
		t.Errorf("Keychain StickerID = %d, want 36", item.Keychains[0].StickerID)
	}
}

func TestGenCodeFromLinkFromHex(t *testing.T) {
	url, err := Generate(7, 474, 306, 0.22540508, nil)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	hex := strings.TrimPrefix(url, InspectBase)
	code, err := GenCodeFromLink(hex, "!gen")
	if err != nil {
		t.Fatalf("GenCodeFromLink() error: %v", err)
	}
	if !strings.HasPrefix(code, "!gen 7 474 306") {
		t.Errorf("GenCodeFromLink() = %q, want prefix %q", code, "!gen 7 474 306")
	}
}

func TestGenCodeFromLinkFromFullURL(t *testing.T) {
	url, err := Generate(7, 474, 306, 0.22540508, nil)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	code, err := GenCodeFromLink(url, "!gen")
	if err != nil {
		t.Fatalf("GenCodeFromLink() error: %v", err)
	}
	if !strings.HasPrefix(code, "!gen 7 474 306") {
		t.Errorf("GenCodeFromLink() = %q, want prefix %q", code, "!gen 7 474 306")
	}
}

func TestGenerateRoundtrip(t *testing.T) {
	url, err := Generate(7, 474, 306, 0.22540508, nil)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if !strings.HasPrefix(url, InspectBase) {
		t.Errorf("URL does not start with InspectBase: %q", url)
	}

	hex := strings.TrimPrefix(url, InspectBase)
	item, err := Deserialize(hex)
	if err != nil {
		t.Fatalf("Deserialize() error: %v", err)
	}
	if item.DefIndex != 7 {
		t.Errorf("DefIndex = %d, want 7", item.DefIndex)
	}
	if item.PaintIndex != 474 {
		t.Errorf("PaintIndex = %d, want 474", item.PaintIndex)
	}
}
