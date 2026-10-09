package models

// Untuk apa berkas ini: KEPUTUSAN PENYETUJU - flow action `ViewTransferDtl` (Section `ShowTransfer`, tombol "Submit"
// = finishAssignment) lalu pasca-proses `KomitePost` S1 -> `KomitePostAdjustment` (TT 2). Fungsi di sini MURNI: dari
// keadaan kasus komite, klaim induk, dan isian layar, menyusun rencana tulisan langkah demi langkah. Layanan
// menjalankannya di SATU transaksi; langkah yang butuh basis data (nomor S16.5-S16.9, retro S21, efek keluar)
// ditandai dan dikerjakan layanan.
//
// ⚠️ Langkah ter-remark / tak terjangkau (tidak dibangun):
//   - S16.1-S16.4 (`GETTanggalClosing_SQL`, `GenerateNoAcceptTreaty`), S39 (`ASMForceCaseClose`) - blok `//`.
//   - S7 (S6 dilewati, tulis `ComiteeClaim(<LAST>)`, bila adjustment induk SUDAH subjectivity): Pega mengirim ulang
//     baris subjectivity ke komite (`AddKomiteTreatyChild_ACT` S16-S19/S31), tetapi Claim Prop menautkan `KOMITE_ID`
//     sekali (`AND KOMITE_ID IS NULL`, UNIQUE) - tak terjangkau sampai penyerahan ulang diputuskan (OQ-KCP-06).
//   - `FlagOnGoingCommitte` (S26.1, S27) dibuang (keputusan 26-28, 19-09-2026).
//   - S20 / S37 `UpdateWorkObject` tidak diekspor; S18 / S19 / S36 Obj-Save = tulisan transaksi ini.

import (
	"strings"
	"time"

	"nusantarare/inti/backend/kontrak"
)

// PesanKosong - pesan wajib isi (Property-Set-Messages bawaan Pega; teks sama dengan Claim Prop).
const PesanKosong = "Value cannot be blank"

// PesanDiLuarPilihan - kode di luar daftar pilihan dropdown (pertahanan server; layar hanya menawarkan daftar pilihan).
// Teks validasi tabel bawaan Pega tidak ada di ekspor - bukan VERBATIM.
const PesanDiLuarPilihan = "Value is not in the list"

// Label isian `ShowTransfer` (VERBATIM).
const (
	LabelAcceptStatus     = "Are you sure to accept this document?"
	LabelSubjectivity     = "Subjectivity ?"
	LabelSubjectivityNote = "Subjectivity Note"
	LabelProposeClose     = "Propose To Close Case"
	LabelProposeReserved  = "Propose To Reserved"
	LabelNote             = "Note"
)

// Keputusan - isian layar `ShowTransfer` satu Submit.
type Keputusan struct {
	// AcceptStatus - `.AcceptStatus` (wajib; 1 Approve, 2 Reject - SetDataAcceptationTreaty_Act S2-S3).
	AcceptStatus string `json:"acceptStatus"`
	// Comment - `.Comment` "Note" (wajib).
	Comment string `json:"comment"`
	// IsSubjectivity - `.IsSubjectivity` "Subjectivity ?": tampil bila AcceptStatus = 1 (TT 2), nonaktif bila
	// `.KomiteCount != '1'`.
	IsSubjectivity bool `json:"isSubjectivity"`
	// SubjectivityNote - `.SubjectivityNote` (dropdown kode 1..7, SubjectivityNote.xml): tampil + wajib bila Subjectivity,
	// nonaktif bila KomiteCount != '1'.
	SubjectivityNote string `json:"subjectivityNote"`
	// UsulTutup / UsulCadang - `.Adjustment.IsProposeClose` / `.IsPropReserved` (TT 2), nonaktif bila
	// KomiteCount != '1'.
	UsulTutup  bool `json:"usulTutup"`
	UsulCadang bool `json:"usulCadang"`
}

// IsianTerbuka - isian bertingkat-1 (Subjectivity, catatannya, dua Propose) terbuka hanya di tingkat 1
// (`pyDisabledWhen .KomiteCount!='1'`).
func IsianTerbuka(k Kasus) bool { return k.Count == 1 }

