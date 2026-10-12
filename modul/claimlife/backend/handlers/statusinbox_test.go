package handlers

// Kolom Work Status kotak masuk (10-10-2026): case baru (STS_REJECT header NULL) = Outstanding, kata Inggris.

import (
	"testing"

	"nusantarare/inti/backend/kontrak"
)

func TestKataStatusInbox(t *testing.T) {
	for kode, mau := range map[string]string{
		"": "Outstanding", " ": "Outstanding", kontrak.KodeOutstanding: "Outstanding",
		kontrak.KodeAksep: "Accepted", kontrak.KodeDitolak: "Rejected", "9": "Unknown",
	} {
		if got := kataStatusInbox(kode); got != mau {
			t.Errorf("%q: %q, mau %q", kode, got, mau)
		}
	}
}
