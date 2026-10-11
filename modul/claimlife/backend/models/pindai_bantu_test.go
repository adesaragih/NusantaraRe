package models_test

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

// Lingkup penjaga milik Claim Life (keputusan work owner 10-10-2026 "mulai kerja Claim Life").
//
// ⛔ Sejak satu folder per modul, akar pindaian di atas = SELURUH APP_RNM, sehingga aturan internal Claim Life
// (literal kode status, satu rumah `Type`, kolom tak dibawa, total peserta) ikut menuduh kode modul lain yang punya
// kosakata sendiri - `inti/backend/master`, bordereaux, nbfacin. Lingkup barunya: SELURUH berkas `modul/claimlife/`,
// ditambah berkas modul lain yang KODENYA (komentar dibuang) menyebut tabel milik Claim Life - di sanalah aturan
// Claim Life memang berlaku lintas modul (mis. Komite Claim Life menulis `T_GENERAL_CLAIM`).
//
// Tabel milik Claim Life = yang dibuat migrasi `modul/claimlife/backend/migrations` (001-023) + tabel datar warisan
// `OS_AKSEPTASI_KLAIM_LIFE`, KECUALI tabel Komite (`T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`, migrasi 013): milik konteks
// Komite Claim Life (penjaga `TestNolPenyimpanKeputusanKomiteDiKonteksIni`). `T_LOG_SERVICE_RNM` (015) SENGAJA tidak ikut: antrean layanan bersama yang ditulis banyak
// modul (Master Product Name Life, Endorsement Life, Komite) dengan kosakatanya sendiri.
var polaTabelClaimLife = regexp.MustCompile(`(?i)\b(T_WORK_CLAIM|T_GENERAL_CLAIM|T_CLAIMLF_[A-Z0-9_]+|OS_AKSEPTASI_KLAIM_LIFE)\b`)

// milikClaimLife - berkas di folder modul Claim Life.
func milikClaimLife(jalur string) bool {
	return strings.Contains("/"+filepath.ToSlash(jalur), "/modul/claimlife/")
}

// dalamLingkupClaimLife - berkas Claim Life, atau berkas lain yang kodenya (tanpa komentar) menyebut tabel Claim Life.
func dalamLingkupClaimLife(jalur, kodeTanpaKomentar string) bool {
	return milikClaimLife(jalur) || polaTabelClaimLife.MatchString(kodeTanpaKomentar)
}

// relAkar - jalur relatif APP_RNM apa adanya (untuk pesan; relLapisan membuang nama modul).
func relAkar(jalur string) string {
	rel := filepath.ToSlash(jalur)
	for strings.HasPrefix(rel, "../") {
		rel = strings.TrimPrefix(rel, "../")
	}
	return rel
}
