# Langkah work owner — menyalakan Disease Life (`diseaselife`, migrasi modul 080-081 + slot 951)

> ## ⛔ LANGKAH PERTAMA, SEBELUM APA PUN: CADANGAN P3
>
> **Tanpa cadangan = TIDAK BOLEH D2 (hapus baris uji) dan TIDAK BOLEH `-migrate`.** Seluruh `DISEASE_LIFE`
> (**97.586 baris**), `M_DISEASE_LIFE` (3 baris JSON lama), dan nilai `M_DISEASE_LIFE_SEQ`.
>
> ```powershell
> mkdir "D:\NUSARE DEV\NUSARE\cadangan-diseaselife-20261008"     # <tanggal> hari menjalankan
> cd "D:\NUSARE DEV\NUSARE\cadangan-diseaselife-20261008"
> sqlplus POOLDATA@<db> @"D:\NUSARE DEV\NUSARE\NusantaraRe\modul\diseaselife\docs\sql\dis_p3_cadangan.sql"
> ```
>
> Hasil: `DISEASE_LIFE_LENGKAP.csv` (**97.586 baris**), `M_DISEASE_LIFE_LENGKAP.csv` (3), `M_DISEASE_LIFE_SEQ.txt`
> (LAST_NUMBER 2052). **Buka ketiganya dan periksa sebelum lanjut.** Urutan resmi (P1 → P2 → P3 → bukti → D2 → P2
> ulang → `-migrate`) ada di bawah; kotak ini ada di atas supaya tidak terlewat.

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini (executor tidak terhubung ke Oracle; pembacaan DB ditolak pengaman sesi). Keputusan work owner
> 08-10-2026 K0, D1-D4, K5 (`../MODUL.md`). Skema DEV: POOLDATA. Angka DEV (97.586 baris, ID tertinggi 197585, ID kembar
> 102051, `M_DISEASE_LIFE_SEQ` 2052) = PEMBANDING dari WO 08-10-2026, bukan angka tetap.

Keadaan awal: `DISEASE_LIFE` (ID VARCHAR2(100) NULLABLE, ICD_CODE VARCHAR2(100), DISEASE VARCHAR2(1000)) **tanpa PK,
tanpa indeks**; prosedur `PEGA_DISEASE_LIFE` (penulis Pega, VALID); objek lama `M_DISEASE_LIFE` / `PEGA_M_DISEASE_LIFE`
/ `M_DISEASE_LIFE_SEQ` (TIDAK disentuh). Pembaca lain: `claimlife` (pencarian diagnosa - nama dan kolom tetap).

**Semua SQL = berkas SQL\*Plus** di `modul/diseaselife/docs/sql/`. Jalankan sebagai BERKAS (`@"…\dis_xxx.sql"`) dari
folder cadangan - hanya lewat berkas `SET TERMOUT OFF` berlaku. SQL\*Plus 12.2+ atau SQLcl. ⛔ CSV cadangan dapat memuat
kunci identitas operator: simpan di tempat cadangan DBA, jangan disalin ke dokumen / tiket.

| Berkas | Langkah | Sifat |
| --- | --- | --- |
| `dis_p3_cadangan.sql` | (a) P3 cadangan seluruh DISEASE_LIFE + M_DISEASE_LIFE + nilai M_DISEASE_LIFE_SEQ | baca-saja |
| `dis_p1_kunci.sql` / `dis_p1_kunci_dba.sql` | (a) P1 kunci dan sesi | baca-saja |
| `dis_p2_acuan.sql` | (a) P2 acuan - TIGA kali | baca-saja |
| `dis_bukti.sql` | (a) bukti D1-D3 (ID tertinggi, rumus Pega bertabrakan, ID kembar / NULL, ICD kembar, huruf) | baca-saja |
| `disease_hapus_baris_uji.sql` | (a) **D2** hapus baris uji 102051 `TEST123 / Sakit` | **MENULIS** (satu baris, COMMIT) |
| `dis_c_verifikasi.sql` | (c) verifikasi + CSV sesudah | baca-saja |
| `dis_d_hak.sql` | (d) hak menu (K5) | MENULIS M_LOGIN_GO_MENU |
| `dis_keadaan.sql` | Pemulihan - penentu keadaan | baca-saja |

