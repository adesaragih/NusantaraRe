# Pertanyaan untuk **IAM**
## Migrasi Treaty Inward — Realisasi & Endorsement · modul *NB Treaty In*

Anda pemegang identitas dan hak akses. Lima pertanyaan di bawah semuanya tentang **siapa boleh apa** — dan sistem lama menjawabnya dengan cara yang tidak dapat kami pindahkan apa adanya.

> ⭐ **5 pertanyaan untuk Anda**, diambil dari **40 pertanyaan** seluruh proyek modul ini. ⛔ Sisanya bukan urusan Anda dan tidak disertakan.
>
> **P11** dan **P13** menahan pembangunan wewenang di sistem baru. Tiga sisanya tentang nama orang yang tertulis langsung di dalam aturan.
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

## 1. P11 — Peran seseorang disimpan di kolom NOMOR TELEPON. Di mana peran yang sebenarnya?

⛔ Satu-satunya penentu peran yang dapat kami baca adalah **kolom nomor telepon** pada data
pengguna. Isinya bukan nomor telepon, melainkan **kode seperti `TREATY1` dan `SPVTREATY1`**.
Percabangan alur diputuskan dari kode itu.

**Konteks:** sistem baru butuh daftar peran yang sebenarnya. Kami tidak dapat menurunkannya dari
kolom yang dipakai untuk hal lain.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh daftar** — peran apa saja yang ada di proses treaty
inward, **siapa saja pemegangnya sekarang**, dan **apa arti tiap kode** yang muncul di kolom itu
*(`TREATY1`, `TREATY2`, `SPVTREATY1`, `SPVTREATY2`)*. ⚠️ Serta: **siapa yang berwenang mengubah
isi kolom itu hari ini.**

**Dampak bila salah:** ⛔ orang yang tidak berwenang **dapat menyetujui**, atau orang yang
berwenang **terkunci di luar**.

rujukan: `OperatorID.pyTelephone` — `When\IsTreaty1.xml`, `When\IsSPVTreaty1.xml` — OQ-027

**Jawaban:**

> **Ikuti rekomendasi.** `[keputusan work owner]` 2026-09-22
>
> **1. Peran menjadi medan tersendiri.** Sistem baru **tidak** menyimpan kode peran di kolom nomor
> telepon. Kolom telepon kembali hanya berisi nomor telepon. Peran disimpan di medan/tabel peran
> yang sebenarnya.
>
> **2. Dua peran dicatat:** `TREATY1` (staf treaty) dan `SPVTREATY1` (supervisor treaty).
>
> `[terverifikasi]` Kedua kode itu satu-satunya yang pernah dibandingkan terhadap
> `OperatorID.pyTelephone` di modul ini — 2 berkas, dipakai oleh `When\IsTreaty1.xml` dan
> `When\IsSPVTreaty1.xml`. Perintah audit:
> `grep -rhoE 'pyTelephone[^<]{0,60}' --include='*.xml' "NB Treaty In" | grep -oE '"[A-Za-z0-9_]{2,20}"' | sort -u`
>
> ⚠️ `[terbuka]` **Kelengkapan daftar peran ini TIDAK dapat diverifikasi dari korpus.**
> Konteks Identity & Access berstatus **absent** — nol berkas, nol rule otorisasi di seluruh
> 9.430 berkas. Jadi "dua peran" adalah **yang terbaca**, bukan terbukti **yang ada**. Bila IAM
> kelak menyebut peran ketiga, keputusan ini ditambah, bukan dibantah.
>
> ⚠️ Penetapan peran ke pengguna, pembuatan peran baru, dan pencabutannya **belum dibahas sama
> sekali** — tidak ada bahannya di korpus.

---

## 2. P13 — Siapa yang memegang pekerjaan di tiap tahap?

Alur menyebut **lima nama antrean** — Admin, Kepala Seksi, Ketua Kelompok, Kepala Departemen, dan
Direktur. ⛔ Tetapi **hubungan antara tahap dan antrean tidak tertulis** — penentuan antreannya
dilakukan secara khusus, dan kolom yang biasanya menyimpannya **kosong**.

