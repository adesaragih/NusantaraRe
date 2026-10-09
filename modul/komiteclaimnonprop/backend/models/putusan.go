package models

// Untuk apa berkas ini: KEPUTUSAN PENYETUJU - flow action `ViewTransferDtl` (Section `ShowTransfer`) lalu pasca-proses
// `KomitePostAdjustment` (korpus `Komite Claim Non Prop`). Fungsi di sini MURNI: dari keadaan kasus komite, klaim
// induk, dan isian layar, menyusun rencana tulisan langkah demi langkah. Layanan menjalankannya di SATU transaksi;
// langkah yang butuh basis data (nomor S14.8-S14.11, OS / CLAIMXOL2 / OS subjectivity, efek keluar) ditandai dan
// dikerjakan layanan. Pola disalin dari Komite Claim Prop (bukan impor).
//
// Tidak dibangun / diubah (PARITAS):
//   - S1 / S31 jalur CWP (`IsCloseFile` / `IsReject` kasus komite -> `KomitePostAdjustmentCWP`): kasus komite CWP tidak
//     dapat dilahirkan Claim Non Prop (OQ-CNP-36).
//   - S3-S4 pemutus bukan KomiteID tingkat berjalan: XML hanya memberi pesan lalu lanjut; di sini ditolak (OQ-CNP-16
//     bawaan, pola workbasket Komite Claim Prop - pemutus = anggota workbasket tingkat berjalan).
//   - S14.3-S14.7, S14.14, S19.1 ter-remark (`//`); S14.20 `GetBase64Attachment` dan PDF persetujuan `GenerateAccCNP_act`
//     (stream `AccClaimKomite_HTML` tidak diekspor, OQ-CNP-22).
//   - S25 `UpdateWorkObject` tidak diekspor; S14.16 / S26 Obj-Save + S27 Commit = tulisan transaksi ini.

import (
	"strings"
	"time"

	"nusantarare/inti/backend/kontrak"
)

// PesanKosong - pesan wajib isi (Property-Set-Messages bawaan Pega; teks sama dengan Komite Claim Prop).
const PesanKosong = "Value cannot be blank"

// PesanDiLuarPilihan - kode di luar daftar pilihan dropdown (pertahanan server; layar hanya menawarkan daftar pilihan).
// Teks validasi tabel bawaan Pega tidak ada di ekspor - bukan VERBATIM.
const PesanDiLuarPilihan = "Value is not in the list"

// Label isian `ShowTransfer` (VERBATIM).
const (
	LabelAcceptStatus     = "Are you sure to accept this document?"
	LabelSubjectivity     = "Subjectivity"
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
	// IsSubjectivity - `.IsSubjectivity` "Subjectivity": tampil bila AcceptStatus = 1, nonaktif bila
	// `.KomiteCount != '1'`.
	IsSubjectivity bool `json:"isSubjectivity"`
	// SubjectivityNote - `.SubjectivityNote` (dropdown kode 1..7, SubjectivityNote.xml): tampil + wajib bila Subjectivity,
	// nonaktif bila KomiteCount != '1'.
	SubjectivityNote string `json:"subjectivityNote"`
	// UsulTutup / UsulCadang - `.Adjustment.IsProposeClose` / `.IsPropReserved`, nonaktif bila KomiteCount != '1'.
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

// boolPega - nilai kotak centang Pega ("true" / "false"; `IsSubjectivity == true` KomitePostAdjustment S6-S7).
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
	// Pemutus - akun pelaku yang memutus baris ini (KOMITE_OPERATORID ditimpa: KomiteID workbasket -> akun pemutus -
	// jejak siapa yang memutus dan nama pengirim email, keputusan work owner 09-10-2026); kosong = tidak ditimpa (S26.1).
	Pemutus string
}

// Penerima email `SendEmailKlaim_KMT`.
const (
	EmailPenyetujuBerikut = "berikut"  // bukan tingkat akhir, disetujui -> anggota KomiteCount+1
	EmailPembuatSetuju    = "approval" // tingkat akhir disetujui -> pembuat "(Approval)"
)

// Tolak tanpa email: KomitePostAdjustment S12 melompat ke EXIT (S20), melewati S19 SendEmailKlaim_KMT - cabang
// "(Reject)" stream itu tidak terjangkau di Non Prop.

