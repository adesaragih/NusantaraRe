package repository

// Baca/simpan case NB untuk layar Inward Facultative (tiket 31) - butir 78. Case =
// baris T_WORK_POLIS ber-LINI LiniFacIn; blok General ke T_GENERAL_POLIS (182) dan
// T_QUOTATIONDATA (183) - tabel flat rancangan, dibuat sebagian. Seluruh nilai
// parameter terikat; tidak ada `SELECT *`.

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
	// TabelGeneralPolis, TabelQuotationData - tabel flat rancangan (migrasi 182, 183).
	TabelGeneralPolis  = "T_GENERAL_POLIS"
	TabelQuotationData = "T_QUOTATIONDATA"
	// SequenceQuotationData - pembangkit T_QUOTATIONDATA.ID (migrasi 183).
	SequenceQuotationData = "SEQ_T_QUOTATIONDATA"
)

// ErrKasusTidakAda - tidak ada baris T_WORK_POLIS ber-ID itu dengan LINI Fac In.
var ErrKasusTidakAda = errors.New("repository: case NB tidak ditemukan")

// PembacaKasus - baca satu case NB.
type PembacaKasus interface {
	BacaKasus(ctx context.Context, id string) (models.Kasus, error)
}

// PenulisGeneral - simpan blok General di dalam transaksi pemanggil. Nilai kosong = NULL.
type PenulisGeneral interface {
	SimpanGeneral(ctx context.Context, tx *db.Tx, id string, g models.General) error
}

// PenyimpanKasus - PembacaKasus + PenulisGeneral.
type PenyimpanKasus interface {
	PembacaKasus
	PenulisGeneral
}

// KasusOracle - PenyimpanKasus atas Oracle.
type KasusOracle struct{ db *db.DB }

// NewKasusOracle merakit pembaca/penulis case NB.
func NewKasusOracle(d *db.DB) *KasusOracle { return &KasusOracle{db: d} }

// kolomBacaKasus - kolom sqlBacaKasus, BERNAMA (dibaca lewat kunci, bukan urutan).
// "INSUREDNAME" = subkueri MAX: T_M_ACCOUNT tidak ber-PK di DDL-nya (A88).
var kolomBacaKasus = []string{"w.ID", "w.POSITION", "w.STATUS_WORK",
	"o.ESTIMATED_CLOSING_DATE", "o.BUSINESS_PROSPECT_NAME", "o.ACCOUNT_ID", "o.INSURED_ID", "o.GROUP_BUSINESS_ID",
	"o.GROUP_BUSINESS", "o.CLASS_OF_BUSINESS", "o.TYPE_OF_INWARD", "o.TYPE_OF_FACULTATIVE", "o.PHASE", "o.STAGE",
	"o.OPPORTUNITY_SOURCE", "o.BUSINESS_STATUS", "o.DESCRIPTION", "INSUREDNAME",
	"g.ID", "g.START_DATE_TIME", "g.OFFERING_DATE", "g.END_DATE_TIME", "g.FOLLOWING",
	"q.NO_OFFER_SLIP", "q.QQ_NAME", "q.POLICY_TYPE", "q.MOID", "q.EDM_DAY", "q.TYPE_FACULTATIVE", "q.SOURCE_OF_BUSINESS", "q.SOB_NAME",
	"q.CEDING_CO_NAME", "q.GROUP_NAME"}

// kolomTanggalKasus - satu-satunya kolom DATE di kolomBacaKasus.
const kolomTanggalKasus = "o.ESTIMATED_CLOSING_DATE"

// sqlBacaKasus - satu baris: case + opportunity + nama tertanggung + General, dibatasi LINI.
func sqlBacaKasus(work, opp, akun, general, quo string) string {
	pilih := make([]string, len(kolomBacaKasus))
	for i, k := range kolomBacaKasus {
		pilih[i] = k
		if k == "INSUREDNAME" {
			pilih[i] = "(SELECT MAX(a.INSUREDNAME) FROM " + akun + " a WHERE a.ID = o.ACCOUNT_ID)"
		}
	}
	return fmt.Sprintf(`SELECT %s
FROM %s w
LEFT JOIN %s o ON o.ID = w.ID
LEFT JOIN %s g ON g.ID = w.ID
LEFT JOIN %s q ON q.PARENT_ID = w.ID
WHERE w.ID = :1 AND w.LINI = :2`, strings.Join(pilih, ", "), work, opp, general, quo)
}

