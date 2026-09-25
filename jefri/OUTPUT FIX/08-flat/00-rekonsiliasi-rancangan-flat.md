# Rekonsiliasi rancangan tabel flat — lama (5 contoh) vs baru (114 contoh)

> **24 September 2026.** Keluaran **Tugas 1**: sembilan titik bentrok diperiksa ke korpus, lalu
> direkomendasikan. ⛔ **Tidak ada yang diputuskan di sini.** Setelah dokumen ini: **berhenti dan
> tunggu work owner.**
>
> Dokumen yang direkonsiliasi:
> **lama** `08-flat\01-rancangan-tabel-flat-dari-contoh-nb.md` (24.414 B, 5 contoh) ·
> **baru** `Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx` (18 lembar, 114 contoh, keputusan V-1…V-50)
>
> ⛔ **Dokumen lama belum digantikan.** Ia dipakai menyusun tiket; penggantinya hanya sah setelah
> kesembilan butir diputuskan (`PANDUAN-KERJA` §7).
>
> ⚠️ **BACA BERSAMA `00a-verifikasi-independen-tugas-1.md` (24-09-2026).** Verifikasi independen
> mereproduksi 25 dari 27 angka dokumen ini, tetapi **mengubah dua butir**: butir **7** berbalik
> (V-31 bertahan — lihat baris relasi 58 diagram rumah) dan butir **3** tidak lagi memblokir
> (pertentangannya semu). Ia juga menemukan bahwa **dua contoh adalah duplikat byte-identik**
> (populasi efektif **112**, bukan 114) dan bahwa berkas sumber `Diagram-Skema-Tabel-NusantaraRe.xlsx`
> **sudah dihapus** dari `DDL\`.

---

## 0. Gerbang — nomor yang dibaca, bukan ditebak

`[terverifikasi]`

| | Nilai | Perintah audit |
| --- | --- | --- |
| Judul `K-` di register | **60** | `[regex]::Matches($t,'(?m)^##\s+K-(\d{3})')` |
| Rentang | **K-001 … K-062** | idem |
| **`K-` berikutnya** | **K-063** | idem |
| ADR terakhir | **`0006-mata-uang-unknown-eksplisit.md`** | `[IO.Directory]::GetFiles('OUTPUT\adr','*.md')` |
| **ADR berikutnya** | **0007** | idem |

⚠️ `K-017` dan `K-022` **tidak punya judul `##`** di register meski rentangnya utuh 001–062. Tidak
memblokir pekerjaan ini; dicatat agar diperiksa terpisah.

---

## 1. Apa yang saya ukur sendiri, apa yang saya terima apa adanya

⛔ Pemisahan ini mengikat pembacaan seluruh dokumen.

**Diukur ulang di sesi ini** — 114 berkas `DDL\CONTOH\*.xml`, `DDL\*.txt`, dan isi xlsx:
jumlah berkas contoh dan siklusnya · bentuk `<ID>` akar · seluruh tag ber-pola `Old` · DDL
`JSON_POLIS` dan `HISTORYAKSEPTASIPRODUCTION` · 1.290 baris lembar `Kolom` · 173 baris lembar
`Daftar Relasi` · seluruh nama tabel di `Diagram-Skema-Tabel-NusantaraRe.xlsx`.

**DITERIMA apa adanya dari rancangan baru, TIDAK saya ukur ulang** — dan karena itu **tidak boleh
dibaca sebagai `[terverifikasi]` oleh saya**: V-41 (117 dari 1.422 baris spreading) · V-47 (325
penunjuk lokasi, nol menggantung) · R1 (6 kolom indeks-diri, nol penyimpangan) · R1b (1.094 sama, 2
beda) · V-28 (uji aritmetika 106 contoh, selisih pembulatan 2,6 per sepuluh juta) · R5 (53 tabel
tunggal) · R7 (509 → 126) · cakupan kolom V-48.

---

## 2. ⛔ TABEL KEPUTUSAN — sembilan butir

