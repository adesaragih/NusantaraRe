# Pertanyaan terbuka — layar Treaty In disamakan dengan Pega (5 Oktober 2026)

**Diajukan kepada pemilik proses.** Lima pertanyaan, dan tiap satunya punya pengukuran di
belakangnya. Yang **tidak** dilakukan ronde ini: menebak jawabannya.

---

## 1 · Empat label `Accounting Mode` — rule-nya tidak diekspor

Briefing meminta teks tampil, bukan nilai simpan: *"yang tampil itu **label**, bukan nilai
tersimpan."* Setuju — dan labelnya **tidak ada di korpus**.

**Yang terukur.** Kedua dropdown memakai `pyListSource = associated`, artinya daftar
pilihannya datang dari **rule propertinya sendiri** (`Rule-Obj-Property`, `pyTableOption`),
bukan dari Section. Dan korpus `D:\XML_NURE` **tidak punya folder `Property`** — 1.048 berkas
XML disapu, nol di antaranya `Rule-Obj-Property` untuk kedua properti ini.

**Yang membuktikan bentuknya.** Dua `Rule-Obj-Property` memang ada, di
`_migration-docs\treaty-in-adjustment\ekspor-tambahan\`, dan keduanya memperlihatkan persis
tempat label itu tinggal:

```xml
<pyPromptTableList REPEATINGTYPE="PageList">
  <rowdata REPEATINGINDEX="1">
    <pyLocalizedValue>Internal</pyLocalizedValue>   <pyStandardValue>1</pyStandardValue>
  <rowdata REPEATINGINDEX="2">
    <pyLocalizedValue>External</pyLocalizedValue>   <pyStandardValue>2</pyStandardValue>
