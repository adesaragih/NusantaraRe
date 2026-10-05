package services_test

// Aturan Company Detail TANPA Oracle (gudang tiruan). Sumber setiap aturan: dokumentasi paket services.

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/companydetail/backend/models"
	"nusantarare/modul/companydetail/backend/services"
	"nusantarare/modul/companydetail/backend/tiruan"
)

var (
	ctx   = context.Background()
	admin = inti.Pelaku{AkunID: "UJI-ADMIN"}
	jam   = time.Date(2026, 10, 4, 2, 3, 4, 500e6, time.UTC)
)

const anak = "ASM-SFAGIS-WORK-ORG ORG-115"

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi).DenganJam(func() time.Time { return jam })
}

func isianBaru() models.Isian {
	return models.Isian{Nama: " UJI Organisasi Baru ", Title: "PT.", NPWP: "UJI-NPWP-1", Country: "001",
		BusinessField: "01", Note: "UJI catatan",
		PIC: []models.PIC{{Nama: "UJI Kontak A", Position: "Direktur", Gender: "1", Email: "uji.a@contoh.id",
			DateOfBirth: "15-01-1990", Phone: "0812"}, {Nama: "UJI Kontak B", Position: "Staff"}},
		Alamat: []models.Alamat{{Jenis: "2", Address: "UJI Jalan Baru", Telfax: []models.Telfax{
			{Jenis: "5", Code: "021", No: "UJI-TELP-1"}, {Jenis: "3", No: "UJI-HP-1"}}},
			{Jenis: "1", Address: "UJI Jalan Rumah"}}}
}

// Create: nomor ORG dari SEQ_CLIENT_ORG (tiruan mulai 121 = M_CLIENT 120 + 1, seperti migrasi 810), COUNTRYNAME dari
// NATION, PIC baru PIC-1/PIC-2, tanggal lahir disimpan YYYYMMDD, satu baris alamat per nomor dengan jejak pelaku.
func TestCreateOrganisasiBaru(t *testing.T) {
	g := tiruan.Contoh()
	d, err := layanan(g).Tambah(ctx, admin, isianBaru())
	if err != nil {
		t.Fatal(err)
	}
	if d.ID != "ASM-SFAGIS-WORK-ORG ORG-121" || d.IDView != "ORG-121" || g.SeqOrg != 122 {
		t.Errorf("nomor ORG %s %s (sequence berikut %d), mau ORG-121 dari sequence", d.ID, d.IDView, g.SeqOrg)
	}
	o := g.Disisip[0]
	if o.Nama != "UJI ORGANISASI BARU" || o.Title != "PT." || o.CountryName != "UJI NEGARA SATU" || o.BusinessField != "01" ||
		o.CreatedBy != "UJI-ADMIN" || o.UpdatedBy != "UJI-ADMIN" || o.Note != "UJI catatan" {
		t.Errorf("baris organisasi %+v", o)
	}
	pic := g.PIC[d.ID]
	if len(pic) != 2 || pic[0].UserIdentifier != "PIC-1" || pic[1].UserIdentifier != "PIC-2" || pic[0].DateOfBirth != "19900115" {
		t.Errorf("PIC tersimpan %+v", pic)
	}
	if d.PIC[0].DateOfBirth != "15-01-1990" {
		t.Errorf("tanggal lahir ke layar %q", d.PIC[0].DateOfBirth)
	}
	baris := g.Alamat[d.ID]
	if len(baris) != 3 {
		t.Fatalf("baris alamat %d, mau 3 (dua nomor + satu alamat tanpa nomor): %+v", len(baris), baris)
	}
	if baris[0].TelfaxType != "5" || baris[0].TelfaxCode != "021" || baris[1].TelfaxType != "3" || baris[2].TelfaxNo != "" ||
		baris[0].Address != baris[1].Address || baris[0].PxCreateOperator != "UJI-ADMIN" ||
		baris[0].PxCreateDateTime != "20261004T020304.500 GMT" {
		t.Errorf("baris alamat %+v", baris)
	}
	if len(d.Alamat) != 2 || len(d.Alamat[0].Telfax) != 2 || len(d.Alamat[1].Telfax) != 0 {
		t.Errorf("alamat ke layar %+v", d.Alamat)
	}
}

