package services_test

// Uji pelaksana penyimpanan lampiran + pekerja latar (OQ-TCO-08/09) - TANPA
// Oracle, TANPA layanan sungguhan. Garam palsu (awalan UJI-).

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

type tokenTersimpanUji struct {
	app, token, pengguna string
	sampai               time.Time
}

// penyimpanTokenUji meniru `T_FOLDER_IMAGE` + `GCP_IMAGE`.
type penyimpanTokenUji struct {
	app       string
	appErr    error
	lama      string
	habis     time.Time
	simpanErr error
	dibaca    int
	disimpan  []tokenTersimpanUji
}

func (p *penyimpanTokenUji) AppStorage(context.Context) (string, error) {
	p.dibaca++
	return p.app, p.appErr
}
func (p *penyimpanTokenUji) TokenBerlaku(_ context.Context, _ *repository.Tx, _ string, _ time.Time) (string, time.Time, error) {
	p.dibaca++
	return p.lama, p.habis, nil
}
func (p *penyimpanTokenUji) SimpanToken(_ context.Context, _ *repository.Tx, app, token, pengguna string, sampai time.Time) error {
	if p.simpanErr != nil {
		return p.simpanErr
	}
	p.disimpan = append(p.disimpan, tokenTersimpanUji{app, token, pengguna, sampai})
	return nil
}

// `GET_TOKEN_STORAGE` ditiru: pakai ulang yang berlaku BESERTA kedaluwarsanya,
// atau terbitkan baru `ASMAPP + garam + stempel` dengan umur satu menit.
func TestSumberTokenStoragePakaiUlangAtauTerbitkan(t *testing.T) {
	saat := time.Date(2026, 9, 29, 9, 0, 0, 0, time.UTC)
	jam := func() time.Time { return saat }
	ctx := context.Background()

	lama := &penyimpanTokenUji{app: appUjiStorage, lama: tokenUjiStorage, habis: saat.Add(40 * time.Second)}
	tok, exp, err := services.NewSumberTokenStorageTCO(transaksiUji, lama, garamUjiStorage, jam).TokenBaru(ctx)
	if err != nil || tok != tokenUjiStorage || !exp.Equal(saat.Add(40*time.Second)) || len(lama.disimpan) != 0 {
		t.Errorf("pakai ulang: %q %v %v, disimpan %d", tok, exp, err, len(lama.disimpan))
	}

	baru := &penyimpanTokenUji{app: appUjiStorage}
	tok, exp, err = services.NewSumberTokenStorageTCO(transaksiUji, baru, garamUjiStorage, jam).TokenBaru(ctx)
	mau, _ := services.RakitToken(garamUjiStorage, saat)
	if err != nil || tok != mau || !exp.Equal(saat.Add(time.Minute)) {
		t.Fatalf("terbitkan: %v %v", exp, err)
	}
	if len(baru.disimpan) != 1 || baru.disimpan[0] != (tokenTersimpanUji{appUjiStorage, mau, "Job", saat.Add(time.Minute)}) {
		t.Errorf("baris GCP_IMAGE: %d baris", len(baru.disimpan))
	}
}

func TestSumberTokenStorageGagalTerangTanpaMembocorkan(t *testing.T) {
	jam := func() time.Time { return time.Date(2026, 9, 29, 9, 0, 0, 0, time.UTC) }
	ctx := context.Background()
	for _, garam := range []string{"", "   "} {
		p := &penyimpanTokenUji{app: appUjiStorage}
		if _, _, err := services.NewSumberTokenStorageTCO(transaksiUji, p, garam, jam).TokenBaru(ctx); !errors.Is(err, services.ErrGaramTokenKosong) || p.dibaca != 0 {
			t.Errorf("garam %q: %v, dibaca %d", garam, err, p.dibaca)
		}
	}
	if _, _, err := services.NewSumberTokenStorageTCO(transaksiUji, nil, garamUjiStorage, jam).TokenBaru(ctx); !errors.Is(err, repository.ErrTanpaOracle) {
		t.Errorf("tanpa penyimpan: %v", err)
	}
	if _, _, err := services.NewSumberTokenStorageTCO(transaksiUji, &penyimpanTokenUji{}, garamUjiStorage, jam).TokenBaru(ctx); !errors.Is(err, services.ErrAppNameKosong) {
		t.Errorf("App kosong: %v", err)
	}
	p := &penyimpanTokenUji{app: appUjiStorage, simpanErr: errors.New("repository: menyimpan token penyimpanan untuk \"UJI-APP\"")}
	_, _, err := services.NewSumberTokenStorageTCO(transaksiUji, p, garamUjiStorage, jam).TokenBaru(ctx)
	if err == nil {
		t.Fatal("galat simpan ditelan")
	}
	mau, _ := services.RakitToken(garamUjiStorage, jam())
	if strings.Contains(err.Error(), garamUjiStorage) || strings.Contains(err.Error(), mau) {
		t.Errorf("galat memuat garam/token: %v", err)
	}
}

