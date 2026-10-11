package services_test

import (
	"errors"
	"io"
	"testing"

	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/services"
	"nusantarare/modul/bordereaux/backend/tiruan"
)

// berkasDiPembuat - berkas UJI di tangan pembuatnya (Edit / ViewStage 0 bagi pembuat).
func berkasDiPembuat(g *tiruan.Gudang) string {
	g.Header["BDX-UJI.2"] = models.Header{BdxID: "BDX-UJI.2", UserInput: "UJI-MAKER", Position: "UJI-MAKER", Type: "PREMIUM",
		TypeBusiness: "FIRE"}
	return "BDX-UJI.2"
}

func layananLampiran(g *tiruan.Gudang) (*services.Layanan, *tiruan.Penyimpanan) {
	p := tiruan.BaruPenyimpanan()
	return layanan(g).DenganPenyimpanan(p), p
}

// Upload File / Delete hanya untuk berkas yang boleh di-Edit - pengecualian IT Developer Pega dibuang (keputusan work
// owner 08-10-2026: "kalo cuman view jangan bisa upload dan delete").
func TestBolehLampiranHanyaSaatEdit(t *testing.T) {
	diPembuat := models.Header{UserInput: "UJI-MAKER", Position: "UJI-MAKER"}
	selesai := models.Header{UserInput: "UJI-MAKER", Position: "", StatusAksep: models.StatusResolveComplete}
	superLihat := services.Aktor{AkunID: "UJI-SUPER", Superadmin: true}
	for _, k := range []struct {
		nama string
		a    services.Aktor
		h    models.Header
		mau  bool
	}{
		{"pembuat, berkas di tangannya", maker, diPembuat, true},
		{"pembuat, sudah Resolve-Complete", maker, selesai, false},
		{"checker", checker, models.Header{UserInput: "UJI-MAKER", Position: models.PosisiChecker}, false},
		{"bukan pembuat", services.Aktor{AkunID: "UJI-LAIN", Penuh: true}, diPembuat, false},
		{"superadmin mengambil alih berkas di tangan pembuat", super, diPembuat, true},
		{"superadmin, berkas selesai (hanya View)", super, selesai, false},
		{"superadmin, berkas di Checker (hanya View)", super, models.Header{UserInput: "UJI-MAKER", Position: models.PosisiChecker}, false},
		{"superadmin ber-menu View only", superLihat, diPembuat, false},
		{"pembuat ber-menu View only", services.Aktor{AkunID: "UJI-MAKER"}, diPembuat, false},
		{"tanpa akun", tamu, diPembuat, false},
	} {
		if got := services.BolehLampiran(k.a, k.h); got != k.mau {
			t.Errorf("%s: %v, mau %v", k.nama, got, k.mau)
		}
	}
	g := tiruan.Contoh()
	id := berkasDiPembuat(g)
	r, err := layanan(g).Buka(ctx, maker, id)
	if err != nil || !r.Hak.Lampiran {
		t.Errorf("Buka: hak lampiran %+v %v", r.Hak, err)
	}
}

func TestUnggahLampiranSepertiAttachDocBdxPost(t *testing.T) {
	g := tiruan.Contoh()
	id := berkasDiPembuat(g)
	l, p := layananLampiran(g)
	if _, err := l.UnggahLampiran(ctx, checker, id, "00001", "a.pdf", []byte("isi")); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("checker: %v", err)
	}
	if _, err := l.UnggahLampiran(ctx, maker, id, "99999", "a.pdf", []byte("isi")); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("kategori asing: %v", err)
	}
	if _, err := l.UnggahLampiran(ctx, maker, id, "00001", "  ", []byte("isi")); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("tanpa nama: %v", err)
	}
	if _, err := l.UnggahLampiran(ctx, maker, "BDX-TIDAK-ADA", "00001", "a.pdf", []byte("isi")); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("berkas tidak ada: %v", err)
	}
	b, err := l.UnggahLampiran(ctx, maker, id, "00001", `C:\fakepath\Nota SOA.PDF`, []byte("isi pdf"))
	if err != nil {
		t.Fatal(err)
	}
	mau := models.Lampiran{ID: "20261004140507123", BdxID: id, KategoriID: "00001", Kategori: "UJI SOA", FileName: "Nota SOA.PDF",
		Ekstensi: "pdf", Username: "UJI-MAKER", StorageID: "UJI-IMG-1"}
	if b != mau || g.Lampiran[b.ID] != mau {
		t.Errorf("lampiran %+v", b)
	}
	m := p.Masuk[0]
	if m.Folder != "Contract" || m.Durasi != 3600 || m.Ext != "pdf" || m.NamaFile != "Nota SOA.PDF" || m.Pengguna != "UJI-MAKER" {
		t.Errorf("InsertGoogleStorage_Act %+v", m)
	}
	if _, ada := p.Objek["UJI-IMG-1"]; !ada {
		t.Error("Insert_T_Storage_SQL tidak dijalankan")
	}
	// Milidetik yang sama: ID berikutnya (PK M_ATTACHMENTBORDEREAUX).
	b2, err := l.UnggahLampiran(ctx, super, id, "00000", "b.xlsx", []byte("isi xlsx"))
	if err != nil || b2.ID != "20261004140507124" || b2.Kategori != "Others" {
		t.Errorf("unggah kedua %+v %v", b2, err)
	}
	p.Gagal = penyimpanan.ErrStorageGagal
	if _, err := l.UnggahLampiran(ctx, maker, id, "00001", "c.pdf", []byte("isi")); !errors.Is(err, penyimpanan.ErrStorageGagal) || len(g.Lampiran) != 2 {
		t.Errorf("penyimpanan gagal: %v, %d lampiran", err, len(g.Lampiran))
	}
}

