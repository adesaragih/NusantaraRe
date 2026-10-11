package models

// Untuk apa berkas ini: KEPUTUSAN PENYETUJU - Section `ShowTransfer` (isian) lalu pasca-proses `KomitePostAct`
// (korpus `Komite Claim FacIn`): S2 TT2 `KomitePost_Adjustment`, S4 TT3 `KomitePost_Reject`, S5 TT4
// `KomitePost_CloseClaim` (S3 TT1 survey ter-remark). Fungsi di sini MURNI: dari keadaan kasus komite, klaim induk, dan
// isian layar, menyusun rencana tulisan langkah demi langkah. Layanan menjalankannya di SATU transaksi; langkah yang
// butuh basis data / outbox ditandai.
//
// PERBAIKAN kelainan XML (prompt §5, OQ-CFI-03):
//   - butir 2: `KomitePost_Adjustment` S12.1.1 `.AcceptanceStatus = "1"` (tanda `=` tunggal) dibaca pembandingan;
//   - butir 3: S7.2.1.9.1 transisi ke label yang tidak ada - tanpa lompatan, langkah berikut berjalan;
//   - butir 6: Obj-Save klaim + komite per iterasi (S7.2.1.10-S7.2.1.14, S26-S27) = satu transaksi Submit.
//
// Tidak dibangun (PARITAS): `ClaimData.ObjectList(o).IsKomiteApprove` / `.KomiteList` / `FlagOnGoingCommitte` /
// `IDObjectUpdate` / `ClaimComitee` klaim induk (tanpa kolom, tanpa pembaca), `SetDataForInformation_Act` S29 (halaman
// informasi tanpa pembaca), `ApprovalCommite` (KomitePost_Reject S7, bahan stream PDF yang tidak diekspor).

import (
	"strings"
	"time"

	"nusantarare/inti/backend/kontrak"
)

// PesanKosong - pesan wajib isi (Property-Set-Messages bawaan Pega; teks sama dengan Komite Claim Prop / Non Prop).
const PesanKosong = "Value cannot be blank"

// Label isian `ShowTransfer` LS45 (VERBATIM).
const (
	LabelAcceptStatus    = "Are you sure to accept this document?"
	LabelProposeClose    = "Propose To Close Case"
	LabelProposeReserved = "Propose To Reserved"
	LabelNote            = "Note"
)

// Keputusan - isian layar `ShowTransfer` satu Submit.
type Keputusan struct {
	// AcceptStatus - `.AcceptStatus` (wajib; 1 Approve, 2 Reject).
	AcceptStatus string `json:"acceptStatus"`
	// Comment - `.Comment` "Note" (wajib).
	Comment string `json:"comment"`
	// UsulTutup / UsulCadang - `.Adjustment.IsProposeClose` / `.IsPropReserved` (TT2, nonaktif bila KomiteCount != '1').
	UsulTutup  bool `json:"usulTutup"`
	UsulCadang bool `json:"usulCadang"`
}

// IsianTerbuka - dua Propose terbuka hanya di tingkat 1 TT2 (`VIS .TransferType = 2`, `NA .KomiteCount != '1'`).
func IsianTerbuka(k Kasus) bool { return k.TransferType == TransferAdjustment && k.Count == 1 }

// PeriksaIsian = validasi klien Section ShowTransfer (LS45 `REQ=true`).
func PeriksaIsian(kep Keputusan) []string {
	var p []string
	if kep.AcceptStatus != KeputusanSetuju && kep.AcceptStatus != KeputusanTolak {
		p = append(p, LabelAcceptStatus+" "+PesanKosong)
	}
	if strings.TrimSpace(kep.Comment) == "" {
		p = append(p, LabelNote+": "+PesanKosong)
	}
	return p
}

// UbahAnggota - tulisan atas satu baris tangga.
type UbahAnggota struct {
	ID        string
	Keputusan string
	Komentar  string
	// IsiKomentar - false = komentar tidak disentuh (S7.2.1.9.1.1 menolak sisa tanpa komentar).
	IsiKomentar bool
	Tanggal     time.Time
	// Pemutus - akun pelaku yang memutus baris ini (KOMITE_OPERATORID ditimpa: KomiteID workbasket -> akun pemutus,
	// pola Komite Claim Prop 09-10-2026); kosong = tidak ditimpa.
	Pemutus string
}

// Jenis email `SendEmailKlaim_KMT` / KomitePost_Reject S13.10 (penerima dirakit saat dikirim).
const (
	EmailPenyetujuBerikut = "berikut"  // S12: bukan tingkat akhir, disetujui -> ComiteeClaim(KomiteCount + 1)
	EmailPembuatSetuju    = "approval" // S14: tingkat akhir disetujui -> pembuat "(Approval)"
	EmailPembuatTolak     = "reject"   // S15: ditolak -> pembuat "(Reject)"
)

