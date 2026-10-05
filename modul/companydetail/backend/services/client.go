package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/companydetail/backend/models"
	"nusantarare/modul/companydetail/backend/repository"
)

// Lebar kolom VARCHAR2 dalam BYTE (katalog DEV 03/04-10-2026; kolom tambahan dari migrasi 805-807).
const (
	LebarNama        = 4000 // CLIENT.NAME
	LebarTitle       = 10   // CLIENT.TITLE
	LebarNPWP        = 1000 // CLIENT.NPWP
	LebarCountry     = 20   // CLIENT.COUNTRY
	LebarCountryName = 50   // CLIENT.COUNTRYNAME
	LebarBU          = 1000 // CLIENT.BU_ID
	LebarGroupName   = 1000 // CLIENT.GROUPNAME
	LebarNote        = 4000 // CLIENT.NOTE
	LebarPelaku      = 64   // CLIENT.CREATED_BY / UPDATED_BY
	LebarPICNama     = 500  // CLIENT_PICLIST.NICKNAME
	LebarPICEmail    = 500  // CLIENT_PICLIST.EMAIL
	LebarPICPosisi   = 100  // CLIENT_PICLIST.POSITION
	LebarPICTelepon  = 100  // CLIENT_PICLIST.PHONENUMBER
	LebarAlamat      = 500  // CLIENT_ADDRESS.ASMADDRESS
	LebarJenis       = 10   // CLIENT_ADDRESS.ASMADDRESSTYPE, TELFAX_TYPE, TELFAX_CODE
	LebarTelfaxNo    = 50   // CLIENT_ADDRESS.TELFAX_NO
	LebarOperator    = 100  // CLIENT_ADDRESS.PXCREATEOPERATOR
)

// Ukuran halaman daftar organisasi.
const (
	UkuranHalaman  = 50
	UkuranMaksimum = 200
)

// Galat layanan. Pesan untuk layar (bahasa Inggris) mengikuti sesudah `: `.
var (
	// ErrTidakAda - organisasi tidak ada (404).
	ErrTidakAda = errors.New("organization not found")
	// ErrMasukanTidakSah - isian ditolak (422).
	ErrMasukanTidakSah = errors.New("invalid input")
	// ErrTanpaPelaku - permintaan tanpa akun pelaku (401).
	ErrTanpaPelaku = errors.New("user account is required")
	// ErrBelumDimigrasi - migrasi 800-810 belum dijalankan (503).
	ErrBelumDimigrasi = repository.ErrBelumDimigrasi
)

func tolak(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrMasukanTidakSah, fmt.Sprintf(format, a...))
}

// Transaksi menjalankan fn di dalam satu transaksi dan menutupnya.
type Transaksi func(ctx context.Context, fn func(tx *dbTx) error) error

// Gudang adalah semua yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	CariOrg(ctx context.Context, kueri string, offset, ukuran int) ([]models.BarisDaftar, int, error)
	AmbilOrg(ctx context.Context, tx *dbTx, id string) (models.Organisasi, error)
	// NomorOrgBerikut - SEQ_CLIENT_ORG.NEXTVAL: dua Create bersamaan tidak pernah mendapat nomor ORG yang sama.
	NomorOrgBerikut(ctx context.Context, tx *dbTx) (int64, error)
	SisipOrg(ctx context.Context, tx *dbTx, o models.Organisasi) error
	PerbaruiOrg(ctx context.Context, tx *dbTx, o models.Organisasi) error
	DaftarPIC(ctx context.Context, tx *dbTx, clientID string) ([]models.PIC, error)
	GantiPIC(ctx context.Context, tx *dbTx, clientID string, pic []models.PIC) error
	DaftarBarisAlamat(ctx context.Context, tx *dbTx, clientID string) ([]models.BarisAlamat, error)
	GantiAlamat(ctx context.Context, tx *dbTx, clientID string, baris []models.BarisAlamat) error
	DaftarPilihan(ctx context.Context) ([]models.Pilihan, error)
	DaftarNegara(ctx context.Context) ([]models.Negara, error)
	// DaftarAkunAktif - akun login aktif M_LOGIN_GO (pilihan PIC Name).
	DaftarAkunAktif(ctx context.Context) ([]models.Akun, error)
	CariInduk(ctx context.Context, kueri, kecuali string) ([]models.BarisDaftar, error)
	// DaftarNamaOrg - seluruh organisasi bernama (pemeriksaan nama sama/mirip).
	DaftarNamaOrg(ctx context.Context) ([]models.BarisDaftar, error)
	// SiapkanOrgLama - isi popup Copy Old (dokumen M_CLIENT yang datanya belum pindah).
	SiapkanOrgLama(ctx context.Context) ([]models.OrgLama, error)
	// ReferensiPindah - title dan negara untuk salinan Copy Old; dibaca SEKALI per Process Copy.
	ReferensiPindah(ctx context.Context) (repository.ReferensiPindah, error)
	// SalinOrgLama - salin satu organisasi dokumen di transaksi pemanggil; false = tidak ada yang perlu ditulis.
	SalinOrgLama(ctx context.Context, tx *dbTx, id string, r repository.ReferensiPindah) (bool, error)
}

