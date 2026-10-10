# Langkah work owner — R/I Comm Life SATU tabel per jenis data (RALAT R1, migrasi inti 931-934)

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini. Keputusan work owner 08-10-2026 (`modul/ricommlife/MODUL.md`, RALAT R1). Skema DEV: POOLDATA.
> Angka DEV (ringkasan 3, `M_RICOMM_LIFE` 0, `RICOMM_LIFE` 2) = PEMBANDING dari WO 08-10-2026, bukan angka tetap.
> Pola sama dengan riratelife (`modul/riratelife/docs/LANGKAH-WO-RIRATELIFE-DETAIL-FLAT.md`).

Keadaan awal yang diandaikan: 924/925 sudah jalan; `M_RICOMM_LIFE_SUMMARY` = ID PK + JSONDATA (IS JSON) + view
`RICOMM_LIFE_SUMMARY`; `M_RICOMM_LIFE` = ID PK + JSONDATA (IS JSON); tabel flat `RICOMM_LIFE` (924) = rincian aplikasi.

**Berkas SQL\*Plus** di `modul/ricommlife/docs/sql/` (semuanya baca-saja). Jalankan sebagai BERKAS dari folder tempat
cadangan disimpan, mis. `@"D:\NUSARE DEV\NUSARE\NusantaraRe\modul\ricommlife\docs\sql\satu_tabel_b_cadangan.sql"` -
hanya lewat berkas `SET TERMOUT OFF` berlaku. SQL\*Plus 12.2+ atau SQLcl. ⛔ Berkas CSV dan keluaran memuat OPERATORID
(nama orang): simpan di tempat cadangan DBA, jangan disalin ke dokumen / tiket.

| Berkas | Dipakai di |
| --- | --- |
| `satu_tabel_a1_kunci.sql` / `satu_tabel_a1_kunci_dba.sql` | (a)1 kunci dan sesi |
| `satu_tabel_a2_acuan.sql` | (a)2 dan (c)0 - angka acuan, ditambahkan ke `satu_tabel_acuan.txt` |
| `satu_tabel_b_cadangan.sql` | (b) - empat CSV |
| `satu_tabel_d_verifikasi.sql` | (d) - `satu_tabel_verifikasi.txt` + dua CSV sesudah |
| `satu_tabel_keadaan.sql` | Pemulihan - penentu keadaan |

## (a) Prasyarat dan angka acuan — SEBELUM apa pun dibuang

### ⛔ (a)0 Hentikan SEMUA penulis

1. **Pega untuk R/I Comm dihentikan** (penyimpan R/I Comm, termasuk prosedur `PEGA_M_RICOMM_LIFE` /
   `PEGA_M_RICOMM_LIFE_SUMMARY` - masih VALID sampai 932/934 membuang JSONDATA).
2. **Backend :8080 dihentikan** dan tetap mati sampai (e).

Alasannya: 932/934 mengisi kolom lalu membuang JSONDATA dalam pernyataan terpisah tanpa transaksi. Baris yang disisip
DI ANTARA keduanya mendapat kolom NULL, JSON-nya ikut terbuang, dan tidak ada di cadangan CSV.

### (a)1 Tidak ada sesi yang mengunci ketiga tabel

`@satu_tabel_a1_kunci.sql` sebagai pemilik POOLDATA: kueri 1 (V$LOCKED_OBJECT ⨝ ALL_OBJECTS ⨝ V$SESSION) harus **nol
baris**; kueri 2 menampilkan sesi POOLDATA yang masih tersambung - pastikan tidak ada sesi Pega atau backend. ORA-00942
pada V$ = tidak punya hak: minta **DBA** menjalankan `satu_tabel_a1_kunci_dba.sql` (DBA_DML_LOCKS, harus nol baris) dan
memeriksa V$SESSION. Pemeriksaan ini hanya melihat transaksi yang terbuka saat itu; jaminannya tetap (a)0.

### (a)2 Angka acuan sesaat sebelum cadangan

