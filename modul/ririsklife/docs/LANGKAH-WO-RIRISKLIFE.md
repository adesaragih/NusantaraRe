# Langkah work owner — menyalakan R/I Risk (`ririsklife`, migrasi inti 935-941)

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini (executor tidak terhubung ke Oracle). Keputusan work owner 08-10-2026 K1-K4 (`../MODUL.md`).
> Skema DEV: POOLDATA. Angka DEV (117 ringkasan, 16.398 rincian, 2 yatim) = PEMBANDING dari WO 08-10-2026, bukan angka
> tetap. Pola sama dengan `modul/ricommlife/docs/LANGKAH-WO-RICOMMLIFE-SATU-TABEL.md`.

Keadaan awal: `M_RIRISK_LIFE_SUMMARY` (ID PK + JSONDATA) + view `RIRISK_LIFE_SUMMARY`; `M_RIRISK_LIFE` (ID VARCHAR2(6)
tanpa PK, JSONDATA, kolom datar warisan BASI, indeks `INDEX4`) + view `RIRISK_LIFE`; `M_RIRISK_LIFE_TEMP` (TIDAK
disentuh); prosedur `PEGA_M_RIRISK_LIFE*`; skema `GL` punya SELECT atas `M_RIRISK_LIFE`.

**Berkas SQL\*Plus** di `modul/ririsklife/docs/sql/`. Jalankan sebagai BERKAS dari folder cadangan, mis.
`@"D:\NUSARE DEV\NUSARE\NusantaraRe\modul\ririsklife\docs\sql\ririsk_a4_cadangan.sql"` - hanya lewat berkas
`SET TERMOUT OFF` berlaku. SQL\*Plus 12.2+ atau SQLcl. ⛔ Keluaran dan CSV memuat OPERATORID (nama orang): simpan di
tempat cadangan DBA, jangan disalin ke dokumen / tiket.

| Berkas | Langkah | Sifat |
| --- | --- | --- |
| `ririsk_a1_kunci.sql` / `ririsk_a1_kunci_dba.sql` | (a) P1 kunci dan sesi | baca-saja |
| `ririsk_a2_acuan.sql` | (a) P1 acuan - DUA kali | baca-saja |
| `ririsk_a3_gl_dba.sql` | (a) P2 pemakaian skema GL (DBA) | baca-saja |
| `ririsk_a4_cadangan.sql` | (a) P3 cadangan CSV | baca-saja |
| `ririsk_bukti_konversi.sql` | (a) bukti konversi RISK di dua NLS | baca-saja (ALTER SESSION saja) |
| `ririsk_c_verifikasi.sql` | (c) verifikasi + CSV sesudah | baca-saja |
| `ririsk_d_hak_superadmin.sql` | (d) hak menu | MENULIS M_LOGIN_GO_MENU |
| `ririsk_keadaan.sql` | Pemulihan - penentu keadaan | baca-saja |

## (a) Prasyarat P1-P3 — SEBELUM `-migrate`

### ⛔ P1 — hentikan penulis, kunci, acuan dua kali

1. **Hentikan Pega R/I Risk** (penyimpan R/I Risk / unggah `M_RIRISK_LIFE_TEMP` -> prosedur `PEGA_M_RIRISK_LIFE*` -
   masih VALID sampai 937/940) dan **backend :8080** (tetap mati sampai (e)). Alasannya: 937/940 mengisi kolom lalu
   membuang JSONDATA dalam pernyataan terpisah tanpa transaksi; baris yang disisip di antaranya kehilangan isi JSON-nya
   dan tidak ada di cadangan.
2. `@ririsk_a1_kunci.sql`: kueri 1 (V$LOCKED_OBJECT ⨝ ALL_OBJECTS ⨝ V$SESSION) **nol baris**; kueri 2 = sesi POOLDATA /
   GL yang masih tersambung - pastikan tidak ada Pega / backend. ORA-00942 pada V$ = minta DBA menjalankan
   `ririsk_a1_kunci_dba.sql` (DBA_DML_LOCKS, nol baris).
