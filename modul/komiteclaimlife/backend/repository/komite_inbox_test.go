package repository

// Inbox Komite - tiket 01 Komite Claim Life. TANPA Oracle.

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
)

func penandaBerurut(q string) []string {
	var k []string
	for _, m := range regexp.MustCompile(`:([a-z]+)\b`).FindAllStringSubmatch(q, -1) {
		k = append(k, m[1])
	}
	return k
}

// TestInboxKomitePenandaBerurutSesuaiArgumen - driver mengikat berurutan.
//
// ⛔ `Ambil` mengirim argumen menurut urutan MUNCULNYA penanda. Bila teks
// SQL-nya diubah tanpa argumennya, `:menunggu` kedua menerima nilai
// `:tutup` - dan inbox diam-diam kosong.
func TestInboxKomitePenandaBerurutSesuaiArgumen(t *testing.T) {
	baris := penandaBerurut(sqlInboxKomite("G", "W", "L", "C", "A"))
	mau := []string{"akun", "menunggu", "menunggu", "tutup", "setuju", "offset", "ukuran"}
	if !reflect.DeepEqual(baris, mau) {
		t.Errorf("penanda inbox %v, mau %v", baris, mau)
	}
	cacah := penandaBerurut(sqlCacahInboxKomite("G", "W", "L", "C", "A"))
	if !reflect.DeepEqual(cacah, mau[:5]) {
		t.Errorf("penanda cacah %v, mau %v", cacah, mau[:5])
	}
}

// TestInboxKomiteHanyaAnggotaBerjalan - `KomiteRouter`, AC 10 spec.
func TestInboxKomiteHanyaAnggotaBerjalan(t *testing.T) {
	q := sqlInboxKomite("G", "W", "L", "C", "A")
	for _, s := range []string{
		"l.KOMITE_OPERATORID = :akun",
		"l.KOMITE_APPROVAL = :menunggu",
		"l.KOMITE_URUT = (SELECT MIN(l2.KOMITE_URUT)",
		"l2.KOMITE_APPROVAL = :menunggu",
		"w.STATUS_WORK <> :tutup",
	} {
		if !strings.Contains(q, s) {
			t.Errorf("inbox komite kehilangan %q:\n%s", s, q)
		}
	}
	if ApprovalKomiteMenunggu != "0" {
		t.Errorf("approval menunggu %q, mau \"0\" (CreateKMTLife_Act b912)", ApprovalKomiteMenunggu)
	}
}

// TestKasusKomiteTidakMembacaEmail - alamat orang tidak dibaca tanpa perlu.
func TestKasusKomiteTidakMembacaEmail(t *testing.T) {
	for _, q := range []string{
		sqlInboxKomite("G", "W", "L", "C", "A"),
		sqlKasusKomite("G", "W", "L", "C", "A"),
		sqlTanggaKasus("L"),
	} {
		if strings.Contains(q, "EMAIL") {
			t.Errorf("pembaca komite membaca kolom email:\n%s", q)
		}
	}
}

// TestUangKomiteLewatTeks - ADR-U-0003.
func TestUangKomiteLewatTeks(t *testing.T) {
	for _, q := range []string{kolomBarisKomite, sqlKasusKomite("G", "W", "L", "C", "A")} {
		if !strings.Contains(q, "TO_CHAR(a.CLAIM_AMOUNT, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')") {
			t.Errorf("CLAIM_AMOUNT tidak dibaca lewat TO_CHAR ber-NLS:\n%s", q)
		}
	}
}

// TestKeputusanKomiteBersyaratDuaBaris - tiket 02.
//
// ⛔ Anak tangga hanya ditulis bila MASIH menunggu dan milik pelaku; kepala
// hanya dimajukan dari count yang dibaca. Dua keputusan serentak: satu kalah.
func TestKeputusanKomiteBersyaratDuaBaris(t *testing.T) {
	q := sqlCatatAnakTangga("L")
	for _, s := range []string{"DATA_KOMITE_ID = :4", "KOMITE_URUT = :5",
		"KOMITE_OPERATORID = :6", "KOMITE_APPROVAL = :7"} {
		if !strings.Contains(q, s) {
			t.Errorf("penulisan anak tangga tanpa %q:\n%s", s, q)
		}
	}
	if !strings.Contains(sqlMajukanTangga("G"), "WHERE ID = :3 AND KOMITE_COUNT = :4") {
		t.Errorf("kepala komite dimajukan tanpa kunci count:\n%s", sqlMajukanTangga("G"))
	}
}

