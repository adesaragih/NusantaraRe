package repository

// Case NB baru dari tombol Create opportunity (tiket 29) - butir 76, keputusan work
// owner 03-10-2026. Tiga tulisan dalam SATU transaksi milik pemanggil (services):
//
//  1. nomor dari SEQ_WORK_POLIS_NB (migrasi 181) - "lanjut dari nomor terakhir Pega"
//     (butir 74.2): angka awalnya diisi work owner/DBA sebelum migrasi dijalankan;
//  2. baris T_WORK_POLIS - tabel yang ADA, milik premiumlistlife (050, 057/059/063),
//     dipakai bersama sesuai K-064; penandanya LINI = LiniFacIn (butir 76.2);
//  3. baris T_NB_OPPORTUNITY (migrasi 180, butir 76.3) - berbagi PK dengan T_WORK_POLIS,
//     tanpa constraint FK (pola T_PREMIUM_LIST premiumlistlife 050/051).
//
// Pembuat = PENGENAL AKUN pelaku di CREATE_OP dan CREATE_OP_NAME, bukan nama orang
// (ADR-U-0030; pola premiumlistlife polis_kasus.go). Seluruh nilai parameter terikat.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelWorkPolis - tabel kerja polis lintas lini (K-064), ejaan migrasi 050.
	TabelWorkPolis = "T_WORK_POLIS"
	// TabelOpportunity - isian opportunity, migrasi 180.
	TabelOpportunity = "T_NB_OPPORTUNITY"
	// SequenceCaseNB - penghitung nomor case NB, migrasi 181.
	SequenceCaseNB = "SEQ_WORK_POLIS_NB"
	// AwalanCaseNB - nomor case NB Pega berbentuk NB-<angka> tanpa tahun (tiket 29).
	AwalanCaseNB = "NB-"
	// LiniFacIn - penanda lini Fac In di T_WORK_POLIS.LINI (butir 76.2, preseden
	// T_WORK_CLAIM.LINI K-064).
	LiniFacIn = "FAC"
)

// Keadaan awal case NB yang baru dibuat - butir 77, keputusan work owner 03-10-2026
// (menggantikan A80 "kosong"), teks VERBATIM: "saat berhasil create position = Offer,
// FLAG_ONGOING_POLICY=0, STATUS_WORK=Pending-Policy." Kolomnya milik premiumlistlife:
// POSITION (050), FLAG_ONGOING_POLICY VARCHAR2(1) (057), STATUS_WORK (059/063).
const (
	PosisiAwalCaseNB = "Offer"
	FlagAwalCaseNB   = "0"
	StatusAwalCaseNB = "Pending-Policy"
)

// PenulisCaseNB - tiga langkah pembuatan case NB; transaksinya milik pemanggil.
type PenulisCaseNB interface {
	// PengenalBerikut - nomor case berikutnya, sudah berawalan (NB-<n>).
	PengenalBerikut(ctx context.Context, tx *db.Tx) (string, error)
	// SisipCase - baris T_WORK_POLIS berkeadaan awal butir 77 (Offer / 0 / Pending-Policy).
	SisipCase(ctx context.Context, tx *db.Tx, id, pembuat string) error
	// SisipOpportunity - baris T_NB_OPPORTUNITY ber-ID sama.
	SisipOpportunity(ctx context.Context, tx *db.Tx, id string, o models.Opportunity) error
}

// CaseNBOracle - PenulisCaseNB atas Oracle.
type CaseNBOracle struct{ db *db.DB }

// NewCaseNBOracle merakit penulis case NB.
func NewCaseNBOracle(d *db.DB) *CaseNBOracle { return &CaseNBOracle{db: d} }

// RakitPengenalCaseNB - awalan + angka sequence, TANPA padding (bentuk nomor NB Pega,
// sepola RakitPengenalWorkPolis premiumlistlife).
func RakitPengenalCaseNB(urut string) string { return AwalanCaseNB + strings.TrimSpace(urut) }

