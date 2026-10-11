package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// urlUjiOfficeH - URL penyimpanan tiruan (disusun lewat url.URL: ADR-U-0013 melarang literal alamat di kode).
func urlUjiOfficeH(kueri string) string {
	return (&url.URL{Scheme: "https", Host: "penyimpanan.uji", Path: "/objek", RawQuery: kueri}).String()
}

// lampiranNyataH - catatan T_STORAGE_IMAGE ber-URLPUBLIC (penyimpanan nyata).
type lampiranNyataH struct{ *lampiranTiruanH }

func (l lampiranNyataH) LinkStorage(ctx context.Context, tx *db.Tx, imageID string) (models.LinkStorage, bool, error) {
	k, ada, err := l.lampiranTiruanH.LinkStorage(ctx, tx, imageID)
	k.URLPublic = urlUjiOfficeH("")
	return k, ada, err
}

// storagePenautH - klien penyimpanan nyata tiruan (Tautan).
type storagePenautH struct{ services.KlienStorageLampiran }

func (storagePenautH) Tautan(context.Context, models.LinkStorage) (string, error) {
	return urlUjiOfficeH("tanda=UJI"), nil
}

// TestOfficeLampiranHandler - GET …/office: 200 {url} + Cache-Control no-store; bukan kantor 400; stub lokal 409; belum tersimpan
// 409; tidak ada 404; tanpa identitas 401.
func TestOfficeLampiranHandler(t *testing.T) {
	jalur := "/api/nbfacin/kasus/UJI-NB-1/lampiran/berkas/"
	svc, repo, _ := layananLampiranH(t, "0")
	repo.dok["X1"] = models.DokumenPolis{ID: "X1", IDPega: "ASM-FW-GISFW-WORK UJI-NB-1", NamaFile: "a.xlsx", Mime: "xlsx", TStorageID: imageUjiH}
	repo.dok["X2"] = models.DokumenPolis{ID: "X2", IDPega: "ASM-FW-GISFW-WORK UJI-NB-1", NamaFile: "b.docx", Mime: "docx", TStorageID: "BELUM"}
	for nama, c := range map[string]struct {
		id   string
		kode int
	}{"pdf": {"20261008010405007", 400}, "stub lokal": {"X1", 409}, "belum tersimpan": {"X2", 409}, "tidak ada": {"X9", 404}} {
		if kode, isi := mintaPeran(t, svc, "GET", jalur+c.id+"/office", "", "UJI-MKT", "ReasFacInMarketing"); kode != c.kode {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	if kode, _ := mintaPeran(t, svc, "GET", jalur+"X1/office", "", "", ""); kode != 401 {
		t.Errorf("401: %d", kode)
	}
	svc.DenganLampiranUji(lampiranNyataH{repo}, &antreanDiamH{}, storagePenautH{services.StorageLampiranLokal(t.TempDir())}, t.TempDir())
	w := mintaRekam(t, svc, "GET", jalur+"X1/office", "UJI-UW", "ReasFacInUnderwriting")
	if w.Code != 200 || strings.TrimSpace(w.Body.String()) != `{"url":"`+urlUjiOfficeH("tanda=UJI")+`"}` || w.Header().Get("Cache-Control") != "no-store" {
		t.Errorf("200: %d %s %v", w.Code, w.Body.String(), w.Header())
	}
}

// mintaRekam - mintaPeran dengan rekaman utuh (header jawaban).
func mintaRekam(t *testing.T, svc *services.Service, metode, jalur, pelaku, peran string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, true)
	r := httptest.NewRequest(metode, jalur, nil)
	r.Header.Set("X-Pelaku", pelaku)
	r.Header.Set("X-Peran", peran)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w
}
