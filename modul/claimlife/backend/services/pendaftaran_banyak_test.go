package services

// Pendaftaran banyak case - satu case per PL number + mata uang (keputusan work owner 10-10-2026), TANPA Oracle:
// pembaca peserta, transaksi, pengenal, dan penyimpan diganti jahitan. Fixture UJI-.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/repository"
)

// sumberUji - baris sumber tiruan: PL -> sertifikat -> peserta (PL, polis, mata uang dari SUMBER).
type sumberUji map[string]map[string]models.Peserta

func pesertaUji(pl, polis, uang, sertifikat string) models.Peserta {
	return models.Peserta{SumberID: "UJI-SRC-" + pl + "-" + sertifikat, NomorPremiList: pl, NomorPolis: polis,
		NomorSertifikat: sertifikat, MataUang: uang, IdPega: "NBLF-" + pl}
}

func (s sumberUji) tambah(ps ...models.Peserta) sumberUji {
	for _, p := range ps {
		if s[p.NomorPremiList] == nil {
			s[p.NomorPremiList] = map[string]models.Peserta{}
		}
		s[p.NomorPremiList][p.NomorSertifikat] = p
	}
	return s
}

// catatanUji - apa saja yang terjadi pada "basis data" tiruan.
type catatanUji struct {
	transaksi   int // transaksi dibuka
	nomor       int // nomor diambil
	nomorDiTx   []int
	simpanDiTx  []int
	disimpan    []models.PohonKlaim // yang di-commit (transaksi berhasil)
	gagalSimpan int                 // simpan ke-n (1-based) gagal; 0 = tidak
	// aktif - jawaban tiruan PesertaAktifDiKlaimLain per PL (saringan status/cabut milik SQL, diuji di repository
	// dan uji -tags db); kueriAktif mencatat PL tiap kueri.
	aktif      map[string][]repository.PesertaAktifInApp
	kueriAktif []string
	// polis - jawaban tiruan Ringkas per KUNCI (nil = polisUji); kunciRingkas mencatat kunci tiap panggilan; kode
	// mencatat kode bisnis tiap nomor.
	polis        map[string]kontrak.PolisRingkas
	kunciRingkas []string
	kode         []string
}

// polisUji - header tiruan per IDPEGA (pengenal kasus versi polis): Type dan kode bisnis berbeda supaya case beda PL
// terbukti beda Type.
var polisUji = map[string]kontrak.PolisRingkas{
	"NBLF-UJI-PL-A": {NomorPolis: "", Status: kontrak.StatusWorkSelesai, Type: "QP", BusinessCode: "L1", ProdKe: 1},
	"NBLF-UJI-PL-B": {NomorPolis: "", Status: kontrak.StatusWorkSelesai, Type: "QR", BusinessCode: "L2", ProdKe: 1},
	"NBLF-UJI-PL-C": {NomorPolis: "", Status: kontrak.StatusWorkSelesai, Type: "QP", BusinessCode: "L1", ProdKe: 1},
}

// pembacaUji - pembaca polis tiruan: Ringkas(NO_POLIS) TIDAK boleh dipanggil claimlife lagi (dicatat sebagai "NO_POLIS:");
// RingkasDariIdPega mencatat kunci, tidak ditemukan = ErrPolisNomorTakDitemukan.
type pembacaUji struct{ c *catatanUji }

func (p pembacaUji) Ringkas(_ context.Context, kunci string) (kontrak.PolisRingkas, error) {
	p.c.kunciRingkas = append(p.c.kunciRingkas, "NO_POLIS:"+kunci)
	return kontrak.PolisRingkas{}, kontrak.ErrPolisNomorTakDitemukan
}

func (p pembacaUji) RingkasDariIdPega(_ context.Context, idPega string) (kontrak.PolisRingkas, error) {
	p.c.kunciRingkas = append(p.c.kunciRingkas, idPega)
	peta := p.c.polis
	if peta == nil {
		peta = polisUji
	}
	if v, ada := peta[idPega]; ada {
		return v, nil
	}
	return kontrak.PolisRingkas{}, kontrak.ErrPolisNomorTakDitemukan
}

