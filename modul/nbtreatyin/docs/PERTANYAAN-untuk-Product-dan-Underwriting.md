# Pertanyaan untuk **Product + Underwriting**
## Migrasi Treaty Inward — Realisasi & Endorsement · modul *NB Treaty In*

> ⛔⛔ **RALAT MENYELURUH ATAS BERKAS INI — P18 DITARIK.** `[penyimpangan sadar]` 2026-09-22
>
> Berkas ini memuat **5 pernyataan** yang bersandar pada premis ⛔ *"isi langkah penetapan nilai
> tidak ikut terekspor"* — di baris **541 · 623 · 791 · 793 · 795**. **Premis itu keliru.**
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


Anda pemilik aturan bisnisnya. Empat belas pertanyaan di bawah adalah keputusan yang **hanya Anda yang tahu** — kami dapat melihat apa yang sistem lama lakukan, tetapi tidak apakah itu memang yang dimaksudkan.

> ⭐ **17 pertanyaan untuk Anda** *(⭐ 3 baru dari ronde 4)*, diambil dari **48 pertanyaan** seluruh proyek modul ini. ⛔ Sisanya bukan urusan Anda dan tidak disertakan.
>
> **P20** dan **P21** menyentuh uang dan dijawab bersama Finance. **P6** dan **P5** menyentuh jenjang persetujuan. **P35** menyentuh masa berlaku kontrak.
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

## 1. P20 — Dua puluh nomor polis tertulis langsung di dalam aturan penghitung uang. Apakah itu disengaja?  ⭐ *(dijawab bersama **Finance** — pertanyaan yang sama ada di lembar mereka)*

⛔ Di dalam tiga aturan sistem lama terdapat **dua puluh nomor polis tertentu yang ditulis
langsung**, bukan dibaca dari data. ⚠️ **Dua dari tiga aturan itu menghitung nilai pertanggungan
dan premi.** Nomor-nomornya bertanggal **2023, 2024, dan 2025** — jadi ditambahkan bertahap
selama tiga tahun.

**Konteks:** artinya dua puluh penutupan itu **dihitung dengan cara berbeda** dari semua penutupan
lain, dan perbedaannya tidak tercatat di mana pun selain di dalam aturan itu.

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** memang perlakuan khusus
yang disengaja, jelaskan alasannya · **(b)** tambalan sementara untuk memperbaiki data yang salah,
sudah tidak perlu · **(c)** tidak tahu, perlu diperiksa satu per satu.
⭐ Bila (a) atau (c): ⭐ **kami perlu daftar itu dibaca bersama**, sebab sistem baru harus tahu
apakah dua puluh polis ini tetap istimewa.

**Dampak bila salah:** ⛔⛔ dua puluh penutupan berpindah ke perhitungan yang berbeda dari
sekarang, **tanpa satu pun pesan galat** yang memberi tahu.

rujukan: `Activity\SumTSIPremiSpreadedRNM_Act.xml` 24× · `Activity\ProtectFIREMBUPA_Act.xml` 12× ·
`When\IsErrorSpreading.xml` 11×; 47 kemunculan, 20 nomor unik — ronde 2 baru #2

**Jawaban:**

> **Dipindahkan ke modul Fac In — rantai ini bukan milik Treaty.** `[keputusan work owner]` 2026-09-22
>
> Kedua puluh nomor polis itu berada di dalam rantai spreading/premi yang **berasal dari Fac In**.
> Rantai itu tampak di folder NB Treaty In karena layarnya **ditanam bersama** (embed), bukan karena
> logikanya milik Treaty. Butir ini **tidak dijawab di lingkup NB Treaty In**; dipindahkan ke daftar
> pertanyaan modul Fac In.
>
> **Bukti — berkasnya identik sampai byte terakhir.** `[terverifikasi]`
>
> | Aturan | NB Treaty In | NB FacIn | Kelas pemilik |
> | --- | ---: | ---: | --- |
> | `Section\ListSuggest` | 218.982 B | 218.982 B | `ASM-FW-GISFW-Data-PolicyTreatyIn` |
> | `Activity\Protection_Act` | 343.271 B | 343.271 B | `ASM-FW-GISFW-Work` |
> | `Activity\CheckSpreadingProtect_ACT` | 306.703 B | 306.703 B | `ASM-FW-GISFW-Work` |
> | `Activity\cekSpreadingFactIn` | — | — | `ASM-FW-GISFW-Data-Coverage` |
> | `Activity\SumTSIPremiSpreadedRNM_Act` | — | — | `ASM-FW-GISFW-Data-Coverage` |
>
> Penguat: **`EDM Treaty In` memiliki `Protection_Act` sendiri** — 157.304 B, kelas
> `ASM-FW-GISFW-Data-PolicyTreatyIn`, jauh lebih kecil. Versi khusus Treaty memang ada, dan
> NB Treaty In tidak memakainya.
>
> Penguat lain: nama pemanggilnya `cekSpreadingFactIn` — "cek Spreading **FacIn**"; halaman
> `OfferFacIn` muncul 1.635 kali di dalam folder NB Treaty In; dan `SumTSIPremiSpreadedRNM_Act`
> juga dipanggil `IsThereAnyObjectLocation_Act` yang **hanya ada di tiga modul Fac**.
>
> **Sebaran nomor per tahun**, untuk dibawa ke modul Fac: 2023 — 12 · 2024 — 35 · 2025 — 72.
> 119 kemunculan, 20 nomor unik, di `SumTSIPremiSpreadedRNM_Act` (96), `ProtectFIREMBUPA_Act` (12),
> `When\IsErrorSpreading.xml` (11). Pertumbuhannya menunjukkan **kebiasaan berulang**, bukan
> tambalan sekali jalan; di modul Fac nanti pertanyaannya adalah **ciri bersama apa** yang membuat
> kedua puluh polis itu diperlakukan berbeda.
>
> `[terbuka]` **Belum dipastikan dengan pengamatan.** Tidak ada satu pun syarat di sepanjang rantai
> pemanggilan — langkah 21, 17, dan 11 seluruhnya berjalan selalu — dan titik masuknya adalah aksi
> refresh pada kelas Treaty. Secara struktur rantai itu **terjangkau** dari berkas Treaty.
> Dugaan yang menjelaskan keduanya: berkas Treaty tidak memiliki halaman `Coverage` berisi, sehingga
> rantai berjalan di atas halaman kosong dan tidak mengerjakan apa pun. **Wajib dibuktikan** dengan
> satu pengamatan di sistem berjalan sebelum 290 langkah `Data-Coverage` dikeluarkan dari lingkup.

