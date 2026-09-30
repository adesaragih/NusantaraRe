package services_test

// Uji transport penyimpanan lampiran (OQ-TCO-08) - SERVER TIRUAN LOKAL.
//
// ⛔ Tidak ada layanan sungguhan yang dipanggil: setiap alamat berasal dari
// `httptest.NewTLSServer` saat uji berjalan; nol alamat literal di berkas ini.
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

	"nusantarare/inti/backend/layanan"
	"nusantarare/modul/treatycontractout/repository"
	"nusantarare/modul/treatycontractout/services"
)

const (
	appUjiStorage   = "UJI-APP"
	tokenUjiStorage = "UJI-TOKEN-PALSU"
	garamUjiStorage = "UJI-GARAM-PALSU"
	kunciUjiStorage = "ABCDEF0123456789"
)

func appUjiStorageTCO(context.Context) (string, error) { return appUjiStorage, nil }

type permintaanStorageDiterima struct {
	jalur, metode, jenis, jengkal string
	badan                         map[string]any
}

// layananStorageUji meniru layanan penyimpanan seperti Pega memakainya:
// `/upload` menyimpan di `Folder + Namafile`, `/geturl` memberi URL bertanda
// tangan `/objek/{folder}{nama}`, `/delete` membuang `Namafile` = jalur penuh.
type layananStorageUji struct {
	mu          sync.Mutex
	objek       map[string][]byte
	diterima    []permintaanStorageDiterima
	paksa       map[string]int
	paksaSekali map[string]int
	urlGanti    func(jalurObjek string) string
	urlKosong   bool
	srv         *httptest.Server
}

func layananStorageTiruan(t *testing.T) *layananStorageUji {
	t.Helper()
	l := &layananStorageUji{objek: map[string][]byte{}, paksa: map[string]int{}, paksaSekali: map[string]int{}}
	l.srv = httptest.NewTLSServer(http.HandlerFunc(l.layani))
	t.Cleanup(l.srv.Close)
	return l
}

func (l *layananStorageUji) alamat(jalur string) string { return l.srv.URL + jalur }

func (l *layananStorageUji) setel(fn func()) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fn()
}

func (l *layananStorageUji) salinDiterima() []permintaanStorageDiterima {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]permintaanStorageDiterima(nil), l.diterima...)
}

func inangDari(alamat string) string {
	u, _ := url.Parse(alamat)
	return u.Host
}

func (l *layananStorageUji) layani(w http.ResponseWriter, r *http.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var badan map[string]any
	_ = json.NewDecoder(r.Body).Decode(&badan)
	l.diterima = append(l.diterima, permintaanStorageDiterima{jalur: r.URL.Path, metode: r.Method,
		jenis: r.Header.Get("Content-Type"), jengkal: r.Header.Get("Range"), badan: badan})
	if kode, ada := l.paksaSekali[r.URL.Path]; ada {
		delete(l.paksaSekali, r.URL.Path)
		w.WriteHeader(kode)
		return
	}
	if kode, ada := l.paksa[r.URL.Path]; ada {
		w.WriteHeader(kode)
		return
	}
	if r.URL.Path == "/alih" {
		http.Redirect(w, r, "/objek/"+services.FolderStorageTCO+kunciUjiStorage, http.StatusFound)
		return
	}
	if jalur, ada := strings.CutPrefix(r.URL.Path, "/objek/"); ada {
		isi, ada := l.objek[jalur]
		if !ada {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("Range") != "" && len(isi) > 0 {
			w.WriteHeader(http.StatusPartialContent)
			_, _ = w.Write(isi[:1])
			return
		}
		_, _ = w.Write(isi)
		return
	}
	folder, _ := badan["Folder"].(string)
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
		l.objek[folder+nama] = data
		u := l.alamat("/objek/" + folder + nama)
		if l.urlKosong {
			u = ""
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"URLImage": u, "exp": expUjiStorage})
	case "/geturl":
		// Seperti URL bertanda tangan: diberikan walau objeknya tidak ada.
		u := l.alamat("/objek/" + folder + nama)
		if l.urlGanti != nil {
			u = l.urlGanti(folder + nama)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"URLImage": u, "exp": expUjiStorage,
			"appfolder": "UJI-FOLDER", "DateTime": "09/29/2026 09:00:00"})
	case "/delete":
		delete(l.objek, nama)
		_ = json.NewEncoder(w).Encode(map[string]string{})
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (l *layananStorageUji) isiObjek(jalur string) ([]byte, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ada := l.objek[jalur]
	return b, ada
}