type penomorCatat struct{ c *catatanUji }

func (p penomorCatat) NomorBerikut(_ context.Context, _ *db.Tx, kode string, _ time.Time) (string, error) {
	p.c.nomor++
	p.c.kode = append(p.c.kode, kode)
	p.c.nomorDiTx = append(p.c.nomorDiTx, p.c.transaksi)
	return fmt.Sprintf("UJI-CLMLF-%05d", p.c.nomor), nil
}

// pendaftaranUji merakit Pendaftaran berjahitan: transaksi tiruan meng-"commit" hanya bila fn berhasil.
func pendaftaranUji(sumber sumberUji, c *catatanUji) *Pendaftaran {
	pengenal := 0
	var tertunda []models.PohonKlaim
	simpan := 0
	return &Pendaftaran{jahit: &jahitanDaftar{
		bacaPeserta: func(_ context.Context, pl string, sertifikat []string) ([]models.Peserta, error) {
			var out []models.Peserta
			for _, no := range sertifikat {
				ps, ada := sumber[pl][no]
				if !ada {
					return nil, fmt.Errorf("%w: peserta sertifikat %q pada premium list %q tidak ada atau sudah batal/delete",
						repository.ErrPesertaKlaimTidakAda, no, pl)
				}
				out = append(out, ps)
			}
			return out, nil
		},
		// Fungsi kunci polis yang SESUNGGUHNYA (IDPEGA) atas pembaca tiruan.
		ringkas: func(ctx context.Context, idPega []string) (kontrak.PolisRingkas, string, error) {
			return RingkasPolisIdPega(ctx, pembacaUji{c}, idPega)
		},
		pesertaAktif: func(_ context.Context, pl string, sertifikat []string) ([]repository.PesertaAktifInApp, error) {
			c.kueriAktif = append(c.kueriAktif, pl)
			var out []repository.PesertaAktifInApp
			for _, b := range c.aktif[pl] {
				for _, s := range sertifikat {
					if s == b.Sertifikat {
						out = append(out, b)
					}
				}
			}
			return out, nil
		},
		transaksi: func(_ context.Context, fn func(tx *db.Tx) error) error {
			c.transaksi++
			tertunda = nil
			if err := fn(nil); err != nil {
				tertunda = nil // rollback
				return err
			}
			c.disimpan = append(c.disimpan, tertunda...)
			return nil
		},
		pengenal: func(context.Context, *db.Tx) (string, error) {
			pengenal++
			return fmt.Sprintf("UJI-CLMLF-W%d", pengenal), nil
		},
		simpan: func(_ context.Context, _ *db.Tx, p models.PohonKlaim) error {
			simpan++
			c.simpanDiTx = append(c.simpanDiTx, c.transaksi)
			if simpan == c.gagalSimpan {
				return errors.New("UJI: ORA-00001 tiruan")
			}
			tertunda = append(tertunda, p)
			return nil
		},
	}}
}

func mintaKelompok(k ...KelompokDaftar) PermintaanDaftar {
	return PermintaanDaftar{Kelompok: k}
}

func sumberDuaPL() sumberUji {
	return sumberUji{}.tambah(
		pesertaUji("UJI-PL-A", "UJI-POL-A", "IDR", "001"), pesertaUji("UJI-PL-A", "UJI-POL-A", "IDR", "002"),
		pesertaUji("UJI-PL-B", "UJI-POL-B", "IDR", "101"), pesertaUji("UJI-PL-B", "UJI-POL-B", "IDR", "102"),
		pesertaUji("UJI-PL-B", "UJI-POL-B", "IDR", "103"), pesertaUji("UJI-PL-B", "UJI-POL-B", "IDR", "104"),
		pesertaUji("UJI-PL-C", "UJI-POL-C", "IDR", "201"), pesertaUji("UJI-PL-C", "UJI-POL-C", "USD", "202"),
	)
}

