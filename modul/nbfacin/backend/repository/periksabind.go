package repository

// Pemeriksa bind sebelum INSERT objek (bug DEV 10-10-2026: Save Spreading Aneka ORA-01722 di T_SPREADINGLIST masih muncul
// setelah kolomArrayBertipe). Setiap nilai yang masuk ke angkaMasuk(...) wajib cocok PolaAngkaMasuk, dan setiap bind polos ke
// kolom nomor (ID, PARENT_ID, SEQ_NO - NUMBER, NOT NULL) wajib bilangan bulat; bila tidak, galat ErrBindTidakSah menyebut
// tabel, kolom, baris, dan BENTUK nilai (nilai uang mentah tidak dicetak). Bila Oracle masih menolak padahal semua bind lolos,
// masalahnya di tipe / driver, bukan isi. Hanya menolak yang memang akan ditolak mask Oracle (PolaAngkaMasuk = mask
// fmtAngkaMasuk, TestPolaAngkaMasuk); jumlah nilai per baris wajib = jumlah bind teks SQL.

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ErrBindTidakSah - nilai bind INSERT yang akan ditolak Oracle (ORA-01722) atau baris yang tidak sejalan dengan teks SQL.
var ErrBindTidakSah = errors.New("repository: nilai bind tidak sah")

// kolomNomor - kolom NUMBER(19) / NUMBER(5) berbind polos di INSERT objek (ID / induk / urutan; NOT NULL, migrasi 186-198 / 214).
// Kecuali PARENT_ID tabel ber-induk T_GENERAL_POLIS (indukTeks) - VARCHAR2(32) ID case, tidak diperiksa.
var kolomNomor = map[string]bool{"ID": true, "PARENT_ID": true, "SEQ_NO": true}

// indukTeks - tabel objek yang PARENT_ID-nya VARCHAR2(32) = T_GENERAL_POLIS.ID: T_LOCATIONLIST (migrasi 186 L17), T_CARGOLIST
// (256, tiket 74).
var indukTeks = map[string]bool{TabelLocationList: true, TabelCargoList: true}

var (
	polaBindPolos  = regexp.MustCompile(`^:(\d+)$`)
	polaBindAngka  = regexp.MustCompile(`^TO_NUMBER\(:(\d+), 'FM`)
	polaBindApapun = regexp.MustCompile(`:(\d+)`)
	polaNomor      = regexp.MustCompile(`^[0-9]{1,19}$`)
)

// kolomBind - kolom tujuan satu bind INSERT; angka = lewat angkaMasuk, selain itu kolom nomor berbind polos.
type kolomBind struct {
	nama  string
	angka bool
}

// petaBindSisip - `INSERT INTO t (k1, …) VALUES (e1, …)`: nomor bind (0-based) -> kolom yang diperiksa, dan jumlah bind teks
// SQL (nomor terbesar). Ekspresi selain bind polos / angkaMasuk (subkueri, NVL, literal) tidak dipetakan.
func petaBindSisip(q, tabel string) (map[int]kolomBind, int, error) {
	n := 0
	for _, m := range polaBindApapun.FindAllStringSubmatch(q, -1) {
		if k, _ := strconv.Atoi(m[1]); k > n {
			n = k
		}
	}
	buka := strings.Index(q, "(")
	tengah := strings.Index(q, ") VALUES (")
	if buka < 0 || tengah < buka || !strings.HasSuffix(q, ")") {
		return nil, n, fmt.Errorf("%w: teks SQL bukan INSERT … VALUES (…)", ErrBindTidakSah)
	}
	kolom := strings.Split(q[buka+1:tengah], ",")
	nilai := pisahPuncak(q[tengah+len(") VALUES (") : len(q)-1])
	if len(kolom) != len(nilai) {
		return nil, n, fmt.Errorf("%w: %d kolom, %d nilai di teks SQL", ErrBindTidakSah, len(kolom), len(nilai))
	}
	peta := map[int]kolomBind{}
	for i, e := range nilai {
		nama := strings.TrimSpace(kolom[i])
		e = strings.TrimSpace(e)
		if m := polaBindAngka.FindStringSubmatch(e); m != nil {
			k, _ := strconv.Atoi(m[1])
			peta[k-1] = kolomBind{nama: nama, angka: true}
		} else if m := polaBindPolos.FindStringSubmatch(e); m != nil && kolomNomor[nama] &&
			!(nama == "PARENT_ID" && indukTeks[tabel]) {
			k, _ := strconv.Atoi(m[1])
			peta[k-1] = kolomBind{nama: nama}
		}
	}
	return peta, n, nil
}

