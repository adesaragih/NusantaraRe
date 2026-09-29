package services_test

// Uji layanan klausul - TANPA Oracle (tiket 08).

import (
	"context"
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

type gudangKlausulUji struct {
	baris map[string]models.KlausulTreaty
	urut  int
	dobel string
}

func (g *gudangKlausulUji) Daftar(_ context.Context, tahunID, descID, parent string) ([]models.KlausulTreaty, error) {
	var hasil []models.KlausulTreaty
	for _, k := range g.baris {
		if k.TreatyYearID == tahunID && k.TreatyDescID == descID && k.ParentReinsTypeID == parent {
			hasil = append(hasil, k)
		}
	}
	sort.Slice(hasil, func(i, j int) bool { return hasil[i].ID < hasil[j].ID })
	return hasil, nil
}
func (g *gudangKlausulUji) Ambil(_ context.Context, tahunID, id string) (models.KlausulTreaty, error) {
	k, ada := g.baris[id]
	if !ada || k.TreatyYearID != tahunID {
		return models.KlausulTreaty{}, repository.ErrKlausulTidakAda
	}
	return k, nil
}
func (g *gudangKlausulUji) Induk(_ context.Context, tahunID, descID, reins string) (models.KlausulTreaty, error) {
	for _, k := range g.baris {
		if k.TreatyYearID == tahunID && k.TreatyDescID == descID && k.ParentReinsTypeID == "00" && k.ReinsTypeID == reins {
			return k, nil
		}
	}
	return models.KlausulTreaty{}, repository.ErrKlausulTidakAda
}
func (g *gudangKlausulUji) PctAnakLain(_ context.Context, _ *repository.Tx, tahunID, descID, parent, kecuali string) ([]*apd.Decimal, error) {
	var hasil []*apd.Decimal
	for id, k := range g.baris {
		if k.TreatyYearID == tahunID && k.TreatyDescID == descID && k.ParentReinsTypeID == parent && id != kecuali {
			hasil = append(hasil, k.Pct)
		}
	}
	return hasil, nil
}
func (g *gudangKlausulUji) CariDobel(context.Context, *repository.Tx, models.KlausulTreaty, []string) (string, error) {
	return g.dobel, nil
}
func (g *gudangKlausulUji) Sisip(_ context.Context, _ *repository.Tx, k models.KlausulTreaty) (string, error) {
	g.urut++
	k.ID = "1000000" + string(rune('0'+g.urut))
	g.baris[k.ID] = k
	return k.ID, nil
}
func (g *gudangKlausulUji) Perbarui(_ context.Context, _ *repository.Tx, k models.KlausulTreaty) error {
	g.baris[k.ID] = k
	return nil
}

type masterKlausulUji struct{}

func (masterKlausulUji) JenisKlausul(_ context.Context, isXOL string) ([]repository.JenisKlausulMasterTCO, error) {
	semua := []repository.JenisKlausulMasterTCO{
		{ID: "10001", DescName: "UJI TREATY LIMIT", IsXOL: "0"}, {ID: "10009", DescName: "UJI EPI", IsXOL: "0"},
		{ID: "10013", DescName: "UJI EXCLUSION", IsXOL: "0"}, {ID: "10017", DescName: "UJI LIMIT MB", IsXOL: "1"},
		{ID: "10099", DescName: "UJI JENIS BARU", IsXOL: "1"},
	}
	var hasil []repository.JenisKlausulMasterTCO
	for _, j := range semua {
		if isXOL == "" || j.IsXOL == isXOL {
			hasil = append(hasil, j)
		}
	}
	return hasil, nil
}
func (masterKlausulUji) CariPilihan(_ context.Context, master, _ string) ([]repository.PilihanMasterTCO, error) {
	return []repository.PilihanMasterTCO{{ID: "UJI-1", Nama: "UJI " + master}}, nil
}
func (masterKlausulUji) AmbilPilihan(_ context.Context, master, id string) (repository.PilihanMasterTCO, error) {
	if id != "UJI-1" {
		return repository.PilihanMasterTCO{}, repository.ErrPilihanMasterTidakAda
	}
	return repository.PilihanMasterTCO{ID: id, Nama: "UJI NAMA " + master}, nil
}

type tahunKlausulUji struct{ dikunci *int }

func (tahunKlausulUji) Ambil(_ context.Context, id string) (models.TahunTreaty, error) {
	if id != "1000001" {
		return models.TahunTreaty{}, repository.ErrTahunTreatyTidakAda
	}
	return models.TahunTreaty{ID: id, TreatyYear: "2026", TreatyGroupID: "10001", TreatyGroupName: "UJI GRUP"}, nil
}
func (t tahunKlausulUji) Kunci(context.Context, *repository.Tx, string) error {
	*t.dikunci++
	return nil
}

func layananKlausul(g *gudangKlausulUji, dikunci *int) *services.KlausulTCO {
	return services.New(nil).KlausulTCO().DenganGudang(g).DenganMaster(masterKlausulUji{}).
		DenganTahun(tahunKlausulUji{dikunci: dikunci}).
		DenganJenis(jenisKontrakUji{{ID: "10003", Note: "UJI QUOTA SHARE"}, {ID: "10005", Note: "UJI SURPLUS"}}).
		DenganKurs(kursKlausulUji{}).DenganTransaksi(transaksiUji).DenganJam(jamUji)
}

func gudangKlausulKosong() *gudangKlausulUji {
	return &gudangKlausulUji{baris: map[string]models.KlausulTreaty{}}
}

// epi - induk EPI; `Usd` TIDAK dikirim: turunan `Rp / Kurs` (tiket 11).
func epi(reins, rp string) services.KlausulMasuk {
	return services.KlausulMasuk{DescID: "10009", Medan: map[string]string{"ReinsTypeID": reins, "Rp": rp}}
}

// kursKlausulUji - kurs berlaku 12.5 (angka bulat supaya turunan mudah dibaca);
// `kosong` meniru periode tanpa baris kurs.
type kursKlausulUji struct{ kosong bool }

func (k kursKlausulUji) Berlaku(_ context.Context, tahun models.TahunTreaty) (models.KursTCO, error) {
	if k.kosong {
		return models.KursTCO{}, models.GalatKursTidakAda{TreatyYear: tahun.TreatyYear}
	}
	return models.KursTCO{ToIDR: apd.New(125, -1), IDCurrency: "UJI-USD", Currency: "USD", Quarter: "0"}, nil
}

func TestKlausulTanpaIdentitasDanBawaan(t *testing.T) {
	n := 0
	l := layananKlausul(gudangKlausulKosong(), &n)
	if _, err := l.Simpan(context.Background(), services.Pelaku{}, "1000001", epi("10003", "1")); !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Errorf("identitas: %v", err)
	}
	if _, err := services.New(nil).KlausulTCO().JenisKlausul(context.Background(), pelakuUjiTCO, ""); !errors.Is(err, services.ErrGudangKlausulBelumDisuntik) {
		t.Errorf("bawaan: %v", err)
	}
}

