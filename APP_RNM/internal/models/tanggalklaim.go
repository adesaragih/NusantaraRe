package models

// Tiga tanggal klaim dialog Edit Date - sensus Claim Life §3.1.
//
// Untuk apa berkas ini: gerbang tahap dialog `Edit Date` selain DOL. DOL
// sendiri sudah ditulis tiket 06 (`services/dol.go` Set).
//
// `[terverifikasi]` `Claim Life/Section/EditDateClaimLife_Section.xml`
// (pecahan `sed 's/></>\n</g'`):
//
//	DATE_OF_LOSS         b796   change -> ValidasiDOL_Act b837 (tiket 06)
//	CLAIM_RECEIVED_DATE  b1076  change -> ValidasiClaimReceived_Act b1120
//	COMPLETE_DATE        b1387  caption "DOCUMENT COMPLETE DATE"; postValue saja
//	CONFIRMATION_DATE    b1626  postValue saja
//	Save b1910 -> UpdateDateClaimLife_Act b1929 -> closeContainer -> refresh
//
// Keempat isian berprasyarat baca-saja YANG SAMA (b1000, b1313, b1550, b1788):
//
//	pyWorkPage.pyPosition!='ReasLifeAdmin' || pyWorkPage.ClaimData.PremiumListSummary.CLAIM_NO!=''
//
// ⚠️ `[terbuka - OQ-M1]` separuh kedua, `CLAIM_NO != ''`, TIDAK ditiru. Pega
// menomori klaim saat Save Outstanding (`GetSequenceNumber_SQL` dari
// `SaveOutStandingLife_Act` b8057); aplikasi ini menomorinya saat PENDAFTARAN
// (`services/pendaftaran.go`). Menirunya huruf demi huruf membuat keempat
// isian terkunci pada SETIAP klaim - fitur yang tidak pernah dapat dipakai.
// Maksudnya ("sesudah Save Outstanding, tanggal tidak diubah lagi") menunggu
// padanan Save Outstanding, yang belum punya rute.
//
// ⚠️ `[terbuka - OQ-M9]` `ValidasiClaimReceived_Act` b1120 TIDAK dijalankan
// saat simpan. Aturan DAN ambangnya sudah ada (`services.AmbangKlaim.Hitung`,
// butir ba/bh) tetapi nol pemanggil; penandanya (`.MAXCLAIM_RECEIVED`) tidak
// punya kolom, dan `Hitung` menuntut kedua penanda sekaligus (STNC ikut).
// Di Pega penandanya pun tidak memblokir (OQ-G).
//
// Dibaca sesudah: tahap.go, diagnosa.go.

import "time"

// TanggalKlaim adalah tiga tanggal dialog Edit Date selain DOL.
//
// nil berarti KOSONG - kolomnya ditulis NULL, sama seperti Pega menulis
// isian yang dikosongkan (ADR-U-0027: kosong bukan tanggal nol).
type TanggalKlaim struct {
	TerimaKlaim    *time.Time // CLAIM_RECEIVED_DATE
	DokumenLengkap *time.Time // COMPLETE_DATE
	Konfirmasi     *time.Time // CONFIRMATION_DATE
}

// TahapBolehUbahTanggalKlaim menjawab gerbang tahap dialog Edit Date - MURNI.
//
// Dua syarat, keduanya dari XML: tombol `Edit Date` (b14115) berdiri di grid
// peserta `ClaimLifeDetailGCNM`, jadi tahapnya harus membuka grid itu
// (`TahapBergridPeserta`); dan `pyPosition == 'ReasLifeAdmin'` - tahapnya
// dipegang Admin. Irisannya hanya Outstanding Claim.
func TahapBolehUbahTanggalKlaim(t Tahap) bool {
	peran, ada := PeranPemegangTahap(t)
	return ada && peran == PeranAdminLife && TahapBergridPeserta(t)
}
