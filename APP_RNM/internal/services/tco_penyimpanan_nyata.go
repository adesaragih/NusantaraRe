package services

// Pelaksana penyimpanan lampiran + pekerja latar - OQ-TCO-08/09 Treaty
// Contract Out [keputusan work owner 29-09-2026: OQ-TCO-08 "sekarang",
// OQ-TCO-09 "perlu"].
//
// Untuk apa berkas ini:
//
//  1. `PenyimpananLampiranTCO` memilih pelaksana efek penyimpanan menurut
//     `PELAKSANA_STORAGE` (dibaca `internal/config`, dipasang `cmd/api`):
//     `stub` (BAWAAN) = folder lokal; `nyata` = rangkaian jarak jauh -
//     alamat dari `M_LINK_SERVICE` saat jalan, token dari `GCP_IMAGE` /
//     rumus `GET_TOKEN_STORAGE` (`RakitToken`, garam `STORAGE_TOKEN_SALT`),
//     transport `tco_pengirim_storage.go`.
//  2. `JalankanPekerja` menjalankan `JalankanAntrean` berkala (interval dari
//     env; nol = mati, bawaan di uji). Coba ulang dan anti-dobel memakai yang
//     sudah ada: `Backoff`/`percobaanMaksimum` outbox, kunci idempoten
//     `IMAGEID`, dan `FOR UPDATE SKIP LOCKED` pemungutan.
//
// ⛔ Garam tidak pernah dicetak, dicatat, atau masuk pesan galat; token juga.
//
// Dibaca sesudah: tco_penyimpanan.go, tokenstorage.go, tco_lampiran.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/internal/repository"
)

// AkunPekerjaLampiranTCO - pelaku jejak kerja antrean oleh pekerja latar.
const AkunPekerjaLampiranTCO = "SISTEM-PEKERJA-TCO"

// batasPutaranPekerjaTCO - efek paling banyak per ketukan pekerja.
const batasPutaranPekerjaTCO = 50

// pengaturanPenyimpananTCO dipasang SEKALI saat proses menyala.
//
// ⛔ Garam TIDAK disimpan di sini - hanya sumber token di dalam cache yang
// memegangnya.
type pengaturanPenyimpananTCO struct {
	nyata bool
	cache *CacheTokenTCO
}

// DenganPenyimpananLampiranTCO memasang pelaksana penyimpanan lampiran.
//
// ⚠️ Dipanggil sekali dari `cmd/api`: `nyata` diputuskan `internal/config`
// (`PELAKSANA_STORAGE`), garam = `config.StorageTokenSalt`.
func (s *Service) DenganPenyimpananLampiranTCO(nyata bool, garam string) *Service {
	salin := *s
	p := &pengaturanPenyimpananTCO{nyata: nyata}
	if p.nyata {
		var penyimp PenyimpanTokenStorageTCO
		if salin.db != nil {
			penyimp = penyimpanTokenOracleTCO{db: salin.db}
		}
		p.cache = NewCacheTokenTCO(NewSumberTokenStorageTCO(salin.DalamTransaksi, penyimp, garam, time.Now),
			time.Now, MarginTokenTCO)
	}
	salin.penyimpananTCO = p
	return &salin
}

// PenyimpananLampiranNyataTCO menyatakan pelaksana nyata terpasang.
func (s *Service) PenyimpananLampiranNyataTCO() bool {
	return s.penyimpananTCO != nil && s.penyimpananTCO.nyata
}

// PenyimpananLampiranTCO - klien penyimpanan yang handler dan pekerja pakai.
func PenyimpananLampiranTCO(svc *Service) KlienPenyimpananTCO {
	if !svc.PenyimpananLampiranNyataTCO() {
		return PenyimpananLokalTCO(svc)
	}
	app := func(ctx context.Context) (string, error) {
		if svc.db == nil {
			return "", repository.ErrTanpaOracle
		}
		return svc.db.AppStorageTCO(ctx)
	}
	// ⚠️ `ResolverLinkServiceOracle` menuntut db; tanpa Oracle gagal terang.
	var resolver ResolverEndpoint = resolverTanpaOracleTCO{}
	if svc.db != nil {
		resolver = ResolverLinkServiceOracle(svc)
	}
	return NewPenyimpananJarakJauhTCO(resolver, svc.penyimpananTCO.cache, NewPengirimBerkasHTTPTCO(nil, app))
}

type resolverTanpaOracleTCO struct{}

func (resolverTanpaOracleTCO) Resolve(context.Context, KunciLayanan) (string, error) {
	return "", repository.ErrTanpaOracle
}

// PenyimpanTokenStorageTCO - bacaan dan tulisan token yang sumber token pakai.
type PenyimpanTokenStorageTCO interface {
	AppStorage(ctx context.Context) (string, error)
	TokenBerlaku(ctx context.Context, tx *repository.Tx, app string, saat time.Time,
		sisaMinimum time.Duration) (string, time.Duration, error)
	SimpanToken(ctx context.Context, tx *repository.Tx, app, token, pengguna string, sampai time.Time) error
}