// BacaKasus - lihat PembacaKasus. ErrKasusTidakAda bila tidak ada.
func (r *KasusOracle) BacaKasus(ctx context.Context, id string) (models.Kasus, error) {
	var t [5]string
	for i, n := range []string{TabelWorkPolis, TabelOpportunity, TabelAkun, TabelGeneralPolis, TabelQuotationData} {
		q, err := r.db.Qualify(n)
		if err != nil {
			return models.Kasus{}, err
		}
		t[i] = q
	}
	teks := map[string]*sql.NullString{}
	var tutup sql.NullTime
	tujuan := make([]any, len(kolomBacaKasus))
	for i, k := range kolomBacaKasus {
		if k == kolomTanggalKasus {
			tujuan[i] = &tutup
			continue
		}
		teks[k] = &sql.NullString{}
		tujuan[i] = teks[k]
	}
	err := r.db.QueryRowContext(ctx, sqlBacaKasus(t[0], t[1], t[2], t[3], t[4]), id, LiniFacIn).Scan(tujuan...)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Kasus{}, ErrKasusTidakAda
	}
	if err != nil {
		return models.Kasus{}, fmt.Errorf("repository: baca case NB: %w", err)
	}
	v := func(k string) string { return teks[k].String }
	k := models.Kasus{CaseID: v("w.ID"), Position: v("w.POSITION"), StatusWork: v("w.STATUS_WORK"),
		InsuredName: v("INSUREDNAME"),
		Opportunity: models.Opportunity{BusinessProspectName: v("o.BUSINESS_PROSPECT_NAME"), AccountID: v("o.ACCOUNT_ID"),
			InsuredID: v("o.INSURED_ID"), GroupBusinessID: v("o.GROUP_BUSINESS_ID"), GroupBusiness: v("o.GROUP_BUSINESS"),
			ClassOfBusiness: v("o.CLASS_OF_BUSINESS"), TypeOfInward: v("o.TYPE_OF_INWARD"),
			TypeOfFacultative: v("o.TYPE_OF_FACULTATIVE"), Phase: v("o.PHASE"), Stage: v("o.STAGE"),
			OpportunitySource: v("o.OPPORTUNITY_SOURCE"), BusinessStatus: v("o.BUSINESS_STATUS"),
			Description: v("o.DESCRIPTION")},
		General: models.General{Tersimpan: teks["g.ID"].Valid, StartDateTime: v("g.START_DATE_TIME"),
			OfferingDate: v("g.OFFERING_DATE"), EndDateTime: v("g.END_DATE_TIME"), OldPolicyNumber: v("g.FOLLOWING"),
			ReffNumber: v("q.NO_OFFER_SLIP"), QQName: v("q.QQ_NAME"), PolicyType: v("q.POLICY_TYPE"),
			MarketingID: v("q.MOID"), Day: v("q.EDM_DAY"), TypeFacultative: v("q.TYPE_FACULTATIVE"),
			SourceOfBusinessID: v("q.SOURCE_OF_BUSINESS"), SourceOfBusiness: v("q.SOB_NAME"), CedingCoName: v("q.CEDING_CO_NAME"), GroupName: v("q.GROUP_NAME")}}
	if tutup.Valid {
		k.Opportunity.EstimatedClosingDate = tutup.Time
	}
	ceding, err := r.db.Qualify(TabelCedingCoList)
	if err != nil {
		return models.Kasus{}, err
	}
	if k.General.CedingList, err = r.bacaCeding(ctx, ceding, t[4], id); err != nil {
		return models.Kasus{}, err
	}
	return k, nil
}

func sqlSentuhCase(work string) string {
	return "UPDATE " + work + " SET TGL_UPDATE = SYSDATE WHERE ID = :1 AND LINI = :2"
}

func sqlUbahGeneral(general string) string {
	return "UPDATE " + general + " SET START_DATE_TIME = :1, OFFERING_DATE = :2, END_DATE_TIME = :3 WHERE ID = :4"
}

func sqlSisipGeneral(general string) string {
	return "INSERT INTO " + general + " (ID, START_DATE_TIME, OFFERING_DATE, END_DATE_TIME) VALUES (:1, :2, :3, :4)"
}

