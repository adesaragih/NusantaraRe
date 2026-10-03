# Pertanyaan untuk **pengembang Pega lama**
## Migrasi Treaty Inward — Realisasi & Endorsement · modul *NB Treaty In*

> ⛔⛔ **RALAT MENYELURUH ATAS BERKAS INI — P18 DITARIK.** `[penyimpangan sadar]` 2026-09-22
>
> Berkas ini memuat **4 pernyataan** yang bersandar pada premis ⛔ *"isi langkah penetapan nilai
> tidak ikut terekspor"* — di baris **88 · 210 · 211 · 599**. **Premis itu keliru.**
>
> Isi langkah **ada di dalam ekspor**, di tag `PropertiesName`/`PropertiesValue` — **tanpa awalan
> `py`**. Tim migrasi memeriksa varian **dengan** awalan, menemukan nol, dan mempercayainya.
> ⭐ **Sebabnya dua huruf.** **938 dari 939** langkah `Property-Set` membawa isinya sendiri;
> **2.481** pasangan nama=nilai terisi; **nol** tanda terpotong. Lihat
> `VERIFIKASI-P18.md` dan `KOREKSI-P18-DIJALANKAN.md`.
>
> ⚠️ **Kalimat-kalimat itu sengaja TIDAK disunting satu per satu**, sebab sebagian berada di dalam
> **jawaban yang sudah diberikan pemiliknya** — mengubah jawaban orang lain bukan wewenang tim
> migrasi. Ralat ini berlaku atas seluruhnya.

---


Anda yang membangun atau merawat sistem lama. Tiga belas pertanyaan di bawah adalah hal-hal yang **terbaca di berkasnya tetapi tidak dapat kami artikan** — bukan karena berkasnya kurang, tetapi karena maksudnya hanya ada di kepala orang.

> ⭐ **15 pertanyaan untuk Anda** *(⭐ 2 baru dari ronde 4)*, diambil dari **48 pertanyaan** seluruh proyek modul ini. ⛔ Sisanya bukan urusan Anda dan tidak disertakan.
>
> **P14**, **P15**, dan **P27** saling terkait dan menahan pemetaan data. **P23**, **P31**, dan **P32** adalah tiga tempat di mana sistem lama dapat salah tanpa memberi tanda.
>
> ⭐ **Diurutkan dari yang paling menahan pekerjaan**, bukan menurut nomor.

## Cara menjawab

- ⭐ **Jawab langsung di bawah tiap pertanyaan.** Tidak perlu format khusus.
- ⭐ **Baris `rujukan:` paling bawah tiap pertanyaan boleh Anda abaikan** — itu catatan teknis
  untuk tim migrasi, bukan bagian dari pertanyaan.
- ⚠️ **Tidak tahu adalah jawaban yang sah**, dan lebih berguna daripada perkiraan. Tulis
  *"tidak tahu"* dan sebutkan siapa yang mungkin tahu.
- ⚠️ **Jangan menebak.** Sebuah tebakan yang masuk ke spesifikasi akan dibangun sebagai fakta.
- ⛔ Nomor pertanyaan **jangan diubah** — ia dipakai untuk melacak jawaban ke temuan aslinya.

---

## 1. P14 — Satu wadah keterangan dipakai bersama oleh 18 aturan, dan isinya berbeda-beda. Bagaimana cara membacanya?

Ada sebuah **wadah keterangan bernomor** — slot 1, slot 2, slot 3, dan seterusnya — yang dipakai
**bersama oleh 18 aturan berbeda**. ⛔ **Slot yang sama membawa hal yang berbeda di aturan yang
berbeda**: slot ke-3 membawa **seluruh muatan data dalam bentuk teks** di satu aturan, dan membawa
**posisi/jabatan** di aturan lain.

**Konteks:** kami hendak memetakan setiap slot ke kolom yang tepat. Bila pemetaannya dibuat satu
kali untuk semua, ⛔ **ia akan salah di sebagian besar aturan**.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang per-aturan, tidak ada makna
tetap · **(b)** ada dokumen pemetaan, **butuh berkas** · **(c)** ada aturan tak tertulis,
jelaskan. ⭐ Bila (a): **aturan mana yang mengisi slot itu sebelum penyimpanan dijalankan.**

**Dampak bila salah:** ⛔⛔ satu pemetaan yang salah **merusak enam aturan lain** yang memakai slot
yang sama — dan kesalahannya **tidak menimbulkan pesan galat**, hanya data yang salah tempat.

rujukan: halaman `InputData` 15 slot / 18 rule; `InputData.CARI3` di 7 rule —
OQ-059, butir baru #1

**Jawaban:**

