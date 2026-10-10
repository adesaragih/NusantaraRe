# Langkah work owner — menyalakan Benefit (`benefitlife`, migrasi inti 942-945)

> ## ⛔ LANGKAH PERTAMA, SEBELUM APA PUN: CADANGAN P3
>
> **Dua migrasi terakhir dijalankan TANPA cadangan. Jangan terulang. Tanpa cadangan ID + JSONDATA = TIDAK BOLEH
> `-migrate`.** 944 MEMBUANG JSONDATA (satu-satunya tempat kunci `Number`, `pxObjClass`, `px*` dan bentuk JSON asli).
>
> ```powershell
> mkdir "D:\NUSARE DEV\NUSARE\cadangan-benefitlife-20261008"     # <tanggal> hari menjalankan
> cd "D:\NUSARE DEV\NUSARE\cadangan-benefitlife-20261008"
> sqlplus POOLDATA@<db> @"D:\NUSARE DEV\NUSARE\NusantaraRe\modul\benefitlife\docs\sql\benefit_p3_cadangan.sql"
> ```
>
> Hasil: `M_BENEFIT_LIFE_LENGKAP.csv` (ID + JSONDATA utuh) dan `BENEFIT_LIFE_SEBELUM.csv` (ID, BENEFIT lewat view),
> keduanya **10 baris** (cacah dicetak di akhir berkas). **Buka keduanya dan periksa sebelum lanjut.** Urutan resmi
> (P1 → P2 → P3) ada di bawah: P3 dijalankan SESUDAH penulis dihentikan (P1) dan acuan pertama (P2); kotak ini ada di
> atas supaya tidak terlewat.

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini (executor tidak terhubung ke Oracle). Keputusan work owner 08-10-2026 K1-K6 (`../MODUL.md`).
> Skema DEV: POOLDATA. Angka DEV (10 baris, ID 100001-100011 tanpa 100003, sequence last_number 12, Benefit maks 48
> byte) = PEMBANDING dari WO 08-10-2026, bukan angka tetap. Pola sama dengan
> `modul/ririsklife/docs/LANGKAH-WO-RIRISKLIFE.md`.

Keadaan awal: `M_BENEFIT_LIFE` (ID VARCHAR2(10) NOT NULL PK `SYS_C009031`, JSONDATA CLOB `ENSURE_M_BENEFIT_LIFE_JSON`) +
view `BENEFIT_LIFE` (`SELECT a.ID, a.JSONDATA.Benefit FROM M_BENEFIT_LIFE a`); sequence `M_BENEFIT_LIFE_SEQ`; prosedur
`PEGA_M_BENEFIT_LIFE`. Tidak ada pembaca lain, tidak ada grant ke skema lain (fakta WO). Objek `M_BENEFIT`, `MBENEFIT*`,
`VJ_*BENEFIT*`, `VH_*BENEFIT*`, `HCD_*BENEFIT*`, `T_BENEFIT`, `DET_GROUP_BENEFIT`, `M_PLAN_BENEFIT`, `M_PROPERTY_BENEFIT`
milik aplikasi lain - TIDAK disentuh.

**Semua SQL = berkas SQL\*Plus** di `modul/benefitlife/docs/sql/`. Jalankan sebagai BERKAS (`@"…\benefit_xxx.sql"`) dari
folder cadangan - hanya lewat berkas `SET TERMOUT OFF` berlaku. SQL\*Plus 12.2+ atau SQLcl. ⛔ CSV cadangan memuat
kunci `pxCreateOperator` / `pxCreateOpName` (identitas orang): simpan di tempat cadangan DBA, jangan disalin ke
dokumen / tiket.

| Berkas | Langkah | Sifat |
| --- | --- | --- |
| `benefit_p3_cadangan.sql` | (a) P3 cadangan ID + JSONDATA | baca-saja |
| `benefit_p1_kunci.sql` / `benefit_p1_kunci_dba.sql` | (a) P1 kunci dan sesi | baca-saja |
| `benefit_p2_acuan.sql` | (a) P2 acuan - DUA kali | baca-saja |
| `benefit_bukti_k2.sql` | (a) bukti K2: ekspresi 944 = view, 10 nilai identik | baca-saja |
| `benefit_c_verifikasi.sql` | (c) verifikasi + CSV sesudah | baca-saja |
| `benefit_d_hak.sql` | (d) hak menu (K5) | MENULIS M_LOGIN_GO_MENU |
| `benefit_keadaan.sql` | Pemulihan - penentu keadaan | baca-saja |