// Layanan memegang seluruh aturan modul ini di atas satu Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
	jam    func() time.Time
}

// BaruLayanan menyusun Layanan - dipakai uji dengan gudang tiruan dan transaksi tiruan (fn(nil)).
func BaruLayanan(g Gudang, tx Transaksi) *Layanan { return &Layanan{gudang: g, tx: tx, jam: time.Now} }

// DenganJam mengganti jam layanan (uji).
func (l *Layanan) DenganJam(jam func() time.Time) *Layanan {
	l.jam = jam
	return l
}

// Halaman adalah satu halaman daftar organisasi.
type Halaman struct {
	Daftar  []models.BarisDaftar `json:"daftar"`
	Total   int                  `json:"total"`
	Halaman int                  `json:"halaman"`
	Ukuran  int                  `json:"ukuran"`
}

// PilihanForm adalah isi dropdown form. Setiap daftar memuat juga pilihan tidak aktif (untuk menampilkan nilai
// lama); dropdown hanya menawarkan yang aktif.
type PilihanForm struct {
	Title         []models.Pilihan `json:"title"`
	BusinessField []models.Pilihan `json:"businessField"`
	Position      []models.Pilihan `json:"position"`
	AddressType   []models.Pilihan `json:"addressType"`
	Telfax        []models.Pilihan `json:"telfax"`
	KodeArea      []models.Pilihan `json:"kodeArea"`
	Gender        []models.Pilihan `json:"gender"`
	Negara        []models.Negara  `json:"negara"`
	// Akun - pilihan PIC Name: akun login aktif M_LOGIN_GO (perintah work owner 05-10-2026).
	Akun []models.Akun `json:"akun"`
}

// referensi - pilihan yang dipakai pemeriksa isian.
type referensi struct {
	enum   map[string]map[string]models.Pilihan // jenis -> kode -> pilihan
	negara map[string]models.Negara             // CLIENT.COUNTRY -> negara
	akun   []models.Akun                        // akun login aktif
	// jabatanAkun - NAME akun login aktif -> JOB_POSITION-nya (nama kembar boleh berbeda jabatan).
	jabatanAkun map[string]map[string]bool
}

func (r referensi) ada(jenis, kode string) (models.Pilihan, bool) {
	p, ok := r.enum[jenis][kode]
	return p, ok
}

func (r referensi) aktif(jenis, kode string) bool {
	p, ok := r.ada(jenis, kode)
	return ok && p.Aktif
}

func (l *Layanan) muatReferensi(ctx context.Context) (referensi, []models.Pilihan, []models.Negara, error) {
	pil, err := l.gudang.DaftarPilihan(ctx)
	if err != nil {
		return referensi{}, nil, nil, err
	}
	neg, err := l.gudang.DaftarNegara(ctx)
	if err != nil {
		return referensi{}, nil, nil, err
	}
	akun, err := l.gudang.DaftarAkunAktif(ctx)
	if err != nil {
		return referensi{}, nil, nil, err
	}
	r := referensi{enum: map[string]map[string]models.Pilihan{}, negara: map[string]models.Negara{}, akun: akun,
		jabatanAkun: map[string]map[string]bool{}}
	for _, a := range akun {
		n := strings.TrimSpace(a.Nama)
		if r.jabatanAkun[n] == nil {
			r.jabatanAkun[n] = map[string]bool{}
		}
		r.jabatanAkun[n][strings.TrimSpace(a.Jabatan)] = true
	}
	for _, p := range pil {
		if r.enum[p.Jenis] == nil {
			r.enum[p.Jenis] = map[string]models.Pilihan{}
		}
		r.enum[p.Jenis][p.Kode] = p
	}
	for _, n := range neg {
		r.negara[n.Kode()] = n
	}
	return r, pil, neg, nil
}