> **(a) Memang per-aturan, tidak ada makna tetap.** `[keputusan work owner]` 2026-09-22
>
> `InputData` bukan struktur data melainkan **kantong parameter**. Arti tiap slot ditentukan
> oleh perintah SQL yang dipanggil, bukan oleh wadahnya. Tidak ada — dan tidak perlu dicari —
> satu pemetaan yang berlaku untuk semua aturan.
>
> **Bukti dari korpus.** `[terverifikasi]` Sepuluh aturan di NB Treaty In mengisi `InputData`.
> Slot yang sama membawa hal yang berbeda:
>
> | Slot | Dipakai | Jenis isi berbeda |
> | --- | ---: | --- |
> | `CARI1` | 9 aturan | 5 — ID grup treaty · kunci internal berkas · awalan nomor berkas · nomor penawaran · mata uang |
> | `CARI2` | 5 aturan | 5 — hasil hitungan · nomor polis · pesan galat · nomor urut baris · tanggal mulai |
> | `CARI3` | 4 aturan | 4 — jenis konfirmasi ceding · teks tetap `"Policy"` · muatan JSON penuh · tanggal akhir |
>
> Slot 3 membawa satu tanggal di satu aturan dan **seluruh muatan data dalam bentuk JSON** di
> aturan lain. Itu menutup kemungkinan adanya makna tetap.
>
> **Sepuluh aturan pengisinya** — semuanya di `NB Treaty In\Activity\`:
> `FetchTreatyGroupOldID` · `GeneratePolicyNoTreaty_Act` · `ProtectCurrencyTSI_Act` ·
> `SaveJsonPolisTreatyIn_Act` · `SaveViewSuggest` · `SetCurrency_act` · `SetTreatyCurrencyID` ·
> `TreatyNonPropOutSetSpreading` · `TreatyNonPropSetSpreading` · `TreatyRealizationCheckDuplicate`
>
> **Akibatnya untuk migrasi: masalah ini hilang, tidak dipindahkan.**
> Karena tidak ada lagi pemanggilan procedure dari SQL, setiap pemanggilan ke Oracle ditulis
> sebagai query dengan **parameter bernama**, bukan slot bernomor. Begitu parameter punya nama,
> satu slot tidak bisa lagi berarti dua hal, dan cacat ini tidak dapat terulang.
>
> **RALAT** `[penyimpangan sadar]` 2026-09-22 — paragraf penutup jawaban ini semula berbunyi:
>
> > *"Untuk menulis query penggantinya, kita tetap perlu tahu slot nomor berapa masuk ke kolom apa
> > pada setiap pemanggilan. Itu hanya ada di dalam naskah SQL-nya, dan tidak ada satu pun naskah
> > SQL di korpus — diperiksa, nol blok. Karena itu kebutuhan ini dipindahkan menjadi permintaan
> > berkas kepada pemilik export Pega, dicatat sebagai P41."*
>
> **Keliru.** Naskah SQL ada di dalam ekspor, pada tag `pyBrowseSQL` di rule tipe `RDBList` —
> **41 pernyataan di NB Treaty In, 1.151 di seluruh korpus.** Penelusuran yang menghasilkan
> "nol blok" hanya mencakup berkas yang menyebut `IsApproved`, lalu kesimpulannya diperluas tanpa
> dasar. **P41 sudah ditarik** dari lembar pemilik export Pega.
>
> Yang benar: **52 ikatan slot-ke-kolom terbaca langsung dari naskah SQL** yang sudah kita miliki.
> Sisanya adalah slot pada papan hitung sementara yang tidak pernah menjadi parameter basis data —
> arti slot itu hanya ada di langkah Property-Set, sehingga pemulihannya bergantung pada **P18**.
> Rinciannya dicatat di **P27**.

---

## 2. P15 — Sebuah fungsi yang membentuk seluruh muatan data tidak ada naskahnya. Apa isinya?

Ada satu fungsi buatan sendiri yang **mengubah data satu kasus menjadi satu teks besar**, dan teks
itulah yang disimpan ke basis data sebagai muatan utama. ⛔ **Naskah fungsi itu tidak ada di dalam
ekspor mana pun.**

**Konteks:** data lama tersimpan dalam bentuk yang dihasilkan fungsi ini. Sistem baru harus dapat
**membacanya kembali**.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — naskah fungsi tersebut, **atau** beberapa
**contoh nyata** hasilnya dari basis data *(3–5 contoh sudah cukup untuk menurunkan bentuknya)*.

**Dampak bila salah:** ⛔ **seluruh data lama tidak terbaca** oleh sistem baru.

rujukan: `@ASM.GetPageJSONString()` di `Activity\SaveJsonPolisTreatyIn_Act.xml` langkah 6 — OQ-012

**Jawaban:**

> **Naskah diterima, dan fungsinya tidak dimigrasi.** `[keputusan work owner]` 2026-09-22
>
> Isi `@ASM.GetPageJSONString()`, diberikan work owner:
>
> ```java
> PublicAPI tools = (PegaAPI) ThreadContainer.get().getPublicAPI();
> ClipboardPage stepPage = tools.getStepPage();
> String retValue = stepPage.getJSON(false);
> return retValue;
> ```
>
> **Artinya: fungsi ini tidak memilih apa pun.** Ia mengubah **seluruh halaman langkah** menjadi
> JSON apa adanya. Tidak ada penyaringan medan, tidak ada pemetaan, tidak ada bentuk yang dirancang.
>
> Itu menjelaskan mengapa isi kolom `JSONDATA` tidak dapat diterka dari aturan mana pun: bentuknya
> bukan skema, melainkan **potret halaman kerja pada saat penyimpanan** — apa pun yang kebetulan ada
> di sana ikut masuk.
>
> **Keputusan: penyimpanan pindah ke tabel relasional. JSON tidak dipakai lagi, baik untuk membaca
> maupun menulis.** Melengkapi keputusan pada **P29** yang semula hanya menyangkut pembacaan.
> Fungsi ini tidak dimigrasi.
>
> ---
>
> **Dua akibat yang perlu dicatat.**
>
> ~~`[terbuka]`~~ ✅ **P1 TERJAWAB 2026-09-22.** Butir ini semula berbunyi:
> > *"`[terbuka]` **P1 menjadi lebih penting, bukan berkurang.** Karena JSON tidak ditulis lagi,
> > sistem baru harus menulis langsung ke tabel-tabel yang selama ini diisi
> > `POOLDATA.PEGA_JSON_POLIS_TREATYIN`. Untuk itu kita tetap perlu tahu **tabel dan kolom mana
> > saja yang disentuh procedure itu** — permintaan P1 berdiri utuh untuk ketiga procedure."*
>
> ⭐ **Naskahnya diterima — dan bukan tiga, melainkan EMPAT.** `PEGA_TREATY_IN` menulis **dua**
> tabel *(`M_TREATY_IN` dan `POOLDATA.TREATY_IN` berisi 20 kolom datar)*;
> `PEGA_JSON_POLIS_TREATYIN` menulis **satu** tabel delapan kolom; `PROC_GENERATE_SEQUENCE_NUMBER`
> tidak menulis data bisnis; `PEGA_DELETE_ERROR_KONVERSI` — yang keempat — menghapus.
> Rinciannya di `PERTANYAAN-untuk-DBA.md` §P1.
>
> `[terbuka]` **Isi `JSONDATA` lama kemungkinan memuat jauh lebih banyak daripada data bisnis.**
> Karena `getJSON(false)` menyalin seluruh halaman, dokumen lama dapat memuat properti internal
> Pega dan — perlu diperiksa — **nilai berupa nama orang**. Bila kelak ada pemindahan data dari
> kolom itu, isinya wajib disaring lebih dulu, bukan disalin utuh.

---

## 3. P27 — Wadah slot bernomor: sebagian artinya tertulis di catatan. Apakah catatan itu masih benar, dan di mana sisanya?

Melanjutkan pertanyaan sebelumnya tentang wadah keterangan bernomor: ⭐ **kami menemukan sebagian
artinya tertulis di dalam catatan pengembang** — slot 44 untuk pengurangan, 31 untuk mata uang,
47 untuk saldo terutang, 13 untuk premi neto, 32 untuk premi bruto.

**Konteks:** catatan bisa basi tanpa ada yang tahu, dan kami hanya menemukan arti **sebagian**
slot — bukan semuanya.

**Bentuk jawaban yang diharapkan:** ⭐ **konfirmasi + kelengkapan** — apakah kelima arti di atas
masih benar, dan **di mana arti slot lainnya dapat dibaca**. Bila tidak ada tempatnya,
katakan begitu.

**Dampak bila salah:** ⛔ nilai uang masuk ke kolom yang salah, dan **kesalahannya tidak
menimbulkan pesan galat** — hanya angka yang salah tempat.

rujukan: `pyStepsDescription` — `Deduction (CARI44[deduction])` ·
`Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])` ·
`Gross Premi (CARI32, CARI31[CURRENCY])` — OQ-059, ronde 1 baru #1, ronde 2 pendalaman

**Jawaban:**

> **Terjawab dari korpus — tidak perlu ditanyakan ke siapa pun.** `[terverifikasi]` 2026-09-22
>
> **RALAT lebih dulu.** `[penyimpangan sadar]` Jawaban P14 menyatakan *"tidak ada satu pun naskah
> SQL di korpus — diperiksa, nol blok"*, dan atas dasar itu P14 dialihkan menjadi permintaan berkas
> **P41** kepada pemilik export Pega. **Pernyataan itu keliru.** Naskah SQL ada, tersimpan di tag
> `pyBrowseSQL` pada rule tipe `RDBList`: **41 pernyataan di NB Treaty In, 1.151 di seluruh korpus.**
> Penelusuran yang menghasilkan "nol blok" hanya mencakup berkas yang menyebut `IsApproved`, lalu
> kesimpulannya diperluas tanpa dasar. **P41 sudah ditarik.**
>
> **Dua koreksi atas pertanyaan ini:**
>
> Pertama, slot yang benar-benar punya arti tertulis ada **empat**, bukan lima: `CARI13`
> (NET PREMIUM), `CARI31` (CURRENCY), `CARI44` (deduction), `CARI47` (BALANCEDUETO). `CARI32` dan
> `CARI39` memang muncul di catatan tetapi **tanpa label** — artinya disimpulkan dari judul langkah,
> bukan ditulis.
>
> Kedua, ini bukan satu wadah melainkan **34 halaman berbeda** yang memakai kebiasaan slot bernomor
> yang sama, 598 sebutan, nomor slot tertinggi 55. Dihitung sebagai pasangan halaman+slot:
> **110 slot berbeda.** Diverifikasi dua cara.
>
> **Jawabannya: kedua jenis slot itu berbeda sifat, dan hanya satu yang berkaitan dengan basis data.**
>
> | Jenis | Contoh halaman | Arti dibaca dari |
> | --- | --- | --- |
> | parameter basis data | `InputData`, `ParamSeq`, `InsertHistory`, `OldID`, `ParamInProdCase` | naskah SQL — **sudah ada di korpus** |
> | papan hitung sementara | `InputXOL`, `ProtectSpreading`, `ProtectDeductible`, `ProtectMarine` | langkah Property-Set di dalam Activity |
>
> Dari naskah SQL, **52 ikatan slot-ke-kolom terbaca langsung**, misalnya
> `InputData.CARI1` = `NOOFFER` pada `TreatyRealizationCheckDuplicate`, `InputData.CARI2` =
> `BEGINDATE`, `ParamInProdCase.CARI16` = `zipcode`, `TempOccupation.CARI1` = `OCCUPATIONCODE`.
> Halaman yang sama dipakai untuk kolom berbeda di SQL berbeda — menegaskan jawaban (a) pada P14.
>
> `[terverifikasi]` **Keempat slot yang berlabel di catatan justru bukan parameter basis data.**
> `CARI13`, `CARI31`, `CARI44`, `CARI47` hidup di halaman `InputXOL`, dan `InputXOL` **tidak muncul
> sebagai parameter pada satu pun dari 1.151 pernyataan SQL di seluruh korpus.** Begitu pula
> `ProtectSpreading`, `ProtectDeductible`, dan `ProtectMarine`.
>
> `[terbuka]` Akibatnya, arti slot pada papan hitung sementara **tidak dapat dipulihkan dari SQL**.
> Ia hanya ada di dalam langkah Property-Set — yang isinya tidak ikut terekspor. Pemulihannya
> bergantung sepenuhnya pada **P18**, bukan pada pertanyaan kepada pengembang Pega lama.
> Catatan pengembang yang empat itu tetap berguna sebagai **uji silang** ketika P18 dijawab.

---

## 4. P23 — Pada sepertiga penggolong, keterangan yang tertulis berbeda dari yang dikerjakan aturan. Mana yang benar?

Ada 75 aturan kecil yang menggolongkan jenis bisnis. Masing-masing punya **dua hal**: keterangan
yang ditulis manusia, dan syarat yang benar-benar dijalankan. ⛔ **Pada 18 dari 52 yang punya
keduanya — sepertiga — keduanya tidak berbicara tentang hal yang sama.** Satu contoh: keterangannya
berbunyi *"Kode Bisnis 02 dan 58"*, sedangkan yang dijalankan memeriksa **jenis** bisnis bernama
*"Bonding"*.

**Konteks:** kami menulis spesifikasi dari aturan ini. Bila kami membaca keterangannya, kami salah
pada sepertiga kasus; bila kami membaca yang dijalankan, mungkin keterangannya yang mencerminkan
maksud sebenarnya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** yang dijalankan selalu benar,
keterangannya basi · **(b)** keterangannya menyatakan maksud, yang dijalankan mungkin bug ·
**(c)** kasus per kasus, ⭐ **perlu ditinjau bersama 18 aturan itu**.

**Dampak bila salah:** ⛔ sistem baru menggolongkan bisnis **berbeda dari sistem lama**, dan
penggolongan itu menentukan treaty, premi, dan laporan.

rujukan: 52 `When` punya `pyConditionString` + bentuk bertanda kurung; **18 tidak cocok (35 %)`;
contoh `When\IsBondingAndCustomBonds.xml` — ronde 2 baru #4

