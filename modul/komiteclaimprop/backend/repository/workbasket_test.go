package repository

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// Keputusan work owner 09-10-2026 (tangga komite PROP ke workbasket, migrasi claimprop 537): daftar kerja = baris
// berjalan ber-KomiteID akun ATAU workbasket aktif pelaku, tanpa kasus yang tingkat lainnya sudah diputus pelaku;
// keputusan menimpa KOMITE_OPERATORID dengan akun pemutus; T_WORK_CLAIM.POSITION mengikuti tingkat berjalan; email
// tingkat berikut ke semua anggota workbasket.

// cacahBind - cacah penampung `:n` (driver mengikat menurut urutan kemunculan; format tanggal juga memuat ':').
func cacahBind(q string) int { return len(regexp.MustCompile(`:\d+`).FindAllString(q, -1)) }

func TestDaftarKerjaMenyaringWorkbasketDanRangkap(t *testing.T) {
	q := sqlDaftarKerja("S.G", "S.W", "S.L", "S.C", "S.A", 2)
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(q, "l.KOMITE_OPERATORID IN (:2, :3)") {
		t.Errorf("baris berjalan ber-KomiteID workbasket pelaku: %s", q)
	}
	if !strings.Contains(q, "NOT EXISTS") || !strings.Contains(q, "l3.KOMITE_OPERATORID = :7") {
		t.Errorf("kasus yang tingkat lainnya sudah diputus pelaku tidak tampil: %s", q)
	}
	if n := cacahBind(q); n != 9 {
		t.Errorf("bind %d, mau 9 (akun, 2 peran, 3 keputusan, akun, LINI, awalan)", n)
	}
	if q0 := sqlDaftarKerja("S.G", "S.W", "S.L", "S.C", "S.A", 0); strings.Contains(q0, " IN (") || cacahBind(q0) != 7 {
		t.Errorf("tanpa workbasket: hanya akun (7 bind): %s", q0)
	}
}

func TestKeputusanMenimpaPemutusDanPosisiMengikutiTingkat(t *testing.T) {
	if q := sqlTulisAnggota("S.L", true); !strings.Contains(q, "KOMITE_OPERATORID = NVL(:4, KOMITE_OPERATORID)") ||
		!strings.Contains(q, "KOMITE_APPROVAL = :7") {
		t.Errorf("baris yang diputus menyimpan akun pemutus: %s", q)
	}
	if q := sqlTulisAnggota("S.L", false); strings.Contains(q, "KOMITE_OPERATORID") {
		t.Errorf("tolak otomatis S26.1 tidak menimpa KomiteID: %s", q)
	}
	if q := sqlSentuhKasus("S.W"); !strings.Contains(q, "POSITION = :2") {
		t.Errorf("kasus berjalan: POSITION = workbasket tingkat berikut: %s", q)
	}
	if q := sqlTutupKasus("S.W"); !strings.Contains(q, "POSITION = NULL") {
		t.Errorf("kasus selesai: POSITION dikosongkan: %s", q)
	}
	q := sqlEmailAnggotaWorkbasket("S.LWB", "S.WB", "S.LOGIN")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(q, "w.IS_ACTIVE = 1") || !strings.Contains(q, "m.IS_ACTIVE = '1'") || cacahBind(q) != 1 {
		t.Errorf("email anggota aktif workbasket aktif: %s", q)
	}
}
