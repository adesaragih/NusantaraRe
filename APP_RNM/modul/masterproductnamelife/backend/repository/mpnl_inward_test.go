package repository

// Penulis sisi inward (paket 4) - perakit JSONDATA inward dan SQL-nya, tanpa Oracle.

import (
	"strings"
	"testing"

	go_ora "github.com/sijms/go-ora/v2"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

func inwardUji() models.Produk {
	p := produkUji()
	p.Inward = models.ProdukInward{ID: "100044", ProductID: "100044", PolicyHolder: "UJI-ORG-1",
		PolicyHolderName: "UJI PEMEGANG", Insured: "UJI TERTANGGUNG", Begin: "2026-03-01", Mature: "2027-02-28",
		STNC: "2026-03-26", CedingLimit: "150000000.123456789", MinAge: "22", MaxAge: "70", MaxExpiredClaim: "180",
		MaxDataReceive: "90", Currency: "IDR", CurrencyID: "1", ExpiryAge: "75", Payment: "1",
		SubjectTo: "baris 1\nbaris 2 <b>"}
	return p
}

func TestRakitInwardKunciPegaDanKunciView(t *testing.T) {
	teks, err := RakitInward(inwardUji(), "")
	if err != nil {
		t.Fatal(err)
	}
	obj, err := uraiObjek(teks)
	if err != nil {
		t.Fatalf("JSON inward tidak sah: %v\n%s", err, teks)
	}
	for _, k := range KunciViewInward {
		if _, ada := obj[k]; !ada {
			t.Errorf("kunci view PRODUCTINWARD_LIFE %q tidak ada di JSON hasil simpan", k)
		}
	}
	for _, w := range []string{`"ID":"100044"`, `"PRODUCTID":"100044"`, `"POLICYHODER":"UJI-ORG-1"`, `"BEGIN":"01/03/2026"`,
		`"MATURE":"28/02/2027"`, `"STNC":"26/03/2026"`, `"CEDINGLIMIT":"150000000.123456789"`, `"EXPIRYAGE":"75"`,
		`"CURRENCYID":"1"`, `"SUBJECTTO":"baris 1\nbaris 2 <b>"`} {
		if !strings.Contains(teks, w) {
			t.Errorf("JSON inward tanpa %s:\n%s", w, teks)
		}
	}
}

func TestRakitInwardMempertahankanKunciLama(t *testing.T) {
	teks, err := RakitInward(inwardUji(), `{"ID":"100009","PRODUCTID":"100044","KUNCILAMA":[1,2],"MAXAGE":"60"}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(teks, `"KUNCILAMA":[1,2]`) || !strings.Contains(teks, `"MAXAGE":"70"`) {
		t.Errorf("kunci lama dipertahankan, kunci dikelola ditimpa:\n%s", teks)
	}
}

func TestBolakBalikInward(t *testing.T) {
	p := inwardUji()
	teks, err := RakitInward(p, "")
	if err != nil {
		t.Fatal(err)
	}
	q, err := UraiProduk(p.ID, `{}`, p.Inward.ID, teks)
	if err != nil {
		t.Fatal(err)
	}
	if q.Inward != p.Inward {
		t.Errorf("inward berubah menyeberang JSON:\n%+v\n%+v", q.Inward, p.Inward)
	}
}

func TestSQLInward(t *testing.T) {
	i := rata(sqlSisipInward("S.M_PRODUCTINWARD_LIFE"))
	if !strings.Contains(i, "INSERT INTO S.M_PRODUCTINWARD_LIFE (ID, JSONDATA) VALUES (:1, :2)") {
		t.Errorf("INSERT inward: %s", i)
	}
	u := rata(sqlPerbaruiInward("S.M_PRODUCTINWARD_LIFE"))
	if !strings.Contains(u, "UPDATE S.M_PRODUCTINWARD_LIFE SET JSONDATA = :1 WHERE ID = :2") {
		t.Errorf("UPDATE inward dikunci ID barisnya: %s", u)
	}
	args := argPerbaruiInward("100009", "{}")
	if c, ok := args[0].(go_ora.Clob); !ok || !c.Valid || args[1] != "100009" {
		t.Errorf("JSONDATA CLOB lalu ID: %#v", args)
	}
}
