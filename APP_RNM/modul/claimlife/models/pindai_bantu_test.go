package models_test

// Pembantu pindai seluruh aplikasi - refactor bentuk B (30-09-2026).
//
// Penjaga di paket ini dulu memindai "..", yaitu `internal/` saja. Sejak kode
// bersama pindah ke `inti/` dan modul ke `modul/<nama>/`, akar itu diam-diam
// menyempit. Kini yang dipindai akar aplikasi; folder yang bukan sumber Go
// aplikasi dilewati, dan jalur dinormalkan per lapisan.

import (
	"path/filepath"
	"strings"
)

// akarAplikasiPindai menunjuk folder APP_RNM dari folder paket ini.
const akarAplikasiPindai = "../../.."

// lewatiFolderPindai - folder yang bukan sumber Go aplikasi.
func lewatiFolderPindai(nama string) bool {
	switch nama {
	case "frontend", "node_modules", "bin", ".git", "dist", "unggahan":
		return true
	}
	return false
}

// relLapisan - jalur relatif akar aplikasi tanpa awalan `internal/` atau
// `modul/<nama>/`: "models/x.go" berarti lapisan models modul MANA PUN.
func relLapisan(jalur string) string {
	rel := filepath.ToSlash(jalur)
	for strings.HasPrefix(rel, "../") {
		rel = strings.TrimPrefix(rel, "../")
	}
	rel = strings.TrimPrefix(rel, "internal/")
	if strings.HasPrefix(rel, "modul/") {
		sisa := strings.TrimPrefix(rel, "modul/")
		if i := strings.Index(sisa, "/"); i >= 0 {
			rel = sisa[i+1:]
		}
	}
	return rel
}