**Jawaban:**

> **(a) Yang dijalankan selalu benar; keterangannya basi.** `[keputusan work owner]` 2026-09-22
>
> Spesifikasi ditulis dari **syarat yang dijalankan**. Keterangan yang ditulis manusia
> (`pyConditionString`) **diabaikan seluruhnya** — tidak dipakai bahkan sebagai petunjuk.
>
> **Hitungan diperbaiki.** `[terverifikasi]` Pertanyaan ini semula mencatat "18 dari 52 tidak
> cocok (35 %)". Pemilahan ulang atas 75 rule `When` di NB Treaty In:
>
> | | Jumlah |
> | --- | ---: |
> | tanpa keterangan — masih teks bawaan `[Double click to add condition]` | 21 |
> | syarat tidak ikut terekspor | 4 |
> | keterangan cocok dengan yang dijalankan | 41 |
> | **bertentangan** | **9** |
> | | **75** |
>
> Jadi **9 dari 50** yang punya keduanya — **18 %**, bukan 35 %. Yang 21 itu bukan salah,
> melainkan tidak pernah diisi.
>
> **Dasar keputusan — jejak salin-tempel.** `[terverifikasi]` Keterangan yang sama persis,
> `Kode Bisnis = "24"`, dipakai di **lima rule berbeda** yang menjalankan syarat yang sama sekali
> berbeda: `IsHE` (jenis `"HE"`), `IsMarineCargo` (jenis `"MarineCargo"`), `IsMarineHullOffshore`
> (jenis `"Aneka"` + kode `10157`), `IsProductsLiability` (`"Aneka"` + `10021`),
> `IsProfessionalLiability` (`"Aneka"` + `10023`). Lima rule tidak mungkin semuanya berarti kode 24.
> Pola yang sama pada `IsBondingAndCustomBonds` (keterangan kode `"02"`, dijalankan jenis
> `"Bonding"`) dan `IsTravel` (keterangan kode `"77"`, dijalankan jenis `"Travel"`).
>
> **Sembilan rule yang bertentangan:** `IsBondingAndCustomBonds` · `IsHE` · `IsLife` ·
> `IsMarineCargo` · `IsMarineHullOffshore` · `IsOperatorLife` · `IsProductsLiability` ·
> `IsProfessionalLiability` · `IsTravel`
>
> `[terbuka]` **`IsOperatorLife`** tetap ditandai terbuka. Keterangannya menyebut *grup akses*
> `"GISFW:AuctionUsers"`, yang dijalankan memeriksa *grup kerja* `"ReasLife"` — dua konsep berbeda,
> bukan sekadar nama basi. Ada kemungkinan maksud rule ini pernah diubah. Tidak memblokir: sistem
> baru mengikuti yang dijalankan. Bila kelak ditemukan pengguna yang seharusnya lolos tetapi
> tertahan di jalur Life, keputusan ini wajib ditinjau ulang.
>
> `[terbuka]` **`IsLife`** menyimpang dari keterangannya dengan cara yang sama — keterangan
> menyebut kode `"10164"`, yang dijalankan memeriksa `BusinessOldId` `L1`-`L16`. Konsisten dengan
> keputusan P38 ("ikuti apa adanya"); dicatat agar tercatat bahwa penyimpangan itu disadari.

---

## 5. P31 — Ketika data kontrak gagal dibaca, sistem hanya mencatat di log lalu LANJUT. Disengaja?

Di enam tempat, sistem lama membongkar data kontrak dari satu kolom. Bila pembongkaran itu
**gagal** — misalnya isinya rusak — ⛔ **sistem hanya menulis catatan di log, lalu melanjutkan
langkah berikutnya.** Tidak ada pesan kepada pengguna, dan tidak ada penghentian.

