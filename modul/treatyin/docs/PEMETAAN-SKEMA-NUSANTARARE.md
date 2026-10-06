# Pemetaan `Diagram-Skema-Tabel-NusantaraRe.xlsx` ↔ tabel modul `treatyin`

**Acuan:** `D:\NUSANTARA RE APP\Diagram-Skema-Tabel-NusantaraRe.xlsx`, 2 Oktober 2026.
**Dipetakan:** 5 Oktober 2026. ⛔ **Nol migrasi, nol `CREATE`/`ALTER`/`RENAME`/`DROP`, nol
perubahan jalur baca, model, layar, maupun uji.** Keluaran ronde ini hanya dokumen ini.

---

## 0 · ⭐⭐ TEMUAN YANG MENENTUKAN SELURUH SISANYA

**Berkas acuan dan modul `treatyin` memodelkan OBJEK PEGA YANG BERBEDA.**

| | Berkas acuan | Modul `treatyin` |
| --- | --- | --- |
| Objek sumber | **`pyWorkPage.PolicyTreatyIn`** | **`TreatyIn.*`** (`ASM-FW-GISFW-Int-TREATY_IN`) |
| Tabel sumber | `POOLDATA.json_polis` | `TREATY_IN` + `M_TREATY_IN.JSONDATA` |
| Yang dimodelkan | **POLIS** Treaty In (NB dan EDM) | **KONTRAK MASTER** Treaty In |
| Kunci alami | `(NOPOLIS, PRODKE)` | `TREATY_IN.ID` |

**Buktinya di berkas acuan sendiri**, bukan tafsiran:

| Lembar · sel | Bunyi |
| --- | --- |
| `NB Treaty In NonProp` **F10** | *"diambil dari `pyWorkPage.PolicyTreatyIn` — 79 medan skalar tingkat atas"* |
| `NB Treaty In NonProp` **F11** | *"dan `POOLDATA.json_polis` — 7 kolom yang sudah datar"* |
| `NB Treaty In Prop` **B85** | *"**TABEL LAMA YANG SUDAH DATAR — di luar pohon, tidak dirancang ulang**"* |
| `NB Treaty In Prop` **F98** | di bawah judul itu: *"HISTORYAKSEPTASIPEGA · JSON_POLIS_MONITORING · **TREATY_IN** · ACHIEVEMENT · **M_TREATY_IN**"* |
| `NB Treaty In NonProp` **F118** | *"dibaca saja: … **M_TREATY_IN_EDM** …"* |

⭐ **Kedua tabel yang menjadi SATU-SATUNYA sumber modul `treatyin` — `TREATY_IN` dan
`M_TREATY_IN` — berdiri di berkas acuan sebagai "tabel lama yang sudah datar, DI LUAR POHON,
tidak dirancang ulang".**

Itu menjelaskan setiap angka nol di bawah sekaligus: berkas acuan **tidak memodelkan** apa yang
modul ini modelkan. Nol padanan bukan karena nama berbeda — melainkan karena lingkupnya berbeda.

> ## ⭐ TERJAWAB 5 Oktober 2026 — kemungkinan **2**
>
> Pemilik proses memutuskan `treatyin` **ADA DI DALAM** lingkup, dan pengujian terhadap data
> menyisihkan dua dari tiga kemungkinan: irisan kunci puncak hanya **14 dari 312/140** (tujuh di
> antaranya perabot Pega), dan satu kontrak master memegang sampai **385 polis**. Rinciannya di
> [`RENCANA-PINDAH-SKEMA.md`](RENCANA-PINDAH-SKEMA.md) §0.
>
> Ketiga daftar di bawah karena itu berarti **"pekerjaan belum dimulai"**, bukan "beda lingkup".

⛔ **Dan ini DINYATAKAN, bukan diputuskan.** Dua kemungkinan, dan dari dalam kode keduanya
terlihat sama persis:

1. Modul `treatyin` memang **di luar lingkup** model acuan, dan tetap membaca `M_TREATY_IN`
   selamanya. Berkas acuan sendiri berkata begitu.
2. Modul `treatyin` **kelak dipindahkan** ke model acuan, dan berkas ini belum memuat
   pemetaannya.

**Pemilik proses yang memutuskan.** Keputusan itu menentukan apakah ketiga daftar di bawah
adalah "beda lingkup" atau "pekerjaan yang belum dimulai".

