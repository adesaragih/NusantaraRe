package repository

// Tab Object FIRE per halaman (tiket 75 F3, K75-3): case besar (hasil unggah CSV tanpa batas) tidak dibaca utuh. Halaman = lokasi
// urut SEQ_NO (lalu ID) dengan OFFSET / FETCH; pohon anak dibaca dengan kueri baca objek yang sama, dibatasi ke T_LOCATIONLIST.ID
// halaman itu (lingkupBaca). Pencarian `cari` `[keputusan agent]`: memuat (tanpa beda huruf besar-kecil) di OBJECT_NO /
// OBJECT_NAME / alamat risiko (ASM_ADDRESS) / kode pos (ASM_ZIP_CODE).

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/modul/nbfacin/backend/models"
)

// pembatasLokasi - pola tunggal kueri baca objek yang membatasi lokasi ke case.
const pembatasLokasi = "l.PARENT_ID = :1"

// lingkupBaca - lingkup bacaan pohon objek FIRE: semua lokasi case `id` (lokasi nil) atau hanya T_LOCATIONLIST.ID `lokasi`.
type lingkupBaca struct {
	id     string
	lokasi []string
}

// kueri - q apa adanya (semua lokasi), atau dengan `AND l.ID IN (:2, …)` sesudah satu-satunya pembatasLokasi. go-ora mengikat
// menurut URUTAN kemunculan: :1 tetap case, lalu ID lokasi.
func (b lingkupBaca) kueri(q string) (string, []any, error) {
	arg := []any{b.id}
	if b.lokasi == nil {
		return q, arg, nil
	}
	if n := strings.Count(q, pembatasLokasi); n != 1 {
		return "", nil, fmt.Errorf("repository: kueri objek memuat %d pola %q (harus 1)", n, pembatasLokasi)
	}
	if len(b.lokasi) == 0 {
		return "", nil, fmt.Errorf("repository: lingkup lokasi kosong")
	}
	bind := make([]string, len(b.lokasi))
	for i, x := range b.lokasi {
		bind[i] = ":" + strconv.Itoa(i+2)
		arg = append(arg, x)
	}
	return strings.Replace(q, pembatasLokasi, pembatasLokasi+" AND l.ID IN ("+strings.Join(bind, ", ")+")", 1), arg, nil
}

// BatasUkuranHalamanObjek - ukuran halaman tab Object terbesar (≤ 1000 ekspresi IN Oracle, ORA-01795).
const BatasUkuranHalamanObjek = 500

// PembacaHalamanObjek - tiket 75 F3.
type PembacaHalamanObjek interface {
	// JumlahLokasi - jumlah T_LOCATIONLIST case; ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
	JumlahLokasi(ctx context.Context, id string) (int, error)
	// BacaObjekHalaman - halaman ke-`halaman` (1..) berukuran `ukuran` (1..BatasUkuranHalamanObjek) lokasi yang cocok `cari`
	// (kosong = semua). ErrKasusTidakAda bila case tidak ada.
	BacaObjekHalaman(ctx context.Context, id string, halaman, ukuran int, cari string) (models.HalamanObjek, error)
}

// sqlDariLokasi - FROM + WHERE lokasi case, dengan saringan pencarian (4 bind :2..:5) bila cari.
func sqlDariLokasi(t tabelObjek, cari bool) string {
	q := " FROM " + t.loc + " l LEFT JOIN " + t.prop + " p ON p.PARENT_ID = l.ID LEFT JOIN " + t.risk +
		" r ON r.PARENT_ID = p.ID WHERE l.PARENT_ID = :1"
	if cari {
		q += ` AND (UPPER(p.OBJECT_NO) LIKE :2 ESCAPE '\' OR UPPER(p.OBJECT_NAME) LIKE :3 ESCAPE '\'` +
			` OR UPPER(r.ASM_ADDRESS) LIKE :4 ESCAPE '\' OR UPPER(r.ASM_ZIP_CODE) LIKE :5 ESCAPE '\')`
	}
	return q
}

func sqlJumlahLokasiCari(t tabelObjek, cari bool) string {
	return "SELECT COUNT(*)" + sqlDariLokasi(t, cari)
}