---

## 2. P21 — Dua aturan uang hanya tertulis sebagai catatan. Apakah keduanya benar-benar berlaku?  ⭐ *(dijawab bersama **Finance** — pertanyaan yang sama ada di lembar mereka)*

Di dalam catatan pengembang pada langkah perhitungan, ada **dua aturan uang** yang dinyatakan
sebagai kewajiban: bahwa **premi polis induk harus nol** dalam keadaan tertentu, dan bahwa
**total bagian premi yang diberikan kepada pemberi bisnis harus sama dengan premi perusahaan**.
⛔ Keduanya **hanya ada sebagai catatan** — kami tidak dapat memastikan apakah sistem benar-benar
menolak ketika keduanya dilanggar.

**Konteks:** kalau ini benar-benar kewajiban, sistem baru harus menegakkannya; kalau bukan,
membangunnya akan menolak data yang sah.

**Bentuk jawaban yang diharapkan:** untuk masing-masing, pilihan ganda — **(a)** kewajiban keras,
harus ditolak bila dilanggar · **(b)** peringatan saja, boleh dilanjutkan · **(c)** bukan aturan,
hanya catatan lama. ⭐ Bila (a): **siapa yang boleh menyimpang**, kalau ada.

**Dampak bila salah:** ⛔ ketidakseimbangan premi **lolos tanpa terdeteksi**, atau sebaliknya
penutupan yang sah **ditolak**.

rujukan: `pyStepsDescription` — `set protect kalau master polis premi harus 0` ·
`TOTAL PREMI SHARE CEDANT HARUS = PREMI RNM` — ronde 2 baru #10

**Jawaban:**

> **Dipindahkan ke modul Fac In — sama seperti P20.** `[keputusan work owner]` 2026-09-22
>
> Kedua aturan uang itu berada di rantai yang sama dengan P20, dan rantai itu **tidak terjangkau
> dari NB Treaty In**. Butir ini tidak dijawab di lingkup NB Treaty In.
>
> **Bukti jangkauan, sesudah `Protection_Act` versi benar dimasukkan 2026-09-22.** `[terverifikasi]`
>
> | Aturan | Pemanggil di NB Treaty In |
> | --- | --- |
> | `ProtectShareCedant_Act` | **tidak ada — yatim** |
> | `ProtectPremiPolicy_Act` | hanya `CheckSpreadingProtect_ACT`, yang sendirinya **yatim** |
>
> Versi `Protection_Act` yang semula ada di folder ini ternyata **berkas yang salah dimasukkan** —
> versi berkelas `ASM-FW-GISFW-Work` milik Fac, 37 langkah, 10 pemanggilan. Versi Treaty yang benar
> berkelas `ASM-FW-GISFW-Data-PolicyTreatyIn`, 16 langkah, satu pemanggilan
> (`ProtectionNonProp_Act`), dan **tidak memanggil satu pun** aturan pada rantai spreading.
>
> Penguat: kedua berkas **identik byte** dengan versi di NB FacIn dan RNW Fac In —
> `ProtectPremiPolicy_Act` 207.132 B, `ProtectShareCedant_Act` 130.178 B — dan keduanya berkelas
> `ASM-FW-GISFW-Work`, bukan kelas Treaty.
>
> **Yang ikut pindah ke modul Fac In**, dua aturan uang yang hanya tertulis sebagai catatan
> pengembang:
>
> | Catatan | Di dalam |
> | --- | --- |
> | `set protect kalau master polis premi harus 0` | `ProtectPremiPolicy_Act` |
> | `TOTAL PREMI SHARE CEDANT HARUS = PREMI RNM` | `ProtectShareCedant_Act` |
>
> Pertanyaan aslinya tetap berlaku di sana: apakah keduanya **kewajiban keras** yang harus ditolak
> bila dilanggar, **peringatan saja**, atau **bukan aturan** — dan bila keras, siapa yang boleh
> menyimpang.
>
> `[terbuka]` Satu butir ikut pindah: `ProtectShareCedant_Act` juga memuat catatan
> *"set precision 4 buat cek nilai total premi share cedant"*. Presisi 4 desimal itu perlu diperiksa
> bersama ketidakseragaman 4-lawan-8 desimal yang tercatat pada **P46**, ketika modul Fac digarap.

---

## 3. P6 — Ada dua aturan berbeda dengan nama sama untuk memutuskan "sudah disetujui atau belum". Mana yang berlaku?

Gerbang **pertama** seluruh alur adalah pemeriksaan "apakah ini sudah benar/disetujui". Di sistem
lama ada **dua aturan terpisah dengan nama yang persis sama** untuk itu, dibuat **berselisih
sekitar 26 menit** pada hari yang sama, dan ⛔ **kami tidak dapat menentukan mana yang sebenarnya
dipakai.**

**Konteks:** enam titik percabangan memakai pemeriksaan ini; bila kami memilih yang salah, arah
seluruh alur bisa berubah.

**Bentuk jawaban yang diharapkan:** **butuh berkas atau demonstrasi** — satu contoh kasus nyata
beserta keterangan apakah ia lolos gerbang pertama, sehingga kami dapat mencocokkan. Atau, bila
diketahui: **aturan mana yang berlaku**, berbentuk daftar syarat.

**Dampak bila salah:** ⛔ pengajuan yang seharusnya **berhenti** malah **diteruskan**, atau
sebaliknya pengajuan sah **tertahan**.