| # | Titik bentrok | Rekomendasi | Alasan singkat | Keyakinan |
| :-: | --- | --- | --- | --- |
| **1** | Nama tabel `FLAT_*` vs `T_*` | ✅ **BARU menang** | `[terverifikasi]` skema rumah Nusantara Re memakai **43 nama `T_*`, nol `FLAT_*`** — dan `T_WORK_POLIS` + `T_GENERAL_POLIS` **sudah ada di sana** | **terverifikasi** |
| **2** | Kunci alami dipertahankan vs V-47 (14 dibuang, 34 jadi FK) | ✅ **BARU menang, dengan syarat** | Arah V-47 benar; tetapi angkanya **tidak saya ukur ulang**. Syarat: kolom yang buktinya bercampur tetap **tidak** disentuh | **dugaan** |
| **3** | `pxListSubscript` dipertahankan vs V-41 `SEQ_NO` | ⚠️ **TUNDA — rancangan baru bertentangan dengan dirinya sendiri** | V-41 menghapus `pxListSubscript`; **V-40b menyatakan `PX_LIST_SUBSCRIPT` WAJIB disimpan**. Dua pernyataan ini tidak dapat berlaku bersama | **terverifikasi** (pertentangannya) |
| **4** | Pasangan `*Old` dipertahankan vs V-28b tidak disimpan | ⚠️ **BARU menang untuk kaidahnya, TAPI ada risiko migrasi** | `[terverifikasi]` 9 tag berakhiran `Old` terisi **hanya di EDM, nol di NB**. ⛔ Tetapi nilai lama hanya dapat dibaca bila **versi sebelumnya ikut termigrasi** | **terverifikasi** (kaidah) · **dugaan** (risiko) |
| **5** | FK komposit `IDPEGA+PRODKE+NOPOLIS` vs `IDPEGA` saja | ✅ **BARU menang** | `[terverifikasi]` `JSON_POLIS` ber-**`PRIMARY KEY (IDPEGA)`** tunggal → dua versi dengan IDPEGA sama **mustahil**. Ini menutup `[pertanyaan terbuka]` §3.1 dokumen lama | **terverifikasi** |
| **6** | `PRODKE VARCHAR2(5)` vs `PROD_KE NUMBER(5)` | ✅ **BARU menang** | `[terverifikasi]` tipe lama memang `VARCHAR2(5)`. Alasan V-38 sahih. ⚠️ Lihat butir 5b: mungkin **tidak diperlukan sama sekali** | **terverifikasi** |
| **7** | `FLAT_VIEW_SUGGEST` vs V-31 pakai `HISTORYAKSEPTASIPRODUCTION` | ⛔ **DIBUKA ULANG — keduanya melewatkan satu fakta** | `[terverifikasi]` **`T_VIEW_SUGGEST` SUDAH ADA** di skema rumah, dan diagram itu sendiri bertanya *"[terbuka] Apakah FAC memakai T_VIEW_SUGGEST"* | **terverifikasi** |
| **8** | `FLAT_TOTAL_*` sebagai tabel vs V-28 dibuang | ⛔ **BENTROK DENGAN ADR-0001 — perlu keputusan** | V-28 menyimpulkan rekonsiliasi *"perlu **ambang toleransi**"*; ADR-0001 **menolak toleransi secara eksplisit** | **terverifikasi** |
| **9** | `*_CCY` per kolom uang vs tidak diterapkan | ⛔ **LAMA menang — celah terbesar rancangan baru** | `[terverifikasi]` **65 dari 75 tabel** tidak punya kolom mata uang; **61 pasangan (grup, tabel)** yang memuat kolom uang **tidak punya mata uang di mana pun sepanjang rantai induknya** | **terverifikasi** |

---

## 3. Rincian dan bukti

### 3.1 Butir 1 — penamaan · ✅ rancangan baru

`[terverifikasi]` `DDL\CONTOH\Diagram-Skema-Tabel-NusantaraRe.xlsx` (10 lembar) memuat **43** nama
berawalan `T_` dan **0** berawalan `FLAT_`.