// Status kasus klaim (`TempMainWork.CNPStatusCase`).
const (
	StatusKlaimDiterima = "CLAIM ACCEPTED" // S14.17
	StatusKlaimDitolak  = "CLAIM REJECTED" // S20
)

// Rencana - hasil satu Submit, urut langkah KomitePostAdjustment.
type Rencana struct {
	// Tangga - S6 / S7 (keputusan tingkat berjalan) dan S11.3 (tolak: baris tingkat berjalan dan sesudahnya).
	Tangga []UbahAnggota
	// Count / AcceptStatus / UsulTutup / UsulCadang - header kasus komite sesudah S20 / S29.
	Count        int
	AcceptStatus string
	UsulTutup    string
	UsulCadang   string
	// SubjectivitySimpan / SubjectivityNoteSimpan - isian Subjectivity tingkat 1 di header kasus komite.
	SubjectivitySimpan     string
	SubjectivityNoteSimpan string
	// Selesai - S28 `KomiteCount >= KomiteLoop` -> ASMForceCaseClose Resolved-Completed.
	Selesai bool
	// Klaim - tulisan balik ke kasus klaim induk (lewat kontrak).
	Klaim kontrak.UbahanKlaimTreaty
	// TingkatAkhirSetuju - `KomiteCount==KomiteLoop && AcceptStatus=="1"`.
	TingkatAkhirSetuju bool
	// Subjectivity - `pyWorkPage.IsSubjectivity` yang dipakai langkah tingkat akhir.
	Subjectivity bool
	// TerbitkanNomor - S14.9-S14.11 (tingkat akhir disetujui, bukan subjectivity, AcceptedNo kosong).
	TerbitkanNomor bool
	// SimpanOS - S14.18 InsertOSKlaimCNP (OS_AKSEPTASI_KLAIM + CLAIMXOL2 + JSON_KLAIM); SimpanOSSubjectivity - S17
	// InsertOSSubjectivityCNP; Konversi - S14.22 (IsPEGAPROD); Kasir - S19.3 (IsPEGAPROD).
	SimpanOS             bool
	SimpanOSSubjectivity bool
	Konversi             bool
	Kasir                bool
	// StatusRiwayat - S22 `InsertHistory.CARI5`: ACCEPT / REJECT.
	StatusRiwayat string
	// Email - S19.2 (`SendEmailKlaim_KMT`, bergerbang IsPEGAPROD): jenis penerima; kosong = tanpa email (tolak).
	Email string
	// Posisi - T_WORK_CLAIM.POSITION sesudah Submit: KomiteID baris menunggu pertama yang tersisa; kosong = selesai.
	Posisi string
}

