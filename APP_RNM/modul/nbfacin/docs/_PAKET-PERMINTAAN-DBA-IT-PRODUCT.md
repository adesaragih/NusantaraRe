# Paket Permintaan — DBA · IT · Product

> Dokumen ini dipecah per pihak. Kirim bagian yang relevan ke pemiliknya agar semua berjalan paralel.
> Semua permintaan sudah terverifikasi dari korpus + DDL yang sudah masuk. Status per butir mengacu
> ke `_CHECKLIST-KESIAPAN.md`.
>
> **Aturan data:** minta **definisi/struktur dan nilai referensi saja**. Jangan sertakan data
> pelanggan. Kolom yang memuat nama orang (`NAMA`, `LOGIN`) tidak perlu dikirim.
>
> Legenda: 🟢 query siap-jalan · 🟡 sudah dipetakan, tinggal dieksekusi · ⬜ perlu disiapkan pihaknya
>
> ⛔ **Jangan pernah meminta atau mengirim kolom `USERNAME`, `PASSWORD`, `NAMA`, `LOGIN`,
> `OPERATORID`.** Kredensial dan identitas orang tidak masuk artefak mana pun (`CLAUDE.md` §3 butir 5,
> §4.4). **Nama kolomnya boleh disebut; isinya tidak.**

---

# 🚀 PAKET SIAP KIRIM — status **1 Oktober 2026** (pengerjaan tiket NB-01 … NB-17)

> Bagian ini **menggantikan** paket 21 September sebagai daftar kerja untuk modul `nbfacin`; paket lama di bawahnya
> **tidak dihapus**. Setiap butir menyebut berkas pertanyaan dan butir register (`KEPUTUSAN-30-09-2026.md`) asalnya.
> ⛔ Aturan data tetap: tanpa `NAMA`, `LOGIN`, nomor kasus, atau data pelanggan — cukup hitungan, struktur, dan angka.

### Untuk DBA

| # | Permintaan | Bentuk jawaban | Memblokir |
| :-: | --- | --- | --- |
| **D-1** | **Satu query hitungan** produksi NB Fac In: kasus PA ber-`CalculateMethod_FacIn = 2`; kasus PA ber-`DiscountType = Percent`; coverage MBU dengan `Loading` terisi; coverage MBU dengan `ProRatePercent` ≠ 100 | empat angka hitungan. Hitungan nol = jalur tidak terpakai (cukup ditandai). Tidak nol → satu contoh per jalur, **daftar-izin angka saja** (TSI, Rate, ProRatePercent, PctShortPeriod, Discount/DiscountPercentage, Loading, Premium) | rekonsiliasi nyata NB-04 (metode 2, diskon persen) dan bentuk bersarang NB-05 (`PERTANYAAN-LANJUTAN.md` butir 2–3; butir 35) |
| ~~D-2~~ | ~~CSV tabel limit `M_LIMIT_*`~~ | ✅ **diterima 01-10-2026** (butir 42) | — |

### Untuk pemilik aplikasi Pega

| # | Pertanyaan | Bentuk jawaban | Memblokir |
| :-: | --- | --- | --- |
| **G-1** | Kasus MBU NB #5 (fixture `mbu_mata_uang.json`): `CurrencyList.Premium` tersimpan 309.703.644,8, padahal langkah 2.6 `FillPremiMBU_FacIn` atas data coverage tersimpan memberi 201.301.628. **Jalur mana yang menulis nilai itu?** (unggah? activity lain? nilai sisa sebelum coverage diubah?) | nama activity/flow, atau "data sisa" | rekonsiliasi agregat NB-07 — satu-satunya selisih di kerangka tiket 16 (butir 54) |

### Untuk Underwriting

| # | Pertanyaan | Bentuk jawaban | Memblokir |
| :-: | --- | --- | --- |
| **U-1** | Keputusan **Revise** mematikan `ConfirmBinding` dan `ReceivedRiSlip` sama seperti **Reject** (`NB FacIn\DataTransform\SetReviseProposal.xml` L179/L209). Disengaja? | ya / tidak | A10 — memperluas K-014; port sekarang mengikuti sistem lama |

### Untuk work owner (konfirmasi keputusan agent)

| # | Pilihan agent | Butir |
| :-: | --- | --- |
| **W-1** | A16–A19 — bentuk fixture de-identifikasi (pisah identitas vs teks bebas; `ObjectName`/`Others` tidak dikosongkan; bagian alamat dikosongkan; salinan angka `Remarks` hilang) | 38 |
| **W-2** | A31 — agregat dibandingkan sebagai nilai desimal eksak, bukan teks (dasarnya dikuatkan butir 52) | NB-07 |
| **W-3** | Tiket baru rumus premi FIRE, ANEKA, BONDING, GOLF, MARINE CARGO — dibuat manusia lewat `/mattpocock-skills:to-tickets` (agent tidak dapat memanggilnya) | prasyarat tahap 2 tiket 16 |

### Untuk pemilik folder `dastin/`

| # | Temuan | Tindakan |
| :-: | --- | --- |
| **X-1** | Lima nilai **login operator** dari predikat `IsGroup` / `IsGroupCreate` / `IsSPVCreate` tersalin di 7 berkas `dastin/_migration-docs/` (dicari dengan nilai dari korpus, tanpa mencetaknya) — melanggar `CLAUDE.md` §4 butir 10. Di `nbfacin` sudah dibersihkan (A14) | pemilik folder membersihkan; agent tidak menyunting folder milik orang lain |

---

# 🚀 PAKET SIAP KIRIM — status **21 September 2026**

