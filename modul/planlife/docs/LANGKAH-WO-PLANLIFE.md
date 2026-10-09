# Langkah work owner — menyalakan Plan (`planlife`, migrasi inti 946-949)

> ## ⛔ LANGKAH PERTAMA, SEBELUM APA PUN: CADANGAN P3
>
> **Tanpa cadangan ID + JSONDATA = TIDAK BOLEH `-migrate`.** 948 MEMBUANG JSONDATA (satu-satunya tempat kunci `px*`,
> `pyRuleHarness`, dan bentuk JSON asli).
>
> ```powershell
> mkdir "D:\NUSARE DEV\NUSARE\cadangan-planlife-20261008"     # <tanggal> hari menjalankan
> cd "D:\NUSARE DEV\NUSARE\cadangan-planlife-20261008"
> sqlplus POOLDATA@<db> @"D:\NUSARE DEV\NUSARE\NusantaraRe\modul\planlife\docs\sql\plan_p3_cadangan.sql"
> ```
>
> Hasil: `M_PRODUCT_TYPE_LIFE_LENGKAP.csv` (ID + JSONDATA utuh) dan `PRODUCT_TYPE_LIFE_SEBELUM.csv` (6 kolom view),
> keduanya **32 baris**. **Buka keduanya dan periksa sebelum lanjut.** Urutan resmi (P1 → P2 → P3) ada di bawah; kotak
> ini di atas supaya tidak terlewat.

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini (executor tidak terhubung ke Oracle). Keputusan work owner 08-10-2026 K1-K7 (`../MODUL.md`).
> Angka DEV (32 baris, sequence last_number 44, CoverName 34 / Business 29 / Benefit 48 byte) = PEMBANDING, bukan angka
> tetap. Pola `modul/benefitlife/docs/LANGKAH-WO-BENEFITLIFE.md`.

Keadaan awal: `M_PRODUCT_TYPE_LIFE` (ID VARCHAR2(6) NULLABLE TANPA PK, JSONDATA CLOB `ENSURE_M_PRODUCT_TYPE_LIFE`) + view
`PRODUCT_TYPE_LIFE`; sequence `M_PRODUCT_TYPE_LIFE_SEQ`; prosedur `PEGA_M_PRODUCT_TYPE_LIFE`. Pembaca lain:
`masterproductnamelife` (autocomplete Plan). **Jangan disentuh:** `PLAN_LIFE_SUMMARY`, `PEGA_M_PLAN_LIFE_SUMMARY`,
`BUSINESS`, `M_BUSINESS`, `BENEFIT_LIFE`.

**Semua SQL = berkas SQL\*Plus** di `modul/planlife/docs/sql/` (`@"…\plan_xxx.sql"`, dari folder cadangan). ⛔ CSV
cadangan memuat kunci `px*` (identitas operator): simpan di tempat cadangan DBA, jangan disalin ke dokumen / tiket.

| Berkas | Langkah | Sifat |
| --- | --- | --- |
| `plan_p3_cadangan.sql` | (a) P3 cadangan ID + JSONDATA + view | baca-saja |
| `plan_p1_kunci.sql` / `plan_p1_kunci_dba.sql` | (a) P1 kunci (biasa / DBA) | baca-saja |
| `plan_p2_acuan.sql` | (a) P2 acuan - DUA kali | baca-saja |
| `plan_bukti_k1.sql` | (a) bukti K2: 5 ekspresi = view, ID = JSONDATA.ID, lebar, master | baca-saja |
| `plan_c_verifikasi.sql` | (c) verifikasi + CSV sesudah | baca-saja |
| `plan_d_hak.sql` | (d) hak menu (K7) | MENULIS M_LOGIN_GO_MENU |
| `plan_keadaan.sql` | Pemulihan - penentu keadaan | baca-saja |

## (a) Prasyarat — WAJIB SEBELUM `-migrate`

### ⛔ P1 — hentikan penulis, cek kunci

1. **Hentikan backend :8080 dan Pega Plan** (penyimpan `InboxProductType` → `PEGA_M_PRODUCT_TYPE_LIFE`, VALID sampai
   948). Backend tetap mati sampai (e).
