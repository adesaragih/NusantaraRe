package services_test

// Uji transport penyimpanan lampiran (OQ-TCO-08) - SERVER TIRUAN LOKAL.
//
// ⛔ Tidak ada layanan sungguhan yang dipanggil: setiap alamat berasal dari
// `httptest.NewServer` saat uji berjalan; nol alamat literal di berkas ini.
// Token dan garam di sini palsu (awalan UJI-).

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

const (
	appUjiStorage   = "UJI-APP"
	tokenUjiStorage = "UJI-TOKEN-PALSU"
	garamUjiStorage = "UJI-GARAM-PALSU"
	kunciUjiStorage = "ABCDEF0123456789"
)

func appUjiStorageTCO(context.Context) (string, error) { return appUjiStorage, nil }

type permintaanStorageDiterima struct {
	jalur, metode, jenis string
	badan                map[string]any
}

// layananStorageUji meniru layanan penyimpanan: `/upload`, `/geturl`,
// `/delete`, dan objek bertanda tangan di `/objek/{kunci}`.
type layananStorageUji struct {
	mu        sync.Mutex
	objek     map[string][]byte
	diterima  []permintaanStorageDiterima
	paksa     map[string]int
	urlKosong bool
	srv       *httptest.Server
}

func layananStorageTiruan(t *testing.T) *layananStorageUji {
	t.Helper()
	l := &layananStorageUji{objek: map[string][]byte{}, paksa: map[string]int{}}
	l.srv = httptest.NewServer(http.HandlerFunc(l.layani))
	t.Cleanup(l.srv.Close)
	return l
}

func (l *layananStorageUji) alamat(jalur string) string { return l.srv.URL + jalur }

func inangDari(alamat string) string {
	u, _ := url.Parse(alamat)
	return u.Host
}

func (l *layananStorageUji) layani(w http.ResponseWriter, r *http.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if kode, ada := l.paksa[r.URL.Path]; ada {
		w.WriteHeader(kode)
		return
	}
	if nama, ada := strings.CutPrefix(r.URL.Path, "/objek/"); ada {
		isi, ada := l.objek[nama]
		if !ada {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write(isi)
		return
	}
	var badan map[string]any
	_ = json.NewDecoder(r.Body).Decode(&badan)
	l.diterima = append(l.diterima, permintaanStorageDiterima{jalur: r.URL.Path, metode: r.Method,
		jenis: r.Header.Get("Content-Type"), badan: badan})
	nama, _ := badan["Namafile"].(string)
	w.Header().Set("Content-Type", "application/json")
	switch r.URL.Path {
	case "/upload":
		gambar, _ := badan["Image"].(string)
		data, err := base64.StdEncoding.DecodeString(gambar)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		l.objek[nama] = data
		u := l.alamat("/objek/" + nama)
		if l.urlKosong {
			u = ""
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"URLImage": u})
	case "/geturl":
		// Seperti URL bertanda tangan: diberikan walau objeknya tidak ada.
		_ = json.NewEncoder(w).Encode(map[string]string{"URLImage": l.alamat("/objek/" + nama)})
	case "/delete":
		delete(l.objek, nama)
		_ = json.NewEncoder(w).Encode(map[string]string{})
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (l *layananStorageUji) cacahObjek() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.objek)
}

// Bentuk `UploadDoc` (InsertGoogleStorage_Act b1339-b1641) dan pengulangan ke
// objek yang SAMA (AC 58).
func TestPengirimStorageUnggahBentukUploadDoc(t *testing.T) {
	l := layananStorageTiruan(t)
	p := services.NewPengirimBerkasHTTPTCO(l.srv.Client(), appUjiStorageTCO)
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		if err := p.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage,
			strings.NewReader("ISI-PDF"), "application/pdf"); err != nil {
			t.Fatal(err)
		}
	}
	if l.cacahObjek() != 1 || string(l.objek[kunciUjiStorage]) != "ISI-PDF" {
		t.Fatalf("objek: %d, isi %q", l.cacahObjek(), l.objek[kunciUjiStorage])
	}
	d := l.diterima[0]
	if d.metode != http.MethodPost || d.jenis != "application/json" {
		t.Errorf("metode %q jenis %q", d.metode, d.jenis)
	}
	mau := map[string]any{"App": appUjiStorage, "Kodestring": tokenUjiStorage, "Namafile": kunciUjiStorage,
		"Folder": services.FolderStorageTCO, "Durasi": float64(services.DurasiURLStorageTCO),
		"MimeType": "application/pdf", "ext": "pdf"}
	for k, v := range mau {
		if d.badan[k] != v {
			t.Errorf("medan %s = %v, mau %v", k, d.badan[k], v)
		}
	}
}

