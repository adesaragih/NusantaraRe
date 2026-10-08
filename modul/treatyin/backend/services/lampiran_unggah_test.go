package services_test

// Uji tombol `Upload file` panel Attachment — `TreatySaveAttachment` +
// `InsertGoogleStorage_Act` — di atas gudang dan penyimpanan tiruan.

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func (g *gudangTiruan) NamaAplikasiSimpanan(context.Context) (string, error) { return "rnmtest", nil }

func (g *gudangTiruan) CatatLampiran(_ context.Context, l models.LampiranBaru) (string, error) {
	g.lampiranBaru = append(g.lampiranBaru, l)
	return "20261008120000123", nil
}

type simpananTiruan struct {
	diminta []services.PermintaanSimpanan
	galat   error
	// geturl dan delete.
	urlBaru []services.PermintaanSimpanan
	dihapus []services.PermintaanSimpanan
}

func (s *simpananTiruan) URLBaru(_ context.Context, p services.PermintaanSimpanan) (services.JawabanSimpanan, error) {
	s.urlBaru = append(s.urlBaru, p)
	if s.galat != nil {
		return services.JawabanSimpanan{}, s.galat
	}
	return services.JawabanSimpanan{URLImage: "https://storage.googleapis.com/rnmtest/baru", Exp: "2030-01-01T00:00:00Z",
		AppFolder: "gs://rnmtest/" + p.Folder + p.Namafile, DateTime: "01/01/2030 00:00:00"}, nil
}

func (s *simpananTiruan) Hapus(_ context.Context, p services.PermintaanSimpanan) error {
	s.dihapus = append(s.dihapus, p)
	return s.galat
}

func (s *simpananTiruan) Unggah(_ context.Context, p services.PermintaanSimpanan) (services.JawabanSimpanan, error) {
	s.diminta = append(s.diminta, p)
	if s.galat != nil {
		return services.JawabanSimpanan{}, s.galat
	}
	return services.JawabanSimpanan{
		URLImage:  "https://storage.googleapis.com/rnmtest/" + p.Folder + p.Namafile,
		Exp:       "2026-10-08T12:30:00Z",
		AppFolder: "gs://rnmtest/" + p.Folder + p.Namafile,
	}, nil
}

func gudangLampiran() *gudangTiruan {
	g := gudangSimpan()
	g.kepalaTreatyIn["1001001"]["ProportionType"] = "NonProportional"
	g.katalogKategori = map[string]string{
		"00002": "Approval Email",
		"00007": "Pega Proportional Calculation /Perhitungan Pega Proportional",
	}
	return g
}

