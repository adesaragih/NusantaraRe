# Verifikasi independen atas `00-rekonsiliasi-rancangan-flat.md`

> **24 September 2026.** Seluruh angka di bawah **diukur ulang sendiri** dari korpus, bukan disalin
> dari dokumen yang diverifikasi. ⛔ **Tidak ada keputusan diambil di sini.**
>
> Yang diverifikasi: `08-flat\00-rekonsiliasi-rancangan-flat.md` (24-09-2026 15.37).
> Bahan: `DDL\CONTOH\*.xml` (114) · `DDL\*.txt` · `Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx` ·
> `Diagram-Skema-Tabel-NusantaraRe.xlsx` (lihat §1 — berkasnya sudah tidak ada di tempatnya).

---

## 0. Ringkasan

**Tidak satu pun temuan pokok laporan itu runtuh.** Dari 27 angka yang dapat diukur ulang, **25
tereproduksi persis**, 2 meleset tipis tanpa mengubah arah. Butir 4, 5, 6, 9 dan §4.1 berdiri utuh.

Yang berubah:

| | Butir | Status sesudah verifikasi |
| :-: | --- | --- |
| **A** | **Butir 7** | ⛔ **Rekomendasi berbalik.** Pembukaan ulang bersandar pada salah baca. **V-31 bertahan.** |
| **B** | **Butir 3** | ✅ **Tidak perlu ditunda lagi.** Pertentangannya semu; diselesaikan oleh pengukuran baru. |
| **C** | Butir 8 | Pertentangan **nyata** — tetapi ada jalan **D** yang tidak mengubah ADR-0001. |
| **D** | Provenans | ⛔ Sumber butir 1, 5b, dan 7 **dihapus hari ini**. Dapat dipulihkan. |
| **E** | Korpus | ⛔ **Dua contoh adalah duplikat byte-identik.** Populasi efektif **112**, bukan 114. |

---

## 1. ⛔ Sumber butir 1, 5b dan 7 sudah tidak ada di tempatnya

