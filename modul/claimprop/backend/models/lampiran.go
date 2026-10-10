package models

// Untuk apa berkas ini: LAMPIRAN KLAIM - activity `GCNMSaveAttachments` (kelas Work-, diekspor work owner 08-10-2026 ke
// korpus Claim Prop) dan master `T_KATEGORI_DOC_KLAIM` (kategoriifile.xls, migrasi 535/536).
//
//	GCNMSaveAttachments S1    setiap berkas `dragDropFileUpload.pxResults`
//	                    S1.2  Local.Category = TempInputParam.pyCategory (kosong -> .pyCategory berkas);
//	                          Local.InsKey = pxCoverInsKey (kosong -> pzInsKey) = kasus klaim
//	                    S1.7  InsertDocument_Act (MIME .pyFileType, KATEGORI_1 Local.Category, NAMAFILE .pyFileName,
//	                          IDPEGA Local.InsKey; KATEGORI_2 / NOAKSEP / NOPREKAS / PAYMENTDATE kosong)
//	                    S1.1, S1.3-S1.6 ter-remark (`//`) - lampiran bawaan Pega (Link-Attachment) TIDAK dibuat
//	AttachmentProtect_ACT S3 `AttachCategory.pxResults` (.ID, .CountAttach) = master kategori PROP x cacah dokumen
//	                          klaim per KATEGORI_1

import (
	"fmt"
	"time"
)

// TypeKlaimProp - T_KATEGORI_DOC_KLAIM.TYPE_KLAIM kategori Claim Prop.
const TypeKlaimProp = "PROP"

// InsertDocument_Act S4 (`InsertGoogleStorage_Act`: Folder "Claim", Durasi 1800).
const (
	FolderLampiranKlaim = "Claim"
	DurasiLampiranKlaim = 1800
)

// KategoriLampiran - satu baris `AttachCategory.pxResults`: kategori master dan cacah berkas klaim ini.
type KategoriLampiran struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	CountAttach int    `json:"countAttach"`
}

// CacahLampiran - `AttachCategory.pxResults` sebagai ID -> CountAttach (masukan AttachmentProtect).
func CacahLampiran(k []KategoriLampiran) map[string]int {
	out := make(map[string]int, len(k))
	for _, x := range k {
		out[x.ID] = x.CountAttach
	}
	return out
}

// Lampiran - satu dokumen klaim tersimpan. Jendela View File mengikuti NB Treaty In (perintah work owner 09-10-2026,
// popup `ReasViewAttachment`: View / View Office Online, Delete) dan screenshot layar Pega View File (File Name,
// Category, Create Date `DD/MM/YYYY HH24:MI`, Action, Delete); Attached By = PXCREATEOPERATOR.
type Lampiran struct {
	ID       string `json:"id"`
	NamaFile string `json:"namaFile"`
	Kategori string `json:"kategori"`
	// Note - KATEGORI_2 (kolom Note popup NB).
	Note string `json:"note"`
	// MIME - ekstensi huruf kecil (InsertDocument_Act S3).
	MIME     string `json:"mime"`
	Tanggal  string `json:"tanggal"`
	Operator string `json:"operator"`
	// AdaObjek - objek penyimpanan tercatat (syarat View / View Office Online, ReasViewAttachment b2951).
	AdaObjek bool `json:"adaObjek"`
	// StorageID - T_STORAGE_ID (IMAGEID penyimpanan); tidak dikirim ke layar.
	StorageID string `json:"-"`
}

// EkstensiOffice - View Office Online hanya untuk berkas Excel, Word, PowerPoint (ReasViewAttachment b2951).
var EkstensiOffice = map[string]bool{"xls": true, "xlsx": true, "doc": true, "docx": true, "ppt": true, "pptx": true}

// DurasiLihatLampiran - GetBase64Attachment S5.1 `GetUrlGoogleStorage_Act` Durasi 1800.
const DurasiLihatLampiran = 1800

// BarisDokumenKlaim - InsertDocument_Act S3 (+ T_STORAGE_ID S4).
type BarisDokumenKlaim struct {
	ID        string    // S3 `@CurrentDate("yyyyMMddhhmmssSSS","Asia/Jakarta")`
	Tanggal   time.Time // S3 `@CurrentDateTime()`
	IDPega    string
	NamaFile  string
	MIME      string // S3 `@toLowerCase(Param.MIME)`
	Kategori1 string
	StorageID string
	Operator  string // S3 `pxRequestor.pyUserIdentifier`
}

// IDDokumenKlaim = InsertDocument_Act S3 `@CurrentDate("yyyyMMddhhmmssSSS","Asia/Jakarta")` - pola Java `hh` = jam
// 12-an (katalog DEV 08-10-2026: "20260521032231511" bertanggal 15:22:31), milidetik tiga angka.
func IDDokumenKlaim(t time.Time) string {
	t = t.In(Jakarta)
	return t.Format("20060102030405") + fmt.Sprintf("%03d", t.Nanosecond()/int(time.Millisecond))
}
