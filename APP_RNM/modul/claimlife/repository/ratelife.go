package repository

// Rate retro per produk - OQ-M7 DITUTUP 29-09-2026 (GILIRAN-17), izin sempit
// seperti butir bh.
//
// Untuk apa berkas ini: membaca baris `RATE_LIFE` satu `OUTWARDRATEID` -
// masukan rate `services.HitungSpreading` (`BarisRate`).
//
// `[terverifikasi]` `Claim Life/RDBList/GetRateRetro.xml` b84:
//
//	SELECT AGE AS CARI1, CONTRACT AS CARI2, GENDER AS CARI3, RATE AS CARI4
//	  FROM POOLDATA.RATE_LIFE WHERE IDUSEDBY = {InputData.CARI3}
//
// dipanggil `SpreadingClaimLife_Act` langkah 6 (b1692, `RateRetroTemp` b1699,
// b1750) dengan `CARI3 = .OUTWARDRATEID` (langkah 5.1 b1502, b1527-b1528).
//
// `[keputusan work owner 29-09-2026 - lembar keputusan]` Izinnya TEPAT: satu
// pembaca, READ-ONLY, lima kolom bernama, berkunci `IDUSEDBY`. `ID` ikut
// dibaca supaya urutannya tetap (Pega tanpa ORDER BY) dan baris ganda dapat
// ditelusuri. `USEDBY` (nama pemakai) dan `TYPE` (NULL di setiap baris
// `[data DBA]`) TIDAK dibaca. Nol `JSONDATA`, nol tulisan.
//
// ⛔ `OUTWARDRATEID` tetap `[terbuka - DBA]`: di Pega ia hasil mengurai
// `M_PRODUCT_LIFE.JSONDATA` (`GetProductLife` b84-86, Java b1016), yang AC 38
// larang. Karena itu Spreading TETAP tidak dipanggil, dan pembaca ini nol
// pemanggil produksi sampai sumbernya ada (daftar serah terima).
//
// ⚠️ Nilainya TEKS apa adanya. `RATE` di view berdesimal koma pada sebagian
// besar baris, dan 564 nilai tidak bersih `[data DBA]`
// (KATALOG-TABEL-PESERTA-DAN-TREATY.md) - konversinya milik penyambung
// Spreading kelak, sekali saat masuk (ADR-U-0022), dengan galat per baris.
//
// Dibaca sesudah: ambangproduk.go (pola izin yang sama).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// KolomRateLife adalah LIMA kolom yang izin OQ-M7 cakup - dibaca penjaga
// `TestMasterViewTidakDisentuh`; kolom keenam akan berbunyi.
var KolomRateLife = []string{"ID", "AGE", "CONTRACT", "GENDER", "RATE"}

// NamaViewRateLife - view yang dibaca, ditulis UTUH (lihat ambangproduk.go).
const NamaViewRateLife = "RATE_LIFE"

// ErrRateLifeTanpaPengenal - `IDUSEDBY` kosong tidak pernah membaca view.
var ErrRateLifeTanpaPengenal = errors.New(
	"repository: pengenal rate (OUTWARDRATEID) kosong; RATE_LIFE tidak dibaca")

// BarisRateLife adalah satu baris `RATE_LIFE`, TEKS apa adanya.
type BarisRateLife struct {
	ID           string
	Umur         string // AGE
	Kontrak      string // CONTRACT
	JenisKelamin string // GENDER
	Rate         string // RATE - desimal koma sebagian besar baris
}

// RateLife membaca rate retro. READ-ONLY, seluruhnya.
type RateLife struct{ db *db.DB }

// NewRateLife menyusunnya.
func NewRateLife(db *db.DB) *RateLife { return &RateLife{db: db} }

// sqlRateLife merakit pembacaan kelima kolom.
func sqlRateLife(view string) string {
	return fmt.Sprintf(`SELECT r.%s, r.%s, r.%s, r.%s, r.%s FROM %s r
	  WHERE r.IDUSEDBY = :1 ORDER BY r.ID`,
		KolomRateLife[0], KolomRateLife[1], KolomRateLife[2], KolomRateLife[3], KolomRateLife[4], view)
}

// Baca membaca seluruh baris satu `IDUSEDBY`.
//
// ⛔ SATU-SATUNYA fungsi yang menyentuh view itu. Penjaga master mengunci
// namanya; pembaca kedua akan berbunyi.
func (r *RateLife) Baca(ctx context.Context, idUsedBy string) ([]BarisRateLife, error) {
	if strings.TrimSpace(idUsedBy) == "" {
		return nil, ErrRateLifeTanpaPengenal
	}
	view, err := r.db.Qualify(NamaViewRateLife)
	if err != nil {
		return nil, err
	}
	q := sqlRateLife(view)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, strings.TrimSpace(idUsedBy))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca view rate: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []BarisRateLife
	for rows.Next() {
		var id, umur, kontrak, jk, rate sql.NullString
		if err := rows.Scan(&id, &umur, &kontrak, &jk, &rate); err != nil {
			return nil, fmt.Errorf("repository: membaca satu baris rate: %w", err)
		}
		out = append(out, BarisRateLife{ID: id.String, Umur: umur.String, Kontrak: kontrak.String,
			JenisKelamin: jk.String, Rate: rate.String})
	}
	return out, rows.Err()
}