func (l *layananStorageUji) cacahObjek() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.objek)
}

// Bentuk `UploadDoc` (InsertGoogleStorage_Act b1339-b1641), `ext` dari nama
// berkas asli (b587), dan pengulangan ke objek yang SAMA (AC 58).
func TestPengirimStorageUnggahBentukUploadDoc(t *testing.T) {
	l := layananStorageTiruan(t)
	p := services.NewPengirimBerkasHTTPTCO(l.srv.Client(), appUjiStorageTCO)
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		if _, err := p.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage,
			strings.NewReader("ISI-PDF"), "application/pdf", ".PDF"); err != nil {
			t.Fatal(err)
		}
	}
	jalur := services.FolderStorageTCO + kunciUjiStorage
	if isi, _ := l.isiObjek(jalur); l.cacahObjek() != 1 || string(isi) != "ISI-PDF" {
		t.Fatalf("objek: %d, isi %q", l.cacahObjek(), isi)
	}
	if !strings.HasSuffix(services.FolderStorageTCO, "/") {
		t.Error("folder tanpa garis miring penutup (Pega b1407-b1408)")
	}
	d := l.salinDiterima()[0]
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
	// Tanpa ekstensi: medan `ext` tidak dikirim, TIDAK ditebak dari MIME.
	if _, err := p.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage,
		strings.NewReader("x"), "image/jpeg", ""); err != nil {
		t.Fatal(err)
	}
	if b := l.salinDiterima()[2].badan; b["ext"] != nil {
		t.Errorf("ext ditebak: %v", b["ext"])
	}
}

func TestPengirimStorageAmbilPeriksaBuang(t *testing.T) {
	l := layananStorageTiruan(t)
	p := services.NewPengirimBerkasHTTPTCO(l.srv.Client(), appUjiStorageTCO)
	ctx := context.Background()
	if _, err := p.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage, strings.NewReader("ISI"), "text/plain", "txt"); err != nil {
		t.Fatal(err)
	}
	rc, _, err := p.Ambil(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage)
	if err != nil {
		t.Fatal(err)
	}
	isi, _ := io.ReadAll(rc)
	_ = rc.Close()
	if string(isi) != "ISI" {
		t.Errorf("isi %q", isi)
	}
	if ada, _, err := p.Periksa(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage); err != nil || !ada {
		t.Errorf("periksa sebelum buang: %v %v", ada, err)
	}
	if d := l.salinDiterima(); d[len(d)-1].jengkal != "bytes=0-0" {
		t.Errorf("periksa mengunduh objek utuh (Range %q)", d[len(d)-1].jengkal)
	}
	if err := p.Buang(ctx, l.alamat("/delete"), tokenUjiStorage, kunciUjiStorage); err != nil {
		t.Fatal(err)
	}
	d := l.salinDiterima()
	b := d[len(d)-1].badan
	// DeleteGoogleStorage_Act b1091: Namafile = jalur objek PENUH, tanpa Folder.
	if b["App"] != appUjiStorage || b["Kodestring"] != tokenUjiStorage ||
		b["Namafile"] != services.FolderStorageTCO+kunciUjiStorage || b["Image"] != nil || b["Folder"] != nil {
		t.Errorf("badan delete: %v", b)
	}
	if l.cacahObjek() != 0 {
		t.Fatal("objek tidak terhapus - delete menunjuk jalur lain")
	}
	if ada, _, err := p.Periksa(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage); err != nil || ada {
		t.Errorf("periksa sesudah buang: %v %v", ada, err)
	}
	if _, _, err := p.Ambil(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage); !errors.Is(err, services.ErrBerkasTidakAdaDiPenyimpanan) {
		t.Errorf("ambil sesudah buang: %v", err)
	}
}