`[terverifikasi]` `Diagram-Skema-Tabel-NusantaraRe.xlsx` **tidak ada** di mana pun di bawah
`D:\migrasi\RNM\`. Ia **dihapus ke Recycle Bin hari ini**:

| Salinan | Dihapus | Bukti |
| --- | --- | --- |
| `DDL\Diagram-Skema-Tabel-NusantaraRe.xlsx` | 24-09-2026 **15.31.49** | `$ISIG7TE.xlsx` |
| `DDL\CONTOH\Diagram-Skema-Tabel-NusantaraRe.xlsx` | 24-09-2026 **15.31.57** | `$IUUGGHS.xlsx` |

Kedua stempel waktu itu **persis** sama dengan `LastWriteTime` folder `DDL\` dan `DDL\CONTOH\`, dan
keduanya **enam menit sebelum** laporan disimpan (15.37).

`[terverifikasi]` Isinya **tidak hilang dan tidak berubah**. Tiga salinan ber-SHA256 identik
`1DCDCDA0A8F7BD49…` — dua di Recycle Bin, satu di `C:\Users\Administrator\Downloads\RNM\DDL\`.
Seluruh verifikasi butir 1, 5b dan 7 di bawah dikerjakan atas salinan itu.

⚠️ **Tindakan yang diperlukan:** kembalikan berkas itu ke `D:\migrasi\RNM\DDL\`. Selama ia tidak di
sana, tiga butir bersandar pada sumber yang **tidak dapat diaudit ulang** oleh siapa pun — melanggar
`CLAUDE.md` §3.1.

Tiga berkas lain di `Claude outputs\` ikut dihapus 15.32 (`Rancangan-Tabel-Flat-FacIn`,
`Tabel-Flat-per-Grup-Bisnis-1`, `-2`). `[terverifikasi]` `-2` **byte-identik** dengan
`Tabel-Flat-per-Grup-Bisnis.xlsx` yang masih hidup (`1D0341E2C30FC06B…`), jadi laporan membaca isi
yang sama dengan yang saya baca.

```powershell
$sh=New-Object -ComObject Shell.Application; $rb=$sh.Namespace(10)
$rb.Items() | ForEach-Object { "{0} | {1} | {2}" -f $rb.GetDetailsOf($_,0),$rb.GetDetailsOf($_,1),$rb.GetDetailsOf($_,2) }
```

---

## 2. Yang tereproduksi persis

`[terverifikasi]` Diukur ulang sendiri, seluruhnya cocok:

| Klaim laporan | Nilai saya | |
| --- | --- | :-: |
| Judul `K-` = 60, rentang K-001…K-062, berikutnya **K-063** | 60 · K-001…K-062 | ✅ |
| ADR terakhir `0006`, berikutnya **0007** | idem | ✅ |
| 114 contoh = **101 NB + 13 EDM + 0 RNW** | 101 / 13 / 0 | ✅ |
| **Butir 1** — 43 nama `T_*`, **0** `FLAT_*` | 43 · 0 | ✅ |
| **Butir 1** — `T_WORK_TREATY_IN` tidak ada di diagram (rujukan V-48 salah nama) | 0 kemunculan | ✅ |
| **Butir 5** — `JSON_POLIS` ber-`PRIMARY KEY (IDPEGA)` tunggal | terbaca | ✅ |
| **Butir 6** — `PRODKE VARCHAR2(5)` | terbaca | ✅ |
| **Butir 7** — `HISTORYAKSEPTASIPRODUCTION` memuat `NOURUT PIC KETERANGAN TGL_INP APPROVAL` dari 15 kolom | 15 kolom, kelimanya ada | ✅ |
| **Butir 4** — 9 tag berakhiran `Old`, **nol di NB** | 9 tag, 0 NB | ✅ |
| **Butir 4** — `OldID` 101 NB · `OldTSI` 15 NB · `EDMOldCommision`/`Payment` 4 NB · `EDMOldPremi` 2 NB | sama persis | ✅ |
| **§4.1** — `<ID>` akar: **38** angka / **50** kelas+nomor / **26** tidak ada | 38 / 50 / 26 | ✅ |
| **Butir 9** — 75 tabel · 10 ber-mata-uang · 65 tanpa · 156 kolom uang · 131 di tabel tanpa · 16 tertelusur | sama persis | ✅ |
| **Butir 3** — kutipan V-41 dan V-40b | verbatim | ✅ |
| **Butir 8** — kutipan V-28 dan ADR-0001 | verbatim | ✅ |

Dua yang meleset tipis, **tanpa mengubah arah**:

- **Butir 9** — pasangan (grup, tabel) **76**, bukan 77; tak tertelusur **60**, bukan 61. Beda satu.
- **§1 laporan** — "173 baris lembar `Daftar Relasi`" keliru: lembar itu punya **220** baris relasi
  bernomor. Bukan angka penyangga kesimpulan; rekonstruksi saya memakai seluruh 220 dan tetap
  menghasilkan 16 tertelusur.

⚠️ **K-017 dan K-022 tidak perlu diperiksa terpisah.** Registernya sudah menjelaskan sendiri di
baris 1047–1050: K-022 **dicadangkan** untuk keputusan riwayat akseptasi (entrinya belum ditulis),
K-017 **kosong permanen**. Peringatan di §0 laporan boleh dicoret.

---

## 3. ⛔ Butir 7 — rekomendasi berbalik, V-31 bertahan

Laporan membuka ulang butir 7 atas dua dasar. `[terverifikasi]` **Keduanya salah baca.**

### 3.1 Treaty In — kerabat Fac In yang berbentuk polis — memakai tabel lama, bukan `T_VIEW_SUGGEST`

Baris relasi **58** diagram rumah, yang **tidak disebut** laporan:

```
58 | Treaty In | T_GENERAL_POLIS | POOLDATA.HISTORYAKSEPTASIPRODUCTION | IDPEGA | 1:N | di Go | index
   | dari SuggestList · SaveViewSuggest -> InsertViewSuggest_SQL
