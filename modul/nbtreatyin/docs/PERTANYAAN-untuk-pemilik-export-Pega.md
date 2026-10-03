# Pertanyaan untuk **pemilik export Pega**
## Migrasi Treaty Inward — Realisasi & Endorsement · modul *NB Treaty In*

Anda yang menyiapkan berkas ekspor sistem lama untuk kami. Pertanyaan di bawah semuanya tentang hal yang **tidak ikut terkirim** — bukan tentang isinya, tetapi tentang apakah ia dapat dikirim ulang.

> ⭐ **3 pertanyaan untuk Anda** — **P30**, **P19**, dan **P42**. ⛔ Sisanya bukan urusan Anda dan tidak disertakan.
>
> ⛔⛔ **DUA BUTIR SUDAH DITARIK oleh tim migrasi**, dan keduanya dibiarkan di tempatnya supaya terbaca: **P18** *(nomor 1)* dan **P41** *(nomor 4)*. ⭐ **Anda tidak perlu mengerjakan apa pun untuk keduanya.**
>
> ⚠️ **P18 sebelumnya disebut yang terberat di seluruh proyek.** Kalimat lamanya berbunyi: ⛔ *"Ketiganya menahan. **P18** adalah yang terberat di seluruh proyek: tanpa jawabannya, 94 dari setiap 100 langkah perhitungan tidak dapat ditiru."* — **keliru, dan kekeliruannya milik tim migrasi.** Isinya ada di ekspor Anda sejak awal.
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

## 1. P18 — ⛔ DITARIK OLEH TIM MIGRASI

