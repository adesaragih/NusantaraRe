package penyimpanan

// Uji keempat activity terhadap server HTTPS tiruan lokal (`httptest`) dan catatan di memori - tidak ada layanan
// penyimpanan, Gemini, maupun Oracle sungguhan yang dipanggil.

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
)

const tokenUji = "UJI-TOKEN-RAHASIA"

type catatanTiruan struct {
	mu       sync.Mutex
	app      string
	objek    map[string]Objek
	perbarui []Objek
}

func (c *catatanTiruan) NamaAplikasi(context.Context) (string, error) { return c.app, nil }
func (c *catatanTiruan) Ambil(_ context.Context, id string) (Objek, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	o, ada := c.objek[id]
	return o, ada, nil
}
func (c *catatanTiruan) Catat(_ context.Context, _ *db.Tx, o Objek) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.objek[o.ImageID] = o
	return nil
}
func (c *catatanTiruan) Perbarui(_ context.Context, o Objek) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.perbarui = append(c.perbarui, o)
	c.objek[o.ImageID] = o
	return nil
}
func (c *catatanTiruan) Hapus(_ context.Context, _ *db.Tx, id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.objek, id)
	return nil
}

// layananTiruan merekam setiap POST per titik (`/upload`, `/geturl`, `/delete`, `/getAI`) dan menjawab `jawab`.
type layananTiruan struct {
	mu     sync.Mutex
	badan  map[string][]map[string]any
	jawab  map[string]any
	status int
	isi    string
}

func ujiPenyimpanan(t *testing.T) (*Penyimpanan, *catatanTiruan, *layananTiruan, *httptest.Server) {
	t.Helper()
	l := &layananTiruan{badan: map[string][]map[string]any{}, jawab: map[string]any{}, status: http.StatusOK}
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet { // URL bertanda tangan
			if r.URL.Path == "/hilang" {
				http.NotFound(w, r)
				return
			}
			if r.URL.Path == "/alih" {
				http.Redirect(w, r, "/berkas", http.StatusFound)
				return
			}
			_, _ = io.WriteString(w, l.isi)
			return
		}
		var b map[string]any
		_ = json.NewDecoder(r.Body).Decode(&b)
		l.mu.Lock()
		l.badan[r.URL.Path] = append(l.badan[r.URL.Path], b)
		status, j := l.status, l.jawab[r.URL.Path]
		l.mu.Unlock()
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(j)
	}))
	t.Cleanup(srv.Close)
	c := &catatanTiruan{app: "uji-bucket", objek: map[string]Objek{}}
	token := func(_ context.Context, app, _ string) (string, error) {
		if app == "" {
			return "", layanan.ErrAppNameKosong
		}
		return tokenUji, nil
	}
	alamat := func(_ context.Context, k layanan.KunciLayanan) (string, error) {
		if k.Kategori1 != layanan.KategoriGoogle {
			return "", layanan.ErrAlamatLayananTidakAda
		}
		return srv.URL + "/" + k.Kategori2, nil
	}
	jam := func() time.Time { return time.Date(2026, 10, 8, 14, 5, 9, 7_000_000, zonaJakarta) }
	return Baru(c, token, alamat, srv.Client(), jam), c, l, srv
}

func TestUnggahMengirimSepertiXML(t *testing.T) {
	p, _, l, _ := ujiPenyimpanan(t)
	l.jawab["/upload"] = map[string]any{"URLImage": "URL-TIRUAN-1", "exp": "2026-10-08T08:05:09.000Z",
		"appfolder": "gs-jalur-tiruan"}
	o, err := p.Unggah(context.Background(), MasukUnggah{Folder: "Contract", NamaFile: "Laporan.PDF", Isi: []byte("isi berkas uji"),
		Durasi: 3600, Pengguna: "UJI-MAKER"})
	if err != nil {
		t.Fatal(err)
	}
	b := l.badan["/upload"]
	if len(b) != 1 {
		t.Fatalf("upload %d kali, mau 1", len(b))
	}
	mau := map[string]any{"App": "uji-bucket", "Kodestring": tokenUji, "Durasi": float64(3600), "Folder": "Contract/Doc/2026/10/",
		"Namafile": "20261008-020509-7 - Laporan.PDF", "Image": base64.StdEncoding.EncodeToString([]byte("isi berkas uji")),
		"ext": "pdf", "MimeType": "application/pdf"}
	for k, v := range mau {
		if b[0][k] != v {
			t.Errorf("%s = %#v, mau %#v", k, b[0][k], v)
		}
	}
	if len(o.ImageID) != 32 || strings.ToUpper(o.ImageID) != o.ImageID {
		t.Errorf("ImageID %q bukan 32 heksa huruf besar", o.ImageID)
	}
	if o.Exp != "08/10/2026 08:05:09" || o.FileName != "20261008-020509-7 - Laporan.PDF" || o.AppName != "uji-bucket" ||
		o.AppFolder != "gs-jalur-tiruan" || o.URLPublic != "URL-TIRUAN-1" {
		t.Errorf("objek %+v", o)
	}
}