**Konteks:** artinya perhitungan berikutnya dapat berjalan di atas data yang **kosong atau
separuh terisi**, dan hasilnya tampak normal.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** disengaja, kegagalan memang boleh
diabaikan · **(b)** seharusnya berhenti dan memberi pesan · **(c)** tidak tahu.
⭐ Dan: **seberapa sering catatan galat itu muncul** di log sistem lama, kalau bisa dilihat.

**Dampak bila salah:** ⛔⛔ angka yang salah **terlihat benar**, dan tidak ada satu pun tanda bahwa
datanya tidak lengkap.

rujukan: `catch(Exception e){ oLog.error("ReloadSection:Expection : "+...) }` — tanpa
penghentian; ⚠️ pesan log menyebut konteks lain *(section reload)* — ronde 3 baru #3

**Jawaban:**

> **(a) Hentikan, tampilkan galat ke pengguna.** `[keputusan work owner]` 2026-09-22
> ⚠️ **PENYIMPANGAN SADAR** dari perilaku Pega.
>
> Bila pembacaan data kontrak gagal, sistem baru **berhenti** dan menampilkan galat. Proses **tidak
> diteruskan** dengan data separuh terisi, dan tidak ada kasus yang tersimpan dari pembacaan yang
> gagal.
>
> `[terverifikasi]` Perilaku Pega sekarang: **enam langkah Java identik** membaca kolom data
> kontrak lalu `adoptJSONObject` ke halaman `TreatyIn`. Bila gagal, galatnya **hanya ditulis ke
> log** dan aktivitas **LANJUT**. Tidak ada penghentian, tidak ada pemberitahuan.
>
> ⚠️ **Ini penyimpangan yang disengaja, dan alasannya dicatat:** data separuh yang lolos diam-diam
> tidak pernah muncul sebagai galat — ia muncul sebagai angka yang salah berbulan-bulan kemudian.
> Di sistem yang mencatat komitmen keuangan, kegagalan yang terlihat lebih murah daripada kegagalan
> yang tersembunyi.
>
> **Yang wajib ada di implementasi:** galat disampaikan ke pengguna dengan sebab yang dapat
> ditindaklanjuti, bukan pesan teknis; dan kegagalan tetap tercatat di log seperti sekarang —
> penghentian **menambah**, tidak menggantikan pencatatan.
>
> ⚠️ `[terbuka]` Test migrasi wajib memeriksa: adakah kasus lama di basis data yang tersimpan
> dengan data kontrak kosong akibat perilaku lama ini. Bila ada, jumlahnya perlu diketahui sebelum
> migrasi, sebab kasus semacam itu **tidak akan dapat dibuat ulang** di sistem baru.

---

## 6. P32 — Dua cara penulisan tanggal dipakai bersamaan. Mana yang benar di mana?

Sistem lama menulis tanggal dengan **dua susunan berbeda** dalam modul yang sama: satu
**hari-bulan-tahun**, satu lagi **bulan-hari-tahun berikut jam**. Keduanya dipakai untuk medan
yang berbeda, dalam aturan yang berdekatan.

**Konteks:** tanggal 3 Februari dan 2 Maret **tidak dapat dibedakan** bila susunannya tertukar,
dan kekeliruan seperti ini **tidak menimbulkan pesan galat**.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh daftar** — medan mana memakai susunan mana.
⭐ Atau, bila lebih mudah: **contoh nilai nyata** dari tiga medan tanggal di basis data, supaya
kami dapat menurunkan susunannya sendiri.

**Dampak bila salah:** ⛔⛔ tanggal mulai dan berakhirnya kontrak **bergeser sampai sebelas bulan**,
dan tidak ada yang melihatnya sampai ada yang menagih.

rujukan: `@DateTime.CurrentDate("dd/MM/yyyy","")` lawan
`@DateTime.CurrentDate("MM/dd/yyyy hh:mm a","")` di `InputPolicyTreatyIn_preDT` dan
`DeptHeadTreatyInUW_preDT` — ronde 3 baru #5

**Jawaban:**

> **(a) Seragamkan, simpan sebagai tipe tanggal.** `[keputusan work owner]` 2026-09-22
> ⚠️ **PENYIMPANGAN SADAR** dari perilaku Pega.
>
> Sistem baru menyimpan tanggal sebagai **tipe tanggal Oracle**, bukan teks, dan memakai **satu
> format** di seluruh sistem. Format tampilan diurus di lapisan layar, bukan ikut tersimpan.
>
> `[terverifikasi]` Perilaku Pega sekarang: dua format berbeda dipakai di berkas yang sama —
> `@DateTime.CurrentDate("dd/MM/yyyy","")` untuk tanggal mulai dan
> `@DateTime.CurrentDate("MM/dd/yyyy hh:mm a","")` untuk tanggal akhir, keduanya di
> `DataTransform\InputPolicyTreatyIn_preDT.xml`. `[data DBA]` Dan penyimpanannya bertipe teks:
> `POOLDATA.TANGGAL_CLOSING.TANGGAL` adalah `VARCHAR2(10)`.
>
> ⛔ `[terbuka]` **Konsekuensi migrasi yang belum terselesaikan, dan ia tidak kecil.**
> Data lama tersimpan sebagai teks dengan **dua format tercampur**, dan sebagian nilainya
> **ambigu secara mutlak**: `05/06` dapat berarti 5 Juni atau 6 Mei, dan tidak ada di dalam nilai
> itu sendiri yang membedakannya.
>
> Sebelum migrasi dijalankan, dibutuhkan: **(1)** cara menentukan format per baris — misalnya dari
> kolom pendamping, tanggal pembuatan baris, atau rentang yang mustahil (`13/xx` pasti `dd/MM`);
> **(2)** hitungan berapa baris yang tetap ambigu sesudah cara itu diterapkan; dan **(3)** keputusan
> apa yang dilakukan atas baris yang tetap ambigu.
>
> ⛔ Menebak format per baris berarti menggeser tanggal berlaku kontrak sampai sebelas bulan.
> Butir ini **wajib selesai sebelum migrasi data**, bukan sesudah.

---

## 7. P24 — Bendera "sudah disetujui" dibandingkan sebagai teks di satu tempat dan sebagai angka di tempat lain. Mana yang dimaksud?

Pemeriksaan **pertama** yang dilewati setiap pengajuan membaca sebuah penanda "sudah disetujui".
⛔ Di satu tempat penanda itu dibandingkan **sebagai teks** — nol di dalam tanda kutip — dan di
tempat lain **sebagai angka** — nol tanpa kutip.

**Konteks:** bila penanda itu belum pernah diisi, "kosong" dan "nol" adalah dua hal berbeda dalam
satu perbandingan dan satu hal yang sama dalam perbandingan lainnya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** penanda ini angka, pembandingan teks
adalah kekeliruan · **(b)** penanda ini teks · **(c)** boleh keduanya, tidak pernah jadi masalah.
⭐ Dan: **apa arti penanda ini saat belum pernah diisi** — belum diperiksa, atau ditolak?

**Dampak bila salah:** ⛔ pengajuan yang belum diperiksa **dianggap ditolak**, atau sebaliknya
**diteruskan** seolah sudah disetujui.

rujukan: `.IsApproved = '0'` berkutip 2× vs `.IsApproved==0` telanjang 2× di layar;
`When\isApproved.xml` menguji `= 1` — ronde 2 baru #5

**Jawaban:**