3. `@ririsk_a2_acuan.sql` (ditambahkan ke `ririsk_acuan.txt`): cacah + ID tertinggi (117 / 16.398), sequence, sidik
   isi ringkasan dan rincian lewat view, profil RISK (11.063 berkoma, 89 bertitik, `BUKAN_ANGKA` **0**), ID rincian
   NULL / ganda (**0** - PK 940 menolaknya), yatim (2) dan ringkasan tanpa rincian (2) - **TIDAK dihapus**, nilai yang
   melewati lebar (nol baris), **NULLABLE asli JSONDATA**, constraint, **kolom setiap indeks `M_RIRISK_LIFE` (INDEX4)**,
   status `PEGA_*`, grant. Satu angka di luar harapan = BERHENTI dan laporkan.

### ⛔ P2 — pemakaian skema GL (DBA)

`@ririsk_a3_gl_dba.sql` sebagai DBA: dependensi, kode, view, dan sinonim GL / siapa pun yang menyebut `M_RIRISK_LIFE*` /
`RIRISK_LIFE*`. **Semua harus nol baris. Satu pemakai saja = BERHENTI, jangan `-migrate`, laporkan ke WO** - RENAME
(938) dan pembuangan JSONDATA (940) akan mematahkannya.

### ⛔ P3 — cadangan CSV

