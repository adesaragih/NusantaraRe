# Spesifikasi penyimpanan relasional — EDM Treaty In

> ⛔ **~~`T_POLIS_BREAKDOWN_SPREAD`~~ ⛔ DIBATALKAN 23-09-2026.** `[keputusan work owner]` Tabel itu **tidak ada**. Keempat medannya turunan: dua dari master `POOLDATA.PROPORTIONALARRG`, dua dihitung dari `TotalPremium` dan `TotalClaim` yang sudah tersimpan. Rujukan di bawah dicoret, bunyinya tidak dihapus.

---


## Migrasi Treaty Inward — Endorsemen · lapisan penyimpanan

---

> ## ⛔⛔ RALAT 23 September 2026 — sesudah data guide dan satu dokumen nyata
>
> Spec ini ditulis **sebelum** data guide dan satu dokumen EDM sungguhan diterima. Lima hal
> berubah. Yang berlaku adalah **`rancangan-tabel-datar-treaty-in.md` Bab 4quinque**, bukan yang
> tertulis di bawah.
>
> | # | Yang tertulis di spec ini | Yang berlaku |
> | ---: | --- | --- |
> | 1 | skala kolom uang **minimal 9 desimal** | ⭐ **`NUMBER(38,8)`** *(semula ~~`NUMBER(20,8)`~~ — dinaikkan 23-09-2026 sore; **desimalnya tetap delapan**, yang berubah hanya sisi kiri koma)* |
> | 2 | AC uang "tersimpan tanpa kehilangan satu digit pun" | utuh **sampai 8 desimal** |
> | 3 | migrasi "tidak mengubah satu pun nilai uang" | **membulatkan** yang berdesimal > 8 |
> | 4 | — | ⭐ **tabel baru ~~`T_POLIS_BREAKDOWN_SPREAD`~~ ⛔** ← `BreakDownSpreadList`, 4 medan, ada di tingkat polis dan `OldData` |
> | 5 | rumus selisih `[dugaan]` | ⭐ **`[terverifikasi]` dari data produksi** — `baru − lama.nilai`, terbukti pada 6 medan |
>
> ⭐ Dan satu medan baru yang belum ada di spec ini: **`Remark`** *(panjang 128)*.
> `Show` dan `ViewState` tetap tidak dimigrasi.
>
> ⚠️ **Data guide berkedudukan PELENGKAP, bukan sumber kebenaran** — satu dokumen tunggal memuat
> **95 jalur** yang tidak ada di dalamnya. Ia sah untuk `o:length`; ⛔ tidak sah untuk membuktikan
> ketiadaan sebuah medan.

---


## 1 · Cara membaca berkas ini

### 1.1 ⭐⭐ Aturan terpenting berkas ini

**EDM tidak merancang ulang satu pun tabel.** Sembilan tabel dasar sudah ditetapkan
`..\nb-treaty-in\spec-penyimpanan-relasional.md`, dan **bentuknya final**. Berkas ini hanya
menyatakan **apa yang berbeda**:

| | Yang ditambahkan EDM |
| --- | --- |
| **kolom** | **21** kolom yang hanya terisi di EDM |
| **baris** | `PRODKE ≥ 1` · `OLD_POLIS_ID` **terisi** |
| **tabel proyeksi** | `T_POLIS_DIFFERENCE` dan anaknya — ⭐ **nol baris di NB, hidup di EDM** |

**ID-6b** ⛔ **DICABUT 23-09-2026.** `[keputusan work owner]` Tabel sebaran tambahan **tidak ada**. Bunyi lamanya dikutip di blok kepala berkas ini, tidak dihapus.

⭐ Ada di tingkat polis **dan** di dalam `OldData` — jadi ia ikut aturan generasi seperti tabel anak lainnya. ⚠️ Sensus korpus **melewatkannya**: tidak satu pun aturan merujuk anggotanya. ⛔ `[terbuka]` Apa bedanya secara dagang dari `T_POLIS_SPREADING` — spreading punya `CURRENCY` dan `CURRENCY_ID`, breakdown tidak.

⇒ ~~**EDM Prop 11 tabel · EDM NonProp 15 tabel.**~~ ⛔ ⭐ **EDM Prop 10 tabel · EDM NonProp 14 tabel** — kembali ke cacah semula sesudah tabel sebaran tambahan dibatalkan 23-09-2026 sore.

⛔ **Nol tabel dasar diubah bentuknya.** Bila spec ini menuntutnya, itu **temuan** dan ditulis di
Bab 10 — bukan diubah diam-diam. `[terverifikasi]` Bab 10.1: **tidak ada tuntutan seperti itu.**

### 1.2 Ringkasan — dihitung POLA TEKS lebih dulu

> ⚠️ **Jendela hitung: seluruh berkas ini MINUS blok ringkasan 1.2 sendiri** — **812 baris** ⚠️ *(RALAT — semula **766**, diukur sebelum Bab 10.7 dan telemetri ditambahkan; angka lama tidak dihapus)*.
>
> ⭐ **Cacah dilakukan dengan pola teks SEBELUM blok ini ditulis**, bukan cacah manual yang
> diperiksa belakangan. ⚠️ Pada spec NB, cacah manual meleset pada **lima dari enam** penanda —
> `[terverifikasi]` ditaksir 47 padahal 82. Kekeliruan itu tidak diulang di sini.
>
> | Ukuran | Cara A *(pola teks — berlaku)* | Cara B *(keberurutan nomor)* | |
> | --- | ---: | ---: | :---: |
> | User story | **31** | **31** | ✅ |
> | Acceptance criteria | **58** | **58** | ✅ | *(~~55~~ · +3 pada koreksi 23-09)*
> | Implementation decision *(`ID-n`)* | **49** | **49** | ✅ | *(~~47~~ · ~~+ID-6b~~ ⛔ dicabut, +ID-27b, ID-27c)*
> | ⛔ AC tanpa penanda | **0** | — | ✅ |
> | ⛔ AC tanpa rujukan `ID-n` | **0** | — | ✅ |
> | Butir `[terbuka]` aktif *(Bab 9)* | **9** | — | ✅ |
> | Butir `[penyimpangan sadar]` | **4** | — | ✅ |
>
> **Sebaran penanda** — kemunculan label di dalam jendela: `[keputusan work owner]` **46** · `[terverifikasi]` **38** · `[terbuka]` **6** · `[data DBA]` **3** · `[penyimpangan sadar]` **2** · `[dugaan]` **1**

### 1.3 Penanda

| Penanda | Arti |
| --- | --- |
| `[terverifikasi]` | ada **jalur berkas + tipe rule + nama rule**, dan perintah ujinya |
| `[keputusan work owner]` | diputuskan yang berwenang — ⛔ **tidak diubah spec ini** |
| `[penyimpangan sadar]` | ⚠️ **sengaja berbeda dari Pega**, alasannya tertulis |
| `[terbuka]` | belum terjawab — ⛔ **tidak ditutup spec ini** |
| `[data DBA]` | hanya dapat dijawab dari basis data |
| `[dugaan]` | ⛔ **tidak dipakai sebagai dasar keputusan mana pun** |

### 1.4 ⚠️ Peringatan yang berlaku untuk seluruh berkas

⚠️⚠️ **Setiap cacah medan di berkas ini adalah BATAS BAWAH, bukan total.** Sensus berasal dari
rujukan di dalam aturan; medan yang ada di dokumen tetapi tidak pernah dirujuk satu aturan pun
**tidak tertangkap cara mana pun**. Empat jebakan yang sudah menjerat proyek ini ada di Bab 10.5.

⚠️ **Peringatan dari spec NB tetap berlaku dan tidak diulang penuh di sini** —
`TREATYINPRODUCTION.DEDUCTION1` menampung **dua satuan**; `LAYER*` pada polis proporsional bernilai
`"0"` bukan kosong; ada **empat titik sisip** dengan jalur XOL menyisip **per mata uang**.
⇒ Rujuk `..\nb-treaty-in\spec-penyimpanan-relasional.md`.

---

## 2 · Problem Statement

Endorsemen mengubah polis yang sudah terbit. Di sistem lama, perubahan itu disimpan sebagai
**dokumen baru** yang membawa salinan beku nilai lama di dalam dirinya, ditambah **halaman selisih**
yang dihitung ulang setiap kali layar dibuka.

Memindahkannya ke tabel relasional menghadapkan empat persoalan yang tidak ada di polis baru:

1. ⛔ **Rantai generasi dapat bercabang tanpa ada yang menyadarinya.** `[terverifikasi]`
   `RDBList\SelectProdKe` membaca nomor generasi dengan
   `select PRODKE … order by prodke desc`, lalu pemanggilnya menambah satu. ⛔ **Tanpa penguncian.**
   Dua endorsemen serentak atas polis yang sama membaca angka yang sama, dan keduanya tersimpan.
2. ⛔ **Selisih dihitung dua cara yang berbeda di sisi proporsional**, dan salah satunya
   mengurangi terhadap **selisih** generasi lampau, bukan terhadap **nilainya**. Akibat aritmetiknya
   besar dan tidak menimbulkan pesan galat.
3. ⛔ **Nilai lama membawa galat presisi** yang sudah permanen di produksi. Endorsemen berlapis
   menumpuknya pada tiap putaran.