> **Ikuti sistem lama apa adanya.** `[keputusan work owner]` 2026-09-22
>
> Aturan nilai: **`0` berarti ditolak; nilai apa pun selain itu berarti disetujui.**
> Ini menyalin persis tabel keputusan yang hidup, `DecisionTable\isApproved.xml`
> (`= 0` → `No`, bawaan → `YES`). Penanda `IsApproved` tetap dua keadaan; tidak
> dijadikan enum, tidak ditambah keadaan baru.
>
> Untuk bagian **teks atau angka** dari pertanyaan ini: jawabannya **teks**.
> `[terverifikasi]` Kolom tabel keputusan itu ber-`pyColumnDataType` = `text`,
> jadi perbandingan yang hidup memang perbandingan teks. Rule `When\isApproved.xml`
> yang menguji `= 1` bukan aturan yang dipakai — tidak satu pun kotak Decision di
> `Flow\InputRealizationTreatyIn.xml` menyambung ke sana.
>
> **Perilaku penolakan, ditetapkan menurut tingkat:**
>
> | Siapa menolak | Akibat |
> | --- | --- |
> | Admin | berkas langsung diselesaikan sebagai **ditolak** |
> | Atasan (Sec Head, Dept Head) | berkas **kembali ke admin** |
>
> Baris kedua sama dengan sistem lama: cabang `No` dari kotak `Is Correct?` sesudah
> `Acceptance by Head. Treaty` dan `Acceptance by Dept. Head` memang menuju
> `Input Realitation`. `[terverifikasi]`
>
> Baris pertama **berbeda** dari sistem lama, dan ini disengaja. `[penyimpangan sadar]`
> Di Pega, cabang `No` milik admin menuju `End3` — satu-satunya kotak End di seluruh
> alur, yang juga menjadi muara jalur sukses dari `HIT SERVICE ARASAPAS`. Jadi berkas
> yang ditolak admin dan berkas yang selesai normal berakhir di tempat yang sama.
> Sistem baru memisahkannya: penolakan admin diselesaikan sebagai ditolak, bukan
> sebagai selesai. Nilai `Resolved-Rejected` sudah ada di data — dipakai sebagai
> penyaring di `ReportDefinition\GetListOpportunity.xml`.
>
> `[terbuka]` **Nilai kosong dan `NULL` terbaca sebagai DISETUJUI**, sebagai akibat
> langsung dari aturan "selain `0` berarti disetujui". Keputusan ini diambil tanpa
> mengetahui berapa banyak baris semacam itu ada di Oracle. Permintaan cacah ke DBA
> tetap berdiri, dan bila ternyata ada berkas **yang masih berjalan** dengan nilai
> kosong, keputusan ini **wajib ditinjau ulang** sebelum go-live.
>
> `[terbuka]` Nama tabel dan kolom Oracle untuk penanda ini belum diketahui — tidak
> ada satu pun blok SQL atau Java di korpus NB Treaty In yang menyebut `IsApproved`
> (diperiksa, nol). Wajib didapat dari DBA sebelum pekerjaan data dimulai.

---

## 8. P25 — Wewenang di beberapa tempat ditentukan oleh POSISI seseorang dalam sebuah daftar. Apakah urutan itu dijamin?

Beberapa aturan memutuskan apa yang boleh dilihat seseorang dengan melihat **antrean kerja
nomor satu** atau **nomor dua** dalam daftar antrean pengguna itu — ⛔ **bukan dengan menyebut
nama antreannya.**

**Konteks:** kalau seseorang ditambahkan ke antrean baru, atau daftarnya diurutkan ulang,
⛔ **perilaku aturan berubah tanpa ada yang menyentuh aturan itu.**

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** urutannya memang dijamin, jelaskan
bagaimana · **(b)** tidak dijamin, ini memang rapuh · **(c)** tidak tahu. ⭐ Bila (a) atau (c):
⭐ **antrean mana yang seharusnya dimaksud** pada posisi satu dan posisi dua.

**Dampak bila salah:** ⛔ wewenang seseorang **berpindah diam-diam** ketika keanggotaan antreannya
berubah karena alasan yang sama sekali lain.

rujukan: `pxRequestor.OperatorID.pyWorkBasketList(1)` 3× di `When`;
`OperatorID.pyWorkBasketList(2).pyWorkBasketName` 1× di layar — ronde 2 baru #7

**Jawaban:**

> **(b) Tidak dijamin — memang rapuh, dan tidak dimigrasi.** `[keputusan work owner]` 2026-09-22
>
> Pencarian berdasarkan **nomor urut** antrean tidak dibawa ke sistem baru. Pertanyaannya
> diganti: bukan *"antrean nomor dua Anda namanya apa"*, melainkan *"apakah pengguna ini
> anggota antrean X"*. Hasilnya sama selama urutan daftar kebetulan benar, dan berhenti
> salah ketika urutan itu berubah karena sebab lain.
>
> Tiga tempat di NB Treaty In yang terkena: `Section\SFAPortal_OpportunitiesList.xml`
> (posisi 2), `DataTransform\InputPolicyTreatyIn_preDT.xml` (posisi 2), dan
> `When\IsUW.xml` (posisi 1). `[terverifikasi]` Di seluruh korpus polanya ada di
> 41 berkas pada 9 modul — perubahan yang sama berlaku di sana ketika modul itu digarap.
>
> Perubahan ini tidak memerlukan informasi dari pihak mana pun dan bisa langsung dikerjakan.
>
> ---
>
> **Pertanyaan turunan — `When\IsUW.xml` hanya menguji nama antrean Fac:**
>
> > **(a) Layar itu memang dipakai bersama Fac dan Treaty.** `[keputusan work owner]` 2026-09-22
>
> Empat belas elemen layar yang dikendalikan `IsUW` — 7 di `Section\DetailPolicyTreatyIn.xml`
> dan 7 di `Section\GeneralPolicyTreatyIn.xml` — memang **hanya tampil untuk pengguna Fac**.
> Itu disengaja, bukan sisa salin-tempel. Elemen-elemen itu **tetap dibangun**.
>
> Konsekuensi yang perlu diketahui sejak sekarang: `[terbuka]`
> Karena elemen-elemen itu tidak pernah tampil bagi pengguna Treaty, **14 elemen tersebut
> tidak dapat diuji selama pengujian penerimaan NB Treaty In.** Pembuktiannya baru mungkin
> ketika modul Fac digarap. Sampai saat itu elemen-elemen tersebut berstatus dibangun tetapi
> belum terbukti, dan harus masuk daftar uji modul Fac, bukan dianggap selesai di sini.
>
> `[dugaan]` Pembacaan ini didukung oleh kelas layarnya, `ASM-FW-GISFW-Data-PolicyTreatyIn`,
> yang memang kelas bersama. Tidak diperiksa di sistem berjalan.

---

## 9. P33 — Satu medan "nama operator" diisi dari dua sumber berbeda. Mana yang dimaksud?

Medan yang mencatat **nama operator** diisi dari **dua tempat berbeda** di dua aturan yang
berdekatan: satu memakai **pengenal akun**, satu memakai **nama tampilan orangnya**.

**Konteks:** keduanya terlihat serupa di layar tetapi berbeda isinya, dan medan ini masuk ke
jejak audit.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** seharusnya pengenal akun ·
**(b)** seharusnya nama tampilan · **(c)** memang berbeda per tahap, jelaskan.

**Dampak bila salah:** ⚠️ jejak audit memuat **dua jenis pengenal di satu kolom**, dan mencari
riwayat seseorang menjadi tidak dapat diandalkan.

rujukan: `DeptHeadTreatyInUW_preDT` lgk 2 `= OperatorID.pyUserIdentifier` lawan
`InputPolicyTreatyIn_preDT` lgk 10 `= OperatorID.pyUserName` — ronde 3 baru #8

