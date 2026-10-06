package repository

// Pindah JSON -> flat ringkasan (keputusan work owner 06-10-2026 K-F1/K-F2; alat `backend/alat/pindahflat`, pola
// ricommlife / masterproductnamelife): `M_RATE_LIFE_SUMMARY` (JSON warisan; data DEV MASIH BERUBAH - cacahnya dibaca
// saat berjalan, tidak pernah ditanam) -> tabel flat `RATE_LIFE_SUMMARY` (migrasi inti 926), kolom view selain FLAG, isi
// APA ADANYA (teks JSON tanpa dipangkas, seperti `a.JSONDATA.X`) - termasuk TYPE. FLAG TIDAK disalin (RALAT R5: tidak
// digunakan, keputusan work owner 06-10-2026; nilainya tetap di JSON).
//
//	-uji       (bawaan) hanya SELECT; nol tulisan; laporan AGREGAT: cacah sumber dan flat SAAT itu, baru / berubah /
//	           sama / konflik / dilewati, panjang maksimum (byte) tiap kolom lawan lebar kolom flat, ID baris - tidak
//	           pernah nilainya.
//	-jalankan  ditolak bila IS_PEGA_PROD=true, bila ada nilai yang TIDAK MUAT (nol pemotongan), dan bila ada kegagalan.
//	           SATU koneksi terkunci (`db.Koneksi`, SesiPindah) untuk seluruh putaran: setelan sesi lalu transaksi yang
//	           lebih dulu MENGUNCI tabel flat.
//
// ATURAN DELTA (keputusan work owner 06-10-2026: cutover delta) - dibandingkan per ID terhadap isi tabel flat SAAT itu:
//
//	baru      ID belum ada di flat                                  -> disisip
//	          ... kecuali `-sejak` diberikan dan MODIFIEDDATE sumber TIDAK lebih baru dari batas itu (atau kosong /
//	          tidak terbaca): `dilewati` - ringkasan lama yang tidak ada di flat dianggap DIHAPUS aplikasi, tidak
//	          dihidupkan lagi, dilaporkan.
//	sama      kelima kolom sama persis                              -> dilewati
//	berubah   beda, dan MODIFIEDDATE sumber LEBIH BARU dari flat     -> flat diperbarui kelima kolom (Pega lebih baru)
//	konflik   beda, dan MODIFIEDDATE flat sama / lebih baru / salah satunya kosong atau tidak terbaca -> TIDAK ditimpa,
//	          dilaporkan (tulisan aplikasi sesudah pemindahan tidak pernah ditimpa diam-diam)
//	flat saja ID hanya di flat (ringkasan baru aplikasi)            -> dibiarkan
//
// Konflik dan dilewati tidak menahan `-jalankan` (yang lain tetap ditulis); keduanya disebut ID-nya untuk diperiksa WO.
// Laporan menyebut `batas delta berikutnya` = MODIFIEDDATE sumber terbaru, untuk `-sejak` putaran berikutnya.
// M_RATE_LIFE_SUMMARY tidak pernah ditulis.

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
)

var (
	// ErrPindahDiProduksi - `-jalankan` di lingkungan IS_PEGA_PROD=true.
	ErrPindahDiProduksi = errors.New("repository: moving R/I rate summaries to the flat table is refused when IS_PEGA_PROD=true")
	// ErrPindahTidakLolos - rencana memuat nilai tidak muat atau kegagalan; nol tulisan.
	ErrPindahTidakLolos = errors.New("repository: the move plan did not pass; nothing was written")
	// ErrSejakTidakTerbaca - `-sejak` bukan stempel Pega.
	ErrSejakTidakTerbaca = errors.New("repository: -sejak must be a Pega timestamp like 20261006T040628.169 GMT")
)

// LebarKolomFlat - lebar (byte) kolom tabel flat `RATE_LIFE_SUMMARY` menurut migrasi 926 - diikat uji ke DDL-nya.
var LebarKolomFlat = map[string]int{"ID": 10, "USEDBY": 500, "TYPE": 100, "MODIFIEDDATE": 50, "OPERATORID": 200}