func TestUnggahTanpaAppfolderMerakitJalur(t *testing.T) {
	p, _, l, _ := ujiPenyimpanan(t)
	l.jawab["/upload"] = map[string]any{"URLImage": "u"}
	o, err := p.Unggah(context.Background(), MasukUnggah{Folder: "Contract", NamaFile: "a.docx", Isi: []byte("isi berkas uji")})
	if err != nil {
		t.Fatal(err)
	}
	if o.AppFolder != awalanGS("uji-bucket")+"Contract/Doc/2026/10/20261008-020509-7 - a.docx" {
		t.Errorf("AppFolder %q", o.AppFolder)
	}
}

func TestUnggahMenolakBerkasYangPegaLewati(t *testing.T) {
	p, _, l, _ := ujiPenyimpanan(t)
	for nama, m := range map[string]MasukUnggah{
		"tanpa ekstensi": {NamaFile: "laporan", Isi: []byte("isi berkas uji")},
		"jenis asing":    {NamaFile: "program.exe", Isi: []byte("isi berkas uji")},
		"isi kosong":     {NamaFile: "a.pdf"},
		"isi < 10 b64":   {NamaFile: "a.pdf", Isi: []byte("abc")},
	} {
		if _, err := p.Unggah(context.Background(), m); !errors.Is(err, ErrBerkasDitolak) {
			t.Errorf("%s: %v, mau ErrBerkasDitolak", nama, err)
		}
	}
	if len(l.badan) != 0 {
		t.Errorf("layanan dipanggil untuk berkas yang ditolak: %v", l.badan)
	}
}

func TestUnggahURLKosongGagalTerang(t *testing.T) {
	p, _, _, _ := ujiPenyimpanan(t)
	if _, err := p.Unggah(context.Background(), MasukUnggah{NamaFile: "a.pdf", Isi: []byte("isi berkas uji")}); !errors.Is(err, ErrStorageGagal) {
		t.Errorf("%v, mau ErrStorageGagal", err)
	}
}

func TestGalatTidakMemuatAlamatAtauToken(t *testing.T) {
	p, _, l, srv := ujiPenyimpanan(t)
	l.status = http.StatusInternalServerError
	_, err := p.Unggah(context.Background(), MasukUnggah{NamaFile: "a.pdf", Isi: []byte("isi berkas uji")})
	if !errors.Is(err, ErrStorageGagal) {
		t.Fatalf("%v, mau ErrStorageGagal", err)
	}
	if strings.Contains(err.Error(), srv.URL) || strings.Contains(err.Error(), tokenUji) {
		t.Errorf("galat membocorkan alamat atau token: %v", err)
	}
	p2 := Baru(&catatanTiruan{app: "uji-bucket", objek: map[string]Objek{}},
		func(context.Context, string, string) (string, error) { return tokenUji, nil },
		func(context.Context, layanan.KunciLayanan) (string, error) {
			return "", layanan.ErrAlamatLayananTidakAda
		}, srv.Client(), nil)
	if _, err := p2.Unggah(context.Background(), MasukUnggah{NamaFile: "a.pdf", Isi: []byte("isi berkas uji")}); !errors.Is(err, ErrStorageBelumSiap) {
		t.Errorf("alamat tidak ada: %v, mau ErrStorageBelumSiap", err)
	}
	p3 := Baru(&catatanTiruan{objek: map[string]Objek{}}, nil, nil, srv.Client(), nil)
	if _, err := p3.Unggah(context.Background(), MasukUnggah{NamaFile: "a.pdf", Isi: []byte("isi berkas uji")}); !errors.Is(err, ErrStorageBelumSiap) {
		t.Errorf("APPNAME kosong: %v, mau ErrStorageBelumSiap", err)
	}
}