// ⛔ 404 TITIK LAYANAN bukan "berkas tidak ada" - ia juga jawaban jalur
// M_LINK_SERVICE yang salah; "tidak ada" hanya dari URL bertanda tangan.
func TestPengirimStorage404TitikLayananBukanBerkasHilang(t *testing.T) {
	l := layananStorageTiruan(t)
	p := services.NewPengirimBerkasHTTPTCO(l.srv.Client(), appUjiStorageTCO)
	ctx := context.Background()
	if err := p.Buang(ctx, l.alamat("/salah-jalur"), tokenUjiStorage, kunciUjiStorage); errors.Is(err, services.ErrBerkasTidakAdaDiPenyimpanan) ||
		!errors.Is(err, services.ErrStorageGagalTCO) {
		t.Errorf("delete 404: %v", err)
	}
	if ada, _, err := p.Periksa(ctx, l.alamat("/salah-jalur"), tokenUjiStorage, kunciUjiStorage); ada || !errors.Is(err, services.ErrStorageGagalTCO) {
		t.Errorf("geturl 404: %v %v", ada, err)
	}
}

// ⛔ URL bertanda tangan hanya https, dan pengalihannya tidak diikuti.
func TestPengirimStorageURLBertandaTerjaga(t *testing.T) {
	l := layananStorageTiruan(t)
	p := services.NewPengirimBerkasHTTPTCO(l.srv.Client(), appUjiStorageTCO)
	ctx := context.Background()
	if _, err := p.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage, strings.NewReader("ISI"), "", ""); err != nil {
		t.Fatal(err)
	}
	polos := strings.Replace(l.srv.URL, "https", "http", 1)
	l.setel(func() { l.urlGanti = func(j string) string { return polos + "/objek/" + j } })
	if _, _, err := p.Ambil(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage); !errors.Is(err, services.ErrStorageJawabanRusakTCO) {
		t.Errorf("URLImage polos: %v", err)
	}
	l.setel(func() { l.urlGanti = func(string) string { return l.alamat("/alih") } })
	sebelum := len(l.salinDiterima())
	if _, _, err := p.Ambil(ctx, l.alamat("/geturl"), tokenUjiStorage, kunciUjiStorage); !errors.Is(err, services.ErrStorageGagalTCO) {
		t.Errorf("pengalihan: %v", err)
	}
	for _, d := range l.salinDiterima()[sebelum:] {
		if strings.HasPrefix(d.jalur, "/objek/") {
			t.Errorf("pengalihan diikuti ke %s", d.jalur)
		}
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
		http.StatusNotFound:            services.ErrStorageGagalTCO,
		http.StatusUnauthorized:        services.ErrStorageTokenDitolakTCO,
		http.StatusForbidden:           services.ErrStorageTokenDitolakTCO,
		http.StatusTooManyRequests:     services.ErrStorageGagalTCO,
		http.StatusInternalServerError: services.ErrStorageGagalTCO,
	} {
		l.setel(func() { l.paksa["/upload"] = kode })
		_, err := p.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage, strings.NewReader("x"), "", "")
		if !errors.Is(err, mau) {
			t.Errorf("status %d: %v, mau %v", kode, err, mau)
			continue
		}
		bersihDariRahasia(t, err, inangDari(l.srv.URL))
	}
	l.setel(func() {
		delete(l.paksa, "/upload")
		l.urlKosong = true
	})
	if _, err := p.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage, strings.NewReader("x"), "", ""); !errors.Is(err, services.ErrStorageJawabanRusakTCO) {
		t.Errorf("URLImage kosong (b2902): %v", err)
	}

	mati := httptest.NewTLSServer(http.NotFoundHandler())
	alamatMati := mati.URL + "/upload"
	mati.Close()
	_, err := p.Kirim(ctx, alamatMati, tokenUjiStorage, kunciUjiStorage, strings.NewReader("x"), "", "")
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
	_, err = pl.Kirim(ctx, lambat.URL+"/upload", tokenUjiStorage, kunciUjiStorage, strings.NewReader("x"), "", "")
	if !errors.Is(err, services.ErrStorageTakTerjangkauTCO) || !strings.Contains(err.Error(), "timeout") {
		t.Errorf("batas waktu: %v", err)
	}
	bersihDariRahasia(t, err, inangDari(lambat.URL))

	tanpaApp := services.NewPengirimBerkasHTTPTCO(l.srv.Client(), func(context.Context) (string, error) {
		return "", repository.ErrAppStorageKosongTCO
	})
	sebelum := len(l.salinDiterima())
	if _, err := tanpaApp.Kirim(ctx, l.alamat("/upload"), tokenUjiStorage, kunciUjiStorage, strings.NewReader("x"), "", ""); !errors.Is(err, repository.ErrAppStorageKosongTCO) {
		t.Errorf("tanpa App: %v", err)
	}
	if len(l.salinDiterima()) != sebelum {
		t.Error("permintaan terkirim padahal App belum ada")
	}
}