⚠️ `rancangan-tabel-datar-treaty-in.md` §4q.9 menguatkan kemungkinan **1**:
*"`M_TREATY_IN.JSONDATA` — **di luar lingkup sekarang**"* `[keputusan work owner]`.

---

## 1 · Daftar PADANAN

⛔ **NOL padanan yang memenuhi syarat bukti.**

Aturannya: *"Mirip pada nama bukan padanan; sebut buktinya — medan sumber yang sama, induk yang
sama, atau kunci alami yang sama."* Diterapkan atas ke-51 tabel kita lawan ke-40 tabel acuan,
**nol pasangan lulus** — sebab medan sumbernya berbeda secara sistematis (`TreatyIn.*` lawan
`PolicyTreatyIn.*`, §0).

### 1.1 Calon yang MIRIP tetapi TIDAK lulus — disebut supaya tidak dicari dua kali

| Tabel kita | Calon acuan | Lembar · sel | Mengapa TIDAK lulus |
| --- | --- | --- | --- |
| `LAYER` · `KELOMPOK_LAYER` · `KELAS_BISNIS_LAYER` | `T_POLIS_XOL_LAYER` | NB NonProp **K131**, R82–R87 | sumbernya `TreatyXOLList().ValueList()`; kita `TreatyIn.Limits[]` — objek berbeda |
| `PENYEBARAN` · `RINCIAN_PENYEBARAN` | `T_POLIS_SPREADING` | NB NonProp **K129**, J71–J74 | sumbernya `PolicyTreatyIn.SpreadingRiskList()`; kita `Limits[].Detail[].SpreadingList[]` |
| `TERMIN` · `RINCIAN_ANGSURAN` | `T_POLIS_INSTALMENT` · `_DETAIL` | NB NonProp **K126**, **K127** | sumbernya `PolicyTreatyIn.ListInstallment()`; kita `TreatyIn.Installment[]` |
| `NILAI_SELISIH` · `NILAI_SEBELUM_PRO_RATE` | `T_POLIS_DIFFERENCE` dkk. | EDM NonProp **K182–K185** | induk berbeda: mereka `POLIS_ID` + `NOURUT`; kita `ID_VERSI_KONTRAK` + besaran + mata uang |
| `CATATAN_PERSETUJUAN` | `T_VIEW_SUGGEST` | Daftar Relasi **D18**, **D31** | sumbernya `PolicyTreatyIn.SuggestList`; kita `M_TREATYIN_COMMENT` (kini `T_VIEW_COMMENT`) |
| `T_VIEW_COMMENT` (migrasi 436) | `T_VIEW_SUGGEST` | Daftar Relasi **D18** | ⚠️ **nama hampir sama, tabel berbeda** — lihat §4.3 |

⚠️ Kelima baris pertama **mungkin menjadi padanan** bila jawaban §0 adalah kemungkinan 2. Hari
ini bukti yang disyaratkan tidak ada, jadi keduanya dicatat sebagai calon — bukan padanan.

---

## 2 · NOL PADANAN, MILIK KITA — ke-51 tabel

Sumber: `modul/treatyin/docs/DAFTAR-TABEL-DIBUAT.txt`.

⛔ **Untuk tiap kelompok dinyatakan dua kemungkinan, tidak diputuskan satu pun** — dari dalam
kode "model baru tidak memerlukannya" dan "berkasnya belum memuatnya" terlihat identik.

### 2.1 Sembilan tabel tab — dinamai ulang migrasi `436`

`T_TREATY_REPORTING_PERIOD` · `T_TREATY_PORTFOLIO` · `T_TREATY_ACCUMULATION` · `T_TREATY_EGNPI` ·
`T_TREATY_RETENTION` · `T_TREATY_INSTALLMENT` · `T_TREATY_INSTALLMENT_ITEM` · `T_VIEW_COMMENT` ·
`M_TREATYIN_COINSCALE`

**Padanan: NOL.** Awalan `T_TREATY_` **nol kali** di seluruh berkas acuan (terukur, 10 lembar).

- *Kemungkinan A* — model acuan memang tidak memodelkan tab kontrak master, sebab ia memodelkan
  polis. Didukung §0.
- *Kemungkinan B* — berkas belum memuatnya, dan nama `T_TREATY_*` kelak bertabrakan atau
  diselaraskan.

⚠️ Keputusannya sekaligus menutup `TestTCONolNamaTabelBaruDiKode` — lihat §4.3.

