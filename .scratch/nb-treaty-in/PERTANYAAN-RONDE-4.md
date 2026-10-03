# Pertanyaan baru — ronde 4
## Migrasi Treaty Inward — Realisasi & Endorsement · modul *NB Treaty In*

> ⛔⛔ **RALAT MENYELURUH ATAS BERKAS INI — P18 DITARIK.** `[penyimpangan sadar]` 2026-09-22
>
> Berkas ini memuat **5 pernyataan** yang bersandar pada premis ⛔ *"isi langkah penetapan nilai
> tidak ikut terekspor"* — di baris **118 · 249 · 251 · 253**. **Premis itu keliru.**
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


⭐ **Delapan pertanyaan baru, P42–P49.** P1–P41 sudah terpakai; P41 sudah ditarik.
⛔ **Berkas ini hanya berisi pertanyaan** — temuan dan buktinya ada di `grilling-ronde-4.md`.

**Cara menjawab:** tulis langsung di bawah tiap pertanyaan. Baris `rujukan:` boleh diabaikan —
itu catatan teknis, bukan bagian dari pertanyaan. ⚠️ *"Tidak tahu"* adalah jawaban yang sah dan
lebih berguna daripada perkiraan. ⛔ Nomor pertanyaan jangan diubah.

| Pemilik | Pertanyaan |
| --- | --- |
| pemilik export Pega | **P42** |
| Finance | **P43** |
| Product & Underwriting | **P44 · P45 · P46** |
| pengembang Pega lama | **P47 · P48** |
| DBA | **P49** |
| IAM | ⭐ **tidak ada ronde ini** |

---

# Untuk pemilik export Pega

## P42 — Sembilan belas medan tampil di layar, tetapi tidak satu pun aturan mengisinya. Dari mana datangnya?

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

---

# Untuk Finance

## P43 — Delapan angka uang tampil di layar tanpa ada perhitungan yang menghasilkannya. Apakah angkanya benar?

Di layar kontrak terdapat **delapan medan uang** — batas, retensi, premi bersih, dan pendapatan
premi diperkirakan, masing-masing berikut mata uangnya. ⛔ **Tidak satu pun aturan perhitungan
yang kami terima menghasilkan angka-angka itu.**

Artinya salah satu dari dua hal: angka itu **datang jadi dari basis data** tanpa dihitung ulang,
atau ia **dihitung di tempat yang belum kami lihat**.

**Konteks:** kami perlu tahu apakah angka yang dilihat pengguna itu **dihitung** atau **disalin**,
sebab keduanya berperilaku berbeda ketika data dasarnya berubah.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** angka itu disalin apa adanya dari data
kontrak, tidak pernah dihitung ulang · **(b)** dihitung, dan seharusnya berubah bila data dasarnya
berubah · **(c)** sebagian disalin sebagian dihitung, sebutkan mana.
⭐ Dan satu hal yang sangat membantu: **apakah pernah ditemukan angka di layar berbeda dari angka
di laporan?**

**Dampak bila salah:** ⛔ sistem baru **menghitung ulang** angka yang seharusnya tetap, atau
sebaliknya **membekukan** angka yang seharusnya mengikuti perubahan — dan selisihnya baru terlihat
saat penagihan.

rujukan: `pxCurrency` 113 medan di 31 berkas layar; delapan properti uang dalam daftar 19
properti hanya-layar — ronde 4 Bab A.6

**Jawaban:**

> **(a) Angka disalin apa adanya, tidak pernah dihitung ulang.** `[keputusan work owner]` 2026-09-22
>
> Sistem baru **menampilkan nilai yang tersimpan**, tidak menghitung ulang dari data dasar.
>
> **Pertanyaan pendamping — "pernahkah angka di layar berbeda dari angka di laporan?" — dijawab
> TIDAK.** Belum pernah ditemukan selisih. Artinya layar dan laporan selama ini membaca sumber yang
> sama, dan menyalin apa adanya tidak melestarikan penyimpangan apa pun.
>
> **Dasar dari korpus.** `[terverifikasi]` Kedelapan medan uang adalah **kolom tersimpan** pada view
> `POOLDATA.TREATYINDETAILJOINEDM` — `LIMITCURRENCY`/`LIMITVALUE`,
> `RETENTIONCURRENCY`/`RETENTIONVALUE`, `EPICURRENCY`/`EPIVALUE`,
> `NETPREMICURRENCY`/`NETPREMIVALUE`. Tidak satu pun `Activity` atau `DataTransform` di NB Treaty In
> menghasilkannya. Lihat **P29** dan **P42**.
>
> **Akibat yang menguntungkan:** butir ini **tidak lagi bergantung pada P18**. Bila angka harus
> dihitung ulang, rumusnya wajib diketahui lebih dulu — dan rumus itu justru yang belum terkirim.
> Dengan menyalin apa adanya, penyajian kedelapan medan dapat ditulis ke dalam spec sekarang juga.
>
> `[terbuka]` Nilai disimpan dengan presisi penuh, tanpa pembulatan di lapisan repository
> (lihat **P29**). Format penyajian di layar — jumlah desimal dan pemisah ribuan — belum ditetapkan
> dan perlu dipastikan sebelum bab Acceptance Criteria ditulis. **Tidak memblokir.**

