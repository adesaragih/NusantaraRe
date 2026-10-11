package repository_test

// ⛔ FILTER YANG TIDAK DIKIRIM HARUS DILEWATI — `BrowseTreatyArrangement_
// ParentReinsMasterTrt` punya delapan filter (`B AND C AND D AND (E OR A) AND
// F AND G AND H`) dan NOL di antaranya ber-`pyUseNullIfEmpty`.
//
// Parameter yang BENAR-BENAR dikirim, dibaca ulang dari ekspor 8 Oktober 2026:
//
//	Section/Share.xml        (Spreading Type XOL)
//	    TreatyYear · TreatyGroupID=.TreatyGroupList(1).TreatyGroupID ·
//	    TreatyDescID="10001" · StartDate · ReinsTypeID="10246"
//	Section/DetailShare.xml  (Spreading Type Prop)
//	    sama, TreatyGroupID=.TreatyGroupID
//	kedua dropdown Reins Type MANUAL
//	    TreatyGroupID=TempSprd.TreatyGroupID — halaman yang tak pernah diisi,
//	    jadi filter grup DILEWATI; TreatyDescID tetap "10001"
//	Activity/FetchQSfromMaster(XOL)
//	    TreatyGroupID ✓ · TreatyDescID "10001" · StartDate · ReinsTypeID ""
//
// ⚠️ PEMBACAAN PERTAMA KELIRU, dan uji ini ada supaya kekeliruan itu tidak
// terulang: jendela baca yang terlalu sempit memotong daftar parameter dan
// menyisakan dua yang terakhir, sehingga `TreatyGroupID` dan `TreatyDescID`
// terbaca "tidak dikirim" lalu sempat dibuang dari layar. Akibatnya dropdown
// terisi dengan induk DI LUAR Treaty Group barisnya — dan `FetchQSfromMaster`
// mencari induk itu dengan filter grup yang sama, jadi Spreading Type
// terpilih sementara grid spreadingnya tetap kosong.

import (
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/repository"
)

func TestSaringanIndukSpreadingMelewatiParameterKosong(t *testing.T) {
	// ⚠️ Uji ini menjaga FUNGSINYA, bukan apa yang layar kirim. Layar sendiri
	// sengaja TIDAK mengirim grup sejak 8 Oktober 2026 (keputusan pemilik
	// proses, lihat `fetchQS`) — tapi kemampuan menyaring per grup harus tetap
	// benar, sebab yang dibuang adalah pemakaiannya, bukan filternya.
	//
	// Grup DAN desc dikirim: keduanya terpasang.
	w, arg := repository.SaringanIndukSpreading("10007", "10001", "20250301", repository.IndukDikecualikanDropdown)
	for _, wajib := range []string{"TREATYGROUPID", "TREATYDESCID", "NVL(p.REINSTYPEID"} {
		if !strings.Contains(w, wajib) {
			t.Errorf("dropdown kehilangan %s:\n%s", wajib, w)
		}
	}
	if len(arg) != 5 {
		t.Errorf("argumen = %v, mau [grup desc mulai mulai kecuali]", arg)
	}

	// Dropdown Reins Type MANUAL: `TempSprd.TreatyGroupID` tak pernah diisi,
	// jadi filter grup dilewati — tapi desc tetap dikirim.
	w, arg = repository.SaringanIndukSpreading("", "10001", "20250301", repository.IndukDikecualikanDropdown)
	if strings.Contains(w, "TREATYGROUPID") {
		t.Errorf("grup kosong tetap disaring — filternya harus DILEWATI:\n%s", w)
	}
	if !strings.Contains(w, "TREATYDESCID") {
		t.Errorf("dropdown manual kehilangan TREATYDESCID:\n%s", w)
	}
	if len(arg) != 4 || arg[0] != "10001" {
		t.Errorf("argumen = %v, mau [10001 mulai mulai kecuali]", arg)
	}

	// D, E dan A tidak berparameter — ketiganya selalu berlaku.
	for _, wajib := range []string{"PARENTREINSTYPEID = '00'", "'TRT'", "'ORS'", "STARTDATE", "ENDDATE"} {
		if !strings.Contains(w, wajib) {
			t.Errorf("kehilangan %s:\n%s", wajib, w)
		}
	}

	// ⛔ `TreatyDescID` BERPARAMETER, bukan konstanta yang dipakukan di SQL:
	// RD-nya `C .TreatyDescID = Param.TreatyDescID`. Memakukannya membuat
	// pemanggil yang tidak mengirimnya ikut tersaring.
	w, arg = repository.SaringanIndukSpreading("", "", "20250301", "")
	for _, dilarang := range []string{"TREATYGROUPID", "TREATYDESCID", "NVL(p.REINSTYPEID"} {
		if strings.Contains(w, dilarang) {
			t.Errorf("%s disaring padahal parameternya kosong:\n%s", dilarang, w)
		}
	}
	if len(arg) != 2 {
		t.Errorf("argumen = %v, mau dua tanggal saja", arg)
	}
}

// ⭐ 9 Oktober 2026 — Commencement panel New Adjustment tiba sebagai bentuk
// kabel / kotak tanggal; dibandingkan sebagai TEKS `YYYYMMDD`, ia harus
// dinormalkan dulu atau seluruh susunan tersaring habis.
func TestSaringanIndukSpreadingMenormalkanTanggal(t *testing.T) {
	for _, masuk := range []string{"20261008", "08-10-2026", "2026-10-08", "08/10/2026", "20261008T000000.000 GMT"} {
		_, arg := repository.SaringanIndukSpreading("10007", "10001", masuk, "")
		if len(arg) < 4 || arg[2] != "20261008" || arg[3] != "20261008" {
			t.Errorf("%q → %v", masuk, arg)
		}
	}
}
