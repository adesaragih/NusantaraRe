package models

import "testing"

// Uji dua validasi tanggal produk - ValidasiClaimReceived_Act dan
// ValidasiSTNC_Act, dibaca sebagai pohon langkah utuh.
//
//	Activity/ValidasiClaimReceived_Act.xml
//	  b246 RDB-List ProductNameInward  -> RDBList/GetProductName.xml:84
//	  b451 local.errmsg = "Max Claim invalid"
//	  b498 local.DOL               = @addCalendar(.DATE_OF_LOSS, 0,0,0,0,0,0,0)
//	  b519 local.ClaimReceivedDate = @addCalendar(.CLAIM_RECEIVED_DATE, 0,...)
//	  b540 .MAXEXPIREDCLAIM = ProductNameInward.pxResults(1).MAXEXPIREDCLAIM
//	  b561 TempDetail.CARI2 = @DateTimeDifference(local.DOL, local.ClaimReceivedDate, "D")
//	  b582 .MAXCLAIM_RECEIVED = @if(TempDetail.CARI2 <= .MAXEXPIREDCLAIM, "",
//	                                @FormatDateTime(local.ClaimReceivedDate,
//	                                                "dd/MM/YYYY","Asia/Jakarta",""))
//
//	Activity/ValidasiSTNC_Act.xml - bentuk yang SAMA
//	  b479 local.errmsg = "STNC invalid"
//	  b577 .MAXDATARECEIVED = ProductNameInward.pxResults(1).MAXDATARECEIVE
//	  b598 Local.DateDif = @DateTimeDifference(local.EffectiveDate, local.ReceivedDate, "D")
//	  b627 .STNC = @if(Local.DateDif <= .MAXDATARECEIVED, "", @FormatDateTime(...))

func TestSelisihHari(t *testing.T) {
	kasus := []struct {
		nama, dari, ke string
		mau            int
		sah            bool
	}{
		{"hari yang sama nol", "2026-03-01", "2026-03-01", 0, true},
		{"maju sepuluh hari", "2026-03-01", "2026-03-11", 10, true},
		{"melewati batas bulan", "2026-02-25", "2026-03-02", 5, true},
		{"tahun kabisat 2028", "2028-02-28", "2028-03-01", 2, true},
		{"tahun biasa 2026", "2026-02-28", "2026-03-01", 1, true},
		// ⛔ Mundur bernilai NEGATIF, tidak dijadikan nilai mutlak. Lihat
		// TestLubangSelisihNegatif: itu lubang warisan yang ditiru, bukan
		// ditambal diam-diam.
		{"mundur bernilai negatif", "2026-03-11", "2026-03-01", -10, true},
		{"tanggal berjam tetap dihitung per hari", "2026-03-01 23:00:00", "2026-03-02 01:00:00", 1, true},
		{"tanggal kosong tidak sah", "", "2026-03-01", 0, false},
		{"tanggal ngawur tidak sah", "bukan tanggal", "2026-03-01", 0, false},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			got, err := SelisihHari(k.dari, k.ke)
			if (err == nil) != k.sah {
				t.Fatalf("SelisihHari(%q,%q) galat=%v, sah=%v", k.dari, k.ke, err, k.sah)
			}
			if k.sah && got != k.mau {
				t.Errorf("SelisihHari(%q,%q) = %d, mau %d", k.dari, k.ke, got, k.mau)
			}
		})
	}
}

func TestPenandaBatasHari(t *testing.T) {
	// Bentuk penanda: KOSONG berarti sah; terisi berarti TIDAK sah dan isinya
	// adalah tanggal yang melanggar, dd/MM/yyyy.
	kasus := []struct {
		nama, dari, ke, batas string
		mau                   string
	}{
		{"di bawah batas: kosong", "2026-03-01", "2026-03-11", "30", ""},
		{"tepat di batas: kosong", "2026-03-01", "2026-03-31", "30", ""},
		{"lewat satu hari: tanggalnya", "2026-03-01", "2026-04-01", "30", "01/04/2026"},
		{"batas nol, hari sama: kosong", "2026-03-01", "2026-03-01", "0", ""},
		{"batas nol, besoknya: tanggalnya", "2026-03-01", "2026-03-02", "0", "02/03/2026"},
		{"hari satu digit tetap dua digit", "2026-01-01", "2026-01-09", "3", "09/01/2026"},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			got, err := PenandaBatasHari(k.dari, k.ke, k.batas)
			if err != nil {
				t.Fatalf("galat: %v", err)
			}
			if got != k.mau {
				t.Errorf("PenandaBatasHari(%q,%q,%q) = %q, mau %q",
					k.dari, k.ke, k.batas, got, k.mau)
			}
		})
	}
}

func TestPenandaBatasHariBatasKosong(t *testing.T) {
	// ⛔ Ambang yang KOSONG bukan ambang nol. Di Pega ia datang dari
	// PRODUCTINWARD_LIFE lewat PolicyDataLife.ProductNameID; selama modul
	// PremiumList Life belum ada, ambangnya memang belum diketahui.
	//
	// Menganggapnya nol akan membuat SETIAP klaim yang diterima sehari
	// sesudah kejadian ditandai tidak sah - larangan massal yang lahir dari
	// data yang belum ada. Karena itu ambang kosong menjawab GALAT, dan
	// pemanggilnya yang memutuskan (ADR-U-0027: kosong bukan nol).
	for _, batas := range []string{"", "   "} {
		if _, err := PenandaBatasHari("2026-03-01", "2026-06-01", batas); err == nil {
			t.Errorf("ambang %q diterima; ia harus menolak, bukan menganggapnya nol", batas)
		}
	}
}

func TestPenandaBatasHariAmbangBukanAngka(t *testing.T) {
	if _, err := PenandaBatasHari("2026-03-01", "2026-03-05", "tiga puluh"); err == nil {
		t.Error("ambang bukan angka diterima")
	}
}

func TestLubangSelisihNegatif(t *testing.T) {
	// ⚠️ LUBANG WARISAN, DITIRU DENGAN SENGAJA dan dilaporkan (OQ-G).
	//
	// XML membandingkan `selisih <= batas` tanpa lantai bawah. Klaim yang
	// diterima SEBELUM tanggal kejadiannya menghasilkan selisih negatif, dan
	// negatif selalu <= batas - sehingga ia lolos sebagai sah.
	//
	// Test ini ada supaya perilaku itu TERLIHAT dan tidak berubah diam-diam.
	// Bila kelak work owner memutuskan menambal lubangnya, test inilah yang
	// gagal lebih dulu dan menagih keputusannya.
	got, err := PenandaBatasHari("2026-06-01", "2026-03-01", "30")
	if err != nil {
		t.Fatalf("galat: %v", err)
	}
	if got != "" {
		t.Errorf("penanda = %q; XML meloloskan selisih negatif, dan itu yang ditiru", got)
	}
}

func TestPesanValidasiVerbatim(t *testing.T) {
	// ⛔ VERBATIM dari rule-nya, bukan terjemahan. Pemakai dan pengembang
	// Pega harus dapat mencari kalimat yang sama di kedua sistem.
	if PesanMaxClaimTidakSah != "Max Claim invalid" {
		t.Errorf("pesan = %q, mau b451 apa adanya", PesanMaxClaimTidakSah)
	}
	if PesanSTNCTidakSah != "STNC invalid" {
		t.Errorf("pesan = %q, mau b479 apa adanya", PesanSTNCTidakSah)
	}
}