2. `@plan_p1_kunci.sql` → kueri 1 **nol baris**. Akun POOLDATA DEV tidak punya hak `V$LOCKED_OBJECT` (ORA-00942): minta
   DBA menjalankan `@plan_p1_kunci_dba.sql` (`DBA_DML_LOCKS` + `V$LOCKED_OBJECT`, keduanya nol baris).

### ⛔ P2 — acuan, DUA kali

`@plan_p2_acuan.sql` (ke `plan_acuan.txt`): **32 baris**, `N_ID` = `N_ID_UNIK` = 32, sequence, sidik `H_BARIS` lewat
view, kueri 3 **nol baris**, NULLABLE asli, constraint (tanpa PK), objek / dependensi / grant. Sekali sebelum P3, sekali
lagi tepat sebelum `-migrate` ((c)0). Beda = ada penulis hidup - ulangi P1-P3.

### ⛔ P3 — cadangan (kotak di atas)

### Bukti K2 — sebelum `-migrate`

`@plan_bukti_k1.sql` → `plan_bukti_k1.txt`: A setiap `N_*` = **32**; B `SAMA_SEMUA` = **32**; C, F, G **nol baris**; D
maksimum ≤ lebar (ID ≤ 6, nama ≤ 200, ID master ≤ 10); E lebar `BUSINESS.ID` / `BENEFIT_LIFE.ID` (bila `BUSINESS.ID`
lebih lebar dari 10 dan ada nilai > 10: BERHENTI, laporkan); `N_BUSINESS_009` = 21. **Kirim keluarannya ke executor /
laporan** (menunggu WO). Satu angka berbeda = BERHENTI.

## (b) `-migrate`

```powershell
. .\muat-env.ps1                 # ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate
```

Log yang benar: `dijalankan:` 946, 947, 948, 949 berurutan, **tanpa** "dilewati, objeknya sudah ada".

**Bila 948 berhenti dengan `ORA-12899: value too large for column "POOLDATA"."PRODUCT_TYPE_LIFE"."COVERNAME" (actual:
201, maximum: 200)`**: itu pemeriksaan K1.4 / K2 yang BEKERJA (ID NULL / ganda / ≠ JSONDATA.ID, atau isi kolom ≠ JSON).
Tidak ada yang berubah (JSONDATA utuh, PK belum, 948 tidak tercatat). JANGAN mengulang; laporkan dengan keluaran
`plan_bukti_k1.sql` kueri C. ORA-12899 dengan nilai lain dari 201 = isi JSON melewati lebar kolom (kueri D).

## (c) Verifikasi — baca-saja

`@plan_c_verifikasi.sql` → `plan_verifikasi.txt`:

1. kolom `ID` VARCHAR2(6) **NULLABLE = N**, `COVERNAME` 200, `BUSINESS` 200, `BUSINESSID` 10, `BENEFIT` 200,
   `BENEFITID` 10, urutan view lama; nol JSONDATA; nol kolom setengah terbuang;
2. `PRODUCT_TYPE_LIFE` = **TABLE**; `M_PRODUCT_TYPE_LIFE` **tidak ada**; **nol VIEW**; sequence tetap; prosedur Pega
   INVALID (K3);
3. **32 baris**, ID unik, **sidik `H_BARIS` SAMA dengan acuan** P2 kueri 2;
4. `PK_PRODUCT_TYPE_LIFE` (P, ENABLED) ada; `ENSURE_M_PRODUCT_TYPE_LIFE` hilang;
5. `T_MIGRASI` 946-949 (empat baris); `M_NAV_MENU` `planlife` 'Plan' MASTER TREATY **13** '1', STATUS_AKTIF sama
   dengan `benefitlife`;
6. kueri pembaca Product Name Life = 32 baris.

MINUS dua arah:

```powershell
$a = Get-Content 'PRODUCT_TYPE_LIFE_SEBELUM.csv' -Encoding UTF8; $b = Get-Content 'PRODUCT_TYPE_LIFE_SESUDAH.csv' -Encoding UTF8
"$($a.Count) baris sebelum, $($b.Count) sesudah"
Compare-Object $a $b -CaseSensitive | Group-Object SideIndicator | Select-Object Name, Count   # kosong = sama
```

## (d) Hak menu (K7)