4. ⚠️ **Halaman selisih dibaca bukan hanya oleh layar**, tetapi juga oleh SQL langsung. Menjadikan
   selisih semata hasil hitung di aplikasi akan memutus pembaca itu.

Dan satu yang menguntungkan: ⭐ **bentuk tabelnya sudah ada.** Sembilan tabel dasar ditetapkan spec
NB. Yang tersisa adalah menyatakan perbedaannya dengan tepat.

---

## 3 · Solution

Endorsemen disimpan sebagai **generasi** di dalam tabel yang sama dengan polis baru. Tidak ada tabel
"endorsemen"; yang ada baris dengan `PRODKE ≥ 1` dan `OLD_POLIS_ID` yang menunjuk generasi
sebelumnya.

| Persoalan | Cara menyelesaikannya |
| --- | --- |
| rantai bercabang | `UNIQUE (OLD_POLIS_ID)` — satu generasi hanya boleh punya **satu** penerus |
| nomor generasi bentrok | `UNIQUE (NOPOLIS, PRODKE)` — yang kedua **gagal**, bukan bentrok diam-diam |
| dua rumus selisih | ⭐ **satu rumus untuk semua**: `selisih.X = baris_ini.X − baris(OLD_POLIS_ID).X` |
| pembaca SQL langsung | **tabel proyeksi** `T_POLIS_DIFFERENCE`, ditulis Go, dibaca siapa saja |
| angka lama bergalat | ⭐ **disalin apa adanya**, tidak dihitung ulang — laporan lama tetap dapat direkonsiliasi |
| nilai lama harus beku | baris generasi lampau **tidak boleh disunting** |

⭐ **`OldData` tidak menjadi tabel.** Ia baris yang ditunjuk `OLD_POLIS_ID`.

---

## 4 · User Stories

### Generasi dan rantai

1. Sebagai **admin reasuransi**, saya ingin endorsemen tersimpan sebagai generasi baru dari polis
   yang sama, supaya riwayat perubahan terbaca dari satu rantai, bukan dari dokumen yang tersebar.
2. Sebagai **admin reasuransi**, saya ingin nomor generasi tidak pernah terpakai dua kali untuk
   polis yang sama, supaya dua endorsemen serentak tidak menimpa pekerjaan satu sama lain.
3. Sebagai **admin reasuransi**, saya ingin endorsemen kedua yang serentak **ditolak terang-terangan**,
   supaya saya tahu harus mengulang — bukan menemukan angkanya salah berbulan kemudian.
4. Sebagai **auditor internal**, saya ingin satu generasi hanya punya satu penerus, supaya
   pertanyaan *"selisih mana yang benar"* tidak pernah muncul.
5. Sebagai **auditor internal**, saya ingin nilai generasi lampau tidak dapat disunting, supaya
   selisih yang sudah disetujui tetap dapat dipertanggungjawabkan.
6. Sebagai **underwriter**, saya ingin endorsemen berlapis menghitung selisih terhadap generasi
   tepat sebelumnya, supaya premi tidak diakui dua kali.
7. Sebagai **pemilik sistem**, saya ingin endorsemen yang kehilangan baris rincian ditolak, supaya
   generasi baru tidak diam-diam menghapus komitmen yang masih berlaku.

### Nomor urut baris

8. Sebagai **underwriter**, saya ingin baris rincian generasi baru dapat dipasangkan dengan baris
   generasi lama, supaya selisih per baris berarti.
9. Sebagai **underwriter**, saya ingin baris baru ditambahkan di belakang tanpa mengganggu nomor
   baris lama, supaya pasangan lama tidak bergeser.
10. Sebagai **pemilik data**, saya ingin tahu bila sebuah pasangan bergeser, supaya anomali terlihat
    dan bukan dianggap derau yang wajar.

### Selisih

11. Sebagai **underwriter**, saya ingin selisih dihitung dengan **satu** rumus, supaya hasilnya
    tidak bergantung pada apakah ini endorsemen pertama atau kelima.
12. Sebagai **Finance**, saya ingin medan persentase **disalin**, bukan dikurangi, supaya angka
    persentase tidak berubah makna.
13. Sebagai **Finance**, saya ingin arah utang-piutang diturunkan dari **tanda** selisih, supaya
    tidak perlu ditetapkan manual.
14. Sebagai **pembaca laporan**, saya ingin membaca angka selisih lewat SQL langsung, supaya saya
    tidak harus membuka layar aplikasi untuk merekonsiliasi.
15. Sebagai **pengembang sistem baru**, saya ingin rumus selisih tinggal di **satu** tempat, supaya
    ia tidak bercabang antara aplikasi dan basis data.
16. Sebagai **pemilik data**, saya ingin dapat membangun ulang tabel selisih, supaya cacat
    perhitungan dapat diperbaiki tanpa menyentuh data sumber.
17. Sebagai **pemilik data**, saya ingin pembangunan ulang itu **tidak menyentuh angka historis**,
    supaya satu perintah tidak menghapus rekonsiliasi bertahun-tahun.

### Pembatalan

18. Sebagai **admin reasuransi**, saya ingin pembatalan tercatat sebagai generasi bernilai nol,
    bukan sebagai penghapusan baris, supaya jejaknya tetap ada.
19. Sebagai **auditor internal**, saya ingin pembatalan terbaca dari jenis endorsemennya, supaya
    tidak perlu menebak dari nilai nol.

### Migrasi

20. Sebagai **Finance**, saya ingin angka lama **disalin apa adanya**, lengkap dengan galat
    presisinya, supaya laporan lama masih cocok.
21. Sebagai **pemilik data**, saya ingin tahu baris mana yang lahir dari rumus lama yang berlapis,
    supaya angka itu dapat diperlakukan berbeda bila kelak diputuskan demikian.
22. Sebagai **pemilik data**, saya ingin penanda migrasi hanya berarti untuk baris hasil migrasi,
    supaya baris baru tidak ikut tertandai.
23. Sebagai **pengembang sistem baru**, saya ingin pemuat migrasi menulis lewat antarmuka yang sama
    dengan jalur biasa, supaya data lama melewati pemeriksaan yang sama.

### Bentuk dan perbedaan dari polis baru

24. Sebagai **pengembang sistem baru**, saya ingin kolom yang hanya dipakai endorsemen tetap kosong
    pada polis baru, supaya tidak ada yang mengira kolomnya mati.
25. Sebagai **underwriter**, saya ingin penyebaran risiko endorsemen dihitung dengan aturannya
    sendiri, supaya hasilnya sama dengan sistem lama.
26. Sebagai **underwriter**, saya ingin rincian angsuran bertingkat tetap tersimpan pada endorsemen
    non-proporsional, supaya cicilan per lapisan tidak hilang.
27. Sebagai **pengembang sistem baru**, saya ingin bentuk selisih angsuran **satu tingkat**, sesuai
    sistem lama, supaya tidak ada tabel yang dibuat tanpa isi.

### Keutuhan dan kepatuhan

28. Sebagai **pemilik sistem**, saya ingin seluruh penyimpanan satu generasi berada dalam **satu
    transaksi**, supaya kegagalan tidak meninggalkan generasi separuh jadi.
29. Sebagai **pemilik data**, saya ingin setiap query menulis skema basis data eksplisit, supaya
    tidak ada tabel terbaca dari skema yang salah.
30. Sebagai **Finance**, saya ingin nilai uang tidak pernah `float`, supaya sistem baru tidak
    menambah galat baru.
31. Sebagai **auditor internal**, saya ingin pelaku setiap perubahan tercatat dari identitas login,
    supaya kolom pelaku tidak kosong seperti di sistem lama.

---

## 5 · Implementation Decisions

### 5.1 Modul, antarmuka, dan seam

**ID-1** `[keputusan work owner]` Arah ketergantungan tetap `handlers → services → repository`
*(CLAUDE.md §5)*. `repository` tidak memanggil `services`.

**ID-2** Satu antarmuka `repository` menjadi **satu-satunya pintu** ke penyimpanan polis treaty —
sama dengan spec NB, **tidak ditambah** untuk endorsemen.

**ID-3** Pemuat data migrasi **menulis lewat antarmuka yang sama**, bukan jalur tersendiri.

**ID-4** ⭐ Rumus selisih dihitung di **`services`**, bukan di Oracle. `repository` hanya mengambil
**dua baris**: baris ini dan baris yang ditunjuk `OLD_POLIS_ID`.

### 5.2 Tabel

**ID-5** ⭐ **Sembilan tabel dasar diwarisi tanpa perubahan bentuk** — `T_WORK_POLIS` ·
`T_GENERAL_POLIS` · `T_POLIS_QUOTATION` · `T_POLIS_CEDING` · `T_POLIS_INSTALMENT` ·
`T_POLIS_INSTALMENT_DETAIL` · `T_POLIS_SPREADING` · `T_POLIS_XOL` · `T_POLIS_XOL_LAYER`, ditambah
`POOLDATA.HISTORYAKSEPTASIPRODUCTION` yang **tidak dibuat ulang**.

**ID-6** Tabel **proyeksi** yang ditambahkan EDM:

```
T_GENERAL_POLIS
 └ T_POLIS_DIFFERENCE                  1:1  ← POLIS_ID, UNIK
    ├ T_POLIS_DIFFERENCE_SPREADING     1:N
    ├ T_POLIS_DIFFERENCE_INSTALMENT    1:N
    └ T_POLIS_XOL_LAYER_DIFFERENCE     1:N  ← hanya NonProp
```

