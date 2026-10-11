package models

// Gerbang kelayakan endorsement - `Activity/SetErrorBatalEndorsement_Act.xml`.
//
// ⚠️ Nama rule Pega berbohong (spec §3): ia gerbang kelayakan PEMBUATAN
// endorsement, bukan pembatalan. Kelima pemeriksaan berjalan berurutan dan
// MENGUMPULKAN pesan (`Property-Set-Messages` 3.1-3.10, nol keluar dini);
// 3.11 b2359 menyalakan `Protect.CARI1 = 1` bila ada satu pesan pun, dan
// tombol `Submit` b4226 hanya tampil bila `Protect.CARI1 = '0'`.

// Pesan VERBATIM langkah 1 b285.
const (
	// PesanSudahBatal - b312 (`Local.errmsg`), gerbang 4.
	PesanSudahBatal = "Sudah Di endorsement Batal"
	// PesanEDMBelumSelesai - b359 (`Local.ErrMsgCreate`), gerbang 3.
	PesanEDMBelumSelesai = "There's EDM with this policy no that haven't finish yet!"
	// PesanPolisTidakSah - b380 (`Local.ErrMsgPol`), gerbang 1.
	PesanPolisTidakSah = "Invalid Policy No !"
	// PesanPolisTidakAda - b401 (`Local.ErrMsgNull`), gerbang 2.
	PesanPolisTidakAda = "Policy no Not Found !"
	// PesanSudahDibayar - b422 (`Local.ErrMsgPembayaran`), gerbang 5.
	PesanSudahDibayar = "There's already payment with this policy no"
)

// FaktaGerbang adalah jawaban basis data atas kelima pemeriksaan.
type FaktaGerbang struct {
	// PolisAda - 3.2/3.3 `GetPL_NumberLife` (versi polis di salah satu sumber).
	PolisAda bool
	// AdaKasusTerbuka - 3.4-3.6 RD `FilterProteksiEDMLife` (OQ-EDM-002).
	AdaKasusTerbuka bool
	// SudahBatal - 3.7/3.8 `GetEdmTypeLife` memuat `3`.
	SudahBatal bool
	// SudahDibayar - 3.9 `SearcStatusBayarArasaps_SQL` (hanya dibaca bila
	// `EdmType=3`, prakondisi 3.10 b2319).
	SudahDibayar bool
}

// PesanGerbang menyusun pesan kelima gerbang, urut langkah korpus.
//
// ⚠️ Nomor polis kosong memberi DUA pesan, seperti Pega: 3.1 b801
// (`param.Nopolis==""`) lalu 3.3 b1125 - `GetPL_NumberLife` atas nomor kosong
// tidak menemukan apa pun. Pemeriksaan basis data lain tidak berjalan (jawabannya
// pasti kosong).
func PesanGerbang(nomorPolis, edmType string, f FaktaGerbang) []string {
	var pesan []string
	if nomorPolis == "" {
		return []string{PesanPolisTidakSah, PesanPolisTidakAda}
	}
	if !f.PolisAda {
		pesan = append(pesan, PesanPolisTidakAda)
	}
	if f.AdaKasusTerbuka {
		pesan = append(pesan, PesanEDMBelumSelesai)
	}
	if f.SudahBatal {
		pesan = append(pesan, PesanSudahBatal)
	}
	if f.SudahDibayar && edmType == EdmTypeBatal {
		pesan = append(pesan, PesanSudahDibayar)
	}
	return pesan
}

// Kelayakan adalah jawaban gerbang untuk layar `EndorsmentLife_Section`.
// GerbangPembayaranAktif - saklar gerbang 5 (pembayaran Arasapas untuk EDM Type 3 / Batal,
// `SetErrorBatalEndorsement_Act` 3.10). DINONAKTIFKAN SEMENTARA (keputusan work owner 10-10-2026: "untuk
// sementara pengecekan di hapus dulu, tapi nanti akan digunakan lagi") sampai izin baca
// `ARASAPAS.DETAIL_INVOICE` (OQ-EDM-012) tersedia - selama mati, skema Arasapas sama sekali tidak dibaca.
// Ubah ke `true` untuk menyalakan kembali. `var` (bukan `const`) hanya supaya uji dapat menyalakannya.
var GerbangPembayaranAktif = false

type Kelayakan struct {
	// Pesan - kosong = boleh dibuat.
	Pesan []string `json:"pesan"`
	// Boleh - `Protect.CARI1 = '0'` **dan** `EdmType` `1`/`3` (syarat tombol
	// `Submit` b4226).
	Boleh bool `json:"boleh"`
	// Policy No bernomor endorsement (keputusan work owner 10-10-2026): NomorPolis = hasil
	// `NomorPolisDariMasukan` yang dipakai semua gerbang; AdaAkhiran/AkhiranMasukan = akhiran yang diketik;
	// PolisAda = versi polis ditemukan; NomorVersiTerakhir/AkhiranVersiTerakhir = versi PRODKE terbesar
	// (`VersiBerjalan`, `<polis>/%02d`) - kosong bila belum pernah di-endorse (versi new business).
	NomorPolis           string `json:"nomorPolis"`
	AdaAkhiran           bool   `json:"adaAkhiran"`
	AkhiranMasukan       string `json:"akhiranMasukan"`
	PolisAda             bool   `json:"polisAda"`
	NomorVersiTerakhir   string `json:"nomorVersiTerakhir"`
	AkhiranVersiTerakhir string `json:"akhiranVersiTerakhir"`
}

// PesanCSVPlan - `SaveCSVEDMLife` 4 b800, VERBATIM (berbahasa Indonesia di korpus). ⚠️ TIDAK LAGI dipakai unggahan
// sejak 10-10-2026 (keputusan work owner): PLAN / POLICY_HOLDER kini diambil dari polis, nilai CSV diabaikan
// (`RapikanBarisCSV`). Konstanta dipertahankan sebagai bukti korpus.
const PesanCSVPlan = "Plan di CSV tidak sesuai, mohon di cek kembali"

// KodeJurnalPelunasan - `IVD_JR_ID` baris pelunasan Arasapas
// (`RDBList/SearcStatusBayarArasaps_SQL.xml` b58 `IVD_JR_ID ='5'`; spec §13
// `[keputusan work owner]`: pembayaran/pelunasan). Kode bernilai satu digit
// tinggal di models (penjaga kode status Claim Life).
const KodeJurnalPelunasan = "5"