// TestInboxKomiteMengikutiIsKomiteLoop - Tolak di tengah tidak jatuh ke tingkat berikut.
func TestInboxKomiteMengikutiIsKomiteLoop(t *testing.T) {
	q := sqlInboxKomite("G", "W", "L", "C", "A")
	if !strings.Contains(q, "g.ACCEPT_STATUS IS NULL") ||
		!strings.Contains(q, "g.ACCEPT_STATUS = :setuju AND g.KOMITE_COUNT <= g.KOMITE_LOOP") {
		t.Errorf("inbox komite tidak menegakkan IsKomiteLoop:\n%s", q)
	}
}

// TestEskalasiMengosongkanBukanMemutuskan - tiket 03.
func TestEskalasiMengosongkanBukanMemutuskan(t *testing.T) {
	q := sqlLewatiAnakTangga("L")
	if !strings.Contains(q, "SET KOMITE_APPROVAL = NULL") || !strings.Contains(q, "KOMITE_APPROVAL = :3") {
		t.Errorf("eskalasi tidak mengosongkan anak tangga yang menunggu:\n%s", q)
	}
	g := sqlNaikkanTingkat("G")
	if strings.Contains(g, "ACCEPT_STATUS") || !strings.Contains(g, "KOMITE_COUNT = :3") {
		t.Errorf("eskalasi menyentuh ACCEPT_STATUS atau tanpa kunci count:\n%s", g)
	}
}

// TestKasusKomiteMembacaPesertaDanNomorDiperiksa - tiket 04a.
func TestKasusKomiteMembacaPesertaDanNomorDiperiksa(t *testing.T) {
	if !strings.Contains(sqlKasusKomite("G", "W", "L", "C", "A"), "a.PREMIUM_LIST_DETAIL_ID") {
		t.Error("kasus komite tidak membaca peserta pemilik baris")
	}
	if !strings.Contains(sqlNomorAksepDiAdjustment("A"), "WHERE ACCEPTED_NO = :1") {
		t.Error("pemeriksa keunikan nomor akseptasi tidak menyaring ACCEPTED_NO")
	}
}

// TestRekamAkhirSatuJalurDanDipagari - tiket 04b.
func TestRekamAkhirSatuJalurDanDipagari(t *testing.T) {
	q := sqlRekamAkhirWarisan("D")
	if !strings.Contains(q, "SET STS_REJECT = :1, NO_ACCEPTATION = :2, ACCEPTATION_DATE = :3") ||
		!strings.Contains(q, "WHERE ID = :4") || strings.Contains(strings.ToUpper(q), "INSERT") {
		t.Errorf("rekam akhir bukan UPDATE baris datar yang ada:\n%s", q)
	}
	r := &InboxKomite{}
	for _, s := range []string{"", "0", "3", "Aksep"} {
		if err := r.RekamAkhirWarisan(context.Background(), &db.Tx{}, "A", s, "", time.Now()); !errors.Is(err, ErrStatusAkhirKomiteTidakSah) {
			t.Errorf("status %q: %v", s, err)
		}
	}
}

// OQ-K-05 (GILIRAN-17): langkah 5.1 - SATU UPDATE bersyarat atas seluruh
// tingkat yang MEMUTUS; tingkat yang dilewati eskalasi (NULL) tidak disentuh.
func TestSQLTimpaTanggaTolakAkhir(t *testing.T) {
	q := sqlTimpaTanggaTolakAkhir("S.L")
	for _, mau := range []string{"UPDATE S.L", "SET KOMITE_APPROVAL = :1, KOMITE_COMMENT = NULL, DATE_APPROVE = :2",
		"WHERE DATA_KOMITE_ID = :3 AND KOMITE_APPROVAL IS NOT NULL"} {
		if !strings.Contains(q, mau) {
			t.Errorf("tanpa %q:\n%s", mau, q)
		}
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Error(err)
	}
	b := sqlTanggaSebelumDitimpa("S.L")
	for _, mau := range []string{"SELECT KOMITE_URUT, KOMITE_OPERATORID, KOMITE_APPROVAL, KOMITE_COMMENT",
		"FROM S.L WHERE DATA_KOMITE_ID = :1", "ORDER BY KOMITE_URUT"} {
		if !strings.Contains(b, mau) {
			t.Errorf("tangga tanpa %q:\n%s", mau, b)
		}
	}
	r := NewInboxKomite(nil)
	if _, err := r.TanggaSebelumDitimpa(context.Background(), nil, "K"); err == nil {
		t.Error("membaca tangga tanpa transaksi diterima")
	}
	if err := r.TimpaTanggaTolakAkhir(context.Background(), nil, "K", time.Now()); err == nil {
		t.Error("menimpa tangga tanpa transaksi diterima")
	}
}