func periksaPelaku(p inti.Pelaku) error {
	if strings.TrimSpace(p.AkunID) == "" {
		return ErrTanpaPelaku
	}
	if len(p.AkunID) > LebarPelaku {
		return tolak("your user account is longer than %d bytes", LebarPelaku)
	}
	return nil
}

// Daftar membaca satu halaman organisasi; `kueri` mencari nama, ORG ID, dan NPWP.
func (l *Layanan) Daftar(ctx context.Context, kueri string, halaman, ukuran int) (Halaman, error) {
	if ukuran <= 0 {
		ukuran = UkuranHalaman
	}
	if ukuran > UkuranMaksimum {
		ukuran = UkuranMaksimum
	}
	if halaman < 1 {
		halaman = 1
	}
	baris, total, err := l.gudang.CariOrg(ctx, strings.TrimSpace(kueri), (halaman-1)*ukuran, ukuran)
	if err != nil {
		return Halaman{}, err
	}
	return Halaman{Daftar: baris, Total: total, Halaman: halaman, Ukuran: ukuran}, nil
}

// Pilihan menyusun isi dropdown form.
func (l *Layanan) Pilihan(ctx context.Context) (PilihanForm, error) {
	r, pil, neg, err := l.muatReferensi(ctx)
	if err != nil {
		return PilihanForm{}, err
	}
	f := PilihanForm{Title: []models.Pilihan{}, BusinessField: []models.Pilihan{}, Position: []models.Pilihan{},
		AddressType: []models.Pilihan{}, Telfax: []models.Pilihan{}, KodeArea: []models.Pilihan{},
		Gender: []models.Pilihan{}, Negara: neg, Akun: r.akun}
	tujuan := map[string]*[]models.Pilihan{models.JenisTitle: &f.Title, models.JenisBidangUsaha: &f.BusinessField,
		models.JenisPosisi: &f.Position, models.JenisAlamat: &f.AddressType, models.JenisTelfax: &f.Telfax,
		models.JenisKodeHP: &f.KodeArea, models.JenisGender: &f.Gender}
	for _, p := range pil {
		if t, ok := tujuan[p.Jenis]; ok {
			*t = append(*t, p)
		}
	}
	if f.Negara == nil {
		f.Negara = []models.Negara{}
	}
	if f.Akun == nil {
		f.Akun = []models.Akun{}
	}
	return f, nil
}

// CariInduk - pilihan Parent organization: organisasi bernama mirip `kueri`, selain `kecuali`.
func (l *Layanan) CariInduk(ctx context.Context, kueri, kecuali string) ([]models.BarisDaftar, error) {
	return l.gudang.CariInduk(ctx, strings.TrimSpace(kueri), strings.TrimSpace(kecuali))
}

// Ambil membaca satu organisasi lengkap.
func (l *Layanan) Ambil(ctx context.Context, id string) (models.Detail, error) {
	return l.ambil(ctx, nil, id)
}

func (l *Layanan) ambil(ctx context.Context, tx *dbTx, id string) (models.Detail, error) {
	o, err := l.gudang.AmbilOrg(ctx, tx, id)
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Detail{}, ErrTidakAda
	}
	if err != nil {
		return models.Detail{}, err
	}
	pic, err := l.gudang.DaftarPIC(ctx, tx, id)
	if err != nil {
		return models.Detail{}, err
	}
	baris, err := l.gudang.DaftarBarisAlamat(ctx, tx, id)
	if err != nil {
		return models.Detail{}, err
	}
	return SusunDetail(o, pic, baris), nil
}