⛔ **`T_WORK_POLIS` dan `T_GENERAL_POLIS` bukan nama baru** — keduanya sudah dipakai lembar
**Treaty In** (NB & EDM, Prop & NonProp) dan **PremiumList**. Baris relasi 52 diagram itu berbunyi:
`T_WORK_POLIS → T_GENERAL_POLIS · SHARED PK (ID = ID) · 1:1 · T_WORK_POLIS dipakai bersama PremiumList`.

Jadi rancangan baru **bukan mencipta akar**, melainkan **menyambung ke akar yang sudah ada**.
`FLAT_*` akan memecah skema rumah jadi dua gaya.

📝 **Koreksi rujukan, bukan koreksi keputusan.** Catatan V-48 menyebut *"meniru **T_WORK_TREATY_IN**
di lembar Treaty In NonProp"*. `[terverifikasi]` **nama itu tidak ada** di diagram — yang ada
`T_WORK_POLIS`. Keputusannya tetap benar; **rujukannya yang salah nama**.

⚠️ **Pertanyaan yang lahir dari sini dan belum dijawab siapa pun:** apakah Fac In memakai **tabel
fisik yang sama** `T_WORK_POLIS`/`T_GENERAL_POLIS` dengan Treaty In, atau tabel sendiri yang kebetulan
senama? Bila sama, kolomnya harus direkonsiliasi dengan lembar Treaty In — dan itu **belum dilakukan**:
rancangan baru menurunkan kolomnya **hanya** dari 114 contoh Fac In.

```powershell
# 43 token T_*, 0 FLAT_*
Add-Type -AssemblyName System.IO.Compression.FileSystem
$zip=[IO.Compression.ZipFile]::OpenRead('D:\migrasi\RNM\DDL\CONTOH\Diagram-Skema-Tabel-NusantaraRe.xlsx')
$tok=@{}
foreach($e in $zip.Entries){ if($e.FullName -notlike 'xl/worksheets/*'){continue}
  $s=$e.Open();$r=New-Object IO.StreamReader($s,[Text.Encoding]::UTF8);$t=$r.ReadToEnd();$r.Close();$s.Close()
  foreach($m in [regex]::Matches($t,'\bT_[A-Z0-9_]{2,}\b')){ $tok[$m.Value]=1 } }
$zip.Dispose(); $tok.Count
```

### 3.2 Butir 5 — kunci versi · ✅ rancangan baru · **menutup satu `[pertanyaan terbuka]`**

Dokumen lama §3.1 menandai ini terbuka: *"perlu contoh dua versi dengan IDPEGA sama"*.

`[terverifikasi]` `DDL\JSON_POLIS.txt` menjawabnya **tanpa perlu contoh itu**:

```
CONSTRAINT "PK_JSON_POLIS" PRIMARY KEY ("IDPEGA")
CONSTRAINT "JSON_POLIS_U01" UNIQUE ("NOPOLIS", "IDPEGA")
"PRODKE" VARCHAR2(5)
```

⛔ `IDPEGA` adalah **primary key tunggal**. Dua baris dengan `IDPEGA` sama **tidak mungkin ada** di
tabel itu. Maka versi kedua sebuah polis **wajib** ber-`IDPEGA` berbeda — pilihan "PK efektif
(IDPEGA, PRODKE)" di dokumen lama **tertutup oleh DDL**, bukan oleh selera.

⚠️ **Yang tetap TIDAK terbukti:** tidak ada satu pun pasangan dua versi polis yang sama di dalam 114
contoh. `[terverifikasi]` 101 NB + 13 EDM + **0 RNW**; ke-13 EDM memakai nomor work object EDM-nnnnn
sendiri. Jadi kesimpulan di atas bersandar pada **DDL**, bukan pada contoh.

### 3.3 Butir 5b — ⛔ **temuan yang tidak ada di daftar sembilan**: skema rumah sudah punya jawabannya

