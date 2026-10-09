package models

// Lampiran berkas Bordereaux - section `AttachmentsBdx` (grid kategori) dan `AttachmentDetailBdx` (popup berkas per
// kategori), tabel warisan `M_ATTACHMENTBORDEREAUX` dan master `M_KATEGORIBORDEREAUX`.

// KategoriLampiran - satu baris grid `AttachmentsBdx` (`GetKategoryDocBDX_SQL`: MASTERID = ID, TYPE = NOTE,
// STATUSAKSEP = jumlah lampiran berkas ini di kategori itu).
type KategoriLampiran struct {
	ID    string `json:"id"`
	Nama  string `json:"nama"`
	Cacah int    `json:"cacah"`
}

// Lampiran - satu baris `M_ATTACHMENTBORDEREAUX` (`AttachDocumentBdx_SQL`).
type Lampiran struct {
	ID         string `json:"id"`
	BdxID      string `json:"-"`
	KategoriID string `json:"kategoriId"`
	// Kategori - `CATEGORY` (kolom Type popup, b2962): nama kategori saat diunggah.
	Kategori string `json:"kategori"`
	FileName string `json:"fileName"`
	// Ekstensi - `FILEMIMETYPE` (`.pyFileMimeType`): ekstensi berkas huruf kecil, mis. `xlsx`.
	Ekstensi string `json:"ekstensi"`
	Username string `json:"username"`
	// StorageID - `T_STORAGE_ID` = `T_STORAGE_IMAGE.IMAGEID`; tidak dikirim ke layar.
	StorageID string `json:"-"`
}