// Rencana - hasil satu Submit, urut langkah KomitePost_*.
type Rencana struct {
	// Tangga - keputusan tingkat berjalan (S7.2.1.8 / Reject S6) dan penolakan sisa (S7.2.1.9.1.1).
	Tangga []UbahAnggota
	// Kepala - KOMITE_COUNT (S14 / S24, Reject S16), ACCEPT_STATUS, usul tingkat 1, KOMITE_LOOP.
	Kepala Kepala
	// Selesai - decision IsKomiteLoop salah (`AcceptStatus == "1" && KomiteCount <= KomiteLoop` tidak terpenuhi) -> END52.
	Selesai bool
	// Posisi - T_WORK_CLAIM.POSITION sesudah Submit (KomiteID tingkat berikut); kosong = selesai.
	Posisi string
	// Klaim - tulisan balik ke kasus klaim induk (lewat kontrak).
	Klaim kontrak.UbahanKlaimFacIn
	// TingkatAkhirSetuju - TT2 `KomiteCount == KomiteLoop && AcceptStatus == "1"`.
	TingkatAkhirSetuju bool
	// TerbitkanNomor - S7.2.1.2.4-S7.2.1.2.7; Tanda - SaveAcceptation_KMT (S7.2.1.16); SimpanOS - S8; JSONKlaim - S16 /
	// Reject S11; LogAkseptasi - S18-S19 (bukan Fac Retro); Kasir - S25 (HitServiceToKasirKMT_Act).
	TerbitkanNomor bool
	Tanda          bool
	SimpanOS       bool
	JSONKlaim      bool
	LogAkseptasi   bool
	Kasir          bool
	// Konversi - STSREJECT `KonversiKlaim_Act` ("1" TT2 S17 bukan Fac Retro, "2" TT3 S15, "4" TT4 S12.4); kosong = tidak.
	Konversi string
	// StatusRiwayat - HISTORYAKSEPTASIPEGA S20 (TT2): ACCEPT / REJECT; kosong = tanpa baris (TT3 / TT4).
	StatusRiwayat string
	// SubProgres - S22-S23 SUBPROGRESSCLAIM.POSITION2 (TT2 tingkat akhir atau tolak).
	SubProgres string
	// Email - jenis email (kosong = tanpa email); Subjek - "Akseptasi" (TT2) / "Reject" / "Close" (TT3 / TT4).
	Email string
	// OSTolak - TT3 S12 (SaveReject_ACT_KMT); OSTutup - TT4 S12.1-S12.3; KlaimDitolak - TT3 S14 CLAIMREJECTED.
	OSTolak, OSTutup, KlaimDitolak bool
	// TutupKlaim - pxForceCaseClose klaim induk (TT3 S17 Resolved-Rejected, TT4 S14 Resolved-Completed).
	TutupKlaim string
	// PDF - stream HTML dokumen PDF yang diterbitkan Pega (tidak diekspor di korpus - berkas tidak dibuat, OQ-KCFI-01).
	PDF string
}

// teksRiwayat - S3 / S4 (TT2), Reject / Close S2-S4: "Accepted by <IDKomite> - <KMT>" / "Rejected by ...".
func teksRiwayat(status, jabatan, kmt string) string {
	switch status {
	case KeputusanSetuju:
		return "Accepted by " + jabatan + " - " + kmt
	case KeputusanTolak:
		return "Rejected by " + jabatan + " - " + kmt
	}
	return ""
}

// bolPenanda - kotak centang usul -> penanda tersimpan (`KOMITE_USUL_*` / IsCloseFile klaim '1' / '0').
func bolPenanda(b bool) string {
	if b {
		return UsulYa
	}
	return UsulTidak
}

// KomiteIDTingkat - KomiteID anggota tangga tingkat `n` (KomiteRouter: assignment tingkat berjalan; SendEmailKlaim_KMT
// S12 `ComiteeClaim(KomiteCount + 1)`).
func KomiteIDTingkat(k Kasus, n int) string {
	for _, a := range k.Tangga {
		if a.Urut == n {
			return a.OperatorID
		}
	}
	return ""
}

// Rencanakan menyusun rencana Submit `kep` oleh `akun` atas kasus `k` (tangga SUDAH berisi perluasan KCF-02 bila ada;
// `k.Loop` = cacahnya) dan klaim induk `kl`. Pemeriksaan pemegang dan isian dilakukan sebelumnya (layanan).
func Rencanakan(k Kasus, kl kontrak.KlaimFacIn, kep Keputusan, akun string, saat time.Time) Rencana {
	if k.TransferType == TransferReject || k.TransferType == TransferClose {
		return rencanaTutup(k, kep, akun, saat)
	}
	return rencanaAdjustment(k, kl, kep, akun, saat)
}