`[terverifikasi]` Baris relasi **53** diagram rumah:

```
Treaty In | T_WORK_POLIS | T_GENERAL_POLIS | OLD_POLIS_ID | 1:1 | di Go | UNIK
          | nullable · penunjuk ke generasi SEBELUMNYA · pengganti OldData · UNIK melarang percabangan
```

⛔ **Treaty In menaut versi dengan penunjuk langsung `OLD_POLIS_ID`**, bukan dengan mencari
`PROD_KE` tertinggi. Rancangan Fac In baru (V-38, V-40) memilih mekanisme **yang berbeda** untuk
masalah **yang sama**, di basis data **yang sama**.

📌 Bila `OLD_POLIS_ID` diadopsi, **butir 6 sebagian besar gugur** — tipe `PROD_KE` tidak lagi
menentukan pencarian versi. Ini **usulan**, bukan keputusan.

### 3.4 Butir 4 — before-image · kaidahnya terbukti, risikonya nyata

`[terverifikasi]` Sapuan seluruh 114 contoh, dihitung per **berkas yang tag-nya terisi**:

| Pola | Tag | NB | RNW | EDM |
| --- | --- | ---: | ---: | ---: |
| **berakhiran `Old`** | `TSIOld` · `PremiumOld` | 0 | 0 | 10 · 10 |
| | `RateOld` | 0 | 0 | 9 |
| | `TotalGrossPremiOld` · `TotalPremiumNusantaraReOld` · `TSIObjectItemOld` | 0 | 0 | 6 |
| | `PremiNusantaraReOld` | 0 | 0 | 2 |
| | `PremiRpOld` · `PremiumGrossDiscountFleetOld` | 0 | 0 | 1 |
| **berawalan `Old`** | `OldID` | **101** | 0 | 13 |
| | `OldTSI` | **15** | 0 | 2 |
| | `OldPolicyNo` | 0 | 0 | 13 |
| **pola `EDMOld`** | `EDMOldCommision` · `EDMOldPayment` | **4** | 0 | 13 |
| | `EDMOldPremi` | **2** | 0 | 13 |
| | `EDMOldBrokerage` · `EDMOldDeduction2` · `EDMOldPPh` · `EDMOldPPN` | 0 | 0 | 12 |

✅ **Kaidah V-28b tereproduksi**: seluruh 9 tag **berakhiran** `Old` nol di NB. `OldTSI` **15 NB**
tereproduksi persis.
📝 **Satu angka berbeda:** V-28b menyebut `OldID` di **93** berkas NB; saya mengukur **101**.
Kemungkinan definisinya berbeda (saya menghitung tag bernama `OldID` **di mana pun** dengan nilai tak
kosong). **Perlu satu klarifikasi**, bukan perubahan kaidah.
✅ **V-45 tereproduksi:** `EDMOld*` memang terisi di berkas NB — **4 berkas** pada
`EDMOldCommision`/`EDMOldPayment`, **2** pada `EDMOldPremi`.

⛔ **Risiko yang belum tercatat di kedua dokumen.** Membuang field berakhiran `Old` berarti nilai
sebelum endorsement **hanya** dapat dibaca dari baris versi sebelumnya. Untuk **data lama yang
dimigrasikan**, itu berlaku hanya bila versi sebelumnya **ikut dimuat**. Bila migrasi memuat versi
terakhir saja, nilai lama **hilang permanen** — padahal ia ada di kolom yang dibuang. Ini pertanyaan
lingkup migrasi, dan **bukan keputusan saya**.

### 3.5 Butir 7 — ⛔ **dibuka ulang**

`[terverifikasi]` `HISTORYAKSEPTASIPRODUCTION` memang memuat kelima kolom padanan yang disebut V-31 —
`NOURUT`, `PIC`, `KETERANGAN`, `TGL_INP`, `APPROVAL` (dari 15 kolom seluruhnya).

⛔ **Tetapi `T_VIEW_SUGGEST` SUDAH ADA sebagai tabel di skema rumah**, dengan pola yang justru sudah
dirancang untuk dipakai lintas-lini:

