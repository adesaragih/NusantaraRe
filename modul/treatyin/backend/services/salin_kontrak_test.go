package services_test

// Uji tombol `Copy` daftar kontrak — `Activity/TreatyInCopy.xml` di atas
// gudang tiruan. Rantainya di `salin_kontrak.go`.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// gudangSalin - kontrak sumber `1001001` yang Resolve Complete, Position
// kosong (syarat tampil cell 993), dengan riwayat, revisi, dan tab tersimpan.
func gudangSalin() *gudangTiruan {
	return &gudangTiruan{
		kepalaTreatyIn: map[string]map[string]any{
			"1001001": {
				"TreatyContractName": "SUMBER", "Ceding": "ASURANSI A", "CedingID": "C1",
				"LeadingReinsSource": "DIRECT", "LeadingReinsSourceID": "S1",
				"Position": "", "StatusAkseptasi": models.StatusTuntas, "PositionUsername": "LAMA",
			},
		},
		dokumenTersimpan: map[string]map[string]any{
			"1001001": {
				"Portfolio":     []any{map[string]any{"Description": "tersimpan"}},
				"CommentList":   []any{map[string]any{"Suggest": "Create Revision"}, map[string]any{"Suggest": "Accept"}},
				"RevisionState": "1", "ViewState": "1", "OLDID": "",
			},
		},
		pemegangPosisi: map[string][]string{models.PosisiSecHead: {"SEC1"}},
	}
}

func TestTerapkanSalinPersisTreatyInCopy(t *testing.T) {
	doc := map[string]any{
		"ID": "1001001", "CommentList": []any{map[string]any{"Suggest": "lama"}},
		"RevisionState": "1", "Position": "", "StatusAkseptasi": models.StatusTuntas, "ViewState": "1",
		"TreatyContractName": "SUMBER",
	}
	services.TerapkanSalin(doc, "1001001", "ADESAMUEL", "20261008T030000.000 GMT")
	// [1]–[3]: riwayat lama dibuang, SATU komentar salinan.
	l, _ := doc["CommentList"].([]any)
	if len(l) != 1 {
		t.Fatalf("CommentList %v, mau satu baris", doc["CommentList"])
	}
	k := l[0].(map[string]any)
	if k["Suggest"] != "Copied from ID 1001001" || k["OperatorName"] != "ADESAMUEL" || k["Date"] != "20261008T030000.000 GMT" {
		t.Errorf("komentar salinan %v", k)
	}
	if _, ada := k["IsApproved"]; ada {
		t.Error("[3] tidak menyetel IsApproved")
	}
	// [3] OLDID = ID lama; [4] ID baru (kosong → `idKontrakBaru`).
	if doc["OLDID"] != "1001001" {
		t.Errorf("OLDID %v", doc["OLDID"])
	}
	if _, ada := doc["ID"]; ada {
		t.Error("ID sumber tidak boleh terbawa — UnknownId = pengenal baru")
	}
	if doc["Position"] != models.PosisiAdmin || doc["StatusAkseptasi"] != "" || doc["RevisionState"] != "" || doc["ViewState"] != "0" {
		t.Errorf("keadaan salinan %v", doc)
	}
	// Isi kontrak lain tetap — clipboard disalin utuh.
	if doc["TreatyContractName"] != "SUMBER" {
		t.Errorf("isi hilang: %v", doc)
	}
}

