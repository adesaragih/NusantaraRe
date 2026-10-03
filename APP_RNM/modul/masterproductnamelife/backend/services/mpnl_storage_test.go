package services

// Penyimpanan berkas NYATA (keputusan work owner 03-10-2026 "untuk document masih belum berfungsi, ikuti dari XML
// nya aja") diuji terhadap layanan TIRUAN lokal (`httptest`, TLS). ⛔ Layanan sungguhan tidak pernah dipanggil uji;
// alamatnya hanya alamat server tiruan saat jalan.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/layanan"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

// jamStorage - "sekarang" uji (UTC).
var jamStorage = time.Date(2026, 10, 1, 8, 10, 0, 0, time.UTC)

// permintaanTercatat - satu panggilan titik layanan tiruan.
type permintaanTercatat struct {
	jalur, tipe, mentah string
	badan               map[string]any
}

// storageTiruan - `ServiceGoogle` tiruan: `/upload`, `/geturl`, `/delete`, dan objek bertanda tangan `/objek/<jalur>`.
type storageTiruan struct {
	srv   *httptest.Server
	mu    sync.Mutex
	minta []permintaanTercatat
	objek map[string][]byte
	// statusLayanan - jawaban titik layanan selain 200; urlKosong - `URLImage` kosong.
	statusLayanan int
	urlKosong     bool
	// tanpaAppfolder - jawaban geturl/upload tanpa `appfolder`.
	tanpaAppfolder bool
	ambilObjek     int
	// alihkan - setiap jawaban 307/302 ke `/tangkap`; tertangkap - `/tangkap` dipanggil (pengalihan diikuti).
	alihkan    bool
	tertangkap int
}

// gsObjek - `appfolder` jawaban layanan: awalan gs+App lalu jalur objek APA ADANYA (spasi tidak di-encode - bentuk
// APPFOLDER lampiran Pega di DEV).
func gsObjek(app, jalur string) string { return awalanGS(app) + jalur }

func baruStorageTiruan(t *testing.T) *storageTiruan {
	t.Helper()
	s := &storageTiruan{objek: map[string][]byte{}}
	s.srv = httptest.NewTLSServer(http.HandlerFunc(s.layani))
	t.Cleanup(s.srv.Close)
	return s
}