> ⛔⛔ **DITARIK.** `[penyimpangan sadar]` 2026-09-22
>
> ⭐⭐ **Anda tidak perlu mengerjakan apa pun untuk butir ini.** Tidak ada ekspor ulang yang
> diminta. Lampirannya juga dibatalkan.
>
> ⚠️ **Butir ini ditarik oleh tim migrasi, bukan dijawab oleh Anda.** Tidak seorang pun di luar
> tim migrasi pernah menjawabnya. **Kekeliruannya milik tim migrasi sepenuhnya.**
>
> **Apa yang keliru:** isi langkah penetapan nilai **ada di dalam ekspor yang Anda kirim** — di tag
> `PropertiesName` dan `PropertiesValue`, **tanpa awalan `py`**, di dalam `pyParamArray` milik tiap
> langkah. Tim migrasi memeriksa `pyPropRef`, `pyPropertiesName`, dan `pyPropertiesValue` —
> **dengan** awalan `py`. Ketiganya memang kosong, dan kekosongan itu terlanjur dipercaya tiga
> ronde berturut-turut.
>
> ⭐ **Sebabnya dua huruf: `py`.**
>
> **Yang terukur sesudah diperiksa ulang** — **938 dari 939** langkah `Property-Set` membawa isinya
> sendiri *(99,9 %)*; **2.481** pasangan nama=nilai terisi; **707** pasangan pada ke-51 aturan yang
> semula diminta; **nol** tanda nilai terpotong pada tujuh uji keutuhan. Rinciannya di
> `VERIFIKASI-P18.md`.
>
> ⚠️ **Akibat kedua yang didalilkan butir ini juga gugur.** Dari **34** medan wajib-sekaligus-
> terkunci, **28** punya langkah pengisi; tiga sisanya — `.ExcessLoss`, `.OutstandingClaim`,
> `.SalvageValue` — **diketik underwriter** di layar sebelumnya. Layar Kepala Departemen
> **dapat disimpan**.
>
> ---
>
> **Bunyi pertanyaan yang ditarik — dikutip utuh, tidak dihapus:**
>
> > ## 1. P18 — Isi 268 langkah penghitungan tidak ada di dalam berkas yang kami terima. Bisa dikirim ulang?
>
> > Sistem lama punya langkah-langkah yang **menetapkan nilai** — hasil perhitungan premi, bagian,
> > pengurangan, dan pajak. Kami dapat melihat **bahwa** langkah itu ada, dan **urutannya**,
> > ⛔ **tetapi tidak satu pun rumus atau nilai yang ditetapkannya ikut terkirim.** Yang ada hanya
> > salinan sementara yang tertinggal dari layar penyunting, dan salinan itu **terpotong di tengah**
> > sehingga tidak dapat dipakai.
>
> > **Konteks:** tanpa isi langkah-langkah ini, perhitungan uang sistem lama tidak dapat ditiru —
> > hanya ditebak.
>
> > ⭐ **Lingkupnya sudah kami persempit lebih dulu, supaya Anda tidak mengerjakan yang tidak kami
> > butuhkan.** Folder yang kami terima memuat 92 aturan dengan 938 langkah, tetapi sebagian besar
> > bukan milik modul ini:
>
> > | Tahap penyaringan | Aturan | Langkah |
> > | --- | ---: | ---: |
> > | seluruh isi folder | 92 | 938 |
> > | dibuang — tidak terjangkau dari titik masuk mana pun | −28 | −561 |
> > | dibuang — tidak dibangun ulang karena penyimpanan JSON ditinggalkan | −5 | −109 |
> > | ⭐ **yang kami minta** | ⭐ **51** | ⭐ **268** |
>
> > ⭐ **Turun 72 % dari angka yang mungkin Anda dengar sebelumnya.**
>
> > **Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — ekspor ulang **51 aturan** pada lampiran
> > terpisah, dengan **isi langkah disertakan**.
>
> > ⭐ Bila tidak seluruhnya dapat dikerjakan, kerjakan **dari urutan teratas lampiran**: sepuluh
> > teratas sudah mencakup **56 %** dari seluruh langkah, dua puluh teratas **76 %**.
>
> > ⚠️ Bila ekspor jenis ini memang **tidak pernah** menyertakan isi langkah, katakan begitu, dan
> > sebutkan **cara lain yang tersedia** *(cetakan layar aturan, dokumentasi, atau akses baca ke
> > sistem lama)*.
>
> > **Dampak bila salah:** ⛔⛔ seluruh perhitungan uang dibangun dari tebakan, dan **selisihnya baru
> > muncul di laporan keuangan berbulan-bulan kemudian**.
>
> > ⛔ **Satu akibat lagi yang baru diketahui:** sembilan medan pada layar Kepala Departemen
> > **wajib diisi tetapi terkunci** — pengguna tidak dapat mengetiknya. Yang mengisinya adalah
> > langkah-langkah ini. Selama isinya belum ada, ⛔ **layar Kepala Departemen tidak akan dapat
> > disimpan sama sekali**, sehingga butir ini menahan bukan hanya perhitungan uang tetapi juga
> > alur persetujuan.
>
> > rujukan: langkah `Property-Set` di 51 `Activity` terjangkau; `PropertiesValue` = 74 dengan salinan
> > editor, **0** tanpa; `pyStepsCallParams` hanya placeholder; penjangkauan ditelusuri dari `Flow`,
> > `Section`, `Harness`, `FlowAction`, `DataTransform` — ronde 2 baru #1, lingkup diperbarui 2026-09-22
>
> > ⭐ **Lampiran:** `LAMPIRAN-P18-ATURAN-DIMINTA.md` — 51 nama aturan, jumlah langkah masing-masing,
> > diurutkan dari yang terberat.
>
> > **Jawaban:**
>
> > > *(tulis di sini)*
>
> > ---
>
>
> ---
>
> ⛔ **Tidak ada yang perlu dikirim. Nol aturan, nol ekspor ulang, nol cetakan layar.**

**Jawaban:**

> ⛔ **Tidak perlu dijawab.** Butir ini ditarik sebelum Anda sempat menjawabnya.

---

## 2. P30 — Dua langkah lagi yang dipanggil tetapi tidak ikut terkirim, dan keduanya di jalur penting

Selain langkah yang sudah kami tanyakan, ada **dua langkah lagi** yang dipanggil sistem lama
tetapi **tidak ada di dalam berkas yang kami terima**. Satu dipanggil oleh langkah yang
**menyimpan data kontrak**; satu lagi dipanggil oleh langkah yang **menjumlahkan nilai
pertanggungan dan premi**.

**Konteks:** keduanya di jalur yang paling penting — menyimpan, dan menghitung uang.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — kedua langkah itu. Bila ekspornya memang
tidak menyertakan langkah dari bagian lain sistem, sebutkan **bagian mana** keduanya tinggal,
supaya kami dapat memintanya terpisah.

**Dampak bila salah:** ⛔ dua langkah yang selama ini berjalan **tidak dibangun**, dan tidak ada
yang tahu ia pernah ada — sebab yang hilang **tidak meninggalkan jejak**.

rujukan: `CopyToPolicy` dipanggil `SaveJsonPolisTreatyIn_Act`;
`SumTSIPremiSpreadRNMMultiCob_Act` dipanggil `SumTSIPremiSpreadedRNM_Act`; keduanya di luar
66 target yang ada — ronde 3 baru #2

