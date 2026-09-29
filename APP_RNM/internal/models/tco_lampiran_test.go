package models

import "testing"

// Status lampiran diturunkan dari DUA sumber: kunci penyimpanan yang sudah
// dipastikan, dan nasib efek outbox-nya. Terkirim menang atas apa pun yang
// outbox katakan - berkas yang sudah ada di penyimpanan tidak "gagal".
func TestStatusLampiranTCO(t *testing.T) {
	kasus := []struct {
		storage  string
		menyerah bool
		mau      string
	}{
		{"", false, StatusLampiranTertunda},
		{"", true, StatusLampiranGagal},
		{"ABC", false, StatusLampiranTerkirim},
		{"ABC", true, StatusLampiranTerkirim},
		{"   ", false, StatusLampiranTertunda},
	}
	for _, k := range kasus {
		if dapat := StatusLampiranTCO(k.storage, k.menyerah); dapat != k.mau {
			t.Errorf("StatusLampiranTCO(%q, %v) = %q, mau %q", k.storage, k.menyerah, dapat, k.mau)
		}
	}
}

// Kategori disimpan VERBATIM seperti master menuliskannya, meski pemakai
// mengetik dengan huruf dan spasi berbeda.
func TestKategoriLampiranSah(t *testing.T) {
	master := []string{"CLAUSES", "R/I SLIP", "OTHERS"}
	if k, ok := KategoriLampiranSah(master, "  r/i slip "); !ok || k != "R/I SLIP" {
		t.Errorf("r/i slip -> %q %v", k, ok)
	}
	for _, salah := range []string{"", "   ", "SLIP", "R/I"} {
		if _, ok := KategoriLampiranSah(master, salah); ok {
			t.Errorf("%q lolos padahal bukan kategori master", salah)
		}
	}
}

// Nama berkas antrean hanya berisi kunci + akhiran yang disaring: nama
// unggahan tidak pernah menentukan jalur di disk.
func TestNamaBerkasAntreLampiranTCO(t *testing.T) {
	kasus := map[[2]string]string{
		{"ABC123", "kontrak.PDF"}:              "ABC123.pdf",
		{"ABC123", "tanpa-akhiran"}:            "ABC123",
		{"ABC123", `..\..\rahasia.exe`}:        "ABC123.exe",
		{"ABC123", "a.b/../../x.txt"}:          "ABC123.txt",
		{"ABC123", "aneh.p%d$f"}:               "ABC123",
		{"ABC123", "panjang.abcdefghijklmnop"}: "ABC123",
	}
	for masuk, mau := range kasus {
		if dapat := NamaBerkasAntreLampiranTCO(masuk[0], masuk[1]); dapat != mau {
			t.Errorf("NamaBerkasAntreLampiranTCO(%q, %q) = %q, mau %q", masuk[0], masuk[1], dapat, mau)
		}
	}
}

// Entri zip "Download All" berawalan ID lampiran (dua berkas senama tidak
// saling timpa) dan tidak pernah memuat jalur.
func TestNamaEntriZipLampiranTCO(t *testing.T) {
	kasus := map[[2]string]string{
		{"1000000001", "kontrak.pdf"}:      "1000000001_kontrak.pdf",
		{"1000000002", `C:\tmp\slip.xlsx`}: "1000000002_slip.xlsx",
		{"1000000003", "../../etc/passwd"}: "1000000003_passwd",
		{"1000000004", ""}:                 "1000000004_berkas",
		{"1000000005", ".."}:               "1000000005_berkas",
	}
	for masuk, mau := range kasus {
		if dapat := NamaEntriZipLampiranTCO(masuk[0], masuk[1]); dapat != mau {
			t.Errorf("NamaEntriZipLampiranTCO(%q, %q) = %q, mau %q", masuk[0], masuk[1], dapat, mau)
		}
	}
}

// OQ-TCO-26 (lanjutan 4): `exp` jawaban penyimpanan diubah seperti Pega
// (`GetUrlGoogleStorage_Act` b2146/b2211, `InsertGoogleStorage_Act` b2366/b2431):
// buang `-`/`:`, baca sebagai GMT, tulis `dd/MM/yyyy HH:mm:ss` - bentuk To_date SQL-nya.
func TestExpStorageTCOSepertiPega(t *testing.T) {
	for masuk, mau := range map[string]string{
		"2026-09-29T10:05:07.123Z": "29/09/2026 10:05:07",
		"2026-09-29T10:05:07Z":     "29/09/2026 10:05:07",
		"20260929T100507.000 GMT":  "29/09/2026 10:05:07",
		"":                         "",
		"29/09/2026 10:05:07":      "",
		"kapan-kapan":              "",
	} {
		if dapat := ExpStorageTCO(masuk); dapat != mau {
			t.Errorf("%q -> %q, mau %q", masuk, dapat, mau)
		}
	}
}

// Temuan /code-review lanjutan 4: `DateTime` geturl masuk To_date
// `MM/DD/YYYY HH24:MI:SS` - bentuk lain jadi kosong (NULL), bukan ORA-01843
// yang menggagalkan seluruh `Update_T_Storage_SQL`.
func TestTanggalUploadStorageTCOHanyaBentukToDate(t *testing.T) {
	for masuk, mau := range map[string]string{
		"09/29/2026 09:00:00":  "09/29/2026 09:00:00",
		" 09/29/2026 09:00:00": "09/29/2026 09:00:00",
		"2026-09-29T09:00:00Z": "",
		"29/09/2026 09:00:00":  "",
		"09/29/2026":           "",
		"":                     "",
	} {
		if dapat := TanggalUploadStorageTCO(masuk); dapat != mau {
			t.Errorf("%q -> %q, mau %q", masuk, dapat, mau)
		}
	}
}