func TestKategoriDaftarUnduhOffice(t *testing.T) {
	g := tiruan.Contoh()
	id := berkasDiPembuat(g)
	l, _ := layananLampiran(g)
	pdf, _ := l.UnggahLampiran(ctx, maker, id, "00001", "a.pdf", []byte("ISI PDF"))
	xlsx, _ := l.UnggahLampiran(ctx, maker, id, "00001", "b.xlsx", []byte("ISI XLSX"))
	k, err := l.KategoriLampiran(ctx, maker, id)
	if err != nil || len(k) != 2 || k[0].Cacah != 0 || k[1].Cacah != 2 || k[1].Nama != "UJI SOA" {
		t.Errorf("kategori %+v %v", k, err)
	}
	d, err := l.DaftarLampiran(ctx, checker, id, "00001")
	if err != nil || len(d) != 2 || d[0].FileName != "a.pdf" {
		t.Errorf("daftar %+v %v", d, err)
	}
	f, err := l.UnduhLampiran(ctx, checker, id, "00001", pdf.ID)
	if err != nil {
		t.Fatal(err)
	}
	isi, _ := io.ReadAll(f.Isi)
	if string(isi) != "ISI PDF" || f.Nama != "a.pdf" || f.Mime != "application/pdf" {
		t.Errorf("unduh %q %+v", isi, f)
	}
	if _, err := l.UnduhLampiran(ctx, checker, id, "00000", pdf.ID); !errors.Is(err, services.ErrLampiranTidakAda) {
		t.Errorf("kategori lain: %v", err)
	}
	if _, err := l.TautanOffice(ctx, checker, id, "00001", pdf.ID); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("office untuk pdf: %v", err)
	}
	if u, err := l.TautanOffice(ctx, checker, id, "00001", xlsx.ID); err != nil || u != "URL-UJI-IMG-2" {
		t.Errorf("office xlsx %q %v", u, err)
	}
}

func TestHapusLampiranSepertiDeleteAttachmentBdx(t *testing.T) {
	g := tiruan.Contoh()
	id := berkasDiPembuat(g)
	l, p := layananLampiran(g)
	b, _ := l.UnggahLampiran(ctx, maker, id, "00001", "a.pdf", []byte("ISI PDF"))
	if err := l.HapusLampiran(ctx, checker, id, "00001", b.ID); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("checker: %v", err)
	}
	p.Gagal = penyimpanan.ErrStorageGagal
	if err := l.HapusLampiran(ctx, maker, id, "00001", b.ID); !errors.Is(err, penyimpanan.ErrStorageGagal) {
		t.Errorf("penyimpanan gagal: %v", err)
	}
	if _, ada := g.Lampiran[b.ID]; !ada {
		t.Error("hapus jarak jauh gagal tetapi baris lampiran terhapus")
	}
	p.Gagal = nil
	if err := l.HapusLampiran(ctx, maker, id, "00001", b.ID); err != nil {
		t.Fatal(err)
	}
	if _, ada := g.Lampiran[b.ID]; ada || len(p.Objek) != 0 || len(p.Hapus) != 1 || p.Hapus[0] != b.StorageID {
		t.Errorf("sesudah hapus: lampiran %v objek %v hapus %v", g.Lampiran, p.Objek, p.Hapus)
	}
	// Baris lama tanpa T_STORAGE_ID: langsung dihapus, tanpa DeleteGoogleStorage_Act.
	g.Lampiran["LAMA"] = models.Lampiran{ID: "LAMA", BdxID: id, KategoriID: "00001", FileName: "lama.pdf"}
	if err := l.HapusLampiran(ctx, maker, id, "00001", "LAMA"); err != nil || len(p.Hapus) != 1 {
		t.Errorf("baris lama: %v, hapus %v", err, p.Hapus)
	}
}
