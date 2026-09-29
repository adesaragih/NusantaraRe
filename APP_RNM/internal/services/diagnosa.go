package services

// Diagnosa yang dilekatkan pada peserta - butir bd, §2 A bagian 2.
//
// Untuk apa berkas ini: ketiga tombol grid `.DiagnoseList` - `Add` b4690,
// `Choose` b2509 (lewat `SetDisease`), dan `Delete` b6160.
//
// Dibaca sesudah: penyakit.go (katalognya) dan statusbaris.go (yang
// memutuskan `STS_REJECT` peserta, yang di sini hanya dicerminkan).
//
// ⛔ TIGA GERBANG, dan ketiganya dari pohon XML - bukan dari selera:
//
//  1. TAHAP. Grid diagnosa hidup di `ClaimLifeDetailGCNM`, dan pintu
//     masuknya `ViewClaimDetailLifeGCNM` - `pyEditAction` grid peserta di
//     `InputOSClaimLife` b18252, `MedicalCheckClaimLife` b17416, dan
//     `InputAkseptasiClaimLife` b17387. TIGA tahap, bukan satu, dan bukan
//     Input Register.
//  2. PEMEGANG. Orang menyunting pekerjaan yang SEDANG IA PEGANG - pola yang
//     sama dengan `tahap.go` dan `tutup.go`.
//  3. `STS_REJECT` PESERTA. Satu kalimat di empat tempat (b4682, b5059,
//     b5870, b6152) - lihat `models.DiagnosaTerkunci`.
//
// ⚠️ Gerbang 1 dan 2 TIDAK ada sebagai `pyWhenName` di rule mana pun -
// ketiga pemuatnya berprivilese kosong. Ia disimpulkan dari LETAK: sebuah
// layar yang hanya dapat dibuka dari tiga tahap hanya dapat dipakai dari
// tiga tahap. Kesimpulan itu dinyatakan di sini supaya dapat dibantah, bukan
// disembunyikan di dalam kode.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
)

var (
	// ErrDiagnosaTerkunci - peserta sudah diputus, diagnosanya beku.
	//
	// Padanan `pyDisabledWhen` - di Pega kontrolnya MATI, di sini
	// permintaannya DITOLAK. Layar bukan penjaga: tombol yang mati hanya
	// menahan jari, bukan permintaan HTTP.
	ErrDiagnosaTerkunci = errors.New("services: peserta sudah diputus, diagnosanya tidak dapat diubah")
	// ErrTahapTidakBergridPeserta - tahap ini tidak membuka layar detail.
	ErrTahapTidakBergridPeserta = errors.New("services: tahap ini tidak membuka grid peserta")
	// ErrNilaiDiagnosaKepanjangan - satu nilai melebihi lebar kolomnya.
	ErrNilaiDiagnosaKepanjangan = errors.New("services: nilai diagnosa melebihi lebar kolomnya")
)

// DiagnosaPeserta melayani ketiga tombol grid diagnosa.
type DiagnosaPeserta struct{ svc *Service }

// Diagnosa menyusun layanannya.
func (s *Service) Diagnosa() *DiagnosaPeserta { return &DiagnosaPeserta{svc: s} }

// konteksDiagnosa adalah hasil ketiga gerbang, dipakai bersama oleh ketiga
// rute.
//
// ⛔ SATU tempat, bukan tiga. Gerbang yang disalin ke tiga pintu adalah tiga
// kesempatan untuk berbeda - dan yang berbeda akan menjadi pintu yang
// terbuka sendiri. Pelajaran `PastikanKasusTerbuka` (butir bb), yang sengaja
// diletakkan di badan `ubah` alih-alih di `Tolak` dan `Ubah`.
type konteksDiagnosa struct {
	baca    *repository.KlaimLife
	diag    *repository.Diagnosa
	peserta models.Peserta
}