### 2.2 Entitas ERD kontrak master — 33 tabel

`BAGIAN` · `BAHAYA` · `BATAS_PER_BAHAYA` · `CATATAN_PERSETUJUAN` · `DETAIL_PROPORSIONAL` ·
`DOKUMEN_KONTRAK` · `EGNPI` · `JENIS_POTONGAN` · `JENIS_REASURANSI` · `KELAS_BISNIS` ·
`KELAS_BISNIS_KELOMPOK` · `KELAS_BISNIS_LAYER` · `KELOMPOK_LAYER` · `KELOMPOK_TREATY` ·
`KONTRAK` · `LAYER` · `NILAI_CADANGAN_PREMI` · `NILAI_MDP` · `NILAI_MDP_MINIMUM` ·
`NILAI_PENYEBARAN` · `NILAI_PREMI_BRUTO` · `NILAI_PREMI_BRUTO_MINIMUM` · `PEMULIHAN_LIMIT` ·
`PENCAPAIAN` · `PENYEBARAN` · `PERIODE_AKUMULASI` · `PERIODE_PELAPORAN` · `PORTOFOLIO` ·
`POTONGAN` · `RETENSI_CEDANT` · `RINCIAN_ANGSURAN` · `RINCIAN_PENYEBARAN` · `SKALA_KOASURANSI` ·
`TERMIN` · `VERSI_KONTRAK`

**Padanan: NOL** menurut syarat bukti. Keduanya memodelkan hal yang namanya mirip pada
lapisan dagang, dari dua objek Pega yang berbeda.

### 2.3 Dua tabel modul Adjustment

`NILAI_SELISIH` · `NILAI_SEBELUM_PRO_RATE`

⭐ **Di sinilah padanan paling mungkin kelak muncul**, dan berkas acuan memecahkan soal yang sama
dengan struktur yang berbeda — lihat §5.1 (`OLD_POLIS_ID`) dan §5.2 (`NOURUT`).

### 2.4 Tabel infrastruktur migrasi — 4 tabel

`ARSIP_MUATAN_KELUAR` · `JEJAK_PERUBAHAN` · `MIGRASI_KORELASI` · `MIGRASI_NILAI_DITOLAK` ·
`MIGRASI_PENDARATAN`

⭐ Untuk kelima ini **kemungkinan A hampir pasti**: keduanya milik proses migrasi kita sendiri,
dan berkas acuan memodelkan keadaan akhir, bukan perkakas perpindahannya.

---

## 3 · NOL PADANAN, MILIK BERKAS ITU — ke-40 tabel

⚠️ **Terukur 40, bukan 44.** Briefing menyebut 44; sapuan atas kesepuluh lembar dengan pola
`T_[A-Z0-9_]+` menemukan **40** nama berbeda. Ke-20 token huruf besar lain yang ikut tersaring
pola longgar seluruhnya **nama kolom** (`POLIS_ID`, `OLD_POLIS_ID`, `QUOTATION_ID`, …), bukan
tabel. **Ukuran saya yang dipakai; silakan diadu.**

