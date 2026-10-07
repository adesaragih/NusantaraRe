package handlers_test

// Copy Old NB Treaty In (perintah work owner 07-10-2026 "nb ttreatyin tobol copy untuk data lama mana?" -> "langsung
// anda kerjakan!") - seam HTTP di atas gudang tiruan + JSON_POLIS tiruan. Gerbang superadmin = pemegang menu Kelola
// User dengan menu NB Treaty In ber-hak PENUH (pola Bordereaux / Copy Old EDM Treaty In). Fixture UJI- (dokumen dari
// `models/testdata/dokumen_uji_prop.json`).

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/menu"
	"nusantarare/modul/nbtreatyin/backend/handlers"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
	"nusantarare/modul/nbtreatyin/backend/tiruan"
)

// gudangLama - tiruan modul + JSON_POLIS + dua penulis tambahan pemuat (`repository/lama.go`).
type gudangLama struct {
	*tiruan.Gudang
	urut  []string
	dok   map[string]models.BarisJSONPolis
	datar map[string]models.KolomDatarLama
	// tanpaProduksi - ROWID yang IDPEGA-nya tidak ada di tabel kerja Pega / TREATYINPRODUCTION.
	tanpaProduksi map[string]bool
}

var _ services.GudangPemuat = (*gudangLama)(nil)

func (g *gudangLama) HitungJSONPolisLain(context.Context) (int, error) { return 0, nil }
func (g *gudangLama) KunciJSONPolis(context.Context) ([]string, error) { return g.urut, nil }

// KunciJSONPolisCopyOld - saringan WO 07-10-2026 (kasus Pega ada DAN sudah berproduksi): baris di `tanpaProduksi`
// tidak lolos.
func (g *gudangLama) KunciJSONPolisCopyOld(context.Context) ([]string, error) {
	var out []string
	for _, k := range g.urut {
		if !g.tanpaProduksi[k] {
			out = append(out, k)
		}
	}
	return out, nil
}
func (g *gudangLama) BacaJSONPolis(_ context.Context, k string) (models.BarisJSONPolis, error) {
	b, ada := g.dok[k]
	if !ada {
		return b, fmt.Errorf("uji: ROWID %s tidak ada", k)
	}
	return b, nil
}
func (g *gudangLama) AdaKasus(_ context.Context, _ *db.Tx, id string) (bool, error) {
	_, ada := g.Kasus[id]
	return ada, nil
}
func (g *gudangLama) SetelKolomDatarLama(_ context.Context, _ *db.Tx, id string, k models.KolomDatarLama) error {
	g.datar[id] = k
	return nil
}
func (g *gudangLama) SalinUsulanLama(ctx context.Context, tx *db.Tx, idPega string, baris []models.UsulanProduksi) (models.NasibUsulan, error) {
	if len(baris) == 0 {
		return models.UsulanTanpaBaris, nil
	}
	return models.UsulanDisalin, g.CatatUsulan(ctx, tx, idPega, baris)
}

const polisUjiLama = "UJI-QP.T1.10.2017.00001"

func gudangLamaUji(t *testing.T) *gudangLama {
	t.Helper()
	dok, err := os.ReadFile("../models/testdata/dokumen_uji_prop.json")
	if err != nil {
		t.Fatal(err)
	}
	baris := func(idpega, prodke string, data []byte) models.BarisJSONPolis {
		return models.BarisJSONPolis{IDPega: idpega, NoPolis: polisUjiLama, ProdKe: prodke, TglInput: "2017-10-02 08:00:00",
			TglProd: "2017-10-02 08:00:00", Username: "UJI-AKUN", DataJSON: data}
	}
	g := &gudangLama{Gudang: tiruan.Baru(), dok: map[string]models.BarisJSONPolis{}, datar: map[string]models.KolomDatarLama{}}
	for _, x := range []struct {
		rowid string
		b     models.BarisJSONPolis
	}{
		{"R1", baris("ASM-FW-GISFW-WORK-NB NB-990001", "0", dok)},
		{"R2", baris("ASM-FW-GISFW-WORK-NB NB-990002", "0", []byte(`{`))},
		{"R3", baris("NB-5", "0", nil)}, // tulisan Utility1 aplikasi baru
		{"R4", baris("ASM-FW-GISFW-WORK-LAIN LAIN-1", "0", []byte(`{"pxObjClass":"UJI-Lain"}`))},
		{"R5", baris("ASM-FW-GISFW-WORK-NB NB-990005", "1", dok)}, // generasi endorsemen - milik EDM
		{"R6", baris("ASM-FW-GISFW-WORK-NB NB-990006", "0", dok)}, // sudah di tabel flat
		{"R7", baris("ASM-FW-GISFW-WORK-NB NB-990007", "0", dok)}, // belum berproduksi (saringan WO)
	} {
		g.urut = append(g.urut, x.rowid)
		g.dok[x.rowid] = x.b
	}
	g.Kasus["NB-990006"] = models.Kasus{ID: "NB-990006"}
	g.tanpaProduksi = map[string]bool{"R7": true}
	return g
}

func mintaMenu(h http.Handler, metode, jalur, badan, akun string, lihat []string, menuAkun ...string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(metode, handlers.Prefix+jalur, strings.NewReader(badan))
	r.Header.Set("X-Pelaku", akun)
	ctx := inti.DenganMenuLihat(inti.DenganAksesMenu(r.Context(), menuAkun), lihat)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r.WithContext(ctx))
	return w
}

