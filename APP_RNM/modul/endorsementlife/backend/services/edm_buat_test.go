package services_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/jejak"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

// jejakTiruan mencatat setiap transisi yang direkam.
type jejakTiruan struct {
	catatan []jejak.CatatanJejak
	galat   error
}

func (j *jejakTiruan) Rekam(_ context.Context, _ *db.Tx, c jejak.CatatanJejak) error {
	if j.galat != nil {
		return j.galat
	}
	j.catatan = append(j.catatan, c)
	return nil
}

var jamUji = time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)

func layananJejak(g *tiruan.Gudang) (*services.Layanan, *jejakTiruan) {
	j := &jejakTiruan{}
	return services.BaruLayanan(g, tiruan.Transaksi, nil).DenganJejak(j).DenganJam(func() time.Time { return jamUji }), j
}

// gudangPolisNB - polis new business sistem baru `UJI-PL-NB` (versi 1): tiga
// peserta, satu di antaranya ber-spreading + retro; dan satu polis warisan.
func gudangPolisNB() *tiruan.Gudang {
	g := tiruan.Baru()
	g.Polis["NBLF-7"] = &tiruan.Polis{ID: "NBLF-7", ProdKe: 1, Kepala: map[string]string{
		"TYPE": "QR", "TYPE_CEDING": "1", "BUSINESS_NAME": "UJI-COB", "CEDING_CO_NAME": "UJI-CEDING", "PRODUCT_NAME": "UJI-PRODUK"}}
	for i, c := range []string{"UJI-C1", "UJI-C2", "UJI-C3"} {
		g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: "UJI-NB-D" + string(rune('1'+i)), PolisID: "NBLF-7", PLNumber: "UJI-PL-NB",
			Nilai: map[string]string{"CERTIFICATE_NO": c, "SUM_INSURED": "1000.5", "PLAN": "UJI-PLAN"}})
	}
	g.Spreading = append(g.Spreading, &tiruan.Spreading{ID: "UJI-S1", DetailID: "UJI-NB-D1", Nama: "UJI-QS", Share: "50",
		Retro: []models.SpreadingRetro{{ReinsurerName: "UJI-REAS", PercentShare: "100"}}})
	g.PolisWarisan = append(g.PolisWarisan,
		&tiruan.PolisWarisan{IDPega: "UJI-IDPEGA-NB", NoPolis: "UJI-PL-W", ProdKe: 0,
			Kepala: map[string]string{"Type": "TR", "DateReceived": "20240115", "WPC": "bukan-tanggal"}})
	g.PesertaWarisan = append(g.PesertaWarisan,
		&tiruan.PesertaWarisan{ID: "UJI-W1", PLNumber: "UJI-PL-W", IDPega: "UJI-IDPEGA-NB", Nilai: map[string]string{"CERTIFICATE_NO": "UJI-WC1"}},
		&tiruan.PesertaWarisan{ID: "UJI-W2", PLNumber: "UJI-PL-W", IDPega: "UJI-IDPEGA-LAIN", Nilai: map[string]string{"CERTIFICATE_NO": "UJI-WC2"}})
	return g
}

