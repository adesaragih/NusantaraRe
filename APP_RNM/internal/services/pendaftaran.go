package services

// Pendaftaran klaim Life - tiket 02.
//
// Untuk apa berkas ini: membentuk satu klaim baru dari peserta yang dipilih
// pengguna, dalam SATU transaksi. Tiga tempat ditulis sekaligus - baris work
// object, header klaim, dan peserta terpilih - ditambah baris datar warisan,
// sebab dua sistem masih berjalan berdampingan (ADR-U-0042).
//
// Dibaca sesudah: services.go dan repository/pohonklaim.go.
//
// Istilah:
//   - work object : satu baris T_WORK_CLAIM; akar tangga klaim di Pega.
//   - nomor klaim : nomor BISNIS yang dibaca manusia, berbeda dari pengenal.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// ErrPenomorBelumDiputuskan menandai penomoran yang belum boleh ditulis.
var ErrPenomorBelumDiputuskan = errors.New(
	"services: cara membentuk nomor klaim belum diputuskan work owner")

// ErrPermintaanTidakSah menandai permintaan pendaftaran yang tidak lengkap.
var ErrPermintaanTidakSah = errors.New("services: permintaan pendaftaran tidak sah")

// Penomor membentuk NOMOR BISNIS klaim - bukan pengenal barisnya.
//
// ⛔ Kenapa ini antarmuka dan bukan fungsi: keputusan work owner o
// (26-09-2026) melarang memanggil stored procedure, sedangkan AC 2, 3, dan
// 7-11 tiket 02 masih MENUNTUT procedure itu. Sampai teks AC-nya ditulis
// ulang, menulis logikanya berarti memilih salah satu pihak diam-diam.
//
// Yang dilakukan di sini: memisahkan tempatnya, supaya seluruh jalur lain
// dapat dibangun dan diuji, dan supaya yang belum diputuskan terlihat sebagai
// satu galat terang - bukan sebagai nomor karangan yang tampak benar.
type Penomor interface {
	NomorBerikut(ctx context.Context, tx *repository.Tx, kodeBisnis string, saat time.Time) (string, error)
}

// PenomorBelumDiputuskan adalah implementasi bawaan sampai butir o dijawab.
type PenomorBelumDiputuskan struct{}

// NomorBerikut selalu gagal, dengan pesan yang menyebut apa yang ditunggu.
func (PenomorBelumDiputuskan) NomorBerikut(context.Context, *repository.Tx, string, time.Time) (string, error) {
	return "", fmt.Errorf("%w: butir o melarang memanggil "+
		"POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER, sedangkan AC 2, 3, dan 7-11 tiket 02 "+
		"masih menuntutnya. Teks AC itu perlu ditulis ulang work owner lebih dulu",
		ErrPenomorBelumDiputuskan)
}

// PermintaanDaftar adalah isi satu pendaftaran klaim.
type PermintaanDaftar struct {
	NomorPremiList string
	NomorPolis     string
	Type           string
	KodeBisnis     string
	MataUang       string
	// Sertifikat adalah NOMOR SERTIFIKAT peserta yang dipilih pengguna.
	//
	// ⛔ Hanya nomornya. Nilai polis - uang, tanggal valuasi, share - dibaca
	// server sendiri dari sumbernya, sebab nilai itu menentukan angka klaim
	// dan jendela DOL. Menerimanya dari badan permintaan berarti siapa pun
	// yang dapat mengirim permintaan dapat menentukannya.
	Sertifikat []string
}

// Periksa memastikan permintaan cukup untuk membentuk klaim.
//
// Dipisah dari Daftar supaya dapat diuji tanpa Oracle, dan supaya jalur yang
// DITOLAK punya ujinya sendiri.
func (p PermintaanDaftar) Periksa() error {
	kosong := func(s string) bool { return strings.TrimSpace(s) == "" }
	switch {
	case kosong(p.NomorPremiList):
		return fmt.Errorf("%w: nomor premium list wajib terisi", ErrPermintaanTidakSah)
	case kosong(p.Type):
		return fmt.Errorf("%w: Type wajib terisi; ia yang menentukan jendela DOL", ErrPermintaanTidakSah)
	case kosong(p.KodeBisnis):
		return fmt.Errorf("%w: kode bisnis wajib terisi; ia yang menentukan prefix nomor", ErrPermintaanTidakSah)
	case len(p.Sertifikat) == 0:
		return fmt.Errorf("%w: nol peserta dipilih; klaim tanpa peserta tidak berarti apa-apa",
			ErrPermintaanTidakSah)
	}
	return nil
}

// Pendaftaran adalah layanan pendaftaran klaim baru.
type Pendaftaran struct {
	svc     *Service
	penomor Penomor
}

// Pendaftaran menyusun layanan itu dengan penomor bawaan.
func (s *Service) Pendaftaran() *Pendaftaran {
	return &Pendaftaran{svc: s, penomor: PenomorBelumDiputuskan{}}
}

