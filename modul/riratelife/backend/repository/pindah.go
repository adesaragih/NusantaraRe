package repository

// Pindah JSON -> flat ringkasan (keputusan work owner 06-10-2026 K-F1/K-F2; alat `backend/alat/pindahflat`, pola
// ricommlife / masterproductnamelife): `M_RATE_LIFE_SUMMARY` (JSON, 339 baris DEV) -> tabel flat `RATE_LIFE_SUMMARY`
// (migrasi inti 926), SEMUA enam kolom view, isi APA ADANYA (teks JSON tanpa dipangkas, seperti `a.JSONDATA.X`).
//
//	-uji       (bawaan) hanya SELECT; nol tulisan; laporan AGREGAT: cacah, panjang maksimum (byte) tiap kolom di
//	           sumber lawan lebar kolom flat, ID baris yang gagal - tidak pernah nilainya.
//	-jalankan  ditolak bila IS_PEGA_PROD=true, bila ada nilai yang TIDAK MUAT (nol pemotongan), bila ada kegagalan lain,
//	           dan bila tabel flat memuat ID yang sama dengan isi BERBEDA (tulisan aplikasi tidak pernah ditimpa). SATU
//	           koneksi terkunci (`db.Koneksi`, SesiPindah) untuk seluruh putaran: setelan sesi lalu transaksi yang lebih
//	           dulu MENGUNCI tabel flat; baris yang sudah ada dan sama dilewati - aman diulang (mis. sesudah Pega
//	           menulis ringkasan baru, risiko cutover MODUL.md).
//
// M_RATE_LIFE_SUMMARY tidak pernah ditulis.

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"nusantarare/inti/backend/db"
)

var (
	// ErrPindahDiProduksi - `-jalankan` di lingkungan IS_PEGA_PROD=true.
	ErrPindahDiProduksi = errors.New("repository: moving R/I rate summaries to the flat table is refused when IS_PEGA_PROD=true")
	// ErrPindahTidakLolos - rencana memuat nilai tidak muat, kegagalan, atau baris flat berbeda; nol tulisan.
	ErrPindahTidakLolos = errors.New("repository: the move plan did not pass; nothing was written")
)

// LebarKolomFlat - lebar (byte) kolom tabel flat `RATE_LIFE_SUMMARY` menurut migrasi 926 - diikat uji ke DDL-nya.
var LebarKolomFlat = map[string]int{"ID": 10, "USEDBY": 500, "TYPE": 100, "MODIFIEDDATE": 50, "OPERATORID": 200, "FLAG": 100}

// PeriksaMode - `-jalankan` ditolak bila IS_PEGA_PROD=true; uji kering selalu boleh.
func PeriksaMode(jalankan, pegaProduksi bool) error {
	if jalankan && pegaProduksi {
		return ErrPindahDiProduksi
	}
	return nil
}

// BarisJSON - satu baris `M_RATE_LIFE_SUMMARY`.
type BarisJSON struct{ ID, JSON string }

// RingkasanFlat - satu baris tabel flat, keenam kolom view (nil = NULL).
type RingkasanFlat struct {
	ID                                           string
	UsedBy, Type, ModifiedDate, OperatorID, Flag *string
}

