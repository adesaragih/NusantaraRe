package services_test

// Uji penyimpanan lampiran - tiket 12 (AC 59, 60, 62) + penjaga statik.

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/services"
)

// --- stub lokal ------------------------------------------------------------

func TestPenyimpananLokalMenimpaKunciYangSama(t *testing.T) {
	folder := t.TempDir()
	p := services.PenyimpananLokalDi(folder)
	ctx := context.Background()
	for _, isi := range []string{"PERTAMA", "KEDUA"} {
		if err := p.Simpan(ctx, "ABCDEF0123456789", strings.NewReader(isi), "application/pdf"); err != nil {
			t.Fatal(err)
		}
	}
	entri, _ := os.ReadDir(folder)
	if len(entri) != 1 {
		t.Fatalf("folder berisi %d berkas, mau 1 (kunci sama ditimpa)", len(entri))
	}
	rc, err := p.Buka(ctx, "ABCDEF0123456789")
	if err != nil {
		t.Fatal(err)
	}
	d, _ := io.ReadAll(rc)
	_ = rc.Close()
	if string(d) != "KEDUA" {
		t.Errorf("isi %q", d)
	}
	if err := p.Hapus(ctx, "ABCDEF0123456789"); err != nil {
		t.Fatal(err)
	}
	if err := p.Hapus(ctx, "ABCDEF0123456789"); !errors.Is(err, services.ErrBerkasTidakAdaDiPenyimpanan) {
		t.Errorf("hapus kedua: %v", err)
	}
	if _, err := p.Buka(ctx, "ABCDEF0123456789"); !errors.Is(err, services.ErrBerkasTidakAdaDiPenyimpanan) {
		t.Errorf("buka sesudah hapus: %v", err)
	}
	if ada, err := p.Ada(ctx, "ABCDEF0123456789"); ada || err != nil {
		t.Errorf("ada sesudah hapus: %v %v", ada, err)
	}
}

func TestPenyimpananLokalMenolakKunciBerjalur(t *testing.T) {
	p := services.PenyimpananLokalDi(t.TempDir())
	for _, k := range []string{"../../keluar", `..\x`, "a/b", "", "PENDEK"} {
		if err := p.Simpan(context.Background(), k, strings.NewReader("x"), ""); !errors.Is(err, services.ErrPermintaanTidakSah) {
			t.Errorf("kunci %q: %v", k, err)
		}
	}
}

func TestPenyimpananLokalTanpaFolderGagalTerang(t *testing.T) {
	p := services.PenyimpananLokalTCO(services.New(nil))
	if err := p.Simpan(context.Background(), "ABCDEF0123456789", strings.NewReader("x"), ""); !errors.Is(err, services.ErrUnggahanDirBelumDisetel) {
		t.Errorf("tanpa UNGGAHAN_DIR: %v", err)
	}
}

// --- rangkaian jarak jauh --------------------------------------------------

// resolverLampiranUji meniru `M_LINK_SERVICE`: isinya dapat diganti di tengah uji.
type resolverLampiranUji struct {
	alamat  map[services.KunciLayanan]string
	diminta []services.KunciLayanan
}

func (r *resolverLampiranUji) Resolve(_ context.Context, k services.KunciLayanan) (string, error) {
	r.diminta = append(r.diminta, k)
	a, ada := r.alamat[k]
	if !ada {
		return "", services.ErrEndpointTidakDitemukan
	}
	return a, nil
}

type pengirimLampiranUji struct {
	alamat []string
	token  []string
}

func (p *pengirimLampiranUji) catat(alamat, token string) {
	p.alamat = append(p.alamat, alamat)
	p.token = append(p.token, token)
}
func (p *pengirimLampiranUji) Kirim(_ context.Context, alamat, token, _ string, isi io.Reader, _ string) error {
	p.catat(alamat, token)
	_, err := io.Copy(io.Discard, isi)
	return err
}
func (p *pengirimLampiranUji) Ambil(_ context.Context, alamat, token, _ string) (io.ReadCloser, error) {
	p.catat(alamat, token)
	return io.NopCloser(strings.NewReader("ISI")), nil
}
func (p *pengirimLampiranUji) Buang(_ context.Context, alamat, token, _ string) error {
	p.catat(alamat, token)
	return nil
}
func (p *pengirimLampiranUji) Periksa(_ context.Context, alamat, token, _ string) (bool, error) {
	p.catat(alamat, token)
	return true, nil
}

type sumberTokenLampiranUji struct {
	terbit int
	gagal  error
	jam    func() time.Time
}

