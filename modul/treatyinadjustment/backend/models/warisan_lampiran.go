package models

// Lampiran kontrak - panel Attachment, dari `POOLDATA.M_ATTACHMENTTREATY_2`.
//
// ⛔ Tabel WARISAN yang hidup, 43 baris. Nol tabel baru, nol migrasi di
// modul ini.
//
// ⚠️ Bentuknya SAMA dengan `modul/treatyin/backend/models`, dan itu
// disengaja: keduanya memotret tabel yang sama. Yang dilarang adalah
// mengimpor model modul sebelah - `TestModulTidakMengimporModulLain`
// menolaknya, dan alasannya berdiri sendiri. Nama tabel dan nama kolomnya
// milik sistem lama; bentuk tampilnya boleh menyimpang per modul.
type BarisLampiranWarisan struct {
	ID           string `json:"id"`
	KodeKategori string `json:"kodeKategori"`
	NamaKategori string `json:"namaKategori"`
	NamaBerkas   string `json:"namaBerkas"`
	JenisMime    string `json:"jenisMime"`
	IDSimpanan   string `json:"idSimpanan"`
	Diunggah     string `json:"diunggah"`
	Pengunggah   string `json:"pengunggah"`
}

// BarisKategoriLampiran - satu baris panel Attachment: Category + Count.
//
// ⚠️ SELURUH kategori tampil, termasuk yang nol berkas: kolom `Count` tidak
// akan pernah berbunyi `0` kalau barisnya disembunyikan saat kosong.
type BarisKategoriLampiran struct {
	Kode  string `json:"kode"`
	Nama  string `json:"nama"`
	Cacah int    `json:"cacah"`
	// Dipastikan - kode ini bernama di katalog (`M_KATEGORIMASTERTREATY`,
	// lalu data lampiran). Sejak 8 Oktober 2026 kesebelas kode bernama;
	// `false` hanya untuk kode di data yang tidak ada di katalog mana pun.
	Dipastikan bool `json:"dipastikan"`
}

// BarisRiwayatWarisan - satu baris panel History.
//
// Keempat medannya berpadanan satu-ke-satu dengan kolom layar lama:
// Date · PIC · Approval · Comment, dari `T_VIEW_COMMENT`.
type BarisRiwayatWarisan struct {
	Tanggal   string `json:"tanggal"`
	Operator  string `json:"operator"`
	Disetujui string `json:"disetujui"`
	Catatan   string `json:"catatan"`
}

// BarisPolisMaster - satu baris panel `Existing Policy for Master ID` layar
// Adjustment: `Policy No` (`NOPOLIS`) · `Pega ID` (`IDPEGA` tanpa 18 aksara
// awalnya), dari `TREATYINPRODUCTION`.
type BarisPolisMaster struct {
	NomorPolis string `json:"nomorPolis"`
	PegaID     string `json:"pegaID"`
}
