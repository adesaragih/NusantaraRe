package login

// Untuk apa berkas ini: PEMUAT AKUN dari ekspor operator Pega (perintah work owner 11-10-2026: "buatkan akun, ikuti
// sheet LOGIN; untuk ID contact yang sudah ada cek pada sheet CONTACT, jika belum ada tambahin; set default password
// 123456 untuk akun baru"). Masukan = sheet LOGIN dan sheet CONTACT yang disimpan sebagai CSV; dijalankan operator lewat
// `cmd/api -muat-pengguna` - uji-kering bawaan, menulis hanya dengan `jalankan`. Data orang tidak pernah ditanam di
// repo.
//
// Keputusan work owner 11-10-2026:
//   - divisi / unit dipetakan dari nama divisi Pega (`petaDivisiPega`); unit Pega (Unit / Life / NonLife / Admin) tidak
//     dipakai;
//   - email kembar dipakai akun PERTAMA di berkas, berikutnya dikosongkan; email tidak sah dikosongkan;
//   - akun yang sudah ada (username tanpa beda huruf) dilewati - tidak diubah;
//   - akun baru: password `SandiBawaan`, TIDAK wajib ganti, tanpa menu / workbasket (diatur di Kelola User);
//   - CONTACT_ID = kode contact sheet CONTACT (`CON-n` SFAGIS lama) bila ada, selain itu sequence (CON-1001 dst.).

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// polaIDKontak - bentuk CONTACT_ID (migrasi 905).
var polaIDKontak = regexp.MustCompile(`^CON-[0-9]+$`)

// Judul kolom sheet LOGIN (ekspor operator Pega) dan sheet CONTACT.
const (
	kolomEmail      = "PYEMAILADDRESS"
	kolomDivisi     = "PYORGDIVISION"
	kolomOrganisasi = "PYORGANIZATION"
	kolomJabatan    = "PYPOSITION"
	kolomAkun       = "PYUSERIDENTIFIER"
	kolomNama       = "PYUSERNAME"
	kolomKontak     = "CONTACT"
	kolomAksesLogin = "AKSES LOGIN"
)

// petaDivisiPega - divisi operator Pega (huruf kecil) -> divisi dan unit master (keputusan work owner 11-10-2026
// "petakan otomatis"). Divisi Pega lain (DivisionHead, Supervisor, Manager, Statistics) = tanpa divisi dan unit.
var petaDivisiPega = map[string][2]string{
	"marketing":    {"MKT", ""},
	"underwriting": {"TECH", "UW"},
	"claim":        {"TECH", "CLM"},
	"it":           {"IT", "ITAPP"},
	"director":     {"EXEC", "BOD"},
}

// Status satu baris pemuat akun.
const (
	StatusMuatSiap     = "siap"      // uji-kering: lolos semua pemeriksaan
	StatusMuatDibuat   = "dibuat"    // -jalankan: akun ditulis
	StatusMuatSudahAda = "sudah ada" // username sudah dipakai (tanpa beda huruf) - dilewati
	StatusMuatGagal    = "gagal"     // tidak lolos pemeriksaan - tidak ditulis
)

// BarisMuat - satu akun sheet LOGIN yang sudah disusun.
type BarisMuat struct {
	// Baris - nomor baris CSV (judul = 1).
	Baris int
	Akun  AkunBaru
	// Catatan - perubahan yang diambil pemuat (email dikosongkan, divisi tanpa padanan).
	Catatan []string
}

// HasilMuat - hasil satu baris.
type HasilMuat struct {
	Baris                      int
	ID, Nama, IDKontak, Status string
	Catatan                    []string
}

// BacaCSV membaca CSV berjudul (judul = baris pertama, huruf besar, spasi tepi dibuang): pemisah koma, titik koma, atau
// tab dikenali dari baris judul; BOM UTF-8 dibuang. Baris kosong dilewati.
func BacaCSV(r io.Reader) ([]map[string]string, error) {
	isi, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	isi = bytes.TrimPrefix(isi, []byte("\xef\xbb\xbf"))
	judulBaris, _, _ := bytes.Cut(isi, []byte("\n"))
	c := csv.NewReader(bytes.NewReader(isi))
	c.FieldsPerRecord = -1
	c.Comma = ','
	for _, p := range []rune{';', '\t'} {
		if bytes.Count(judulBaris, []byte(string(p))) > bytes.Count(judulBaris, []byte(string(c.Comma))) {
			c.Comma = p
		}
	}
	semua, err := c.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("login: membaca CSV: %w", err)
	}
	if len(semua) == 0 {
		return nil, fmt.Errorf("login: CSV kosong")
	}
	judul := make([]string, len(semua[0]))
	for i, j := range semua[0] {
		judul[i] = strings.ToUpper(strings.TrimSpace(j))
	}
	var out []map[string]string
	for _, rec := range semua[1:] {
		m, isi := map[string]string{}, false
		for i, v := range rec {
			if i < len(judul) {
				m[judul[i]] = strings.TrimSpace(v)
				isi = isi || strings.TrimSpace(v) != ""
			}
		}
		if isi {
			out = append(out, m)
		}
	}
	return out, nil
}