func sqlUbahQuotation(quo string) string {
	return "UPDATE " + quo + " SET NO_OFFER_SLIP = :1, QQ_NAME = :2, POLICY_TYPE = :3, MOID = :4, EDM_DAY = :5," +
		" TYPE_FACULTATIVE = :6, SOURCE_OF_BUSINESS = :7, SOB_NAME = :8 WHERE PARENT_ID = :9"
}

func sqlSisipQuotation(quo string) string {
	return "INSERT INTO " + quo + " (ID, PARENT_ID, NO_OFFER_SLIP, QQ_NAME, POLICY_TYPE, MOID, EDM_DAY, TYPE_FACULTATIVE," +
		" SOURCE_OF_BUSINESS, SOB_NAME) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10)"
}

// SimpanGeneral - lihat PenulisGeneral. Urutan: sentuh T_WORK_POLIS (TGL_UPDATE; sekaligus
// mengunci baris dan memastikan case Fac In ada), lalu UPDATE-atau-INSERT T_GENERAL_POLIS
// dan T_QUOTATIONDATA, lalu daftar Ceding Co (tiket 34). SOB dan ceding ditulis dengan nama dari
// AGENT (tiket 33/34); Group Name dan Following tampil saja, tidak ditulis.
func (r *KasusOracle) SimpanGeneral(ctx context.Context, tx *db.Tx, id string, g models.General) error {
	work, err := r.db.Qualify(TabelWorkPolis)
	if err != nil {
		return err
	}
	general, err := r.db.Qualify(TabelGeneralPolis)
	if err != nil {
		return err
	}
	quo, err := r.db.Qualify(TabelQuotationData)
	if err != nil {
		return err
	}
	n, err := r.ubah(ctx, tx, sqlSentuhCase(work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrKasusTidakAda
	}
	k := db.KosongJadiNil
	tgl := []any{k(g.StartDateTime), k(g.OfferingDate), k(g.EndDateTime)}
	if n, err = r.ubah(ctx, tx, sqlUbahGeneral(general), "mengubah "+TabelGeneralPolis, append(tgl, id)...); err != nil {
		return err
	}
	if n == 0 {
		if err := r.sisip(ctx, tx, sqlSisipGeneral(general), TabelGeneralPolis, append([]any{id}, tgl...)...); err != nil {
			return err
		}
	}
	// Tiket 33 (E-4): nama SOB dari AGENT menurut kode, bukan dari klien; kode kosong =
	// kode dan nama dikosongkan.
	namaSob := ""
	if g.SourceOfBusinessID != "" {
		if namaSob, err = namaAgent(ctx, r.db, tx, g.SourceOfBusinessID); err != nil {
			return err
		}
	}
	isi := []any{k(g.ReffNumber), k(g.QQName), k(g.PolicyType), k(g.MarketingID), k(g.Day), k(g.TypeFacultative),
		k(g.SourceOfBusinessID), k(namaSob)}
	if n, err = r.ubah(ctx, tx, sqlUbahQuotation(quo), "mengubah "+TabelQuotationData, append(isi, id)...); err != nil {
		return err
	}
	if n == 0 {
		urut, err := r.db.NomorBerikut(ctx, tx, SequenceQuotationData)
		if err != nil {
			return err
		}
		if err := r.sisip(ctx, tx, sqlSisipQuotation(quo), TabelQuotationData, append([]any{urut, id}, isi...)...); err != nil {
			return err
		}
	}
	// Tiket 34: daftar Ceding Co diganti utuh (kosong = dikosongkan).
	return r.simpanCeding(ctx, tx, quo, id, g.CedingIDs)
}

func (r *KasusOracle) ubah(ctx context.Context, tx *db.Tx, q, aksi string, arg ...any) (int64, error) {
	hasil, err := jalankan(ctx, tx, q, aksi, arg...)
	if err != nil {
		return 0, err
	}
	return hasil.RowsAffected()
}

func (r *KasusOracle) sisip(ctx context.Context, tx *db.Tx, q, tabel string, arg ...any) error {
	hasil, err := jalankan(ctx, tx, q, "menyisipkan "+tabel, arg...)
	if err != nil {
		return err
	}
	return db.PastikanSatuBaris(hasil, tabel)
}