// SusunDetail menyusun satu organisasi untuk layar: tanggal lahir PIC ke bentuk kabel, baris alamat dikelompokkan
// per alamat (urutan kemunculan pertama).
func SusunDetail(o models.Organisasi, pic []models.PIC, baris []models.BarisAlamat) models.Detail {
	d := models.Detail{Organisasi: o, PIC: make([]models.PIC, 0, len(pic)), Alamat: KelompokAlamat(baris)}
	for _, p := range pic {
		p.DateOfBirth = TanggalKeKabel(p.DateOfBirth)
		d.PIC = append(d.PIC, p)
	}
	return d
}

// KelompokAlamat menyatukan baris CLIENT_ADDRESS beralamat sama menjadi satu alamat; baris bernomor menjadi
// Phone and Fax-nya.
func KelompokAlamat(baris []models.BarisAlamat) []models.Alamat {
	out := []models.Alamat{}
	letak := map[string]int{}
	for _, b := range baris {
		i, ada := letak[b.Address]
		if !ada {
			i = len(out)
			letak[b.Address] = i
			out = append(out, models.Alamat{Asal: b.Address, Jenis: b.Type, Address: b.Address, Telfax: []models.Telfax{}})
		}
		if out[i].Jenis == "" {
			out[i].Jenis = b.Type
		}
		if b.TelfaxType != "" || b.TelfaxCode != "" || b.TelfaxNo != "" {
			out[i].Telfax = append(out[i].Telfax, models.Telfax{Jenis: b.TelfaxType, Code: b.TelfaxCode, No: b.TelfaxNo})
		}
	}
	return out
}

// TanggalKeKabel - `YYYYMMDD` (bentuk simpan) ke `DD-MM-YYYY`; bentuk lain dikembalikan apa adanya.
func TanggalKeKabel(s string) string {
	t, err := time.Parse("20060102", s)
	if err != nil {
		return s
	}
	return t.Format("02-01-2006")
}

// TanggalDariKabel - `DD-MM-YYYY` ke `YYYYMMDD`; false = bukan tanggal sah.
func TanggalDariKabel(s string) (string, bool) {
	t, err := time.Parse("02-01-2006", s)
	if err != nil {
		return "", false
	}
	return t.Format("20060102"), true
}

var polaEmail = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// polaKodeArea - kode area "Others" yang diisi sendiri: angka, boleh diawali + (mis. `0999`, `+62`).
var polaKodeArea = regexp.MustCompile(`^\+?[0-9]+$`)

// polaNomorPIC - USERIDENTIFIER buatan aplikasi ini.
var polaNomorPIC = regexp.MustCompile(`^` + models.AwalanPIC + `([0-9]+)$`)

func rapikan(isi models.Isian) models.Isian {
	// Organization Name selalu huruf besar (perintah work owner 04-10-2026: "langsung uppercase untuk nama company
	// detail") - juga saat diubah, juga bila API dipanggil tanpa layar.
	isi.Nama = strings.ToUpper(strings.TrimSpace(isi.Nama))
	isi.Title = strings.TrimSpace(isi.Title)
	isi.NPWP = strings.TrimSpace(isi.NPWP)
	isi.Country = strings.TrimSpace(isi.Country)
	isi.BusinessField = strings.TrimSpace(isi.BusinessField)
	isi.ParentID = strings.TrimSpace(isi.ParentID)
	isi.Note = strings.TrimSpace(isi.Note)
	return isi
}

