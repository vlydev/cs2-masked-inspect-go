package cs2inspect

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// InspectBase is the Steam inspect URL prefix.
const InspectBase = "steam://rungame/730/76561202255233023/+csgo_econ_action_preview%20"

// formatFloat formats a float32, stripping trailing zeros (max 8 decimal places).
func formatFloat(v float32) string {
	s := strconv.FormatFloat(float64(v), 'f', 8, 32)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		return "0"
	}
	return s
}

func serializeStickerPairs(stickers []Sticker, padTo int) []string {
	result := []string{}
	filtered := []Sticker{}
	for _, s := range stickers {
		if s.StickerID != 0 {
			filtered = append(filtered, s)
		}
	}

	if padTo >= 0 {
		slotMap := map[uint32]Sticker{}
		for _, s := range filtered {
			slotMap[s.Slot] = s
		}
		for slot := 0; slot < padTo; slot++ {
			if s, ok := slotMap[uint32(slot)]; ok {
				var wear float32
				if s.Wear != nil {
					wear = *s.Wear
				}
				result = append(result, fmt.Sprintf("%d", s.StickerID), formatFloat(wear))
			} else {
				result = append(result, "0", "0")
			}
		}
	} else {
		// Sort by slot
		sorted := make([]Sticker, len(filtered))
		copy(sorted, filtered)
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[j].Slot < sorted[i].Slot {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
		for _, s := range sorted {
			var wear float32
			if s.Wear != nil {
				wear = *s.Wear
			}
			result = append(result, fmt.Sprintf("%d", s.StickerID), formatFloat(wear))
		}
	}

	return result
}

// ToGenCode converts an ItemPreviewData to a gen code string.
//
// The prefix is typically "!gen" or "!g". The gen code format is:
//
//	!gen {defindex} {paintindex} {paintseed} {paintwear} [{sticker pairs x5}] [{keychain pairs}]
func ToGenCode(item *ItemPreviewData, prefix string) string {
	var wearStr string
	if item.PaintWear != nil {
		wearStr = formatFloat(*item.PaintWear)
	} else {
		wearStr = "0"
	}

	parts := []string{
		fmt.Sprintf("%d", item.DefIndex),
		fmt.Sprintf("%d", item.PaintIndex),
		fmt.Sprintf("%d", item.PaintSeed),
		wearStr,
	}

	hasStickers := false
	for _, s := range item.Stickers {
		if s.StickerID != 0 {
			hasStickers = true
			break
		}
	}
	hasKeychains := false
	for _, s := range item.Keychains {
		if s.StickerID != 0 {
			hasKeychains = true
			break
		}
	}

	if hasStickers || hasKeychains {
		parts = append(parts, serializeStickerPairs(item.Stickers, 5)...)
		parts = append(parts, serializeStickerPairs(item.Keychains, -1)...)
	}

	payload := strings.Join(parts, " ")
	if prefix != "" {
		return prefix + " " + payload
	}
	return payload
}

// GenerateOptions holds optional parameters for Generate.
type GenerateOptions struct {
	Rarity    uint32
	Quality   uint32
	Stickers  []Sticker
	Keychains []Sticker
}

// Generate creates a full Steam inspect URL from item parameters.
//
// Returns an error if serialization fails (e.g., PaintWear out of range).
func Generate(defIndex, paintIndex, paintSeed uint32, paintWear float32, opts *GenerateOptions) (string, error) {
	if opts == nil {
		opts = &GenerateOptions{}
	}
	pw := paintWear
	data := &ItemPreviewData{
		DefIndex:   defIndex,
		PaintIndex: paintIndex,
		PaintSeed:  paintSeed,
		PaintWear:  &pw,
		Rarity:     opts.Rarity,
		Quality:    opts.Quality,
		Stickers:   opts.Stickers,
		Keychains:  opts.Keychains,
	}
	hex, err := Serialize(data)
	if err != nil {
		return "", err
	}
	return InspectBase + hex, nil
}

// GenCodeFromLink generates a gen code string from an existing CS2 inspect link.
//
// Deserializes the inspect link and converts the item data to gen code format.
// Returns an error if deserialization fails.
func GenCodeFromLink(hexOrUrl string, prefix string) (string, error) {
	item, err := Deserialize(hexOrUrl)
	if err != nil {
		return "", err
	}
	return ToGenCode(item, prefix), nil
}

// ParseGenCode parses a gen code string into an ItemPreviewData.
//
// Accepts codes like:
//
//	"!gen 7 474 306 0.22540508"
//	"7 941 2 0.22540508 0 0 0 0 7203 0 0 0 0 0 36 0"
//
// Returns an error if the code has fewer than 4 tokens.
func ParseGenCode(genCode string) (*ItemPreviewData, error) {
	tokens := strings.Fields(strings.TrimSpace(genCode))

	// Skip leading !-prefixed command
	if len(tokens) > 0 && strings.HasPrefix(tokens[0], "!") {
		tokens = tokens[1:]
	}

	if len(tokens) < 4 {
		return nil, fmt.Errorf("gen code must have at least 4 tokens, got: %q", genCode)
	}

	defIndex64, err := strconv.ParseUint(tokens[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid defindex: %w", err)
	}
	paintIndex64, err := strconv.ParseUint(tokens[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid paintindex: %w", err)
	}
	paintSeed64, err := strconv.ParseUint(tokens[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid paintseed: %w", err)
	}
	paintWear64, err := strconv.ParseFloat(tokens[3], 32)
	if err != nil {
		return nil, fmt.Errorf("invalid paintwear: %w", err)
	}

	rest := tokens[4:]
	stickers := []Sticker{}
	keychains := []Sticker{}

	if len(rest) >= 10 {
		stickerTokens := rest[:10]
		for slot := 0; slot < 5; slot++ {
			sid, _ := strconv.ParseUint(stickerTokens[slot*2], 10, 32)
			wear64, _ := strconv.ParseFloat(stickerTokens[slot*2+1], 32)
			wear := float32(wear64)
			if sid != 0 {
				s := Sticker{Slot: uint32(slot), StickerID: uint32(sid), Wear: &wear}
				stickers = append(stickers, s)
			}
		}
		rest = rest[10:]
	}

	for i := 0; i+1 < len(rest); i += 2 {
		sid, _ := strconv.ParseUint(rest[i], 10, 32)
		wear64, _ := strconv.ParseFloat(rest[i+1], 32)
		wear := float32(wear64)
		if sid != 0 {
			kc := Sticker{Slot: uint32(i / 2), StickerID: uint32(sid), Wear: &wear}
			keychains = append(keychains, kc)
		}
	}

	pw := float32(paintWear64)
	_ = math.MaxFloat32 // ensure math is used

	return &ItemPreviewData{
		DefIndex:   uint32(defIndex64),
		PaintIndex: uint32(paintIndex64),
		PaintSeed:  uint32(paintSeed64),
		PaintWear:  &pw,
		Stickers:   stickers,
		Keychains:  keychains,
	}, nil
}