func (s *sumberTokenLampiranUji) TokenBaru(context.Context) (string, time.Time, error) {
	if s.gagal != nil {
		return "", time.Time{}, s.gagal
	}
	s.terbit++
	return "TOKEN-UJI-" + string(rune('A'+s.terbit-1)), s.jam().Add(time.Minute), nil
}

// AC 59: alamat datang dari resolver SAAT JALAN, per operasi dengan kuncinya
// sendiri, dan penggantian isi tabel berlaku tanpa menyusun ulang apa pun.
func TestPenyimpananJarakJauhAlamatDariResolverSaatJalan(t *testing.T) {
	saat := time.Date(2026, 9, 29, 9, 0, 0, 0, time.UTC)
	jam := func() time.Time { return saat }
	r := &resolverLampiranUji{alamat: map[services.KunciLayanan]string{
		services.KunciUnggahBerkas: "alamat-uji-unggah-1",
		services.KunciURLBerkas:    "alamat-uji-url",
		services.KunciHapusBerkas:  "alamat-uji-hapus",
	}}
	kirim := &pengirimLampiranUji{}
	p := services.NewPenyimpananJarakJauhTCO(r,
		services.NewCacheTokenTCO(&sumberTokenLampiranUji{jam: jam}, jam, services.MarginTokenTCO), kirim)
	ctx := context.Background()
	if err := p.Simpan(ctx, "ABCDEF0123456789", strings.NewReader("x"), ""); err != nil {
		t.Fatal(err)
	}
	r.alamat[services.KunciUnggahBerkas] = "alamat-uji-unggah-2"
	if err := p.Simpan(ctx, "ABCDEF0123456789", strings.NewReader("x"), ""); err != nil {
		t.Fatal(err)
	}
	if err := p.Hapus(ctx, "ABCDEF0123456789"); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Ada(ctx, "ABCDEF0123456789"); err != nil {
		t.Fatal(err)
	}
	mau := []string{"alamat-uji-unggah-1", "alamat-uji-unggah-2", "alamat-uji-hapus", "alamat-uji-url"}
	if strings.Join(kirim.alamat, ",") != strings.Join(mau, ",") {
		t.Errorf("alamat terpakai %v, mau %v", kirim.alamat, mau)
	}
	// Kunci yang belum ada di tabel: gagal permanen, tidak berputar.
	delete(r.alamat, services.KunciURLBerkas)
	_, err := p.Buka(ctx, "ABCDEF0123456789")
	if !errors.Is(err, services.ErrEndpointTidakDitemukan) || services.LayakDicobaUlang(err) {
		t.Errorf("kunci tidak ada: %v", err)
	}
}

func TestPenyimpananJarakJauhTanpaPengirimBelumDisetujui(t *testing.T) {
	p := services.NewPenyimpananJarakJauhTCO(&resolverLampiranUji{}, nil, nil)
	err := p.Simpan(context.Background(), "ABCDEF0123456789", strings.NewReader("x"), "")
	if !errors.Is(err, services.ErrPenyimpananBelumDisetujui) || services.LayakDicobaUlang(err) {
		t.Errorf("tanpa pengirim: %v", err)
	}
}

// AC 60: token dipakai ulang, diperbarui SEBELUM kedaluwarsa, dan kegagalan
// mengambilnya terlihat - transport tidak dipanggil sama sekali.
func TestCacheTokenDiperbaruiSebelumKedaluwarsa(t *testing.T) {
	saat := time.Date(2026, 9, 29, 9, 0, 0, 0, time.UTC)
	jam := func() time.Time { return saat }
	sumber := &sumberTokenLampiranUji{jam: jam}
	c := services.NewCacheTokenTCO(sumber, jam, 15*time.Second)
	ctx := context.Background()
	t1, _ := c.Token(ctx)
	saat = saat.Add(30 * time.Second)
	t2, _ := c.Token(ctx)
	if t1 != t2 || sumber.terbit != 1 {
		t.Errorf("token tidak dipakai ulang: %q %q, terbit %d", t1, t2, sumber.terbit)
	}
	saat = saat.Add(20 * time.Second) // 50 detik: sisa 10 detik < margin 15 detik
	t3, _ := c.Token(ctx)
	if t3 == t1 || sumber.terbit != 2 {
		t.Errorf("token tidak diperbarui sebelum kedaluwarsa: %q, terbit %d", t3, sumber.terbit)
	}
}