// Contoh work owner: 2 peserta PL A + 4 peserta PL B = 2 case, satu transaksi nomor + satu transaksi penulisan.
func TestDuaPLDuaCase(t *testing.T) {
	c := &catatanUji{}
	hasil, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", NomorPolis: "UJI-POL-A", Sertifikat: []string{"001", "002"}},
		KelompokDaftar{NomorPremiList: "UJI-PL-B", NomorPolis: "UJI-POL-B", Sertifikat: []string{"101", "102", "103", "104"}},
	))
	if err != nil {
		t.Fatal(err)
	}
	if len(hasil) != 2 || len(c.disimpan) != 2 {
		t.Fatalf("case %d, tersimpan %d, mau 2", len(hasil), len(c.disimpan))
	}
	for i, mau := range []struct {
		pl, polis string
		n         int
	}{{"UJI-PL-A", "UJI-POL-A", 2}, {"UJI-PL-B", "UJI-POL-B", 4}} {
		k := hasil[i].Klaim
		if k.Peserta[0].NomorPremiList != mau.pl || k.NomorPolis != mau.polis || len(k.Peserta) != mau.n ||
			k.ClaimRetro.Currency != "IDR" || k.Type != polisUji["NBLF-"+mau.pl].Type {
			t.Errorf("case %d: PL %s polis %s peserta %d uang %s type %s", i, k.Peserta[0].NomorPremiList, k.NomorPolis,
				len(k.Peserta), k.ClaimRetro.Currency, k.Type)
		}
	}
	// RALAT 10-10-2026: Register TIDAK menomori - CLAIM_NO kosong sampai Save to RNM; pengenal CLMLF- berbeda.
	if hasil[0].Klaim.NomorKlaim != "" || hasil[1].Klaim.NomorKlaim != "" || hasil[0].Work.ID == hasil[1].Work.ID {
		t.Error("dua case berbagi nomor atau pengenal")
	}
	// Nomor: dua-duanya di transaksi PERTAMA; penulisan: dua-duanya di transaksi KEDUA (AC tiket 02).
	// SATU transaksi untuk semua case, nol nomor diambil.
	if c.transaksi != 1 || c.nomor != 0 || fmt.Sprint(c.simpanDiTx) != "[1 1]" {
		t.Errorf("transaksi %d, nomor %d, simpan di %v", c.transaksi, c.nomor, c.simpanDiTx)
	}
	// 10-10-2026: Type dan kode bisnis per case dari header IDPEGA pesertanya; tidak ada Ringkas(NO_POLIS).
	if hasil[0].Klaim.Type != "QP" || hasil[1].Klaim.Type != "QR" {
		t.Errorf("Type %s/%s", hasil[0].Klaim.Type, hasil[1].Klaim.Type)
	}
	// Bawaan Register 11-10-2026: Non-Catastrophe / Non-KPR.
	for _, h := range hasil {
		if h.Klaim.StsKatastrofe != NonKatastrofe || h.Klaim.IsKPR != NonKPR {
			t.Errorf("bawaan Register %q/%q, mau Non-Catastrophe/Non-KPR", h.Klaim.StsKatastrofe, h.Klaim.IsKPR)
		}
	}
	if fmt.Sprint(c.kunciRingkas) != "[NBLF-UJI-PL-A NBLF-UJI-PL-B]" {
		t.Errorf("header dibaca dengan %v, mau IDPEGA peserta", c.kunciRingkas)
	}
}

// Satu PL dengan dua mata uang = dua case (kunci = PL + mata uang SUMBER).
func TestSatuPLDuaMataUangDuaCase(t *testing.T) {
	c := &catatanUji{}
	hasil, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-C", MataUang: "IDR", Sertifikat: []string{"201", "202"}}))
	if err != nil {
		t.Fatal(err)
	}
	if len(hasil) != 2 || hasil[0].Klaim.ClaimRetro.Currency != "IDR" || hasil[1].Klaim.ClaimRetro.Currency != "USD" {
		t.Fatalf("case %+v", hasil)
	}
}

