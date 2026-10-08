package models

// Tombol `Add Revision` / `Add Adjustment Premium` layar Adjustment —
// `Section/InputTreatyInAdjustment.xml` @500554 / @515429.
//
// ⭐ Rantainya di ekspor:
//
//	Add …    DataTransform `TreatyCreateEDM` (EDMState "1" / "3"), lalu modal
//	         `PickerTreatyInMasterRevisi` / `PickerTreatyInMaster`
//	Choose   Activity `TreatyInEDMSetValue` — muat master, setel nilai EDM,
//	         salin ke `OLDDATA`, nomor revisi baru, SALIN LAMPIRAN, SIMPAN
//
// ⛔ Keputusan pemilik proses 7 Oktober 2026: `Choose` menyusun DRAF DI
// LAYAR. Langkah [8] (salin lampiran) dan [9] (`SaveTreatyIn_EDM_Act`)
// adalah tulisan, dan menunggu jalur Save.

// BarisMasterPilihan - satu baris grid picker (`TempMasterList.pxResults`,
// diisi `TreatyLoadMasterJoinEdm` / `TreatyLoadMasterJoinEdmXOL`).
//
// ⛔ Nilainya APA ADANYA dari kolom `TREATY_IN` / `TREATY_IN_EDM` — nol
// tafsir; tanggal diformat di layar.
type BarisMasterPilihan struct {
	ID              string `json:"id"`              // CARI1
	NamaKontrak     string `json:"namaKontrak"`     // CARI2
	SifatProporsi   string `json:"sifatProporsi"`   // CARI3
	AsalBisnis      string `json:"asalBisnis"`      // CARI4
	Cedant          string `json:"cedant"`          // CARI5
	TanggalMulai    string `json:"tanggalMulai"`    // CARI6
	TanggalBerakhir string `json:"tanggalBerakhir"` // CARI7
}

// MasukanDraf - parameter `TreatyInEDMSetValue` yang tombol `Choose` kirim.
//
//	Revisi   ID=.CARI1  InternalType=TreatyIn.EDMState  MaterialType=TreatyIn.EDMMaterialType
//	Premi    ID=.CARI1  InternalType='3'                MaterialType='1'
type MasukanDraf struct {
	ID           string `json:"id"`
	InternalType string `json:"internalType"`
	MaterialType string `json:"materialType"`
}