// kolomOpportunity - urutan kolom T_NB_OPPORTUNITY = urutan bind argSisipOpportunity.
var kolomOpportunity = []string{"ID", "ESTIMATED_CLOSING_DATE", "BUSINESS_PROSPECT_NAME", "ACCOUNT_ID", "INSURED_ID",
	"GROUP_BUSINESS_ID", "GROUP_BUSINESS", "CLASS_OF_BUSINESS", "TYPE_OF_INWARD", "TYPE_OF_FACULTATIVE", "PHASE", "STAGE",
	"OPPORTUNITY_SOURCE", "BUSINESS_STATUS", "DESCRIPTION"}

func sqlSisipCase(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, LINI, POSITION, FLAG_ONGOING_POLICY, STATUS_WORK, CREATE_OP, CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE)`+
		` VALUES (:1, :2, :3, :4, :5, :6, :7, SYSDATE, SYSDATE)`, tabel)
}

func argSisipCase(id, pembuat string) []any {
	return []any{id, LiniFacIn, PosisiAwalCaseNB, FlagAwalCaseNB, StatusAwalCaseNB, pembuat, pembuat}
}

func sqlSisipOpportunity(tabel string) string {
	bind := make([]string, len(kolomOpportunity))
	for i := range bind {
		bind[i] = fmt.Sprintf(":%d", i+1)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tabel, strings.Join(kolomOpportunity, ", "), strings.Join(bind, ", "))
}

func argSisipOpportunity(id string, o models.Opportunity) []any {
	k := db.KosongJadiNil
	return []any{id, o.EstimatedClosingDate, k(o.BusinessProspectName), k(o.AccountID), k(o.InsuredID),
		k(o.GroupBusinessID), k(o.GroupBusiness), k(o.ClassOfBusiness), k(o.TypeOfInward), k(o.TypeOfFacultative),
		k(o.Phase), k(o.Stage), k(o.OpportunitySource), k(o.BusinessStatus), k(o.Description)}
}

// PengenalBerikut - lihat PenulisCaseNB.
func (r *CaseNBOracle) PengenalBerikut(ctx context.Context, tx *db.Tx) (string, error) {
	urut, err := r.db.NomorBerikut(ctx, tx, SequenceCaseNB)
	if err != nil {
		return "", err
	}
	return RakitPengenalCaseNB(urut), nil
}

// SisipCase - lihat PenulisCaseNB.
func (r *CaseNBOracle) SisipCase(ctx context.Context, tx *db.Tx, id, pembuat string) error {
	return r.sisip(ctx, tx, TabelWorkPolis, sqlSisipCase, argSisipCase(id, pembuat))
}

// SisipOpportunity - lihat PenulisCaseNB.
func (r *CaseNBOracle) SisipOpportunity(ctx context.Context, tx *db.Tx, id string, o models.Opportunity) error {
	return r.sisip(ctx, tx, TabelOpportunity, sqlSisipOpportunity, argSisipOpportunity(id, o))
}

func (r *CaseNBOracle) sisip(ctx context.Context, tx *db.Tx, tabel string, sql func(string) string, arg []any) error {
	q, err := r.db.Qualify(tabel)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sql(q), "menyisipkan "+tabel, arg...)
	if err != nil {
		return err
	}
	return db.PastikanSatuBaris(hasil, tabel)
}

// jalankan - satu pernyataan tulis di transaksi pemanggil: PeriksaSQL, ExecContext,
// galat berlabel `aksi` (mis. "menyisipkan T_NB_OPPORTUNITY"). Dipakai casenb.go dan kasus.go.
func jalankan(ctx context.Context, tx *db.Tx, q, aksi string, arg ...any) (sql.Result, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	hasil, err := tx.ExecContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: %s: %w", aksi, err)
	}
	return hasil, nil
}