**Konteks:** tanpa pemetaan ini, kami tidak dapat menyatakan siapa memegang giliran di tahap mana.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh tabel sederhana** — untuk setiap tahap
*("Input Realisasi", "Akseptasi Kepala Treaty", "Akseptasi Kepala Departemen")*, **antrean mana**
yang menerimanya.

**Dampak bila salah:** pekerjaan masuk ke **kotak yang salah** dan tidak ada yang merasa
bertanggung jawab.

rujukan: `<pyWorkBasket>` kosong, `<pyRouteTo>Custom`; nama antrean di `<pyRuleName>` — OQ-024

**Jawaban:**

> **Tangga tiga anak tangga.** `[keputusan work owner]` 2026-09-22
>
> ```
> ReasTreatyInAdmin  ->  ReasTreatyInSecHead  ->  ReasTreatyInDeptHead
> ```
>
> Itu seluruh tangga persetujuan realisasi treaty inward di sistem baru. **Tidak ada tangga
> keempat.**
>
> **Dua posisi DIBUANG, tidak dibangun di sistem baru:**
> - `ReasTreatyInGroupLeader`
> - `ReasTreatyInDirector`
>
> ⚠️ **Angka yang menyertai keputusan ini, supaya ia dapat ditinjau ulang kelak.**
> `[terverifikasi]` Jejak kelima posisi di `NB Treaty In`, dihitung per berkas yang menyebutnya:
>
> | Posisi | Berkas di NB Treaty In | Berkas di seluruh korpus |
> | --- | ---: | ---: |
> | `ReasTreatyInAdmin` | 5 | 41 |
> | `ReasTreatyInSecHead` | 5 | 19 |
> | `ReasTreatyInGroupLeader` | 2 | 18 |
> | `ReasTreatyInDeptHead` | 2 | 10 |
> | `ReasTreatyInDirector` | 2 | 18 |
>
> ⛔ **Kedua posisi yang dibuang punya jejak sama banyak dengan `DeptHead` yang dipertahankan**
> — dua berkas. Jadi korpus **tidak mendukung dan tidak membantah** pembuangan itu; ia sepenuhnya
> keputusan bisnis. Bila kelak ditemukan kasus berjalan yang melewati GroupLeader atau Director,
> keputusan ini wajib ditinjau ulang.
>
> ⭐ Keputusan ini menutup butir yang lahir dari grilling `Treaty In`: *"Siapa menempatkan kasus ke
> GroupLeader — jalur samping yang melompati Kepala Departemen"*. Jawabannya: jalur itu tidak
> dipakai.
>
> ⚠️ **Tidak bertabrakan dengan P5.** P5 mencabut sub-graf `Klaim -> Manajer Klaim ->
> Acceptance by Dir.` yang tanpa garis masuk. Keputusan ini mencabut posisi `ReasTreatyInDirector`
> di jalur treaty. Keduanya hal berbeda, dan keduanya kini sama-sama tidak dibangun.
>
> ⚠️ `[terbuka]` Pemetaan tiga posisi ini ke dua kode peran `TREATY1` / `SPVTREATY1` (P11)
> **belum ditentukan** — tiga posisi tidak dapat dibedakan oleh dua kode.

---

## 3. P12 — Ada nama orang tertentu tertulis di dalam aturan. Apakah itu masih benar?

Di beberapa aturan, wewenang diputuskan dengan **mencocokkan nama orang tertentu** yang tertulis
langsung di dalam aturan itu. ⛔ **Nilainya tidak kami salin ke berkas mana pun.**

**Konteks:** aturan seperti ini berhenti bekerja saat orangnya pindah jabatan atau keluar, dan
tidak seorang pun mendapat pemberitahuan.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** ganti dengan peran, sebutkan perannya ·
**(b)** memang harus orang tertentu, jelaskan alasannya · **(c)** sudah tidak relevan, hapus.

**Dampak bila salah:** ⛔ wewenang tetap menempel pada orang yang sudah **tidak berhak**.

rujukan: `pyWorkPage.pxCreateOperator` di `When\IsSPVCreate.xml`;
`OperatorID.pyUserIdentifier` di 12 berkas — OQ-021

**Jawaban:**

