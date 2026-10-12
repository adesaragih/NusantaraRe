package services

// Pembantu lingkup penjaga statik paket internal `services` (totalpeserta_statik_test.go) - pasangan dari
// pindai_bantu_test.go paket services_test (paket berbeda, jadi ditulis di sini).

import (
	"path/filepath"
	"regexp"
	"strings"
)

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

// relAkar - jalur relatif APP_RNM apa adanya (untuk pesan; relLapisan membuang nama modul).
func relAkar(jalur string) string {
	rel := filepath.ToSlash(jalur)
	for strings.HasPrefix(rel, "../") {
		rel = strings.TrimPrefix(rel, "../")
	}
	return rel
}

// polaSasaranClaimLife - pernyataan SQL yang SASARANNYA tabel Claim Life (11-10-2026). Menyebut saja tidak cukup:
// migrasi claimprop / claimfacin / claimnonprop menautkan tabel MEREKA SENDIRI ke header klaim bersama lewat
// `FOREIGN KEY (CLAIM_ID) REFERENCES {skema}.T_GENERAL_CLAIM (ID)`; kolom tabel itu keputusan modulnya, bukan skema
// Claim Life. Yang diperiksa: CREATE / ALTER / DROP / COMMENT ON / INSERT / UPDATE / MERGE atas tabel Claim Life.
var polaSasaranClaimLife = regexp.MustCompile(`(?is)^\s*(CREATE\s+(GLOBAL\s+TEMPORARY\s+)?TABLE|ALTER\s+TABLE|DROP\s+TABLE|` +
	`COMMENT\s+ON\s+(TABLE|COLUMN)|INSERT\s+INTO|UPDATE|MERGE\s+INTO)\s+("?[A-Z0-9_{}]+"?\.)?"?` +
	`(T_WORK_CLAIM|T_GENERAL_CLAIM|T_CLAIMLF_[A-Z0-9_]+|OS_AKSEPTASI_KLAIM_LIFE)\b`)

// pernyataanSQL memecah isi migrasi menjadi pernyataan: baris komentar `--` dibuang, pemisahnya `;` ATAU baris `/`
// sendirian (gaya SQL*Plus yang dipakai migrasi modul lain - memecah hanya pada `;` membuat seluruh berkas satu
// pernyataan).
func pernyataanSQL(isi string) []string {
	var hasil []string
	var kini strings.Builder
	tutup := func() {
		for _, p := range strings.Split(kini.String(), ";") {
			if strings.TrimSpace(p) != "" {
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

// sqlLingkupClaimLife - bagian migrasi yang tunduk pada aturan skema Claim Life: berkas Claim Life UTUH; berkas modul
// lain hanya pernyataan yang sasarannya tabel Claim Life (polaSasaranClaimLife).
func sqlLingkupClaimLife(jalur, isi string) []string {
	if milikClaimLife(jalur) {
		return []string{isi}
	}
	var hasil []string
	for _, p := range pernyataanSQL(isi) {
		if polaSasaranClaimLife.MatchString(p) {
			hasil = append(hasil, p)
		}
	}
	return hasil
}
