package services

// Untuk apa berkas ini: PENUTUPAN DAN PENOLAKAN - tombol "Close Claim" layar Input Adjustment (local action
// PreventRejectClaim -> `CloseClaim`) dan tombol "Reject Claim" layar Input Register (local action RejectSurveyClaim,
// pra-proses `SetRejectClaim_pre`). Jalur komite TT3 / TT4 = OQ-CFI-27 (models/tutup.go).

import (
	"nusantarare/modul/claimfacin/backend/models"
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
