package repository

import (
	"os"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// TestTahapDanCatatanPosisi - tiket 77: tahap dari FLAG_ONGOING_POLICY (Go = SQL), POSITION_NOTE dari PosisiBaru services ('Offer' ->
// ReasFacInMarketing, Decline / "" -> NULL, lainnya apa adanya).
func TestTahapDanCatatanPosisi(t *testing.T) {
	for flag, mau := range map[string]string{"0": "Offer", "2": "Binding", "1": "Policy", "": "Policy", "9": "Policy"} {
		if got := TahapDariFlag(flag); got != mau {
			t.Errorf("flag %q: %q", flag, got)
		}
	}
	if q := sqlTahapDariFlag("F"); q != "CASE F WHEN '0' THEN 'Offer' WHEN '2' THEN 'Binding' ELSE 'Policy' END" {
		t.Errorf("%q", q)
	}
	for masuk, mau := range map[string]any{"Offer": "ReasFacInMarketing", "Decline": nil, "": nil, "ReasFacInTeamLeader": "ReasFacInTeamLeader",
		"ReasFacInMarketing": "ReasFacInMarketing", "ReasFacInJuniorUnderwritingB": "ReasFacInJuniorUnderwritingB"} {
		if got := CatatanPosisiDari(masuk); got != mau {
			t.Errorf("catatan %q: %v", masuk, got)
		}
	}
}

// TestPosisiEfektifBarisCampuran - tiket 77: kotak masuk / penjaga membaca kosakata lama dari baris lama (belum dikonversi), baru,
// dan Decline (aturan kontrak 9d: cadangan POSITION lama hanya untuk case aktif).
func TestPosisiEfektifBarisCampuran(t *testing.T) {
	for _, c := range []struct{ nama, note, posisi, flag, status, mau string }{
		{"baru marketing Offer", "ReasFacInMarketing", "Offer", "0", "", "Offer"},
		{"baru marketing Binding", "ReasFacInMarketing", "Binding", "2", "", "ReasFacInMarketing"},
		{"baru Team Leader", "ReasFacInTeamLeader", "Offer", "0", "", "ReasFacInTeamLeader"},
		{"baru TL Policy", "ReasFacInTeamLeader", "Policy", "1", "", "ReasFacInTeamLeader"},
		{"baru selesai Policy", "ReasFacInTeamLeader", "Policy", "1", "Resolved-Completed", "ReasFacInTeamLeader"},
		{"baru Decline", "", "Offer", "0", "Resolved-Decline", "Decline"},
		{"baru Decline Binding", "", "Binding", "2", "Resolved-Decline", "Decline"},
		{"lama marketing", "", "Offer", "0", "", "Offer"},
		{"lama Team Leader", "", "ReasFacInTeamLeader", "0", "Pending-Policy", "ReasFacInTeamLeader"},
		{"lama Binding", "", "ReasFacInMarketing", "2", "", "ReasFacInMarketing"},
		{"lama Reject TL tahap Offer", "", "ReasFacInMarketing", "0", "", "Offer"},
		{"lama Decline", "", "Decline", "0", "Resolved-Decline", "Decline"},
		{"lama selesai", "", "ReasFacInUnderwriting", "1", "Resolved-Completed", "ReasFacInUnderwriting"},
	} {
		if got := PosisiEfektif(c.note, c.posisi, c.flag, c.status); got != c.mau {
			t.Errorf("%s: %q mau %q", c.nama, got, c.mau)
		}
	}
	if q := sqlPosisiEfektif("w"); q != "CASE WHEN w.POSITION_NOTE IS NULL AND w.STATUS_WORK LIKE 'Resolved-%' AND w.POSITION IN "+
		"('Offer', 'Binding', 'Policy') THEN CASE WHEN w.STATUS_WORK = 'Resolved-Decline' THEN 'Decline' END WHEN NVL(w.POSITION_NOTE, "+
		"w.POSITION) = 'ReasFacInMarketing' AND w.FLAG_ONGOING_POLICY = '0' THEN 'Offer' ELSE NVL(w.POSITION_NOTE, w.POSITION) END" {
		t.Errorf("%q", q)
	}
}

// TestSQLCatatanPosisi - tiket 77: tiap pembaca / penulis punya varian ber-posisi efektif yang HANYA mengganti kolom posisi (teks lama
// = uji lama, tanpa 254); bind tidak bergeser; sinkron tahap; uji kolom lolos PeriksaSQL.
func TestSQLCatatanPosisi(t *testing.T) {
	e := sqlPosisiEfektif("w")
	tb := tabelSubmit{work: "W", mo: "M", quo: "Q", general: "G"}
	for nama, c := range map[string]struct {
		lama, baru string
		bind       int
	}{
		"baca kasus": {sqlBacaKasus("W", "O", "A", "G", "Q", "S"), sqlBacaKasusCatatan("W", "O", "A", "G", "Q", "S"), 2},
		"tahap":      {sqlTahapKasus("W", "G"), sqlTahapKasusCatatan("W", "G"), 2},
		"konteks":    {sqlKonteksSubmit(tb), sqlKonteksSubmitCatatan(tb), 2},
	} {
		if strings.Count(c.baru, e) != 1 || strings.Contains(c.baru, "w.POSITION,") || strings.Replace(c.baru, e, "w.POSITION", 1) != c.lama ||
			bindBerurutan(t, nama, c.baru) != c.bind {
			t.Errorf("%s:\n%s", nama, c.baru)
		}
	}
	if q := sqlPosisiSubmitCatatan("W"); q != "UPDATE W w SET POSITION_NOTE = :1, STATUS_WORK = NVL(:2, STATUS_WORK), TGL_UPDATE = SYSDATE"+
		" WHERE w.ID = :3 AND w.LINI = :4 AND "+e+" = :5 AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT LIKE 'Resolved-%')" ||
		bindBerurutan(t, "posisi", q) != 5 {
		t.Errorf("%q", q)
	}
	if q := sqlKepalaScoringCatatan("W"); q != "SELECT "+e+", w.FLAG_ONGOING_POLICY FROM W w WHERE w.ID = :1 AND w.LINI = :2" {
		t.Errorf("%q", q)
	}
	if q := sqlSentuhTahapCatatan("W"); q != "UPDATE W w SET TGL_UPDATE = SYSDATE WHERE w.ID = :1 AND w.LINI = :2 AND "+e+
		" = :3 AND w.FLAG_ONGOING_POLICY = :4" {
		t.Errorf("%q", q)
	}
	if q := sqlSinkronTahap("W"); q != "UPDATE W SET POSITION = CASE FLAG_ONGOING_POLICY WHEN '0' THEN 'Offer' WHEN '2' THEN 'Binding'"+
		" ELSE 'Policy' END WHERE ID = :1 AND LINI = :2" {
		t.Errorf("%q", q)
	}
	q := sqlSisipCaseCatatan("W")
	arg := argSisipCaseCatatan("NB-1", "UJI-USER")
	if bindBerurutan(t, "sisip", q) != len(arg) || !strings.Contains(q, "(ID, LINI, POSITION, POSITION_NOTE, FLAG_ONGOING_POLICY,") ||
		arg[2] != "Offer" || arg[3] != "ReasFacInMarketing" || arg[4] != FlagAwalCaseNB {
		t.Errorf("%q %v", q, arg)
	}
	if err := db.PeriksaSQL(sqlUjiCatatanPosisi("W")); err != nil || sqlUjiCatatanPosisi("W") != "SELECT POSITION_NOTE FROM W WHERE ROWNUM < 1" {
		t.Errorf("uji kolom: %v", err)
	}
	if gantiSekali("a X b X", "X", "Y") != "a X b X" || gantiSekali("a X b", "X", "Y") != "a Y b" {
		t.Error("gantiSekali")
	}
}

// TestKotakMasukCatatanPosisi - tiket 77: kotak masuk memakai posisi efektif di saringan antrean + team dan GROUP BY bila 254 ada;
// tanpa 254 teks lama (w.POSITION).
func TestKotakMasukCatatanPosisi(t *testing.T) {
	s := SaringanKotakMasuk{Tim: []AntreanTim{{Posisi: []string{"Offer", "ReasFacInMarketing"}, Team: []string{"1"}}},
		Lain: []string{"ReasFacInTeamLeader"}}
	tb := tabelPortal{work: "W", opp: "O", akun: "A", quo: "Q", mo: "M", general: "G"}
	lama, argLama := sqlHitungKotakMasuk(tb, s)
	tb.posisi = sqlPosisiEfektif("w")
	baru, argBaru := sqlHitungKotakMasuk(tb, s)
	if strings.Count(lama, "w.POSITION") != 4 || strings.Count(baru, sqlPosisiEfektif("w")) != 4 ||
		strings.ReplaceAll(baru, sqlPosisiEfektif("w"), "w.POSITION") != lama || len(argLama) != len(argBaru) {
		t.Errorf("lama:\n%s\nbaru:\n%s", lama, baru)
	}
}

// TestTulisSubmitCatatanPosisi - tiket 77 (teks sumber; tanpa driver tiruan): TulisSubmit memilih varian menurut deteksi kolom, NOTE
// dari CatatanPosisiDari, dan menyinkronkan tahap SESUDAH binding / policy / validitas / total (flag akhir) dan SEBELUM kunci
// korespondensi; tanpa 254 tidak ada sinkron.
func TestTulisSubmitCatatanPosisi(t *testing.T) {
	b, err := os.ReadFile("submit.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, x := range []string{"catatan := adaCatatanPosisi(ctx, r.db, t.work)",
		"qPosisi, posisiBaru = sqlPosisiSubmitCatatan(t.work), CatatanPosisiDari(w.PosisiBaru)",
		"if catatan { // tiket 77"} {
		if !strings.Contains(s, x) {
			t.Errorf("submit.go tanpa %q", x)
		}
	}
	i, j, k := strings.Index(s, "if b := w.Binding; b != nil {"), strings.Index(s, "sqlSinkronTahap(t.work)"),
		strings.Index(s, "sqlKunciPenerima(t.pen, t.kor), \"kunci penerima")
	if i < 0 || j < i || k < j || strings.Index(s, "sqlDateOffer(t.general)") > j {
		t.Errorf("urutan binding %d, sinkron %d, kunci %d", i, j, k)
	}
}

// TestPembacaMemakaiDeteksi - tiket 77 (teks sumber): setiap pembaca / penulis posisi memilih varian menurut adaCatatanPosisi.
func TestPembacaMemakaiDeteksi(t *testing.T) {
	for berkas, mau := range map[string]int{"kasus.go": 2, "submit.go": 2, "kotakmasuk.go": 2, "scoring.go": 1, "binding.go": 1, "casenb.go": 1} {
		b, err := os.ReadFile(berkas)
		if err != nil {
			t.Fatal(err)
		}
		if n := strings.Count(string(b), "adaCatatanPosisi(ctx, r.db,"); n != mau {
			t.Errorf("%s: %d pemanggilan adaCatatanPosisi, mau %d", berkas, n, mau)
		}
	}
}

// TestKonversiMempertahankanPosisiEfektif - docs/alat/konversi_position_note_254*.sql: untuk setiap kombinasi baris lama yang
// dihasilkan sistem sebelum tiket 77 (posisi x flag x status), posisi efektif SEBELUM = SESUDAH konversi (NOTE = CASE POSITION ...,
// POSITION = tahap flag) dan baris hasil konversi tidak tercakup lagi (idempoten). Teks UPDATE skrip = aturan yang diuji di sini;
// cakupan LINI 'FAC' saja (FACEDM milik endorsmentfacin).
func TestKonversiMempertahankanPosisiEfektif(t *testing.T) {
	konversi := func(posisi string) string {
		switch posisi {
		case "Offer":
			return "ReasFacInMarketing"
		case "Decline":
			return ""
		}
		return posisi
	}
	cakup := func(note, posisi, status string) bool {
		return note == "" && !(strings.HasPrefix(status, "Resolved-") && (posisi == "Offer" || posisi == "Binding" || posisi == "Policy"))
	}
	lama := []struct{ posisi, flag, status string }{
		{"Offer", "0", ""}, {"Offer", "0", "Pending-Policy"}, {"ReasFacInMarketing", "0", ""}, {"ReasFacInMarketing", "2", ""},
		{"ReasFacInTeamLeader", "0", ""}, {"ReasFacInTeamLeader", "1", ""}, {"ReasFacInJuniorUnderwritingB", "0", ""},
		{"ReasFacInUnderwritingFinancial", "1", ""}, {"Decline", "0", "Resolved-Decline"}, {"Decline", "2", "Resolved-Decline"},
		{"ReasFacInTeamLeader", "1", "Resolved-Completed"},
	}
	for _, b := range lama {
		if !cakup("", b.posisi, b.status) {
			t.Errorf("%+v tidak tercakup konversi", b)
			continue
		}
		sebelum := PosisiEfektif("", b.posisi, b.flag, b.status)
		sesudah := PosisiEfektif(konversi(b.posisi), TahapDariFlag(b.flag), b.flag, b.status)
		if sebelum != sesudah {
			t.Errorf("%+v: sebelum %q sesudah %q", b, sebelum, sesudah)
		}
		if cakup(konversi(b.posisi), TahapDariFlag(b.flag), b.status) {
			t.Errorf("%+v: tercakup lagi sesudah konversi (tidak idempoten)", b)
		}
	}
	for _, berkas := range []string{"konversi_position_note_254.sql", "konversi_position_note_254_pooldata.sql"} {
		b, err := os.ReadFile("../../docs/alat/" + berkas)
		if err != nil {
			t.Fatal(err)
		}
		s := strings.ReplaceAll(string(b), "\r", "")
		var kode []string
		for _, baris := range strings.Split(s, "\n") {
			if x := strings.TrimSpace(baris); x != "" && !strings.HasPrefix(x, "--") {
				kode = append(kode, x)
			}
		}
		k := strings.Join(kode, " ")
		if strings.Contains(strings.ToUpper(k), "COMMIT") || strings.Count(k, "UPDATE ") != 1 || strings.Contains(k, "DELETE") ||
			!strings.Contains(k, "SET w.POSITION_NOTE = CASE w.POSITION WHEN 'Offer' THEN 'ReasFacInMarketing' WHEN 'Decline' THEN NULL ELSE w.POSITION END,") ||
			!strings.Contains(k, "w.POSITION = CASE w.FLAG_ONGOING_POLICY WHEN '0' THEN 'Offer' WHEN '2' THEN 'Binding' ELSE 'Policy' END") ||
			strings.Count(k, "w.LINI = 'FAC' AND w.ID LIKE 'NB-%' AND w.POSITION_NOTE IS NULL AND NOT (w.STATUS_WORK LIKE 'Resolved-%' AND w.POSITION IN ('Offer', 'Binding', 'Policy'))") != 4 ||
			strings.Count(k, sqlPosisiEfektif("w")) != 4 || strings.Contains(k, "FACEDM") {
			t.Errorf("%s:\n%s", berkas, k)
		}
	}
}

// TestJalurTulisPulangPergi - tiket 77: tiap jalur tulis posisi (buat case, Submit semua tahap, Binding, Decline, Selesai Policy):
// PosisiBaru services + flag AKHIR -> POSITION_NOTE (CatatanPosisiDari) dan POSITION (TahapDariFlag, sinkron di akhir transaksi) ->
// posisi efektif yang dibaca balik = kosakata services (marketing tahap Offer dibaca "Offer"; services menyamakan keduanya).
func TestJalurTulisPulangPergi(t *testing.T) {
	for _, c := range []struct{ jalur, posisiBaru, flag, status, note, tahap, baca string }{
		{"buat case", "Offer", "0", "", "ReasFacInMarketing", "Offer", "Offer"},
		{"Submit marketing Accept", "ReasFacInTeamLeader", "0", "", "ReasFacInTeamLeader", "Offer", "ReasFacInTeamLeader"},
		{"Submit TL Accept", "ReasFacInJuniorUnderwritingB", "0", "", "ReasFacInJuniorUnderwritingB", "Offer", "ReasFacInJuniorUnderwritingB"},
		{"Submit TL Reject", "ReasFacInMarketing", "0", "", "ReasFacInMarketing", "Offer", "Offer"},
		{"Submit marketing Reject", "Offer", "0", "", "ReasFacInMarketing", "Offer", "Offer"},
		{"Accept UW -> Binding (Transition38)", "ReasFacInMarketing", "2", "", "ReasFacInMarketing", "Binding", "ReasFacInMarketing"},
		{"Binding Accept -> Policy", "ReasFacInTeamLeader", "1", "", "ReasFacInTeamLeader", "Policy", "ReasFacInTeamLeader"},
		{"Binding Banding -> Offer", "ReasFacInSeniorUnderwriting", "0", "", "ReasFacInSeniorUnderwriting", "Offer", "ReasFacInSeniorUnderwriting"},
		{"TL Policy Reject -> Binding", "ReasFacInMarketing", "2", "", "ReasFacInMarketing", "Binding", "ReasFacInMarketing"},
		{"Decline tahap Offer", "Decline", "0", "Resolved-Decline", "", "Offer", "Decline"},
		{"Decline tahap Binding", "Decline", "2", "Resolved-Decline", "", "Binding", "Decline"},
		{"Selesai Policy", "ReasFacInTeamLeader", "1", "Resolved-Completed", "ReasFacInTeamLeader", "Policy", "ReasFacInTeamLeader"},
	} {
		note, _ := CatatanPosisiDari(c.posisiBaru).(string)
		tahap := TahapDariFlag(c.flag)
		if note != c.note || tahap != c.tahap {
			t.Errorf("%s: NOTE %q POSITION %q", c.jalur, note, tahap)
		}
		if got := PosisiEfektif(note, tahap, c.flag, c.status); got != c.baca {
			t.Errorf("%s: dibaca %q mau %q", c.jalur, got, c.baca)
		}
	}
}
