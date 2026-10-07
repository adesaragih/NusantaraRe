package models

// Uji seam 3 - penanda persetujuan, tangga, nomor polis, batas Sec Head
// (tiket 02, 03, 06, 09, 10). Nilai harapan dari connector Flow dan tabel
// keputusan di INVENTARIS-XML.md bab 4 dan 10.

import (
	"testing"
	"time"
)

func TestIsApprovedNolDitolak(t *testing.T) { // AC 1
	if Disetujui("0") {
		t.Fatal(`"0" harus ditolak`)
	}
}

func TestIsApprovedSelainNolDisetujuiTermasukKosong(t *testing.T) { // AC 2, 84
	for _, v := range []string{"", "1", "2", "x"} {
		if !Disetujui(v) {
			t.Errorf("%q harus disetujui (bawaan YES)", v)
		}
	}
}

func TestIsApprovedDibandingkanSebagaiTeks(t *testing.T) { // AC 3, 4
	// pembanding numerik akan menolak " 0" dan "0.0"; tabel keputusan
	// berkolom text tidak.
	for _, v := range []string{" 0", "0.0", "00"} {
		if !Disetujui(v) {
			t.Errorf("%q bukan teks \"0\" - harus disetujui", v)
		}
	}
	// When/isApproved menguji "= 1": nilai kosong TIDAK disetujui di sana.
	// Tabel keputusan yang hidup menyetujuinya.
	if !Disetujui("") {
		t.Fatal("perilaku When/isApproved (kosong = tolak) dipakai - harus tabel keputusan")
	}
}

func TestTanggaAdminMenolakSelesaiDitolak(t *testing.T) { // AC 5
	tr, err := Langkah(PosisiAdmin, "0", false)
	if err != nil {
		t.Fatal(err)
	}
	if tr.StatusTutup != StatusDitolak || tr.PosisiBaru != "" {
		t.Fatalf("admin menolak: %+v, harap ditutup Resolved-Rejected", tr)
	}
}

func TestTanggaAtasanMenolakKembaliKeAdmin(t *testing.T) { // AC 6
	for _, pos := range []string{PosisiSecHead, PosisiDeptHead} {
		tr, err := Langkah(pos, "0", true)
		if err != nil {
			t.Fatal(err)
		}
		if tr.PosisiBaru != PosisiAdmin || tr.Ditutup() || !tr.KembaliKePembuat {
			t.Errorf("%s menolak: %+v, harap kembali ke admin, NBStatus menunjuk pembuat", pos, tr)
		}
	}
}

func TestTanggaAdminMenyetujuiNaikKeSecHead(t *testing.T) { // AC 7
	tr, _ := Langkah(PosisiAdmin, "1", true)
	if tr.PosisiBaru != PosisiSecHead {
		t.Fatalf("admin menyetujui -> %q, harap Sec Head (tidak boleh melompat)", tr.PosisiBaru)
	}
}

// AC 8 - keputusan work owner, bertentangan dengan XML (Decision13 /
// CekLimitTreatyAcc_Act); WO diikuti: Sec Head menyetujui SELALU naik ke Dept
// Head, bernomor polis atau tidak. RALAT di tiket 03.
func TestTanggaSecHeadMenyetujuiSelaluNaikKeDeptHead(t *testing.T) { // AC 8
	for _, bernomor := range []bool{true, false} {
		for _, ia := range []string{"1", "", "2"} {
			tr, _ := Langkah(PosisiSecHead, ia, bernomor)
			if tr.PosisiBaru != PosisiDeptHead || tr.Ditutup() {
				t.Fatalf("Sec Head menyetujui (IsApproved %q, bernomor %v) -> %+v, harap Dept Head", ia, bernomor, tr)
			}
		}
	}
}

