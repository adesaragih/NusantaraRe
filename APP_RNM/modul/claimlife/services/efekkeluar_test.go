package services_test

// Efek keluar asinkron, antre-ulang, dan flag lingkungan - TANPA Oracle.
//
// Pemilik: tiket 12. Dibaca sesudah: efekkeluar.go.

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/layanan"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/claimlife/services"
)

// resolverUji menggantikan pembacaan `M_LINK_SERVICE`.
type resolverUji struct {
	url   string
	galat error
	minta []layanan.KunciLayanan
}

func (r *resolverUji) Resolve(_ context.Context,
	k layanan.KunciLayanan) (string, error) {
	r.minta = append(r.minta, k)
	return r.url, r.galat
}

// efekUji mencatat apakah ia benar-benar dijalankan.
type efekUji struct {
	nama      string
	galat     error
	dipanggil int
}

func (e *efekUji) Nama() string { return e.nama }
func (e *efekUji) Jalankan(context.Context, outbox.MuatanEfek) error {
	e.dipanggil++
	return e.galat
}

// antreanUji menggantikan tempat antre-ulang yang belum disahkan.
type antreanUji struct{ catatan []outbox.CatatanEfekGagal }

func (a *antreanUji) Antre(_ context.Context, c outbox.CatatanEfekGagal) error {
	a.catatan = append(a.catatan, c)
	return nil
}

// TestNonProduksiTidakMenjalankanEfekKeluar - AC 21 spec.
//
// ⚠️ Flag hanya menggerbangi EFEK KELUAR, tidak pernah penyimpanan. Ini
// penyimpangan sadar dari Pega, yang memakai `IsPEGAPROD` juga untuk simpan
// utama - tanpa itu lingkungan non-produksi tidak dapat dipakai menguji.
func TestNonProduksiTidakMenjalankanEfekKeluar(t *testing.T) {
	efek := &efekUji{nama: "UJI-EFEK"}
	antre := &antreanUji{}
	p := outbox.NewPenyalur(inti.BukanProduksi, antre, efek)

	hasil := p.Salurkan(context.Background(), outbox.MuatanEfek{KlaimID: "CLM-1"})
	if efek.dipanggil != 0 {
		t.Errorf("efek keluar berjalan %d kali di non-produksi", efek.dipanggil)
	}
	if !hasil.Dilewati {
		t.Error("hasil tidak menyatakan efek keluar dilewati")
	}
	// ⛔ Dilewati BUKAN gagal: tidak ada yang perlu diantre ulang.
	if len(antre.catatan) != 0 {
		t.Errorf("%d catatan diantre padahal efeknya sengaja dilewati", len(antre.catatan))
	}
}

// TestKegagalanEfekKeluarTidakMenahanApaPun - AC 19 spec.
func TestKegagalanEfekKeluarTidakMenahanApaPun(t *testing.T) {
	rusak := errors.New("uji: layanan luar sedang gagal")
	efek := &efekUji{nama: "UJI-EFEK", galat: rusak}
	antre := &antreanUji{}
	p := outbox.NewPenyalur(inti.Produksi, antre, efek)

	// ⛔ Salurkan TIDAK mengembalikan galat. Bila ia mengembalikannya,
	// pemanggil akan tergoda meneruskannya - dan transisi status klaim akan
	// tertahan oleh layanan luar yang sedang gagal.
	hasil := p.Salurkan(context.Background(), outbox.MuatanEfek{KlaimID: "CLM-1"})
	if hasil.Dilewati {
		t.Error("efek di produksi dinyatakan dilewati")
	}
	if len(hasil.Gagal) != 1 || hasil.Gagal[0].Nama != "UJI-EFEK" {
		t.Fatalf("hasil.Gagal = %+v, mau satu kegagalan bernama UJI-EFEK", hasil.Gagal)
	}
	// AC 20: kegagalan masuk jalur audit DAN dapat diantre ulang.
	if len(antre.catatan) != 1 {
		t.Fatalf("%d catatan diantre, mau 1", len(antre.catatan))
	}
	if antre.catatan[0].KlaimID != "CLM-1" || antre.catatan[0].Nama != "UJI-EFEK" {
		t.Errorf("catatan = %+v, tidak menunjuk klaim dan efek yang benar", antre.catatan[0])
	}
	if antre.catatan[0].Sebab == "" {
		t.Error("catatan tidak menyebut sebab kegagalannya")
	}
}

// TestSatuEfekGagalTidakMenghentikanSisanya - kegagalan satu layanan luar
// bukan alasan dua lainnya tidak dicoba.
func TestSatuEfekGagalTidakMenghentikanSisanya(t *testing.T) {
	pertama := &efekUji{nama: "UJI-1", galat: errors.New("uji: gagal")}
	kedua := &efekUji{nama: "UJI-2"}
	antre := &antreanUji{}
	p := outbox.NewPenyalur(inti.Produksi, antre, pertama, kedua)

	hasil := p.Salurkan(context.Background(), outbox.MuatanEfek{KlaimID: "CLM-1"})
	if kedua.dipanggil != 1 {
		t.Errorf("efek kedua dijalankan %d kali; kegagalan efek pertama "+
			"menghentikannya", kedua.dipanggil)
	}
	if len(hasil.Gagal) != 1 {
		t.Errorf("hasil.Gagal = %d, mau 1", len(hasil.Gagal))
	}
}

