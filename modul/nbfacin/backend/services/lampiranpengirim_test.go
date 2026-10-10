package services

import (
	"context"
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

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/layanan"
	"nusantarare/modul/nbfacin/backend/models"
)

// layananStorageTiruan - server httptest TLS: titik upload / geturl / delete + objek bertanda tangan.
type layananStorageTiruan struct {
	mu       sync.Mutex
	srv      *httptest.Server
	minta    map[string][]map[string]any
	status   map[string]int
	jawaban  map[string]map[string]string
	objek    map[string]string // jalur objek -> isi
	unduhan  []string
	tolakURL bool
}

func layananStorageUji(t *testing.T) *layananStorageTiruan {
	t.Helper()
	l := &layananStorageTiruan{minta: map[string][]map[string]any{}, status: map[string]int{},
		jawaban: map[string]map[string]string{}, objek: map[string]string{}}
	l.srv = httptest.NewTLSServer(http.HandlerFunc(l.layani))
	t.Cleanup(l.srv.Close)
	return l
}

func (l *layananStorageTiruan) layani(w http.ResponseWriter, r *http.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if strings.HasPrefix(r.URL.Path, "/objek/") {
		l.unduhan = append(l.unduhan, r.URL.Path+"|"+r.Header.Get("Range"))
		isi, ada := l.objek[strings.TrimPrefix(r.URL.Path, "/objek/")]
		if !ada {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("Range") != "" {
			w.WriteHeader(http.StatusPartialContent)
		}
		_, _ = io.WriteString(w, isi)
		return
	}
	if r.URL.Path == "/alih" {
		http.Redirect(w, r, "/objek/x", http.StatusFound)
		return
	}
	titik := strings.TrimPrefix(r.URL.Path, "/")
	var badan map[string]any
	_ = json.NewDecoder(r.Body).Decode(&badan)
	l.minta[titik] = append(l.minta[titik], badan)
	if s := l.status[titik]; s != 0 {
		w.WriteHeader(s)
		return
	}
	j := l.jawaban[titik]
	if j == nil {
		j = map[string]string{}
	}
	_ = json.NewEncoder(w).Encode(j)
}

type resolverTiruan struct{ dasar string }

func (r resolverTiruan) Resolve(_ context.Context, k layanan.KunciLayanan) (string, error) {
	return r.dasar + "/" + k.Kategori2, nil
}

type sumberTokenTiruan struct{ n *int }

func (s sumberTokenTiruan) TokenBaru(context.Context) (string, time.Time, error) {
	*s.n++
	return "TOKEN-UJI", time.Now().Add(time.Hour), nil
}

type pencatatTiruan struct{ o []models.ObjekStorage }

func (p *pencatatTiruan) PerbaruiObjek(_ context.Context, o models.ObjekStorage) error {
	p.o = append(p.o, o)
	return nil
}

func storageJarakJauhUji(t *testing.T, l *layananStorageTiruan, kini time.Time) (KlienStorageLampiran, *int, *pencatatTiruan) {
	t.Helper()
	n := 0
	pc := &pencatatTiruan{}
	app := func(context.Context) (string, error) { return "UJIAPP", nil }
	return NewStorageLampiranJarakJauh(resolverTiruan{l.srv.URL}, NewCacheTokenStorage(sumberTokenTiruan{&n}, nil, MarginTokenStorage),
		app, l.srv.Client(), pc, func() time.Time { return kini }), &n, pc
}

var lokUji = models.LokasiObjek{ImageID: "ABCDEF0123456789", Folder: "Policy/Doc/2026/10/", Namafile: "20261008-010405-7 - a.pdf"}

// TestStorageJarakJauhSimpan - badan InsertGoogleStorage_Act (App, Kodestring, Durasi ANGKA 1800, Folder, Namafile, Image
// base64, ext, MimeType); jawaban -> objek (exp diubah bentuk To_date); URLImage kosong = rusak; status -> jenis galat.
func TestStorageJarakJauhSimpan(t *testing.T) {
	ctx := context.Background()
	l := layananStorageUji(t)
	appFolder := models.AppFolderUji("UJIAPP", "Policy/Doc/2026/10/20261008-010405-7 - a.pdf")
	l.jawaban["upload"] = map[string]string{"URLImage": l.srv.URL + "/objek/a", "exp": "2026-10-08T13:34:05.000Z",
		"appfolder": appFolder}
	st, nToken, _ := storageJarakJauhUji(t, l, time.Now())
	o, err := st.Simpan(ctx, lokUji, strings.NewReader("ISI"), "pdf", "application/pdf")
	mau := models.ObjekStorage{ImageID: lokUji.ImageID, URLPublic: l.srv.URL + "/objek/a",
		AppFolder: appFolder, Exp: "08/10/2026 13:34:05",
		Namafile: lokUji.Namafile, App: "UJIAPP"}
	if err != nil || o != mau {
		t.Fatalf("%v %+v", err, o)
	}
	b := l.minta["upload"][0]
	if b["App"] != "UJIAPP" || b["Kodestring"] != "TOKEN-UJI" || b["Durasi"] != float64(1800) || b["Folder"] != lokUji.Folder ||
		b["Namafile"] != lokUji.Namafile || b["Image"] != "SVNJ" || b["ext"] != "pdf" || b["MimeType"] != "application/pdf" {
		t.Errorf("badan %+v", b)
	}
	// Token dipakai ulang.
	_, _ = st.Simpan(ctx, lokUji, strings.NewReader("ISI"), "pdf", "application/pdf")
	if *nToken != 1 {
		t.Errorf("token diterbitkan %d kali", *nToken)
	}
	l.jawaban["upload"] = map[string]string{"URLImage": ""}
	if _, err := st.Simpan(ctx, lokUji, strings.NewReader("ISI"), "pdf", "application/pdf"); !errors.Is(err, ErrStorageJawabanRusak) {
		t.Errorf("URLImage kosong: %v", err)
	}
	for kode, mau := range map[int]error{400: ErrStorageMenolakPermintaan, 422: ErrStorageMenolakPermintaan,
		401: ErrStorageTokenDitolak, 403: ErrStorageTokenDitolak, 500: ErrStorageGagal, 404: ErrStorageGagal} {
		l.status["upload"] = kode
		_, err := st.Simpan(ctx, lokUji, strings.NewReader("ISI"), "pdf", "application/pdf")
		if !errors.Is(err, mau) || strings.Contains(err.Error(), l.srv.URL) || strings.Contains(err.Error(), "TOKEN-UJI") {
			t.Errorf("%d: %v", kode, err)
		}
	}
	if *nToken < 2 {
		t.Errorf("token ditolak tidak dilupakan (%d)", *nToken)
	}
}

// TestStorageJarakJauhBuka - URL tersimpan berlaku: tanpa geturl; lewat: geturl (Folder dari APPFOLDER, Namafile, Durasi)
// + penyegaran; EXPDATE kosong: tanpa URL; objek 404: tidak ada; URL bukan https ditolak; pengalihan tidak diikuti.
func TestStorageJarakJauhBuka(t *testing.T) {
	ctx := context.Background()
	l := layananStorageUji(t)
	l.objek["a"] = "ISI-OBJEK"
	kini := time.Date(2026, 10, 8, 6, 0, 0, 0, time.UTC)
	st, _, pc := storageJarakJauhUji(t, l, kini)
	link := models.LinkStorage{ImageID: lokUji.ImageID, URLPublic: l.srv.URL + "/objek/a", ExpUTC: "20261008063000",
		AppFolder: models.AppFolderUji("UJIAPP", "Policy/Doc/2026/10/"+lokUji.Namafile), App: "UJIAPP", Namafile: lokUji.Namafile}
	rc, err := st.Buka(ctx, link)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(rc)
	_ = rc.Close()
	if string(b) != "ISI-OBJEK" || len(l.minta["geturl"]) != 0 {
		t.Errorf("tersimpan: %q %d", b, len(l.minta["geturl"]))
	}
	// Lewat (sisa < 1 menit).
	link.ExpUTC = "20261008060030"
	l.jawaban["geturl"] = map[string]string{"URLImage": l.srv.URL + "/objek/a", "exp": "20261008T070000", "DateTime": "10/08/2026 13:00:00"}
	rc, err = st.Buka(ctx, link)
	if err != nil {
		t.Fatal(err)
	}
	_ = rc.Close()
	g := l.minta["geturl"]
	if len(g) != 1 || g[0]["Folder"] != "Policy/Doc/2026/10/" || g[0]["Namafile"] != lokUji.Namafile || g[0]["Durasi"] != float64(1800) ||
		len(pc.o) != 1 || pc.o[0].Exp != "08/10/2026 07:00:00" || pc.o[0].TanggalUpload != "10/08/2026 13:00:00" {
		t.Errorf("geturl %+v / %+v", g, pc.o)
	}
	if _, err := st.Buka(ctx, models.LinkStorage{ImageID: "X", URLPublic: l.srv.URL}); !errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
		t.Errorf("EXPDATE kosong: %v", err)
	}
	link.URLPublic, link.ExpUTC = l.srv.URL+"/objek/hilang", "20261008063000"
	if _, err := st.Buka(ctx, link); !errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
		t.Errorf("objek 404: %v", err)
	}
	link.ExpUTC = "20261008060030"
	l.jawaban["geturl"] = map[string]string{"URLImage": (&url.URL{Scheme: "http", Host: "polos", Path: "/objek/a"}).String()}
	if _, err := st.Buka(ctx, link); !errors.Is(err, ErrStorageJawabanRusak) {
		t.Errorf("bukan https: %v", err)
	}
	l.jawaban["geturl"] = map[string]string{"URLImage": l.srv.URL + "/alih"}
	if _, err := st.Buka(ctx, link); !errors.Is(err, ErrStorageGagal) {
		t.Errorf("pengalihan: %v", err)
	}
}

