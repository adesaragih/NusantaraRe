package services_test

// Uji `TreatyInEDMSetValue` [8] `TreatyRevisionCopyAttachment` — lampiran
// master disalin (diunggah ulang) ke pengenal penyesuaian baru pada tulisan
// PERTAMA draf, di atas gudang dan penyimpanan tiruan.

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func gudangSalinLampiran() *gudangTiruan {
	g := gudangEDM()
	// Baris sumber APA ADANYA — termasuk ekor CR LF nama kategori dan
	// FILEMIMETYPE huruf besar: keduanya disalin verbatim ([2.8]).
	g.lampiran = []models.BarisLampiranWarisan{
		{ID: "L1", KodeKategori: "00004", NamaKategori: "Slip\r\n", NamaBerkas: "Slip 2025.pdf", JenisMime: "PDF", IDSimpanan: "IMG1"},
		{ID: "L2", KodeKategori: "00002", NamaKategori: "Approval Email", NamaBerkas: "Bordero.xlsx", JenisMime: "", IDSimpanan: "IMG2"},
	}
	g.objekSimpanan = map[string]models.ObjekSimpanan{
		"IMG1": {ImageID: "IMG1", URLPublik: "https://storage.googleapis.com/rnmtest/slip?sig=1", Exp: "01/01/2099 00:00:00",
			App: "rnmtest", AppFolder: "gs://rnmtest/Contract/Doc/2025/08/x - Slip 2025.pdf", NamaObjek: "x - Slip 2025.pdf"},
		"IMG2": {ImageID: "IMG2", URLPublik: "https://storage.googleapis.com/rnmtest/bordero?sig=1", Exp: "01/01/2099 00:00:00",
			App: "rnmtest", AppFolder: "gs://rnmtest/Contract/Doc/2025/08/y - Bordero.xlsx", NamaObjek: "y - Bordero.xlsx"},
	}
	return g
}

func masukanDraf(oldID string) services.MasukanPenyesuaian {
	return services.MasukanPenyesuaian{
		ID: "1001001/R02", Draf: true,
		Baru: services.SisiKiriman{Medan: map[string]any{"OLDID": oldID, "EDMState": "2"}},
		Lama: &services.SisiKiriman{Medan: map[string]any{"RNMShare": "4"}},
	}
}

func TestSaveDrafMenyalinLampiranMasterSepertiPega(t *testing.T) {
	g, s := gudangSalinLampiran(), &simpananTiruan{}
	h, err := services.LayananDenganSimpanan(g, s).SimpanPenyesuaian(context.Background(), admin, masukanDraf("1001001/R01"))
	if err != nil {
		t.Fatal(err)
	}
	if h.SalinanLampiran == nil || h.SalinanLampiran.Tersalin != 2 || h.SalinanLampiran.Sumber != "1001001/R01" || len(h.SalinanLampiran.Berkas) != 2 {
		t.Fatalf("salinan %+v", h.SalinanLampiran)
	}
	// [2.3]/[2.5] isi diambil dari URL bertanda tangan objek SUMBER.
	if len(s.diambil) != 2 || s.diambil[0] != "https://storage.googleapis.com/rnmtest/slip?sig=1" {
		t.Errorf("diambil %v", s.diambil)
	}
	// [2.7] InsertGoogleStorage_Act: Folder Contract, Namafile nama asli,
	// Image = Base64 isi, Ext = FILEMIMETYPE huruf kecil (kosong → dari nama).
	isi := base64.StdEncoding.EncodeToString([]byte("isi-berkas"))
	if len(s.diminta) != 2 {
		t.Fatalf("unggah %d", len(s.diminta))
	}
	p1, p2 := s.diminta[0], s.diminta[1]
	if !strings.HasPrefix(p1.Folder, "Contract/Doc/") || p1.Durasi != 1800 || p1.Image != isi || p1.Ext != "pdf" ||
		p1.MimeType != "application/pdf" || !strings.HasSuffix(p1.Namafile, " - Slip 2025.pdf") {
		t.Errorf("unggah 1 %+v", p1)
	}
	if p2.Ext != "xlsx" || !strings.HasSuffix(p2.Namafile, " - Bordero.xlsx") {
		t.Errorf("unggah 2 %+v", p2)
	}
	// [2.8] InsertAttachment2_Sql: TREATYID BARU, kolom sumber verbatim,
	// USERNAME operator, T_STORAGE_ID = objek BARU (bukan IMG1).
	if len(g.lampiranBaru) != 2 {
		t.Fatalf("tercatat %d", len(g.lampiranBaru))
	}
	l := g.lampiranBaru[0]
	if l.IDKontrak != "1001001/R02" || l.KodeKategori != "00004" || l.NamaKategori != "Slip\r\n" || l.NamaBerkas != "Slip 2025.pdf" ||
		l.Ekstensi != "PDF" || l.Pengguna != "ADESAMUEL" || len(l.ImageID) != 32 || l.ImageID == "IMG1" || l.NamaObjek != p1.Namafile {
		t.Errorf("lampiran %+v", l)
	}
	if g.lampiranBaru[1].Ekstensi != "" || g.lampiranBaru[1].KodeKategori != "00002" {
		t.Errorf("lampiran 2 %+v", g.lampiranBaru[1])
	}
	// Kepala tersimpan LEBIH DULU, dan pesannya tetap pesan prosedur.
	if len(g.disimpanEDM) != 1 || h.Pesan != "Data Sudah Disimpan Dengan ID : 1001001/R02" {
		t.Errorf("simpan %d pesan %q", len(g.disimpanEDM), h.Pesan)
	}
}

