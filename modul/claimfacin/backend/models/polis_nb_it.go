package models

// Tab Loss Record dan Correspondence tampilan polis NB hanya-baca (perintah work owner 11-10-2026 "tampilan view polis
// yang belum tergenerate tolong diperbaiki. cek tampilan nya di nb facin!"; keputusan work owner 11-10-2026 "Ikut NB:
// hanya akun IT"): seperti nbfacin (tiket 55 / 56) kedua tab HANYA untuk akun ber-M_LOGIN_GO.JOB_POSITION "IT", ditegakkan
// di layanan (403), dan datanya TIDAK ikut halaman polis (`OfferFacIn`) yang dibaca semua pelaku - dibaca lewat rute
// sendiri. Sumber = tabel nbfacin ber-induk kasus NB (`T_WORK_POLIS.ID`): T_LOSSRECORD -> T_LOSSRECORDLIST (migrasi
// nbfacin tiket 55), T_CORRESPONDENCELIST -> T_RECIPIENTLIST (tiket 56). Nilai TERSIMPAN apa adanya (pola keputusan
// work owner Premium Deduction / Payment / Scoring): Uw Year, Loss Amount ( IDR ), dan Loss Ratio hasil Save NB terakhir
// tidak dihitung ulang.

import "strings"

// JabatanTabIT - M_LOGIN_GO.JOB_POSITION pemilik tab Loss Record / Correspondence (nbfacin `jabatanTabIT`).
const JabatanTabIT = "IT"

// JabatanIT - jabatan akun = "IT" (tanpa beda huruf, spasi tepi diabaikan - nbfacin `jabatanIT`).
func JabatanIT(jabatan string) bool {
	return strings.EqualFold(strings.TrimSpace(jabatan), JabatanTabIT)
}

// LossRecordNB - kepala T_LOSSRECORD + baris T_LOSSRECORDLIST satu kasus NB (teks desimal bertitik; "" = kosong).
type LossRecordNB struct {
	PremiPar   string `json:"premiPar"`
	PremiMb    string `json:"premiMb"`
	PremiPl    string `json:"premiPl"`
	PremiOther string `json:"premiOther"`
	// Gross1 .. Net5 - LOSS_RATIO{1,3,5}_YEAR_PERCENT / NET_LOSS_RATIO{1,3,5}_YEAR.
	Gross1 string `json:"gross1"`
	Gross3 string `json:"gross3"`
	Gross5 string `json:"gross5"`
	Net1   string `json:"net1"`
	Net3   string `json:"net3"`
	Net5   string `json:"net5"`
	// Baris - urut SEQ_NO, ID.
	Baris []BarisLossRecordNB `json:"baris"`
}

// BarisLossRecordNB - satu baris T_LOSSRECORDLIST (DateOfLoss teks tersimpan: stempel Pega / tanggal).
type BarisLossRecordNB struct {
	UwYear          string `json:"uwYear"`
	ClassOfBusiness string `json:"classOfBusiness"`
	DateOfLoss      string `json:"dateOfLoss"`
	CauseOfLoss     string `json:"causeOfLoss"`
	Currency        string `json:"currency"`
	LossAmount      string `json:"lossAmount"`
	LossAmountIDR   string `json:"lossAmountIdr"`
}

// KorespondensiNB - satu T_CORRESPONDENCELIST + penerimanya (urut SEQ_NO, ID). Tanggal / DikirimPada = teks tersimpan;
// Message = HTML tersimpan (layar membersihkannya sebelum ditampilkan, aturan nbfacin `bersihkanHtml`).
type KorespondensiNB struct {
	Tanggal    string       `json:"tanggal"`
	Sender     string       `json:"sender"`
	SenderMail string       `json:"senderMail"`
	Subject    string       `json:"subject"`
	Message    string       `json:"message"`
	Terkirim   bool         `json:"terkirim"`
	Penerima   []PenerimaNB `json:"penerima"`
}

// PenerimaNB - satu T_RECIPIENTLIST.
type PenerimaNB struct {
	Recipient     string `json:"recipient"`
	RecipientMail string `json:"recipientMail"`
}

// StatusKorespondensiTerkirim - IS_EDITABLE korespondensi terkirim (nbfacin `statusTerkirim`; draf "1").
const StatusKorespondensiTerkirim = "0"

// PolisNBIT - isi kedua tab IT satu polis; Kasus = ID kasus NB ("" = polis tanpa kasus NB sistem baru: kedua tab kosong).
type PolisNBIT struct {
	Kasus         string            `json:"-"`
	LossRecord    LossRecordNB      `json:"lossRecord"`
	Korespondensi []KorespondensiNB `json:"korespondensi"`
}
