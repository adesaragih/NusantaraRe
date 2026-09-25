# FINDING-005 — Tanggal kerugian dibandingkan terhadap tanggal akhir treaty dalam format berbeda

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

**Jenis**: laporan kondisi sistem lama
**Status**: TERBUKTI DARI XML — tidak menunggu data Oracle
**Ruang lingkup bukti**: `D:\XML_NURE\Claim Non Prop` saja
**Tanggal**: 18 September 2026

---

## 1. Dua format, satu perbandingan

### 1.1 `.DateOfLoss` berformat `YYYYMMDD`

Dibuktikan dua kali, dari dua arah berbeda:

| Bukti | Lokasi |
|---|---|
| Diurai per posisi: `@substring(.DateOfLoss,0,4)` = tahun, `(4,6)` = bulan, `(6,8)` = tanggal | beberapa Section/Activity |
| Disambung menjadi DateTime Pega: `@FormatDateTime((pyWorkPage.ClaimData.DateOfLoss + "T050000.000 GMT"), "dd/MM/yyyy", ...)` | `Activity\GenerateCACNP_Act.xml` baris 1416, `Activity\GenerateCFS_act.xml` baris 1146 |

Format `YYYYMMDDTHHMMSS.mmm GMT` adalah bentuk DateTime internal Pega. Penyambungan itu hanya sah bila bagian kirinya `YYYYMMDD`.

Dikuatkan dari sisi basis data: kolom `DATEOFLOSS` pada tabel work Pega bertipe `VARCHAR2(8 BYTE)` — delapan karakter, persis panjang `YYYYMMDD`. *(Sumber: disampaikan pengguna; DDL tabel work belum masuk ke berkas `.xls` — lihat catatan di `BLUEPRINT.md` bagian 13.)*

### 1.2 `.EndDateTreaty` berformat `dd/MM/yyyy`

| Bukti | Lokasi |
|---|---|
| `@FormatDateTime(...,"dd/MM/yyyy","Asia/Jakarta","in_ID") == pyWorkPage.ClaimData.EndDateTreaty` | `Activity\CheckPeriodPolicy_Act.xml` baris 401 |
| `@CompareDates(@FormatDateTime(...,"dd/MM/yyyy",...), pyWorkPage.ClaimData.EndDateTreaty)` | `Activity\CheckPeriodPolicy_Act.xml` baris 314 |
| `@FormatDateTime(...,"dd/MM/yyyy",...) > pyWorkPage.ClaimData.EndDateTreaty` | `Activity\SetEndDate_Act.xml` baris 1643 |
| Dideklarasikan `TYPE="STRING"` | definisi parameter |

Tiga tempat mengubah nilai lain ke `dd/MM/yyyy` **lebih dulu** sebelum membandingkannya dengan `.EndDateTreaty`. Itu memperlihatkan format yang diharapkan.

### 1.3 Perbandingan yang tidak menyamakan format

```
pyWorkPage.ClaimData.DateOfLoss > pyWorkPage.ClaimData.EndDateTreaty
```
`Activity\CheckDateDOL_Act.xml` baris **4011** dan **4611**.

Kedua sisi adalah teks. Yang dibandingkan: `"20260315"` terhadap `"31/12/2026"` — **perbandingan leksikografis atas dua format yang berbeda**, bukan perbandingan tanggal.

## 2. Akibatnya

Hasil perbandingan ditentukan oleh **digit pertama** masing-masing teks, bukan oleh tanggalnya:

| `.DateOfLoss` | `.EndDateTreaty` | Perbandingan teks | Yang benar secara tanggal |
|---|---|---|---|
| `20260315` (15 Mar 2026) | `31/12/2026` | `"2" < "3"` → **false** | 15 Mar 2026 < 31 Des 2026 → false |
| `20260315` (15 Mar 2026) | `01/01/2026` | `"2" > "0"` → **true** | 15 Mar 2026 > 1 Jan 2026 → true |
| `20251115` (15 Nov 2025) | `01/01/2026` | `"2" > "0"` → **true** | 15 Nov 2025 < 1 Jan 2026 → **false** |

Baris ketiga adalah kesalahannya: hasilnya benar hanya **kebetulan**, ketika digit pertama tanggal treaty kebetulan selaras. Secara umum, **hasil perbandingan bergantung pada digit pertama hari dalam bulan pada tanggal akhir treaty** — `0`, `1`, `2`, atau `3` — dan sama sekali tidak bergantung pada tahunnya.

Rule yang memuatnya bernama `CheckDateDOL_Act` — pemeriksa tanggal kerugian. Yaitu tepat penjaga yang seharusnya menolak klaim dengan tanggal kejadian di luar masa berlaku treaty.

## 3. Pola yang sama di rule tetangga

`Activity\CheckPeriodPolicy_Act.xml` memuat keduanya sekaligus:

| Baris | Ekspresi | Menyamakan format? |
|---|---|---|
| 314, 401 | `@FormatDateTime(...,"dd/MM/yyyy",...)` dibandingkan ke `.EndDateTreaty` | **ya** |
| 868 | `pyWorkPage.ClaimData.PolicyData.EndDateTime > pyWorkPage.ClaimData.EndDateTreaty` | **tidak** — DateTime Pega dibandingkan langsung ke teks `dd/MM/yyyy` |

Jadi di dalam satu rule, sebagian perbandingan menyamakan format dan sebagian tidak. Ini menunjukkan penyebabnya bukan ketidaktahuan, melainkan ketidakseragaman.

## 4. Batas klaim

1. **Belum diuji apakah kesalahannya pernah berakibat di produksi.** Yang terbukti adalah perbandingannya salah bentuk. Berapa klaim yang lolos atau tertolak karenanya adalah pertanyaan data → **REQ-020**.
2. **Nilai `.EndDateTreaty` yang sebenarnya tersimpan belum diperiksa.** Format `dd/MM/yyyy` disimpulkan dari cara nilai lain diubah sebelum dibandingkan dengannya, bukan dari contoh isinya.
3. **Penulis `.DateOfLoss` tidak ditemukan.** Sejalan dengan temuan sebelumnya bahwa `.DateOfLoss` tidak pernah ditulis oleh activity mana pun di folder ini (ADR-0001).

## 5. Hubungan dengan sistem baru

Tanggal disimpan sebagai tipe tanggal, bukan teks, dan seluruh perbandingan dilakukan atas tipe tanggal. Itu sudah menjadi konsekuensi wajar dari ADR-0001 dan tidak memerlukan keputusan baru.

Yang memerlukan keputusan: **apa yang dilakukan terhadap klaim lama yang lolos penjaga ini karena perbandingannya salah.** Itu pertanyaan migrasi data, bukan pertanyaan rancangan.

---

## 6. Tambahan — asal-usul format `dd/MM/yyyy` ditemukan

Versi pertama menyimpulkan format `.EndDateTreaty` dari cara nilai lain dinormalkan sebelum dibandingkan dengannya. Pembacaan `pengetahuan/ddl/VIEW_V_POLIS.sql` menunjukkan sumbernya langsung.

`V_POLIS` — 26 kolom, seluruhnya proyeksi atas **satu** tabel `json_polis`, dengan 32 pemanggilan `JSON_VALUE`. Kolom tanggalnya dibentuk begini:

```sql
SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 7, 2)
|| '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 5, 2)
|| '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 1, 4)
|| ' ' || SUBSTR(..., 10, 2) || ':' || SUBSTR(..., 12, 2) || ':' || SUBSTR(..., 14, 2)
```

Jadi rantainya lengkap sekarang:

| Tahap | Bentuk |
|---|---|
| Tersimpan di JSON | `YYYYMMDDHHMISS` (teks) |
| Diubah oleh `V_POLIS` dengan `SUBSTR` + `\|\|` | `dd/MM/yyyy HH:MI:SS` (teks) |
| Dibaca Pega sebagai `.EndDateTreaty` | teks `dd/MM/yyyy` |
| `.DateOfLoss` | tetap `YYYYMMDD` |

**Penyebabnya bukan kesalahan di satu rule.** Basis data mengubah format tanggal menjadi teks bergaya Indonesia di lapisan view, sementara `.DateOfLoss` tidak melewati lapisan itu dan tetap dalam bentuk aslinya. Perbandingan di `CheckDateDOL_Act` mempertemukan keduanya.

Tidak ada satu pun kolom tanggal di `V_POLIS` yang bertipe `DATE`. Semuanya teks hasil sambungan.

### 6.1 Dua bentuk JSON di satu tabel, dan satu di antaranya terpotong

`V_POLIS` bercabang pada `a.data_json.QuotationData.BusinessFac`:

| Nilai | Jalur JSON yang dibaca | Hasilnya |
|---|---|---|
| `'F'` (Facultative) | `$.PolicyData.StartDateTime`, `$.PolicyData.EndDateTime` | `dd/MM/yyyy HH:MI:SS` — lengkap |
| `'T'` (Treaty) | `$.StartDate`, `$.EndDate` | `dd/MM/yyyy HH` — **berhenti di jam, tanpa menit dan detik** |

Dua hal terbaca dari sini:

1. **Satu tabel `json_polis` menampung dua bentuk dokumen yang berbeda**, dibedakan `BusinessFac`. Jalur JSON-nya tidak sama.
2. **Cabang Treaty menghasilkan teks tanggal yang lebih pendek.** Karena perbandingan atasnya adalah perbandingan teks, dua nilai dengan panjang berbeda tidak dapat dibandingkan secara andal. Klaim non-proporsional berjalan di jalur Treaty — jalur yang terpotong itu.

**Batas klaim**: apakah pemotongan itu berakibat bergantung pada perbandingan mana yang benar-benar dijalankan atas kolom tersebut. Belum ditelusuri seluruhnya.