// Kelompok dan mata uang kiriman klien DIABAIKAN: satu PL dipecah klien jadi dua kelompok = satu case; mata uang
// klien yang salah kalah oleh sumber; dua PL dalam SATU kelompok klien (bentuk lama berisi sertifikat PL lain tidak
// mungkin - pembacaannya per PL) tetap dikelompokkan per PL sumber.
func TestKelompokKlienDiabaikan(t *testing.T) {
	c := &catatanUji{}
	hasil, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", MataUang: "USD", Sertifikat: []string{"001"}},
		KelompokDaftar{NomorPremiList: "UJI-PL-B", MataUang: "USD", Sertifikat: []string{"101"}},
		KelompokDaftar{NomorPremiList: "UJI-PL-A", MataUang: "EUR", Sertifikat: []string{"002"}},
	))
	if err != nil {
		t.Fatal(err)
	}
	if len(hasil) != 2 || len(hasil[0].Klaim.Peserta) != 2 || hasil[0].Klaim.ClaimRetro.Currency != "IDR" ||
		hasil[1].Klaim.Peserta[0].NomorPremiList != "UJI-PL-B" {
		t.Fatalf("pengelompokan klien dipercaya: %d case", len(hasil))
	}
}

// Semua atau tidak sama sekali (1): satu kelompok gagal DIVALIDASI = nol nomor, nol transaksi, nol case; pesannya
// menyebut PL dan sertifikat.
func TestSatuKelompokGagalNolCase(t *testing.T) {
	c := &catatanUji{}
	_, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"001"}},
		KelompokDaftar{NomorPremiList: "UJI-PL-B", Sertifikat: []string{"101", "999"}},
	))
	if !errors.Is(err, ErrKelompokTidakSah) || !strings.Contains(err.Error(), "999") || !strings.Contains(err.Error(), "UJI-PL-B") {
		t.Fatalf("galat %v", err)
	}
	if c.nomor != 0 || c.transaksi != 0 || len(c.disimpan) != 0 {
		t.Errorf("nomor %d, transaksi %d, tersimpan %d - mau nol semua", c.nomor, c.transaksi, len(c.disimpan))
	}
}

// Semua atau tidak sama sekali (2): case KEDUA gagal ditulis = transaksi penulisan dibatalkan, nol case tersimpan.
// ⚠️ Kedua nomor sudah terbakar (transaksi pertama commit) - konsekuensi yang dinyatakan di Daftar.
func TestPenulisanCaseKeduaGagalNolCase(t *testing.T) {
	c := &catatanUji{gagalSimpan: 2}
	_, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"001"}},
		KelompokDaftar{NomorPremiList: "UJI-PL-B", Sertifikat: []string{"101"}},
	))
	if err == nil || !strings.Contains(err.Error(), "UJI-PL-B") {
		t.Fatalf("galat %v", err)
	}
	if len(c.disimpan) != 0 {
		t.Errorf("%d case tersimpan padahal case kedua gagal", len(c.disimpan))
	}
	// RALAT 10-10-2026: tanpa penomoran di Register, tidak ada nomor yang hangus.
	if c.nomor != 0 {
		t.Errorf("nomor diambil %d, mau 0", c.nomor)
	}
}

