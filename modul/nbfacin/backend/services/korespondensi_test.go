package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// korespondensiTiruan - penyimpan dalam memori dengan aturan draf / terkirim yang sama dengan repository Oracle.
type korespondensiTiruan struct {
	daftar      []models.Korespondensi
	urut        int
	galatTandai error
}

func (k *korespondensiTiruan) DaftarKorespondensi(_ context.Context, id string) ([]models.Korespondensi, error) {
	if id != "UJI-NB-1" {
		return nil, repository.ErrKasusTidakAda
	}
	return append([]models.Korespondensi{}, k.daftar...), nil
}

func (k *korespondensiTiruan) cari(no string) int {
	for i := range k.daftar {
		if k.daftar[i].No == no {
			return i
		}
	}
	return -1
}

func (k *korespondensiTiruan) TambahKorespondensi(_ context.Context, _ *db.Tx, _ string, x models.Korespondensi, _ string) (string, error) {
	k.urut++
	x.No = strings.Repeat("1", k.urut)
	x.Penerima = []models.PenerimaKorespondensi{}
	k.daftar = append(k.daftar, x)
	return x.No, nil
}

func (k *korespondensiTiruan) UbahKorespondensi(_ context.Context, _ *db.Tx, _, no string, isi models.IsiKorespondensi) error {
	i := k.cari(no)
	switch {
	case i < 0:
		return repository.ErrKorespondensiTidakAda
	case k.daftar[i].Terkirim:
		return repository.ErrKorespondensiTerkirim
	}
	k.daftar[i].Subject, k.daftar[i].Message, k.daftar[i].Penerima = isi.Subject, isi.Message, isi.Penerima
	return nil
}

func (k *korespondensiTiruan) HapusKorespondensi(_ context.Context, _ *db.Tx, _, no string) error {
	i := k.cari(no)
	switch {
	case i < 0:
		return repository.ErrKorespondensiTidakAda
	case k.daftar[i].Terkirim:
		return repository.ErrKorespondensiTerkirim
	}
	k.daftar = append(k.daftar[:i], k.daftar[i+1:]...)
	return nil
}