// isiOrg mengisi kolom layar baris organisasi. `lama` nil = baris baru; nilai yang TIDAK diubah tidak diperiksa
// ulang terhadap daftar pilihan (nilai lama Pega boleh tetap).
func (l *Layanan) isiOrg(ctx context.Context, tx *dbTx, o *models.Organisasi, isi models.Isian, r referensi,
	lama *models.Organisasi) error {
	// Edit: Organization Name tidak boleh diubah (perintah work owner 05-10-2026). Nama yang dikirim harus sama
	// (tanpa beda huruf/spasi) dan NAME tersimpan tetap apa adanya - 4.472 nama lama Pega (DEV) yang memuat
	// PT/CV/PD/UD ikut tetap.
	if lama != nil {
		if isi.Nama != strings.ToUpper(strings.TrimSpace(lama.Nama)) {
			return tolak("Organization Name cannot be changed")
		}
		o.Nama = lama.Nama
	} else {
		if isi.Nama == "" {
			return tolak("Organization Name is required")
		}
		if len(isi.Nama) > LebarNama {
			return tolak("Organization Name is longer than %d bytes", LebarNama)
		}
		// Title tidak boleh ada di nama (work owner 04-10-2026: "ga boleh disimpan dong").
		if t := TitleDiNama(isi.Nama, kataTitle(r)); t != "" {
			return tolak("Organization Name must not contain the title %s - choose it in Title", t)
		}
		o.Nama = isi.Nama
	}

	o.Title = ""
	if isi.Title != "" {
		sah := lama != nil && isi.Title == lama.Title
		for _, p := range r.enum[models.JenisTitle] {
			sah = sah || (p.Aktif && p.Label == isi.Title)
		}
		if !sah {
			return tolak("Title %s is not in the title list", isi.Title)
		}
		if len(isi.Title) > LebarTitle {
			return tolak("Title %s is longer than %d bytes", isi.Title, LebarTitle)
		}
		o.Title = isi.Title
	}

	if len(isi.NPWP) > LebarNPWP {
		return tolak("NPWP is longer than %d bytes", LebarNPWP)
	}
	o.NPWP = isi.NPWP

	if isi.Country == "" {
		return tolak("COUNTRY is required")
	}
	if n, ada := r.negara[isi.Country]; ada {
		if len(n.Kode()) > LebarCountry || len(n.Nama) > LebarCountryName {
			return tolak("COUNTRY %s does not fit the client table", n.Nama)
		}
		o.Country, o.CountryName = n.Kode(), n.Nama
	} else if lama == nil || isi.Country != lama.Country {
		return tolak("COUNTRY %s is not in the country list (NATION)", isi.Country)
	}

	if isi.BusinessField == "" {
		return tolak("Business Field is required")
	}
	if !r.aktif(models.JenisBidangUsaha, isi.BusinessField) && (lama == nil || isi.BusinessField != lama.BusinessField) {
		return tolak("Business Field %s is not in the business field list", isi.BusinessField)
	}
	if len(isi.BusinessField) > LebarBU {
		return tolak("Business Field %s is longer than %d bytes", isi.BusinessField, LebarBU)
	}
	o.BusinessField = isi.BusinessField

	switch {
	case lama != nil && isi.ParentID == lama.ParentID:
		// Tidak diubah: GROUPNAME lama (juga yang tanpa PARENT_ID) tetap.
	case isi.ParentID == "":
		o.ParentID, o.ParentName = "", ""
	case isi.ParentID == o.ID:
		return tolak("An organization cannot be its own parent")
	default:
		induk, err := l.gudang.AmbilOrg(ctx, tx, isi.ParentID)
		if errors.Is(err, repository.ErrTidakAda) {
			return tolak("Parent organization %s not found", isi.ParentID)
		}
		if err != nil {
			return err
		}
		if len(induk.Nama) > LebarGroupName {
			return tolak("the name of Parent organization %s is longer than %d bytes", induk.IDView, LebarGroupName)
		}
		o.ParentID, o.ParentName = induk.ID, induk.Nama
	}

	if len(isi.Note) > LebarNote {
		return tolak("Note is longer than %d bytes", LebarNote)
	}
	o.Note = isi.Note
	return nil
}