// pagari menjalankan KETIGA gerbang dan menyiapkan pembacanya.
//
// Urutannya bukan gaya: identitas, bentuk permintaan, kasus terbuka, tahap,
// pemegang, lalu keadaan peserta. Permintaan tanpa identitas tidak berhak
// tahu apakah klaimnya ada; permintaan tanpa pengenal peserta harus dijawab
// "pengenal wajib diisi", bukan "ORACLE_DSN belum dikonfigurasi" - cacat
// yang pernah nyata di `statusbaris.go` dan ditangkap
// `TestUbahStatusMenjagaPagarnya`.
func (d *DiagnosaPeserta) pagari(ctx context.Context, pelaku inti.Pelaku,
	klaimID, pesertaID string) (konteksDiagnosa, error) {

	var k konteksDiagnosa
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return k, err
	}
	if strings.TrimSpace(klaimID) == "" || strings.TrimSpace(pesertaID) == "" {
		return k, fmt.Errorf("%w: pengenal klaim dan peserta wajib diisi",
			galat.ErrPermintaanTidakSah)
	}
	if d == nil || d.svc == nil || !d.svc.PunyaDatabase() {
		return k, db.ErrTanpaOracle
	}
	// ⛔ BUTIR bb - kasus tertutup tidak dapat diubah lagi. Penjaga statik
	// `TestSetiapLayananPengubahMemeriksaKasusTerbuka` menagih baris ini.
	if err := d.svc.PastikanKasusTerbuka(ctx, klaimID); err != nil {
		return k, err
	}

	k.baca = repository.NewKlaimLife(d.svc.DB())
	k.diag = repository.NewDiagnosa(d.svc.DB())

	// Gerbang 1 dan 2 - tahap, lalu pemegangnya. Kolom TAHAP menang;
	// PY_POSITION cadangan untuk baris lama (butir at).
	tahap, err := tahapKasus(ctx, k.baca, klaimID)
	if err != nil {
		return k, err
	}
	if !models.TahapBergridPeserta(tahap) {
		return k, fmt.Errorf("%w: %s (yang membukanya %v)",
			ErrTahapTidakBergridPeserta, tahap, models.TahapPembukaGridPeserta())
	}
	peranTahap, ada := models.PeranPemegangTahap(tahap)
	if !ada {
		return k, fmt.Errorf("%w: tahap %q", ErrPeranTahapBelumDiputuskan, tahap)
	}
	if err := inti.WajibPeran(pelaku, peranTahap); err != nil {
		return k, err
	}

	// Gerbang 3 - keadaan peserta.
	peserta, err := k.baca.AmbilPeserta(ctx, klaimID)
	if err != nil {
		return k, err
	}
	for _, p := range peserta {
		if p.ID == pesertaID {
			k.peserta = p
		}
	}
	if k.peserta.ID == "" {
		return k, fmt.Errorf("%w: peserta %q bukan milik klaim %q",
			galat.ErrPermintaanTidakSah, pesertaID, klaimID)
	}
	if models.DiagnosaTerkunci(k.peserta.KodeStatus) {
		return k, fmt.Errorf("%w: peserta %q berkode %q",
			ErrDiagnosaTerkunci, pesertaID, k.peserta.KodeStatus)
	}
	return k, nil
}

// Tambah menyisipkan satu baris diagnosa kosong - `Add` b4690.
//
// `[terverifikasi]` b4700 `pyAction addRow`, b4715 `pyPosition After` - baris
// baru di EKOR. Diikuti b4730 `refresh`, BUKAN `save`: di Pega baris itu
// belum menetap sampai sesuatu yang lain menyimpannya.
//
// ⚠️ PENYIMPANGAN SADAR, dan arahnya disengaja. Di sini baris itu LANGSUNG
// menetap. Sebabnya: halaman Pega hidup di sesi server dan dapat menampung
// baris yang belum tersimpan; HTTP tidak punya tempat setara, dan menirunya
// berarti menahan baris kosong di memori layar sampai `Choose` datang - yang
// membuat `Delete` atas baris yang belum ada harus dikarang tersendiri, dan
// membuat dua pemakai pada peserta yang sama melihat urutan yang berbeda.
// Yang hilang: baris kosong yang ditambahkan lalu ditinggalkan akan
// tersimpan. Yang didapat: `URUTAN` yang sama untuk semua yang melihatnya.
func (d *DiagnosaPeserta) Tambah(ctx context.Context, pelaku inti.Pelaku,
	klaimID, pesertaID string) (models.Diagnosa, error) {

	k, err := d.pagari(ctx, pelaku, klaimID, pesertaID)
	if err != nil {
		return models.Diagnosa{}, err
	}
	perPeserta, err := k.diag.AmbilDiagnosa(ctx, klaimID)
	if err != nil {
		return models.Diagnosa{}, err
	}
	urutan := models.UrutanBerikutnya(perPeserta[pesertaID])
	// ⚠️ DIBACA DI LUAR TRANSAKSI, jadi dua `Add` serentak pada peserta
	// yang sama dapat memperoleh `URUTAN` yang sama. Diterima, dan sebabnya
	// dinyatakan: pembacaan `ORDER BY URUTAN, ID` tetap memberi susunan yang
	// PASTI - `ID` yang memutus seri - dan `RapatkanUrutan` menomori ulang
	// pada penghapusan berikutnya. Menguncinya *(`SELECT MAX ... FOR UPDATE`)*
	// akan menyerialkan setiap penambahan diagnosa demi kerapian nomor yang
	// tidak dilihat siapa pun; di Pega pun subscript PageList hanya berarti
	// urutan tampil.

	var lahir models.Diagnosa
	err = d.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		var e error
		lahir, e = k.diag.SisipDiagnosa(ctx, tx, pesertaID, urutan,
			k.peserta.KodeStatus)
		return e
	})
	if err != nil {
		return models.Diagnosa{}, err
	}
	return lahir, nil
}