// AC 26: jenis dari master; aturan dilekatkan dari kode; jenis master tanpa
// aturan ditandai, bukan dilepas tanpa validasi.
func TestJenisKlausulDariMasterDenganAturan(t *testing.T) {
	n := 0
	d, err := layananKlausul(gudangKlausulKosong(), &n).JenisKlausul(context.Background(), pelakuUjiTCO, "0")
	if err != nil || len(d) != 3 {
		t.Fatalf("non-XOL: %+v %v", d, err)
	}
	for _, j := range d {
		if j.ID == "10009" && (len(j.Aturan) != 2 || j.Aturan[0].Jenis != "EPI" || j.Aturan[1].Jenis != "EpiList") {
			t.Errorf("EPI aturan: %+v", j.Aturan)
		}
		if j.ID == "10013" && len(j.Aturan) != 4 {
			t.Errorf("ExclutionTreaty subjenis %d, mau 4", len(j.Aturan))
		}
	}
	x, _ := layananKlausul(gudangKlausulKosong(), &n).JenisKlausul(context.Background(), pelakuUjiTCO, "1")
	for _, j := range x {
		if j.ID == "10099" && j.Catatan == "" {
			t.Error("jenis master tanpa aturan tidak ditandai")
		}
		if j.ID == "10017" && j.Aturan[0].Ditahan == "" {
			t.Error("LimitMB tidak ditahan")
		}
	}
}