---

# Untuk Product & Underwriting

## P44 — Sembilan puluh satu bagian layar dimatikan, tetapi delapan puluh di antaranya kosong. Yang mana yang masih perlu?

Kami sebelumnya melaporkan bahwa **91 bagian layar** diberi syarat tampil yang tidak mungkin
pernah benar — setara "tampilkan jika 1 sama dengan 2". ⭐ **Sekarang kami tahu isi masing-masing**,
dan gambarannya jauh lebih ringan daripada kesannya:

- ⭐ **80 dari 91 tidak menyembunyikan satu medan pun** — hanya wadah kosong atau label.
- **11 sisanya menyembunyikan 128 medan**, dan ⭐ **104 dari 128 itu berada di empat bagian saja**,
  di dua layar yang nyaris kembar — masing-masing menyembunyikan **26 medan**.

**Konteks:** pertanyaannya menjadi jauh lebih sempit. Kami tidak perlu keputusan tentang 91 hal —
⭐ **hanya tentang sebelas**, dan terutama tentang **empat bagian berisi 26 medan** itu.

**Bentuk jawaban yang diharapkan:** untuk **empat bagian besar** *(daftar jenis usaha dan sumber
bisnis, dan kembarannya untuk retro)* — pilihan ganda: **(a)** memang sudah tidak dipakai, boleh
hilang · **(b)** dimatikan sementara dan seharusnya kembali · **(c)** tidak tahu.
⭐ Untuk 80 yang kosong: **boleh kami abaikan seluruhnya?**

**Dampak bila salah:** ⛔ dua puluh enam medan yang masih dibutuhkan **hilang** dari sistem baru,
atau kami membangun dan merawat bagian yang **tidak pernah dipakai siapa pun**.

rujukan: 91 elemen *(`1=2` 78 · `1=2`/`NEVER` pada `pyContainerVisibleWhen` 11 · 2 gabungan
`&& NEVER`)*; sel di bawahnya `FIELD` 128 · `LABEL` 132 · `LAYOUT` 20 · `SUB_SECTION` 5;
empat elemen 26-medan di `BusinessAndSOBList` dan `BusinessAndSOBListRetro` — ronde 4 Bab B.2

**Jawaban:**

> **Dipersempit: 80 elemen dibuang tanpa ditanyakan; hanya 4 yang dibawa ke Product & Underwriting.**
> `[keputusan work owner]` 2026-09-22
>
> **Sebabnya terbukti struktural, bukan kebetulan.** `[terverifikasi]` Kedua tag pengendali
> kemunculan mengerjakan hal yang berbeda:
>
> | Tag | Mustahil | Menempel pada | Menyembunyikan medan? |
> | --- | ---: | --- | --- |
> | `pyCondition` | 80 | sel tunggal — umumnya label dan hiasan | **tidak satu pun** |
> | `pyContainerVisibleWhen` | 11 | wadah yang memuat medan | **128 medan** |
>
> Karena itu 80 elemen `pyCondition` **tidak perlu ditanyakan kepada siapa pun**: membangunnya
> berbiaya nol, membuangnya juga berbiaya nol. Dibuang dari daftar pekerjaan.
>
> `[terbuka]` **Yang dibawa ke Product & Underwriting hanya empat wadah, berisi 104 medan:**
>
> | Medan | Berkas |
> | ---: | --- |
> | 26 | `BusinessAndSOBList` — wadah pertama |
> | 26 | `BusinessAndSOBList` — wadah kedua |
> | 26 | `BusinessAndSOBListRetro` — wadah pertama |
> | 26 | `BusinessAndSOBListRetro` — wadah kedua |
>
> Pertanyaannya satu kalimat: **"Empat bagian layar ini berisi 104 medan dan sudah dimatikan.
> Masih dipakai, atau sudah ditinggalkan?"**
>
> `[terbuka]` Tujuh wadah sisanya menyembunyikan 24 medan — `Installments_ReadOnly` (6),
> `DetailPolicyTreatyIn` (4), `GeneralPolicyTreatyIn` (4), `SpreadingRiskList` (4),
> `DetailPolicyTreatyInNonProportionalEDM` (2), `SFAPortalOpportunities` dan
> `SFAPortalOpportunitiesHeader`. Dilampirkan sebagai daftar untuk ditandai ya/tidak, **tanpa
> pembahasan** — bobotnya terlalu kecil untuk memakan waktu pertemuan.
>
> **Bila jawabannya "sudah ditinggalkan": 128 medan tidak perlu dibangun.**

