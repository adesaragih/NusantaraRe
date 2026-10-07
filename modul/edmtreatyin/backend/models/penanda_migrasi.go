package models

// Untuk apa berkas ini: DUA PENANDA MIGRASI - tiket EDM 09 (spec-penyimpanan-relasional.md ID-33..ID-37, AC 40-43).
// Fungsi MURNI: menandai baris proyeksi selisih hasil migrasi (SUMBER 'PEGA') TANPA mengubah satu angka pun -
// halaman hanya dibaca.
//
//	PASANGAN_BERGESER  (ID-35, AC 41) = 1 bila kunci dagang baris ke-n generasi baru berbeda dari baris ke-n
//	                   generasi lama. Pasangan baris antar generasi = NOURUT (posisi, ID-34 / AC 40), sama dengan
//	                   Pega `.pxListSubscript` (`EDMTCalculateTreatyDifference` 2.1 / 3.1, `CalculateDifferenceEDM_act`
//	                   subskrip induk c dan lapisan l); kunci dagang turun pangkat menjadi PEMERIKSA:
//	                     spreading  TreatyType                      (T_POLIS_DIFFERENCE_SPREADING, 361)
//	                     angsuran   InstallmentNo                   (T_POLIS_DIFFERENCE_INSTALMENT, 362)
//	                     lapisan    Layer + LayerPart + IDCurrency  (T_POLIS_XOL_LAYER_DIFFERENCE, 363)
//	                   Baris tanpa pasangan (baris tambahan, NOURUT = maks + 1, ID-17) TIDAK bergeser.
//	RUMUS_BERLAPIS     (ID-36, AC 42) = 1 pada SETIAP baris selisih spreading / angsuran bila generasi sebelumnya
//	                   ber-EDMNo terisi: `EDMTCalculateTreatyDifference` langkah 1 / 4 memilih varian lewat
//	                   `.OldData.EDMNo==""` (salah -> lompat `HasEDMNo`, langkah 4-6 `baru - OldData.TreatyDifference`).
//	                   Kolomnya hanya ada di 361 dan 362 (363 tidak; varian XOL `CalculateDifferenceEDM_act` tanpa
//	                   pemilih EDMNo); induk T_POLIS_DIFFERENCE (360) tanpa kolom penanda - cakupan induk butir WO
//	                   (koreksi 06-10 tiket 09 butir 3), tidak diputuskan di sini.
//
// ⛔ ID-37 / AC 43: hanya untuk baris SUMBER 'PEGA' - penulisnya (`repository.SetelPenandaMigrasi`) menyaring SUMBER
// di SQL; baris 'GO' tidak pernah ditandai.

import "strings"

// PenandaBaris - penanda satu baris selisih (NOURUT = indeks + 1).
type PenandaBaris struct {
	PasanganBergeser bool
	RumusBerlapis    bool
}

// PenandaMigrasi - penanda seluruh baris selisih satu generasi, urut NOURUT tabel proyeksinya.
type PenandaMigrasi struct {
	Spreading []PenandaBaris
	Angsuran  []PenandaBaris
	// Lapisan - urut daftar datar (mata uang c, lapisan l) = `DatarSelisihLapisan`; RumusBerlapis selalu salah.
	Lapisan []PenandaBaris
}