// TestStorageJarakJauhHapusAda - delete: Namafile = APPFOLDER tanpa "gs://<App>/" (tanpa Folder); Ada: geturl + Range 1 byte.
func TestStorageJarakJauhHapusAda(t *testing.T) {
	ctx := context.Background()
	l := layananStorageUji(t)
	st, _, _ := storageJarakJauhUji(t, l, time.Now())
	link := models.LinkStorage{ImageID: lokUji.ImageID, AppFolder: models.AppFolderUji("UJIAPP", "Policy/Doc/2026/10/"+lokUji.Namafile), App: "UJIAPP",
		Namafile: lokUji.Namafile}
	if err := st.Hapus(ctx, link); err != nil {
		t.Fatal(err)
	}
	if d := l.minta["delete"]; len(d) != 1 || d[0]["Namafile"] != "Policy/Doc/2026/10/"+lokUji.Namafile || d[0]["Folder"] != nil ||
		d[0]["App"] != "UJIAPP" {
		t.Errorf("delete %+v", d)
	}
	if err := st.Hapus(ctx, models.LinkStorage{ImageID: "X"}); !errors.Is(err, ErrBerkasTidakAdaDiPenyimpanan) {
		t.Errorf("APPFOLDER kosong: %v", err)
	}
	l.jawaban["geturl"] = map[string]string{"URLImage": l.srv.URL + "/objek/b", "exp": "20261008T070000"}
	if ada, _, err := st.Ada(ctx, lokUji); err != nil || ada {
		t.Errorf("tidak ada: %v %v", ada, err)
	}
	l.objek["b"] = "ISI"
	ada, o, err := st.Ada(ctx, lokUji)
	if err != nil || !ada || o.Exp != "08/10/2026 07:00:00" || o.URLPublic == "" || l.unduhan[len(l.unduhan)-1] != "/objek/b|bytes=0-0" {
		t.Errorf("ada: %v %v %+v %q", ada, err, o, l.unduhan)
	}
}

