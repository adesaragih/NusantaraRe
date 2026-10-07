package repository

// Tab Inw Fac Cedant Panels kasus FIRE (tiket 49): Share Cedant Type (T_GENERAL_POLIS.SHARE_CEDANT_TYPE), Share of
// Ceding (T_QUOTATIONDATA.SHARE_OF_CEDING, teks Pega), dan `.OfferFacIn.CedingCedantList` (T_CEDINGCEDANTLIST) - migrasi
// 199. Ceding Co tab General dibaca dari T_CEDINGCOLIST (tiket 34) sebagai sumber Add pertama (AddCedantList_act).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelCedingCedantList - grid cedant kasus (migrasi 199).
	TabelCedingCedantList    = "T_CEDINGCEDANTLIST"
	sequenceCedingCedantList = "SEQ_T_CEDINGCEDANTLIST"
	// lebarNamaCedant - T_CEDINGCEDANTLIST.CEDING_CO_NAME (migrasi 199).
	lebarNamaCedant = 500
)

// PenyimpanCedant - data tab Cedant.
type PenyimpanCedant interface {
	// BacaCedant - ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
	BacaCedant(ctx context.Context, id string) (models.KasusCedant, error)
	// TulisCedant - Share Cedant Type, Share of Ceding (bila tidak nil), dan grid cedant diganti utuh di transaksi
	// pemanggil; nama ceding dari AGENT menurut kode. ErrKasusTidakAda bila case tidak ada, ErrCedingTidakSah /
	// ErrCedingTerlaluPanjang bila kode tidak lolos AGENT / nama melebihi lebar kolom.
	TulisCedant(ctx context.Context, tx *db.Tx, id string, k models.SimpanCedant) error
}

// CedantOracle - PenyimpanCedant atas Oracle.
type CedantOracle struct{ db *db.DB }

// NewCedantOracle merakit penyimpan tab Cedant.
func NewCedantOracle(d *db.DB) *CedantOracle { return &CedantOracle{db: d} }

type tabelCedant struct{ work, general, quo, ceding, cedant string }

func (r *CedantOracle) tabel() (tabelCedant, error) {
	var t tabelCedant
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelGeneralPolis, &t.general}, {TabelQuotationData, &t.quo},
		{TabelCedingCoList, &t.ceding}, {TabelCedingCedantList, &t.cedant}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

// sqlBacaKepalaCedant - satu baris per case Fac In (General / Quotation boleh belum ada).
func sqlBacaKepalaCedant(t tabelCedant) string {
	return "SELECT TO_CHAR(g.SHARE_CEDANT_TYPE), q.SOB_NAME, q.CEDING_CO, q.CEDING_CO_NAME, w.FLAG_ONGOING_POLICY FROM " + t.work + " w LEFT JOIN " +
		t.general + " g ON g.ID = w.ID LEFT JOIN " + t.quo + " q ON q.PARENT_ID = w.ID WHERE w.ID = :1 AND w.LINI = :2"
}

func sqlBacaCedant(cedant string) string {
	return "SELECT CEDING_CO, CEDING_CO_NAME, " + angkaKeluar("SHARE_CEDING") + " FROM " + cedant +
		" WHERE PARENT_ID = :1 ORDER BY SEQ_NO"
}

func sqlUbahShareCedantType(general string) string {
	return "UPDATE " + general + " SET SHARE_CEDANT_TYPE = :1 WHERE ID = :2"
}

func sqlUbahShareOfCeding(quo string) string {
	return "UPDATE " + quo + " SET SHARE_OF_CEDING = :1 WHERE PARENT_ID = :2"
}

func sqlSisipQuotationShareOfCeding(quo string) string {
	return "INSERT INTO " + quo + " (ID, PARENT_ID, SHARE_OF_CEDING) VALUES (:1, :2, :3)"
}

func sqlHapusCedant(cedant string) string { return "DELETE FROM " + cedant + " WHERE PARENT_ID = :1" }

func sqlSisipCedant(cedant string) string {
	return "INSERT INTO " + cedant + " (ID, PARENT_ID, SEQ_NO, ROW_UID, CEDING_CO, CEDING_CO_NAME, SHARE_CEDING) VALUES (:1, :2," +
		" :3, :4, :5, :6, " + angkaMasuk(":7") + ")"
}

// BacaCedant - lihat PenyimpanCedant.
func (r *CedantOracle) BacaCedant(ctx context.Context, id string) (models.KasusCedant, error) {
	var k models.KasusCedant
	t, err := r.tabel()
	if err != nil {
		return k, err
	}
	var tipe, sob, kode, nama, flag sql.NullString
	err = r.db.QueryRowContext(ctx, sqlBacaKepalaCedant(t), id, LiniFacIn).Scan(&tipe, &sob, &kode, &nama, &flag)
	if errors.Is(err, sql.ErrNoRows) {
		return k, ErrKasusTidakAda
	}
	if err != nil {
		return k, fmt.Errorf("repository: baca kepala tab Cedant: %w", err)
	}
	k.ShareCedantType, k.SobName, k.FlagOnGoingPolicy = tipe.String, sob.String, flag.String
	if k.Cedant, err = r.bacaCedant(ctx, t.cedant, id); err != nil {
		return k, err
	}
	// AddCedantList_act: CedingCoList berisi -> seluruhnya; kosong -> satu baris CedingCo / CedingCoName.
	baris, err := r.db.QueryContext(ctx, sqlBacaCeding(t.ceding, t.quo), id)
	if err != nil {
		return k, fmt.Errorf("repository: baca %s: %w", TabelCedingCoList, err)
	}
	defer baris.Close()
	k.CedingUmum = []models.BarisCedant{}
	for baris.Next() {
		var c, n sql.NullString
		if err := baris.Scan(&c, &n); err != nil {
			return k, fmt.Errorf("repository: %s: %w", TabelCedingCoList, err)
		}
		k.CedingUmum = append(k.CedingUmum, models.BarisCedant{CedingCo: c.String, CedingCoName: n.String})
	}
	if err := baris.Err(); err != nil {
		return k, fmt.Errorf("repository: %s: %w", TabelCedingCoList, err)
	}
	if len(k.CedingUmum) == 0 && (kode.String != "" || nama.String != "") {
		k.CedingUmum = pecahGabunganCeding(kode.String, nama.String)
	}
	return k, nil
}

