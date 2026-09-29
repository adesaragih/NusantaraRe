package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/jejak"
)

// ErrPekerjaTanpaEfek diteruskan apa adanya dari repository: tidak ada yang
// jatuh tempo. Keadaan NORMAL.
var ErrPekerjaTanpaEfek = ErrEfekTidakAda

// Modul dan lini yang menandai baris outbox milik modul ini.
//
// ⚠️ Outbox-nya LINTAS MODUL. Tanpa kedua penanda ini, pekerja Komite dan
// pekerja Claim Life akan saling memungut pekerjaan yang bukan miliknya.
const (
	LiniLife        = "LIFE"
	backoffAwal     = 30 * time.Second
	backoffMaksimum = 6 * time.Hour
	// percobaanMaksimum - sesudahnya efek menyerah dan menunggu manusia.
	//
	// ⚠️ Ada BATASnya, dan itu disengaja. Antrean tanpa batas percobaan
	// menyimpan baris yang dicoba selamanya, dan kegagalan sungguhan tidak
	// pernah terlihat sebagai kegagalan oleh siapa pun.
	PercobaanMaksimum = 8
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

// Kosakata jejak audit untuk kegagalan efek keluar.
//
// ⛔ AC 20 tiket 12 menuntut kegagalan tercatat di JALUR AUDIT, bukan hanya
// di log layanan. Outbox saja tidak cukup: outbox adalah antrean kerja,
// dan orang yang bertanya "apa yang terjadi pada klaim ini" membaca jejak
// auditnya, bukan antrean pekerja.
//
// ⚠️ Hanya kegagalan PERMANEN yang masuk jejak. Kegagalan yang masih akan
// dicoba lagi belum menjadi sejarah klaim - mencatatnya berarti membanjiri
// jejak audit dengan delapan baris untuk satu email yang akhirnya terkirim.
const (
	// TahapEfekKeluar mengisi kolom `DARI` - asal transisinya.
	TahapEfekKeluar = "efek-keluar"
	// TahapEfekMenyerah mengisi kolom `KE`.
	TahapEfekMenyerah = "gagal-permanen"
)

// muatanOutbox adalah bentuk JSON yang disimpan di kolom `MUATAN`.
//
// ⛔ PENGENAL DAN WAKTU SAJA. Nol nama orang, nol email, nol alamat, nol
// kredensial - kolom ini dibaca lagi entah kapan oleh entah siapa. Penjaga
// statik menegakkan bentuknya.
type MuatanOutbox struct {
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
	svc   inti.Akar
	pohon *Penyimpan
	jejak jejak.Jejak
	// modul mengisi kolom `MODUL` - sejak tiket 06 PremiumList Life, outbox
	// ini dipakai dua modul.
	modul string
}

// AntreanEfekOracleModul sama dengan `AntreanEfekOracle`, untuk modul lain.
//
// ⚠️ Tiket 06 PremiumList Life. Satu tabel outbox untuk dua modul, dipisah
// kolom `MODUL`. `PekerjaEfek` memungut HANYA modulnya sendiri
// (`sqlPungutEfek`); baris PREMIUMLISTLIFE menunggu pekerja PremiumList, yang
// belum ada - ia tetap terlihat berstatus `antre`, tidak dipungut orang lain.
func AntreanEfekOracleModul(svc inti.Akar, modul string) Antrean {
	return antreanOracle{svc: svc, pohon: NewPenyimpan(svc.DB()),
		jejak: jejak.PerekamJejakOracle(svc), modul: modul}
}

// Antre menyimpan satu kegagalan untuk dicoba lagi - atau untuk ditunggu
// manusia, bila ia tidak layak dicoba ulang.
func (a antreanOracle) Antre(ctx context.Context, c CatatanEfekGagal) error {
	if !a.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	muatan, err := json.Marshal(MuatanOutbox{
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
	return a.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		id, err := a.pohon.AntreEfek(ctx, tx, LiniLife, a.modul,
			c.Nama, rujukan, string(muatan), c.Waktu)
		if err != nil {
			return err
		}
		if c.LayakUlang {
			// Percobaan pertama dijadwalkan; `AntreEfek` menulisnya jatuh
			// tempo seketika, jadi jadwalnya digeser di sini.
			return a.pohon.TuntaskanEfek(ctx, tx, id, StatusEfekAntre,
				c.Waktu.Add(Backoff(1)), c.Sebab, c.Waktu)
		}
		// ⛔ Permanen: nol jadwal. Ia menunggu manusia, bukan menunggu waktu.
		if err := a.pohon.TuntaskanEfek(ctx, tx, id,
			StatusEfekGagalPermanen, time.Time{}, c.Sebab,
			c.Waktu); err != nil {
			return err
		}
		// ⛔ Dan ia masuk JALUR AUDIT, dalam transaksi yang sama (AC 20
		// tiket 12). Transaksi terpisah berarti ada saat ketika outbox
		// berkata menyerah sementara jejaknya belum mengatakan apa pun.
		return a.rekamMenyerah(ctx, tx, c)
	})
}

// rekamMenyerah menulis satu kegagalan permanen ke jejak audit.
//
// ⚠️ `Dari`/`Ke` TIDAK memuat pesan galatnya. Pesan galat dapat menyebut
// nama objek basis data, alamat, bahkan nilai kolom; rinciannya tinggal di
// kolom `GALAT_TERAKHIR` outbox. Yang masuk jejak: efek apa, dan bahwa ia
// menyerah.
func (a antreanOracle) rekamMenyerah(ctx context.Context, tx *db.Tx,
	c CatatanEfekGagal) error {

	if a.jejak == nil {
		return jejak.ErrJejakBelumDiputuskan
	}
	return a.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
		AdjustmentID: c.AdjustmentID,
		KlaimID:      c.KlaimID,
		Dari:         TahapEfekKeluar + ":" + c.Nama,
		Ke:           TahapEfekMenyerah,
		AkunID:       c.AkunID,
		Waktu:        c.Waktu,
	})
}