// DenganPenomor mengganti penomornya - dipakai test, dan kelak oleh butir o.
func (p *Pendaftaran) DenganPenomor(n Penomor) *Pendaftaran {
	return &Pendaftaran{svc: p.svc, penomor: n}
}

// Daftar membentuk satu klaim baru dalam SATU transaksi.
//
// Urutannya disengaja: permintaan diperiksa lebih dulu, baru nomor diambil.
// Nomor yang sudah terbentuk tidak dapat dikembalikan ke urutannya, jadi
// mengambilnya sebelum validasi berarti membuang satu nomor setiap kali
// pengguna salah mengisi borang.
func (p *Pendaftaran) Daftar(ctx context.Context, pelaku Pelaku, minta PermintaanDaftar) (
	models.PohonKlaim, error) {
	var hasil models.PohonKlaim

	// ⛔ FAIL-CLOSED atas pelaku anonim. Tanpa ini, permintaan tanpa pelaku -
	// keadaan BAWAAN saat stub mati, yaitu keadaan produksi - tetap menulis
	// klaim, hanya dengan CREATE_OP_NAME kosong. Itu bukan pagar yang gagal
	// tertutup melainkan pagar yang gagal ANONIM: jejaknya hilang, dan tidak
	// ada yang tahu siapa mendaftarkan klaim itu.
	//
	// ⛔ Identitas, lalu PERAN. Gerbang perannya dipasang tiket 07 - komentar
	// di sini sebelumnya menandainya `[terbuka - tiket 07]`, dan tiket itu
	// sekarang sudah lewat. `[terverifikasi]` `Flow/Register_Flow.xml`:
	// `Assignment2` (Input Register) dipegang `ReasLifeAdmin`.
	//
	// ⚠️ Pendaftaran MENULIS status Outstanding lewat TandaiOutstanding, jadi
	// ia salah satu jalur pengubah status - dan jalur pengubah status yang
	// tidak bergerbang adalah lubang, bukan kelonggaran.
	if err := WajibIdentitas(pelaku); err != nil {
		return hasil, fmt.Errorf("%w: klaim mencatat siapa yang membuatnya", err)
	}
	if err := WajibPeran(pelaku, PeranInputRegister); err != nil {
		return hasil, err
	}
	if err := minta.Periksa(); err != nil {
		return hasil, err
	}
	if !p.svc.PunyaDatabase() {
		return hasil, repository.ErrTanpaOracle
	}

	err := p.svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		pohon := repository.NewPohonKlaim(p.svc.db)

		// ⛔ Urutannya: NOMOR dulu, baru pengenal. Alasan yang sama dengan
		// "validasi sebelum nomor" di atas, dan ia berlaku persis di sini:
		// SEQ_WORK_CLAIM bersifat non-transaksional, jadi angka yang sudah
		// diambil TIDAK kembali saat transaksi dibatalkan. Selama butir o
		// belum diputuskan, penomoran SELALU gagal - dan urutan yang terbalik
		// membakar satu pengenal pada setiap permintaan yang pasti ditolak.
		nomor, err := p.penomor.NomorBerikut(ctx, tx, minta.KodeBisnis, time.Now())
		if err != nil {
			return err
		}
		pengenal, err := pohon.PengenalWorkBerikut(ctx, tx, repository.AwalanKlaim)
		if err != nil {
			return err
		}

		// ⛔ Peserta dibaca ULANG dari sumbernya, bukan diterima dari klien.
		peserta, err := repository.NewPesertaPolis(p.svc.db).
			AmbilUntukKlaim(ctx, minta.NomorPremiList, minta.Sertifikat)
		if err != nil {
			return err
		}

		hasil = models.PohonKlaim{
			Work: models.WorkClaim{
				ID:   pengenal,
				Lini: models.LiniLife,
				Type: minta.Type,
				// `[keputusan work owner 26-09-2026, butir ae1]` CASEID =
				// pengenal work object.
				//
				// ⚠️ Ronde sebelumnya saya menolak ini dengan alasan "satu
				// nilai dua arti" - dan itu keliru: di Pega pun CASEID ADALAH
				// pengenal work object-nya, jadi menyamakannya adalah paritas
				// dengan sistem berjalan. Membiarkannya kosong justru yang
				// merusak: baris datar warisan klaim baru tidak dapat
				// dikelompokkan hilir yang membaca per CASEID, dan Hapus serta
				// CacahBarisLama memakai sumbu itu.
				CaseID:       pengenal,
				CreateOpName: pelaku.AkunID,
			},
			Klaim: models.Klaim{
				ID:         pengenal,
				NomorKlaim: nomor,
				NomorPolis: minta.NomorPolis,
				ClaimRetro: models.Money{Currency: minta.MataUang},
				Peserta:    peserta,
			},
		}
		return pohon.Simpan(ctx, tx, hasil)
	})
	if err != nil {
		return models.PohonKlaim{}, err
	}
	return hasil, nil
}