**Jawaban:**

> **Separuh tertutup: `CopyToPolicy` tidak dipakai.** `[keputusan work owner]` 2026-09-22
>
> Pemanggilan `CopyToPolicy` dari `Activity\SaveJsonPolisTreatyIn_Act` **tidak digunakan** dan
> **tidak dimigrasi**. Tidak perlu diminta kepada pemilik export.
>
> `[terverifikasi]` Berkasnya memang tidak ada di NB Treaty In; yang ada hanya di tiga modul
> Fac In — `Endorsment Fac In`, `NB FacIn`, `RNW Fac In`. Di NB Treaty In namanya disebut satu kali,
> yaitu pada pemanggilan yang kini dinyatakan tidak dipakai itu.
>
> ---
>
> `[terbuka]` **`SumTSIPremiSpreadRNMMultiCob_Act` masih dicari.**
>
> Dipanggil oleh `Activity\SumTSIPremiSpreadedRNM_Act` — penjumlah nilai pertanggungan dan premi.
> **Tidak ada di seluruh korpus 21 modul**, padahal dipanggil oleh empat modul: `NB Treaty In`,
> `NB FacIn`, `RNW Fac In`, `Endorsment Fac In`.
>
> Berbeda dari `CopyToPolicy`, tidak ada satu pun salinan di modul mana pun untuk dibandingkan,
> dan letaknya di rantai uang. **Permintaan ke pemilik export tetap berdiri untuk butir ini:**
> di bagian mana aturan itu tinggal, supaya dapat diminta terpisah.

---

## 3. P19 — Isi dua tabel keputusan juga tidak terkirim, dan satu di antaranya masih sering diubah

Ada dua **tabel keputusan** — daftar "kalau begini maka begitu" — dan ⛔ **baris-barisnya tidak
ikut terkirim.** Yang tersisa hanya namanya. Salah satunya **terakhir diubah 22 September 2025**
dan punya **delapan versi sebelumnya**, jadi ia jelas masih dipakai dan masih disesuaikan.

**Konteks:** salah satu dari keduanya adalah **pemeriksaan pertama** yang dilewati setiap
pengajuan.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — kedua tabel keputusan itu beserta
barisnya. Bila tidak bisa, **cetakan layarnya** sudah cukup.

**Dampak bila salah:** ⛔ penggolongan bisnis menjadi berbeda, dan itu menentukan **treaty mana
yang dipakai** untuk sebuah penutupan.

rujukan: `DecisionTable\isApproved.xml` · `DecisionTable\BusinessType_DeT.xml`, nol
`pyCriteriaValue`/`pyResult` — OQ-026, ronde 2 baru #3

**Jawaban:**