func TestUnggahLampiranMencatatObjekDanBarisSepertiPega(t *testing.T) {
	g, s := gudangLampiran(), &simpananTiruan{}
	h, err := services.LayananDenganSimpanan(g, s).UnggahLampiran(context.Background(), admin, services.MasukanUnggahLampiran{
		IDKontrak: "1001001", KodeKategori: "00002",
		Berkas: []services.BerkasUnggah{{Nama: "Summary Treaty Leader 2025.PDF", Isi: []byte("%PDF-1.4 isi berkas uji")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Berkas) != 1 || !h.Berkas[0].Berhasil {
		t.Fatalf("hasil %+v", h.Berkas)
	}
	// [8] Set Data: Folder Contract/Doc/YYYY/MM/, Durasi 1800, ext huruf kecil.
	p := s.diminta[0]
	if p.App != "rnmtest" || p.Durasi != 1800 || !strings.HasPrefix(p.Folder, "Contract/Doc/") || !strings.HasSuffix(p.Folder, "/") ||
		p.Ext != "pdf" || p.MimeType != "application/pdf" || !strings.HasSuffix(p.Namafile, " - Summary Treaty Leader 2025.PDF") {
		t.Errorf("permintaan %+v", p)
	}
	// [1.6] InsertAttachment2_Sql: CATEGORY = nama tampil, FILEMIMETYPE = ext.
	l := g.lampiranBaru[0]
	if l.IDKontrak != "1001001" || l.KodeKategori != "00002" || l.NamaKategori != "Approval Email" ||
		l.NamaBerkas != "Summary Treaty Leader 2025.PDF" || l.Ekstensi != "pdf" || l.Pengguna != "ADESAMUEL" ||
		len(l.ImageID) != 32 || l.App != "rnmtest" || l.Exp != "08/10/2026 12:30:00" || l.NamaObjek != p.Namafile {
		t.Errorf("lampiran %+v", l)
	}
}

func TestUnggahLampiranNonPropMemakaiNamaKategoriNonProp(t *testing.T) {
	g, s := gudangLampiran(), &simpananTiruan{}
	if _, err := services.LayananDenganSimpanan(g, s).UnggahLampiran(context.Background(), admin, services.MasukanUnggahLampiran{
		IDKontrak: "1001001", KodeKategori: "00007",
		Berkas: []services.BerkasUnggah{{Nama: "hitung.xlsx", Isi: []byte("PK isi lembar kerja")}},
	}); err != nil {
		t.Fatal(err)
	}
	if got := g.lampiranBaru[0].NamaKategori; got != "Pega Non Proportional Calculation /Perhitungan Pega Non Proportional" {
		t.Errorf("CATEGORY %q", got)
	}
}

// [5] Exit-Activity: ext kosong / tak dikenal → berkas DILEWATI (dilaporkan),
// berkas lain tetap diunggah.
func TestUnggahLampiranJenisTakDikenalDilaporkanLainTetapJalan(t *testing.T) {
	g, s := gudangLampiran(), &simpananTiruan{}
	h, err := services.LayananDenganSimpanan(g, s).UnggahLampiran(context.Background(), admin, services.MasukanUnggahLampiran{
		IDKontrak: "1001001", KodeKategori: "00002",
		Berkas: []services.BerkasUnggah{
			{Nama: "tanpaekstensi", Isi: []byte("isi isi isi isi")},
			{Nama: "surat.pdf", Isi: []byte("%PDF-1.4 isi berkas uji")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if h.Berkas[0].Berhasil || !h.Berkas[1].Berhasil || len(g.lampiranBaru) != 1 || len(s.diminta) != 1 {
		t.Errorf("hasil %+v, tercatat %d, terkirim %d", h.Berkas, len(g.lampiranBaru), len(s.diminta))
	}
}

func TestUnggahLampiranDitolakTanpaKontrakAtauKategori(t *testing.T) {
	g, s := gudangLampiran(), &simpananTiruan{}
	l := services.LayananDenganSimpanan(g, s)
	berkas := []services.BerkasUnggah{{Nama: "a.pdf", Isi: []byte("%PDF-1.4 isi berkas uji")}}
	if _, err := l.UnggahLampiran(context.Background(), admin, services.MasukanUnggahLampiran{IDKontrak: "9999999", KodeKategori: "00002", Berkas: berkas}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("kontrak tak ada: %v", err)
	}
	if _, err := l.UnggahLampiran(context.Background(), admin, services.MasukanUnggahLampiran{IDKontrak: "1001001", KodeKategori: "99999", Berkas: berkas}); !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("kategori tak ada: %v", err)
	}
	if len(s.diminta) != 0 || len(g.lampiranBaru) != 0 {
		t.Error("ada yang terkirim/tercatat")
	}
}

// [12] URLImage kosong / layanan gagal → NOL baris tercatat.
func TestUnggahLampiranGagalLayananNolCatatan(t *testing.T) {
	g, s := gudangLampiran(), &simpananTiruan{galat: services.ErrSimpananGagal}
	_, err := services.LayananDenganSimpanan(g, s).UnggahLampiran(context.Background(), admin, services.MasukanUnggahLampiran{
		IDKontrak: "1001001", KodeKategori: "00002",
		Berkas: []services.BerkasUnggah{{Nama: "a.pdf", Isi: []byte("%PDF-1.4 isi berkas uji")}},
	})
	if !errors.Is(err, services.ErrSimpananGagal) || len(g.lampiranBaru) != 0 {
		t.Errorf("galat %v, tercatat %d", err, len(g.lampiranBaru))
	}
}

func TestNamaObjekLampiranFormatJava(t *testing.T) {
	kini := time.Date(2025, 8, 22, 18, 41, 50, 98*int(time.Millisecond), time.FixedZone("WIB", 7*3600))
	if got := services.NamaObjekLampiran(kini, "a.pdf"); got != "20250822-064150-98 - a.pdf" {
		t.Errorf("nama objek %q", got)
	}
}
