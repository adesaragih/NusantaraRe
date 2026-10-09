# Langkah work owner — rincian R/I Rate Life menjadi tabel flat `M_RATE_LIFE` (RALAT R7)

> Untuk work owner / DBA. **Dijalankan WO, bukan executor** — tidak satu pun perintah di bawah pernah dijalankan saat
> menulis berkas ini. Keputusan work owner 07-10-2026 (`modul/riratelife/MODUL.md`, RALAT R7). Skema DEV: POOLDATA.
> Angka DEV (98.306 baris, 539 yatim) = PEMBANDING dari WO 07-10-2026, bukan angka tetap.

Keadaan awal yang diandaikan: 927/928 sudah jalan (ringkasan satu tabel); `POOLDATA.M_RATE_LIFE` = `ID` PK +
`JSONDATA` CLOB (`ENSURE_M_RATE_LIFE_JSON`); view `POOLDATA.RATE_LIFE` 8 kolom.

**Berkas SQL\*Plus** di `modul/riratelife/docs/sql/` (semuanya baca-saja). Jalankan sebagai BERKAS dari folder tempat
cadangan disimpan, mis. `@"D:\NUSARE DEV\NUSARE\NusantaraRe\modul\riratelife\docs\sql\detail_flat_a3_cadangan.sql"` -
hanya lewat berkas `SET TERMOUT OFF` berlaku (98 ribu baris tidak dicetak ke layar). SQL\*Plus 12.2+ atau SQLcl.

| Berkas | Dipakai di |
| --- | --- |
| `detail_flat_a1_kunci.sql` / `detail_flat_a1_kunci_dba.sql` | (a)1 kunci dan sesi |
| `detail_flat_a2_acuan.sql` | (a)2 dan (b)0 - angka acuan, ditambahkan ke `detail_flat_acuan.txt` |
| `detail_flat_a3_cadangan.sql` | (a)3 - `M_RATE_LIFE_JSONDATA.csv`, `M_RATE_LIFE_SEBELUM.csv` |
| `detail_flat_c_verifikasi.sql` | (c) - `detail_flat_verifikasi.txt`, `M_RATE_LIFE_SESUDAH.csv` |
| `detail_flat_keadaan.sql` | Pemulihan - penentu keadaan |

## (a) Prasyarat, cadangan, angka acuan — SEBELUM apa pun dibuang

### ⛔ (a)0 Prasyarat - hentikan SEMUA penulis `M_RATE_LIFE`

1. **Pega untuk R/I Rate dihentikan** (aplikasi / agen yang menyimpan R/I Rate, termasuk prosedur `PEGA_M_RATE_LIFE` -
   prosedur itu masih VALID sampai 930 membuang JSONDATA).
2. **Backend :8080 dihentikan** dan tetap mati sampai (d).

Alasannya: 930 mengisi kolom (UPDATE) lalu membuang JSONDATA dalam dua pernyataan terpisah tanpa transaksi. Baris yang
disisip Pega/backend DI ANTARA keduanya mendapat kolom NULL, JSON-nya ikut terbuang, dan tidak ada di cadangan CSV -
hilang permanen.

### (a)1 Tidak ada sesi yang mengunci `M_RATE_LIFE`

`@detail_flat_a1_kunci.sql` sebagai pemilik POOLDATA. Kueri 1 (V$LOCKED_OBJECT ⨝ ALL_OBJECTS ⨝ V$SESSION) harus
**nol baris**; kueri 2 menampilkan sesi POOLDATA yang masih tersambung (nama mesin / program / modul) - pastikan tidak
ada sesi Pega atau backend. ORA-00942 pada V$ = pemilik POOLDATA tidak punya hak: minta **DBA** menjalankan
`detail_flat_a1_kunci_dba.sql` (DBA_DML_LOCKS, harus nol baris) dan memeriksa V$SESSION. Pemeriksaan ini hanya melihat
transaksi yang TERBUKA saat itu; jaminan sesungguhnya tetap (a)0.