// AC 33, 35: induk EPI - gerbang menyebut medan; nama jenis & reasuransi dari
// master; TREATYYEARID tahun; sentinel "00"; jejak; tahun dikunci.
func TestKlausulIndukEPI(t *testing.T) {
	g, n := gudangKlausulKosong(), 0
	l := layananKlausul(g, &n)
	_, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", epi("10003", ""))
	if !errors.Is(err, models.ErrKlausulMedanWajib) || !strings.Contains(err.Error(), "Rp") {
		t.Fatalf("wajib: %v", err)
	}
	// Tiket 11: `Usd` induk hanya dibaca di form Pega - klien tidak mengirimnya.
	kirimUsd := epi("10003", "1000")
	kirimUsd.Medan["Usd"] = "70"
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", kirimUsd); !errors.Is(err, models.ErrMedanBukanMilikJenis) {
		t.Errorf("Usd dari klien: %v", err)
	}
	h, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", epi("10003", "1.000.000,5"))
	if !errors.Is(err, models.ErrPersenBukanDesimal) {
		t.Errorf("desimal ganda: %v", err)
	}
	h, err = l.Simpan(context.Background(), pelakuUjiTCO, "1000001", epi("10003", "1000000,5"))
	if err != nil {
		t.Fatal(err)
	}
	k := h.Klausul
	if k.TreatyYearID != "1000001" || k.TreatyYear != "2026" || k.TreatyDescName != "UJI EPI" || k.ReinsTypeName != "UJI QUOTA SHARE" ||
		k.ParentReinsTypeID != "00" || k.Medan["Rp"] != "1000000.5" || k.Medan["Usd"] != "80000.04000000" ||
		k.Kurs != "12.5" || n != 1 {
		t.Errorf("induk: %+v kunci %d", k, n)
	}
}

// HitungRpUsd + batas total: anak EpiList - Rp/Usd TURUNAN induk; klien tidak
// boleh mengirimnya; total > 100 ditolak.
func TestKlausulAnakTurunanDanTotal(t *testing.T) {
	g, n := gudangKlausulKosong(), 0
	l := layananKlausul(g, &n)
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", epi("10003", "1000")); err != nil {
		t.Fatal(err)
	}
	anak := func(reins, pct string) services.KlausulMasuk {
		return services.KlausulMasuk{DescID: "10009", Anak: true, ParentReinsTypeID: "10003",
			Medan: map[string]string{"ReinsTypeID": reins, "Pct": pct}}
	}
	h, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", anak("10005", "25"))
	if err != nil {
		t.Fatal(err)
	}
	if h.Klausul.Medan["Rp"] != "250.00000000" || h.Klausul.Medan["Usd"] != "20.00000000" || h.TotalPct != "25" || h.Peringatan != "" {
		t.Errorf("turunan: %+v total %s peringatan %q", h.Klausul.Medan, h.TotalPct, h.Peringatan)
	}
	m := anak("10003", "10")
	m.Medan["Rp"] = "1"
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", m); !errors.Is(err, models.ErrMedanBukanMilikJenis) {
		t.Errorf("Rp dari klien: %v", err)
	}
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", anak("10003", "75.00000001")); !errors.Is(err, models.ErrTotalPctAnakMelebihi100) {
		t.Errorf("total > 100: %v", err)
	}
	yatim := anak("10003", "10")
	yatim.ParentReinsTypeID = "10005"
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", yatim); !errors.Is(err, services.ErrKlausulTidakAda) {
		t.Errorf("anak tanpa induk: %v", err)
	}
	d, err := l.Daftar(context.Background(), pelakuUjiTCO, "1000001", "10009", "10003")
	if err != nil || d.Total != 1 || d.TotalPct != "25" || d.Peringatan != "" {
		t.Errorf("daftar anak: %+v %v", d, err)
	}
}