func TestPengirimStorageAmbilPeriksaBuang(t *testing.T) {
	l := layananStorageTiruan(t)
	p := services.NewPengirimBerkasHTTPTCO(l.srv.Client(), appUjiStorageTCO)
	ctx := context.Background()
	if err := p.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage, strings.NewReader("ISI"), "text/plain"); err != nil {
		t.Fatal(err)
	}
	rc, err := p.Ambil(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage)
	if err != nil {
		t.Fatal(err)
	}
	isi, _ := io.ReadAll(rc)
	_ = rc.Close()
	if string(isi) != "ISI" {
		t.Errorf("isi %q", isi)
	}
	if ada, err := p.Periksa(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage); err != nil || !ada {
		t.Errorf("periksa sebelum buang: %v %v", ada, err)
	}
	if err := p.Buang(ctx, l.alamat("/delete"), tokenUjiStorage, kunciUjiStorage); err != nil {
		t.Fatal(err)
	}
	b := l.diterima[len(l.diterima)-1].badan
	if b["App"] != appUjiStorage || b["Kodestring"] != tokenUjiStorage || b["Namafile"] != kunciUjiStorage ||
		b["Image"] != nil || b["Folder"] != nil {
		t.Errorf("badan delete (DeleteGoogleStorage_Act b1019-b1091): %v", b)
	}
	if ada, err := p.Periksa(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage); err != nil || ada {
		t.Errorf("periksa sesudah buang: %v %v", ada, err)
	}
	if _, err := p.Ambil(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage); !errors.Is(err, services.ErrBerkasTidakAdaDiPenyimpanan) {
		t.Errorf("ambil sesudah buang: %v", err)
	}
	l.mu.Lock()
	l.paksa["/geturl"] = http.StatusNotFound
	l.mu.Unlock()
	if ada, err := p.Periksa(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage); err != nil || ada {
		t.Errorf("periksa geturl 404: %v %v", ada, err)
	}
}

// bersihDariRahasia - ⛔ galat tidak pernah memuat alamat, inang, atau token.
func bersihDariRahasia(t *testing.T, err error, inang string) {
	t.Helper()
	for _, larangan := range []string{inang, tokenUjiStorage, garamUjiStorage} {
		if larangan != "" && strings.Contains(err.Error(), larangan) {
			t.Errorf("galat memuat rahasia/alamat: %v", err)
		}
	}
}