### (a)2 Angka acuan sesaat sebelum cadangan

`@detail_flat_a2_acuan.sql` - ditambahkan ke `detail_flat_acuan.txt`: waktu baca, `N` + `ID_TERTINGGI` (`MAX(LPAD(ID,
10, '0'))`) + `LAST_NUMBER` SEQ_M_RATE_LIFE, sidik isi (cacah tidak-NULL dan `SUM(ORA_HASH)` per kolom, `H_BARIS`
per baris), `YATIM`, status `PEGA_M_RATE_LIFE`, nilai yang melewati lebar 929 (harus nol baris; bila ada: BERHENTI -
930 akan mati di ORA-12899), dan **NULLABLE asli JSONDATA** (dipakai jalur mundur). Kueri sidik = satu pindai penuh
(± 98 ribu baris JSON).

### (a)3 Cadangan CSV

`@detail_flat_a3_cadangan.sql` → `M_RATE_LIFE_JSONDATA.csv` (ID + JSONDATA - satu-satunya salinan `FLAG`, RALAT R5, dan
bentuk JSON asli) dan `M_RATE_LIFE_SEBELUM.csv` (8 kolom view, acuan (c)3b). SQL Developer sebagai ganti: hasil kueri
yang sama → klik kanan → Export → csv. Periksa kedua berkas terbuka dan cacah barisnya = `N` (a)2.

## (b) Migrasi 929 + 930

**(b)0 TEPAT sebelum `-migrate`:** jalankan lagi `@detail_flat_a2_acuan.sql`. `N`, `ID_TERTINGGI`, `LAST_NUMBER`, dan
seluruh sidik harus SAMA dengan (a)2. **Bila ada yang berubah: ada penulis yang masih hidup - ulangi (a)0-(a)3.**

PowerShell, folder akar repo yang memuat kode RALAT R7:

```powershell
. .\muat-env.ps1                 # dot-source; ORACLE_SCHEMA=POOLDATA, IS_PEGA_PROD=false (PANDUAN-MENJALANKAN.txt bab 3)
go run ./cmd/api -migrate        # = target Makefile `migrate`
```

Versi cmd.exe: `call muat-env.cmd` lalu `go run ./cmd/api -migrate`. Log yang benar: `dijalankan: 929_m_rate_life_kolom`
lalu `dijalankan: 930_m_rate_life_satu_tabel`. 929 = satu ALTER … ADD; 930 = blok berpelindung katalog UPDATE ketujuh
kolom dari JSONDATA (satu pernyataan, detik hingga beberapa menit, undo puluhan MB), blok berpelindung katalog
`DROP COLUMN JSONDATA CASCADE CONSTRAINTS`, CREATE INDEX, DROP VIEW `RATE_LIFE` terakhir.

## (c) Verifikasi — baca-saja

`@detail_flat_c_verifikasi.sql` → `detail_flat_verifikasi.txt`:

1. tepat 8 kolom `ID, IDUSEDBY, USEDBY, TYPE, GENDER, CONTRACT, AGE, RATE` (VARCHAR2, nol JSONDATA) dan
   `SETENGAH_TERBUANG` = 0;
2. `N`, `ID_TERTINGGI`, `YATIM` = (a)2/(b)0;
3. **3a** sidik isi (cacah tidak-NULL, `SUM(ORA_HASH)` per kolom dan `H_BARIS`) = (a)2/(b)0;
4. view `RATE_LIFE` dan constraint `ENSURE_M_RATE_LIFE_JSON` hilang (nol baris), indeks `IX_M_RATE_LIFE_IDUSEDBY`
   (IDUSEDBY) ada;
5. T_MIGRASI memuat 929 dan 930;
6. **catat** status `PEGA_M_RATE_LIFE` (INVALID diharapkan, diterima WO; jangan dikompilasi ulang).

