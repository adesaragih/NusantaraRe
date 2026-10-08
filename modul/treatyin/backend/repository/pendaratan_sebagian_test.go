package repository

// Penjaga cacat KEHILANGAN DATA tombol Save — 7 Oktober 2026.
//
// ---------------------------------------------------------------------
// ⛔ CACAT YANG DIJAGA DI SINI
// ---------------------------------------------------------------------
// `MuatKontrak` mengosongkan ke-31 tabel pendaratan satu kontrak lebih dulu,
// lalu menyisipkan ulang dari dokumen. Benar bagi PEMUAT MASSAL, yang
// dokumennya selalu kontrak utuh.
//
// Tombol Save memakai ulang fungsi itu dengan dokumen yang SEBAGIAN: layar
// hanya mengirim properti tab yang PERNAH DIBUKA, sebab tab dirender hanya
// ketika aktif. Maka satu penekanan Save menghapus baris seluruh tab yang
// tidak dibuka, dan layar memuat ulang keadaan yang sudah terlanjur kosong.
//
// Laporan pemakai: *"saat melakukan inputan tiba tiba datanya hilang dan
// semua di reset malah tidak tersimpan"*.
//
// ⚠️ Uji di berkas ini NOL basis data — yang dijaga KEPUTUSANNYA, dan
// keputusan itu sengaja dipisahkan menjadi fungsi murni justru supaya ia
// dapat diuji tanpa Oracle.

import (
	"os"
	"sort"
	"strings"
	"testing"
)

// Tab Limits dan anak-anaknya — yang pemakai laporkan hilang.
var tabelPohonLimitsUji = []string{
	"T_TREATY_LIMITS",
	"T_TREATY_LIMIT_DETAIL",
	"T_TREATY_LIMIT_GROUP",
	"T_TREATY_LIMIT_COB",
	"T_TREATY_LIMIT_ACHIEVEMENT",
	"T_TREATY_LIMIT_GROUP_COB",
	"T_TREATY_LIMIT_AMOUNT",
	"T_TREATY_LIMIT_MEASURE",
}

// ⛔ INILAH CACATNYA. Dokumen tanpa kunci `Limits` TIDAK berwenang atas
// satu pun tabel pohon Limits — termasuk setiap anaknya.
func TestDokumenTanpaLimitsTidakMenyentuhTabelLimits(t *testing.T) {
	// Dokumen khas Save ketika pemakai hanya menyunting kepala: 19 medan
	// kepala, nol larik tab.
	doc := map[string]any{
		"ID": "12345", "TreatyContractName": "X", "ProportionType": "Proportional",
	}
	dialamatkan := TabelDialamatkan(doc)
	for _, tabel := range tabelPohonLimitsUji {
		if dialamatkan[tabel] {
			t.Errorf("%s dialamatkan oleh dokumen yang NOL kunci Limits — barisnya akan terhapus", tabel)
		}
	}
	for _, p := range PetaPendaratan {
		if p.Akar {
			continue
		}
		if dialamatkan[p.Tabel] {
			t.Errorf("%s dialamatkan padahal dokumen nol larik", p.Tabel)
		}
	}
}

// ⭐ LARIK KOSONG TETAP BERWENANG — dan ini bukan kehalusan.
//
// Pemakai yang membuang seluruh baris grid lalu menekan Save memang
// bermaksud mengosongkannya. Membedakan "kunci tidak ada" dari "larik
// kosong" adalah SELURUH inti perbaikan ini; menguji panjang lariknya
// alih-alih keberadaan kuncinya akan membuat grid mustahil dikosongkan.
func TestLarikKosongTetapBerwenang(t *testing.T) {
	dialamatkan := TabelDialamatkan(map[string]any{"Limits": []any{}})
	for _, tabel := range tabelPohonLimitsUji {
		if !dialamatkan[tabel] {
			t.Errorf("%s TIDAK dialamatkan oleh `Limits: []` — grid menjadi mustahil dikosongkan", tabel)
		}
	}
}

// ⛔ ANAK IKUT INDUKNYA, selalu.
//
// Baris anak menunjuk pengenal baris induk yang baru lahir. Membiarkan anak
// hidup sementara induknya diganti meninggalkan baris yatim yang menunjuk
// pengenal yang sudah tidak ada.
func TestAnakSelaluIkutInduknya(t *testing.T) {
	for _, doc := range []map[string]any{
		{"Limits": []any{}},
		{"Share": []any{}},
		{"Installment": []any{}},
		{"FacultativeShareList": []any{}},
	} {
		dialamatkan := TabelDialamatkan(doc)
		for _, p := range PetaPendaratan {
			if p.Induk == "" {
				continue
			}
			if dialamatkan[p.Tabel] != dialamatkan[p.Induk] {
				t.Errorf("anak %s = %v tetapi induk %s = %v",
					p.Tabel, dialamatkan[p.Tabel], p.Induk, dialamatkan[p.Induk])
			}
		}
	}
}