---

## P45 — Tiga puluh delapan medan dikunci permanen. Haruskah tetap terkunci?

**Tiga puluh delapan medan** diberi setelan hanya-baca yang **selalu berlaku** — tidak pernah dapat
diisi siapa pun, dalam keadaan apa pun. ⭐ **Tiga puluh enam di antaranya ada di dua layar Kepala
Departemen.**

Yang terbaca namanya, sebagai contoh: **"Statement Period"**, **"Survey Report"**, dan
**"Statement Date"**.

**Konteks:** sebelumnya kami menduga seluruh layar Kepala Departemen memang dibuat hanya-lihat.
⭐ **Ternyata tidak** — yang dikunci adalah **medan tertentu, satu per satu**, dan dua di antaranya
bahkan berada di layar biasa, bukan layar Kepala Departemen.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** benar, Kepala Departemen memang hanya
melihat, tidak mengubah · **(b)** seharusnya ada yang dapat diubah, sebutkan medan mana ·
**(c)** penguncian ini peninggalan, sudah tidak dimaksudkan.
⭐ Khususnya: **apakah "Statement Period" dan "Statement Date" memang tidak boleh diubah?**

**Dampak bila salah:** ⛔ pengguna **tidak dapat memperbaiki** keterangan yang keliru dan harus
menempuh jalan memutar, atau sebaliknya medan yang seharusnya terkunci **menjadi dapat diubah**.

rujukan: `pyReadOnlyCondition` selalu-benar **38**, pada elemen `pyUserData` yang induknya sel
`pyType=FIELD`; `DetailDeptHeadTreatyIn_UW` 18 · `GeneralDeptHeadTreatyIn_UW` 18 ·
`DetailPolicyTreatyIn` 1 · `GeneralPolicyTreatyIn` 1 — ronde 4 Bab B.3

**Jawaban:**

> **Ikuti apa adanya — 38 medan tetap terkunci permanen.** `[keputusan work owner]` 2026-09-22
>
> Perilaku sistem lama disalin tanpa perubahan. Tidak ada medan yang dibuka, tidak ada kewajiban
> yang dicabut.
>
> **Sensus dipastikan ulang — 38 tepat.** `[terverifikasi]`
>
> | Berkas | Terkunci |
> | --- | ---: |
> | `DetailDeptHeadTreatyIn_UW` | 18 |
> | `GeneralDeptHeadTreatyIn_UW` | 18 |
> | `DetailPolicyTreatyIn` | 1 |
> | `GeneralPolicyTreatyIn` | 1 |
>
> `pyReadOnlyCondition` yang selalu benar: `1==1` 32x, `1=1` 2x, `1 = 1` 2x, `ALWAYS` 2x.
> Setiap kunci mengenai **tepat satu medan**, bukan bagian atau tab.
>
> Medan yang terkunci di kedua layar Dept Head: `Deduction2` `ExcessLoss` `NetPremium`
> `OutstandingClaim` `OveriddingCommOgp` `OveriddingCommOnp` `PPHValue` `ResultOgp1` `ResultOgp2`
> `ResultOnp1` `ResultOnp2` `RiCommOgp` `RiCommOnp` `SalvageValue` `StatementType`.
>
> ---
>
> `[terbuka]` **Sembilan medan wajib diisi TETAPI terkunci.** Disilangkan dengan daftar medan wajib
> pada **P47**: `Deduction2` · `ExcessLoss` · `OutstandingClaim` · `OveriddingCommOgp` ·
> `OveriddingCommOnp` · `ResultOnp1` · `RiCommOgp` · `RiCommOnp` · `SalvageValue`.
>
> Kesembilannya wajib diisi tetapi tidak dapat diketik pengguna. Itu hanya dapat bekerja bila ada
> yang **mengisinya secara otomatis**. Rumus pengisinya berada di dalam **268 langkah yang belum
> terkirim (P18)**.
>
> **Akibat yang harus diketahui sebelum go-live:** selama P18 kosong, sistem baru akan menampilkan
> sembilan medan wajib yang selamanya kosong, dan layar Dept Head **tidak akan dapat disimpan**.
> Ini menjadikan P18 penahan bukan hanya bagi perhitungan uang, tetapi juga bagi **alur persetujuan
> tingkat Dept Head**.
>
> Enam medan terkunci lainnya — `NetPremium` `PPHValue` `ResultOgp1` `ResultOgp2` `ResultOnp2`
> `StatementType` — tidak wajib, sehingga tidak menahan penyimpanan.

