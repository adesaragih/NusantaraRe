package services

// Untuk apa berkas ini: LAYANAN - galat yang dipetakan handlers, daftar kerja penyetuju (worklist KomiteRouter), dan
// pembukaan satu kasus (flow action `ViewTransferDtl`: pra-proses `SetKomiteList_Act` + Section `ShowTransfer`).

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimprop/backend/models"
	"nusantarare/modul/komiteclaimprop/backend/repository"
)

// Galat yang dipetakan handlers ke kode HTTP.
var (
	ErrTanpaOracle        = db.ErrTanpaOracle
	ErrKasusTidakAda      = repository.ErrKasusTidakAda
	ErrPermintaanTidakSah = galat.ErrPermintaanTidakSah
	// ErrKeputusanBersamaan - kasus berubah sejak dibaca (409).
	ErrKeputusanBersamaan = repository.ErrKeputusanBersamaan
	// ErrKasusTertutup - kasus komite (atau klaim induknya) sudah Resolved-Completed (409).
	ErrKasusTertutup = errors.New("services: kasus komite sudah selesai")
	// ErrBukanPemegang - pelaku bukan pemilik baris tangga tingkat berjalan (KomiteRouter S6.1, keputusan 30).
	ErrBukanPemegang = fmt.Errorf("%w: bukan penyetuju tingkat berjalan kasus komite ini", inti.ErrTanpaWewenang)
	// ErrAksiTertutup - tombol Submit nonaktif (`pyDisabledWhen pyWorkPage.Adjustment.AcceptedNo != ''`) (409).
	ErrAksiTertutup = errors.New("services: aksi ini tidak tersedia di layar kasus komite saat ini")
	// ErrKlaimIndukTidakAda - kasus komite menunjuk klaim induk / baris adjustment yang tidak ada (409; butir terbuka 14:
	// menolak terang, bukan diam seperti Pega).
	ErrKlaimIndukTidakAda = errors.New("services: klaim induk atau baris adjustment kasus komite ini tidak ditemukan")
	// ErrKlaimBelumDisambung - kontrak Claim Prop tidak pernah disambung (salah rakit).
	ErrKlaimBelumDisambung = errors.New(
		"services: kontrak Claim Prop untuk Komite belum disambung (inti/backend/kontrak.KlaimTreatyKomite)")
)

// GalatValidasi - pesan wajib isi layar (pyRequired / pyRequiredWhen). Handlers menjawabnya 422.
type GalatValidasi struct{ Pesan []string }

func (e *GalatValidasi) Error() string {
	return "services: validasi layar komite gagal: " + strings.Join(e.Pesan, "; ")
}

// Layanan adalah pintu aturan dagang Komite Claim Prop.
type Layanan struct {
	g        Gudang
	a        Acuan
	klaim    kontrak.KlaimTreatyKomite
	jam      func() time.Time
	produksi bool
	kasir    models.KonfigurasiKasir
	surel    models.KonfigurasiEmail
}

// Baru menyusun layanan; `g` nil = tanpa Oracle (setiap aksi 503).
func Baru(g Gudang, a Acuan, k kontrak.KlaimTreatyKomite, jam func() time.Time, produksi bool) *Layanan {
	if jam == nil {
		jam = time.Now
	}
	return &Layanan{g: g, a: a, klaim: k, jam: jam, produksi: produksi}
}

// DenganEmail menyetel akun notifikasi dan CC email komite (konfigurasi berdokumen, MODUL.md).
func (l *Layanan) DenganEmail(c models.KonfigurasiEmail) *Layanan {
	salin := *l
	salin.surel = c
	return &salin
}

// DenganKasir menyetel kode tetap muatan Kasir (konfigurasi berdokumen, MODUL.md).
func (l *Layanan) DenganKasir(c models.KonfigurasiKasir) *Layanan {
	salin := *l
	salin.kasir = c
	return &salin
}