⇒ **EDM Prop 10 tabel · EDM NonProp 14 tabel.**

**ID-7** ⭐ **`T_POLIS_XOL_DIFFERENCE` sengaja tidak ada.** `[terverifikasi]` Induk
`TreatyXOLDifferenceList` dirujuk hanya lewat **delapan medan skalar**, masing-masing **satu kali**
— `LayerType` `Layer` `LayerPartType` `LayerPart` `Currency` `IDCurrency` `DueTo` `NetPremi` —
sementara anaknya `.ValueList` dirujuk **85 kali** dengan **17 medan**. ⇒ Induknya salinan penanda
yang sudah ada di anaknya; ⛔ **tidak membawa informasi baru.**

⚠️ **Ralat kecil atas brief ronde ini:** brief menyebut *"tujuh medan kunci"*. Terukur **delapan**,
dan yang kedelapan — `NetPremi` — **bukan kunci melainkan uang**. ⭐ Kesimpulannya tidak berubah:
nilainya tetap ada di `.ValueList`, sehingga induknya tetap tidak perlu tabel.

**ID-8** ⛔ **`OldData` TIDAK menjadi tabel.** Ia baris yang ditunjuk `OLD_POLIS_ID`.
`[terverifikasi]` Rujukan `OldData` **di dalam** `OldData` = **0**. Susunan bersarang itu dibuat
`Activity\CreateEDMT` langkah 10 dan **tidak pernah dibaca**.

### 5.3 Generasi dan rantai

**ID-9** `[keputusan work owner]` Pada EDM, `PRODKE ≥ 1` dan `OLD_POLIS_ID` **terisi**, menunjuk
`T_WORK_POLIS.ID` generasi sebelumnya. Pada NB, `PRODKE = 0` dan `OLD_POLIS_ID` kosong.

**ID-10** ⭐ `UNIQUE (OLD_POLIS_ID)` **melarang percabangan** — satu generasi hanya boleh punya satu
penerus. ⛔ Tanpanya rantai menjadi pohon, dan pertanyaan *"selisih mana yang benar"* tidak
terjawab.

**ID-11** `UNIQUE (NOPOLIS, PRODKE)`. ⚠️ **Ini memperbaiki bahaya nyata, bukan kehati-hatian
teoretis.** `[terverifikasi]` `RDBList\SelectProdKe` berbunyi:

```
select PRODKE as CARI1 from json_polis
where substr(nopolis,1,24) = {…OldData.PolicyNo} order by prodke desc
```

dan `Activity\SetEDMTNoPolis` langkah 3 menambah satu:
`ProdKe = @toInt(OutputData.pxResults(1).CARI1) + 1`.
⛔ **Tidak ada `FOR UPDATE`, tidak ada penguncian.** Dua endorsemen serentak membaca angka yang
sama. Di sistem baru yang kedua **gagal**, bukan bentrok diam-diam.

⚠️ Dan naskah itu merujuk `json_polis` **tanpa skema** — satu bukti lagi bahwa sistem lama tidak
menaati ketetapan skema eksplisit *(P3)*. Ketetapannya **tetap berlaku untuk sistem baru**.

**ID-12** `[keputusan work owner]` `EDM_NO = NOPOLIS + "/E" + PRODKE dua digit` —
`Activity\SetEDMTNoPolis` *(P55)*.

**ID-13** `[keputusan work owner]` Baris generasi lampau **tidak boleh disunting** — itulah
pembekuan `OldData` *(P58)*.

**ID-14** `[keputusan work owner]` Endorsemen berlapis *(P57)* adalah **senarai berantai**. Selisih
dihitung terhadap **generasi tepat sebelumnya**.

**ID-15** ⭐ **Aturan keutuhan:** generasi `n+1` wajib memuat **setiap `NOURUT`** yang ada di
generasi `n`. Yang hilang berarti endorsemen tidak sah, dan **`services` menolaknya**.

### 5.4 `NOURUT` di EDM

**ID-16** ⭐ `[keputusan work owner]` **Di EDM baris TIDAK BISA DIHAPUS.** `NOURUT` lama terbawa apa
adanya. ⚠️ Ini **berbeda dari NB**, tempat baris boleh dihapus dan `NOURUT` dinomori ulang rapat
`1..n` — aman di sana karena belum ada generasi untuk dipasangkan.

**ID-17** Baris baru ditambahkan di belakang: `NOURUT = maksimum + 1`.

**ID-18** ⚠️ Karena tidak ada penghapusan, `PASANGAN_BERGESER = 1` berarti **anomali sungguhan**,
bukan derau yang wajar.

**ID-19** `[terverifikasi]` Pembatalan *(P56, `Activity\SetEDMTCancel`)* **menolkan seluruh kolom
uang** ⇒ generasi baru berisi nol, ⛔ **bukan penghapusan baris**.

### 5.5 Dua puluh satu kolom yang hanya terisi di EDM

`[terverifikasi]` Enam belas tidak pernah dirujuk NB; lima lagi dirujuk tetapi **tidak pernah
ditulis** NB.

```
EDM_NO · EDM_TYPE · PROD_KE · QUARTAL · YEAR_OF_QUARTAL · STATEMENT_TYPE
ID_NEW_BISNIS · GROSS_PREMIUM · PREMI_ONP · RI_COMM_ONP · RESULT_ONP1 · RESULT_ONP2
OVERIDDING_COMM_ONP · RESULT_OGP1 · RESULT_OGP2 · OVERIDDING_COMM_OGP
OUTSTANDING_CLAIM · SALVAGE_VALUE · EXCESS_LOSS     ← NB merujuk, tak pernah menulis
```

**ID-20** Di tabel anak: `T_POLIS_QUOTATION.OLD_POLICY_NO` dan
`T_POLIS_SPREADING.SPLIT_RNM_SHARE_PCT` — keduanya **hanya EDM yang menulis**.

**ID-21** ⭐ Dan **empat kolom hanya dipakai NB** — `IS_EDM_INPUT_ON_NB` ·
`IS_APPROVEDTO_DEPT_HEAD` · `HAS_FAC_OUT` · `SHARE_CURRENCY`. ⛔ **Bukan kolom mati**; ia kosong di
EDM karena memang bukan urusannya.

**ID-22** ⛔ `ViewState` ditulis EDM lewat `DataTransform` — **keadaan layar, TIDAK dimigrasi.**

**ID-27b** ⭐ **`REMARK`** *(panjang 128)* — medan tingkat atas yang sensus korpus lewatkan dan baru terlihat dari data guide. Dimigrasi. `Show` dan `ViewState` tetap **tidak** dimigrasi.

**ID-27c** ⚠️ **Data guide berkedudukan PELENGKAP, bukan sumber kebenaran.** `[terverifikasi]` Satu dokumen tunggal memuat **95 jalur** yang tidak ada di dalamnya — termasuk `QuotationData.CedingCoList`, seluruh `BreakDownSpreadList`, `Claim`, `ExcessLoss`, `ResultOgp1/2`, `SalvageValue`, `MasterID`. Indeksnya tampaknya basi; DBA diminta menyegarkannya.

⭐ Ia **sah** untuk `o:length` — panjang maksimum per medan, yang menghapus kebutuhan menebak ukuran kolom. ⛔ Ia **tidak sah** untuk membuktikan ketiadaan sebuah medan.

### 5.6 Tabel proyeksi — bukan tabel sumber

**ID-23** ⛔ **`V_POLIS_DIFFERENCE` DIBATALKAN.** View berarti rumus ditulis **dua kali** — di Go
dan di SQL — dan cepat atau lambat keduanya bercabang. ⛔ **Jangan dihidupkan kembali diam-diam.**

**ID-24** Tabel proyeksi dibangun **karena ada pembaca SQL langsung**: `[keputusan work owner]`
angka selisih dibaca lewat layar aplikasi **dan** lewat SQL sendiri.

**ID-25** Kolom `SUMBER`: `'PEGA'` = hasil migrasi, **beku** · `'GO'` = dihitung sistem baru.

**ID-26** ⭐ **Tiga aturan yang membuatnya bukan tabel sumber:**
① ia **hanya ditulis Go**, dalam transaksi yang sama dengan generasinya *(P2)*;
② ia **boleh dihapus total dan dibangun ulang** — ⛔ **hanya baris `'GO'`**;
③ bila isinya berbeda dari hasil hitung ulang, ⭐ **tabelnya yang salah**, bukan operannya.

**ID-27** Kunci penyaring disertakan supaya pembaca SQL tidak perlu join balik: `NOPOLIS` ·
`PRODKE` · `EDM_NO` · `IDPEGA`. Gaya penamaan mengikuti `TREATYINPRODUCTION`.

### 5.7 Rumus selisih

**ID-28** ⭐⭐ `[keputusan work owner]` **SATU RUMUS UNTUK SEMUA:**

```
selisih.X = baris_ini.X − baris(OLD_POLIS_ID).X
```

Tanpa memandang endorsemen pertama atau berlapis, proporsional atau non-proporsional.

⭐⭐ **TERBUKTI DARI DATA PRODUKSI 23-09-2026** — naik dari `[dugaan]` ke `[terverifikasi]`. Satu dokumen EDM sungguhan memperlihatkan rumusnya pada **enam medan sekaligus**:

