package services

// Antre-ulang efek keluar dan pekerjanya - butir aq, A2.
//
// Untuk apa berkas ini: menutup stub `AntreanBelumDiputuskan`. Kegagalan efek
// keluar disimpan ke `T_EFEK_KELUAR` dengan jadwal percobaan berikutnya, dan
// seorang pekerja memungutnya kembali satu per satu.
//
// ⛔ Kegagalan PERMANEN tidak pernah dijadwalkan ulang. `LayakDicobaUlang`
// sudah memisahkannya; yang berubah di sini hanya akibatnya - `gagal-permanen`
// tanpa jadwal. Antrean yang menjadwalkan ulang kegagalan konfigurasi akan
// memutari keadaan yang tidak mungkin berubah karena dicoba lagi.
//
// Dibaca sesudah: efekkeluar.go.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"nusantarare/internal/repository"
)

// ErrPekerjaTanpaEfek diteruskan apa adanya dari repository: tidak ada yang
// jatuh tempo. Keadaan NORMAL.
var ErrPekerjaTanpaEfek = repository.ErrEfekTidakAda

// Modul dan lini yang menandai baris outbox milik modul ini.
//
// ⚠️ Outbox-nya LINTAS MODUL. Tanpa kedua penanda ini, pekerja Komite dan
// pekerja Claim Life akan saling memungut pekerjaan yang bukan miliknya.
const (
	LiniLife        = "LIFE"
	ModulClaimLife  = "CLAIMLIFE"
	backoffAwal     = 30 * time.Second
	backoffMaksimum = 6 * time.Hour
	// percobaanMaksimum - sesudahnya efek menyerah dan menunggu manusia.
	//
	// ⚠️ Ada BATASnya, dan itu disengaja. Antrean tanpa batas percobaan
	// menyimpan baris yang dicoba selamanya, dan kegagalan sungguhan tidak
	// pernah terlihat sebagai kegagalan oleh siapa pun.
	percobaanMaksimum = 8
)

// Backoff menghitung jeda sebelum percobaan ke-(n+1).
//
// ⛔ Berlipat dua dan BERPLAFON. Tanpa plafon, percobaan ke-20 dijadwalkan
// bertahun-tahun ke depan - yang sama saja dengan hilang.
//
// ⚠️ Deterministik, tanpa jitter. Jitter berguna melawan kawanan pekerja yang
// bangun serentak; di sini pekerjanya satu, dan ketidakpastian membuat
// perilakunya tidak dapat diuji. Bila kelak ada banyak pekerja, jitter
// ditambahkan sebagai keputusan tersendiri, bukan diselundupkan sekarang.
func Backoff(percobaan int) time.Duration {
	if percobaan < 1 {
		percobaan = 1
	}
	jeda := backoffAwal
	for i := 1; i < percobaan; i++ {
		jeda *= 2
		if jeda >= backoffMaksimum {
			return backoffMaksimum
		}
	}
	return jeda
}

// muatanOutbox adalah bentuk JSON yang disimpan di kolom `MUATAN`.
//
// ⛔ PENGENAL DAN WAKTU SAJA. Nol nama orang, nol email, nol alamat, nol
// kredensial - kolom ini dibaca lagi entah kapan oleh entah siapa. Penjaga
// statik menegakkan bentuknya.
type muatanOutbox struct {
	KlaimID      string `json:"klaim_id"`
	AdjustmentID string `json:"adjustment_id"`
	AkunID       string `json:"akun_id"`
	Waktu        string `json:"waktu"`
}

// antreanOracle menyimpan kegagalan efek keluar ke outbox.
//
// ⚠️ TANPA `*repository.Tx`, dan itu bukan kelalaian: kegagalan efek keluar
// terjadi SESUDAH transaksi klaimnya selesai. Ia membuka transaksinya sendiri,
// yang pendek dan berdiri sendiri.
type antreanOracle struct {
	svc   *Service
	pohon *repository.PohonKlaim
}

// AntreanEfekOracle menyusun antre-ulang yang menulis ke `T_EFEK_KELUAR`.
func AntreanEfekOracle(svc *Service) Antrean {
	return antreanOracle{svc: svc, pohon: repository.NewPohonKlaim(svc.db)}
}