// Nomor ORG = NEXTVAL apa adanya, bukan nomor tertinggi + 1: dua Create berurutan 121 lalu 122, dan sequence yang
// sudah di depan (500) dipakai apa adanya walau nomor tertinggi di tabel lebih kecil.
func TestCreateNomorDariSequence(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	satu, err := l.Tambah(ctx, admin, isianBaru())
	if err != nil {
		t.Fatal(err)
	}
	isi := isianBaru()
	isi.Nama, isi.NPWP = "UJI Organisasi Kedua", "UJI-NPWP-2"
	dua, err := l.Tambah(ctx, admin, isi)
	if err != nil {
		t.Fatal(err)
	}
	if satu.IDView != "ORG-121" || dua.IDView != "ORG-122" {
		t.Errorf("berurutan: %s lalu %s, mau ORG-121 lalu ORG-122", satu.IDView, dua.IDView)
	}
	g = tiruan.Contoh()
	g.SeqOrg = 500
	d, err := layanan(g).Tambah(ctx, admin, isianBaru())
	if err != nil {
		t.Fatal(err)
	}
	if d.IDView != "ORG-500" {
		t.Errorf("sequence di depan: %s, mau ORG-500", d.IDView)
	}
}

// Negara tanpa OLDID disimpan dengan ID NATION-nya.
func TestCreateNegaraTanpaOldID(t *testing.T) {
	g := tiruan.Contoh()
	isi := isianBaru()
	isi.Country = "100902"
	d, err := layanan(g).Tambah(ctx, admin, isi)
	if err != nil {
		t.Fatal(err)
	}
	if d.Country != "100902" || d.CountryName != "UJI NEGARA BARU" {
		t.Errorf("country %q %q", d.Country, d.CountryName)
	}
}

func TestCreateDitolakTanpaTulisan(t *testing.T) {
	for _, k := range []struct {
		nama  string
		ubah  func(*models.Isian)
		pesan string
	}{
		{"tanpa nama", func(i *models.Isian) { i.Nama = " " }, "Organization Name is required"},
		{"title di nama", func(i *models.Isian) { i.Nama = "P.T. UJI Baru" }, "Organization Name must not contain the title PT."},
		{"title di belakang", func(i *models.Isian) { i.Nama = "UJI Baru, CV" }, "must not contain the title CV."},
		{"tanpa country", func(i *models.Isian) { i.Country = "" }, "COUNTRY is required"},
		{"country asing", func(i *models.Isian) { i.Country = "999" }, "COUNTRY 999 is not in the country list"},
		{"tanpa business field", func(i *models.Isian) { i.BusinessField = "" }, "Business Field is required"},
		{"business field asing", func(i *models.Isian) { i.BusinessField = "ZZ" }, "Business Field ZZ is not in"},
		{"title nonaktif", func(i *models.Isian) { i.Title = "TN." }, "Title TN. is not in the title list"},
		{"parent tidak ada", func(i *models.Isian) { i.ParentID = "UJI-TIDAK-ADA" }, "Parent organization UJI-TIDAK-ADA not found"},
		{"PIC tanpa nama", func(i *models.Isian) { i.PIC[1].Nama = "" }, "PIC row 2: Name is required"},
		{"PIC tanpa position", func(i *models.Isian) { i.PIC[0].Position = "" }, "PIC row 1: UJI Kontak A has no Job Position"},
		{"PIC email salah", func(i *models.Isian) { i.PIC[0].Email = "uji-tanpa-at" }, "PIC row 1: Email format is not valid"},
		{"PIC tanggal salah", func(i *models.Isian) { i.PIC[0].DateOfBirth = "1990-01-15" }, "Date of birth must be DD-MM-YYYY"},
		{"PIC gender asing", func(i *models.Isian) { i.PIC[0].Gender = "9" }, "Gender 9 is not in the gender list"},
		{"PIC asing", func(i *models.Isian) { i.PIC[0].UserIdentifier = "PIC-7" }, "PIC row 1 does not belong"},
		{"alamat tanpa type", func(i *models.Isian) { i.Alamat[1].Jenis = "" }, "Address row 2: Type is required"},
		{"type alamat nonaktif", func(i *models.Isian) { i.Alamat[1].Jenis = "3" }, "Type 3 is not in the address type list"},
		{"alamat kosong", func(i *models.Isian) { i.Alamat[0].Address = " " }, "Address row 1: Address is required"},
		{"alamat ganda", func(i *models.Isian) { i.Alamat[1].Address = "UJI Jalan Baru" }, "Address row 2 repeats the address of row 1"},
		{"telfax nonaktif", func(i *models.Isian) { i.Alamat[0].Telfax[0].Jenis = "2" }, "type 2 is not in the phone and fax type list"},
		{"telfax EMAIL dihapus", func(i *models.Isian) { i.Alamat[0].Telfax[0].Jenis = "6" }, "type 6 is not in the phone and fax type list"},
		{"telfax tanpa nomor", func(i *models.Isian) { i.Alamat[0].Telfax[1].No = "" }, "Phone and Fax 2: the number is required"},
		{"kode area Others bukan angka", func(i *models.Isian) { i.Alamat[0].Telfax[0].Code = "02A" }, "area code 02A must be digits"},
		{"kode area Others terlalu panjang", func(i *models.Isian) { i.Alamat[0].Telfax[0].Code = "+6202199999" }, "longer than 10 bytes"},
		{"asal alamat asing", func(i *models.Isian) { i.Alamat[0].Asal = "UJI Jalan Lain" }, "Address row 1 does not belong"},
	} {
		g := tiruan.Contoh()
		isi := isianBaru()
		k.ubah(&isi)
		_, err := layanan(g).Tambah(ctx, admin, isi)
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), k.pesan) {
			t.Errorf("%s: galat %v, mau %q", k.nama, err, k.pesan)
		}
		if len(g.Disisip) != 0 || g.SeqOrg != 121 {
			t.Errorf("%s: ditolak tetapi menulis atau mengambil nomor ORG", k.nama)
		}
	}
	if _, err := layanan(tiruan.Contoh()).Tambah(ctx, inti.Pelaku{}, isianBaru()); !errors.Is(err, services.ErrTanpaPelaku) {
		t.Errorf("tanpa pelaku: %v", err)
	}
}