```

| Properti | `pyTableOption` | Label yang diekspor |
| --- | --- | --- |
| `EDMState` | `PromptList` | `1 = Internal` · `2 = External` |
| `EDMMaterialType` | `PromptList` | `1 = Material` · `2 = Non Material` |
| `AccountingMode` | — | ⛔ **rule-nya tidak diekspor** |
| `AccountingModeNonProp` | — | ⛔ **rule-nya tidak diekspor** |
| `Bordeaux` | — | ⛔ **rule-nya tidak diekspor** |

⛔ **Jadi nilainya ditampilkan apa adanya** — `underwriting` · `accounting` · `loss` · `risk` ·
`reporting` · `nonreporting`. Mengarang "Underwriting Year" dan "Accounting Year" akan
menuliskan tebakan ke layar yang orang percayai.

**Yang diminta:** ekspor tiga `Rule-Obj-Property` itu dengan cara yang **persis sama** seperti
`EDMState.xml` sudah diekspor. Begitu `pyPromptTableList`-nya ada, pemetaannya dipasang di satu
tempat (`FORM_KONTRAK` di `labels.ts`) dan selesai.

---

## 2 · `05/1` di medan Treaty Year — itu bukan tahun

Briefing: *"JSONDATA menyimpan tahun biasa; layar Pega menampilkan `05/1`."*

**Rumusnya ketemu**, dan ia tidak menghasilkan `05/1`. `DataTransform/TreatyInSetTreatyYear.xml`
(`pyRuleAvailable = Yes`):

```
TreatyIn.TreatyYear  := @substring(TreatyIn.Commencement,0,4)
TreatyIn.Termination := @addCalendar(TreatyIn.Commencement,"1","0","0","0","0","0","0")
```

**Diadu dengan POOLDATA:**

| Yang diuji | Cocok |
| --- | ---: |
| `TREATYYEAR = SUBSTR(COMMENCEMENT,1,4)` | **1.849** dari 1.854 |
| `TERMINATION = COMMENCEMENT + 12 bulan` | **9** dari 1.854 |

⭐ Baris pertama mengesahkan rumusnya. Baris kedua membuktikan keduanya **nilai awal saat
masuk**, bukan aturan yang terus berlaku — kalau `Termination` sungguh terikat, ia akan cocok
pada ribuan baris.

⛔ **Dan `05/1` tetap tidak terjelaskan.** Pada kontrak yang briefing tunjuk (`1001856`),
`TREATYYEAR` di POOLDATA bernilai **`2026`**, dengan `COMMENCEMENT = 20260917`. Jadi `05/1` di
tangkapan layar datang dari tempat lain.

**Yang ditanyakan:** apakah `05/1` benar-benar di medan **Treaty Year**, atau di medan
sebelahnya? Dan kalau benar di sana — medan itu di Pega menampilkan apa?

---

## 3 · `RNM Share` — menurut ekspor ia PANEL, bukan tab

Briefing menyebutnya tab ke-11 yang hilang dari tangkapan layar non-prop. Ekspornya berkata
lain, dan ronde ini **tidak menghapusnya** — §0 melarang menghapus tab mana pun.

**Yang terukur** di `Section/TreatyInTabsNonProportional.xml` (bersih 4.203.858 bita):

| | |
| --- | ---: |
| wadah `pyHeaderType = TABBED` | **11** |
| di antaranya berjudul `RNM Share` | **0** |

Kedua judul `RNM Share` (@2.093.710, @2.133.893) ber-`BAR`, dan keduanya jatuh **di antara** tab
`Share` (TABBED @1.695.720) dan tab `Retro` (TABBED @3.291.822) — wilayah tab Share, bersama
`Summarry of RNM Share` @2.324.385 dan `Total All Layers RNM Share` @2.486.274.

⚠️ *"BAR berarti bukan tab"* **tidak** berlaku umum: di `TreatyInTabsProportional.xml` nol wadah
ber-TABBED, dan kesebelas tabnya justru `BAR`. Yang membedakan pemakaian **di dalam satu
berkas**.

⭐ **Dan hitungannya cocok:** 11 tab − `Value Difference` (syaratnya tidak terpenuhi) = **10**,
persis yang tangkapan layar perlihatkan. Dengan `RNM Share` ikut, daftar kita 12.

**Yang ditanyakan:** `RNM Share` turun menjadi panel di dalam tab `Share`, atau ia memang tab di
Pega yang berjalan dan ekspornya yang tidak mewakilinya?

---

## 4 · `EDMMaterialType` nol di seluruh 1.854 dokumen

Syarat tab `Value Difference` dibaca utuh dari ekspor — `pyContainerVisibleWhen` yang sama muncul
di **tiga** berkas:

```
TreatyIn.EDMState != 3 && TreatyIn.EDMMaterialType == 1
```

**Diadu dengan POOLDATA, 1.854 dokumen:**

| Kunci | Ada di |
| --- | ---: |
| `EDMState` | **2** dokumen (keduanya proporsional, bernilai `0`) |
| `EDMMaterialType` | **0** dokumen |

⛔ Jadi lewat `JSONDATA` syarat itu **tidak pernah** terpenuhi, dan tab `Value Difference` tidak
pernah tampil — konsisten dengan tangkapan layar.

⚠️ **Itu tidak berarti tabnya mati.** Pega menilai syaratnya atas *clipboard*, yang dapat diisi
Activity saat jalan tanpa pernah tersimpan. Yang terukur hanya bahwa nilainya tidak **tersimpan**.

**Yang ditanyakan:** `EDMState` dan `EDMMaterialType` diisi dari mana saat Pega jalan — dan
adakah kontrak yang tab `Value Difference`-nya sungguh pernah terlihat? Satu nomor kontrak sudah
cukup.

⚠️ Pertanyaan yang sama berlaku bagi **`Retro` proporsional**: syaratnya
`TreatyIn.IsMultipleRetro`, dan kelima kontrak yang bernilai `"true"` (`1000493`, `1000755`,
`1001043`, `1001405`, `1001853`) **seluruhnya NON-proporsional**. Nol kontrak proporsional dari
1.079 bernilai `"true"`. Jadi tab `Retro` proporsional pun tidak pernah tampil hari ini.

---

## 5 · `1001856` adalah kontrak DRAFT — tangkapan layarnya kosong karena isinya kosong

Briefing meminta buktinya dijalankan atas kontrak ini. Dijalankan, dan hasilnya perlu dinyatakan:

| | `1001856` | kontrak nyata |
| --- | ---: | ---: |
| nama kontrak | `dasdas` | *nama kontrak sebenarnya* |
| panjang `JSONDATA` | **1.764** bita | 33.000 – 56.000 bita |
| kunci kepala yang ada | **0** dari 8 | 8 dari 8 |

Nol `Bordeaux`, nol `AccountingMode`, nol `AccountingModeNonProp`, nol `TreatyYear` di
dokumennya. `COMMENCEMENT = 20260917`.

⛔ **Jadi medan yang kosong di tangkapan layar itu kosong karena KONTRAKNYA kosong**, bukan
karena layar membaca salah. Bukti §0 tetap dijalankan atasnya (uji
`TestKontrakBuktiAdalahDraftKosong` menjaga sifat itu), tetapi ia **tidak dapat** membuktikan
medan mana yang terisi benar.

**Yang ditanyakan:** adakah satu nomor kontrak non-proporsional yang **terisi lengkap** untuk
dipakai sebagai kontrak bukti ronde berikutnya? Saran dari data: `1001853` (`WHOLE ACCOUNT RISK
AND CATASTROPHE EXCESS OF LOSS REINSURANCE 2025`, `JSONDATA` 52.824 bita).

---

## 6 · `RNM as Treaty Leader` bersyarat `TreatyMasterInEDM`, dan syaratnya salah di setiap kontrak

`When/TreatyMasterInEDM.xml` (`pyRuleAvailable = Yes`, memo *"added edmstate"*):

```
TreatyIn.EDMState = "1"   atau   = "2"   atau   = "3"
```

Sel `TreatyIn.TreatyLeader` @191.054 berpenjaga rule itu (`pyCondition` @248.313), dan `EDMState`
ada di **2** dari 1.854 dokumen — keduanya bernilai `0`. Jadi menurut ekspor, kotak centang
`RNM as Treaty Leader` **tidak pernah tampil**.

⛔ **Medannya TIDAK disembunyikan.** `JSONDATA.TreatyLeader` ada di **659** kontrak dengan nilai
sungguhan (`true` 103 · `false` 556); menyembunyikan medan yang punya data nyata di sepertiga
kontrak, atas dasar properti yang asalnya belum jelas (§4), adalah tebakan yang mahal.

**Yang ditanyakan:** apakah kotak centang itu memang hanya tampil pada kontrak ber-EDM — dan
kalau ya, 659 kontrak itu memperoleh nilainya dari mana?

---

# ⭐ RALAT 5 Oktober 2026 — dokumen desain menutup §2 dan §3 di atas

`D:\NUSANTARA RE APP\design treaty in.docx` (43 tangkapan layar, kontrak `1001846` dan
`1001841`) menjawab dua pertanyaan di atas. Keduanya **dibiarkan tertulis** supaya jalan
berpikirnya terbaca; jawabannya di sini.

| Pertanyaan | Jawaban |
| --- | --- |
| §2 `05/1` di medan Treaty Year | ⛔ **Premisnya gugur.** Gambar 01 memperlihatkan `Treaty Year 2025` — tahun polos. `05/1` artefak form yang belum tersimpan, dan pemilik proses menyatakannya sendiri. Rumus `@substring(Commencement,0,4)` tetap berlaku sebagai nilai awal. |
| §3 `RNM Share` panel atau tab | ⭐ **Sub-tab di dalam `Share`**, kedua cabang — gambar `16`/`17` (prop) dan `34` (non-prop). Sejalan dengan pengukuran ekspor. Dibangun sebagai sub-tab; **nol tab dihapus**. |

⚠️ §1 (label `Accounting Mode`) **separuh terjawab** — lihat §7 di bawah.
⚠️ §4, §5, dan §6 **tetap terbuka**, dan gambar justru memperkuat §6: lihat §8 di bawah.

---

## 7 · Label `risk` untuk `AccountingModeNonProp` — tiga dari empat ketemu

Gambar menutup tiga padanan, dan satu tetap kosong:

| Tersimpan | Label | Sumber |
| --- | --- | --- |
| `underwriting` | **Underwriting Year** | ⭐ gambar `01`, kontrak 1001846 |
| `accounting` | **Accounting Year** | ⭐ pasangan yang pemilik proses nyatakan sendiri |
| `loss` | **Loss Occuring** | ⭐ gambar `26`, kontrak 1001841 |
| `risk` | ⛔ **tidak diketahui** | 21 dari 1.854 kontrak |

⛔ **`risk` tidak ditebak**, dan pencariannya dicatat supaya tidak diulang:

- Nol dari 43 gambar memperlihatkan kontrak ber-`risk`.
- Sapuan 1.048 berkas XML `D:\XML_NURE` menemukan **"Risk Attaching" tepat sekali** — di
  `Claim Non Prop\Activity\SetEndDate_Act.xml`, di dalam `pyMemo` **berbahasa Indonesia**:
  *"Memperbaiki kondisi IF Risk Attaching periode mulai polis (ceding) tidak boleh melebihi
  periode treaty"*. Itu catatan pengembang di **modul lain**, bukan label properti ini.
- `Rule-Obj-Property` untuk `AccountingModeNonProp` tetap **tidak diekspor**.

⚠️ Uji `TestRiskTidakDitebak` menjaga nilainya tampil apa adanya, dan ia **sengaja merah** pada
hari seseorang memasang tebakan.

**Yang diminta:** buka satu kontrak ber-`risk` di Pega dan kirim satu tangkapan layar medan
`Accounting Mode`-nya — atau ekspor `Rule-Obj-Property` `AccountingModeNonProp` seperti
`EDMState.xml` sudah diekspor.

**Ketiga kontrak ber-`risk` yang paling mudah dibuka** (dari 21 yang terukur): `1000033`,
`1000221`, `1000227`. Daftar lengkapnya: 1000033 · 1000221 · 1000227 · 1000228 · 1000264 ·
1000298 · 1000391 · 1000402 · 1000414 · 1000418 · 1000419 · 1000439 · 1000450 · 1000459 ·
1000679 · 1000680 · 1000694 · 1000780 · 1001013 · 1001014 · 1001101.

### 7.1 Dan `nonreporting` juga belum terlihat

Padanan `reporting` -> **"Reporting"** terbaca di gambar `01`. Pasangannya **`nonreporting`
tidak muncul di satu pun dari 43 gambar** — kedua kontrak contoh bernilai `reporting` atau tidak
punya medannya.

⛔ "Non Reporting" terdengar jelas dan tetap tebakan; **687 kontrak** memakainya. Satu tangkapan
layar sudah cukup menutupnya.

---

## 8 · Tab `Limits` proporsional bersarang TIGA tingkat — dan datanya belum punya rumah

Empat belas gambar (`03`–`15`) untuk satu tab, dan sebabnya struktural. Pega merender:

```
Limits → Kind of Treaty → Treaty Type → Treaty Group → Class of Business
       → 100% Limit / Retention / Cession to R/I  (masing-masing berdaftar mata uang)
       → 11 sub-tab (Event Limits … Achievement)