```
OldData.NetPremium               130.463.146,76
NetPremium (generasi baru)                 0,00
TreatyDifference.NetPremium     −130.463.146,76     = baru − lama     ✓
```

Idem `PremiOgp` −200.712.533,47 · `ResultOgp1` −70.249.386,71 · `Claim` −10.502.027,00 · `ListInstallment.Premium` −119.961.119,76 · `SpreadingRiskList.PremiumSpreaded` −130.463.146,76.

⭐ **Persentase disalin, tidak dikurangi** — juga terbukti: `TotalSharePercentagePremium` bernilai `"100"` di ketiga tempat *(baru, lama, selisih)*.

⭐ Dan penjelasan varian ganjil itu ikut ketemu: `OldData.TreatyDifference` isinya **kosong** — `{pxObjClass, ListInstallment:[], SpreadingRiskList:[]}`. Jadi varian kedua sistem lama memang mengurangi terhadap **nol**, bukan terhadap angka yang bermakna. Itu menjelaskan kenapa aritmetiknya tampak ganjil, dan menguatkan keputusan untuk tidak menirunya.

**ID-29** `[terverifikasi]` Empat aturan turunannya, dari `Activity\EDMTCalculateTreatyDifference`:
**uang** dikurangi · **persen** **disalin, tidak dikurangi** · **kunci** disalin ·
**`DUE_TO`** diturunkan dari **tanda** selisih.

**ID-30** ⚠️ `[penyimpangan sadar]` **1 — varian rumus berlapis TIDAK ditiru.**
`[keputusan work owner]` 23-09-2026.

`[terverifikasi]` Sistem lama memuat **dua** blok rumus di sisi proporsional, dan keterangan
langkahnya berbunyi apa adanya:

> **langkah 1** — *"Value Difference = New value - Old Value"*
> `…TreatyDifference.GrossPremium = …GrossPremium − …OldData.GrossPremium`
>
> **langkah 4** — *"Value Difference = New value - Old Value (for data with .OldData.EDMNo)"*
> `…TreatyDifference.GrossPremium = …GrossPremium − …OldData.TreatyDifference.GrossPremium`

Akibat aritmetiknya, dengan angka contoh: nilai baru **180**, nilai lama **150**, selisih lama
**50** ⇒ varian kedua menghasilkan `180 − 50 = **130**`, padahal yang benar `180 − 150 = **30**`.

⛔ **Varian kedua tidak dibangun.** ⭐ ID-28 berlaku untuk seluruh generasi.

⚠️ **Dan satu hal yang tidak dapat saya buktikan dari ekspor:** `[terbuka]` kedua blok itu
**sama-sama bergerbang `pyStepsPreCondition = "true"`**, dengan parameter percabangan kosong.
⇒ **Pemilih antara keduanya hanya tertulis di keterangan langkah**, bukan di gerbang yang
dijalankan. Mana yang benar-benar berjalan **tidak dapat dinyatakan dari korpus**. ⭐ Ini tidak
menahan: varian kedua tidak ditiru apa pun jawabannya.

**ID-31** ⭐ `[terverifikasi]` Sisi non-proporsional — `Activity\CalculateDifferenceEDM_act` —
memang **sudah seragam**: pengurangannya selalu terhadap **nilai** generasi lampau.
⇒ ID-28 **menyamakan proporsional dengan non-proporsional**, bukan mengubah keduanya.

**ID-32** `[terverifikasi]` `DUE_TO` pada selisih XOL: `local.duetovalue` dijumlahkan sepanjang
`ValueList`, lalu `@If(> 0, "DUE TO US", "DUE TO YOU")`.

### 5.8 Migrasi

**ID-33** ⭐ `[keputusan work owner]` **Nilai lama TIDAK PERNAH dihitung ulang** — disalin apa adanya
dari dokumen lama, **lengkap dengan galat presisinya**, supaya laporan lama dapat direkonsiliasi.

**ID-34** Pemasangan baris antar generasi lewat **`NOURUT`**, bukan kunci dagang. ⭐ Kunci dagang
turun pangkat menjadi **pemeriksa**.

**ID-35** `PASANGAN_BERGESER = 1` bila kunci dagang `baru[n] ≠ lama[n]`. Dihitung saat migrasi
⛔ **tanpa menyentuh satu pun angka**.

**ID-36** `RUMUS_BERLAPIS = 1` bila nilainya lahir dari rumus lama `baru − selisih lama`. Dapat
ditentukan pasti: **setiap baris selisih yang generasi sebelumnya punya `EDMNo` terisi**.

**ID-37** Kedua penanda **hanya berarti untuk `SUMBER = 'PEGA'`**.

**ID-38** ⛔ Bangun ulang proyeksi **hanya menyentuh baris `'GO'`**. ⛔ Tanpa penjaga ini, satu
perintah rebuild menimpa seluruh angka historis dan **tidak dapat dikembalikan**.

### 5.9 Yang berbeda dari NB pada tabel yang sama

**ID-39** ⭐ **P60 — `T_POLIS_SPREADING`.** `[terverifikasi]` EDM menetapkan baris pertama
`SPLIT_RNM_SHARE_PCT = 100` lalu membagi `SplitRNMSharePct / RNMShare` dengan presisi **20**;
NB membagi `100 / jumlah baris` dengan presisi **10**. `Activity\CountSpreading_Act` EDM
**8 langkah** lawan **7** di NB. ⭐ **Ditiru apa adanya**, dan perbedaannya **wajib punya test
tersendiri**.

**ID-40** `T_POLIS_INSTALMENT_DETAIL` **hidup di EDM non-proporsional**. ⭐ Dan di EDM sarang itu
juga datang dari **salinan master kontrak**:
`EDMChooseBusiness_Act → FillMasterInstallment → SetInstallmentValue`, menyalin dari
`pyWorkPage.TreatyIn.Installment(n).InstallmentList`.
⚠️ ⇒ ketergantungan pada dokumen kontrak **tetap ada**, walau dokumen itu di luar lingkup.

**ID-41** `T_POLIS_QUOTATION.EDM_TYPE` menentukan **jenis endorsemen**, termasuk pembatalan *(P56)*,
dipilih di awal saat berkas dibuat.

**ID-42** ⭐ `[terverifikasi]` **`TreatyDifference.ListInstallment` hanya SATU tingkat** — nol
rujukan `…ListInstallment(n).InstallmentList`, **12** rujukan ke medan skalarnya — walau di sisi
sumber non-proporsional sarang itu ada. ⇒ ⛔ **Tidak ada `T_POLIS_DIFFERENCE_INSTALMENT_DETAIL`.**

### 5.10 Ketetapan lama yang tetap berlaku

**ID-43** Setiap query menulis skema `POOLDATA.` **eksplisit** *(P3)*.
**ID-44** Seluruh urutan penyimpanan satu generasi dibungkus **satu transaksi** *(P2)*.
**ID-45** `OPERATORID` dari identitas akses login; `PIC` dari nama tampilan *(P4, P33)*.
**ID-46** Uang **tidak pernah `float`** *(ADR-0003)*; ⭐ **`NUMBER(38,8)`** *(semula ~~`NUMBER(20,8)`~~ — dinaikkan 23-09-2026 sore; **desimalnya tetap delapan**, yang berubah hanya sisi kiri koma)*

⚠️ Bunyi lama dikutip: ~~*"skala minimal sembilan desimal; galat lama diikuti apa adanya"*~~. Nilai berdesimal lebih dari delapan **dibulatkan saat dimuat** — `ShareValue` bersimpan 24 desimal, `PremiumSpreaded` 20 desimal. ⭐ Ekor galat `2,76 × 10⁻⁷` tetap terlihat; ia jatuh di desimal ketujuh. ⛔ `[terbuka]` Dua belas digit di depan koma belum diuji terhadap nilai terbesar yang pernah tersimpan.
**ID-47** Kode dan penanda **tetap teks**; teks kosong `""` yang masuk kolom angka atau tanggal
menjadi **`NULL`**, bukan `0`.

---

## 6 · Acceptance Criteria

⛔⛔ **Bab ini tidak memutuskan apa pun.** Ia menyatakan ulang keputusan yang sudah ada dalam bentuk
yang dapat diuji dari luar.

### Generasi dan rantai

1. `[keputusan work owner]` Baris endorsemen tersimpan dengan `PRODKE ≥ 1`. Test yang menemukan
   `PRODKE = 0` pada endorsemen **gagal**. *(ID-9)*
2. `[keputusan work owner]` `OLD_POLIS_ID` pada endorsemen **terisi** dan menunjuk generasi
   sebelumnya. Test yang menemukannya kosong **gagal**. *(ID-9)*
3. `[keputusan work owner]` Dua baris tidak boleh berbagi `OLD_POLIS_ID` yang sama. Test yang
   berhasil menyimpan penerus kedua **gagal**. *(ID-10)*
4. `[terverifikasi]` Dua endorsemen serentak atas polis yang sama: yang kedua **ditolak**. Test yang
   menemukan keduanya tersimpan **gagal**. *(ID-11)*
5. `[keputusan work owner]` Nomor endorsemen berbentuk `NOPOLIS + "/E" + dua digit`. Test yang
   menghasilkan bentuk lain **gagal**. *(ID-12)*
6. `[keputusan work owner]` Menyunting baris generasi yang sudah punya penerus **ditolak**. Test
   yang berhasil mengubahnya **gagal**. *(ID-13)*
