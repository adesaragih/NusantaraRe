package models

// Untuk apa berkas ini: BENTUK GENERASI ENDORSEMEN yang disimpan - aturan keutuhan antar generasi dan penyesuaian
// bentuk terhadap penjaga tabel dasar NB.

import "fmt"

// Proporsional - generasi berjenis Proportional (`QuotationData.ProportionalType`, cadangan Quotation).
func Proporsional(h *Halaman) bool {
	jenis := h.Ambil(HalamanPolis + ".QuotationData.ProportionalType")
	if jenis == "" {
		jenis = h.Ambil(HalamanQuotation + ".ProportionalType")
	}
	return jenis == JenisProporsional
}

// HitungSelisihGenerasi - selisih generasi PROPORSIONAL dihitung ulang dari data baru dan OldData SETIAP generasi
// ditulis (tinjauan kode 06-10-2026). Di Pega hasil tombol "Calculate Value Difference" menempel di pyWorkPage dan
// ikut tersimpan - tanpa tombol, selisih tersimpan kosong / basi dan produksi (selisih) salah. Rancangan proyeksi
// menuntut tabelnya selalu = hitung ulang (diagram EDM Prop J98 "bila beda dari hitung ulang, TABELNYA yang salah";
// spec-penyimpanan ID-26). Rumus: `EDMTCalculateTreatyDifference` (satu rumus, ID-28/30). NonProp: selisih lapisan
// XOL dihitung `CalculateDifferenceEDM` di rantai Choose Business (membaca master) - tidak di sini.
func HitungSelisihGenerasi(h *Halaman) error {
	if !Proporsional(h) {
		return nil
	}
	return EDMTCalculateTreatyDifference(h)
}

// RapikanBentukSimpan - generasi PROPORSIONAL tidak membawa baris XOL maupun rincian angsuran bersarang
// (ketetapan tabel dasar NB `PeriksaBentukSimpan`, spec-penyimpanan ID-26 / ID-35, AC 31 / 33).
//
// ⚠️ Di XML rantai `EDMChooseBusiness_Act` berjalan SAMA untuk proporsional: `InputPolicyTreatyEDMDetail_NP`
// (SetValueEDM_Act 9) mengisi `TreatyXOLList` dari master, `CalculateDifferenceEDM_act` (6) mengisi
// `TreatyXOLDifferenceList`, `FillSpreading` (9) dan `FillPaymentInstallmentEDMT` (12) membacanya. Sesudah rantai
// itu kedua daftar XOL tidak dibaca rule proporsional mana pun (tab Old / New / Value Difference dan produksi blok 9
// memakai TreatyDifference; AddPremi hanya `.IsNewPolicyNonProp = 1`) - di generasi proporsional keduanya HASIL
// ANTARA dan dibuang sebelum disimpan, beserta rincian `ListInstallment(n).InstallmentList` (grid angsuran tab New
// Data satu tingkat, `PropNewData2`).
func RapikanBentukSimpan(h *Halaman) {
	if !Proporsional(h) {
		return
	}
	for i := range h.AmbilDaftar(TabelXOL.Daftar) {
		h.SetelDaftar(JalurAnak(TabelXOL.Daftar, i+1, TabelLayerXOL.Daftar), nil)
	}
	h.SetelDaftar(TabelXOL.Daftar, nil)
	for i := range h.AmbilDaftar(DaftarSelisihXOL) {
		h.SetelDaftar(JalurAnak(DaftarSelisihXOL, i+1, "ValueList"), nil)
	}
	h.SetelDaftar(DaftarSelisihXOL, nil)
	for i := range h.AmbilDaftar(DaftarAngsuran) {
		h.SetelDaftar(JalurAnak(DaftarAngsuran, i+1, TabelAngsuranRinci.Daftar), nil)
	}
}

// PesanShareSpreadingKosong - ⛔ BUKAN teks XML: baris spreading tanpa % Share (keputusan work owner 07-10-2026,
// rekomendasi b). `n` berbasis 1.
func PesanShareSpreadingKosong(n int) string {
	return fmt.Sprintf("Spreading row %d: %% Share is required", n)
}

// PesanBarisSpreadingHilang - ⛔ BUKAN teks XML: pesan aturan keutuhan generasi (`[keputusan work owner]`
// spec-penyimpanan ID-15, AC 8). `n` berbasis 1.
func PesanBarisSpreadingHilang(n int) string {
	return fmt.Sprintf("Spreading row %d of the previous policy generation is missing", n)
}

// BarisSpreadingHilang = aturan keutuhan ID-15 (AC 8): generasi n+1 wajib memuat setiap baris (NOURUT = posisi)
// spreading generasi n. Di EDM baris tidak dapat dihapus (ID-16); baris yang tidak ada di belakang daftar baru
// berarti endorsemen itu diam-diam menghapus komitmen yang masih berlaku.
//
// ⚠️ Hanya SpreadingRiskList: `FillSpreading` (EDMChooseBusiness_Act 9) menyusun data baru dari BARIS PERTAMA
// generasi lama saja, sehingga baris lain harus ditambahkan admin (Add, NOURUT maks + 1) sebelum Submit. Rincian
// angsuran tidak diperiksa - jumlah termin dapat diubah di layar (`.Installment` -> FillPaymentInstallment);
// pertentangan AC 8 lawan XML itu dicatat sebagai pertanyaan. Polis NonProp baru (`.IsNewPolicyNonProp = 1`)
// tidak diperiksa: grid spreading-nya (AddPremi S17) hanya-baca tanpa Add, jadi baris tidak dapat ditambahkan.
func BarisSpreadingHilang(h *Halaman) []string {
	if PolisNonPropBaru(h) {
		return nil
	}
	lama := len(h.AmbilDaftar(od + "SpreadingRiskList"))
	baru := len(h.AmbilDaftar(DaftarSpreading))
	var out []string
	for n := baru + 1; n <= lama; n++ {
		out = append(out, PesanBarisSpreadingHilang(n))
	}
	return out
}