func (s *storageTiruan) layani(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r.URL.Path == "/tangkap" {
		s.tertangkap++
		_ = json.NewEncoder(w).Encode(map[string]string{"URLImage": s.srv.URL + "/tangkap"})
		return
	}
	if s.alihkan {
		w.Header().Set("Location", s.srv.URL+"/tangkap")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/objek/") {
		s.ambilObjek++
		isi, ada := s.objek[strings.TrimPrefix(r.URL.Path, "/objek/")]
		if !ada {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write(isi)
		return
	}
	mentah, _ := io.ReadAll(r.Body)
	var badan map[string]any
	_ = json.Unmarshal(mentah, &badan)
	s.minta = append(s.minta, permintaanTercatat{jalur: r.Method + " " + r.URL.Path, tipe: r.Header.Get("Content-Type"),
		mentah: string(mentah), badan: badan})
	if s.statusLayanan != 0 {
		w.WriteHeader(s.statusLayanan)
		return
	}
	teks := func(k string) string { v, _ := badan[k].(string); return v }
	jawab := map[string]string{}
	switch r.URL.Path {
	case "/upload":
		jalur := teks("Folder") + teks("Namafile")
		isi, _ := base64.StdEncoding.DecodeString(teks("Image"))
		s.objek[jalur] = isi
		jawab = map[string]string{"URLImage": s.srv.URL + "/objek/" + jalur, "exp": "2026-10-01T08:40:00.123Z",
			"appfolder": gsObjek(teks("App"), jalur), "DateTime": "10/01/2026 08:10:00"}
		if s.tanpaAppfolder {
			delete(jawab, "appfolder")
		}
	case "/geturl":
		jalur := teks("Folder") + teks("Namafile")
		jawab = map[string]string{"URLImage": s.srv.URL + "/objek/" + jalur + "?baru", "exp": "2026-10-01T08:40:00Z",
			"appfolder": gsObjek(teks("App"), jalur), "DateTime": "10/01/2026 08:10:00"}
		if s.tanpaAppfolder {
			delete(jawab, "appfolder")
		}
	case "/delete":
		delete(s.objek, teks("Namafile"))
	default:
		http.NotFound(w, r)
		return
	}
	if s.urlKosong {
		jawab["URLImage"] = ""
	}
	_ = json.NewEncoder(w).Encode(jawab)
}

// penyimpanan - penyimpanan nyata di atas tiruan: alamat `M_LINK_SERVICE` = server tiruan + KATEGORI_2.
func (s *storageTiruan) penyimpanan(t *testing.T) (PenyimpananBerkas, string) {
	t.Helper()
	dir := t.TempDir()
	alamat := func(_ context.Context, k layanan.KunciLayanan) (string, error) {
		if k.Kategori1 != "Google" {
			t.Errorf("KATEGORI_1 %q, bukan Google (GetLinkService)", k.Kategori1)
		}
		return s.srv.URL + "/" + k.Kategori2, nil
	}
	token := func(_ context.Context, app string) (string, error) {
		if app != "UJI-APP" {
			t.Errorf("token diminta untuk App %q", app)
		}
		return "UJI-TOKEN", nil
	}
	return PenyimpananGoogle(dir, alamat, token, s.srv.Client(), func() time.Time { return jamStorage }), dir
}

func (s *storageTiruan) daftarMinta() []permintaanTercatat {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]permintaanTercatat(nil), s.minta...)
}

const namaUji = "20261001-031000-0 - UJI Nota.PDF"

func objekUji() models.ObjekPenyimpanan {
	return models.ObjekPenyimpanan{ImageID: "ABC123", AppFolder: "Contract/Doc/2026/10/", FileName: namaUji,
		AppName: "UJI-APP", DurasiDetik: 1800}
}

func TestStorageGoogleKirimMeniruInsertGoogleStorage(t *testing.T) {
	s := baruStorageTiruan(t)
	p, _ := s.penyimpanan(t)
	ctx := context.Background()
	if err := p.SimpanAntrean(ctx, "ABC123", strings.NewReader("isi pdf")); err != nil {
		t.Fatal(err)
	}
	hasil, err := p.Kirim(ctx, objekUji(), "pdf", "application/pdf")
	if err != nil {
		t.Fatal(err)
	}
	m := s.daftarMinta()
	if len(m) != 1 || m[0].jalur != "POST /upload" || m[0].tipe != "application/json" {
		t.Fatalf("satu POST JSON ke titik upload (ServiceGoogle b148): %+v", m)
	}
	b := m[0].badan
	// `InsertGoogleStorage_Act` 8 b1339 (Set Data) + 9 b1572 (SET JSON: Durasi tanpa kutip, b1641).
	if b["App"] != "UJI-APP" || b["Kodestring"] != "UJI-TOKEN" || b["Folder"] != "Contract/Doc/2026/10/" ||
		b["Namafile"] != namaUji || b["Image"] != base64.StdEncoding.EncodeToString([]byte("isi pdf")) ||
		b["ext"] != "pdf" || b["MimeType"] != "application/pdf" || !strings.Contains(m[0].mentah, `"Durasi":1800`) {
		t.Errorf("badan UploadDoc: %s", m[0].mentah)
	}
	// b2366–b2515: URLImage, appfolder, exp (GMT, dd/MM/yyyy HH:mm:ss), Namafile, App → Insert_T_Storage_SQL.
	jalur := "Contract/Doc/2026/10/" + namaUji
	if hasil.ImageID != "ABC123" || hasil.URLPublic != s.srv.URL+"/objek/"+jalur || hasil.AppFolder != gsObjek("UJI-APP", jalur) ||
		hasil.Exp != "01/10/2026 08:40:00" || hasil.FileName != namaUji || hasil.AppName != "UJI-APP" {
		t.Errorf("objek T_STORAGE_IMAGE: %+v", hasil)
	}
	if string(s.objek[jalur]) != "isi pdf" {
		t.Errorf("isi terkirim: %q", s.objek[jalur])
	}
}

