package repository

// Outbox efek keluar - butir aq, A2.
//
// Untuk apa berkas ini: menyimpan efek keluar DALAM TRANSAKSI aksi bisnisnya,
// lalu memberi worker cara memungutnya satu per satu tanpa dua worker memungut
// baris yang sama.
//
// ⛔ Inti pola outbox ada pada SATU transaksi. `AntreEfek` menerima `*Tx` dan
// menuntutnya bukan nil: bila klaim tersimpan, efeknya PASTI terantre; bila
// transaksinya batal, antreannya ikut batal. Tidak ada saat ketika salah
// satunya ada tanpa yang lain.
//
// ⚠️ `PungutEfek` justru TIDAK boleh berbagi transaksi dengan aksi bisnis - ia
// berjalan belakangan, di proses yang mungkin lain.
//
// Dibaca sesudah: linkservice.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Status outbox - himpunan tertutup, TEKS (ADR-U-0022).
const (
	// StatusEfekAntre - menunggu dipungut.
	StatusEfekAntre = "antre"
	// StatusEfekJalan - sudah dipungut, sedang dikerjakan.
	StatusEfekJalan = "jalan"
	// StatusEfekSelesai - berhasil.
	StatusEfekSelesai = "selesai"
	// StatusEfekGagalPermanen - tidak akan dicoba lagi tanpa manusia.
	StatusEfekGagalPermanen = "gagal-permanen"
)

// ErrEfekTidakAda - tidak ada efek yang jatuh tempo.
//
// ⚠️ Keadaan NORMAL, bukan kegagalan. Worker yang menganggapnya galat akan
// mencatat galat setiap detik sepanjang malam.
var ErrEfekTidakAda = errors.New("repository: tidak ada efek keluar yang jatuh tempo")

// BarisEfekKeluar adalah satu efek yang menunggu dijalankan.
//
// ⚠️ `Muatan` TEKS JSON, bukan struktur. Ia menyeberang batas proses dan batas
// waktu; tipe Go yang berubah minggu depan tidak boleh membuat baris yang
// sudah terantre tak terbaca.
type BarisEfekKeluar struct {
	ID        string
	Lini      string
	Modul     string
	Jenis     string
	Rujukan   string
	Muatan    string
	Percobaan int
}