**3b. MINUS dua arah atas 8 kolom terhadap cadangan, TANPA tabel baru.** View sudah tidak ada dan cadangan hanya ada
sebagai berkas, jadi MINUS dilakukan di luar basis data atas dua CSV yang dibentuk sama (kolom, `ORDER BY ID`, `MARKUP
CSV ON QUOTE ON`). Berkas verifikasi di atas sudah menulis `M_RATE_LIFE_SESUDAH.csv`; lalu PowerShell:

```powershell
$a = Get-Content .\M_RATE_LIFE_SEBELUM.csv -Encoding UTF8
$b = Get-Content .\M_RATE_LIFE_SESUDAH.csv -Encoding UTF8
"$($a.Count) baris sebelum, $($b.Count) baris sesudah"
Compare-Object $a $b -CaseSensitive | Group-Object SideIndicator | Select-Object Name, Count   # kosong = sama
```

`<=` = sebelum MINUS sesudah, `=>` = sesudah MINUS sebelum; keluaran kosong = MINUS dua arah nol. `-CaseSensitive`
wajib (tanpa itu `u` dan `U` dianggap sama). Sidik 3a = pemeriksaan kedua yang tidak bergantung berkas: cacah tidak-NULL
menangkap nilai hilang/muncul, `SUM(ORA_HASH)` menangkap nilai yang berubah; bila ada satu angka beda, CSV menunjukkan
barisnya.

## (d) Restart backend dan periksa layar

```powershell
go run ./cmd/api                 # = target Makefile `run-api`, jendela yang sudah memuat env
```

- **R/I Rate Life:** Rate Detail sebuah ringkasan menampilkan rincian yang sama (RATE apa adanya, berkoma/bertitik);
  tambah dan ubah satu baris uji lalu kembalikan; Upload; Edit nama ringkasan mengganti `USEDBY` rinciannya.
- **Claim Life:** spreading retro sebuah klaim berproduk retro menghitung rate seperti sebelumnya (tidak ada galat
  `M_RATE_LIFE` / `RATE_LIFE` di log).
- **PremiumList Life:** rate produk dan Hitung QR sebuah polis menampilkan rate seperti sebelumnya.
- **Product Name Life:** dialog `View Rate` menampilkan rincian sebuah R/I Rate.
- (Contract Retro Life `Rate List` membaca objek yang sama - periksa bila sempat.)

Pega R/I Rate TIDAK dinyalakan lagi untuk menyimpan rate (`PEGA_M_RATE_LIFE` INVALID - diterima WO).

## Pemulihan — 929 / 930 gagal di tengah

Pelari menjalankan pernyataan TANPA transaksi; DDL langsung tetap. Backend dan Pega tetap mati sampai keadaan pulih.
Tentukan keadaan dengan `@detail_flat_keadaan.sql` (`KOLOM_BARU`, `ADA_JSONDATA`, `SETENGAH_TERBUANG`, `ADA_VIEW`,
`ADA_INDEKS`, `T929`, `T930`, lalu `N` / `N_IDUSEDBY` / `N_RATE`), lalu ikuti SATU baris:

| # | Keadaan | Maju (tindakan) | Mundur boleh? |
| --- | --- | --- | --- |
| K0 | `KOLOM_BARU` 0, `T929` 0 | 929 tidak mengubah apa pun (satu pernyataan). Perbaiki sebabnya (pesan ORA), ulangi (b) | tidak perlu |
| K1 | `KOLOM_BARU` 7, `T929` 0 (929 jalan, tidak tercatat) | jangan ulangi 929 (ORA-01430): `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('929_m_rate_life_kolom', SYSDATE); COMMIT;` lalu baca lagi keadaan | lihat K2 |
| K2 | `KOLOM_BARU` 7, `ADA_JSONDATA` 1, `SETENGAH_TERBUANG` 0, `T930` 0 | 930 belum mulai, atau UPDATE-nya dibatalkan utuh (mis. ORA-12899 → BERHENTI, laporkan ID dari (a)2 kueri 4, jangan memperlebar tanpa WO), atau UPDATE selesai tetapi JSONDATA belum dibuang. Ulangi (b) | YA, karena JSONDATA utuh: DBA menjalankan `929_m_rate_life_kolom_down.sql` (pelindungnya lolos), lalu `DELETE FROM POOLDATA.T_MIGRASI WHERE NAMA = '929_m_rate_life_kolom'; COMMIT;`. `930_down` TIDAK (930 tidak tercatat) |
| K3 | `SETENGAH_TERBUANG` 1 (DROP COLUMN terputus - kolom ditandai UNUSED, tidak terlihat di ALL_TAB_COLUMNS) | `ALTER TABLE POOLDATA.M_RATE_LIFE DROP COLUMNS CONTINUE;` lalu baca lagi keadaan → K4 | **TIDAK** - JSONDATA tidak bisa kembali |
| K4 | `KOLOM_BARU` 7, `ADA_JSONDATA` 0, `SETENGAH_TERBUANG` 0, `T930` 0 | **langkah MAJU WAJIB diselesaikan.** `ADA_VIEW` 1: ulangi (b) (kedua blok melewati diri, CREATE INDEX dibuat atau ORA-00955 ditoleransi, DROP VIEW). `ADA_VIEW` 0: jangan ulangi (b) (DROP VIEW mati di ORA-00942); pastikan (c)1-4 benar, lalu `INSERT INTO POOLDATA.T_MIGRASI (NAMA, DIJALANKAN_PADA) VALUES ('930_m_rate_life_satu_tabel', SYSDATE); COMMIT;` | **TIDAK.** `929_down` DILARANG (isi rincian hanya tinggal di kolom); pelindungnya gagal keras ORA-00904 bila dicoba. Selesaikan maju dulu, baru jalur mundur penuh bila perlu |
| K5 | `T930` 1 | selesai - lanjut (c) | hanya lewat Jalur mundur di bawah |

Isi kolom pada K4 sudah benar karena blok pengisi berjalan SEBELUM blok pembuang; buktikan dengan (c)3a/3b.

## Jalur mundur

Hanya dari K5 (930 TERCATAT di T_MIGRASI). Hentikan backend. DBA menjalankan pernyataan
`inti/backend/migrations/930_m_rate_life_satu_tabel_down.sql` lalu `929_m_rate_life_kolom_down.sql` berurutan (ganti
`{skema}` dengan POOLDATA), lalu `DELETE FROM POOLDATA.T_MIGRASI WHERE NAMA IN ('930_m_rate_life_satu_tabel',
'929_m_rate_life_kolom'); COMMIT;` dan backend dari kode sebelum RALAT R7.

- 930_down aman diulang (indeks, JSONDATA, constraint lewat blok berpelindung katalog; CREATE VIEW terakhir - bila
  hanya view yang sudah ada, lewati pernyataan itu).
- 929_down mulai dengan pelindung gagal-keras (`UPDATE … SET JSONDATA = JSONDATA WHERE 1 = 0`): tanpa JSONDATA ia mati
  ORA-00904 sebelum membuang kolom.
- JSONDATA dibangun ulang sebagai teks ketujuh kunci view, **NULLABLE**: repo tidak memuat bukti NOT NULL (fakta WO
  07-10-2026 hanya "ID PK + JSONDATA CLOB IS JSON"). Bila (a)2 kueri 5 menunjukkan `NULLABLE = N`, DBA menambah
  `ALTER TABLE POOLDATA.M_RATE_LIFE MODIFY (JSONDATA NOT NULL);` sesudah mundur.
- **`FLAG` dan bentuk JSON asli hanya dari cadangan `M_RATE_LIFE_JSONDATA.csv`** (per ID).
- `PEGA_M_RATE_LIFE` perlu dikompilasi ulang DBA.
- ⛔ Jangan memakai `-migrate-down` di DEV (dipagari skema uji, membongkar SELURUH migrasi).