## (a) Prasyarat P1-P3 — WAJIB SEBELUM `-migrate`

### ⛔ P1 — hentikan penulis, cek kunci

1. **Hentikan Pega Benefit** (penyimpan `InboxBenefit` → prosedur `PEGA_M_BENEFIT_LIFE`, masih VALID sampai 944) dan
   **backend :8080** (tetap mati sampai (e)). Alasannya: 944 mengisi BENEFIT lalu membuang JSONDATA dalam pernyataan
   terpisah tanpa transaksi; baris yang disisip di antaranya kehilangan isinya dan tidak ada di cadangan.
2. `@benefit_p1_kunci.sql`: kueri 1 (`V$LOCKED_OBJECT` ⨝ `ALL_OBJECTS` ⨝ `V$SESSION`) **nol baris**; kueri 2 = sesi
   POOLDATA yang masih tersambung - pastikan tidak ada Pega / backend. ORA-00942 pada V$ = minta DBA menjalankan
   `benefit_p1_kunci_dba.sql` (`DBA_DML_LOCKS`, nol baris).

### ⛔ P2 — acuan baris, DUA kali

`@benefit_p2_acuan.sql` (ditambahkan ke `benefit_acuan.txt`): **10 baris**, ID tertinggi, sequence, sidik
`SUM(ORA_HASH(ID || '|' || Benefit))` lewat view, panjang maksimum (48), nilai > 200 byte (**nol baris**), NULLABLE asli
JSONDATA, constraint, objek / dependensi (hanya view + prosedur Pega) / grant (nol). Jalankan **sekali sebelum P3** dan
**sekali lagi tepat sebelum `-migrate`** ((c)0). Satu angka di luar harapan, atau dua keluaran berbeda = BERHENTI.

### ⛔ P3 — cadangan ID + JSONDATA

