package repository

// Tabel kolom bersama untuk tabel anak rancangan yang medannya banyak (T_COVERAGELIST tiket 43, T_DEDUCTIBLELIST tiket
// 45): SQL sisip, bind, kolom baca, dan pemindaian baris DIBANGUN dari SATU daftar kolom per tabel - tidak ada pemasangan
// nama ke nilai menurut urutan yang ditulis dua kali (CLAUDE.md §4a jebakan 6).

import (
	"database/sql"
	"strconv"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
)

// jenisKolom - cara satu kolom data ditulis / dibaca.
type jenisKolom int

const (
	jenisTeks      jenisKolom = iota // VARCHAR2 teks
	jenisAngka                       // NUMBER(38,8) - angkaMasuk / angkaKeluar
	jenisTeksAngka                   // VARCHAR2 berisi teks desimal (rancangan) - ikatDesimal / bacaDesimal tanpa TO_NUMBER
	// jenisKode - NUMBER(5) berisi kode bulat (pyStandardValue): teks di aplikasi, TO_NUMBER / TO_CHAR polos di SQL. Aman
	// tanpa format NLS HANYA karena isinya digit saja (tanpa pemisah desimal / ribuan) - services memeriksanya dulu
	// (polaKodeDeductible).
	jenisKode
)

// kolomData - satu kolom data tabel dan medan struct T-nya: `teks` untuk jenisTeks / jenisKode, `des` untuk jenisAngka /
// jenisTeksAngka (tepat satu yang terisi).
type kolomData[T any] struct {
	nama  string
	jenis jenisKolom
	teks  func(*T) *string
	des   func(*T) **apd.Decimal
}

func kt[T any](nama string, f func(*T) *string) kolomData[T] {
	return kolomData[T]{nama: nama, jenis: jenisTeks, teks: f}
}

func kk[T any](nama string, f func(*T) *string) kolomData[T] {
	return kolomData[T]{nama: nama, jenis: jenisKode, teks: f}
}

func ka[T any](nama string, jenis jenisKolom, f func(*T) **apd.Decimal) kolomData[T] {
	return kolomData[T]{nama: nama, jenis: jenis, des: f}
}

// kunciBaca - ekspresi terpilih kolom ber-alias `alias` (juga kunci peta hasil pindai).
func (k kolomData[T]) kunciBaca(alias string) string {
	switch k.jenis {
	case jenisAngka:
		return angkaKeluar(alias + "." + k.nama)
	case jenisKode:
		return "TO_CHAR(" + alias + "." + k.nama + ")"
	}
	return alias + "." + k.nama
}

// nilaiSisip - placeholder bind ke-n kolom pada INSERT.
func (k kolomData[T]) nilaiSisip(n int) string {
	b := ":" + strconv.Itoa(n)
	switch k.jenis {
	case jenisAngka:
		return angkaMasuk(b)
	case jenisKode:
		return "TO_NUMBER(" + b + ")"
	}
	return b
}

// argData - bind kolom (urutan daftar yang sama dengan sqlSisipData): teks kosong -> NULL.
func argData[T any](kolom []kolomData[T], v T) []any {
	arg := make([]any, 0, len(kolom))
	for _, k := range kolom {
		if k.teks != nil {
			arg = append(arg, db.KosongJadiNil(*k.teks(&v)))
		} else {
			arg = append(arg, ikatDesimal(*k.des(&v)))
		}
	}
	return arg
}

// kolomSisipData - nama kolom dan placeholder kolom data, bind mulai `awal`.
func kolomSisipData[T any](kolom []kolomData[T], awal int) (nama, nilai []string) {
	for i, k := range kolom {
		nama = append(nama, k.nama)
		nilai = append(nilai, k.nilaiSisip(awal+i))
	}
	return nama, nilai
}

// kolomBacaData - ekspresi terpilih kolom data ber-alias.
func kolomBacaData[T any](kolom []kolomData[T], alias string) []string {
	hasil := make([]string, 0, len(kolom))
	for _, k := range kolom {
		hasil = append(hasil, k.kunciBaca(alias))
	}
	return hasil
}

// pindaiBaris - satu baris hasil kueri -> peta kunci ekspresi -> nilai teks.
func pindaiBaris(baris *sql.Rows, kunci []string) (map[string]*sql.NullString, error) {
	teks := make(map[string]*sql.NullString, len(kunci))
	tujuan := make([]any, len(kunci))
	for i, k := range kunci {
		teks[k] = &sql.NullString{}
		tujuan[i] = teks[k]
	}
	return teks, baris.Scan(tujuan...)
}

// isiData - nilai teks hasil pindai -> struct T (desimal lewat bacaDesimal; `label` untuk pesan galat).
func isiData[T any](kolom []kolomData[T], alias, label string, teks map[string]*sql.NullString) (T, error) {
	var v T
	for _, k := range kolom {
		n := teks[k.kunciBaca(alias)]
		if k.teks != nil {
			*k.teks(&v) = n.String
			continue
		}
		d, err := bacaDesimal(label, k.nama, n)
		if err != nil {
			return v, err
		}
		*k.des(&v) = d
	}
	return v, nil
}
