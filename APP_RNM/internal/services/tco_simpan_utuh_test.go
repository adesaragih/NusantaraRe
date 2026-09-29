package services_test

// Uji orkestrator simpan utuh - TANPA Oracle (tiket 09). Atomisitas nyata
// (rollback, identitas sementara) dibuktikan uji `db`; di sini dibuktikan
// bahwa SELURUH penulis berjalan di SATU transaksi, galat baris ke-N keluar
// dari transaksi itu (sehingga ia dibatalkan), dan identitas ditetapkan hanya
// bila seluruh baris lolos.

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

type transaksiRekamUji struct {
	dipanggil int
	gagal     bool
}

func (r *transaksiRekamUji) jalankan(_ context.Context, fn func(*repository.Tx) error) error {
	r.dipanggil++
	err := fn(nil)
	r.gagal = err != nil
	return err
}

type penetapUji struct {
	dipanggil int
	peta      map[string]string
	jejak     []string
}

func (p *penetapUji) tetapkan(context.Context, *repository.Tx) (map[string]string, error) {
	p.dipanggil++
	return p.peta, nil
}

func (p *penetapUji) rekam(_ context.Context, _ *repository.Tx, akun, tabel, baris, aksi, ket string, _ time.Time) error {
	p.jejak = append(p.jejak, strings.Join([]string{akun, tabel, baris, aksi, ket}, "|"))
	return nil
}

type rakitanUtuhUji struct {
	kontrak   *gudangKontrakUji
	reinsurer *gudangReinsurerUji
	security  *gudangSecurityUji
	business  *gudangBusinessUji
	klausul   *gudangKlausulUji
	tx        *transaksiRekamUji
	penetap   *penetapUji
	layanan   *services.SimpanUtuhTCO
}

func rakitUtuh() *rakitanUtuhUji {
	r := &rakitanUtuhUji{
		kontrak: &gudangKontrakUji{baris: map[string]models.KontrakTreaty{
			"1000003": {ID: "1000003", IDTreatyYear: "1000001", ReinsTypeID: "10003", ReinsTypeName: "UJI QUOTA SHARE"},
		}},
		reinsurer: &gudangReinsurerUji{baris: map[string]models.ReinsurerTreaty{}},
		security:  gudangSecurityKosong(), business: gudangBusinessKosong(), klausul: gudangKlausulKosong(),
		tx: &transaksiRekamUji{}, penetap: &penetapUji{peta: map[string]string{}},
	}
	n := 0
	r.layanan = services.New(nil).SimpanUtuhTCO().
		DenganKontrak(layananKontrak(r.kontrak)).
		DenganReinsurer(layananReinsurer(r.reinsurer)).
		DenganSecurity(layananSecurity(r.security, r.reinsurer)).
		DenganBusiness(layananBusiness(r.business)).
		DenganKlausul(layananKlausul(r.klausul, &n)).
		DenganPenetapIdentitas(r.penetap.tetapkan, r.penetap.rekam).
		DenganTransaksi(r.tx.jalankan).DenganJam(jamUji)
	return r
}

func bundelUtuh() services.KontrakUtuhMasuk {
	m := kontrakMasuk()
	m.ID = "1000003"
	return services.KontrakUtuhMasuk{
		Kontrak: m,
		Reinsurer: []services.ReinsurerUtuhMasuk{{
			ReinsurerMasuk: services.ReinsurerMasuk{ReinsurerID: "UJI-R1", PctShare: "40", Ricomm: "10"},
			Security:       []services.SecurityMasuk{{ReasSecurity: "UJI-R2", PctShare: "25"}},
		}},
		Business: []services.BusinessMasuk{{BizCode: "UJI-B1", IsActive: "1"}},
		Klausul: []services.KlausulMasuk{epi("10003", "1000"), {DescID: "10009", Anak: true, ParentReinsTypeID: "10003",
			Medan: map[string]string{"ReinsTypeID": "10003", "Pct": "25"}}},
	}
}

// AC 37, 41, "commit sekali": seluruh penulis di SATU transaksi; jejak utuh;
// identitas ditetapkan sekali, di akhir; status "1".
func TestSimpanUtuhSeluruhnyaSatuTransaksi(t *testing.T) {
	r := rakitUtuh()
	h, err := r.layanan.Simpan(context.Background(), pelakuUjiTCO, "1000001", bundelUtuh())
	if err != nil {
		t.Fatal(err)
	}
	if r.tx.dipanggil != 1 || r.tx.gagal || r.penetap.dipanggil != 1 || !models.StatusSimpanSuksesTCO(&h.Status) {
		t.Errorf("transaksi %d gagal %v tetapkan %d status %q", r.tx.dipanggil, r.tx.gagal, r.penetap.dipanggil, h.Status)
	}
	if h.Kontrak.ID != "1000003" || len(h.Reinsurer) != 1 || len(h.Reinsurer[0].Security) != 1 ||
		h.Reinsurer[0].Security[0].ReasID != h.Reinsurer[0].Reinsurer.ID || len(h.Business) != 1 || len(h.Klausul) != 2 {
		t.Errorf("hasil: %+v", h)
	}
	if len(r.reinsurer.baris) != 1 || len(r.security.baris) != 1 || len(r.business.baris) != 1 || len(r.klausul.baris) != 2 ||
		r.kontrak.diperbaru != 1 {
		t.Errorf("tulisan: reas %d sec %d biz %d klausul %d kontrak %d", len(r.reinsurer.baris), len(r.security.baris),
			len(r.business.baris), len(r.klausul.baris), r.kontrak.diperbaru)
	}
	if len(r.penetap.jejak) != 1 || !strings.Contains(r.penetap.jejak[0],
		"T_TREATYCONTRACT|1000003|simpan|simpan utuh: 1 kontrak, 1 reinsurer, 1 security, 1 business, 2 klausul") {
		t.Errorf("jejak utuh: %v", r.penetap.jejak)
	}
}