---

## P46 — Sebuah medan terisi otomatis dengan kata "Inclusive". Apa artinya, dan kapan berlaku?

Layar sistem lama hampir **tidak pernah** mengisi medan dengan nilai awal — hanya **sebelas medan**
di seluruh modul yang punya nilai bawaan. Sembilan di antaranya diisi angka **nol**; dua sisanya
diisi kata **"Inclusive"**.

**Konteks:** kata itu satu-satunya nilai bawaan yang berupa **istilah bisnis**, bukan angka. Kami
tidak dapat menurunkan artinya dari berkas, dan nilai bawaan menentukan apa yang tersimpan ketika
pengguna tidak mengubah apa pun.

**Bentuk jawaban yang diharapkan:** satu paragraf — **apa arti "Inclusive"** pada medan itu,
**apa nilai lain yang mungkin**, dan **apakah "Inclusive" memang bawaan yang benar** atau
seharusnya kosong sampai pengguna memilih.

**Dampak bila salah:** ⚠️ kontrak yang pengguna tidak menyentuh medan itu **tersimpan dengan
pilihan yang tidak pernah ia buat**.

rujukan: `pyDefaultValue` berisi pada 11 sel `FIELD`; nilai `0` **9×**, `Inclusive` **2×**
*(`DetailPolicyTreatyIn` + 1 berkas)* — ronde 4 Bab A.5

**Jawaban:**

> **Ikuti apa adanya — perilaku sistem lama disalin utuh, termasuk ketidakseragamannya.**
> `[keputusan work owner]` 2026-09-22
>
> Medan `PolicyTreatyIn.TypeTax` bernilai awal `"Inclusive"` dan menggerakkan rumus potongan pajak
> pada brokerage:
>
> ```
> BrokerageFeeSebenarnya = @if(TypeTax == "Inclusive", Deduction / (102.2/100), Deduction)
> ```
>
> Bila `"Inclusive"`, potongan dibagi 1,022 — pajak 2,2 % dikeluarkan dari angkanya. Bila bukan,
> dipakai apa adanya. Rumus ini ada di **12 tempat pada 7 Activity**; angka `2.2` juga muncul di
> 4 berkas layar. `[terverifikasi]`
>
> **Tiga sifat yang ikut disalin, seluruhnya disadari:**
>
> `[terverifikasi]` **`Exclusive` tidak pernah tertulis di mana pun** — nol kemunculan di seluruh
> modul NB Treaty In. Cabang "bukan Inclusive" hanya tercapai bila medan berisi sesuatu yang lain,
> termasuk kosong.
>
> `[terverifikasi]` **Pembanding berupa teks persis.** `TypeTax == "Inclusive"`. Huruf kecil, spasi
> di belakang, atau ejaan lain jatuh ke cabang kedua dan brokerage **tidak** dibagi 1,022 —
> selisih 2,2 % pada setiap potongan, tanpa pesan galat.
>
> `[terbuka]` **Presisi tidak seragam untuk rumus yang sama.** Empat desimal di
> `InputPolicyTreatyInDetail_NonProp` dan `InputPolicyTreatyOutDetail_NonProp`; delapan desimal di
> `InputPolicyTreatyInDetail_preACT`, `InsertToTreatyXOLList`, `InsertToTreatyOutXOLList`, dan
> `SetPPNPPH`. Akibatnya **kontrak yang sama dapat menghasilkan brokerage yang sedikit berbeda,
> bergantung layar mana yang menyimpannya.**
>
> Keputusan "ikuti apa adanya" berarti ketidakseragaman ini **ditiru**, bukan diperbaiki. Perbedaan
> nilainya sudah ada di data sekarang, bukan risiko baru. Dicatat sebagai butir terbuka karena bila
> kelak ditemukan selisih brokerage yang dipersoalkan, **di sinilah asalnya**, dan keputusan ini
> wajib ditinjau ulang bersama Finance.
>
> Angka 2,2 % tetap ditulis di tempat pemakaiannya, tidak dijadikan setelan — sesuai keputusan yang
> sama. Bila tarif berubah, perubahannya menyentuh 12 tempat.

