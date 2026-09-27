package models

import "testing"

// Uji gerbang `Close Claim` — ProtectCloseClaim_act, dibaca sebagai pohon.
//
//	Activity/ProtectCloseClaim_act.xml
//	  b236 Property-Set          ProtectLife.CARI1 = ""
//	  b370 Property-Set  halaman pyWorkPage...PremiumListDetail
//	       b812 ULANG(EMBEDDED) atas seluruh peserta
//	       b398 local.idx  = .pxListSubscript
//	       b445 local.name = .NAME_OF_INSURED
//	    b484 Property-Set   prasyarat b608 `.STS_REJECT!=1`
//	                        WhenTrue=2 lanjut / WhenFalse=3 lewati
//	         b511 ProtectLife.CARI1 = 1
//	         b557 local.errmsg = <nama> + " is not approved yet, on list " + <idx>
//	    b646 Page-Set-Messages  prasyarat b756 `ProtectLife.CARI1==1`
//	  b836 Call FinishAssignment prasyarat b992 `ProtectLife.CARI1==""`
//
// Bacaannya: penugasan hanya DISELESAIKAN bila tidak satu pun peserta
// ditandai. `STS_REJECT == 1` berarti DIAKSEP - `models.KodeAksep` = "1" -
// jadi namanya menyesatkan tetapi maknanya jelas dari kodenya.

func TestBolehTutupKlaim(t *testing.T) {
	t.Run("seluruh peserta diaksep: boleh", func(t *testing.T) {
		baris := []BarisTutup{
			{Urutan: 1, NomorSertifikat: "UJI-001", KodeStatus: KodeAksep},
			{Urutan: 2, NomorSertifikat: "UJI-002", KodeStatus: KodeAksep},
		}
		if !BolehTutupKlaim(baris) {
			t.Error("seluruhnya diaksep tetapi ditolak")
		}
		if got := PenghalangTutupKlaim(baris); len(got) != 0 {
			t.Errorf("penghalang = %v, mau kosong", got)
		}
	})

	t.Run("satu peserta belum diaksep: tertahan", func(t *testing.T) {
		baris := []BarisTutup{
			{Urutan: 1, NomorSertifikat: "UJI-001", KodeStatus: KodeAksep},
			{Urutan: 2, NomorSertifikat: "UJI-002", KodeStatus: KodeOutstanding},
		}
		if BolehTutupKlaim(baris) {
			t.Error("satu peserta belum diaksep tetapi tutup diizinkan")
		}
		p := PenghalangTutupKlaim(baris)
		if len(p) != 1 || p[0].Urutan != 2 || p[0].NomorSertifikat != "UJI-002" {
			t.Fatalf("penghalang = %v, mau baris 2 saja", p)
		}
	})

	t.Run("peserta DITOLAK pun menahan, bukan hanya yang outstanding", func(t *testing.T) {
		// ⛔ Prasyaratnya `.STS_REJECT != 1`, bukan `== 0`. Status DITOLAK
		// ("2") juga bukan 1, jadi ia menahan pula. Memperlakukan "ditolak"
		// sebagai "selesai" akan menutup klaim yang barisnya belum diputus
		// ulang - dan XML tidak pernah mengatakannya.
		baris := []BarisTutup{{Urutan: 1, NomorSertifikat: "UJI-001", KodeStatus: KodeDitolak}}
		if BolehTutupKlaim(baris) {
			t.Error("peserta ditolak seharusnya tetap menahan tutup")
		}
	})

	t.Run("status kosong menahan", func(t *testing.T) {
		baris := []BarisTutup{{Urutan: 1, NomorSertifikat: "UJI-001", KodeStatus: ""}}
		if BolehTutupKlaim(baris) {
			t.Error("status kosong seharusnya menahan")
		}
	})

	t.Run("SELURUH penghalang dilaporkan, bukan yang pertama saja", func(t *testing.T) {
		// Pega memasang pesan di dalam loop, sekali per iterasi yang
		// tertandai. Melaporkan satu saja memaksa pemakai menutup klaim
		// berulang kali, menemukan satu penghalang baru tiap kali.
		baris := []BarisTutup{
			{Urutan: 1, NomorSertifikat: "UJI-001", KodeStatus: ""},
			{Urutan: 2, NomorSertifikat: "UJI-002", KodeStatus: KodeAksep},
			{Urutan: 3, NomorSertifikat: "UJI-003", KodeStatus: KodeDitolak},
		}
		p := PenghalangTutupKlaim(baris)
		if len(p) != 2 {
			t.Fatalf("penghalang = %d, mau 2", len(p))
		}
		if p[0].Urutan != 1 || p[1].Urutan != 3 {
			t.Errorf("urutan penghalang = %v, mau baris 1 dan 3", p)
		}
	})

	t.Run("klaim tanpa peserta: boleh, sama seperti Pega", func(t *testing.T) {
		// ⚠️ Ditiru apa adanya dan DINYATAKAN: loopnya tidak pernah berjalan,
		// `ProtectLife.CARI1` tetap kosong, dan `Call FinishAssignment`
		// berprasyarat kosong - jadi Pega mengizinkannya. Kita tidak
		// menambahkan larangan yang XML tidak punya.
		if !BolehTutupKlaim(nil) {
			t.Error("klaim tanpa peserta ditolak; Pega mengizinkannya")
		}
	})
}