// PelaksanaEfek menjalankan satu baris outbox yang sudah dipungut.
//
// ⛔ Antarmuka, bukan `switch` di dalam pekerja. Pekerja tahu cara memungut
// dan menuntaskan; ia TIDAK tahu apa arti "berkas" atau "arasapas".
type PelaksanaEfek interface {
	Laksanakan(ctx context.Context, tx *db.Tx,
		b BarisEfekKeluar) error
}

// ⚠️ `tx` MASUK 27-09-2026, bersama butir be. Sebabnya: efek yang
// menuntaskan dirinya sendiri harus menulis DI DALAM transaksi yang sama
// dengan penuntasan barisnya - `storage-unggah` mengisi `T_STORAGE_ID` dan
// menulis kartu berkas, dan kedua tulisan itu tidak boleh berdiri sendirian
// bila penuntasannya gagal. Pola yang sama dengan `Jejak.Rekam`.
//
// ⛔ Pekerja tetap TIDAK tahu apa arti "berkas": ia menyerahkan
// transaksinya, bukan maknanya.

// PekerjaEfek memungut satu efek yang jatuh tempo dan menuntaskannya.
type PekerjaEfek struct {
	svc       inti.Akar
	pohon     *Penyimpan
	pelaksana PelaksanaEfek
	// penjejak dipinjam dari antrean: DUA jalur dapat menyerah - yang
	// pertama saat efeknya gagal permanen sejak awal, yang kedua saat
	// jatah percobaannya habis di sini. Keduanya wajib meninggalkan
	// jejak yang sama bentuknya, jadi keduanya memanggil kode yang sama.
	penjejak antreanOracle
	// modul menyaring baris yang dipungut - hanya milik modul pekerja ini.
	modul string
}