**Jawaban:**

> **(b) Seharusnya nama tampilan.** `[keputusan work owner]` 2026-09-22
>
> Medan `OperatorName` berarti **satu hal saja**: nama tampilan orangnya. Pengisian dari pengenal
> akun di `DataTransform\DeptHeadTreatyInUW_preDT.xml` adalah **bug**, bukan perbedaan maksud
> antar tahap, dan tidak ditiru di sistem baru.
>
> **Jejak lengkapnya.** `[terverifikasi]`
>
> ```
> .PolicyTreatyIn.OperatorName = OperatorID.pyUserName        <- InputPolicyTreatyIn_preDT (langkah admin)
> .PolicyTreatyIn.OperatorName = OperatorID.pyUserIdentifier  <- DeptHeadTreatyInUW_preDT   (langkah Dept Head)
> .OperatorName -> InputData.CARI4 -> kolom PIC
>                                     tabel POOLDATA.historyakseptasiproduction
> ```
>
> Dasar keputusan: INSERT yang sama juga mengisi kolom terpisah `AKSES_LOGIN` dari
> `OperatorID.pyUserIdentifier` — lihat `RDBList\InsertViewSuggest_SQL.xml`. Tabelnya **sudah punya
> kolom sendiri untuk pengenal akun**, jadi `PIC` memang diperuntukkan bagi nama orangnya. Bila
> tidak, kedua kolom itu akan menyimpan hal yang sama.
>
> Di sistem baru: nama tampilan ke `PIC`, pengenal akun ke `AKSES_LOGIN`. Satu medan, satu arti.
>
> `[terbuka]` **Data lama pada kolom `PIC` berisi campuran** dua jenis nilai dan tidak dapat
> dirapikan tanpa memetakan nama ke akun satu per satu. **Tidak memblokir:** penelusuran riwayat
> seseorang dilakukan lewat `AKSES_LOGIN`, yang selalu diisi dari sumber yang sama sepanjang
> seluruh rentang data lama. Pencarian lewat `PIC` tidak boleh diandalkan untuk data sebelum
> go-live.

---

## 10. P34 — Tujuh belas langkah MENGHAPUS pesan kesalahan. Kenapa?

Sistem lama memasang pesan kesalahan di **122 tempat**, dan ⛔ **menghapus pesan kesalahan di
17 tempat.**

**Konteks:** menghapus pesan kesalahan berarti melanjutkan sesuatu yang sesaat sebelumnya
ditandai salah. Itu bisa sah — misalnya membersihkan pesan lama sebelum memeriksa ulang — atau
bisa berarti validasi dilewati.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** pembersihan sebelum pemeriksaan ulang,
wajar · **(b)** memang untuk melewati validasi tertentu, sebutkan mana · **(c)** campuran, perlu
ditinjau.

**Dampak bila salah:** ⚠️ sistem baru **menolak** data yang selama ini diterima, atau **menerima**
data yang seharusnya ditolak.

rujukan: `Page-Set-Messages` 105 + `Property-Set-Messages` 17 lawan `Page-Clear-Messages` 17
di 92 `Activity` — ronde 3 baru #12

**Jawaban:**

> **(c) Campuran — dan peninjauannya dipersempit dari 17 menjadi 5.** `[keputusan work owner]` 2026-09-22
>
> Pemilahan dilakukan berdasarkan **urutan langkah**: apakah penghapusan terjadi sebelum atau
> sesudah pesan dipasang di dalam aktivitas yang sama. `[terverifikasi]`
>
> | | Jumlah | Perlakuan |
> | --- | ---: | --- |
> | menghapus **sebelum** pesan apa pun dipasang — pembersihan awal | 12 | diterima apa adanya, langsung ditulis ulang |
> | menghapus **sesudah** pesan dipasang | 5 | **ditahan** sampai dipastikan |
>
> Tiga dari yang 12 menulis alasannya sendiri di catatan langkah: *"Clear spreading error message"*,
> *"Init remove error message"*.
>
> **Lima yang ditahan** — semuanya menghapus di langkah terakhir aktivitasnya, sesudah validasi
> memasang pesan, dan semuanya berada di rantai spreading/premi:
>
> | Aktivitas | Pesan dipasang di langkah | Dihapus di langkah |
> | --- | --- | ---: |
> | `Protection_Act` | 9, 11, 14, 16, 17, 29, 34, 35, 36 | 37 dari 37 |
> | `CheckSpreadingProtect_ACT` | 11, 12, 15, 25, 28, 30 | 31 dari 31 |
> | `ProtectPremiPolicy_Act` | 3, 11, 16, 17, 18 | 19 dari 19 |
> | `SumTSIPremiSpreadedRNM_Act` | 38, 40 | 51 dari 52 |
> | `cekSpreadingFactIn` | 8, 10, 18 | 19 dari 19 |
>
> `[terverifikasi]` **Tidak satu pun dari 17 langkah itu bersyarat.** Seluruhnya berjalan selalu —
> nol langkah memuat `pyStepsPreCondParamsWhen`.
>
> `[dugaan]` **Belum dapat dibuktikan** bahwa penghapusan mengenai halaman yang sama dengan tempat
> pesannya dipasang. Tag halaman langkah `pyStepsPage` **tidak ada sama sekali di ekspor** — nol
> kemunculan di 92 Activity NB Treaty In. Bila penghapusan mengenai halaman yang berbeda, kelima
> kasus itu bisa jadi wajar juga. Hal ini akan terjawab sendiri bila **P18** dipenuhi: ekspor ulang
> yang sama membawa halaman langkah sekaligus isi Property-Set.
>
> `[terbuka]` **Pertanyaan lanjutan kepada pengembang Pega lama, dipersempit menjadi satu kalimat:**
> *"Di lima aktivitas ini, penghapusan pesan di langkah terakhir itu disengaja atau sisa uji coba?"*
>
> Bila jawabannya "sisa uji coba", kelimanya tidak dimigrasi dan validasinya menjadi **aktif** di
> sistem baru. Akibatnya sistem baru akan **menolak sebagian data yang selama ini lolos**. Itu wajib
> diketahui sebelum go-live, bukan sesudah.

---

## 11. P16 — Satu perhitungan membaca tabel internal Pega secara langsung. Apa penggantinya?

Sebuah perhitungan membaca **tabel kerja internal sistem lama** secara langsung, seolah tabel
biasa. ⛔ **Tabel itu akan hilang** begitu sistem lama ditinggalkan.

**Konteks:** apa pun yang bergantung pada perhitungan itu akan **berhenti bekerja tanpa pesan
galat**.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** perhitungan ini masih dipakai,
jelaskan untuk apa · **(b)** sudah tidak dipakai · **(c)** tidak tahu. ⭐ Bila (a): **angka apa
yang sebenarnya dicari** dari situ.

**Dampak bila salah:** sebuah angka di layar menjadi **selalu nol** dan tidak ada yang menyadarinya.

rujukan: `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` di `RDBList\GetCountClaim.xml` — butir baru #5

**Jawaban:**