> **(a) Ganti dengan peran — nama orang TIDAK dimigrasi.** `[keputusan work owner]` 2026-09-22
>
> Identitas orang yang ditanam langsung di dalam aturan **tidak dibawa** ke sistem baru. Setiap
> tempat yang sekarang memeriksa "apakah pengguna ini orang tertentu" diganti menjadi "apakah
> pengguna ini memegang peran tertentu".
>
> `[terverifikasi]` Pola ini ada di **12 berkas** modul `NB Treaty In`, lewat
> `OperatorID.pyUserIdentifier` yang dibandingkan dengan nilai tetap. Sesuai CLAUDE.md §4 butir 10,
> **nilai namanya tidak disalin** ke artefak mana pun — yang dicatat hanya medan, berkas, dan
> jumlahnya.
>
> ⛔ `[terbuka]` **Keputusan ini jelas, tetapi belum dapat dijalankan sepenuhnya.** Mengganti nama
> dengan peran menuntut pemetaan **nama -> peran** untuk tiap tempat, dan pemetaan itu **tidak ada
> di korpus**. Yang dibutuhkan: untuk tiap dari 12 tempat itu, peran mana yang dimaksudkan.
> Pemiliknya `[IAM]` bersama `[work owner]`.
>
> ⚠️ Sampai pemetaan itu ada, kedua belas tempat ditandai **tertunda** di spec. ⛔ Menebak perannya
> berarti memberi atau mencabut wewenang atas dasar tebakan — dilarang.
>
> ⚠️ Terkait P11 dan P13: peran yang tersedia baru `TREATY1` dan `SPVTREATY1`, sedangkan tangga
> punya tiga posisi. Pemetaan 12 tempat ini **tidak dapat diselesaikan** sebelum ketidakcocokan itu
> diselesaikan.

---

## 4. P28 — Nama orang juga dipakai di layar, lewat medan keempat yang belum pernah kami laporkan

Selain tempat-tempat yang sudah kami tanyakan, ⛔ **nama orang tertentu juga menentukan apa yang
muncul di layar** — di **empat layar**, **dua belas tempat**, memuat **empat nama berbeda**.
⛔ **Nilainya tidak kami salin ke berkas mana pun.** Dan salah satu layar memakai **medan identitas
yang berbeda lagi** dari yang pernah kami laporkan.

**Konteks:** pada dua layar, nama yang sama dipakai **dua arah** — ada bagian yang muncul hanya
untuk orang itu, dan bagian lain yang muncul untuk **semua kecuali** orang itu. ⭐ Mengeluarkan
orang itu mengubah **dua** perilaku sekaligus.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** ganti dengan peran, sebutkan perannya ·
**(b)** memang harus orang tertentu, jelaskan alasannya · **(c)** sudah tidak relevan.
⭐ Dan: **apa yang seharusnya terjadi** pada bagian layar itu untuk orang lain.

**Dampak bila salah:** ⛔ bagian layar **hilang** bagi orang yang berhak, atau **terlihat** oleh
yang tidak berhak.

rujukan: `pyUserIdentifier` di `DetailDeptHeadTreatyIn_UW` 3× · `GeneralDeptHeadTreatyIn_UW` 3× ·
`ListSuggest` 4×; ⚠️ `pxInsName` di `DetailPoliciesNonProportional` 2× — OQ-021 diperluas,
ronde 2 baru *(pola guard keempat)*

**Jawaban:**