// Rencanakan menyusun rencana Submit `kep` oleh `akun` (nama tampilan `nama` = `OperatorID.pyUserName`) atas kasus `k`
// (klaim induk `kl`, baris akseptasi posisi `kl.Adjustment`). Pemeriksaan pemegang dan isian dilakukan sebelumnya
// (layanan).
func Rencanakan(k Kasus, kl kontrak.KlaimTreaty, kep Keputusan, akun, nama string, saat time.Time) Rencana {
	r := Rencana{Count: k.Count, UsulTutup: k.UsulTutup, UsulCadang: k.UsulCadang, AcceptStatus: kep.AcceptStatus,
		Klaim: kontrak.UbahanKlaimTreaty{Header: map[string]string{}, Adjustment: map[string]string{}}}
	if r.UsulTutup == "" {
		r.UsulTutup = UsulTidak
	}
	if r.UsulCadang == "" {
		r.UsulCadang = UsulTidak
	}
	// Isian tingkat 1 (Subjectivity, catatannya, dua Propose) - ditulis layar ke pyWorkPage (halaman kerja yang bertahan
	// antar tingkat), nonaktif di tingkat lain: tingkat lain memakai nilai yang tersimpan di header.
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
		// S7 (akseptasi induk SUDAH subjectivity): S6 dilewati, keputusan ditulis ke `ComiteeClaim(<LAST>)` - baris
		// TERAKHIR tangga kasus komite ini (KomiteList dan ComiteeClaim satu baris di sistem baru).
		sasaran = len(k.Tangga) - 1
	}
	jabatan := ""
	if berjalan >= 0 {
		jabatan = k.Tangga[berjalan].Jabatan
	}
	if sasaran >= 0 { // S6 / S7 (+ S11.1): keputusan, komentar, @CurrentDateTime()
		r.Tangga = append(r.Tangga, UbahAnggota{ID: k.Tangga[sasaran].ID, Keputusan: kep.AcceptStatus,
			Komentar: kep.Comment, IsiKomentar: true, Tanggal: saat, Pemutus: akun})
	}
	// S8-S10: InsertChronology_DT "Accepted by " / "Rejected by " + OperatorID.pyUserName.
	teks := ""
	if setuju {
		teks = "Accepted by " + nama
	} else if tolak {
		teks = "Rejected by " + nama
	}
	if teks != "" {
		r.Klaim.Riwayat = append(r.Klaim.Riwayat, kontrak.RiwayatKlaimTreaty{Teks: teks, Pelaku: akun, Tingkat: jabatan,
			Saat: saat})
	}
	count := k.Count
	if tolak {
		// S11.3: baris tangga tingkat berjalan dan sesudahnya ditolak, komentar dan tanggal ikut ditulis; S11.4
		// AcceptanceStatus akseptasi induk := 2; S12 -> EXIT; S20 KomiteCount := KomiteLoop, CNPStatusCase.
		for i, a := range k.Tangga {
			if i == sasaran || a.Urut < k.Count {
				continue
			}
			r.Tangga = append(r.Tangga, UbahAnggota{ID: a.ID, Keputusan: KeputusanTolak, Komentar: kep.Comment,
				IsiKomentar: true, Tanggal: saat})
		}
		r.Klaim.Adjustment["AcceptanceStatus"] = KeputusanTolak
		r.Klaim.Header["CNPStatusCase"] = StatusKlaimDitolak
		count = k.Loop
	}
	if setuju {
		akhir := count == k.Loop
		r.TingkatAkhirSetuju = akhir
		if akhir && !subj { // S14
			r.TerbitkanNomor = adjKlaim["AcceptedNo"] == ""
			// S14.12: TempMainWork.IsCloseFile := Adjustment.IsProposeClose ("1" / "0" - Claim Non Prop membaca "1").
			r.Klaim.Header["IsCloseFile"] = r.UsulTutup
			r.Klaim.Header["CNPStatusCase"] = StatusKlaimDiterima // S14.17
			r.SimpanOS = true                                     // S14.18
			r.Konversi = true                                     // S14.22
			r.Kasir = true                                        // S19.3
		}
		if subj { // S16
			r.Klaim.Adjustment["IsKomite"] = "0"
		}
		if akhir { // S17 (OS subjectivity), S18 (TempMainWork.ClaimData.IsSubjectivity: tanpa kolom / pembaca, PARITAS)
			r.SimpanOSSubjectivity = true
			r.Klaim.Adjustment["IsSubjectivity"] = boolPega(subj)
			r.Klaim.Adjustment["SubjectivityNote"] = note
		}
		// S19.2 SendEmailKlaim_KMT (S12 penyetuju berikut / S14 pembuat).
		r.Email = EmailPenyetujuBerikut
		if akhir {
			r.Email = EmailPembuatSetuju
		}
	}
	// S22-S23: HISTORYAKSEPTASIPEGA ACCEPT / REJECT / kode apa adanya.
	switch {
	case setuju:
		r.StatusRiwayat = StatusRiwayatACC
	case tolak:
		r.StatusRiwayat = StatusRiwayatREJ
	default:
		r.StatusRiwayat = kep.AcceptStatus
	}
	r.Selesai = count >= k.Loop // S28 (sebelum S29)
	r.Count = count + 1         // S29
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

// AdjustmentDiterima - S14.13: AcceptedNo, AcceptedDate := now, AcceptanceStatus := 1.
func (r *Rencana) AdjustmentDiterima(nomor string, saat time.Time) {
	if r.Klaim.Adjustment == nil {
		r.Klaim.Adjustment = map[string]string{}
	}
	r.Klaim.Adjustment["AcceptedNo"] = nomor
	r.Klaim.Adjustment["AcceptanceStatus"] = KeputusanSetuju
	r.Klaim.Adjustment["AcceptedDate"] = FormatWaktu(saat)
}
