package models

// Lampiran produk (paket 8, tiket 08–09, PARITAS §6) - baris
// `M_ATTACHMENTPRODUCTNAME` (`InsertAttachProdName_Sql` b84) beserta status
// pengirimannya ke penyimpanan (outbox P5; stub lokal atau penyimpanan nyata - OQ-MPNL-10 dibalik 03-10-2026).

// Status lampiran - terbaca lewat API (AC 36).
const (
	// StatusTerunggah - objeknya tercatat di `T_STORAGE_IMAGE` (`Insert_T_Storage_SQL` b85).
	StatusTerunggah = "terunggah"
	// StatusGagal - percobaan terakhir gagal; galatnya terbaca, dapat diulang.
	StatusGagal = "gagal"
	// StatusBelum - terantre, belum dikirim.
	StatusBelum = "belum"
)

// Lampiran - satu baris grid `TempData.AttachmentList` (`GetAttachmentProdName_Sql` b84).
type Lampiran struct {
	ID string `json:"id"` // `TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF3')`
	// ProdukID - kolom `TREATYID` (= `ProductName.ID`).
	ProdukID string `json:"produkId"`
	Category string `json:"category"` // `"File"` (`ProductNameSaveAttachment` 2.1 b475)
	FileName string `json:"fileName"` // `File Name` b68426 (`.pyFileName`)
	// FileMimeType - `.pyFileMimeType`: EKSTENSI berkas (huruf kecil), bentuk yang
	// dibandingkan syarat `View Office Online` b69291 (`'xls'`, `'docx'`, ...).
	FileMimeType string `json:"fileMimeType"`
	UserName     string `json:"userName"`
	// StorageID - kolom `T_STORAGE_ID` (= `IMAGEID`, `.type` di grid).
	StorageID string `json:"storageId"`
	Status    string `json:"status"`
	// Galat - kalimat kegagalan pengiriman terakhir (status `gagal`).
	Galat string `json:"galat,omitempty"`
}

// ObjekPenyimpanan - satu baris `T_STORAGE_IMAGE` (`Insert_T_Storage_SQL` b85, `GetLinkStorage_SQL` b85,
// `Update_T_Storage_SQL` b85).
type ObjekPenyimpanan struct {
	ImageID string
	// AppFolder - sebelum dikirim: folder objek `Folder + "/Doc/" + YYYY + "/" + MM + "/"` (`InsertGoogleStorage_Act`
	// 8 b1339); sesudah dikirim: `appfolder` jawaban layanan (b2473) - jalur objek PENUH berawalan skema gs + App.
	// Objek yang dicatat stub lokal menyimpan foldernya saja.
	AppFolder string
	FileName  string // `yyyyMMdd-hhmmss-S - <nama>` (8 b1339)
	AppName   string // `T_FOLDER_IMAGE.APPNAME` (`GetAppName_SQL` b58)
	// DurasiDetik - `Durasi="1800"` (`ProductNameSaveAttachment` 2.4 b904) yang dikirim ke layanan.
	DurasiDetik int
	// URLPublic - `URLImage` jawaban layanan (URL bertanda tangan, b2452); kosong = dicatat stub lokal.
	URLPublic string
	// Exp - `EXPDATE` berbentuk `DD/MM/YYYY HH24:MI:SS` GMT (To_date `Insert_T_Storage_SQL`); kosong = NULL.
	Exp string
	// TanggalUpload - `DateTime` jawaban geturl, `MM/DD/YYYY HH24:MI:SS` (`Update_T_Storage_SQL`); kosong = NULL.
	TanggalUpload string
}
