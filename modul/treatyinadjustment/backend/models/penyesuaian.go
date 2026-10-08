package models

// Penyesuaian WARISAN — layar `InputTreatyInAdjustment` sistem lama.
//
// ⛔ SUMBERNYA `POOLDATA.TREATY_IN_EDM` + `POOLDATA.M_TREATY_IN_EDM`, dan
// keduanya DIUKUR 5 Oktober 2026 sebelum dipilih:
//
//	TREATY_IN_EDM     280 baris  ID `1000080/R01` …, OLDID menunjuk asalnya
//	M_TREATY_IN_EDM   280 baris  JSONDATA = halaman `TreatyIn` utuh
//	                             280 dari 280 membawa kunci `OLDDATA`
//
// ⚠️ BUKAN `KONTRAK`/`VERSI_KONTRAK` (migrasi 401/440): keduanya NOL baris,
// jadi jalur versi-dasar migrasi 440 belum punya satu pun data untuk
// ditampilkan. Penyesuaian yang hidup ada di tabel warisan di atas.
//
// ⛔ Kedua panel Old dan New membaca HALAMAN YANG SAMA, dan itu yang ekspor
// katakan: `Section/TreatyInNONProportionalOldData.xml` mengikat
// `TreatyIn.OLDDATA.*` (mis. @34625 `TreatyIn.OLDDATA.TreatyContractName`),
// `Section/TreatyInNONProportional.xml` mengikat `TreatyIn.*`. Satu dokumen,
// dua halaman di dalamnya — bukan dua baris yang digabung.

// BarisPenyesuaian - satu baris grid daftar (mode `DATASHOW != 1`,
// `Section/InputTreatyInAdjustment.xml` @481804).
//
// ⛔ Nilainya APA ADANYA dari kolom `TREATY_IN_EDM`. `JenisPenyesuaian` dan
// `JenisMaterial` adalah KODE (`1`/`2`/`3`, `1`/`2`); teks pilihannya
// berasal dari daftar `associated` milik rule properti yang TIDAK ikut
// diekspor, jadi ia tidak dikarang di sini.
type BarisPenyesuaian struct {
	ID               string `json:"id"`
	IDAsal           string `json:"idAsal"`
	JenisPenyesuaian string `json:"jenisPenyesuaian"`
	JenisMaterial    string `json:"jenisMaterial"`
	NamaKontrak      string `json:"namaKontrak"`
	SifatProporsi    string `json:"sifatProporsi"`
	AsalBisnis       string `json:"asalBisnis"`
	Cedant           string `json:"cedant"`
	TanggalMulai     string `json:"tanggalMulai"`
	TanggalBerakhir  string `json:"tanggalBerakhir"`
	Posisi           string `json:"posisi"`
	StatusAkseptasi  string `json:"statusAkseptasi"`
}

// SisiPenyesuaian - satu halaman di dalam dokumen: `TreatyIn` (New) atau
// `TreatyIn.OLDDATA` (Old).
//
// ⛔ `Medan` HANYA memuat kunci yang ADA di dokumen. Kunci yang tidak ada
// tidak dimasukkan dengan nilai kosong — "tidak ada di sistem lama" dan
// "kosong" berarti hal yang berbeda di layar, dan peta yang mengisi
// keduanya dengan `""` menghapus bedanya tanpa suara.
//
// ⛔ `Larik` memuat larik yang ADA di dokumen; larik yang tidak ada tidak
// dimasukkan sama sekali, larik kosong masuk sebagai larik kosong.
type SisiPenyesuaian struct {
	Medan map[string]string              `json:"medan"`
	Larik map[string][]map[string]string `json:"larik"`
	// Pohon - larik AKAR yang punya larik anak (Limits, Share,
	// FacultativeShareList, Installment) sebagai simpul BERSARANG: nilai teks
	// atau larik simpul. Dibaca rumus tombol; grid membaca `Larik`.
	Pohon map[string][]map[string]any `json:"pohon"`
}

// Penyesuaian - satu penyesuaian utuh: pengenalnya, sisi New, dan sisi Old.
//
// ⚠️ Medan kepala layar (`ID Revision`, `Reinsurance Type`, `Adjustment
// Type`, …) dibaca dari `Baru.Medan`, BUKAN dari kolom `TREATY_IN_EDM`:
// layar lama mengikat `TreatyIn.*`, dan dua sumber untuk satu medan adalah
// dua kebenaran yang suatu hari berselisih. `ID` dan `IDAsal` di sini hanya
// pengenal baris yang dibuka.
type Penyesuaian struct {
	ID     string          `json:"id"`
	IDAsal string          `json:"idAsal"`
	Baru   SisiPenyesuaian `json:"baru"`
	Lama   SisiPenyesuaian `json:"lama"`
}
