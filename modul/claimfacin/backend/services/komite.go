package services

// Untuk apa berkas ini: PENYERAHAN KE KOMITE TT2 - tombol "Send to Committe" panel InputAdjustment (`SendPICProtect_Act`
// + harness Comittee ber-pra-proses `SetRemarksKomite`) dan tombol "Send Claim to Committee" pop-up Comittee
// (`DraftGenerateDLAFacin_Act` + `CreateKMTNo_Act`). Kelahiran kasus komite lewat `repository/komite.go`.
//
// Gerbang `Protect.CARI1 = 1 And Protect.CARI2 = 1` DITEGAKKAN di server: dihitung ulang sebelum tata aksi
// "KirimKomite" dievaluasi (`praTata`), bukan diambil dari kiriman layar (pola Claim Prop).

import (
	"context"
	"fmt"
	"strings"

	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// kunciModalKomite - kunci modal pop-up Comittee adjustment (o, i, a).
func kunciModalKomite(o, i, a int) string {
	return models.KunciPanel(models.ModalKomite, models.DaftarAdj(o, i), a)
}

// proteksiKomite = SendPICProtect_Act (+ CekPremiLunas_Act 19): penanda Protect.CARI1 / CARI2 dan pesannya.
func (l *Layanan) proteksiKomite(ctx context.Context, h *models.Halaman, o, i, a int) error {
	lim, _, err := l.a.LimitDirekturUtama(ctx) // 5
	if err != nil {
		return err
	}
	kat, err := l.g.KategoriLampiran(ctx, h.Ambil("pyID")) // 3 AttachCategory.pxResults
	if err != nil {
		return err
	}
	lolos, err := models.ProteksiKomite(h, o, i, a, models.CacahLampiran(kat), lim)
	if err != nil {
		return err
	}
	b, err := models.Adj(h, o, i, a)
	if err != nil {
		return err
	}
	premi := true
	if models.PerluCekPremi(b) { // 19.1
		inv, cur, nopol := models.KunciPremi(h, b)
		saldo, err := l.a.SaldoPremi(ctx, inv, cur) // CekPremiLunas_Act 4
		if err != nil {
			return err
		}
		belum, err := models.PremiBelumLunas(saldo) // 5-6
		if err != nil {
			return err
		}
		if belum {
			buka, err := l.a.AdaProteksiPremi(ctx, nopol) // 7
			if err != nil {
				return err
			}
			if !buka {
				premi = false
				models.PesanPremi(h) // 8
			}
		}
	}
	h.Setel(models.JalurProtect1, nolSatu(lolos))
	h.Setel(models.JalurProtect2, nolSatu(premi)) // 19.2
	return nil
}

func nolSatu(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// aksiBukaKomite - tombol "Send to Committe" (NA `.IsKomite = 1 || .SpreadingAdjustment(1).TreatyType = ”`):
// SendPICProtect_Act, lalu - hanya bila CARI1 dan CARI2 = 1 - pop-up Comittee (pra-proses SetRemarksKomite). Pesan
// proteksi ditampilkan tanpa membatalkan aksi (penanda Protect tetap dibawa layar).
func aksiBukaKomite(j *jalanAksi) error {
	if err := j.l.proteksiKomite(j.ctx, j.h, j.o, j.i, j.e); err != nil {
		return err
	}
	if j.h.Ambil(models.JalurProtect1) != "1" || j.h.Ambil(models.JalurProtect2) != "1" {
		return nil
	}
	if err := models.SetRemarksKomite(j.k, j.h, j.o, j.i, j.e); err != nil {
		return err
	}
	j.bukaModal = kunciModalKomite(j.o, j.i, j.e)
	return nil
}

// aksiKirimKomite - tombol "Send Claim to Committee" pop-up Comittee:
//
//	gerbang     Protect.CARI1 / CARI2 (dihitung ulang, `praTata`) dan tombol tampil pada isian sesudah digabung
//	DraftGenerateDLAFacin_Act   models.DraftDLA (adjustment retro)
//	CreateKMTNo_Act 8-10        tangga = roster calon SetListKomite_act; T_WORK_CLAIM KMT- + T_GENERAL_KOMITE + tangga
//	CreateKMTNo_Act 12-17       models.TandaiKirimKomite; InsertProgressClaim "Auto Create" / "Waiting Committee"
//	CreateKMTNo_Act 14          SendEmailKlaim (outbox surel, hanya produksi)
func aksiKirimKomite(j *jalanAksi) error {
	h, o, i, a := j.h, j.o, j.i, j.e
	// pembuka "Send to Committe" (NA IsKomite = 1 atau Spreading Adjustment kosong) - aksi ini melewati aksiTerbuka
	if !pembukaTerbuka(j, models.KunciPanel(models.PanelAdj, models.DaftarAdj(o, i), a), "SendPICProtect", 0) {
		return fmt.Errorf("%w: Send to Committe adjustment %d tertutup", ErrAksiTertutup, a)
	}
	if h.Ambil(models.JalurProtect1) != "1" || h.Ambil(models.JalurProtect2) != "1" {
		if err := validasi(h); err != nil {
			return err
		}
		return fmt.Errorf("%w: penyerahan komite ditolak (lampiran / rekening / premi)", ErrAksiTertutup)
	}
	ts := models.Evaluasi(h, models.LayarKomite(o, i, a, func(*models.Halaman) bool { return true }), false)
	if !models.AksiTerbuka(ts, "KirimKomite", 0) {
		if err := wajibTerisi(j, ts); err != nil {
			return err
		}
		return fmt.Errorf("%w: tombol Send Claim to Committee tidak tampil untuk isian ini", ErrAksiTertutup)
	}
	if err := wajibTerisi(j, ts); err != nil {
		return err
	}
	b, err := models.Adj(h, o, i, a)
	if err != nil {
		return err
	}
	if b["IsKomite"] == "1" || b[models.PropKomiteID] != "" {
		return ErrAdjustmentDiKomite
	}
	if err := models.DraftDLA(j.k, h, o, i, a); err != nil {
		return err
	}
	roster, err := j.l.a.RosterKomite(j.ctx)
	if err != nil {
		return err
	}
	calon, err := models.RosterKomiteCalon(h, o, i, a, roster) // CreateKMTNo_Act 8 SetListKomite_act
	if err != nil {
		return err
	}
	if len(calon) == 0 {
		return &GalatValidasi{Pesan: []string{"Roster komite (EMAILKOMITE STS_KLAIM FACIN) kosong untuk nilai adjustment ini"}}
	}
	if b[models.PropID] == "" { // baris baru di aksi ini: ID stabil dulu
		models.BuangTurunan(h)
		if err := j.l.g.SimpanHalaman(j.ctx, j.tx, j.kasus.ID, h); err != nil {
			return err
		}
	}
	anggota := make([]repository.AnggotaTangga, 0, len(calon))
	for n, r := range calon {
		anggota = append(anggota, repository.AnggotaTangga{Urut: n + 1, OperatorID: r.OperatorID, Jabatan: r.Jabatan,
			Email: r.Email})
	}
	nama, err := j.l.a.NamaPelaku(j.ctx, j.k.Pelaku)
	if err != nil {
		return err
	}
	kmt, err := j.l.g.BuatKasusKomite(j.ctx, j.tx, j.kasus.ID, b[models.PropID], j.k.Pelaku, nama, anggota,
		j.k.Sekarang) // 10 pxAddChildWork
	if err != nil {
		return err
	}
	if err := j.l.g.SetelKomiteAdjustment(j.ctx, j.tx, b[models.PropID], kmt, ""); err != nil {
		return err
	}
	if err := models.TandaiKirimKomite(j.k, h, o, i, a, kmt); err != nil { // 12-15, 17
		return err
	}
	if err := j.l.tulisProgres(j.ctx, j.tx, j.k, j.kasus.ID, h, models.ParamProgres{CaseID: kmt, // 16
		Progres1: models.PesanPembayaran(b["PaymentType"]), Progres2: models.ProgresMenungguKmt,
		Comment: models.ProgresAutoCreate}); err != nil {
		return err
	}
	j.info = models.OQDokumenPDF
	return j.antre(JenisEfekEmail, kmt, map[string]string{"klaim": j.kasus.ID, "komite": kmt}) // 14
}

// praTata - perhitungan server sebelum tata aksi dievaluasi (gerbang yang tidak boleh diambil dari kiriman layar).
var praTata = map[string]func(j *jalanAksi) error{
	"KirimKomite": func(j *jalanAksi) error {
		o, i, a, ok := konteksAdj(j.r.Konteks)
		if !ok || a == 0 || !strings.HasPrefix(j.r.Konteks, models.ModalKomite+":") {
			return fmt.Errorf("%w: konteks %q", ErrPermintaanTidakSah, j.r.Konteks)
		}
		return j.l.proteksiKomite(j.ctx, j.h, o, i, a)
	},
}

// aksiPascaGabung - aksi yang keterbukaannya diperiksa penangannya SESUDAH isian layar digabung (tombol bergantung pada
// medan yang baru diketik di pop-up yang sama).
var aksiPascaGabung = map[string]bool{"KirimKomite": true}