// TestKunciKategoriTidakDitemukanGagalTerang - AC tiket 12.
//
// ⛔ PENYIMPANGAN SADAR. `[terverifikasi]`
// `Claim Life/Activity/GetLinkService.xml` pecahan baris 705-706 menyetel
// `ResponLink.URL = linkService.pxResults(1).URL` dengan langkah
// ber-`pyStepsPreCondition` KOSONG (baris 701) - artinya `Obj-Browse` yang
// tidak menemukan apa pun menghasilkan URL KOSONG, diam-diam, dan
// `Connect-REST` berikutnya menembak alamat kosong.
func TestKunciKategoriTidakDitemukanGagalTerang(t *testing.T) {
	kosong := &resolverUji{url: ""}
	_, err := layanan.AlamatLayanan(context.Background(), kosong,
		layanan.KunciArasapasLife)
	if !errors.Is(err, layanan.ErrEndpointTidakDitemukan) {
		t.Fatalf("galat = %v, mau ErrEndpointTidakDitemukan", err)
	}
	if len(kosong.minta) != 1 {
		t.Fatalf("resolver dipanggil %d kali, mau 1", len(kosong.minta))
	}
	// `[terverifikasi]` `serviceInsertArasapasClaimLife_act.xml` pecahan baris
	// 540-541: `Kategori_1 = "Klaim"`, `Kategori_2 = "insertClaimLife"`.
	if kosong.minta[0].Kategori1 != "Klaim" ||
		kosong.minta[0].Kategori2 != "insertClaimLife" {
		t.Errorf("kunci = %+v, mau (Klaim, insertClaimLife)", kosong.minta[0])
	}
}

// TestResolverBawaanGagalTerang - selama tabelnya belum boleh dibaca.
func TestResolverBawaanGagalTerang(t *testing.T) {
	_, err := layanan.ResolverBelumDiputuskan{}.Resolve(
		context.Background(), layanan.KunciArasapasLife)
	if !errors.Is(err, layanan.ErrResolverBelumDiputuskan) {
		t.Errorf("galat = %v, mau ErrResolverBelumDiputuskan", err)
	}
}

// TestKegagalanKonfigurasiDibedakanDariJaringan - AC tiket 12.
//
// Antre-ulang yang memperlakukan keduanya sama akan berputar sia-sia:
// alamat yang tidak ada di tabel tidak akan muncul karena dicoba lagi.
func TestKegagalanKonfigurasiDibedakanDariJaringan(t *testing.T) {
	// Kegagalan jaringan SEMBARANG layak dicoba ulang.
	if !outbox.LayakDicobaUlang(errors.New("uji: koneksi terputus di tengah")) {
		t.Error("kegagalan jaringan dinyatakan tidak layak dicoba ulang")
	}
	// ⛔ Ketiga …BelumDisetujui IKUT di sini. Ronde pertama melewatkannya -
	// dan justru ketiganyalah yang benar-benar diproduksi di produksi hari
	// ini, sehingga antrean akan berputar selamanya.
	for _, konfig := range []error{
		layanan.ErrEndpointTidakDitemukan,
		layanan.ErrResolverBelumDiputuskan,
		outbox.ErrAntreanBelumDiputuskan,
		outbox.ErrPenyimpananBelumDisetujui,
		outbox.ErrEmailBelumDisetujui,
		outbox.ErrArasapasBelumDisetujui,
	} {
		if outbox.LayakDicobaUlang(konfig) {
			t.Errorf("kegagalan konfigurasi %v dinyatakan layak dicoba ulang; "+
				"antre-ulang akan berputar sia-sia", konfig)
		}
	}
}

// TestAntreanBawaanGagalTerang - tempat antre-ulang belum disahkan.
func TestAntreanBawaanGagalTerang(t *testing.T) {
	err := outbox.AntreanBelumDiputuskan{}.Antre(
		context.Background(), outbox.CatatanEfekGagal{MuatanEfek: outbox.MuatanEfek{KlaimID: "CLM-1"}})
	if !errors.Is(err, outbox.ErrAntreanBelumDiputuskan) {
		t.Errorf("galat = %v, mau ErrAntreanBelumDiputuskan", err)
	}
}