> **(a) Masih dipakai, dan angka yang dicari sudah jelas.** `[keputusan work owner]` 2026-09-22
>
> Yang dicari: **berapa banyak berkas klaim yang terhubung ke id reasuransi ini.** Bila lebih dari
> nol, sebuah peringatan dipasang ketika pengguna hendak membuat penawaran ganda.
>
> **Rantai pemanggilannya utuh** — jadi perhitungan ini hidup, bukan sisa. `[terverifikasi]`
>
> ```
> TreatyRealizationCheckXOLList -> SetTreatyIn_Act -> CheckDuplicateOffer -> GetCountClaim
> ```
>
> `Activity\CheckDuplicateOffer.xml` berisi 8 langkah. Langkah 6 memanggil `GetCountClaim` dan
> menaruh hasilnya di halaman `HasilClaim`; langkah 8 adalah `Page-Set-Messages` dengan syarat
> `HasilClaim.pxResults(1).CARI1 > 0`.
>
> Naskah SQL-nya, dari `RDBList\GetCountClaim.xml`:
>
> ```
> select count(1) as "CARI1" from DATAPEGA.PC_ASM_FW_GCNMFW_WORK
> where masterid = {TreatyWarning.CAIREINSFACIN}
> ```
>
> `PC_ASM_FW_GCNMFW_WORK` adalah tabel kerja internal Pega untuk kerangka `GCNMFW` — bukan `GISFW`
> yang dipakai modul ini — dan hilang bersama Pega.
>
> **Penggantinya:** hitung jumlah klaim dari **penyimpanan klaim yang baru**, bukan dari tabel kerja
> Pega. Modul klaim termasuk yang dimigrasi, jadi sumbernya akan tersedia. Tidak perlu ditanyakan
> kepada siapa pun.
>
> `[terbuka]` **Urutan pengerjaan.** Peringatan ini baru dapat bekerja **setelah modul klaim selesai
> dimigrasi**. Sampai saat itu ia akan selalu diam — persis jenis kegagalan yang disebut di
> pertanyaan ini: angka menjadi selalu nol dan tidak ada yang menyadarinya. Karena itu butir ini
> masuk **daftar uji lintas modul**, dan tidak boleh ditandai selesai pada pengujian NB Treaty In.

---

## 12. P26 — Satu nilai pengaturan tampak salah ketik. Apakah ia bekerja?

Di empat tempat, sebuah pengaturan yang mengunci medan diberi nilai yang **tampaknya salah
ketik** — satu huruf hilang dari kata yang seharusnya. Di tempat lain nilai yang sama ditulis
dengan benar.

**Konteks:** kalau sistem lama **mengabaikan** nilai yang salah ketik itu, maka keempat medan itu
sebenarnya **tidak pernah terkunci** — padahal seseorang berniat menguncinya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** salah ketik, penguncian tidak bekerja ·
**(b)** nilai itu sah, bekerja normal · **(c)** tidak tahu, perlu dicoba di sistem lama.

**Dampak bila salah:** ⚠️ medan yang seharusnya terkunci **dapat disunting**, atau sebaliknya kami
mengunci medan yang selama ini bebas.

rujukan: `pyDisabledNew` = `truewhn` 4×, di samping `true` 17× dan `always` 12× — ronde 2 baru #8

**Jawaban:**

> **(a) Memang salah ketik, dan penguncian lewat setelan itu tidak bekerja — tetapi tidak berdampak.**
> `[keputusan work owner]` 2026-09-22
>
> Medannya **sudah terkunci** oleh setelan lain pada kendali yang sama. Niat orang yang menulisnya
> tetap tercapai.
>
> **Buktinya.** `[terverifikasi]` Keempat kemunculan `truewhn` menjaga medan yang sama —
> `.InstallmentPercentage` (persentase angsuran) — di empat Section: `DetailPolicyTreatyIn`,
> `GeneralPolicyTreatyIn`, `DetailDeptHeadTreatyIn_UW`, `GeneralDeptHeadTreatyIn_UW`.
>
> Isi blok kendalinya, dibaca utuh sebagai satu blok (41 kunci), bukan dari teks di sekitarnya:
>
> ```
> pyValue        = .InstallmentPercentage
> pyFormat       = pxNumber
> pyReadOnly     = true
> pyDisabledNew  = truewhn
> pyRequired     = false
> ```
>
> `pyReadOnly` dan `pyDisabledNew` berada pada kendali yang **sama**. Sebaran nilai `pyDisabledNew`
> di NB Treaty In: `false` 539x, `true` 17x, `always` 12x, `truewhn` 4x — `truewhn` bukan nilai yang
> dikenal.
>
> **Untuk sistem baru:** jadikan `InstallmentPercentage` hanya-baca di keempat layar itu.
> Tidak perlu dicoba di sistem lama, tidak perlu ditanyakan kepada siapa pun.

---

## 13. P17 — Aturan yang hanya hidup di layar belum kami baca sama sekali. Apa yang ada di sana?

⛔ **Kami belum membuka satu pun dari 31 berkas layar** modul ini — dan **kelima berkas terbesar
seluruh modul ada di antaranya**. Artinya pemeriksaan yang dilakukan **di layar** — medan wajib,
tombol yang hanya muncul untuk jabatan tertentu, perhitungan yang berjalan saat mengetik — ⛔ **belum
terbaca.**

**Konteks:** kami menyatakan ini terang-terangan agar tidak ada yang menyangka ronde ini sudah
meliput modulnya. ⭐ Ronde berikutnya perlu arahan **layar mana yang paling penting**.

**Bentuk jawaban yang diharapkan:** ⭐ **urutan prioritas** — dari daftar layar utama
*(input polis treaty, rincian polis, rincian non-proporsional, layar Kepala Departemen)*, mana yang
**paling banyak memuat aturan bisnis**, bukan hanya tampilan.

**Dampak bila salah:** ronde berikutnya membaca **8 MB layar** dan menemukan sebagian besarnya
tata letak, sementara aturan yang penting tetap terlewat.

rujukan: 25 `Section` + 6 `Harness`, **nol dibuka**; lima terbesar semuanya `Section` — butir baru #6

**Jawaban:**

> **TIDAK BERLAKU LAGI — sudah dikerjakan sendiri, tidak perlu ditanyakan.** `[keputusan work owner]` 2026-09-22
>
> Pertanyaan ini meminta **urutan prioritas** pembacaan lapisan layar. Grilling ronde 4 sudah
> mengurai **seluruhnya**: 25 `Section` + 6 `Harness`, **16.635.834 byte, 100 %, nol berkas tidak
> terurai.** `[terverifikasi]`
>
> Yang dihasilkan sudah menjadi pertanyaan tersendiri — **P44** sampai **P47** — jadi tidak ada lagi
> yang perlu diminta dari pengembang Pega lama untuk butir ini.
>
> Sisa pekerjaannya bukan pertanyaan kepada siapa pun, melainkan **penulisan spec**: syarat tampil,
> medan wajib, medan hanya-baca, validasi sisi layar, dan urutan pengisian dituangkan ke dalam
> bab User Stories dan Acceptance Criteria.
>
> Ditutup bukan karena dijawab, melainkan karena **pertanyaannya tidak berlaku lagi**.

---

## 14. P47 — Seratus enam puluh setelan "wajib isi", tetapi tidak satu pun menempel pada medan. Di mana kewajibannya berlaku?  ⭐ *(BARU — ronde 4)*


Di ke-31 layar terdapat **160 setelan "wajib isi" yang berisi**. ⛔ **Tetapi tidak satu pun dari
870 medan** di layar-layar itu membawanya. Artinya setelan-setelan itu menempel pada **sesuatu yang
lain** — dan kami tidak dapat memastikan pada apa.