// TeksPenanda - nilai kolom penanda VARCHAR2(1): "1" / "0" (baris PEGA selalu TERISI - AC 42).
func TeksPenanda(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// Cacah - jumlah baris bergeser per tabel dan jumlah baris berlapis.
func (p PenandaMigrasi) Cacah() (spreading, angsuran, lapisan, berlapis int) {
	for _, b := range p.Spreading {
		spreading += pmSatu(b.PasanganBergeser)
		berlapis += pmSatu(b.RumusBerlapis)
	}
	for _, b := range p.Angsuran {
		angsuran += pmSatu(b.PasanganBergeser)
		berlapis += pmSatu(b.RumusBerlapis)
	}
	for _, b := range p.Lapisan {
		lapisan += pmSatu(b.PasanganBergeser)
	}
	return
}

func pmSatu(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Kunci dagang pemeriksa pasangan (kolom katalog, supaya bentuk pembandingnya mengikuti golongan kolom).
var (
	pmKunciSpreading = []Kolom{kKode("TreatyType", "TREATY_TYPE", 64)}
	pmKunciAngsuran  = []Kolom{kCacah("InstallmentNo", "INSTALLMENT_NO")}
	pmKunciLapisan   = []Kolom{kKode("Layer", "LAYER", 64), kKode("LayerPart", "LAYER_PART", 64),
		kKode("IDCurrency", "ID_CURRENCY", 64)}
)

// pmBentukKunci - bentuk pembanding satu nilai kunci: kode / teks APA ADANYA (nol di depan bermakna, ID-47);
// cacah tanpa nol di depan (NUMBER(10) membaca "01" sebagai 1 - dokumen dan tabel sama-sama sah).
func pmBentukKunci(k Kolom, v string) string {
	if k.Golongan != GolCacah {
		return v
	}
	s := strings.TrimSpace(v)
	if s == "" || strings.Trim(s, "0123456789") != "" {
		return s
	}
	if s = strings.TrimLeft(s, "0"); s == "" {
		return "0"
	}
	return s
}

// pmBergeser - baris ke-i kedua daftar ADA dan kunci dagangnya berbeda.
func pmBergeser(baru, lama []Baris, i int, kunci []Kolom) bool {
	if i >= len(baru) || i >= len(lama) || baru[i] == nil || lama[i] == nil {
		return false
	}
	for _, k := range kunci {
		if pmBentukKunci(k, baru[i][k.Properti]) != pmBentukKunci(k, lama[i][k.Properti]) {
			return true
		}
	}
	return false
}

// HitungPenandaMigrasi menghitung kedua penanda untuk setiap baris selisih generasi `baru` (halaman dokumen yang
// akan disimpan: daftar data baru + `TreatyDifference.*` + `TreatyXOLDifferenceList`) terhadap generasi `lama`
// (baris OLD_POLIS_ID: `PolicyTreatyIn.*`, `PolicyTreatyIn.EDMNo` = NOENDORS). Kedua halaman TIDAK diubah.
func HitungPenandaMigrasi(baru, lama *Halaman) PenandaMigrasi {
	var p PenandaMigrasi
	berlapis := strings.TrimSpace(lama.Ambil(pmJalurEDMNo)) != ""
	for i := range baru.AmbilDaftar(DaftarSelisihSpreading) {
		p.Spreading = append(p.Spreading, PenandaBaris{
			PasanganBergeser: pmBergeser(baru.AmbilDaftar(DaftarSpreading), lama.AmbilDaftar(DaftarSpreading), i, pmKunciSpreading),
			RumusBerlapis:    berlapis,
		})
	}
	for i := range baru.AmbilDaftar(DaftarSelisihAngsuran) {
		p.Angsuran = append(p.Angsuran, PenandaBaris{
			PasanganBergeser: pmBergeser(baru.AmbilDaftar(DaftarAngsuran), lama.AmbilDaftar(DaftarAngsuran), i, pmKunciAngsuran),
			RumusBerlapis:    berlapis,
		})
	}
	// Lapisan: subskrip (c, l) daftar selisih = subskrip data baru (CalculateDifferenceEDM_act 1 / 1.2); generasi
	// lama pada subskrip yang sama (`OldData.TreatyXOLList(c).ValueList(l)`).
	for c := range baru.AmbilDaftar(DaftarSelisihXOL) {
		jalurBaru := JalurAnak(TabelXOL.Daftar, c+1, TabelLayerXOL.Daftar)
		for l := range baru.AmbilDaftar(JalurAnak(DaftarSelisihXOL, c+1, TabelSelisihLapisan.Daftar)) {
			p.Lapisan = append(p.Lapisan, PenandaBaris{
				PasanganBergeser: pmBergeser(baru.AmbilDaftar(jalurBaru), lama.AmbilDaftar(jalurBaru), l, pmKunciLapisan),
			})
		}
	}
	return p
}