rujukan: `When\isApproved.xml` dan `DecisionTable\isApproved.xml`, keduanya
`ASM-FW-GISFW-WORK!ISAPPROVED` — OQ-026, butir baru #7

**Jawaban:**

> **Yang dipakai adalah tabel keputusan, bukan rule `When`.** `[terverifikasi]` 2026-09-22
> Terjawab dari korpus — tidak perlu dijawab oleh Product & Underwriting.
>
> Ada dua rule bernama `isApproved` di NB Treaty In: satu `When` dan satu `DecisionTable`.
> Keduanya berbunyi **berbeda**:
>
> | Rule | Bunyi |
> | --- | --- |
> | `When\isApproved.xml` | `IsApproved = 1` |
> | `DecisionTable\isApproved.xml` | `= 0` -> `No`; selain itu -> `YES` |
>
> **Yang hidup adalah DecisionTable.** Seluruh **enam** kotak Decision di
> `Flow\InputRealizationTreatyIn.xml` — semuanya berlabel `Is Correct?` — menyambung ke
> `pxLinkedClassTo` = `Rule-Declare-DecisionTable`. Tidak satu pun menyambung ke rule `When`.
>
> Perbedaannya nyata, bukan sekadar gaya penulisan: rule `When` memperlakukan nilai kosong sebagai
> **tidak disetujui**, sedangkan DecisionTable memperlakukannya sebagai **disetujui**. Membaca rule
> yang salah akan membalik perilaku pada setiap berkas yang penandanya belum pernah diisi.
>
> Spesifikasi ditulis dari DecisionTable. Rule `When` tidak dimigrasi.
>
> Lihat juga **P24**, tempat aturan nilai ini ditetapkan, dan **P19**, tempat isi kedua tabel
> keputusan dicatat lengkap.

---

## 4. P5 — Apakah persetujuan Direktur dan jalur klaim benar-benar berjalan?

Di dalam gambar alur kerja ada **lima kotak** yang menggambarkan sebuah jalur: "Klaim" →
"Manajer Klaim" → "Persetujuan Direktur". ⛔ **Tidak ada satu pun garis yang masuk ke jalur itu.**
⚠️ Jadi entah jalur itu **sudah lama tidak dipakai**, atau **garisnya hilang** ketika sistem lama
diekspor.

**Konteks:** kalau jalur itu hidup, sistem baru harus punya jenjang persetujuan Direktur; kalau
mati, membangunnya berarti menambah jenjang yang tidak pernah ada.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** jalur ini masih dipakai, jelaskan
kapan · **(b)** sudah lama tidak dipakai · **(c)** tidak tahu, perlu diperiksa di sistem berjalan.
⭐ Bila (a): **berapa kali dipakai dalam setahun terakhir.**

**Dampak bila salah:** ⛔ satu jenjang persetujuan **hilang** dari sistem baru, atau sebuah jenjang
**dibangun dan diwajibkan** padahal tidak pernah ada.

rujukan: `Flow\InputRealizationTreatyIn.xml` — `Decision5`/`Assignment5`/`Assignment1`, sub-graf
tanpa connector masuk — OQ-023

**Jawaban:**

> **(b) Sudah lama tidak dipakai.** `[keputusan work owner]` 2026-09-22
>
> Jalur `Claim` → `Claim Manager` → `Acceptance by Dir.` **tidak dibangun** di sistem baru.
> Jenjang persetujuan Direktur **tidak ada** dalam alur realisasi treaty inward.
>
> ⚠️ **Catatan teknis yang menyertai — bukan bagian dari jawaban work owner:**
> Gerbang jalur itu, `When\isClaimTreaty.xml`, syaratnya **lengkap dan spesifik**, bukan rangka
> kosong: `ClaimPaymentType` = `"Claim"` / `"Salvage"` · `ClaimType` = `"XOL"` ·
> `Claim != 0` · `ExcessLoss != 0` · `SalvageValue != 0` · `ProportionalType != "Proportional"`.
>
> `[dugaan]` Jalur ini kemungkinan dibangun ketika klaim **non-proporsional** masih ditangani di
> dalam alur realisasi treaty, lalu pindah ke modul `Claim Non Prop` dan kotaknya ditinggalkan
> tanpa dicabut. Itu menjelaskan garis penghubung yang hilang sementara syarat gerbangnya utuh.
>
> `[terbuka]` **Belum diperiksa di sistem berjalan** apakah ada kasus realisasi treaty inward yang
> pernah mencapai `Acceptance by Dir.` dalam 12 bulan terakhir. Keputusan diambil tanpa
> pemeriksaan itu. Bila kelak ditemukan kasus semacam itu, keputusan ini **wajib ditinjau ulang**.

---

## 5. P7 — Apa arti kode-kode yang menentukan percabangan?

Beberapa percabangan alur diputuskan oleh **kode tanpa keterangan**. Yang kami temukan:
angka **1** berarti "lolos"; teks **`TREATYINDEPTHEAD`** disimpan **di kolom nomor surat** dan
dipakai sebagai penanda arah; ⛔ kode bisnis **`40`** menghentikan satu langkah penyimpanan; dan
sebuah penanda bernilai **`EDM`** atau **`POLICY`** membedakan dua cara penyimpanan.

**Konteks:** kode-kode ini menggerakkan alur, jadi artinya menentukan perilaku, bukan tampilan.

**Bentuk jawaban yang diharapkan:** satu baris per kode — **apa artinya**, dan **nilai lain apa
yang mungkin muncul**. Contoh: `40` ⇒ *"lini usaha …, dikecualikan karena …"*.

**Dampak bila salah:** ⛔ pengajuan diarahkan ke jenjang yang salah, atau **dilewatkan dari
penyimpanan** tanpa ada yang menyadarinya.

rujukan: `.IsApproved`, `LetterNo`, `.Quotation.BusinessCode != "40"`, `param.isFOR` — OQ-020

**Jawaban:**