```
T_VIEW_SUGGEST punya DUA induk - PREMIUM_LIST_ID dan CLAIM_ID, keduanya nullable,
dengan CHECK ( (PREMIUM_LIST_ID IS NULL) <> (CLAIM_ID IS NULL) ). Tepat satu terisi.
```

Dan diagram rumah itu **mengajukan pertanyaannya sendiri**: `[terbuka] Apakah FAC memakai
T_VIEW_SUGGEST.`

📌 Bila Fac In ikut memakainya dengan induk ketiga (`POLIS_ID`, nullable, CHECK diperluas), maka
**peringatan V-31 ikut selesai**: `DateTransfer` dan `IsCedingConfirm` — yang V-31 catat **tidak punya
padanan** di tabel lama — mendapat tempat tanpa memaksa `IsCedingConfirm` masuk kolom `POSISI`.

⚠️ Memasukkan `IsCedingConfirm` ke `POSISI` (usul V-32) **bertabrakan dengan V-48**, yang memakai
`POSISI` untuk posisi tangga akseptasi. Satu kolom, dua arti.

### 3.6 Butir 8 — ⛔ bentrok dengan ADR-0001

V-28 (lembar `BACA-INI`) menutup dengan:

> *"hitung ulang di Go tidak menghasilkan angka identik dengan yang Pega simpan — selisih pembulatan
> sampai 2,6 per sepuluh juta pada penjumlahan 195 baris — jadi rekonsiliasi dengan baris
> FACINPRODUCTION lama perlu **ambang toleransi**, bukan perbandingan sama persis."*

`adr/0001-rekonsiliasi-eksak-bertahap.md` menolaknya secara eksplisit:

> *"**Toleransi seragam (mis. ±0,01) — ditolak.** … Selisih sekecil apa pun berarti ada salah-port,
> dan toleransi justru menyembunyikannya."*

⛔ **Keduanya tidak dapat berlaku bersama.** Membuang tabel `Total*` memaksa **hitung ulang**;
hitung ulang memaksa **toleransi**; toleransi dilarang ADR-0001.

Tiga jalan keluar — **pilihan work owner, bukan saya**:

| | Jalan | Akibat |
| :-: | --- | --- |
| **A** | Tabel `Total*` tetap dibuang; **ADR-0001 diamandemen** untuk kelas nilai turunan ini | Keputusan V-28 utuh; ukuran keberhasilan migrasi berubah |
| **B** | Nilai `Total*` **disimpan apa adanya** sebagai angka terekam (bukan dihitung ulang), tabelnya boleh menyusut jadi kolom | ADR-0001 utuh; V-28 disempitkan dari "dibuang" jadi "tidak dihitung ulang" |
| **C** | Dibuang untuk data **baru**, disimpan untuk data **lama yang dimigrasikan** | Rekonsiliasi paralel run tetap eksak; skema jadi dua perlakuan |

⚠️ Saya **tidak mengukur ulang** uji aritmetika V-28 maupun angka selisih 2,6 per sepuluh juta.
Yang saya verifikasi adalah **pertentangan tekstualnya** dengan ADR-0001.

### 3.7 Butir 9 — ⛔ celah terbesar · **dokumen lama yang benar**

Prinsip dokumen lama: *"Kolom `*_CCY VARCHAR2(10)` menyertai tiap kolom uang."*
Dasar mengikatnya: **K-010/K-012** dan **ADR-0004** (`Money{Amount, Currency}`).

`[terverifikasi]` Atas 1.290 baris lembar `Kolom` dan 173 baris lembar `Daftar Relasi`:

| Ukuran | Nilai |
| --- | ---: |
| Tabel | **75** |
| Tabel **punya** kolom mata uang | **10** |
| Tabel **tanpa** kolom mata uang | **65** |
| Kolom uang terdeteksi (`NUMBER` + nama uang) | **156** |
| — berada di tabel tanpa kolom mata uang | **131** |
| Pasangan (grup, tabel) ber-kolom uang tanpa mata uang sendiri | **77** |
| — mata uang **dapat** ditelusuri ke leluhur | **16** |
| — ⛔ **tidak dapat sama sekali** | **61** |

