# Pertanyaan untuk **Finance**
## Migrasi Treaty Inward — Realisasi & Endorsement · modul *NB Treaty In*

> ⛔⛔ **RALAT MENYELURUH ATAS BERKAS INI — P18 DITARIK.** `[penyimpangan sadar]` 2026-09-22
>
> Berkas ini memuat **1 pernyataan** yang bersandar pada premis ⛔ *"isi langkah penetapan nilai
> tidak ikut terekspor"* — di baris **192**. **Premis itu keliru.**
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


Dua pertanyaan di bawah menyentuh **angka uang**, dan keduanya dijawab bersama Product + Underwriting. Keduanya juga ada di lembar mereka — Anda tidak perlu menyalin apa pun, cukup menjawab bagian yang menyentuh nilai.

> ⭐ **3 pertanyaan untuk Anda** *(⭐ 1 baru dari ronde 4)*, diambil dari **48 pertanyaan** seluruh proyek modul ini. ⛔ Sisanya bukan urusan Anda dan tidak disertakan.
>
> ⛔ **Keduanya menahan pemeriksaan paritas perhitungan.** Tanpa jawabannya, kami tidak dapat membuktikan bahwa sistem baru menghitung sama dengan sistem lama.
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

## 1. P20 — Dua puluh nomor polis tertulis langsung di dalam aturan penghitung uang. Apakah itu disengaja?  ⭐ *(dijawab bersama **Product + Underwriting** — pertanyaan yang sama ada di lembar mereka)*

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

## 2. P21 — Dua aturan uang hanya tertulis sebagai catatan. Apakah keduanya benar-benar berlaku?  ⭐ *(dijawab bersama **Product + Underwriting** — pertanyaan yang sama ada di lembar mereka)*

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

## 3. P43 — Delapan angka uang tampil di layar tanpa ada perhitungan yang menghasilkannya. Apakah angkanya benar?  ⭐ *(BARU — ronde 4)*


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

## Sesudah Anda menjawab

⭐ Kembalikan lembar ini apa adanya — **tidak perlu dirapikan**. Jawaban setengah lengkap tetap berguna; yang tidak berguna adalah lembar yang ditahan sampai lengkap.

⚠️ **Yang tertahan selama pertanyaan ini belum terjawab:** ⛔ **Keduanya menahan pemeriksaan paritas perhitungan.** Tanpa jawabannya, kami tidak dapat membuktikan bahwa sistem baru menghitung sama dengan sistem lama. Bagian lain pekerjaan tetap berjalan.

---

*Disusun dari pembacaan berkas ekspor sistem lama, tiga putaran. Tidak ada pertanyaan di lembar ini yang dijawab sendiri oleh tim migrasi.*