// Ubah: PIC lama dikenali dari USERIDENTIFIER-nya (yang dibuang dari grid hilang), PIC baru = nomor tertinggi + 1;
// alamat yang diubah teksnya membawa kolom tersembunyi dari baris lamanya; jenis nomor lama (FAX) boleh tetap.
func TestUbahPICDanAlamat(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	d, err := l.Ambil(ctx, anak)
	if err != nil {
		t.Fatal(err)
	}
	if d.PIC[1].DateOfBirth != "15-01-1990" || len(d.Alamat) != 2 || len(d.Alamat[0].Telfax) != 2 || d.Alamat[0].Asal != "UJI Jalan Satu" {
		t.Fatalf("detail awal %+v", d)
	}
	isi := models.Isian{Nama: d.Nama, Title: d.Title, NPWP: d.NPWP, Country: d.Country, BusinessField: d.BusinessField,
		ParentID: d.ParentID, Note: "UJI ubah",
		PIC: []models.PIC{d.PIC[1], {Nama: "UJI Kontak Baru", Position: "Manager"}},
		Alamat: []models.Alamat{{Asal: "UJI Jalan Satu", Jenis: "2", Address: "UJI Jalan Satu Baru",
			Telfax: []models.Telfax{d.Alamat[0].Telfax[0]}}}}
	hasil, err := l.Ubah(ctx, inti.Pelaku{AkunID: "UJI-EDITOR"}, anak, isi)
	if err != nil {
		t.Fatal(err)
	}
	ids := []string{}
	for _, p := range g.PIC[anak] {
		ids = append(ids, p.UserIdentifier)
	}
	if !reflect.DeepEqual(ids, []string{"PIC-2", "PIC-3"}) || g.PIC[anak][0].DateOfBirth != "19900115" {
		t.Errorf("PIC sesudah ubah %+v", g.PIC[anak])
	}
	baris := g.Alamat[anak]
	if len(baris) != 1 || baris[0].Address != "UJI Jalan Satu Baru" || baris[0].City != "UJIKOTA" || baris[0].ZipCode != "10000" ||
		baris[0].PxCreateOperator != "UJI-PEGA" || baris[0].TelfaxType != "2" {
		t.Errorf("baris alamat sesudah ubah %+v", baris)
	}
	o := g.Org[anak]
	if o.UpdatedBy != "UJI-EDITOR" || o.Note != "UJI ubah" || o.ParentName != "UJI Induk Grup" || hasil.IDView != "ORG-115" {
		t.Errorf("organisasi sesudah ubah %+v", o)
	}
}