// pisahPuncak - pisah pada koma di luar kurung dan di luar literal '…'.
func pisahPuncak(s string) []string {
	var hasil []string
	dalam, kutip, awal := 0, false, 0
	for i, r := range s {
		switch {
		case r == '\'':
			kutip = !kutip
		case kutip:
		case r == '(':
			dalam++
		case r == ')':
			dalam--
		case r == ',' && dalam == 0:
			hasil = append(hasil, s[awal:i])
			awal = i + 1
		}
	}
	return append(hasil, s[awal:])
}

// periksaBindBaris - baris (nilai per bind) untuk teks SQL q; awal = nomor baris pertama (0-based) di kelompoknya.
func periksaBindBaris(q, tabel string, baris [][]any, awal int) error {
	peta, n, err := petaBindSisip(q, tabel)
	if err != nil {
		return fmt.Errorf("%w (%s)", err, tabel)
	}
	for r, arg := range baris {
		if len(arg) != n {
			return fmt.Errorf("%w: %s baris %d: %d nilai, teks SQL %d bind", ErrBindTidakSah, tabel, awal+r+1, len(arg), n)
		}
		for i, k := range peta {
			if bentuk, sah := periksaNilaiBind(arg[i], k.angka); !sah {
				return fmt.Errorf("%w: %s kolom %s baris %d: %s", ErrBindTidakSah, tabel, k.nama, awal+r+1, bentuk)
			}
		}
	}
	return nil
}

// periksaNilaiBind - angka: NULL / teks cocok PolaAngkaMasuk; nomor: bilangan bulat (Go int atau teks digit) tidak NULL.
func periksaNilaiBind(x any, angka bool) (string, bool) {
	if n, ok := x.(sql.NullString); ok {
		if !n.Valid {
			x = nil
		} else {
			x = n.String
		}
	}
	switch v := x.(type) {
	case nil:
		return "NULL (kolom NOT NULL)", angka
	case string:
		if angka && PolaAngkaMasuk.MatchString(v) || !angka && polaNomor.MatchString(v) {
			return "", true
		}
		return bentukTeks(v), false
	case int, int64, int32:
		return "", !angka
	}
	return fmt.Sprintf("tipe Go %T", x), false
}

// bentukTeks - bentuk nilai tanpa nilainya: kosong / memuat E / koma / spasi / karakter lain, panjang bulat dan desimal.
func bentukTeks(s string) string {
	if s == "" {
		return "teks kosong"
	}
	var ciri []string
	if strings.ContainsAny(s, "eE") {
		ciri = append(ciri, "memuat E")
	}
	if strings.Contains(s, ",") {
		ciri = append(ciri, "memuat koma")
	}
	if strings.ContainsAny(s, " \t") {
		ciri = append(ciri, "memuat spasi")
	}
	if strings.HasPrefix(s, "-") || strings.HasPrefix(s, "+") {
		ciri = append(ciri, "bertanda "+s[:1])
	}
	t := strings.TrimLeft(s, "+-")
	bulat, desimal, adaTitik := strings.Cut(t, ".")
	if strings.Count(t, ".") > 1 {
		ciri = append(ciri, fmt.Sprintf("%d titik", strings.Count(t, ".")))
	}
	lain := 0
	for _, r := range t {
		if (r < '0' || r > '9') && r != '.' {
			lain++
		}
	}
	if lain > 0 {
		ciri = append(ciri, fmt.Sprintf("%d karakter bukan angka", lain))
	}
	p := fmt.Sprintf("panjang bulat %d", len(bulat))
	if adaTitik {
		p += fmt.Sprintf(", desimal %d", len(desimal))
	}
	return strings.Join(append(ciri, p), ", ")
}
