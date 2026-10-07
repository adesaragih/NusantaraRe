package models

// Untuk apa berkas ini: MUATAN KONVERSI ke Arasapas - Utility2 `Activity/serviceInsertArasapas_act` (kelas Work,
// ruleset 01-01-95; 18 langkah aktif, 2 `//`) sesudah endorsemen selesai. Pola: salinan
// `modul/nbtreatyin/backend/models/konversi.go`; bentuk muatan menurut korpus EDM, yang - berbeda dari NB -
// MEMUAT rule Connect-REST-nya:
//
//	langkah 3-4  OfferFacIn.PolicyData.PolicyNo = RDB GetPolicyNoByCaseId (`SELECT nopolis FROM json_polis WHERE
//	             idpega = pzInsKey`) - nomor polis induk generasi ini
//	langkah 5    GetLinkService(Kategori_1 "Production", Kategori_2 "convertJsonNusareToProduction") -> LinkService
//	langkah 6    Connect-REST `convertJsonNusareToProduction` POST, parameter kueri noPolis = .OfferFacIn.PolicyData.
//	             PolicyNo, caseId = .pzInsKey, tglInput = .pxCreateDateTime; respons -> .StatusService
//	             (When IsSuccessHitService: StsKonversiFacIn = 1 && (StsKonversiFacOut = "" || = 1))
//	langkah 8    pesan gagal kecuali IsTreatyIn (Quotation.BusinessFac = "T") - kasus treaty TIDAK diberi pesan
//	langkah 11-20 STS_KONVERSI, PEGA_DELETE_ERROR_KONVERSI, INSERTJSONPOLISMONITORING (JSON halaman OfferFacIn), UPDATE
//	             ERR_NOTE, surel ke alamat tertanam; langkah 18 ASMForceCaseClose `//`
//
// ⛔ `[terbuka]` Sambungan nyata TIDAK dibangun (ketetapan NB `services/konversi.go`: pengirim bawaan gagal terang,
// menyambung layanan luar butuh persetujuan manusia). Langkah 11-20 menunggu jawaban work owner (daftar pertanyaan).

import "strings"

// PesanGagalKonversi - VERBATIM `FlagErrorKonversi`.
const PesanGagalKonversi = "Gagal Konversi, Silahkan Coba Lagi atau Hub IT !"

// MuatanKonversi adalah tiga parameter kueri Connect-REST `convertJsonNusareToProduction`.
type MuatanKonversi struct {
	NoPolis  string `json:"noPolis"`
	CaseID   string `json:"caseId"`
	TglInput string `json:"tglInput"`
}

// RakitMuatanKonversi menyusun muatan dari halaman kasus yang selesai: noPolis = PolicyNo (json_polis.NOPOLIS
// baris IDPEGA ini), caseId = ID kasus (`KunciInstans`), tglInput = waktu lahir kasus (pxCreateDateTime) APA
// ADANYA (T_WORK_POLIS.TGL_CREATE "2006-01-02 15:04:05") - format DateTime yang diharapkan layanan tidak ada di
// korpus (butir terbuka bersama sambungannya).
func RakitMuatanKonversi(id string, h *Halaman, tglCreate string) MuatanKonversi {
	return MuatanKonversi{NoPolis: h.Ambil(HalamanPolis + ".PolicyNo"), CaseID: KunciInstans(id), TglInput: strings.TrimSpace(tglCreate)}
}