Lihat kotak di atas berkas ini: `@benefit_p3_cadangan.sql` dari folder `cadangan-benefitlife-<tanggal>\`. Dua CSV, 10
baris masing-masing. **Tidak ada cadangan = tidak boleh `-migrate`.**

### Bukti K2 — sebelum `-migrate`

`@benefit_bukti_k2.sql` → `benefit_bukti_k2.txt`: ekspresi pengisian 944 (`m.JSONDATA.Benefit`, notasi titik PERSIS
teks view) sebagai SELECT, dibandingkan dengan view. Harus: A `N` = `N_EKSPRESI` = `N_VIEW` = `N_JSON_VALUE` =
`N_ADA_KUNCI` = `N_SAMA` = **10**; B sepuluh baris `SAMA`; C dan D **nol baris**; E `MAKS_BYTE` ≤ 200 (48); F
`N_NUMBER` 5 / `N_NUMBER_BEDA_ID` 0. G (cacah Benefit yang belum huruf besar) = bahan pertanyaan terbuka T1 di
`PARITAS-LAYAR-DAN-AKSI.md`. **Kirim berkas keluarannya ke executor / laporan** (butir yang menunggu WO). Satu angka
berbeda = BERHENTI, jangan `-migrate`.

## (b) `-migrate`

```powershell
. .\muat-env.ps1                 # dot-source; ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate        # = target Makefile `migrate`
```

Log yang benar: `dijalankan:` 942, 943, 944, 945 berurutan. 942 = DROP VIEW (berpelindung ALL_VIEWS) + RENAME
(berpelindung tabel sumber); 943 = satu ALTER … ADD; 944 = isi BENEFIT dari JSONDATA → pemeriksaan K2 → buang JSONDATA;
945 = baris menu.

**Bila 944 berhenti dengan `ORA-01407: cannot update ("POOLDATA"."BENEFIT_LIFE"."ID") to NULL`**: itu pemeriksaan K2
yang BEKERJA - ada baris yang JSON-nya memuat Benefit tetapi BENEFIT-nya NULL / berbeda. Tidak ada yang berubah
(JSONDATA utuh, 944 tidak tercatat). JANGAN mengulang; laporkan ke WO dengan keluaran `benefit_bukti_k2.sql` kueri C.

## (c) Verifikasi — baca-saja

`@benefit_c_verifikasi.sql` → `benefit_verifikasi.txt`:

1. `BENEFIT_LIFE` = `ID` VARCHAR2(10) NOT NULL, `BENEFIT` VARCHAR2(200); nol JSONDATA; `ALL_PARTIAL_DROP_TABS` nol baris;
2. `BENEFIT_LIFE` = **TABLE**; `M_BENEFIT_LIFE` **tidak ada**; **nol VIEW** `BENEFIT_LIFE`; `M_BENEFIT_LIFE_SEQ` tetap;
   `PEGA_M_BENEFIT_LIFE` INVALID (diharapkan, K3 - jangan dikompilasi ulang);
3. **10 baris**, ID tertinggi sama, **sidik `H_BARIS` SAMA dengan acuan** (a) P2 kueri 2;
4. PK `SYS_C009031` ada; `ENSURE_M_BENEFIT_LIFE_JSON` hilang;
5. `T_MIGRASI` 942, 943, 944, 945 (empat baris); `M_NAV_MENU` `benefitlife` 'Benefit' MASTER TREATY **12** '1' dengan
   STATUS_AKTIF sama seperti baris `ririsklife`.

MINUS dua arah (CSV sesudah ditulis berkas verifikasi):

```powershell
$a = Get-Content 'BENEFIT_LIFE_SEBELUM.csv' -Encoding UTF8; $b = Get-Content 'BENEFIT_LIFE_SESUDAH.csv' -Encoding UTF8
"$($a.Count) baris sebelum, $($b.Count) sesudah"
Compare-Object $a $b -CaseSensitive | Group-Object SideIndicator | Select-Object Name, Count   # kosong = sama
```

`<=` = sebelum MINUS sesudah, `=>` = sesudah MINUS sebelum; kosong = MINUS dua arah nol (`-CaseSensitive` wajib).

## (d) Hak menu (K5)

`@benefit_d_hak.sql`: setiap akun yang memegang `riratelife`, `ricommlife`, atau `ririsklife` mendapat `benefitlife` -
`PENUH` bila salah satunya PENUH, selain itu `LIHAT`; aman diulang (`NOT EXISTS`). Periksa `AKUN_BENEFITLIFE` =
`AKUN_SUMBER`. Akun lain: centang menu "Benefit" di Kelola User.

## (e) Restart backend dan periksa layar

```powershell
go run ./cmd/api                 # = target Makefile `run-api`, jendela yang sudah memuat env
```

- **Benefit** (MASTER TREATY, urutan 12, sesudah R/I Risk): judul "INSURANCE BENEFIT"; grid ID / Benefit 10 baris, ID
  menurun (100011 di atas); saring; Save satu Benefit uji (ID baru `100012` bila sequence belum dipakai), Edit Benefit
  itu, Cancel; teks diketik huruf kecil tersimpan huruf besar.
- Akun View only: grid tampil tanpa form dan tanpa Edit.

Pega Benefit TIDAK dinyalakan lagi untuk menyimpan (prosedur INVALID - diterima WO K3).

## Pemulihan — 942-945 gagal di tengah

Pelari tanpa transaksi. Backend dan Pega tetap mati sampai pulih. `@benefit_keadaan.sql` memberi `ADA_VIEW`,
`ADA_SUMBER`, `ADA_TARGET`, `ADA_BENEFIT`, `ADA_JSONDATA`, `SETENGAH_TERBUANG`, `T942`-`T945`, `ADA_MENU`. Ikuti SATU
baris:

| # | Keadaan | Maju (tindakan) | Mundur boleh? |
| --- | --- | --- | --- |
| K0 | mati SEBELUM RENAME: `ADA_SUMBER` 1, `ADA_TARGET` 0 (`ADA_VIEW` 1 atau 0 - view mungkin sudah dibuang), `T942` 0 | 942 aman diulang (DROP VIEW hanya bila VIEW; RENAME hanya bila sumber ada). Ulangi (b) | tidak perlu; bila tidak jadi maju dan `ADA_VIEW` 0: jalankan pernyataan CREATE VIEW terakhir `942_benefit_life_ganti_nama_down.sql` saja |
| K1 | SESUDAH RENAME: `ADA_TARGET` 1, `ADA_SUMBER` 0, `ADA_BENEFIT` 0 | bila `T942` 0: kedua blok 942 melewati diri - ulangi (b) (942 tercatat, lalu 943 dst.) | YA: `942_down` (pelindungnya lolos selama JSONDATA ada) + hapus catatan 942 bila ada |
| K2 | SESUDAH ADD: `ADA_BENEFIT` 1, `ADA_JSONDATA` 1, `T943` 0 **atau** 944 berhenti sebelum DROP (`T944` 0, `ADA_JSONDATA` 1) | `T943` 0: jangan ulangi 943 (ORA-01430) - `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('943_benefit_life_kolom', SYSDATE); COMMIT;` lalu ulangi (b). 944 berhenti ORA-12899 / **ORA-01407 (pemeriksaan K2)** = BERHENTI, laporkan ke WO; galat lain: ulangi (b) (isi menulis nilai yang sama) | YA (JSONDATA utuh): `943_down` lalu `942_down`, hapus catatannya |
| K3 | SESUDAH DROP JSONDATA: `ADA_JSONDATA` 0 (`SETENGAH_TERBUANG` 1 = `ALTER TABLE POOLDATA.BENEFIT_LIFE DROP COLUMNS CONTINUE;` dulu), `T944` 0 | **MAJU WAJIB diselesaikan**: ulangi (b) (ketiga blok 944 melewati diri, 944 tercatat, lalu 945). Pastikan (c) | **TIDAK.** `943_down` / `942_down` DILARANG (isi hanya tinggal di kolom); pelindungnya gagal keras ORA-00904 |
| K4 | SESUDAH MENU: `T945` 0 tetapi `ADA_MENU` 1 | catat 945 manual: `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('945_m_nav_menu_benefitlife', SYSDATE); COMMIT;`. `ADA_MENU` 0 = ulangi (b) (INSERT `NOT EXISTS`) | `945_down` (hapus hak lalu baris menu) |

Nama langkah: `942_benefit_life_ganti_nama`, `943_benefit_life_kolom`, `944_benefit_life_satu_tabel`,
`945_m_nav_menu_benefitlife`.

## Jalur mundur

Hanya bila semua langkah yang dibongkar TERCATAT. Hentikan backend. DBA menjalankan berkas `_down` di
`inti/backend/migrations/` MENURUN: 945, 944, 943, 942 (ganti `{skema}` dengan POOLDATA), lalu hapus catatan
`T_MIGRASI` keempatnya; backend dari kode sebelum modul ini.

- `944_down` aman diulang (blok berpelindung): JSONDATA dibangun ulang `JSON_OBJECT('Benefit' VALUE …)` (kunci huruf
  persis Pega) + constraint IS JSON bernama asli. Kunci `Number`, `pxObjClass`, `px*` hanya dari cadangan P3.
- `943_down` mulai dengan pelindung `UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0` (ORA-00904 bila JSONDATA tidak ada).
- `942_down`: pelindung yang sama, RENAME balik berpelindung, lalu CREATE VIEW PERSIS teks DEV (terakhir; view yang
  sudah ada = ORA-00955, lewati pernyataan itu).
- JSONDATA dikembalikan **NULLABLE** - bila acuan (a) P2 kueri 4 menunjukkan `NULLABLE = N`, DBA menambah
  `ALTER TABLE POOLDATA.M_BENEFIT_LIFE MODIFY (JSONDATA NOT NULL);` sesudah mundur.
- Prosedur `PEGA_M_BENEFIT_LIFE` perlu dikompilasi ulang DBA. ⛔ Jangan memakai `-migrate-down` di DEV.