type penyimpanTokenOracleTCO struct{ db *repository.DB }

func (p penyimpanTokenOracleTCO) AppStorage(ctx context.Context) (string, error) {
	return p.db.AppStorageTCO(ctx)
}
func (p penyimpanTokenOracleTCO) TokenBerlaku(ctx context.Context, tx *repository.Tx, app string, saat time.Time,
	sisaMinimum time.Duration) (string, time.Duration, error) {
	return p.db.TokenStorageBerlakuTCO(ctx, tx, app, saat, sisaMinimum)
}
func (p penyimpanTokenOracleTCO) SimpanToken(ctx context.Context, tx *repository.Tx, app, token, pengguna string, sampai time.Time) error {
	return repository.NewPohonKlaim(p.db).SimpanToken(ctx, tx, app, token, pengguna, sampai)
}

// sumberTokenStorageTCO - `GET_TOKEN_STORAGE` ditiru (lihat `tokenstorage.go`):
// token berlaku dipakai ulang BESERTA sisa umurnya; bila tidak ada, token baru
// dirakit dengan garam dan disimpan dengan umur satu menit.
//
// ⚠️ PENYIMPANGAN SADAR KECIL: Pega memakai ulang token apa pun yang
// `INPUTDATE > SYSDATE`. Di sini hanya yang sisa umurnya > `MarginTokenTCO`
// (AC 60) - tanpa itu cache yang menyegarkan di jendela margin menerima token
// yang sama yang hampir mati, berulang, satu transaksi per panggilan.
type sumberTokenStorageTCO struct {
	transaksi PenjalanTransaksiTCO
	penyimp   PenyimpanTokenStorageTCO
	garam     string
	jam       func() time.Time
}

// PenjalanTransaksiTCO - bentuk `(*Service).DalamTransaksi`; uji memasang
// penjalan tiruan yang memanggil fn dengan tx nil.
type PenjalanTransaksiTCO func(ctx context.Context, fn func(tx *repository.Tx) error) error

// NewSumberTokenStorageTCO menyusun sumber token. `penyimp` nil = tanpa Oracle.
func NewSumberTokenStorageTCO(transaksi PenjalanTransaksiTCO, penyimp PenyimpanTokenStorageTCO,
	garam string, jam func() time.Time) SumberTokenTCO {
	return sumberTokenStorageTCO{transaksi: transaksi, penyimp: penyimp, garam: garam, jam: jam}
}

// TokenBaru menerbitkan (atau memakai ulang) token dalam SATU transaksi.
func (s sumberTokenStorageTCO) TokenBaru(ctx context.Context) (string, time.Time, error) {
	if strings.TrimSpace(s.garam) == "" {
		return "", time.Time{}, ErrGaramTokenKosong
	}
	penyimp := s.penyimp
	if penyimp == nil || s.transaksi == nil {
		return "", time.Time{}, repository.ErrTanpaOracle
	}
	jam := s.jam
	if jam == nil {
		jam = time.Now
	}
	app, err := penyimp.AppStorage(ctx)
	if err != nil {
		return "", time.Time{}, err
	}
	if strings.TrimSpace(app) == "" {
		return "", time.Time{}, ErrAppNameKosong
	}
	var tok string
	var sampai time.Time
	err = s.transaksi(ctx, func(tx *repository.Tx) error {
		saat := jam()
		lama, sisa, err := penyimp.TokenBerlaku(ctx, tx, app, saat, MarginTokenTCO)
		if err != nil {
			return err
		}
		if lama != "" {
			tok, sampai = lama, saat.Add(sisa)
			return nil
		}
		baru, err := RakitToken(s.garam, saat)
		if err != nil {
			return err
		}
		// ⛔ Galat repositori sudah tanpa token (lihat repository/tokenstorage.go).
		if err := penyimp.SimpanToken(ctx, tx, app, baru, penggunaTokenBawaan, saat.Add(umurToken)); err != nil {
			return err
		}
		tok, sampai = baru, saat.Add(umurToken)
		return nil
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return tok, sampai, nil
}

// JalankanPekerja menjalankan antrean lampiran setiap `interval` sampai ctx
// selesai. Interval <= 0 = pekerja mati (bawaan; juga di uji).
//
// `catat` menerima ringkasan per ketukan yang menjalankan sesuatu - jumlah dan
// galat ringkas, tidak pernah alamat, token, atau garam.
func (l *LampiranTahunTCO) JalankanPekerja(ctx context.Context, interval time.Duration, catat func(string)) {
	if interval <= 0 {
		return
	}
	if catat == nil {
		catat = func(string) {}
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			n, err := l.JalankanAntrean(ctx, AkunPekerjaLampiranTCO, batasPutaranPekerjaTCO)
			if err != nil && !errors.Is(err, context.Canceled) {
				catat(fmt.Sprintf("pekerja lampiran TCO: %d efek dijalankan, lalu berhenti: %v", n, err))
			} else if n > 0 {
				catat(fmt.Sprintf("pekerja lampiran TCO: %d efek dijalankan", n))
			}
		}
	}
}
