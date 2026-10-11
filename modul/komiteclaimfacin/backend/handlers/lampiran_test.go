package handlers_test

// Lampiran klaim induk di layar komite (perintah work owner 11-10-2026): hanya kategori yang berisi, dokumen bisa
// diunduh / dilihat / View Office Online, tanpa unggah.

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimfacin/backend/models"
	"nusantarare/modul/komiteclaimfacin/backend/tiruan"
)

const pelihatLampiran = "UJI-PELIHAT"

func siapLampiran(t *testing.T) *uji {
	t.Helper()
	u := siap(t, false)
	u.g.SiapkanTT2("1000", saatUji)
	u.g.Klaim.SetelLampiran(tiruan.KlaimUji, kontrak.DaftarLampiranKlaim{
		Kategori: []kontrak.KategoriLampiranKlaim{
			{ID: "ADU", Label: "Approval", Jumlah: 1},
			{ID: "AcceptanceNote", Label: "AcceptanceNote", Jumlah: 0},
			{ID: "File", Label: "File", Jumlah: 2},
			{ID: "LOD", Label: "LOD", Jumlah: 0},
		},
		Lampiran: []kontrak.LampiranKlaim{
			{ID: "UJI-DOK-3", NamaFile: "uji-tabel.xlsx", Kategori: "File", MIME: "xlsx", AdaObjek: true},
			{ID: "UJI-DOK-2", NamaFile: "uji-surat.pdf", Kategori: "ADU", MIME: "pdf", Tanggal: "11/10/2026 09:30",
				Operator: "UJI-ADMIN", AdaObjek: true},
			{ID: "UJI-DOK-1", NamaFile: "uji-tanpa-objek.pdf", Kategori: "File", MIME: "pdf"},
			{ID: "UJI-DOK-0", NamaFile: "uji-luar-master.pdf", Kategori: "BukanMaster", MIME: "pdf", AdaObjek: true},
		},
	}, map[string]string{"UJI-DOK-2": "%PDF-UJI"})
	return u
}

func TestLampiranKomiteHanyaKategoriBerisi(t *testing.T) {
	u := siapLampiran(t)
	w := u.minta("GET", "/kasus/"+tiruan.KomiteUji+"/lampiran", pelihatLampiran, "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("lampiran: %d %s", w.Code, w.Body.String())
	}
	var out models.LampiranKomite
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	mau := []models.KategoriLampiran{{ID: "ADU", Label: "Approval", CountAttach: 1}, {ID: "File", Label: "File",
		CountAttach: 2}}
	if !reflect.DeepEqual(out.Kategori, mau) {
		t.Fatalf("kategori %+v, mau hanya yang berisi %+v", out.Kategori, mau)
	}
	var ids []string
	for _, a := range out.Lampiran {
		ids = append(ids, a.ID)
	}
	if strings.Join(ids, ",") != "UJI-DOK-3,UJI-DOK-2,UJI-DOK-1" {
		t.Fatalf("dokumen %v (urutan penyedia, tanpa kategori di luar grid)", ids)
	}
	if a := out.Lampiran[1]; a.Tanggal != "11/10/2026 09:30" || a.Operator != "UJI-ADMIN" || !a.AdaObjek || a.MIME != "pdf" {
		t.Fatalf("dokumen %+v", a)
	}
}

func TestLampiranKomiteKosongBukanNull(t *testing.T) {
	u := siap(t, false)
	u.g.SiapkanTT2("1000", saatUji)
	w := u.minta("GET", "/kasus/"+tiruan.KomiteUji+"/lampiran", pelihatLampiran, "", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"kategori":[]`) ||
		!strings.Contains(w.Body.String(), `"lampiran":[]`) {
		t.Fatalf("lampiran kosong: %d %s", w.Code, w.Body.String())
	}
}

func TestLampiranKomiteUnduhDanOffice(t *testing.T) {
	u := siapLampiran(t)
	jalur := "/kasus/" + tiruan.KomiteUji + "/lampiran/"
	w := u.minta("GET", jalur+"UJI-DOK-2/isi", pelihatLampiran, "", nil)
	if w.Code != http.StatusOK || w.Body.String() != "%PDF-UJI" ||
		!strings.Contains(w.Header().Get("Content-Disposition"), "uji-surat.pdf") ||
		w.Header().Get("Content-Type") != "application/pdf" {
		t.Fatalf("unduh: %d %q %v", w.Code, w.Body.String(), w.Header())
	}
	for _, c := range []struct {
		jalur string
		mau   int
	}{
		{"UJI-DOK-9/isi", http.StatusNotFound},               // bukan dokumen klaim ini
		{"UJI-DOK-1/isi", http.StatusConflict},               // objek penyimpanan tidak tercatat
		{"UJI-DOK-2/office", http.StatusUnprocessableEntity}, // pdf bukan Office
	} {
		if w := u.minta("GET", jalur+c.jalur, pelihatLampiran, "", nil); w.Code != c.mau {
			t.Errorf("%s: %d (mau %d) %s", c.jalur, w.Code, c.mau, w.Body.String())
		}
	}
	w = u.minta("GET", jalur+"UJI-DOK-3/office", pelihatLampiran, "", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), tiruan.AwalanTautanOffice+"UJI-DOK-3") {
		t.Fatalf("office: %d %s", w.Code, w.Body.String())
	}
}

func TestLampiranKomiteWajibIdentitasDanKasus(t *testing.T) {
	u := siapLampiran(t)
	if w := u.minta("GET", "/kasus/"+tiruan.KomiteUji+"/lampiran", "", "", nil); w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d", w.Code)
	}
	if w := u.minta("GET", "/kasus/KMT-UJI999/lampiran", pelihatLampiran, "", nil); w.Code != http.StatusNotFound {
		t.Errorf("kasus komite tidak ada: %d", w.Code)
	}
	if w := u.minta("GET", "/kasus/CLM-UJI001/lampiran", pelihatLampiran, "", nil); w.Code != http.StatusNotFound {
		t.Errorf("bukan kasus komite: %d", w.Code)
	}
}
