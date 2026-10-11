package services_test

// Uji tombol `Revision` (`SetTreatyIn_Act` ber-revisionstate) dan tangga
// revisinya — di atas gudang tiruan.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func gudangTuntas() *gudangTiruan {
	g := gudangSimpan()
	g.dokumenTersimpan["1001001"]["StatusAkseptasi"] = models.StatusTuntas
	g.dokumenTersimpan["1001001"]["CommentList"] = []any{map[string]any{"OperatorName": "ADESAMUEL", "IsApproved": "Accept"}}
	return g
}

func TestRevisionMenyetelKeadaanRevisiLaluMenyimpan(t *testing.T) {
	g := gudangTuntas()
	h, err := services.LayananDengan(g).MulaiRevisi(context.Background(), admin, services.MasukanRevisi{IDKontrak: "1001001"})
	if err != nil {
		t.Fatal(err)
	}
	d := g.disimpan[0].Dokumen
	// [6]–[7]: Position TIDAK diubah.
	if d["RevisionState"] != "1" || d["ViewState"] != "1" || d["StatusAkseptasi"] != "" ||
		d["PositionUsername"] != "ADESAMUEL" || d["Position"] != "" {
		t.Errorf("keadaan revisi %v %v %q %v %q", d["RevisionState"], d["ViewState"], d["StatusAkseptasi"], d["PositionUsername"], d["Position"])
	}
	// [9]: "Create Revision", tanpa IsApproved.
	k, _ := d["CommentList"].([]any)
	el, _ := k[len(k)-1].(map[string]any)
	if len(k) != 2 || el["Suggest"] != services.SuggestBuatRevisi || el["OperatorName"] != "ADESAMUEL" || el["Date"] == "" {
		t.Errorf("komentar %v", k)
	}
	if _, ada := el["IsApproved"]; ada {
		t.Errorf("IsApproved tidak diisi SetTreatyIn_Act [9]: %v", el)
	}
	if h.ID != "1001001" || h.Status != "" {
		t.Errorf("hasil %+v", h)
	}
}

func TestRevisionHanyaUntukAdminDanKontrakTuntas(t *testing.T) {
	ctx := context.Background()
	if _, err := services.LayananDengan(gudangTuntas()).MulaiRevisi(ctx, secHead, services.MasukanRevisi{IDKontrak: "1001001"}); !errors.Is(err, services.ErrBukanPemegangPosisi) {
		t.Errorf("Sec Head: %v", err)
	}
	// Belum tuntas.
	if _, err := services.LayananDengan(gudangSimpan()).MulaiRevisi(ctx, admin, services.MasukanRevisi{IDKontrak: "1001001"}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("belum tuntas: %v", err)
	}
	// Tuntas tetapi masih menunggu di sebuah posisi.
	g := gudangTuntas()
	g.dokumenTersimpan["1001001"]["Position"] = models.PosisiSecHead
	if _, err := services.LayananDengan(g).MulaiRevisi(ctx, admin, services.MasukanRevisi{IDKontrak: "1001001"}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("berposisi: %v", err)
	}
}

// ⛔ Tanpa kolom 448 keadaan revisi hilang saat disimpan — DITOLAK, nol tulisan.
func TestRevisionDitolakSebelumMigrasi448(t *testing.T) {
	g := gudangTuntas()
	g.belumTerpasang = []string{"RevisionState", "ViewState"}
	_, err := services.LayananDengan(g).MulaiRevisi(context.Background(), admin, services.MasukanRevisi{IDKontrak: "1001001"})
	if !errors.Is(err, services.ErrTombolDitolak) || len(g.disimpan) != 0 {
		t.Errorf("galat %v, disimpan %d", err, len(g.disimpan))
	}
}

// Layar tidak dapat menyetel keadaan revisi sendiri.
func TestSaveMengabaikanRevisionStateKirimanLayar(t *testing.T) {
	g := gudangSimpan()
	_, err := services.LayananDengan(g).SimpanKontrak(context.Background(), admin, services.MasukanSimpan{
		IDKontrak: "1001001", Dokumen: map[string]any{"RevisionState": "1", "ViewState": "1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	d := g.disimpan[0].Dokumen
	if _, ada := d["RevisionState"]; ada {
		t.Errorf("RevisionState dari layar diterima: %v", d["RevisionState"])
	}
}

// ⭐ Submit revisi: cabang revisi `Akseptasi_DT` — Admin → Sec Head, lalu
// Sec Head Accept MENUNTASKAN dan mengosongkan keadaan revisi.
func TestSubmitRevisiTuntasDiSecHead(t *testing.T) {
	g := gudangSimpan()
	g.dokumenTersimpan["1001001"]["RevisionState"] = "1"
	g.dokumenTersimpan["1001001"]["ViewState"] = "1"
	l := services.LayananDengan(g)
	if _, err := l.KirimKontrak(context.Background(), admin, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"}, Aksi: services.AksiSubmit,
	}); err != nil {
		t.Fatal(err)
	}
	d := g.disimpan[0].Dokumen
	if d["Position"] != models.PosisiSecHead || d["StatusAkseptasi"] != "Accept" || d["RevisionState"] != "1" {
		t.Errorf("naik %v %v %v", d["Position"], d["StatusAkseptasi"], d["RevisionState"])
	}
	// Sec Head menerima — tangga revisi berhenti di sini.
	g.dokumenTersimpan["1001001"] = d
	if _, err := l.KirimKontrak(context.Background(), secHead, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"}, Aksi: services.AksiAkseptasi, Pilihan: "Accept",
	}); err != nil {
		t.Fatal(err)
	}
	d = g.disimpan[1].Dokumen
	if d["Position"] != "" || d["StatusAkseptasi"] != models.StatusTuntas || d["RevisionState"] != "" || d["ViewState"] != "" {
		t.Errorf("tuntas %q %v %q %q", d["Position"], d["StatusAkseptasi"], d["RevisionState"], d["ViewState"])
	}
}
