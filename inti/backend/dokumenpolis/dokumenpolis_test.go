package dokumenpolis_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/dokumenpolis"
	"nusantarare/inti/backend/penyimpanan"
)

// catatanTiruan - DOCUMENT_POLIS dan CATEGORY_ATTACH_REAS di memori (IDPEGA per dokumen).
type catatanTiruan struct {
	mu       sync.Mutex
	kategori []string
	dok      map[string]dokumenpolis.Dokumen
	idPega   map[string]string
}

func (c *catatanTiruan) cocok(baca []string, id string) bool {
	for _, b := range baca {
		if c.idPega[id] == b {
			return true
		}
	}
	return false
}

func (c *catatanTiruan) Kategori(_ context.Context, baca []string) ([]dokumenpolis.Kategori, error) {
	out := []dokumenpolis.Kategori{}
	for _, k := range c.kategori {
		n := 0
		for id, d := range c.dok {
			if d.Kategori == k && c.cocok(baca, id) {
				n++
			}
		}
		out = append(out, dokumenpolis.Kategori{Nama: k, Cacah: n})
	}
	return out, nil
}
func (c *catatanTiruan) AdaKategori(_ context.Context, nama string) (bool, error) {
	for _, k := range c.kategori {
		if k == nama {
			return true, nil
		}
	}
	return false, nil
}
func (c *catatanTiruan) Daftar(_ context.Context, baca []string, kategori string) ([]dokumenpolis.Dokumen, error) {
	out := []dokumenpolis.Dokumen{}
	for id, d := range c.dok {
		if d.Kategori == kategori && c.cocok(baca, id) {
			out = append(out, d)
		}
	}
	return out, nil
}
func (c *catatanTiruan) Ambil(_ context.Context, baca []string, id string) (dokumenpolis.Dokumen, bool, error) {
	d, ada := c.dok[id]
	if !ada || !c.cocok(baca, id) {
		return dokumenpolis.Dokumen{}, false, nil
	}
	return d, true, nil
}
func (c *catatanTiruan) Sisip(_ context.Context, _ *db.Tx, d dokumenpolis.Dokumen, idPega string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ada := c.dok[d.ID]; ada {
		return dokumenpolis.ErrIDTerpakai
	}
	c.dok[d.ID], c.idPega[d.ID] = d, idPega
	return nil
}
func (c *catatanTiruan) Hapus(_ context.Context, _ *db.Tx, baca []string, id string) error {
	if c.cocok(baca, id) {
		delete(c.dok, id)
	}
	return nil
}

// berkasTiruan - penyimpanan di memori.
type berkasTiruan struct {
	isi   map[string][]byte
	objek map[string]bool
	masuk []penyimpanan.MasukUnggah
	hapus []string
	gagal error
	n     int
}

func (b *berkasTiruan) Unggah(_ context.Context, m penyimpanan.MasukUnggah) (penyimpanan.Objek, error) {
	if b.gagal != nil {
		return penyimpanan.Objek{}, b.gagal
	}
	b.n++
	id := "UJI-IMG-" + string(rune('A'-1+b.n))
	b.masuk = append(b.masuk, m)
	b.isi[id] = m.Isi
	return penyimpanan.Objek{ImageID: id}, nil
}
func (b *berkasTiruan) Catat(_ context.Context, _ *db.Tx, o penyimpanan.Objek) error {
	b.objek[o.ImageID] = true
	return nil
}
func (b *berkasTiruan) Buka(_ context.Context, id string, _ int, _ string) (io.ReadCloser, error) {
	if !b.objek[id] {
		return nil, penyimpanan.ErrObjekTidakAda
	}
	return io.NopCloser(bytes.NewReader(b.isi[id])), nil
}
func (b *berkasTiruan) Tautan(_ context.Context, id string, _ int, _ string) (string, error) {
	return "URL-" + id, nil
}
func (b *berkasTiruan) HapusObjek(_ context.Context, id, _ string) error {
	if b.gagal != nil {
		return b.gagal
	}
	b.hapus = append(b.hapus, id)
	return nil
}
func (b *berkasTiruan) HapusCatatan(_ context.Context, _ *db.Tx, id string) error {
	delete(b.objek, id)
	return nil
}