> Bagian ini **menggantikan** ringkasan 17 September di bawahnya sebagai daftar kerja. Ringkasan lama
> **tidak dihapus** (`PANDUAN-KERJA` §7) — ia tetap merekam keadaan saat itu.
>
> Daftar ini diturunkan **dari status yang tertulis di dokumen ini sendiri**, lalu tiap butir
> diperiksa ulang terhadap `D:\migrasi\RNM\DDL\` dan korpus. Kolom *dasar* menyebut kalimat dokumen
> atau pengukuran yang mendasarinya.

| # | Butir terbuka | Pemilik | Memblokir |
| :-: | --- | --- | --- |
| **T-1** | Isi tabel `M_LINK_SERVICE` — ⛔ **hanya** kolom `URL`, `KATEGORI_1`, `KATEGORI_2` | **DBA** | **F13** (cetak RI Slip + email) · seluruh pemanggilan endpoint |
| **T-2** | Dua pertanyaan `HISTORYAKSEPTASIPEGA` (pemangkasan berkala · pengisi `ID_KOMITE`) | **DBA** | jalur banding · flag reject |
| **T-3** | Konfirmasi status **8 nama** yang tiba sebagai `CREATE TABLE`, bukan kode prosedur | **DBA** | ketepatan daftar A.1 — **bukan** jalur produksi |
| **T-4** | Isi baris tabel `OPENPROTEKSI_EDM` | **DBA** | klep 4 gerbang EDM (**K-049**) |
| **T-5** | **10** rule `DecisionTable` diekspor ulang beserta seluruh barisnya | **IT** | arah keputusan UW · `10-audit\08` §4 butir 5 |
| **T-6** | **10** rule golongan **G2** diekspor ulang | **IT** | F10 · F11 · rekonsiliasi |
| **T-7** | Konfirmasi apakah `CountRateRetroCov` kelas `ASM-FW-GISFW-Data-Cargo` pernah ada | **IT / admin Pega** | **K-060** pertanyaan terbuka 1 · **F06** |
| **T-8** | Satu ekspor **rule** produksi pada satu titik waktu (P-5) | **IT** | 11 activity yang tampak hilang · 6 cabang tertangguh |
| **T-9** | Arti `BusinessCode` (98 kode) + `BusinessOldId` (87 kode) | **Product** | klasifikasi produk · layar input |
| **T-10** | ⭐ **BARU 25 Sept (K-070).** Nilai **`PRODKE` tertinggi** di `JSON_POLIS` + berapa baris yang dua digit | **DBA** | aturan pemilihan "versi terakhir" pada program pemuatan (**J-2**) |

**Yang TIDAK perlu dikirim lagi** — tertutup, dicatat agar tidak diminta dua kali:
**P-1** (DDL + kode prosedur) · **P-2** · **P-4** (dropdown) · **P-7** (enam berkas enumerasi, ditutup
**K-029**) · **P-8** · ~~**P-9**~~ (dibatalkan, K-032) · **P-10** kecuali T-4 · **P-11** · **P-12**.

**Tindak lanjut internal (tanpa pihak luar):** baca `DDL\TABLEOFLIMIT.xls` — lihat catatan P-4.

📌 **Daftar periksa UI Pega terpisah** ada di `10-audit\11-daftar-periksa-ui-pega.md` — pertanyaan yang
**hanya** dapat dijawab dengan membuka rule di layar Pega, bukan dengan ekspor atau query.

---

## RINGKASAN — status per 17 September 2026 *(digantikan bagian di atas; dipertahankan sebagai jejak)*

| Paket | Pihak | Status | Membuka modul |
| --- | --- | --- | --- |
| P-1 DDL produksi + lookup | DBA | ✅ masuk (26/26) | `models`, `repository` |
| P-2 Pengukuran D1/D2 | DBA | ✅ hasil masuk — Kiro perlu baca | konfirmasi K-018 + mata uang |
| P-3 Riwayat akseptasi | DBA | 🟡 DDL masuk; 2 tanya sisa | jalur banding + flag reject |
| P-4 Isi dropdown | DBA | 🟡 tuntas kecuali re-save `TABLEOFLIMIT` | layar input |
| P-5 Ekspor produksi tunggal | IT | 🟡 5 sampel kasus masuk; ekspor rule menyusul | Special Acceptance + 6 cabang |
| P-6 DecisionTable | IT | 🟡 2 masuk (`IsUWAccepted`,`isApproved`) | arah keputusan UW |
| P-7 Arti enumerasi | Product | 🟡 6 berkas masuk — Kiro perlu baca | klasifikasi produk, layar input |
| P-8 Mesin spreading | IT + UW | ✅ masuk + masih dipakai | `services/spreading` |
| ~~P-9 Ekspor rule RNW lengkap~~ | ~~IT~~ | ⛔ **DIBATALKAN** — dijawab work owner (K-032 tertutup) | — |

**Yang benar-benar tersisa dari pihak luar:** re-save `TABLEOFLIMIT` (DBA), 2 pertanyaan
`HISTORYAKSEPTASIPEGA` (DBA), ekspor rule produksi bila diperlukan verifikasi activity hilang (IT).

✅ **P-9 tidak jadi dikirim** — ketiga rujukan menggantung RNW ternyata bukan ekspor tidak lengkap,
melainkan dua rule usang dan satu yang sudah ditambahkan ke `DDL\`.
**Tindak lanjut Kiro (tanpa pihak luar):** baca `P-2 D1/D2.xls` + 6 berkas enumerasi P-7.

---

# A. UNTUK DBA

## T-1 · Isi tabel `M_LINK_SERVICE` — ⬜ **BARU, belum pernah diminta** (21 Sept 2026)

**Yang kami minta:** isi baris tabel `POOLDATA.M_LINK_SERVICE`, **hanya tiga kolom**:

```sql
SELECT URL, KATEGORI_1, KATEGORI_2
FROM   POOLDATA.M_LINK_SERVICE
ORDER  BY KATEGORI_1, KATEGORI_2;
```

⛔ **`USERNAME` dan `PASSWORD` JANGAN disertakan.** Kami tidak meminta dan tidak akan menyimpannya.
Endpoint dan kredensial di sistem baru berasal dari **konfigurasi/env var**, tidak pernah literal di
kode (`CLAUDE.md` §4.4). Yang kami perlukan hanya **daftar kunci dan URL** agar bentuk konfigurasinya
dapat dirancang.

**Mengapa:** `[terverifikasi]` `DDL\M_LINK_SERVICE.txt` (528 B) hanya memuat `CREATE TABLE` — **nol**
`INSERT`, **nol** `SELECT`. Strukturnya sudah kami punya:
`URL VARCHAR2(1000) NOT NULL`, `KATEGORI_1 VARCHAR2(100)`, `KATEGORI_2 VARCHAR2(100)`,
`USERNAME VARCHAR2(50)`, `PASSWORD VARCHAR2(50)`. **Isinya yang belum ada.**

**Memblokir:** **F13** (cetak RI Slip + kirim email) dan setiap pemanggilan layanan luar, karena
endpoint di-resolve lewat pasangan `KATEGORI_1` + `KATEGORI_2` saat runtime.

```powershell
# menghasilkan 528 B / INSERT=0 / SELECT=0 / CREATE TABLE=True
$t=[IO.File]::ReadAllText('D:\migrasi\RNM\DDL\M_LINK_SERVICE.txt')
"byte=$($t.Length) INSERT=$(([regex]::Matches($t,'(?i)\bINSERT\s+INTO\b')).Count)" +
" SELECT=$(([regex]::Matches($t,'(?i)\bSELECT\b')).Count)" +
" CREATETABLE=$(([regex]'(?i)\bCREATE\s+TABLE\b').IsMatch($t))"
```

---

## T-3 · Konfirmasi **8 nama** yang tiba sebagai tabel, bukan kode prosedur — ⬜ **BARU** (21 Sept 2026)

**Yang kami minta:** satu konfirmasi tertulis — apakah kedelapan nama di bawah **ada sebagai
PROCEDURE/FUNCTION** di `POOLDATA`, atau memang **hanya tabel**?

`FACINOFFER` · `FACINLIFE` · `FACINSPREADLIFE` · `FACINPERFORMANCE` · `ERRORFACINPROD` ·
`TREATYPRODUCTION_BACKUP` · `HISTORYAKSEPTASIPRODUCTION` · `MONITORING_PROD_LOG`

```sql
SELECT OBJECT_NAME, OBJECT_TYPE
FROM   ALL_OBJECTS
WHERE  OWNER = 'POOLDATA'
AND    OBJECT_NAME IN ('FACINOFFER','FACINLIFE','FACINSPREADLIFE','FACINPERFORMANCE',
                       'ERRORFACINPROD','TREATYPRODUCTION_BACKUP',
                       'HISTORYAKSEPTASIPRODUCTION','MONITORING_PROD_LOG');