7. `[keputusan work owner]` Selisih dihitung terhadap generasi **tepat sebelumnya**. Test pada
   generasi ketiga yang menemukan selisih terhadap generasi pertama **gagal**. *(ID-14)*
8. `[keputusan work owner]` Generasi `n+1` yang kehilangan salah satu `NOURUT` milik generasi `n`
   **ditolak** di `services`. Test yang berhasil menyimpannya **gagal**. *(ID-15)*
9. `[terverifikasi]` Tidak ada tabel bernama `OldData` atau semacamnya; nilai lama dibaca lewat
   `OLD_POLIS_ID`. Test yang menemukan tabel salinan **gagal**. *(ID-8)*

### `NOURUT`

10. `[keputusan work owner]` Pada endorsemen, baris rincian **tidak dapat dihapus**. Test yang
    berhasil menghapusnya **gagal**. *(ID-16)*
11. `[keputusan work owner]` `NOURUT` baris lama **terbawa apa adanya** ke generasi berikutnya. Test
    yang menemukan penomoran ulang **gagal**. *(ID-16)*
12. `[keputusan work owner]` Baris baru mendapat `NOURUT = maksimum + 1`. Test yang menemukan
    penyisipan di tengah **gagal**. *(ID-17)*
13. `[terverifikasi]` `PASANGAN_BERGESER = 1` diperlakukan sebagai **anomali**, bukan keadaan wajar.
    Test yang mengabaikannya **gagal**. *(ID-18)*

### Pembatalan

14. `[terverifikasi]` Pembatalan menghasilkan generasi baru dengan kolom uang **bernilai nol**. Test
    yang menemukan penghapusan baris **gagal**. *(ID-19)*
15. `[keputusan work owner]` Jenis endorsemen — termasuk pembatalan — tersimpan di kolomnya sendiri
    dan dipilih di awal. Test yang menurunkannya dari nilai nol **gagal**. *(ID-41)*

### Selisih

16. `[keputusan work owner]` Selisih dihitung `baris_ini.X − baris(OLD_POLIS_ID).X` untuk **seluruh**
    generasi. Test yang menemukan rumus kedua **gagal**. *(ID-28, ID-30)*
17. `[terverifikasi]` Medan **uang** pada tabel selisih berisi hasil pengurangan. Test yang
    menemukan salinan **gagal**. *(ID-29)*
18. `[terverifikasi]` Medan **persentase** pada tabel selisih **disalin**, tidak dikurangi. Test yang
    menemukan pengurangan **gagal**. *(ID-29)*
19. `[terverifikasi]` Medan **kunci** pada tabel selisih disalin apa adanya. Test yang menemukan
    perubahan **gagal**. *(ID-29)*
20. `[terverifikasi]` `DUE_TO` diturunkan dari **tanda** selisih, bukan disalin. Test yang menemukan
    salinan **gagal**. *(ID-29)*
21. `[terverifikasi]` `DUE_TO` pada selisih XOL diturunkan dari **jumlah** sepanjang daftar lapisan,
    bukan dari satu lapisan. Test yang menurunkannya dari lapisan pertama **gagal**. *(ID-32)*
22. `[terverifikasi]` Sisi non-proporsional menghasilkan angka yang **sama** sebelum dan sesudah
    penyeragaman rumus. Test yang menemukan perubahan **gagal** — ID-28 tidak mengubah sisi ini.
    *(ID-31)*

### Tabel proyeksi

23. `[keputusan work owner]` Rumus selisih tidak pernah ditulis di dalam basis data. Test yang
    menemukan view atau prosedur yang menghitung selisih **gagal**. *(ID-23)*
24. `[keputusan work owner]` Tabel selisih **hanya ditulis** oleh lapisan aplikasi. Test yang
    menemukan penulisan dari jalur lain **gagal**. *(ID-26)*
25. `[keputusan work owner]` Tabel selisih ditulis **dalam transaksi yang sama** dengan generasinya.
    Test yang menemukan generasi tersimpan tanpa selisihnya **gagal**. *(ID-26, ID-44)*
26. `[keputusan work owner]` Membangun ulang tabel selisih **hanya menyentuh baris `SUMBER = 'GO'`**.
    Test yang menemukan satu baris `'PEGA'` berubah **gagal**. *(ID-26, ID-38)*
27. `[keputusan work owner]` Bila isi tabel selisih berbeda dari hasil hitung ulang, **tabelnya**
    yang dinyatakan salah. Test yang menyimpulkan sebaliknya **gagal**. *(ID-26)*
28. `[keputusan work owner]` Tabel selisih memuat `NOPOLIS` `PRODKE` `EDM_NO` `IDPEGA` sehingga
    pembaca SQL dapat menyaring tanpa join balik. Test yang memerlukan join **gagal**. *(ID-27)*
29. `[terverifikasi]` Tidak ada tabel selisih untuk **induk** XOL. Test yang menemukannya **gagal**.
    *(ID-7)*
30. `[terverifikasi]` Tidak ada tabel selisih untuk rincian angsuran. Test yang menemukannya
    **gagal**. *(ID-42)*
31. `[keputusan work owner]` `repository` mengambil **dua baris** untuk menghitung selisih, bukan
    menjalankan agregasi di basis data. Test yang menemukan agregasi **gagal**. *(ID-4)*

### Kolom yang khas

32. `[terverifikasi]` Dua puluh satu kolom khas endorsemen **kosong** pada baris polis baru. Test
    yang menemukannya terisi **gagal**. *(ID-20)*
33. `[terverifikasi]` Empat kolom khas polis baru **kosong** pada baris endorsemen, dan **tidak**
    dihapus dari skema. Test yang menemukan kolomnya hilang **gagal**. *(ID-21)*
34. `[keputusan work owner]` Keadaan layar **tidak tersimpan** di tabel mana pun. Test yang
    menemukan kolomnya **gagal**. *(ID-22)*

### Perbedaan dari polis baru pada tabel yang sama

35. `[terverifikasi]` Penyebaran risiko endorsemen memakai presisi **20**; polis baru memakai **10**.
    Test yang menyeragamkan keduanya **gagal**. *(ID-39)*
36. `[terverifikasi]` Baris pertama penyebaran pada endorsemen menerima nilai bawaan **100** ketika
    persentase pembagiannya kosong. Test yang menemukan pembagian dengan nol **gagal**. *(ID-39)*
37. `[terverifikasi]` Rincian angsuran bertingkat **tersimpan** pada endorsemen non-proporsional.
    Test yang menemukannya hilang **gagal**. *(ID-40)*
38. `[terverifikasi]` Rincian angsuran yang berasal dari salinan master kontrak tersimpan sama
    seperti yang diketik. Test yang membedakan keduanya **gagal**. *(ID-40)*

### Migrasi

39. `[keputusan work owner]` Nilai hasil migrasi **tidak dihitung ulang** — tersimpan persis seperti
    di dokumen lama, termasuk digit galatnya. Test yang menemukan pembulatan **gagal**. *(ID-33)*
40. `[keputusan work owner]` Pemasangan baris antar generasi memakai `NOURUT`. Test yang memasangkan
    menurut kunci dagang **gagal**. *(ID-34)*
41. `[terverifikasi]` `PASANGAN_BERGESER` dihitung **tanpa mengubah satu pun angka**. Test yang
    menemukan nilai berubah saat penandaan **gagal**. *(ID-35)*
42. `[terverifikasi]` `RUMUS_BERLAPIS = 1` pada setiap baris selisih yang generasi sebelumnya punya
    nomor endorsemen terisi. Test yang menemukannya kosong **gagal**. *(ID-36)*
43. `[keputusan work owner]` Kedua penanda migrasi **hanya** terisi pada baris `SUMBER = 'PEGA'`.
    Test yang menemukannya pada baris `'GO'` **gagal**. *(ID-37)*
44. `[keputusan work owner]` Pemuat migrasi menulis lewat antarmuka `repository` yang sama. Test yang
    menemukan jalur tulis tersendiri **gagal**. *(ID-3)*

### Keutuhan dan kepatuhan

45. `[keputusan work owner]` Seluruh penyimpanan satu generasi berada dalam **satu transaksi**. Test
    yang menemukan generasi separuh tersimpan **gagal**. *(ID-44)*
46. `[keputusan work owner]` Setiap query menulis skema `POOLDATA.` eksplisit. Test yang menemukan
    rujukan tabel tanpa skema **gagal**. *(ID-43)*
47. `[keputusan work owner]` Kolom pelaku diisi dari identitas akses login. Test yang menemukannya
    kosong **gagal**. *(ID-45)*
48. `[keputusan work owner]` Tidak ada kolom uang bertipe `float`. Test yang menemukan satu pun
    **gagal**. *(ID-46)*
49. `[keputusan work owner]` ⭐ Kolom uang menerima **delapan** angka di belakang koma. Nilai berdesimal lebih dari delapan **dibulatkan saat dimuat, bukan ditolak**. Test yang menemukan penolakan pemuatan **gagal**. *(ID-46)*

