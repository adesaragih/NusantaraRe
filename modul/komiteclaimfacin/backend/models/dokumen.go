package models

// Untuk apa berkas ini: BARIS DOKUMEN KLAIM - `PrintPDFAccep_MultiAksep_KMT` S16-S22 (nama berkas, kategori) dan S27
// `InsertDocument_Act` (kelas `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`): PDF akseptasi disimpan lewat `inti/backend/penyimpanan`
// (InsertGoogleStorage_Act, Folder "Claim") lalu dicatat di tabel warisan DOCUMENT_CLAIM (pola Komite Claim Prop,
// disalin bukan diimpor).

import (
	"fmt"
	"time"
)

// Kategori / folder dokumen akseptasi (S16-S22 `param.AttachmentCategory`, InsertDocument_Act S4 Folder).
const (
	KategoriDokumenAkseptasi = "AcceptanceNote"
	FolderDokumenKlaim       = "Claim"
	MIMEDokumenAkseptasi     = "pdf" // S27 `MIME="pdf"`
)

// DurasiDokumenKlaim - InsertDocument_Act S4 `Durasi=1800` (detik).
const DurasiDokumenKlaim = 1800

// BarisDokumenKlaim - satu baris DOCUMENT_CLAIM (InsertDocument_Act S3, T_STORAGE_ID S4).
type BarisDokumenKlaim struct {
	ID        string    // S3 `@CurrentDate("yyyyMMddhhmmssSSS","Asia/Jakarta")`
	Tanggal   time.Time // S3 `@CurrentDateTime()`
	IDPega    string    // S27 `TempMainPage.pzInsKey` = ID kasus klaim (KunciInstans)
	NamaFile  string
	MIME      string
	Kategori1 string
	StorageID string // S4 T_STORAGE_IMAGE.IMAGEID
	Operator  string // pxCreateOperator
}

// IDDokumenKlaim = S3 `yyyyMMddhhmmssSSS` (hh 12 jam, Asia/Jakarta - VERBATIM pola Java).
func IDDokumenKlaim(t time.Time) string {
	t = t.In(Jakarta)
	return t.Format("20060102030405") + fmt.Sprintf("%03d", t.Nanosecond()/int(time.Millisecond))
}

// NamaBerkasAkseptasi = S16-S22 `"Persetujuan " + "  " + TempMainPage.ClaimData.ClaimNo + " AcceptNo " +
// TempAcceptedNo.Policy.AccountNo + ".pdf"` (`ClaimData.ClaimNo` = ID kasus klaim, `Policy.AccountNo` = AcceptedNo).
func NamaBerkasAkseptasi(claimNo, acceptedNo string) string {
	return "Persetujuan " + "  " + claimNo + " AcceptNo " + acceptedNo + ".pdf"
}