// Bawaan STUB; `nyata` hanya bila disetel (OQ-TCO-08, keputusan work owner).
func TestPenyimpananLampiranPilihanPelaksana(t *testing.T) {
	svc := services.New(nil)
	if svc.PenyimpananLampiranNyataTCO() {
		t.Fatal("bawaan bukan stub")
	}
	if _, jauh := services.PenyimpananLampiranTCO(svc).(*services.PenyimpananJarakJauhTCO); jauh {
		t.Error("bawaan memasang rangkaian jarak jauh")
	}
	for _, v := range []string{"", "stub", "STUB", "lain"} {
		s := svc.DenganPenyimpananLampiranTCO(v, garamUjiStorage)
		if s.PenyimpananLampiranNyataTCO() {
			t.Errorf("%q menyalakan pelaksana nyata", v)
		}
	}
	for _, v := range []string{"nyata", " NYATA "} {
		s := svc.DenganPenyimpananLampiranTCO(v, garamUjiStorage)
		k, jauh := services.PenyimpananLampiranTCO(s).(*services.PenyimpananJarakJauhTCO)
		if !s.PenyimpananLampiranNyataTCO() || !jauh {
			t.Errorf("%q: pelaksana nyata tidak terpasang", v)
			continue
		}
		// Tanpa Oracle tidak ada alamat yang dapat di-resolve: gagal, tanpa panggilan keluar.
		if err := k.Simpan(context.Background(), kunciUjiStorage, strings.NewReader("x"), ""); err == nil {
			t.Errorf("%q tanpa Oracle: simpan berhasil", v)
		}
	}
	if svc.PenyimpananLampiranNyataTCO() {
		t.Error("DenganPenyimpananLampiranTCO mengubah Service asal")
	}
}

// Keadaan permanen pelaksana nyata tidak diputar ulang oleh antrean.
func TestLampiranGalatStorageNyataPermanen(t *testing.T) {
	for _, gagal := range []error{
		fmt.Errorf("%w: status 400", services.ErrStorageMenolakPermintaanTCO),
		fmt.Errorf("%w: %w", services.ErrTokenPenyimpananGagal, services.ErrGaramTokenKosong),
		fmt.Errorf("%w: %w", services.ErrTokenPenyimpananGagal, repository.ErrAppStorageKosongTCO),
	} {
		r := rakitanLampiran(t)
		r.simpan.setelGagal(gagal, false)
		h := r.unggah(t, "1000001", "a.pdf", "ISI")
		if h.Lampiran.Status != models.StatusLampiranGagal || r.gudang.jejakBeraksi(repository.AksiJejakMenyerah) != 1 {
			t.Errorf("%v: status %q", gagal, h.Lampiran.Status)
		}
	}
	for _, sementara := range []error{services.ErrStorageTakTerjangkauTCO, services.ErrStorageGagalTCO} {
		r := rakitanLampiran(t)
		r.simpan.setelGagal(sementara, false)
		if h := r.unggah(t, "1000001", "a.pdf", "ISI"); h.Lampiran.Status != models.StatusLampiranTertunda {
			t.Errorf("%v: status %q, mau tertunda", sementara, h.Lampiran.Status)
		}
	}
}

// OQ-TCO-09: interval <= 0 = mati; interval > 0 memungut antrean sampai ctx selesai.
func TestPekerjaLampiranTCO(t *testing.T) {
	for _, d := range []time.Duration{0, -time.Second} {
		selesai := make(chan struct{})
		go func() {
			rakitanLampiran(t).l.JalankanPekerja(context.Background(), d, nil)
			close(selesai)
		}()
		select {
		case <-selesai:
		case <-time.After(2 * time.Second):
			t.Fatalf("interval %v: pekerja tidak langsung kembali", d)
		}
	}

	r := rakitanLampiran(t)
	r.simpan.setelGagal(errors.New("sementara"), false)
	h := r.unggah(t, "1000001", "a.pdf", "ISI")
	r.simpan.setelGagal(nil, false)
	r.maju(24 * time.Hour)
	ctx, batal := context.WithCancel(context.Background())
	catatan := make(chan string, 8)
	selesai := make(chan struct{})
	go func() {
		r.l.JalankanPekerja(ctx, 5*time.Millisecond, func(s string) { catatan <- s })
		close(selesai)
	}()
	select {
	case s := <-catatan:
		if !strings.Contains(s, "1 efek dijalankan") {
			t.Errorf("catatan: %q", s)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("pekerja tidak menjalankan antrean")
	}
	batal()
	select {
	case <-selesai:
	case <-time.After(2 * time.Second):
		t.Fatal("pekerja tidak berhenti bersama ctx")
	}
	b, err := r.l.Daftar(context.Background(), pelakuUjiTCO, "1000001")
	if err != nil || len(b) != 1 || b[0].ID != h.Lampiran.ID || b[0].Status != models.StatusLampiranTerkirim {
		t.Errorf("sesudah pekerja: %+v %v", b, err)
	}
	if r.simpan.cacah() != 1 {
		t.Errorf("berkas di penyimpanan %d, mau 1", r.simpan.cacah())
	}
}
