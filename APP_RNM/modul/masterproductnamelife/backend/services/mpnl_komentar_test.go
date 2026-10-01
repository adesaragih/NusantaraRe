package services_test

// Komentar, dokumen klaim, underwriting finansial, lien clause (paket 7, tiket 07).

import (
	"context"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

var jamUji = time.Date(2026, 10, 1, 1, 2, 3, 456_000_000, time.UTC)

func TestSetiapSimpanMenambahSatuKomentarAddCommentList(t *testing.T) {
	l, g := layananMaster()
	l = l.DenganJam(func() time.Time { return jamUji })
	m := produkMasuk()
	m.Umum.Comment = "UJI pertimbangan"
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.CommentList) != 1 {
		t.Fatalf("satu baris: %+v", p.CommentList)
	}
	k := p.CommentList[0]
	if k.Date != "20261001T010203.456 GMT" || k.OperatorName != "UJI-PELAKU" || k.Suggest != "UJI pertimbangan" || k.IsApproved != "" {
		t.Errorf("AddCommentList_Act b235: Date @CurrentDateTime, OperatorName pxInsName, Suggest comment: %+v", k)
	}
	// Ubah: riwayat lama tetap di depan, baris baru di belakang - juga bila komentar kosong (OQ-MPNL-14).
	m2 := produkMasuk()
	m2.ID = p.ID
	m2.Umum.Comment = ""
	m2.CommentList = []models.BarisKomentar{{Suggest: "DIUBAH KLIEN"}}
	p2, err := l.SimpanProduk(context.Background(), pelakuUji, m2, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(p2.CommentList) != 2 || p2.CommentList[0].Suggest != "UJI pertimbangan" || p2.CommentList[1].Suggest != "" {
		t.Errorf("komentar tidak dapat diubah klien; setiap simpan menambah satu baris: %+v", p2.CommentList)
	}
	if !strings.Contains(g.Umum[p.ID], `"Comment":""`) {
		t.Errorf("medan Comment halaman ikut tersimpan seperti Pega: %s", g.Umum[p.ID])
	}
}

func TestDokumenLienDanUnderwritingFinansial(t *testing.T) {
	l, _ := layananMaster()
	m := produkMasuk()
	m.DocumentClaim = []models.BarisDokumen{{Document: "UJI KTP"}, {Document: "UJI SURAT"}}
	m.LienClause = []models.BarisLien{{Usia: "1", Manfaat: "50"}}
	m.FinancialUnderwriting = []models.BarisFinUW{{MinInsured: "1.000,5", MaxInsured: "2", Employee: "Y", NonEmployee: "N"}}
	_, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if !strings.Contains(services.Pesan(err), `FINANCIAL UNDERWRITING row 1: Min Insured "1.000,5" contains both a comma and a dot`) {
		t.Errorf("angka underwriting finansial: %v", err)
	}
	m.FinancialUnderwriting = []models.BarisFinUW{{MinInsured: "5", MaxInsured: "2"}}
	_, err = l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if !strings.Contains(services.Pesan(err), "FINANCIAL UNDERWRITING row 1: Min Insured must not be greater than Max Insured") {
		t.Errorf("rentang underwriting finansial (R17): %v", err)
	}
	m.FinancialUnderwriting = []models.BarisFinUW{{MinInsured: "1", MaxInsured: "2", Employee: "Y", NonEmployee: "N"}}
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.DocumentClaim) != 2 || p.DocumentClaim[1].Document != "UJI SURAT" || p.LienClause[0].Manfaat != "50" ||
		p.FinancialUnderwriting[0].NonEmployee != "N" {
		t.Errorf("ketiga daftar tersimpan berurutan: %+v %+v %+v", p.DocumentClaim, p.LienClause, p.FinancialUnderwriting)
	}
}