> ⛔ **DISELARASKAN 23-09-2026 sore.** Bunyi lama dikutip: ~~*"menerima sekurangnya **sembilan** angka di belakang koma; test yang menemukan pemotongan gagal"*~~.
>
> ⚠️ Pertentangan ini **ditemukan pelaksana ronde tiket**, bukan oleh penyusun spec: AC 49 menuntut sembilan sementara **AC 58 pada spec yang sama** menetapkan delapan dan menyatakan kelebihannya dibulatkan. Keduanya merujuk `ID-46`.
>
> ⭐ Akibatnya terukur: nilai produksi berdesimal sembilan **berubah** pada skala delapan — diterima menurut AC 58, ditolak menurut bunyi lama AC 49. Yang berlaku: **AC 58**.
>
> ⛔ Spec NB sudah diselaraskan lebih dulu *(AC 19, 20, 55)*; AC 49 ini yang tertinggal.
50. `[keputusan work owner]` Kode dan penanda tersimpan sebagai **teks**, nol di depan utuh. Test
    yang menemukannya sebagai bilangan **gagal**. *(ID-47)*
51. `[keputusan work owner]` Teks kosong yang masuk kolom angka atau tanggal tersimpan **`NULL`**,
    bukan `0`. Test yang menemukan `0` **gagal**. *(ID-47)*
52. `[keputusan work owner]` Arah ketergantungan `handlers → services → repository` tidak pernah
    dibalik. Test yang menemukan pemanggilan ke arah sebaliknya **gagal**. *(ID-1)*
53. `[keputusan work owner]` Hanya ada **satu** antarmuka `repository` untuk penyimpanan polis
    treaty. Test yang menemukan antarmuka kedua khusus endorsemen **gagal**. *(ID-2)*

### Bentuk tabel

54. `[terverifikasi]` Bentuk kesembilan tabel dasar **tidak berubah** oleh spec ini. Test yang
    menemukan kolom dasar bertambah, hilang, atau berganti tipe **gagal**. *(ID-5)*
55. `[terverifikasi]` ⭐ Endorsemen proporsional memakai **10** tabel; non-proporsional
    **14**. Test yang menemukan jumlah lain **gagal**. *(ID-5, ID-6)*

    > ⚠️ **Angka ini berubah dua kali dalam satu hari, lalu kembali ke semula.**
    > | Waktu | Prop | NonProp | Sebab |
    > | --- | ---: | ---: | --- |
    > | semula | 10 | 14 | — |
    > | 23-09 pagi | ~~11~~ | ~~15~~ | tabel sebaran tambahan dikira tabel dasar baru |
    > | 23-09 sore | **10** | **14** | tabel itu **dibatalkan** — keempat medannya turunan |
    >
    > ⭐ Yang berlaku baris terakhir. Kedua bunyi sebelumnya dikutip, tidak dihapus.
56. ~~`[keputusan work owner]`~~ ⛔ **DITARIK DARI LINGKUP 23-09-2026 sore.** `[keputusan work owner]` Sebaran tambahan **tidak dimigrasi sama sekali** — bukan tabelnya saja yang tidak dibuat, **perhitungannya pun tidak dibawa** ke sistem baru. Nol tabel, nol rumus, nol layar.

> ⛔ **Dua bunyi sebelumnya dikutip, keduanya ditarik.**
> 1. *(pagi)* ~~"`T_POLIS_BREAKDOWN_SPREAD` menyimpan empat medan — `TREATY_TYPE` `SHARE_PERCENTAGE` `PREMIUM_SPREADED` `CLAIM_SPREADED` — dan tidak memiliki `CURRENCY` maupun `CURRENCY_ID`"~~
> 2. *(sore, koreksi pertama)* ~~"Tidak ada tabel yang menyimpan sebaran tambahan. Angkanya **dihitung** — persentase dari master dikali total yang tersimpan"~~
>
> ⚠️ Koreksi pertama itu **masih terlalu jauh**: ia mempertahankan perhitungannya di lapisan `services`. Work owner menyatakan perhitungannya pun tidak perlu dibawa.

⭐ **Cacah AC berlaku turun 58 → 57.** Nomor 56 tidak dipakai ulang, supaya rujukan lama tidak salah arah.
57. `[terverifikasi]` `T_GENERAL_POLIS` menyimpan `REMARK` dengan panjang sekurangnya 128. Test
    yang menemukan medan itu hilang saat pulang-pergi **gagal**. *(ID-27b)*
58. `[keputusan work owner]` Kolom uang bertipe `NUMBER(38,8)` *(semula ~~`NUMBER(20,8)`~~)*. Nilai berdesimal lebih dari delapan
    **dibulatkan, bukan ditolak**. Test yang menemukan kegagalan pemuatan pada `ShareValue`
    berdesimal 24 **gagal**. *(ID-46)*

## 7 · Testing Decisions

### 7.1 Apa yang membuat sebuah test baik di sini

Test menguji **perilaku dari luar**. Test yang pecah ketika sebuah fungsi dipecah dua, padahal
keluarannya sama, adalah test yang salah tulis.

### 7.2 Seam — satu, sama dengan spec NB

**`repository`.** Pemuat migrasi menulis lewat antarmuka yang **sama**, bukan jalur tersendiri.
⛔ Tidak ada seam kedua yang ditambahkan untuk endorsemen.

### 7.3 ⚠️ Pulang-pergi saja tidak cukup

Bila tulis dan baca sama-sama salah secara simetris, uji pulang-pergi tetap hijau. ⇒ **Sebagian
test memeriksa nilai kolom langsung**, bukan hanya hasil baca-ulang.

### 7.4 ⭐ Yang wajib punya test tersendiri di EDM

| # | Yang diuji | Kenapa |
| ---: | --- | --- |
| 1 | rantai `OLD_POLIS_ID` **berlapis tiga generasi** | selisih terhadap generasi tepat sebelumnya hanya terbukti mulai generasi ketiga |
| 2 | **penolakan percabangan** | dua penerus untuk satu generasi |
| 3 | **penolakan generasi yang kehilangan `NOURUT`** | aturan keutuhan |
| 4 | **pembatalan menghasilkan generasi nol**, bukan penghapusan | P56 |
| 5 | **bangun ulang proyeksi tidak menyentuh baris `'PEGA'`** | satu perintah dapat menghapus rekonsiliasi bertahun-tahun |
| 6 | **presisi 20 lawan 10** pada penyebaran | perbedaan sadar dari NB, P60 |
| 7 | **dua endorsemen serentak** | `SelectProdKe` tanpa penguncian |

### 7.5 ⛔ Prior art tidak dapat disebut

⛔ **Repo implementasi terlarang, dan kode Go belum ada.** ⇒ **Tidak ada prior art test yang dapat
dirujuk.** ⚠️ Menyebut satu pun berarti mengarangnya.

⭐ Yang boleh dirujuk hanya **preseden dokumen**: `..\nb-treaty-in\spec-penyimpanan-relasional.md`
Bab 7, dan `..\claim-life\spec-penyimpanan-relasional.md` — keduanya memakai seam `repository` yang
sama.

---

## 8 · Out of Scope

| # | Yang dikeluarkan | Sebab |
| ---: | --- | --- |
| 1 | **Sembilan tabel dasar** | ⭐ sudah dispec NB — dinyatakan ulang perbedaannya, ⛔ bukan dirancang ulang |
| 2 | `T_POLIS_XOL_DIFFERENCE` | induk hanya salinan penanda yang sudah ada di anaknya *(ID-7)* |
| 3 | `V_POLIS_DIFFERENCE` | ⛔ **dibatalkan** — rumus di lapisan aplikasi, bukan di Oracle *(ID-23)* |
| 4 | `T_POLIS_DIFFERENCE_INSTALMENT_DETAIL` | selisih angsuran hanya satu tingkat *(ID-42)* |
| 5 | Jalur konversi **FacOut** | EDM tidak punya FacOut, tidak pernah berjalan *(P50)* |
| 6 | Langkah penutupan paksa berkas | tetap mati *(P52)* |
| 7 | Jalur penomoran kedua | ⚠️ lihat Bab 10.2 — **bukan aturan mati**, tetapi tidak dimigrasi *(P55)* |
| 8 | `OldData` di dalam `OldData` | nol rujukan di korpus *(ID-8)* |
| 9 | Dokumen kontrak sebagai penyimpanan | di luar lingkup — ⚠️ **tetapi ketergantungannya tetap ada** *(ID-40)* |
| 10 | Keadaan layar | bukan data bisnis *(ID-22)* |
| 11 | Perhitungan premi, komisi, pajak | bukan persoalan penyimpanan — ada di `spec.md` EDM |
| 12 | DDL dan presisi fisik | ⛔ dicocokkan DBA di dalam tiket |

⚠️ **Peringatan yang dibawa dari spec NB, tidak diulang penuh:** `TREATYINPRODUCTION.DEDUCTION1`
menampung **dua satuan** — persen dari jalur proporsional, uang dari jalur XOL, sehingga pembaca SQL
**wajib menyaring jenis kontrak lebih dulu**; `LAYER*` pada polis proporsional bernilai `"0"` bukan
kosong; ada **empat titik sisip** dengan jalur XOL menyisip **per mata uang**.

---

## 9 · Butir `[terbuka]` — daftar penuh

⛔⛔ **Spec ini tidak menutup satu pun butir di bawah.**

### 9.1 Diwarisi dari spec NB

| # | Butir | Pemilik |
| ---: | --- | --- |
| 1 | ⭐ **Panduan bentuk dokumen dari DBA** — satu-satunya yang dapat menutup selisih **74 lawan 79** medan | `[data DBA]` |
| 2 | **Presisi fisik kolom uang** — skala pasti | `[data DBA]` |
| 3 | ⚠️ **Treaty Out masuk lingkup atau tidak** — folder EDM juga memuat `InsertToTreatyOutXOLList` dan `InsertToTreatyOutXOLListEDMOldData` | `[work owner]` |
| 4 | **Migrasi dokumen lama** — dipindahkan, atau dibaca lewat jalur lama | `[work owner]` |

