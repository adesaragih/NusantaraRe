package services

// Untuk apa berkas ini: LAYANAN - galat yang dipetakan handlers, daftar kerja penyetuju (worklist KomiteRouter, tabel
// komite di bawah inbox Claim Fac In), dan pembukaan satu kasus (flow action `ViewTransferDtl`: pra-proses
// `SetValueKomite` + `ApprovalKomite_Act` lalu Section `ShowTransfer`). Pola disalin dari Komite Claim Prop / Non Prop.

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
	"nusantarare/modul/komiteclaimfacin/backend/models"
	"nusantarare/modul/komiteclaimfacin/backend/repository"
)

// Galat yang dipetakan handlers ke kode HTTP.
var (
	ErrTanpaOracle        = db.ErrTanpaOracle
	ErrKasusTidakAda      = repository.ErrKasusTidakAda
	ErrPermintaanTidakSah = galat.ErrPermintaanTidakSah
	// ErrKeputusanBersamaan - kasus berubah sejak dibaca (409).
	ErrKeputusanBersamaan = repository.ErrKeputusanBersamaan
	// ErrKasusTertutup - kasus komite (atau klaim induknya) sudah selesai, atau adjustment-nya sudah diputus (409).
	ErrKasusTertutup = errors.New("services: kasus komite sudah selesai")
	// ErrBukanPemegang - pelaku bukan anggota workbasket / akun tingkat berjalan (KCF-01; SetProteksiSubmiteKomite
	// diperbaiki prompt §5 butir 1).
	ErrBukanPemegang = fmt.Errorf("%w: bukan penyetuju tingkat berjalan kasus komite ini", inti.ErrTanpaWewenang)
	// ErrKlaimIndukTidakAda - kasus komite menunjuk klaim induk / baris adjustment yang tidak ada (409).
	ErrKlaimIndukTidakAda = errors.New("services: klaim induk atau baris adjustment kasus komite ini tidak ditemukan")
	// ErrKlaimBelumDisambung - kontrak Claim Fac In tidak pernah disambung (salah rakit).
	ErrKlaimBelumDisambung = errors.New(
		"services: kontrak Claim Fac In untuk Komite belum disambung (inti/backend/kontrak.KlaimFacInKomite)")
)

// GalatValidasi - pesan wajib isi layar (LS45 `REQ=true`). Handlers menjawabnya 422.
type GalatValidasi struct{ Pesan []string }

func (e *GalatValidasi) Error() string {
	return "services: validasi layar komite gagal: " + strings.Join(e.Pesan, "; ")
}

// Layanan adalah pintu aturan dagang Komite Claim Fac In.
type Layanan struct {
	g        Gudang
	a        Acuan
	klaim    kontrak.KlaimFacInKomite
	jam      func() time.Time
	produksi bool
	kasir    models.KonfigurasiKasir
	surel    models.KonfigurasiEmail
	berkas   PenyimpananBerkas
}