// TreatyTestChildTotal_Act: hanya TreatyLimitChild, dan hanya PERINGATAN.
func TestKlausulPeringatanSpreadingTreatyLimitChild(t *testing.T) {
	g, n := gudangKlausulKosong(), 0
	l := layananKlausul(g, &n)
	induk := services.KlausulMasuk{DescID: "10001", Medan: map[string]string{"ReinsTypeID": "10003", "Rp": "100"}}
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", induk); err != nil {
		t.Fatal(err)
	}
	anak := services.KlausulMasuk{DescID: "10001", Anak: true, ParentReinsTypeID: "10003",
		Medan: map[string]string{"ReinsTypeID": "10005", "Pct": "60"}}
	h, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", anak)
	if err != nil || h.Peringatan != "Please make sure spreading is 100%" {
		t.Fatalf("60%%: %+v %v", h, err)
	}
	anak.Medan = map[string]string{"ReinsTypeID": "10003", "Pct": "40"}
	if h, err = l.Simpan(context.Background(), pelakuUjiTCO, "1000001", anak); err != nil || h.Peringatan != "" || h.TotalPct != "100" {
		t.Errorf("100%%: %+v %v", h, err)
	}
}

// AC 30: "sudah pernah diinput" ditolak dengan pesan jelas; AC 36 ditahan.
func TestKlausulDobelDanDitahan(t *testing.T) {
	g, n := gudangKlausulKosong(), 0
	g.dobel = "10000009"
	_, err := layananKlausul(g, &n).Simpan(context.Background(), pelakuUjiTCO, "1000001", epi("10003", "1"))
	if !errors.Is(err, services.ErrKlausulDobel) || !strings.Contains(err.Error(), "Data sudah pernah di Input") || !strings.Contains(err.Error(), "10000009") {
		t.Errorf("dobel: %v", err)
	}
	for _, desc := range []string{"10017", "10002"} {
		_, err := layananKlausul(gudangKlausulKosong(), &n).Simpan(context.Background(), pelakuUjiTCO, "1000001",
			services.KlausulMasuk{DescID: desc, Medan: map[string]string{}})
		if !errors.Is(err, models.ErrKlausulDitahan) {
			t.Errorf("%s: %v", desc, err)
		}
	}
	if _, err := layananKlausul(gudangKlausulKosong(), &n).Simpan(context.Background(), pelakuUjiTCO, "1000001",
		services.KlausulMasuk{DescID: "10099", Medan: map[string]string{}}); !errors.Is(err, models.ErrKlausulJenisTakDikenal) {
		t.Errorf("jenis tanpa aturan: %v", err)
	}
}

// ExclutionTreaty/Occupation: nama dari master; ID asing ditolak.
func TestKlausulExclusionOccupation(t *testing.T) {
	g, n := gudangKlausulKosong(), 0
	m := services.KlausulMasuk{DescID: "10013", Subjenis: "Occupation", Medan: map[string]string{
		"ID_Occupation": "UJI-1", "Occupation": "KARANGAN", "Line": "A", "Usd": "1", "Rp": "15000"}}
	h, err := layananKlausul(g, &n).Simpan(context.Background(), pelakuUjiTCO, "1000001", m)
	if err != nil || h.Klausul.Medan["Occupation"] != "UJI NAMA OCCUPATION" || h.Klausul.Subjenis != "Occupation" {
		t.Errorf("occupation: %+v %v", h, err)
	}
	m.Medan["ID_Occupation"] = "UJI-9"
	if _, err := layananKlausul(gudangKlausulKosong(), &n).Simpan(context.Background(), pelakuUjiTCO, "1000001", m); !errors.Is(err, services.ErrPilihanDiLuarMaster) {
		t.Errorf("occupation asing: %v", err)
	}
}