// Parent organization: diganti = GROUPNAME disalin dari induk baru; diri sendiri ditolak; dikosongkan = keduanya kosong.
func TestUbahParent(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	isi := func(parent string) models.Isian {
		o := g.Org[anak]
		return models.Isian{Nama: o.Nama, Title: o.Title, Country: o.Country, BusinessField: o.BusinessField, ParentID: parent}
	}
	if _, err := l.Ubah(ctx, admin, anak, isi(anak)); err == nil || !strings.Contains(err.Error(), "cannot be its own parent") {
		t.Errorf("induk diri sendiri: %v", err)
	}
	if _, err := l.Ubah(ctx, admin, anak, isi("")); err != nil {
		t.Fatal(err)
	}
	if o := g.Org[anak]; o.ParentID != "" || o.ParentName != "" {
		t.Errorf("parent dikosongkan %+v", o)
	}
	if _, err := l.Ubah(ctx, admin, anak, isi("ASM-SFAGIS-WORK-ORG ORG-100")); err != nil {
		t.Fatal(err)
	}
	if o := g.Org[anak]; o.ParentName != "UJI Induk Grup" {
		t.Errorf("GROUPNAME %q", o.ParentName)
	}
	if _, err := l.Ubah(ctx, admin, "UJI-TIDAK-ADA", isi("")); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("organisasi tidak ada: %v", err)
	}
}

// Nilai lama Pega yang tidak lagi aktif (Title TN.) boleh tetap bila tidak diubah.
func TestUbahNilaiLamaTetapBoleh(t *testing.T) {
	g := tiruan.Contoh()
	o := g.Org[anak]
	o.Title = "TN."
	g.Org[anak] = o
	isi := models.Isian{Nama: o.Nama, Title: "TN.", Country: o.Country, BusinessField: o.BusinessField, ParentID: o.ParentID}
	if _, err := layanan(g).Ubah(ctx, admin, anak, isi); err != nil {
		t.Errorf("title lama tetap: %v", err)
	}
}

// Daftar: halaman, saringan, dan batas ukuran.
func TestDaftarHalamanDanCari(t *testing.T) {
	l := layanan(tiruan.Contoh())
	h, err := l.Daftar(ctx, "", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if h.Total != 2 || h.Halaman != 1 || h.Ukuran != services.UkuranHalaman || h.Daftar[0].Nama != "UJI Anak Usaha" {
		t.Errorf("daftar %+v", h)
	}
	if h, _ := l.Daftar(ctx, "npwp-115", 1, 1); h.Total != 1 || h.Daftar[0].IDView != "ORG-115" {
		t.Errorf("cari NPWP %+v", h)
	}
	if h, _ := l.Daftar(ctx, "", 2, 1); len(h.Daftar) != 1 || h.Daftar[0].IDView != "ORG-100" {
		t.Errorf("halaman 2 %+v", h)
	}
	if h, _ := l.Daftar(ctx, "", 1, 1000); h.Ukuran != services.UkuranMaksimum {
		t.Errorf("ukuran maksimum %d", h.Ukuran)
	}
}

// Pilihan dikelompokkan per jenis; negara dari NATION.
func TestPilihanPerJenis(t *testing.T) {
	p, err := layanan(tiruan.Contoh()).Pilihan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Title) != 3 || len(p.Telfax) != 4 || len(p.AddressType) != 3 || len(p.Gender) != 2 || len(p.Negara) != 2 {
		t.Errorf("pilihan %+v", p)
	}
}

func TestTanggalKabel(t *testing.T) {
	if s := services.TanggalKeKabel("19900115"); s != "15-01-1990" {
		t.Errorf("ke kabel %q", s)
	}
	if s := services.TanggalKeKabel("lain"); s != "lain" {
		t.Errorf("bentuk lain %q", s)
	}
	if s, ok := services.TanggalDariKabel("31-02-1990"); ok {
		t.Errorf("tanggal mustahil diterima: %q", s)
	}
}