// 10-10-2026: kunci polis = IDPEGA. Nomor polis kiriman klien tidak memblokir, dan satu PL berisi dua POLICY_NO ceding
// tetap SATU case (header klaim = POLICY_NO peserta pertama, arti kolom tidak berubah).
func TestNomorPolisKlienTidakMemblokir(t *testing.T) {
	c := &catatanUji{}
	hasil, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", NomorPolis: "UJI-POL-LAIN", Sertifikat: []string{"001"}}))
	if err != nil || len(hasil) != 1 || hasil[0].Klaim.NomorPolis != "UJI-POL-A" {
		t.Fatalf("nomor polis klien memblokir / dipakai: %v", err)
	}
	sumber := sumberDuaPL().tambah(pesertaUji("UJI-PL-A", "UJI-POL-X", "IDR", "003"))
	hasil, err = pendaftaranUji(sumber, &catatanUji{}).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"001", "003"}}))
	if err != nil || len(hasil) != 1 || len(hasil[0].Klaim.Peserta) != 2 || hasil[0].Klaim.NomorPolis != "UJI-POL-A" {
		t.Errorf("PL dua POLICY_NO: %v %d", err, len(hasil))
	}
}

// Peserta satu case ber-IDPEGA berbeda (sebagian versi NB, sebagian versi EDM): header versi TERBARU (PROD_KE) dipakai.
func TestCaseBedaIdPegaPakaiVersiTerbaru(t *testing.T) {
	sumber := sumberDuaPL()
	ps := pesertaUji("UJI-PL-A", "UJI-POL-A", "IDR", "002")
	ps.IdPega = "EDML-UJI-1"
	sumber.tambah(ps)
	c := &catatanUji{polis: map[string]kontrak.PolisRingkas{
		"NBLF-UJI-PL-A": {Status: kontrak.StatusWorkSelesai, Type: "QP", BusinessCode: "L1", ProdKe: 1},
		// EDM terkonfirmasi: STATUSS Resolved-Completed (edm_putusan.go) - boleh.
		"EDML-UJI-1": {Status: kontrak.StatusWorkSelesai, Type: "QR", BusinessCode: "L9", ProdKe: 2},
	}}
	hasil, err := pendaftaranUji(sumber, c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"001", "002"}}))
	if err != nil || len(hasil) != 1 || hasil[0].Klaim.Type != "QR" {
		t.Fatalf("versi terbaru: %v %+v", err, hasil)
	}
	for _, k := range c.kunciRingkas {
		if strings.HasPrefix(k, "NO_POLIS:") {
			t.Errorf("claimlife masih memanggil Ringkas(NO_POLIS): %v", c.kunciRingkas)
		}
	}
}

// Header IDPEGA tidak ditemukan, tanpa Type / kode bisnis, atau peserta tanpa IDPEGA = 422 berpesan, nol nomor, nol case;
// semua case diperiksa sekaligus.
func TestPolisTanpaTypeNolNomor(t *testing.T) {
	c := &catatanUji{polis: map[string]kontrak.PolisRingkas{"NBLF-UJI-PL-A": {Type: "", BusinessCode: "L1"}}}
	_, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"001"}},
		KelompokDaftar{NomorPremiList: "UJI-PL-B", Sertifikat: []string{"101"}}))
	if !errors.Is(err, ErrKelompokTidakSah) || !strings.Contains(err.Error(), "(NBLF-UJI-PL-A, PL UJI-PL-A) has no Type") ||
		!strings.Contains(err.Error(), "Policy data not found for this insured (NBLF-UJI-PL-B).") {
		t.Fatalf("galat %v", err)
	}
	if c.nomor != 0 || c.transaksi != 0 || len(c.disimpan) != 0 {
		t.Errorf("nomor %d, transaksi %d, tersimpan %d", c.nomor, c.transaksi, len(c.disimpan))
	}
	sumber := sumberDuaPL()
	tanpa := pesertaUji("UJI-PL-A", "UJI-POL-A", "IDR", "009")
	tanpa.IdPega = ""
	sumber.tambah(tanpa)
	c2 := &catatanUji{}
	_, err = pendaftaranUji(sumber, c2).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"009"}}))
	if !errors.Is(err, ErrKelompokTidakSah) || !strings.Contains(err.Error(), "Insured Certificate No \"009\" (PL UJI-PL-A) has no IDPEGA") || c2.nomor != 0 {
		t.Errorf("peserta tanpa IDPEGA: %v", err)
	}
}