```

Dan keempat lembar Treaty In (NB/EDM × Prop/NonProp) menuliskan hal yang sama:
`POOLDATA.HISTORYAKSEPTASIPRODUCTION 1:N · 15 kolom · SUDAH datar`.

⛔ Jadi untuk **objek berbentuk polis**, skema rumah sudah memutuskan persis seperti V-31:
`SuggestList` → `HISTORYAKSEPTASIPRODUCTION`. V-31 **bukan menyimpang dari skema rumah — ia sejalan
dengannya.**

### 3.2 `T_VIEW_SUGGEST` bukan tabel polis, dan pertanyaan `[terbuka]` itu soal KLAIM

`[terverifikasi]` Kedua induk `T_VIEW_SUGGEST` adalah konteks klaim/premium, bukan polis:

```
15 | PremiumList | T_PREMIUM_LIST  | T_VIEW_SUGGEST | PREMIUM_LIST_ID | 1:N
28 | Claim Prop  | T_GENERAL_CLAIM | T_VIEW_SUGGEST | CLAIM_ID        | 1:N
```

Dan kalimat `[terbuka] Apakah FAC memakai T_VIEW_SUGGEST` berada di lembar **`Claim Fac In`**
(baris 113) — lembar **klaim**, bukan underwriting. Lembar yang sama, baris 3, menggolongkan
`T_VIEW_SUGGEST` sebagai *"Tabel yang HANYA dipakai Prop"*. Pertanyaan itu menanyakan lini **klaim
Fac**, bukan Fac In.

### 3.3 Yang benar-benar tersisa terbuka

Bukan "T_VIEW_SUGGEST atau HISTORYAKSEPTASIPRODUCTION", melainkan **satu kolom saja**:

- `DateTransfer` — **sudah selesai**, bukan masalah lagi: **V-36** memutuskan hanya `DateSuggest`
  yang mengisi `TGL_INP`, `DateTransfer` tidak disimpan. V-32 yang memetakannya ke `TGL_INP`
  **sudah kedaluwarsa** oleh V-36.
- `IsCedingConfirm` → `POSISI` (V-32) — **masih bentrok**, dan lebih parah dari yang dicatat
  laporan. V-48 menjadikan `T_WORK_POLIS.POSISI` **cerminan baris terakhir**
  `HISTORYAKSEPTASIPRODUCTION.POSISI`. Bila `IsCedingConfirm` ditulis ke `POSISI`, maka
  `T_WORK_POLIS.POSISI` akan mencerminkan konfirmasi ceding, **bukan posisi tangga akseptasi**.
  V-48c sendiri sudah mencatat daftar nilai `POSISI` yang sah **belum pernah terlihat** — hanya DDL,
  nol baris data.

📌 **Usul:** butir 7 ditutup mengikuti V-31; yang naik ke daftar keputusan adalah **`IsCedingConfirm`
saja** — kolom sendiri, atau `POSISI`, tetapi tidak dua arti dalam satu kolom.

---

## 4. ✅ Butir 3 — pertentangannya semu; dapat diputuskan sekarang

Laporan menunda karena V-41 (hapus `pxListSubscript`) dan V-40b (`PX_LIST_SUBSCRIPT` **WAJIB**
disimpan) tampak tidak dapat berlaku bersama. `[terverifikasi]` Pengukuran baru atas **seluruh 114
contoh** menyelesaikannya:

| | Nilai |
| --- | ---: |
| Baris yang **punya** `pxListSubscript` | **18.140** |
| — nilainya **sama dengan posisi baris** | **18.140** |
| — nilainya **berbeda** dari posisi baris | **0** |
| — **ada tapi kosong** | **0** |

Dan baris yang **tidak punya tag itu sama sekali**:

| Wadah | Baris | Punya | **Tidak punya** |
| --- | ---: | ---: | ---: |
| `SpreadingList` | 4.509 | 4.375 | **134** |
| `CoverageList` | 4.926 | 4.362 | **564** |
| `DeductibleList` | 6.535 | 6.189 | **346** |
| `PropertyItemList` | 929 | 823 | **106** |
| `LocationList` | 533 | 459 | **74** |
| `AnekaList` | 188 | 155 | **33** |
| `ViewSuggest` | 922 | 821 | **101** |

Dua akibat yang menutup butir 3:

1. **Menyimpan `PX_LIST_SUBSCRIPT` apa adanya justru MERUSAK tujuan V-40b.** V-40b memerlukannya
   sebagai kunci penjodohan selisih endorsement — tetapi di **134 baris `SpreadingList`, 106
   `PropertyItemList`, 74 `LocationList`, 33 `AnekaList`** tag itu **tidak ada**, jadi kuncinya akan
   `NULL` tepat di tempat ia dibutuhkan.
2. **`SEQ_NO` membawa keterangan yang identik dan terdefinisi di setiap baris.** Di mana pun
   `pxListSubscript` ada, nilainya **selalu** posisi baris — 18.140 dari 18.140, nol penyimpangan.

⛔ Jadi V-41 dan V-40b **tidak** saling meniadakan: V-40b menuntut **posisi baris**, V-41 membuang
**nama kolom Pega-nya**. Dugaan laporan benar, dan kini `[terverifikasi]`, bukan `[dugaan]`.

📌 **Usul:** adopsi `SEQ_NO` (V-41), lalu **perbaiki bunyi V-40 dan V-40b** supaya menyebut `SEQ_NO`
sebagai kunci penjodohan di dalam objek. Tanpa perbaikan bunyi itu, dua keputusan yang sah akan
terus terbaca bertentangan.

⚠️ **Angka V-41 tidak tereproduksi:** V-41 menyebut *"kosong di 117 dari 1.422 baris SpreadingList"*
dan *"1.823 baris diperiksa"*. Saya mengukur **134 dari 4.509** dan **18.140 baris**. Populasinya
berbeda (V-41 tampaknya diukur pada korpus 105/106 contoh). **Isinya bertahan, angkanya perlu
ditulis ulang.**

---

## 5. ⛔ Dua contoh adalah duplikat byte-identik — populasi efektif 112

`[terverifikasi]` Dua pasang berkas di `DDL\CONTOH\` ber-**SHA256 sama persis**:

| Pasangan | Ukuran | SHA256 (20 pertama) |
| --- | ---: | --- |
| `NB-172576 (AS. KREDIT)` ≡ `NB-176005 (AS. KREDIT)` | 14.050 B | `595D9FDD2F2CD86832BD` |
| `NB-181622 (MARINE CARGO)` ≡ `NB-184183 (MARINE CARGO)` | 35.659 B | `20050C0D3DD6BB30D5D2` |

Bukan sebagian medan yang tersalin — **seluruh isinya sama**: 242 dari 242 dan 561 dari 561 daun
terisi identik, nol berbeda. Tidak ada duplikat lain di antara 114 berkas.

Akibatnya:

1. **Populasi efektif 112, bukan 114.** Setiap statistik berbentuk "N dari 114" atau "N dari 106"
   **menghitung dua kasus dua kali** — termasuk cakupan kolom V-48 (PIC 114, IDNewBisnis 97, …),
   V-49a (38 dari 114), dan hitungan butir 4 saya sendiri (`OldID` 101 NB → **99 kasus berbeda**).
2. **Tafsir V-45 tidak didukung bukti.** V-45 menduga *"penyalinan data dari EDM lain"* pada 4 berkas
   NB yang mengisi `EDMOld*`. Karena `NB-181622` ≡ `NB-184183`, itu **3 kasus berbeda, bukan 4** —
   dan mekanismenya bukan penyalinan medan di Pega melainkan **berkas yang sama diekspor dua kali
   dengan nama berbeda**. Sama untuk "penyimpangan `IDNewBisnis`" di V-48: kedua penyimpangan itu
   satu gejala yang sama.
3. **Perbaikan murah:** minta ekspor ulang `NB-176005` dan `NB-184183`. Bila hasilnya berbeda, ini
   cacat pengumpulan contoh dan bukan cacat data Pega.

```powershell
Get-ChildItem 'D:\migrasi\RNM\DDL\CONTOH' -Filter *.xml | Get-FileHash -Algorithm SHA256 |
  Group-Object Hash | Where-Object Count -gt 1