⛔ Yang paling menentukan:

- **`T_SPREADINGLIST`** — butir terhalus, pemberi makan `FACINPRODUCTION`, 5 kolom uang. Mata uang
  dapat ditelusuri **hanya** untuk **FIRE** dan **Aneka**. Untuk **Life, MarineCargo, MBUCar, PA:
  tidak ada di seluruh rantai** sampai `T_WORK_POLIS`.
- **`T_COVERAGELIST`** — simpul konvergen, **17 kolom uang**. Sama: hanya FIRE dan Aneka yang
  tertelusur.
- **`T_FR_CURRENCYLIST`** — **18 kolom uang**, dan namanya sendiri *CurrencyList*, tetapi **tidak
  punya kolom kode mata uang**.
- Hampir seluruh cabang **`T_FR_*`** (Fac Retro) tidak tertelusur di grup mana pun.

⚠️ Konsekuensinya persis yang dilarang **ADR-0006**: nilai uang tanpa mata uang yang **tidak
eksplisit `Unknown`** melainkan **tidak ada sama sekali** — kegagalannya senyap.

```powershell
# 75 tabel / 10 punya mata uang / 156 kolom uang / 61 pasangan tak tertelusur
# (lembar Kolom dan Daftar Relasi diekstrak dari xlsx, lalu rantai induk ditelusuri per grup)
$rxCcy =[regex]'(?i)CURRENCY|_CCY$|^CCY|CURR_ID'
$rxUang=[regex]'(?i)PREMI|TSI|AMOUNT|PAYMENT|DISCOUNT|COMMISION|COMMISSION|BROKERAGE|PPN|PPH|DEDUCTION|STAMP|CLAIM|SPREADED|RETENTION|SUBLIMIT'
$rxBukan=[regex]'(?i)PERCENT|PCT|_RATE|RATE_|^RATE$|_NO$|^NO_|IDX|INDEX|_ID$|^IS_|FLAG|TYPE|NAME|NOTE|DATE|_REF_'
# kolom uang = JENIS != sistem, TIPE ~ NUMBER, nama cocok $rxUang dan tidak cocok $rxBukan
```

⚠️ **Batas ukuran ini:** "kolom uang" dikenali dari **nama dan tipe**, bukan dari arti bisnisnya.
Ambangnya bisa meleset di kedua arah. Angka **61** adalah **batas bawah kekhawatiran**, bukan angka
final — tetapi arahnya tidak berubah oleh penyetelan ambang.

### 3.8 Butir 3 — ⛔ rancangan baru bertentangan dengan dirinya sendiri

`[terverifikasi]` Dua pernyataan di lembar `BACA-INI` yang sama:

| Nomor | Bunyi |
| --- | --- |
| **V-41** | *"pxListSubscript tidak membawa keterangan apa pun di luar posisi, dan sebagai kolom data ia **dihapus**"* |
| **V-40b** | *"Karena itu **PX_LIST_SUBSCRIPT WAJIB disimpan**, dan ITEM_TYPE dipakai sebagai pemeriksa silang"* |

Keduanya tidak dapat berlaku bersama. Pembacaan yang paling mungkin mendamaikan: V-40b sebenarnya
menuntut **posisi baris**, yang setelah V-41 bernama `SEQ_NO` — sehingga yang dihapus hanyalah
**namanya**, bukan **isinya**. ⚠️ Itu **`[dugaan]` saya**, bukan yang tertulis. **Perlu satu kalimat
penegas dari work owner**, karena V-40 menjadikan posisi baris sebagai **kunci penjodohan selisih
endorsement** — kalau kolomnya benar-benar hilang, penjodohan itu ikut hilang.

### 3.9 Butir 2 dan 6 — ringkas

