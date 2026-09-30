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

	"nusantarare/inti/db"
	"nusantarare/inti/layanan"
	"nusantarare/inti/outbox"
	"nusantarare/inti/unggah"
	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/repository"
	"nusantarare/modul/treatycontractout/services"
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
	sisa      time.Duration
	minimum   []time.Duration
	simpanErr error
	dibaca    int
	disimpan  []tokenTersimpanUji
}

func (p *penyimpanTokenUji) AppStorage(context.Context) (string, error) {
	p.dibaca++
	return p.app, p.appErr
}

// TokenBerlaku meniru `INPUTDATE > saat + sisaMinimum` kueri Oracle-nya.
func (p *penyimpanTokenUji) TokenBerlaku(_ context.Context, _ *db.Tx, _ string, _ time.Time,
	sisaMinimum time.Duration) (string, time.Duration, error) {
	p.dibaca++
	p.minimum = append(p.minimum, sisaMinimum)
	if p.lama == "" || p.sisa <= sisaMinimum {
		return "", 0, nil
	}
	return p.lama, p.sisa, nil
}
func (p *penyimpanTokenUji) SimpanToken(_ context.Context, _ *db.Tx, app, token, pengguna string, sampai time.Time) error {
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

	lama := &penyimpanTokenUji{app: appUjiStorage, lama: tokenUjiStorage, sisa: 40 * time.Second}
	tok, exp, err := services.NewSumberTokenStorageTCO(transaksiUji, lama, garamUjiStorage, jam).TokenBaru(ctx)
	if err != nil || tok != tokenUjiStorage || !exp.Equal(saat.Add(40*time.Second)) || len(lama.disimpan) != 0 {
		t.Errorf("pakai ulang: %q %v %v, disimpan %d", tok, exp, err, len(lama.disimpan))
	}
	if len(lama.minimum) != 1 || lama.minimum[0] != services.MarginTokenTCO {
		t.Errorf("ambang pakai ulang %v, mau MarginTokenTCO", lama.minimum)
	}

	// AC 60: token yang sisa umurnya di dalam margin TIDAK dipakai ulang -
	// tanpa ini cache menerima token hampir mati yang sama berulang kali.
	hampir := &penyimpanTokenUji{app: appUjiStorage, lama: tokenUjiStorage, sisa: 10 * time.Second}
	tok, exp, err = services.NewSumberTokenStorageTCO(transaksiUji, hampir, garamUjiStorage, jam).TokenBaru(ctx)
	if err != nil || tok == tokenUjiStorage || !exp.Equal(saat.Add(time.Minute)) || len(hampir.disimpan) != 1 {
		t.Errorf("token hampir mati dipakai ulang: %v %v, disimpan %d", exp, err, len(hampir.disimpan))
	}

	baru := &penyimpanTokenUji{app: appUjiStorage}
	tok, exp, err = services.NewSumberTokenStorageTCO(transaksiUji, baru, garamUjiStorage, jam).TokenBaru(ctx)
	mau, _ := layanan.RakitToken(garamUjiStorage, saat)
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
		if _, _, err := services.NewSumberTokenStorageTCO(transaksiUji, p, garam, jam).TokenBaru(ctx); !errors.Is(err, layanan.ErrGaramTokenKosong) || p.dibaca != 0 {
			t.Errorf("garam %q: %v, dibaca %d", garam, err, p.dibaca)
		}
	}
	if _, _, err := services.NewSumberTokenStorageTCO(transaksiUji, nil, garamUjiStorage, jam).TokenBaru(ctx); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("tanpa penyimpan: %v", err)
	}
	if _, _, err := services.NewSumberTokenStorageTCO(transaksiUji, &penyimpanTokenUji{}, garamUjiStorage, jam).TokenBaru(ctx); !errors.Is(err, layanan.ErrAppNameKosong) {
		t.Errorf("App kosong: %v", err)
	}
	p := &penyimpanTokenUji{app: appUjiStorage, simpanErr: errors.New("repository: menyimpan token penyimpanan untuk \"UJI-APP\"")}
	_, _, err := services.NewSumberTokenStorageTCO(transaksiUji, p, garamUjiStorage, jam).TokenBaru(ctx)
	if err == nil {
		t.Fatal("galat simpan ditelan")
	}
	mau, _ := layanan.RakitToken(garamUjiStorage, jam())
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
	if s := svc.DenganPenyimpananLampiranTCO(false, garamUjiStorage); s.PenyimpananLampiranNyataTCO() {
		t.Error("false menyalakan pelaksana nyata")
	}
	s := svc.DenganPenyimpananLampiranTCO(true, garamUjiStorage)
	k, jauh := services.PenyimpananLampiranTCO(s).(*services.PenyimpananJarakJauhTCO)
	if !s.PenyimpananLampiranNyataTCO() || !jauh {
		t.Fatal("pelaksana nyata tidak terpasang")
	}
	// Tanpa Oracle tidak ada alamat yang dapat di-resolve: gagal terang, tanpa panggilan keluar.
	if _, err := k.Simpan(context.Background(), kunciUjiStorage, strings.NewReader("x"), "", ""); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: %v", err)
	}
	if svc.PenyimpananLampiranNyataTCO() {
		t.Error("DenganPenyimpananLampiranTCO mengubah Service asal")
	}
}

// Keadaan permanen pelaksana nyata tidak diputar ulang oleh antrean.
func TestLampiranGalatStorageNyataPermanen(t *testing.T) {
	for _, gagal := range []error{
		fmt.Errorf("%w: status 400", services.ErrStorageMenolakPermintaanTCO),
		fmt.Errorf("%w: %w", services.ErrTokenPenyimpananGagal, layanan.ErrGaramTokenKosong),
		fmt.Errorf("%w: %w", services.ErrTokenPenyimpananGagal, repository.ErrAppStorageKosongTCO),
	} {
		r := rakitanLampiran(t)
		r.simpan.setelGagal(gagal, false)
		h := r.unggah(t, "1000001", "a.pdf", "ISI")
		if h.Lampiran.Status != models.StatusLampiranGagal || r.antrean.cacahStatus(unggah.JenisEfekStorageUnggah, outbox.StatusEfekGagalPermanen) != 1 {
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
		if !strings.Contains(s, "1 effect(s) run") {
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

// `ext` UploadDoc dari NAMA BERKAS ASLI, huruf kecil tanpa titik
// (InsertGoogleStorage_Act b587 `@toLowerCase(Param.Ext)`) - tidak ditebak dari MIME.
func TestLampiranEkstensiDariNamaBerkas(t *testing.T) {
	r := rakitanLampiran(t)
	r.unggah(t, "1000001", "Kontrak.Final.PDF", "ISI")
	r.unggah(t, "1000001", "tanpa-titik", "ISI")
	r.simpan.mu.Lock()
	defer r.simpan.mu.Unlock()
	if len(r.simpan.ekstensi) != 2 || r.simpan.ekstensi[0] != "pdf" || r.simpan.ekstensi[1] != "" {
		t.Errorf("ekstensi %q, mau [pdf \"\"]", r.simpan.ekstensi)
	}
}
