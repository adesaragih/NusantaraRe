package services_test

// Pembantu pindai seluruh aplikasi - refactor bentuk B (30-09-2026).
//
// Penjaga di paket ini dulu memindai "..", yaitu `internal/` saja. Sejak kode
// bersama pindah ke `inti/` dan modul ke `modul/<nama>/`, akar itu diam-diam
// menyempit. Kini yang dipindai akar aplikasi; folder yang bukan sumber Go
// aplikasi dilewati, dan jalur dinormalkan per lapisan.

import (
	"path/filepath"
	"regexp"
	"strings"
)

// akarAplikasiPindai menunjuk folder APP_RNM dari folder paket ini.
const akarAplikasiPindai = "../../../.." // APP_RNM, dari modul/<nama>/backend/<paket>

// lewatiFolderPindai - folder yang bukan sumber Go aplikasi.
func lewatiFolderPindai(nama string) bool {
	switch nama {
	case "frontend", "node_modules", "bin", ".git", "dist", "unggahan":
		return true
	}
	return false
}

// relLapisan - jalur relatif akar aplikasi tanpa awalan `internal/` atau
// `modul/<nama>/backend/`: "models/x.go" berarti lapisan models modul MANA PUN.
// (`backend/` sejak struktur tim satu folder per modul, 30-09-2026.)
func relLapisan(jalur string) string {
	rel := filepath.ToSlash(jalur)
	for strings.HasPrefix(rel, "../") {
		rel = strings.TrimPrefix(rel, "../")
	}
	rel = strings.TrimPrefix(rel, "internal/")
	if strings.HasPrefix(rel, "modul/") {
		sisa := strings.TrimPrefix(rel, "modul/")
		if i := strings.Index(sisa, "/"); i >= 0 {
			rel = strings.TrimPrefix(sisa[i+1:], "backend/")
		}
	}
	return rel
}

// ---------------------------------------------------------------------------
// Lingkup penjaga milik Claim Life (11-10-2026) - salinan paket services_test dari lingkup_claimlife_test.go (paket
// internal; nama tak berekspor tidak menyeberang). Lihat penjelasan lengkap di sana.
// ---------------------------------------------------------------------------

var polaTabelClaimLifeLuar = regexp.MustCompile(`(?i)\b(T_WORK_CLAIM|T_GENERAL_CLAIM|T_CLAIMLF_[A-Z0-9_]+|OS_AKSEPTASI_KLAIM_LIFE)\b`)

var polaSasaranClaimLifeLuar = regexp.MustCompile(`(?is)^\s*(CREATE\s+(GLOBAL\s+TEMPORARY\s+)?TABLE|ALTER\s+TABLE|DROP\s+TABLE|` +
	`COMMENT\s+ON\s+(TABLE|COLUMN)|INSERT\s+INTO|UPDATE|MERGE\s+INTO)\s+("?[A-Z0-9_{}]+"?\.)?"?` +
	`(T_WORK_CLAIM|T_GENERAL_CLAIM|T_CLAIMLF_[A-Z0-9_]+|OS_AKSEPTASI_KLAIM_LIFE)\b`)

// milikClaimLifeLuar - berkas di folder modul Claim Life.
func milikClaimLifeLuar(jalur string) bool {
	return strings.Contains("/"+filepath.ToSlash(jalur), "/modul/claimlife/")
}

// goLingkupClaimLife - berkas Go Claim Life, atau berkas lain yang KODENYA (tanpa komentar) menyebut tabel Claim Life.
func goLingkupClaimLife(jalur, kodeTanpaKomentar string) bool {
	return milikClaimLifeLuar(jalur) || polaTabelClaimLifeLuar.MatchString(kodeTanpaKomentar)
}

// sqlLingkupClaimLifeLuar - berkas migrasi Claim Life UTUH; modul lain hanya pernyataan bersasaran tabel Claim Life.
func sqlLingkupClaimLifeLuar(jalur, isi string) []string {
	if milikClaimLifeLuar(jalur) {
		return []string{isi}
	}
	var hasil []string
	var kini strings.Builder
	tutup := func() {
		for _, p := range strings.Split(kini.String(), ";") {
			if polaSasaranClaimLifeLuar.MatchString(p) {
				hasil = append(hasil, p)
			}
		}
		kini.Reset()
	}
	for _, baris := range strings.Split(strings.ReplaceAll(isi, "\r\n", "\n"), "\n") {
		t := strings.TrimSpace(baris)
		if strings.HasPrefix(t, "--") {
			continue
		}
		if t == "/" {
			tutup()
			continue
		}
		kini.WriteString(baris)
		kini.WriteString("\n")
	}
	tutup()
	return hasil
}