**Konteks:** kami tidak ingin menyimpulkan bahwa modul ini **tidak punya medan wajib** hanya karena
kami mencari di tempat yang salah. Ini persis jenis kekeliruan yang sudah beberapa kali menjerat
kami.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** setelan itu memang menempel pada
wadah/baris, bukan pada medan, jelaskan bagaimana ia berlaku · **(b)** ia peninggalan yang tidak
lagi bekerja · **(c)** ada tempat lain yang seharusnya kami periksa, sebutkan.
⭐ Dan yang paling membantu: ⭐ **daftar medan yang benar-benar wajib diisi** menurut sistem lama —
walau hanya untuk satu layar utama.

**Dampak bila salah:** ⛔ sistem baru **menerima** kontrak yang selama ini ditolak karena ada medan
kosong, atau sebaliknya **menolak** yang selama ini diterima.

rujukan: `pyRequired` 160 berisi di 31 berkas layar *(ronde 2 §2.4)* lawan 0 pada `pyUserData` dari
870 sel `pyType=FIELD` *(ronde 4 Bab A.3)*

**Jawaban:**

> **(a) Setelan itu MEMANG menempel pada medan. Premis pertanyaan keliru.** `[terverifikasi]` 2026-09-22
> Terjawab dari korpus — tidak perlu ditanyakan kepada siapa pun.
>
> **RALAT.** Pertanyaan ini menyatakan *"tidak satu pun dari 870 medan membawanya"*. Penelusuran
> yang menghasilkan kesimpulan itu mencari di `pyUserData`. Setelan `pyRequired` tidak tinggal di
> sana — ia tinggal di dalam blok sel layar, berdampingan dengan `pyValue` yang menyebut nama
> medannya.
>
> Diperiksa ulang: dari 160 `pyRequired=true`, **seluruhnya dapat dipasangkan dengan nama medan**.
> Pemasangan dilakukan dalam jendela blok yang sama, bukan menurut urutan.
>
> **27 medan wajib berbeda, tersebar di 6 layar:**
>
> | Layar | Jumlah | Medan |
> | --- | ---: | --- |
> | `GeneralPolicyTreatyIn` | 22 | `Claim` `ClaimPaymentType` `ClaimType` `Deduction1` `Deduction2` `EndDate` `ID` `IDCurrency` `OutstandingClaim` `OveriddingCommOgp` `OveriddingCommOnp` `PremiOgp` `PremiOnp` `Quartal` `RiCommOgp` `RiCommOnp` `SalvageValue` `StartDate` `StatementDate` `TypeTax` `YearOfQuartal` |
> | `DetailPolicyTreatyIn` | 21 | sama, tanpa `ID` |
> | `DetailDeptHeadTreatyIn_UW` | 16 | `Claim` `Deduction1` `Deduction2` `EndDate` `ExcessLoss` `OutstandingClaim` `OveriddingCommOgp` `OveriddingCommOnp` `PremiOgp` `PremiOnp` `ResultOnp1` `RiCommOgp` `RiCommOnp` `SalvageValue` `StartDate` `StatementDate` |
> | `GeneralDeptHeadTreatyIn_UW` | 16 | sama dengan di atas |
> | `ListSuggest` | 3 | `IsApproved` `ProductionDate` `Suggest` |
> | `InputHistoricalSurveyReportDtl` | 1 | `DateofSurvey` |
>
> **Yang perlu diperhatikan saat menulis spec:**
>
> Layar Dept Head **tidak** mewajibkan `ClaimPaymentType`, `ClaimType`, `IDCurrency`, `Quartal`,
> `TypeTax`, `YearOfQuartal` — enam medan yang wajib di layar admin. Sebaliknya layar Dept Head
> mewajibkan `ResultOnp1`, yang tidak wajib di layar admin.
>
> Jadi kewajiban **berbeda menurut tingkat**, bukan seragam. Itu harus dipertahankan di sistem baru,
> dan tidak boleh disederhanakan menjadi satu daftar tunggal.
>
> `[terbuka]` Beberapa medan wajib itu juga punya syarat `pyRequiredWhen` — kewajiban bersyarat,
> bukan mutlak. Contoh yang sudah diketahui: `ListSuggest.IsApproved` wajib hanya bila dua pengguna
> tertentu yang membuka (lihat **P25**, tempat ID operator ditulis mati diganti peran). Syarat
> bersyarat itu perlu disalin satu per satu ketika bab Acceptance Criteria ditulis.

---

## 15. P48 — DITARIK  ⛔ *(lahir ronde 4, ditarik pada hari yang sama)*


> **RALAT.** `[terverifikasi]` 2026-09-22
>
> Butir ini menyatakan bahwa `DecisionTable\BusinessType_DeT.xml` kini tinggal **sembilan baris**
> dengan nomor tertinggi 34, dan menyimpulkan sebagian besar barisnya pernah dihapus.
> **Tidak ada baris yang hilang. Tabelnya utuh, 36 baris, seluruh nilainya terbaca.**
>
> Bunyi pernyataan yang ditarik:
> > *"Salah satu tabel penggolong jenis usaha kini memuat sembilan baris, tetapi nomor barisnya
> > melompat-lompat dan yang tertinggi bernomor 34. Itu tanda bahwa tabel ini pernah jauh lebih
> > panjang dan sebagian besar barisnya dihapus sepanjang waktu."*
>
> **Sebabnya.** Deretan `-1, 2, 3, 5, 7, 22, 29, 30, 34` bukan nomor baris tabel. Itu isi
> `pyRowNum` di dalam `pyOrConditions` — penanda **baris mana saja yang membawa daftar OR**,
> berbasis nol. Delapan baris membawa daftar OR pada kolom `BusinessOldId`; satu entri `-1`
> adalah grup kosong milik kolom `GroupPanel`. Ini persis jebakan #3 yang tercantum di prompt
> ronde 4.
>
> **Hitungan ulang, tiga cara, pada berkas yang sama:**
>
> ```
> <pyResults>      1 blok : 36 rowdata, 36 berisi
> <pyCondition>    2 blok : 36 rowdata/36 berisi (GroupPanel)
>                           36 rowdata/22 berisi (BusinessOldId)
> <pyOrConditions> kolom-2 : 106 nilai di dalam 8 grup
> ```
>
> Sensus nilai: 22 sel tunggal + 106 nilai OR = **128 kode `BusinessOldId`, seluruhnya unik**.
> Hasilnya 36 baris dengan **34 jenis usaha berbeda** — `Medicare` dan `FireStyle1` masing-masing
> muncul dua kali pada panel yang berbeda. Angka 34 itulah yang terbaca sebagai "pernah 34 baris".
>
> Uji silang: baris berbasis-nol 34 menghasilkan `"Life"` dengan daftar `L1`-`L18`, cocok dengan
> `When\IsLife.xml` yang menguji `BusinessOldId` `L1`-`L16`.
>
> Isi lengkap tabel sudah tercatat pada jawaban **P19** di lembar pemilik export Pega.
> **Tidak ada yang perlu ditanyakan kepada siapa pun untuk butir ini.**

---

## Sesudah Anda menjawab

⭐ Kembalikan lembar ini apa adanya — **tidak perlu dirapikan**. Jawaban setengah lengkap tetap berguna; yang tidak berguna adalah lembar yang ditahan sampai lengkap.

⚠️ **Yang tertahan selama pertanyaan ini belum terjawab:** **P14**, **P15**, dan **P27** saling terkait dan menahan pemetaan data. **P23**, **P31**, dan **P32** adalah tiga tempat di mana sistem lama dapat salah tanpa memberi tanda. Bagian lain pekerjaan tetap berjalan.

---

*Disusun dari pembacaan berkas ekspor sistem lama, tiga putaran. Tidak ada pertanyaan di lembar ini yang dijawab sendiri oleh tim migrasi.*
