package services_test

// Regresi temuan /code-review 01-10-2026 (sesudah paket 11).

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

// #1 (dipindah 02-10-2026) - baris inward JSON milik produk lain kini urusan ALAT PINDAH: lihat repository
// `TestRekonsiliasiInwardMilikProdukLainTidakBersilang`. Di tabel flat sisi inward ada di baris induk produknya sendiri.

// #3 (diganti 02-10-2026) - `Copy` (`CopyProduct` menyalin halaman utuh) menyalin baris induk DAN seluruh anak produk
// asal: daftar milik server (komentar, outward) dari baris tersimpan, daftar form dari kiriman klien; kunci JSON tak
// dikelola tidak ada lagi (D2).
func TestCopyMenyalinIndukDanSeluruhAnak(t *testing.T) {
	l, g := layananSalin()
	m := produkMasuk()
	m.SalinanDari = "100007"
	m.DocumentClaim = []models.BarisDokumen{{Document: "UJI DOK"}}
	m.UnderwritingLimit = []models.BarisUWLimit{{MinInsured: "1", MaxInsured: "2", MinAge: "18", MaxAge: "65"}}
	m.CommentList = []models.BarisKomentar{{Suggest: "DIKARANG KLIEN"}}
	m.OutwardList = []models.BarisOutward{{ReinsTypeID: "KLIEN"}}
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil {
		t.Fatal(err)
	}
	s := g.Produk[p.ID]
	if s.Inward.ProductID != p.ID || len(s.DocumentClaim) != 1 || len(s.UnderwritingLimit) != 1 ||
		s.UnderwritingLimit[0].MaxAge != "65" {
		t.Errorf("induk + anak form tersimpan di salinan: %+v", s)
	}
	if len(s.CommentList) != 2 || s.CommentList[0].Suggest != "UJI KOMENTAR ASAL" || len(s.OutwardList) != 1 ||
		s.OutwardList[0].ReinsTypeID != "10200" {
		t.Errorf("COMMENT dan OUTWARD produk asal tersalin, kiriman klien diabaikan: %+v %+v", s.CommentList, s.OutwardList)
	}
}

// #4 - produk lama yang PLAN LIST-nya tidak disentuh tetap dapat disimpan; disentuh = diperiksa.
func TestProteksiPlanHanyaBilaDaftarPlanBerubah(t *testing.T) {
	l, g := layananMaster()
	g.Plan = []models.JenisPlan{{ID: "P1", CoverName: "UJI COVER", Business: "UJI BIZ", Benefit: "UJI MANFAAT"},
		{ID: "P2", CoverName: "UJI COVER DUA", Business: "UJI BIZ 2", Benefit: "UJI MANFAAT 2"}}
	g.IsiJSON("100007", `{"ID":"100007","PlanList":[{"Plan":"UJI COVER","PlanID":"P1","Name":"UJI BIZ","Benefit":"UJI MANFAAT",`+
		`"RIRATE":"","RIRATEID":""}]}`, "")
	simpan := func(m models.Produk) error {
		m.ID = "100007"
		_, err := l.SimpanProduk(context.Background(), pelakuUji, m, false)
		return err
	}
	lama := models.BarisPlan{Plan: "UJI COVER", PlanID: "P1", Name: "UJI BIZ", Benefit: "UJI MANFAAT"}
	m := produkMasuk()
	m.PlanList = []models.BarisPlan{lama}
	if err := simpan(m); err != nil {
		t.Errorf("baris plan warisan tanpa R/I Rate, tidak disentuh: simpan tetap jalan: %v", err)
	}
	m.PlanList = []models.BarisPlan{lama, {Plan: "UJI COVER DUA", PlanID: "P2"}}
	if err := simpan(m); err == nil || !strings.Contains(err.Error(), "PLAN LIST row 1: "+services.PesanRIRatePlanKosong) {
		t.Errorf("PLAN LIST berubah: ProteksiPlanListLife berjalan atas SEMUA baris: %v", err)
	}
}

// #13 - duplikat diperiksa atas nama plan SESUDAH diseragamkan dengan master.
func TestPlanSamaSesudahDiseragamkanMaster(t *testing.T) {
	_, simpan := layananPlan()
	m := produkMasuk()
	m.PlanList = []models.BarisPlan{planTersimpan, {Plan: "teks lain", PlanID: "P1", RIRate: "UJI RATE", RIRateID: "R1"}}
	if err := simpan(m); err == nil || !strings.Contains(err.Error(), "PLAN LIST row 2: "+services.PesanPlanSama) {
		t.Errorf("dua baris PlanID sama = plan sama: %v", err)
	}
}