**Butir 2.** Arah V-47 sejalan dengan skema rumah, yang memakai FK bernama (`CLAIM_ID`,
`OBJECT_ITEM_ID`, `PREMIUM_LIST_ID`) dan bukan nomor posisi. ⚠️ Angka "14 dibuang / 34 jadi FK /
325 penunjuk nol menggantung" **tidak saya ukur ulang** → **`[dugaan]`**. Syarat yang saya
rekomendasikan: **10 kolom yang V-50 tahan karena buktinya bercampur tetap tidak disentuh**.

**Butir 6.** `[terverifikasi]` `JSON_POLIS.PRODKE` memang `VARCHAR2(5)`. Alasan V-38 sahih secara
umum. ⚠️ Tetapi apakah `PROD_KE` masih dipakai untuk **mencari versi** bergantung pada butir 5b —
bila `OLD_POLIS_ID` diadopsi, tipe `PROD_KE` tinggal soal kerapian.
⚠️ Apakah seluruh nilai `PRODKE` lama benar-benar angka **tidak dapat diperiksa dari korpus** —
isinya ada di basis data, bukan di ekspor. **Butir DBA.**

---

## 4. Tiga temuan di luar daftar sembilan

### 4.1 ⛔ `<ID>` di akar memuat DUA hal berbeda

`[terverifikasi]` Atas 114 contoh, elemen `<ID>` anak-langsung akar:

| Bentuk | Berkas | Nilai unik |
| --- | ---: | ---: |
| **angka murni** (mis. 6 digit) | **38** | 38 |
| **kelas Pega + spasi + nomor work object** | **50** | **48** |
| tidak ada elemennya | **26** | — |

✅ Angka **38** **tereproduksi persis** dengan V-49a (*"SOURCE_ID … terisi di 38/114"*).

⛔ **Tetapi 50 berkas lain memuat bentuk yang sama sekali berbeda** — kunci objek kerja. Rancangan
baru memetakan `ID` akar → **`SOURCE_ID`** secara seragam. Untuk kelima puluh berkas itu, `SOURCE_ID`
akan terisi **kunci kerja**, bukan rujukan master. **Dua ruang nilai dalam satu kolom.**

📌 **Dua nilai bentuk-kelas berulang** di 48 nilai unik dari 50 berkas — pasangannya sama persis
dengan anomali penyalinan data yang sudah dicatat V-45: `NB-181622` ↔ `NB-184183` dan
`NB-172576` ↔ `NB-176005`.

### 4.2 ⛔ Tidak ada satu pun contoh RNW

`[terverifikasi]` 114 contoh = **101 NB + 13 EDM + 0 RNW**.

Dokumen lama dibangun dari 5 contoh yang **memuat 1 RNW** (`P-5 RNW-10579 (FIRE).txt`), dan justru
dari situ ia menyimpulkan *"RNW dan EDM memakai struktur yang SAMA dengan NB"*.

⛔ Rancangan baru **tidak dapat menguji ulang kesimpulan itu** — populasinya nol. Untuk siklus RNW,
rancangan baru **bukan lebih kuat** daripada yang lama; ia **tidak punya bukti sama sekali**.
Setiap pernyataan V-28b/V-45 yang berbunyi *"nol di NB **dan RNW**"* **hampa** untuk bagian RNW.

⚠️ Ini **tidak membatalkan** rancangan baru. Ia menandai satu lubang yang dapat ditutup murah:
**minta beberapa contoh RNW**.

### 4.3 Berkas contoh memuat data pribadi nyata

