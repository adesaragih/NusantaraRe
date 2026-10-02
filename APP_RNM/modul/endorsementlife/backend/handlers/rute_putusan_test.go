package handlers_test

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/handlers"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

func TestRutePutuskan(t *testing.T) {
	g := gudangSatuKasus()
	g.Polis["NBLF-1"] = &tiruan.Polis{ID: "NBLF-1", ProdKe: 1, Kepala: map[string]string{"TYPE": "QR"}}
	g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: "UJI-NB-1", PolisID: "NBLF-1", PLNumber: "UJI-PL-1", Nilai: map[string]string{}})
	g.Polis["EDMLF-1"].ProdKe = 2
	g.Peserta[0].Nilai["CURRENCY"] = "IDR"
	h := routerBuat(g)
	jalur := handlers.Prefix + "/kasus/EDMLF-1/putuskan"
	if w := minta(t, h, "POST", jalur, `{"status":"1"}`, false); w.Code != http.StatusUnauthorized {
		t.Fatalf("tanpa identitas %d", w.Code)
	}
	if w := minta(t, h, "POST", jalur, `{"status":"1"}`, true); w.Code != http.StatusConflict {
		t.Fatalf("belum Save %d %s", w.Code, w.Body)
	}
	if w := minta(t, h, "POST", handlers.Prefix+"/kasus/EDMLF-1/simpan", `{}`, true); w.Code != http.StatusOK {
		t.Fatalf("simpan %d %s", w.Code, w.Body)
	}
	if w := minta(t, h, "POST", jalur, `{"status":"3"}`, true); w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status asing %d %s", w.Code, w.Body)
	}
	if w := minta(t, h, "POST", jalur, `{"status":`, true); w.Code != http.StatusBadRequest {
		t.Fatalf("badan rusak %d", w.Code)
	}
	w := minta(t, h, "POST", jalur, `{"status":"1","comment":"UJI"}`, true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"status":"Resolved-Completed","noEndorsement":"UJI-PL-1/01"`) {
		t.Fatalf("confirm %d %s", w.Code, w.Body)
	}
	if g.Polis["EDMLF-1"].Status != models.StatusKasusSelesai {
		t.Error("kasus tidak resmi")
	}
	// K5 (01-10-2026): peserta kasus ikut tertulis ke tabel warisan - cacahnya di jawaban Confirm.
	if !strings.Contains(w.Body.String(), `"pesertaWarisan":1`) || len(g.PesertaWarisanTertulis) != 1 {
		t.Errorf("peserta warisan %s (%d baris)", w.Body, len(g.PesertaWarisanTertulis))
	}
	w = minta(t, h, "GET", handlers.Prefix+"/kasus/EDMLF-1", "", true)
	if !strings.Contains(w.Body.String(), `"riwayat":[{"no":1,`) || !strings.Contains(w.Body.String(), `"status":"Accept","comment":"UJI"`) {
		t.Errorf("riwayat tidak terbaca %s", w.Body)
	}
	if w := minta(t, h, "POST", jalur, `{"status":"2"}`, true); w.Code != http.StatusConflict {
		t.Fatalf("putuskan ulang %d %s", w.Code, w.Body)
	}
}