// sqlAntreEfek merakit penulisannya.
func sqlAntreEfek(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
		(ID, LINI, MODUL, JENIS_EFEK, RUJUKAN, MUATAN, STATUS, PERCOBAAN,
		 JADWAL_BERIKUT, DIBUAT)
		VALUES (:1,:2,:3,:4,:5,:6,:7,0,:8,:9)`, tabel)
}

// AntreEfek menulis satu efek keluar ke outbox, di dalam transaksi pemanggil.
//
// ⛔ `tx` nil DITOLAK, tidak diam-diam diganti sambungan biasa. Outbox yang
// ditulis di luar transaksi bisnisnya bukan outbox - ia tabel log biasa yang
// dapat berbeda isi dari klaimnya.
func (r *PohonKlaim) AntreEfek(ctx context.Context, tx *Tx,
	lini, modul, jenis, rujukan, muatan string, saat time.Time) (string, error) {

	if tx == nil {
		return "", fmt.Errorf("repository: AntreEfek menuntut transaksi; " +
			"outbox di luar transaksi bisnisnya tidak menjamin apa pun")
	}
	id, err := r.nomorBerikut(ctx, tx, "SEQ_LOG_SERVICE_RNM")
	if err != nil {
		return "", err
	}
	tabel, err := r.db.Qualify("T_LOG_SERVICE_RNM")
	if err != nil {
		return "", err
	}
	q := sqlAntreEfek(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	// ⛔ `JADWAL_BERIKUT` diisi `saat` - BUKAN dibiarkan kosong. Worker
	// memilih `JADWAL_BERIKUT <= sekarang`; baris berjadwal kosong tidak
	// pernah terpilih, dan efeknya tidak pernah jalan tanpa satu pun galat.
	hasil, err := tx.tx.ExecContext(ctx, q, id, kosongJadiNil(lini),
		kosongJadiNil(modul), kosongJadiNil(jenis), kosongJadiNil(rujukan),
		kosongJadiNil(muatan), StatusEfekAntre, saat, saat)
	if err != nil {
		return "", fmt.Errorf("repository: mengantre efek keluar: %w", err)
	}
	if err := pastikanSatuBaris(hasil, "pengantrean efek keluar"); err != nil {
		return "", err
	}
	return id, nil
}

// sqlPungutEfek merakit pemungutannya.
//
// ⛔ `FOR UPDATE SKIP LOCKED` adalah yang membuat dua worker tidak memungut
// baris yang sama: yang kedua MELEWATI baris terkunci alih-alih menunggunya.
// Tanpa `SKIP LOCKED` worker kedua menunggu yang pertama selesai dan
// paralelismenya hilang; tanpa `FOR UPDATE` keduanya memungut baris yang sama
// dan efeknya berjalan dua kali.
//
// ⛔ `AND MODUL = :3` - sejak tiket 06 PremiumList Life outbox ini dipakai
// DUA modul. Tanpa penyaring itu worker Claim Life memungut baris polis,
// pelaksananya tidak mengenal jenisnya, dan baris itu ditandai GAGAL PERMANEN
// oleh worker yang bukan pemiliknya - kegagalan palsu yang menutupi yang asli.
func sqlPungutEfek(tabel string) string {
	return fmt.Sprintf(`SELECT ID, LINI, MODUL, JENIS_EFEK, RUJUKAN, MUATAN,
			   PERCOBAAN
		  FROM %s
		 WHERE STATUS = :1
		   AND JADWAL_BERIKUT <= :2
		   AND MODUL = :3
		 ORDER BY JADWAL_BERIKUT ASC
		 FETCH FIRST 1 ROWS ONLY
		 FOR UPDATE SKIP LOCKED`, tabel)
}

// sqlTandaiJalan merakit penandaannya.
func sqlTandaiJalan(tabel string) string {
	return fmt.Sprintf(`UPDATE %s
		   SET STATUS = :1, PERCOBAAN = PERCOBAAN + 1, DIPERBARUI = :2
		 WHERE ID = :3 AND STATUS = :4`, tabel)
}

// PungutEfek mengambil SATU efek yang jatuh tempo dan menandainya jalan.
//
// ⛔ Pemungutan dan penandaannya SATU transaksi. Memungut lalu menandai di
// transaksi lain berarti ada jendela ketika baris sudah dipungut tetapi masih
// berstatus `antre` - dan worker kedua memungutnya lagi.
//
// ⚠️ `PERCOBAAN` dinaikkan SAAT DIPUNGUT, bukan saat gagal. Efek yang membuat
// worker-nya mati di tengah jalan tetap menghabiskan jatahnya; bila tidak, ia
// akan membunuh worker berikutnya, dan berikutnya lagi, selamanya.
func (r *PohonKlaim) PungutEfek(ctx context.Context, tx *Tx,
	modulPekerja string, saat time.Time) (BarisEfekKeluar, error) {

	if tx == nil {
		return BarisEfekKeluar{}, fmt.Errorf("repository: PungutEfek menuntut " +
			"transaksi; kuncian SKIP LOCKED hidup selama transaksinya saja")
	}
	tabel, err := r.db.Qualify("T_LOG_SERVICE_RNM")
	if err != nil {
		return BarisEfekKeluar{}, err
	}
	q := sqlPungutEfek(tabel)
	if err := PeriksaSQL(q); err != nil {
		return BarisEfekKeluar{}, err
	}
	var (
		b                                   BarisEfekKeluar
		lini, modul, jenis, rujukan, muatan sql.NullString
		percobaan                           sql.NullInt64
	)
	err = tx.tx.QueryRowContext(ctx, q, StatusEfekAntre, saat, modulPekerja).Scan(
		&b.ID, &lini, &modul, &jenis, &rujukan, &muatan, &percobaan)
	if err == sql.ErrNoRows {
		return BarisEfekKeluar{}, ErrEfekTidakAda
	}
	if err != nil {
		return BarisEfekKeluar{}, fmt.Errorf("repository: memungut efek keluar: %w", err)
	}
	b.Lini, b.Modul, b.Jenis = lini.String, modul.String, jenis.String
	b.Rujukan, b.Muatan = rujukan.String, muatan.String
	b.Percobaan = int(percobaan.Int64)

	tandai := sqlTandaiJalan(tabel)
	if err := PeriksaSQL(tandai); err != nil {
		return BarisEfekKeluar{}, err
	}
	// ⛔ `AND STATUS = :4` pengaman ganda: seandainya `SKIP LOCKED` tidak
	// berlaku (Oracle lama, atau query diubah kelak), penandaan yang mengenai
	// NOL baris membuat pemungutan ganda GAGAL TERANG alih-alih menjalankan
	// efeknya dua kali.
	hasil, err := tx.tx.ExecContext(ctx, tandai, StatusEfekJalan, saat, b.ID,
		StatusEfekAntre)
	if err != nil {
		return BarisEfekKeluar{}, fmt.Errorf("repository: menandai efek jalan: %w", err)
	}
	if err := pastikanSatuBaris(hasil, "penandaan efek jalan"); err != nil {
		return BarisEfekKeluar{}, err
	}
	return b, nil
}

// sqlTuntaskanEfek merakit penuntasannya.
func sqlTuntaskanEfek(tabel string) string {
	return fmt.Sprintf(`UPDATE %s
		   SET STATUS = :1, JADWAL_BERIKUT = :2, GALAT_TERAKHIR = :3,
			   DIPERBARUI = :4
		 WHERE ID = :5`, tabel)
}

// panjangGalatTersimpan adalah lebar kolom `GALAT_TERAKHIR`.
const panjangGalatTersimpan = 4000

// TuntaskanEfek menutup sebuah efek: selesai, dijadwal ulang, atau menyerah.
//
// ⚠️ `galat` DIPANGKAS pada lebar kolomnya. Galat yang lebih panjang membuat
// seluruh penuntasan gagal - dan efek yang tidak dapat dituntaskan akan
// berstatus `jalan` selamanya, tak terpungut oleh siapa pun.
//
// ⚠️ `jadwal` kosong ditulis KOSONG, bukan waktu nol tahun 1. Efek yang sudah
// selesai tidak punya jadwal berikutnya, dan kosong berbeda dari nol
// (ADR-U-0027).
func (r *PohonKlaim) TuntaskanEfek(ctx context.Context, tx *Tx, id, status string,
	jadwal time.Time, galat string, saat time.Time) error {

	if tx == nil {
		return fmt.Errorf("repository: TuntaskanEfek menuntut transaksi")
	}
	tabel, err := r.db.Qualify("T_LOG_SERVICE_RNM")
	if err != nil {
		return err
	}
	q := sqlTuntaskanEfek(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	if len(galat) > panjangGalatTersimpan {
		galat = galat[:panjangGalatTersimpan]
	}
	var waktuJadwal any
	if !jadwal.IsZero() {
		waktuJadwal = jadwal
	}
	hasil, err := tx.tx.ExecContext(ctx, q, status, waktuJadwal,
		kosongJadiNil(galat), saat, id)
	if err != nil {
		return fmt.Errorf("repository: menuntaskan efek keluar: %w", err)
	}
	return pastikanSatuBaris(hasil, "penuntasan efek keluar")
}