// pecahGabunganCeding - T_QUOTATIONDATA.CEDING_CO / CEDING_CO_NAME menyimpan gabungan `;` (tiket 34, gabungCeding), bukan
// `QuotationData.CedingCo` tunggal Pega: dipecah per kode. Kode dan nama berasal dari baris dan urutan gabung yang sama;
// bila jumlah potongannya berbeda (nama ber-`;`), nama dikosongkan - Choose Ceding / Save mengisinya dari AGENT.
func pecahGabunganCeding(kode, nama string) []models.BarisCedant {
	k, n := strings.Split(kode, pemisahCeding), strings.Split(nama, pemisahCeding)
	hasil := make([]models.BarisCedant, 0, len(k))
	for i, c := range k {
		b := models.BarisCedant{CedingCo: c}
		if len(n) == len(k) {
			b.CedingCoName = n[i]
		}
		hasil = append(hasil, b)
	}
	return hasil
}

func (r *CedantOracle) bacaCedant(ctx context.Context, cedant, id string) ([]models.BarisCedant, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaCedant(cedant), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelCedingCedantList, err)
	}
	defer baris.Close()
	hasil := []models.BarisCedant{}
	for baris.Next() {
		var c, n, s sql.NullString
		if err := baris.Scan(&c, &n, &s); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelCedingCedantList, err)
		}
		share, err := bacaDesimal(id, TabelCedingCedantList+".SHARE_CEDING", &s)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, models.BarisCedant{CedingCo: c.String, CedingCoName: n.String, ShareCeding: share})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelCedingCedantList, err)
	}
	return hasil, nil
}

// TulisCedant - lihat PenyimpanCedant. Urutan: sentuh T_WORK_POLIS (kunci + case ada), pastikan General, Share Cedant
// Type, Share of Ceding (UPDATE-atau-INSERT baris Quotation), lalu grid diganti utuh.
func (r *CedantOracle) TulisCedant(ctx context.Context, tx *db.Tx, id string, k models.SimpanCedant) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return err
	}
	var tipe any // NULL = belum dipilih
	if k.ShareCedantType != "" {
		n, err := strconv.Atoi(k.ShareCedantType)
		if err != nil {
			return fmt.Errorf("repository: SHARE_CEDANT_TYPE %q bukan angka", k.ShareCedantType)
		}
		tipe = n
	}
	if _, err := jalankan(ctx, tx, sqlUbahShareCedantType(t.general), "mengubah Share Cedant Type", tipe, id); err != nil {
		return err
	}
	if k.ShareOfCeding != nil {
		hasil, err := jalankan(ctx, tx, sqlUbahShareOfCeding(t.quo), "mengubah Share of Ceding", db.KosongJadiNil(*k.ShareOfCeding), id)
		if err != nil {
			return err
		}
		n, err := hasil.RowsAffected()
		if err != nil {
			return err
		}
		if n == 0 {
			urut, err := r.db.NomorBerikut(ctx, tx, SequenceQuotationData)
			if err != nil {
				return err
			}
			if _, err := jalankan(ctx, tx, sqlSisipQuotationShareOfCeding(t.quo), "menyisip "+TabelQuotationData, urut, id,
				db.KosongJadiNil(*k.ShareOfCeding)); err != nil {
				return err
			}
		}
	}
	if _, err := jalankan(ctx, tx, sqlHapusCedant(t.cedant), "menghapus "+TabelCedingCedantList, id); err != nil {
		return err
	}
	for i, c := range k.Cedant {
		// Choose Ceding (SetDataSobCeding_Act cabang CedingCedant): kode = AGENT.ID, nama = AGENT.CLIENTNAME - nama dari
		// AGENT menurut kode (pola tiket 33/34, E-4), bukan dari klien; kode kosong = baris Add belum dipilih (C-4).
		nama := ""
		if c.CedingCo != "" {
			if nama, err = namaAgent(ctx, r.db, tx, c.CedingCo); errors.Is(err, ErrSOBTidakSah) {
				return fmt.Errorf("%w: cedant[%d]", ErrCedingTidakSah, i)
			} else if err != nil {
				return err
			}
			if len(nama) > lebarNamaCedant {
				return fmt.Errorf("%w: nama ceding cedant[%d] melebihi %d bita", ErrCedingTerlaluPanjang, i, lebarNamaCedant)
			}
		}
		urut, err := r.db.NomorBerikut(ctx, tx, sequenceCedingCedantList)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		if _, err := jalankan(ctx, tx, sqlSisipCedant(t.cedant), "menyisip "+TabelCedingCedantList, urut, id, i+1, uid,
			db.KosongJadiNil(c.CedingCo), db.KosongJadiNil(nama), ikatDesimal(c.ShareCeding)); err != nil {
			return err
		}
	}
	return nil
}