func TestKelayakanMengumpulkanPesanVERBATIM(t *testing.T) {
	g := gudangPolisNB()
	l, _ := layananJejak(g)
	ctx := context.Background()
	k, err := l.Kelayakan(ctx, pelakuUji, services.MasukanKelayakan{NomorPolis: " ", EdmType: "1"})
	if err != nil || !reflect.DeepEqual(k.Pesan, []string{models.PesanPolisTidakSah, models.PesanPolisTidakAda}) || k.Boleh {
		t.Fatalf("polis kosong: %+v %v (Pega: 3.1 lalu 3.3)", k, err)
	}
	k, _ = l.Kelayakan(ctx, pelakuUji, services.MasukanKelayakan{NomorPolis: "UJI-TIDAK-ADA", EdmType: "1"})
	if !reflect.DeepEqual(k.Pesan, []string{models.PesanPolisTidakAda}) {
		t.Errorf("polis tak dikenal: %v", k.Pesan)
	}
	k, _ = l.Kelayakan(ctx, pelakuUji, services.MasukanKelayakan{NomorPolis: "UJI-PL-NB", EdmType: "1"})
	if len(k.Pesan) != 0 || !k.Boleh {
		t.Errorf("polis NB sah: %+v", k)
	}
	if k, _ = l.Kelayakan(ctx, pelakuUji, services.MasukanKelayakan{NomorPolis: "UJI-PL-NB", EdmType: "2"}); k.Boleh {
		t.Error("EdmType 2 tidak boleh menyalakan Submit (b4226)")
	}
	// Gerbang 5: dibayar hanya menolak Batal; Arasapas tidak dibaca untuk Perubahan Data.
	g.Dibayar["UJI-PL-NB"] = true
	if k, _ = l.Kelayakan(ctx, pelakuUji, services.MasukanKelayakan{NomorPolis: "UJI-PL-NB", EdmType: "1"}); len(k.Pesan) != 0 {
		t.Errorf("Perubahan Data atas polis dibayar ditolak: %v (AC 5)", k.Pesan)
	}
	g.Dibayar = map[string]bool{"UJI-PLNB": true}
	if k, _ = l.Kelayakan(ctx, pelakuUji, services.MasukanKelayakan{NomorPolis: "UJI-PL-NB", EdmType: "3"}); len(k.Pesan) != 0 {
		t.Errorf("nomor invoice bukan nomor polis tanpa titik: %v", k.Pesan)
	}
	g.Dibayar = map[string]bool{"UJI-PL-NB": true} // nomor uji tanpa titik: invoice = nomor polis
	if k, _ = l.Kelayakan(ctx, pelakuUji, services.MasukanKelayakan{NomorPolis: "UJI-PL-NB", EdmType: "3"}); !reflect.DeepEqual(k.Pesan, []string{models.PesanSudahDibayar}) {
		t.Errorf("Batal atas polis dibayar: %v", k.Pesan)
	}
	g.GalatArasapas = repository.ErrArasapasTakTerbaca
	if _, err := l.Kelayakan(ctx, pelakuUji, services.MasukanKelayakan{NomorPolis: "UJI-PL-NB", EdmType: "3"}); !errors.Is(err, services.ErrArasapas) {
		t.Errorf("Arasapas tak terbaca: %v (gagal tertutup, OQ-EDM-012)", err)
	}
	if _, err := l.Kelayakan(ctx, pelakuUji, services.MasukanKelayakan{NomorPolis: "UJI-PL-NB", EdmType: "1"}); err != nil {
		t.Errorf("Perubahan Data tidak boleh menyentuh Arasapas: %v", err)
	}
}

func TestBuatKasusMenyalinVersiNBBesertaSpreading(t *testing.T) {
	g := gudangPolisNB()
	g.Peserta[2].EdmStatus = models.StatusDelete // baris Delete versi lama: tidak tersalin (R10)
	l, j := layananJejak(g)
	h, err := l.BuatKasus(context.Background(), pelakuUji, services.MasukanKasus{
		NomorPolis: "UJI-PL-NB", EdmType: "1", EdmDate: "2026-10-01", Deskripsi: "UJI-CATATAN"})
	if err != nil {
		t.Fatal(err)
	}
	if h.ID != "EDMLF-1" || h.Peserta != 2 || h.Spreading != 1 || h.SpreadingRetro != 1 {
		t.Fatalf("hasil = %+v", h)
	}
	k := g.Polis["EDMLF-1"]
	if k.OldPolicyNo != "UJI-PL-NB" || k.ProdKe != 2 || k.EdmNote != "UJI-CATATAN" || k.Pembuat != pelakuUji.AkunID ||
		k.Kepala["PRODUCT_NAME"] != "UJI-PRODUK" || k.Kepala["TYPE"] != "QR" || k.NoPolis != "" || k.Status != "" {
		t.Errorf("kepala kasus = %+v", k)
	}
	for _, d := range g.Peserta {
		if d.PolisID != "EDMLF-1" {
			continue
		}
		if d.EdmStatus != models.StatusOld || d.ParentID == "" || d.Nilai["SUM_INSURED"] != "1000.5" {
			t.Errorf("peserta salinan = %+v", d)
		}
	}
	if len(j.catatan) != 1 || j.catatan[0].KlaimID != "EDMLF-1" || j.catatan[0].Ke != models.TahapInputEDMLife ||
		j.catatan[0].AkunID != pelakuUji.AkunID || !j.catatan[0].Waktu.Equal(jamUji) {
		t.Errorf("jejak = %+v", j.catatan)
	}
	// Gerbang 3 kini menolak kasus kedua atas polis yang sama.
	_, err = l.BuatKasus(context.Background(), pelakuUji, services.MasukanKasus{NomorPolis: "UJI-PL-NB", EdmType: "1"})
	var gk services.GalatKelayakan
	if !errors.As(err, &gk) || !reflect.DeepEqual(gk.Pesan, []string{models.PesanEDMBelumSelesai}) {
		t.Errorf("kasus kedua: %v", err)
	}
}