```

**Mengapa:** `[terverifikasi]` `_CHECKLIST-KESIAPAN.md` A.1 mencentang **33 nama** sebagai "isi kode
prosedur masuk". Dari 33 berkas `DDL\*.txt` itu, **25 memuat `CREATE … PROCEDURE/FUNCTION`** dan
**8 hanya memuat `CREATE TABLE`**. `[terverifikasi]` kedelapannya juga **bukan** rule `RDBList` di
korpus (0 dari 8 di ketiga folder), jadi **dugaan paling wajar adalah keduanya memang tabel** dan
daftar 35 semula terlalu luas — ⚠️ tetapi itu **`[dugaan]`**, dan satu baris jawaban menutupnya.

⛔ **Ini BUKAN blocker jalur produksi.** Jalur produksi tidak menunggu butir ini; yang dikoreksi hanya
ketepatan daftar A.1.

```powershell
# menghasilkan 33 / 25 / 8
$rxP=[regex]'(?i)\bCREATE\s+(OR\s+REPLACE\s+)?((NON)?EDITIONABLE\s+)?(PROCEDURE|FUNCTION|PACKAGE)\b'
# CATATAN: kata kunci EDITIONABLE berada DI ANTARA "REPLACE" dan "PROCEDURE".
#   Regex tanpa bagian itu melaporkan 9, bukan 25 - kekeliruan ukur, bukan berkas hilang.
```

---

## T-10 · Nilai `PRODKE` tertinggi di `JSON_POLIS` — ⬜ **BARU** (25 Sept 2026)

**Yang kami minta:** tiga angka. **Tanpa nomor polis, tanpa data pelanggan.**

```sql
SELECT MAX(TO_NUMBER(PRODKE)) AS prodke_tertinggi,
       COUNT(CASE WHEN LENGTH(TRIM(PRODKE)) >= 2 THEN 1 END) AS baris_dua_digit,
       COUNT(*) AS total_baris
FROM   POOLDATA.JSON_POLIS
WHERE  PRODKE IS NOT NULL;
```

**Mengapa:** `[terverifikasi]` `DDL\JSON_POLIS.txt` baris 7 — **`"PRODKE" VARCHAR2(5)`**, yaitu
**teks, bukan angka**. Program pemuatan harus memilih "versi terakhir" tiap polis; bila ia
mengurutkan `PRODKE` sebagai teks, **`'9'` terbaca lebih besar daripada `'10'`** dan generasi yang
dimuat menjadi **salah secara diam-diam**. Jebakan yang sama bentuknya dengan
`FACINOFFER.RATE VARCHAR2(100)` (`CLAUDE.md` §4.1).

⚠️ `[terverifikasi]` `PRODKE` **nol kemunculan** pada 115 contoh `DDL\CONTOH\` — ia hanya hidup
sebagai kolom Oracle, sehingga **korpus tidak dapat menjawab ini sendiri**.

📌 Bila `prodke_tertinggi` **≤ 9**, jebakan itu belum pernah terwujud dan risikonya rendah untuk data
yang ada sekarang — ⛔ tetapi program pemuatan **tetap** harus mengurutkan secara numerik, karena
data akan terus bertambah.

⚠️ **Bobot turun sesudah K-071 (25 Sept).** Work owner menetapkan `PRODKE` dan `TGL_INPUT`
**sama-sama tertinggi**, sehingga loader mengurutkan dengan `TGL_INPUT` (bertipe `DATE`, bebas
jebakan) dan memakai `PRODKE` hanya sebagai **pemeriksa silang**. ⛔ Butir ini **tidak lagi
memblokir**; angkanya tetap berguna untuk mengetahui seberapa sering jebakan urut-teks berpeluang
terwujud, dan untuk menguji apakah anggapan "keduanya sama-sama tertinggi" bertahan.

---

## P-1 · Isi kode prosedur + DDL tabel produksi  ✅ MASUK (status 17 Sept 2026)

### P-1a · DDL tabel penyimpanan/produksi — ✅ 11 dari 11 masuk
`JSON_POLIS`, `JSON_OFFER`, `FACINPRODUCTION`, `FACOUTPRODUCTION`, `FACINPRODUCTION_BACKUP`,
`TREATYINPRODUCTION`, `TREATYINPRODUCTION_BACKUP`, `CEDING_FACINPRODUCTION`, `JSON_POLIS_MONITORING`,
`JSON_KLAIM`, `C_COUNTER_PRODKE` — semua ada di `DDL\`.

### P-1b · DDL tabel lookup/master — ✅ 15 dari 15 masuk
`M_CURRENCYSTANDARD`, `M_CLIENT`, `M_CITY`, `M_DISTRICT`, `M_NATION`, `M_PROVINCE`, `RISKADDRESS`,
`SHIP`, `MARKETINGOFFICER`, `M_ACCUMULATION`, `M_ACCUMULATEDTYPE`, `M_ACCUMULATION_LIFE`,
`M_TREATY_IN`, `M_SITE_DATABASE`, `M_LINK_SERVICE` — semua ada di `DDL\`.

> ⚠️ Yang tersisa dari jalur produksi bukan DDL, melainkan **kode `ALL_SOURCE` 32 prosedur** (sudah
> masuk sebelumnya). DDL tipe kolom kini lengkap → `models` terbuka.
> Catatan: `TREATYINPRODUCTION*` masuk; ingat cakupan Treaty terbatas (K-005) saat porting.

> ### ✅ Ditegaskan ulang 21 September 2026 — `ALL_SOURCE` **BUKAN** butir terbuka
>
> `[terverifikasi]` Ketiga-puluh-tiga nama yang dicentang `_CHECKLIST-KESIAPAN.md` A.1 **seluruhnya
> ada** sebagai berkas di `D:\migrasi\RNM\DDL\*.txt` (**33 ada, 0 hilang**); **25** memuat kode
> `PROCEDURE`/`FUNCTION`, **8** hanya `CREATE TABLE` (→ **T-3**). Di seluruh 80 berkas `DDL\*.txt`:
> **25 kode prosedur/fungsi · 46 `CREATE TABLE` · 9 lain**.
>
> ⛔ **Dokumen lain yang masih menyebut `ALL_SOURCE` sebagai penunggu jalur produksi sudah basi**
> (mis. `05-tickets\00-INDEKS.md`, `08-flat\01`, beberapa berkas `04-spec\`). Yang benar-benar menahan
> jalur produksi adalah **struktur tabel flat** (keputusan work owner), **bukan** `ALL_SOURCE`.
> Butir ini **tidak dikirim ke DBA**.

## P-2 · Pengukuran D1 & D2  ✅ HASIL MASUK & DIBACA (status 17 Sept 2026)

Hasil dilampirkan `DDL\D1.xml` & `DDL\D2.xml` — sudah dibaca Kiro:
- **D1** `[terverifikasi]`: total **691.925** baris, **4** tanpa mata uang (0,0006%). Menguatkan
  K-012 (`Unknown` keadaan langka-tapi-nyata; jangan default ke IDR).
- **D2** `[terverifikasi]`: sebaran RATE campuran — banyak nilai < 0,1 (khas ‰) berdampingan dengan
  nilai ~0,5–1,4 (khas %). Menguatkan K-018 (satuan berbeda per COB, bukan seragam). D2 belum
  dipecah per COB, jadi ini konfirmasi silang, bukan pemetaan baru.
- ⚠️ **Temuan baru:** RATE disimpan **pakai koma desimal** (`0,0244`). Parser Go wajib menangani
  koma, bukan titik. Kandidat bug port — perlu dicatat sebagai kontrak parsing.

**D1 — berapa banyak nilai uang tersimpan tanpa mata uang:**
```sql
SELECT COUNT(*) AS total,
       SUM(CASE WHEN CURR_ID IS NULL OR TRIM(CURR_ID) = '' THEN 1 ELSE 0 END) AS tanpa_mata_uang