`@satu_tabel_a2_acuan.sql` → `satu_tabel_acuan.txt`: cacah + ID tertinggi ketiga tabel, `LAST_NUMBER` kedua sequence,
sidik isi ringkasan (lewat view) dan `RICOMM_LIFE` (cacah tidak-NULL, `SUM(ORA_HASH)` per kolom / per baris), yatim,
status prosedur `PEGA_*`, nilai ringkasan yang melewati lebar 931 (harus nol baris), **NULLABLE asli JSONDATA** kedua
tabel, dan **constraint `RICOMM_LIFE` beserta SEARCH_CONDITION** (`SYS_C0015437` diharapkan `"ID" IS NOT NULL`; bila
isinya aturan LAIN: BERHENTI dan laporkan - 934 tidak menyalinnya).

⛔ **`N_JSON_RINCIAN` (`M_RICOMM_LIFE`) harus 0** (fakta DEV). Bila tidak nol: BERHENTI dan laporkan ke WO - 934 tetap
mengisi kolom baris itu dari JSON, tetapi rincian Pega di luar tabel flat belum pernah diputuskan.

## (b) Cadangan CSV

`@satu_tabel_b_cadangan.sql` → `M_RICOMM_LIFE_SUMMARY_JSONDATA.csv` (ID + JSONDATA - satu-satunya salinan `pxObjClass`
dan bentuk JSON asli), `M_RICOMM_LIFE_JSONDATA.csv` (ID + JSONDATA rincian; DEV kosong), `RICOMM_LIFE_SEBELUM.csv`
(6 kolom flat), `M_RICOMM_LIFE_SUMMARY_SEBELUM.csv` (4 kolom view). Periksa berkas terbuka dan cacah barisnya = (a)2.

## (c) Migrasi 931-934

**(c)0 TEPAT sebelum `-migrate`:** jalankan lagi `@satu_tabel_a2_acuan.sql`. Semua angka harus SAMA dengan (a)2.
**Bila ada yang berubah: ada penulis yang masih hidup - ulangi (a)0-(b).**

PowerShell, folder akar repo yang memuat kode RALAT R1:

```powershell
. .\muat-env.ps1                 # dot-source; ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate        # = target Makefile `migrate`
```

Log yang benar: `dijalankan:` 931, 932, 933, 934 berurutan. 931/933 = satu ALTER … ADD; 932 = blok isi dari JSONDATA,
blok buang JSONDATA, CREATE INDEX nama, DROP VIEW terakhir; 934 = blok isi dari JSON, blok buang JSONDATA, UPDATE dari
`RICOMM_LIFE`, INSERT `NOT EXISTS`, CREATE INDEX, DROP TABLE `RICOMM_LIFE` terakhir.

## (d) Verifikasi — baca-saja

`@satu_tabel_d_verifikasi.sql` → `satu_tabel_verifikasi.txt`:

1. `M_RICOMM_LIFE_SUMMARY` tepat 4 kolom (ID, USEDBY VARCHAR2(200), MODIFIEDDATE VARCHAR2(50), OPERATORID
   VARCHAR2(200)); `M_RICOMM_LIFE` tepat 6 kolom (ID, IDUSEDBY VARCHAR2(10), USEDBY VARCHAR2(200), CONTRACT NUMBER(5),
   YEAR NUMBER(5), COMM NUMBER(38,8)); nol JSONDATA; `ALL_PARTIAL_DROP_TABS` nol baris;
2. `N_RINGKASAN` = acuan (pembanding 3); `N_RINCIAN` = `N_FLAT` acuan (pembanding 2) + `N_JSON_RINCIAN` (0); yatim = acuan;
3. sidik isi = acuan (a)2 kueri 2 (ringkasan) dan 3 (rincian);
4. `RICOMM_LIFE` dan `RICOMM_LIFE_SUMMARY` tidak ada; constraint IS JSON hilang; `IX_M_RICOMM_LIFE_SUMMARY_NAMA` dan
   `IX_M_RICOMM_LIFE_IDUSEDBY` ada;
5. T_MIGRASI memuat 931, 932, 933, 934;
6. **catat** status `PEGA_M_RICOMM_LIFE` / `PEGA_M_RICOMM_LIFE_SUMMARY` (INVALID diharapkan, diterima WO).