// Kunci penghitung nomor klaim Life.
//
// `[data DBA]` `GENERATE_SEQUENCE_NUMBER` berkunci `(CLASS, JENIS, TAHUN)`.
// `CLASS` dan `JENIS` menentukan seri; `TAHUN` memisahkan tahun buku.
const (
	ClassNomorKlaimLife = "CLAIMLIFE"
	JenisNomorKlaimLife = "K"
)

// ⛔ KONSTANTA `AwalanNomorKlaim` DIBUANG, ronde 27-09-2026.
//
// Ia bernilai `"RNML-"` - benar untuk lingkungan yang kebetulan dipakai
// saat kode ditulis, dan AC tiket 02 melarangnya dengan kalimat yang tidak
// dapat disalahartikan: *"Prefix diperoleh lewat lookup ke
// `POOLDATA.KODE_PRODUKSI` (`TYPE='LIFE'`), tidak ditanam sebagai konstanta
// di kode."* Awalan itu MILIK basis data.
//
// `[terverifikasi]` `Claim Life/RDBList/GetKodeProdLife_SQL.xml` baris 85.
// Bentuk `<awalan>K<kode bisnis>.MM.YYYY.<5 digit>` sejajar dengan nomor
// akseptasi `RNML-A…`/`RNML-AR…` (`Generate_NoAccept_Life.xml` baris 85);
// huruf tengahnya yang membedakan seri - `K` klaim, `A` akseptasi.
//
// ⚠️ Nomor AKSEPTASI masih merakit awalannya sendiri (`AwalanAkseptasiGross`
// / `AwalanAkseptasiRetro` di akseptasi.go). Itu dinyatakan, bukan
// dilewatkan: awalan akseptasi TIDAK ada di `KODE_PRODUKSI` - tabel itu
// hanya punya `RNML-`, sedangkan akseptasi menuntut `RNML-A` dan `RNML-AR`.
// Menyatukannya menuntut keputusan work owner tentang dari mana huruf
// `A`/`AR` datang, dan itu `[terbuka]`.

// RakitNomorKlaim menyusun nomor bisnis klaim.
//
// ⚠️ LIMA digit seperti `LPAD(v_seq,5,'0')`, dan urut yang sudah lebih panjang
// TIDAK dipotong - memotongnya menerbitkan nomor yang bertabrakan.
//
// ⚠️ `awalan` DITERIMA, tidak diambil sendiri: fungsi ini murni supaya
// bentuk nomornya dapat diuji tanpa Oracle. Yang mengambilnya penomornya.
func RakitNomorKlaim(awalan, kodeBisnis, mmYYYY string, urut int) string {
	u := strconv.Itoa(urut)
	if len(u) < 5 {
		u = strings.Repeat("0", 5-len(u)) + u
	}
	return awalan + JenisNomorKlaimLife + kodeBisnis + "." + mmYYYY + "." + u
}

// penomorCounter menulis ulang `PROC_GENERATE_SEQUENCE_NUMBER` di Go.
//
// ⛔ `[keputusan work owner]` butir **o1**: procedure tidak dipanggil. Yang
// tetap di Oracle hanya `SELECT … FOR UPDATE`, sebab kunci baris memang milik
// basis data.
type penomorCounter struct{ pohon *repository.PohonKlaim }

// PenomorCounterOracle menyusun penomor yang memakai penghitung Oracle.
func PenomorCounterOracle(svc *Service) Penomor {
	return penomorCounter{pohon: repository.NewPohonKlaim(svc.db)}
}

// NomorBerikut menerbitkan satu nomor klaim baru.
//
// Urutannya: awalan → hari tutup buku → periode → kunci baris → naikkan →
// rakit. Empat langkah tengahnya persis procedure-nya; awalan di depan
// sebab nomor tanpa awalan bukan nomor.
func (p penomorCounter) NomorBerikut(ctx context.Context, tx *repository.Tx,
	kodeBisnis string, saat time.Time) (string, error) {

	if strings.TrimSpace(kodeBisnis) == "" {
		return "", fmt.Errorf("%w: kode bisnis kosong; nomor klaim memuatnya",
			ErrPermintaanTidakSah)
	}
	// ⛔ Awalannya di-LOOKUP, bukan konstanta - lihat catatan di atas.
	awalan, err := p.pohon.AwalanProduksi(ctx, tx, repository.TipeKodeProduksiLife)
	if err != nil {
		return "", err
	}
	hariClosing, err := p.pohon.HariClosing(ctx, tx)
	if err != nil {
		return "", err
	}
	periode := repository.HitungPeriodeNomor(saat, hariClosing)
	urut, err := p.pohon.UrutNomorBerikut(ctx, tx,
		ClassNomorKlaimLife, JenisNomorKlaimLife, periode, saat)
	if err != nil {
		return "", err
	}
	return RakitNomorKlaim(awalan, kodeBisnis, periode.MMYYYY, urut), nil
}
