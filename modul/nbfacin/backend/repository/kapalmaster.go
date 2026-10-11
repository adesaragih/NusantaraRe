package repository

// Master kapal untuk unggah CSV Marine Cargo (tiket 74 P6, keputusan work owner 11-10-2026 K74-16). Tabel POOLDATA.SHIP (`[dugaan
// kuat]` kelas Int-SHIP, lihat lookupaneka.go), skema lewat db.Qualify seperti CariKapal.
//
// `[terverifikasi]` ProtectUploadMarineCargo_act (NB FacIn):
//   - langkah 6 BrowseShip_RD: `.ID = Param.ID` (Param.ID = @toUpperCase(DATASHIP)), filter lain kosong, G/H `.Remark != "Not Active"
//     OR .Remark IS NULL`; ada hasil -> "Jangan gunakan ID kapal dikolom DATASHIP pada CSV" (AdaKapalID);
//   - langkah 7 BrowseShip2_RD: `.Name = Param.Name` (pyCaseInsensitive, Param.Name = @toUpperCase(DATASHIP)), filter G/H yang sama,
//     urut `.ID` DESC, pxResults(1) dipakai InsertUploadMarineCargo_act 3.10 (KapalDariNama);
//   - langkah 7.5 tidak ada -> RDB-List UpdateMasterShip = `POOLDATA.RDBMASTERSHIP(InputShip.*, OperatorID.pyUserIdentifier, …);
//     COMMIT;` dengan InputShip.Name = UPPER(DATASHIP), Remark "Incomplete", Type "Marine Cargo", parameter lain kosong.
// RDBMASTERSHIP (`D:\migrasi\RNM\DDL\RDBMASTERSHIP.txt`) DIPORT tanpa CALL (ADR-0043): SELECT duplikatnya membandingkan NAME DAN
// REGISTER / FLAGNAME / YEARMAKE1 / SHIPNAME dengan `=` - parameter kosong = NULL dan `= NULL` tidak pernah benar, jadi untuk
// pemanggil ini prosedur SELALU menyisip `[terverifikasi]`: ID = '1' || LPAD(TO_CHAR(ship_seq.NEXTVAL), 9, '0') (LPAD memotong dari
// kiri bila lebih dari 9), NAME, REMARK, TYPE, OPERATOR_ID, INPUT_TIME = SYSDATE; kolom lain NULL. Penulisan POOLDATA.SHIP disetujui
// work owner 11-10-2026 (K74-16). Transaksi milik pemanggil (Pega men-COMMIT per panggilan - services memakai transaksi sendiri).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// sequenceShip - `ship_seq` prosedur POOLDATA.RDBMASTERSHIP (tanpa awalan skema = skema pemilik prosedur).
	sequenceShipMaster = "SHIP_SEQ"
	// RemarkKapalBaru / TipeKapalBaru - nilai InputShip langkah 7.5.1 ProtectUploadMarineCargo_act.
	RemarkKapalBaru = "Incomplete"
	TipeKapalBaru   = "Marine Cargo"
)

// PenyimpanKapalMaster - kapal POOLDATA.SHIP untuk unggah CSV Marine Cargo.
type PenyimpanKapalMaster interface {
	// AdaKapalID - BrowseShip_RD: kapal aktif ber-ID `id` ada.
	AdaKapalID(ctx context.Context, id string) (bool, error)
	// KapalDariNama - BrowseShip2_RD pxResults(1): kapal aktif bernama `nama` (tak peka huruf), ID terbesar.
	KapalDariNama(ctx context.Context, nama string) (models.KapalCargo, bool, error)
	// DaftarKapal - port RDBMASTERSHIP cabang INSERT (lihat kepala berkas) di transaksi tx; jawab kapal baru.
	DaftarKapal(ctx context.Context, tx *db.Tx, nama, operator string) (models.KapalCargo, error)
}

// KapalMasterOracle - PenyimpanKapalMaster atas Oracle.
type KapalMasterOracle struct{ db *db.DB }

// NewKapalMasterOracle merakit penyimpan master kapal.
func NewKapalMasterOracle(d *db.DB) *KapalMasterOracle { return &KapalMasterOracle{db: d} }

// sqlKapalAktifID - :1 ID, :2 literal Remark tidak aktif.
func sqlKapalAktifID(tabel string) string {
	return "SELECT COUNT(*) FROM " + tabel + " WHERE ID = :1 AND (REMARK != :2 OR REMARK IS NULL) AND ROWNUM = 1"
}

