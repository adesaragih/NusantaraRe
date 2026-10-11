# Langkah work owner — menyalakan Cover Life (`coverlife`, migrasi modul 085-086 + slot 957)

> ## ⛔ LANGKAH PERTAMA, SEBELUM APA PUN: CADANGAN P3
>
> **Tanpa cadangan ID + JSONDATA = TIDAK BOLEH `-migrate`.** 086 MEMBUANG JSONDATA (satu-satunya tempat kunci
> `pxObjClass`, `pyRuleHarness` dan bentuk JSON asli) dan view `COVER_LIFE`.
>
> ```powershell
> mkdir "D:\NUSARE DEV\NUSARE\cadangan-coverlife-20261008"     # <tanggal> hari menjalankan
> cd "D:\NUSARE DEV\NUSARE\cadangan-coverlife-20261008"
> sqlplus POOLDATA@<db> @"D:\NUSARE DEV\NUSARE\NusantaraRe\modul\coverlife\docs\sql\cov_p3_cadangan.sql"
> ```
>
> Hasil: `M_COVER_LIFE_LENGKAP.csv` (ID + JSONDATA utuh) dan `COVER_LIFE_SEBELUM.csv` (isi view: ID, COVER, NOTE),
> keduanya **4 baris** (cacah dicetak di akhir berkas). **Buka keduanya dan periksa sebelum lanjut.** Urutan resmi
> (P1 → P2 → P3) ada di bawah; kotak ini ada di atas supaya tidak terlewat.

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini (executor tidak terhubung ke Oracle; pembacaan DB ditolak pengaman sesi). Keputusan work owner
> 08-10-2026 K0, C1-C4, K5 (`../MODUL.md`). Skema DEV: POOLDATA. Angka DEV (4 baris 100001-100004, NOTE NULL 4/4,
> sequence last_number 5, maks 17 byte) = PEMBANDING dari WO 08-10-2026, bukan angka tetap.

Keadaan awal: `M_COVER_LIFE` (ID VARCHAR2(10) NOT NULL PK `SYS_C009203`, JSONDATA CLOB `ENSURE_M_COVER_LIFE_JSON`) +
view `COVER_LIFE` (`SELECT ID, a.JSONDATA.Cover, a.JSONDATA.Note FROM M_COVER_LIFE a`); sequence `M_COVER_LIFE_SEQ`;
prosedur `PEGA_M_COVER_LIFE`. Tidak ada pembaca lain. `COVERAGE`, `COVERAGE_FACIN`, `COVERAGETRAVEL*`, `COVERNOTE*` milik
aplikasi lain - TIDAK disentuh. **Nama tabel TETAP** (C1, tanpa RENAME).

**Semua SQL = berkas SQL\*Plus** di `modul/coverlife/docs/sql/`. Jalankan sebagai BERKAS (`@"…\cov_xxx.sql"`) dari folder
cadangan - hanya lewat berkas `SET TERMOUT OFF` berlaku. SQL\*Plus 12.2+ atau SQLcl. ⛔ CSV cadangan dapat memuat kunci
identitas operator: simpan di tempat cadangan DBA, jangan disalin ke dokumen / tiket.

| Berkas | Langkah | Sifat |
| --- | --- | --- |
| `cov_p3_cadangan.sql` | (a) P3 cadangan ID + JSONDATA + isi view | baca-saja |
| `cov_p1_kunci.sql` / `cov_p1_kunci_dba.sql` | (a) P1 kunci dan sesi | baca-saja |
| `cov_p2_acuan.sql` | (a) P2 acuan - DUA kali | baca-saja |
| `cov_bukti.sql` | (a) bukti C1: ekspresi 086 = view, 4 baris (NOTE NULL di keduanya); C2: ID berikutnya belum terpakai | baca-saja |
| `cov_c_verifikasi.sql` | (c) verifikasi + CSV sesudah | baca-saja |
| `cov_d_hak.sql` | (d) hak menu (K5) | MENULIS M_LOGIN_GO_MENU |
| `cov_keadaan.sql` | Pemulihan - penentu keadaan | baca-saja |