// PeriksaMode - `-jalankan` ditolak bila IS_PEGA_PROD=true; uji kering selalu boleh.
func PeriksaMode(jalankan, pegaProduksi bool) error {
	if jalankan && pegaProduksi {
		return ErrPindahDiProduksi
	}
	return nil
}

// BarisJSON - satu baris `M_RATE_LIFE_SUMMARY`.
type BarisJSON struct{ ID, JSON string }

// RingkasanFlat - satu baris tabel flat, kelima kolom (nil = NULL).
type RingkasanFlat struct {
	ID                                     string
	UsedBy, Type, ModifiedDate, OperatorID *string
}

func samaTeks(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

// Sama - kelima kolom sama persis.
func (r RingkasanFlat) Sama(l RingkasanFlat) bool {
	return r.ID == l.ID && samaTeks(r.UsedBy, l.UsedBy) && samaTeks(r.Type, l.Type) && samaTeks(r.ModifiedDate, l.ModifiedDate) &&
		samaTeks(r.OperatorID, l.OperatorID)
}

// WaktuPega - stempel MODIFIEDDATE (`20261006T040628.169 GMT`, atau `YYYYMMDD`) -> waktu; ok=false bila kosong /
// bentuk lain.
func WaktuPega(s *string) (time.Time, bool) {
	if s == nil {
		return time.Time{}, false
	}
	t := strings.TrimSpace(*s)
	if w, err := time.Parse("20060102T150405.000 MST", t); err == nil {
		return w, true
	}
	if w, err := time.Parse("20060102", t); err == nil {
		return w, true
	}
	return time.Time{}, false
}

// LaporanPindah - laporan agregat satu putaran.
type LaporanPindah struct {
	Mode string
	// CacahSumber / CacahFlat - cacah baris SAAT dijalankan (sebelum menulis); CacahFlatSesudah - sesudah commit.
	CacahSumber, CacahFlat, CacahFlatSesudah int
	Baru, Berubah, Sama                      []string
	// Konflik - beda, flat sama / lebih baru / stempel tak terbaca: TIDAK ditimpa.
	Konflik []string
	// Dilewati - tidak ada di flat, MODIFIEDDATE sumber tidak lebih baru dari `-sejak`: dianggap dihapus aplikasi.
	Dilewati []string
	FlatSaja int
	// Sejak - batas `-sejak` yang dipakai ("" = tanpa batas); BatasBerikut - MODIFIEDDATE sumber terbaru.
	Sejak, BatasBerikut string
	// PanjangMaks - kolom -> panjang maksimum (byte) di sumber.
	PanjangMaks map[string]int
	// TidakMuat - "ID KOLOM: n byte > lebar" (tanpa nilai).
	TidakMuat []string
	// Gagal - "ID KOLOM: jenis" (tanpa nilai).
	Gagal   []string
	Ditulis bool
}

// BolehDitulis - nol tidak muat, nol kegagalan. Konflik dan dilewati dilaporkan, tidak menahan.
func (l LaporanPindah) BolehDitulis() bool { return len(l.TidakMuat) == 0 && len(l.Gagal) == 0 }

func daftarID(id []string) string {
	if len(id) == 0 {
		return ""
	}
	return " " + strings.Join(id, " ")
}

// Teks - laporan untuk operator (agregat; ID dan nama kolom, tidak pernah nilai).
func (l LaporanPindah) Teks() string {
	var b strings.Builder
	fmt.Fprintf(&b, "pindahflat R/I Rate Life ringkasan (%s)\n", l.Mode)
	fmt.Fprintf(&b, "  cacah saat dijalankan: sumber M_RATE_LIFE_SUMMARY %d, flat RATE_LIFE_SUMMARY %d\n", l.CacahSumber, l.CacahFlat)
	if l.Sejak != "" {
		fmt.Fprintf(&b, "  -sejak: %s\n", l.Sejak)
	}
	fmt.Fprintf(&b, "  baru (disisip): %d\n", len(l.Baru))
	fmt.Fprintf(&b, "  berubah (sumber lebih baru, diperbarui): %d%s\n", len(l.Berubah), daftarID(l.Berubah))
	fmt.Fprintf(&b, "  sama (dilewati): %d\n", len(l.Sama))
	fmt.Fprintf(&b, "  konflik (TIDAK ditimpa, periksa): %d%s\n", len(l.Konflik), daftarID(l.Konflik))
	fmt.Fprintf(&b, "  dilewati, tidak lebih baru dari -sejak dan tidak ada di flat (dihapus aplikasi?): %d%s\n", len(l.Dilewati), daftarID(l.Dilewati))
	fmt.Fprintf(&b, "  hanya di tabel flat (dibiarkan): %d\n", l.FlatSaja)
	fmt.Fprintf(&b, "  panjang maksimum sumber (byte) / lebar kolom:\n")
	for _, k := range KolomRingkasanFlat {
		fmt.Fprintf(&b, "    %-13s %4d / %d\n", k, l.PanjangMaks[k], LebarKolomFlat[k])
	}
	fmt.Fprintf(&b, "  tidak muat: %d\n", len(l.TidakMuat))
	for _, g := range l.TidakMuat {
		fmt.Fprintf(&b, "    %s\n", g)
	}
	fmt.Fprintf(&b, "  gagal: %d\n", len(l.Gagal))
	for _, g := range l.Gagal {
		fmt.Fprintf(&b, "    %s\n", g)
	}
	if l.BatasBerikut != "" {
		fmt.Fprintf(&b, "  batas delta berikutnya (-sejak): %s\n", l.BatasBerikut)
	}
	if l.Ditulis {
		fmt.Fprintf(&b, "  cacah flat sesudah: %d\n", l.CacahFlatSesudah)
	}
	fmt.Fprintf(&b, "  ditulis: %v\n", l.Ditulis)
	return b.String()
}

// NilaiJSON - nilai kunci seperti `a.JSONDATA.<kunci>` view: teks APA ADANYA (tanpa dipangkas), angka dan boolean
// sebagai teks literalnya; kunci tidak ada / null / teks kosong / objek / larik = nil. ok=false bila JSONDATA bukan
// objek JSON.
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

// RencanaPindah - laporan, baris yang disisip, dan baris yang diperbarui (aturan delta di kepala berkas), dari sumber
// JSON dan isi tabel flat saat ini. sejak = "" tanpa batas. Murni.
func RencanaPindah(sumber []BarisJSON, flat []RingkasanFlat, sejak string) (LaporanPindah, []RingkasanFlat, []RingkasanFlat, error) {
	lap := LaporanPindah{CacahSumber: len(sumber), CacahFlat: len(flat), PanjangMaks: map[string]int{}, Sejak: strings.TrimSpace(sejak)}
	var batas time.Time
	if lap.Sejak != "" {
		w, ok := WaktuPega(&lap.Sejak)
		if !ok {
			return lap, nil, nil, ErrSejakTidakTerbaca
		}
		batas = w
	}
	ada := map[string]RingkasanFlat{}
	for _, f := range flat {
		ada[f.ID] = f
	}
	dipakai := map[string]bool{}
	var sisip, ubah []RingkasanFlat
	var terbaru time.Time
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
		}{{"USEDBY", &r.UsedBy}, {"TYPE", &r.Type}, {"MODIFIEDDATE", &r.ModifiedDate}, {"OPERATORID", &r.OperatorID}} {
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
		wSumber, okSumber := WaktuPega(r.ModifiedDate)
		if okSumber && wSumber.After(terbaru) {
			terbaru, lap.BatasBerikut = wSumber, *r.ModifiedDate
		}
		f, sudah := ada[id]
		switch {
		case !sudah && lap.Sejak != "" && (!okSumber || !wSumber.After(batas)):
			lap.Dilewati = append(lap.Dilewati, id)
		case !sudah:
			lap.Baru = append(lap.Baru, id)
			sisip = append(sisip, r)
		case f.Sama(r):
			lap.Sama = append(lap.Sama, id)
		default:
			wFlat, okFlat := WaktuPega(f.ModifiedDate)
			if okSumber && okFlat && wSumber.After(wFlat) {
				lap.Berubah = append(lap.Berubah, id)
				ubah = append(ubah, r)
			} else {
				lap.Konflik = append(lap.Konflik, id)
			}
		}
	}
	for id := range ada {
		if !dipakai[id] {
			lap.FlatSaja++
		}
	}
	for _, d := range [][]string{lap.Baru, lap.Berubah, lap.Sama, lap.Konflik, lap.Dilewati} {
		sort.Strings(d)
	}
	return lap, sisip, ubah, nil
}

