package models

// Baris layar daftar kontrak — ronde layar 1.
//
// ⛔ Kesembilan kolomnya DIBACA dari `Section/InputTreatyInOffer.xml`, bukan
// dipilih di sini. Tiap kolom membawa jejaknya di komentar: nama rule
// `pyCaption` dan posisi bita di dalam ekspor 2026-09, supaya pembaca
// berikutnya dapat memeriksanya tanpa mempercayai berkas ini.
//
// ⚠️ DUA kolom tidak punya rumah di model baru, dan itu dinyatakan — bukan
// diisi tebakan. Lihat `PosisiKe` dan pasangan nama cedant/asal bisnis.

// BarisDaftarKontrak adalah satu baris layar daftar.
type BarisDaftarKontrak struct {
	// `pyCaption ID` — kolom pertama layar lama.
	ID int64 `json:"id"`
	// `pyCaption Contract Name` @7257527 — di sistem lama `.TreatyContractName`
	// hidup di kepala kontrak; di model baru ia kolom VERSI_KONTRAK, jadi yang
	// ditampilkan nama pada versi TERAKHIR.
	NamaKontrak string `json:"namaKontrak"`
	// `pyCaption Reinsurance Type` @7236037 — `TreatyIn.ProportionType`.
	SifatProporsi string `json:"sifatProporsi"`
	// `pyCaption Source of Business` @6540552 — `TreatyIn.LeadingReinsSource`.
	// ⚠️ Model baru menyimpan PENGENALNYA saja: `ERD.md` §2.8 menyatakan
	// `ASAL_BISNIS` ada DI LUAR skema ini, jadi namanya tidak dapat dibaca
	// tanpa modul yang memilikinya. Pengenal, bukan nama yang dikarang.
	IDAsalBisnis int64 `json:"idAsalBisnis"`
	// `pyCaption Ceding` @6503393 — `TreatyIn.Ceding`, tampil `.ClientName`.
	// ⚠️ Sama: `CEDANT` di luar skema (`ERD.md` §2.8).
	IDCedant int64 `json:"idCedant"`
	// `pyCaption Commencement` @6540017 — `TreatyIn.Commencement`.
	TanggalMulai string `json:"tanggalMulai"`
	// `pyCaption Termination` @6523491 — `TreatyIn.Termination`.
	TanggalBerakhir string `json:"tanggalBerakhir"`
	// `pyCaption Status Accept` @7249608 — `.StatusAkseptasi` di sistem lama.
	// Padanannya di model baru keadaan siklus hidup versi terakhir.
	KeadaanSiklusHidup string `json:"keadaanSiklusHidup"`
	// `pyCaption Position To` @7232843 — `.PositionUsername` di sistem lama.
	//
	// ⛔ NOL kolom di model baru menyimpannya, dan ia SENGAJA dikirim kosong.
	// Tiket 45 (daftar keadaan) dan 49-53 (jalur persetujuan) yang akan
	// menentukan siapa pemegang antrian sebuah versi; sampai itu ada,
	// mengisinya berarti mengarang. Kosong di sini berarti "belum ada yang
	// menghasilkannya", dan layar menyatakannya begitu.
	PosisiKe string `json:"posisiKe"`
}