// penyimpanTokenTiruan - GCP_IMAGE dalam memori.
type penyimpanTokenTiruan struct {
	app, lama string
	simpan    []string
}

func (p *penyimpanTokenTiruan) AppStorage(context.Context) (string, error) { return p.app, nil }
func (p *penyimpanTokenTiruan) TokenStorageBerlaku(context.Context, *db.Tx, string, time.Time, time.Duration) (string, time.Duration, error) {
	if p.lama != "" {
		return p.lama, 40 * time.Second, nil
	}
	return "", 0, nil
}
func (p *penyimpanTokenTiruan) SimpanToken(_ context.Context, _ *db.Tx, app, token, pengguna string, _ time.Time) error {
	p.simpan = append(p.simpan, app+"|"+pengguna+"|"+token)
	return nil
}

// TestSumberTokenStorage - GET_TOKEN_STORAGE tanpa CALL: token berlaku dipakai ulang; bila tidak ada dirakit + disimpan
// (USERINPUT "Job"); garam / APPNAME kosong ditolak.
func TestSumberTokenStorage(t *testing.T) {
	ctx := context.Background()
	kini := time.Date(2026, 10, 8, 6, 0, 0, 0, time.UTC)
	p := &penyimpanTokenTiruan{app: "UJIAPP", lama: "LAMA"}
	s := sumberTokenStorage{transaksi: tanpaTx, penyimp: p, garam: "GARAM-UJI", jam: func() time.Time { return kini }}
	tok, exp, err := s.TokenBaru(ctx)
	if err != nil || tok != "LAMA" || !exp.Equal(kini.Add(40*time.Second)) || len(p.simpan) != 0 {
		t.Errorf("lama: %v %q %v", err, tok, exp)
	}
	p.lama = ""
	tok, exp, err = s.TokenBaru(ctx)
	baru, _ := layanan.RakitToken("GARAM-UJI", kini)
	if err != nil || tok != baru || !exp.Equal(kini.Add(layanan.UmurToken)) || len(p.simpan) != 1 || p.simpan[0] != "UJIAPP|Job|"+baru {
		t.Errorf("baru: %v %q %v %q", err, tok, exp, p.simpan)
	}
	if _, _, err := (sumberTokenStorage{transaksi: tanpaTx, penyimp: p}).TokenBaru(ctx); !errors.Is(err, layanan.ErrGaramTokenKosong) {
		t.Errorf("garam kosong: %v", err)
	}
	p.app = " "
	if _, _, err := s.TokenBaru(ctx); !errors.Is(err, layanan.ErrAppNameKosong) {
		t.Errorf("APPNAME kosong: %v", err)
	}
}