// sqlIDLokasiHalaman - ID lokasi satu halaman; OFFSET / FETCH = dua bind sesudah saringan.
func sqlIDLokasiHalaman(t tabelObjek, cari bool) string {
	n := 2
	if cari {
		n = 6
	}
	return "SELECT TO_CHAR(l.ID)" + sqlDariLokasi(t, cari) + " ORDER BY l.SEQ_NO, l.ID OFFSET :" + strconv.Itoa(n) +
		" ROWS FETCH NEXT :" + strconv.Itoa(n+1) + " ROWS ONLY"
}

// polaCari - nilai LIKE "memuat": huruf besar, `\ % _` di-escape.
func polaCari(cari string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(strings.ToUpper(strings.TrimSpace(cari))) + "%"
}

// argCari - bind :1 (case) + empat pola cari bila ada.
func argCari(id, cari string) []any {
	arg := []any{id}
	if strings.TrimSpace(cari) != "" {
		p := polaCari(cari)
		arg = append(arg, p, p, p, p)
	}
	return arg
}

// cekKasusObjek - ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
func (r *ObjekOracle) cekKasusObjek(ctx context.Context, t tabelObjek, id string) error {
	var n int
	if err := r.db.QueryRowContext(ctx, sqlAdaKasus(t.work), id, LiniFacIn).Scan(&n); err != nil {
		return fmt.Errorf("repository: cek case NB: %w", err)
	}
	if n == 0 {
		return ErrKasusTidakAda
	}
	return nil
}

// JumlahLokasi - lihat PembacaHalamanObjek.
func (r *ObjekOracle) JumlahLokasi(ctx context.Context, id string) (int, error) {
	t, err := r.tabel()
	if err != nil {
		return 0, err
	}
	if err := r.cekKasusObjek(ctx, t, id); err != nil {
		return 0, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, sqlJumlahLokasiCari(t, false), id).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: jumlah lokasi: %w", err)
	}
	return n, nil
}

// BacaObjekHalaman - lihat PembacaHalamanObjek.
func (r *ObjekOracle) BacaObjekHalaman(ctx context.Context, id string, halaman, ukuran int, cari string) (models.HalamanObjek, error) {
	h := models.HalamanObjek{Lokasi: []models.ObjekFire{}, IDLokasi: []string{}}
	if halaman < 1 || ukuran < 1 || ukuran > BatasUkuranHalamanObjek {
		return h, fmt.Errorf("repository: halaman %d / ukuran %d tidak sah", halaman, ukuran)
	}
	t, err := r.tabel()
	if err != nil {
		return h, err
	}
	if err := r.cekKasusObjek(ctx, t, id); err != nil {
		return h, err
	}
	ada := strings.TrimSpace(cari) != ""
	arg := argCari(id, cari)
	if err := r.db.QueryRowContext(ctx, sqlJumlahLokasiCari(t, ada), arg...).Scan(&h.Total); err != nil {
		return h, fmt.Errorf("repository: jumlah lokasi halaman: %w", err)
	}
	baris, err := r.db.QueryContext(ctx, sqlIDLokasiHalaman(t, ada), append(arg, (halaman-1)*ukuran, ukuran)...)
	if err != nil {
		return h, fmt.Errorf("repository: ID lokasi halaman: %w", err)
	}
	defer baris.Close()
	var idLokasi []string
	for baris.Next() {
		var x string
		if err := baris.Scan(&x); err != nil {
			return h, fmt.Errorf("repository: ID lokasi halaman: %w", err)
		}
		idLokasi = append(idLokasi, x)
	}
	if err := baris.Err(); err != nil {
		return h, fmt.Errorf("repository: ID lokasi halaman: %w", err)
	}
	baris.Close()
	if len(idLokasi) == 0 {
		return h, nil
	}
	lokasi, idBaca, err := r.bacaPohon(ctx, t, lingkupBaca{id: id, lokasi: idLokasi})
	if err != nil {
		return h, err
	}
	h.Lokasi, h.IDLokasi = lokasi, idBaca
	return h, nil
}