func TestStorageGoogleKirimGagalTerangTanpaAlamatDanToken(t *testing.T) {
	ctx := context.Background()
	s := baruStorageTiruan(t)
	p, _ := s.penyimpanan(t)
	if _, err := p.Kirim(ctx, objekUji(), "pdf", "application/pdf"); !errors.Is(err, ErrBerkasSumberHilang) {
		t.Errorf("antrean tidak ada: %v", err)
	}
	_ = p.SimpanAntrean(ctx, "ABC123", strings.NewReader("x"))
	s.urlKosong = true
	if _, err := p.Kirim(ctx, objekUji(), "pdf", "application/pdf"); !errors.Is(err, ErrStorageGagal) {
		t.Errorf("URLImage kosong = gagal (b2902): %v", err)
	}
	inang := s.srv.Listener.Addr().String()
	s.urlKosong, s.statusLayanan = false, http.StatusUnauthorized
	_, err := p.Kirim(ctx, objekUji(), "pdf", "application/pdf")
	if !errors.Is(err, ErrStorageGagal) || strings.Contains(err.Error(), inang) || strings.Contains(err.Error(), "UJI-TOKEN") {
		t.Errorf("401: galat tanpa alamat dan token: %v", err)
	}
	s.srv.Close()
	s.statusLayanan = 0
	if _, err := p.Kirim(ctx, objekUji(), "pdf", "application/pdf"); !errors.Is(err, ErrStorageGagal) ||
		strings.Contains(err.Error(), inang) {
		t.Errorf("tak terjangkau: galat tanpa alamat: %v", err)
	}
	gagalToken := PenyimpananGoogle(t.TempDir(), func(context.Context, layanan.KunciLayanan) (string, error) { return "x", nil },
		func(context.Context, string) (string, error) { return "", layanan.ErrGaramTokenKosong }, nil, nil)
	_ = gagalToken.SimpanAntrean(ctx, "ABC123", strings.NewReader("x"))
	if _, err := gagalToken.Kirim(ctx, objekUji(), "pdf", "application/pdf"); !errors.Is(err, ErrStorageBelumSiap) ||
		!errors.Is(err, layanan.ErrGaramTokenKosong) {
		t.Errorf("token tidak dapat diterbitkan = penyimpanan belum siap (503): %v", err)
	}
}