// ⭐ Tabel ber-`LarikGabung` di tingkat pertama berwenang bila SALAH SATU
// lariknya dibawa — bukan hanya yang pertama.
func TestLarikGabungCukupSalahSatu(t *testing.T) {
	for _, p := range PetaPendaratan {
		if p.Induk != "" || len(p.LarikGabung) == 0 {
			continue
		}
		for _, nama := range p.LarikGabung {
			if !TabelDialamatkan(map[string]any{nama: []any{}})[p.Tabel] {
				t.Errorf("%s tidak dialamatkan oleh larik %q", p.Tabel, nama)
			}
		}
	}
}

// ⛔ TABEL AKAR TIDAK PERNAH DIHAPUS lewat jalur Save — ia DIGABUNG.
//
// `T_TREATY_REVISION` memuat skalar dari banyak tab sekaligus (`Exclusions`,
// `Information`, `TotalLimitsROL`, …). Menghapus lalu menyisipkan ulang dari
// dokumen sebagian akan meng-NULL-kan kolom tab yang tidak dibuka — cacat
// yang sama, hanya berpindah dari BARIS ke KOLOM.
func TestTabelAkarTidakPernahDihapusJalurSave(t *testing.T) {
	pilih := TabelDialamatkan(map[string]any{"Limits": []any{}})
	var akar []string
	for _, p := range PetaPendaratan {
		if !p.Akar {
			continue
		}
		akar = append(akar, p.Tabel)
		if !pilih[p.Tabel] {
			t.Errorf("%s harus DIALAMATKAN (ia digabung, bukan dilewati)", p.Tabel)
		}
		if BolehDihapus(p, pilih) {
			t.Errorf("%s BOLEH dihapus jalur Save — kolom tab yang tak dibuka akan hilang", p.Tabel)
		}
	}
	if len(akar) == 0 {
		t.Fatal("nol tabel akar di peta — uji ini kehilangan sasarannya")
	}
	t.Logf("tabel akar yang dijaga: %s", strings.Join(akar, ", "))
}

// ⭐ PEMUAT MASSAL TIDAK BERUBAH. `pilih == nil` berarti seluruhnya dibuang —
// termasuk tabel akar. Itulah yang membuatnya idempoten, dan perbaikan
// tombol Save tidak boleh mencabutnya.
func TestPemuatMassalTetapMembuangSeluruhnya(t *testing.T) {
	for _, p := range PetaPendaratan {
		if !BolehDihapus(p, nil) {
			t.Errorf("%s tidak dibuang pemuat massal — idempotennya hilang", p.Tabel)
		}
	}
}

// ⭐ CERMIN: dokumen UTUH berwenang atas SELURUH tabel. Bila sebuah tabel
// tidak pernah dapat dialamatkan, ia mustahil ditulis lewat Save — dan itu
// kehilangan data yang lain lagi, hanya tanpa suara.
func TestSetiapTabelDapatDialamatkan(t *testing.T) {
	utuh := map[string]any{}
	for _, p := range PetaPendaratan {
		for _, nama := range LarikPeta(p) {
			utuh[nama] = []any{}
		}
	}
	dialamatkan := TabelDialamatkan(utuh)
	var luput []string
	for _, p := range PetaPendaratan {
		if !dialamatkan[p.Tabel] {
			luput = append(luput, p.Tabel)
		}
	}
	sort.Strings(luput)
	if len(luput) > 0 {
		t.Fatalf("tabel yang mustahil ditulis lewat Save:\n  %s", strings.Join(luput, "\n  "))
	}
}

// ⛔ JALUR SAVE WAJIB MEMANGGIL `MuatKontrakSebagian`.
//
// ⚠️ Seluruh uji di atas menjaga KEPUTUSANNYA; yang satu ini menjaga
// bahwa keputusan itu benar-benar DIPAKAI. Tanpa uji ini, menukar satu
// nama fungsi di `simpanDalam` mengembalikan cacatnya sepenuhnya sementara
// delapan uji lain tetap hijau.
func TestJalurSaveMemakaiPemuatSebagian(t *testing.T) {
	b, err := os.ReadFile("simpan_kontrak.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(b)
	if !strings.Contains(src, "g.MuatKontrakSebagian(ctx, tx, id, doc)") {
		t.Fatal("jalur Save tidak memanggil `MuatKontrakSebagian`")
	}
	if strings.Contains(src, "g.MuatKontrak(ctx, tx, id, doc)") {
		t.Fatal("jalur Save masih memanggil `MuatKontrak` — ia mengosongkan ke-31 tabel")
	}
}