---

# Untuk pengembang Pega lama

## P47 — Seratus enam puluh setelan "wajib isi", tetapi tidak satu pun menempel pada medan. Di mana kewajibannya berlaku?

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

## P48 — DITARIK

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

# Untuk DBA

## P49 — Satu kolom riwayat tidak pernah diisi dari sistem lama. Apakah ada yang lain yang mengisinya?

Kami sudah menanyakan kolom penampung identitas operator pada tabel riwayat akseptasi **(P4)**.
⭐ **Sekarang kami dapat menjawab separuhnya sendiri, dan jawabannya tegas:**

⛔ **Dari seluruh 1.151 pernyataan SQL di seluruh sistem lama — bukan hanya modul ini — tidak satu
pun menulis ke kolom itu.** Tidak satu pun bahkan menyebutnya sebagai kolom. Ke-67 kemunculan
namanya ternyata **nama halaman pengguna**, bukan nama kolom.

**Konteks:** yang kami periksa hanyalah pernyataan SQL yang ditulis di dalam sistem lama.
⛔ Kami **tidak** dapat memeriksa program yang tersimpan di dalam basis data, pemicu basis data,
pekerjaan terjadwal, atau sistem lain di luar Pega.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** kolom itu memang tidak pernah terisi,
boleh ditinggalkan · **(b)** diisi dari luar Pega, sebutkan oleh apa · **(c)** dulu terisi, sekarang
tidak.
⭐ Yang paling menentukan, dan cepat: ⭐ **berapa baris yang kolom itu tidak kosong** di basis data
sekarang. Nol berarti (a); lebih dari nol berarti (b) atau (c).

**Dampak bila salah:** ⚠️ sebuah keterangan jejak audit yang selama ini ada **hilang** dari sistem
baru, atau sebaliknya kami membawa kolom mati selamanya.

rujukan: sisiran 1.151 `pyBrowseSQL` di 21 modul; 67 kemunculan teks `OPERATORID`, **seluruhnya**
berbentuk `OperatorID.pyXxx`; ⛔ 0 penulisan, 0 rujukan kolom. ⚠️ Lingkup **tidak** mencakup stored
procedure, pemicu, penjadwal, atau sistem luar — ronde 4 Bab D

**Jawaban:**

> **DITUTUP SEBAGAI RANGKAP — lihat P4.** `[keputusan work owner]` 2026-09-22
>
> Butir ini menanyakan hal yang sama dengan **P4**, dan ditulis pada ronde 4 tanpa mengetahui P4
> sudah ada. Jawaban yang berlaku ada di P4: kolom `OPERATORID` **dipertahankan** dan diisi dari
> **identitas akses login**.
>
> Temuan ronde 4 yang menyertainya tetap sah dan sudah pindah ke P4: dari 1.151 pernyataan SQL di
> 21 modul, **nol** menulis ke kolom itu, dan ke-67 kemunculan namanya seluruhnya berbentuk
> `OperatorID.pyXxx` — nama halaman Pega, bukan nama kolom. `[terverifikasi]`
>
> `[terbuka]` **Satu hitungan masih layak diminta, tetapi tidak memblokir**, dan ditumpangkan ke
> surat **P1** sebagai satu baris tambahan, bukan sebagai pertanyaan tersendiri:
>
> ⚠️ **Catatan 2026-09-22.** **P1 sudah terjawab** — naskah empat procedure diterima — tetapi
> **hitungan di bawah tidak ikut terjawab**. Ia tetap `[terbuka]`, tetap tidak memblokir, dan
> tercatat di lembar DBA sebagai **P49**.
>
> ```sql
> SELECT COUNT(*) total, COUNT(operatorid) terisi FROM POOLDATA.HISTORYAKSEPTASIPEGA;
> ```
>
> Nol berarti tidak ada tindak lanjut. Bukan nol berarti ada penulis **di luar Pega** — pemicu basis
> data, penjadwal, atau sistem lain — yang belum masuk peta migrasi dan wajib ditelusuri sebelum
> go-live.

---

## Sesudah Anda menjawab

⭐ Kembalikan lembar ini apa adanya — **tidak perlu dirapikan, tidak perlu lengkap**. Jawaban
sebagian tetap berguna; yang tidak berguna adalah lembar yang ditahan sampai lengkap.

---

*Disusun dari pembacaan berkas ekspor sistem lama, empat putaran. Tidak ada pertanyaan di lembar
ini yang dijawab sendiri oleh tim migrasi.*