`[terverifikasi]` Berkas di `DDL\CONTOH\` memuat nilai nama orang pada sekurangnya `PIC`,
`CommentCeding`, dan metadata `pxCreateOpName`. V-48 sudah menandai `PIC` sebagai data pribadi.

⛔ **Konsekuensi kerja:** setiap skrip yang menyapu folder ini **tidak boleh** mencetak nilai —
hanya nama tag dan hitungan. Seluruh pengukuran di dokumen ini mematuhi itu; **nol nilai data pribadi
disalin ke sini** (K-025).

---

## 5. Yang TIDAK dapat saya verifikasi

| # | Klaim | Sebab |
| :-: | --- | --- |
| 1 | Angka V-41 · V-47 · R1 · R1b · R5 · R7 | Tidak diukur ulang di sesi ini — butuh sapuan tersendiri per kolom |
| 2 | Uji aritmetika V-28 dan selisih 2,6 per sepuluh juta | Tidak diukur ulang; yang saya verifikasi hanya pertentangannya dengan ADR-0001 |
| 3 | Apakah Fac In memakai **tabel fisik yang sama** dengan Treaty In | Diagram rumah tidak menyatakan; kolom Treaty In belum direkonsiliasi dengan 114 contoh Fac In |
| 4 | Apakah seluruh nilai `PRODKE` lama berupa angka | Isinya di basis data, bukan di ekspor — **butir DBA** |
| 5 | Apakah migrasi akan memuat **seluruh** versi polis atau versi terakhir saja | Menentukan apakah risiko butir 4 nyata; keputusan lingkup, bukan fakta korpus |
| 6 | Apakah struktur RNW sama dengan NB | **Nol contoh RNW** di populasi baru |
| 7 | Maksud sebenarnya V-40b terhadap `PX_LIST_SUBSCRIPT` | Dua pernyataan saling bertentangan di dokumen yang sama |
| 8 | Apakah 26 berkas tanpa `<ID>` akar memang tidak punya, atau strukturnya berbeda | Diukur sebagai "tidak ada anak-langsung bernama `ID`"; bentuk lain belum ditelusuri |

---

## 6. Yang saya minta diputuskan

⛔ **Sebelum `08-flat\02-…` boleh ditulis**, kesembilan butir §2 perlu jawaban. Yang **paling
memblokir**, berurutan:

1. **Butir 9** — apakah `*_CCY` (atau `CURRENCY_CODE` per tabel) diterapkan menyeluruh. Menyentuh
   K-010, K-012, ADR-0004, ADR-0006.
2. **Butir 8** — jalan **A**, **B**, atau **C**. Menyentuh ADR-0001, yaitu **ukuran keberhasilan
   migrasi**.
3. **Butir 7** — `T_VIEW_SUGGEST` rumah, `HISTORYAKSEPTASIPRODUCTION`, atau keduanya.
4. **Butir 5b** — `OLD_POLIS_ID` seperti Treaty In, atau pencarian `PROD_KE` seperti V-38.
5. **Butir 3** — apakah `PX_LIST_SUBSCRIPT` disimpan atau tidak.
6. **Butir 1** lanjutan — tabel fisik bersama dengan Treaty In, atau terpisah.

Butir **2**, **4**, **6** dapat mengikuti rekomendasi §2 tanpa pembahasan tambahan bila work owner
setuju.

---

## 7. Perintah audit pokok

```powershell
# 114 contoh; 101 NB / 13 EDM / 0 RNW
$f=@([IO.Directory]::GetFiles('D:\migrasi\RNM\DDL\CONTOH','*.xml')); $f.Count

# tag pola Old per siklus (hanya HITUNGAN, tidak pernah nilainya)
#   berakhiran Old : 9 tag, seluruhnya 0 di NB
#   OldID 101 NB / OldTSI 15 NB / EDMOldCommision 4 NB

# bentuk <ID> akar : 38 angka murni / 50 kelas+nomor / 26 tidak ada
foreach($x in $f){ $d=New-Object Xml.XmlDocument; $d.Load($x)
  foreach($c in $d.DocumentElement.ChildNodes){ if($c.Name -eq 'ID'){ $c.InnerText.Trim() } } }

# kunci versi
Select-String -Path 'D:\migrasi\RNM\DDL\JSON_POLIS.txt' -Pattern 'PRIMARY KEY|UNIQUE|PRODKE'
```

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan. Tidak ada keputusan diambil di dokumen ini.*