func TestTanggaDeptHeadMenyetujuiSelesai(t *testing.T) { // AC 9
	tr, _ := Langkah(PosisiDeptHead, "1", true)
	if tr.StatusTutup != StatusSelesai || !tr.Simpan || tr.PosisiBaru != "" {
		t.Fatalf("Dept Head menyetujui bernomor -> %+v, harap selesai; tidak ada jenjang keempat", tr)
	}
	tr, _ = Langkah(PosisiDeptHead, "1", false)
	if tr.PosisiBaru != PosisiDeptHead {
		t.Fatalf("Dept Head tanpa nomor polis -> %+v, harap tetap di Dept Head (NopolisEmpty)", tr)
	}
}

func TestPosisiBuanganTidakAda(t *testing.T) { // AC 10
	for _, p := range []string{"ReasTreatyInGroupLeader", "ReasTreatyInDirector"} {
		if AdalahPosisiTangga(p) {
			t.Errorf("%s harus tidak ada di sistem baru", p)
		}
		if _, err := Langkah(p, "1", true); err == nil {
			t.Errorf("berkas di %s tidak boleh punya langkah", p)
		}
	}
	for _, pos := range PosisiTangga {
		for _, ia := range []string{"0", "1", ""} {
			for _, n := range []bool{true, false} {
				tr, _ := Langkah(pos, ia, n)
				if tr.PosisiBaru != "" && !AdalahPosisiTangga(tr.PosisiBaru) {
					t.Fatalf("berkas dirutekan ke %q", tr.PosisiBaru)
				}
			}
		}
	}
}

func TestTeksNBStatus(t *testing.T) { // AC 44: nama dari data, bukan tertanam
	if got := TeksNBStatusDitolakOleh("Uji Pengguna"); got != "NB WAS DECLINED BY  UJI PENGGUNA" {
		t.Fatalf("dapat %q", got)
	}
	if got := TeksNBStatusKotakMasuk("uji-satu"); got != "NB IS IN UJI-SATU'S INBOX" {
		t.Fatalf("dapat %q", got)
	}
}

func TestPraprosesAdminTanggalKosongHariIni(t *testing.T) { // AC 34, 35, 42
	h := HalamanBaru()
	h.Setel("Quotation.ProportionalType", JenisProporsional)
	h.Setel("Quotation.MarketingName", "UJI-MO")
	sekarang := time.Date(2026, 10, 3, 9, 30, 0, 0, time.UTC)
	PraprosesAdmin(h, sekarang, "Uji Pengguna")
	if h.Ambil("PolicyTreatyIn.StartDate") != "2026-10-03" {
		t.Errorf("StartDate = %q", h.Ambil("PolicyTreatyIn.StartDate"))
	}
	// P35: tanggal akhir kosong = HARI INI, bukan +1 tahun
	if h.Ambil("PolicyTreatyIn.EndDate") != "2026-10-03" {
		t.Errorf("EndDate = %q, harap hari ini", h.Ambil("PolicyTreatyIn.EndDate"))
	}
	if h.Ambil("PolicyTreatyIn.StatementDate") == "" {
		t.Error("StatementDate kosong harus diisi")
	}
	if h.Ambil("PolicyTreatyIn.OperatorName") != "Uji Pengguna" {
		t.Error("OperatorName harus nama tampilan")
	}
	if h.Ambil("PolicyTreatyIn.IsNewPolicyNonProp") != "0" {
		t.Error("proporsional: IsNewPolicyNonProp harus 0")
	}
	if h.Ambil("PolicyTreatyIn.QuotationData.MarketingName") != "UJI-MO" {
		t.Error("QuotationData harus salinan Quotation")
	}
}

func TestPraprosesAdminTanggalTerisiTidakDitimpa(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.StartDate", "2026-01-01")
	h.Setel("PolicyTreatyIn.EndDate", "2026-12-31")
	PraprosesAdmin(h, time.Now(), "x")
	if h.Ambil("PolicyTreatyIn.StartDate") != "2026-01-01" || h.Ambil("PolicyTreatyIn.EndDate") != "2026-12-31" {
		t.Fatal("tanggal yang sudah terisi tidak boleh ditimpa")
	}
}