// `GetUrlGoogleStorage_Act`: URL tersimpan dipakai selama EXPDATE belum lewat; selain itu `geturl` dengan Folder =
// APPFOLDER tanpa Namafile tanpa awalan gs+App (b1260, b1325), Durasi 1800 (`DownloadAttProdName_Act` 6 b999).
func TestStorageGoogleBukaMemakaiURLTersimpanAtauMemintaYangBaru(t *testing.T) {
	ctx := context.Background()
	s := baruStorageTiruan(t)
	p, _ := s.penyimpanan(t)
	jalur := "Contract/Doc/2025/08/" + namaUji
	s.objek[jalur] = []byte("isi lama")
	o := models.ObjekPenyimpanan{ImageID: "ABC123", URLPublic: s.srv.URL + "/objek/" + jalur, AppFolder: gsObjek("UJI-APP", jalur),
		Exp: "01/10/2026 08:40:00", FileName: namaUji, AppName: "UJI-APP"}

	isi, baru, err := p.Buka(ctx, o)
	if err != nil {
		t.Fatal(err)
	}
	d, _ := io.ReadAll(isi)
	_ = isi.Close()
	if string(d) != "isi lama" || baru != nil || len(s.daftarMinta()) != 0 {
		t.Errorf("URL belum kedaluwarsa dipakai ulang tanpa geturl: %q %+v %d", d, baru, len(s.daftarMinta()))
	}

	o.Exp = "22/08/2025 07:00:00"
	isi, baru, err = p.Buka(ctx, o)
	if err != nil {
		t.Fatal(err)
	}
	d, _ = io.ReadAll(isi)
	_ = isi.Close()
	m := s.daftarMinta()
	if string(d) != "isi lama" || len(m) != 1 || m[0].jalur != "POST /geturl" {
		t.Fatalf("kedaluwarsa: geturl lalu isi: %q %+v", d, m)
	}
	b := m[0].badan
	if b["App"] != "UJI-APP" || b["Kodestring"] != "UJI-TOKEN" || b["Folder"] != "Contract/Doc/2025/08/" || b["Namafile"] != namaUji ||
		!strings.Contains(m[0].mentah, `"Durasi":1800`) || b["Image"] != nil {
		t.Errorf("badan geturl (b1346–b1431): %s", m[0].mentah)
	}
	// b2125–b2295 → Update_T_Storage_SQL: URLImage, appfolder, exp, DateTime.
	if baru == nil || baru.URLPublic != s.srv.URL+"/objek/"+jalur+"?baru" || baru.Exp != "01/10/2026 08:40:00" ||
		baru.TanggalUpload != "10/01/2026 08:10:00" || baru.AppFolder != gsObjek("UJI-APP", jalur) || baru.ImageID != "ABC123" {
		t.Errorf("objek diperbarui: %+v", baru)
	}

	// Jawaban tanpa appfolder tidak mengosongkan APPFOLDER - hapus membutuhkannya (penyimpangan sadar kecil).
	s.tanpaAppfolder = true
	_, baru, err = p.Buka(ctx, o)
	if err != nil || baru == nil || baru.AppFolder != o.AppFolder {
		t.Errorf("appfolder lama dipertahankan: %+v %v", baru, err)
	}
}

func TestStorageGoogleBukaGagalTerang(t *testing.T) {
	ctx := context.Background()
	s := baruStorageTiruan(t)
	p, _ := s.penyimpanan(t)
	hilang := models.ObjekPenyimpanan{ImageID: "ABC123", URLPublic: s.srv.URL + "/objek/tidak-ada", Exp: "01/10/2026 08:40:00",
		FileName: namaUji, AppName: "UJI-APP", AppFolder: gsObjek("UJI-APP", "tidak-ada")}
	if _, _, err := p.Buka(ctx, hilang); !errors.Is(err, ErrBerkasTidakDiStorage) {
		t.Errorf("404 URL bertanda tangan = berkas tidak ada: %v", err)
	}
	polos := hilang
	polos.URLPublic = strings.Replace(s.srv.URL, "https", "http", 1) + "/objek/tidak-ada"
	if _, _, err := p.Buka(ctx, polos); !errors.Is(err, ErrStorageGagal) || s.ambilObjek != 1 {
		t.Errorf("URL tanpa sandi tidak dibuka: %v (dibuka %d)", err, s.ambilObjek)
	}
}