func TestBuatKasusDariVersiNBWarisan(t *testing.T) {
	g := gudangPolisNB()
	l, _ := layananJejak(g)
	h, err := l.BuatKasus(context.Background(), pelakuUji, services.MasukanKasus{NomorPolis: "UJI-PL-W", EdmType: "3"})
	if err != nil {
		t.Fatal(err)
	}
	if h.Peserta != 1 || h.Spreading != 0 {
		t.Fatalf("salinan warisan = %+v (IDPEGA versi itu saja)", h)
	}
	k := g.Polis[h.ID]
	if k.ProdKe != 2 || k.Kepala["TYPE"] != "TR" || k.Kepala["DATE_RECEIVED"] != "2024-01-15" || k.Kepala["WPC"] != "" {
		t.Errorf("kepala warisan = %+v (PRODKE kosong = versi 1; tanggal rusak dikosongkan)", k)
	}
	for _, d := range g.Peserta {
		if d.PolisID == h.ID && d.ParentID != "" {
			t.Errorf("peserta warisan ber-PARENT_ID %q (OQ-EDM-007)", d.ParentID)
		}
	}
}

func TestBuatKasusDitolak(t *testing.T) {
	ctx := context.Background()
	g := gudangPolisNB()
	l, j := layananJejak(g)
	for _, m := range []services.MasukanKasus{
		{NomorPolis: "UJI-PL-NB", EdmType: "2"},
		{NomorPolis: "UJI-PL-NB", EdmType: "1", EdmDate: "01/10/2026"},
	} {
		if _, err := l.BuatKasus(ctx, pelakuUji, m); !errors.Is(err, services.ErrMasukanTidakSah) {
			t.Errorf("%+v: %v", m, err)
		}
	}
	var gk services.GalatKelayakan
	if _, err := l.BuatKasus(ctx, pelakuUji, services.MasukanKasus{NomorPolis: "UJI-TIDAK-ADA", EdmType: "1"}); !errors.As(err, &gk) {
		t.Errorf("polis tak dikenal: %v", err)
	}
	// Versi berjalan polis warisan adalah endorsement sistem lama → ditunda OQ-EDM-016.
	g.PolisWarisan = append(g.PolisWarisan, &tiruan.PolisWarisan{IDPega: "UJI-IDPEGA-EDM", NoPolis: "UJI-PL-W", ProdKe: 2, EdmType: "1"})
	if _, err := l.BuatKasus(ctx, pelakuUji, services.MasukanKasus{NomorPolis: "UJI-PL-W", EdmType: "1"}); !errors.Is(err, services.ErrSumberWarisanEDM) {
		t.Errorf("sumber warisan EDM: %v", err)
	}
	// Polis yang versi terakhirnya Batal ditolak gerbang 4.
	g.PolisWarisan = append(g.PolisWarisan, &tiruan.PolisWarisan{IDPega: "UJI-IDPEGA-BTL", NoPolis: "UJI-PL-B", ProdKe: 3, EdmType: "3"})
	if _, err := l.BuatKasus(ctx, pelakuUji, services.MasukanKasus{NomorPolis: "UJI-PL-B", EdmType: "1"}); !errors.As(err, &gk) ||
		!reflect.DeepEqual(gk.Pesan, []string{models.PesanSudahBatal}) {
		t.Errorf("polis sudah Batal: %v", err)
	}
	if len(j.catatan) != 0 {
		t.Errorf("jejak tertulis untuk kasus yang ditolak: %+v", j.catatan)
	}
	// Jejak gagal = pembuatan gagal (ADR-U-0007: transisi tanpa jejak tidak boleh lolos).
	l2 := services.BaruLayanan(gudangPolisNB(), tiruan.Transaksi, nil).DenganJejak(&jejakTiruan{galat: errors.New("UJI-jejak")})
	if _, err := l2.BuatKasus(ctx, pelakuUji, services.MasukanKasus{NomorPolis: "UJI-PL-NB", EdmType: "1"}); err == nil {
		t.Error("galat jejak ditelan")
	}
	// Jejak bawaan gagal terang.
	l3 := services.BaruLayanan(gudangPolisNB(), tiruan.Transaksi, nil)
	if _, err := l3.BuatKasus(ctx, pelakuUji, services.MasukanKasus{NomorPolis: "UJI-PL-NB", EdmType: "1"}); !errors.Is(err, jejak.ErrJejakBelumDiputuskan) {
		t.Errorf("jejak bawaan: %v", err)
	}
}