func TestNonProporsionalDitandai(t *testing.T) { // AC 75
	h := HalamanBaru()
	h.Setel("Quotation.ProportionalType", JenisNonProporsional)
	PraprosesAdmin(h, time.Now(), "x")
	if h.Ambil("PolicyTreatyIn.IsNewPolicyNonProp") != "1" {
		t.Fatal("non-proporsional: IsNewPolicyNonProp harus 1")
	}
}

func TestPraprosesAtasanMengosongkanPenanda(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.IsApproved", "1")
	h.Setel("PolicyTreatyIn.Suggest", "lama")
	PraprosesAtasan(h, time.Now(), "Uji Atasan")
	if h.Ambil("PolicyTreatyIn.IsApproved") != "" || h.Ambil("PolicyTreatyIn.Suggest") != "" {
		t.Fatal("pra-proses atasan mengosongkan IsApproved dan Suggest")
	}
	if h.Ambil("PolicyTreatyIn.OperatorName") != "Uji Atasan" { // P33 - bukan pengenal akun
		t.Fatal("OperatorName di tahap atasan harus nama tampilan")
	}
}

func TestHasFacOutHanyaDuaKode(t *testing.T) { // AC 76
	h := HalamanBaru()
	h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "10001"}})
	TetapkanHasFacOut(h)
	if h.Ambil("PolicyTreatyIn.HasFacOut") != "0" {
		t.Fatal("kode lain tidak boleh menandai")
	}
	for _, kode := range []string{"10015", "10218"} {
		h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "10001"}, {"TreatyType": kode}})
		TetapkanHasFacOut(h)
		if h.Ambil("PolicyTreatyIn.HasFacOut") != "1" {
			t.Errorf("kode %s harus menandai", kode)
		}
	}
}

func TestTipeNomorPolis(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.DueTo", "1")
	if TipeNomorPolis(h) != "QR" {
		t.Fatal("DueTo 1 -> QR")
	}
	h.Setel("PolicyTreatyIn.DueTo", "0")
	if TipeNomorPolis(h) != "QP" {
		t.Fatal("DueTo 0 -> QP")
	}
	h.Setel("PolicyTreatyIn.ClaimType", "XOL Retro")
	if TipeNomorPolis(h) != "TP" {
		t.Fatal("XOL Retro -> TP menimpa")
	}
}

func TestRakitNomorPolis(t *testing.T) { // AC 73
	got := RakitNomorPolis("UJI", "QR", "005", "10.2026", 42)
	if got != "UJIQR.T005.10.2026.00042" {
		t.Fatalf("dapat %q", got)
	}
}

func TestTanggalProduksiNomor(t *testing.T) {
	jkt := zonaJakarta()
	sekarang := time.Date(2026, 10, 3, 10, 0, 0, 0, jkt)
	// hari 3 <= closing 25: tetap sekarang
	if got := TanggalProduksiNomor(sekarang, sekarang, 25); !got.Equal(sekarang) {
		t.Fatalf("dapat %v", got)
	}
	// hari 27 > closing 25: tanggal 1 bulan berikut 05:00 GMT
	akhir := time.Date(2026, 10, 27, 10, 0, 0, 0, jkt)
	if got := TanggalProduksiNomor(akhir, akhir, 25); !got.Equal(time.Date(2026, 11, 1, 5, 0, 0, 0, time.UTC)) {
		t.Fatalf("dapat %v, harap 2026-11-01 05:00 GMT", got)
	}
	// statement dua hari di depan: dipakai
	depan := sekarang.Add(48 * time.Hour)
	if got := TanggalProduksiNomor(sekarang, depan, 25); !got.Equal(depan) {
		t.Fatalf("statement di depan: dapat %v", got)
	}
	// statement 3 jam di depan: @toInt memotong ke nol - tidak dipakai
	if got := TanggalProduksiNomor(sekarang, sekarang.Add(3*time.Hour), 25); !got.Equal(sekarang) {
		t.Fatalf("statement 3 jam di depan: dapat %v", got)
	}
}