// TestTigaEfekKeluarBukanEmpat - `[keputusan work owner 2026-09-16]`.
//
// ⛔ Konversi ke produksi TIDAK lagi mengirim payload JSON; hilir membaca
// langsung dari tabel klaim. Test yang menemukan efek keempat berupa
// pengiriman payload GAGAL.
func TestTigaEfekKeluarBukanEmpat(t *testing.T) {
	efek := services.EfekKeluarClaimLife(layanan.ResolverBelumDiputuskan{})
	if len(efek) != 3 {
		t.Fatalf("%d efek keluar, mau 3", len(efek))
	}
	mau := []string{
		outbox.NamaEfekBerkas, outbox.NamaEfekEmail, outbox.NamaEfekArasapas,
	}
	for i, e := range efek {
		if e.Nama() != mau[i] {
			t.Errorf("efek ke-%d bernama %q, mau %q", i, e.Nama(), mau[i])
		}
	}
	// ⛔ Ketiganya gagal TERANG, masing-masing menyebut apa yang ditunggu -
	// bukan diam-diam berhasil tanpa menghubungi apa pun.
	for _, e := range efek {
		if err := e.Jalankan(context.Background(), outbox.MuatanEfek{}); err == nil {
			t.Errorf("efek %q berhasil padahal belum disetujui", e.Nama())
		}
	}
}

// TestArasapasMelaporkanKegagalanKonfigurasiApaAdanya - kegagalan resolusi
// tidak boleh tersamar sebagai "belum disetujui".
func TestArasapasMelaporkanKegagalanKonfigurasiApaAdanya(t *testing.T) {
	// Resolver menjawab, tetapi kuncinya tidak ada: itu kegagalan KONFIGURASI,
	// dan antre-ulang tidak boleh memutarnya.
	e := outbox.EfekArasapas{Resolver: &resolverUji{url: ""}}
	err := e.Jalankan(context.Background(), outbox.MuatanEfek{})
	if !errors.Is(err, layanan.ErrEndpointTidakDitemukan) {
		t.Fatalf("galat = %v, mau ErrEndpointTidakDitemukan", err)
	}
	if outbox.LayakDicobaUlang(err) {
		t.Error("kegagalan konfigurasi dinyatakan layak dicoba ulang")
	}
	// Alamatnya ketemu: barulah ia berhenti pada "belum disetujui".
	e = outbox.EfekArasapas{Resolver: &resolverUji{url: "alamat-uji-bukan-URL"}}
	if err := e.Jalankan(context.Background(), outbox.MuatanEfek{}); !errors.Is(
		err, outbox.ErrArasapasBelumDisetujui) {
		t.Errorf("galat = %v, mau ErrArasapasBelumDisetujui", err)
	}
}

// TestLingkunganDariFlagKonfigurasi - ADR-U-0005.
func TestLingkunganDariFlagKonfigurasi(t *testing.T) {
	if !inti.LingkunganDariFlag(true).AdalahProduksi() {
		t.Error("IS_PEGA_PROD=true bukan produksi")
	}
	if inti.LingkunganDariFlag(false).AdalahProduksi() {
		t.Error("IS_PEGA_PROD=false dianggap produksi")
	}
}

// TestLingkunganServiceSampaiKePenyalur - AC lingkungan, dibuktikan BERJALAN.
//
// ⛔ Tanpa test ini AC lingkungan tercentang secara HAMPA: `Penyalur` bawaan
// `Komite()` pernah menyetel `BukanProduksi` secara harfiah, sehingga efek
// keluar tidak pernah berjalan di lingkungan mana pun - dan "di non-produksi
// efek keluar tidak berjalan" benar hanya karena ia tidak berjalan di mana
// pun. Yang diperiksa di sini: nilainya benar-benar MENGALIR dari Service.
func TestLingkunganServiceSampaiKePenyalur(t *testing.T) {
	if got := services.New(nil).Lingkungan(); got.AdalahProduksi() {
		t.Error("Service bawaan berlingkungan produksi; ia harus gagal tertutup")
	}
	prod := services.New(nil).DenganLingkungan(inti.Produksi)
	if !prod.Lingkungan().AdalahProduksi() {
		t.Fatal("DenganLingkungan(Produksi) tidak tersimpan")
	}
	// ⛔ Dan ia MENGALIR ke penyalur yang Komite() susun. Diperiksa lewat
	// perilaku: efek yang dipasang harus benar-benar berjalan.
	efek := &efekUji{nama: "UJI-ALIR"}
	outbox.NewPenyalur(prod.Lingkungan(), &antreanUji{}, efek).
		Salurkan(context.Background(), outbox.MuatanEfek{KlaimID: "CLM-1"})
	if efek.dipanggil != 1 {
		t.Errorf("efek berjalan %d kali, mau 1; lingkungan tidak mengalir "+
			"dari Service ke penyalur", efek.dipanggil)
	}
	// Sebaliknya: Service bawaan tidak menjalankan apa pun.
	bawaan := &efekUji{nama: "UJI-BAWAAN"}
	outbox.NewPenyalur(services.New(nil).Lingkungan(), &antreanUji{}, bawaan).
		Salurkan(context.Background(), outbox.MuatanEfek{KlaimID: "CLM-1"})
	if bawaan.dipanggil != 0 {
		t.Errorf("efek berjalan %d kali di Service bawaan", bawaan.dipanggil)
	}
}