`@plan_d_hak.sql`: setiap pemegang `riratelife`, `ricommlife`, `ririsklife`, atau `benefitlife` mendapat `planlife`
(PENUH bila salah satunya PENUH, selain itu LIHAT); aman diulang. `AKUN_PLANLIFE` = `AKUN_SUMBER`.

## (e) Restart backend dan periksa layar

```powershell
go run ./cmd/api
```

- **Plan** (MASTER TREATY, urutan 13): grid Plan Name / Business / Benefit 10 baris; autocomplete Business (OLDID / Note
  / ID, 21 pilihan) dan Benefit; Save satu plan uji (ID baru `100044` bila sequence belum dipakai), Edit, New; teks yang
  tidak ada di daftar ditolak berkalimat.
- **Product Name Life**: autocomplete `Plan Name` tetap menampilkan dan memilih Plan.
- Akun View only: grid tanpa form dan Edit.

## Pemulihan — 946-949 gagal di tengah

Pelari tanpa transaksi. Backend dan Pega tetap mati sampai pulih. `@plan_keadaan.sql` memberi `ADA_VIEW`, `ADA_SUMBER`,
`ADA_TARGET`, `KOLOM_BARU`, `ADA_JSONDATA`, `ADA_PK`, `SETENGAH_TERBUANG`, `T946`-`T949`, `ADA_MENU`.

| # | Keadaan | Maju (tindakan) | Mundur boleh? |
| --- | --- | --- | --- |
| K0 | mati SEBELUM RENAME: `ADA_SUMBER` 1, `ADA_TARGET` 0 (`ADA_VIEW` 1 / 0), `T946` 0 | 946 aman diulang. Ulangi (b) | tidak perlu; bila batal maju dan `ADA_VIEW` 0: jalankan CREATE VIEW terakhir `946_product_type_life_ganti_nama_down.sql` saja |
| K1 | SESUDAH RENAME: `ADA_TARGET` 1, `ADA_SUMBER` 0, `KOLOM_BARU` 0 | ulangi (b) (blok 946 melewati diri) | YA: `946_down` + hapus catatan 946 bila ada |
| K2 | SESUDAH ADD: `KOLOM_BARU` 5, `ADA_JSONDATA` 1, `T947` 0 **atau** 948 berhenti sebelum DROP (`T948` 0, `ADA_JSONDATA` 1, `ADA_PK` 0 / 1) | `T947` 0: jangan ulangi 947 (ORA-01430) - `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('947_product_type_life_kolom', SYSDATE); COMMIT;` lalu ulangi (b). 948 berhenti **ORA-12899 actual 201 (pemeriksaan) / ORA-02437 / ORA-01449 (PK)** = BERHENTI, laporkan; galat lain: ulangi (b) | YA (JSONDATA utuh): `948_down` bila `ADA_PK` 1, lalu `947_down`, `946_down` |
| K3 | SESUDAH DROP JSONDATA: `ADA_JSONDATA` 0 (`SETENGAH_TERBUANG` 1 = `ALTER TABLE POOLDATA.PRODUCT_TYPE_LIFE DROP COLUMNS CONTINUE;` dulu), `T948` 0 | **MAJU WAJIB**: ulangi (b) (blok 948 melewati diri; PK berpelindung). Pastikan (c) | **TIDAK** - `947_down` / `946_down` DILARANG; pelindungnya gagal ORA-00904 |
| K4 | SESUDAH MENU: `T949` 0, `ADA_MENU` 1 | `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('949_m_nav_menu_planlife', SYSDATE); COMMIT;`. `ADA_MENU` 0 = ulangi (b) | `949_down` |

## Jalur mundur

Hanya bila semua langkah yang dibongkar TERCATAT. Hentikan backend. DBA menjalankan `_down` MENURUN: 949, 948, 947, 946
(ganti `{skema}` dengan POOLDATA), lalu hapus catatan `T_MIGRASI`. `948_down` membuang PK (ID kembali NULLABLE - periksa
`ALL_TAB_COLUMNS.NULLABLE = 'Y'`), membangun ulang JSONDATA enam kunci huruf Pega + constraint IS JSON bernama asli;
`947_down` / `946_down` berpelindung ORA-00904; `946_down` membuat ulang view PERSIS teks DEV. Prosedur Pega perlu
dikompilasi ulang DBA. ⛔ Jangan `-migrate-down` di DEV.