func TestPesanPenghalangTutup(t *testing.T) {
	// ⛔ PENYIMPANGAN SADAR, dan sebabnya bukan gaya.
	//
	// Pega menyusun pesannya dari `.NAME_OF_INSURED` (b445, b557). Kita
	// SENGAJA tidak menyimpan nama tertanggung: `kolomSalin` tiket 02
	// meninggalkan NAME_OF_INSURED dan POLICY_HOLDER justru supaya data
	// pribadi tidak berganda ke tabel klaim. Nama itu karena itu TIDAK ADA
	// untuk dipakai saat menutup.
	//
	// Yang menggantikannya adalah nomor sertifikat - penunjuk baris yang
	// sama persis, yang memang kita simpan, dan yang tidak memuat nama
	// siapa pun. Kalimat sisanya VERBATIM.
	p := Penghalang{Urutan: 7, NomorSertifikat: "UJI-042"}
	mau := "UJI-042 is not approved yet, on list 7"
	if got := p.Pesan(); got != mau {
		t.Errorf("Pesan() = %q, mau %q", got, mau)
	}
}

func TestPesanPenghalangTidakMemuatNama(t *testing.T) {
	// Penjaga atas keputusan di atas: bila kelak seseorang menambahkan medan
	// nama ke Penghalang, uji ini yang gagal lebih dulu dan menagih
	// keputusannya - bukan tabel klaim yang diam-diam berisi nama orang.
	p := Penghalang{Urutan: 1, NomorSertifikat: "UJI-001"}
	if p.Pesan() == "" {
		t.Fatal("pesan kosong")
	}
	// Bentuknya hanya menerima dua hal, dan keduanya bukan nama.
	if got := (Penghalang{Urutan: 2, NomorSertifikat: "UJI-002"}).Pesan(); got !=
		"UJI-002 is not approved yet, on list 2" {
		t.Errorf("bentuk pesan berubah: %q", got)
	}
}

// ⛔ Kode status dibandingkan PERSIS, tanpa merapikan spasi.
//
// Ronde pertama memakai strings.TrimSpace "untuk aman", dan itu justru
// MELONGGARKAN gerbang uang: " 1 " akan menutup klaim di Go padahal Pega -
// yang membandingkan `.STS_REJECT != 1` apa adanya (b608) - menahannya.
// Gerbang yang lebih longgar daripada aslinya menutup klaim yang di sistem
// lama tidak akan pernah tertutup.
func TestKodeStatusDibandingkanPersis(t *testing.T) {
	for _, mirip := range []string{" 1", "1 ", " 1 ", "01", "1.0", "\t1"} {
		baris := []BarisTutup{{Urutan: 1, NomorSertifikat: "UJI-001", KodeStatus: mirip}}
		if BolehTutupKlaim(baris) {
			t.Errorf("kode %q diterima sebagai diaksep; hanya %q yang berarti diaksep",
				mirip, KodeAksep)
		}
	}
	// Dan yang persis tetap diterima.
	baris := []BarisTutup{{Urutan: 1, NomorSertifikat: "UJI-001", KodeStatus: KodeAksep}}
	if !BolehTutupKlaim(baris) {
		t.Error("kode yang persis ditolak")
	}
}