// RingkasPolisIdPega: pembaca tanpa kemampuan IDPEGA gagal terang (tidak jatuh ke NO_POLIS); daftar kosong = tidak ditemukan.
func TestRingkasPolisIdPegaAturan(t *testing.T) {
	c := &catatanUji{}
	if _, _, err := RingkasPolisIdPega(context.Background(), hanyaNoPolis{}, []string{"NBLF-UJI-1"}); !errors.Is(err, ErrPembacaTanpaIdPega) {
		t.Errorf("pembaca tanpa IDPEGA: %v", err)
	}
	if _, _, err := RingkasPolisIdPega(context.Background(), pembacaUji{c}, nil); !errors.Is(err, kontrak.ErrPolisNomorTakDitemukan) {
		t.Errorf("tanpa IDPEGA: %v", err)
	}
	got := IdPegaPeserta([]models.Peserta{{IdPega: " "}, {IdPega: "NBLF-1"}, {IdPega: "NBLF-1"}, {IdPega: "EDML-2"}})
	if fmt.Sprint(got) != "[NBLF-1 EDML-2]" {
		t.Errorf("IdPegaPeserta %v", got)
	}
}

// hanyaNoPolis - pembaca yang hanya memenuhi kontrak.PembacaPolis lama.
type hanyaNoPolis struct{}

func (hanyaNoPolis) Ringkas(context.Context, string) (kontrak.PolisRingkas, error) {
	return kontrak.PolisRingkas{Type: "QP"}, nil
}

// Sertifikat dobel di dalam permintaan ditolak SEBELUM pembacaan (400), juga di dua kelompok.
func TestSertifikatDobelDalamPermintaan(t *testing.T) {
	c := &catatanUji{}
	_, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"001"}},
		KelompokDaftar{NomorPremiList: " UJI-PL-A ", Sertifikat: []string{"001 "}}))
	if !errors.Is(err, galat.ErrPermintaanTidakSah) || !strings.Contains(err.Error(), "001") {
		t.Fatalf("galat %v", err)
	}
	if c.nomor != 0 || c.transaksi != 0 {
		t.Error("nomor / transaksi berjalan")
	}
}

// Bentuk lama (medan tunggal, tanpa Kelompok) tetap jalan = satu kelompok.
func TestBentukLamaTetapSatuCase(t *testing.T) {
	c := &catatanUji{}
	hasil, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), PermintaanDaftar{
		NomorPremiList: "UJI-PL-B", NomorPolis: "UJI-POL-B", MataUang: "IDR",
		Sertifikat: []string{"101", "102"}})
	if err != nil || len(hasil) != 1 || len(hasil[0].Klaim.Peserta) != 2 || hasil[0].Klaim.NomorPolis != "UJI-POL-B" {
		t.Fatalf("bentuk lama %v %d", err, len(hasil))
	}
}

// KelompokkanPeserta menolak hasil baca yang tidak berpasangan dengan kelompoknya (galat pemanggil, bukan data).
func TestKelompokkanPesertaPanjangTakSama(t *testing.T) {
	isi := KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"001"}}
	if _, err := KelompokkanPeserta([]KelompokDaftar{isi}, nil); err == nil {
		t.Error("jumlah hasil baca tidak sama dengan kelompok harus galat")
	}
}