func TestTautanMemakaiURLYangMasihBerlaku(t *testing.T) {
	p, c, l, _ := ujiPenyimpanan(t)
	c.objek["IMG1"] = Objek{ImageID: "IMG1", URLPublic: "URL-LAMA", Exp: "08/10/2026 15:00:00", AppName: "uji-bucket"}
	u, err := p.Tautan(context.Background(), "IMG1", 3600, "UJI")
	if err != nil || u != "URL-LAMA" {
		t.Fatalf("%q %v", u, err)
	}
	if len(l.badan) != 0 {
		t.Errorf("geturl dipanggil padahal EXPDATE belum lewat")
	}
}

func TestTautanMemintaURLBaruSaatKedaluwarsa(t *testing.T) {
	p, c, l, _ := ujiPenyimpanan(t)
	jalur := awalanGS("uji-bucket") + "Contract/Doc/2026/10/f.pdf"
	c.objek["IMG1"] = Objek{ImageID: "IMG1", URLPublic: "URL-LAMA", Exp: "08/10/2026 14:00:00", AppName: "uji-bucket",
		AppFolder: jalur, FileName: "f.pdf"}
	l.jawab["/geturl"] = map[string]any{"URLImage": "URL-BARU", "exp": "2026-10-08T09:05:09Z", "DateTime": "10/08/2026 14:05:09"}
	u, err := p.Tautan(context.Background(), "IMG1", 3600, "UJI")
	if err != nil || u != "URL-BARU" {
		t.Fatalf("%q %v", u, err)
	}
	b := l.badan["/geturl"][0]
	if b["Folder"] != "Contract/Doc/2026/10/" || b["Namafile"] != "f.pdf" || b["Durasi"] != float64(3600) || b["App"] != "uji-bucket" {
		t.Errorf("badan geturl %v", b)
	}
	if len(c.perbarui) != 1 || c.perbarui[0].URLPublic != "URL-BARU" || c.perbarui[0].Exp != "08/10/2026 09:05:09" ||
		c.perbarui[0].AppFolder != jalur || c.perbarui[0].TanggalUpload != "10/08/2026 14:05:09" {
		t.Errorf("Update_T_Storage_SQL %+v", c.perbarui)
	}
	// 6.6 dilewati: URLImage kosong -> URL tersimpan, tanpa UPDATE.
	l.jawab["/geturl"] = map[string]any{}
	c.objek["IMG2"] = Objek{ImageID: "IMG2", URLPublic: "URL-LAMA-2", AppName: "uji-bucket", AppFolder: jalur}
	if u, err := p.Tautan(context.Background(), "IMG2", 3600, "UJI"); err != nil || u != "URL-LAMA-2" {
		t.Errorf("%q %v", u, err)
	}
	if _, err := p.Tautan(context.Background(), "TIDAK-ADA", 3600, "UJI"); !errors.Is(err, ErrObjekTidakAda) {
		t.Errorf("%v, mau ErrObjekTidakAda", err)
	}
}

func TestBukaHanyaHttpsTanpaPengalihan(t *testing.T) {
	p, c, l, srv := ujiPenyimpanan(t)
	l.isi = "isi dari penyimpanan"
	berlaku := "08/10/2026 23:00:00"
	c.objek["OK"] = Objek{ImageID: "OK", URLPublic: srv.URL + "/berkas", Exp: berlaku}
	c.objek["HILANG"] = Objek{ImageID: "HILANG", URLPublic: srv.URL + "/hilang", Exp: berlaku}
	c.objek["ALIH"] = Objek{ImageID: "ALIH", URLPublic: srv.URL + "/alih", Exp: berlaku}
	polos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = io.WriteString(w, "x") }))
	defer polos.Close()
	c.objek["POLOS"] = Objek{ImageID: "POLOS", URLPublic: polos.URL, Exp: berlaku}
	isi, err := p.Buka(context.Background(), "OK", 3600, "UJI")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(isi)
	_ = isi.Close()
	if string(b) != "isi dari penyimpanan" {
		t.Errorf("isi %q", b)
	}
	if _, err := p.Buka(context.Background(), "HILANG", 3600, "UJI"); !errors.Is(err, ErrBerkasTidakDiStorage) {
		t.Errorf("404: %v", err)
	}
	if _, err := p.Buka(context.Background(), "ALIH", 3600, "UJI"); !errors.Is(err, ErrStorageGagal) {
		t.Errorf("pengalihan diikuti: %v", err)
	}
	if _, err := p.Buka(context.Background(), "POLOS", 3600, "UJI"); !errors.Is(err, ErrStorageGagal) {
		t.Errorf("alamat bukan https dibuka: %v", err)
	}
}