func daftarLama(t *testing.T, h http.Handler) []models.DokumenLama {
	t.Helper()
	w := mintaMenu(h, "GET", "/lama", "", "UJI-SUPER", nil, handlers.KodeMenu, menu.KodeKelolaUser)
	var d []models.DokumenLama
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &d) != nil {
		t.Fatalf("daftar Copy Old %d %s", w.Code, w.Body.String())
	}
	return d
}

func TestCopyOldNBDaftar(t *testing.T) {
	h := handlers.Router(services.Baru(gudangLamaUji(t), jam), true)
	d := daftarLama(t, h)
	if len(d) != 2 || d[0].ID != "NB-990001" || d[1].ID != "NB-990002" {
		t.Fatalf("popup hanya dokumen NB lama yang belum di tabel flat: %+v", d)
	}
	if !d[0].BolehDisalin || d[0].NoPolis != polisUjiLama || d[0].TglProd == "" {
		t.Fatalf("baris dokumen sah %+v", d[0])
	}
	if d[1].BolehDisalin || len(d[1].Alasan) != 1 || d[1].Alasan[0] != "the old JSON cannot be read" {
		t.Fatalf("dokumen rusak: %+v", d[1])
	}
}

func TestCopyOldNBHanyaSuperadmin(t *testing.T) {
	h := handlers.Router(services.Baru(gudangLamaUji(t), jam), true)
	nb := handlers.KodeMenu
	if w := mintaMenu(h, "GET", "/hak", "", "UJI-ADMIN", nil, nb); w.Code != 200 || !strings.Contains(w.Body.String(), `"copyOld":false`) {
		t.Fatalf("hak admin %d %s", w.Code, w.Body.String())
	}
	if w := mintaMenu(h, "GET", "/lama", "", "UJI-ADMIN", nil, nb); w.Code != http.StatusForbidden {
		t.Fatalf("bukan superadmin %d", w.Code)
	}
	if w := mintaMenu(h, "POST", "/lama/salin", `{"ids":["NB-990001"]}`, "UJI-ADMIN", nil, nb); w.Code != http.StatusForbidden {
		t.Fatalf("bukan superadmin menyalin %d", w.Code)
	}
	if w := mintaMenu(h, "GET", "/hak", "", "UJI-SUPER", nil, nb, menu.KodeKelolaUser); !strings.Contains(w.Body.String(), `"copyOld":true`) {
		t.Fatalf("hak superadmin %s", w.Body.String())
	}
	// View only berlaku juga bagi superadmin (keputusan work owner 05-10-2026)
	if w := mintaMenu(h, "GET", "/lama", "", "UJI-SUPER", []string{nb}, nb, menu.KodeKelolaUser); w.Code != http.StatusForbidden {
		t.Fatalf("superadmin View only %d", w.Code)
	}
	// gudang tanpa pembaca JSON_POLIS (tiruan biasa): tombol tidak tampil walau superadmin
	polos := handlers.Router(services.Baru(tiruan.Baru(), jam), true)
	if w := mintaMenu(polos, "GET", "/hak", "", "UJI-SUPER", nil, nb, menu.KodeKelolaUser); !strings.Contains(w.Body.String(), `"copyOld":false`) {
		t.Fatalf("tanpa pembaca JSON_POLIS %s", w.Body.String())
	}
}

func TestCopyOldNBProcessCopy(t *testing.T) {
	g := gudangLamaUji(t)
	h := handlers.Router(services.Baru(g, jam), true)
	if w := mintaMenu(h, "POST", "/lama/salin", `{"ids":[" "]}`, "UJI-SUPER", nil, handlers.KodeMenu, menu.KodeKelolaUser); w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("tanpa pilihan %d", w.Code)
	}
	w := mintaMenu(h, "POST", "/lama/salin", `{"ids":["NB-990002","NB-990001","NB-404","NB-990001"]}`, "UJI-SUPER", nil,
		handlers.KodeMenu, menu.KodeKelolaUser)
	var j models.JawabanSalinLama
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &j) != nil {
		t.Fatalf("Process Copy %d %s", w.Code, w.Body.String())
	}
	status := map[string]string{}
	for _, x := range j.Hasil {
		status[x.ID] = x.Status
		for _, p := range x.Pesan {
			if strings.Contains(p, "UJI-QP") {
				t.Fatalf("pesan layar memuat nomor polis: %q", p)
			}
		}
	}
	if j.Disalin != 1 || status["NB-990001"] != models.SalinDisalin || status["NB-990002"] != models.SalinDitolak ||
		status["NB-404"] != models.SalinSudahAda || len(j.Hasil) != 3 {
		t.Fatalf("hasil %+v", j)
	}
	if k := g.Kasus["NB-990001"]; k.StatusWork != models.StatusSelesai {
		t.Fatalf("kasus tersalin lewat pemuat ditutup Resolved-Completed: %+v", k)
	}
	if g.datar["NB-990001"].Username != "UJI-AKUN" {
		t.Fatalf("kolom datar json_polis: %+v", g.datar["NB-990001"])
	}
	// sesudah disalin: hilang dari popup
	for _, d := range daftarLama(t, h) {
		if d.ID == "NB-990001" {
			t.Fatal("dokumen tersalin masih di popup")
		}
	}
}