func (k *korespondensiTiruan) TandaiTerkirim(ctx context.Context, _ *db.Tx, _, no, waktu string) error {
	if k.galatTandai != nil {
		return k.galatTandai
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	i := k.cari(no)
	k.daftar[i].Terkirim, k.daftar[i].DikirimPada = true, waktu
	return nil
}

func (k *korespondensiTiruan) AkunLogin(_ context.Context, akun string) (string, string, bool, error) {
	if akun == "UJI-USER" {
		return "UJI PENGIRIM", "uji.pengirim@contoh.test", true, nil
	}
	return "", "", false, nil
}

func (k *korespondensiTiruan) SaranPenerima(_ context.Context, kata string) ([]models.SaranAkumulasi, error) {
	return []models.SaranAkumulasi{{ID: "UJI-A", Label: "UJI " + kata, Ekstra: "uji.a@contoh.test"}}, nil
}

// surelTiruan - pengirim tiruan: TIDAK PERNAH menyambung jaringan.
type surelTiruan struct {
	kirim []models.Surel
	galat error
}

func (s *surelTiruan) KirimSurel(_ context.Context, x models.Surel) error {
	if s.galat != nil {
		return s.galat
	}
	s.kirim = append(s.kirim, x)
	return nil
}

var pelakuKor = inti.Pelaku{AkunID: "UJI-USER"}

func layananKorespondensi(kt *korespondensiTiruan, st *surelTiruan) *Service {
	jam := func() time.Time { return time.Date(2026, 10, 6, 3, 4, 0, 0, time.UTC) } // 10:04 WIB
	svc := Baru(nil).DenganKorespondensi(kt, kt).DenganJabatan(jabatanTiruan{"UJI-USER": "IT", "UJI-UW": "Underwriter"}).
		DenganTransaksi(tanpaTx).DenganJam(jam)
	if st != nil {
		svc.DenganPengirimSurel(st)
	}
	return svc
}

// TestTambahSimpanHapusKorespondensi - Add (tanggal WIB, pengirim akun login, draf kosong), Save (pesan dibersihkan,
// penerima kosong dibuang, subject dipangkas), 404 nomor asing, 401 tanpa sesi, Hapus draf.
func TestTambahSimpanHapusKorespondensi(t *testing.T) {
	ctx := context.Background()
	kt := &korespondensiTiruan{}
	svc := layananKorespondensi(kt, nil)
	k, err := svc.TambahKorespondensi(ctx, pelakuKor, "UJI-NB-1")
	if err != nil || k.Tanggal != "06-10-2026 10:04" || k.Sender != "UJI PENGIRIM" || k.SenderMail != "uji.pengirim@contoh.test" ||
		k.Terkirim || len(k.Penerima) != 0 || k.Subject != "" {
		t.Fatalf("tambah: %v %+v", err, k)
	}
	if _, err := svc.TambahKorespondensi(ctx, inti.Pelaku{}, "UJI-NB-1"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("401: %v", err)
	}
	isi := models.IsiKorespondensi{Subject: "  UJI subjek ", Message: `<p onclick="x">halo</p><img src=x onerror=alert(1)>`,
		Penerima: []models.PenerimaKorespondensi{{Recipient: " UJI A ", RecipientMail: " uji.a@contoh.test "}, {}}}
	k, err = svc.SimpanKorespondensi(ctx, pelakuKor, "UJI-NB-1", k.No, isi)
	if err != nil || k.Subject != "UJI subjek" || k.Message != "<p>halo</p>" ||
		!reflect.DeepEqual(k.Penerima, []models.PenerimaKorespondensi{{Recipient: "UJI A", RecipientMail: "uji.a@contoh.test"}}) {
		t.Fatalf("simpan: %v %+v", err, k)
	}
	for _, no := range []string{"999", "abc", ""} {
		if _, err := svc.SimpanKorespondensi(ctx, pelakuKor, "UJI-NB-1", no, isi); !errors.Is(err, ErrKorespondensiTidakAda) {
			t.Errorf("404 %q: %v", no, err)
		}
	}
	for nama, salah := range map[string]models.IsiKorespondensi{
		"subject CRLF":  {Subject: "a\r\nBcc: x@contoh.test"},
		"subject 1001":  {Subject: strings.Repeat("a", 1001)},
		"pesan > 4000":  {Message: strings.Repeat("a", 4001)},
		"penerima CRLF": {Penerima: []models.PenerimaKorespondensi{{RecipientMail: "a@contoh.test\nBcc: b@contoh.test"}}},
	} {
		if _, err := svc.SimpanKorespondensi(ctx, pelakuKor, "UJI-NB-1", k.No, salah); !errors.Is(err, ErrMasukanKorespondensi) {
			t.Errorf("400 %s: %v", nama, err)
		}
	}
	if err := svc.HapusKorespondensi(ctx, pelakuKor, "UJI-NB-1", k.No); err != nil || len(kt.daftar) != 0 {
		t.Errorf("hapus: %v %d", err, len(kt.daftar))
	}
}

// TestKirimKorespondensi - To = penerima unik (tanpa beda huruf), Cc + Reply-To = pengirim; terkirim terkunci (409 Save /
// Hapus / Kirim ulang); 400 tanpa penerima ber-email / alamat tidak sah; 503 SMTP belum dikonfigurasi; 502 gagal - isi
// TETAP tersimpan dan status draf.
func TestKirimKorespondensi(t *testing.T) {
	ctx := context.Background()
	kt, st := &korespondensiTiruan{}, &surelTiruan{}
	svc := layananKorespondensi(kt, st)
	k, _ := svc.TambahKorespondensi(ctx, pelakuKor, "UJI-NB-1")
	isi := models.IsiKorespondensi{Subject: "UJI subjek", Message: "<b>isi</b>", Penerima: []models.PenerimaKorespondensi{
		{Recipient: "UJI A", RecipientMail: "uji.a@contoh.test"}, {Recipient: "UJI A lagi", RecipientMail: "UJI.A@contoh.test"},
		{Recipient: "UJI B", RecipientMail: "uji.b@contoh.test"}, {Recipient: "UJI tanpa email"}}}
	hasil, err := svc.KirimKorespondensi(ctx, pelakuKor, "UJI-NB-1", k.No, isi)
	mau := models.Surel{To: []string{"uji.a@contoh.test", "uji.b@contoh.test"}, Cc: []string{"uji.pengirim@contoh.test"},
		ReplyTo: "uji.pengirim@contoh.test", Subject: "UJI subjek", HTML: "<b>isi</b>"}
	if err != nil || !hasil.Terkirim || hasil.DikirimPada != "06-10-2026 10:04" || len(st.kirim) != 1 || !reflect.DeepEqual(st.kirim[0], mau) {
		t.Fatalf("kirim: %v %+v %+v", err, hasil, st.kirim)
	}
	if _, err := svc.KirimKorespondensi(ctx, pelakuKor, "UJI-NB-1", k.No, isi); !errors.Is(err, ErrKorespondensiTerkirim) || len(st.kirim) != 1 {
		t.Errorf("kirim ulang 409: %v", err)
	}
	if _, err := svc.SimpanKorespondensi(ctx, pelakuKor, "UJI-NB-1", k.No, isi); !errors.Is(err, ErrKorespondensiTerkirim) {
		t.Errorf("simpan terkirim 409: %v", err)
	}
	if err := svc.HapusKorespondensi(ctx, pelakuKor, "UJI-NB-1", k.No); !errors.Is(err, ErrKorespondensiTerkirim) {
		t.Errorf("hapus terkirim 409: %v", err)
	}
	// 409 lebih dulu dari 400: isian tidak sah atas korespondensi terkirim tetap 409; nomor asing tetap 404.
	salah := models.IsiKorespondensi{Subject: strings.Repeat("a", 1001)}
	if _, err := svc.SimpanKorespondensi(ctx, pelakuKor, "UJI-NB-1", k.No, salah); !errors.Is(err, ErrKorespondensiTerkirim) {
		t.Errorf("409 sebelum 400: %v", err)
	}
	if _, err := svc.SimpanKorespondensi(ctx, pelakuKor, "UJI-NB-1", "999", salah); !errors.Is(err, ErrKorespondensiTidakAda) {
		t.Errorf("404 sebelum 400: %v", err)
	}
	k2, _ := svc.TambahKorespondensi(ctx, pelakuKor, "UJI-NB-1")
	for nama, p := range map[string][]models.PenerimaKorespondensi{
		"tanpa email":  {{Recipient: "UJI tanpa email"}},
		"alamat rusak": {{RecipientMail: "bukan alamat"}},
		"nama tampil":  {{RecipientMail: "UJI <uji.a@contoh.test>"}},
	} {
		_, err := svc.KirimKorespondensi(ctx, pelakuKor, "UJI-NB-1", k2.No, models.IsiKorespondensi{Subject: "UJI " + nama, Penerima: p})
		if !errors.Is(err, ErrMasukanKorespondensi) || len(st.kirim) != 1 {
			t.Errorf("400 %s: %v", nama, err)
		}
		if kt.daftar[1].Subject != "UJI "+nama {
			t.Errorf("isi harus tersimpan walau 400: %+v", kt.daftar[1])
		}
	}
	ok := models.IsiKorespondensi{Subject: "UJI 503", Penerima: []models.PenerimaKorespondensi{{RecipientMail: "uji.a@contoh.test"}}}
	if _, err := layananKorespondensi(kt, nil).KirimKorespondensi(ctx, pelakuKor, "UJI-NB-1", k2.No, ok); !errors.Is(err, ErrSMTPBelumDikonfigurasi) ||
		kt.daftar[1].Subject != "UJI 503" || kt.daftar[1].Terkirim {
		t.Errorf("503: %v %+v", err, kt.daftar[1])
	}
	st.galat = errors.New("UJI server menolak")
	ok.Subject = "UJI 502"
	if _, err := svc.KirimKorespondensi(ctx, pelakuKor, "UJI-NB-1", k2.No, ok); !errors.Is(err, ErrKirimSurelGagal) ||
		kt.daftar[1].Subject != "UJI 502" || kt.daftar[1].Terkirim {
		t.Errorf("502: %v %+v", err, kt.daftar[1])
	}
	// Server tanpa TLS (TLS wajib) -> galat tersendiri 502, tidak terkirim, isi tersimpan, status draf.
	st.galat = fmt.Errorf("uji: %w", repository.ErrSMTPTanpaTLS)
	ok.Subject = "UJI tanpa TLS"
	if _, err := svc.KirimKorespondensi(ctx, pelakuKor, "UJI-NB-1", k2.No, ok); !errors.Is(err, ErrSMTPTanpaTLS) || errors.Is(err, ErrKirimSurelGagal) ||
		kt.daftar[1].Subject != "UJI tanpa TLS" || kt.daftar[1].Terkirim {
		t.Errorf("tanpa TLS: %v %+v", err, kt.daftar[1])
	}
	// Email terkirim tetapi status gagal ditulis -> galat tersendiri (jangan kirim ulang); permintaan dibatalkan sesudah
	// kirim tetap menulis status (context.WithoutCancel).
	st.galat = nil
	kt.galatTandai = errors.New("UJI oracle putus")
	if _, err := svc.KirimKorespondensi(ctx, pelakuKor, "UJI-NB-1", k2.No, ok); !errors.Is(err, ErrKorespondensiStatusGagal) || kt.daftar[1].Terkirim {
		t.Errorf("status gagal: %v", err)
	}
	kt.galatTandai = nil
	batal, henti := context.WithCancel(ctx)
	svc.DenganPengirimSurel(pengirimMembatalkan{henti})
	if _, err := svc.KirimKorespondensi(batal, pelakuKor, "UJI-NB-1", k2.No, ok); !kt.daftar[1].Terkirim {
		t.Errorf("status harus tertulis walau dibatalkan sesudah kirim: %v", err)
	}
	// DenganSMTP: host / dari kosong = tidak dipasang (503); terisi = pengirim SMTP (tidak disambung di uji).
	if svc.DenganSMTP(repository.KonfigSMTP{Dari: "uji@contoh.test"}).surel != nil || svc.DenganSMTP(repository.KonfigSMTP{Host: "uji.host"}).surel != nil ||
		svc.DenganSMTP(repository.KonfigSMTP{Host: "uji.host", Port: 587, Dari: "uji@contoh.test"}).surel == nil {
		t.Error("DenganSMTP")
	}
}

// pengirimMembatalkan - berhasil "mengirim" lalu membatalkan permintaan (klien menutup koneksi sesudah email keluar).
type pengirimMembatalkan struct{ henti context.CancelFunc }

func (p pengirimMembatalkan) KirimSurel(context.Context, models.Surel) error { p.henti(); return nil }

// TestAksesKorespondensi - keputusan work owner 06-10-2026 (tab hanya JOB_POSITION "IT", aturan Loss Record): jalur
// DITOLAK di keenam metode - tanpa identitas 401, jabatan lain / akun tak ada 403, galat baca jabatan diteruskan; akses
// diperiksa SEBELUM status baris (korespondensi terkirim / nomor asing tetap 403, bukan 409 / 404) dan tanpa email keluar.
func TestAksesKorespondensi(t *testing.T) {
	ctx := context.Background()
	kt, st := &korespondensiTiruan{}, &surelTiruan{}
	svc := layananKorespondensi(kt, st)
	k, _ := svc.TambahKorespondensi(ctx, pelakuKor, "UJI-NB-1")
	isi := models.IsiKorespondensi{Subject: "UJI", Penerima: []models.PenerimaKorespondensi{{RecipientMail: "uji.a@contoh.test"}}}
	if _, err := svc.KirimKorespondensi(ctx, pelakuKor, "UJI-NB-1", k.No, isi); err != nil {
		t.Fatal(err)
	}
	metode := map[string]func(inti.Pelaku, string) error{
		"daftar": func(p inti.Pelaku, _ string) error { _, err := svc.DaftarKorespondensi(ctx, p, "UJI-NB-1"); return err },
		"tambah": func(p inti.Pelaku, _ string) error { _, err := svc.TambahKorespondensi(ctx, p, "UJI-NB-1"); return err },
		"simpan": func(p inti.Pelaku, no string) error {
			_, err := svc.SimpanKorespondensi(ctx, p, "UJI-NB-1", no, isi)
			return err
		},
		"hapus": func(p inti.Pelaku, no string) error { return svc.HapusKorespondensi(ctx, p, "UJI-NB-1", no) },
		"kirim": func(p inti.Pelaku, no string) error {
			_, err := svc.KirimKorespondensi(ctx, p, "UJI-NB-1", no, isi)
			return err
		},
		"saran": func(p inti.Pelaku, _ string) error { _, err := svc.SaranPenerimaKorespondensi(ctx, p, "a"); return err },
	}
	for nama, m := range metode {
		for _, no := range []string{k.No, "999"} {
			if err := m(inti.Pelaku{}, no); !errors.Is(err, inti.ErrTanpaIdentitas) {
				t.Errorf("%s %s 401: %v", nama, no, err)
			}
			for _, p := range []inti.Pelaku{{AkunID: "UJI-UW"}, {AkunID: "UJI-TAK-ADA"}} {
				if err := m(p, no); !errors.Is(err, inti.ErrTanpaWewenang) {
					t.Errorf("%s %s 403 %s: %v", nama, no, p.AkunID, err)
				}
			}
		}
		if err := m(inti.Pelaku{AkunID: "UJI-GALAT"}, k.No); err == nil || errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("%s: galat baca jabatan harus diteruskan: %v", nama, err)
		}
	}
	if len(kt.daftar) != 1 || len(st.kirim) != 1 {
		t.Errorf("jalur ditolak tidak boleh menulis / mengirim: %d baris, %d email", len(kt.daftar), len(st.kirim))
	}
	if !svc.AksesTabIT(ctx, pelakuKor) || svc.AksesTabIT(ctx, inti.Pelaku{AkunID: "UJI-UW"}) {
		t.Error("AksesTabIT")
	}
}

// TestDaftarKorespondensiBersihkanUlang -Message dibersihkan ulang saat dibaca; pembersih idempoten atas keluarannya.
func TestDaftarKorespondensiBersihkanUlang(t *testing.T) {
	kt := &korespondensiTiruan{daftar: []models.Korespondensi{{No: "1", Message: `<p onclick="x">halo</p><script>alert(1)</script>`}}}
	d, err := layananKorespondensi(kt, nil).DaftarKorespondensi(context.Background(), pelakuKor, "UJI-NB-1")
	if err != nil || d.Daftar[0].Message != "<p>halo</p>" {
		t.Fatalf("%v %+v", err, d)
	}
	for _, s := range []string{`<p onclick="x">a<b>b</b></p><img src=x onerror=y>`, `<a href="javascript:x">k</a><ul><li>z</li></ul>`, "teks & <i>miring</i>"} {
		if satu := BersihkanHTMLKorespondensi(s); BersihkanHTMLKorespondensi(satu) != satu {
			t.Errorf("tidak idempoten: %q -> %q -> %q", s, satu, BersihkanHTMLKorespondensi(satu))
		}
	}
}