// Antre menyimpan satu kegagalan untuk dicoba lagi - atau untuk ditunggu
// manusia, bila ia tidak layak dicoba ulang.
func (a antreanOracle) Antre(ctx context.Context, c CatatanEfekGagal) error {
	if !a.svc.PunyaDatabase() {
		return repository.ErrTanpaOracle
	}
	muatan, err := json.Marshal(muatanOutbox{
		KlaimID:      c.KlaimID,
		AdjustmentID: c.AdjustmentID,
		AkunID:       c.AkunID,
		Waktu:        c.Waktu.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("services: merakit muatan outbox: %w", err)
	}
	// ⛔ Rujukannya baris adjustment bila ada, klaim bila tidak. Baris outbox
	// tanpa rujukan tidak dapat ditelusuri balik ke pekerjaan yang gagal.
	rujukan := c.AdjustmentID
	if rujukan == "" {
		rujukan = c.KlaimID
	}
	return a.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		id, err := a.pohon.AntreEfek(ctx, tx, LiniLife, ModulClaimLife,
			c.Nama, rujukan, string(muatan), c.Waktu)
		if err != nil {
			return err
		}
		if c.LayakUlang {
			// Percobaan pertama dijadwalkan; `AntreEfek` menulisnya jatuh
			// tempo seketika, jadi jadwalnya digeser di sini.
			return a.pohon.TuntaskanEfek(ctx, tx, id, repository.StatusEfekAntre,
				c.Waktu.Add(Backoff(1)), c.Sebab, c.Waktu)
		}
		// ⛔ Permanen: nol jadwal. Ia menunggu manusia, bukan menunggu waktu.
		return a.pohon.TuntaskanEfek(ctx, tx, id,
			repository.StatusEfekGagalPermanen, time.Time{}, c.Sebab, c.Waktu)
	})
}

// PelaksanaEfek menjalankan satu baris outbox yang sudah dipungut.
//
// ⛔ Antarmuka, bukan `switch` di dalam pekerja. Pekerja tahu cara memungut
// dan menuntaskan; ia TIDAK tahu apa arti "berkas" atau "arasapas".
type PelaksanaEfek interface {
	Laksanakan(ctx context.Context, b repository.BarisEfekKeluar) error
}

// PekerjaEfek memungut satu efek yang jatuh tempo dan menuntaskannya.
type PekerjaEfek struct {
	svc       *Service
	pohon     *repository.PohonKlaim
	pelaksana PelaksanaEfek
}

// NewPekerjaEfek menyusun pekerjanya.
func NewPekerjaEfek(svc *Service, p PelaksanaEfek) *PekerjaEfek {
	return &PekerjaEfek{svc: svc, pohon: repository.NewPohonKlaim(svc.db),
		pelaksana: p}
}

// SatuPutaran memungut SATU efek, menjalankannya, lalu menuntaskannya.
//
// Mengembalikan `ErrPekerjaTanpaEfek` bila tidak ada yang jatuh tempo.
//
// ⛔ Seluruhnya SATU transaksi: pungut, jalankan, tuntaskan. Menuntaskan di
// transaksi lain berarti efek yang berhasil dapat gagal dicatat berhasil - dan
// dijalankan lagi oleh putaran berikutnya.
//
// ⚠️ Konsekuensi yang diterima: transaksi terbuka selama efeknya berjalan,
// termasuk selama panggilan jaringan. Itu ditanggung karena alternatifnya
// menjalankan efek dua kali, dan efek keluar tidak selalu idempoten.
func (w *PekerjaEfek) SatuPutaran(ctx context.Context, saat time.Time) error {
	if !w.svc.PunyaDatabase() {
		return repository.ErrTanpaOracle
	}
	if w.pelaksana == nil {
		return errors.New("services: pekerja efek tanpa pelaksana")
	}
	return w.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		baris, err := w.pohon.PungutEfek(ctx, tx, saat)
		if err != nil {
			return err
		}
		jalanErr := w.pelaksana.Laksanakan(ctx, baris)
		if jalanErr == nil {
			return w.pohon.TuntaskanEfek(ctx, tx, baris.ID,
				repository.StatusEfekSelesai, time.Time{}, "", saat)
		}
		status, jadwal := repository.StatusEfekAntre,
			saat.Add(Backoff(baris.Percobaan+1))
		// ⛔ DUA jalan menuju menyerah: tidak layak dicoba ulang, ATAU jatah
		// percobaannya habis. Ronde pertama hanya punya yang pertama, dan
		// kegagalan jaringan yang tak kunjung pulih akan berputar selamanya.
		if !LayakDicobaUlang(jalanErr) || baris.Percobaan >= percobaanMaksimum {
			status, jadwal = repository.StatusEfekGagalPermanen, time.Time{}
		}
		return w.pohon.TuntaskanEfek(ctx, tx, baris.ID, status, jadwal,
			jalanErr.Error(), saat)
	})
}

// PenyalurClaimLifeOracle menyusun penyalur dengan KETIGA ketergantungan yang
// sudah punya jalan ke Oracle - antre-ulang dan resolver endpoint.
//
// ⛔ Lingkungannya tetap datang dari `Service`, bukan dari sini. Penyalur yang
// memutuskan sendiri "ini produksi" adalah penyalur yang suatu hari mengirim
// email kepada orang sungguhan dari mesin pengembang.
//
// ⚠️ Ketiga EFEKnya masih gagal terang (`…BelumDisetujui`) - menghubungkan
// storage, email, dan Arasapas nyata menuntut persetujuan manusia. Yang
// ditutup butir **aq** adalah TEMPAT kegagalannya mendarat, bukan izin
// memanggil layanannya.
func PenyalurClaimLifeOracle(svc *Service) *Penyalur {
	return NewPenyalur(svc.lingkungan, AntreanEfekOracle(svc),
		EfekKeluarClaimLife(ResolverLinkServiceOracle(svc))...)
}