// susunPIC memeriksa grid PIC dan menyusun barisnya. `lama` = PIC organisasi ini sebelum diubah (nil = baru).
// PIC lama dikenali dari USERIDENTIFIER-nya; PIC baru mendapat `PIC-n` (nomor tertinggi organisasi ini + 1), urutan
// baris layar dipertahankan.
func susunPIC(masuk, lama []models.PIC, r referensi) ([]models.PIC, error) {
	lamaPer := map[string]models.PIC{}
	maks := 0
	catatNomor := func(uid string) {
		if m := polaNomorPIC.FindStringSubmatch(uid); m != nil {
			if n, err := strconv.Atoi(m[1]); err == nil && n > maks {
				maks = n
			}
		}
	}
	for _, p := range lama {
		if p.UserIdentifier != "" {
			lamaPer[p.UserIdentifier] = p
			catatNomor(p.UserIdentifier)
		}
	}
	out := make([]models.PIC, 0, len(masuk))
	dipakai := map[string]int{}
	for i, p := range masuk {
		ke := i + 1
		p.UserIdentifier = strings.TrimSpace(p.UserIdentifier)
		p.Nama = strings.TrimSpace(p.Nama)
		p.Position = strings.TrimSpace(p.Position)
		p.Gender = strings.TrimSpace(p.Gender)
		p.Email = strings.TrimSpace(p.Email)
		p.DateOfBirth = strings.TrimSpace(p.DateOfBirth)
		p.Phone = strings.TrimSpace(p.Phone)
		if p.Nama == "" {
			return nil, tolak("PIC row %d: Name is required", ke)
		}
		if len(p.Nama) > LebarPICNama {
			return nil, tolak("PIC row %d: Name is longer than %d bytes", ke, LebarPICNama)
		}
		lamaPIC, adaLama := lamaPer[p.UserIdentifier]
		if p.UserIdentifier != "" && !adaLama {
			return nil, tolak("PIC row %d does not belong to this organization", ke)
		}
		// Name dipilih dari akun login aktif dan Position = JOB_POSITION akun itu (perintah work owner 05-10-2026); PIC
		// lama Pega yang namanya tidak diubah tetap boleh.
		if !(adaLama && p.Nama == strings.TrimSpace(lamaPIC.Nama)) {
			jabatan, ada := r.jabatanAkun[p.Nama]
			switch {
			case !ada:
				return nil, tolak("PIC row %d: Name %s is not an active user in the login list", ke, p.Nama)
			case p.Position == "":
				return nil, tolak("PIC row %d: %s has no Job Position - fill it in Kelola User", ke, p.Nama)
			case !jabatan[p.Position]:
				return nil, tolak("PIC row %d: Position must be the Job Position of %s in Kelola User", ke, p.Nama)
			}
		}
		if p.Position == "" {
			return nil, tolak("PIC row %d: Position is required", ke)
		}
		if len(p.Position) > LebarPICPosisi {
			return nil, tolak("PIC row %d: Position is longer than %d bytes", ke, LebarPICPosisi)
		}
		if p.Gender != "" {
			if _, ada := r.ada(models.JenisGender, p.Gender); !ada {
				return nil, tolak("PIC row %d: Gender %s is not in the gender list", ke, p.Gender)
			}
		}
		if p.Email != "" && (len(p.Email) > LebarPICEmail || !polaEmail.MatchString(p.Email)) {
			return nil, tolak("PIC row %d: Email format is not valid", ke)
		}
		if len(p.Phone) > LebarPICTelepon {
			return nil, tolak("PIC row %d: Phone number is longer than %d bytes", ke, LebarPICTelepon)
		}
		if sebelum, ganda := dipakai[p.UserIdentifier]; ganda && p.UserIdentifier != "" {
			return nil, tolak("PIC row %d repeats PIC row %d", ke, sebelum)
		}
		switch {
		case p.DateOfBirth == "":
		case adaLama && p.DateOfBirth == TanggalKeKabel(lamaPIC.DateOfBirth):
			p.DateOfBirth = lamaPIC.DateOfBirth
		default:
			simpan, sah := TanggalDariKabel(p.DateOfBirth)
			if !sah {
				return nil, tolak("PIC row %d: Date of birth must be DD-MM-YYYY", ke)
			}
			p.DateOfBirth = simpan
		}
		if p.UserIdentifier != "" {
			dipakai[p.UserIdentifier] = ke
		}
		out = append(out, p)
	}
	// PIC baru diberi nomor SESUDAH seluruh nomor lama tercatat.
	for i := range out {
		if out[i].UserIdentifier == "" {
			maks++
			out[i].UserIdentifier = models.AwalanPIC + strconv.Itoa(maks)
		}
	}
	return out, nil
}

// FormatWaktuPega - bentuk PXCREATEDATETIME baris alamat (bentuk tanggal-waktu Pega, GMT).
func FormatWaktuPega(t time.Time) string { return t.UTC().Format("20060102T150405.000") + " GMT" }