// AdaGudang - layanan tersambung ke penyimpanan.
func (l *Layanan) AdaGudang() bool { return l != nil && l.g != nil && l.a != nil }

func (l *Layanan) siap(p inti.Pelaku) error {
	if !l.AdaGudang() {
		return ErrTanpaOracle
	}
	if l.klaim == nil {
		return ErrKlaimBelumDisambung
	}
	return inti.WajibIdentitas(p)
}

// galatKontrak menerjemahkan galat kontrak Claim Prop.
func galatKontrak(err error) error {
	switch {
	case errors.Is(err, kontrak.ErrKlaimTreatyTertutup):
		return fmt.Errorf("%w: kasus klaim induk sudah ditutup", ErrKasusTertutup)
	case errors.Is(err, kontrak.ErrKlaimTreatyTidakAda), errors.Is(err, kontrak.ErrAdjustmentTreatyTidakAda):
		return fmt.Errorf("%w: %v", ErrKlaimIndukTidakAda, err)
	}
	return err
}

// DaftarKerja - worklist KomiteRouter pelaku (halaman awal).
func (l *Layanan) DaftarKerja(ctx context.Context, p inti.Pelaku) ([]models.BarisKerja, error) {
	if err := l.siap(p); err != nil {
		return nil, err
	}
	d, err := l.g.DaftarKerja(ctx, p.AkunID)
	for i := range d {
		d[i].StatusBaris = models.LabelStatusBaris(d[i].StatusBaris)
	}
	return d, err
}

// muat - kasus + tangga + klaim induk (tx nil = layar).
func (l *Layanan) muat(ctx context.Context, tx *db.Tx, id string, kunci bool) (models.Kasus, kontrak.KlaimTreaty, error) {
	if !strings.HasPrefix(id, models.AwalanKomite) {
		return models.Kasus{}, kontrak.KlaimTreaty{}, fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
	}
	k, err := l.g.BacaKasus(ctx, tx, id, kunci)
	if err != nil {
		return models.Kasus{}, kontrak.KlaimTreaty{}, err
	}
	if k.Tangga, err = l.g.BacaTangga(ctx, tx, id); err != nil {
		return models.Kasus{}, kontrak.KlaimTreaty{}, err
	}
	if kunci && !k.Tertutup() {
		if err := l.klaim.KunciKlaimTreaty(ctx, tx, k.KlaimID); err != nil { // S4
			return models.Kasus{}, kontrak.KlaimTreaty{}, galatKontrak(err)
		}
	}
	kl, err := l.klaim.BacaKlaimTreaty(ctx, tx, k.KlaimID, k.AdjustmentID)
	if err != nil {
		return models.Kasus{}, kontrak.KlaimTreaty{}, galatKontrak(err)
	}
	return k, kl, nil
}

// BukaKasus = flow action `ViewTransferDtl`: pra-proses `SetKomiteList_Act` (total per mata uang, TT 2) lalu Section
// `ShowTransfer`. Siapa pun yang berhak membuka modul ini boleh melihat; hanya pemegang yang boleh Submit.
func (l *Layanan) BukaKasus(ctx context.Context, p inti.Pelaku, id string) (models.Layar, error) {
	if err := l.siap(p); err != nil {
		return models.Layar{}, err
	}
	k, kl, err := l.muat(ctx, nil, id, false)
	if err != nil {
		return models.Layar{}, err
	}
	if k.Count == 1 && models.AdjustmentKlaim(kl)["IsSubjectivity"] == "true" { // S16 kirim ulang subjectivity
		if k.KomentarAwal, err = l.g.KomentarAwal(ctx, k.KlaimID, k.AdjustmentID, k.ID); err != nil {
			return models.Layar{}, err
		}
	}
	total, err := models.TotalPenyesuaian(kl.Daftar[models.DaftarAdjustment])
	if err != nil {
		return models.Layar{}, err
	}
	return models.SusunLayar(k, kl, total, p.AkunID), nil
}