// rangkaianNyataUji: resolver (M_LINK_SERVICE tiruan) -> cache token -> sumber
// token GET_TOKEN_STORAGE tiruan -> transport HTTP -> server tiruan lokal.
func rangkaianNyataUji(l *layananStorageUji, hapus string) (*services.PenyimpananJarakJauhTCO, *penyimpanTokenUji, time.Time) {
	saat := time.Date(2026, 9, 29, 9, 0, 0, 0, time.UTC)
	jam := func() time.Time { return saat }
	r := &resolverLampiranUji{alamat: map[layanan.KunciLayanan]string{
		layanan.KunciUnggahBerkas: l.alamat("/upload"),
		layanan.KunciURLBerkas:    l.alamat("/geturl"),
		layanan.KunciHapusBerkas:  l.alamat(hapus),
	}}
	gudangToken := &penyimpanTokenUji{app: appUjiStorage}
	sumber := services.NewSumberTokenStorageTCO(transaksiUji, gudangToken, garamUjiStorage, jam)
	return services.NewPenyimpananJarakJauhTCO(r, services.NewCacheTokenTCO(sumber, jam, services.MarginTokenTCO),
		services.NewPengirimBerkasHTTPTCO(l.srv.Client(), appUjiStorageTCO)), gudangToken, saat
}

func TestPenyimpananNyataUjungKeUjung(t *testing.T) {
	l := layananStorageTiruan(t)
	p, gudangToken, saat := rangkaianNyataUji(l, "/delete")
	ctx := context.Background()
	if _, err := p.Simpan(ctx, kunciUjiStorage, strings.NewReader("ISI"), "application/pdf", "pdf"); err != nil {
		t.Fatal(err)
	}
	if ada, err := p.Ada(ctx, kunciUjiStorage); err != nil || !ada {
		t.Fatalf("ada: %v %v", ada, err)
	}
	if err := p.Hapus(ctx, kunciUjiStorage); err != nil {
		t.Fatal(err)
	}
	mau, _ := layanan.RakitToken(garamUjiStorage, saat)
	for _, d := range l.salinDiterima() {
		if !strings.HasPrefix(d.jalur, "/objek/") && d.badan["Kodestring"] != mau {
			t.Errorf("%s: Kodestring bukan token rakitan", d.jalur)
		}
	}
	if len(gudangToken.disimpan) != 1 {
		t.Errorf("token diterbitkan %d kali, mau 1 (dipakai ulang oleh cache)", len(gudangToken.disimpan))
	}
	if l.cacahObjek() != 0 {
		t.Error("objek masih ada sesudah hapus")
	}
	// Berkas yang sudah tidak ada: dibaca dari objeknya, delete tidak dikirim.
	sebelum := len(l.salinDiterima())
	if err := p.Hapus(ctx, kunciUjiStorage); !errors.Is(err, services.ErrBerkasTidakAdaDiPenyimpanan) {
		t.Errorf("hapus berkas yang sudah tidak ada: %v", err)
	}
	for _, d := range l.salinDiterima()[sebelum:] {
		if d.jalur == "/delete" {
			t.Error("delete dikirim untuk objek yang tidak ada")
		}
	}
}