> `[keputusan work owner]` 2026-09-22
>
> **1. `LetterNo` = `"TREATYINDEPTHEAD"` — medan dipinjam untuk routing.**
> ✅ **Dipisahkan di sistem baru.** Medan "nomor surat" yang dipakai menyimpan penanda arah adalah
> warisan, bukan aturan bisnis. Sistem baru memakai medan tersendiri untuk routing, dan medan nomor
> surat kembali hanya berisi nomor surat.
>
> **2. `isFOR` = `EDM` / `POLICY`.**
> ✅ **`EDM` = endorsemen · `POLICY` = polis baru.** Keduanya menempuh cara penyimpanan berbeda.
>
> ⚠️ **Sisa P7 dialihkan, bukan dijawab di sini.** Arti kode lini bisnis — termasuk `"40"` yang
> menghentikan langkah simpan JSON — menjadi **permintaan data ke DBA**, bukan pertanyaan bisnis:
> `SELECT ID, OLDID, NOTE, GROUPPANEL, BusinessGroupID FROM POOLDATA.BUSINESS;`
>
> `[terverifikasi]` Sembilan arti kode sudah dipulihkan dari nama rule pengujinya sendiri:
> 10021 Products Liability · 10023 Professional Liability · 10032 dan 10166 Engineering ·
> 10106 Billboard/Neon Syariah · 10157 Marine Hull Offshore · 10167 Yield Shortfall ·
> 10168 Crime · 10169 Environmental · `BusinessOldId` L1–L16 Life.
> `[terbuka]` `"40"` berada pada **skala penomoran berbeda** (dua digit, bukan lima) dan belum
> terjelaskan. `[dugaan]` ia menandai Life atau kelompok yang memuat Life, sebab langkah yang
> dihentikannya adalah simpan JSON polis treaty sementara Life punya jalur JSON sendiri.

---

## 6. P35 — Bila tanggal akhir kontrak dibiarkan kosong, sistem mengisinya dengan HARI INI. Benarkah begitu?

Ketika sebuah kontrak dibuka dan **tanggal akhirnya belum diisi**, sistem lama mengisinya dengan
**tanggal hari itu** — ⭐ **sama dengan tanggal mulainya**. Sementara di bagian lain sistem ada
perhitungan yang menetapkan masa berlaku **satu tahun**.

**Konteks:** kalau yang berlaku adalah yang pertama, maka kontrak yang tanggal akhirnya terlupakan
**tersimpan sebagai kontrak yang berakhir di hari ia dibuat.**

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang harus satu tahun, pengisian
hari-ini adalah kekeliruan · **(b)** hari-ini benar, dan pengguna wajib mengubahnya ·
**(c)** tergantung jenis kontrak, jelaskan mana.

**Dampak bila salah:** ⛔⛔ **masa berlaku kontrak salah**, dan itu menentukan premi, klaim yang
ditanggung, dan pelaporan.

rujukan: `InputPolicyTreatyIn_preDT` lgk 2 — gerbang `.EndDate==""` ⇒
`@DateTime.CurrentDate("dd/MM/yyyy","")`; bandingkan `SystemSetOneYear_DT` ⇒
`@DateTime.addCalendar(.StartDate,1,0,0,0,0,0,0)` — ronde 3 baru #4

**Jawaban:**

> **(b) Kosong → hari ini. Perilaku lama dipertahankan apa adanya.** `[keputusan work owner]` 2026-09-22
>
> Bila tanggal akhir dikosongkan, sistem baru mengisinya dengan **tanggal hari ini**, sama seperti
> sistem lama — walau hasilnya kontrak bermasa berlaku nol hari.
>
> ⚠️ **Konsekuensi yang dicatat, bukan dibantah:** migrasi data lama akan menghasilkan kontrak
> ber-`EndDate` = `StartDate`. Laporan apa pun yang menyaring kontrak "masih berlaku" akan
> mengeluarkannya sejak hari pertama. Ini perilaku yang disengaja, bukan cacat migrasi.
>
> ⚠️ `[terbuka]` Bila demikian, peran `SystemSetOneYear_DT` (`.EndDate` =
> `@DateTime.addCalendar(.StartDate,1,0,0,0,0,0,0)`) menjadi tidak jelas — kapan aturan satu-tahun
> itu berjalan, dan apakah ia menimpa pengisian hari-ini atau sebaliknya. **Belum ditentukan.**
>
> `[terverifikasi]` Kedua aturan memakai format tanggal berbeda: `dd/MM/yyyy` dan `MM/dd/yyyy`.
> Penyeragamannya dibahas terpisah di P32.

---

## 7. P36 — Ada DUA penanda persetujuan, bukan satu. Apa bedanya?

Kami menemukan **penanda persetujuan kedua** yang belum pernah kami laporkan — namanya menyebut
**persetujuan kepada Kepala Departemen**, terpisah dari penanda persetujuan umum. Keduanya
**dikosongkan bersamaan** ketika kasus masuk tahap Kepala Departemen.

**Konteks:** dua penanda berarti dua keadaan yang dapat berbeda, dan kami tidak dapat menurunkan
bedanya dari aturan.

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** dua tahap persetujuan yang
berbeda, jelaskan mana milik siapa · **(b)** salah satunya sudah tidak dipakai ·
**(c)** satu turunan dari yang lain.

**Dampak bila salah:** ⛔ persetujuan tercatat di penanda yang salah, dan kasus **maju atau tertahan
keliru**.

rujukan: `DeptHeadTreatyInUW_preDT` lgk 4–5 — `.PolicyTreatyIn.IsApproved = ""` dan
`.PolicyTreatyIn.isApprovedtoDeptHead = ""` — ronde 3 baru #6

**Jawaban:**