```

---

## 6. §4.1 dikuatkan — `<ID>` akar adalah KUNCI KERJA, bukan rujukan master

`[terverifikasi]` Atas ke-50 berkas ber-`<ID>` bentuk-kelas:

| | Nilai |
| --- | ---: |
| Kelas Pega — **satu-satunya** yang muncul: `ASM-FW-GISFW-WORK` | 50 |
| Ekor = **nomor work object berkas itu sendiri** | **48** |
| Ekor menunjuk work object **lain** | **2** (persis dua pasang duplikat §5) |
| Nilai unik | 48 |

Jadi bentuknya bukan sekadar "mirip kunci kerja" — **ia memang kunci kerja berkas itu sendiri**
(bentuk `pzInsKey` seperti dicatat V-48).

⚠️ **Tambahan yang membatasi tafsir:** ke-50 nilai itu **tidak muncul sebagai isi medan lain mana
pun** di berkas yang sama — jadi ia **bukan** salinan `PolicyMasterIDPega`. Klaim V-49a *"tidak satu
pun nilainya muncul sebagai medan lain"* juga tereproduksi untuk ke-38 yang berbentuk angka: **0
dari 38**.

⛔ Kesimpulan laporan berdiri: memetakan `<ID>` akar → `SOURCE_ID` secara seragam menaruh **kunci
kerja** ke kolom rujukan master di **50 dari 114** berkas. V-49a hanya menghitung yang **38**, jadi
ke-50 itu tidak pernah masuk pertimbangan.

---

## 7. Butir 8 — pertentangan nyata, tetapi ada jalan **D**

`[terverifikasi]` Kedua kutipan benar dan memang tidak dapat berlaku bersama.

Yang belum dipertimbangkan laporan: **ADR-0001 sudah meramalkan gejala V-28.** Bunyinya —

> *"Bila urutan operasi dan presisi per-langkah direproduksi apa adanya, `decimal` bersifat
> deterministik dan hasilnya **harus** identik. Selisih sekecil apa pun berarti ada salah-port."*

Selisih 2,6 per sepuluh juta pada penjumlahan 195 baris adalah **tepat** bentuk sisa yang timbul
bila urutan pembulatan tidak direproduksi. Menurut ADR-0001 itu **bukan alasan memakai toleransi —
itu gejala port yang belum tepat**, dan ADR-0005 adalah mekanisme yang seharusnya menutupnya.

| | Jalan | Akibat |
| :-: | --- | --- |
| A | Buang tabel `Total*`, **amandemen ADR-0001** | Ukuran keberhasilan migrasi berubah |
| B | Simpan nilai `Total*` apa adanya, **tanpa hitung ulang** | ADR-0001 utuh |
| C | Beda perlakuan data baru vs lama | Skema jadi dua perlakuan |
| **D** | **Reproduksi urutan dan presisi pembulatan Pega**, lalu bandingkan eksak | ADR-0001 utuh, tabel `Total*` tetap boleh dibuang |

⚠️ Jalan D bergantung pada satu fakta yang **belum diketahui**: apakah uji aritmetika V-28
mereproduksi pembulatan per-langkah, atau menjumlah sekali di akhir. **Pertanyaan itu perlu diajukan
kepada penyusun V-28 sebelum A/B/C dipilih** — bila jawabannya "menjumlah sekali", maka tidak ada
pertentangan sama sekali dan tidak ada ADR yang perlu diubah.

📌 Bila harus memilih sekarang tanpa jawaban itu: **B** — satu-satunya yang tidak mengubah ukuran
keberhasilan migrasi maupun menambah perlakuan ganda.

---

## 8. Tambahan untuk butir 9 dan 5b

**Butir 9 — tabel akar sendiri ikut terdampak.** Laporan menyebut `T_SPREADINGLIST`,
`T_COVERAGELIST` dan `T_FR_CURRENCYLIST`. `[terverifikasi]` **`T_GENERAL_POLIS` sendiri punya 9 kolom
uang tanpa mata uang yang dapat ditelusuri — di keenam grup.** Itu tabel akar penawaran, bukan cabang.

Rincian lengkap 24 tabel tak tertelusur ada di §9.

**Butir 5b — ada preseden di skema Oracle yang hidup.** `[terverifikasi]` `JSON_POLIS` sudah memuat
kolom **`OLDNOPOLIS VARCHAR2(100)`**. Jadi gagasan `OLD_POLIS_ID` (penunjuk langsung ke generasi
sebelumnya) **bukan hal baru bagi basis data ini** — bentuk yang serupa sudah ada di tabel yang sama
yang menyimpan `PRODKE`. Ini memperkuat arah 5b; **bukan** keputusan.

---

## 9. Rincian butir 9 — 24 tabel, mata uang tak tertelusur

`[terverifikasi]` Tabel ber-kolom uang yang tidak punya kolom mata uang sendiri **dan** tidak punya
leluhur ber-mata-uang, per grup pemakai:

| Tabel | Kolom uang | Grup |
| --- | ---: | --- |
| `T_FR_CURRENCYLIST` | 18 | Aneka FIRE MarineCargo PA |
| `T_COVERAGELIST` | 17 | Life MBUCar MarineCargo PA |
| `T_FR_COVERAGELIST` | 10 | Aneka MarineCargo PA |
| `T_GENERAL_POLIS` | **9** | **keenam grup** |
| `T_FR_PAYMENT` | 9 | MarineCargo PA |
| `T_FR_LISTINSTALLMENT` | 8 | MarineCargo PA |
| `T_SPREADINGLIST` | 5 | Life MBUCar MarineCargo PA |
| `T_FR_SPREADINGLIST` | 5 | Aneka MarineCargo PA |
| `T_RETROLIST` | 5 | Life |
| `T_FR_NETPERCURRENCY` | 4 | Aneka FIRE MarineCargo PA |
| `T_COVERAGEDATALIST` · `T_FR_ANEKALIST` · `T_FR_OFFEREDPAYMENT` · `T_FR_PERSONLIST` | 3 | — |
| `T_ADDITIONALCOVERAGE` · `T_CEDING_CURRENCYLIST` · `T_FR_CARGOLIST` · `T_FR_FACOUTOBJECTLIST` · `T_FR_LOCATIONLIST` · `T_FR_SECURITYREINSURER` · `T_FR_TOTALTSIPREMIRETRO` | 2 | — |
| `T_FR_PRINTRISLIP` · `T_FR_PROPERTY` · `T_PERSONLIST` | 1 | — |

⚠️ Batas ukuran sama dengan laporan: "kolom uang" dikenali dari **nama dan tipe**, bukan arti
bisnisnya. Arahnya tidak berubah oleh penyetelan ambang.

---

## 10. Yang tetap TIDAK dapat diverifikasi

| # | Klaim | Sebab |
| :-: | --- | --- |
| 1 | Angka V-47 · R1 · R1b · R5 · R7 | Tidak diukur ulang di sesi ini |
| 2 | Uji aritmetika V-28 dan selisih 2,6 per sepuluh juta | Tidak diukur ulang; lihat §7 — yang diperlukan justru **metodenya**, bukan angkanya |
| 3 | Apakah Fac In memakai **tabel fisik yang sama** dengan Treaty In | Diagram tidak menyatakan; kolom Treaty In belum direkonsiliasi dengan 114 contoh |
| 4 | Apakah seluruh nilai `PRODKE` lama berupa angka | Isinya di basis data — **butir DBA** |
| 5 | Apakah migrasi memuat seluruh versi polis atau versi terakhir saja | Keputusan lingkup |
| 6 | Apakah struktur RNW sama dengan NB | **Nol contoh RNW** — tetap lubang terbesar kedua |
| 7 | Daftar nilai `POSISI` yang sah | V-48c: hanya DDL, nol baris data |

---

## 11. Urutan yang saya sarankan

Berubah dari urutan laporan, karena dua butir kini tidak lagi memblokir:

1. **Butir 9** — tetap nomor satu. Menyentuh K-010, K-012, ADR-0004, ADR-0006.
2. **Butir 8** — ajukan **satu pertanyaan** ke penyusun V-28 (§7) sebelum memilih A/B/C/D.
3. **Butir 5b** — `OLD_POLIS_ID` (didukung preseden `JSON_POLIS.OLDNOPOLIS`) atau pencarian `PROD_KE`.
4. **Butir 1 lanjutan** — tabel fisik bersama Treaty In atau terpisah.
5. **`IsCedingConfirm`** — kolom sendiri atau `POSISI` (menggantikan butir 7 lama).

**Tidak lagi memblokir:** butir **3** (selesai, §4) dan butir **7** (selesai mengikuti V-31, §3).
Butir **2**, **4**, **6** tetap dapat mengikuti rekomendasi laporan.

**Dua tindakan di luar keputusan, dapat dikerjakan sekarang:**
`[a]` kembalikan `Diagram-Skema-Tabel-NusantaraRe.xlsx` ke `DDL\` (§1);
`[b]` minta ekspor ulang `NB-176005` dan `NB-184183` (§5).

---

## 12. Perintah audit pokok

```powershell
# 114 contoh; 101 NB / 13 EDM / 0 RNW; duplikat byte-identik
$f=@([IO.Directory]::GetFiles('D:\migrasi\RNM\DDL\CONTOH','*.xml')); $f.Count
Get-ChildItem 'D:\migrasi\RNM\DDL\CONTOH' -Filter *.xml | Get-FileHash -Algorithm SHA256 |
  Group-Object Hash | Where-Object Count -gt 1

# pxListSubscript: 18.140 punya / 0 kosong / 0 beda dari posisi  (skrip penuh: lihat catatan sesi)
# baris tanpa tag: SpreadingList 134, CoverageList 564, DeductibleList 346,
#                  PropertyItemList 106, LocationList 74, AnekaList 33, ViewSuggest 101

# kunci versi dan preseden OLD_POLIS_ID
Select-String -Path 'D:\migrasi\RNM\DDL\JSON_POLIS.txt' -Pattern 'PRIMARY KEY|UNIQUE|PRODKE|OLDNOPOLIS'

# 15 kolom HISTORYAKSEPTASIPRODUCTION
Select-String -Path 'D:\migrasi\RNM\DDL\HISTORYAKSEPTASIPRODUCTION.txt' -Pattern '"[A-Z0-9_]+"\s+(VARCHAR2|NUMBER|DATE|CHAR)'
```

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan. Seluruh skrip yang menyapu `DDL\CONTOH\`
hanya mengeluarkan hitungan, nama tag, nama kelas Pega, dan nomor work object. Tidak ada keputusan
diambil di dokumen ini.*
