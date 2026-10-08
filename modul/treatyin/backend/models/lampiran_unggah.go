package models

// LampiranBaru - satu berkas yang SUDAH terkirim ke penyimpanan dan siap
// dicatat: baris `T_STORAGE_IMAGE` (`Insert_T_Storage_SQL`) dan baris
// `M_ATTACHMENTTREATY_2` (`InsertAttachment2_Sql`), satu transaksi.
type LampiranBaru struct {
	// `M_ATTACHMENTTREATY_2`
	IDKontrak    string // TREATYID — `TreatyIn.ID`
	KodeKategori string // CATEGORY_ID — `StatusDoc.CARI40`
	NamaKategori string // CATEGORY — `DropFile.pyCategory` (= `StatusDoc.CARI41`)
	NamaBerkas   string // FILENAME — `DropFile.pyFileName`, nama asli
	Ekstensi     string // FILEMIMETYPE — `DropFile.pyFileMimeType` (data: `pdf`, `xlsx`)
	Pengguna     string // USERNAME — `OperatorID.pyUserIdentifier`

	// `T_STORAGE_IMAGE`
	ImageID   string // IMAGEID — `GenerateImageID_SQL`
	URLPublik string // URLPUBLIC — `UploadDoc.Response.URLImage`
	AppFolder string // APPFOLDER — `UploadDoc.Response.appfolder`
	Exp       string // EXPDATE — `DD/MM/YYYY HH24:MI:SS`, kosong = NULL
	NamaObjek string // FILENAME — `UploadDoc.Namafile` yang dikirim
	App       string // APPNAME — `T_FOLDER_IMAGE.APPNAME`
}

// ObjekSimpanan - satu baris `T_STORAGE_IMAGE` (`GetLinkStorage_SQL`).
type ObjekSimpanan struct {
	ImageID   string
	URLPublik string // URLPUBLIC — URL bertanda tangan
	AppFolder string // APPFOLDER — jalur objek PENUH berawalan skema gs + App, lalu Folder + Namafile
	Exp       string // EXPDATE — `DD/MM/YYYY HH24:MI:SS`
	App       string // APPNAME
	NamaObjek string // FILENAME
}

// PerubahanKategori - satu baris yang `Change Category` → `Save` tulis
// (`ChangeKateAttachment2_Sql`).
type PerubahanKategori struct {
	IDLampiran   string
	KodeKategori string
	NamaKategori string
}