| Tabel | Lembar · sel (tiga pertama) |
| --- | --- |
| `T_CLAIMLF_ADJUSTMENT` | Daftar Relasi C8, C13, D14 |
| `T_CLAIMLF_ADJUSTMENT_SPREADING` | Daftar Relasi C9, D8 |
| `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` | Daftar Relasi D9 |
| `T_CLAIMLF_PREMIUMLIST_DETAIL` | Daftar Relasi C7, C10, D6 |
| `T_CLAIM_ADJUSTMENT` | Claim Fac In B83; Daftar Relasi C33, C34 |
| `T_CLAIM_ADJ_LOSS_ALLOCATION` | Daftar Relasi D35 |
| `T_CLAIM_ADJ_QUOTA_SHARE` | Claim Fac In B86; Daftar Relasi D34, D49 |
| `T_CLAIM_ADJ_SPREADING` | Claim Fac In B85; Daftar Relasi D33, D48 |
| `T_CLAIM_BREAK_QS` | Claim Fac In B82; Daftar Relasi D29, D46 |
| `T_CLAIM_CLAIM_AMOUNT` | Daftar Relasi D27 |
| `T_CLAIM_ESTIMATION` | Claim Fac In B80; Daftar Relasi D25, D44 |
| `T_CLAIM_FAC_RETRO` | Claim Fac In B84; Daftar Relasi D30, D50 |
| `T_CLAIM_INTEREST` | Daftar Relasi D26 |
| `T_CLAIM_OBJECT` | Daftar Relasi C43, D42 |
| `T_CLAIM_OBJECT_ITEM` | Daftar Relasi C44, C45, C46 |
| `T_CLAIM_SPREADING` | Claim Fac In B81; Daftar Relasi D28, D45 |
| `T_GENERAL_CLAIM` | Daftar Relasi C25, C26, C27 |
| `T_GENERAL_KOMITE` | Daftar Relasi C12, C39, C54 |
| **`T_GENERAL_POLIS`** | Daftar Relasi C74, C76, C77; **NB Prop K108** · **NB NonProp K123** · **EDM Prop K147** · **EDM NonProp K172** |
| `T_KOMITE_KOMITELIST` | Daftar Relasi D12, D39, D54 |
| ⛔ `T_POLIS_BREAKDOWN_SPREAD` | Daftar Relasi D77 — **DIBATALKAN 23-09-2026** (NB NonProp B4, J63) |
| **`T_POLIS_CEDING`** | Daftar Relasi D75; NB Prop K110 · NB NonProp K125 · EDM Prop K149 · EDM NonProp K174 |
| **`T_POLIS_DIFFERENCE`** | Daftar Relasi C88–C90; NB Prop K115 · EDM NonProp K182 |
| **`T_POLIS_DIFFERENCE_INSTALMENT`** | Daftar Relasi D89; EDM NonProp K184 |
| **`T_POLIS_DIFFERENCE_SPREADING`** | Daftar Relasi D88; EDM NonProp K183 |
| **`T_POLIS_INSTALMENT`** | Daftar Relasi C82, D76; NB NonProp K126 |
| **`T_POLIS_INSTALMENT_DETAIL`** | Daftar Relasi D82; NB NonProp K127 — **hanya NonProp** |
| **`T_POLIS_QUOTATION`** | Daftar Relasi C75, D74; NB NonProp K124 |
| **`T_POLIS_SPREADING`** | Daftar Relasi D78; NB NonProp K129 |
| **`T_POLIS_XOL`** | Daftar Relasi C84, D83; NB NonProp K130 — **hanya NonProp** |
| **`T_POLIS_XOL_LAYER`** | Daftar Relasi D84; NB NonProp K131 — **hanya NonProp** |
| **`T_POLIS_XOL_LAYER_DIFFERENCE`** | Daftar Relasi D90; EDM NonProp K185 |
| `T_PREMIUM_LIST` | Daftar Relasi C16–C18 |
| `T_PREMIUM_LIST_DETAIL` | Daftar Relasi C19, C21, D17 |
| `T_PREMIUM_LIST_SPREADING` | Daftar Relasi C20, D19 |
| `T_PREMIUM_LIST_SPREADING_RETRO` | Daftar Relasi D20 |
| `T_PREMIUM_LIST_SUMMARY` | Daftar Relasi D16; PremiumList NB + EDM B51 |
| `T_VIEW_SUGGEST` | Daftar Relasi D18, D31 |
| `T_WORK_CLAIM` | Daftar Relasi C11, C14, C23 |
| **`T_WORK_POLIS`** | Daftar Relasi C15, C72, C73; NB Prop K107 · NB NonProp K122 · EDM Prop K146 · EDM NonProp K171 |

**Bertebal = 13 tabel yang menyentuh Treaty In.** Sisanya milik Claim Life, Claim Fac In,
Komite, dan PremiumList.

---

## 4 · BENTROK

### 4.1 ⭐ `LAYER` — **BUKAN bentrok; keduanya SEPAKAT**

Briefing meminta ini diperiksa, dan jawabannya menghapus satu kekhawatiran.

> **Acuan** (`NB Treaty In NonProp` **F26–F27**):
> *"⛔ LAYER · LAYER_TYPE · LAYER_PART · LAYER_PART_TYPE **DICORET dari tabel ini** — di sistem
> lama ia PANTULAN, dibaca balik dari kolom tabel lewat `pxResults(1)`, dan `pxResults(1)` hanya
> mengambil baris PERTAMA ⇒ nilainya **layer pertama saja**, bukan ringkasan seluruh layer."*
>
> Dan **R85**: *"LAYER · LAYER_TYPE · LAYER_PART · LAYER_PART_TYPE **ADA DI SINI SAJA** — 12
> rujukan masing-masing, di induk NOL."* (tabel `T_POLIS_XOL_LAYER`)