```

Layar kita merender **satu grid datar** dari `M_TREATY_IN2`, yang memuat **satu baris per
layer** — bukan pohon ini. Entitas yang hilang: `Kind of Treaty`, `Treaty Type`, dan
`Class of Business` per treaty group.

**Yang ditanyakan:**

1. Ketiga entitas itu tersimpan di mana? Nol jejaknya di `M_TREATY_IN2` maupun di kunci puncak
   `JSONDATA` yang sudah dipetakan.
2. Apakah tab Limits proporsional masuk lingkup migrasi, atau ia menunggu tiket tersendiri?
   Empat belas gambar untuk satu tab menyarankan yang kedua.

---

## 9 · ⚠️ Aturan angka tunggal TERBANTAH oleh gambarnya sendiri

§7 briefing menetapkan satu jumlah desimal per jenis (uang 4 · persen 2 · persen share 8).
Gambar memperlihatkan jumlah desimal yang **berbeda-beda per medan**, dan dua medan bernilai
sama di baris yang sama pun berbeda:

| Gambar | Medan | Tampil | Desimal |
| --- | --- | --- | ---: |
| 29 | EGNPI `Amount` | `137.849.315.068,00` | 2 |
| 29 | EGNPI `Amount in IDR` | `137.849.315.068` | **0** |
| 29 | `Proportion %` di grid | `100,00` | 2 |
| 29 | `Proportion %` di rincian | `100,0000000000` | **10** |
| 30 | Layers `100% Limits ( IDR )` | `750.000.000,00` | 2 |
| 30 | `Summary of Limit` `100% Limit (IDR)` | `750.000.000` | **0** |
| 32 | `Adjustment Rate %` | `0,367` | **3** |
| 38 | `% Installment` | `25,00` | 2 |
| 38 | `% Total` | `100,0000` | **4** |

Dan **nol di ekor TIDAK dibuang** di layar lama: `1,00` tetap `1,00` (gambar 01), `5.345,00`
tetap `5.345,00` (gambar 05). Aturan kita membuangnya.

⛔ **Nol baris aturan angka diubah ronde ini**, dan sebabnya: yang menentukan di Pega adalah
`pyDecimalPlaces` tiap kontrol, dan **nilai itu tidak ikut diekspor**. Yang dapat disimpulkan
dari gambar adalah bahwa aturan tunggal itu salah — bukan aturan benarnya apa. Menebaknya akan
mengubah setiap angka di setiap grid modul ini sekaligus.

**Yang ditanyakan:**

1. Apakah jumlah desimal memang per medan, dan kalau ya — daftarnya diambil dari mana?
2. Apakah nol di ekor dipertahankan (`1,00`) seperti layar lama, atau dibuang (`1`) seperti
   aturan kita hari ini? Keduanya tidak dapat berlaku sekaligus.

---

## 10 · Sebelas nama kategori lampiran ketemu — pasangan KODE-nya tetap tidak

Gambar `24` (prop) dan `42` (non-prop) memberi kesebelas namanya, dan menemukan hal yang
sebelumnya tidak diketahui: **daftarnya BERCABANG**. Sepuluh nama identik; yang kesepuluh
berbeda — `Pega Proportional Calculation` lawan `Pega Non Proportional Calculation`.

⛔ **Ini TIDAK menutup `PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md`.** Yang terbuka adalah
pasangan **kode↔nama**, dan gambar hanya memberi namanya. Urutannya di layar alfabetis, dan
ronde 4 Oktober 2026 melarang keras menyimpulkan kode dari urutan abjad.

**Yang diminta:** satu tangkapan layar Pega yang memperlihatkan kode kategori di sampingnya,
atau isi tabel kategori lampiran.

---

## 11 · ⛔ §23 DIHENTIKAN — bacaan "padankan sampai presisi kolom" TIDAK cocok dengan ronde 69

Briefing 5 Oktober 2026 (§23a) meminta: *"Baca keputusan ronde 69 itu sendiri, dan kalau bacaan
ini tidak cocok dengan bunyinya — berhenti dan katakan."*

**Dibaca. Tidak cocok. Berhenti.**

### 11.1 Bunyi ronde 69 sesungguhnya

Keputusan itu tinggal di `inti/frontend/lib/format.ts` baris 69–79, di atas `formatNumber`:

> **# Nol di EKOR tidak ditampilkan (ronde 69)**
>
> Versi sebelumnya memadankan pecahan ke **panjang tetap**, sehingga 2484250 tampil
> `"2.484.250,000000"` **pada batas 6** dan setiap nilai bulat memperoleh ekor nol. Pemilik
> proyek menolaknya dengan menunjuk sistem existing:
>
> > `"DI PEGA 2484250 MAKA DI SISTEM BARU JANGAN 2484250.000000000"`
>
> Yang dipangkas hanya nol yang TIDAK BERMAKNA — tidak satu digit berarti pun hilang, dan nilai
> tersimpan tidak disentuh.

Dan di baris 63–65, **sebelum** itu:

> `desimal` adalah **BATAS ATAS, bukan panjang tetap**. Nilainya ditiru per jenis:
> **6 uang realisasi · 3 share/komisi · 2 limit/TSI · 1 Pct arrangement** ·
> `DESIMAL_TAK_DIBATASI` untuk Rp/Usd arrangement.

### 11.2 ⛔ Mengapa bacaan briefing tidak cocok

Briefing mengusulkan: *"aturannya bukan 'jangan pernah memadankan nol' melainkan 'padankan sampai
presisi kolomnya, bukan sampai batas global'."*

**Premisnya keliru: ronde 69 tidak pernah memadankan ke batas GLOBAL.** Ia sudah memadankan ke
batas **per jenis kolom** — 6 untuk uang realisasi, 3 untuk share, 2 untuk limit, 1 untuk Pct.
Angka `6` dalam kalimat *"pada batas 6"* **adalah** presisi kolom uang realisasi, bukan satu
batas untuk seluruh layar.

Jadi yang pemilik proyek tolak pada ronde 69 persis perilaku yang §23 kini minta dibangun:
**memadankan nol sampai presisi kolom**.

| | Ronde 69 | §23 ronde ini |
| --- | --- | --- |
| Presisi | per jenis kolom (6 · 3 · 2 · 1) | per kolom (2 · 0 · 10 · 3 · 4) |
| Nol di ekor | **dibuang** | **dipertahankan** |
| Contoh | `2484250` → `2.484.250` | `1` → `1,00` |

⭐ Yang berbeda hanya **sumber angkanya** (jenis lawan kolom). Yang **bertentangan** adalah
perlakuan nol di ekor — dan itu justru butir yang pemilik proyek nyatakan dengan kalimat
langsung.

### 11.3 ⚠️ Satu bacaan yang MUNGKIN cocok — dan ia milik pemilik proses, bukan milik saya

Ronde 69 menyebut layarnya sendiri: *"Dipakai untuk `To IDR`, `Limit (%)`, dan kolom `Rp`/`Usd`
Arrangement"*, dan `formatSel` menyebut *"Laporan Realisasi (view existing), Limit Treaty In, dan
grid XOL"*. **Nol di antaranya layar Treaty In yang ronde ini kerjakan.**

Jadi keduanya dapat berdiri bersama **bila** ronde 69 dibaca sebagai keputusan untuk **ketiga
layar itu**, bukan untuk seluruh aplikasi — dan §23 berlaku untuk layar Treaty In saja.

⛔ **Saya tidak memilih bacaan itu sendiri.** Ia mengubah jangkauan sebuah keputusan pemilik
proyek, dan §23a melarang menimpa keputusan lama diam-diam. **Nol baris aturan angka diubah
ronde ini.**

### 11.4 Yang ditanyakan — satu pertanyaan, dua jawaban yang mungkin

1. **Ronde 69 berlaku untuk seluruh aplikasi** → maka §23 bertentangan dengannya, dan yang satu
   harus dicabut. Mana?
2. **Ronde 69 berlaku untuk tiga layar yang ia sebut** (Laporan Realisasi · Limit Treaty In ·
   grid XOL) → maka §23 dapat dibangun untuk layar Treaty In tanpa menyentuh ketiganya, dan
   pemadanannya di lapis modul seperti §23b minta.

⚠️ Bila jawabannya **2**, satu hal masih kurang sebelum kodenya dapat ditulis: §23 memberi
desimal untuk **5 kolom** yang gambar perlihatkan. Modul ini punya **puluhan** kolom angka.
Sisanya *"tetap memakai aturan lama"* — tetapi aturan lama membuang nol di ekor, sehingga satu
grid akan memuat kolom yang memadankan nol **di sebelah** kolom yang membuangnya. Itu bentuk yang
tidak satu pun dari 43 gambar perlihatkan.

Yang diminta untuk menutupnya: **ekspor `pyDecimalPlaces` tiap kontrol** — permintaan yang sama
yang sudah ditagih di §12 di bawah.

### 11.5 ⭐ §23b sudah terjawab sebelum dibangun

*"Pemadanan dikerjakan di mana, dan bukti `format.ts` tidak tersentuh."*

Bila §23 kelak dibangun, pemadanannya di **lapis modul** — `selAngka` di
`modul/treatyin/frontend/pages/FormKontrakTreatyIn.tsx`, sesudah `formatNumber` dipanggil. Itu
memadankan hasil, bukan memformat ulang, dan nol baris `format.ts` berubah.

`format.ts` **tidak tersentuh ronde ini**: stempel waktunya tetap 4 Oktober 2026 17:22, sama
dengan `styles.css` dan `dasar.tsx`.

---

## 12 · Empat permintaan ekspor — satu daftar, satu bentuk contoh

Keempatnya `Rule-Obj-Property`, dan ketiganya sudah ditagih di §1 dan §7.1. Dikumpulkan di sini
supaya satu permintaan menutup seluruhnya:

| Rule | Untuk | Yang dibutuhkan |
| --- | --- | --- |
| `AccountingMode` | label `underwriting` / `accounting` | `pyPromptTableList` |
| `AccountingModeNonProp` | label `loss` / **`risk`** | `pyPromptTableList` |
| `Bordeaux` | label `reporting` / **`nonreporting`** | `pyPromptTableList` |
| **tiap kontrol angka** | `pyDecimalPlaces` per kolom | setelan kontrolnya |

⭐ **Bentuk contohnya sudah ada di repositori**:
`D:\XML_NURE\_migration-docs\treaty-in-adjustment\ekspor-tambahan\EDMState.xml` memperlihatkan
persis blok yang dibutuhkan:

```xml
<pyPromptTableList REPEATINGTYPE="PageList">
  <rowdata REPEATINGINDEX="1">
    <pyLocalizedValue>Internal</pyLocalizedValue>   <pyStandardValue>1</pyStandardValue>
  <rowdata REPEATINGINDEX="2">
    <pyLocalizedValue>External</pyLocalizedValue>   <pyStandardValue>2</pyStandardValue>
