// Package services memuat aturan dagang, batas transaksi, dan penegakan
// wewenang.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// pernah mengimpor handlers.
//
// Fase 0 - scaffold: nol aturan dagang. Yang ada di sini adalah batas
// transaksi dan bentuk penegakan wewenang yang harus diikuti setiap tiket.
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/internal/repository"
)

var (
	// ErrTanpaWewenang dikembalikan bila peran tidak mencukupi.
	ErrTanpaWewenang = errors.New("services: wewenang tidak mencukupi")
	// ErrWajibIsi dikembalikan bila medan wajib kosong. Seluruh kolom
	// nullable di database; kewajiban isi ditegakkan DI SINI (ADR-U-0027).
	ErrWajibIsi = errors.New("services: medan wajib belum terisi")
)

// Service adalah akar seluruh layanan.
type Service struct {
	db *repository.DB
	// lingkungan menggerbangi EFEK KELUAR saja, tidak pernah penyimpanan.
	//
	// ⛔ Bawaannya BUKAN produksi, dan itu disengaja: gagal tertutup, bukan
	// gagal terbuka. Proses yang lupa menyetelnya tidak akan mengirim email
	// kepada orang sungguhan.
	lingkungan Lingkungan
	// unggahanDir adalah folder lokal berkas unggahan (butir be). Kosong
	// berarti unggahan GAGAL TERANG - bukan bawaan diam-diam.
	unggahanDir string
	// penyimpananTCO - pelaksana penyimpanan lampiran Treaty Contract Out
	// (OQ-TCO-08). nil = stub; lihat tco_penyimpanan_nyata.go.
	penyimpananTCO *pengaturanPenyimpananTCO
}

// DenganUnggahanDir menyetel folder berkas unggahan.
//
// ⚠️ Dipanggil SEKALI saat proses menyala, dari `cmd/api`, dengan nilai
// dari `config.UnggahanDir`. Tanpa pemanggilan itu tidak ada unggahan yang
// pernah berhasil di mana pun - dan tidak adanya pemanggil itulah yang
// membuat AC lingkungan sempat tercentang secara hampa (lihat
// DenganLingkungan).
func (s *Service) DenganUnggahanDir(dir string) *Service {
	salin := *s
	salin.unggahanDir = dir
	return &salin
}

// New membuat Service. db boleh nil bila proses berjalan tanpa Oracle;
// layanan yang memerlukannya akan menolak saat dipanggil.
func New(db *repository.DB) *Service {
	return &Service{db: db, lingkungan: BukanProduksi}
}

// DenganLingkungan menyetel lingkungan efek keluarnya.
//
// ⚠️ Dipanggil SEKALI saat proses menyala, dari `cmd/api`, dengan nilai dari
// `config.IsPegaProd` (ADR-U-0005). Tanpa pemanggilan itu tidak ada efek
// keluar yang pernah berjalan di mana pun - dan tidak adanya pemanggil itulah
// yang membuat AC lingkungan sempat tercentang secara hampa.
func (s *Service) DenganLingkungan(l Lingkungan) *Service {
	salin := *s
	salin.lingkungan = l
	return &salin
}

// Lingkungan menyebut lingkungan efek keluar proses ini.
func (s *Service) Lingkungan() Lingkungan { return s.lingkungan }

// PunyaDatabase menyatakan apakah proses dikonfigurasi menyentuh Oracle.
func (s *Service) PunyaDatabase() bool { return s != nil && s.db != nil }

// CekKesehatan memeriksa apakah Oracle terjangkau.
//
// Ia ada DI SINI, bukan di handlers, supaya arah handlers -> services ->
// repository tidak dipotong: handlers tidak pernah memegang koneksi.
// Mengembalikan ErrTanpaOracle bila database memang tidak dikonfigurasi -
// itu keadaan yang sah di Fase 0, bukan kegagalan.
func (s *Service) CekKesehatan(ctx context.Context) error {
	if !s.PunyaDatabase() {
		return repository.ErrTanpaOracle
	}
	return s.db.Ping(ctx)
}

// SkemaAktif mengembalikan nama skema yang dipakai, atau teks kosong bila
// tidak ada database.
func (s *Service) SkemaAktif() string {
	if !s.PunyaDatabase() {
		return ""
	}
	return s.db.Skema()
}

// DalamTransaksi menjalankan fn di dalam satu transaksi.
//
// Transaksi dibuka dan ditutup DI SINI, bukan di dalam teks SQL
// (ADR-U-0029). Galat apa pun dari fn membatalkan seluruhnya.
func (s *Service) DalamTransaksi(ctx context.Context, fn func(tx *repository.Tx) error) (err error) {
	if s.db == nil {
		return repository.ErrTanpaOracle
	}
	tx, err := s.db.Mulai(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// Pelaku adalah identitas yang meminta sebuah tindakan.
//
// ⛔ Nol nama orang di kode (ADR-U-0030, CLAUDE.md bab 4 butir 10). Yang
// dibawa hanya pengenal akun dan peran; nama tidak pernah menjadi dasar
// keputusan dan tidak pernah ditulis ke kode atau fixture.
type Pelaku struct {
	AkunID string
	Peran  []string
}

// PunyaPeran memeriksa satu peran.
//
// Sumber peran adalah SATU tabel (ADR-U-0030). Tabel itu belum ada di Fase 0;
// pembacanya lahir bersama tiket yang memerlukannya, dan pemeriksaan tetap
// dilakukan di lapisan ini - tidak pernah di handlers, tidak pernah di SQL.
func (p Pelaku) PunyaPeran(peran string) bool {
	for _, x := range p.Peran {
		if x == peran {
			return true
		}
	}
	return false
}

// ErrTanpaIdentitas menandai permintaan tanpa pengenal akun.
//
// ⛔ DIPISAH dari ErrTanpaWewenang 26-09-2026. Satu galat yang berarti dua hal
// - "aku tidak tahu kamu siapa" dan "aku tahu kamu siapa, tetapi kamu tidak
// boleh" - memaksa pemanggil menebak, dan dua handler memang menerjemahkannya
// ke dua kode HTTP yang berbeda untuk galat yang sama. 401 dan 403 menjawab
// pertanyaan yang berbeda.
var ErrTanpaIdentitas = errors.New("services: permintaan tanpa identitas pelaku")

// WajibIdentitas menolak permintaan tanpa pengenal akun.
func WajibIdentitas(p Pelaku) error {
	if strings.TrimSpace(p.AkunID) == "" {
		return ErrTanpaIdentitas
	}
	return nil
}

// WajibPeran mengembalikan galat bila peran tidak dimiliki.
//
// Jalur yang DITOLAK wajib punya uji tersendiri, bukan hanya jalur yang
// berhasil (brief bab 5).
func WajibPeran(p Pelaku, peran string) error {
	if !p.PunyaPeran(peran) {
		return fmt.Errorf("%w: perlu peran %q", ErrTanpaWewenang, peran)
	}
	return nil
}