// rencanaAdjustment = KomitePost_Adjustment (TT2).
func rencanaAdjustment(k Kasus, kl kontrak.KlaimFacIn, kep Keputusan, akun string, saat time.Time) Rencana {
	r := Rencana{Klaim: kontrak.UbahanKlaimFacIn{Header: map[string]string{}, Adjustment: map[string]string{},
		Objek: map[string]string{}, Item: map[string]string{}}}
	setuju, tolak := kep.AcceptStatus == KeputusanSetuju, kep.AcceptStatus == KeputusanTolak
	usulTutup, usulCadang := k.UsulTutup, k.UsulCadang
	if IsianTerbuka(k) { // LS45 dua Propose (tingkat 1), disimpan di kepala kasus komite
		usulTutup, usulCadang = bolPenanda(kep.UsulTutup), bolPenanda(kep.UsulCadang)
	}
	if usulTutup == "" {
		usulTutup = UsulTidak
	}
	if usulCadang == "" {
		usulCadang = UsulTidak
	}
	adj := Adjustment(kl)
	akhir := k.Count == k.Loop
	jabatan := ""
	if i := k.barisBerjalan(); i >= 0 {
		jabatan = k.Tangga[i].Jabatan
		// S7.2.1.4 / S7.2.1.5 / S7.2.1.6 (ComiteeClaim(count)) dan S7.2.1.8 (KomiteList(count)) - satu baris tangga.
		r.Tangga = append(r.Tangga, UbahAnggota{ID: k.Tangga[i].ID, Keputusan: kep.AcceptStatus, Komentar: kep.Comment,
			IsiKomentar: true, Tanggal: saat, Pemutus: akun})
	}
	if t := teksRiwayat(kep.AcceptStatus, jabatan, k.ID); t != "" { // S3-S5
		r.Klaim.Riwayat = append(r.Klaim.Riwayat, kontrak.RiwayatKlaimFacIn{Teks: t, Pelaku: akun, Tingkat: jabatan,
			Saat: saat})
	}
	if setuju && akhir { // S7.2.1.2 nomor, S7.2.1.5 AcceptanceStatus / AcceptedNo / AcceptedDate / Notes
		r.TingkatAkhirSetuju, r.TerbitkanNomor = true, true
		r.Klaim.Adjustment["AcceptanceStatus"] = KeputusanSetuju
		r.Klaim.Adjustment["AcceptedDate"] = FormatWaktu(saat)
		r.Klaim.Adjustment["Notes"] = kep.Comment
	}
	if tolak { // S7.2.1.6 lalu S7.2.1.9: sisa tangga menunggu ditolak (tanggal, tanpa komentar)
		r.Klaim.Adjustment["AcceptanceStatus"] = KeputusanTolak
		for _, a := range k.Tangga {
			if a.Keputusan != KeputusanMenunggu || (len(r.Tangga) > 0 && a.ID == r.Tangga[0].ID) {
				continue
			}
			r.Tangga = append(r.Tangga, UbahAnggota{ID: a.ID, Keputusan: KeputusanTolak, Tanggal: saat})
		}
	}
	if akhir { // S7.2.1.7 (tingkat akhir, setuju atau tolak)
		r.Klaim.Adjustment["IsApproved"] = kep.AcceptStatus
		r.Klaim.Header["AktifButton"] = "0"
	}
	if r.TingkatAkhirSetuju && adj["IsPrintAccept"] == "" { // S7.2.1.16 SaveAcceptation_KMT
		r.Tanda = true
		ob, it, ad := TandaAkseptasi(kl)
		for p, v := range ob {
			r.Klaim.Objek[p] = v
		}
		for p, v := range it {
			r.Klaim.Item[p] = v
		}
		for p, v := range ad {
			r.Klaim.Adjustment[p] = v
		}
	}
	r.SimpanOS = r.TingkatAkhirSetuju // S7.2.1.17 + S8 (sekali per KMT, prompt §5 butir 5)
	if akhir {                        // S12 (tingkat akhir): penanda cetak akseptasi
		if tolak {
			r.Klaim.Objek["IsPrintAccept"] = "1" // @if(AcceptStatus == "2", "1", "")
		} else {
			r.Klaim.Objek["IsPrintAccept"] = ""
		}
		if setuju { // S12.1.1 `.KomiteNo == KMT && .AcceptanceStatus == "1"`
			r.Klaim.Adjustment["IsPrintAccept"] = "1"
		}
	}
	// S13 (pre=false: setiap tingkat) usul penutupan / cadangan dari kepala kasus komite.
	r.Klaim.Header["ClaimData.IsCloseFile"] = usulTutup
	r.Klaim.Header["ClaimData.IsReservedClaim"] = usulCadang
	count := k.Count
	if tolak { // S14
		count = k.Loop
	}
	retro := adj["IsFacRetro"] == "1" // S7.2.1.1 Local.FacRetro (sebelum SaveAcceptation_KMT)
	if r.TingkatAkhirSetuju {
		r.JSONKlaim = true // S16
		if !retro {        // S17 / S19
			r.Konversi, r.LogAkseptasi = StsOSAkseptasi, true
		}
		r.Kasir = true          // S25 (penyaring DirectToKasir / StatusKasir / panjang nomor di layanan)
		r.PDF = StreamAkseptasi // S9 PrintPDFAccep_MultiAksep_KMT (kategori AcceptanceNote)
	}
	switch { // S20
	case setuju:
		r.StatusRiwayat = StatusRiwayatACC
	case tolak:
		r.StatusRiwayat = StatusRiwayatREJ
	}
	if count == k.Loop { // S22-S23 (`KomiteCount == KomiteLoop || AcceptStatus == "2"`)
		r.SubProgres = SubProgresDitolak
		if setuju {
			r.SubProgres = SubProgresDiterima
		}
	}
	switch { // 7.2.1.15 SendEmailKlaim_KMT S12 / S14 / S15
	case setuju && !akhir:
		r.Email = EmailPenyetujuBerikut
	case setuju:
		r.Email = EmailPembuatSetuju
	case tolak:
		r.Email = EmailPembuatTolak
	}
	count++ // S24 (pre=false: selalu)
	r.Selesai = !(setuju && count <= k.Loop)
	if !r.Selesai {
		r.Posisi = KomiteIDTingkat(k, count)
	}
	r.Kepala = Kepala{Count: count, Loop: k.Loop, AcceptStatus: kep.AcceptStatus, UsulTutup: usulTutup,
		UsulCadang: usulCadang}
	return r
}

