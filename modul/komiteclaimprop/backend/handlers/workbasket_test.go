package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nusantarare/modul/komiteclaimprop/backend/handlers"
	"nusantarare/modul/komiteclaimprop/backend/models"
	"nusantarare/modul/komiteclaimprop/backend/services"
	"nusantarare/modul/komiteclaimprop/backend/tiruan"
)

// Keputusan work owner 09-10-2026 (tangga komite PROP ke workbasket, migrasi claimprop 537): penyetuju = anggota
// workbasket tingkat berjalan; baris yang diputus menyimpan akun pemutus; T_WORK_CLAIM.POSITION = workbasket tingkat
// berjalan; email tingkat berikut ke semua anggota workbasket. TANPA larangan rangkap (WO: "1 akun memang tidak boleh
// memiliki 2 jabatan dalam komite" - dijaga pengaturan akun; SUPERADMIN memegang semua untuk pengujian).
func TestTanggaWorkbasketPemegangPosisiEmail(t *testing.T) {
	u := siap(t, 2, false)
	tangga := u.g.Tangga[tiruan.KomiteUji]
	tangga[0].OperatorID, tangga[1].OperatorID = "UJI-WB-1", "UJI-WB-2"
	u.a.Nama["UJI-A"] = "UJI Anggota A"
	u.a.AnggotaWB = map[string][]string{"UJI-WB-2": {"uji.b@contoh.invalid", "uji.c@contoh.invalid"}}

	minta := func(metode, jalur, pelaku, peran string, badan any) *httptest.ResponseRecorder {
		t.Helper()
		var b []byte
		if badan != nil {
			b, _ = json.Marshal(badan)
		}
		r := httptest.NewRequest(metode, handlers.Prefix+jalur, bytes.NewReader(b))
		r.Header.Set("X-Pelaku", pelaku)
		r.Header.Set("X-Peran", peran)
		w := httptest.NewRecorder()
		u.srv.ServeHTTP(w, r)
		return w
	}
	kerja := func(pelaku, peran string) []models.BarisKerja {
		t.Helper()
		w := minta("GET", "/kasus", pelaku, peran, nil)
		var out []models.BarisKerja
		if w.Code != http.StatusOK || json.Unmarshal(w.Body.Bytes(), &out) != nil {
			t.Fatalf("daftar kerja %s: %d %s", pelaku, w.Code, w.Body.String())
		}
		return out
	}
	const keduanya = "UJI-WB-1,UJI-WB-2"

	if d := kerja("UJI-A", keduanya); len(d) != 1 || d[0].Tingkat != 1 {
		t.Fatalf("anggota workbasket tingkat 1: %+v", d)
	}
	if d := kerja("UJI-B", "UJI-WB-2"); len(d) != 0 {
		t.Fatalf("anggota workbasket tingkat 2 belum giliran: %+v", d)
	}
	if w := minta("POST", "/kasus/"+tiruan.KomiteUji+"/putuskan", "UJI-A", keduanya, setuju("UJI 1")); w.Code != http.StatusOK {
		t.Fatalf("anggota UJI-WB-1 memutus tingkat 1: %d %s", w.Code, w.Body.String())
	}
	if op := u.g.Tangga[tiruan.KomiteUji][0].OperatorID; op != "UJI-A" {
		t.Fatalf("baris tingkat 1 menyimpan akun pemutus: %q", op)
	}
	if p := u.g.Posisi[tiruan.KomiteUji]; p != "UJI-WB-2" {
		t.Fatalf("POSITION sesudah tingkat 1 = %q, mau UJI-WB-2", p)
	}

	surel, err := u.l.SusunEmailKomite(context.Background(), tiruan.KomiteUji, map[string]string{
		services.IsiJenis: models.EmailPenyetujuBerikut, services.IsiAnggota: u.g.Tangga[tiruan.KomiteUji][0].ID,
		services.IsiPenerima: "UJI-WB-2"})
	if err != nil {
		t.Fatal(err)
	}
	if surel.Kepada != "uji.b@contoh.invalid, uji.c@contoh.invalid" {
		t.Fatalf("email tingkat berikut ke semua anggota workbasket: %q", surel.Kepada)
	}

	if d := kerja("UJI-B", "UJI-WB-2"); len(d) != 1 || d[0].Tingkat != 2 {
		t.Fatalf("anggota lain workbasket tingkat 2: %+v", d)
	}
	// akun yang memegang dua workbasket (pengujian) tetap melihat dan memutus tingkat berikutnya
	if d := kerja("UJI-A", keduanya); len(d) != 1 || d[0].Tingkat != 2 {
		t.Fatalf("UJI-A pemegang UJI-WB-2 juga: tingkat 2 tampil: %+v", d)
	}
	if w := minta("POST", "/kasus/"+tiruan.KomiteUji+"/putuskan", "UJI-A", keduanya, setuju("UJI 2")); w.Code != http.StatusOK {
		t.Fatalf("UJI-A memutus tingkat 2: %d %s", w.Code, w.Body.String())
	}
	if p, ada := u.g.Posisi[tiruan.KomiteUji]; !ada || p != "" {
		t.Fatalf("kasus selesai: POSITION dikosongkan, dapat %q", p)
	}
}