// Objek yang dicatat stub lokal (URLPUBLIC kosong) dibaca dari folder stub - riwayat sebelum penyambungan tetap terbaca.
func TestStorageGoogleObjekStubDariFolderLokal(t *testing.T) {
	ctx := context.Background()
	s := baruStorageTiruan(t)
	p, dir := s.penyimpanan(t)
	lokal := PenyimpananLokal(dir)
	_ = lokal.SimpanAntrean(ctx, "STUB1", strings.NewReader("isi stub"))
	o := models.ObjekPenyimpanan{ImageID: "STUB1", AppFolder: "Contract/Doc/2026/09/", FileName: "x - UJI.pdf", AppName: "UJI-APP"}
	if _, err := lokal.Kirim(ctx, o, "pdf", "application/pdf"); err != nil {
		t.Fatal(err)
	}
	isi, baru, err := p.Buka(ctx, o)
	if err != nil {
		t.Fatal(err)
	}
	d, _ := io.ReadAll(isi)
	_ = isi.Close()
	if string(d) != "isi stub" || baru != nil || len(s.daftarMinta()) != 0 {
		t.Errorf("objek stub: %q %+v %d", d, baru, len(s.daftarMinta()))
	}
	if err := p.Hapus(ctx, "STUB1", &o); err != nil || len(s.daftarMinta()) != 0 {
		t.Errorf("hapus objek stub tanpa memanggil layanan: %v %d", err, len(s.daftarMinta()))
	}
	if _, _, err := lokal.Buka(ctx, o); err == nil {
		t.Error("berkas stub terhapus")
	}
}

// `DeleteGoogleStorage_Act`: Namafile = APPFOLDER tanpa awalan gs+App (b1091), tanpa Folder; antrean lokal ikut dibuang.
func TestStorageGoogleHapusMeniruDeleteGoogleStorage(t *testing.T) {
	ctx := context.Background()
	s := baruStorageTiruan(t)
	p, _ := s.penyimpanan(t)
	jalur := "Contract/Doc/2025/08/" + namaUji
	s.objek[jalur] = []byte("x")
	_ = p.SimpanAntrean(ctx, "ABC123", strings.NewReader("antre"))
	o := models.ObjekPenyimpanan{ImageID: "ABC123", URLPublic: s.srv.URL + "/objek/" + jalur, AppFolder: gsObjek("UJI-APP", jalur),
		FileName: namaUji, AppName: "UJI-APP"}
	if err := p.Hapus(ctx, "ABC123", &o); err != nil {
		t.Fatal(err)
	}
	m := s.daftarMinta()
	if len(m) != 1 || m[0].jalur != "POST /delete" || m[0].badan["Namafile"] != jalur || m[0].badan["App"] != "UJI-APP" ||
		m[0].badan["Kodestring"] != "UJI-TOKEN" || m[0].badan["Folder"] != nil {
		t.Errorf("badan delete: %+v", m)
	}
	if _, ada := s.objek[jalur]; ada {
		t.Error("objek terhapus di penyimpanan")
	}
	if _, err := p.Kirim(ctx, objekUji(), "pdf", "application/pdf"); !errors.Is(err, ErrBerkasSumberHilang) {
		t.Errorf("antrean lokal ikut dibuang: %v", err)
	}
	if err := p.Hapus(ctx, "BELUM1", nil); err != nil || len(s.daftarMinta()) != 1 {
		t.Errorf("belum terkirim: hanya antrean, tanpa layanan: %v", err)
	}
	s.statusLayanan = http.StatusInternalServerError
	if err := p.Hapus(ctx, "ABC123", &o); !errors.Is(err, ErrStorageGagal) {
		t.Errorf("hapus gagal = rekam tetap (DeleteAttacProdName_act 2 b411 StepStatusFail): %v", err)
	}
}

func TestExpPega(t *testing.T) {
	// `@substring(@pxReplaceAllViaRegex(exp, "[-:]", ""), 0, 19) + " GMT"` → `dd/MM/yyyy HH:mm:ss` GMT (b2366, b2431).
	for masuk, keluar := range map[string]string{
		"2026-10-01T08:40:00.123Z": "01/10/2026 08:40:00",
		"2026-10-01T08:40:00Z":     "01/10/2026 08:40:00",
		"":                         "",
		"besok":                    "",
	} {
		if g := expPega(masuk); g != keluar {
			t.Errorf("%q → %q, bukan %q", masuk, g, keluar)
		}
	}
}

