package backend

// Uji invarian atas TEKS DDL - bukan atas Oracle (`L-3`).

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

func seluruhMigrasi(t *testing.T) map[string]string {
	t.Helper()
	entri, err := fs.ReadDir(berkasMigrasi, "migrations")
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, e := range entri {
		isi, err := fs.ReadFile(berkasMigrasi, "migrations/"+e.Name())
		if err != nil {
			t.Fatal(err)
		}
		out[e.Name()] = string(isi)
	}
	if len(out) == 0 {
		t.Fatal("nol berkas migrasi terbaca; pembacanya yang rusak")
	}
	return out
}

// tanpaKomentar membuang baris komentar `--`, menyisakan PERNYATAAN SQL saja -
// berkas migrasi modul ini MENJELASKAN keputusannya di dalam komentar, dan
// sapuan atas teks mentah akan menemukan kata-katanya di prosanya sendiri.
func tanpaKomentar(sql string) string {
	var b strings.Builder
	for _, baris := range strings.Split(sql, "\n") {
		if strings.HasPrefix(strings.TrimSpace(baris), "--") {
			continue
		}
		b.WriteString(baris)
		b.WriteByte('\n')
	}
	return b.String()
}

func majuSaja(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for nama, isi := range seluruhMigrasi(t) {
		if !strings.HasSuffix(nama, "_down.sql") {
			out[nama] = isi
		}
	}
	return out
}

// ⛔ Modul ini TIDAK membuat satu tabel pun, dan itu keputusan - bukan keadaan
// sementara. Model datanya satu dengan Treaty In; yang dipisahkan 25-09-2026
// hanya papan tiketnya. Tabel baru di sini berarti model datanya bercabang.
// tabelMilikModulIni adalah SATU-SATUNYA tabel yang modul ini boleh buat.
//
// ⛔ Daftar-IZIN, bukan pencabutan. Penjaga di bawah lahir menyatakan
// *"modul ini hanya mengubah tabel milik `treatyin`"*, dan ia MENOLAK
// migrasi `442` ketika pemindahan pertama kali dicoba - riwayat penolakan
// itu tercatat di kepala `treatyin/.../426_nilai_selisih.sql`.
//
// Keputusan pemilik proses 4 Oktober 2026
// (`modul/treatyin/docs/KEPUTUSAN-PENYELARASAN-REPO.md` §19) menggantikan
// aturan itu UNTUK DUA TABEL INI SAJA: tiket 76 dan 77 ada di papan modul
// ini, dan sejak ronde itu tabelnya ada di sini juga.
//
// ⛔ Yang TIDAK berubah: tabel ketiga tetap ditolak. Melonggarkan penjaga
// menjadi "modul ini boleh membuat tabel" akan membuang pagar yang masih
// diperlukan - model datanya memang sebagian besar masih satu dengan
// `treatyin`, dan percabangan diam-diam adalah persis yang dicegah.
//
// ⚠️ DITAGIH ke pemilik spec: `treaty-in/SPEC-MODEL-DATA.md` masih
// menyatakan model datanya SATU, dan pesan penjaga lama mengutipnya.
// Pernyataan itu kini tidak lagi benar seluruhnya.
var tabelMilikModulIni = map[string]bool{
	"NILAI_SELISIH":          true,
	"NILAI_SEBELUM_PRO_RATE": true,
}