// Draf yang dikirim lewat Submit pun tulisan PERTAMA — Pega menyalin saat
// `Choose`, apa pun tombol sesudahnya.
func TestSubmitDrafJugaMenyalinLampiran(t *testing.T) {
	g, s := gudangSalinLampiran(), &simpananTiruan{}
	m := masukanDraf("1001001/R01")
	h, err := services.LayananDenganSimpanan(g, s).KirimPenyesuaian(context.Background(), admin, services.MasukanKirimPenyesuaian{
		MasukanPenyesuaian: m, Aksi: services.AksiSubmit,
	})
	if err != nil {
		t.Fatal(err)
	}
	if h.SalinanLampiran == nil || h.SalinanLampiran.Tersalin != 2 || len(g.lampiranBaru) != 2 {
		t.Errorf("salinan %+v tercatat %d", h.SalinanLampiran, len(g.lampiranBaru))
	}
}

// Penyesuaian TERSIMPAN tidak menyalin lagi — salinan tepat sekali.
func TestSavePenyesuaianTersimpanTidakMenyalinLampiran(t *testing.T) {
	g, s := gudangSalinLampiran(), &simpananTiruan{}
	h, err := services.LayananDenganSimpanan(g, s).SimpanPenyesuaian(context.Background(), admin, services.MasukanPenyesuaian{ID: "1001001/R01"})
	if err != nil {
		t.Fatal(err)
	}
	if h.SalinanLampiran != nil || len(s.diminta) != 0 || len(g.lampiranBaru) != 0 {
		t.Errorf("salinan %+v unggah %d tercatat %d", h.SalinanLampiran, len(s.diminta), len(g.lampiranBaru))
	}
}

// Draf yang pengenalnya sudah ada DITOLAK sebelum apa pun disalin.
func TestDrafDitolakNolSalinan(t *testing.T) {
	g, s := gudangSalinLampiran(), &simpananTiruan{}
	m := masukanDraf("1001001")
	m.ID = "1001001/R01"
	if _, err := services.LayananDenganSimpanan(g, s).SimpanPenyesuaian(context.Background(), admin, m); err == nil {
		t.Fatal("draf berpengenal tersimpan diterima")
	}
	if len(s.diambil) != 0 || len(s.diminta) != 0 || len(g.lampiranBaru) != 0 {
		t.Errorf("ada yang tersalin: ambil %d unggah %d catat %d", len(s.diambil), len(s.diminta), len(g.lampiranBaru))
	}
}

