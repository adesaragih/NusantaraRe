package services

// Untuk apa berkas ini: PENUTUPAN DAN PENOLAKAN - tombol "Close Claim" layar Input Adjustment (local action
// PreventRejectClaim -> `CloseClaim` atau `SendCloseClaimToKomite` TT4) dan tombol "Reject Claim" layar Input Register
// (local action RejectSurveyClaim, pra-proses `SetRejectClaim_pre`; "Yes" SureRejectClaim -> `SendRejectClaimToKomite2`
// TT3). Kelahiran kasus komite TT3 / TT4 lewat `repository/komite.go` (keputusan work owner 10-10-2026 KCF-03).

import (
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// Kunci modal penutupan / penolakan.
const (
	ModalTutup = "tutup" // PreventRejectClaim
	ModalTolak = "tolak" // RejectSurveyClaim (ClaimComiteeReject + SureRejectClaim_section)
)

// aksiTutupKlaim = "Yes" PreventRejectClaim (AllocationShareSalvage false) -> `CloseClaim`:
//
//	1       ValidationAdjustmentKomite (models.ValidasiTutup; lampiran "CloseClaim" = cacah dokumen klaim)
//	4-5     Remark / Remark_Close, proteksi "Please Print DLA Before Close This Claim"
//	6       PERBAIKAN prompt §5 butir 8: setiap pesan 1-5 MENGHENTIKAN penutupan (Pega membaca Param.ErrorMsg dari
//	        activity terpanggil tanpa halaman parameternya - pesan tampil tetapi klaim tetap ditutup)
//	7       UpdateTotalJob_sql: `ClaimData.UserTeknis` tanpa penulis (medan LS177 InputRegisterDetail VIS 1=2) -
//	        tidak pernah berjalan, tidak dibangun
//	8       OS_AKSEPTASI_KLAIM STS 4 (models.BarisOSTutup) + KonversiKlaim_Act STS 4 (efek hanya produksi)
//	9-11    kronologi "Finish Adjustment (Close Claim)", JSON_KLAIM
//	12      ASMForceCaseClose (PEGA_JSON_KLAIM_PNC StsSimpan = 1 pada setiap keberhasilan) -> Resolved-Completed
func aksiTutupKlaim(j *jalanAksi) error {
	h := j.h
	if err := wajibTerisi(j, models.Evaluasi(h, models.LayarTutup(), false)); err != nil {
		return err
	}
	kat, err := j.l.g.KategoriLampiran(j.ctx, j.kasus.ID)
	if err != nil {
		return err
	}
	if err := models.ValidasiTutup(h, models.CacahLampiran(kat)[models.KategoriTutupKlaim]); err != nil { // 1
		return err
	}
	models.CekDLATutup(h)               // 4-5
	if err := validasi(h); err != nil { // 6
		return err
	}
	if err := j.l.g.SisipOS(j.ctx, j.tx, models.BarisOSTutup(h, j.kasus.ID), j.k.Sekarang); err != nil { // 8.3
		return err
	}
	if err := j.antreKonversi(models.StsOSTutup); err != nil { // 8.4
		return err
	}
	models.SelesaiTutup(j.k, h)                // 9-10
	if err := j.salinJSONKlaim(); err != nil { // 11
		return err
	}
	models.BuangTurunan(h) // 12
	if err := j.l.g.SimpanHalaman(j.ctx, j.tx, j.kasus.ID, h); err != nil {
		return err
	}
	if err := j.l.g.TutupKasus(j.ctx, j.tx, j.kasus.ID, j.kasus.Tahap, j.k.Sekarang); err != nil {
		return err
	}
	// 12 ASMForceCaseClose CloseAllSubCases=true: kasus komite KMT- klaim ini yang masih menunggu ikut ditutup
	if err := j.l.g.TutupKomiteAnak(j.ctx, j.tx, j.kasus.ID, "", models.StatusSelesai, j.k.Sekarang); err != nil {
		return err
	}
	j.selesai = true
	return nil
}

// aksiBukaTolak - tombol "Reject Claim" layar Input Register: pra-proses SetRejectClaim_pre lalu pop-up
// ClaimComiteeReject.
func aksiBukaTolak(j *jalanAksi) error {
	models.SetRejectClaimPre(j.k, j.h, j.kasus.PembuatNama)
	j.bukaModal = ModalTolak
	return nil
}

// aksiKirimTolakKomite = "Yes" SureRejectClaim_section -> `SendRejectClaimToKomite2` (TT3).
func aksiKirimTolakKomite(j *jalanAksi) error { return kirimKomiteTutup(j, models.TransferTolak) }

// aksiKirimTutupKomite = "Yes" PreventRejectClaim (AllocationShareSalvage true) -> `SendCloseClaimToKomite` (TT4).
func aksiKirimTutupKomite(j *jalanAksi) error { return kirimKomiteTutup(j, models.TransferTutup) }

// kirimKomiteTutup = SendRejectClaimToKomite2 / SendCloseClaimToKomite:
//
//	2        Remark / Remark_Close := Remarks pop-up (models.SalinCatatanTutup)
//	4-6      TT3: sudah ada akseptasi / estimasi belum Face Claim; TT4: adjustment belum diputus - pesan VERBATIM,
//	         kasus komite tidak lahir
//	7.1-7.4  kasus komite KMT- tanpa adjustment, TransferType 3 / 4, tangga satu tingkat ReasClaimDeptHead (KCF-03;
//	         7.2 / 7.3 akun + email orang tertulis mati diganti workbasket, prompt §5 butir 7)
//	7.5-7.7  Komite.CARI1, kronologi "Request Reject claim " / "Request close claim without payment " + KMT
//	7.9      SendEmailKlaimRejectClose (outbox, hanya produksi; CC / BCC orang tidak disalin)
//
// ClaimComitee klaim induk (7.2 / 7.5) tidak ditulis: tanpa kolom (tangga dibaca dari tabel komite). Penjaga ganda
// `models.PesanTutupKomiteGanda` = penyimpangan sadar (PARITAS).
func kirimKomiteTutup(j *jalanAksi, transfer string) error {
	h := j.h
	layar := models.LayarTolak()
	if transfer == models.TransferTutup {
		layar = models.LayarTutup()
	}
	if err := wajibTerisi(j, models.Evaluasi(h, layar, false)); err != nil {
		return err
	}
	models.SalinCatatanTutup(h) // 2
	galat := models.PeriksaTolakKomite(h)
	if transfer == models.TransferTutup {
		galat = models.PeriksaTutupKomite(h)
	}
	if galat != "" { // 6 Page-Set-Messages
		h.TambahPesan("", galat)
		return validasi(h)
	}
	ada, err := j.l.g.AdaKomiteTutupTerbuka(j.ctx, j.tx, j.kasus.ID)
	if err != nil {
		return err
	}
	if ada {
		h.TambahPesan("", models.PesanTutupKomiteGanda)
		return validasi(h)
	}
	nama, err := j.l.a.NamaPelaku(j.ctx, j.k.Pelaku)
	if err != nil {
		return err
	}
	anggota := []repository.AnggotaTangga{{Urut: 1, OperatorID: models.WorkbasketTutupKomite,
		Jabatan: models.JabatanTutupKomite}}
	kmt, err := j.l.g.BuatKasusKomite(j.ctx, j.tx, j.kasus.ID, "", transfer, j.k.Pelaku, nama, anggota,
		j.k.Sekarang) // 7.4 pxAddChildWork
	if err != nil {
		return err
	}
	models.TandaiKirimTutup(j.k, h, transfer, kmt) // 7.5-7.7
	if transfer == models.TransferTolak {
		j.bukaModal = ModalTolak
	} else {
		j.bukaModal = ModalTutup
	}
	return j.antre(JenisEfekEmail, kmt, map[string]string{"klaim": j.kasus.ID, "komite": kmt}) // 7.9
}