func transaksi(_ context.Context, fn func(tx *db.Tx) error) error { return fn(nil) }

var jamUji = time.Date(2026, 10, 8, 14, 5, 7, 123*int(time.Millisecond), time.FixedZone("WIB", 7*3600))

func layanan() (*dokumenpolis.Layanan, *catatanTiruan, *berkasTiruan) {
	c := &catatanTiruan{kategori: []string{"R/I SLIP", "UW ANALYSIS"}, dok: map[string]dokumenpolis.Dokumen{}, idPega: map[string]string{}}
	b := &berkasTiruan{isi: map[string][]byte{}, objek: map[string]bool{}}
	return dokumenpolis.Baru(c, b, transaksi, func() time.Time { return jamUji }), c, b
}

func TestKasusDariMembacaKunciBaruDanPzInsKeyPega(t *testing.T) {
	k := dokumenpolis.KasusDari("NB-22445", "NB-22445", true)
	if k.Tulis != "NB-22445" || len(k.Baca) != 2 || k.Baca[1] != "ASM-FW-GISFW-WORK NB-22445" || !k.BolehUbah {
		t.Errorf("%+v", k)
	}
	// Salinan Copy Old: ID kasus = pzInsKey utuh - tidak digandakan.
	k = dokumenpolis.KasusDari("ASM-FW-GISFW-WORK NB-7", "NB-7", false)
	if len(k.Baca) != 1 || k.Tulis != "ASM-FW-GISFW-WORK NB-7" {
		t.Errorf("%+v", k)
	}
}

func TestIDMengikutiXML(t *testing.T) {
	if got := dokumenpolis.IDDokumen(jamUji); got != "20261008020507123" {
		t.Errorf("%s (hh 12 jam seperti InsertDocument_Act)", got)
	}
}

func TestUnggahSepertiInsertDocumentAct(t *testing.T) {
	l, c, b := layanan()
	ctx := context.Background()
	tutup := dokumenpolis.KasusDari("NB-1", "NB-1", false)
	if _, err := l.Unggah(ctx, tutup, "UJI", "UW ANALYSIS", "a.pdf", []byte("ISI")); !errors.Is(err, dokumenpolis.ErrDilarang) {
		t.Errorf("kasus Resolve: %v", err)
	}
	k := dokumenpolis.KasusDari("NB-1", "NB-1", true)
	if _, err := l.Unggah(ctx, k, "UJI", "TIDAK ADA", "a.pdf", []byte("ISI")); !errors.Is(err, dokumenpolis.ErrMasukanTidakSah) {
		t.Errorf("kategori asing: %v", err)
	}
	d, err := l.Unggah(ctx, k, "UJI-ADMIN", "R/I SLIP", `C:\fakepath\Slip Final.PDF`, []byte("ISI PDF"))
	if err != nil {
		t.Fatal(err)
	}
	if d.ID != "20261008020507123" || d.NamaFile != "Slip Final.PDF" || d.Ekstensi != "pdf" || d.Kategori != "R/I SLIP" ||
		d.Pengunggah != "UJI-ADMIN" || d.StorageID != "UJI-IMG-A" || c.idPega[d.ID] != "NB-1" {
		t.Errorf("dokumen %+v idpega %q", d, c.idPega[d.ID])
	}
	if m := b.masuk[0]; m.Folder != "Policy" || m.Durasi != 1800 || m.Ext != "pdf" || m.Pengguna != "UJI-ADMIN" {
		t.Errorf("InsertGoogleStorage_Act %+v", m)
	}
	d2, err := l.Unggah(ctx, k, "UJI-ADMIN", "R/I SLIP", "b.xlsx", []byte("ISI XLSX"))
	if err != nil || d2.ID != "20261008020507124" {
		t.Errorf("ID bentrok -> milidetik berikutnya: %+v %v", d2, err)
	}
	b.gagal = penyimpanan.ErrStorageGagal
	if _, err := l.Unggah(ctx, k, "UJI", "R/I SLIP", "c.pdf", []byte("ISI")); !errors.Is(err, penyimpanan.ErrStorageGagal) || len(c.dok) != 2 {
		t.Errorf("penyimpanan gagal: %v, %d dokumen", err, len(c.dok))
	}
}

