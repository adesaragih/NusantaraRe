package repository

// Sequence nomor ORG (migrasi 810) TANPA Oracle: bentuk blok, pola nomor = format ID modul, jalur mundur, NEXTVAL.

import (
	"errors"
	"os"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/companydetail/backend/models"
)

// 810 membuat SEQ_CLIENT_ORG mulai dari nomor ORG tertinggi CLIENT.IDVIEW dan M_CLIENT.ID + 1 - dengan pola yang
// SAMA dengan format ID yang ditulis Create (models.AwalanID, models.AwalanIDView).
func TestMigrasi810SequenceMulaiSesudahNomorTertinggi(t *testing.T) {
	maju, err := migrasi.PernyataanLangkah(os.DirFS(".."), "810_seq_client_org.sql")
	if err != nil {
		t.Fatal(err)
	}
	if len(maju) != 1 {
		t.Fatalf("810: %d pernyataan, mau satu blok", len(maju))
	}
	s, ok := migrasi.BacaSequenceDariKueri(maju[0])
	if !ok {
		t.Fatalf("810 bukan blok sequence dari kueri:\n%s", maju[0])
	}
	if s.Nama != SequenceOrg {
		t.Errorf("sequence %s, mau %s", s.Nama, SequenceOrg)
	}
	q := strings.Join(strings.Fields(s.Kueri), " ")
	for _, k := range []string{
		"SELECT NVL(MAX(N), 0) + 1 FROM (",
		"FROM {skema}.CLIENT WHERE REGEXP_LIKE(IDVIEW, '^" + models.AwalanIDView + "[0-9]+$')",
		"UNION ALL",
		"FROM {skema}.M_CLIENT WHERE REGEXP_LIKE(ID, '^" + models.AwalanID + models.AwalanIDView + "[0-9]+$')",
		"REGEXP_LIKE(ID, '" + polaID + "')",
	} {
		if !strings.Contains(q, k) {
			t.Errorf("nilai awal tanpa %q:\n%s", k, q)
		}
	}
	if s.Opsi != "INCREMENT BY 1 NOCACHE NOCYCLE" {
		t.Errorf("opsi %q", s.Opsi)
	}
	mundur, err := migrasi.PernyataanLangkah(os.DirFS(".."), "810_seq_client_org_down.sql")
	if err != nil {
		t.Fatal(err)
	}
	if len(mundur) != 1 || mundur[0] != "DROP SEQUENCE {skema}."+SequenceOrg {
		t.Errorf("810 mundur: %q", mundur)
	}
}

// Nomor ORG dari NEXTVAL; sequence yang belum dibuat (810 belum dijalankan) = belum dimigrasi, bukan galat 500.
func TestNomorOrgDariSequence(t *testing.T) {
	if q := sqlNomorOrgBerikut("S.SEQ_CLIENT_ORG"); q != "SELECT S.SEQ_CLIENT_ORG.NEXTVAL FROM DUAL" {
		t.Errorf("nomor berikut: %s", q)
	}
	if err := bungkus(errors.New("ORA-02289: sequence does not exist"), "mengambil nomor ORG"); !errors.Is(err, ErrBelumDimigrasi) {
		t.Errorf("ORA-02289: %v", err)
	}
}