FROM   POOLDATA.FACINOFFER;
```

**D2 — sebaran nilai RATE (membuktikan satuan ‰ vs %):**
```sql
SELECT RATE, COUNT(*) AS jumlah
FROM   POOLDATA.FACINOFFER
WHERE  RATE IS NOT NULL
GROUP  BY RATE
ORDER  BY COUNT(*) DESC
FETCH FIRST 30 ROWS ONLY;
```
> Catatan: `FACINOFFER.RATE` bertipe `VARCHAR2(100)` (teks) — mohon kirim apa adanya, termasuk
> koma/titik desimalnya.

## P-3 · Riwayat akseptasi `HISTORYAKSEPTASIPEGA`  🟡 DDL MASUK (status 17 Sept 2026)

- [x] **DDL** `HISTORYAKSEPTASIPEGA` — masuk di `DDL\HISTORYAKSEPTASIPEGA.txt`
- [ ] Konfirmasi apakah tabel diarsip/dipangkas berkala
- [ ] Penjelasan **siapa/proses apa yang mengisi kolom `ID_KOMITE`** (bila kolom itu ada & relevan)

> Catatan: `DATAPEGA.pc_History_ASM_FW_GISFW_Work` **tidak** diperlukan (musnah bersama Pega).

> ### ⬜ MASIH PERLU DIKIRIM — ditegaskan 21 September 2026 (= **T-2**)
>
> Kedua kotak di atas **masih kosong**; tidak ada jawaban yang tercatat di dokumen ini maupun di
> `DDL\`. Syarat "bila kolom itu ada & relevan" kini **terpenuhi**:
>
> `[terverifikasi]` `DDL\HISTORYAKSEPTASIPEGA.txt` (572 B) memuat **7 kolom** —
> `ID_PEGA`, `TGL_TRANSFER`, `STATUS`, `USERNAME`, `WORKBASKET`, **`ID_KOMITE`**, `OPERATORID`.
> Kolom `ID_KOMITE` **ADA**, jadi pertanyaan kedua berlaku dan tetap terbuka.
>
> ⛔ **Kami tidak meminta isi kolom `USERNAME` maupun `OPERATORID`** — hanya penjelasan proses yang
> mengisi `ID_KOMITE`. Nama kolom disebut; nilainya tidak diminta dan tidak akan disalin.
>
> **Memblokir:** jalur banding dan flag reject dibaca dari tabel ini untuk **mengambil keputusan**,
> bukan sekadar audit (`CLAUDE.md` §4.3) — jadi perilaku pemangkasan berkala mengubah rancangan.

## P-4 · Isi dropdown  🟡 → hampir tuntas (status diperbarui 17 Sept 2026)

Status per sumber (dicek terhadap `D:\migrasi\RNM\DDL\`):

**✅ Sudah masuk:**
- Wilayah: `CITY`, `DISTRICT`, `RW`
- Master bisnis: `LST_BANK_GROUP`, `OCCUPATION`, `BRANDDETAIL`, `M_KLAUSUL_PLAN`, `M_TYPE_PROPERTY_PLAN`

**⛔ Dikonfirmasi TIDAK ADA lagi (kelewat saat penghapusan — bukan tertunggak):**
- `V_ZIPCODE` · `V_COINS` · `VIEW_BENEFIT_PROPERTY`
- `VJ_M_TYPE_PROPERTY_PLAN` — tidak dipakai lagi
- Tabel enumerasi utama `Data-Enumeration` — **tidak dipakai lagi**; bila masih ada panggilan di
  Section, itu sisa yang kelewat saat penghapusan. Tidak perlu diekstrak.

**✅ `TABLEOFLIMIT` — struktur XML masuk (17 Sept 2026):** `DDL\TABLEOFLIMIT.xml` terbaca.
`[terverifikasi]` kolom: `ID`, `BIZCODE`, `NOTE`, `TAHUN`, `CATEGORY`, `DESCRIPTION`, `PCTLIMIT`.
Berisi 1 record contoh (sampel struktur, bukan isi penuh). Cukup untuk merancang bentuk flat.
Catatan: `PCTLIMIT` memakai koma desimal (`70,000`) — sama seperti `RATE`, perlu parser koma.

**Tidak ada lagi butir dropdown yang tertunggak.**

> ### ✅ P-4 TERTUTUP untuk DBA — 21 September 2026 · **tindak lanjut internal tersisa**
>
> Hipotesis "P-4 re-save `TABLEOFLIMIT` masih perlu dikirim" **tidak didukung dokumen ini sendiri**:
> badan P-4 sudah berbunyi *"✅ `TABLEOFLIMIT` — struktur XML masuk (17 Sept 2026)"* dan
> *"Tidak ada lagi butir dropdown yang tertunggak."* Baris 🟡 di ringkasan 17 September **basi**.
>
> `[terverifikasi]` `DDL\TABLEOFLIMIT.xml` (651 B) memuat **tepat satu** `<DATA_RECORD>` — sesuai
> keterangan "1 record contoh (sampel struktur, bukan isi penuh)".
>
> ⚠️ **Yang belum tercatat di dokumen ini:** `DDL\TABLEOFLIMIT.xls` **juga ada** — **18.944 B**,
> berkas OLE Excel asli (tanda tangan `D0 CF 11 E0`), berdampingan dengan tujuh `.xls` tabel limit
> lain yang memang berisi **isi penuh**. `[dugaan]` isi penuh `TABLEOFLIMIT` **sudah masuk** lewat
> berkas itu dan hanya **belum dibaca**.
>
> ⛔ **Jangan minta ulang ke DBA sebelum berkas itu dibaca.** Ini **tindak lanjut internal**, bukan
> permintaan pihak luar.
>
> ```powershell
> # 651 B, 1 DATA_RECORD  |  .xls 18944 B, sig D0 CF 11 E0 (OLE compound)
> $x=New-Object Xml.XmlDocument; $x.Load('D:\migrasi\RNM\DDL\TABLEOFLIMIT.xml')
> ($x.DocumentElement.ChildNodes | Where-Object { $_.Name -eq 'DATA_RECORD' }).Count
> ([IO.File]::ReadAllBytes('D:\migrasi\RNM\DDL\TABLEOFLIMIT.xls')).Length
> ```

---

# B. UNTUK IT

## P-5 · Ekspor produksi tunggal  🟡 SAMPEL MASUK (status 17 Sept 2026)

Lima berkas kasus produksi sudah masuk di `DDL\`, mencakup ketiga siklus + beberapa COB:
`EDM-13445 (FIRE)`, `NB-176005 (AS. KREDIT)`, `NB-181231 (FIRE)`, `NB-184233 (MARINE CARGO)`,
`RNW-10579 (FIRE)`.

- [ ] Bila memungkinkan: satu ekspor **rule** produksi (bukan hanya kasus) pada satu titik waktu,
      untuk memverifikasi 11 activity yang tampak hilang + 6 cabang tertangguh (Special Acceptance,
      tangga putaran kedua). Sampel kasus di atas berguna untuk rekonsiliasi, tetapi verifikasi
      activity "hilang" tetap butuh ekspor rule.

> ### ⬜ MASIH PERLU DIKIRIM — 21 September 2026 (= **T-8**), tetapi **prioritasnya turun**
>
> Kotaknya masih kosong dan tidak ada ekspor rule produksi yang tercatat masuk, jadi butir ini
> **tetap terbuka**.
>
> ⚠️ **Premisnya berubah.** Sebagian dari "activity yang tampak hilang" berdiri di atas angka
> **95 rule berbeda versi**, yang kini diketahui **mencampur dua hal berbeda**: drift sejati (**G2**)
> dan **rule kembar** (**G3** — nama sama, kelas berbeda, keduanya sah). Rincian:
> `10-audit\10-rule-kembar-pzinskey.md` dan `_ARSIP-lintas-siklus\_BACA-INI.md` Koreksi R1.
>
> ⛔ **Jangan disimpulkan bahwa P-8 gugur** — ia tidak terbantah, hanya **tidak lagi mendesak**.
> Bentuk minimum yang benar-benar dibutuhkan sekarang adalah **T-6** (10 rule G2) dan **T-5**
> (10 DecisionTable) di bawah, yang jauh lebih kecil dan dapat dipenuhi lebih dulu.

## P-6 · DecisionTable beserta seluruh barisnya  🟡 SEBAGIAN MASUK (status 17 Sept 2026)

- [x] `IsUWAccepted.xml` — masuk (18 KB, berisi)
- [x] `isApproved.xml` — masuk (16 KB, berisi)
- [ ] DecisionTable lain (bila masih ada yang relevan) — menyusul bila diperlukan

> ### 📝 Dikoreksi 21 September 2026 — **4 masuk, bukan 2**; **10 tersisa** (= **T-5**)
>
> `[terverifikasi]` `DDL\` memuat **empat** rule `DecisionTable`, bukan dua: `IsUWAccepted.xml`,
> `isApproved.xml`, dan — **belum pernah dicatat di dokumen ini** — **`InputEDMLife.xml`** dan
> **`OfferFacRetro.xml`**. Kedua yang terakhir **tidak ada** di folder `DecisionTable\` korpus mana
> pun, jadi keduanya tambahan murni.
>
> **Daftar lengkapnya ada di T-5** pada bagian D di bawah.

## P-8 · Mesin spreading terbesar  ✅ MASUK + dikonfirmasi (status 17 Sept 2026)

- [x] `GetKapasitasTreaty.xml` (±1,02 MB) sudah **diekspor ulang** dan ada di `DDL\`
- [x] **Dikonfirmasi work owner: masih dipakai** → diport apa adanya (reproduksi perilaku, CLAUDE.md §1)

---

## P-9 · ~~Ekspor rule **RNW (Renewal)** yang lengkap~~  ⛔ **DIBATALKAN — JANGAN DIKIRIM**

> ⛔ **Gugur 18 September 2026 — dijawab work owner sebelum sempat dikirim.**
> Lihat **K-032 (tertutup)** di `00-KEPUTUSAN-WORK-OWNER.md`.
>
> Dugaan "ekspor RNW tidak lengkap" **tidak terbukti**. Ketiganya punya sebab masing-masing:
>
> | Rule | Jawaban |
> | --- | --- |
> | `GenerateNoPolicy` | **usang** — simpan produksi NB & renewal memakai activity yang sama (`SaveJsonPolicyFacIn_Act`, `[terverifikasi]` identik byte-per-byte). Panggilan tersisa = sisa kelewat |
> | `SumTreatyCapacity_Act` | ✅ **sudah ditambahkan ke `DDL\`** — tidak menggantung lagi |
> | `SearchJobID` | **usang** — sudah lama dihapus. Panggilan tersisa = sisa kelewat |
>
> **Tidak ada yang perlu diminta ke IT.** Badan permintaan di bawah dipertahankan sebagai jejak, bukan
> untuk dijalankan.

### ~~Yang kami minta~~ *(arsip)*

**Ekspor rule siklus Renewal (`RNW Fac In`) secara LENGKAP dari lingkungan produksi, pada satu titik
waktu** — analog ekspor NB yang sudah kami terima.

Tujuannya satu dan spesifik: memastikan apakah **tiga rule di bawah memang tidak ada di ruleset
renewal**, atau **sekadar tidak ikut terekspor**.

**Definisi rule saja. Tanpa data pelanggan.**

### Mengapa kami menduga ekspornya tidak lengkap

Ini bukan tuduhan — ini kesimpulan sementara dari satu angka yang mengejutkan kami sendiri.

`[terverifikasi]` Dari **1.927** berkas di ekspor RNW, **1.907 identik byte-per-byte** dengan berkas
bernama sama di ekspor NB — **nol berbeda**, metadata ekspor termasuk. Artinya kedua ekspor berasal
dari **keadaan ruleset yang sama**.

Justru karena itu, sebuah berkas yang hadir di NB tetapi tidak di RNW menjadi **lebih** mencurigakan
sebagai kelalaian ekspor, bukan kurang.

⛔ Kami **belum menyimpulkan** ketiganya usang. Aturan proyek kami melarang menyimpulkan "tidak ada di
ekspor" sebagai "tidak dipakai" — justru permintaan inilah yang menjawabnya.

### Tiga rule yang dirujuk tetapi berkasnya tidak ada di ekspor RNW

Ketiganya **ada di ekspor NB**, dan **dirujuk dari berkas yang ada di ekspor RNW**.

---

#### 1. `GenerateNoPolicy` — jenis rule **RDBList** ⛔ **PALING BERDAMPAK**

**Ini penghasil nomor polis.** Dirujuk lewat `<RequestType>` oleh dua activity yang ada di RNW:

| Perujuk (ada di RNW) | Baris | Isi tag |
| --- | ---: | --- |
| `Activity\SaveJsonPolicyFacIn_Act` | 2204 | `<RequestType>GenerateNoPolicy</RequestType>` |
| `Activity\GenerateNopolis_Act` | 3179 | `<RequestType>GenerateNoPolicy</RequestType>` |

**Akibat bila benar-benar tidak ada:** jalur simpan produksi renewal kehilangan penghasil nomor
polisnya. Ini yang paling kami perlukan lebih dulu.

---

#### 2. `SumTreatyCapacity_Act` — jenis rule **Activity**

Dirujuk lewat langkah `Call` oleh dua activity yang ada di RNW:

| Perujuk (ada di RNW) | Baris | Isi tag |
| --- | ---: | --- |
| `Activity\SetValidateDate_Act` | 1129 | `<pyStepsActivityName>Call SumTreatyCapacity_Act</pyStepsActivityName>` |
| `Activity\CountASMShareTotal_ACT` | 444 | `<pyStepsActivityName>Call SumTreatyCapacity_Act</pyStepsActivityName>` |

**Akibat bila benar-benar tidak ada:** cabang perhitungan kapasitas treaty terputus di siklus renewal.

---

#### 3. `SearchJobID` — jenis rule **RDBList**

Dirujuk lewat `<RequestType>` oleh satu activity yang ada di RNW:

| Perujuk (ada di RNW) | Baris | Isi tag |
| --- | ---: | --- |
| `Activity\UploadCSVPerson_PostAct` | 2404 | `<RequestType>SearchJobID</RequestType>` |

*Butir ini berasal dari pemeriksaan discovery dan **belum diverifikasi ulang** secara terpisah —
kami cantumkan agar sekalian terjawab.*

---

### Bentuk jawaban yang cukup

Salah satu dari dua ini sudah memadai:

1. **Ekspor rule RNW lengkap** (berkas `.xml`), atau
2. **Ketiga rule di atas saja**, beserta **konfirmasi tertulis** apakah masing-masing ada di ruleset
   renewal.

Bila salah satu memang **tidak ada** di ruleset renewal, konfirmasi itu sendiri sudah menjawab — kami
tidak memerlukan berkasnya.

# C. UNTUK PRODUCT

## P-7 · Arti kode enumerasi  🟡 BERKAS MASUK (status 17 Sept 2026)

Berkas sudah dilampirkan di `DDL\`: `Type.xml`, `EdmType.xml`, `TypeDeductible.xml` +
`TypeDeductible2.xml`, `TeamGroup.xml`, `ProRateType.xml`, `B2B.xls`.

- [x] Berkas keenam enumerasi masuk (Type, EdmType, DeductibleType×2, TeamGroup, ProRateType, IsB2B)
- [x] **`B2B.xml` dibaca** `[terverifikasi]`: nilai `IsB2B` = `ASM`, `KBRU`, `BDX`, `SRB`. Karena
      `ASM` valid, rule `IsPKSASM` (`IsB2B="ASM"`) **dapat berhenti `panic`** — kondisinya kini
      punya nilai referensi nyata. → kandidat keputusan penutup.
- [x] **Tindak lanjut Kiro:** baca `Type.xml`/`EdmType.xml`/`TypeDeductible*.xml`/`TeamGroup.xml`/
      `ProRateType.xml` (arti kode) → catat keputusan bila menutup `[pertanyaan terbuka]`.
      ✅ **SELESAI** — ditutup **K-029** (17 September 2026); kosakatanya di `steering\GLOSARIUM.md`.

> ### ✅ P-7 TIDAK PERLU DIKIRIM LAGI — 21 September 2026 · **satu sisa, butir baru**
>
> Keenam berkas enumerasi sudah masuk **dan sudah dibaca**; artinya ditutup **K-029**. `GLOSARIUM.md`
> menyatakannya: *"Enumerasi — artinya terverifikasi, seluruhnya tertutup (17 September 2026)"*.
> **Tidak ada yang perlu diminta ulang untuk keenam berkas itu.**
>
> ⬜ **Yang tersisa untuk Product adalah butir lain** (= **T-9**), yang dicatat `GLOSARIUM.md` sebagai
> satu-satunya yang belum terverifikasi:
>
> | Istilah | Nilai terbaca | Arti | Memblokir |
> | --- | ---: | :-: | --- |
> | `BusinessCode` | **98** kode | ☐ | klasifikasi produk · layar input |
> | `BusinessOldId` | **87** kode | ☐ | idem |
>
> **Bentuk jawaban yang cukup:** tabel `kode → arti` untuk keduanya. Tanpa data pelanggan.
> **Jangan ditebak** (`CLAUDE.md` §3 butir 4) — sampai terjawab, keduanya tetap `panic` bila tercapai.

---

# D. UNTUK IT / ADMINISTRATOR PEGA — **ekspor ulang rule** (21 September 2026)

> Bagian **baru**. Tiga butir, semuanya berupa **definisi rule** atau **konfirmasi tertulis**.
> **Tanpa data pelanggan, tanpa kredensial.**

## T-5 · Sepuluh `DecisionTable` diekspor ulang beserta **seluruh barisnya**

`[terverifikasi]` Korpus memuat **12 identitas `DecisionTable` unik** (nama + `pyClassName`). Dua di
antaranya sudah ada di `DDL\`. **Sepuluh sisanya belum:**

| # | Rule | `pyClassName` | Ada di folder |
| ---: | --- | --- | --- |
| 1 | `MappingCoverage` | `ASM-FW-GISFW-Work` | NB · RNW · EDM |
| 2 | `MappingAdditionalCoverageIndex` | `ASM-FW-GISFW-Data-Coverage` | NB · RNW · EDM |
| 3 | `MappingOutgoIndex` | `ASM-FW-GISFW-Data-Split` | NB · RNW · EDM |
| 4 | `MappingOutgoIndex2` | `ASM-FW-GISFW-Data-Split` | NB · RNW · EDM |
| 5 | `SetUploadHubAW1` | `ASM-FW-GISFW-Data-BatchPAList` | NB · RNW · EDM |
| 6 | `SetUploadHubAW2` | `ASM-FW-GISFW-Data-BatchPAList` | NB · RNW · EDM |
| 7 | `SetUploadHubAW3` | `ASM-FW-GISFW-Data-BatchPAList` | NB · RNW · EDM |
| 8 | `LicensePlatRegion_DeT` | `ASM-FW-GISFW-Data-Vehicle` | NB · RNW · EDM |
| 9 | `GetMimeType` | `ASM-FW-GISFW-Int-T_STORAGE_IMAGE` | NB · RNW · EDM |
| 10 | `BusinessType_DeT` | `ASM-FW-GISFW-Work` | **NB saja** |

⛔ **Sertakan `pyClassName` pada tiap ekspor.** Nama berkas saja **tidak** menentukan identitas rule —
dua berkas bernama sama di kelas berbeda adalah **dua rule** (Koreksi R1).

**Mengapa:** `10-audit\07` mencatat tipe rule ini **100 % berbeda antar folder** (10 dari 10 identitas
bersama), dan `10-audit\08` §4 butir 5 menandainya **belum dipilah G2/G3**. Ekspor ulang menutup
keduanya sekaligus: yang sama kelasnya → drift yang harus diselesaikan; yang beda kelasnya → rule
kembar yang **keduanya diport**.

**Memblokir:** arah keputusan UW (`IsUWAccepted`/`isApproved` sudah aman) dan pemetaan coverage /
spreading yang dipakai F04, F05, F08.

## T-6 · Sepuluh rule golongan **G2** diekspor ulang

`[terverifikasi]` `10-audit\08` §3A memuat **21 baris G2** atas **11 berkas unik**. Satu di antaranya
— `InsertFacoutProduction` — **sudah tiba** di `DDL\` (`[terverifikasi]` kelas
`ASM-FW-GISFW-Data-FacOffer`, basis `pzInsKey` sama dengan ketiga salinan korpus, commit
`20260831T101051` = **paling baru**). **Sepuluh sisanya:**

| # | Rule | Tipe | `pyClassName` (sama di kedua salinan) | Kelas beda |
| ---: | --- | --- | --- | :-: |
| 1 | **`CountEdmAdjTSI_Act`** | Activity | `ASM-FW-GISFW-Data-Aneka` | **a — uang** ⛔ prioritas |
| 2 | `Protection_Act` | Activity | `ASM-FW-GISFW-Work` | b — gerbang |
| 3 | `serviceInsertArasapasEDM_act` | Activity | `ASM-FW-GISFW-Work` | b — gerbang |
| 4 | `SaveFacinProdFireNB_Act` | Activity | `ASM-FW-GISFW-Work` | c — produksi |
| 5 | `SaveFacinProdEDMMarineCargo_Act` | Activity | `ASM-FW-GISFW-Work` | c — produksi |
| 6 | `InsertFacoutProductionEDM` | Activity | `ASM-FW-GISFW-Data-FacOffer` | c — produksi |
| 7 | `InsertFacoutProductionEDMCurr` | Activity | `ASM-FW-GISFW-Data-FacOffer` | c — produksi |
| 8 | `DeleteDataProduction` | RDBList | `ASM-FW-GISFW-Int-policyjson` | c — produksi |
| 9 | `InsertHistoryAkseptasiPega_Sql` | RDBList | `ASM-FW-GISFW-int-policyjson` ⚠️ **`i` kecil** | c — produksi |
| 10 | `INSERTJSON_JSONPOLISMONITORING_FACIN` | RDBList | `ASM-FW-GISFW-Int-policyjson` | c — produksi |

⚠️ Ejaan kelas baris 9 memakai **`int`** huruf kecil, berbeda dari `Int` pada dua RDBList lain.
`[terverifikasi]` perbandingannya peka huruf dan **konsisten antar folder** — **bukan** salah ketik
kami. Ejaan itu **diport apa adanya**.

⛔ **Yang TIDAK kami minta:** rule golongan **G3** (nama sama, kelas berbeda). Mengekspor ulang tidak
menyelesaikan apa pun di sana — keduanya rule berbeda dan **keduanya diport**. **39 dari 60 baris**
daftar lama sudah dikeluarkan atas dasar itu.

**Memblokir:** F10, F11, dan rekonsiliasi paralel run — `CountEdmAdjTSI_Act` menyentuh jalur uang.

## T-7 · Konfirmasi: apakah `CountRateRetroCov` kelas `ASM-FW-GISFW-Data-Cargo` **pernah ada**?

**Yang kami minta:** satu jawaban — **ada / pernah ada lalu dihapus / tidak pernah ada**. Bila ada,
mohon ekspornya.

**Mengapa:** `[terverifikasi]` indeks `Embed-Reference-Rule` pada **7 berkas** mencatat
`CountRateRetroCov` pada **tiga** kelas — masing-masing **7 kali**:

| Kelas | Berkas yang mengindeks |
| --- | ---: |
| `ASM-FW-GISFW-Data-PropertyItem` | 7 |
| `ASM-FW-GISFW-Data-Aneka` | 7 |
| **`ASM-FW-GISFW-Data-Cargo`** | **7** |

Dan tiap pemanggil produksi (`InsertFacoutProduction`, `InsertFacoutProductionEDM`) memuat **tiga**
langkah `Call CountRateRetroCov` yang seragam. ⛔ **Rule berkelas `Data-Cargo` tidak ada** di ketiga
folder korpus maupun di `DDL\`.

⚠️ **Work owner menyatakan hanya ada DUA rule.** Kami **tidak menyimpulkan** salah satu: indeks
rujukan merekam resolusi **saat berkas terakhir disimpan**, sehingga ia bukan bukti rule itu masih
ada — tetapi juga tidak menjelaskan mengapa ada **tiga** langkah `Call`.

**Pertanyaan tambahan yang sama pentingnya:** apakah langkah `Call` **ketiga** masih dapat tercapai,
dan bila tercapai, **ke rule mana** ia resolve? Rincian: `10-audit\09-pemanggil-countrateretrocov.md`.

**Memblokir:** `K-060` pertanyaan terbuka 1 · tiket **F06**.

---

---

# TAMBAHAN — BLOCKER SIKLUS ENDORSEMENT (EDM)

> Ditambahkan 19 September 2026 setelah discovery EDM (E-1…E-6) selesai. Lingkup EDM final = **708
> berkas** (354 berbeda + 354 EDM-only) setelah amandemen K-043. Butir di bawah adalah blocker jalur
> produksi & tangga akseptasi EDM yang menunggu pihak luar. Semua `[terverifikasi]` dari korpus
> `Endorsment Fac In\` kecuali yang ditandai lain.

## RINGKASAN TAMBAHAN EDM

| Paket | Pihak | Status | Membuka |
| --- | --- | --- | --- |
| P-10 Stored procedure + skema produksi EDM | DBA | ✅ **tertutup** (jawaban 19 Sept) | jalur produksi EDM → tabel flat |
| P-11 Arti `pyStepsPreCondition=false` + versi rule | IT | ✅ **tertutup** | validitas pemetaan alur EDM |
| P-12 Empat keputusan lingkup EDM | Work owner | ⬜ perlu diputuskan | lingkup spec EDM |

**Tidak ada lagi yang tersisa dari pihak luar (DBA/IT) untuk EDM.** P-10 (DBA) ✅ dan P-11 (IT) ✅
tertutup. Yang tersisa hanya **P-12 — empat keputusan work owner (W-1…W-4)** sebelum spec EDM disusun.

---

## P-10 · Stored procedure + skema tabel produksi EDM  ✅ TERTUTUP (dijawab work owner + DBA, 19 Sept 2026)

Jalur konversi produksi EDM memanggil objek yang isinya tidak ada di korpus. **Seluruh butir kini
terjawab** — tidak ada lagi yang perlu diminta ke DBA untuk EDM.

| Objek / pertanyaan | Jawaban | Status |
| --- | --- | :-: |
| `POOLDATA.INSERTJSONPOLIS` | **Tidak dipakai di sistem baru** — insert dibuat ke **tabel flat** bersama NB & RNW setelah proses ini | ✅ |
| `POOLDATA.PEGA_DELETE_ERROR_KONVERSI` | File **sudah diperbarui** — `[terverifikasi Kiro]` ada di `DDL\PEGA_DELETE_ERROR_KONVERSI.txt` (3.019 B, berisi definisi) | ✅ |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | File **sudah diperbarui** — `[terverifikasi Kiro]` ada di `DDL\PROC_GENERATE_SEQUENCE_NUMBER.txt` (3.015 B, berisi definisi) | ✅ |
| `MachingDataFacin_Sql` | **Tidak dipakai di sistem baru** | ✅ |
| Skema `facinproduction` vs `pooldata.facinproduction` (E-Q29) | **Objek yang SAMA** → guard idempotensi berfungsi, tidak ada risiko duplikasi | ✅ tutup |
| Kerapatan `PRODKE` (E-Q30) | `GetProdKeOldData_SQL` **sudah sesuai** (pakai `ORDER BY TGL_INPUT DESC`, tidak rentan lubang) | ✅ tutup |
| Tabel `OPENPROTEKSI_EDM` (`Param.Type` 1–4) | **DIKOREKSI (K-049):** rule pengecekan `BrowseOpenProteksiEdm_RD` **DIPAKAI saat buat EDM** (klep 4 gerbang) → rule masuk lingkup EDM. **Tabelnya** dibagi dengan menu lain; isi tabel di luar korpus → dari DBA/menu lain. Klep hanya tutup 4 dari 6 gerbang. | 🟡 rule masuk lingkup; isi tabel dari DBA |

> ### ⬜ Satu baris MASIH TERBUKA — ditegaskan 21 September 2026 (= **T-4**)
>
> P-10 tertutup **kecuali** baris terakhir tabel di atas, yang masih bertanda 🟡. Kalimat yang
> mendasarinya, dikutip dari baris itu: *"**Tabelnya** dibagi dengan menu lain; isi tabel di luar
> korpus → dari DBA/menu lain. Klep hanya tutup 4 dari 6 gerbang."*
>
> **Yang kami minta:** isi baris tabel `OPENPROTEKSI_EDM` untuk `Param.Type` **1–4**
> (kode + label + status aktif). Tanpa data pelanggan.
> **Memblokir:** klep 4 gerbang pembuatan EDM (**K-049**).

> Konsekuensi lingkup: jalur simpan-produksi EDM di sistem baru **tidak** memanggil stored procedure
> Pega — ia menulis ke **tabel flat** yang dirancang bersama NB & RNW (butir bagian E `_CHECKLIST`).
> `INSERTJSONPOLIS` & `MachingDataFacin_Sql` **tidak diport**; keduanya konteks Pega yang hilang
> bersama Pega (CLAUDE.md §4.3 `DATAPEGA.PC_*`). Ini **keputusan work owner**, dicatat agar tidak
> tampak sebagai penghilangan diam-diam (§1).

---

## P-11 · Arti `pyStepsPreCondition=false` + versi rule produksi  ✅ TERTUTUP (status 19 Sept 2026)

**1. Semantik `<pyStepsPreCondition>false</pyStepsPreCondition>` (E-Q16) — ✅ TERTUTUP (diverifikasi work owner via UI Pega).**
`[terverifikasi]` muncul **617 kali di 173 berkas** EDM; **47** di antaranya menggerbangi cabang lini bisnis.

**Kesimpulan: `false` TIDAK berarti langkah di-skip. Langkah tetap dijalankan.** Diverifikasi work
owner langsung di UI Pega: langkah `SaveEDMToJsonPolicy_Act` "Get ProdKe dari JSON_POLIS"
(`pyStepsPreCondition=false` di XML) tampil **aktif** di form Pega — `Loop / When / RDB-List`,
`RequestType=GetProdKeOldData_SQL` — bukan langkah tercoret/disabled.

> Arti tag: `false` di serialisasi XML hanya menandai atribut kondisi-langkah bernilai false, **bukan**
> menonaktifkan langkah. **Konsekuensi (baik): seluruh pemetaan alur EDM valid apa adanya** — tidak ada
> langkah yang ternyata mati. SetOldData mengisi `*Old`, SaveEDMToJsonPolicy mengambil ProdKe, dan
> penulisan spreading/selisih semuanya **tetap tereksekusi**.

> ⚠️ Catatan terpisah (tetap benar): `SetOldData` **hanya jalan untuk menu endorsement** — tapi itu
> karena **langkah gerbang `IF[IsEDM] then=2 else=6`** di awal activity (E-3), bukan karena
> `pyStepsPreCondition`. Dua mekanisme berbeda; keduanya kini jelas.

**2. Rule produksi EDM — ✅ TERTUTUP (dijawab work owner).** Nomor polis endorsement + insert tabel
ada di **`Activity\SaveEDMToJsonPolicy_Act`**; untuk NB & Renewal ada di
**`Activity\SaveJsonPolicyFacIn_Act`**. Ini menutup rujukan `GenerateEndorsementNo`/`GenerateEDMNoLife`:
bukan rule hilang yang perlu panic, melainkan proses generate-nopolis yang **tertanam di activity EDM
tersendiri**. Tidak ada ekspor rule tambahan yang diperlukan untuk butir ini.

> Catatan: verifikasi **versi rule** lintas-ekspor (`SetToInbox_ACT`/`SetBanding_ACT`, arsip §10 #1)
> tetap bergantung pada ekspor rule produksi tunggal **P-5** yang sudah diminta untuk NB — tidak perlu
> permintaan terpisah untuk EDM.

---

## P-12 · Empat keputusan lingkup EDM  🟡 3 DIJAWAB, W-4 menunggu (status 19 Sept 2026)

**✅ SEMUA P-12 TERTUTUP (19 Sept 2026) — dicatat di register K-044, K-045, K-046.**

| # | Keputusan | Jawaban work owner | Status |
| --- | --- | --- | :-: |
| W-1 | Lingkup Life (E-Q26/E-Q31) | **EDM Life IKUT dikerjakan** → K-044 | ✅ |
| W-2 | Pilih-tertanggung EDM (E-Q25) | **Dibuang** — tertanggung dari polis yang di-endors; K-038 berlaku EDM → K-044 | ✅ |
| W-3 | Rule tangga akseptasi di-remark (E-Q28) | **Putaran 2 (`GetLimitAkseptasi_Act2`) mati** (di-remark di kedua pemanggil); **Putaran 1 tetap hidup**; Bonding tetap terbuka → K-045 | ✅ |
| W-4 | 12 kode usang | **9 diport apa adanya, 1 diperbaiki (A.5 `PremiNusantaraReOld`→angka 0), 2 bukan keanehan** → K-046 | ✅ |

> **⚠️ W-3 — sumber bukti & konsekuensi (jujur dicatat).** Keputusan "sudah tidak dipakai" berasal
> dari **observasi work owner di UI Pega**, **bukan** dari korpus XML. Bukti: di
> `InputOfferFacInEngineerUW_preACT` (`ASM-FW-GISFW-Work`, RS GISFW:01-01-87), **step 5 & 6 berlabel
> `//` (di-remark)** — step 5 `Call GetLimitAkseptasi_Act` ("Putaran 1"), step 6
> `Call GetLimitAkseptasi_Act2` ("Putaran 2"). `[terverifikasi Kiro]` di ekspor XML langkah ini
> **tetap muncul tanpa penanda disable yang terbaca** — status remark (`//`) tidak terserialisasi jelas
> ke XML. Sama seperti kasus `pyStepsPreCondition`: **UI Pega mengalahkan pembacaan XML.**
>
> **⚠️ Remark ini LOKAL — bukan vonis usang menyeluruh (dipertegas work owner).** Remark `//` hanya
> berlaku untuk pemanggilan **di dalam `InputOfferFacInEngineerUW_preACT`**. Rule-nya bisa **tetap
> dipakai dari pemanggil lain.** `[terverifikasi Kiro]` pemanggil lain di korpus EDM:
>
> | Rule | Pemanggil lain (selain yang di-remark) | Kesimpulan aman |
> | --- | --- | --- |
> | `GetLimitAkseptasi_Act` (Putaran 1) | `InputAddendumFacIn_PreAct`, `InputAddendumFacultativeIn` (flow utama), `SetValidateDate_PostAct` | **MASIH HIDUP** — 3 pemanggil lain; jangan divonis mati |
> | `GetLimitAkseptasi_Act2` (Putaran 2) | hanya `SetValidateDate_PostAct` | perlu cek: di situ di-remark juga atau tidak |
> | `GetLimitAkseptasi1SA/2SA_Act` | `GetLimitAkseptasi_Act` | tergantung status Putaran 1/2 |
>
> **Yang boleh disimpulkan sekarang (dari UI Pega work owner):** di `InputOfferFacInEngineerUW_preACT`,
> Putaran 1 & 2 di-remark. **Yang BELUM boleh disimpulkan:** rule tangga akseptasi mati di seluruh EDM
> — `GetLimitAkseptasi_Act` justru terbukti masih dipanggil 3 tempat lain. Vonis usang menyeluruh
> butuh pengecekan setiap pemanggil di UI Pega (aturan K-006: jangan simpulkan usang tanpa bukti
> menyeluruh). **Cabang `IsBonding` & tangga akseptasi TIDAK dihapus** — statusnya tetap terbuka
> sampai seluruh pemanggil dicek.