// Keputusan work owner 03-10-2026 "selalu nyata, ikut XML": Layanan Oracle SELALU memakai penyimpanan Google
// (tanpa saklar `PELAKSANA_STORAGE`); garam hanya bahan token baru.
func TestPerakitanSelaluPenyimpananNyata(t *testing.T) {
	s := New(nil).DenganUnggahanDir(t.TempDir())
	for _, svc := range []*Service{{Dasar: s}, (&Service{Dasar: s}).DenganGaramToken("UJI-GARAM")} {
		p, ok := svc.penyimpanan().(penyimpananGoogle)
		if !ok || p.lokal == nil || p.klien.Timeout != BatasWaktuStorage {
			t.Fatalf("penyimpanan Google berantrean lokal, batas waktu ServiceGoogle: %T", svc.penyimpanan())
		}
		if _, err := p.token(context.Background(), "UJI-APP"); err == nil {
			t.Error("tanpa Oracle: token gagal terang")
		}
		if _, err := p.alamat(context.Background(), layanan.KunciUnggahBerkas); err == nil {
			t.Error("tanpa Oracle: alamat gagal terang")
		}
	}
}

// `View Office Online` (`DownloadAttProdName_Act` 6 b953 + 7 b1080): URL bertanda tangan saja - tanpa mengunduh isinya;
// objek stub tidak punya URL.
func TestStorageGoogleTautan(t *testing.T) {
	ctx := context.Background()
	s := baruStorageTiruan(t)
	p, _ := s.penyimpanan(t)
	jalur := "Contract/Doc/2025/08/" + namaUji
	o := models.ObjekPenyimpanan{ImageID: "ABC123", URLPublic: s.srv.URL + "/objek/" + jalur, AppFolder: gsObjek("UJI-APP", jalur),
		Exp: "01/10/2026 08:40:00", FileName: namaUji, AppName: "UJI-APP"}
	u, baru, err := p.Tautan(ctx, o)
	if err != nil || u != o.URLPublic || baru != nil || len(s.daftarMinta()) != 0 || s.ambilObjek != 0 {
		t.Errorf("URL tersimpan: %q %+v %v (geturl %d, unduh %d)", u, baru, err, len(s.daftarMinta()), s.ambilObjek)
	}
	o.Exp = "22/08/2025 07:00:00"
	u, baru, err = p.Tautan(ctx, o)
	if err != nil || baru == nil || u != baru.URLPublic || len(s.daftarMinta()) != 1 || s.ambilObjek != 0 {
		t.Errorf("kedaluwarsa: geturl, tanpa unduh: %q %+v %v", u, baru, err)
	}
	if _, _, err := p.Tautan(ctx, models.ObjekPenyimpanan{ImageID: "STUB1", FileName: "x"}); !errors.Is(err, ErrOfficeStub) {
		t.Errorf("objek stub: %v", err)
	}
	if _, _, err := PenyimpananLokal(t.TempDir()).Tautan(ctx, o); !errors.Is(err, ErrOfficeStub) {
		t.Errorf("stub lokal: %v", err)
	}
}

// Jawaban upload tanpa `appfolder`: APPFOLDER dirakit dari awalan gs+App, Folder, dan Namafile - hapus dan geturl
// tetap dapat menamai objeknya.
func TestStorageGoogleKirimTanpaAppfolderDirakit(t *testing.T) {
	ctx := context.Background()
	s := baruStorageTiruan(t)
	p, _ := s.penyimpanan(t)
	_ = p.SimpanAntrean(ctx, "ABC123", strings.NewReader("x"))
	s.tanpaAppfolder = true
	hasil, err := p.Kirim(ctx, objekUji(), "pdf", "application/pdf")
	if err != nil || hasil.AppFolder != gsObjek("UJI-APP", "Contract/Doc/2026/10/"+namaUji) {
		t.Errorf("APPFOLDER dirakit: %q %v", hasil.AppFolder, err)
	}
	if err := p.Hapus(ctx, "ABC123", &models.ObjekPenyimpanan{URLPublic: "x", AppName: "UJI-APP"}); !errors.Is(err, ErrStorageGagal) {
		t.Errorf("APPFOLDER kosong: hapus ditolak, bukan Namafile kosong ke layanan: %v", err)
	}
}