// Baru menyusun layanan; `g` nil = tanpa Oracle (setiap aksi 503).
func Baru(g Gudang, a Acuan, k kontrak.KlaimFacInKomite, jam func() time.Time, produksi bool) *Layanan {
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

// DenganPenyimpanan memasang penyimpanan berkas dokumen akseptasi (Oracle: `penyimpanan.Oracle`; uji: tiruan).
func (l *Layanan) DenganPenyimpanan(p PenyimpananBerkas) *Layanan {
	salin := *l
	salin.berkas = p
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

// galatKontrak menerjemahkan galat kontrak Claim Fac In.
func galatKontrak(err error) error {
	switch {
	case errors.Is(err, kontrak.ErrKlaimFacInTertutup):
		return fmt.Errorf("%w: kasus klaim induk sudah ditutup", ErrKasusTertutup)
	case errors.Is(err, kontrak.ErrKlaimFacInTidakAda), errors.Is(err, kontrak.ErrAdjustmentFacInTidakAda):
		return fmt.Errorf("%w: %v", ErrKlaimIndukTidakAda, err)
	}
	return err
}

// DaftarKerja - worklist KomiteRouter pelaku (tabel komite inbox Claim Fac In).
func (l *Layanan) DaftarKerja(ctx context.Context, p inti.Pelaku) ([]models.BarisKerja, error) {
	if err := l.siap(p); err != nil {
		return nil, err
	}
	d, err := l.g.DaftarKerja(ctx, p.AkunID, models.PeranKerja(p.Peran))
	for i := range d {
		d[i].StatusBaris = models.LabelStatusBaris(d[i].StatusBaris)
	}
	return d, err
}

// muat - kasus + tangga + klaim induk (tx nil = layar).
func (l *Layanan) muat(ctx context.Context, tx *db.Tx, id string, kunci bool) (models.Kasus, kontrak.KlaimFacIn, error) {
	if !strings.HasPrefix(id, models.AwalanKomite) {
		return models.Kasus{}, kontrak.KlaimFacIn{}, fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
	}
	k, err := l.g.BacaKasus(ctx, tx, id, kunci)
	if err != nil {
		return models.Kasus{}, kontrak.KlaimFacIn{}, err
	}
	if k.Tangga, err = l.g.BacaTangga(ctx, tx, id); err != nil {
		return models.Kasus{}, kontrak.KlaimFacIn{}, err
	}
	if kunci && !k.Tertutup() {
		if err := l.klaim.KunciKlaimFacIn(ctx, tx, k.KlaimID); err != nil { // KomitePost_* S1 / S2
			return models.Kasus{}, kontrak.KlaimFacIn{}, galatKontrak(err)
		}
	}
	kl, err := l.klaim.BacaKlaimFacIn(ctx, tx, k.KlaimID, k.AdjustmentID)
	if err != nil {
		return models.Kasus{}, kontrak.KlaimFacIn{}, galatKontrak(err)
	}
	return k, kl, nil
}

// praProses = SetValueKomite + ApprovalKomite_Act S1-S3; kurs baris tanpa `.CurrencyDol` dari CurrencyStandard.
func (l *Layanan) praProses(ctx context.Context, k models.Kasus, kl kontrak.KlaimFacIn) (models.PraProses, error) {
	kurs := map[string]string{}
	if k.TransferType == models.TransferAdjustment {
		for o := range kl.Daftar[models.DaftarObjek] {
			for i := range kl.Daftar[models.DaftarItem(o+1)] {
				for _, b := range kl.Daftar[models.DaftarAdj(o+1, i+1)] {
					cur := b["CurrencyID"]
					if b["CurrencyDol"] != "" || cur == "" {
						continue
					}
					if _, ada := kurs[cur]; ada {
						continue
					}
					v, err := l.a.KursStandar(ctx, cur)
					if err != nil {
						return models.PraProses{}, err
					}
					kurs[cur] = v
				}
			}
		}
	}
	return models.SetValueKomite(k, kl, kurs)
}

// perluasan = ApprovalKomite_Act S4-S8 (KCF-02): anggota tangga baru bagi kasus tingkat 1.
func (l *Layanan) perluasan(ctx context.Context, k models.Kasus, pr models.PraProses, peran []string) ([]models.Anggota,
	error) {
	if k.TransferType != models.TransferAdjustment || k.Count != 1 || len(k.Tangga) != 1 || pr.Retro {
		return nil, nil
	}
	roster, err := l.a.RosterKomite(ctx)
	if err != nil {
		return nil, err
	}
	return models.PerluasTangga(k, pr, roster, models.AnggotaSPVB(peran))
}

// BukaKasus = flow action `ViewTransferDtl`: pra-proses lalu Section `ShowTransfer`. Perluasan tangga tingkat 1
// DITAMPILKAN (dihitung untuk pelaku yang membuka; KCF-02 "tanpa menulis saat GET"). Siapa pun yang berhak membuka modul
// ini boleh melihat; hanya pemegang yang boleh Submit.
func (l *Layanan) BukaKasus(ctx context.Context, p inti.Pelaku, id string) (models.Layar, error) {
	if err := l.siap(p); err != nil {
		return models.Layar{}, err
	}
	k, kl, err := l.muat(ctx, nil, id, false)
	if err != nil {
		return models.Layar{}, err
	}
	pr, err := l.praProses(ctx, k, kl)
	if err != nil {
		return models.Layar{}, err
	}
	var baru []models.Anggota
	if !k.Tertutup() {
		peran := []string(nil)
		if k.Pemegang(p.AkunID, p.Peran) { // pita SPV B hanya untuk pemutus tingkat 1 yang membuka
			peran = p.Peran
		}
		if baru, err = l.perluasan(ctx, k, pr, peran); err != nil {
			return models.Layar{}, err
		}
	}
	jenis := map[string]string{}
	if k.TransferType == models.TransferAdjustment {
		if jenis, err = l.a.NamaJenisReas(ctx); err != nil {
			return models.Layar{}, err
		}
	}
	return models.SusunLayar(k, kl, pr, baru, jenis, p.AkunID, p.Peran), nil
}
