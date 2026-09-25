# PROMPT — spec penyimpanan relasional, EDM Treaty In

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
>
> ⛔ **Prasyarat:** `.scratch\nb-treaty-in\spec-penyimpanan-relasional.md` sudah ada.
> Bila belum, **berhenti** — spec EDM membangun di atasnya, bukan di sampingnya.
>
> ⚠️ **Skill `to-spec` tidak dapat dipanggil sendiri oleh agen** (CLAUDE.md §8). Ketikkan
> `/to-spec` sebagai manusia, lalu berkas ini menjadi briefnya.

---

## 0. LINGKUP — DIKUNCI

Hanya **`EDM Treaty In`** — endorsemen.

`D:\XML\RNM_BRD\` **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang** — dan proyeknya belum punya kode Go, jadi prior art test
**tidak dapat disebut**. Katakan begitu, jangan dikarang.

Keluaran satu berkas:
**`OUTPUT_HASIL_RNM\.scratch\edm-treaty-in\spec-penyimpanan-relasional.md`**

⛔ `grilling-ronde-1.md` tersegel. ⛔ `spec.md` EDM yang sudah ada jangan disunting.
⛔ **`.scratch\nb-treaty-in\spec-penyimpanan-relasional.md` jangan disunting** — ia sumber, bukan
sasaran.

### ⭐⭐ Aturan terpenting ronde ini

**EDM TIDAK MERANCANG ULANG SATU PUN TABEL.** Sembilan tabel dasar sudah ditetapkan spec NB dan
**bentuknya final**. Yang ditambahkan EDM hanya:

| | |
| --- | --- |
| **kolom** | 21 kolom yang hanya terisi di EDM |
| **baris** | `PRODKE ≥ 1` · `OLD_POLIS_ID` terisi |
| **tabel proyeksi** | `T_POLIS_DIFFERENCE` dan anaknya — **nol baris di NB, hidup di EDM** |

Bila spec EDM menuntut perubahan bentuk tabel dasar, **itu temuan** — tulis di Further Notes,
jangan diam-diam mengubahnya.

---

## 1. SUMBER DAN URUTAN KEWENANGAN

| # | Sumber | Sifat |
| ---: | --- | --- |
| 1 | `.scratch\nb-treaty-in\spec-penyimpanan-relasional.md` | ⭐ **bentuk tabel dasar — final, tidak boleh diubah** |
| 2 | `.scratch\nb-treaty-in\rancangan-tabel-datar-treaty-in.md` Bab **4bis, 4ter, 4quater** | keputusan work owner 22–23 Sep 2026, mengikat |
| 3 | `.scratch\edm-treaty-in\PERTANYAAN-RONDE-1.md` — **P50–P60**, seluruhnya terjawab | keputusan work owner |
| 4 | `.scratch\nb-treaty-in\PERTANYAAN-untuk-*.md` — P1–P49 | berlaku lintas modul |
| 5 | `Diagram-Skema-Tabel-NusantaraRe.xlsx` sheet **`EDM Treaty In Prop`** dan **`EDM Treaty In NonProp`** | bentuk tabel dan relasinya |
| 6 | `.scratch\edm-treaty-in\KEADAAN-EDM-TREATY-IN.md` | keadaan terukur |
| 7 | korpus XML | selalu boleh dipakai membuktikan ulang |

⚠️ Sumber 2 memuat **tiga bab ralat** yang mengutip bunyi lamanya. Yang berlaku selalu yang terbaru.

Bentuk rumah: `.scratch\nb-treaty-in\spec-penyimpanan-relasional.md` *(46 KB, 11 bab)* — spec EDM
mengikuti bentuk yang sama persis.

---

## 2. TABEL YANG DISPEC

### Yang diwarisi — jangan dirancang ulang

```
T_WORK_POLIS · T_GENERAL_POLIS · T_POLIS_QUOTATION · T_POLIS_CEDING
T_POLIS_INSTALMENT · T_POLIS_INSTALMENT_DETAIL · T_POLIS_SPREADING
T_POLIS_XOL · T_POLIS_XOL_LAYER · POOLDATA.HISTORYAKSEPTASIPRODUCTION
```

Spec EDM **menyatakan ulang** apa yang berbeda pada tabel-tabel itu — kolom yang terisi, baris yang
lahir — tanpa mengubah bentuknya.

### Yang ditambahkan — tabel PROYEKSI

```
T_GENERAL_POLIS
 └ T_POLIS_DIFFERENCE                    1:1  ← POLIS_ID, UNIK
    ├ T_POLIS_DIFFERENCE_SPREADING       1:N
    ├ T_POLIS_DIFFERENCE_INSTALMENT      1:N
    └ T_POLIS_XOL_LAYER_DIFFERENCE       1:N  ← hanya NonProp