> **Yang dibangun** (`repository/warisan_layer_dokumen.go`): pohon `Limits[]` → `.Detail[]` →
> `.COBList[]`, **setiap** elemen dibaca — 4.210 elemen layer, 5.716 silang layer × treaty group.

⭐ **Keduanya menolak hal yang sama**: nilai layer di tingkat induk/polis, yang hanya pantulan
baris pertama. Keduanya menyimpan layer di tingkat anaknya sendiri dan membaca **seluruhnya**.
Yang acuan sebut `T_POLIS_XOL_LAYER`, kita baca sebagai elemen `Limits[]`.

⛔ **Jadi ini bukan bentrok, dan tidak perlu diputuskan.**

### 4.2 ⚠️ PRESISI ANGKA — bentrok yang NYATA, keduanya dikutip

> **Acuan** (`NB Treaty In NonProp` **F20**):
> *"uang · persen → angka presisi tetap, skala **MINIMAL 9 desimal** (P29: galat lama diikuti apa
> adanya, tidak dibulatkan)"*
>
> Dan `rancangan-tabel-datar-treaty-in.md` §4q.1:
> *"Nilai lama dipindahkan **utuh termasuk ekor galatnya**. Tidak dibulatkan saat migrasi."*
> *"skala kolom uang: **minimal 9 desimal** — `NUMBER(p,9)` atau lebih"*
> *"pembandingan uang: ⛔ **tidak boleh sama-persis** — bertoleransi, atau dibandingkan
> terbulatkan"*

> **Yang dibangun** (KEPUTUSAN §23/§24, 5 Oktober 2026):
> desimal **per kolom**, dibaca dari gambar desain — `Value to IDR` **2** · EGNPI `Amount` **2** ·
> `Amount in IDR` **0** · Maximum Retention `Amount` **0** · Installment `Amount`/`Pct` **2**;
> nol di ekor **DIPERTAHANKAN** sampai presisi kolomnya.

**Yang perlu diperhatikan sebelum memilih** — dan ini pengamatan, bukan rekomendasi:

⚠️ Ketiga kutipan acuan berbicara tentang **kolom, migrasi, dan pembandingan** — yaitu
**penyimpanan**. §24 berbicara tentang **tampilan**. Secara mekanis keduanya dapat berdiri
bersama: kolom `NUMBER(p,9)` yang menyimpan `2484250.000000001` tetap dapat ditampilkan
`2.484.250,00`.

⛔ **Tetapi ketegangannya nyata dan harus dinyatakan**: tampilan 2 desimal **menyembunyikan**
ekor galat yang penyimpanan sengaja pertahankan. Siapa yang memeriksa selisih di layar tidak
akan melihatnya. P29 tidak mengatakan apa pun tentang tampilan, jadi ia tidak melarangnya — dan
juga tidak mengizinkannya.

**Nol rekomendasi.** Pemilik proses yang memutuskan.

### 4.3 ⚠️ NAMA TABEL TAB — bentrok yang menunggu §0

> **Acuan**: awalan `T_TREATY_` muncul **NOL KALI** di kesepuluh lembar. Nol padanan untuk
> kesembilan tabel tab.

> **Yang dijalankan**: migrasi **`436`** memberi nama `T_TREATY_*` dan `T_VIEW_COMMENT` kepada
> kesembilannya, dan **sudah dijalankan** terhadap POOLDATA.

Dua kemungkinan, apa adanya:

1. **Nama diganti lagi** mengikuti acuan baru — menuntut migrasi berikutnya, dan menuntut
   jawaban apa nama yang benar (acuan tidak memberinya).
2. **Acuan memang tidak menjangkau tabel tab** — §0 mendukung ini — dan `T_TREATY_*` berdiri
   sebagai nama modul ini sendiri.

⚠️ Keputusan ini **sekaligus menutup** `TestTCONolNamaTabelBaruDiKode`, yang hari ini merah
karena menganggap awalan `T_TREATY_*` milik tabel warisan modul lain. ⛔ Penjaganya **tidak
disentuh** ronde ini.

⚠️ Dan satu tabrakan nama yang perlu diketahui lebih dulu: migrasi `436` menamai
`M_TREATYIN_COMMENT` menjadi **`T_VIEW_COMMENT`**, sementara acuan punya **`T_VIEW_SUGGEST`**
(Daftar Relasi D18, D31) — **dua tabel berbeda dengan pola nama yang sama**. Milik kita memuat
komentar kontrak master; milik acuan memuat `PolicyTreatyIn.SuggestList`.