// PetaKontak - sheet CONTACT: AKSES LOGIN (huruf besar) -> CONTACT_ID (`CON-n`, kata terakhir kolom CONTACT, mis.
// "ASM-SFAGIS-WORK-CONTACT CON-49"). Baris tanpa akses login dilewati; CONTACT bukan `CON-angka` = galat.
func PetaKontak(baris []map[string]string) (map[string]string, error) {
	out := map[string]string{}
	for i, b := range baris {
		akses := strings.ToUpper(strings.TrimSpace(b[kolomAksesLogin]))
		if akses == "" {
			continue
		}
		kata := strings.Fields(b[kolomKontak])
		if len(kata) == 0 || !polaIDKontak.MatchString(kata[len(kata)-1]) {
			return nil, fmt.Errorf("login: sheet CONTACT baris %d: %q bukan kode CON-angka", i+2, b[kolomKontak])
		}
		out[akses] = kata[len(kata)-1]
	}
	return out, nil
}

// SusunMuat menyusun akun dari sheet LOGIN (`login`, urutan berkas) dan peta contact (`PetaKontak`): email dirapikan
// (tidak sah / kembar dikosongkan - akun pertama di berkas yang memakainya), divisi / unit dari `petaDivisiPega`.
func SusunMuat(login []map[string]string, kontak map[string]string) []BarisMuat {
	pemilikEmail := map[string]int{}
	out := make([]BarisMuat, 0, len(login))
	for i, b := range login {
		bm := BarisMuat{Baris: i + 2}
		a := AkunBaru{ID: strings.TrimSpace(b[kolomAkun]), Nama: strings.TrimSpace(b[kolomNama]),
			Organisasi: strings.TrimSpace(b[kolomOrganisasi])}
		a.Kontak = RapikanKontak(Kontak{Email: b[kolomEmail], Jabatan: b[kolomJabatan]})
		if e := a.Email; e != "" {
			switch baris, ada := pemilikEmail[e]; {
			case PeriksaKontak(Kontak{Email: e}) != nil:
				bm.Catatan = append(bm.Catatan, fmt.Sprintf("email %q tidak sah - dikosongkan", e))
				a.Email = ""
			case ada:
				bm.Catatan = append(bm.Catatan, fmt.Sprintf("email %q sudah dipakai baris %d - dikosongkan", e, baris))
				a.Email = ""
			default:
				pemilikEmail[e] = bm.Baris
			}
		}
		if div := strings.TrimSpace(b[kolomDivisi]); div != "" {
			if p, ada := petaDivisiPega[strings.ToLower(div)]; ada {
				a.Divisi, a.Unit = p[0], p[1]
			} else {
				bm.Catatan = append(bm.Catatan, fmt.Sprintf("divisi Pega %q tanpa padanan - divisi / unit kosong", div))
			}
		}
		a.IDKontak = kontak[strings.ToUpper(a.ID)]
		bm.Akun = a
		out = append(out, bm)
	}
	return out
}

// MuatPengguna memeriksa (dan dengan `jalankan` menulis) akun `baris` satu per satu, masing-masing transaksinya sendiri;
// galat satu akun tidak menghentikan yang lain. Akun yang sudah ada dilewati; email yang sudah dipakai akun di basis
// data dikosongkan; CONTACT_ID tetap yang sudah dipakai = gagal.
func (l *Layanan) MuatPengguna(ctx context.Context, baris []BarisMuat, jalankan bool) ([]HasilMuat, error) {
	diproses := map[string]int{}
	out := make([]HasilMuat, 0, len(baris))
	for _, b := range baris {
		a := b.Akun
		h := HasilMuat{Baris: b.Baris, ID: a.ID, Nama: a.Nama, IDKontak: a.IDKontak,
			Catatan: append([]string{}, b.Catatan...)}
		kunci := strings.ToLower(a.ID)
		if sebelum, ada := diproses[kunci]; ada {
			h.Status = StatusMuatGagal
			h.Catatan = append(h.Catatan, fmt.Sprintf("username kembar dengan baris %d di berkas", sebelum))
			out = append(out, h)
			continue
		}
		diproses[kunci] = b.Baris
		if pemakai, err := l.gudang.PemakaiUsername(ctx, a.ID); err != nil {
			return out, err
		} else if len(pemakai) > 0 {
			h.Status = StatusMuatSudahAda
			h.Catatan = append(h.Catatan, "akun "+strings.Join(pemakai, ", ")+" sudah ada - dilewati")
			out = append(out, h)
			continue
		}
		if a.Email != "" {
			if pemakai, err := l.gudang.PemakaiEmail(ctx, a.Email); err != nil {
				return out, err
			} else if len(pemakai) > 0 {
				h.Catatan = append(h.Catatan, fmt.Sprintf("email %q sudah dipakai akun %s - dikosongkan", a.Email,
					strings.Join(pemakai, ", ")))
				a.Email = ""
			}
		}
		if err := l.periksaBaru(ctx, &a); err != nil {
			h.Status = StatusMuatGagal
			h.Catatan = append(h.Catatan, strings.TrimPrefix(err.Error(), "login: "))
			out = append(out, h)
			continue
		}
		if !jalankan {
			h.Status = StatusMuatSiap
			out = append(out, h)
			continue
		}
		hash, err := HashSandi(SandiBawaan)
		if err != nil {
			return out, err
		}
		if err := l.gudang.BuatAkun(ctx, a, hash, false); err != nil {
			h.Status = StatusMuatGagal
			h.Catatan = append(h.Catatan, strings.TrimPrefix(err.Error(), "login: "))
			out = append(out, h)
			continue
		}
		h.Status = StatusMuatDibuat
		out = append(out, h)
	}
	return out, nil
}
