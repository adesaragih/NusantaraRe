package models_test

import (
	"reflect"
	"testing"

	"nusantarare/modul/claimprop/backend/models"
)

// tampakDiLayar - padanan `saringKlien` layar: unsur yang syarat `pyIsClientWhen`-nya terpenuhi `nilai`, diratakan;
// hasilnya ID tombol, label, dan jalur medan yang tampak.
func tampakDiLayar(ts []models.Tata, nilai map[string]string) (tombol, label, medan []string) {
	var jalan func([]models.Tata)
	jalan = func(xs []models.Tata) {
		for _, x := range xs {
			if s := x.TampilJika; s != nil && (nilai[s.Jalur] == s.Nilai) == s.Beda {
				continue
			}
			switch x.Jenis {
			case models.JenisTombol:
				tombol = append(tombol, x.ID)
			case models.JenisLabel:
				label = append(label, x.Label)
			case models.JenisMedan:
				medan = append(medan, x.Jalur)
			}
			jalan(x.Anak)
		}
	}
	jalan(ts)
	return
}

// Pop-up PreventRejectClaimProp: kondisi layout ber-`pyIsClientWhen = true` (server mengirim semua baris, layar memilih);
// perintah work owner 10-10-2026: kartu sukses bernomor sesudah Yes tanpa bayar, kartu "already submitted" selama
// pengajuan Reject / Close Without Payment masih menunggu komite, teks bahasa Inggris.
func TestLayarTutupKlaimKondisiDiLayar(t *testing.T) {
	ts := models.Evaluasi(models.HalamanBaru(), models.LayarTutupKlaim(), false)
	cwp, baru, tunggu, jenis := models.PropCentangTutupTanpaBayar, models.PropKomiteBaru, models.PropKomiteMenunggu,
		models.PropJenisKomiteMenunggu
	kasus := []struct {
		nama          string
		nilai         map[string]string
		tombol, label []string
		centang       bool
	}{
		{"awal", map[string]string{}, []string{"CloseYes", "CloseNo2"}, []string{models.TanyaTutup}, true},
		{"dicentang", map[string]string{cwp: "true"}, []string{"CloseWithoutPayment", "CloseNo"},
			[]string{models.TanyaTutupTanpaBayar}, true},
		{"reject menunggu", map[string]string{tunggu: "TKMT-UJI1", jenis: models.TransferTolak}, []string{"MenungguTutup"},
			[]string{models.JudulMenungguTolak, models.TeksMenungguTolak}, false},
		{"close menunggu", map[string]string{tunggu: "TKMT-UJI1", jenis: models.TransferTutup}, []string{"MenungguTutup"},
			[]string{models.JudulMenungguTutup, models.TeksMenungguTutup}, false},
		// sesudah Yes: kartu sukses menang atas penanda menunggu (kasus komite yang baru lahir ikut terbaca)
		{"sukses", map[string]string{cwp: "true", baru: "TKMT-UJI1", tunggu: "TKMT-UJI1", jenis: models.TransferTutup},
			[]string{"CloseSukses"}, []string{models.JudulSuksesTutup, models.TeksSuksesTutup}, false},
	}
	for _, c := range kasus {
		tb, lb, md := tampakDiLayar(ts, c.nilai)
		if !reflect.DeepEqual(tb, c.tombol) || !reflect.DeepEqual(lb, c.label) {
			t.Errorf("%s: tombol %v label %v, mau %v %v", c.nama, tb, lb, c.tombol, c.label)
		}
		ada := false
		for _, m := range md {
			ada = ada || m == cwp
		}
		if ada != c.centang {
			t.Errorf("%s: centang Close Without Payment tampak %v, mau %v", c.nama, ada, c.centang)
		}
	}
}

// Pop-up ClaimRejectTreaty: isian + Submit hanya bila tidak ada pengajuan yang menunggu komite.
func TestLayarTolakKlaimKondisiDiLayar(t *testing.T) {
	ts := models.Evaluasi(models.HalamanBaru(), models.LayarTolakKlaim(), false)
	tb, _, md := tampakDiLayar(ts, map[string]string{})
	if !reflect.DeepEqual(tb, []string{"SubmitTolak", "CancelTolak"}) || len(md) != 5 {
		t.Errorf("awal: tombol %v medan %v", tb, md)
	}
	tb, lb, md := tampakDiLayar(ts, map[string]string{models.PropKomiteMenunggu: "TKMT-UJI1",
		models.PropJenisKomiteMenunggu: models.TransferTolak})
	if !reflect.DeepEqual(tb, []string{"MenungguTolak"}) ||
		!reflect.DeepEqual(lb, []string{models.JudulMenungguTolak, models.TeksMenungguTolak}) ||
		!reflect.DeepEqual(md, []string{models.PropKomiteMenunggu}) {
		t.Errorf("menunggu: tombol %v label %v medan %v", tb, lb, md)
	}
}

func TestPesanSudahDiajukan(t *testing.T) {
	if p := models.PesanSudahDiajukan(models.TransferTolak, "TKMT-000013"); p !=
		"A reject request for this claim has already been submitted (TKMT-000013) and is still waiting for the committee's decision." {
		t.Errorf("reject: %q", p)
	}
	if p := models.PesanSudahDiajukan(models.TransferTutup, "TKMT-000014"); p !=
		"A close request for this claim has already been submitted (TKMT-000014) and is still waiting for the committee's decision." {
		t.Errorf("close: %q", p)
	}
}