**MINUS dua arah terhadap cadangan, TANPA tabel baru** - CSV sesudah (ditulis berkas verifikasi) dibandingkan dengan CSV
sebelum yang dibentuk sama:

```powershell
foreach ($p in @(@('RICOMM_LIFE_SEBELUM.csv', 'M_RICOMM_LIFE_SESUDAH.csv'),
                 @('M_RICOMM_LIFE_SUMMARY_SEBELUM.csv', 'M_RICOMM_LIFE_SUMMARY_SESUDAH.csv'))) {
  $a = Get-Content $p[0] -Encoding UTF8; $b = Get-Content $p[1] -Encoding UTF8
  "$($p[1]): $($a.Count) baris sebelum, $($b.Count) sesudah"
  Compare-Object $a $b -CaseSensitive | Group-Object SideIndicator | Select-Object Name, Count   # kosong = sama
}
```

`<=` = sebelum MINUS sesudah, `=>` = sesudah MINUS sebelum; keluaran kosong = MINUS dua arah nol (`-CaseSensitive`
wajib). Sidik (d)3 = pemeriksaan kedua yang tidak bergantung berkas.

## (e) Restart backend dan periksa layar

```powershell
go run ./cmd/api                 # = target Makefile `run-api`, jendela yang sudah memuat env
```

**R/I Comm Life**: grid ringkasan (jumlah = (d)2, filter, urut, tanggal); Detail sebuah ringkasan (rincian sama, angka
koma desimal); Save ringkasan baru uji dan Edit nama (rincian ikut berganti nama); tambah / ubah satu rincian; Upload
(pratinjau + simpan) satu berkas uji; Delete ringkasan uji beserta rinciannya. Pega R/I Comm TIDAK dinyalakan lagi untuk
menyimpan (prosedur INVALID - diterima WO).

## Pemulihan — 931-934 gagal di tengah

Pelari menjalankan pernyataan TANPA transaksi. Backend dan Pega tetap mati sampai pulih. `@satu_tabel_keadaan.sql`
memberi dua baris keadaan - ringkasan (931/932: `T931`, `T932`, `ADA_VIEW`) dan rincian (933/934: `T933`, `T934`,
`ADA_FLAT`) - dengan `KOLOM_BARU`, `ADA_JSONDATA`, `SETENGAH_TERBUANG`, `ADA_INDEKS`. Ikuti SATU baris per pasangan
(X = 931 / 933, Y = 932 / 934, OBJEK = view `RICOMM_LIFE_SUMMARY` / tabel `RICOMM_LIFE`):

| # | Keadaan | Maju (tindakan) | Mundur boleh? |
| --- | --- | --- | --- |
| K0 | `KOLOM_BARU` 0, `TX` 0 | X tidak mengubah apa pun (satu pernyataan). Perbaiki sebabnya (ORA), ulangi (c) | tidak perlu |
| K1 | `KOLOM_BARU` penuh (3 / 5), `TX` 0 | jangan ulangi X (ORA-01430): `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('<nama X>', SYSDATE); COMMIT;` lalu baca lagi keadaan | lihat K2 |
| K2 | `KOLOM_BARU` penuh, `ADA_JSONDATA` 1, `SETENGAH_TERBUANG` 0, `TY` 0 | Y belum mulai, atau blok isi dibatalkan utuh (ORA-12899 / ORA-01722: BERHENTI, laporkan ID dari (a)2, jangan memperlebar tanpa WO), atau isi selesai tetapi JSONDATA belum dibuang. Ulangi (c) | YA (JSONDATA utuh): DBA menjalankan `X_down` (pelindungnya lolos), lalu `DELETE FROM POOLDATA.T_MIGRASI WHERE NAMA = '<nama X>'; COMMIT;`. `Y_down` TIDAK (Y tidak tercatat) |
| K3 | `SETENGAH_TERBUANG` 1 (DROP COLUMN terputus) | `ALTER TABLE POOLDATA.<tabel> DROP COLUMNS CONTINUE;` lalu baca lagi → K4 | **TIDAK** |
| K4 | `ADA_JSONDATA` 0, `SETENGAH_TERBUANG` 0, `TY` 0 | **langkah MAJU WAJIB diselesaikan.** OBJEK masih ada: ulangi (c) (blok melewati diri; UPDATE dari flat dan INSERT `NOT EXISTS` aman diulang; CREATE INDEX ada = ORA-00955 ditoleransi; DROP terakhir). OBJEK sudah tidak ada: jangan ulangi (c); pastikan (d)1-4 benar (934: `N_FLAT_BELUM_PINDAH` sebelum DROP = 0 terbukti dari (d)2/(d)3), lalu `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('<nama Y>', SYSDATE); COMMIT;` | **TIDAK.** `X_down` DILARANG (isi hanya tinggal di kolom); pelindungnya gagal keras ORA-00904 bila dicoba |
| K5 | `TY` 1 | selesai - lanjut (d) | hanya lewat Jalur mundur |