func samaTeks(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

// Sama - keenam kolom sama persis.
func (r RingkasanFlat) Sama(l RingkasanFlat) bool {
	return r.ID == l.ID && samaTeks(r.UsedBy, l.UsedBy) && samaTeks(r.Type, l.Type) && samaTeks(r.ModifiedDate, l.ModifiedDate) &&
		samaTeks(r.OperatorID, l.OperatorID) && samaTeks(r.Flag, l.Flag)
}

// LaporanPindah - laporan agregat satu putaran.
type LaporanPindah struct {
	Mode        string
	Sumber      int
	AkanDitulis int
	SudahSama   int
	FlatSaja    int
	// PanjangMaks - kolom -> panjang maksimum (byte) di sumber.
	PanjangMaks map[string]int
	// TidakMuat - "ID KOLOM: n byte > lebar" (tanpa nilai).
	TidakMuat []string
	// Berbeda - ID sumber yang sudah ada di tabel flat dengan isi BERBEDA.
	Berbeda []string
	// Gagal - "ID KOLOM: jenis" (tanpa nilai).
	Gagal   []string
	Ditulis bool
}

// BolehDitulis - nol tidak muat, nol kegagalan, nol baris flat berbeda.
func (l LaporanPindah) BolehDitulis() bool {
	return len(l.TidakMuat) == 0 && len(l.Gagal) == 0 && len(l.Berbeda) == 0
}

// Teks - laporan untuk operator (agregat).
func (l LaporanPindah) Teks() string {
	var b strings.Builder
	fmt.Fprintf(&b, "pindahflat R/I Rate Life ringkasan (%s)\n", l.Mode)
	fmt.Fprintf(&b, "  sumber M_RATE_LIFE_SUMMARY: %d\n", l.Sumber)
	fmt.Fprintf(&b, "  baris yang (akan) ditulis: %d\n", l.AkanDitulis)
	fmt.Fprintf(&b, "  sudah ada dan sama (dilewati): %d\n", l.SudahSama)
	fmt.Fprintf(&b, "  hanya di tabel flat (dibiarkan): %d\n", l.FlatSaja)
	fmt.Fprintf(&b, "  panjang maksimum sumber (byte) / lebar kolom:\n")
	for _, k := range KolomViewRingkasan {
		fmt.Fprintf(&b, "    %-13s %4d / %d\n", k, l.PanjangMaks[k], LebarKolomFlat[k])
	}
	fmt.Fprintf(&b, "  tidak muat: %d\n", len(l.TidakMuat))
	for _, g := range l.TidakMuat {
		fmt.Fprintf(&b, "    %s\n", g)
	}
	fmt.Fprintf(&b, "  tabel flat berbeda: %d %s\n", len(l.Berbeda), strings.Join(l.Berbeda, " "))
	fmt.Fprintf(&b, "  gagal: %d\n", len(l.Gagal))
	for _, g := range l.Gagal {
		fmt.Fprintf(&b, "    %s\n", g)
	}
	fmt.Fprintf(&b, "  ditulis: %v\n", l.Ditulis)
	return b.String()
}

// NilaiJSON - nilai kunci seperti `a.JSONDATA.<kunci>` view: teks APA ADANYA (tanpa dipangkas), angka dan boolean
// sebagai teks literalnya; kunci tidak ada / null / teks kosong / objek / larik = nil. ok=false bila JSONDATA bukan objek JSON.
func NilaiJSON(jsonData, kunci string) (*string, bool) {
	obj := map[string]json.RawMessage{}
	if t := strings.TrimSpace(jsonData); t == "" {
		return nil, true
	}
	if err := json.Unmarshal([]byte(jsonData), &obj); err != nil || obj == nil {
		return nil, false
	}
	r, ada := obj[kunci]
	if !ada {
		return nil, true
	}
	var s string
	if json.Unmarshal(r, &s) == nil {
		if s == "" {
			return nil, true // Oracle: teks kosong = NULL (view pun menjawab NULL)
		}
		return &s, true
	}
	t := strings.TrimSpace(string(r))
	if t == "null" || strings.HasPrefix(t, "{") || strings.HasPrefix(t, "[") {
		return nil, true
	}
	return &t, true
}

// RencanaPindah - laporan dan baris yang akan disisipkan, dari sumber JSON dan isi tabel flat saat ini. Murni.
func RencanaPindah(sumber []BarisJSON, flat []RingkasanFlat) (LaporanPindah, []RingkasanFlat) {
	lap := LaporanPindah{Sumber: len(sumber), PanjangMaks: map[string]int{}}
	ada := map[string]RingkasanFlat{}
	for _, f := range flat {
		ada[f.ID] = f
	}
	dipakai := map[string]bool{}
	var tulis []RingkasanFlat
	for _, b := range sumber {
		id := b.ID
		gagal := false
		catat := func(kolom string, v *string) {
			if v == nil {
				return
			}
			n := len(*v)
			lap.PanjangMaks[kolom] = max(lap.PanjangMaks[kolom], n)
			if n > LebarKolomFlat[kolom] {
				lap.TidakMuat = append(lap.TidakMuat, fmt.Sprintf("%s %s: %d byte > %d", id, kolom, n, LebarKolomFlat[kolom]))
				gagal = true
			}
		}
		if strings.TrimSpace(id) == "" {
			lap.Gagal = append(lap.Gagal, "(kosong) ID: kosong")
			continue
		}
		catat("ID", &id)
		r := RingkasanFlat{ID: id}
		jsonSah := true
		for _, x := range []struct {
			kolom string
			ke    **string
		}{{"USEDBY", &r.UsedBy}, {"TYPE", &r.Type}, {"MODIFIEDDATE", &r.ModifiedDate}, {"OPERATORID", &r.OperatorID}, {"FLAG", &r.Flag}} {
			v, ok := NilaiJSON(b.JSON, x.kolom)
			if !ok {
				jsonSah = false
				break
			}
			*x.ke = v
			catat(x.kolom, v)
		}
		if !jsonSah {
			lap.Gagal = append(lap.Gagal, id+" JSONDATA: bukan objek JSON")
			continue
		}
		if gagal {
			continue
		}
		if dipakai[id] {
			lap.Gagal = append(lap.Gagal, id+" ID: kembar di sumber")
			continue
		}
		dipakai[id] = true
		if f, sudah := ada[id]; sudah {
			if f.Sama(r) {
				lap.SudahSama++
			} else {
				lap.Berbeda = append(lap.Berbeda, id)
			}
			continue
		}
		tulis = append(tulis, r)
	}
	for id := range ada {
		if !dipakai[id] {
			lap.FlatSaja++
		}
	}
	sort.Strings(lap.Berbeda)
	lap.AkanDitulis = len(tulis)
	return lap, tulis
}

// SqlSesiNLS - setelan sesi alat pindah (kolom semua teks; lapis kedua untuk konversi implisit), di koneksi yang SAMA
// dengan transaksinya (SesiPindah).
const SqlSesiNLS = `ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '.,'`

// SqlSumberJSON - seluruh ringkasan JSON warisan (alat pindah).
func SqlSumberJSON(t string) string {
	return fmt.Sprintf(`SELECT ID, %s FROM %s ORDER BY ID`, KolomJSON, t)
}

// SqlSemuaRingkasanFlat - seluruh tabel flat, keenam kolom.
func SqlSemuaRingkasanFlat(t string) string {
	return fmt.Sprintf(`SELECT ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG FROM %s ORDER BY ID`, t)
}

// SqlSisipRingkasanFlat - satu baris pindahan, keenam kolom apa adanya.
func SqlSisipRingkasanFlat(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID, FLAG) VALUES (:1, :2, :3, :4, :5, :6)`, t)
}

// SqlKunciRingkasan - kunci tabel flat sepanjang transaksi alat pindah.
func SqlKunciRingkasan(t string) string { return fmt.Sprintf(`LOCK TABLE %s IN EXCLUSIVE MODE`, t) }

// SesiPindah - SATU koneksi terkunci (`db.Koneksi`): setelan sesi lalu transaksi di koneksi yang SAMA. Pemanggil
// menutup keduanya.
func (g *Gudang) SesiPindah(ctx context.Context) (*db.Koneksi, *db.Tx, error) {
	if err := db.PeriksaSQL(SqlSesiNLS); err != nil {
		return nil, nil, err
	}
	kon, err := g.db.Koneksi(ctx)
	if err != nil {
		return nil, nil, err
	}
	if _, err := kon.ExecContext(ctx, SqlSesiNLS); err != nil {
		_ = kon.Close()
		return nil, nil, bungkus(err, "menyetel sesi pindah")
	}
	tx, err := kon.Mulai(ctx)
	if err != nil {
		_ = kon.Close()
		return nil, nil, err
	}
	return kon, tx, nil
}

func nilaiArg(v *string) any {
	if v == nil {
		return nil
	}
	return *v
}

func teksNull(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

// bacaRencana - sumber JSON dan isi tabel flat di transaksi tx.
func (g *Gudang) bacaRencana(ctx context.Context, tx *db.Tx, n nama) ([]BarisJSON, []RingkasanFlat, error) {
	q := SqlSumberJSON(n.ringkasanJSON)
	if err := siap(TabelRingkasanJSON, q); err != nil {
		return nil, nil, err
	}
	rows, err := g.dari(tx).QueryContext(ctx, q)
	if err != nil {
		return nil, nil, bungkus(err, "membaca M_RATE_LIFE_SUMMARY")
	}
	var sumber []BarisJSON
	for rows.Next() {
		var id sql.NullString
		var p db.PindaiTeksPanjang
		if err := rows.Scan(&id, &p); err != nil {
			_ = rows.Close()
			return nil, nil, bungkus(err, "memindai M_RATE_LIFE_SUMMARY")
		}
		sumber = append(sumber, BarisJSON{ID: id.String, JSON: p.Teks()})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, nil, bungkus(err, "membaca M_RATE_LIFE_SUMMARY")
	}
	_ = rows.Close()
	q = SqlSemuaRingkasanFlat(n.tabelRingkasan)
	if err := siap(TabelRingkasan, q); err != nil {
		return nil, nil, err
	}
	rows, err = g.dari(tx).QueryContext(ctx, q)
	if err != nil {
		return nil, nil, bungkus(err, "membaca RATE_LIFE_SUMMARY")
	}
	defer func() { _ = rows.Close() }()
	var flat []RingkasanFlat
	for rows.Next() {
		var v [6]sql.NullString
		if err := rows.Scan(&v[0], &v[1], &v[2], &v[3], &v[4], &v[5]); err != nil {
			return nil, nil, bungkus(err, "memindai RATE_LIFE_SUMMARY")
		}
		flat = append(flat, RingkasanFlat{ID: v[0].String, UsedBy: teksNull(v[1]), Type: teksNull(v[2]),
			ModifiedDate: teksNull(v[3]), OperatorID: teksNull(v[4]), Flag: teksNull(v[5])})
	}
	return sumber, flat, bungkus(rows.Err(), "membaca RATE_LIFE_SUMMARY")
}

// PindahFlat - satu putaran alat pindah (lihat kepala berkas): SELURUHNYA di satu koneksi (SesiPindah).
func (g *Gudang) PindahFlat(ctx context.Context, jalankan bool) (LaporanPindah, error) {
	mode := "uji"
	if jalankan {
		mode = "jalankan"
	}
	if err := PeriksaMode(jalankan, g.db.PegaProduksi()); err != nil {
		return LaporanPindah{Mode: mode}, err
	}
	n, err := g.nama()
	if err != nil {
		return LaporanPindah{Mode: mode}, err
	}
	kon, tx, err := g.SesiPindah(ctx)
	if err != nil {
		return LaporanPindah{Mode: mode}, err
	}
	defer func() { _ = kon.Close() }()
	defer func() { _ = tx.Rollback() }()
	if jalankan {
		if _, err := g.tulis(ctx, tx, TabelRingkasan, SqlKunciRingkasan(n.tabelRingkasan), "mengunci RATE_LIFE_SUMMARY"); err != nil {
			return LaporanPindah{Mode: mode}, err
		}
	}
	sumber, flat, err := g.bacaRencana(ctx, tx, n)
	if err != nil {
		return LaporanPindah{Mode: mode}, err
	}
	lap, tulis := RencanaPindah(sumber, flat)
	lap.Mode = mode
	if !jalankan {
		return lap, nil
	}
	if !lap.BolehDitulis() {
		return lap, ErrPindahTidakLolos
	}
	for _, r := range tulis {
		if _, err := g.tulis(ctx, tx, TabelRingkasan, SqlSisipRingkasanFlat(n.tabelRingkasan), "memindah ringkasan", r.ID,
			nilaiArg(r.UsedBy), nilaiArg(r.Type), nilaiArg(r.ModifiedDate), nilaiArg(r.OperatorID), nilaiArg(r.Flag)); err != nil {
			return lap, err
		}
	}
	if err := tx.Commit(); err != nil {
		return lap, fmt.Errorf("repository: menutup transaksi pindah: %w", err)
	}
	lap.Ditulis = true
	return lap, nil
}
