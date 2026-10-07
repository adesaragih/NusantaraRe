package handlers_test

// Uji seam 1 - NBStatus saat Submit (`[keputusan work owner 06-10-2026]`): pemegang AKTIF workbasket tujuan
// (M_LOGIN_GO_WORKBASKET x M_LOGIN_GO) satu orang -> "NB IS IN <NAMA AKUN>'S INBOX"; lebih dari satu ->
// "NB IS IN <M_WORKBASKET.NAME>'S INBOX". Huruf besar (`@toUpperCase`). Fixture UJI-.

import (
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestNBStatusKotakMasukDariPemegangWorkbasket(t *testing.T) {
	u := baru(t)
	u.g.KotakMasuk = map[string]models.PemegangKotakMasuk{
		models.PosisiSecHead:  {Jumlah: 1, NamaAkun: "Uji Sec Satu", NamaWorkbasket: "UJI Treaty Inward - Section Head"},
		models.PosisiDeptHead: {Jumlah: 3, NamaAkun: "Uji Dept", NamaWorkbasket: "UJI Treaty Inward - Department Head"},
	}
	id := u.buat()
	if kode, isi := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin submit: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("NBStatus"); got != "NB IS IN UJI SEC SATU'S INBOX" {
		t.Errorf("satu pemegang Sec Head: NBStatus %q", got)
	}
	if kode, isi := u.kirim(id, secHead, putusan("1")); kode != http.StatusOK {
		t.Fatalf("Sec Head submit: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("NBStatus"); got != "NB IS IN UJI TREATY INWARD - DEPARTMENT HEAD'S INBOX" {
		t.Errorf("tiga pemegang Dept Head: NBStatus %q", got)
	}
}

// NBStatus tidak pernah kosong (keputusan work owner 06-10-2026 "petunjuk ke user NB-nya ada di siapa"): atasan
// menolak -> berkas kembali ke PEMBUATNYA (grid admin hanya menampilkan buatan sendiri) -> "NB IS IN <PEMBUAT>'S
// INBOX". XML: connector Decision11/4/2 No menulis NBStatus "".
func TestAtasanMenolakNBStatusMenunjukPembuat(t *testing.T) {
	for _, tolak := range []pelakuUji{secHead, deptHead} {
		u := baru(t)
		u.g.Nama[admin.akun] = "Uji Pembuat"
		id := u.buat()
		if kode, isi := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusOK {
			t.Fatalf("admin submit: %d %s", kode, isi)
		}
		if tolak == deptHead {
			if kode, isi := u.kirim(id, secHead, putusan("1")); kode != http.StatusOK {
				t.Fatalf("Sec Head submit: %d %s", kode, isi)
			}
		}
		if kode, isi := u.kirim(id, tolak, putusan("0")); kode != http.StatusOK {
			t.Fatalf("%s menolak: %d %s", tolak.akun, kode, isi)
		}
		if got := u.g.Halaman[id].Ambil("NBStatus"); u.g.Kasus[id].PositionNote != models.PosisiAdmin || got != "NB IS IN UJI PEMBUAT'S INBOX" {
			t.Errorf("%s menolak: posisi %q NBStatus %q", tolak.akun, u.g.Kasus[id].PositionNote, got)
		}
	}
}

// Berkas lama tanpa pembuat (pemuat Pega: CREATE_OP kosong) ditolak atasan -> aturan pemegang workbasket Admin.
func TestAtasanMenolakBerkasTanpaPembuat(t *testing.T) {
	u := baru(t)
	u.g.KotakMasuk = map[string]models.PemegangKotakMasuk{models.PosisiAdmin: {Jumlah: 2, NamaWorkbasket: "UJI Treaty Inward - Admin"}}
	id := u.buat()
	if kode, isi := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin submit: %d %s", kode, isi)
	}
	k := u.g.Kasus[id]
	k.CreateOp = ""
	u.g.Kasus[id] = k
	if kode, isi := u.kirim(id, secHead, putusan("0")); kode != http.StatusOK {
		t.Fatalf("Sec Head menolak: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("NBStatus"); got != "NB IS IN UJI TREATY INWARD - ADMIN'S INBOX" {
		t.Errorf("tanpa pembuat: NBStatus %q", got)
	}
}