func TestPolisLamaDanRincianPeserta(t *testing.T) {
	ctx := context.Background()
	g := gudangPolisNB()
	g.RekapWarisan = []map[string]string{{"PL_NUMBER": "UJI-PL-NB", "CURRENCY": "IDR", "PREMIUM": "10"}, {"PL_NUMBER": "UJI-LAIN"}}
	l, _ := layananJejak(g)
	h, err := l.BuatKasus(ctx, pelakuUji, services.MasukanKasus{NomorPolis: "UJI-PL-NB", EdmType: "1"})
	if err != nil {
		t.Fatal(err)
	}
	pl, err := l.PolisLama(ctx, pelakuUji, h.ID, 1, 20)
	if err != nil || pl.Sumber != models.SumberAplikasi || pl.ProdKe != 1 || pl.Total != 3 || len(pl.Rekap) != 1 || pl.Tipe != "QR" {
		t.Fatalf("polis lama = %+v, %v", pl, err)
	}
	var salinan string
	for _, d := range g.Peserta {
		if d.PolisID == h.ID && d.ParentID == "UJI-NB-D1" {
			salinan = d.ID
		}
	}
	r, err := l.RincianPeserta(ctx, pelakuUji, h.ID, salinan)
	if err != nil || len(r.Spreading) != 1 || len(r.Spreading[0].Retro) != 1 || r.Spreading[0].TreatyTypeName != "UJI-QS" {
		t.Fatalf("rincian = %+v, %v", r, err)
	}
	if _, err := l.RincianPeserta(ctx, pelakuUji, h.ID, "UJI-NB-D1"); !errors.Is(err, services.ErrPesertaTidakAda) {
		t.Errorf("peserta versi lama terbaca lewat kasus: %v", err)
	}
	// Kasus resmi: polis lama = versi di bawah PROD_KE kasus.
	g.Polis[h.ID].Status = models.StatusKasusSelesai
	g.Polis[h.ID].NoPolis = "UJI-PL-NB"
	if pl, err = l.PolisLama(ctx, pelakuUji, h.ID, 1, 20); err != nil || pl.ProdKe != 1 || !strings.HasPrefix(pl.Peserta[0].ID, "UJI-NB-") {
		t.Errorf("polis lama kasus resmi = %+v, %v", pl, err)
	}
}