// Jalur delete yang salah di M_LINK_SERVICE: galat terlihat, objek tetap ada.
func TestPenyimpananNyataHapusJalurSalahTidakDiam(t *testing.T) {
	l := layananStorageTiruan(t)
	p, _, _ := rangkaianNyataUji(l, "/salah-jalur")
	ctx := context.Background()
	if _, err := p.Simpan(ctx, kunciUjiStorage, strings.NewReader("ISI"), "", ""); err != nil {
		t.Fatal(err)
	}
	err := p.Hapus(ctx, kunciUjiStorage)
	if err == nil || errors.Is(err, services.ErrBerkasTidakAdaDiPenyimpanan) {
		t.Errorf("hapus ke jalur salah: %v", err)
	}
	if l.cacahObjek() != 1 {
		t.Error("objek hilang padahal delete gagal")
	}
}

// Token yang DITOLAK layanan (401/403) dilupakan: percobaan berikutnya menerbitkan yang baru.
func TestPenyimpananNyataTokenDitolakDilupakan(t *testing.T) {
	l := layananStorageTiruan(t)
	p, gudangToken, _ := rangkaianNyataUji(l, "/delete")
	ctx := context.Background()
	l.setel(func() { l.paksaSekali["/upload"] = http.StatusUnauthorized })
	if _, err := p.Simpan(ctx, kunciUjiStorage, strings.NewReader("ISI"), "", ""); !errors.Is(err, services.ErrStorageTokenDitolakTCO) {
		t.Fatalf("401: %v", err)
	}
	if _, err := p.Simpan(ctx, kunciUjiStorage, strings.NewReader("ISI"), "", ""); err != nil {
		t.Fatal(err)
	}
	if len(gudangToken.disimpan) != 2 {
		t.Errorf("token diterbitkan %d kali, mau 2 (yang ditolak dilupakan)", len(gudangToken.disimpan))
	}
}

// expUjiStorage - `exp` jawaban layanan berbentuk ISO (Pega membuang `-`/`:`).
const expUjiStorage = "2026-09-29T10:00:00.000Z"

// OQ-TCO-26 (lanjutan 4): jawaban geturl (`exp` ISO, `appfolder`, `DateTime`)
// sampai ke pencatat dalam bentuk To_date `Update_T_Storage_SQL`; jawaban
// unggah pun mengubah `exp` seperti `InsertGoogleStorage_Act` b2366/b2431.
func TestPenyimpananNyataMenyegarkanObjekSesudahGetURL(t *testing.T) {
	l := layananStorageTiruan(t)
	p, _, _ := rangkaianNyataUji(l, "/delete")
	pc := &pencatatObjekUji{}
	p = p.DenganPencatatObjek(pc, func(string) {})
	ctx := context.Background()
	o, err := p.Simpan(ctx, kunciUjiStorage, strings.NewReader("ISI"), "application/pdf", "pdf")
	if err != nil {
		t.Fatal(err)
	}
	if o.Exp != "29/09/2026 10:00:00" {
		t.Errorf("exp unggah %q, mau bentuk To_date", o.Exp)
	}
	rc, err := p.Buka(ctx, kunciUjiStorage)
	if err != nil {
		t.Fatal(err)
	}
	_ = rc.Close()
	if len(pc.objek) != 1 {
		t.Fatalf("disegarkan %d kali, mau 1", len(pc.objek))
	}
	g := pc.objek[0]
	if g.ImageID != kunciUjiStorage || g.URLPublic != l.alamat("/objek/"+services.FolderStorageTCO+kunciUjiStorage) ||
		g.AppFolder != "UJI-FOLDER" || g.Exp != "29/09/2026 10:00:00" || g.TanggalUpload != "09/29/2026 09:00:00" {
		t.Errorf("objek disegarkan: %+v", g)
	}
}