Nama langkah: `931_m_ricomm_life_summary_kolom`, `932_m_ricomm_life_summary_satu_tabel`, `933_m_ricomm_life_kolom`,
`934_m_ricomm_life_satu_tabel`. 934 pada K4 dengan `ADA_FLAT` 1: rincian yang belum pindah masih di `RICOMM_LIFE`
(`N_FLAT_BELUM_PINDAH`), ulangi (c) memindahkannya.

## Jalur mundur

Hanya dari K5 kedua pasangan. Hentikan backend. DBA menjalankan pernyataan `inti/backend/migrations/` berurutan:
`934_m_ricomm_life_satu_tabel_down.sql`, `933_m_ricomm_life_kolom_down.sql`, `932_m_ricomm_life_summary_satu_tabel_down.sql`,
`931_m_ricomm_life_summary_kolom_down.sql` (ganti `{skema}` dengan POOLDATA), lalu `DELETE FROM POOLDATA.T_MIGRASI WHERE
NAMA IN ('934_m_ricomm_life_satu_tabel', '933_m_ricomm_life_kolom', '932_m_ricomm_life_summary_satu_tabel',
'931_m_ricomm_life_summary_kolom'); COMMIT;` dan backend dari kode sebelum RALAT R1.

- 934_down membangun ulang `RICOMM_LIFE` PERSIS bentuk 924 (kolom, `PK_RICOMM_LIFE`, `IX_RICOMM_LIFE_IDUSEDBY`) berisi
  data, dan JSONDATA rincian + `ENSURE_M_RICOMM_LIFE_JSON`; baris tetap juga ada di `M_RICOMM_LIFE` (sebelum 934: 0 baris
  JSON) - backend lama membaca `RICOMM_LIFE`, ID terpakai tetap aman.
- 932_down membangun ulang JSONDATA ringkasan (kunci USEDBY, MODIFIEDDATE, OPERATORID), `ENSURE_M_RICOMM_LIFE_SUMMARY_JSON`
  (nama asli 33 byte), dan view `RICOMM_LIFE_SUMMARY` persis.
- Kedua `_down` aman diulang (blok berpelindung katalog); CREATE VIEW 932_down terakhir - bila hanya view yang sudah ada,
  lewati pernyataan itu.
- 933_down / 931_down mulai dengan pelindung gagal-keras (`UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0`): tanpa
  JSONDATA ia mati ORA-00904 sebelum membuang kolom.
- JSONDATA dikembalikan **NULLABLE** (repo tidak memuat bukti NOT NULL). Bila (a)2 kueri 6 menunjukkan `NULLABLE = N`,
  DBA menambah `ALTER TABLE POOLDATA.<tabel> MODIFY (JSONDATA NOT NULL);` sesudah mundur.
- `pxObjClass` dan bentuk JSON asli hanya dari cadangan `M_RICOMM_LIFE_SUMMARY_JSONDATA.csv` (per ID).
- Prosedur `PEGA_*` perlu dikompilasi ulang DBA.
- ⛔ Jangan memakai `-migrate-down` di DEV (dipagari skema uji, membongkar SELURUH migrasi).