// Ubah menulis hasil `Choose` - b2509 -> `SetDisease` b2528.
//
// `[terverifikasi]` `Activity/SetDisease.xml` b260-261 dan b307-308 menulis
// DUA kolom dari parameternya; `GROUP_DIAGNOSE` datang dari dropdown b5860
// yang mem-`postValue` b5886 sendiri. Satu rute menulis ketiganya sebab di
// layar keduanya menyunting baris yang sama, dan dua rute berarti baris yang
// separuh tersimpan.
//
// ⛔ `groupDiagnose` diterima apa adanya - `[terbuka - OQ-L]`. Daftar
// pilihannya hidup pada rule properti yang tidak diekspor, jadi satu-satunya
// pemeriksaan yang jujur adalah panjangnya.
func (d *DiagnosaPeserta) Ubah(ctx context.Context, pelaku inti.Pelaku,
	klaimID, pesertaID string, diagID int64, kodeICD, nama, grup string) error {

	k, err := d.pagari(ctx, pelaku, klaimID, pesertaID)
	if err != nil {
		return err
	}
	kodeICD = strings.TrimSpace(kodeICD)
	nama = strings.TrimSpace(nama)
	grup = strings.TrimSpace(grup)
	if kolom, muat := models.PotongDiagnosa(kodeICD, nama, grup); !muat {
		return fmt.Errorf("%w: kolom %s", ErrNilaiDiagnosaKepanjangan, kolom)
	}
	if _, err := d.milikPeserta(ctx, k, klaimID, pesertaID, diagID); err != nil {
		return err
	}
	return d.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		return k.diag.PerbaruiDiagnosa(ctx, tx, diagID, kodeICD, nama, grup)
	})
}

// Hapus membuang satu baris - `Delete` b6160.
//
// `[terverifikasi]` b6170 `deleteRow` DAN b6191 `save` (b6189 `pyActionLabel
// Save`). Dua perilaku pada satu klik, dan yang kedua itulah bedanya dengan
// `Add`: penghapusan MENETAP seketika. Rute ini karena itu tidak menunda.
//
// ⛔ `URUTAN` dirapatkan DI DALAM transaksi yang sama. `deleteRow` di Pega
// menggeser subscript anggota sesudahnya sendiri; kolom tidak bergeser
// sendiri, dan daftar yang berlubang akan menomori baris berikutnya dengan
// angka yang sudah dipakai.
func (d *DiagnosaPeserta) Hapus(ctx context.Context, pelaku inti.Pelaku,
	klaimID, pesertaID string, diagID int64) error {

	k, err := d.pagari(ctx, pelaku, klaimID, pesertaID)
	if err != nil {
		return err
	}
	sasaran, err := d.milikPeserta(ctx, k, klaimID, pesertaID, diagID)
	if err != nil {
		return err
	}
	return d.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		if err := k.diag.HapusDiagnosa(ctx, tx, diagID); err != nil {
			return err
		}
		return k.diag.RapatkanUrutan(ctx, tx, pesertaID, sasaran.Urutan)
	})
}

// milikPeserta memastikan baris diagnosa benar-benar milik peserta itu.
//
// ⛔ Pengenal diagnosa datang dari JALUR URL, dan jalur URL datang dari
// siapa saja. Tanpa pemeriksaan ini, `DELETE .../klaim/A/peserta/B/diagnosa/9`
// akan menghapus diagnosa nomor 9 milik klaim siapa pun - gerbang tahap dan
// gerbang `STS_REJECT` di atas memeriksa peserta B, bukan baris 9.
func (d *DiagnosaPeserta) milikPeserta(ctx context.Context, k konteksDiagnosa,
	klaimID, pesertaID string, diagID int64) (models.Diagnosa, error) {

	perPeserta, err := k.diag.AmbilDiagnosa(ctx, klaimID)
	if err != nil {
		return models.Diagnosa{}, err
	}
	for _, baris := range perPeserta[pesertaID] {
		if baris.ID == diagID {
			return baris, nil
		}
	}
	return models.Diagnosa{}, fmt.Errorf("%w: diagnosa %d bukan milik peserta %q",
		galat.ErrPermintaanTidakSah, diagID, pesertaID)
}

// cerminkanKeDiagnosa menyalin keputusan peserta ke seluruh diagnosanya.
//
// `[terverifikasi]` `Activity/SetSTS_Reject.xml` - lihat
// `repository.Diagnosa.CerminkanStsReject` untuk pohon langkahnya.
//
// ⛔ Ia menerima `tx`, bukan membuka transaksinya sendiri: pencerminan yang
// menyusul di transaksi terpisah dapat gagal sendirian, dan peserta yang
// sudah ditolak akan berdiri di samping diagnosa yang masih tampak
// Outstanding. Dipanggil dari `statusbaris.go`, satu-satunya tempat
// `STS_REJECT` peserta berubah.
//
// ⚠️ Fungsi paket, bukan metode `DiagnosaPeserta` - pemanggilnya sudah berada
// di dalam transaksi dan tidak boleh melewati gerbang `pagari` lagi: gerbang
// itu MENOLAK peserta yang sudah diputus, sedangkan inilah yang memutuskannya.
func cerminkanKeDiagnosa(ctx context.Context, svc *Service, tx *db.Tx,
	pesertaID, sts string) error {

	return repository.NewDiagnosa(svc.DB()).CerminkanStsReject(ctx, tx, pesertaID, sts)
}