> **`isApprovedtoDeptHead` sudah tidak dipakai — TIDAK dimigrasi.** `[keputusan work owner]` 2026-09-22
>
> Sistem baru memakai **satu** penanda persetujuan saja: `IsApproved`. Properti
> `.PolicyTreatyIn.isApprovedtoDeptHead` dan rule bernama sama **tidak dibangun**, dan kolomnya
> tidak dibuat di skema baru.
>
> ⭐ **Keputusan ini sekaligus menutup satu butir yang korpus tidak dapat jawab sendiri.**
> `[terverifikasi]` Rule `isApprovedtoDeptHead` **dirujuk dari 6 berkas di 2 modul** —
> `NB Treaty In` dan `NB FacIn`, masing-masing di `DataTransform\DeptHeadTreatyInUW_preDT.xml`,
> `DeptHeadTreatyIn_UW_postDT.xml`, dan `InboxPolicyTreatyIn_postDT.xml` — tetapi **berkas
> rule-nya tidak ada di seluruh korpus**. Perintah audit:
> `find . -iname 'isApprovedtoDeptHead.xml' | wc -l` -> 0
>
> Sebelum jawaban ini, butir tersebut berdiri sebagai "aturan dipanggil tetapi tidak ada",
> sejajar dengan `SumTSIPremiSpreadRNMMultiCob_Act`. Kini ia **tidak perlu diminta** dalam ekspor
> ulang, sebab memang tidak dipakai lagi.
>
> ⚠️ Berlaku juga untuk `NB FacIn` bila modul itu digrilling kelak — tiga berkasnya merujuk rule
> yang sama.

---

## 8. P8 — Apa yang sebenarnya dikirim ke layanan "ARASAPAS" di akhir proses?

Langkah **terakhir** jalur realisasi memanggil sesuatu bernama "HIT SERVICE ARASAPAS".
⛔ **Isi langkah itu tidak ada di dalam ekspor yang kami terima untuk modul ini** — hanya
pembungkusnya. Salinan isinya ada di **dua modul lain**, dan ⛔ **kami sengaja tidak meminjamnya**,
sebab keduanya terdaftar berisi berbeda.

**Konteks:** ini efek keluar terakhir — bila ia mengirim data ke sistem keuangan, kegagalannya
berarti transaksi tidak tercatat di sana.

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** kirim ke sistem
faktur/keuangan, sebutkan sistemnya · **(b)** kirim ke sistem lain · **(c)** sudah tidak dipakai.
⭐ Dan: **apa yang terjadi sekarang bila pengiriman itu gagal** — dicoba ulang, atau diabaikan?

**Dampak bila salah:** ⛔ transaksi **tidak sampai** ke sistem tujuan, atau ⛔ **terkirim dua kali**.

rujukan: `Flow\InputRealizationTreatyIn.xml` `Utility2` →
`Activity\serviceInsertArasapas_act.xml` — OQ-025 *(BLOCKER)*

**Jawaban:**

> **Ikuti apa adanya.** `[keputusan work owner]` 2026-09-22
>
> Pengiriman ke layanan ARASAPAS di akhir alur realisasi **ditiru persis** seperti sistem lama —
> pemicunya, muatannya, dan penempatannya di akhir alur. Tidak ditambah, tidak dikurangi, tidak
> dipindah ke titik lain.
>
> ⛔ `[terbuka]` **Keputusan ini belum dapat dijalankan hari ini, dan sebabnya bukan kelalaian.**
> `[terverifikasi]` Berkas `Activity\serviceInsertArasapas_act.xml` di `NB Treaty In` hanya
> **18.013 B** dan berisi **satu langkah**: memanggil rule beridentitas
> `ASM-FW-GISFW-WORK / SERVICEINSERTARASAPAS_ACT`. Rule itu **tidak ada di modul ini**.
> Implementasi sebenarnya ada di dua modul lain:
>
> | Modul | Ukuran |
> | --- | ---: |
> | `EDM Treaty In` | 232.000 B |
> | `Endorsment Fac In` | 231.975 B |
>
> Aturan telusur `_METHOD.md` §1.1 butir 4 melarang meminjam varian dari modul lain, dan register
> OQ-011 entri #415 mencatat kedua varian itu **berisi berbeda**. Jadi apa yang dikirim ke ARASAPAS
> **tidak dapat dinyatakan** dari bahan yang ada sekarang.
>
> ⭐ **Yang membuka butir ini: grilling `EDM Treaty In`.** Ia memuat implementasi yang dipanggil,
> dan begitu dibaca, "apa adanya" dapat ditiru tanpa menebak. Sampai saat itu, langkah ARASAPAS
> ditandai **tertunda** di spec — bukan dihapus, bukan dikarang.
>
> ⚠️ `[terbuka]` Apa itu ARASAPAS sebagai sistem — akuntansi, keuangan, atau lainnya — **belum
> dijawab**, dan pertanyaan itu tetap berdiri walau keputusan "ikuti apa adanya" sudah diambil.

---

## 9. P37 — Dua kode angka menandai kontrak "punya penempatan keluar". Apa artinya?

Sebuah aturan memeriksa **dua kode angka jenis kontrak** dan, bila cocok, menandai kontrak itu
**punya penempatan keluar**. Kedua kode itu **hanya angka**, tanpa keterangan.

**Konteks:** penandaan ini menentukan apakah sebuah kontrak masuk jalur penempatan keluar,
dan itu jalur yang berbeda sepenuhnya.

**Bentuk jawaban yang diharapkan:** dua baris — **apa arti masing-masing kode**, dan
⭐ **apakah masih ada kode lain** yang seharusnya ikut ditandai tetapi belum.

**Dampak bila salah:** ⛔ kontrak yang seharusnya masuk jalur penempatan keluar **tidak masuk**,
atau sebaliknya.

rujukan: `TestTreatyToFacStatus` — `.TreatyType=="10015"` dan `.TreatyType=="10218"` ⇒
`Primary.PolicyTreatyIn.HasFacOut = "1"` — ronde 3 baru #10

**Jawaban:**