// SqlSesiNLS - setelan sesi alat pindah (kolom semua teks; lapis kedua untuk konversi implisit), di koneksi yang SAMA
// dengan transaksinya (SesiPindah).
const SqlSesiNLS = `ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '.,'`

// SqlSumberJSON - seluruh ringkasan JSON warisan (alat pindah).
func SqlSumberJSON(t string) string {
	return fmt.Sprintf(`SELECT ID, %s FROM %s ORDER BY ID`, KolomJSON, t)
}

// SqlSemuaRingkasanFlat - seluruh tabel flat, kelima kolom.
func SqlSemuaRingkasanFlat(t string) string {
	return fmt.Sprintf(`SELECT ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID FROM %s ORDER BY ID`, t)
}

// SqlSisipRingkasanFlat - satu baris pindahan, kelima kolom apa adanya.
func SqlSisipRingkasanFlat(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID) VALUES (:1, :2, :3, :4, :5)`, t)
}

// SqlPerbaruiRingkasanFlat - delta `berubah`: kelima kolom dari sumber yang lebih baru.
func SqlPerbaruiRingkasanFlat(t string) string {
	return fmt.Sprintf(`UPDATE %s SET USEDBY = :1, TYPE = :2, MODIFIEDDATE = :3, OPERATORID = :4 WHERE ID = :5`, t)
}

// SqlCacah - cacah baris.
func SqlCacah(t string) string { return fmt.Sprintf(`SELECT COUNT(*) FROM %s`, t) }

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
		var v [5]sql.NullString
		if err := rows.Scan(&v[0], &v[1], &v[2], &v[3], &v[4]); err != nil {
			return nil, nil, bungkus(err, "memindai RATE_LIFE_SUMMARY")
		}
		flat = append(flat, RingkasanFlat{ID: v[0].String, UsedBy: teksNull(v[1]), Type: teksNull(v[2]),
			ModifiedDate: teksNull(v[3]), OperatorID: teksNull(v[4])})
	}
	return sumber, flat, bungkus(rows.Err(), "membaca RATE_LIFE_SUMMARY")
}

// PindahFlat - satu putaran alat pindah (lihat kepala berkas): SELURUHNYA di satu koneksi (SesiPindah). sejak = ""
// pemindahan penuh; selain itu batas delta (`-sejak`).
func (g *Gudang) PindahFlat(ctx context.Context, jalankan bool, sejak string) (LaporanPindah, error) {
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
	lap, sisip, ubah, err := RencanaPindah(sumber, flat, sejak)
	lap.Mode = mode
	if err != nil || !jalankan {
		return lap, err
	}
	if !lap.BolehDitulis() {
		return lap, ErrPindahTidakLolos
	}
	for _, r := range sisip {
		if _, err := g.tulis(ctx, tx, TabelRingkasan, SqlSisipRingkasanFlat(n.tabelRingkasan), "memindah ringkasan", r.ID,
			nilaiArg(r.UsedBy), nilaiArg(r.Type), nilaiArg(r.ModifiedDate), nilaiArg(r.OperatorID)); err != nil {
			return lap, err
		}
	}
	for _, r := range ubah {
		if _, err := g.tulis(ctx, tx, TabelRingkasan, SqlPerbaruiRingkasanFlat(n.tabelRingkasan), "memperbarui ringkasan",
			nilaiArg(r.UsedBy), nilaiArg(r.Type), nilaiArg(r.ModifiedDate), nilaiArg(r.OperatorID), r.ID); err != nil {
			return lap, err
		}
	}
	c, err := g.satuNilai(ctx, tx, TabelRingkasan, SqlCacah(n.tabelRingkasan))
	if err != nil {
		return lap, err
	}
	if err := tx.Commit(); err != nil {
		return lap, fmt.Errorf("repository: menutup transaksi pindah: %w", err)
	}
	lap.CacahFlatSesudah, lap.Ditulis = angka(c), true
	return lap, nil
}