## (a) Prasyarat — WAJIB SEBELUM `-migrate`

### ⛔ P3 — cadangan (PALING PENTING)

Lihat kotak di atas berkas ini: `@dis_p3_cadangan.sql` dari folder `cadangan-diseaselife-<tanggal>\`. **Tidak ada
cadangan = tidak boleh D2, tidak boleh `-migrate`.**

### ⛔ P1 — hentikan penulis, cek kunci

1. **Hentikan Pega Disease** (penyimpan `InboxDisease` → prosedur `PEGA_DISEASE_LIFE`) dan **backend :8080** (tetap
   mati sampai (e)). Alasannya: 080 membaca ID tertinggi lalu membuat sequence dalam satu blok, 081 memvalidasi seluruh
   ID; penulis yang hidup di antaranya dapat menyisip ID kembar / di atas awal sequence.
2. `@dis_p1_kunci.sql`: kueri 1 **nol baris**; kueri 2 = sesi POOLDATA yang masih tersambung - pastikan tidak ada Pega /
   backend. ORA-00942 pada V$ = minta DBA menjalankan `dis_p1_kunci_dba.sql` (`DBA_DML_LOCKS`, nol baris).

### ⛔ P2 — acuan baris, TIGA kali

`@dis_p2_acuan.sql` (ditambahkan ke `dis_acuan.txt`): **97.586 baris**, ID unik 97.585, ID NULL 0, ID angka tertinggi
197585, ICD unik 97.586, sidik `SUM(ORA_HASH(ID || '|' || ICD_CODE || '|' || DISEASE))`, ID kembar = **satu kelompok
102051 x2**, nol constraint, nol indeks, `M_DISEASE_LIFE_SEQ` 2052, `PEGA_DISEASE_LIFE` VALID, nol grant, dan **nol baris
T_MIGRASI bernomor 080-084 / 951** (K0). Jalankan (1) sebelum P3, (2) sesudah P3 sebelum D2 - keduanya SAMA; (3) sesudah
D2 tepat sebelum `-migrate` - cacah **97.585**, ID kembar **nol**. Angka di luar harapan = BERHENTI.

### Bukti — sebelum D2

`@dis_bukti.sql` → `dis_bukti.txt`: A awal sequence **197586**; B rumus Pega `'1' || LPAD(…)` untuk nomor 2051-2056 →
SEMUA sudah terpakai (alasan D1.1); C ID kembar = 102051 (`C718 …` impor + `TEST123 / Sakit` uji); D ID NULL **0**; E ICD
kembar **0**; F huruf kecil hanya baris uji; G baris uji tepat **1**, pemakai teks `TEST123` di tabel klaim **0**; H
panjang maks 6 / 7 / 290. **Kirim keluarannya ke executor / laporan** (butir menunggu WO). Satu angka berbeda = BERHENTI.

### ⛔ D2 — hapus baris uji 102051 (SESUDAH P3, SEBELUM `-migrate`)

```powershell
sqlplus POOLDATA@<db> @"D:\NUSARE DEV\NUSARE\NusantaraRe\modul\diseaselife\docs\sql\disease_hapus_baris_uji.sql"
```

Berkas: cadangkan baris itu (`DISEASE_LIFE_BARIS_UJI.csv`), lalu blok PL/SQL `DELETE … WHERE ID = '102051' AND
ICD_CODE = 'TEST123' AND DISEASE = 'Sakit'`; bila jumlahnya **bukan 1** → ROLLBACK dan berhenti ORA-20001 (nol baris
berubah); tepat 1 → COMMIT. Periksa: baris uji 0, ID 102051 tinggal **1** (baris impor), cacah **97.585**, nol ID
kembar. Lalu **P2 ketiga**. `M_DISEASE_LIFE`, `PEGA_M_DISEASE_LIFE`, `M_DISEASE_LIFE_SEQ` TIDAK disentuh.

## (b) `-migrate`

```powershell
. .\muat-env.ps1                 # dot-source; ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate        # = target Makefile `migrate`
```

Log yang benar: `dijalankan:` **080, 081, 951** (080-081 tampil di antara langkah modul lain bernomor 0xx bila ada yang
baru; 951 sesudah 949), **tanpa** "dilewati, objeknya sudah ada" untuk SEQ_DISEASE_LIFE. 080 = sequence dari ID
tertinggi; 081 = PK; 951 = baris menu.

**Bila 081 berhenti** - tidak ada yang berubah (DDL atomik: nol PK, nol indeks), 081 tidak tercatat:
- `ORA-02437: cannot validate (POOLDATA.PK_DISEASE_LIFE) - primary key violated` = masih ada **ID kembar** (D2 belum
  dijalankan, atau penulis Pega masih hidup). JANGAN menghapus sendiri; jalankan `dis_bukti.sql` kueri C, laporkan ke WO.
- `ORA-01449: column contains NULL values; cannot alter to NOT NULL` = ada **ID NULL** (kueri D). Laporkan ke WO.
- `ORA-02260: table can have only one primary key` = sudah ada PK lain bernama berbeda. Laporkan ke WO.

## (c) Verifikasi — baca-saja

`@dis_c_verifikasi.sql` → `dis_verifikasi.txt`:

1. Kolom DISEASE_LIFE TIDAK berubah (nama, tipe, lebar); ID kini `NULLABLE = N` (karena PK);
2. `PK_DISEASE_LIFE` P / ENABLED / VALIDATED + indeks unik bernama sama;
3. `SEQ_DISEASE_LIFE`: `LAST_NUMBER` (NOCACHE = nomor berikutnya) **> ID angka tertinggi** (`SEQ_DI_ATAS_ID = OK`,
   pembanding 197586) - ⛔ **JANGAN memanggil NEXTVAL**; `M_DISEASE_LIFE_SEQ` tetap 2052;
4. **97.585 baris**, ID kembar 0, ID NULL 0, **sidik `H_BARIS` SAMA dengan acuan KETIGA** (a) P2 kueri 2;
5. `PEGA_DISEASE_LIFE` tetap VALID; `M_DISEASE_LIFE` tetap 3 baris;
6. `T_MIGRASI` memuat `080_seq_disease_life`, `081_disease_life_pk`, `951_menu_diseaselife`; `M_NAV_MENU` `diseaselife`
   'Disease Life' MASTER TREATY **15** '1' dengan STATUS_AKTIF sama seperti baris `causeoflosslife`;
7. kueri pencarian diagnosa claimlife (apa adanya, `A00`) tetap memberi baris.

MINUS dua arah (CSV sebelum = cadangan P3 dikurangi baris uji):

```powershell
$a = Get-Content 'DISEASE_LIFE_LENGKAP.csv' -Encoding UTF8 | Where-Object { $_ -notmatch '^"102051","TEST123","Sakit"$' }
$b = Get-Content 'DISEASE_LIFE_SESUDAH.csv' -Encoding UTF8
"$($a.Count) baris sebelum (tanpa baris uji), $($b.Count) sesudah"
Compare-Object $a $b -CaseSensitive | Group-Object SideIndicator | Select-Object Name, Count   # kosong = sama
```

## (d) Hak menu (K5)

`@dis_d_hak.sql`: setiap akun yang memegang `causeoflosslife` mendapat `diseaselife` dengan HAK yang SAMA (LIHAT /
PENUH); aman diulang (`NOT EXISTS`). Periksa `AKUN_DISEASELIFE` = `AKUN_CAUSEOFLOSSLIFE` dan kueri ketiga nol baris. Akun
lain: centang menu "Disease Life" di Kelola User.

## (e) Restart backend dan periksa layar

```powershell
go run ./cmd/api                 # = target Makefile `run-api`, jendela yang sudah memuat env
```

- **Disease Life** (MASTER TREATY, urutan 15, sesudah Cause Of Loss Life): judul "DISEASE"; grid ID / ICD Code /
  Disease 10 baris, **ID menurun** (197585 di atas), "97585 diseases", halaman berpindah cepat (paging di server);
  saring ICD Code `A00` / Disease `kolera`; klik kepala ICD Code mengurut; Save satu Disease uji (ID baru **197586**,
  ICD / Disease menjadi huruf besar), Edit itu, Cancel; ICD Code kembar (tanpa beda huruf) ditolak.
- Akun View only: grid tampil tanpa form dan tanpa Edit.
- **Claim Life**: pencarian diagnosa tetap menampilkan hasil.

Pega Disease boleh dinyalakan lagi (prosedur `PEGA_DISEASE_LIFE` tetap VALID) - ⚠️ rumus ID-nya tetap bertabrakan dan
kini ditolak PK (ORA-00001 di Pega). Lihat pertanyaan terbuka `PR-DISEASELIFE.md`.

## Pemulihan — D2 / 080-081 / 951 gagal di tengah

Pelari tanpa transaksi. Backend dan Pega tetap mati sampai pulih. `@dis_keadaan.sql` memberi `ADA_SEQ`, `ADA_PK`,
`T080`, `T081`, `T951`, `ADA_MENU`, `N`, `N_ID_NULL`, `N_ID_KEMBAR`, `ADA_BARIS_UJI`. Ikuti SATU baris:

| # | Keadaan | Maju (tindakan) | Mundur boleh? |
| --- | --- | --- | --- |
| K0 | D2 belum / gagal: `ADA_BARIS_UJI` 1, `N_ID_KEMBAR` 1 | jalankan D2 (berkas aman diulang: bukan 1 baris = ROLLBACK). Lalu P2 ketiga, (b) | tidak perlu (nol perubahan) |
| K1 | 080 gagal sebelum sequence: `ADA_SEQ` 0, `T080` 0 | ulangi (b) (blok bertanya ALL_SEQUENCES dulu) | tidak perlu |
| K2 | sequence berdiri, 080 tidak tercatat: `ADA_SEQ` 1, `T080` 0 | ulangi (b): blok melewati diri, 080 tercatat. ⚠️ Bila Pega sempat menulis ID di atas awal sequence: aplikasi menolak tabrakan (409) - DBA dapat `DROP SEQUENCE POOLDATA.SEQ_DISEASE_LIFE;` lalu (b) | YA: `080_down` (DROP SEQUENCE) + hapus catatan 080 bila ada |
| K3 | 081 berhenti ORA-02437 / ORA-01449 / ORA-02260: `ADA_PK` 0, `T081` 0 | BERHENTI, laporkan ke WO dengan `dis_bukti.sql` kueri C / D (data tidak dihapus migrasi). Sesudah WO membereskan: ulangi (b) | tidak perlu (nol perubahan) |
| K4 | PK berdiri, 081 tidak tercatat: `ADA_PK` 1, `T081` 0 | ulangi (b) (blok melewati diri, 081 tercatat) | YA: `081_down` (DROP CONSTRAINT, ID kembali NULLABLE) |
| K5 | menu: `T951` 0 tetapi `ADA_MENU` 1 | catat 951 manual: `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('951_menu_diseaselife', SYSDATE); COMMIT;`. `ADA_MENU` 0 = ulangi (b) | `951_down` (hapus hak lalu baris menu) |

Nama langkah: `080_seq_disease_life`, `081_disease_life_pk`, `951_menu_diseaselife`.

## Jalur mundur

Hanya bila langkah yang dibongkar TERCATAT. Hentikan backend. DBA menjalankan berkas `_down` di
`modul/diseaselife/backend/migrations/` MENURUN: 951, 081, 080 (ganti `{skema}` dengan POOLDATA), lalu hapus catatan
`T_MIGRASI` ketiganya; backend dari kode sebelum modul ini.

- `081_down` (blok berpelindung, aman diulang): DROP CONSTRAINT PK_DISEASE_LIFE - indeksnya ikut; NOT NULL implisit
  ikut terbuang (periksa `NULLABLE = Y`).
- `080_down`: DROP SEQUENCE SEQ_DISEASE_LIFE (baris yang sudah memakai nomornya tetap).
- Baris uji 102051 (D2) TIDAK dikembalikan jalur mundur - hanya dari `DISEASE_LIFE_BARIS_UJI.csv` bila WO memintanya.
- ⛔ Jangan memakai `-migrate-down` di DEV.
- Jatah K0 (premiumlistlife `050-079`, claimlife slot `950`) TIDAK ikut mundur - itu tabel jatah dokumen, bukan objek DB.