> **Ikuti apa adanya.** `[keputusan work owner]` 2026-09-22
>
> Kedua kode angka penanda "punya penempatan keluar" **disalin apa adanya** ke sistem baru.
> Nilainya tidak diterjemahkan, tidak dinormalkan, dan tidak diganti dengan enumerasi bernama.
> Perbandingan di kode Go memakai **nilai literal yang sama persis** seperti Pega.
>
> **Yang berlaku sebagai akibatnya, dan ini disengaja:**
> - tidak ada validasi daftar nilai — nilai di luar yang dikenal tetap diterima dan disimpan;
> - tidak ada aturan bisnis baru yang dibangun di atas arti kode ini;
> - tampilan menampilkan kodenya, bukan keterangannya, kecuali master datanya tersedia.
>
> `[terbuka]` **Arti kedua kode tetap tidak diketahui.** Butir ini **tidak ditutup** — ia hanya
> tidak lagi memblokir penulisan spec, sebab spec meniru perilaku tanpa perlu tahu artinya.
> Bila kelak master datanya diterima dari DBA, arti kode dapat diisi tanpa mengubah perilaku.

---

## 10. P10 — Apakah semua pekerjaan masuk ke antrean BERSAMA, tidak pernah ke antrean pribadi?

Keenam titik penugasan di alur ini memakai **antrean bersama** — ⛔ **tidak satu pun** menugaskan
ke **kotak masuk pribadi** seseorang.

**Konteks:** ini menentukan bentuk kotak masuk di sistem baru: daftar bersama yang siapa pun dalam
peran itu boleh ambil, atau tugas milik satu orang.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** benar, semua antrean bersama ·
**(b)** seharusnya ada yang pribadi, sebutkan tahap mana · **(c)** dulu pribadi, diubah.

**Dampak bila salah:** pekerjaan **menumpuk tanpa pemilik**, atau sebaliknya **terkunci pada satu
orang** yang sedang tidak ada.

rujukan: keenam Assignment `<pyImplementation>WorkBasket` / `<pyRouteTo>Custom` — OQ-028

**Jawaban:**

> **(a) Benar, semua antrean bersama.** `[keputusan work owner]` 2026-09-22
>
> Antrean bersama **dipertahankan** di sistem baru. Tidak ada kotak masuk pribadi; pekerjaan
> menunggu siapa pun yang memegang posisi itu, bukan orang tertentu.
>
> `[terverifikasi]` Bukti korpus bulat tanpa kekecualian: keenam Assignment di
> `Flow\InputRealizationTreatyIn.xml` memakai `ToWorkBasket`, dan **nol** memakai `ToWorklist`.
> Perintah audit:
> `grep -ohE 'ToWorkBasket|ToWorklist' "NB Treaty In/Flow/InputRealizationTreatyIn.xml" | sort | uniq -c`

---

## 11. P22 — Sembilan puluh satu bagian layar dibuat agar tidak pernah muncul, dan tiga puluh delapan dikunci permanen. Masih perlu?

⛔ **Sembilan puluh satu** bagian layar diberi syarat tampil yang **tidak mungkin pernah benar** —
setara dengan menuliskan "tampilkan jika 1 sama dengan 2". Dan ⛔ **tiga puluh delapan** medan
diberi syarat **hanya-baca yang selalu benar**, sehingga tidak pernah dapat disunting.
⭐ **Seluruh tiga puluh delapan itu ada di layar Kepala Departemen.**

**Konteks:** kami perlu tahu mana yang memang tidak boleh dipakai lagi, dan mana yang dimatikan
sementara lalu terlupakan — sebab yang pertama tidak perlu dibangun, yang kedua perlu.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang sudah tidak dipakai, boleh
hilang · **(b)** dimatikan sementara, seharusnya kembali · **(c)** campuran, perlu ditinjau
bersama. ⭐ Untuk layar Kepala Departemen: **apakah jabatan itu memang hanya boleh melihat, tidak
mengubah?**

**Dampak bila salah:** ⛔ medan mati ikut dibangun dan dirawat selamanya, atau medan yang masih
dibutuhkan **hilang** dari sistem baru.

rujukan: `1=2` 72× · `NEVER` 15× · `1==2` 2×; `pyReadOnlyCondition` `1==1` 32× + varian 6×,
seluruhnya di `DetailDeptHeadTreatyIn_UW` dan `GeneralDeptHeadTreatyIn_UW` — ronde 2 baru #6

**Jawaban:**

> **DITUTUP — terwakili sepenuhnya oleh P44 dan P45.** `[keputusan work owner]` 2026-09-22
>
> Butir ini adalah induk dari dua pertanyaan ronde 4, dan keduanya sudah diputuskan:
>
> | Bagian P22 | Diputuskan di |
> | --- | --- |
> | 91 bagian layar yang dimatikan | **P44** — 80 dibuang tanpa ditanyakan; 4 wadah berisi 104 medan dibawa ke Product & Underwriting |
> | 38 medan yang dikunci permanen | **P45** — ikuti apa adanya, seluruhnya tetap terkunci |
>
> ---
>
> **Koreksi atas bunyi pertanyaan.** `[terverifikasi]` P22 menyatakan *"seluruh tiga puluh delapan
> ada di layar Kepala Departemen"*. Yang benar: **36 di layar Kepala Departemen, 2 di layar biasa** —
> satu di `Section\DetailPolicyTreatyIn`, satu di `Section\GeneralPolicyTreatyIn`. Penguncian itu
> karenanya bukan semata-mata sifat jabatan.
>
> ---
>
> **Pertanyaan turunan — "apakah Kepala Departemen hanya boleh melihat, tidak mengubah?"**
>
> > **Tidak.** Kepala Departemen mengisi **tujuh medan** sendiri, lalu memutuskan.
>
> `[terverifikasi]` Disilangkan dari daftar medan wajib (**P47**) dan daftar medan terkunci (**P45**)
> pada layar yang sama:
>
> | | Jumlah |
> | --- | ---: |
> | medan wajib diisi | 16 |
> | di antaranya terkunci | 9 |
> | **wajib dan dapat diketik** | **7** |
>
> Ketujuhnya: `Claim` · `Deduction1` · `EndDate` · `PremiOgp` · `PremiOnp` · `StartDate` ·
> `StatementDate`.
>
> Sembilan medan yang terkunci terisi sendiri oleh perhitungan — rumusnya berada di dalam langkah
> yang diminta pada **P18**. Selama itu belum ada, ketujuh medan ini dapat diisi tetapi layar
> **tetap tidak dapat disimpan**, karena sembilan medan wajib lainnya kosong.
>
> **Untuk perancangan peran di sistem baru:** Kepala Departemen adalah peran **penyunting sekaligus
> pemutus**, bukan pengamat. Tujuh medan itu wewenangnya, dan tidak boleh dijadikan hanya-baca.