func TestHapusObjekMengirimJalurTanpaAwalanGS(t *testing.T) {
	p, c, l, _ := ujiPenyimpanan(t)
	c.objek["IMG1"] = Objek{ImageID: "IMG1", AppName: "uji-bucket", AppFolder: awalanGS("uji-bucket") + "Contract/Doc/2026/10/f.pdf"}
	if err := p.HapusObjek(context.Background(), "IMG1", "UJI"); err != nil {
		t.Fatal(err)
	}
	b := l.badan["/delete"]
	if len(b) != 1 || b[0]["Namafile"] != "Contract/Doc/2026/10/f.pdf" || b[0]["App"] != "uji-bucket" || b[0]["Folder"] != nil {
		t.Errorf("badan delete %v", b)
	}
	if err := p.HapusObjek(context.Background(), "TIDAK-ADA", "UJI"); err != nil || len(l.badan["/delete"]) != 1 {
		t.Errorf("objek tak tercatat: %v, delete %d kali", err, len(l.badan["/delete"]))
	}
	l.status = http.StatusBadGateway
	if err := p.HapusObjek(context.Background(), "IMG1", "UJI"); !errors.Is(err, ErrStorageGagal) {
		t.Errorf("hapus gagal: %v, mau ErrStorageGagal", err)
	}
}

func TestTanyaAI(t *testing.T) {
	p, c, l, _ := ujiPenyimpanan(t)
	l.jawab["/getAI"] = map[string]any{"Hasil": "12345"}
	c.objek["IMG1"] = Objek{ImageID: "IMG1", AppName: "bucket-gambar", AppFolder: "gs-jalur-gambar"}
	h, err := p.TanyaAI(context.Background(), MasukAI{ImageID: "IMG1", Pesan: "kode pos?", Pengguna: "UJI"})
	if err != nil || h != "12345" {
		t.Fatalf("%q %v", h, err)
	}
	h, err = p.TanyaAI(context.Background(), MasukAI{Pesan: "tanpa gambar"})
	if err != nil || h != "12345" {
		t.Fatalf("%q %v", h, err)
	}
	b := l.badan["/getAI"]
	if b[0]["App"] != "bucket-gambar" || b[0]["Tanya"] != "kode pos?" || b[0]["Temp"] != 0.7 || b[0]["Kodestring"] != tokenUji {
		t.Errorf("dengan gambar %v", b[0])
	}
	if n, _ := b[0]["Namafile"].([]any); len(n) != 1 || n[0] != "gs-jalur-gambar" {
		t.Errorf("Namafile %v", b[0]["Namafile"])
	}
	if _, ada := b[1]["Namafile"]; ada || b[1]["App"] != "uji-bucket" {
		t.Errorf("tanpa gambar %v", b[1])
	}
}

func TestPesanBelumSiapMenyebutKunciBukanNilai(t *testing.T) {
	for sebab, mau := range map[error]string{
		layanan.ErrGaramTokenKosong:       "STORAGE_TOKEN_SALT",
		layanan.ErrAppNameKosong:          "T_FOLDER_IMAGE",
		layanan.ErrAlamatLayananTidakAda:  "M_LINK_SERVICE",
		layanan.ErrEndpointTidakDitemukan: "M_LINK_SERVICE",
		errors.New("lain"):                "server log",
	} {
		if p := PesanBelumSiap(belumSiap("x", sebab)); !strings.Contains(p, mau) {
			t.Errorf("%v: %q", sebab, p)
		}
	}
}