// PeriksaIsian = validasi klien Section ShowTransfer (pyRequired / pyRequiredWhen).
func PeriksaIsian(k Kasus, kep Keputusan) []string {
	var p []string
	if kep.AcceptStatus != KeputusanSetuju && kep.AcceptStatus != KeputusanTolak {
		p = append(p, LabelAcceptStatus+" "+PesanKosong)
	}
	if strings.TrimSpace(kep.Comment) == "" {
		p = append(p, LabelNote+": "+PesanKosong)
	}
	if IsianTerbuka(k) && kep.AcceptStatus == KeputusanSetuju && kep.IsSubjectivity {
		if strings.TrimSpace(kep.SubjectivityNote) == "" {
			p = append(p, LabelSubjectivityNote+": "+PesanKosong)
		} else if !SubjectivityNoteSah(kep.SubjectivityNote) { // dropdown SubjectivityNote.xml
			p = append(p, LabelSubjectivityNote+": "+PesanDiLuarPilihan)
		}
	}
	return p
}

// boolPega - nilai kotak centang Pega ("true" / "false") yang dibaca Claim Prop (`IsSubjectivity == "true"`, Section
// InputAcceptation `.ClaimData.IsCloseFile = true`).
func boolPega(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// UbahAnggota - tulisan atas satu baris tangga.
type UbahAnggota struct {
	ID        string
	Keputusan string
	Komentar  string
	// IsiKomentar - false = komentar tidak disentuh (S26.1 menolak sisa tanpa menyentuh komentar).
	IsiKomentar bool
	Tanggal     time.Time
	// Pemutus - akun pelaku yang memutus baris ini (KOMITE_OPERATORID ditimpa: KomiteID workbasket -> akun pemutus,
	// keputusan work owner 09-10-2026); kosong = tidak ditimpa (S26.1).
	Pemutus string
}

// Penerima email `SendEmailKlaim_KMT`.
const (
	EmailPenyetujuBerikut = "berikut"  // S12: bukan tingkat akhir, disetujui -> anggota KomiteCount+1
	EmailPembuatSetuju    = "approval" // S13-S14: tingkat akhir disetujui -> pembuat "(Approval)"
	EmailPembuatTolak     = "reject"   // S13, S15: ditolak -> pembuat "(Reject)"
)

// Rencana - hasil satu Submit, urut langkah KomitePostAdjustment.
type Rencana struct {
	// Tangga - S6 (keputusan tingkat berjalan) dan S26.1 (tolak otomatis sisa yang menunggu).
	Tangga []UbahAnggota
	// Count / AcceptStatus / UsulTutup / UsulCadang - header kasus komite sesudah S25 / S40.
	Count        int
	AcceptStatus string
	UsulTutup    string
	UsulCadang   string
	// SubjectivitySimpan / SubjectivityNoteSimpan - isian Subjectivity tingkat 1 di header kasus komite (migrasi 682).
	SubjectivitySimpan     string
	SubjectivityNoteSimpan string
	// Selesai - Decision `KomiteLoop` (When IsKomiteLoop salah) -> Resolved-Completed.
	Selesai bool
	// Klaim - tulisan balik ke kasus klaim induk (lewat kontrak).
	Klaim kontrak.UbahanKlaimTreaty
	// TingkatAkhirSetuju - `KomiteCount==TotalKomite && AcceptStatus=="1"`.
	TingkatAkhirSetuju bool
	// Subjectivity - `pyWorkPage.IsSubjectivity` yang dipakai langkah tingkat akhir.
	Subjectivity bool
	// TerbitkanNomor - S16 (AcceptedNo kosong, bukan subjectivity) + S16.6-S16.9 (tingkat akhir disetujui).
	TerbitkanNomor bool
	// SimpanAkseptasi - S17 SaveAcceptation_Act; SimpanRetro - S21 SaveAcceptationTreaty_TKMT; Kasir - S34.
	SimpanAkseptasi bool
	SimpanRetro     bool
	Kasir           bool
	// JSONKlaim - S28; Konversi - S29; LogAkseptasi - S30-S31.
	JSONKlaim    bool
	Konversi     bool
	LogAkseptasi bool
	// StatusRiwayat - S32 `InsertHistory.CARI5`: ACCEPT / REJECT.
	StatusRiwayat string
	// Email - S35 (`SendEmailKlaim_KMT`, bergerbang IsPEGAPROD): jenis penerima.
	Email string
	// Posisi - T_WORK_CLAIM.POSITION sesudah Submit: KomiteID baris menunggu pertama yang tersisa; kosong = selesai.
	Posisi string
}

// Rencanakan menyusun rencana Submit `kep` oleh `akun` atas kasus `k` (klaim induk `kl`, baris adjustment posisi
// `kl.Adjustment`). Pemeriksaan pemegang dan isian dilakukan sebelumnya (layanan).
func Rencanakan(k Kasus, kl kontrak.KlaimTreaty, kep Keputusan, akun string, saat time.Time) Rencana {
	r := Rencana{Count: k.Count, UsulTutup: k.UsulTutup, UsulCadang: k.UsulCadang, AcceptStatus: kep.AcceptStatus,
		Klaim: kontrak.UbahanKlaimTreaty{Header: map[string]string{}, Adjustment: map[string]string{}}}
	if r.UsulTutup == "" {
		r.UsulTutup = UsulTidak
	}
	if r.UsulCadang == "" {
		r.UsulCadang = UsulTidak
	}
	// Isian tingkat 1 (Subjectivity, catatannya, dua Propose) - ditulis layar langsung ke pyWorkPage (halaman kerja yang
	// bertahan antar tingkat), nonaktif di tingkat lain: tingkat lain memakai nilai yang tersimpan di header.
	subj, note := k.Subjectivity == UsulYa, k.SubjectivityNote
	if IsianTerbuka(k) {
		r.UsulTutup, r.UsulCadang = UsulTidak, UsulTidak
		if kep.UsulTutup {
			r.UsulTutup = UsulYa
		}
		if kep.UsulCadang {
			r.UsulCadang = UsulYa
		}
		subj, note = kep.AcceptStatus == KeputusanSetuju && kep.IsSubjectivity, ""
		if subj {
			note = kep.SubjectivityNote
		}
	}
	r.SubjectivitySimpan, r.SubjectivityNoteSimpan = UsulTidak, note
	if subj {
		r.SubjectivitySimpan = UsulYa
	}
	r.Subjectivity = subj
	setuju := kep.AcceptStatus == KeputusanSetuju
	tolak := kep.AcceptStatus == KeputusanTolak
	adjKlaim := AdjustmentKlaim(kl)
	berjalan := k.barisBerjalan()
	sasaran := berjalan
	if adjKlaim["IsSubjectivity"] == "true" && len(k.Tangga) > 0 {
		// Kirim ulang baris subjectivity (keputusan work owner 08-10-2026, OQ-KCP-06 "a"): S6 dilewati, S7 menulis
		// `ComiteeClaim(<LAST>)` - anggota yang ditambahkan putaran ini, baris TERAKHIR tangga kasus komite ini.
		sasaran = len(k.Tangga) - 1
	}
	jabatan := ""
	if berjalan >= 0 {
		jabatan = k.Tangga[berjalan].Jabatan // S8-S9: KomiteList(KomiteCount).IDKomite
	}
	if sasaran >= 0 {
		// S6 / S7: keputusan, komentar, @CurrentDateTime() - KomiteList dan ComiteeClaim satu baris di sistem baru.
		r.Tangga = append(r.Tangga, UbahAnggota{ID: k.Tangga[sasaran].ID, Keputusan: kep.AcceptStatus,
			Komentar: kep.Comment, IsiKomentar: true, Tanggal: saat, Pemutus: akun})
	}
	// S8-S10: InsertChronology_DT ("Accepted by " / "Rejected by " + KomiteList(KomiteCount).IDKomite).
	teks := ""
	if setuju {
		teks = "Accepted by " + jabatan
	} else if tolak {
		teks = "Rejected by " + jabatan
	}
	if teks != "" {
		r.Klaim.Riwayat = append(r.Klaim.Riwayat, kontrak.RiwayatKlaimTreaty{Teks: teks, Pelaku: akun, Tingkat: jabatan,
			Saat: saat})
	}
	// S11 (tanpa gerbang): IsCloseFile / IsReservedClaim <- Adjustment.IsProposeClose / IsPropReserved.
	r.Klaim.Header["ClaimData.IsCloseFile"] = boolPega(r.UsulTutup == UsulYa)
	r.Klaim.Header["ClaimData.IsReservedClaim"] = boolPega(r.UsulCadang == UsulYa)
	count := k.Count
	if setuju { // S12: tolak -> lompat ke EXT (S25)
		akhir := count == k.Loop // S13 Local.TotalKomite := KomiteLoop
		r.TingkatAkhirSetuju = akhir
		if akhir {
			r.Klaim.Header["IsAnyAcceptation"] = "1" // S14
			r.Klaim.Adjustment["IsApproved"] = KeputusanSetuju
		}
		// S15: IsApproved <- ComiteeClaim(<LAST>).KomiteAproval - baris terakhir tangga = baris tingkat akhir yang baru
		// diputuskan S6 (di tingkat lain gerbangnya KomiteCount==TotalKomite salah).
		// S16: AcceptedNo kosong DAN bukan subjectivity; S16.6-S16.9 tingkat akhir disetujui.
		r.TerbitkanNomor = akhir && adjKlaim["AcceptedNo"] == "" && !subj
		r.SimpanAkseptasi = akhir && !subj // S17
		if r.SimpanAkseptasi {
			// SaveAcceptation_Act S6: IsOutstanding := 1, IsCFS := "" (S1 TempOpenPage.stsReject := 1 tidak disimpan).
			r.Klaim.Header["IsOutstanding"] = "1"
			r.Klaim.Header["IsCFS"] = ""
		}
		r.SimpanRetro = akhir && !subj // S21
		if akhir {
			r.Klaim.Header["AktifButton"] = "0" // S22
			if subj {
				r.Klaim.Adjustment["IsKomite"] = "0" // S23
			}
			// S24: penanda dan catatan subjectivity ke adjustment dan header induk.
			r.Klaim.Adjustment["IsSubjectivity"] = boolPega(subj)
			if !subj {
				note = ""
			}
			r.Klaim.Adjustment["SubjectivityNote"] = note
			r.Klaim.Header["ClaimData.IsSubjectivity"] = boolPega(subj)
			r.Klaim.Adjustment["Notes"] = kep.Comment                  // S27 (KomiteCount == KomiteLoop)
			r.JSONKlaim, r.Konversi, r.LogAkseptasi = true, true, true // S28, S29, S30-S31
			r.Kasir = !subj                                            // S34
		}
	}
	if tolak {
		// S25 [EXT]: KomiteCount := KomiteLoop, AcceptanceStatus := 2, AktifButton := 0.
		count = k.Loop
		r.Klaim.Adjustment["AcceptanceStatus"] = KeputusanTolak
		r.Klaim.Header["AktifButton"] = "0"
		// S26 / S26.1: setiap baris tangga yang masih menunggu ditolak otomatis (+ tanggal); komentar tidak disentuh.
		for i, a := range k.Tangga {
			if i == sasaran || a.Keputusan != KeputusanMenunggu {
				continue
			}
			r.Tangga = append(r.Tangga, UbahAnggota{ID: a.ID, Keputusan: KeputusanTolak, Tanggal: saat})
		}
		r.Klaim.Adjustment["Notes"] = kep.Comment // S27 (KomiteCount == KomiteLoop sesudah S25)
	}
	// S32-S33: HISTORYAKSEPTASIPEGA ACCEPT / REJECT (tanpa gerbang).
	switch {
	case setuju:
		r.StatusRiwayat = "ACCEPT"
	case tolak:
		r.StatusRiwayat = "REJECT"
	default:
		r.StatusRiwayat = kep.AcceptStatus
	}
	// S35 SendEmailKlaim_KMT (S12 / S14 / S15).
	switch {
	case tolak:
		r.Email = EmailPembuatTolak
	case setuju && count == k.Loop:
		r.Email = EmailPembuatSetuju
	case setuju:
		r.Email = EmailPenyetujuBerikut
	}
	r.Count = count + 1 // S40 (prakondisi nonaktif - selalu)
	r.Selesai = !MasihBerjalan(kep.AcceptStatus, r.Count, k.Loop)
	if !r.Selesai {
		r.Posisi = posisiSesudah(k, r.Tangga)
	}
	if len(r.Klaim.Adjustment) == 0 {
		r.Klaim.Adjustment = nil
	}
	return r
}

// posisiSesudah - KomiteID baris tangga menunggu PERTAMA sesudah ubahan `ubah` (giliran berikut, KomiteRouter S6.1).
func posisiSesudah(k Kasus, ubah []UbahAnggota) string {
	diputus := map[string]bool{}
	for _, u := range ubah {
		diputus[u.ID] = true
	}
	for _, a := range k.Tangga {
		if a.Keputusan == KeputusanMenunggu && !diputus[a.ID] {
			return a.OperatorID
		}
	}
	return ""
}

// Kepala - tulisan kepala kasus komite rencana ini.
func (r Rencana) Kepala() Kepala {
	return Kepala{Count: r.Count, AcceptStatus: r.AcceptStatus, UsulTutup: r.UsulTutup, UsulCadang: r.UsulCadang,
		Subjectivity: r.SubjectivitySimpan, SubjectivityNote: r.SubjectivityNoteSimpan}
}

// AdjustmentDiterima - S16.9: AcceptedNo, AcceptanceStatus := 1, AcceptedDate := now.
func (r *Rencana) AdjustmentDiterima(nomor string, saat time.Time) {
	if r.Klaim.Adjustment == nil {
		r.Klaim.Adjustment = map[string]string{}
	}
	r.Klaim.Adjustment["AcceptedNo"] = nomor
	r.Klaim.Adjustment["AcceptanceStatus"] = KeputusanSetuju
	r.Klaim.Adjustment["AcceptedDate"] = FormatWaktu(saat)
}