## (a) Prasyarat P1-P3 — WAJIB SEBELUM `-migrate`

### ⛔ P3 — cadangan ID + JSONDATA + isi view (PALING PENTING)

Lihat kotak di atas berkas ini: `@cov_p3_cadangan.sql` dari folder `cadangan-coverlife-<tanggal>\`. Dua CSV, 4 baris
masing-masing. **Tidak ada cadangan = tidak boleh `-migrate`.**

### ⛔ P1 — hentikan penulis, cek kunci

1. **Hentikan Pega Cover** (penyimpan `InboxCoverLife` → prosedur `PEGA_M_COVER_LIFE`, masih VALID sampai 086) dan
   **backend :8080** (tetap mati sampai (e)). Alasannya: 086 mengisi COVER / NOTE lalu membuang JSONDATA dalam pernyataan
   terpisah tanpa transaksi; baris yang disisip di antaranya kehilangan isinya.
2. `@cov_p1_kunci.sql`: kueri 1 **nol baris**; kueri 2 = sesi POOLDATA yang masih tersambung - pastikan tidak ada Pega /
   backend. ORA-00942 pada V$ = minta DBA menjalankan `cov_p1_kunci_dba.sql` (`DBA_DML_LOCKS`, nol baris).

### ⛔ P2 — acuan baris, DUA kali

`@cov_p2_acuan.sql` (ditambahkan ke `cov_acuan.txt`): **4 baris**, ID tertinggi 100004, sequence last_number 5, sidik
`SUM(ORA_HASH(ID || '|' || COVER || '|' || NOTE))` lewat view (`N_COVER` 4, `N_NOTE` 0), nilai > lebar kolom (**nol
baris**), panjang maksimum (17), NULLABLE asli JSONDATA, constraint, teks view, objek / dependensi / grant, dan **nol
baris T_MIGRASI bernomor 085-089 / 957** (K0). Jalankan **sekali sebelum P3** dan **sekali lagi tepat sebelum
`-migrate`** ((c)0). Satu angka di luar harapan, atau dua keluaran berbeda = BERHENTI.

### Bukti C1 / C2 — sebelum `-migrate`

`@cov_bukti.sql` → `cov_bukti.txt`: ekspresi pengisian 086 (`m.JSONDATA.Cover`, `m.JSONDATA.Note`, notasi titik PERSIS
teks view) sebagai SELECT, dibandingkan dengan view. Harus: A `N` 4, `N_COVER` = `N_VIEW_COVER` = 4, `N_NOTE` =
`N_VIEW_NOTE` = 0; B empat baris `SAMA` = 1; C **nol baris**; D `MAKS_COVER` ≤ 200 (17) dan nol baris melewati lebar; E
**nol baris** (tanpa kembar); F ID berikutnya `100005` `TERPAKAI` = 0. **Kirim berkas keluarannya ke executor /
laporan** (butir yang menunggu WO). Satu angka berbeda = BERHENTI, jangan `-migrate`.

## (b) `-migrate`

```powershell
. .\muat-env.ps1                 # dot-source; ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate        # = target Makefile `migrate`
```

Log yang benar: `dijalankan:` **085, 086, 957** (085-086 tampil di antara langkah modul lain bernomor 0xx bila ada yang
baru; 957 sesudah 955), **tanpa** "dilewati, objeknya sudah ada" untuk objek COVER. 085 = satu ALTER … ADD; 086 = isi
COVER / NOTE dari JSONDATA → pemeriksaan C1.3 → buang JSONDATA → buang view COVER_LIFE; 957 = baris menu.

**Bila 086 berhenti dengan `ORA-01407: cannot update ("POOLDATA"."M_COVER_LIFE"."ID") to NULL`**: itu pemeriksaan C1.3
yang BEKERJA - ada baris yang kolomnya BERBEDA dari view (NULL lawan terisi, atau isi berbeda; NULL = NULL tidak pernah
menghentikan). Tidak ada yang berubah (JSONDATA dan view utuh, 086 tidak tercatat). JANGAN mengulang; laporkan ke WO
dengan keluaran `cov_bukti.sql` kueri C.

## (c) Verifikasi — baca-saja

`@cov_c_verifikasi.sql` → `cov_verifikasi.txt`:

1. `M_COVER_LIFE` = TABLE dengan kolom PERSIS `ID` VARCHAR2(10) NOT NULL, `COVER` VARCHAR2(200), `NOTE` VARCHAR2(1000)
   (urutan view lama); nol JSONDATA; `ALL_PARTIAL_DROP_TABS` nol baris;
2. **objek `COVER_LIFE` sudah tidak ada** (`N_OBJEK_COVER_LIFE` 0); `M_COVER_LIFE_SEQ` tetap; `PEGA_M_COVER_LIFE` INVALID
   (diharapkan, C2 - jangan dikompilasi ulang);
3. **4 baris**, ID tertinggi sama, **sidik `H_BARIS` SAMA dengan acuan** (a) P2 kueri 2;
4. PK `SYS_C009203` tetap ada; `ENSURE_M_COVER_LIFE_JSON` hilang;
5. `T_MIGRASI` memuat `085_cover_life_kolom`, `086_cover_life_satu_tabel`, `957_menu_coverlife`; `M_NAV_MENU`
   `coverlife` 'Cover Life' MASTER TREATY **16** '1' dengan STATUS_AKTIF sama seperti baris `causeoflosslife`.

MINUS dua arah (isi sama dengan `COVER_LIFE_SEBELUM.csv`):

```powershell
$a = Get-Content 'COVER_LIFE_SEBELUM.csv' -Encoding UTF8; $b = Get-Content 'COVER_LIFE_SESUDAH.csv' -Encoding UTF8
"$($a.Count) baris sebelum, $($b.Count) sesudah"
Compare-Object $a $b -CaseSensitive | Group-Object SideIndicator | Select-Object Name, Count   # kosong = sama
```

## (d) Hak menu (K5)

`@cov_d_hak.sql`: setiap akun yang memegang `causeoflosslife` mendapat `coverlife` dengan HAK yang SAMA (LIHAT /
PENUH); aman diulang (`NOT EXISTS`). Periksa `AKUN_COVERLIFE` = `AKUN_CAUSEOFLOSSLIFE` dan kueri ketiga nol baris. Akun
lain: centang menu "Cover Life" di Kelola User.

## (e) Restart backend dan periksa layar

```powershell
go run ./cmd/api                 # = target Makefile `run-api`, jendela yang sudah memuat env
```

- **Cover Life** (MASTER TREATY, urutan 16, sesudah Disease Life): judul "COVER"; form Cover / Note (tanpa medan ID);
  grid ID / Cover (Note tidak tampil) 50 baris, ID menaik; Save satu Cover uji (ID baru `100005`), Edit itu (Note ikut
  terisi), Cancel; teks tersimpan apa adanya (huruf tidak diubah); Cover kembar (tanpa beda huruf) ditolak.
- Akun View only: grid tampil tanpa form dan tanpa Edit.

Pega Cover TIDAK dinyalakan lagi untuk menyimpan (prosedur INVALID - diterima WO C2).

## Pemulihan — 085-086 / 957 gagal di tengah

Pelari tanpa transaksi. Backend dan Pega tetap mati sampai pulih. `@cov_keadaan.sql` memberi `ADA_VIEW`, `ADA_COVER`,
`ADA_NOTE`, `ADA_JSONDATA`, `SETENGAH_TERBUANG`, `T085`, `T086`, `T957`, `ADA_MENU`. Ikuti SATU baris:

| # | Keadaan | Maju (tindakan) | Mundur boleh? |
| --- | --- | --- | --- |
| K0 | 085 gagal: `ADA_COVER` 0, `T085` 0 | ulangi (b) (ALTER ADD atomik) | tidak perlu |
| K1 | SESUDAH ADD, tidak tercatat: `ADA_COVER` 1, `ADA_NOTE` 1, `T085` 0 | jangan ulangi 085 (ORA-01430) - `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('085_cover_life_kolom', SYSDATE); COMMIT;` lalu ulangi (b) | YA (JSONDATA utuh): `085_down`, hapus catatannya |
| K2 | 086 berhenti sebelum DROP JSONDATA: `ADA_JSONDATA` 1, `ADA_VIEW` 1, `T086` 0 | ORA-12899 / **ORA-01407 (pemeriksaan C1.3)** = BERHENTI, laporkan ke WO; galat lain: ulangi (b) (isi menulis nilai yang sama) | YA: `085_down` (pelindungnya lolos selama JSONDATA ada), hapus catatan 085 |
| K3 | SESUDAH DROP JSONDATA: `ADA_JSONDATA` 0 (`SETENGAH_TERBUANG` 1 = `ALTER TABLE POOLDATA.M_COVER_LIFE DROP COLUMNS CONTINUE;` dulu), `ADA_VIEW` 1 atau 0, `T086` 0 | **MAJU WAJIB diselesaikan**: ulangi (b) (blok 1-3 melewati diri, blok 4 membuang view bila masih ada, 086 tercatat, lalu 957). Pastikan (c) | **TIDAK** tanpa 086 tercatat: `085_down` DILARANG (isi hanya tinggal di kolom); pelindungnya gagal keras ORA-00904 |
| K4 | SESUDAH MENU: `T957` 0 tetapi `ADA_MENU` 1 | catat 957 manual: `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('957_menu_coverlife', SYSDATE); COMMIT;`. `ADA_MENU` 0 = ulangi (b) (INSERT `NOT EXISTS`) | `957_down` (hapus hak lalu baris menu) |

Nama langkah: `085_cover_life_kolom`, `086_cover_life_satu_tabel`, `957_menu_coverlife`.

## Jalur mundur

Hanya bila semua langkah yang dibongkar TERCATAT. Hentikan backend. DBA menjalankan berkas `_down` di
`modul/coverlife/backend/migrations/` MENURUN: 957, 086, 085 (ganti `{skema}` dengan POOLDATA), lalu hapus catatan
`T_MIGRASI` ketiganya; backend dari kode sebelum modul ini.

- `086_down`: pelindung `UPDATE … SET COVER = COVER, NOTE = NOTE WHERE 1 = 0` (ORA-00904 bila kolom sumber tidak ada),
  JSONDATA dibangun ulang `JSON_OBJECT('Cover' VALUE …, 'Note' VALUE … ABSENT ON NULL)` (kunci huruf persis Pega; NOTE
  NULL kembali TANPA kunci - view tetap membacanya NULL) + constraint IS JSON bernama asli, lalu CREATE VIEW `COVER_LIFE`
  PERSIS teks DEV (terakhir; view yang sudah ada = ORA-00955, lewati pernyataan itu). `pxObjClass`, `pyRuleHarness` asli
  hanya dari cadangan P3.
- `085_down` mulai dengan pelindung `UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0` (ORA-00904 bila JSONDATA tidak ada),
  lalu DROP (COVER, NOTE).
- JSONDATA dikembalikan **NULLABLE** - bila acuan (a) P2 kueri 4 menunjukkan `NULLABLE = N`, DBA menambah
  `ALTER TABLE POOLDATA.M_COVER_LIFE MODIFY (JSONDATA NOT NULL);` sesudah mundur.
- Prosedur `PEGA_M_COVER_LIFE` perlu dikompilasi ulang DBA. ⛔ Jangan memakai `-migrate-down` di DEV.
- Jatah K0 (premiumlistlife `050-079`, treatycontractout slot `956`) TIDAK ikut mundur - itu tabel jatah dokumen.