### 9.2 Khas EDM

| # | Butir | Pemilik |
| ---: | --- | --- |
| 5 | **Maksud dagang varian rumus berlapis** — kenapa sistem lama mengurangi terhadap **selisih**, bukan terhadap nilai | `[work owner]` |
| 6 | **Anak `T_POLIS_DIFFERENCE` untuk pembaca SQL** — apakah `_SPREADING` dan `_INSTALMENT` memang dibutuhkan, atau cukup induknya | `[work owner]` |
| 7 | Sesudah **pembatalan** *(P56)*, apakah polis masih boleh di-endorse lagi | `[work owner]` |
| 8 | **Batas berapa kali** satu polis boleh di-endorse *(P57)* — batas teknis 99 | `[work owner]` |

### 9.3 ⭐ Satu yang lahir dari ronde ini

| # | Butir | Pemilik |
| ---: | --- | --- |
| 9 | ⚠️ **Pemilih antara dua rumus selisih tidak terbaca dari ekspor.** Kedua blok bergerbang `"true"` dengan parameter percabangan kosong; pemilihnya hanya tertulis di **keterangan langkah**. ⭐ **Tidak menahan** — varian kedua tidak ditiru apa pun jawabannya | `[pengembang Pega lama]` |

---

## 10 · Further Notes

### 10.1 ⭐ Nol tabel dasar dirancang ulang — diperiksa

`[terverifikasi]` Spec ini **tidak menuntut satu pun perubahan bentuk** pada kesembilan tabel dasar.
Yang dinyatakan hanya: kolom mana terisi, baris mana lahir, dan tabel proyeksi mana yang menyusul.

⚠️ **Satu hal yang mendekati tuntutan, dan sengaja tidak dijadikan perubahan:** ID-16 menyatakan
`NOURUT` **tidak boleh dinomori ulang** di EDM, sementara spec NB membolehkannya. ⭐ Itu **aturan
perilaku, bukan bentuk kolom** — kolomnya sama persis. ⇒ **Bukan perubahan bentuk.**

### 10.2 ⚠️ Pertentangan yang dibawa dari ronde sebelumnya

Brief ronde ini menyebut `Activity\GenerateNoEDMTreaty` sebagai **aturan mati yang syaratnya tak
pernah benar**. `[terverifikasi]` Terukur pada ronde `to-spec` EDM: ia bertipe **`RDBList`**, bukan
`Activity`; ⛔ **tidak ditemukan syarat penjaga apa pun**; dan pemanggilnya
`Activity\GeneratePolicyNoTreatyAddendum_Act` dirujuk dari sebuah layar adendum.

⇒ ⛔ **Ia bukan aturan mati, melainkan jalur penomoran kedua yang hidup**, menghasilkan bentuk
`RNM-E…` dengan lima digit. ⭐ **P55 tetap menang** — satu jalur, dan jalur kedua tidak dimigrasi.
⚠️ Pertentangannya tercatat di `spec.md` EDM Bab 10.1 sebagai butir `[terbuka]`.

### 10.3 Penyimpangan sadar

| # | Penyimpangan | Alasan |
| ---: | --- | --- |
| ⚠️ **1** | **Varian rumus selisih berlapis tidak ditiru** *(ID-30)* | menghasilkan angka yang keliru; `[keputusan work owner]` 23-09-2026 |
| ⚠️ **2** | **Nomor generasi bentrok ditolak**, bukan dibiarkan | sistem lama membaca tanpa penguncian dan menyimpan keduanya *(ID-11)* |
| ⚠️ **3** | **Baris rincian tidak dapat dihapus di endorsemen** | pemasangan antar generasi lewat `NOURUT` menuntutnya *(ID-16)* |
| ⚠️ **4** | **Rumus selisih hanya di lapisan aplikasi**, view dibatalkan | rumus di dua tempat pasti bercabang *(ID-23)* |

### 10.4 Tiga tabel proyeksi — kenapa ada sama sekali

Rumus selisih tinggal di aplikasi *(ID-4)*. Tabelnya ada **bukan** karena rumusnya di sana,
melainkan karena `[keputusan work owner]` angka selisih **dibaca lewat SQL langsung**, di luar layar
aplikasi. ⭐ Ketiga aturan ID-26 yang membuatnya **proyeksi, bukan sumber** adalah harga yang dibayar
untuk itu.

### 10.5 ⛔ Empat jebakan yang sudah menjerat proyek ini

1. ⛔ **Menebak nama tag.** `Activity` memakai `PropertiesName`/`PropertiesValue` **tanpa** awalan
   `py`; `DataTransform` dan `Flow` memakai **dengan** awalan. **Lima property EDM sempat terlewat**
   karena ini.
2. ⛔ **Sapuan berjangkar tidak melihat rujukan relatif** di dalam kalang — anggota daftar dirujuk
   `.Suggest` atau `.CedingCoName`, bukan `PolicyTreatyIn.…`. **Terbukti tiga kali.**
3. ⛔ **Mengambil nilai mayoritas korpus**, bukan nilai di aturan yang bersangkutan — sekali ini
   menghasilkan laporan palsu bahwa dua kolom tertukar.
4. ⛔ **`pyStepsPreCondParams` adalah percabangan lompat**, bukan gerbang hidup/mati.
   ⭐ Ronde ini menemukan akibatnya langsung: kedua blok rumus selisih bergerbang `"true"`, sehingga
   pemilihnya **tidak terbaca** — lihat butir `[terbuka]` **9**.

⚠️ Buang blok `pyExpressionGadget` lebih dulu. ⚠️ Jangan `html.unescape` sebelum mencocokkan pola
struktur. ⚠️ Jalur korpus memuat spasi — pakai Python, bukan loop shell.

### 10.6 Tidak ada ADR baru

⭐ Spec ini **tidak menuntut ADR baru**. **ADR-0003** *(uang bukan `float`)* ditegakkan apa adanya;
ID-23 — rumus di aplikasi, bukan di basis data — adalah penerapan arah ketergantungan
`handlers → services → repository` yang sudah ada *(CLAUDE.md §5)*, bukan keputusan arsitektur baru.

---

### 10.7 ⭐ Peta cakupan — 33 ketetapan mengikat ke Acceptance Criteria

⚠️ **Dapat diperiksa langsung.** Tiap ketetapan brief punya sedikitnya satu AC.

| Ket | Isinya | AC |
| ---: | --- | --- |
| 1 | `PRODKE ≥ 1` · `OLD_POLIS_ID` terisi | **1** · **2** |
| 2 | `UNIQUE (OLD_POLIS_ID)` melarang percabangan | **3** |
| 3 | `UNIQUE (NOPOLIS, PRODKE)` · `SelectProdKe` tanpa penguncian | **4** |
| 4 | bentuk `EDM_NO` | **5** · ⚠️ Bab 10.2 |
| 5 | generasi lampau tidak boleh disunting | **6** |
| 6 | `OldData` **tidak** menjadi tabel | **9** |
| 7 | endorsemen berlapis = senarai berantai | **7** |
| 8 | aturan keutuhan `NOURUT` | **8** |
| 9 | di EDM baris **tidak bisa dihapus** | **10** · **11** |
| 10 | baris baru `NOURUT = maksimum + 1` | **12** |
| 11 | `PASANGAN_BERGESER` = anomali sungguhan | **13** |
| 12 | pembatalan menolkan kolom uang | **14** |
| 13 | rumus di `services`; `repository` ambil dua baris | **31** |
| 14 | `V_POLIS_DIFFERENCE` dibatalkan | **23** |
| 15 | proyeksi ada karena pembaca SQL langsung | **28** |
| 16 | kolom `SUMBER` — `'PEGA'` beku, `'GO'` dihitung | **26** · **43** |
| 17 | tiga aturan yang membuatnya bukan tabel sumber | **24** · **25** · **26** · **27** |
| 18 | kunci penyaring untuk pembaca SQL | **28** |
| 19 | ⭐ **satu rumus untuk semua** | **16** |
| 20 | empat aturan turunan: uang · persen · kunci · `DUE_TO` | **17** · **18** · **19** · **20** |
| 21 | ⚠️ varian rumus berlapis **tidak ditiru** | **16** · ID-30 |
| 22 | sisi non-proporsional sudah seragam | **22** |
| 23 | `DUE_TO` XOL dari jumlah sepanjang daftar lapisan | **21** |
| 24 | nilai lama **tidak dihitung ulang** | **39** |
| 25 | pemasangan lewat `NOURUT`, kunci dagang jadi pemeriksa | **40** |
| 26 | `PASANGAN_BERGESER` dihitung tanpa menyentuh angka | **41** |
| 27 | `RUMUS_BERLAPIS` dapat ditentukan pasti | **42** |
| 28 | kedua penanda hanya untuk `SUMBER = 'PEGA'` | **43** |
| 29 | bangun ulang **hanya baris `'GO'`** | **26** |
| 30 | **P60** presisi 20 lawan 10 | **35** · **36** |
| 31 | rincian angsuran hidup di EDM nonprop · salinan master kontrak | **37** · **38** |
| 32 | `EDM_TYPE` menentukan jenis endorsemen | **15** |
| 33 | selisih angsuran **satu tingkat** | **30** |