func TestDaftarKategoriMembacaLampiranPegaLama(t *testing.T) {
	l, c, b := layanan()
	ctx := context.Background()
	c.dok["LAMA"] = dokumenpolis.Dokumen{ID: "LAMA", NamaFile: "uw.pdf", Ekstensi: "pdf", Kategori: "UW ANALYSIS", StorageID: "IMG-LAMA", AdaObjek: true}
	c.idPega["LAMA"] = "ASM-FW-GISFW-WORK NB-1"
	b.objek["IMG-LAMA"], b.isi["IMG-LAMA"] = true, []byte("ISI LAMA")
	c.dok["LAIN"] = dokumenpolis.Dokumen{ID: "LAIN", Kategori: "UW ANALYSIS"}
	c.idPega["LAIN"] = "ASM-FW-GISFW-WORK NB-2"
	k := dokumenpolis.KasusDari("NB-1", "NB-1", false)
	kat, _ := l.Kategori(ctx, k)
	if len(kat) != 2 || kat[1].Nama != "UW ANALYSIS" || kat[1].Cacah != 1 || kat[0].Cacah != 0 {
		t.Errorf("kategori %+v", kat)
	}
	f, err := l.Unduh(ctx, k, "UJI", "LAMA")
	if err != nil {
		t.Fatal(err)
	}
	isi, _ := io.ReadAll(f.Isi)
	if string(isi) != "ISI LAMA" || f.Mime != "application/pdf" {
		t.Errorf("unduh %q %+v", isi, f)
	}
	if _, err := l.Unduh(ctx, k, "UJI", "LAIN"); !errors.Is(err, dokumenpolis.ErrTidakAda) {
		t.Errorf("dokumen kasus lain terbaca: %v", err)
	}
	if _, err := l.TautanOffice(ctx, k, "UJI", "LAMA"); !errors.Is(err, dokumenpolis.ErrMasukanTidakSah) {
		t.Errorf("office untuk pdf: %v", err)
	}
	if err := l.Hapus(ctx, k, "UJI", "LAMA"); !errors.Is(err, dokumenpolis.ErrDilarang) {
		t.Errorf("hapus kasus Resolve: %v", err)
	}
}

func TestHapusSepertiDeleteDocumentPolis(t *testing.T) {
	l, c, b := layanan()
	ctx := context.Background()
	k := dokumenpolis.KasusDari("NB-1", "NB-1", true)
	d, _ := l.Unggah(ctx, k, "UJI", "UW ANALYSIS", "a.pdf", []byte("ISI"))
	b.gagal = penyimpanan.ErrStorageGagal
	if err := l.Hapus(ctx, k, "UJI", d.ID); !errors.Is(err, penyimpanan.ErrStorageGagal) || len(c.dok) != 1 {
		t.Errorf("hapus jarak jauh gagal: %v, %d dokumen", err, len(c.dok))
	}
	b.gagal = nil
	if err := l.Hapus(ctx, k, "UJI", d.ID); err != nil || len(c.dok) != 0 || len(b.objek) != 0 || b.hapus[0] != d.StorageID {
		t.Errorf("hapus: %v dok %v objek %v", err, c.dok, b.objek)
	}
	// Baris tanpa T_STORAGE_ID: langsung dihapus (`DeleteGoogleStorage_Act` berprasyarat T_STORAGE_ID terisi).
	c.dok["KOSONG"], c.idPega["KOSONG"] = dokumenpolis.Dokumen{ID: "KOSONG", Kategori: "UW ANALYSIS"}, "NB-1"
	if err := l.Hapus(ctx, k, "UJI", "KOSONG"); err != nil || len(b.hapus) != 1 {
		t.Errorf("tanpa objek: %v hapus %v", err, b.hapus)
	}
}