// susunAlamat memeriksa grid Address dan menyusun baris CLIENT_ADDRESS-nya: satu baris per nomor Phone and Fax,
// satu baris kosong-nomor untuk alamat tanpa nomor. Kolom yang tidak ada di layar (kota, kode pos, kecamatan,
// provinsi, RW, jejak pembuatan) dibawa dari baris lama alamat itu (`Asal`); alamat baru mendapat jejak pelaku.
func susunAlamat(masuk []models.Alamat, lama []models.BarisAlamat, r referensi, pelaku string, jam time.Time) ([]models.BarisAlamat, error) {
	lamaPer := map[string]models.BarisAlamat{}
	jenisLama := map[string]map[string]bool{} // alamat -> TELFAX_TYPE lama
	kodeLama := map[string]map[string]bool{}  // alamat -> TELFAX_CODE lama
	for _, b := range lama {
		if _, ada := lamaPer[b.Address]; !ada {
			lamaPer[b.Address] = b
			jenisLama[b.Address], kodeLama[b.Address] = map[string]bool{}, map[string]bool{}
		}
		jenisLama[b.Address][b.TelfaxType] = true
		kodeLama[b.Address][b.TelfaxCode] = true
	}
	if len(pelaku) > LebarOperator {
		return nil, tolak("your user account is longer than %d bytes", LebarOperator)
	}
	out := []models.BarisAlamat{}
	dipakai := map[string]int{}
	for i, a := range masuk {
		ke := i + 1
		a.Asal, a.Jenis, a.Address = strings.TrimSpace(a.Asal), strings.TrimSpace(a.Jenis), strings.TrimSpace(a.Address)
		if a.Address == "" {
			return nil, tolak("Address row %d: Address is required", ke)
		}
		if len(a.Address) > LebarAlamat {
			return nil, tolak("Address row %d: Address is longer than %d bytes", ke, LebarAlamat)
		}
		if sebelum, ganda := dipakai[a.Address]; ganda {
			return nil, tolak("Address row %d repeats the address of row %d", ke, sebelum)
		}
		dipakai[a.Address] = ke
		asal, adaAsal := lamaPer[a.Asal]
		if a.Asal != "" && !adaAsal {
			return nil, tolak("Address row %d does not belong to this organization", ke)
		}
		if a.Jenis == "" {
			return nil, tolak("Address row %d: Type is required", ke)
		}
		if !r.aktif(models.JenisAlamat, a.Jenis) && !(adaAsal && a.Jenis == asal.Type) {
			return nil, tolak("Address row %d: Type %s is not in the address type list", ke, a.Jenis)
		}
		if len(a.Jenis) > LebarJenis {
			return nil, tolak("Address row %d: Type is longer than %d bytes", ke, LebarJenis)
		}
		dasar := models.BarisAlamat{Type: a.Jenis, Address: a.Address, PxCreateOperator: pelaku,
			PxCreateDateTime: FormatWaktuPega(jam)}
		if adaAsal {
			dasar.City, dasar.CityName, dasar.ZipCode = asal.City, asal.CityName, asal.ZipCode
			dasar.DistrictName, dasar.ProvinceName, dasar.RWName = asal.DistrictName, asal.ProvinceName, asal.RWName
			dasar.PxCreateOperator, dasar.PxCreateDateTime = asal.PxCreateOperator, asal.PxCreateDateTime
		}
		if len(a.Telfax) == 0 {
			out = append(out, dasar)
			continue
		}
		for j, t := range a.Telfax {
			nomor := fmt.Sprintf("Address row %d, Phone and Fax %d", ke, j+1)
			t.Jenis, t.Code, t.No = strings.TrimSpace(t.Jenis), strings.TrimSpace(t.Code), strings.TrimSpace(t.No)
			if t.No == "" {
				return nil, tolak("%s: the number is required", nomor)
			}
			if len(t.No) > LebarTelfaxNo {
				return nil, tolak("%s: the number is longer than %d bytes", nomor, LebarTelfaxNo)
			}
			if t.Jenis == "" {
				return nil, tolak("%s: the type is required", nomor)
			}
			if !r.aktif(models.JenisTelfax, t.Jenis) && !(adaAsal && jenisLama[a.Asal][t.Jenis]) {
				return nil, tolak("%s: type %s is not in the phone and fax type list", nomor, t.Jenis)
			}
			if len(t.Jenis) > LebarJenis || len(t.Code) > LebarJenis {
				return nil, tolak("%s: the type or area code is longer than %d bytes", nomor, LebarJenis)
			}
			// Kode area: dari daftar kodehp, atau "Others" diisi sendiri (perintah work owner 05-10-2026) - angka, boleh
			// diawali +. Kode lama Pega baris ini tetap boleh.
			if t.Code != "" {
				_, ada := r.ada(models.JenisKodeHP, t.Code)
				if !ada && !(adaAsal && kodeLama[a.Asal][t.Code]) && !polaKodeArea.MatchString(t.Code) {
					return nil, tolak("%s: area code %s must be digits (a leading + is allowed)", nomor, t.Code)
				}
			}
			b := dasar
			b.TelfaxType, b.TelfaxCode, b.TelfaxNo = t.Jenis, t.Code, t.No
			out = append(out, b)
		}
	}
	return out, nil
}