```

Tiga rule properti diekspor dengan cara yang sama menutup §1, §7, dan §7.1 sekaligus.
`pyDecimalPlaces` menutup §9 dan §11.

---

## 13 · ⭐ §11 TERJAWAB — dan daftar kolom yang desimalnya masih gelap

Pemilik proses menjawab §11 dengan **§24**: ronde 69 berlaku pada ketiga layar yang ia sebut
saja (Laporan Realisasi · Limit Treaty In · grid XOL), dan form Treaty In mengikuti gambar.
Jawaban nomor **2**. §23 dibangun.

⛔ **Yang masih gelap, dan sengaja tidak ditebak:**

| Kelompok kolom | Keadaan |
| --- | --- |
| Seluruh kolom tab **Limits** (14) | nol dari 43 gambar memperlihatkan presisinya |
| Seluruh kolom tab **Share** (11) | idem |
| Seluruh kolom tab **RNM Share** (7) | idem |
| `Proportion %` di RINCIAN baris EGNPI | **10** desimal terbaca (gambar 29) — tetapi layar kita belum punya baris yang dapat dibuka |
| `Adjustment Rate %` · `ROL %` rincian layer | **3** terbaca (gambar 32) — medan rincian, belum dirender |
| `% Total` Installment | **4** terbaca (gambar 38) — medan panel, bukan kolom grid |

Ketiganya yang terbaca **dicatat tanpa dipasang**, sebab tempatnya belum ada. Yang tidak terbaca
memakai aturan lama dan `desimalPadan` mengembalikan `null`.

**Yang menutupnya tetap satu hal:** ekspor `pyDecimalPlaces` tiap kontrol — permintaan ke-4 di
daftar §12.

---

## 14 · `Treaty Type` di tingkat `Kind of Treaty` — calonnya ada, tingkatnya tidak cocok

Gambar `03` memperlihatkan dropdown `Treaty Type` berbunyi **`2025 SPL 66M FAC`**, duduk di
dalam baris `Kind of Treaty` dan **di atas** grid `Treaty Group`.

Calon terkuatnya `Limits[].Detail[].SpreadingType` — 10 nilai yang bentuknya persis sama:

```
2023 QS 150M TRT 523 · QS HR 60M TRT 514 · 2022 QS 145M TRT 498
2020 QS 101M TRT 470 · 2019 QS 101M TRT 250 · … 4 nilai lain
```

⛔ **Tetapi tingkatnya tidak cocok.** `SpreadingType` hidup di `Detail[]` — satu nilai per
**treaty group** — sementara layar menaruhnya satu tingkat di atas, milik seluruh
`Kind of Treaty`. Memasangnya berarti memilih `Detail[0]` dan menyebutnya milik kelompok.

⚠️ Pada kontrak dengan satu treaty group keduanya sama dan tebakan itu tidak terlihat; pada
kontrak bertreaty-group banyak ia menampilkan nilai yang salah untuk semua kecuali yang pertama —
dan tidak ada di layar yang memberi tahu yang mana.

**Yang ditanyakan:**

1. `Treaty Type` di gambar 03 memang `SpreadingType` milik treaty group pertama, atau ia medan
   lain di tingkat `Kind of Treaty` yang belum saya temukan?
2. Kalau yang pertama — apa yang Pega tampilkan ketika satu `Kind of Treaty` punya dua treaty
   group dengan `SpreadingType` berbeda?

---

## 15 · Asal pilihan `Accounting Mode` / `Bordereaux` — disapu sampai database (6 Oktober 2026)

Gejala: dropdown berbunyi *"Underwriting Year (tidak ada di daftar referensi)"*. Sebabnya
kode, bukan data: form mengisi dropdown dengan LABEL hasil `CaraPembukuanTampil`, sementara
pilihannya NILAI tersimpan (`underwriting`). **Diperbaiki:** dropdown memegang nilai tersimpan
(`caraPembukuanAsli`, `bordereauxAsli`); pilihannya pasangan nilai↔label dari
`services.OpsiKepalaKontrak` (rute `GET /api/treaty-in/warisan/opsi-kepala`), dengan label
dari penerjemah yang SAMA.

Asal daftar pilihannya, disapu:

| Tempat | Hasil |
|---|---|
| Section `TreatyInNONProportional.xml` | `pxDropdown`, `pyListSource = associated` — daftar milik rule Property |
| korpus `D:\XML_NURE\Treaty In` (329 berkas XML) | nol `Rule-Obj-Property`; teks "Underwriting Year" NOL kali |
| `DATAPEGA.PR_ASM_FW_GISFW_DATA_ENUMERATI` (1.217 baris, 65 jenis) | nol baris accounting mode; satu-satunya `REPORTING` milik `methodtreatyarr` |
| tabel `POOLDATA` bernama ACCOUNT/REF/LOOKUP/PARAM/MODE/LOV/CODE | nol yang memuatnya |
| skema rule Pega (`PR4_RULE*`) | tidak terjangkau dari DSN aplikasi |

⛔ **Jadi daftar resminya hanya ada di rule Property Pega.** Yang dipakai: domain TERSIMPAN di
1.854 dokumen + label dari tangkapan layar. `risk` dan `nonreporting` tetap berlabel nilainya.
**Yang menutupnya:** ekspor `Rule-Obj-Property` `AccountingMode`, `AccountingModeNonProp`,
`Bordeaux`.

## 16 · Mode Edit — fungsi hidup, TETAPI `Save` belum menyimpan (6 Oktober 2026)

Permintaan pemakai: tombol `Edit` di daftar membuka form yang seluruh fungsinya dapat dipakai
(Add di tiap tabel), `View` baca-saja. **Dibangun:** mode `ubah`/`lihat`; di mode ubah medan,
grid, Rate of Exchange, dan pohon Limits P dapat disunting, Add/Delete bekerja. Di mode lihat
`<fieldset disabled>` mematikan isian dan tombol Add/Delete tidak dirender (ekspor:
`TreatyIn.ViewState !='1'`).

⚠️ **Suntingan hidup di layar saja.** `Save` tetap mati: tabel warisan (`TREATY_IN`,
`M_TREATY_IN`) BACA SAJA, dan jalur simpan ke model baru belum diputuskan. Tombol yang
menjalankan Activity — `Update Total`, `Refresh` Achievement, unggah/unduh lampiran,
`Actions` — juga belum hidup.

**Yang diminta:** keputusan ke mana `Save` menulis.


## 17 · Tombol `Apply` Reporting Period — dibangun; medan kepalanya belum punya sumber (6 Oktober 2026)

⭐ **Apply kini menghitung**, menyalin `Activity/TreatyInSetReport.xml` (dipicu tombol
`TreatyInTabsProportional.xml` @254615, `startdate = ReportingStart`, `autocalculate = true`).
Rumusnya di `services/periode_pelaporan.go`, diukur ulang terhadap **4.548 baris** yang Pega
simpan: **4.524 cocok (99,5%)**, jumlah baris cocok pada 1.135 dari 1.137 kontrak. Lampiran
pemakai (06/10/2026, Quarter Year, 12/12/12 → 17/01, 29/01, 10/02/2027) menjadi uji.

Dua sifat Pega yang ditiru, keduanya terbukti dari data:

- **Jatuh tempo "kurang satu hari"**: Pega menghitung di WIB lalu menyimpan Date dari jam GMT
  (tengah malam WIB = 17:00 GMT hari sebelumnya).
- **Pemotongan akhir bulan menumpuk**: tiap periode ditambah dari `TempDate` sebelumnya.

⚠️ Syarat `ReportingPeriod=="other"` pada langkah 13–14 Activity **tidak aktif**
(`pyStepsPreCondition = false`) — langkahnya berlaku untuk semua periode.

**Yang belum:**

1. ~~**Medan kepala kosong saat tab dibuka.**~~ — ⭐ **DITUTUP 6 Oktober 2026, migrasi `444`.**

   Bunyi butir ini: *"tersimpan hanya di `JSONDATA` — dilarang dibaca layar — dan
   `T_TREATY_REVISION` tidak punya kolomnya."* Separuh pertamanya mengikat; separuh
   keduanya hanya kolom yang belum dibuat, dan `444` membuatnya.

   ⛔ Yang MENUNDA butir ini bukan keputusan yang kurang melainkan PENGUKURAN yang
   belum ada. Sapuan 6 Oktober 2026 atas **seluruh 1.855** dokumen `M_TREATY_IN`
   (nol `ROWNUM`, nol gagal urai) memberi angkanya:

   | Kunci akar | Terisi di |
   | --- | ---: |
   | `ReportingPeriod` | 1.851 |
   | `ReportingStart` · `End` · `Submission` · `Confirmation` · `Settlement` | 1.219 masing-masing |
   | `ReportingInterval` | 300 |

   ⚠️ **Dua pertiga korpus.** Selama ini tab Reporting Period memperlihatkan kepala
   kosong pada seribu dua ratus kontrak yang sungguh punya nilainya — bukan kasus
   pinggiran.

   Migrasi `444` menambahkan ketujuh kolom itu plus sembilan lagi yang lolos saringan
   yang sama (kunci yang menjadi PARAMETER prosedur `POOLDATA.PEGA_TREATY_IN`, yaitu
   yang Pega sendiri perlakukan sebagai data tersimpan). Sudah dijalankan di POOLDATA
   dan pemuatnya sudah mengisinya.
2. ~~**Label pilihan Period**~~ — **terjawab**: Prompt List Property `ReportingPeriod`
   (tangkapan rule dari pemakai): `quarter` → Quarter Year · `half` → Half Year ·
   `month` → Monthly · `other` → Others.
3. Activity tidak memeriksa `Interval`; tanpanya pembagian langkah 12 gagal di Pega. Layar ini
   menolak dengan *"Interval must not be empty"* — kalimat dari label layar, bukan dari Activity.

## 18. Dropdown `Treaty Type` tab Limits (proporsional) — 6 Oktober 2026

**Asal nilai (dari ekspor, bukan tebakan):** `Section/LimitProportional.xml`, medan
`.TreatyTypeID` `pxDropdown`, `pyListSource=reportdefinition` →
`ReportDefinition/BrowseReinsuranceType_RD.xml` (kelas `ASM-FW-GISFW-Int-REINSURANCETYPE`):
nilai `.ID`, label `.Note`, parameter `Flag="active"`, urut `.ID` DESC, maks. 500.
Kelas itu = tabel `POOLDATA.REINSURANCETYPE` — dibuktikan SQL mentah kelas yang sama di
`Claim Non Prop/RDBList/GetReinsuranceTypeBYName_SQL.xml`: `select ID, NOTE as "Note" from
REINSURANCETYPE where … FLAG='active'` (130 baris, 84 `active`, 83 nama unik). Kode
tersimpan `10042`/`10035`/`10037` = SURPLUS / QUOTA SHARE / 2ND SURPLUS.
Saat berubah, `Activity/SetTreatyTypeName_Act.xml` mengisi `.TreatyType` (`Kind of Treaty`)
dengan `.Note` pilihan itu — layar meniru hal yang sama.

**Terpasang:** `GET /api/treaty-in/warisan/jenis-treaty` (baca saja; `REINSURANCETYPE`
masuk `TestWarisanHanyaDibaca`).

**Yang perlu diketahui:** ekspor menandai sel ini `pyEditOptions=Read-only` tanpa syarat.
Di sini dropdown hidup di mode Edit atas permintaan pemakai; di mode View tampil label
`.Note` baca-saja, seperti dropdown read-only Pega.

### 18a. Dropdown lain di tab Limits — mengikuti layar Pega (lampiran pemakai)

| Medan | RD (ekspor) | Tabel | Nilai → label | Saat berubah |
|---|---|---|---|---|
| Treaty Group `.TreatyGroupID` | `BrowseTreatyGroup_RD` | `TREATYGROUP` (33) | `ID` → `TreatyGroupName`, ID DESC | `SetTreatyGroupName_Act` → `.TreatyGroup` |
| Mata uang grid 100% Limit · Retention · Cession · Reserve · PLA · Cash Loss · Claim Coop · EPI · Event Limits | `BrowseCurrencyTreatyIn_RD` | `CURRENCY` tanpa `ITL` | `.Currency` → `.Currency`, tanpa urutan | — |
| Mata uang grid Deduction `.CurrencyID` | `BrowseCurrency_RD` | `CURRENCY` tanpa `ITL` | `.ID` → `.Currency` | `SetCurrName_Act` → `.Currency` |

Rute: `GET /api/treaty-in/warisan/opsi-limits` (ketiganya sekaligus). Tombol baris grid
DetailLimits berlabel `Remove`, baris Kind of Treaty/Treaty Group `Delete` (ekspor). Kind of
Treaty tidak punya isian sendiri (`.TreatyType` diisi dari Treaty Type).

**Belum dibangun:** `FetchQSfromMaster` (Treaty Group berubah → isi SpreadingList dari master
treaty) dan `LimitCalculation` (mata uang/nilai berubah) — keduanya hitungan lintas tab.