// Pembaruan mempertahankan kolom di luar form dan tidak memindah jenis; induk
// beranak tidak boleh berganti jenis reasuransi. Tiket 11: `KURS` = kurs
// berlaku saat disimpan, `Usd` dihitung ulang dari `Rp`.
func TestKlausulPerbarui(t *testing.T) {
	g, n := gudangKlausulKosong(), 0
	l := layananKlausul(g, &n)
	h, _ := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", epi("10003", "100"))
	id := h.Klausul.ID
	lama := g.baris[id]
	lama.Kurs = apd.New(15500, 0)
	lama.LayerType = "UJI-L"
	g.baris[id] = lama
	m := epi("10003", "200")
	m.ID = id
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", m); err != nil {
		t.Fatal(err)
	}
	if g.baris[id].Kurs.Text('f') != "12.5" || g.baris[id].Rp.Text('f') != "200" || g.baris[id].Usd.Text('f') != "16.00000000" ||
		g.baris[id].LayerType != "UJI-L" {
		t.Errorf("perbarui: %+v", g.baris[id])
	}
	pindah := services.KlausulMasuk{ID: id, DescID: "10001", Medan: map[string]string{"ReinsTypeID": "10003", "Rp": "1"}}
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", pindah); !errors.Is(err, services.ErrKlausulJenisBerubah) {
		t.Errorf("pindah jenis: %v", err)
	}
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", services.KlausulMasuk{DescID: "10009", Anak: true,
		ParentReinsTypeID: "10003", Medan: map[string]string{"ReinsTypeID": "10005", "Pct": "10"}}); err != nil {
		t.Fatal(err)
	}
	ganti := epi("10005", "200")
	ganti.ID = id
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", ganti); !errors.Is(err, services.ErrKlausulIndukBeranak) {
		t.Errorf("induk beranak berganti jenis: %v", err)
	}
}

// Tiket 11 (ADR-0015): jenis berkurs - tujuh induk ber-Rp/Usd dan tujuh anak -
// ditolak dengan pesan VERBATIM bila tahun tidak punya kurs; jenis lain tidak.
func TestKlausulTanpaKursDitolak(t *testing.T) {
	n := 0
	l := layananKlausul(gudangKlausulKosong(), &n).DenganKurs(kursKlausulUji{kosong: true})
	_, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", epi("10003", "1000"))
	if !errors.Is(err, services.ErrKursTidakAda) || err.Error() != "Tidak ada Nilai Kurs di Tahun : 2026" {
		t.Errorf("induk tanpa kurs: %v", err)
	}
	_, err = l.Simpan(context.Background(), pelakuUjiTCO, "1000001", services.KlausulMasuk{DescID: "10009", Anak: true,
		ParentReinsTypeID: "10003", Medan: map[string]string{"ReinsTypeID": "10005", "Pct": "10"}})
	if !errors.Is(err, services.ErrKursTidakAda) {
		t.Errorf("anak tanpa kurs: %v", err)
	}
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", services.KlausulMasuk{DescID: "10013", Subjenis: "Occupation",
		Medan: map[string]string{"ID_Occupation": "UJI-1", "Line": "A", "Usd": "1", "Rp": "15000"}}); err != nil {
		t.Errorf("exclusion tidak berkurs di Pega: %v", err)
	}
	d, _ := layananKlausul(gudangKlausulKosong(), &n).JenisKlausul(context.Background(), pelakuUjiTCO, "0")
	for _, j := range d {
		for _, a := range j.Aturan {
			if a.Jenis == "EPI" && (!a.Berkurs || a.Konversi != models.KonversiRpKeUsd || len(a.Turunan) != 1) {
				t.Errorf("aturan EPI: %+v", a)
			}
			if a.Jenis == "ExclutionTreaty" && a.Subjenis == "Occupation" && (a.Berkurs || a.Konversi != models.KonversiDuaArah) {
				t.Errorf("aturan exclusion: %+v", a)
			}
		}
	}
}

// Temuan /code-review: induk yang Rp-nya berubah menghitung ulang anaknya.
func TestKlausulIndukBerubahMenghitungUlangAnak(t *testing.T) {
	g, n := gudangKlausulKosong(), 0
	l := layananKlausul(g, &n)
	h, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", epi("10003", "1000"))
	if err != nil {
		t.Fatal(err)
	}
	a, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", services.KlausulMasuk{DescID: "10009", Anak: true,
		ParentReinsTypeID: "10003", Medan: map[string]string{"ReinsTypeID": "10005", "Pct": "25"}})
	if err != nil || a.Klausul.Medan["Rp"] != "250.00000000" {
		t.Fatalf("anak: %+v %v", a, err)
	}
	ubah := epi("10003", "2000")
	ubah.ID = h.Klausul.ID
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", ubah); err != nil {
		t.Fatal(err)
	}
	c := g.baris[a.Klausul.ID]
	if c.Rp.Text('f') != "500.00000000" || c.Usd.Text('f') != "40.00000000" {
		t.Errorf("anak tidak dihitung ulang: Rp %v Usd %v", c.Rp, c.Usd)
	}
}