Buat folder `cadangan-ririsklife-<tanggal>\` (mis. `cadangan-ririsklife-20261008\`), jalankan SQL\*Plus dari sana,
`@ririsk_a4_cadangan.sql` → `M_RIRISK_LIFE_SUMMARY_LENGKAP.csv` (ID + JSONDATA), `M_RIRISK_LIFE_LENGKAP.csv` (ID +
JSONDATA + kolom datar), `RIRISK_LIFE_SUMMARY_SEBELUM.csv`, `RIRISK_LIFE_SEBELUM.csv`. Periksa keempatnya terbuka dan
cacah barisnya = acuan. **Tidak ada cadangan = tidak boleh `-migrate`.**

### Bukti konversi RISK (K2) — sebelum `-migrate`

`@ririsk_bukti_konversi.sql` → `ririsk_bukti_konversi.txt`: ekspresi 940 sebagai SELECT di NLS titik (A) lalu koma (B).
A dan B harus SAMA: `N_TERKONVERSI` = `N_RISK` = 16.398, `H_HASIL` sama, `KOMA` = `921.9`, `TITIK` = `580.894351210924`,
kueri C nol baris. Kirim hasilnya ke executor / laporan (butir yang menunggu WO).

### (c)0 — acuan kedua, TEPAT sebelum `-migrate`

`@ririsk_a2_acuan.sql` sekali lagi: semua angka SAMA dengan yang pertama. **Berubah = ada penulis hidup - ulangi P1-P3.**

## (b) `-migrate`

```powershell
. .\muat-env.ps1                 # dot-source; ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate        # = target Makefile `migrate`
```

Log yang benar: `dijalankan:` 935, 936, 937, 938, 939, 940, 941 berurutan. 935/938 = DROP VIEW (berpelindung
ALL_VIEWS) + RENAME (berpelindung tabel sumber); 936/939 = satu ALTER … ADD; 937/940 = isi dari JSONDATA, (940: PK),
buang JSONDATA, indeks; 941 = baris menu.

## (c) Verifikasi — baca-saja

`@ririsk_c_verifikasi.sql` → `ririsk_verifikasi.txt`:

1. `RIRISK_LIFE_SUMMARY` = ID, USEDBY VARCHAR2(200), MODIFIEDDATE VARCHAR2(50), OPERATORID VARCHAR2(200); `RIRISK_LIFE`
   = ID VARCHAR2(6), IDUSEDBY VARCHAR2(100), USEDBY VARCHAR2(1000), AGE VARCHAR2(10), YEAR / MONTH / CONTRACT
   VARCHAR2(10), RISK NUMBER (tanpa presisi); nol JSONDATA; `ALL_PARTIAL_DROP_TABS` nol baris;
2. hanya dua TABLE `RIRISK_LIFE*`; `M_RIRISK_LIFE` / `M_RIRISK_LIFE_SUMMARY` dan VIEW tidak ada; `M_RIRISK_LIFE_TEMP`
   tetap ada;
3. cacah = acuan (117 / 16.398);
4. sidik = acuan (RISK NUMBER lewat TM9 = acuan konversi ber-format K2);
5. sampel RISK pecahan (bandingkan dengan `M_RIRISK_LIFE_LENGKAP.csv`: `921,9` → `921.9`);
6. `PK_RIRISK_LIFE` ada, IS JSON hilang, satu indeks berkolom pertama IDUSEDBY (INDEX4 bila kolomnya IDUSEDBY, selain
   itu `IX_RIRISK_LIFE_IDUSEDBY`), `IX_RIRISK_LIFE_SUMMARY_NAMA` ada;
7. `T_MIGRASI` 935-941 (tujuh baris); `M_NAV_MENU` `ririsklife` 'R/I Risk' MASTER TREATY 11 '1';
8. **catat** status `PEGA_M_RIRISK_LIFE*` (INVALID diharapkan, K3) dan grant (ikut RENAME).

**MINUS dua arah, TANPA tabel baru** (CSV sesudah ditulis berkas verifikasi):

```powershell
foreach ($p in @(@('RIRISK_LIFE_SUMMARY_SEBELUM.csv', 'RIRISK_LIFE_SUMMARY_SESUDAH.csv'),
                 @('RIRISK_LIFE_SEBELUM.csv', 'RIRISK_LIFE_SESUDAH.csv'))) {
  $a = Get-Content $p[0] -Encoding UTF8; $b = Get-Content $p[1] -Encoding UTF8
  "$($p[1]): $($a.Count) baris sebelum, $($b.Count) sesudah"
  Compare-Object $a $b -CaseSensitive | Group-Object SideIndicator | Select-Object Name, Count   # kosong = sama
}
```

`<=` = sebelum MINUS sesudah, `=>` = sesudah MINUS sebelum; kosong = MINUS dua arah nol (`-CaseSensitive` wajib).

## (d) Hak SUPERADMIN

`@ririsk_d_hak_superadmin.sql` (pola LANGKAH-WO ricommlife (d)): hak `PENUH` menu `ririsklife` untuk setiap pemegang
`kelolauser`, aman diulang; cacah hak = cacah pemegang `kelolauser`. Akun lain: centang menu "R/I Risk" di Kelola User.

## (e) Restart backend dan periksa layar

```powershell
go run ./cmd/api                 # = target Makefile `run-api`, jendela yang sudah memuat env
```

- **R/I Risk** (MASTER TREATY, urutan 11): grid ringkasan (117, filter, urut, tanggal); Save ringkasan uji baru, Edit
  namanya (rinciannya ikut berganti nama), Cancel; Detail (grid 200 baris; RISK koma desimal; tambah dan EDIT satu baris uji lalu kembalikan); Upload →
  View Upload → Simpan Upload satu berkas uji; Delete ringkasan uji beserta rinciannya.
- **Product Name Life**: `Choose R/I Risk` menampilkan daftar yang sama.
- **PremiumList Life**: View R/I Risk dan Hitung QR (Type QR) sebuah polis menghasilkan angka yang sama dengan sebelum
  migrasi.

Pega R/I Risk TIDAK dinyalakan lagi untuk menyimpan (prosedur INVALID - diterima WO K3).

## Pemulihan — 935-941 gagal di tengah

Pelari tanpa transaksi. Backend dan Pega tetap mati sampai pulih. `@ririsk_keadaan.sql` memberi keadaan ringkasan
(935-937) dan rincian (938-940) + menu (941). Ikuti SATU baris per pasangan (ringkasan: R = 935, K = 936, S = 937;
rincian: R = 938, K = 939, S = 940):

| # | Keadaan | Maju (tindakan) | Mundur boleh? |
| --- | --- | --- | --- |
| K0 | sebelum RENAME: `ADA_SUMBER` 1, `ADA_TARGET` 0 (`ADA_VIEW` 1 atau 0), `TR` 0 | R aman diulang (DROP VIEW hanya bila VIEW; RENAME hanya bila sumber ada). Ulangi (b) | tidak perlu (view yang sudah dibuang dibuat ulang dengan pernyataan CREATE VIEW terakhir `R_down`) |
| K1 | sesudah RENAME: `ADA_TARGET` 1, `ADA_SUMBER` 0, kolom baru belum (`TK` 0) | bila `TR` 0: kedua blok R melewati diri - ulangi (b) (R tercatat). Lalu K | YA: `R_down` (pelindungnya lolos selama JSONDATA ada) + hapus catatan R |
| K2 | sesudah ADD: kolom baru ada, `TK` 0 | jangan ulangi K (ORA-01430): `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('<nama K>', SYSDATE); COMMIT;` lalu ulangi (b) | YA: `K_down`, `R_down` |
| K3 | di tengah isi: `ADA_JSONDATA` 1, `TS` 0 (blok isi selesai atau dibatalkan utuh; ORA-12899 / ORA-01722 / ORA-02437 PK = BERHENTI, laporkan, jangan memperlebar / menghapus tanpa WO) | ulangi (b): blok isi menulis nilai yang sama, PK / indeks berpelindung | YA (JSONDATA utuh): `S_down` bila PK/indeks sudah dibuat, `K_down`, `R_down` |
| K4 | sesudah DROP JSONDATA: `ADA_JSONDATA` 0 (`SETENGAH_TERBUANG` 1 = `ALTER TABLE POOLDATA.<tabel> DROP COLUMNS CONTINUE;` dulu), `TS` 0 | **MAJU WAJIB diselesaikan**: ulangi (b) (blok buang melewati diri, indeks berpelindung ALL_IND_COLUMNS). Bila semua sudah ada tetapi tidak tercatat: pastikan (c), lalu `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('<nama S>', SYSDATE); COMMIT;` | **TIDAK.** `K_down` / `R_down` DILARANG (isi hanya tinggal di kolom); pelindungnya gagal keras ORA-00904 |
| K5 | sesudah menu: `T941` 0 tetapi `ADA_MENU` 1 | catat 941 manual (INSERT T_MIGRASI `941_m_nav_menu_ririsklife`); `ADA_MENU` 0 = ulangi (b) (INSERT `NOT EXISTS`) | `941_down` (hapus hak lalu baris menu) |

Nama langkah: `935_ririsk_life_summary_ganti_nama`, `936_ririsk_life_summary_kolom`, `937_ririsk_life_summary_satu_tabel`,
`938_ririsk_life_ganti_nama`, `939_ririsk_life_kolom`, `940_ririsk_life_satu_tabel`, `941_m_nav_menu_ririsklife`.

## Jalur mundur

Hanya bila semua langkah yang dibongkar TERCATAT. Hentikan backend. DBA menjalankan berkas `_down` di
`inti/backend/migrations/` MENURUN: 941, 940, 939, 938, 937, 936, 935 (ganti `{skema}` dengan POOLDATA), lalu hapus
catatan `T_MIGRASI` ketujuhnya; backend dari kode sebelum modul ini.

- 940_down / 937_down aman diulang (blok berpelindung): JSONDATA dibangun ulang (`JSON_OBJECT … ABSENT ON NULL`) +
  constraint IS JSON bernama asli; PK dan indeks yang dibuat 940 dibuang. RISK dikembalikan sebagai teks bertitik
  (TM9) - bentuk asli (`"921,9"`) dan `pxObjClass` hanya dari cadangan P3.
- 939_down / 936_down mulai dengan pelindung `UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0` (ORA-00904 bila JSONDATA
  tidak ada).
- 938_down / 935_down: pelindung yang sama, RENAME balik berpelindung, lalu CREATE VIEW PERSIS teks DEV (terakhir;
  view yang sudah ada = ORA-00955, lewati pernyataan itu).
- JSONDATA dikembalikan **NULLABLE** - bila acuan (a) kueri 8 menunjukkan `NULLABLE = N`, DBA menambah
  `ALTER TABLE POOLDATA.<tabel> MODIFY (JSONDATA NOT NULL);` sesudah mundur.
- Kolom datar warisan `M_RIRISK_LIFE` tetap berisi nilai dari JSONDATA (dulu basi) - nilai basi asli ada di
  `M_RIRISK_LIFE_LENGKAP.csv`.
- Prosedur `PEGA_*` perlu dikompilasi ulang DBA. ⛔ Jangan memakai `-migrate-down` di DEV.
