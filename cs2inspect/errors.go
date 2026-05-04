package cs2inspect

import "errors"

// ErrMalformedInspectLink is returned (wrapped) by Deserialize when the input
// cannot be a valid inspect-link payload — odd-length hex, non-hex characters,
// payload too short or too long, or proto bytes that fail to parse cleanly.
//
// Callers can use errors.Is(err, ErrMalformedInspectLink) to detect this
// category of failure.
var ErrMalformedInspectLink = errors.New("malformed inspect URL")