---

## Cara mengembalikan jawaban

- **DBA**: kirim hasil query (P-2, **T-1**, **T-4**) sebagai tabel teks; DDL (P-1, P-3) sebagai teks
  `CREATE TABLE`; isi dropdown (P-4) sebagai daftar kode+label; kode prosedur EDM (P-10) sebagai teks
  definisi; **T-2** dan **T-3** sebagai jawaban tertulis singkat.
  Simpan di `D:\migrasi\RNM\DDL\` atau lampiran.
- **IT / administrator Pega**: ekspor (P-5, P-6, P-11, **T-5**, **T-6**) sebagai berkas rule
  **beserta `pyClassName`-nya**; **T-7** sebagai jawaban tertulis singkat; penjelasan spreading (P-8)
  & arti `pyStepsPreCondition` (P-11) sebagai teks.
- **Product**: daftar arti kode (P-7 — **tertutup**, dan **T-9**) sebagai tabel kode → arti.
- **Work owner**: keputusan P-12 (W-1…W-4 — **tertutup**) sebagai jawaban tertulis singkat per butir;
  pengamatan layar Pega mengikuti `10-audit\11-daftar-periksa-ui-pega.md`.

⛔ **Yang tidak pernah kami minta, di paket mana pun:** `USERNAME`, `PASSWORD`, `NAMA`, `LOGIN`,
`OPERATORID`, nomor polis, plat nomor, tanggal lahir, alamat tertanggung, atau dump production.
**Nama kolomnya boleh disebut dalam jawaban; isinya jangan dikirim.**

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan. Semua permintaan terverifikasi dari korpus + DDL.*
