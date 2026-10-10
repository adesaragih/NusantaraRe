package services

// Pemeriksaan "Agent Negative List" — `Activity/TreatyInCheckCedingBlacklist`
// (korpus Treaty In dan Adjustment IDENTIK, isi langkahnya sama):
//
//	[1] Param.ErrMsg = "This name is on Agent Negative List"; TreatyIn.pyErrMsg = ""
//	[2] bila TreatyIn.CedingStatusActive == "inactive":
//	      Property-Set-Messages  Field = TreatyIn.Ceding              Message = Param.ErrMsg
//	[3] bila TreatyIn.SourceStatusActive == "inactive":
//	      Property-Set-Messages  Field = TreatyIn.LeadingReinsSource  Message = Param.ErrMsg
//	[4] bila salah satunya "inactive": TreatyIn.pyErrMsg = "error"
//
// ⭐ KAPAN PEGA MENJALANKANNYA (dibaca dari ekspor, bukan diduga):
//
//   - tombol `Edit` layar daftar (`Section/InputTreatyInOffer.xml`): SetValue
//     → Refresh `SetTreatyIn_Act` → Refresh `TreatyInCheckCedingBlacklist` →
//     `GetMasterTreatyCategory_Act`. `View`, `Copy`, `Revision` TIDAK
//     memanggilnya (`docs/lampiran/INVENTARIS-TOMBOL-EKSPOR.txt` baris 102-130).
//   - event `change` autocomplete Ceding / Business Source
//     (`Section/TreatyInNONProportional.xml` sel 16 / 17).
//   - ⛔ BUKAN tombol `Choose` jendela `TreatyInSearchReinsured` /
//     `TreatyInSearchSoB`: aksinya hanya Refresh `otherSection` lewat
//     DataTransform `TreatyInSetReinsured`, tanpa activity ini.
//
// ⭐ SIFATNYA: pesan MEDAN (Property-Set-Messages pada `.Ceding` /
// `.LeadingReinsSource`), TIDAK menghalangi. `TreatyIn.pyErrMsg = "error"`
// (langkah 4) tidak dibaca aturan mana pun di kedua korpus — hanya ditulis
// (`"empty"` oleh tombol Edit, `""`/`"error"` oleh activity ini).
//
// ⛔ TEMUAN PENTING — di Pega pesan ini TIDAK PERNAH MENYALA dengan data hari ini:
//
//  1. Status yang dibandingkan disalin autocomplete dari `.StatusActive`
//     `BrowseAgentNusaRe_RD`, dan RD itu sendiri menyaring
//     `.StatusActive = "1"` (saringan A, literal). Agen yang dapat dipilih
//     selalu berstatus "1".
//  2. Kolom `AGENT.STATUSACTIVE` terukur 8 Oktober 2026 hanya berisi `'1'`
//     (378) dan `'0'` (51) — nol baris `'inactive'`.
//
// Jadi perbandingan `== "inactive"` dipindahkan APA ADANYA
// (`StatusAgenDaftarNegatif`). Menafsirkan `'0'` sebagai "negative list"
// adalah keputusan pemilik proses, bukan bunyi ekspor — bila diputuskan,
// cukup ganti konstanta itu.
//
// ⚠️ PENYIMPANGAN SUMBER (tak terhindarkan): di Pega, saat Edit, nilai
// `CedingStatusActive` berasal dari dokumen `JSONDATA` kontrak (salinan saat
// agen dipilih). Dokumen itu DILARANG dibaca di aplikasi ini; status dibaca
// dari kolom asalnya, `AGENT.STATUSACTIVE`, untuk pengenal yang tersimpan di
// `TREATY_IN`. Dengan perbandingan literal di atas hasilnya sama: nol pesan.

import (
	"context"
	"strings"

	inti "nusantarare/inti/backend"
)

// PesanDaftarNegatifAgen - `TreatyInCheckCedingBlacklist` langkah 1, apa adanya.
const PesanDaftarNegatifAgen = "This name is on Agent Negative List"

// StatusAgenDaftarNegatif - nilai pembanding langkah 2-4, apa adanya
// (`TreatyIn.CedingStatusActive=="inactive"`).
const StatusAgenDaftarNegatif = "inactive"

// StatusAgenNegatif - hasil periksa SATU medan (Ceding atau Source of Business).
type StatusAgenNegatif struct {
	ID string `json:"id"`
	// Ditemukan - pengenalnya ada di `AGENT`.
	Ditemukan bool `json:"ditemukan"`
	// StatusAktif - `AGENT.STATUSACTIVE` apa adanya.
	StatusAktif   string `json:"statusAktif"`
	DaftarNegatif bool   `json:"daftarNegatif"`
	// Pesan - terisi HANYA bila DaftarNegatif; ditampilkan di bawah medannya.
	Pesan string `json:"pesan,omitempty"`
}

// HasilDaftarNegatifAgen - kedua medan sekaligus, seperti activity-nya.
type HasilDaftarNegatifAgen struct {
	Cedant     StatusAgenNegatif `json:"cedant"`
	AsalBisnis StatusAgenNegatif `json:"asalBisnis"`
}

// PeriksaDaftarNegatifAgen - `TreatyInCheckCedingBlacklist` atas pengenal
// Ceding (`CedingID`) dan Source of Business (`LeadingReinsSourceID`).
// Baca saja; pengenal kosong menghasilkan medan tanpa pesan.
func (l *Layanan) PeriksaDaftarNegatifAgen(ctx context.Context, p inti.Pelaku, idCedant, idAsalBisnis string) (HasilDaftarNegatifAgen, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilDaftarNegatifAgen{}, err
	}
	idCedant, idAsalBisnis = strings.TrimSpace(idCedant), strings.TrimSpace(idAsalBisnis)
	status, err := l.gudang.BacaStatusAktifAgen(ctx, []string{idCedant, idAsalBisnis})
	if err != nil {
		return HasilDaftarNegatifAgen{}, err
	}
	return HasilDaftarNegatifAgen{
		Cedant:     nilaiStatusAgen(idCedant, status),
		AsalBisnis: nilaiStatusAgen(idAsalBisnis, status),
	}, nil
}

// nilaiStatusAgen - langkah 2 / 3 untuk satu medan. Murni.
func nilaiStatusAgen(id string, status map[string]string) StatusAgenNegatif {
	h := StatusAgenNegatif{ID: id}
	if id == "" {
		return h
	}
	h.StatusAktif, h.Ditemukan = status[id]
	// ⛔ `==` Pega atas teks: persis, tanpa memangkas dan tanpa melipat huruf.
	if h.Ditemukan && h.StatusAktif == StatusAgenDaftarNegatif {
		h.DaftarNegatif = true
		h.Pesan = PesanDaftarNegatifAgen
	}
	return h
}
