# Langkah work owner — menyalakan Cause Of Loss Life (`causeoflosslife`, migrasi modul 090-092 + slot 955)

> ## ⛔ LANGKAH PERTAMA, SEBELUM APA PUN: CADANGAN P3
>
> **Tanpa cadangan ID + JSONDATA = TIDAK BOLEH `-migrate`.** 092 MEMBUANG JSONDATA (satu-satunya tempat kunci
> `pxObjClass`, `pyRuleHarness` dan bentuk JSON asli - termasuk `"CauseofLoss":""` baris 100001).
>
> ```powershell
> mkdir "D:\NUSARE DEV\NUSARE\cadangan-causeoflosslife-20261008"     # <tanggal> hari menjalankan
> cd "D:\NUSARE DEV\NUSARE\cadangan-causeoflosslife-20261008"
> sqlplus POOLDATA@<db> @"D:\NUSARE DEV\NUSARE\NusantaraRe\modul\causeoflosslife\docs\sql\col_p3_cadangan.sql"
> ```
>
> Hasil: `M_CAUSEOFLOSS_LIFE_LENGKAP.csv` (ID + JSONDATA utuh) dan `CAUSEOFLOSS_LIFE_SEBELUM.csv` (ID, CAUSEOFLOSS
> lewat view), keduanya **4 baris** (cacah dicetak di akhir berkas). **Buka keduanya dan periksa sebelum lanjut.**
> Urutan resmi (P1 → P2 → P3) ada di bawah: P3 dijalankan SESUDAH penulis dihentikan (P1) dan acuan pertama (P2);
> kotak ini ada di atas supaya tidak terlewat.

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini (executor tidak terhubung ke Oracle; pembacaan DB ditolak pengaman sesi). Keputusan work owner
> 08-10-2026 K0-K6 (`../MODUL.md`). Skema DEV: POOLDATA. Angka DEV (4 baris 100001-100004, 100001 kosong, sequence
> last_number 5, maks 9 byte) = PEMBANDING dari WO 08-10-2026, bukan angka tetap. Pola sama dengan
> `modul/benefitlife/docs/LANGKAH-WO-BENEFITLIFE.md`; bedanya: migrasinya MILIK MODUL (`modul/causeoflosslife/backend/
> migrations/`, K0), bukan inti.

Keadaan awal: `M_CAUSEOFLOSS_LIFE` (ID VARCHAR2(10) NOT NULL PK `SYS_C008825`, JSONDATA CLOB
`ENSURE_M_CAUSEOFLOSS_LIFE_JSON`) + view `CAUSEOFLOSS_LIFE` (`SELECT a.ID, a.JSONDATA.CauseofLoss FROM
M_CAUSEOFLOSS_LIFE a`); sequence `M_CAUSEOFLOSS_LIFE_SEQ`; prosedur `PEGA_M_CAUSEOFLOSS_LIFE`. Pembaca lain: modul
`masterproductnamelife` (pemilih Cause Of Loss - nama dan kolom tetap). Objek `M_CAUSE_OF_LOSS`, `D_CAUSE_OF_LOSS`,
`V_M_CAUSE_OF_LOSS`, `V_D_CAUSE_OF_LOSS*`, `T_LISTCAUSEOFLOSS`, `PEGA_M_CAUSE_OF_LOSS`, `PEGA_D_CAUSE_OF_LOSS` dan
sequence-nya milik aplikasi lain - TIDAK disentuh.

**Semua SQL = berkas SQL\*Plus** di `modul/causeoflosslife/docs/sql/`. Jalankan sebagai BERKAS (`@"…\col_xxx.sql"`)
dari folder cadangan - hanya lewat berkas `SET TERMOUT OFF` berlaku. SQL\*Plus 12.2+ atau SQLcl. ⛔ CSV cadangan dapat
memuat kunci identitas operator: simpan di tempat cadangan DBA, jangan disalin ke dokumen / tiket.

| Berkas | Langkah | Sifat |
| --- | --- | --- |
| `col_p3_cadangan.sql` | (a) P3 cadangan ID + JSONDATA | baca-saja |
| `col_p1_kunci.sql` / `col_p1_kunci_dba.sql` | (a) P1 kunci dan sesi | baca-saja |
| `col_p2_acuan.sql` | (a) P2 acuan - DUA kali | baca-saja |
| `col_bukti_k1.sql` | (a) bukti K2: ekspresi 092 = view, 4 baris (100001 NULL di keduanya) | baca-saja |
| `col_c_verifikasi.sql` | (c) verifikasi + CSV sesudah | baca-saja |
| `col_d_hak.sql` | (d) hak menu (K6) | MENULIS M_LOGIN_GO_MENU |
| `col_keadaan.sql` | Pemulihan - penentu keadaan | baca-saja |

