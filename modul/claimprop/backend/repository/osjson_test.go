package repository

import (
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimprop/backend/models"
)

// DATA_JSON ikut disisipkan bila baris membawanya (keputusan work owner 08-10-2026), dan tidak disebut bila kosong.
func TestSisipOSMembawaDataJSON(t *testing.T) {
	isi := "{\n\"Type\":\"0\"\n}\n"
	q, args, err := sqlSisipOS("UJI.OS_AKSEPTASI_KLAIM", models.BarisOS{CaseID: "UJI-C", DataJSON: isi}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(q, "DATA_JSON") {
		t.Fatalf("DATA_JSON tidak disisipkan: %s", q)
	}
	ada := false
	for _, a := range args {
		if s, ok := a.(string); ok && s == isi {
			ada = true
		}
	}
	if !ada {
		t.Fatalf("isi DATA_JSON tidak terikat sebagai argumen: %v", args)
	}

	q, _, err = sqlSisipOS("UJI.OS_AKSEPTASI_KLAIM", models.BarisOS{CaseID: "UJI-C"}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(q, "DATA_JSON") {
		t.Fatalf("DATA_JSON kosong tetap disebut: %s", q)
	}
}