func TestPengirimStorageGalatTanpaAlamatAtauToken(t *testing.T) {
	l := layananStorageTiruan(t)
	p := services.NewPengirimBerkasHTTPTCO(l.srv.Client(), appUjiStorageTCO)
	ctx := context.Background()
	for kode, mau := range map[int]error{
		http.StatusBadRequest:          services.ErrStorageMenolakPermintaanTCO,
		http.StatusUnprocessableEntity: services.ErrStorageMenolakPermintaanTCO,
		http.StatusNotFound:            services.ErrBerkasTidakAdaDiPenyimpanan,
		http.StatusUnauthorized:        services.ErrStorageGagalTCO,
		http.StatusTooManyRequests:     services.ErrStorageGagalTCO,
		http.StatusInternalServerError: services.ErrStorageGagalTCO,
	} {
		l.mu.Lock()
		l.paksa["/upload"] = kode
		l.mu.Unlock()
		err := p.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage, strings.NewReader("x"), "")
		if !errors.Is(err, mau) {
			t.Errorf("status %d: %v, mau %v", kode, err, mau)
			continue
		}
		bersihDariRahasia(t, err, inangDari(l.srv.URL))
	}
	l.mu.Lock()
	delete(l.paksa, "/upload")
	l.urlKosong = true
	l.mu.Unlock()
	if err := p.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage, strings.NewReader("x"), ""); !errors.Is(err, services.ErrStorageJawabanRusakTCO) {
		t.Errorf("URLImage kosong (b2902): %v", err)
	}

	mati := httptest.NewServer(http.NotFoundHandler())
	alamatMati := mati.URL + "/upload"
	mati.Close()
	err := p.Kirim(ctx, alamatMati, tokenUjiStorage, kunciUjiStorage, strings.NewReader("x"), "")
	if !errors.Is(err, services.ErrStorageTakTerjangkauTCO) {
		t.Fatalf("server mati: %v", err)
	}
	bersihDariRahasia(t, err, inangDari(alamatMati))

	// ⚠️ Pelepas eksplisit: server tidak selalu tahu klien memutus, dan
	// `Close` menunggu setiap handler selesai.
	lepas := make(chan struct{})
	lambat := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-lepas:
		}
	}))
	defer func() { close(lepas); lambat.Close() }()
	pl := services.NewPengirimBerkasHTTPTCO(&http.Client{Timeout: 50 * time.Millisecond}, appUjiStorageTCO)
	err = pl.Kirim(ctx, lambat.URL+"/upload", tokenUjiStorage, kunciUjiStorage, strings.NewReader("x"), "")
	if !errors.Is(err, services.ErrStorageTakTerjangkauTCO) || !strings.Contains(err.Error(), "batas waktu") {
		t.Errorf("batas waktu: %v", err)
	}
	bersihDariRahasia(t, err, inangDari(lambat.URL))

	tanpaApp := services.NewPengirimBerkasHTTPTCO(l.srv.Client(), func(context.Context) (string, error) {
		return "", repository.ErrAppStorageKosongTCO
	})
	l.mu.Lock()
	sebelum := len(l.diterima)
	l.mu.Unlock()
	if err := tanpaApp.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage, strings.NewReader("x"), ""); !errors.Is(err, repository.ErrAppStorageKosongTCO) {
		t.Errorf("tanpa App: %v", err)
	}
	if len(l.diterima) != sebelum {
		t.Error("permintaan terkirim padahal App belum ada")
	}
}

// Rangkaian utuh: resolver (M_LINK_SERVICE tiruan) -> cache token -> sumber
// token GET_TOKEN_STORAGE tiruan -> transport HTTP -> server tiruan lokal.
func TestPenyimpananNyataUjungKeUjung(t *testing.T) {
	l := layananStorageTiruan(t)
	saat := time.Date(2026, 9, 29, 9, 0, 0, 0, time.UTC)
	jam := func() time.Time { return saat }
	r := &resolverLampiranUji{alamat: map[services.KunciLayanan]string{
		services.KunciUnggahBerkas: l.alamat("/upload"),
		services.KunciURLBerkas:    l.alamat("/geturl"),
		services.KunciHapusBerkas:  l.alamat("/delete"),
	}}
	gudangToken := &penyimpanTokenUji{app: appUjiStorage}
	sumber := services.NewSumberTokenStorageTCO(transaksiUji, gudangToken, garamUjiStorage, jam)
	p := services.NewPenyimpananJarakJauhTCO(r, services.NewCacheTokenTCO(sumber, jam, services.MarginTokenTCO),
		services.NewPengirimBerkasHTTPTCO(l.srv.Client(), appUjiStorageTCO))
	ctx := context.Background()
	if err := p.Simpan(ctx, kunciUjiStorage, strings.NewReader("ISI"), "application/pdf"); err != nil {
		t.Fatal(err)
	}
	if ada, err := p.Ada(ctx, kunciUjiStorage); err != nil || !ada {
		t.Fatalf("ada: %v %v", ada, err)
	}
	if err := p.Hapus(ctx, kunciUjiStorage); err != nil {
		t.Fatal(err)
	}
	mau, _ := services.RakitToken(garamUjiStorage, saat)
	for _, d := range l.diterima {
		if d.badan["Kodestring"] != mau {
			t.Errorf("%s: Kodestring bukan token rakitan", d.jalur)
		}
	}
	if len(gudangToken.disimpan) != 1 {
		t.Errorf("token diterbitkan %d kali, mau 1 (dipakai ulang oleh cache)", len(gudangToken.disimpan))
	}
	if l.cacahObjek() != 0 {
		t.Error("objek masih ada sesudah hapus")
	}
}