```

⇒ **EDM Prop 10 tabel · EDM NonProp 14 tabel.**

⭐ `T_POLIS_XOL_DIFFERENCE` **sengaja tidak ada.** `[terverifikasi]` Induk
`TreatyXOLDifferenceList` hanya menyalin tujuh medan kunci yang sudah ada di anaknya; tidak membawa
satu pun informasi baru. Nyatakan ini, jangan diam-diam menambahkannya kembali.

---

## 3. BENTUK SPEC

Bab wajib, berurutan — sama persis dengan spec NB:

```
Cara membaca berkas ini    Problem Statement    Solution
User Stories               Implementation Decisions
Acceptance Criteria        Testing Decisions
Out of Scope               Butir [terbuka] — daftar penuh
Further Notes              TELEMETRI EKSEKUSI
```

### Blok ringkasan

Cacah: user story · acceptance criteria · implementation decision · sebaran penanda · butir
`[terbuka]` aktif · butir `[penyimpangan sadar]`.

⚠️ **Dihitung dua cara**, jendela hitung disebutkan, dan **bila berselisih yang berlaku cara pola
teks** — angka cacah manual dikutip, bukan dihapus (CLAUDE.md §4a).

⭐ Pada spec NB, cacah manual **meleset pada lima dari enam penanda** — `[terverifikasi]` ditaksir
47 padahal 82. Jangan ulangi: hitung dengan pola teks **sebelum** menulis ringkasannya.

### Bentuk Acceptance Criteria

```
N. `[terverifikasi]` <pernyataan>. Test yang menemukan <keadaan sebaliknya> **gagal**. *(ID-n)*
```

Setiap butir membawa penanda dan rujukan ke Implementation Decision. Bab AC **tidak memutuskan apa
pun**.

---

## 4. KEPUTUSAN MENGIKAT — masing-masing wajib jadi satu AC atau lebih

### 4.1 Generasi dan rantai endorsemen

| # | Ketetapan |
| ---: | --- |
| 1 | `PRODKE ≥ 1` pada EDM · `OLD_POLIS_ID` **terisi**, menunjuk `T_WORK_POLIS.ID` generasi sebelumnya |
| 2 | ⭐ `UNIQUE (OLD_POLIS_ID)` **melarang percabangan** — satu generasi hanya boleh punya satu penerus. Tanpa ini rantai menjadi pohon dan selisih mana yang benar tidak terjawab |
| 3 | `UNIQUE (NOPOLIS, PRODKE)` — ⚠️ memperbaiki bahaya nyata: `EDM\RDBList\SelectProdKe` mengambil `MAX + 1` **tanpa penguncian**, sehingga dua endorsemen serentak membaca angka sama. Di sistem baru yang kedua **gagal**, bukan bentrok diam-diam |
| 4 | `EDM_NO = NOPOLIS + "/E" + PRODKE dua digit` — `Activity\SetEDMTNoPolis` *(P55)*. ⛔ `GenerateNoEDMTreaty` **aturan mati**, syaratnya tak pernah benar |
| 5 | Baris generasi lampau **tidak boleh disunting** — itulah pembekuan `OldData` *(P58)* |
| 6 | ⭐ **`OldData` TIDAK menjadi tabel.** Ia baris yang ditunjuk `OLD_POLIS_ID`. `[terverifikasi]` rujukan `OldData` di dalam `OldData` = **0** dari 94 jalur — susunan bersarang itu dibuat `CreateEDMT` langkah 10 dan tidak pernah dibaca |
| 7 | Endorsemen berlapis *(P57)* = senarai berantai. Selisih dihitung terhadap **generasi tepat sebelumnya** |
| 8 | ⭐ **Aturan keutuhan:** generasi `n+1` wajib memuat **setiap `NOURUT`** yang ada di generasi `n`. Yang hilang = endorsemen tidak sah, ditolak di `services` |

### 4.2 `NOURUT` di EDM

| # | Ketetapan |
| ---: | --- |
| 9 | ⭐ **Di EDM baris TIDAK BISA DIHAPUS** *(keputusan work owner)*. `NOURUT` lama terbawa apa adanya |
| 10 | Baris baru ditambahkan di belakang: `NOURUT = maksimum + 1` |
| 11 | ⚠️ Karena tidak ada penghapusan, `PASANGAN_BERGESER = 1` berarti **anomali sungguhan**, bukan derau yang wajar |
| 12 | `[terverifikasi]` Pembatalan *(P56, `SetEDMTCancel`)* **menolkan seluruh kolom uang** ⇒ generasi baru berisi nol, **bukan penghapusan baris** |

### 4.3 Dua puluh satu kolom yang hanya terisi di EDM

`[terverifikasi]` Enam belas tidak pernah dirujuk NB; lima lagi dirujuk tetapi **tidak pernah
ditulis** NB.

```
EDM_NO · EDM_TYPE · PROD_KE · QUARTAL · YEAR_OF_QUARTAL · STATEMENT_TYPE
ID_NEW_BISNIS · GROSS_PREMIUM · PREMI_ONP · RI_COMM_ONP · RESULT_ONP1 · RESULT_ONP2
OVERIDDING_COMM_ONP · RESULT_OGP1 · RESULT_OGP2 · OVERIDDING_COMM_OGP · OUTSTANDING_CLAIM
SALVAGE_VALUE · EXCESS_LOSS          ← NB merujuk 8 kali, tak pernah menulis
```

Di tabel anak: `T_POLIS_QUOTATION.OLD_POLICY_NO` dan
`T_POLIS_SPREADING.SPLIT_RNM_SHARE_PCT` — keduanya hanya EDM yang menulis.

⭐ Dan **empat kolom hanya dipakai NB**: `IS_EDM_INPUT_ON_NB` · `IS_APPROVEDTO_DEPT_HEAD` ·
`HAS_FAC_OUT` · `SHARE_CURRENCY`. Sebutkan juga, supaya tidak dikira kolom mati.

⚠️ `ViewState` ditulis EDM lewat `DataTransform` — **keadaan layar, TIDAK dimigrasi.**

### 4.4 Tabel proyeksi — bukan tabel sumber

| # | Ketetapan |
| ---: | --- |
| 13 | ⭐ Rumus selisih dihitung **di Go, lapisan `services`** — Oracle tidak tahu apa itu selisih. `repository` hanya mengambil **dua baris**: baris ini dan baris `OLD_POLIS_ID` |
| 14 | ⛔ **`V_POLIS_DIFFERENCE` DIBATALKAN.** View berarti rumus ditulis **dua kali**, di Go dan di SQL, dan pasti bercabang. Nyatakan pembatalannya, jangan diam-diam menghidupkannya |
| 15 | Tabel proyeksi dibangun **karena ada pembaca SQL langsung** — work owner memastikan angka selisih dibaca lewat layar aplikasi **dan** lewat SQL sendiri |
| 16 | Kolom `SUMBER`: `'PEGA'` = hasil migrasi, **beku** · `'GO'` = dihitung sistem baru |
| 17 | **Tiga aturan yang membuatnya bukan tabel sumber:** ① hanya ditulis Go, dalam transaksi yang sama dengan generasinya *(P2)* · ② boleh dihapus total dan dibangun ulang, **hanya baris `'GO'`** · ③ bila isinya beda dari hasil hitung ulang, **tabelnya yang salah**, bukan operannya |
| 18 | Sertakan kunci yang dipakai pembaca SQL menyaring: `NOPOLIS` `PRODKE` `EDM_NO` `IDPEGA` — supaya tidak perlu join balik. Gaya penamaan mengikuti `TREATYINPRODUCTION` |

### 4.5 Rumus selisih

| # | Ketetapan |
| ---: | --- |
| 19 | ⭐⭐ **SATU RUMUS UNTUK SEMUA:** `selisih.X = baris_ini.X − baris(OLD_POLIS_ID).X` — tanpa memandang endorsemen pertama atau berlapis, prop atau nonprop |
| 20 | Empat aturan turunannya `[terverifikasi]` dari `EDMTCalculateTreatyDifference`: **uang** dikurangi · **persen** disalin, TIDAK dikurangi · **kunci** disalin · **`DUE_TO`** diturunkan dari **TANDA** selisih |
| 21 | `[penyimpangan sadar]` Sistem lama punya **DUA** rumus di sisi prop, dipilih syarat `.OldData.EDMNo==""`: endorsemen pertama `BARU.X − LAMA.X`, berlapis `BARU.X − LAMA.TreatyDifference.X`. Akibatnya `180 − 50 = 130` padahal yang benar `180 − 150 = 30`. ⛔ **Varian kedua TIDAK ditiru** — keputusan work owner 23-09-2026. **Kutip bunyi aslinya** |
| 22 | ⭐ Sisi nonprop *(`CalculateDifferenceEDM_act`)* memang **sudah seragam** — kesepuluh pengurangannya selalu terhadap **nilai** generasi lampau. Keputusan 19 menyamakan prop dengan nonprop, bukan mengubah keduanya |
| 23 | `DUE_TO` pada selisih XOL: `local.duetovalue` dijumlahkan sepanjang `ValueList` lalu `@If(> 0, "DUE TO US", "DUE TO YOU")` |

### 4.6 Migrasi

| # | Ketetapan |
| ---: | --- |
| 24 | ⭐ **Nilai lama TIDAK PERNAH dihitung ulang** — disalin apa adanya dari `DATA_JSON`, lengkap dengan galat presisinya, supaya laporan lama dapat direkonsiliasi |
| 25 | Pemasangan lewat **`NOURUT`**, bukan kunci dagang. Kunci dagang turun pangkat jadi **pemeriksa** |
| 26 | `PASANGAN_BERGESER = 1` bila kunci dagang `baru[n] ≠ lama[n]`. Dihitung saat migrasi **tanpa menyentuh satu pun angka** |
| 27 | `RUMUS_BERLAPIS = 1` bila nilainya lahir dari rumus lama `baru − selisih lama`. Dapat ditentukan pasti: **setiap baris selisih yang generasi sebelumnya punya `EDMNo` terisi** |
| 28 | Kedua penanda **hanya berarti untuk `SUMBER = 'PEGA'`** |
| 29 | Bangun ulang proyeksi **hanya menyentuh baris `'GO'`**. Tanpa penjaga ini, satu perintah rebuild menimpa seluruh angka historis dan tidak bisa dikembalikan |

### 4.7 Yang berbeda dari NB pada tabel yang sama

| # | Ketetapan |
| ---: | --- |
| 30 | ⭐ **P60** — `T_POLIS_SPREADING`: EDM menetapkan baris pertama `SPLIT_RNM_SHARE_PCT = 100` lalu membagi dengan presisi **20**; NB membagi `100 / jumlah baris` presisi **10**. `CountSpreading_Act` EDM 8 langkah lawan 7 di NB. **Ditiru apa adanya**, dan perbedaannya wajib punya test tersendiri |
| 31 | `T_POLIS_INSTALMENT_DETAIL` **hidup di EDM nonprop**. ⭐ Dan di EDM sarang itu juga datang dari **salinan master kontrak**: `EDMChooseBusiness_Act → FillMasterInstallment → SetInstallmentValue` menyalin dari `pyWorkPage.TreatyIn.Installment(n).InstallmentList` |
| 32 | `T_POLIS_QUOTATION.EDM_TYPE` menentukan **jenis endorsemen**, termasuk pembatalan *(P56)*, dipilih di awal |
| 33 | ⭐ **`TreatyDifference.ListInstallment` hanya SATU tingkat** — tidak punya `InstallmentList`, walau di sisi sumber nonprop ada. ⇒ tidak ada `T_POLIS_DIFFERENCE_INSTALMENT_DETAIL` |

### 4.8 Ketetapan lama yang tetap berlaku

`POOLDATA.` eksplisit *(P3)* · satu transaksi *(P2)* · `OPERATORID` dari identitas login dan `PIC`
dari nama tampilan *(P4, P33)* · uang bukan `float` *(ADR-0003)* · skala minimal 9 desimal, galat
lama diikuti apa adanya *(P29)* · kode dan penanda **tetap teks** · `""` → `NULL` ·
`handlers → services → repository` *(CLAUDE.md §5)*.

---

## 5. YANG WAJIB MASUK OUT OF SCOPE

| Yang dikeluarkan | Sebab |
| --- | --- |
| Sembilan tabel dasar | ⭐ **sudah dispec NB** — dinyatakan ulang perbedaannya saja, bukan dirancang ulang |
| `T_POLIS_XOL_DIFFERENCE` | induk `TreatyXOLDifferenceList` hanya salinan kunci anaknya |
| `V_POLIS_DIFFERENCE` | ⛔ **dibatalkan** — rumus di lapisan aplikasi, bukan di Oracle |
| Jalur konversi **FacOut** | EDM tidak punya FacOut, tidak pernah berjalan *(P50)* |
| Langkah 18 `Call ASMForceCaseClose` | tetap mati *(P52)* |
| `Activity\GenerateNoEDMTreaty` | aturan mati, syaratnya tak pernah benar |
| `OldData` di dalam `OldData` | nol rujukan di korpus |
| `M_TREATY_IN.JSONDATA` | kontrak, di luar lingkup — ⚠️ **tetapi `SetInstallmentValue` menyalin darinya**, jadi ketergantungannya wajib disebut |
| `ViewState` | keadaan layar |
| Perhitungan premi, komisi, pajak | bukan persoalan penyimpanan |
| DDL dan presisi fisik | dicocokkan DBA di dalam tiket |

⚠️ **Peringatan yang wajib dibawa dari spec NB, bukan diulang penuh:** `TREATYINPRODUCTION.DEDUCTION1`
menampung **dua satuan** · `LAYER*` pada polis proporsional bernilai `"0"` bukan kosong · ada
**empat titik sisip** dengan jalur XOL menyisip **per mata uang**. Rujuk ke spec NB, jangan salin
seluruhnya.

---

## 6. BUTIR `[terbuka]` — BAWA UTUH, JANGAN DITUTUP

### Diwarisi dari spec NB

| Butir | Pemilik |
| --- | --- |
| `JSON_DATAGUIDE(DATA_JSON)` dari DBA — menutup selisih **74 lawan 79** medan | `[data DBA]` |
| Presisi fisik kolom uang | `[data DBA]` |
| ⚠️ **Treaty Out masuk lingkup atau tidak** — folder EDM juga memuat `InsertToTreatyOutXOLList` dan `InsertToTreatyOutXOLListEDMOldData` | `[work owner]` |
| Migrasi dokumen lama — dipindahkan atau dibaca lewat jalur lama | `[work owner]` |

### Khas EDM

| Butir | Pemilik |
| --- | --- |
| Maksud dagang varian rumus berlapis — **kenapa** sistem lama mengurangi terhadap selisih, bukan nilai | `[work owner]` |
| Anak `T_POLIS_DIFFERENCE` untuk pembaca SQL — apakah `_SPREADING` dan `_INSTALMENT` memang dibutuhkan, atau cukup induknya | `[work owner]` |
| Sesudah **pembatalan** *(P56)*, apakah polis masih boleh di-endorse lagi | `[work owner]` |
| Batas berapa kali satu polis boleh di-endorse *(P57)* | `[work owner]` |

⛔ **Jangan menutup satu pun.**

---

## 7. DISIPLIN

Setiap pernyataan membawa bukti: **jalur berkas + tipe rule + nama rule**.

Penanda: `[terverifikasi]` · `[dugaan]` · `[terbuka]` · `[data DBA]` · `[keputusan work owner]` ·
`[penyimpangan sadar]`.

⛔ Jangan menulis `CREATE TABLE` atau DDL. ⛔ Jangan membuat ADR baru. ⛔ Nol nama orang · nomor
polis tidak ditulis apa adanya · nol contoh JSON tersimpan · nol data nasabah.

### Seam — satu, sama dengan spec NB

`repository`. Pemuat migrasi menulis lewat antarmuka yang **sama**, bukan jalur tersendiri.

⚠️ **Pulang-pergi saja tidak cukup** — kalau tulis dan baca sama-sama salah simetris, ujinya tetap
hijau. Sebagian test memeriksa **nilai kolom langsung**.

⭐ Khas EDM, wajib punya test tersendiri: rantai `OLD_POLIS_ID` berlapis tiga generasi · penolakan
percabangan · penolakan generasi yang kehilangan `NOURUT` · pembatalan menghasilkan generasi nol
bukan penghapusan · bangun ulang proyeksi tidak menyentuh baris `'PEGA'`.

⛔ **Prior art test tidak dapat disebut** — repo implementasi terlarang dan kode Go belum ada.
Katakan begitu. Yang boleh dirujuk hanya preseden dokumen.

### Empat jebakan yang sudah menjerat proyek ini

1. **Menebak nama tag.** `Activity` memakai `PropertiesName` **tanpa** awalan `py`; `DataTransform`
   dan `Flow` memakai **dengan** awalan. Lima property EDM sempat terlewat karena ini.
2. **Sapuan berjangkar tidak melihat rujukan relatif** di dalam kalang. Terbukti tiga kali.
3. **Mengambil nilai mayoritas korpus, bukan nilai di aturan yang bersangkutan** — sekali ini
   menghasilkan laporan palsu bahwa dua kolom tertukar.
4. **`pyStepsPreCondParams` adalah percabangan lompat**, bukan gerbang hidup/mati.

⇒ **Setiap cacah medan adalah BATAS BAWAH, bukan total.** Nyatakan di spec.

Jalur korpus memuat spasi — pakai Python. ⚠️ Buang `pyExpressionGadget` lebih dulu.
⚠️ Jangan `html.unescape` sebelum mencocokkan pola struktur.

---

## 8. BAB WAJIB — TELEMETRI EKSEKUSI

Ukur dari luar:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Bila tidak dilakukan, **katakan begitu** dan pisahkan tegas mana yang **terukur** dan mana yang
**taksiran** — jangan menyajikan taksiran sebagai angka terukur.

---

## 9. YANG MENANDAKAN RONDE INI BERHASIL

- Spec lengkap sebelas bab, seluruh AC berpenanda dan berujuk ke Implementation Decision.
- **Ketiga puluh tiga ketetapan Bab 4** seluruhnya terwakili di Acceptance Criteria.
- ⭐ **Nol tabel dasar dirancang ulang.** Bila spec menuntutnya, ditulis di Further Notes.
- Pembatalan `V_POLIS_DIFFERENCE` dinyatakan, bukan didiamkan.
- `[penyimpangan sadar]` rumus berlapis tertulis **beserta bunyi aslinya** dan contoh angkanya.
- Ketiga aturan tabel proyeksi tertulis utuh, terutama **bangun ulang hanya baris `'GO'`**.
- Delapan butir `[terbuka]` dibawa utuh — empat warisan, empat khas EDM. Nol yang ditutup sendiri.
- Ringkasan dihitung **dengan pola teks lebih dulu**, bukan cacah manual yang diperiksa belakangan.
- Nol `CREATE TABLE`, nol ADR baru, nol nama orang.
- Bab telemetri memisahkan tegas yang terukur dari yang ditaksir.

---

*Disusun 23 September 2026, sesudah spec penyimpanan NB Treaty In selesai dan dua belas butir
penahan ditutup work owner.*