> **TIDAK JADI DIMINTA — baris kedua tabel keputusan ADA di dalam ekspor.** `[terverifikasi]` 2026-09-22
>
> **RALAT.** `[penyimpangan sadar]` Butir ini ditulis atas dasar penelusuran tim migrasi yang
> mencari tag `pyCriteriaValue` dan `pyResult`, lalu menyimpulkan barisnya tidak terkirim.
> **Kesimpulan itu keliru.** Baris tabel keputusan Pega tidak disimpan di tag tersebut, melainkan
> di `pyCondition`, `pyOrConditions`, dan `pyResults`. Keduanya lengkap di dalam ekspor.
>
> **Tidak ada yang perlu Anda kirim untuk butir ini. Mohon abaikan.**
>
> ---
>
> **Isi `DecisionTable\isApproved.xml`:**
>
> | Kolom | Syarat | Hasil |
> | --- | --- | --- |
> | `pyWorkPage.PolicyTreatyIn.IsApproved` | `= 0` | `No` |
> | | bawaan | `YES` |
>
> **Isi `DecisionTable\BusinessType_DeT.xml`:** 36 baris, dua kolom —
> `.Quotation.GroupPanel` dan `.Quotation.BusinessOldId` — menghasilkan `BusinessType`,
> bawaan `"UNKNOWN"`. Evaluasi berhenti di baris pertama yang cocok (`pyEvaluateAllRows` = `no`).
>
> Sensus nilai, diverifikasi: **22 sel kolom-2 terisi tunggal + 106 nilai di dalam 8 daftar OR =
> 128 nilai `BusinessOldId`, seluruhnya unik, tidak ada yang muncul dua kali.** Delapan baris
> ber-daftar-OR ada pada `pyRowNum` berbasis-nol 2, 3, 5, 7, 22, 29, 30, 34.
>
> Enam baris menguji `GroupPanel` saja tanpa menguji `BusinessOldId` — yaitu baris penampung di
> dalam masing-masing panel: `001`->`Medicare`, `003`->`Aneka`, `004`->`MarineCargo`,
> `005`->`Travel`, `006`->`Fire`, `008`->`Medicare`. Karena evaluasi berhenti di baris pertama yang
> cocok, baris penampung itu menangkap sisa yang tidak tertangkap baris di atasnya.
>
> **Uji silang yang menguatkan:** baris berbasis-nol 34 menghasilkan `"Life"` dengan daftar OR
> `L1`-`L18`. Ini cocok dengan `When\IsLife.xml` yang menguji `BusinessOldId` `L1`-`L16` —
> dua rule berbeda, sumber terpisah, hasil bersesuaian.
>
> **Catatan pemastian, 2026-09-22 sore.** `[terverifikasi]` Laporan grilling ronde 4 menyebut
> `BusinessType_DeT` berisi **9 baris** dan menyebut `pyCondition`/`pyOrConditions`/`pyResults`
> kosong. **Keduanya keliru.** Dihitung ulang tiga cara pada berkas yang sama:
>
> ```
> <pyResults>   1 blok : 36 rowdata, 36 berisi
> <pyCondition> 2 blok : 36 rowdata/36 berisi (GroupPanel) · 36 rowdata/22 berisi (BusinessOldId)
> <pyOrConditions> blok kolom-2 : 106 nilai
> ```
>
> Angka **9** berasal dari salah satu dari dua hal, keduanya bukan baris tabel: `<pyLabel>` muncul
> 9 kali (label versi rule — `01-01-52 (Available)` dan seterusnya), dan nilai `GroupPanel` yang
> unik juga berjumlah 9. Jumlah baris tabel tetap **36**.
>
> Yang **benar** dari laporan ronde 4: `pyCriteriaValue` dan `pyResult` memang **nol kemunculan** —
> jadi pengamatan ronde 2 atas kedua tag itu tepat; hanya kesimpulannya yang keliru. Benar pula
> bahwa `pyColumnDataType` kedua kolom adalah **`text`**, yang menguatkan keputusan P24.
>
> `[terbuka]` Tabel ini terakhir diubah 22 September 2025 dan punya delapan versi sebelumnya.
> Isi yang terbaca di sini adalah **versi yang ada di dalam ekspor**. Bila ada perubahan sesudah
> tanggal ekspor, itu belum tercermin. Perlu dipastikan tanggal ekspor korpus sebelum tabel ini
> dijadikan dasar spesifikasi.

---

## 4. P41 — DITARIK

> **RALAT.** `[penyimpangan sadar]` 2026-09-22
> Butir ini ditulis atas dasar pernyataan tim migrasi bahwa naskah SQL tidak ikut terekspor.
> **Pernyataan itu keliru.** Naskah SQL ada di dalam ekspor, di dalam tag `pyBrowseSQL` pada
> rule tipe `RDBList` — **41 pernyataan di NB Treaty In, 1.151 di seluruh korpus.**
>
> Bunyi permintaan yang ditarik:
> > *"Sepuluh perintah SQL yang dipanggil tidak ikut terkirim … naskah perintah basis data
> > yang dipanggil oleh sepuluh aturan berikut."*
>
> **Tidak ada yang perlu Anda kirim untuk butir ini.** Mohon abaikan.
>
> Satu-satunya isi basis data yang memang belum ada adalah **badan tiga stored procedure**
> (`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`, `POOLDATA.PEGA_JSON_POLIS_TREATYIN`,
> `POOLDATA.PEGA_TREATY_IN`). Itu sudah diminta secara terpisah kepada DBA sebagai **P1**,
> bukan kepada pemilik export.

---

## 5. P42 — Sembilan belas medan tampil di layar, tetapi tidak satu pun aturan mengisinya. Dari mana datangnya?  ⭐ *(BARU — ronde 4)*


Layar sistem lama menampilkan **sembilan belas keterangan** yang **tidak disentuh oleh satu pun
aturan** yang kami terima — tidak diisi, tidak dibaca, tidak dihitung. Delapan di antaranya adalah
**nilai uang beserta mata uangnya**: batas, retensi, premi bersih, dan pendapatan premi
diperkirakan.

Kami menduga keterangan itu **datang dari satu kolom basis data yang isinya berupa dokumen**, yang
dibongkar saat layar dibuka — bukan dari aturan. Kalau dugaan itu benar, sistem baru yang dibangun
hanya dari aturan **tidak akan pernah menampilkannya.**