func TestNolTabelBaru(t *testing.T) {
	pola := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+\{skema\}\.(\w+)`)
	dibuat := map[string]string{}
	for nama, isi := range majuSaja(t) {
		for _, m := range pola.FindAllStringSubmatch(tanpaKomentar(isi), -1) {
			tabel := strings.ToUpper(m[1])
			if !tabelMilikModulIni[tabel] {
				t.Errorf("%s membuat tabel %s; modul ini hanya mengubah tabel milik "+
					"`treatyin`, kecuali dua yang KEPUTUSAN §19 pindahkan ke sini", nama, tabel)
				continue
			}
			dibuat[tabel] = nama
		}
	}
	// ⛔ Dan keduanya BENAR-BENAR dibuat di sini. Daftar-izin yang isinya
	// tidak pernah dipakai adalah izin yang diam-diam menjadi lubang.
	for tabel := range tabelMilikModulIni {
		if dibuat[tabel] == "" {
			t.Errorf("tabel %s diizinkan untuk modul ini tetapi nol migrasi membuatnya; "+
				"bila ia kembali ke `treatyin`, cabut izinnya di sini", tabel)
		}
	}
}

// Tiket 01: kolom dasar berdiri dengan kunci asingnya, dan kunci asing itu
// menunjuk VERSI_KONTRAK - tabelnya sendiri.
func TestKolomDasarBerdiriDenganKunciAsingnya(t *testing.T) {
	isi, ok := majuSaja(t)["440_versi_kontrak_dasar.sql"]
	if !ok {
		t.Fatal("migrasi 440 tidak ada")
	}
	sql := tanpaKomentar(isi)
	if !strings.Contains(sql, "ID_VERSI_KONTRAK_DASAR  NUMBER(19)") {
		t.Error("tiket 01: kolom ID_VERSI_KONTRAK_DASAR tidak ditambahkan")
	}
	if !strings.Contains(sql, "CONSTRAINT FK_VERSI_KONTRAK_DASAR") ||
		!strings.Contains(sql, "REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)") {
		t.Error("tiket 01: kunci asing ke VERSI_KONTRAK tidak dipasang")
	}
	// Kolomnya WAJIB boleh kosong: versi pertama sebuah kontrak berdasar
	// kosong (bahan to-spec B-3), dan baris warisan boleh berdasar kosong
	// sampai tiket 10.
	if strings.Contains(sql, "ID_VERSI_KONTRAK_DASAR  NUMBER(19) NOT NULL") {
		t.Error("tiket 01: kolom dasar NOT NULL; versi pertama berdasar KOSONG")
	}
}

// Tiket 05: langkah PERLUAS - nomor urut versi menjadi boleh kosong.
func TestNomorUrutVersiDilonggarkan(t *testing.T) {
	isi, ok := majuSaja(t)["441_nomor_urut_versi_boleh_kosong.sql"]
	if !ok {
		t.Fatal("migrasi 441 tidak ada")
	}
	sql := tanpaKomentar(isi)
	if !strings.Contains(sql, "MODIFY") || !strings.Contains(sql, "NOMOR_URUT_VERSI  NUMBER(10) NULL") {
		t.Error("tiket 05: NOMOR_URUT_VERSI tidak dilonggarkan menjadi NULL")
	}
	// Lebarnya tidak boleh berubah diam-diam: MODIFY yang menyebut tipe lain
	// mengubah dua hal sekaligus, dan yang kedua tidak pernah ditinjau.
	if strings.Contains(sql, "NUMBER(9)") || strings.Contains(sql, "NUMBER(19)") {
		t.Error("tiket 05: MODIFY mengubah LEBAR kolom, bukan hanya keterisiannya")
	}
}

// ADR-0056 (K-4) dan ADR-U-0029: nol trigger, nol procedure, nol COMMIT.
//
// Tiket 03 MEMINTA trigger INV-54; keputusannya ada di
// `docs/KEPUTUSAN-TIKET-02-03.md` §2 - ditegakkan di services, sejalan dengan
// INV-53 yang sebentuk. Uji ini menjaga keputusan itu.
func TestNolTriggerProcedureDanCommit(t *testing.T) {
	for nama, isi := range seluruhMigrasi(t) {
		atas := strings.ToUpper(tanpaKomentar(isi))
		for _, terlarang := range []string{"CREATE TRIGGER", "CREATE OR REPLACE TRIGGER", "CREATE PROCEDURE", "CREATE FUNCTION", "\nCOMMIT"} {
			if strings.Contains(atas, terlarang) {
				t.Errorf("%s memuat %q; aturan bisnis tidak turun ke basis data (ADR-0056)", nama, terlarang)
			}
		}
	}
}

// gabungan menyatukan seluruh langkah maju menjadi satu teks.
func gabungan(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	for _, isi := range majuSaja(t) {
		b.WriteString(isi)
	}
	return b.String()
}

// INV-18 - perilaku hapus DITETAPKAN SADAR, diadu dengan ERD.md §2.
//
// ⚠️ RALAT 2 Oktober 2026. Uji ini pernah bernama TestNolKaskadeHapus dan
// menuntut nol `ON DELETE` dengan alasan "kebijakan modul ini MENOLAK". Itu
// membaca INV-18 terbalik — invariannya menuntut keputusan yang DINYATAKAN,
// bukan bawaan yang kebetulan cocok, dan sumbernya `ERD.md` §2.
//
// Modul ini punya SATU kunci asing sendiri, dan §2.2 menyatakannya
// [hapus: tolak]: "menghapus versi dasar akan membuat seluruh baris selisih
// kehilangan artinya." TOLAK diwujudkan dengan tidak menulis klausa ON DELETE.
func TestPerilakuHapusSesuaiERD(t *testing.T) {
	pola := regexp.MustCompile(`(?is)CONSTRAINT\s+(FK_\w+)\s+FOREIGN KEY\s*\([^)]*\)\s*REFERENCES\s+\{skema\}\.\w+\s*\([^)]*\)([^,\n]*)`)
	// "" berarti TOLAK: diwujudkan dengan TIDAK menulis klausa ON DELETE.
	const tolak = ""
	mau := map[string]string{
		// §2.2 — "menghapus versi dasar akan membuat seluruh baris selisih
		// kehilangan artinya."
		"FK_VERSI_KONTRAK_DASAR": tolak,
		// ----------------------------------------------------------------
		// PINDAH dari modul `treatyin` 4 Oktober 2026 bersama tabelnya —
		// KEPUTUSAN §19. Barisnya IKUT PINDAH, tidak digandakan: dua tempat
		// yang menyatakan perilaku hapus yang sama akan berselisih, dan
		// yang salah satunya basi tidak akan berbunyi.
		//
		// Keduanya dari `ERD.md` §2.6 dan ERD HTML baris 36/37, dan
		// keduanya sepakat `ikut hapus`.
		"FK_NILAI_SELISIH_1":          "ON DELETE CASCADE",
		"FK_NILAI_SEBELUM_PRO_RATE_1": "ON DELETE CASCADE",
	}

	lihat := map[string]bool{}
	for _, m := range pola.FindAllStringSubmatch(tanpaKomentar(gabungan(t)), -1) {
		nama, ekor := m[1], strings.ToUpper(strings.TrimSpace(m[2]))
		harap, dikenal := mau[nama]
		if !dikenal {
			t.Errorf("kunci asing %s tidak ada di tabel ERD uji ini; tiap kunci asing "+
				"BARU wajib menyebut keputusan ERD.md §2-nya", nama)
			continue
		}
		lihat[nama] = true
		if ekor != strings.ToUpper(harap) {
			t.Errorf("%s: perilaku hapus %q, ERD.md menuntut %q", nama, ekor, harap)
		}
	}
	for nama := range mau {
		if !lihat[nama] {
			t.Errorf("kunci asing %s didaftar uji ini tetapi tidak ada di DDL", nama)
		}
	}
}

// Kaskade HANYA pada berkas yang `MODUL.md` daftarkan.
//
// ⚠️ RALAT 4 Oktober 2026. Uji ini pernah bernama `TestNolKaskadeDiModulIni`
// dan menuntut NOL `ON DELETE` di seluruh modul — benar selama modul ini
// tidak punya tabel sendiri. Sejak KEPUTUSAN §19 memindahkan `NILAI_SELISIH`
// dan `NILAI_SEBELUM_PRO_RATE` ke sini, keduanya membawa kaskadenya dari
// `ERD.md` §2.6. Nama lamanya ditulis di sini supaya pencarian atasnya tetap
// sampai ke tempat ini.
//
// ⛔ Yang TIDAK berubah: kaskade tetap menuntut IZIN. Berkas di luar daftar
// tetap ditolak, dan daftarnya hidup di `MODUL.md` — satu tempat, dibaca
// juga oleh penjaga inti `TestKaskadeHanyaPadaRelasiTerdaftar`.
func TestKaskadeHanyaPadaBerkasTerdaftar(t *testing.T) {
	berkaskade := map[string]bool{"442_nilai_selisih.sql": true}
	for nama, isi := range majuSaja(t) {
		ada := strings.Contains(strings.ToUpper(tanpaKomentar(isi)), "ON DELETE")
		if ada && !berkaskade[nama] {
			t.Errorf("%s memuat ON DELETE tanpa izin; daftarnya di MODUL.md bab kaskade", nama)
		}
		if !ada && berkaskade[nama] {
			t.Errorf("%s didaftar berkaskade tetapi nol ON DELETE; daftar yang tidak "+
				"menggigit adalah daftar yang diam-diam kosong", nama)
		}
	}
}

// Yang TIDAK diuji di sini, sebab penjaga inti sudah menegakkannya atas SELURUH
// modul - menyalinnya ke sini berarti dua tempat yang harus sepakat:
//
//   jalur mundur tiap langkah   -> penjaga `TestSetiapLangkahPunyaJalurMundur`
//   nomor di rentang modulnya   -> penjaga `TestSetiapMigrasiDiRentangAtauSlotModulnya`
//   slot menu hanya menyentuh
//   baris menu modulnya sendiri -> penjaga `TestSlotMenuHanyaMenyentuhMenuModulnya`
//
// Yang DI SINI adalah yang khusus modul ini: invarian bernomor dan pernyataan
// keputusan yang penjaga umum tidak dapat mengetahuinya.