// NewPekerjaEfekModul menyusun pekerja untuk modul lain - tiket 07 Komite.
//
// ⚠️ Aditif pada berkas Claim Life: penyaring `MODUL` pemungutan sudah ada
// (`sqlPungutEfek`), dan antrean pencatat menyerahnya membawa modul yang sama.
//
// Refactor bentuk B (30-09-2026): pekerja dirakit LANGSUNG di sini, tidak
// lagi lewat `NewPekerjaEfek` milik Claim Life - paket bersama tidak boleh
// mengenal nama modul mana pun. Hasilnya sama: medan `modul` pekerja dan
// penjejaknya berisi modul yang diminta.
func NewPekerjaEfekModul(svc inti.Akar, p PelaksanaEfek, modul string) *PekerjaEfek {
	return &PekerjaEfek{
		svc:       svc,
		pohon:     NewPenyimpan(svc.DB()),
		pelaksana: p,
		penjejak:  antreanOracle{svc: svc, jejak: jejak.PerekamJejakOracle(svc), modul: modul},
		modul:     modul,
	}
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
		return db.ErrTanpaOracle
	}
	if w.pelaksana == nil {
		return errors.New("services: pekerja efek tanpa pelaksana")
	}
	return w.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		baris, err := w.pohon.PungutEfek(ctx, tx, w.modul, saat)
		if err != nil {
			return err
		}
		jalanErr := w.pelaksana.Laksanakan(ctx, tx, baris)
		if jalanErr == nil {
			return w.pohon.TuntaskanEfek(ctx, tx, baris.ID,
				StatusEfekSelesai, time.Time{}, "", saat)
		}
		status, jadwal := StatusEfekAntre,
			saat.Add(Backoff(baris.Percobaan+1))
		// ⛔ DUA jalan menuju menyerah: tidak layak dicoba ulang, ATAU jatah
		// percobaannya habis. Ronde pertama hanya punya yang pertama, dan
		// kegagalan jaringan yang tak kunjung pulih akan berputar selamanya.
		menyerah := !LayakDicobaUlang(jalanErr) ||
			baris.Percobaan >= PercobaanMaksimum
		if menyerah {
			status, jadwal = StatusEfekGagalPermanen, time.Time{}
		}
		if err := w.pohon.TuntaskanEfek(ctx, tx, baris.ID, status, jadwal,
			jalanErr.Error(), saat); err != nil {
			return err
		}
		if !menyerah {
			return nil
		}
		// ⛔ Jalur menyerah KEDUA - jatah percobaan habis. AC 20 tiket 12
		// berlaku di sini persis seperti di jalur pertama; email yang
		// menyerah sesudah delapan percobaan sama tidak terkirimnya
		// dengan yang menyerah seketika.
		return w.penjejak.rekamMenyerah(ctx, tx, CatatanEfekGagal{
			MuatanEfek: bacaMuatan(baris.Muatan, saat),
			Nama:       baris.Jenis,
			Sebab:      jalanErr.Error(),
		})
	})
}

// bacaMuatan membaca kembali muatan JSON sebuah baris outbox.
//
// ⚠️ Muatan yang TIDAK terbaca tidak menggagalkan penyerahan. Baris outbox
// berumur panjang; bentuk JSON-nya dapat berubah di antara saat ia ditulis
// dan saat ia menyerah. Jejak dengan pengenal kosong masih lebih berguna
// daripada kegagalan yang tidak tercatat sama sekali - dan `RUJUKAN` di
// outbox tetap menyimpan tautannya.
func bacaMuatan(teks string, saat time.Time) MuatanEfek {
	var m MuatanOutbox
	if err := json.Unmarshal([]byte(teks), &m); err != nil {
		return MuatanEfek{Waktu: saat}
	}
	waktu, err := time.Parse(time.RFC3339, m.Waktu)
	if err != nil {
		waktu = saat
	}
	return MuatanEfek{
		KlaimID:      m.KlaimID,
		AdjustmentID: m.AdjustmentID,
		AkunID:       m.AkunID,
		Waktu:        waktu,
	}
}