## (a) Prasyarat P1-P3 — WAJIB SEBELUM `-migrate`

### ⛔ P3 — cadangan ID + JSONDATA (PALING PENTING)

Lihat kotak di atas berkas ini: `@col_p3_cadangan.sql` dari folder `cadangan-causeoflosslife-<tanggal>\`. Dua CSV, 4
baris masing-masing. **Tidak ada cadangan = tidak boleh `-migrate`.**

### ⛔ P1 — hentikan penulis, cek kunci

1. **Hentikan Pega Cause Of Loss** (penyimpan `InboxCauseofLossLife` → prosedur `PEGA_M_CAUSEOFLOSS_LIFE`, masih VALID
   sampai 092) dan **backend :8080** (tetap mati sampai (e)). Alasannya: 092 mengisi CAUSEOFLOSS lalu membuang
   JSONDATA dalam pernyataan terpisah tanpa transaksi; baris yang disisip di antaranya kehilangan isinya.
2. `@col_p1_kunci.sql`: kueri 1 (`V$LOCKED_OBJECT` ⨝ `ALL_OBJECTS` ⨝ `V$SESSION`) **nol baris**; kueri 2 = sesi
   POOLDATA yang masih tersambung - pastikan tidak ada Pega / backend. ORA-00942 pada V$ = minta DBA menjalankan
   `col_p1_kunci_dba.sql` (`DBA_DML_LOCKS`, nol baris).

### ⛔ P2 — acuan baris, DUA kali

`@col_p2_acuan.sql` (ditambahkan ke `col_acuan.txt`): **4 baris**, ID tertinggi 100004, sequence last_number 5, sidik
`SUM(ORA_HASH(ID || '|' || CAUSEOFLOSS))` lewat view (`N_CAUSEOFLOSS` 3), nilai > 200 byte (**nol baris**), panjang
maksimum (9), NULLABLE asli JSONDATA, constraint, teks view, objek / dependensi / grant, dan **nol baris T_MIGRASI
bernomor 090-099 / 955** (K0). Jalankan **sekali sebelum P3** dan **sekali lagi tepat sebelum `-migrate`** ((c)0). Satu
angka di luar harapan, atau dua keluaran berbeda = BERHENTI.

### Bukti K2 — sebelum `-migrate`

`@col_bukti_k1.sql` → `col_bukti_k1.txt`: ekspresi pengisian 092 (`m.JSONDATA.CauseofLoss`, notasi titik PERSIS teks
view) sebagai SELECT, dibandingkan dengan view. Harus: A `N` 4, `N_EKSPRESI` = `N_VIEW` = 3, `N_ADA_KUNCI` 4; B empat
baris `SAMA` = 1 (100001: EKSPRESI dan VIEW sama-sama kosong); C **nol baris**; D `MAKS_BYTE` ≤ 200 (9) dan nol baris
> 200; E **nol baris** (tanpa kembar); F pemakai `CAUSEID` (100001 tidak dipakai - dilaporkan saja). **Kirim berkas
keluarannya ke executor / laporan** (butir yang menunggu WO). Satu angka berbeda = BERHENTI, jangan `-migrate`.

## (b) `-migrate`

```powershell
. .\muat-env.ps1                 # dot-source; ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate        # = target Makefile `migrate`
```

Log yang benar: `dijalankan:` **090, 091, 092, 955** (090-092 tampil di antara langkah modul lain yang bernomor 0xx
bila ada yang baru; 955 sesudah 949), **tanpa** "dilewati, objeknya sudah ada" untuk objek CAUSEOFLOSS. 090 = DROP
VIEW (berpelindung ALL_VIEWS) + RENAME (berpelindung tabel sumber); 091 = satu ALTER … ADD; 092 = isi CAUSEOFLOSS dari
JSONDATA → pemeriksaan K1.4 → buang JSONDATA; 955 = baris menu.

**Bila 092 berhenti dengan `ORA-01407: cannot update ("POOLDATA"."CAUSEOFLOSS_LIFE"."ID") to NULL`**: itu pemeriksaan
K1.4 yang BEKERJA - ada baris yang kolomnya BERBEDA dari view (NULL lawan terisi, atau isi berbeda; NULL = NULL tidak
pernah menghentikan). Tidak ada yang berubah (JSONDATA utuh, 092 tidak tercatat). JANGAN mengulang; laporkan ke WO
dengan keluaran `col_bukti_k1.sql` kueri C.

## (c) Verifikasi — baca-saja

`@col_c_verifikasi.sql` → `col_verifikasi.txt`:

1. `CAUSEOFLOSS_LIFE` = `ID` VARCHAR2(10) NOT NULL, `CAUSEOFLOSS` VARCHAR2(200) NULLABLE; nol JSONDATA;
   `ALL_PARTIAL_DROP_TABS` nol baris;
2. `CAUSEOFLOSS_LIFE` = **TABLE**; `M_CAUSEOFLOSS_LIFE` **tidak ada**; **nol VIEW** `CAUSEOFLOSS_LIFE`;
   `M_CAUSEOFLOSS_LIFE_SEQ` tetap; `PEGA_M_CAUSEOFLOSS_LIFE` INVALID (diharapkan, K3 - jangan dikompilasi ulang);
3. **4 baris**, ID tertinggi sama, **sidik `H_BARIS` SAMA dengan acuan** (a) P2 kueri 2;
4. PK `SYS_C008825` ada; `ENSURE_M_CAUSEOFLOSS_LIFE_JSON` hilang;
5. `T_MIGRASI` memuat keempat nama `090_causeofloss_life_ganti_nama`, `091_causeofloss_life_kolom`,
   `092_causeofloss_life_satu_tabel`, `955_menu_causeoflosslife`; `M_NAV_MENU` `causeoflosslife` 'Cause Of Loss Life'
   MASTER TREATY **14** '1' dengan STATUS_AKTIF sama seperti baris `planlife`;
6. kueri pemilih Product Name Life tetap memberi 3 baris (100001 NULL tidak cocok LIKE - sama seperti lewat view).

MINUS dua arah (CSV sesudah ditulis berkas verifikasi):

```powershell
$a = Get-Content 'CAUSEOFLOSS_LIFE_SEBELUM.csv' -Encoding UTF8; $b = Get-Content 'CAUSEOFLOSS_LIFE_SESUDAH.csv' -Encoding UTF8
"$($a.Count) baris sebelum, $($b.Count) sesudah"
Compare-Object $a $b -CaseSensitive | Group-Object SideIndicator | Select-Object Name, Count   # kosong = sama
```

## (d) Hak menu (K6)

`@col_d_hak.sql`: setiap akun yang memegang `planlife` mendapat `causeoflosslife` dengan HAK yang SAMA (LIHAT /
PENUH); aman diulang (`NOT EXISTS`). Periksa `AKUN_CAUSEOFLOSSLIFE` = `AKUN_PLANLIFE` dan kueri ketiga nol baris. Akun
lain: centang menu "Cause Of Loss Life" di Kelola User.

## (e) Restart backend dan periksa layar

```powershell
go run ./cmd/api                 # = target Makefile `run-api`, jendela yang sudah memuat env
```

- **Cause Of Loss Life** (MASTER TREATY, urutan 14, sesudah Plan): judul "CAUSE OF LOSS"; grid ID / Cause of Loss 10
  baris, ID menaik (100001 kosong di atas); Save satu Cause of Loss uji (ID baru `100005` bila sequence belum dipakai),
  Edit itu, Cancel; teks tersimpan apa adanya (huruf tidak diubah); nama kembar (tanpa beda huruf) ditolak.
- Akun View only: grid tampil tanpa form dan tanpa Edit.
- **Product Name Life**: pemilih Cause Of Loss tetap menampilkan ILLNESS / ACCIDENT / ANY CAUSE.

Pega Cause Of Loss TIDAK dinyalakan lagi untuk menyimpan (prosedur INVALID - diterima WO K3).

## Pemulihan — 090-092 / 955 gagal di tengah

Pelari tanpa transaksi. Backend dan Pega tetap mati sampai pulih. `@col_keadaan.sql` memberi `ADA_VIEW`, `ADA_SUMBER`,
`ADA_TARGET`, `ADA_KOLOM`, `ADA_JSONDATA`, `SETENGAH_TERBUANG`, `T090`-`T092`, `T955`, `ADA_MENU`. Ikuti SATU baris:

| # | Keadaan | Maju (tindakan) | Mundur boleh? |
| --- | --- | --- | --- |
| K0 | mati SEBELUM RENAME: `ADA_SUMBER` 1, `ADA_TARGET` 0 (`ADA_VIEW` 1 atau 0 - view mungkin sudah dibuang), `T090` 0 | 090 aman diulang (DROP VIEW hanya bila VIEW; RENAME hanya bila sumber ada). Ulangi (b) | tidak perlu; bila tidak jadi maju dan `ADA_VIEW` 0: jalankan pernyataan CREATE VIEW terakhir `090_causeofloss_life_ganti_nama_down.sql` saja |
| K1 | SESUDAH RENAME: `ADA_TARGET` 1, `ADA_SUMBER` 0, `ADA_KOLOM` 0 | bila `T090` 0: kedua blok 090 melewati diri - ulangi (b) (090 tercatat, lalu 091 dst.) | YA: `090_down` (pelindungnya lolos selama JSONDATA ada) + hapus catatan 090 bila ada |
| K2 | SESUDAH ADD: `ADA_KOLOM` 1, `ADA_JSONDATA` 1, `T091` 0 **atau** 092 berhenti sebelum DROP (`T092` 0, `ADA_JSONDATA` 1) | `T091` 0: jangan ulangi 091 (ORA-01430) - `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('091_causeofloss_life_kolom', SYSDATE); COMMIT;` lalu ulangi (b). 092 berhenti ORA-12899 / **ORA-01407 (pemeriksaan K1.4)** = BERHENTI, laporkan ke WO; galat lain: ulangi (b) (isi menulis nilai yang sama) | YA (JSONDATA utuh): `091_down` lalu `090_down`, hapus catatannya |
| K3 | SESUDAH DROP JSONDATA: `ADA_JSONDATA` 0 (`SETENGAH_TERBUANG` 1 = `ALTER TABLE POOLDATA.CAUSEOFLOSS_LIFE DROP COLUMNS CONTINUE;` dulu), `T092` 0 | **MAJU WAJIB diselesaikan**: ulangi (b) (ketiga blok 092 melewati diri, 092 tercatat, lalu 955). Pastikan (c) | **TIDAK.** `091_down` / `090_down` DILARANG (isi hanya tinggal di kolom); pelindungnya gagal keras ORA-00904 |
| K4 | SESUDAH MENU: `T955` 0 tetapi `ADA_MENU` 1 | catat 955 manual: `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('955_menu_causeoflosslife', SYSDATE); COMMIT;`. `ADA_MENU` 0 = ulangi (b) (INSERT `NOT EXISTS`) | `955_down` (hapus hak lalu baris menu) |

Nama langkah: `090_causeofloss_life_ganti_nama`, `091_causeofloss_life_kolom`, `092_causeofloss_life_satu_tabel`,
`955_menu_causeoflosslife`.

## Jalur mundur

Hanya bila semua langkah yang dibongkar TERCATAT. Hentikan backend. DBA menjalankan berkas `_down` di
`modul/causeoflosslife/backend/migrations/` MENURUN: 955, 092, 091, 090 (ganti `{skema}` dengan POOLDATA), lalu hapus
catatan `T_MIGRASI` keempatnya; backend dari kode sebelum modul ini.

- `092_down` aman diulang (blok berpelindung): JSONDATA dibangun ulang `JSON_OBJECT('CauseofLoss' VALUE …
  ABSENT ON NULL)` (kunci huruf persis Pega; 100001 kembali TANPA kunci - view tetap membacanya NULL) + constraint IS
  JSON bernama asli. `pxObjClass`, `pyRuleHarness`, dan `"CauseofLoss":""` asli hanya dari cadangan P3.
- `091_down` mulai dengan pelindung `UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0` (ORA-00904 bila JSONDATA tidak ada).
- `090_down`: pelindung yang sama, RENAME balik berpelindung, lalu CREATE VIEW PERSIS teks DEV (terakhir; view yang
  sudah ada = ORA-00955, lewati pernyataan itu).
- JSONDATA dikembalikan **NULLABLE** - bila acuan (a) P2 kueri 4 menunjukkan `NULLABLE = N`, DBA menambah
  `ALTER TABLE POOLDATA.M_CAUSEOFLOSS_LIFE MODIFY (JSONDATA NOT NULL);` sesudah mundur.
- Prosedur `PEGA_M_CAUSEOFLOSS_LIFE` perlu dikompilasi ulang DBA. ⛔ Jangan memakai `-migrate-down` di DEV.
- Jatah K0 (premiumlistlife `050-089` / `954`) TIDAK ikut mundur - itu tabel jatah dokumen, bukan objek DB.