// Tambah membuat organisasi baru: nomor ORG dari SEQ_CLIENT_ORG (migrasi 810), diambil SESUDAH seluruh isian lolos.
func (l *Layanan) Tambah(ctx context.Context, p inti.Pelaku, isi models.Isian) (models.Detail, error) {
	if err := periksaPelaku(p); err != nil {
		return models.Detail{}, err
	}
	r, _, _, err := l.muatReferensi(ctx)
	if err != nil {
		return models.Detail{}, err
	}
	isi = rapikan(isi)
	var hasil models.Detail
	err = l.tx(ctx, func(tx *dbTx) error {
		o := models.Organisasi{CreatedBy: p.AkunID, UpdatedBy: p.AkunID}
		if err := l.isiOrg(ctx, tx, &o, isi, r, nil); err != nil {
			return err
		}
		pic, err := susunPIC(isi.PIC, nil, r)
		if err != nil {
			return err
		}
		baris, err := susunAlamat(isi.Alamat, nil, r, p.AkunID, l.jam())
		if err != nil {
			return err
		}
		n, err := l.gudang.NomorOrgBerikut(ctx, tx)
		if err != nil {
			return err
		}
		o.IDView = models.AwalanIDView + strconv.FormatInt(n, 10)
		o.ID = models.AwalanID + o.IDView
		if err := l.gudang.SisipOrg(ctx, tx, o); err != nil {
			return err
		}
		if err := l.gudang.GantiPIC(ctx, tx, o.ID, pic); err != nil {
			return err
		}
		if err := l.gudang.GantiAlamat(ctx, tx, o.ID, baris); err != nil {
			return err
		}
		hasil, err = l.ambil(ctx, tx, o.ID)
		return err
	})
	return hasil, err
}

// Ubah menulis ulang organisasi, seluruh PIC-nya, dan seluruh alamatnya (baris yang dibuang dari grid dihapus).
func (l *Layanan) Ubah(ctx context.Context, p inti.Pelaku, id string, isi models.Isian) (models.Detail, error) {
	if err := periksaPelaku(p); err != nil {
		return models.Detail{}, err
	}
	r, _, _, err := l.muatReferensi(ctx)
	if err != nil {
		return models.Detail{}, err
	}
	isi = rapikan(isi)
	var hasil models.Detail
	err = l.tx(ctx, func(tx *dbTx) error {
		lama, err := l.gudang.AmbilOrg(ctx, tx, id)
		if errors.Is(err, repository.ErrTidakAda) {
			return ErrTidakAda
		}
		if err != nil {
			return err
		}
		lamaPIC, err := l.gudang.DaftarPIC(ctx, tx, id)
		if err != nil {
			return err
		}
		lamaBaris, err := l.gudang.DaftarBarisAlamat(ctx, tx, id)
		if err != nil {
			return err
		}
		o := lama
		o.UpdatedBy = p.AkunID
		if err := l.isiOrg(ctx, tx, &o, isi, r, &lama); err != nil {
			return err
		}
		pic, err := susunPIC(isi.PIC, lamaPIC, r)
		if err != nil {
			return err
		}
		baris, err := susunAlamat(isi.Alamat, lamaBaris, r, p.AkunID, l.jam())
		if err != nil {
			return err
		}
		if err := l.gudang.PerbaruiOrg(ctx, tx, o); errors.Is(err, repository.ErrTidakAda) {
			return ErrTidakAda
		} else if err != nil {
			return err
		}
		if err := l.gudang.GantiPIC(ctx, tx, id, pic); err != nil {
			return err
		}
		if err := l.gudang.GantiAlamat(ctx, tx, id, baris); err != nil {
			return err
		}
		hasil, err = l.ambil(ctx, tx, id)
		return err
	})
	return hasil, err
}