func TestTokenGagalTerlihatDanTransportTidakDipanggil(t *testing.T) {
	jam := func() time.Time { return time.Date(2026, 9, 29, 9, 0, 0, 0, time.UTC) }
	kirim := &pengirimLampiranUji{}
	p := services.NewPenyimpananJarakJauhTCO(
		&resolverLampiranUji{alamat: map[services.KunciLayanan]string{services.KunciUnggahBerkas: "alamat-uji"}},
		services.NewCacheTokenTCO(&sumberTokenLampiranUji{gagal: errors.New("GET_TOKEN_STORAGE menolak"), jam: jam},
			jam, services.MarginTokenTCO), kirim)
	err := p.Simpan(context.Background(), "ABCDEF0123456789", strings.NewReader("x"), "")
	if !errors.Is(err, services.ErrTokenPenyimpananGagal) {
		t.Errorf("token gagal: %v", err)
	}
	if len(kirim.alamat) != 0 {
		t.Error("transport dipanggil padahal token gagal diambil")
	}
}

// --- penjaga statik --------------------------------------------------------

// berkasProduksiTCO mengumpulkan berkas Go produksi modul ini.
func berkasProduksiTCO(t *testing.T) map[string]string {
	t.Helper()
	hasil := map[string]string{}
	for _, pola := range []string{"tco_*.go", "../repository/tco_*.go", "../handlers/tco_*.go",
		"../handlers/rute_treaty_contract_out.go", "../models/tco_*.go"} {
		cocok, err := filepath.Glob(pola)
		if err != nil {
			t.Fatal(err)
		}
		for _, j := range cocok {
			if strings.HasSuffix(j, "_test.go") {
				continue
			}
			isi, err := os.ReadFile(j)
			if err != nil {
				t.Fatal(err)
			}
			hasil[j] = string(isi)
		}
	}
	if len(hasil) < 15 {
		t.Fatalf("hanya %d berkas modul terbaca; pembacanya yang rusak", len(hasil))
	}
	return hasil
}

// buangKomentarBaris membuang komentar `//` sampai akhir baris.
func buangKomentarBaris(isi string) string {
	var b strings.Builder
	for _, baris := range strings.Split(isi, "\n") {
		if i := strings.Index(baris, "//"); i >= 0 {
			baris = baris[:i]
		}
		b.WriteString(baris)
		b.WriteByte('\n')
	}
	return b.String()
}

// AC 59: nol alamat sebagai literal, konstanta, atau pembacaan env var, dan
// nol klien HTTP, di seluruh berkas produksi modul ini.
func TestTCOLampiranTanpaAlamatLiteral(t *testing.T) {
	terlarang := map[string]*regexp.Regexp{
		// ⚠️ Pola dirakit dari potongan: literalnya sendiri akan ditangkap
		// `TestNolAlamatLayananDiKode`, yang hanya mengecualikan berkasnya sendiri.
		"URL literal":  regexp.MustCompile(regexp.QuoteMeta(":" + "/" + "/")),
		"env var":      regexp.MustCompile(`os\.(Getenv|LookupEnv|Environ)\(`),
		"klien HTTP":   regexp.MustCompile(`http\.(Post|Get|Head|NewRequest|DefaultClient)|http\.Client\{|\.Do\(req`),
		"host literal": regexp.MustCompile(`(?i)storage\.googleapis|googleapis\.com`),
	}
	// ⚠️ SATU berkas boleh memegang klien HTTP - transport OQ-TCO-08
	// [keputusan work owner 29-09-2026]. Pola lain tetap berlaku untuknya.
	const transportDisetujui = "tco_pengirim_storage.go"
	transportTerlihat := false
	for nama, isi := range berkasProduksiTCO(t) {
		for jenis, pola := range terlarang {
			if !pola.MatchString(isi) {
				continue
			}
			if jenis == "klien HTTP" && filepath.Base(nama) == transportDisetujui {
				transportTerlihat = true
				continue
			}
			t.Errorf("%s memuat %s - alamat di-resolve saat jalan dari M_LINK_SERVICE (ADR-0013)", nama, jenis)
		}
	}
	if !transportTerlihat {
		t.Errorf("pengecualian %s tidak terpakai - hapus dari penjaga ini", transportDisetujui)
	}
}

// Lampiran melekat pada TAHUN treaty: kode yang merujuk kunci treaty inward
// atau tabel lampiran warisannya gagal.
func TestTCOLampiranTanpaRujukanTreatyInward(t *testing.T) {
	pola := regexp.MustCompile(`(?i)treatyid\b|treaty_in\b|treatyin\b|m_attachmenttreaty`)
	for nama, isi := range berkasProduksiTCO(t) {
		if !strings.Contains(nama, "lampiran") && !strings.Contains(nama, "penyimpanan") &&
			!strings.Contains(nama, "rute_treaty_contract_out") {
			continue
		}
		if m := pola.FindString(buangKomentarBaris(isi)); m != "" {
			t.Errorf("%s merujuk %q - lampiran melekat pada tahun treaty, bukan treaty inward", nama, m)
		}
	}
}