// ⛔ Pengalihan tidak diikuti - baik jawaban titik layanan (badan memuat token) maupun URL bertanda tangan.
func TestStorageGoogleTanpaPengalihan(t *testing.T) {
	ctx := context.Background()
	s := baruStorageTiruan(t)
	p, _ := s.penyimpanan(t)
	_ = p.SimpanAntrean(ctx, "ABC123", strings.NewReader("x"))
	s.alihkan = true
	if _, err := p.Kirim(ctx, objekUji(), "pdf", "application/pdf"); !errors.Is(err, ErrStorageGagal) {
		t.Errorf("307 titik layanan = gagal: %v", err)
	}
	o := models.ObjekPenyimpanan{ImageID: "ABC123", URLPublic: s.srv.URL + "/objek/x", Exp: "01/10/2026 08:40:00",
		FileName: namaUji, AppName: "UJI-APP", AppFolder: gsObjek("UJI-APP", "x")}
	if _, _, err := p.Buka(ctx, o); !errors.Is(err, ErrStorageGagal) {
		t.Errorf("302 URL bertanda tangan = gagal: %v", err)
	}
	if s.tertangkap != 0 {
		t.Errorf("pengalihan diikuti %d kali", s.tertangkap)
	}
}

// URL tersimpan yang tinggal kurang dari satu menit diminta ulang - unduhan tidak berangkat dengan URL yang hampir mati.
func TestStorageGoogleURLHampirKedaluwarsaDimintaUlang(t *testing.T) {
	ctx := context.Background()
	s := baruStorageTiruan(t)
	p, _ := s.penyimpanan(t)
	jalur := "Contract/Doc/2025/08/" + namaUji
	s.objek[jalur] = []byte("x")
	o := models.ObjekPenyimpanan{ImageID: "ABC123", URLPublic: s.srv.URL + "/objek/" + jalur, AppFolder: gsObjek("UJI-APP", jalur),
		Exp: jamStorage.Add(30 * time.Second).Format(formatExp), FileName: namaUji, AppName: "UJI-APP"}
	isi, baru, err := p.Buka(ctx, o)
	if err != nil {
		t.Fatal(err)
	}
	_ = isi.Close()
	if baru == nil || len(s.daftarMinta()) != 1 {
		t.Errorf("sisa 30 detik: geturl dipanggil: %+v %d", baru, len(s.daftarMinta()))
	}
	if tanggalUploadPega("2026-10-01 08:10:00") != "" || tanggalUploadPega(" 10/01/2026 08:10:00 ") != "10/01/2026 08:10:00" {
		t.Error("DateTime hanya bentuk MM/DD/YYYY HH24:MI:SS (Update_T_Storage_SQL)")
	}
}

// Sebab mentah (teks Oracle) hanya di log; galat bernama `inti/backend/layanan` menyebut kuncinya.
func TestGalatStorageKalimatLayarTanpaSebabMentah(t *testing.T) {
	err := belumSiap("the storage token", errors.New("ORA-01017: UJI-RAHASIA"))
	if !errors.Is(err, ErrStorageBelumSiap) || strings.Contains(Pesan(err), "ORA-") || !strings.Contains(err.Error(), "ORA-01017") ||
		!strings.Contains(Pesan(err), "the storage token is not available") {
		t.Errorf("layar %q / log %q", Pesan(err), err.Error())
	}
	err = belumSiap("the storage token", layanan.ErrGaramTokenKosong)
	if !errors.Is(err, layanan.ErrGaramTokenKosong) || !strings.Contains(Pesan(err), "STORAGE_TOKEN_SALT") {
		t.Errorf("galat bernama tampil: %q", Pesan(err))
	}
}