### 4.4 Bentrok keempat yang ditemukan ronde ini — `T_POLIS_BREAKDOWN_SPREAD`

Berkas acuan memuatnya di `Daftar Relasi` **D77**, tetapi **membatalkannya** di lembar Treaty In
(`NB NonProp` **B4**: *"T_POLIS_BREAKDOWN_SPREAD DIBATALKAN — keempat medannya turunan, tidak
disimpan"*; **J63**: *"⛔ DIBATALKAN 23-09-2026"*).

⚠️ Jadi berkas acuan **tidak konsisten dengan dirinya sendiri**: lembar relasi masih memuatnya,
lembar rancangan mencabutnya. Dicatat agar pembaca berikutnya tidak menghitungnya sebagai tabel
hidup. Bukan bentrok dengan pekerjaan kita — kita tidak punya padanannya.

---

## 5 · ⭐ Dua soal kita yang acuan sudah pecahkan — dicatat, belum diadopsi

### 5.1 `OLD_POLIS_ID` menggantikan `OldData`

> `NB Treaty In NonProp` **F13–F17**:
> *"`OLD_POLIS_ID` → `T_WORK_POLIS.ID` · nullable · penunjuk ke generasi SEBELUMNYA"*
> *"⭐ `OLD_POLIS_ID` adalah pengganti `OldData`. Sistem lama tidak punya penunjuk ini — ia
> mencari tiap kali (`FetchPolisJsonPolis`: order by `PRODKE` desc fetch first 1), dan pencarian
> keduanya memakai `substr(nopolis, 1, 24)` — nomor polis dianggap selalu 24 karakter. FK
> menghapus jebakan itu"*
> *"`UNIQUE (OLD_POLIS_ID)` melarang percabangan · `UNIQUE (NOPOLIS, PRODKE)` mencegah dua
> endorsemen serentak mendapat nomor sama"*
> *"⛔ baris generasi lampau TIDAK BOLEH disunting — itulah pembekuan `OldData` (P58)"*

Itu panel **Old Data / New Data** modul Adjustment, diselesaikan secara struktur.

### 5.2 `NOURUT` — kunci pemasangan antar generasi

> `NB Treaty In NonProp` **F28–F32**:
> *"⭐⭐ ATURAN NOURUT — setiap tabel anak punya kolom `NOURUT`, dan itulah kunci pemasangan antar
> generasi: **selisih[n] = baru[n] − lama[n]**"*
> *"NB (PRODKE 0), selagi disusun: baris boleh DIHAPUS, dan `NOURUT` DINOMORI ULANG supaya rapat
> 1..n — aman, sebab belum ada generasi untuk dipasangkan"*
> *"EDM: baris TIDAK BISA DIHAPUS. `NOURUT` lama terbawa apa adanya, baris baru ditambahkan di
> belakang dengan `NOURUT` = maksimum + 1"*
> *"pembatalan (P56) menolkan nilainya, BARISNYA TETAP ADA"*
> *"⭐ aturan keutuhan yang dapat ditegakkan: generasi n+1 **WAJIB** memuat setiap `NOURUT` yang
> ada di generasi n. Yang hilang = endorsemen tidak sah, ditolak di services"*

Itu **Value Difference**, berikut aturannya.

⚠️ Keduanya **dicatat, bukan diadopsi**. Mengadopsinya menuntut keputusan §0 lebih dulu: kalau
modul `treatyin` di luar lingkup model acuan, `NOURUT` dan `OLD_POLIS_ID` tidak punya tempat di
tabelnya.

---

## 6 · ⭐ Berkas kolom lengkap — **DITEMUKAN**, tetapi bukan di tempat yang disebut

Berkas acuan menunjuk `.scratch/nb-treaty-in/rancangan-tabel-datar-treaty-in.md`
(`NB Treaty In NonProp` **B2**). Jalur itu **tidak ada** di cakram.

⭐ **Isinya ada**, di dalam repositori: **`modul/nbtreatyin/docs/rancangan-tabel-datar-treaty-in.md`**
— 984 baris, lengkap dengan daftar kolom per tabel (§4.1–§4.6) dan 12 keputusan work owner
(§4q.1–§4q.12).

⚠️ **Tetapi daftar kolomnya memakai NAMA TABEL YANG LAMA** — `TREATY_IN_POLIS`,
`TREATY_IN_POLIS_XOL`, `TREATY_IN_POLIS_ANGSURAN`, `TREATY_IN_POLIS_SPREADING` — sementara
berkas acuan memakai `T_GENERAL_POLIS`, `T_POLIS_XOL`, `T_POLIS_INSTALMENT`, `T_POLIS_SPREADING`.
Pergantian namanya terjadi di antara keduanya dan **tidak dicatat di salah satu pun**.

⛔ **Akibatnya pemetaan tingkat KOLOM tetap tidak dapat dikerjakan hari ini** — bukan karena
daftarnya hilang, melainkan karena padanan nama tabel lama↔baru belum ada. Itu **batas yang
dinyatakan**, bukan hasil yang disamarkan.

**Yang diminta:** satu tabel padanan `TREATY_IN_POLIS*` ↔ `T_*`, dari pemilik berkas acuan.

---

## 7 · Permintaan yang terkumpul — SATU daftar

⛔ Jangan dipecah; daftar inilah yang dikirim ke tim Pega dan ke pemilik berkas acuan.

| # | Yang diminta | Untuk | Bentuk contoh |
| --: | --- | --- | --- |
| 1 | Rule `Rule-Obj-Property` **`AccountingMode`** | label `underwriting` / `accounting` | `ekspor-tambahan/EDMState.xml` |
| 2 | Rule `Rule-Obj-Property` **`AccountingModeNonProp`** | label `loss` / **`risk`** (21 kontrak) | idem |
| 3 | Rule `Rule-Obj-Property` **`Bordeaux`** | label `reporting` / **`nonreporting`** (687 kontrak) | idem |
| 4 | **`pyDecimalPlaces`** tiap kontrol angka | 32 kolom Limits/Share/RNM Share yang desimalnya `null` | setelan kontrolnya |
| 5 | ⭐ **Padanan nama `TREATY_IN_POLIS*` ↔ `T_*`** | memetakan tingkat kolom (§6) | satu tabel dua kolom |
| 6 | ⚠️ **Ketidakkonsistenan `T_POLIS_BREAKDOWN_SPREAD`** (§4.4) | agar tidak dihitung sebagai tabel hidup | satu kalimat |
| 7 | ⭐ **Bentuk akar kontrak master** — `RENCANA-PINDAH-SKEMA.md` §0.3 | memblokir seluruh perpindahan | nama tabel + induknya |
| 8 | ⭐ **Cacat dua penjaga lintas-aplikasi** — `RENCANA-PINDAH-SKEMA.md` §5.1 | agar `./inti/...` hijau | menerima `RENAME`, mengecualikannya dari syarat ber-skema |

⚠️ Butir **6 lama** (*"apakah `treatyin` di dalam lingkup"*) **sudah terjawab** — lihat kotak §0.

Bentuk yang permintaan 1–3 butuhkan sudah ada contohnya di
`D:\XML_NURE\_migration-docs\treaty-in-adjustment\ekspor-tambahan\EDMState.xml`:

```xml
<pyPromptTableList REPEATINGTYPE="PageList">
  <rowdata REPEATINGINDEX="1">
    <pyLocalizedValue>Internal</pyLocalizedValue>   <pyStandardValue>1</pyStandardValue>
  <rowdata REPEATINGINDEX="2">
    <pyLocalizedValue>External</pyLocalizedValue>   <pyStandardValue>2</pyStandardValue>
```

---

## 8 · Cara memeriksa ulang berkas ini

⚠️ **Dua perangkap urutan atribut**, dan keduanya membuat berkasnya tampak kosong:

1. `<sheet>` memuat `xmlns:r="…"` **sebelum** `name="…"` → regex `<sheet name="` memberi **nol**
   lembar. Pakai `<sheet\b[^>]*?\bname="([^"]+)"`.
2. `<Relationship>` memuat `Target="…"` **sebelum** `Id="…"` → regex `Id="…"[^>]*Target="…"`
   memberi `KeyError`. Pakai dua pencarian terpisah atau balik urutannya.

Dan: berkasnya **tanpa `sharedStrings.xml`** — teksnya inline di `<is><t>`, jadi pembaca yang
hanya menengok tabel string bersama akan membaca nol sel.

Terukur: 10 lembar · 1.852 sel terisi · **40** nama `T_*` berbeda · `T_TREATY_*` **nol**.