func TestRutePasang(t *testing.T) {
	l, _, b := layanan()
	errKasusTidakAda := errors.New("UJI kasus tidak ada")
	sumber := func(_ context.Context, _ inti.Pelaku, id string) (dokumenpolis.Kasus, error) {
		switch id {
		case "NB-1":
			return dokumenpolis.KasusDari(id, id, true), nil
		case "NB-9":
			return dokumenpolis.KasusDari(id, id, false), nil
		}
		return dokumenpolis.Kasus{}, errKasusTidakAda
	}
	galatModul := func(w http.ResponseWriter, err error) {
		if errors.Is(err, errKasusTidakAda) {
			http.Error(w, "kasus", http.StatusNotFound)
			return
		}
		http.Error(w, "lain", http.StatusInternalServerError)
	}
	mux := http.NewServeMux()
	dokumenpolis.Pasang(mux, "/api/uji/kasus/{id}/lampiran", func() *dokumenpolis.Layanan { return l }, sumber, true, galatModul)
	kirim := func(metode, jalur string, badan io.Reader, tipe string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(metode, jalur, badan)
		if tipe != "" {
			r.Header.Set("Content-Type", tipe)
		}
		r = r.WithContext(inti.DenganPelakuSesi(r.Context(), inti.Pelaku{AkunID: "UJI-ADMIN"}))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w
	}
	unggahBerkas := func(id, kategori, nama string) *httptest.ResponseRecorder {
		var buf bytes.Buffer
		m := multipart.NewWriter(&buf)
		f, _ := m.CreateFormFile("berkas", nama)
		_, _ = f.Write([]byte("ISI PDF"))
		_ = m.Close()
		return kirim("POST", "/api/uji/kasus/"+id+"/lampiran/dokumen?kategori="+kategori, &buf, m.FormDataContentType())
	}

	if w := unggahBerkas("NB-1", "R%2FI%20SLIP", "slip.pdf"); w.Code != 200 || strings.Contains(w.Body.String(), "UJI-IMG") {
		t.Fatalf("unggah: %d %s", w.Code, w.Body.String())
	}
	if w := unggahBerkas("NB-9", "R%2FI%20SLIP", "slip.pdf"); w.Code != 403 {
		t.Errorf("kasus Resolve: %d %s", w.Code, w.Body.String())
	}
	if w := kirim("GET", "/api/uji/kasus/NB-404/lampiran", nil, ""); w.Code != 404 {
		t.Errorf("kasus tidak ada -> galat modul: %d", w.Code)
	}
	w := kirim("GET", "/api/uji/kasus/NB-1/lampiran", nil, "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `{"nama":"R/I SLIP","cacah":1}`) || !strings.Contains(w.Body.String(), `"bolehUbah":true`) {
		t.Errorf("kategori: %d %s", w.Code, w.Body.String())
	}
	w = kirim("GET", "/api/uji/kasus/NB-1/lampiran/dokumen?kategori=R%2FI%20SLIP", nil, "")
	var daftar struct{ Daftar []dokumenpolis.Dokumen }
	_ = json.Unmarshal(w.Body.Bytes(), &daftar)
	if w.Code != 200 || len(daftar.Daftar) != 1 {
		t.Fatalf("daftar: %d %s", w.Code, w.Body.String())
	}
	id := daftar.Daftar[0].ID
	if w := kirim("GET", "/api/uji/kasus/NB-1/lampiran/dokumen/"+id+"/isi", nil, ""); w.Code != 200 || w.Body.String() != "ISI PDF" ||
		!strings.Contains(w.Header().Get("Content-Disposition"), "slip.pdf") {
		t.Errorf("unduh: %d %s", w.Code, w.Body.String())
	}
	if w := kirim("GET", "/api/uji/kasus/NB-1/lampiran/dokumen/"+id+"/office", nil, ""); w.Code != 422 {
		t.Errorf("office pdf: %d", w.Code)
	}
	b.gagal = penyimpanan.ErrStorageBelumSiap
	if w := kirim("POST", "/api/uji/kasus/NB-1/lampiran/dokumen/"+id+"/hapus", nil, ""); w.Code != 503 {
		t.Errorf("penyimpanan belum siap: %d %s", w.Code, w.Body.String())
	}
	b.gagal = nil
	if w := kirim("POST", "/api/uji/kasus/NB-1/lampiran/dokumen/"+id+"/hapus", nil, ""); w.Code != 200 {
		t.Errorf("hapus: %d %s", w.Code, w.Body.String())
	}
	if w := kirim("GET", "/api/uji/kasus/NB-1/lampiran/dokumen/"+id+"/isi", nil, ""); w.Code != 404 {
		t.Errorf("sesudah hapus: %d", w.Code)
	}
}