// AC 37/38: kegagalan pada klausul ke-N keluar dari transaksi (-> dibatalkan),
// pesannya menyebut baris itu, pemetaan HTTP galat aslinya tetap, dan
// identitas TIDAK ditetapkan (nol nomor sequence terpakai).
func TestSimpanUtuhGagalPadaKlausulKeN(t *testing.T) {
	r := rakitUtuh()
	m := bundelUtuh()
	m.Klausul = append(m.Klausul[:1], services.KlausulMasuk{DescID: "10017", Medan: map[string]string{}}, m.Klausul[1])
	_, err := r.layanan.Simpan(context.Background(), pelakuUjiTCO, "1000001", m)
	var g services.GalatSimpanUtuhTCO
	if !errors.As(err, &g) || g.Bagian != "klausul ke-2" || !errors.Is(err, models.ErrKlausulDitahan) ||
		!strings.Contains(err.Error(), "dibatalkan seluruhnya") {
		t.Fatalf("galat: %v", err)
	}
	if !r.tx.gagal || r.penetap.dipanggil != 0 || len(r.penetap.jejak) != 0 {
		t.Errorf("transaksi tidak dibatalkan: gagal %v tetapkan %d jejak %v", r.tx.gagal, r.penetap.dipanggil, r.penetap.jejak)
	}
	// Kegagalan security menyebut reinsurer induknya.
	r2 := rakitUtuh()
	m2 := bundelUtuh()
	m2.Reinsurer[0].Security = append(m2.Reinsurer[0].Security, services.SecurityMasuk{ReasSecurity: "UJI-R2", PctShare: "1"})
	if _, err := r2.layanan.Simpan(context.Background(), pelakuUjiTCO, "1000001", m2); !errors.Is(err, services.ErrSecurityDobel) ||
		!strings.Contains(err.Error(), "security ke-2 pada reinsurer ke-1") {
		t.Errorf("security dobel: %v", err)
	}
}

// Identitas sementara tidak berarti bagi pemakai - disamarkan di pesan.
func TestGalatSimpanUtuhMenyamarkanIdentitasSementara(t *testing.T) {
	err := services.GalatSimpanUtuhTCO{Bagian: "klausul ke-3", Galat: services.GalatKlausulDobel{IDLain: "S1234567890000004T", Jenis: "EPI"}}
	if strings.Contains(err.Error(), "S1234567890000004T") || !strings.Contains(err.Error(), "(baris baru dalam permintaan ini)") ||
		!errors.Is(err, services.ErrKlausulDobel) {
		t.Errorf("%v", err)
	}
}

// Identitas tetap ditulis ke jawaban.
func TestSimpanUtuhMenulisIdentitasTetap(t *testing.T) {
	r := rakitUtuh()
	r.penetap.peta = map[string]string{"1000001": "1999991"} // identitas fake reinsurer pertama
	h, err := r.layanan.Simpan(context.Background(), pelakuUjiTCO, "1000001", bundelUtuh())
	if err != nil || h.Reinsurer[0].Reinsurer.ID != "1999991" || h.Reinsurer[0].Security[0].ReasID != "1999991" {
		t.Errorf("identitas tetap: %+v %v", h.Reinsurer, err)
	}
}

func TestSimpanUtuhGerbang(t *testing.T) {
	r := rakitUtuh()
	if _, err := r.layanan.Simpan(context.Background(), services.Pelaku{}, "1000001", bundelUtuh()); !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Errorf("identitas: %v", err)
	}
	besar := bundelUtuh()
	for i := 0; i < 500; i++ {
		besar.Business = append(besar.Business, services.BusinessMasuk{BizCode: "UJI-B1", IsActive: "1"})
	}
	if _, err := r.layanan.Simpan(context.Background(), pelakuUjiTCO, "1000001", besar); !errors.Is(err, services.ErrSimpanUtuhTidakSah) || r.tx.dipanggil != 0 {
		t.Errorf("batas baris: %v (transaksi %d)", err, r.tx.dipanggil)
	}
	if _, err := services.New(nil).SimpanUtuhTCO().DenganTransaksi(transaksiUji).Simpan(context.Background(), pelakuUjiTCO,
		"1000001", bundelUtuh()); err == nil {
		t.Error("bawaan tanpa penulis harus gagal terang")
	}
}

// Temuan /code-review: master klausul (kurs, jenis reasuransi, TREATYDESC)
// dibaca SEKALI per permintaan simpan utuh, bukan per baris.
type kursHitungUji struct{ panggil int }

func (k *kursHitungUji) Berlaku(ctx context.Context, tahun models.TahunTreaty) (models.KursTCO, error) {
	k.panggil++
	return kursKlausulUji{}.Berlaku(ctx, tahun)
}

func TestSimpanUtuhMembacaKursSekali(t *testing.T) {
	r := rakitUtuh()
	kurs := &kursHitungUji{}
	n := 0
	r.layanan = r.layanan.DenganKlausul(layananKlausul(r.klausul, &n).DenganKurs(kurs))
	m := bundelUtuh()
	m.Klausul = append(m.Klausul, services.KlausulMasuk{DescID: "10001", Medan: map[string]string{"ReinsTypeID": "10003", "Rp": "5"}})
	if _, err := r.layanan.Simpan(context.Background(), pelakuUjiTCO, "1000001", m); err != nil {
		t.Fatal(err)
	}
	if kurs.panggil != 1 {
		t.Errorf("kurs dibaca %d kali untuk 3 klausul berkurs, mau 1", kurs.panggil)
	}
}