// Pemeriksaan nama sebelum Create: title di nama, nama sama sesudah title dan tanda baca dibuang, nama mirip, dan
// organisasi yang sedang diubah tidak dihitung.
func TestPeriksaNama(t *testing.T) {
	l := layanan(tiruan.Contoh())
	h, err := l.Periksa(ctx, "P.T. UJI Anak-Usaha", "")
	if err != nil {
		t.Fatal(err)
	}
	if h.TitleDalamNama != "PT." || len(h.Serupa) == 0 || h.Serupa[0].IDView != "ORG-115" || h.Serupa[0].Jenis != services.NamaSama {
		t.Errorf("nama sama + title: %+v", h)
	}
	if h, _ := l.Periksa(ctx, "UJI Anak Usah", ""); len(h.Serupa) != 1 || h.Serupa[0].Jenis != services.NamaMirip || h.TitleDalamNama != "" {
		t.Errorf("salah ketik satu huruf: %+v", h)
	}
	if h, _ := l.Periksa(ctx, "UJI Induk Grup Indonesia", ""); len(h.Serupa) != 1 || h.Serupa[0].IDView != "ORG-100" {
		t.Errorf("nama memuat nama lain: %+v", h)
	}
	if h, _ := l.Periksa(ctx, "UJI Anak Usaha", anak); len(h.Serupa) != 0 {
		t.Errorf("organisasi yang diubah ikut dihitung: %+v", h)
	}
	if h, _ := l.Periksa(ctx, "Nama Lain Sekali", ""); len(h.Serupa) != 0 || h.TitleDalamNama != "" {
		t.Errorf("nama berbeda: %+v", h)
	}
	if h, _ := l.Periksa(ctx, "CV.", ""); h.TitleDalamNama != "CV." || len(h.Serupa) != 0 {
		t.Errorf("nama hanya title: %+v", h)
	}
}

func TestKemiripanNama(t *testing.T) {
	for _, k := range []struct {
		a, b, jenis string
	}{
		{"ABC", "ABC", services.NamaSama},
		{"ABC", "ABC INDONESIA", ""},
		{"ABCDE", "ABCDE INDONESIA", services.NamaMirip},
		{"ASURANSI MAJU", "ASURANSI MAJOE", services.NamaMirip},
		{"ASURANSI MAJU", "BANK SENTOSA", ""},
		{"", "ABC", ""},
	} {
		if j, _ := services.KemiripanNama(k.a, k.b); j != k.jenis {
			t.Errorf("%q vs %q: %q, mau %q", k.a, k.b, j, k.jenis)
		}
	}
}

// Organization Name tidak boleh diubah saat Edit (perintah work owner 05-10-2026): nama yang dikirim sama (tanpa
// beda huruf/spasi) diterima dan NAME tersimpan tetap apa adanya - juga nama lama Pega yang memuat title; nama lain
// atau kosong ditolak tanpa menulis apa pun.
func TestNamaTidakBolehDiubah(t *testing.T) {
	g := tiruan.Contoh()
	o := g.Org[anak]
	o.Nama = "PT. UJI Anak Usaha"
	g.Org[anak] = o
	isi := func(nama string) models.Isian {
		return models.Isian{Nama: nama, Title: o.Title, Country: o.Country, BusinessField: o.BusinessField,
			ParentID: o.ParentID, Note: "UJI catatan ubah"}
	}
	l := layanan(g)
	if _, err := l.Ubah(ctx, admin, anak, isi(" pt. uji anak usaha ")); err != nil {
		t.Fatalf("nama sama: %v", err)
	}
	if n := g.Org[anak].Nama; n != "PT. UJI Anak Usaha" {
		t.Errorf("NAME tersimpan berubah jadi %q", n)
	}
	for _, nama := range []string{"UJI Anak Usaha Baru", "PT UJI Anak Usaha Baru", ""} {
		sebelum := g.Org[anak]
		_, err := l.Ubah(ctx, admin, anak, models.Isian{Nama: nama, Title: o.Title, Country: o.Country,
			BusinessField: o.BusinessField, ParentID: o.ParentID, Note: "UJI tidak boleh tertulis"})
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), "cannot be changed") {
			t.Errorf("nama %q: %v", nama, err)
		}
		if g.Org[anak] != sebelum {
			t.Errorf("nama %q ditolak tetapi organisasi tertulis %+v", nama, g.Org[anak])
		}
	}
}

