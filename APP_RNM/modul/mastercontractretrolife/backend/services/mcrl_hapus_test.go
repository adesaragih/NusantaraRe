package services_test

// Hapus berjenjang (paket 7, tiket 09): K2 - kaskade di Go, satu transaksi,
// anak lebih dulu, SESUDAH popup; batal = nol baris; tahun abadi.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/models"
	"nusantarare/modul/mastercontractretrolife/backend/services"
	"nusantarare/modul/mastercontractretrolife/backend/tiruan"
)

func gudangPohon() *tiruan.Gudang {
	g := gudangReinsurer()
	g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1"}
	g.Reinsurer["UJI-R2"] = models.Reinsurer{ID: "UJI-R2", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1"}
	g.Security["UJI-S1"] = models.SecurityReinsurer{ID: "UJI-S1", TreatyReinsurerID: "UJI-R1", TreatyContractID: "UJI-K1"}
	g.Security["UJI-S2"] = models.SecurityReinsurer{ID: "UJI-S2", TreatyReinsurerID: "UJI-R1", TreatyContractID: "UJI-K1"}
	g.Security["UJI-S3"] = models.SecurityReinsurer{ID: "UJI-S3", TreatyReinsurerID: "UJI-R2", TreatyContractID: "UJI-K1"}
	g.Business["UJI-B1"] = models.Business{ID: "UJI-B1", TreatyContractID: "UJI-K1"}
	// Kontrak lain - tidak boleh tersentuh.
	g.Kontrak["UJI-K2"] = models.Kontrak{ID: "UJI-K2", IDTreatyYear: "UJI-T1"}
	g.Reinsurer["UJI-R9"] = models.Reinsurer{ID: "UJI-R9", TreatyContractID: "UJI-K2"}
	g.Security["UJI-S9"] = models.SecurityReinsurer{ID: "UJI-S9", TreatyReinsurerID: "UJI-R9", TreatyContractID: "UJI-K2"}
	g.Business["UJI-B9"] = models.Business{ID: "UJI-B9", TreatyContractID: "UJI-K2"}
	return g
}

func TestDampakHapusKontrakMenyebutApaDanBerapa(t *testing.T) {
	j, err := layananUji(gudangPohon()).DampakHapus(context.Background(), pelaku, services.HapusKontrak, "UJI-K1")
	if err != nil {
		t.Fatal(err)
	}
	if j.Dampak != (models.Dampak{Security: 3, Reinsurer: 2, Business: 1}) {
		t.Errorf("dampak kontrak: %+v", j.Dampak)
	}
	j, _ = layananUji(gudangPohon()).DampakHapus(context.Background(), pelaku, services.HapusReinsurer, "UJI-R1")
	if j.Dampak != (models.Dampak{Security: 2}) {
		t.Errorf("dampak reinsurer: %+v", j.Dampak)
	}
	j, _ = layananUji(gudangPohon()).DampakHapus(context.Background(), pelaku, services.HapusBusiness, "UJI-B1")
	if j.Dampak != (models.Dampak{}) {
		t.Errorf("daun tanpa anak: %+v", j.Dampak)
	}
}

func TestHapusKontrakKaskadeTanpaYatimDanPohonLainUtuh(t *testing.T) {
	g := gudangPohon()
	h, err := layananUji(g).Hapus(context.Background(), pelaku, services.HapusKontrak, "UJI-K1",
		models.Dampak{Security: 3, Reinsurer: 2, Business: 1})
	if err != nil {
		t.Fatal(err)
	}
	if h.Pesan != "Data Berhasil di Hapus" || h.Terhapus != (models.Dampak{Security: 3, Reinsurer: 2, Business: 1}) {
		t.Errorf("hasil: %+v", h)
	}
	if _, ada := g.Kontrak["UJI-K1"]; ada || len(g.Reinsurer) != 1 || len(g.Security) != 1 || len(g.Business) != 1 {
		t.Errorf("sisa: kontrak %v, %d reinsurer, %d security, %d business", g.Kontrak, len(g.Reinsurer), len(g.Security), len(g.Business))
	}
	if g.Komit != 1 {
		t.Errorf("%d transaksi, mau 1", g.Komit)
	}
}

// ⛔ Batal / konfirmasi basi = NOL baris terhapus (tiket 09 AC).
func TestHapusDenganDampakBasiNolBarisTerhapus(t *testing.T) {
	g := gudangPohon()
	_, err := layananUji(g).Hapus(context.Background(), pelaku, services.HapusKontrak, "UJI-K1",
		models.Dampak{Security: 2, Reinsurer: 2, Business: 1})
	if !errors.Is(err, services.ErrDampakBerubah) {
		t.Fatalf("dampak basi: %v", err)
	}
	if len(g.Kontrak) != 2 || len(g.Reinsurer) != 3 || len(g.Security) != 4 || len(g.Business) != 2 {
		t.Error("baris terhapus walau konfirmasi basi")
	}
}

func TestHapusGagalDiTengahDipulihkan(t *testing.T) {
	g := gudangPohon()
	g.GagalTulis["HapusKontrak"] = errors.New("UJI ORA-02292")
	_, err := layananUji(g).Hapus(context.Background(), pelaku, services.HapusKontrak, "UJI-K1",
		models.Dampak{Security: 3, Reinsurer: 2, Business: 1})
	if err == nil || len(g.Security) != 4 || len(g.Reinsurer) != 3 {
		t.Errorf("gagal di tengah: %v - anak harus dipulihkan (satu transaksi)", err)
	}
}

// Lapis kedua: penghapus menyentuh lebih dari yang dicacah (baris lahir di
// antara cacah dan hapus) - seluruh hapus dibatalkan.
func TestHapusMenyentuhLebihDariDicacahDibatalkan(t *testing.T) {
	g := gudangPohon()
	g.SelaHapus = func() {
		g.Security["UJI-S-SELA"] = models.SecurityReinsurer{ID: "UJI-S-SELA", TreatyReinsurerID: "UJI-R2"}
	}
	_, err := layananUji(g).Hapus(context.Background(), pelaku, services.HapusKontrak, "UJI-K1",
		models.Dampak{Security: 3, Reinsurer: 2, Business: 1})
	if !errors.Is(err, services.ErrDampakBerubah) {
		t.Fatalf("hapus menyentuh lebih: %v", err)
	}
	if _, ada := g.Kontrak["UJI-K1"]; !ada || len(g.Security) != 4 {
		t.Errorf("hapus tidak dibatalkan: kontrak ada=%v, %d security (mau 4 seperti sebelum)", ada, len(g.Security))
	}
}

func TestHapusReinsurerHanyaSecuritynya(t *testing.T) {
	g := gudangPohon()
	if _, err := layananUji(g).Hapus(context.Background(), pelaku, services.HapusReinsurer, "UJI-R1", models.Dampak{Security: 2}); err != nil {
		t.Fatal(err)
	}
	if _, ada := g.Reinsurer["UJI-R1"]; ada || len(g.Security) != 2 || len(g.Business) != 2 {
		t.Errorf("hapus reinsurer: %d security, %d business", len(g.Security), len(g.Business))
	}
}

func TestHapusDaunDanPesanBusinessVerbatim(t *testing.T) {
	g := gudangPohon()
	h, err := layananUji(g).Hapus(context.Background(), pelaku, services.HapusBusiness, "UJI-B1", models.Dampak{})
	if err != nil || h.Pesan != "Data Dengan ID UJI-B1 Berhasil di Hapus" || len(g.Business) != 1 {
		t.Errorf("hapus business: %+v %v", h, err)
	}
	if _, err := layananUji(g).Hapus(context.Background(), pelaku, services.HapusSecurity, "UJI-S3", models.Dampak{}); err != nil || len(g.Security) != 3 {
		t.Errorf("hapus security: %v", err)
	}
}

func TestHapusYangTidakAda404(t *testing.T) {
	_, err := layananUji(gudangPohon()).Hapus(context.Background(), pelaku, services.HapusReinsurer, "UJI-TAK-ADA", models.Dampak{})
	if !errors.Is(err, services.ErrReinsurerTidakAda) {
		t.Errorf("hapus tak ada: %v", err)
	}
	_, err = layananUji(gudangPohon()).DampakHapus(context.Background(), pelaku, services.JenisHapus("tahun"), "UJI-T1")
	if !errors.Is(err, services.ErrJenisHapusTidakAda) {
		t.Errorf("hapus tahun treaty (abadi): %v", err)
	}
}