⭐ **Tiga puluh tiga dari tiga puluh tiga terwakili.** Ditambah AC **32** · **33** · **34**
*(kolom khas)*, **44**–**53** *(ketetapan lama Bab 4.8)*, dan **54** · **55** *(bentuk tabel tidak
berubah)*.

---

## ⭐ KEPUTUSAN 23-09-2026 — empat butir tertahan ditutup

Seluruhnya `[keputusan work owner]`. Bunyi lama dikutip, tidak dihapus.

### 1 · Batas endorsemen — **tanpa batas**

Satu polis boleh di-endorse **tanpa batas**. Batas teknis warisan **99** tidak dipertahankan;
kolom nomor generasi disimpan sebagai **bilangan bulat**, bukan dua digit teks.

Bunyi lama: ~~*"batas teknis 99, batas dagang belum ditetapkan"*~~.

⭐ Diuji ke korpus sebelum dicatat — **117 berkas** menyebut kolom itu, diperiksa dua pola:

| Temuan | Akibat |
| --- | --- |
| kenaikannya **aritmetika** — `Prodke+1`, bukan penyambungan teks | pelebaran kolom tidak mengubah cara menaikkannya |
| nol `LPAD` · nol `RPAD` · nol `TO_CHAR` bertopeng lebar | nol tampilan berpadding nol yang ikut berubah |
| nilai lama 1..99 | muat tanpa pemetaan ulang |

⛔ Go **tidak** memasang penolakan pada angka tertentu. Tidak ada galat "batas endorsemen
tercapai" yang perlu ditulis.

⚠️ Yang **tidak** dibuka keputusan ini: kueri pencari generasi terakhir memotong nomor polis pada
**24 karakter tetap** — `substr(nopolis,1,24)`. Asumsi keras itu tetap `[terbuka]`, dan tidak
berhubungan dengan nomor generasi.

### 2 · Anak tabel proyeksi — **keduanya dibuat**

`T_POLIS_DIFFERENCE_SPREADING` dan `T_POLIS_DIFFERENCE_INSTALMENT` dibuat keduanya. Pembaca SQL
langsung membutuhkan selisih sampai **rincian per baris**, bukan hanya total tingkat polis.

Dasarnya artefak rancangan `Diagram-Skema-Tabel-NusantaraRe.xlsx`, sheet **EDM Treaty In Prop**,
sel `J85` · `R103` · `R116` — ketiganya sudah tergambar di sana.

⛔ Keduanya ikut aturan bangun ulang: **hanya baris ber-`SUMBER='GO'`** yang boleh dihapus.

### 3 · Endorse sesudah pembatalan — **boleh, dengan peringatan**

Polis yang sudah dibatalkan **masih boleh di-endorse lagi**. Generasi baru boleh ditumpuk di atas
generasi pembatalan.

⛔ **Nol gerbang peran.** Tidak ada alur persetujuan tambahan. Layar memberi **peringatan**
sebelum pengguna melanjutkan — peringatan, bukan penghalang.

| Lapisan | Tugasnya |
| --- | --- |
| `services` | ⛔ **tidak** menolak generasi di atas pembatalan |
| `repository` | mengembalikan **jenis generasi sebelumnya** supaya layar tahu kapan memperingatkan |
| React | menulis peringatannya |

⚠️ Kata *"dulu"* dalam keputusan dicatat apa adanya: bila kelak gerbang peran diperlukan,
keputusan itu diambil terpisah. Peran yang tersedia bila itu terjadi: `DeptHead` **242**
kemunculan di Treaty In, `UW` 26, `Manager` 35 — `KADIV` dan `Director` masing-masing **1**,
praktis tidak ada tingkat di atas `DeptHead`.

⭐ Test wajib: rantai **generasi 1 → 2 → pembatalan → 4** tersimpan utuh dan terbaca utuh.

### 4 · Lingkup pemindahan dokumen lama — **seluruhnya**

Setiap polis, **setiap generasinya**. Tidak ada penyaringan menurut status, tahun buku, atau lini
usaha.

⭐ Alasannya: selisih dihitung terhadap generasi sebelumnya. Bila hanya generasi terakhir
dipindah, **pembandingnya hilang**.

⚠️ Aturan `SUMBER='PEGA'` **beku** sekarang berlaku atas **seluruh** data historis, bukan
sebagian — taruhan perintah bangun ulang naik.

⛔ `[terbuka]` Besaran datanya belum diketahui. Empat angka diminta ke DBA: cacah baris, tahun
paling awal, ukuran total, dan **nomor generasi tertinggi** — yang terakhir sekaligus menguji
butir 1.

---

## ⭐ KEPUTUSAN 23-09-2026 sore — presisi uang dan daftar medan

### Presisi — **`NUMBER(38,8)`**

`[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."*

Bunyi lama: ~~*`NUMBER(20,8)` — dua belas digit di depan koma*~~ ⛔ **tidak cukup.**

⭐ **30 digit di depan koma, delapan di belakang** — batas tertinggi Oracle.

| Bukti korpus | Digit kiri |
| --- | ---: |
| `99999999999999999.99` sentinel "tanpa batas" | **17** |
| `181500000000.00` ambang kapasitas treaty | **12** |
| `150000000000.00` ambang kapasitas treaty | **12** |

⚠️ Ambang dagang nyata **tepat di batas** dua belas digit — nol ruang lega, padahal penjumlahan lintas mata uang dapat melewatinya.

⭐ `NUMBER` di Oracle berpanjang **berubah-ubah**: hanya digit bermakna yang tersimpan. Pelebaran ini **tidak memakan ruang tambahan**, sehingga tidak ada alasan menahannya di angka yang pas-pasan.

⛔ Butir `[data DBA]` presisi **ditutup** — oleh bukti, bukan oleh jawaban DBA.

### Daftar medan — dari Activity dan Section

`[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`, ikuti dari data yang digunakan di Activity dan Section."*

⭐ Daftar medan disusun ulang dari **302 berkas aturan** — **394 medan unik**, **232** di antaranya tampil di layar. Hasilnya: **`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`**.

Data guide turun derajat menjadi **penambal**. Ia sah hanya untuk **panjang maksimum per medan**, dan untuk 24 nama yang tidak tersapu aturan — `OldData` · `ProdKe` · `EDMNo` · `OldPolicyNo` · `BranchCode` dan lainnya, yang ditulis sistem, bukan diketik pengguna.

⛔ Butir `[data DBA]` panduan bentuk dokumen **ditutup**.

⚠️ **Cacah 394 tetap batas bawah, bukan total.** Penampung medan tak dikenal tetap wajib, dan wajib **kosong** sebelum pekerjaan dinyatakan selesai.

---

## TELEMETRI EKSEKUSI

### ⛔ Yang TIDAK diukur, dinyatakan tegas

| Ukuran | Keadaan |
| --- | --- |
| pengukuran **dari luar** *(`claude --print --output-format json`)* | ⛔ **TIDAK dilakukan** — perintah itu memulai sesi **baru dan terpisah**, yang ongkosnya bukan ongkos ronde ini |
| durasi | ⛔ **tidak diukur** |
| biaya | ⛔ **tidak diukur** |
| panggilan alat | ⛔ **tidak diukur terpisah** — baseline-nya tidak diambil |

⛔ **Tidak satu pun dari empat baris di atas ditaksir lalu disajikan sebagai angka terukur.**

### ⭐ Yang TERUKUR — selisih terhadap baseline sesi

| Ukuran | Nilai |
| --- | ---: |
| panggilan model | **15** |
| token keluar | **104.409** |
| token cache ditulis | **174.149** |
| token cache dibaca | **10.416.193** |
| berkas korpus dibaca | **163** |
| byte korpus dibaca | **17.639.259** |
| berkas korpus disunting | **0** |

⚠️ **Angka token sejati tidak terlihat dari dalam sesi.** Yang di atas adalah **pengurangan terhadap
baseline** yang diambil dari catatan sesi pada awal ronde — pendekatan terbaik yang tersedia,
⛔ **bukan** pengukuran dari luar.

### Disiplin ronde ini

| ⛔ Larangan | Dipatuhi? | Bukti |
| --- | :---: | --- |
| nol tabel dasar dirancang ulang | ✅ | Bab 10.1 — diperiksa dan dinyatakan |
| `spec.md` EDM tidak disunting | ✅ | **0** tulis |
| `spec-penyimpanan-relasional.md` NB tidak disunting | ✅ | **0** tulis — ia sumber, bukan sasaran |
| `grilling-ronde-1.md` tidak disentuh | ✅ | **0** tulis |
| jangan menutup butir `[terbuka]` | ✅ | **0** ditutup; **9** terdaftar *(4 warisan · 4 khas EDM · 1 baru)* |
| nol `CREATE TABLE` / DDL | ✅ | **0** |
| nol ADR baru | ✅ | **0** — Bab 10.6 |
| nol nama orang · nomor polis · contoh JSON | ✅ | **0** · **0** · **0** |
| ringkasan dihitung pola teks **lebih dulu** | ✅ | Bab 1.2 |

---

*Disusun 23 September 2026, sesudah spec penyimpanan NB Treaty In selesai dan bentuk kesembilan
tabel dasar ditetapkan final.*