func TestSaveSalinanMelahirkanKontrakBaruDariClipboardUtuh(t *testing.T) {
	g := gudangSalin()
	h, err := services.LayananDengan(g).SimpanSalinan(context.Background(), admin, services.MasukanSalin{
		IDSumber: "1001001",
		Dokumen: map[string]any{
			"TreatyContractName": "SALINAN",
			// ⛔ Milik server / tombol Copy — diabaikan.
			"StatusAkseptasi": "Resolve Complete", "CommentList": []any{}, "OLDID": "999",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.disimpan) != 1 {
		t.Fatalf("%d tulisan, mau 1", len(g.disimpan))
	}
	r := g.disimpan[0]
	if r.ID != "" {
		t.Errorf("ID rencana %q — salinan harus berpengenal BARU", r.ID)
	}
	d := r.Dokumen
	if d["TreatyContractName"] != "SALINAN" || d["Ceding"] != "ASURANSI A" {
		t.Errorf("kepala %v", d)
	}
	// Tab yang tidak dibuka layar tetap tersalin dari sumber.
	if p, _ := d["Portfolio"].([]any); len(p) != 1 {
		t.Errorf("Portfolio sumber hilang: %v", d["Portfolio"])
	}
	if d["OLDID"] != "1001001" || d["Position"] != models.PosisiAdmin || d["StatusAkseptasi"] != "" {
		t.Errorf("OLDID %v Position %v Status %v", d["OLDID"], d["Position"], d["StatusAkseptasi"])
	}
	if l, _ := d["CommentList"].([]any); len(l) != 1 || !strings.HasPrefix(l[0].(map[string]any)["Suggest"].(string), "Copied from ID 1001001") {
		t.Errorf("CommentList %v", d["CommentList"])
	}
	if h.ID != "1001857" || h.Pesan != "Data Sudah Disimpan Dengan ID : 1001857" {
		t.Errorf("hasil %+v", h)
	}
	// ⛔ Sumbernya tidak disentuh.
	if g.kepalaTreatyIn["1001001"]["StatusAkseptasi"] != models.StatusTuntas {
		t.Error("kontrak sumber berubah")
	}
}

// Syarat tampil cell 993 ditegakkan di server — dan penolakan tidak menulis.
func TestSalinTundukSyaratTampilCell993(t *testing.T) {
	ctx := context.Background()
	g := gudangSalin()
	l := services.LayananDengan(g)
	if _, err := l.SimpanSalinan(ctx, secHead, services.MasukanSalin{IDSumber: "1001001"}); !errors.Is(err, services.ErrBukanPemegangPosisi) {
		t.Errorf("bukan Admin: %v", err)
	}
	if _, err := l.BacaDraftSalinan(ctx, secHead, "1001001"); !errors.Is(err, services.ErrBukanPemegangPosisi) {
		t.Errorf("draf bukan Admin: %v", err)
	}
	g.kepalaTreatyIn["1001001"]["StatusAkseptasi"] = ""
	if _, err := l.SimpanSalinan(ctx, admin, services.MasukanSalin{IDSumber: "1001001"}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("sumber belum tuntas: %v", err)
	}
	g.kepalaTreatyIn["1001001"]["StatusAkseptasi"] = models.StatusTuntas
	g.kepalaTreatyIn["1001001"]["Position"] = models.PosisiSecHead
	if _, err := l.BacaDraftSalinan(ctx, admin, "1001001"); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("sumber masih di posisi: %v", err)
	}
	if _, err := l.SimpanSalinan(ctx, admin, services.MasukanSalin{}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("tanpa sumber: %v", err)
	}
	if len(g.disimpan) != 0 {
		t.Errorf("penolakan menulis %d kali", len(g.disimpan))
	}
}

// Submit draf: `TreatyInCheckError` diperiksa SEBELUM salinan ditulis.
func TestSubmitSalinanDitolakTanpaMenulis(t *testing.T) {
	g := gudangSalin()
	_, err := services.LayananDengan(g).SimpanSalinan(context.Background(), admin, services.MasukanSalin{
		IDSumber: "1001001", Aksi: services.AksiSubmit, Dokumen: map[string]any{"LeadingReinsSource": ""},
	})
	if !errors.Is(err, services.ErrTombolDitolak) || !strings.Contains(err.Error(), "Please input Source of Business (SoB)") {
		t.Errorf("galat %v", err)
	}
	if len(g.disimpan) != 0 {
		t.Errorf("Submit yang ditolak menulis %d kali", len(g.disimpan))
	}
	if _, err := services.LayananDengan(g).SimpanSalinan(context.Background(), admin, services.MasukanSalin{
		IDSumber: "1001001", Aksi: services.AksiAkseptasi,
	}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("akseptasi atas draf: %v", err)
	}
}

// Submit draf: salinan lahir, lalu tangga dijalankan atas pengenal barunya.
func TestSubmitSalinanMenyimpanLaluMenaikkanTangga(t *testing.T) {
	g := gudangSalin()
	// Tiruan `SimpanKontrak` tidak menyimpan ke peta bacanya; kontrak baru
	// disemai supaya langkah kedua membaca salinan yang "tersimpan".
	g.kepalaTreatyIn["1001857"] = map[string]any{
		"Ceding": "ASURANSI A", "LeadingReinsSource": "DIRECT", "Position": models.PosisiAdmin, "StatusAkseptasi": "",
	}
	h, err := services.LayananDengan(g).SimpanSalinan(context.Background(), admin, services.MasukanSalin{
		IDSumber: "1001001", Aksi: services.AksiSubmit,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.disimpan) != 2 || g.disimpan[0].ID != "" || g.disimpan[1].ID != "1001857" {
		t.Fatalf("tulisan %+v", g.disimpan)
	}
	if h.Posisi != models.PosisiSecHead {
		t.Errorf("posisi sesudah Submit %q", h.Posisi)
	}
}

// Draf Copy: kontrak sumber tampil sebagai kontrak BARU — nol tulisan.
func TestDraftSalinanTampilSebagaiKontrakBaru(t *testing.T) {
	g := gudangSalin()
	g.kontrakWarisan = models.KontrakWarisan{ID: "1001001", NamaKontrak: "SUMBER", StatusAkseptasi: models.StatusTuntas}
	g.catatan = []models.BarisCatatanWarisan{{Catatan: "Accept"}, {Catatan: "Create Revision"}}
	g.lampiran = []models.BarisLampiranWarisan{{ID: "9", NamaKategori: "Slip"}}
	g.polisProduksi = []models.BarisPolisProduksi{{NomorPolis: "RNM-1"}}
	k, err := services.LayananDengan(g).BacaDraftSalinan(context.Background(), admin, "1001001")
	if err != nil {
		t.Fatal(err)
	}
	if k.ID != "" || k.NamaKontrak != "SUMBER" || k.Posisi != models.PosisiAdmin || k.StatusAkseptasi != "" {
		t.Errorf("kepala draf %+v", k)
	}
	if k.Penampung["RevisionState"] != "" || k.Penampung["ViewState"] != "0" {
		t.Errorf("penampung %v", k.Penampung)
	}
	if len(k.Catatan) != 1 || k.Catatan[0].Catatan != "Copied from ID 1001001" || k.Catatan[0].Operator != admin.AkunID {
		t.Errorf("riwayat draf %+v", k.Catatan)
	}
	if len(k.Lampiran) != 0 || len(k.PolisProduksi) != 0 {
		t.Errorf("lampiran %v polis %v — berkunci ID, UnknownId memberi nol", k.Lampiran, k.PolisProduksi)
	}
	for _, kat := range k.KategoriLampiran {
		if kat.Cacah != 0 {
			t.Errorf("kategori %+v", kat)
		}
	}
	if len(g.disimpan) != 0 {
		t.Error("draf Copy menulis")
	}
}