**Konteks:** kami tidak dapat memastikan sendiri, sebab isi kolom dokumen itu tidak ikut terkirim
kepada kami.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** benar, keterangan itu datang dari
kolom dokumen · **(b)** datang dari bagian sistem lain yang belum kami terima, sebutkan bagian
mana · **(c)** medan itu memang kosong dan tidak pernah terisi.
⭐ Bila **(a)** atau **(b)**: ⭐ **butuh contoh** — tiga isi kolom itu apa adanya.

**Dampak bila salah:** ⛔ delapan medan uang **selalu kosong** di sistem baru, dan ⛔ **tidak ada
pesan galat** yang memberi tahu — hanya angka yang hilang dari layar.

rujukan: 19 properti layar tanpa kemunculan di 92 `Activity` + 12 `DataTransform` + 75 `When` +
41 `RDBList`; `LIMITVALUE`, `RETENTIONVALUE`, `EPIVALUE`, `NETPREMIVALUE` + empat mata uang
pasangannya; bandingkan kolom `BrowseTreatyInDetail` dan mekanisme `adoptJSONObject` ronde 3 §2.1

**Jawaban:**

> **Terjawab — bukan dari kolom dokumen, melainkan dari view relasional.** `[keputusan work owner]` 2026-09-22
>
> Kesembilan belas medan itu **kolom nyata** pada view `POOLDATA.TREATYINDETAILJOINEDM`, bukan
> hasil perhitungan. Itulah sebabnya tidak satu pun `Activity` atau `DataTransform` menghasilkannya.
>
> Delapan medan uang yang dipersoalkan seluruhnya ada sebagai kolom:
> `LIMITCURRENCY`/`LIMITVALUE` · `RETENTIONCURRENCY`/`RETENTIONVALUE` · `EPICURRENCY`/`EPIVALUE` ·
> `NETPREMICURRENCY`/`NETPREMIVALUE`.
>
> **Uji kecukupan:** `ReportDefinition\BrowseTreatyInDetail.xml` merujuk 33 medan; ke-33-nya ada di
> view; nol yang hilang. `[terverifikasi]`
>
> **Koreksi atas bunyi pertanyaan.** Kalimat "tidak disentuh oleh satu pun aturan" terlalu keras.
> Keempat medan uang dirujuk **33-37 kali** di `ReportDefinition`, sebagai `pyTargetProperty` pada
> kelas integrasi `ASM-FW-GISFW-Int-TREATYINDETAIL` dan kerabatnya. Yang benar: **tidak ada Activity
> atau DataTransform yang menghitungnya** — medan itu dibaca, bukan dihasilkan.
>
> `[terbuka]` Pada sisi Pega keempat medan uang ber-`pyStringType` = **`Text`** — uang disimpan
> sebagai teks, bukan bilangan. Tipe Oracle-nya perlu dipastikan sebelum penguraian di Go
> ditetapkan. Dicatat juga pada P29.


# Untuk Finance

---

## Sesudah Anda menjawab

⭐ Kembalikan lembar ini apa adanya — **tidak perlu dirapikan**. Jawaban setengah lengkap tetap berguna; yang tidak berguna adalah lembar yang ditahan sampai lengkap.

⚠️ **Yang tertahan selama pertanyaan ini belum terjawab:** **P30** — di bagian mana `SumTSIPremiSpreadRNMMultiCob_Act` tinggal — masih menahan sebagian tiket rantai perhitungan uang. Bagian lain pekerjaan tetap berjalan.

⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat ini semula berbunyi: ⛔ *"Ketiganya menahan. **P18** adalah yang terberat di seluruh proyek: tanpa jawabannya, 94 dari setiap 100 langkah perhitungan tidak dapat ditiru."* **P18 ditarik.**

---

*Disusun dari pembacaan berkas ekspor sistem lama, empat putaran.*

⚠️ **Catatan kejujuran.** Kalimat penutup lembar ini semula berbunyi: *"Tidak ada pertanyaan di lembar ini yang dijawab sendiri oleh tim migrasi."* Itu **tetap benar** — tidak satu pun **dijawab** sendiri. Tetapi **dua di antaranya DITARIK** oleh tim migrasi sesudah pembacaan ulang membuktikan pertanyaannya sendiri keliru: **P41** pada 2026-09-22, dan **P18** pada hari yang sama. Penarikan bukan penjawaban, dan bedanya sengaja dibuat terbaca.