// Aturan in-app Register (keputusan work owner 10-10-2026): peserta aktif di klaim lain = SELURUH pendaftaran batal
// (multi-PL ikut batal), nol nomor, nol case, 409; SEMUA bentrok disebut, Outstanding / Accepted.
func TestPesertaAktifDiKlaimLainMembatalkanSemua(t *testing.T) {
	c := &catatanUji{aktif: map[string][]repository.PesertaAktifInApp{
		"UJI-PL-B": {
			{NomorPremiList: "UJI-PL-B", Sertifikat: "101", KlaimID: "UJI-CLMLF-W7", NomorKlaim: "UJI-CLMLF-00007", KodeStatus: ""},
			{NomorPremiList: "UJI-PL-B", Sertifikat: "103", KlaimID: "UJI-CLMLF-W8", NomorKlaim: "", KodeStatus: kontrak.KodeAksep},
		},
	}}
	_, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"001"}},
		KelompokDaftar{NomorPremiList: "UJI-PL-B", Sertifikat: []string{"101", "102", "103"}},
	))
	var g *GalatPesertaSudahDiklaim
	if !errors.Is(err, ErrPesertaSudahDiklaim) || !errors.As(err, &g) || len(g.Bentrok) != 2 {
		t.Fatalf("galat %v", err)
	}
	mau := `Certificate No "101" (PL UJI-PL-B) is already in claim UJI-CLMLF-00007 (Outstanding); ` +
		`Certificate No "103" (PL UJI-PL-B) is already in claim UJI-CLMLF-W8 (Accepted)`
	if err.Error() != mau {
		t.Errorf("pesan:\n%s\nmau:\n%s", err.Error(), mau)
	}
	if c.nomor != 0 || c.transaksi != 0 || len(c.disimpan) != 0 {
		t.Errorf("nomor %d, transaksi %d, tersimpan %d - mau nol semua", c.nomor, c.transaksi, len(c.disimpan))
	}
}

// Satu kueri per PL: dua kelompok klien ber-PL sama digabung; tanpa bentrok pendaftaran berjalan.
func TestPesertaAktifSatuKueriPerPL(t *testing.T) {
	c := &catatanUji{}
	hasil, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
		KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"001"}},
		KelompokDaftar{NomorPremiList: "UJI-PL-B", Sertifikat: []string{"101"}},
		KelompokDaftar{NomorPremiList: " UJI-PL-A", Sertifikat: []string{"002"}},
	))
	if err != nil || len(hasil) != 2 {
		t.Fatalf("%v %d", err, len(hasil))
	}
	if fmt.Sprint(c.kueriAktif) != "[UJI-PL-A UJI-PL-B]" {
		t.Errorf("kueri per PL: %v", c.kueriAktif)
	}
}

// Keputusan work owner 10-10-2026 "polis harus sudah resolve completed": status efektif header bukan Resolved-Completed
// (Pending / Accept / Decline / Closed / kosong) = 422 berpesan PL + status, nol nomor; multi-PL dengan satu PL belum
// selesai = semua batal.
func TestPolisBelumResolvedNolNomor(t *testing.T) {
	for _, status := range []string{"Pending", "Accept", "Decline", "Closed", ""} {
		polis := map[string]kontrak.PolisRingkas{}
		for k, v := range polisUji {
			polis[k] = v
		}
		belum := polis["NBLF-UJI-PL-B"]
		belum.Status = status
		polis["NBLF-UJI-PL-B"] = belum
		c := &catatanUji{polis: polis}
		_, err := pendaftaranUji(sumberDuaPL(), c).Daftar(context.Background(), pelakuUji(), mintaKelompok(
			KelompokDaftar{NomorPremiList: "UJI-PL-A", Sertifikat: []string{"001"}},
			KelompokDaftar{NomorPremiList: "UJI-PL-B", Sertifikat: []string{"101"}}))
		mau := "Policy for PL UJI-PL-B is not Resolved-Completed yet (status: " + status + ")."
		if !errors.Is(err, ErrKelompokTidakSah) || !strings.Contains(err.Error(), mau) || strings.Contains(err.Error(), "PL UJI-PL-A is not") {
			t.Errorf("%q: %v", status, err)
		}
		if c.nomor != 0 || c.transaksi != 0 || len(c.disimpan) != 0 {
			t.Errorf("%q: nomor %d, transaksi %d, tersimpan %d", status, c.nomor, c.transaksi, len(c.disimpan))
		}
	}
}