// rencanaTutup = KomitePost_Reject (TT3) / KomitePost_CloseClaim (TT4): satu tingkat (`ReasClaimDeptHead`, KCF-03).
func rencanaTutup(k Kasus, kep Keputusan, akun string, saat time.Time) Rencana {
	r := Rencana{}
	setuju := kep.AcceptStatus == KeputusanSetuju
	jabatan := ""
	if i := k.barisBerjalan(); i >= 0 { // S6 / S5
		jabatan = k.Tangga[i].Jabatan
		r.Tangga = append(r.Tangga, UbahAnggota{ID: k.Tangga[i].ID, Keputusan: kep.AcceptStatus, Komentar: kep.Comment,
			IsiKomentar: true, Tanggal: saat, Pemutus: akun})
	}
	if t := teksRiwayat(kep.AcceptStatus, jabatan, k.ID); t != "" { // S3-S5 / S2-S4
		r.Klaim.Riwayat = append(r.Klaim.Riwayat, kontrak.RiwayatKlaimFacIn{Teks: t, Pelaku: akun, Tingkat: jabatan,
			Saat: saat})
	}
	r.JSONKlaim = true // S11 / S10 (tanpa syarat)
	r.Email = EmailPembuatTolak
	if setuju {
		r.Email = EmailPembuatSetuju
	}
	if k.TransferType == TransferReject {
		r.PDF = StreamTolak // S13.4-S13.9 "Reject Claim " + NoClaim + ".PDF", kategori CloseClaim (tanpa syarat)
		if setuju {
			r.OSTolak, r.KlaimDitolak, r.Konversi, r.TutupKlaim = true, true, StsOSTolak, kontrak.StatusKlaimDitolak
		}
	} else {
		r.PDF = StreamTutup // S11.4-S11.9 "Close Claim " + NoClaim + NoClaim + ".PDF", kategori CloseClaim
		if setuju {
			r.OSTutup, r.Konversi, r.TutupKlaim = true, StsOSFinal, kontrak.StatusKlaimSelesai
		}
	}
	count := k.Count + 1 // S16 / S13
	r.Selesai = !(setuju && count <= k.Loop)
	if !r.Selesai {
		r.Posisi = KomiteIDTingkat(k, count)
	}
	r.Kepala = Kepala{Count: count, Loop: k.Loop, AcceptStatus: kep.AcceptStatus, UsulTutup: UsulTidak,
		UsulCadang: UsulTidak}
	if k.UsulTutup != "" {
		r.Kepala.UsulTutup = k.UsulTutup
	}
	if k.UsulCadang != "" {
		r.Kepala.UsulCadang = k.UsulCadang
	}
	return r
}

// Stream PDF Pega (VERBATIM; tidak diekspor korpus Komite Claim FacIn - OQ-KCFI-01).
const (
	StreamAkseptasi = "AcceptanceNotePDF"
	StreamTolak     = "CommitteReject_CC"
	StreamTutup     = "CommitteCloseClaim"
)