> **(a) Ganti dengan peran.** `[keputusan work owner]` 2026-09-22
>
> Identitas orang yang menentukan isi layar **tidak dibawa** ke sistem baru. Kedua belas tempat itu
> diganti menjadi pemeriksaan peran, sejalan dengan keputusan P12.
>
> `[terverifikasi]` Sebarannya: **4 layar · 12 tempat · 4 nama berbeda**, dan salah satu layar
> memakai **medan identitas yang berbeda** dari `OperatorID.pyUserIdentifier` yang dilaporkan di
> P12. Sesuai CLAUDE.md §4 butir 10, nilai namanya **tidak disalin** ke artefak mana pun.
>
> ⛔ `[terbuka]` **Rincian per tempat menyusul, dan ia wajib per tempat — bukan satu aturan untuk
> dua belas.** Sebabnya: pada **dua layar**, nama yang sama dipakai **dua arah** — satu bagian
> muncul **hanya** untuk orang itu, bagian lain muncul untuk **semua kecuali** orang itu.
> Mengganti keduanya dengan satu peran yang sama akan **membalik** salah satunya.
>
> Yang dibutuhkan untuk tiap dari 12 tempat: **(1)** peran penggantinya, dan **(2)** arahnya —
> bagian itu muncul untuk peran tersebut, atau muncul untuk semua kecuali peran tersebut.
>
> ⚠️ Sampai rincian itu ada, kedua belas tempat ditandai **tertunda** di spec. Menebak arahnya
> berarti menampilkan atau menyembunyikan bagian layar kepada orang yang salah.
>
> ⚠️ Terikat pada P11 dan P13: daftar peran yang tersedia belum cukup untuk membedakan tiga posisi
> tangga, sehingga pemetaan ini belum dapat diselesaikan.

---

## 5. P40 — Nama seseorang tertulis di dalam teks pemberitahuan yang dibaca pengguna

Selain nama orang yang dipakai untuk menentukan wewenang, kami menemukan **satu nama orang
tertulis di dalam isi pesan pemberitahuan** — pesan yang memberi tahu pengguna di kotak masuk
siapa sebuah pengajuan sedang berada. ⛔ **Nilainya tidak kami salin ke berkas mana pun.**
Pesan-pesan lain di tempat yang sama **mengambil namanya dari data**, bukan menuliskannya.

**Konteks:** ini berbeda dari temuan sebelumnya — di sini nama itu **bukan penjaga akses**,
melainkan **isi pesan**. Jadi bila orangnya berganti, akses tetap jalan, ⛔ **tetapi pesannya
menjadi salah.**

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** seharusnya diambil dari data seperti
pesan lainnya · **(b)** memang harus nama itu, jelaskan alasannya · **(c)** pesan ini sudah tidak
dipakai.

**Dampak bila salah:** ⛔ pengguna diberi tahu bahwa pengajuannya ada di kotak masuk **orang yang
salah**, dan mereka menghubungi orang yang salah.

rujukan: `DeptHeadTreatyIn_UW_postDT`, cabang `PositionNote=="ReasTreatyInGroupLeader"`, `SET`
ke `.NBStatus` — 1 kemunculan; bandingkan 3 cabang lain yang memakai
`@toUpperCase(pyWorkPage.pxCreateOpName)` — ronde 3 baru #7

**Jawaban:**

> **(a) Ambil dari data, samakan dengan pesan lain.** `[keputusan work owner]` 2026-09-22
>
> Nama orang **tidak ditulis di dalam teks pemberitahuan**. Ia diambil dari data pengguna atau data
> kasus, sama seperti pesan-pesan lain di tempat yang sama yang sudah melakukannya dengan benar.
>
> `[terverifikasi]` Satu pesan memuat nama ter-hardcode; pesan lain di tempat sama mengambilnya
> dari data. Nilai namanya **tidak disalin** ke artefak mana pun (CLAUDE.md §4 butir 10).
>
> ⭐ **Butir ini tidak punya sisa terbuka.** Berbeda dari P12 dan P28, ia tidak menyentuh wewenang,
> tidak butuh pemetaan nama->peran, dan polanya sudah ada contohnya di berkas yang sama. Dapat
> langsung ditulis ke spec.

---

## Sesudah Anda menjawab

⭐ Kembalikan lembar ini apa adanya — **tidak perlu dirapikan**. Jawaban setengah lengkap tetap berguna; yang tidak berguna adalah lembar yang ditahan sampai lengkap.

⚠️ **Yang tertahan selama pertanyaan ini belum terjawab:** **P11** dan **P13** menahan pembangunan wewenang di sistem baru. Tiga sisanya tentang nama orang yang tertulis langsung di dalam aturan. Bagian lain pekerjaan tetap berjalan.

---

*Disusun dari pembacaan berkas ekspor sistem lama, tiga putaran. Tidak ada pertanyaan di lembar ini yang dijawab sendiri oleh tim migrasi.*