// Title nonaktif (sapaan orang, mis. NY) tidak dihitung sebagai title di nama.
func TestTitleNonaktifDiNamaBoleh(t *testing.T) {
	isi := isianBaru()
	isi.Nama = "UJI NY Trading"
	if _, err := layanan(tiruan.Contoh()).Tambah(ctx, admin, isi); err != nil {
		t.Errorf("title nonaktif NY ikut dihitung: %v", err)
	}
}

// PIC Name dipilih dari akun login AKTIF M_LOGIN_GO (perintah work owner 05-10-2026: "untuk name pada company detail,
// dropdown dari tabel login"). Nama PIC lama Pega yang tidak diubah tetap boleh walau bukan akun login.
func TestPICNamaDariAkunLogin(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	p, err := l.Pilihan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Akun) == 0 || !reflect.DeepEqual(p.Akun, g.Akun) {
		t.Errorf("pilihan akun %+v", p.Akun)
	}
	isi := isianBaru()
	isi.PIC[1].Nama = "UJI Bukan Akun"
	if _, err := l.Tambah(ctx, admin, isi); !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), "PIC row 2") {
		t.Errorf("nama bukan akun login: %v", err)
	}
	if len(g.Disisip) != 0 {
		t.Errorf("ditolak tetapi tertulis %+v", g.Disisip)
	}
	d, err := l.Ambil(ctx, anak)
	if err != nil {
		t.Fatal(err)
	}
	ubah := models.Isian{Nama: d.Nama, Title: d.Title, Country: d.Country, BusinessField: d.BusinessField,
		ParentID: d.ParentID, PIC: d.PIC}
	if _, err := l.Ubah(ctx, admin, anak, ubah); err != nil {
		t.Errorf("PIC lama bukan akun, tidak diubah: %v", err)
	}
	ubah.PIC[0].Nama = "UJI Kontak Lama Diganti"
	if _, err := l.Ubah(ctx, admin, anak, ubah); !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), "PIC row 1") {
		t.Errorf("PIC lama diganti nama bukan akun: %v", err)
	}
	ubah.PIC[0].Nama = "UJI Kontak A"
	if _, err := l.Ubah(ctx, admin, anak, ubah); err != nil {
		t.Errorf("PIC lama diganti akun login: %v", err)
	}
	if n := g.PIC[anak][0].Nama; n != "UJI Kontak A" {
		t.Errorf("NICKNAME tersimpan %q", n)
	}
}

// PIC Position = JOB_POSITION akun login yang dipilih (perintah work owner 05-10-2026: "Position ambil dari tabel login
// sesuai orang yg dipilih"); akun tanpa JOB_POSITION ditolak dengan petunjuk ke Kelola User.
func TestPICPositionDariAkunLogin(t *testing.T) {
	for _, k := range []struct {
		nama, nama2, posisi2, galat string
	}{
		{"posisi bukan jabatan akun", "UJI Kontak B", "Direktur", "PIC row 2: Position must be the Job Position of UJI Kontak B"},
		{"akun tanpa jabatan", "UJI Tanpa Jabatan", "", "PIC row 2: UJI Tanpa Jabatan has no Job Position - fill it in Kelola User"},
	} {
		g := tiruan.Contoh()
		isi := isianBaru()
		isi.PIC[1].Nama, isi.PIC[1].Position = k.nama2, k.posisi2
		_, err := layanan(g).Tambah(ctx, admin, isi)
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), k.galat) {
			t.Errorf("%s: %v", k.nama, err)
		}
		if len(g.Disisip) != 0 {
			t.Errorf("%s: ditolak tetapi tertulis", k.nama)
		}
	}
}

// Kode area "Others" (perintah work owner 05-10-2026: "tambahin others, bisa isi sendiri"): kode di luar daftar kodehp
// diterima bila berupa angka (boleh diawali +) dan disimpan apa adanya di TELFAX_CODE.
func TestKodeAreaOthers(t *testing.T) {
	for _, kode := range []string{"0999", "+62"} {
		g := tiruan.Contoh()
		isi := isianBaru()
		isi.Alamat[0].Telfax[0].Code = " " + kode + " "
		d, err := layanan(g).Tambah(ctx, admin, isi)
		if err != nil {
			t.Fatalf("kode %s: %v", kode, err)
		}
		if c := g.Alamat[d.ID][0].TelfaxCode; c != kode {
			t.Errorf("TELFAX_CODE %q, mau %q", c, kode)
		}
	}
}