---

## 12. P9 — Mengapa modul realisasi treaty masuk menyentuh data treaty KELUAR?

Modul yang menangani **treaty masuk** ternyata menyentuh tabel **treaty keluar** di **delapan
berkas** — sementara modul yang **namanya** treaty keluar tidak menyentuhnya sama sekali.

**Konteks:** kami perlu tahu apakah ini memang alur bisnis *(misalnya bagian treaty masuk
diteruskan keluar)*, atau peninggalan sejarah.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang satu proses, jelaskan ·
**(b)** peninggalan, tidak dipakai lagi · **(c)** dipakai untuk laporan saja.

**Dampak bila salah:** tanggung jawab treaty keluar **dibangun di tempat yang salah**, dan
perubahan di satu sisi merusak sisi lain.

rujukan: `POOLDATA.M_TREATY_OUT`, class `…INT-TREATYOUTDETAIL` — OQ-022

**Jawaban:**

> **(a) Benar, ini jalur retrosesi.** `[keputusan work owner]` 2026-09-22
>
> Saat Nusantara Re menerima treaty masuk, sebagian ditempatkan kembali keluar (retro), dan layar
> inward perlu menampilkan penempatan keluar itu. **Bukan pelanggaran batas modul.**
>
> Di sistem baru: data treaty outward **tetap dibaca** dari konteks realisasi treaty inward, dan
> **tetap tidak boleh ditulis** dari sini.
>
> `[terverifikasi]` 16 berkas menyebut TreatyOut; **nol** operasi tulis ke tabel treaty out.
> `RDBList\BrowseTreatyOut.xml` dan `RDBList\BrowseTreatyOutDetail.xml` keduanya `SELECT`.
> Penamaannya konsisten retro: `Activity\SetValueRetro_Act.xml`,
> `Section\BusinessAndSOBListRetro.xml`, `Section\DetailPolicyTreatyOutNonProportional.xml`.
>
> `[terbuka]` `Activity\InsertToTreatyOutXOLList.xml` (409.888 B) **bernama** menulis, tetapi tidak
> ada perintah tulis ke tabel treaty out di seluruh modul. `[dugaan]` ia mengisi daftar di layar,
> bukan tabel — belum diverifikasi, dan tidak ditebak lebih jauh.

---

## 13. P38 — Sebuah aturan lini jiwa hidup di dalam modul kontrak reasuransi masuk. Disengaja?

Ada satu aturan yang **tergolong lini asuransi jiwa** tetapi tinggal di dalam modul kontrak
reasuransi masuk, dan ia menulis penanda **"jiwa"** ke catatan lampiran **baris pertama**.

**Konteks:** kalau modul ini tidak menangani lini jiwa, aturan itu tidak seharusnya ada; kalau
menangani, batas antar-modul perlu kami perbaiki.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** modul ini memang juga menangani lini
jiwa, jelaskan kapan · **(b)** peninggalan, tidak dipakai · **(c)** tidak tahu.

**Dampak bila salah:** ⚠️ lampiran tergolong ke lini yang salah, dan penggolongan itu ikut
terbawa ke laporan.

rujukan: `setCategoryAttachment_DT`, class `ASM-FW-GISFW-Work-LIFE`, ⇒
`Attachment.pxResults(1).NOTE = "LIFE"` — ronde 3 baru #13

**Jawaban:**

> **Ikuti apa adanya — `IsLife` dan kedua penjaga dipertahankan.** `[keputusan work owner]` 2026-09-22
>
> ⛔ **RALAT dalam sesi yang sama.** Jawaban pertama atas P38 adalah **"(b) sisa lama, tidak
> dipakai — jangan dimigrasi"**. Kalimat itu dikutip utuh, tidak dihapus. Ia diberikan **sebelum**
> bukti rujukan muncul, dan dicabut setelah bukti berikut disampaikan:
>
> `[terverifikasi]` `When\IsLife.xml` **bukan rule menganggur**. Ia dirujuk dua aktivitas hidup di
> jalur uang: `Activity\ProtectPremiPolicy_Act.xml` (penjaga premi) dan
> `Activity\ProtectShareCedant_Act.xml` (penjaga share cedant). Perintah audit:
> `grep -rli 'IsLife' --include='*.xml' "NB Treaty In"` -> 3 berkas.
>
> **Keputusan akhir:** `IsLife` **dimigrasi apa adanya**, beserta perannya sebagai syarat di dalam
> kedua penjaga. Enam belas nilai `.Quotation.BusinessOldId` (`L1`…`L16`) disalin persis; tidak
> ditambah, tidak dikurangi, tidak dinormalkan.
>
> ⛔ `[terbuka]` **Keputusan ini belum dapat dijalankan sepenuhnya hari ini.** Apa yang sebenarnya
> dilakukan cabang `IsLife` di dalam kedua penjaga berada di dalam langkah `Property-Set` yang
> **isinya tidak terekspor** (268 langkah terjangkau & diminta di modul ini, dari 938 di folder; 25.986 di seluruh korpus). Jadi "apa
> adanya" baru dapat ditiru setelah ekspor ulang diterima — lihat P18.
> Sampai saat itu, kedua penjaga ditandai **tertunda** di spec, bukan ditulis dari tebakan.

---