// penyimpanTokenPalsu - `GCP_IMAGE` di memori.
type penyimpanTokenPalsu struct {
	berlaku        string
	ditanyaSesudah time.Time
	disimpan       []string
	sampai         time.Time
}

func (p *penyimpanTokenPalsu) TokenBerlaku(_ context.Context, _ *db.Tx, _ string, saat time.Time) (string, error) {
	p.ditanyaSesudah = saat
	return p.berlaku, nil
}

func (p *penyimpanTokenPalsu) SimpanToken(_ context.Context, _ *db.Tx, _, token, pengguna string, sampai time.Time) error {
	p.disimpan, p.sampai = append(p.disimpan, token+"/"+pengguna), sampai
	return nil
}

// `GET_TOKEN_STORAGE` ditiru dengan margin 15 detik (preseden TCO `MarginTokenTCO`): token yang tinggal kurang dari itu
// tidak dipakai ulang - unggahan base64 besar tidak berangkat dengan token yang mati di tengah jalan.
func TestTokenStorageBermargin(t *testing.T) {
	ctx := context.Background()
	pt := &penyimpanTokenPalsu{berlaku: "UJI-LAMA"}
	tok, err := tokenStorage(ctx, nil, pt, "UJI-GARAM", "UJI-APP", jamStorage)
	if err != nil || tok != "UJI-LAMA" || !pt.ditanyaSesudah.Equal(jamStorage.Add(MarginToken)) || len(pt.disimpan) != 0 {
		t.Errorf("token berlaku dipakai ulang bila sisa > margin: %q %v %v %v", tok, err, pt.ditanyaSesudah, pt.disimpan)
	}
	pt = &penyimpanTokenPalsu{}
	tok, err = tokenStorage(ctx, nil, pt, "UJI-GARAM", "UJI-APP", jamStorage)
	harap, _ := layanan.RakitToken("UJI-GARAM", jamStorage)
	if err != nil || tok != harap || len(pt.disimpan) != 1 || pt.disimpan[0] != harap+"/"+layanan.PenggunaTokenBawaan ||
		!pt.sampai.Equal(jamStorage.Add(layanan.UmurToken)) {
		t.Errorf("token baru dirakit dan disimpan semenit: %q %v %v %v", tok, err, pt.disimpan, pt.sampai)
	}
	// Garam hanya bahan token BARU: token berlaku dipakai ulang tanpa garam (DEV 03-10-2026: token berlaku +-58 hari).
	if tok, err := tokenStorage(ctx, nil, &penyimpanTokenPalsu{berlaku: "UJI-LAMA"}, "", "UJI-APP", jamStorage); err != nil || tok != "UJI-LAMA" {
		t.Errorf("tanpa garam, token berlaku: %q %v", tok, err)
	}
	if _, err := tokenStorage(ctx, nil, &penyimpanTokenPalsu{}, " ", "UJI-APP", jamStorage); !errors.Is(err, layanan.ErrGaramTokenKosong) {
		t.Errorf("tanpa garam, tanpa token berlaku: %v", err)
	}
	if _, err := tokenStorage(ctx, nil, &penyimpanTokenPalsu{}, "UJI-GARAM", "", jamStorage); !errors.Is(err, layanan.ErrAppNameKosong) {
		t.Errorf("App kosong: %v", err)
	}
}

// Stub tidak menghapus objek yang ada di penyimpanan nyata (URLPUBLIC terisi) - rekam tetap, objek tidak yatim.
func TestStubMenolakHapusObjekPenyimpananNyata(t *testing.T) {
	p := PenyimpananLokal(t.TempDir())
	err := p.Hapus(context.Background(), "ABC123", &models.ObjekPenyimpanan{ImageID: "ABC123", URLPublic: "UJI-URL"})
	if !errors.Is(err, ErrStorageBelumSiap) || !strings.Contains(Pesan(err), "storage service") {
		t.Errorf("stub: %v", err)
	}
}
