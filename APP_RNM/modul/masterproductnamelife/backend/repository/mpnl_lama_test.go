package repository

// Copy Old (permintaan work owner 03-10-2026): daftar produk lama yang belum ada di tabel flat dan status salinnya -
// aturan yang SAMA dengan alat pindah (rekonsiliasi, keputusan OQ-FLAT-07/09, K3).

import (
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

func cariLama(t *testing.T, daftar []models.ProdukLama, id string) models.ProdukLama {
	t.Helper()
	for _, d := range daftar {
		if d.ID == id {
			return d
		}
	}
	t.Fatalf("produk %s tidak ada di daftar %+v", id, daftar)
	return models.ProdukLama{}
}

func TestStatusLama(t *testing.T) {
	umum := []barisJSON{
		barisJ("100044", `{"ID":"100044","PRODUCTNAME":"UJI BERSIH","CEDING":"UJI CEDING","TREATYNUMBER":"UJI-1","INWARDNAME":"UJI TREATY","CREATEOP":"UJI-OP"}`),
		barisJ("100045", `{"ID":"100045","PRODUCTNAME":"UJI SUDAH DISALIN"}`),
		barisJ("100046", `{"ID":"100046","PRODUCTNAME":"UJI MEDAN MATI","TYPE":"UJI-TIPE"}`),
		barisJ("100047", `{"ID":"100047","PRODUCTNAME":"UJI KOMA","RICOMM":"1,5"}`),
		barisJ("100048", `{"ID":"100048","PRODUCTNAME":"UJI SPASI","RICOMM":" 5"}`),
		barisJ("100049", `{"ID":"100049","PRODUCTNAME":"UJI K3","UnderwritingLimitList":[{"MaxInsured":"UJI-BUKAN-ANGKA"}]}`),
		barisJ("100050", `{"ID":"100050","PRODUCTNAME":"UJI TANGGAL"}`),
		barisJ("100051", `{rusak`),
	}
	inward := []barisJSON{
		barisJ("100044", `{"ID":"100044","PRODUCTID":"100044","INSURED":"UJI INSURED"}`),
		barisJ("100050", `{"ID":"100050","PRODUCTID":"100050","MATURE":"1/3/2027"}`),
	}
	daftar, siap := statusLama(umum, inward, map[string]bool{"100045": true})

	if len(daftar) != 7 {
		t.Fatalf("tujuh produk lama belum ada di tabel flat (100045 sudah): %d %+v", len(daftar), daftar)
	}
	for _, d := range daftar {
		if d.ID == "100045" {
			t.Error("produk yang sudah ada di tabel flat tidak didaftar")
		}
	}

	bersih := cariLama(t, daftar, "100044")
	if !bersih.BolehDisalin || len(bersih.Alasan) != 0 || bersih.ProductName != "UJI BERSIH" || bersih.Ceding != "UJI CEDING" ||
		bersih.TreatyNumber != "UJI-1" || bersih.InwardName != "UJI TREATY" || bersih.CreateOp != "UJI-OP" {
		t.Errorf("produk bersih: %+v", bersih)
	}
	if p := siap["100044"]; p.ID != "100044" || p.Inward.Insured != "UJI INSURED" {
		t.Errorf("bentuk flat produk bersih siap ditulis: %+v", p)
	}

	if m := cariLama(t, daftar, "100046"); m.BolehDisalin || !strings.Contains(strings.Join(m.Alasan, " "), "TYPE") {
		t.Errorf("medan tanpa kolom flat yang terisi menolak salin: %+v", m)
	}
	if k := cariLama(t, daftar, "100047"); !k.BolehDisalin || !strings.Contains(strings.Join(k.Catatan, " "), "koma desimal") {
		t.Errorf("koma desimal (diputuskan OQ-FLAT-07) boleh, dicatat: %+v", k)
	}
	if s := cariLama(t, daftar, "100048"); s.BolehDisalin || !strings.Contains(strings.Join(s.Alasan, " "), "spasi tepi") {
		t.Errorf("jenis normalisasi yang belum diputuskan menolak salin: %+v", s)
	}
	if k3 := cariLama(t, daftar, "100049"); !k3.BolehDisalin || !strings.Contains(strings.Join(k3.Catatan, " "), "MAXINSURED") {
		t.Errorf("nilai tidak sah (K3) dikosongkan dan dicatat, salin tetap boleh: %+v", k3)
	}
	if tg := cariLama(t, daftar, "100050"); !tg.BolehDisalin || !strings.Contains(strings.Join(tg.Catatan, " "), "MATURE") {
		t.Errorf("tanggal bentuk lain dikonversi (OQ-FLAT-09) dan dicatat: %+v", tg)
	}
	if r := cariLama(t, daftar, "100051"); r.BolehDisalin || len(r.Alasan) == 0 {
		t.Errorf("JSON rusak menolak salin: %+v", r)
	}
	if _, ada := siap["100046"]; ada {
		t.Error("produk yang ditolak tidak disiapkan untuk ditulis")
	}
	// Alasan dan catatan menyebut kolom, tidak pernah nilainya.
	for _, d := range daftar {
		semua := strings.Join(append(append([]string{}, d.Alasan...), d.Catatan...), " ")
		for _, nilai := range []string{"UJI-TIPE", "1,5", "UJI-BUKAN-ANGKA", "1/3/2027"} {
			if strings.Contains(semua, nilai) {
				t.Errorf("%s: alasan/catatan memuat nilai data %q: %s", d.ID, nilai, semua)
			}
		}
	}
	// Urut ID, seperti grid produk.
	for i := 1; i < len(daftar); i++ {
		if daftar[i-1].ID > daftar[i].ID {
			t.Errorf("urut ID: %s sebelum %s", daftar[i-1].ID, daftar[i].ID)
		}
	}
}