// Berkas yang gagal DILAPORKAN; berkas lain tetap disalin; Save tetap
// berhasil (kepala sudah mengikat).
func TestSalinLampiranBerkasGagalDilaporkanLainTetapJalan(t *testing.T) {
	g, s := gudangSalinLampiran(), &simpananTiruan{}
	delete(g.objekSimpanan, "IMG1")
	h, err := services.LayananDenganSimpanan(g, s).SimpanPenyesuaian(context.Background(), admin, masukanDraf("1001001/R01"))
	if err != nil {
		t.Fatal(err)
	}
	b := h.SalinanLampiran.Berkas
	if h.SalinanLampiran.Tersalin != 1 || b[0].Berhasil || b[0].Pesan == "" || !b[1].Berhasil || len(g.lampiranBaru) != 1 {
		t.Errorf("salinan %+v tercatat %d", h.SalinanLampiran, len(g.lampiranBaru))
	}
}

// Penyimpanan belum disambung → seluruh berkas gagal, Save tetap berhasil.
func TestSalinLampiranTanpaPenyimpananTidakMenggagalkanSave(t *testing.T) {
	g := gudangSalinLampiran()
	h, err := services.LayananDengan(g).SimpanPenyesuaian(context.Background(), admin, masukanDraf("1001001/R01"))
	if err != nil {
		t.Fatal(err)
	}
	if len(g.disimpanEDM) != 1 || h.SalinanLampiran == nil || h.SalinanLampiran.Tersalin != 0 || len(h.SalinanLampiran.Berkas) != 2 ||
		h.SalinanLampiran.Berkas[1].Berhasil {
		t.Errorf("salinan %+v", h.SalinanLampiran)
	}
}

// Pagar OLDID: asal yang tidak sejalan dengan pengenal revisi tidak disalin.
func TestSalinLampiranOLDIDTakSejalanDitolak(t *testing.T) {
	g, s := gudangSalinLampiran(), &simpananTiruan{}
	h, err := services.LayananDenganSimpanan(g, s).SimpanPenyesuaian(context.Background(), admin, masukanDraf("2002002"))
	if err != nil {
		t.Fatal(err)
	}
	if h.SalinanLampiran == nil || h.SalinanLampiran.Pesan == "" || len(s.diambil) != 0 || len(g.lampiranBaru) != 0 {
		t.Errorf("salinan %+v", h.SalinanLampiran)
	}
}

// Master tanpa lampiran, atau OLDID kosong → nol salinan, nol laporan.
func TestSalinLampiranTanpaSumberNil(t *testing.T) {
	g, s := gudangSalinLampiran(), &simpananTiruan{}
	g.lampiran = nil
	h, err := services.LayananDenganSimpanan(g, s).SimpanPenyesuaian(context.Background(), admin, masukanDraf("1001001/R01"))
	if err != nil || h.SalinanLampiran != nil {
		t.Errorf("galat %v salinan %+v", err, h.SalinanLampiran)
	}
	if got := services.LayananDenganSimpanan(gudangSalinLampiran(), s).SalinLampiranRevisi(context.Background(), admin, "1001001/R02", ""); got != nil {
		t.Errorf("OLDID kosong: %+v", got)
	}
}

func TestOLDIDSejalan(t *testing.T) {
	for _, c := range []struct {
		baru, lama string
		ok         bool
	}{
		{"1001001/R01", "1001001", true},
		{"1001001/R03", "1001001/R02", true},
		{"1001001/R01", "1001001/R01", false},
		{"1001001/R01", "2002002", false},
		{"1001001/R01", "", false},
		{"1001001", "1001001/R01", false},
	} {
		if got := services.OLDIDSejalan(c.baru, c.lama); got != c.ok {
			t.Errorf("OLDIDSejalan(%q, %q) = %v", c.baru, c.lama, got)
		}
	}
}
