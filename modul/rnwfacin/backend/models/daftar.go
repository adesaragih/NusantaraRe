package models

import "time"

// Status kerja Pega yang dikecualikan daftar kasus renewal. `[terverifikasi]`
// `D:\migrasi\RNM\RNW Fac In\ReportDefinition\RenewalList_RD.xml`
// (ASM-FW-GISFW-WORK-RENEWAL!RENEWALLIST_RD) `pyContent/pyFilters` filter B dan D:
// `.pyStatusWork != "Resolved-Completed"` / `!= "Resolved-Rejected"`.
// `[pertanyaan terbuka]` padanan status kerja ini di sistem baru.
const (
	StatusKerjaSelesai = "Resolved-Completed"
	StatusKerjaDitolak = "Resolved-Rejected"
)

// BarisDaftarRenewal - satu kasus renewal di daftar (tiket R05). `[terverifikasi]`
// sebelas kolom tampil RenewalList_RD `pyContent/pyFields/pyListFields` (properti
// Pega di komentar tiap medan), ditambah dua properti yang hanya dipakai filter
// (StatusKerja, TeamGroup). MulaiBerlaku/AkhirBerlaku bertipe DateTime di Pega;
// disimpan teks apa adanya selama hanya ditampilkan.
//
// ⚠️ NamaTertanggung, NamaMarketing, dan PembuatID ditampilkan apa adanya, tetapi
// nilainya tidak pernah masuk fixture, test, atau artefak (tiket R05).
type BarisDaftarRenewal struct {
	IDKasus         string    // .pyID
	KunciInstans    string    // .pzInsKey
	DibuatPada      time.Time // .pxCreateDateTime
	PembuatID       string    // .pxCreateOperator (login pembuat)
	OldPolicyNo     string    // .Quotation.OldPolicyNo
	NamaTertanggung string    // .Quotation.InsuredName
	NamaMarketing   string    // .Quotation.MarketingName
	MulaiBerlaku    string    // .OfferFacIn.PolicyData.StartDateTime, apa adanya
	AkhirBerlaku    string    // .OfferFacIn.PolicyData.EndDateTime, apa adanya
	NBStatus        string    // .NBStatus - `[pertanyaan terbuka]` artinya
	NBStatusNew     string    // .NBStatusNew - `[pertanyaan terbuka]` artinya
	StatusKerja     string    // .pyStatusWork (filter B, D)
	TeamGroup       string    // .Quotation.TeamGroup (filter C)
}