// sqlKapalNama - :1 nama (huruf besar), :2 literal Remark tidak aktif; kolom yang diset InsertUploadMarineCargo_act 3.10.
func sqlKapalNama(tabel string) string {
	return "SELECT ID, NAME, SHIPTYPEID, SHIPNAME, FLAGNAME, CONSTRUCTION, GRT, NRT, DWT, YEARMAKE1, YEARMAKE2, REMARK, REGISTER FROM " +
		tabel + " WHERE UPPER(NAME) = UPPER(:1) AND (REMARK != :2 OR REMARK IS NULL) ORDER BY ID DESC FETCH FIRST 1 ROWS ONLY"
}

// sqlDaftarKapal - INSERT RDBMASTERSHIP: :1 ID, :2 NAME, :3 REMARK, :4 TYPE, :5 OPERATOR_ID.
func sqlDaftarKapal(tabel string) string {
	return "INSERT INTO " + tabel + " (ID, NAME, REMARK, TYPE, OPERATOR_ID, INPUT_TIME) VALUES (:1, :2, :3, :4, :5, SYSDATE)"
}

// idKapalBaru - '1' || LPAD(n, 9, '0'): LPAD Oracle memotong ke 9 karakter pertama bila lebih panjang.
func idKapalBaru(n string) string {
	if len(n) >= 9 {
		return "1" + n[:9]
	}
	return "1" + strings.Repeat("0", 9-len(n)) + n
}

// AdaKapalID - lihat PenyimpanKapalMaster.
func (r *KapalMasterOracle) AdaKapalID(ctx context.Context, id string) (bool, error) {
	q, err := r.db.Qualify(TabelShip)
	if err != nil {
		return false, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, sqlKapalAktifID(q), id, remarkTidakAktif).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: baca %s: %w", TabelShip, err)
	}
	return n > 0, nil
}

// KapalDariNama - lihat PenyimpanKapalMaster. GRT / NRT / DWT: "," -> "." (3.10 @toDecimal(@replaceAll(…, ',', '.'))).
func (r *KapalMasterOracle) KapalDariNama(ctx context.Context, nama string) (models.KapalCargo, bool, error) {
	q, err := r.db.Qualify(TabelShip)
	if err != nil {
		return models.KapalCargo{}, false, err
	}
	var v [13]sql.NullString
	tuju := make([]any, len(v))
	for i := range v {
		tuju[i] = &v[i]
	}
	err = r.db.QueryRowContext(ctx, sqlKapalNama(q), nama, remarkTidakAktif).Scan(tuju...)
	if errors.Is(err, sql.ErrNoRows) {
		return models.KapalCargo{}, false, nil
	}
	if err != nil {
		return models.KapalCargo{}, false, fmt.Errorf("repository: baca %s: %w", TabelShip, err)
	}
	koma := func(s sql.NullString) string { return strings.ReplaceAll(s.String, ",", ".") }
	return models.KapalCargo{ID: v[0].String, NmShip: v[1].String, ShipTypeID: v[2].String, ShipTypeName: v[3].String,
		FlagName: v[4].String, Construction: v[5].String, GRT: koma(v[6]), NRT: koma(v[7]), DWT: koma(v[8]), YMake1: v[9].String,
		YMake2: v[10].String, Remark: v[11].String, Register: v[12].String}, true, nil
}

// DaftarKapal - lihat PenyimpanKapalMaster.
func (r *KapalMasterOracle) DaftarKapal(ctx context.Context, tx *db.Tx, nama, operator string) (models.KapalCargo, error) {
	q, err := r.db.Qualify(TabelShip)
	if err != nil {
		return models.KapalCargo{}, err
	}
	n, err := r.db.NomorBerikut(ctx, tx, sequenceShipMaster)
	if err != nil {
		return models.KapalCargo{}, err
	}
	id := idKapalBaru(n)
	if _, err := jalankan(ctx, tx, sqlDaftarKapal(q), "menyisipkan "+TabelShip, id, nama, RemarkKapalBaru, TipeKapalBaru,
		db.KosongJadiNil(operator)); err != nil {
		return models.KapalCargo{}, err
	}
	return models.KapalCargo{ID: id, NmShip: nama, Remark: RemarkKapalBaru}, nil
}