// #15 - `asli` hanya boleh dikembalikan dari baris tersimpan, tidak dikarang. Sejak tabel flat (D2, 02-10-2026) baris
// tersimpan tidak membawa kunci JSON lama sama sekali: `asli` tersimpan selalu kosong, setiap `asli` berisi dari klien
// ditolak.
func TestAsliBarisHarusDariBarisTersimpan(t *testing.T) {
	l, g := layananMaster()
	g.IsiJSON("100007", `{"ID":"100007","LienClause":[{"Usia":"60","Manfaat":"50","UJI_LAMA":"1"}]}`, "")
	asal, err := l.AmbilProduk(context.Background(), pelakuUji, "100007")
	if err != nil || len(asal.LienClause) != 1 || asal.LienClause[0].Asli != "" {
		t.Fatalf("kunci baris JSON lama tidak pindah ke tabel flat: %+v %v", asal.LienClause, err)
	}
	m := produkMasuk()
	m.ID = "100007"
	m.LienClause = []models.BarisLien{asal.LienClause[0], {Usia: "70", Manfaat: "25", Asli: `{"OUTWARDRATEID":"UJI-SUSUP"}`}}
	_, err = l.SimpanProduk(context.Background(), pelakuUji, m, false)
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), "row 2") || strings.Contains(err.Error(), "row 1") {
		t.Errorf("asli karangan ditolak, asli tersimpan diterima: %v", err)
	}
	m.LienClause = m.LienClause[:1]
	m.LienClause[0].Manfaat = "40"
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, false); err != nil {
		t.Fatalf("baris tersimpan yang diubah membawa asli-nya: %v", err)
	}
	if s := g.Produk["100007"].LienClause; len(s) != 1 || s[0].Manfaat != "40" || s[0].Asli != "" {
		t.Errorf("baris LIEN tersimpan tanpa kunci karangan: %+v", s)
	}
}

// #7 - satu kirim ulang = satu percobaan.
func TestUlangiSatuKaliSatuPercobaan(t *testing.T) {
	l, g := layananMaster()
	g.IsiJSON("100007", `{"ID":"100007"}`, "")
	g.AppName = "UJI-APP"
	b := baruBerkasPalsu()
	l = l.DenganJam(func() time.Time { return jamLampiran }).DenganPenyimpanan(b)
	b.Gagal = errors.New("UJI tidak terjangkau")
	a, err := l.UnggahLampiran(context.Background(), pelakuUji, "100007", "UJI.pdf", strings.NewReader("isi"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := l.UlangiLampiran(context.Background(), pelakuUji, "100007", a.ID); err != nil {
		t.Fatal(err)
	}
	if len(g.Outbox) != 1 || g.Outbox[0].Percobaan != 2 || b.kirim != 2 {
		t.Errorf("unggah + satu ulang = dua percobaan, satu efek: %+v kirim=%d", g.Outbox, b.kirim)
	}
}

// #6 - `Download All` tidak gagal karena lampiran lama yang berkasnya tidak di stub.
func TestUnduhSemuaMenyatakanBerkasTakTersedia(t *testing.T) {
	l, b, _ := layananLampiran(t)
	a1, _ := unggah(l, "UJI satu.pdf", "satu")
	a2, _ := unggah(l, "UJI lama.pdf", "lama")
	delete(b.simpan, a2.StorageID) // lampiran Pega: rekam + objek ada, berkas di penyimpanan asal
	var buf bytes.Buffer
	n, err := l.UnduhSemuaLampiran(context.Background(), pelakuUji, "100007", &buf)
	if err != nil || n != 1 {
		t.Fatalf("zip tetap dibuat: %d %v", n, err)
	}
	z, _ := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	nama := map[string]string{}
	for _, f := range z.File {
		r, _ := f.Open()
		isi, _ := io.ReadAll(r)
		_ = r.Close()
		nama[f.Name] = string(isi)
	}
	if nama["UJI satu.pdf"] != "satu" || !strings.Contains(nama[services.NamaDaftarTakTersedia], "UJI lama.pdf") {
		t.Errorf("isi zip: %v", nama)
	}
	if _, err := l.UnduhLampiran(context.Background(), pelakuUji, "100007", a2.ID); !errors.Is(err, services.ErrBerkasTidakDiStub) {
		t.Errorf("unduh satu berkas lama: %v", err)
	}
	delete(b.simpan, a1.StorageID)
	if _, err := l.UnduhSemuaLampiran(context.Background(), pelakuUji, "100007", &bytes.Buffer{}); !errors.Is(err, services.ErrBerkasTidakDiStub) {
		t.Errorf("tidak satu pun tersedia = galat berkalimat: %v", err)
	}
	_ = repository.ErrTidakAda
}

// #10 - batas autocomplete.
func TestCariMasterBerbatas(t *testing.T) {
	l, _ := layananMaster()
	d, err := l.CariMaster(context.Background(), pelakuUji, models.MasterCeding, "uji", 1)
	if err != nil || len(d) != 1 {
		t.Errorf("batas 1: %v %v", d, err)
	}
	if d, _ := l.CariMaster(context.Background(), pelakuUji, models.MasterCeding, "uji", 0); len(d) != 2 {
		t.Errorf("batas 0 = seluruh hasil RD: %v", d)
	}
	if _, err := l.CariMaster(context.Background(), pelakuUji, models.MasterCeding, "", -1); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("batas negatif: %v", err)
	}
}