## 14. P39 — Satu mata uang dikecualikan secara khusus dari dua daftar pilihan. Masih perlu?

Dua daftar pilihan mata uang **mengecualikan satu kode mata uang tertentu** secara khusus, ditulis
langsung di dalam aturan.

**Konteks:** pengecualian yang ditulis langsung tidak dapat diubah tanpa mengubah aturan, dan kami
perlu tahu apakah ia masih dimaksudkan.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** masih perlu dikecualikan, jelaskan
alasannya · **(b)** sudah tidak perlu · **(c)** seharusnya ada mata uang lain yang juga
dikecualikan, sebutkan.

**Dampak bila salah:** ⚠️ mata uang yang masih dipakai **hilang dari pilihan**, atau yang sudah
mati **muncul kembali**.

rujukan: penyaring bernilai `"ITL"` di `BrowseCurrency_RD` dan `BrowseCurrencyTreatyIn_RD` —
ronde 3 baru #15

**Jawaban:**

> **(a) Ikuti apa adanya — pengecualian dipertahankan.** `[keputusan work owner]` 2026-09-22
>
> Kedua daftar pilihan mata uang di sistem baru tetap mengecualikan kode yang sama, disalin persis.
>
> `[terverifikasi]` **Mata uang yang dikecualikan adalah `ITL`.** Penyaring
> `<pyFilterValue>"ITL"</pyFilterValue>` ada di `ReportDefinition\BrowseCurrency_RD.xml` dan
> `ReportDefinition\BrowseCurrencyTreatyIn_RD.xml`, masing-masing dua kali. Perintah audit:
> `grep -oE '<pyFilterValue>"ITL"</pyFilterValue>' "NB Treaty In/ReportDefinition/BrowseCurrency_RD.xml"`
>
> ⚠️ **Konteks yang menyertai, dan ia mengubah cara membaca butir ini:** `ITL` adalah **Lira
> Italia**, mata uang yang digantikan Euro sejak 2002. Jadi pengecualian ini `[dugaan]` bukan
> kebijakan treaty maupun aturan regulator, melainkan **pembersihan mata uang mati** dari daftar
> pilihan. Mempertahankannya tidak berbiaya dan tidak berisiko.
>
> ⭐ `[terbuka]` **Pertanyaan susulan yang lahir dari sini, dan belum ditanyakan ke siapa pun:**
> bila satu mata uang mati disembunyikan dengan penyaring ter-hardcode, apakah ada mata uang mati
> lain di tabel `CURRENCY` yang **tidak** disembunyikan dan masih muncul di pilihan? Cara sehatnya
> di sistem baru adalah penanda aktif/non-aktif di master, bukan penyaring ter-hardcode per layar —
> tetapi itu **usulan desain**, bukan keputusan yang sudah diambil.

---

## 15. P44 — Sembilan puluh satu bagian layar dimatikan, tetapi delapan puluh di antaranya kosong. Yang mana yang masih perlu?  ⭐ *(BARU — ronde 4)*


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

## 16. P45 — Tiga puluh delapan medan dikunci permanen. Haruskah tetap terkunci?  ⭐ *(BARU — ronde 4)*


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

## 17. P46 — Sebuah medan terisi otomatis dengan kata "Inclusive". Apa artinya, dan kapan berlaku?  ⭐ *(BARU — ronde 4)*


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
>
> ---
>
> ⭐⭐ **UJI SILANG DARI DATA PRODUKSI — 2026-09-22.** `[terverifikasi]` `[data DBA]`
>
> Rumus ini semula hanya terbaca dari XML. Ia kini **terbukti dari satu kasus produksi sungguhan**,
> lewat contoh `DATA_JSON` kedua yang diterima work owner *(kasus 2026, MARINE CARGO)*:
>
> ```
> TypeTax                = "Inclusive"
> Deduction1             = 32.516
> BrokerageFeeSebenarnya = 31.8160        <- nilai yang benar-benar tersimpan
>
> 32.516 / 1.022 = 31.81604696673189...
>                  dibulatkan 4 desimal = 31.8160   ✓ COCOK
> ```
>
> ⭐ **Potongan 2,2 % bukan lagi pembacaan rumus, melainkan perilaku terukur.**
>
> ⚠️ **Dan butir `[terbuka]` presisi di atas ikut menajam, tanpa tertutup.** Nilai tersimpan itu
> berhenti di **4 desimal**, bukan 8 — artinya kasus ini melewati cabang **empat desimal**
> *(`InputPolicyTreatyInDetail_NonProp` / `InputPolicyTreatyOutDetail_NonProp`)*, bukan cabang
> delapan desimal. Pada 8 desimal nilainya akan `31.81604697`.
>
> ⛔ **Butir presisi TETAP `[terbuka]`.** Satu kasus menunjukkan **cabang mana yang dilewatinya**,
> bukan **cabang mana yang seharusnya**. Pertanyaan asalnya — apakah dua presisi untuk rumus yang
> sama dikehendaki — masih milik `[Finance]` dan `[Product+Underwriting]`.
>
> ⛔ Contoh produksinya **tidak disimpan** di berkas mana pun; ia memuat nama orang pada medan
> pemasar, operator, dan tertanggung. Yang dikutip di atas hanya **tiga angka** dari medan uang.

---

## Sesudah Anda menjawab

⭐ Kembalikan lembar ini apa adanya — **tidak perlu dirapikan**. Jawaban setengah lengkap tetap berguna; yang tidak berguna adalah lembar yang ditahan sampai lengkap.

⚠️ **Yang tertahan selama pertanyaan ini belum terjawab:** **P20** dan **P21** menyentuh uang dan dijawab bersama Finance. **P6** dan **P5** menyentuh jenjang persetujuan. **P35** menyentuh masa berlaku kontrak. Bagian lain pekerjaan tetap berjalan.

---

*Disusun dari pembacaan berkas ekspor sistem lama, tiga putaran. Tidak ada pertanyaan di lembar ini yang dijawab sendiri oleh tim migrasi.*
