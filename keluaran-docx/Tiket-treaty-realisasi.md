# Tiket - Realisasi dan Endorsemen Treaty

> Dokumen ini memuat **badan tiket lengkap**, disusun per modul lalu per nomor.
> Disusun 25 September 2026 dari berkas tiket proyek migrasi Nusantara Re.

## Matriks status

| Modul | Tiket | Siap | Tertahan | needs-info | wontfix | Lain |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| NB Treaty In | **28** | 20 | 7 | 1 | 0 | 0 |
| EDM Treaty In | **12** | 4 | 8 | 0 | 0 | 0 |
| **Jumlah** | **40** | **24** | **15** | **1** | **0** | **0** |

### Tiket tertahan dan gerbangnya - 16 tiket

| Modul | # | Judul | Tertahan oleh |
| --- | ---: | --- | --- |
| NB Treaty In | 05 | Peran menggantikan nama orang — tertunda sampai pemetaan diterima | pemetaan **nama → peran** dari `[IAM]` dan `[work owner]` |
| NB Treaty In | 13 | Rantai perhitungan uang — karantina P18 DIBUKA, sisa dua penahan lain | - |
| NB Treaty In | 18 | Tipe kolom, konversi masuk, dan presisi uang | - |
| NB Treaty In | 19 | Pemecah dokumen menjadi baris | - |
| NB Treaty In | 22 | Pemuat dokumen lama | - |
| NB Treaty In | 26 | Proyeksi selisih yang dapat dibaca langsung | **20** · **25** · ⛔ `[work owner]` **anak proyeksi dibutuhkan pembaca SQL atau t |
| NB Treaty In | 27 | Perbedaan endorsemen pada tabel yang sama | **19** · **24** · ⛔ `[work owner]` **beda dagang tabel pecahan penyebaran** · ⛔  |
| NB Treaty In | 28 | Penanda migrasi endorsemen dan kepatuhan lapisan | **22** · **26** · ⛔ `[work owner]` **lingkup pemindahan dokumen lama** · ⛔ `[wor |
| EDM Treaty In | 02 | Rantai generasi dan larangan percabangan | - |
| EDM Treaty In | 05 | Pembatalan sebagai generasi bernilai nol | - |
| EDM Treaty In | 07 | Proyeksi selisih yang dapat dibaca langsung | - |
| EDM Treaty In | 08 | Perbedaan perhitungan sebaran dan rincian angsuran | - |
| EDM Treaty In | 09 | Dua penanda migrasi | - |
| EDM Treaty In | 10 | Pemuat migrasi endorsemen | - |
| EDM Treaty In | 11 | Kepatuhan lapisan dan tipe kolom | - |
| EDM Treaty In | 12 | Medan catatan dan penguncian daftar kolom | - |

> **Catatan.** Tiket NB Treaty In **24-28** ditandai **digantikan** oleh sebelas tiket
> penyimpanan EDM Treaty In 01-11, yang lebih halus dan menyebut gerbang tiket NB-nya.
> Keduanya dimuat di sini; keputusan mana yang berlaku milik pemilik pekerjaan.

---

# NB Treaty In

Jumlah tiket: **28**

## NB Treaty In - 01 - Sumber data realisasi treaty — dibaca dari view relasional, gagal baca menghentikan proses

**Status:** ready-for-agent
**Blocked by:** —
**Menutup:** AC 15 · 16 · 17 · 36 · 37 · 38 · 57 · 58 · 89 *(9 AC)* — US 21 · 23 · 24 · 37

#### Hasil & nilai pengguna

Hari ini data kontrak dibaca dengan **membongkar satu dokumen teks** menjadi properti saat
halaman dibuka. `[terverifikasi]` Enam langkah Java identik melakukannya, dan bila pembongkarannya
**gagal**, galatnya **hanya ditulis ke log** — aktivitas **tetap lanjut** dengan halaman kosong
atau separuh terisi. ⛔ Pengguna tidak diberi tahu apa pun.

Sesudah tiket ini, data kontrak dibaca dari **sumber relasional yang bentuknya dapat dinyatakan**,
dan ⭐ bila pembacaan gagal, pengguna **diberi tahu** dan prosesnya **berhenti** — tidak ada lagi
berkas yang tersimpan dari pembacaan yang gagal.

#### Area codebase

- Lapisan repository: pembacaan data realisasi treaty
- Lapisan repository: pembacaan penempatan keluar *(hanya baca)*
- Lapisan service: penanganan kegagalan pembacaan

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pembongkaran dokumen — **tidak dimigrasi** | 6 langkah Java di `Activity\FetchMasterTreatyIn`, `SetTreatyIn_Act`, `InputPolicyTreatyInDetail_*`, `InputPolicyTreatyOutDetail_*` |
| Medan yang dipakai laporan | `ReportDefinition\BrowseTreatyInDetail.xml` — **33 medan** |
| Penempatan keluar, hanya baca | `RDBList\BrowseTreatyOut.xml` · `RDBList\BrowseTreatyOutDetail.xml` — keduanya `SELECT` |

#### ADR terkait

- **ADR-U-0009** — migrasi penuh, tidak ada koeksistensi dua penulis

#### Acceptance criteria

- [ ] **AC 15** — data dibaca dari sumber relasional, bukan dari dokumen
- [ ] **AC 16** — sistem baru **tidak menulis** dokumen
- [ ] **AC 17** — ke-**33** medan yang dipakai laporan tersedia
- [ ] **AC 89** — nol medan yang dipakai tetapi tidak tersedia
- [ ] **AC 36** — kegagalan pembacaan **menghentikan** proses
- [ ] **AC 37** — kegagalan pembacaan **menampilkan galat kepada pengguna**
- [ ] **AC 38** — nol kasus tersimpan dari pembacaan yang gagal
- [ ] **AC 57** — penempatan keluar **dapat dibaca** dari konteks ini
- [ ] **AC 58** — penempatan keluar **tidak pernah ditulis** dari sini

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **P29** | contoh isi dokumen — hanya untuk **migrasi data lama**, bukan untuk pembacaan baru | ⭐ **tidak menahan** |
| **19** | satu aturan **bernama** menulis penempatan keluar tetapi tidak ada perintah tulisnya | tidak menahan |

#### Perintah verifikasi

1. Buka sebuah realisasi treaty — ⭐ seluruh **33** medan laporan terisi.
2. Putus sumber data di tengah pembacaan — ⭐ pengguna **melihat galat**, ⭐ dan **nol** berkas
   tersimpan.
3. Coba menulis penempatan keluar dari konteks ini — ⭐ **ditolak**.

#### Catatan

⚠️ `[penyimpangan sadar]` Menghentikan proses saat gagal baca **berbeda** dari perilaku Pega.
Alasannya tertulis di `spec.md` §5.9: **kegagalan yang terlihat lebih murah daripada yang
tersembunyi**.

## NB Treaty In - 02 - Daur hidup berkas realisasi dan nomor polis — satu nomor, sekali, tanpa bentrok

**Status:** ready-for-agent
**Blocked by:** 01
**Menutup:** AC 31 · 59 · 73 · 74 *(4 AC)* — US 1 · 4 · 5

#### Hasil & nilai pengguna

Hari ini admin treaty membuka penawaran yang sudah disetujui dan melengkapinya menjadi realisasi.
`[terverifikasi]` Nomor polis dibentuk dari sebuah deret, dan ⚠️ salah satu bagiannya diambil dari
slot parameter generik yang **tidak terlihat diisi** di aktivitas mana pun.

Sesudah tiket ini, sebuah realisasi treaty **lahir dari penawaran yang disetujui**, mendapat
**satu nomor polis yang tidak bentrok**, dan ⭐ pengguna **diperingatkan** ketika ia hendak membuat
penawaran yang sudah pernah ada.

#### Area codebase

- Lapisan service: daur hidup berkas realisasi
- Lapisan service: pembentukan nomor polis
- Lapisan handler: peringatan penawaran ganda

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Titik masuk alur | `Flow\InputRealizationTreatyIn.xml` — `pyStartActivity` = `Start1` |
| Pembentukan nomor polis | `RDBList\GenerateNoPolicy.xml` |
| Peringatan penawaran ganda | `Activity\CheckDuplicateOffer.xml` langkah 6 dan 8, lewat `RDBList\GetCountClaim.xml` |
| Rantai pemanggilnya | `TreatyRealizationCheckXOLList` → `SetTreatyIn_Act` → `CheckDuplicateOffer` |

#### ADR terkait

- **ADR-U-0006** — penomoran lewat deret basis data

#### Acceptance criteria

- [ ] **AC 73** — nomor polis memuat awalan tetap, penanda treaty, bulan-tahun, dan nomor urut
      berdigit tetap
- [ ] **AC 74** — nomor dibentuk **sekali** per berkas; tidak berubah pada penyimpanan berikutnya
- [ ] **AC 31** — dua berkas **tidak pernah** bernomor sama
- [ ] **AC 59** — peringatan muncul ketika jumlah berkas klaim terhubung **lebih dari nol**

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **P7** | asal salah satu bagian nomor polis **tidak terlihat diisi** | ⚠️ tidak menahan — polanya diketahui |

#### Perintah verifikasi

1. Realisasikan dua penawaran berbarengan — ⭐ nomor polisnya **berbeda**.
2. Simpan ulang salah satunya — ⭐ nomornya **tidak berubah**.
3. Buat penawaran yang sudah pernah ada — ⭐ **peringatan muncul**.

## NB Treaty In - 03 - Tangga tiga jenjang dan empat cabang putusan — nilai kosong berarti DISETUJUI

**Status:** ready-for-agent
**Blocked by:** 02
**Menutup:** AC 1 · 2 · 3 · 4 · 5 · 6 · 7 · 8 · 9 · 10 · 84 *(11 AC)* — US 7–14

#### Hasil & nilai pengguna

Hari ini persetujuan realisasi treaty melewati **lima posisi**, dan dua di antaranya —
Ketua Kelompok dan Direktur — ⚠️ sudah lama tidak dipakai. ⛔ Lebih rawan lagi: ada **dua aturan
berbeda bernama sama** untuk memutuskan "sudah disetujui atau belum", dan keduanya **berlawanan**
pada nilai kosong.

Sesudah tiket ini, persetujuan berjalan lewat **tiga jenjang** yang jelas, dan ⭐ setiap penolakan
punya akibat yang pasti: **admin menolak ⇒ berkas selesai sebagai ditolak; atasan menolak ⇒ berkas
kembali ke admin.**

#### Area codebase

- Lapisan service: mesin tahap persetujuan
- Lapisan service: penilaian penanda persetujuan

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| ⭐ Aturan **yang hidup** | `DecisionTable\isApproved.xml` — `= 0` ⇒ `No`, bawaan ⇒ `YES`; kolomnya bertipe **teks** |
| ⛔ Aturan yang **tidak dipakai** | `When\isApproved.xml` — menguji `= 1`; **nol** kotak Decision menyambung ke sana |
| Kotak putusan | `Flow\InputRealizationTreatyIn.xml` — **enam** Decision, seluruhnya menyambung ke `Rule-Declare-DecisionTable` |
| Perpindahan tahap | `DataTransform\InboxPolicyTreatyIn_postDT.xml` · `DeptHeadTreatyIn_UW_postDT.xml` |

#### ADR terkait

- **ADR-U-0007** — jejak audit setiap transisi dan setiap jalur balik

#### Acceptance criteria

- [ ] **AC 1** — `0` diperlakukan **ditolak**
- [ ] **AC 2** — nilai apa pun selain `0`, ⭐ **termasuk kosong**, diperlakukan **disetujui**
- [ ] **AC 3** — perbandingan dilakukan sebagai **teks**
- [ ] **AC 4** — aturan berasal dari **tabel keputusan**, bukan dari aturan bernama sama yang menguji `= 1`
- [ ] **AC 5** — admin menolak ⇒ berkas **diselesaikan sebagai ditolak**
- [ ] **AC 6** — atasan menolak ⇒ berkas **kembali ke admin**
- [ ] **AC 7** — admin menyetujui ⇒ naik ke jenjang kedua
- [ ] **AC 8** — jenjang kedua menyetujui ⇒ naik ke jenjang ketiga
- [ ] **AC 9** — jenjang ketiga menyetujui ⇒ realisasi **selesai**; ⛔ tidak ada jenjang keempat
- [ ] **AC 10** — dua posisi yang dibuang **tidak ada**; berkas tidak pernah dirutekan ke sana
- [ ] **AC 84** — nilai kosong pada penanda **tidak** menghentikan alur

#### Perintah verifikasi

1. Ajukan berkas yang penandanya **belum pernah diisi** — ⭐ ia diperlakukan **disetujui**, bukan
   ditolak, dan **tidak** melempar galat.
2. Tolak sebagai admin — ⭐ berkas **selesai**.
3. Tolak sebagai atasan — ⭐ berkas **kembali ke admin**, bukan selesai.
4. Setujui tiga kali berturut-turut — ⭐ realisasi **selesai** di jenjang ketiga.

#### Catatan

⚠️ **Perbedaan dua aturan itu nyata, bukan gaya penulisan.** Aturan yang tidak dipakai
memperlakukan nilai kosong sebagai **tidak disetujui**; yang hidup memperlakukannya **disetujui**.
⛔ Membangun dari yang salah **membalik perilaku** pada setiap berkas yang penandanya belum pernah
diisi.

## NB Treaty In - 04 - Antrean bersama dan pemeriksaan keanggotaan — bukan nomor urut daftar

**Status:** ready-for-agent
**Blocked by:** 03
**Menutup:** AC 11 · 14 · 92 *(3 AC)* — US 3 · 8 · 19

#### Hasil & nilai pengguna

Hari ini pekerjaan menunggu di **antrean bersama**, dan itu memang yang diinginkan. ⛔ Tetapi
wewenang di beberapa tempat ditentukan dengan bertanya *"antrean nomor dua Anda namanya apa"* —
`[terverifikasi]` menunjuk antrean **menurut posisi dalam daftar**, bukan menurut namanya.
⚠️ Menambah seorang pengguna ke antrean baru **mengubah wewenangnya** tanpa ada yang menyentuh
aturan.

Sesudah tiket ini, pekerjaan tetap menunggu di antrean bersama, dan ⭐ pemeriksaan wewenang
bertanya **"apakah pengguna ini anggota antrean X"** — jawabannya tidak berubah ketika daftar
diurutkan ulang.

#### Area codebase

- Lapisan service: penugasan ke antrean
- Lapisan service: pemeriksaan keanggotaan antrean

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Antrean bersama | `Flow\InputRealizationTreatyIn.xml` — **enam** Assignment memakai `ToWorkBasket`, **nol** `ToWorklist` |
| ⛔ Penunjukan menurut posisi | `Section\SFAPortal_OpportunitiesList.xml` *(posisi 2)* · `DataTransform\InputPolicyTreatyIn_preDT.xml` *(posisi 2)* · `When\IsUW.xml` *(posisi 1)* |

#### ADR terkait

- **ADR-U-0002** — RBAC memakai peran yang sudah ada

#### Acceptance criteria

- [ ] **AC 11** — setiap penugasan masuk ke **antrean bersama**; ⛔ tidak ada kotak masuk pribadi
- [ ] **AC 14** — keanggotaan antrean diperiksa **menurut nama antrean**, ⛔ bukan menurut nomor urut
- [ ] **AC 92** — berkas menunggu **posisi**, bukan orang

#### Perintah verifikasi

1. Tambahkan pengguna ke satu antrean lain, lalu urutkan ulang daftarnya — ⭐ wewenangnya
   **tidak berubah**.
2. Ambil berkas sebagai pengguna kedua yang memegang posisi sama — ⭐ **berhasil**.

#### Catatan

⚠️ `[terverifikasi]` Pola penunjukan-menurut-posisi ada di **41 berkas pada 9 modul** di seluruh
korpus. ⭐ Perubahan yang sama berlaku di sana ketika modul itu digarap — **catat, jangan kerjakan
di tiket ini**.

## NB Treaty In - 05 - Peran menggantikan nama orang — tertunda sampai pemetaan diterima

**Status:** needs-info
**Blocked by:** pemetaan **nama → peran** dari `[IAM]` dan `[work owner]`
**Menutup:** AC 12 · 13 · 81 · 82 · 91 *(5 AC)* — US 17 · 18

#### Hasil & nilai pengguna

Hari ini wewenang di **dua belas tempat** ditentukan dengan memeriksa *"apakah pengguna ini orang
tertentu"*. `[terverifikasi]` Nilainya **tidak disalin** ke artefak mana pun. ⛔ Aturan seperti itu
berhenti bekerja ketika orangnya pindah jabatan atau keluar, **tanpa pemberitahuan**. Dan peran
sendiri disimpan di **kolom nomor telepon**.

Sesudah tiket ini, wewenang ditentukan **peran**, dan peran disimpan di medan peran — ⭐ kepergian
seseorang tidak lagi mematahkan alur.

#### Blocker

⛔ **Keputusannya jelas, pelaksanaannya belum bisa.** `[keputusan work owner]` P12 dan P28
menetapkan penggantian dengan peran; ⛔ **pemetaan nama → peran tidak ada di korpus.**

| Yang dibutuhkan | Kenapa |
| --- | --- |
| ⛔ untuk tiap dari **12 tempat**: peran penggantinya | tanpa itu, penggantian adalah tebakan |
| ⛔⛔ untuk **dua layar**: **arahnya** | `[terverifikasi]` nama yang sama dipakai **dua arah** — satu bagian muncul **hanya** untuk orang itu, bagian lain untuk **semua kecuali** orang itu. ⛔ Satu peran untuk keduanya akan **membalik** salah satunya |
| ⛔ peran ketiga | ⭐ tangga punya **tiga** posisi, peran yang terbaca baru **dua** |

#### Area codebase

- Lapisan service: pemeriksaan wewenang
- Lapisan repository: penyimpanan peran pengguna

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Peran di kolom telepon | `When\IsTreaty1.xml` · `When\IsSPVTreaty1.xml` — dibandingkan terhadap `OperatorID.pyTelephone` |
| Guard identitas orang | **12 berkas**, lewat `OperatorID.pyUserIdentifier` |
| Guard di layar | **4 layar · 12 tempat**, salah satunya memakai medan identitas yang berbeda |

#### ADR terkait

- **ADR-U-0002** — RBAC memakai peran yang sudah ada; rangkap peran ditolak

#### Acceptance criteria

- [ ] **AC 12** — wewenang ditentukan **peran**, bukan nama orang
- [ ] **AC 13** — peran disimpan di medan peran; ⛔ kolom nomor telepon kembali berisi nomor telepon
- [ ] **AC 81** — ke-12 tempat ditandai **tertunda**, ⛔ tidak dibangun dengan peran yang ditebak
- [ ] **AC 82** — ⛔ **arah** pemeriksaan pada dua layar **tidak ditebak**
- [ ] **AC 91** — ⛔ peran karangan **tidak dibuat** untuk menutup kekurangan posisi ketiga

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| ⛔ **1** | pemetaan nama → peran untuk 12 tempat, **berikut arahnya** | ⛔ **MENAHAN** |
| ⛔ **2** | peran yang tersedia belum cukup untuk tiga posisi | ⛔ **MENAHAN** |
| **3** | penetapan, pembuatan, dan pencabutan peran belum dibahas | tidak menahan |

#### Catatan

⛔ **Menebak peran berarti memberi atau mencabut wewenang atas dasar tebakan.** Tiket ini ditulis
lengkap supaya siap dikerjakan begitu pemetaannya tiba, ⛔ **bukan supaya dikerjakan sekarang.**

## NB Treaty In - 06 - Penggolongan jenis usaha — 36 baris, berhenti di yang pertama cocok, bawaan UNKNOWN

**Status:** ready-for-agent
**Blocked by:** 01
**Menutup:** AC 19 · 20 · 21 · 22 · 66 · 67 · 75 · 76 *(8 AC)* — US 25 · 26 · 32

#### Hasil & nilai pengguna

Hari ini jenis usaha digolongkan otomatis dari kode bisnisnya lewat sebuah tabel keputusan.
⚠️ `[terverifikasi]` Keterangan yang ditulis manusia pada sebagian aturan penggolong **berbeda dari
syarat yang benar-benar dijalankan** — sembilan di antaranya bertentangan.

Sesudah tiket ini, penggolongan berjalan **persis seperti sistem lama**, ⭐ termasuk berhenti di
baris pertama yang cocok dan memberi nilai bawaan yang jelas ketika tidak ada yang cocok — sehingga
berkas tetap dapat diproses.

#### Area codebase

- Lapisan service: penggolongan jenis usaha
- Fungsi murni: penilaian baris penggolong

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Tabel penggolong | `DecisionTable\BusinessType_DeT.xml` — **36 baris**, dua kolom penguji, **128** nilai kode seluruhnya unik |
| Cara evaluasi | `pyEvaluateAllRows` = `no` ⇒ ⭐ **berhenti di baris pertama yang cocok** |
| Nilai bawaan | ⭐ **`"UNKNOWN"`** |
| Penggolong lini jiwa | `When\IsLife.xml` — **16** nilai kode, dirujuk dua aktivitas penjaga |
| Pembeda jalur simpan | penanda bernilai `EDM` *(endorsemen)* / `POLICY` *(polis baru)* |
| Penanda penempatan keluar | `DataTransform\TestTreatyToFacStatus.xml` — dua kode angka |

#### ADR terkait

- **ADR-U-0005** — aturan lingkungan tidak ditiru sebagai rule

#### Acceptance criteria

- [ ] **AC 19** — penggolongan **berhenti di baris pertama yang cocok**
- [ ] **AC 20** — kode yang tidak cocok baris mana pun menghasilkan **`"UNKNOWN"`**
- [ ] **AC 21** — ke-**128** kode menghasilkan penggolongan yang sama seperti sistem lama
- [ ] **AC 22** — syarat diambil dari **yang dijalankan**, ⛔ bukan dari keterangannya — termasuk
      pada **9** aturan yang bertentangan
- [ ] **AC 66** — penanda `EDM` ⇒ jalur endorsemen; `POLICY` ⇒ polis baru; ⛔ keduanya **tidak**
      menempuh cara penyimpanan yang sama
- [ ] **AC 67** — penggolong lini jiwa dimigrasi apa adanya, **16** nilai kode tanpa perubahan
- [ ] **AC 75** — berkas non-proporsional ditandai sesuai jenis proporsinya
- [ ] **AC 76** — penanda penempatan keluar dipasang **hanya** pada dua kode yang dikenal

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **10** | arti kode lini bisnis, termasuk satu kode berskala penomoran berbeda | tidak menahan |
| **11** | arti dua kode penanda penempatan keluar | tidak menahan |

#### Perintah verifikasi

1. Jalankan ke-**128** kode melalui penggolong — ⭐ hasilnya **sama** dengan sistem lama, satu per satu.
2. Kirim kode yang tidak dikenal — ⭐ hasilnya **`"UNKNOWN"`**, bukan kosong dan bukan galat.
3. Kirim kode yang cocok **dua** baris — ⭐ yang dipakai baris **pertama**.

#### Catatan

⚠️ `[keputusan work owner]` P23 menetapkan **keterangan diabaikan seluruhnya** — tidak dipakai
bahkan sebagai petunjuk. ⛔ Membangun dari keterangan akan salah pada **9 dari 50** penggolong yang
punya keduanya.

## NB Treaty In - 07 - Uang — presisi penuh, persentase bukan uang, pajak brokerage apa adanya

**Status:** ready-for-agent
**Blocked by:** 01
**Menutup:** AC 18 · 23 · 24 · 25 · 26 · 27 · 28 · 85 · 86 *(9 AC)* — US 21 · 22 · 36 · 38 · 39

#### Hasil & nilai pengguna

Hari ini nilai uang ditampilkan di layar apa adanya seperti tersimpan, ⭐ **bukan dihitung ulang**.
⚠️ `[terverifikasi]` Empat medan yang tampak seperti uang sebenarnya **persentase** — `12.5` berarti
12,5 persen. Dan potongan brokerage dibagi **1,022** bila jenis pajaknya *inclusive*, ⛔ dengan
pembanding teks yang **persis**: huruf kecil atau spasi di belakang menggeser angkanya **2,2 %**
tanpa pesan galat.

Sesudah tiket ini, nilai uang tersimpan **berpresisi penuh**, ditampilkan beserta mata uangnya, dan
⭐ rumus pajak brokerage berjalan **persis** seperti sistem lama — termasuk ketidakseragamannya.

#### Area codebase

- Lapisan repository: pembacaan nilai uang dan mata uang pasangannya
- Fungsi murni: rumus potongan pajak brokerage
- Lapisan handler: penyajian

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Rumus pajak brokerage | `@if(TypeTax == "Inclusive", Deduction / (102.2/100), Deduction)` — **12 tempat pada 7 Activity** |
| Nilai bawaan jenis pajak | `"Inclusive"`; ⛔ nilai lawannya **tidak pernah tertulis** di mana pun |
| Medan uang | **kolom tersimpan** pada sumber relasional; ⛔ nol `Activity` atau `DataTransform` menghasilkannya |

#### ADR terkait

- ⭐ **ADR-U-0003** — uang **tidak** direpresentasikan sebagai `float`

#### Acceptance criteria

- [ ] **AC 23** — nilai uang disimpan **berpresisi penuh**
- [ ] **AC 24** — pembulatan **hanya** di titik penyajian, ⛔ tidak pernah di repository
- [ ] **AC 25** — uang **tidak pernah** diwakili tipe pecahan biner
- [ ] **AC 26** — empat medan itu dibaca sebagai **persentase**, bukan jumlah uang
- [ ] **AC 27** — jenis pajak *inclusive* ⇒ potongan dibagi **1,022**
- [ ] **AC 28** — nilai lain — ⭐ **termasuk kosong, huruf kecil, atau berspasi** — ⇒ potongan
      dipakai apa adanya
- [ ] **AC 18** — kedelapan medan uang **ditampilkan apa adanya**, ⛔ tidak dihitung ulang
- [ ] **AC 85** — angka uang ditampilkan beserta **kode mata uang pasangannya**
- [ ] **AC 86** — ⛔ format penyajian **belum ditetapkan**; penyajian tidak dibangun dengan format
      yang ditebak

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | tipe penyimpanan kolom uang — sisi Pega menyimpannya sebagai **teks** | ⚠️ menahan **penguraian**, bukan rumusnya |
| **7** | format penyajian — desimal dan pemisah ribuan | ⚠️ menahan **penyajian** saja |
| **8** | satu medan persentase **dari apa** | tidak menahan |
| **9** | ketidakseragaman presisi 4 lawan 8 desimal | tidak menahan |

#### Perintah verifikasi

1. Jenis pajak `"Inclusive"` — ⭐ potongan **dibagi 1,022**.
2. Jenis pajak `"inclusive"` huruf kecil — ⭐ potongan **tidak** dibagi; selisihnya **2,2 %**.
3. Jenis pajak **kosong** — ⭐ sama seperti butir 2.
4. Simpan nilai berdesimal panjang, baca kembali — ⭐ **tidak berubah**.

#### Catatan

⭐ **Penyajian uang TIDAK bergantung pada P18.** `[keputusan work owner]` P43 menetapkan angka
**disalin apa adanya**, sehingga bagian terbesar uang dapat dibangun sekarang.

> ⚠️ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat penutup semula berbunyi:
> > *"⛔ Yang menunggu P18 adalah **perhitungannya**, ada di tiket **13**."*
>
> ⭐ **P18 ditarik** — perhitungannya pun tidak menunggu P18. Rumus pajak brokerage terbaca penuh
> di ekspor: `@if(TypeTax="Inclusive", @divide(Deduction,@divide(102.2,100,8),8), Deduction)`,
> **PPH 2 %**, **PPN 2,2 %**. Perhitungannya ada di tiket **13**, yang kini tertahan oleh **P30**
> dan **P8** saja.

## NB Treaty In - 08 - Keutuhan penyimpanan — satu transaksi, skema eksplisit, arah ketergantungan

**Status:** ready-for-agent
**Blocked by:** 01
**Menutup:** AC 29 · 30 · 60 · 83 · 90 *(5 AC)* — US 33 · 34

#### Hasil & nilai pengguna

Hari ini penyimpanan realisasi treaty **tidak punya jaminan keutuhan**. `[data DBA]` Dua dari tiga
program penyimpan **menyelesaikan penyimpanannya sendiri**, dan blok pemanggilnya menyelesaikannya
lagi — ⛔ **penyelesaian ganda**. Akibatnya kegagalan pada langkah berikutnya **meninggalkan data
yang sudah permanen**, dan aplikasi tidak dapat membatalkannya.

Sesudah tiket ini, seluruh urutan penyimpanan berada dalam **satu transaksi** — ⭐ kegagalan di
langkah mana pun membatalkan seluruhnya, dan tidak ada lagi berkas yang tersimpan separuh.

#### Area codebase

- Lapisan repository: batas transaksi
- Lapisan service: urutan penyimpanan
- Uji arsitektur: arah ketergantungan

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penyelesaian ganda | `RDBList\SaveTreatyIn.xml` · `RDBList\SavePolisTreatyIn_SQL.xml` — blok memuat penyelesaian, dan program yang dipanggil **juga** |
| Yang **tidak** menyelesaikan sendiri | `RDBList\GetSequenceNumber_SQL.xml` |
| Dua ejaan nama objek | empat nama muncul **dengan dan tanpa** awalan skema di 1.151 naskah SQL korpus |

#### ADR terkait

- **ADR-U-0009** — migrasi penuh; tidak ada koeksistensi dua penulis

#### Acceptance criteria

- [ ] **AC 29** — seluruh urutan penyimpanan berada dalam **satu transaksi**; kegagalan di tengah
      menyisakan **nol** baris
- [ ] **AC 83** — kegagalan menyimpan riwayat **membatalkan seluruh transaksi**
- [ ] **AC 30** — setiap query menyebut **skema secara eksplisit**
- [ ] **AC 90** — keempat nama berejaan ganda diperlakukan sebagai **satu objek**
- [ ] **AC 60** — arah ketergantungan `handlers → services → repository`; ⛔ tidak terbalik, tidak
      memotong lapisan

#### Perintah verifikasi

1. Suntikkan kegagalan di tengah urutan penyimpanan — ⭐ **nol** baris tersisa.
2. Sambung sebagai pengguna **selain** pemilik skema — ⭐ query **tetap menemukan** objeknya.
3. Jalankan uji arsitektur — ⭐ **nol** panggilan dari repository ke service.

#### Catatan

⚠️ `[penyimpangan sadar]` Satu transaksi **lebih ketat** daripada sistem lama, dan itu disengaja —
alasannya tertulis di `spec.md` §5.7.
⛔ Uji keutuhan transaksi **tidak dapat dijalankan dengan tiruan**; ia wajib berjalan lawan basis
data sungguhan.

## NB Treaty In - 09 - Tanggal — satu tipe, satu format, dan pengisian bawaan yang ditiru apa adanya

**Status:** ready-for-agent
**Blocked by:** —
**Menutup:** AC 32 · 33 · 34 · 35 · 69 *(5 AC)* — US 2 · 35 · 44

#### Hasil & nilai pengguna

Hari ini tanggal disimpan **sebagai teks**, dan ⛔ **dua susunan berbeda dipakai di berkas yang
sama** — hari-bulan untuk tanggal mulai, bulan-hari untuk tanggal laporan. `[terverifikasi]`
⚠️ Sebagian nilai lama karenanya **ambigu secara mutlak**: satu nilai dapat berarti dua tanggal
berbeda, dan kekeliruannya **tidak menimbulkan pesan galat**.

Sesudah tiket ini, tanggal disimpan sebagai **tipe tanggal**, dengan **satu format** di seluruh
sistem, dan formatnya diurus di lapisan layar — ⭐ urutan dan perbandingan tanggal menjadi benar.

#### Area codebase

- Lapisan repository: tipe tanggal
- Lapisan service: pengisian tanggal bawaan
- Lapisan handler: format tampilan

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Dua format berdampingan | `DataTransform\InputPolicyTreatyIn_preDT.xml` — dua susunan berbeda di berkas yang sama |
| Pengisian bila kosong | gerbang tanggal-mulai, tanggal-akhir, dan tanggal-laporan kosong ⇒ **tanggal hari ini** |
| Aturan satu tahun | `DataTransform\SystemSetOneYear_DT.xml` — menambah satu tahun ke tanggal mulai |

#### ADR terkait

- **ADR-U-0009** — migrasi penuh

#### Acceptance criteria

- [ ] **AC 32** — tanggal disimpan sebagai **tipe tanggal**, ⛔ bukan teks
- [ ] **AC 33** — **satu format** dipakai di seluruh sistem; ⛔ tidak ada dua format berdampingan
- [ ] **AC 34** — tanggal akhir kosong diisi **tanggal hari ini**, ⛔ bukan ditambah satu tahun
- [ ] **AC 35** — tanggal mulai dan tanggal laporan kosong diisi tanggal hari ini
- [ ] **AC 69** — migrasi menghasilkan kontrak ber-tanggal-akhir **sama dengan tanggal mulai**
      untuk berkas yang tanggal akhirnya kosong

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| ⚠️ **15** | data lama bertanggal **ambigu mutlak** — satu nilai, dua arti | ⚠️ menahan **migrasi**, bukan perilaku baru |
| **12** | kapan aturan satu-tahun berjalan, dan apakah ia menimpa pengisian hari-ini | tidak menahan |

#### Perintah verifikasi

1. Kosongkan tanggal akhir, simpan — ⭐ terisi **tanggal hari ini**, bukan setahun kemudian.
2. Simpan tanggal, baca kembali, urutkan — ⭐ urutannya **benar secara kronologis**.
3. Ubah format tampilan — ⭐ nilai tersimpan **tidak berubah**.

#### Catatan

⚠️ `[penyimpangan sadar]` Tipe tanggal dan format tunggal **berbeda** dari Pega. ⭐ Tetapi
pengisian tanggal-akhir-kosong-jadi-hari-ini **ditiru apa adanya** *(P35)*, walau hasilnya kontrak
bermasa berlaku nol hari — ⛔ **disengaja, bukan cacat migrasi.**

## NB Treaty In - 10 - Jejak audit dan kronologi — identitas akses terpisah dari nama tampilan

**Status:** ready-for-agent
**Blocked by:** 03
**Menutup:** AC 39 · 40 · 41 · 42 · 43 · 44 · 71 · 72 *(8 AC)* — US 15 · 16 · 20 · 40 · 41 · 42

#### Hasil & nilai pengguna

Hari ini riwayat akseptasi mencatat siapa memutuskan apa. ⛔ Tetapi penampung identitas operator
**tidak pernah diisi** dari aplikasi — `[terverifikasi]` **nol dari 1.151** naskah SQL di 21 modul
menyebutnya sebagai kolom. ⚠️ Dan medan "nama operator" diisi dari **dua sumber berbeda** di dua
tahap berdekatan: pengenal akun di satu tahap, nama tampilan di tahap lain.

Sesudah tiket ini, setiap tindakan mencatat **identitas akses login**-nya terpisah dari **nama
tampilan**-nya, dan ⭐ keduanya tidak lagi tertukar — jejaknya tetap benar walau nama tampilan
seseorang berubah.

#### Area codebase

- Lapisan repository: penulisan riwayat akseptasi
- Lapisan service: kronologi dan catatan pengguna

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Pemisahan yang sudah benar | `RDBList\InsertViewSuggest_SQL.xml` — penampung akses-login diisi dari pengenal akun, terpisah dari penampung nama tampilan |
| ⛔ Pengisian yang **bug** | `DataTransform\DeptHeadTreatyInUW_preDT.xml` — medan nama operator diisi dari **pengenal akun** |
| Pengisian yang benar | `DataTransform\InputPolicyTreatyIn_preDT.xml` — dari **nama tampilan** |
| Mekanisme kronologi | `DataTransform\AddToListCommentsPolicyTreatyIn_DT.xml` — empat medan: catatan, penanda persetujuan, tanggal, operator |
| ⛔ Nama orang di dalam teks pesan | `DataTransform\DeptHeadTreatyIn_UW_postDT.xml` — **1** kemunculan; pesan lain di berkas sama mengambilnya dari data |

#### ADR terkait

- ⭐ **ADR-U-0007** — jejak audit merekam siapa + kapan untuk setiap transisi **dan setiap jalur balik**

#### Acceptance criteria

- [ ] **AC 39** — penampung identitas operator diisi dari **identitas akses login**
- [ ] **AC 40** — penampung nama diisi dari **nama tampilan**
- [ ] **AC 41** — penampung identitas operator **terisi** pada setiap penulisan riwayat
- [ ] **AC 42** — medan nama operator diisi dari **nama tampilan di setiap tahap**, termasuk tahap
      jenjang ketiga
- [ ] **AC 43** — setiap perpindahan tahap menulis **satu baris riwayat**
- [ ] **AC 44** — pemberitahuan menyebut nama orang **dari data**; ⛔ tidak tertanam di dalam teks
- [ ] **AC 71** — catatan pengguna tersimpan bersama tanggal dan operatornya
- [ ] **AC 72** — riwayat dapat dibaca **berurutan waktu**

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **16** | berapa baris yang penampung identitas operatornya **sudah terisi** sekarang | tidak menahan — ⚠️ bila bukan nol, ada penulis **di luar aplikasi** yang belum terpeta |

#### Perintah verifikasi

1. Lakukan satu tindakan di tiap jenjang — ⭐ tiga baris riwayat, **masing-masing berisi kedua
   identitas**.
2. Ubah nama tampilan seorang pengguna — ⭐ riwayat lama **tidak berubah**.
3. Baca pemberitahuan — ⭐ namanya **dari data**, dan berubah mengikuti data.

#### Catatan

⚠️ `[penyimpangan sadar]` Mengisi penampung identitas operator adalah **perbaikan jejak audit**,
bukan peniruan — kolomnya selama ini kosong. ⭐ Dan pengisian nama operator dari pengenal akun
adalah **bug**, ⛔ bukan perbedaan maksud antar tahap *(P33)*.

## NB Treaty In - 11 - Layar realisasi — medan wajib, medan terkunci, dan bagian yang tidak dibangun

**Status:** ready-for-agent
**Blocked by:** 02
**Menutup:** AC 45 · 46 · 49 · 50 · 51 · 53 · 54 · 55 · 56 · 77 *(10 AC)* — US 27 · 28 · 29 · 31 · 32

#### Hasil & nilai pengguna

Hari ini layar realisasi menuntut sejumlah medan diisi, mengunci sebagian lainnya, dan
menyembunyikan bagian yang tidak berlaku. ⚠️ `[terverifikasi]` **Delapan puluh** bagian layar
diberi syarat tampil yang **tidak mungkin pernah benar** — dan **tidak satu pun** dari yang 80 itu
menyembunyikan sebuah medan; seluruhnya label dan hiasan.

Sesudah tiket ini, pengguna melihat layar yang **menandai medan wajibnya dengan jelas**, mengunci
yang tidak boleh ia ubah, dan ⭐ **tidak memuat 80 bagian mati** yang selama ini ada tanpa pernah
tampil.

#### Area codebase

- Lapisan handler: validasi medan wajib
- Antarmuka: penandaan medan wajib dan terkunci
- Antarmuka: daftar pilihan mata uang

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Medan wajib | **27** medan berbeda di **6** layar; sebarannya di `spec.md` §5.11 |
| Medan terkunci | **38**, seluruhnya lewat syarat hanya-baca yang **selalu benar**; setiap kunci mengenai **tepat satu medan** |
| Bagian mati | **91**; ⭐ **80** menempel pada sel tunggal, **11** pada wadah yang memuat medan |
| Pengecualian mata uang | `ReportDefinition\BrowseCurrency_RD.xml` · `BrowseCurrencyTreatyIn_RD.xml` |
| Pembersihan pesan galat | **12** tempat yang menghapus **sebelum** pesan dipasang |

#### ADR terkait

- **ADR-U-0002** — RBAC memakai peran yang sudah ada

#### Acceptance criteria

- [ ] **AC 45** — layar menuntut medan wajibnya; **27** medan berbeda di **6** layar
- [ ] **AC 46** — layar jenjang ketiga **tidak** mewajibkan enam medan yang wajib di layar admin
- [ ] **AC 49** — **38** medan terkunci permanen
- [ ] **AC 50** — **36** dari 38 di layar jenjang ketiga; **2** di layar biasa
- [ ] **AC 51** — setiap kunci mengenai **tepat satu medan**, ⛔ bukan bagian atau tab
- [ ] **AC 53** — **80** bagian mati **tidak dibangun**
- [ ] **AC 54** — daftar pilihan mata uang **tidak memuat** kode yang dikecualikan
- [ ] **AC 55** — kode jenis kontrak ditampilkan **apa adanya** bila keterangannya belum tersedia
- [ ] **AC 56** — nilai kode di luar daftar yang dikenal **tetap diterima dan disimpan**
- [ ] **AC 77** — pembersihan pesan galat di awal diterima apa adanya — **12** tempat

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **17** | **empat wadah** berisi **104** medan — masih dipakai atau ditinggalkan | ⚠️ menahan **keempatnya saja**, bukan tiketnya |
| **14** | apakah ada mata uang mati lain yang tidak disembunyikan | tidak menahan |

#### Perintah verifikasi

1. Kosongkan satu medan wajib, simpan — ⭐ **ditolak**, dengan pesan yang menyebut medannya.
2. Coba sunting medan terkunci — ⭐ **tidak bisa**, dan tampil sebagai terkunci.
3. Cari kode mata uang yang dikecualikan di daftar pilihan — ⭐ **tidak ada**.
4. Kirim kode jenis kontrak yang tidak dikenal — ⭐ **tetap tersimpan**, ditampilkan apa adanya.

#### Catatan

⭐ **Delapan puluh bagian mati dibuang tanpa ditanyakan kepada siapa pun** `[keputusan work owner]`
P44 — `[terverifikasi]` sebabnya struktural: membangunnya berbiaya nol, membuangnya juga.
⛔ **Empat wadah berisi 104 medan** menunggu Product & Underwriting dan **tidak dibangun** sebelum
dijawab.

## NB Treaty In - 12 - Layar jenjang ketiga — sembilan medan wajib SEKALIGUS terkunci

**Status:** ready-for-agent
**Blocked by:** 11

> ⭐⭐ **LEPAS DARI BLOKIR.** `[penyimpangan sadar]` 2026-09-22 — dua baris ini semula berbunyi:
> > *"**Status:** blocked · **Blocked by:** 11 · ⛔ **P18** *(isi 268 langkah penetapan nilai)*"*
>
> **P18 ditarik oleh tim migrasi.** Isi langkah ada di dalam ekspor, di tag
> `PropertiesName`/`PropertiesValue` **tanpa awalan `py`**. Lihat `VERIFIKASI-P18.md`.
**Menutup:** AC 47 · 48 · 52 · 78 · 80 *(5 AC)* — US 11 · 30

#### Hasil & nilai pengguna

Hari ini pemegang jenjang ketiga **mengisi tujuh medan sendiri**, lalu memutuskan — ⭐ ia **bukan**
sekadar melihat. ⛔ Tetapi layarnya memuat **sembilan medan yang wajib diisi sekaligus terkunci**:
medan itu wajib, **tidak dapat diisi pengguna**, dan yang seharusnya mengisinya adalah
**perhitungan yang isinya belum terkirim**.

Sesudah tiket ini, pemegang jenjang ketiga dapat mengisi medan keputusannya dan menyimpan berkas.

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat berikut dikutip utuh lalu ditarik:
> > *"⛔ **Selama P18 kosong, layar itu tidak dapat disimpan** — dan itu bukan cacat pembangunan,
> > melainkan akibat bahan yang belum lengkap."*
>
> ⭐ `[terverifikasi]` **Layar itu dapat disimpan.** Dari **34** medan wajib-sekaligus-terkunci di
> seluruh `Section`, **28** punya langkah pengisi yang isinya terbaca. Tiga sisanya — `.ExcessLoss`,
> `.OutstandingClaim`, `.SalvageValue` — **diketik underwriter** di `DetailPolicyTreatyIn` /
> `GeneralPolicyTreatyIn` *(`bacasaja=false`)*, dan hanya **ditampilkan terkunci** di sini.

#### ⭐ Blocker — **NIHIL**

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — bab Blocker semula berbunyi:
> > *"⛔⛔ **Sembilan medan wajib-sekaligus-terkunci** … | ⛔ penyimpanan layar jenjang ketiga |
> > kesembilan medan wajib **tidak dapat diisi** | ⛔ penyelesaian tangga | berkas dapat naik ke
> > jenjang ketiga, ⛔ **tidak dapat keluar** dari sana | … ⭐ **Dua jalan keluar, keduanya bukan
> > keputusan teknis:** ekspor ulang **dengan isi langkah disertakan** *(P18)*, atau **pencabutan
> > kewajiban** atas kesembilan medan oleh Product & Underwriting."*

⭐ **Tidak ada jalan keluar yang perlu ditempuh.** Kesembilan medan itu memang wajib sekaligus
terkunci, dan itu tetap benar — tetapi **yang mengisinya terbaca di ekspor**. `[terverifikasi]`

| Medan | Diisi oleh |
| --- | --- |
| `Deduction2` | `InputPolicyTreatyEDMDetail_NP` · `InputPolicyTreatyInDetail_NonProp` |
| `OveriddingCommOgp` | `CountOGPONP_Act` · `CountOverridingCommOgp_Act` |
| `OveriddingCommOnp` | `CountOGPONP_Act` · `CountOverridingCommOnp_Act` |
| `ResultOnp1` | `CountOGPONP_Act` · `CountResult1Onp_Act` |
| `RiCommOgp` | `CalculatePremi_Act` · `CountOGPONP_Act` |
| `RiCommOnp` | `CountOGPONP_Act` · `CountResult1Onp_Act` |
| ⭐ `ExcessLoss` | **diketik underwriter** di layar `DetailPolicyTreatyIn` / `GeneralPolicyTreatyIn` |
| ⭐ `OutstandingClaim` | idem |
| ⭐ `SalvageValue` | idem |

#### Area codebase

- Antarmuka: layar jenjang ketiga
- Lapisan handler: validasi medan wajib per tingkat

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Medan wajib per tingkat | layar jenjang ketiga **tidak** mewajibkan enam medan yang wajib di layar admin, ⭐ tetapi **mewajibkan satu** yang tidak wajib di sana |
| Medan terkunci | **36** dari 38 ada di dua layar jenjang ketiga |
| Medan yang dapat diisi | **tujuh** |
| Pembersihan pesan yang **ditahan** | **5** tempat yang menghapus **sesudah** validasi memasang pesan; seluruhnya di rantai yang sama |

#### ADR terkait

- **ADR-U-0007** — jejak audit setiap transisi

#### Acceptance criteria

- [ ] **AC 47** — layar jenjang ketiga **mewajibkan** satu medan yang tidak wajib di layar admin
- [ ] **AC 48** — berkas **tidak dapat disimpan** bila medan wajib pada tingkat itu kosong
- [ ] **AC 52** — pemegang jenjang ketiga **dapat mengisi tujuh medan**; ⛔ layarnya **bukan**
      sepenuhnya hanya-baca
- [ ] **AC 78** — **5** tempat yang menghapus pesan **sesudah** validasi memasangnya **ditahan**
      dan tidak dibangun
- [ ] **AC 80** — ⭐ layar **dapat disimpan**; enam medan diisi langkah, tiga diketik underwriter
      di layar sebelumnya *(AC 80 dicabut lalu diganti 2026-09-22 — lihat `spec.md`)*

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| ⭐ ~~**P18**~~ | ~~isi 268 langkah penetapan nilai~~ | ⭐ **DITARIK 2026-09-22** — tidak menahan |
| **13** | **5** tempat penghapus pesan galat sesudah validasi | ⛔ menahan kelimanya |

#### Perintah verifikasi

⭐ **Dapat dijalankan sepenuhnya.**

1. Buka layar jenjang ketiga — ⭐ **tujuh** medan dapat diisi, **36** terkunci.
2. Coba simpan dengan kesembilan medan kosong — ⭐ **ditolak**, dan pesannya menyebut bahwa
   nilainya berasal dari perhitungan yang belum tersedia.

#### Catatan

⭐ **Tiket ini semula ditulis lengkap walau berstatus `blocked`** — supaya begitu P18 dijawab,
pekerjaannya sudah terumus. ⭐⭐ **P18 tidak pernah perlu dijawab**, dan pekerjaannya kini dapat
diambil. Kalimat lama ⛔ *"Jangan ditandai `ready-for-agent` sebelum itu"* — **ditarik**.

## NB Treaty In - 13 - Rantai perhitungan uang — karantina P18 DIBUKA, sisa dua penahan lain

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Kedua penahan gugur. **P30** di luar lingkup Treaty In — seluruh rantai pemanggilnya bekerja pada `OfferFacIn`, nol `PolicyTreatyIn`. **P8** dicabut — bentuk muatannya terbaca dari saudara kelas induk: empat medan, dua panggilan keluar, surat bila gagal.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 6 dan 7.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)* — ⚠️ **bukan lagi karena P18**
~~**Blocked by:** ⛔ **P30** *(aturan penjumlah yang hilang)* · ⛔ **P8** *(muatan efek keluar)*~~ ⛔ **penahan gugur 23-09-2026**

> ⭐⭐ **SISI P18 DIBUKA.** `[penyimpangan sadar]` 2026-09-22 — dua baris ini semula berbunyi:
> > *"# 13: Rantai perhitungan uang — dikarantina sampai bahannya lengkap · **Status:** blocked ·
> > **Blocked by:** ⛔ **P18** *(isi 268 langkah)* · ⛔ **P8** *(muatan efek keluar)*"*
>
> **P18 ditarik oleh tim migrasi.** Isi 268 langkah ada di ekspor — **707 pasangan nama=nilai
> terisi pada 269 langkah**, nol tanda terpotong. Lihat `VERIFIKASI-P18.md`.
>
> ⚠️ **Tiket ini TETAP `blocked`**, tetapi sebabnya berubah: **P30** dan **P8**, bukan P18.
> `[terverifikasi]` `SumTSIPremiSpreadRNMMultiCob_Act` **tetap tidak ada di satu pun dari 21 modul**
> korpus — diperiksa ulang 2026-09-22. ⛔ **P8 tidak dibuka di ronde ini**: keputusan
> "ikuti apa adanya" sudah ada, tetapi sisa `[terbuka]`-nya — ARASAPAS itu sistem apa — masih milik
> `[Product+Underwriting]`.
**Menutup:** AC 79 · 87 *(2 AC)* — US 46

#### Hasil & nilai pengguna

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kedua paragraf berikut dikutip utuh lalu
> ditarik:
> > *"Hari ini seluruh aritmetika uang … tinggal di dalam **langkah penetapan nilai** yang ⛔
> > **isinya tidak ikut dalam ekspor** yang diterima tim migrasi. ⭐ Bentuknya diketahui;
> > **rumusnya tidak**. … ⛔ **Tiket ini belum dapat dikerjakan**, dan ⛔ **tidak boleh dikerjakan
> > dengan rumus yang ditebak.**"*

⭐⭐ **Rumusnya diketahui.** `[terverifikasi]` Seluruh aritmetika uang tinggal di **langkah
penetapan nilai**, dan isinya **terbaca penuh** di ekspor — termasuk pembagi **102,2** untuk
`TypeTax="Inclusive"`, **PPH 2 %**, **PPN 2,2 %**, dan **presisi 8 angka** pada tiap `@divide`.

⚠️ **Yang masih menahan tinggal dua, dan keduanya bukan tentang rumus:** satu aturan penjumlah
yang **hilang dari seluruh korpus** *(P30)*, dan **muatan** efek keluar terakhir *(P8)*.
⛔ **Tetap tidak boleh dikerjakan dengan rumus yang ditebak** — sekarang tidak perlu menebak.

#### Blocker

| Butir | Yang ditunggu | Pemilik |
| --- | --- | --- |
| ⭐ ~~**P18**~~ | ~~isi **268** langkah penetapan nilai yang terjangkau~~ — ⭐ **DITARIK 2026-09-22** | ~~`[pemilik export Pega]`~~ |
| ⛔ **P30** | satu aturan penjumlah yang **dipanggil tetapi tidak ada di seluruh korpus 21 modul** | `[pemilik export Pega]` |
| ⛔ **P8** | **muatan** yang dikirim ke layanan luar di akhir alur | `[Product+Underwriting]` |

#### Yang SUDAH diketahui bentuknya

⭐ **Disebut supaya tidak dicari ulang.** ⚠️ Bab ini semula berjudul maksud *"apa yang TIDAK
menunggu P18"*; sesudah P18 ditarik, **tidak ada lagi yang menunggu P18.**

| Hal | Keadaan |
| --- | --- |
| medan | delapan medan uang **berpasangan mata uang** |
| sumber | ⭐ **kolom tersimpan**, bukan hasil hitung |
| arah aliran | ⭐ **dibaca, bukan dihasilkan** |
| ⭐ penyajian | ⭐ **sudah dapat dibangun** — tiket **07** |
| ⭐ rumus pajak brokerage | ⭐ **terbaca penuh** — tiket **07** |

#### Area codebase

- Lapisan service: rantai perhitungan uang
- Lapisan service: efek keluar di akhir alur

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Langkah penetapan nilai | **268** terjangkau dari titik masuk; ⛔ isinya **tidak terekspor** |
| Aturan yang hilang | dipanggil `Activity\SumTSIPremiSpreadedRNM_Act.xml`; ⛔ **nol salinan** di 21 modul |
| Efek keluar | `Flow\InputRealizationTreatyIn.xml` — `Utility2`, lewat `Activity\serviceInsertArasapas_act.xml` **18.013 B, satu langkah**; ⛔ implementasinya ada di **dua modul lain** dan **tidak boleh dipinjam** |

#### ADR terkait

- ⭐ **ADR-U-0003** — uang tidak `float`
- **ADR-U-0008** — efek keluar asinkron dengan antre-ulang
- **ADR-U-0015** — efek keluar wajib berhasil, transactional outbox

#### Acceptance criteria

- [ ] **AC 79** — ⭐ rantai perhitungan **dibangun dari isi langkah yang terbaca di ekspor**;
      ⛔ nol rumus yang ditebak *(AC 79 dicabut lalu diganti 2026-09-22 — lihat `spec.md`)*
- [ ] **AC 87** — pengiriman ke layanan luar **tidak dibangun** sebelum muatannya diketahui

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| ⭐ ~~**P18**~~ | ~~isi 268 langkah~~ | ⭐ **DITARIK 2026-09-22** — tidak menahan |
| ⛔ **4** | di bagian mana aturan penjumlah yang hilang tinggal | ⛔ **MENAHAN** |
| ⛔ **5** | muatan efek keluar | ⛔ **MENAHAN** |

#### Catatan

⛔⛔ **Menebak rumus di sini berarti membangun angka keuangan dari tebakan** — dan sekarang tidak
ada alasan menebak: rumusnya terbaca.

⚠️ Bab karantina di `spec.md` **§8.2** sudah **DIBUKA**, dan **keempat akibat** yang dulu
didalilkan bila P18 tidak dijawab **gugur seluruhnya**. Bab itu dipertahankan sebagai catatan.

⭐ **Yang dapat diambil sekarang:** seluruh rantai aritmetika uang. ⛔ **Yang masih ditahan:**
pemanggilan `SumTSIPremiSpreadRNMMultiCob_Act` *(P30)* dan pengiriman ke layanan luar *(P8, AC 87)*.

## NB Treaty In - 14 - Lingkup — yang TIDAK dibangun, dan diuji supaya tetap tidak terbangun

**Status:** ready-for-agent
**Blocked by:** —
**Menutup:** AC 61 · 62 · 63 · 64 · 65 · 88 *(6 AC)* — US 45

#### Hasil & nilai pengguna

Hari ini folder modul memuat **lebih banyak aturan daripada pekerjaan modul ini**. `[terverifikasi]`
Ia **dependency closure** — memuat aturan milik modul lain karena pewarisan kelas. ⛔ Membangun dari
isi folder apa adanya akan mewarisi **beban mati** dan pekerjaan yang bukan milik modul ini.

Sesudah tiket ini, sistem baru **tidak memuat** apa yang sudah dinyatakan di luar lingkup, dan
⭐ ada **uji yang menjaganya tetap begitu** — supaya ia tidak masuk kembali diam-diam pada ronde
pembangunan berikutnya.

#### Area codebase

- Uji lingkup: daftar yang tidak boleh terpanggil

#### Rule Pega sumber

| Yang **tidak** dibangun | Sebab | Butir |
| --- | --- | --- |
| **28** aturan yatim, **561** langkah penetapan | ⛔ tidak terjangkau; milik modul Fac | keadaan Bab 2 |
| **6** aturan pembongkar dokumen | penyimpanan pindah ke sumber relasional | P29 · P15 |
| aturan bernama sama yang menguji `= 1` | ⛔ **bukan aturan yang hidup** | P6 |
| penanda persetujuan **kedua** | ⛔ sudah tidak dipakai; ⭐ **berkas rule-nya nol di seluruh korpus** | P36 |
| medan "nomor surat" sebagai penanda arah | routing memakai medan tersendiri | P7 |
| pemanggilan penyalin data yang tidak dipakai | ⛔ tidak digunakan | P30 |

#### ADR terkait

- **ADR-U-0009** — migrasi penuh; tidak ada koeksistensi dua penulis

#### Acceptance criteria

- [ ] **AC 61** — **28** aturan yatim **tidak dimigrasi**; nol terpanggil
- [ ] **AC 62** — **6** aturan pembongkar dokumen **tidak dimigrasi**
- [ ] **AC 63** — aturan bernama sama yang menguji `= 1` **tidak dimigrasi**
- [ ] **AC 64** — penanda persetujuan kedua **tidak dibangun**; penampungnya tidak dibuat
- [ ] **AC 65** — medan "nomor surat" **tidak** dipakai menyimpan penanda arah
- [ ] **AC 88** — nol dari **561** langkah yatim terpanggil

#### Perintah verifikasi

1. Jalankan uji lingkup — ⭐ **nol** dari daftar di atas terpanggil.
2. Cari penanda persetujuan kedua di seluruh sistem baru — ⭐ **tidak ada**.

#### Catatan

⚠️ `[terverifikasi]` **Penanda persetujuan kedua dirujuk juga oleh modul Fac** — tiga berkasnya
merujuk aturan yang sama, yang berkasnya **tidak ada di seluruh korpus**. ⭐ Catat untuk modul itu;
⛔ **jangan kerjakan di tiket ini.**

## NB Treaty In - 15 - Migrasi, paritas, dan jejak keputusan — tidak ada butir terbuka yang ditutup diam-diam

**Status:** ready-for-agent
**Blocked by:** 01 · 08 · 09 · 10
**Menutup:** AC 68 · 70 · 93 · 94 · 95 · 96 *(6 AC)* — US 43 · 47 · 48

#### Hasil & nilai pengguna

Hari ini riwayat kontrak treaty ada di sistem lama, dan ⛔ sebagian nilainya **tidak dapat dibaca
tanpa membongkar dokumen teks**. ⚠️ Sebagian tanggalnya **ambigu secara mutlak**.

Sesudah tiket ini, data lama terbaca di sistem baru, ⭐ setiap penyimpangan dari perilaku lama
**tercatat beserta alasannya**, dan ⭐ setiap butir yang masih terbuka **tertulis di satu tempat** —
sehingga tidak ada yang terlupakan saat go-live.

#### Area codebase

- Alat migrasi: pemindahan data lama
- Uji paritas: perbandingan perilaku lama dan baru
- Uji dokumen: kelengkapan penanda dan rujukan

#### Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Bentuk data lama | dokumen teks hasil potret halaman kerja; ⛔ **bukan skema** |
| Tanggal ambigu | dua susunan tercampur, disimpan sebagai teks |
| Enam penyimpangan sadar | `spec.md` §10.1 |

#### ADR terkait

- ⭐ **ADR-U-0009** — seluruh data dipindahkan; ⛔ tidak ada koeksistensi dua penulis
- **ADR-U-0003** — uang tidak `float`

#### Acceptance criteria

- [ ] **AC 68** — data lama **terbaca** di sistem baru; berkas lama dapat dibuka
- [ ] **AC 70** — setiap penyimpangan dari perilaku lama **tercatat beserta alasannya**
- [ ] **AC 93** — setiap acceptance criterion merujuk **bab asalnya**
- [ ] **AC 94** — setiap acceptance criterion membawa **penanda**
- [ ] **AC 95** — ⛔ butir terbuka **tidak ditutup**; tidak ada yang dinyatakan selesai tanpa
      pemiliknya
- [ ] **AC 96** — ⛔ nilai berupa **nama orang** tidak muncul di artefak mana pun

#### Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| ⚠️ **15** | data lama bertanggal **ambigu mutlak** — satu nilai, dua arti | ⛔ **MENAHAN migrasi tanggal**, bukan tiketnya |
| ✅ ~~**P29**~~ | ~~contoh isi dokumen — dibutuhkan untuk membaca data lama~~ | ✅ **TERJAWAB 2026-09-22** — **tiga** contoh diterima |
| ⛔⛔ **galat uang tersimpan** | ekor `2,76 x 10^-7` pada `NetPremium`, sudah permanen di produksi | ⛔ **MENENTUKAN AMBANG PARITAS** — lihat catatan di bawah |

#### Perintah verifikasi

1. Pindahkan satu berkas lama, buka di sistem baru — ⭐ seluruh medan yang dipakai **terisi**.
2. Sisir spec dan tiket — ⭐ **nol** acceptance criterion tanpa penanda, ⭐ **nol** tanpa rujukan bab.
3. Sisir seluruh artefak — ⭐ **nol** nama orang, **nol** nomor polis apa adanya.
4. Bandingkan hasil penggolongan dan potongan pajak lama lawan baru — ⭐ **sama**.

#### Catatan

⚠️ **Paritas dapat diuji sekarang untuk yang sudah dibangun** — penggolongan, pajak, tanggal,
riwayat. ⛔ **Paritas rantai perhitungan uang menunggu tiket 13.**

---

#### ⛔⛔ Ambang paritas uang — tidak dapat ditetapkan sebelum satu keputusan turun

`[terbuka]` `[terverifikasi]` 2026-09-22 — ⚠️ **dicatat di sini karena tiket inilah yang memiliki
uji paritas**; butir aslinya ada di **P29**, dan **tidak ditutup di sini**.

Data produksi memuat **galat angka uang yang sudah tersimpan permanen**:

```
premium angsuran  148157378.220000069   x 4  =  592629512.880000276
NetPremium        592629512.880000276        <- sama persis
nilai bulat 2 desimal                        =  592629512.88
selisih                                      =  2,76 x 10^-7
```

⛔⛔ **Akibatnya langsung bagi tiket ini:** sistem baru memakai aritmetika desimal yang benar
*(ADR-U-0003, nol `float`)*, sehingga ia **tidak akan menghasilkan ekor itu**. ⛔ **Uji paritas yang
menuntut kesamaan persis akan GAGAL — justru karena sistem barunya benar.**

| Bila keputusan P29 jatuh ke | Ambang paritas yang berlaku |
| --- | --- |
| **a** ikuti apa adanya | ⛔ paritas **persis**; galat lama sengaja ditiru |
| **b** bulatkan saat migrasi | ⚠️ paritas **tidak berlaku** untuk angka historis — angkanya memang berubah |
| **c** simpan apa adanya, bulatkan saat dihitung | ⭐ paritas dengan **toleransi**, dan besarnya toleransi itu yang harus ditetapkan |

⛔ **Ambangnya tidak ditetapkan di sini.** Ia mengikuti keputusan `[work owner]` dan `[Finance]`
pada **P29**. ⚠️ **Yang tidak boleh terjadi:** uji paritas ditulis dengan ambang yang dikarang,
lalu selisihnya diam-diam dianggap wajar.

> ⚠️ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat di atas semula berakhir dengan
> ⛔ *"… yang menunggu P18."* **P18 ditarik.** Tiket 13 kini menunggu **P30** dan **P8**.

## NB Treaty In - 16 - Kerangka penyimpanan polis dan kunci generasi

**Status:** ready-for-agent
**Blocked by:** — *(dapat mulai segera)*
**Menutup:** NB AC **1–7** *(7 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-5..ID-10

#### Hasil & nilai pengguna

Polis treaty tersimpan sebagai baris di dalam tabel relasional, bukan sebagai dokumen teks.
Satu polis baru menempati **satu generasi**, dan generasi itu punya kunci yang mencegah dua berkas
menempati tempat yang sama.

⭐ **Ini tiket pertama rantai penyimpanan** — seluruh tiket berikutnya menulis ke kerangka ini.

#### Yang dibangun

Kerangka tabel akar dan tabel generasi, berbagi kunci utama, beserta ketiga aturan keunikannya:

| Aturan | Yang dicegahnya |
| --- | --- |
| kunci alami `(NOPOLIS, PRODKE)` unik | dua berkas menempati generasi yang sama |
| `OLD_POLIS_ID` unik | satu generasi punya dua penerus — rantai menjadi pohon |
| baris generasi tertutup tidak dapat disunting | nilai yang sudah jadi dasar selisih berubah diam-diam |

⭐ Pada polis baru, `PRODKE = 0` dan `OLD_POLIS_ID` kosong. Tabel akar **lintas-lini** — dipakai
bersama modul lain, dan baris antar modul **tidak saling menunjuk**.

#### Batas — yang TIDAK termasuk

⛔ Tipe dan presisi kolom — tiket **18**.
⛔ Isi tabel anak — tiket **19**.
⛔ Baris endorsemen `PRODKE ≥ 1` — tiket **24**.
⛔ `CREATE TABLE` dan presisi fisik — dicocokkan DBA **di dalam** tiket ini, bukan prasyaratnya.

#### Cara mengujinya

Lewat seam `repository`. Menyimpan dua polis dengan kunci alami sama harus **ditolak oleh basis
data**, bukan oleh pemeriksaan di aplikasi — test memverifikasi penolakan itu datang dari lapisan
penyimpanan.

#### Acceptance criteria

- [ ] **AC 1** — menyimpan dua polis dengan `NOPOLIS` dan `PRODKE` sama **ditolak**
- [ ] **AC 2** — polis baru tersimpan dengan `PRODKE = 0`
- [ ] **AC 3** — `OLD_POLIS_ID` pada polis baru bernilai kosong
- [ ] **AC 4** — dua baris tidak boleh berbagi `OLD_POLIS_ID` yang sama
- [ ] **AC 5** — tabel generasi dan tabel akar **berbagi kunci utama**, tanpa kolom kunci tamu terpisah
- [ ] **AC 6** — menyunting baris generasi yang sudah ditutup **ditolak**
- [ ] **AC 7** — baris antar modul pada tabel akar **tidak saling menunjuk**

## NB Treaty In - 17 - Nomor urut baris anak dan pemasangan antar generasi

**Status:** ready-for-agent
**Blocked by:** **16**
**Menutup:** NB AC **8–11** *(4 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-11..ID-13

#### Hasil & nilai pengguna

Setiap baris rincian — angsuran, penyebaran, lapisan — membawa nomor urut di dalam induknya.
Nomor itulah yang kelak memasangkan baris generasi baru dengan baris generasi lama, sehingga selisih
per baris berarti.

⚠️ **Kunci dagang sengaja TIDAK dipakai untuk memasangkan.** Ia turun pangkat menjadi
**pemeriksa** — bila pasangan menurut nomor urut ternyata berbeda kunci dagangnya, itu ditandai,
bukan dipakai memasangkan ulang.

#### Yang dibangun

Kolom nomor urut pada **setiap** tabel anak, unik di dalam satu induk.

⭐ **Pada polis baru, baris boleh dihapus dan nomornya dinomori ulang rapat `1..n`** — aman karena
belum ada generasi untuk dipasangkan. ⚠️ Begitu generasi ditutup, nomornya **beku**.
⛔ Aturan berbeda berlaku di endorsemen — tiket **24**.

#### Batas — yang TIDAK termasuk

⛔ Perilaku nomor urut di endorsemen — tiket **24**.
⛔ Penanda pasangan bergeser — tiket **28**.

#### Cara mengujinya

Lewat seam `repository`. Uji utama: polis dengan tiga baris angsuran, baris kedua dihapus, lalu
dibaca kembali — nomor urutnya harus rapat, bukan berlubang.

#### Acceptance criteria

- [ ] **AC 8** — setiap tabel anak memiliki kolom nomor urut
- [ ] **AC 9** — menghapus baris di tengah menghasilkan penomoran ulang yang rapat
- [ ] **AC 10** — nomor urut unik di dalam satu induk
- [ ] **AC 11** — pemasangan antar generasi memakai nomor urut, **bukan** kunci dagang

## NB Treaty In - 18 - Tipe kolom, konversi masuk, dan presisi uang

> ## ⭐ PENAHAN GUGUR — 23 September 2026 sore
>
> `[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."* ⭐ Presisi dinaikkan ke **`NUMBER(38,8)`** — **30 digit di depan koma**, delapan di belakang, batas tertinggi Oracle. Dasarnya sapuan korpus: ambang dagang nyata sudah **tepat di batas** 12 digit *(`181500000000.00` · `150000000000.00`)*, dan ada sentinel **17 digit** *(`99999999999999999.99`)*. Penjumlahan lintas mata uang dapat melewati keduanya. ⭐ `NUMBER` di Oracle berpanjang **berubah-ubah** — hanya digit bermakna yang tersimpan, sehingga pelebaran ini **tidak memakan ruang tambahan**. Butir ditutup oleh bukti, bukan oleh DBA.
>
> ⭐ **`blocked` → `ready-for-agent`.** ⛔ **Nol butir `[data DBA]` tersisa di tiket ini.**
>
> Rinciannya: `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026 sore)*
~~**Blocked by:** **16** · ⛔ `[data DBA]` **presisi fisik belum diuji terhadap nilai terbesar** — dua belas digit di depan koma belum dibuktikan cukup~~ ⛔ **gugur 23-09-2026 sore**
**Menutup:** NB AC **12–25** *(14 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-14..ID-20

#### Hasil & nilai pengguna

Nilai yang masuk dari dokumen lama seluruhnya bertipe teks — termasuk uang, tanggal, dan penanda.
Tiket ini menetapkan bagaimana teks itu menjadi kolom bertipe, **sekali saat masuk**, bukan setiap
kali dibaca.

⛔⛔ **Yang paling mudah salah, dan akibatnya besar:** kode `"006"` yang disimpan sebagai bilangan
lalu dibaca balik menjadi `"6"` akan **lolos uji pulang-pergi** tetapi memecahkan penggolong jenis
usaha 36 baris — dan kegagalannya diam.

#### Yang dibangun

| Golongan | Ketetapan |
| --- | --- |
| **uang dan persen** | angka presisi tetap; ⛔ **tidak pernah** bilangan mengambang |
| **kode** | ⭐ **tetap teks** — nol di depan membawa makna |
| **penanda** | ⭐ **tetap teks** — kosong adalah keadaan sah yang **berbeda** dari nol |
| **tanggal** | dua format masuk: delapan digit, dan cap waktu bersufiks zona |
| **teks kosong** | menjadi **kosong**, bukan nol |

⚠️ Pembandingan dua nilai uang memakai **toleransi atau bentuk terbulatkan**, bukan kesamaan
persis — data produksi terbukti membawa galat pecahan.

#### Batas — yang TIDAK termasuk

⛔ `CREATE TABLE` — presisi fisik dicocokkan DBA **di dalam** tiket ini.
⛔ Pemecahan dokumen menjadi baris — tiket **19**.

#### Cara mengujinya

Lewat seam `repository`, dan ⚠️ **sebagian test memeriksa nilai kolom langsung** — pulang-pergi
saja tidak cukup, sebab tulis dan baca yang sama-sama salah simetris tetap hijau.

⭐ Uji yang wajib: kode `"006"` disimpan lalu **dibaca dari kolomnya**, dan penggolong jenis usaha
dijalankan atasnya — hasil bawaan berarti gagal.

#### Acceptance criteria

- [ ] **AC 12** — kode tiga digit berawalan nol tersimpan dan terbaca utuh
- [ ] **AC 13** — kode dua digit berawalan nol tersimpan utuh
- [ ] **AC 14** — penggolong jenis usaha menemukan barisnya, bukan nilai bawaan
- [ ] **AC 15** — penanda kosong tersimpan kosong, **bukan** nol dan bukan tak-bernilai
- [ ] **AC 16** — penanda bernilai nol menghasilkan keputusan ditolak; nilai lain disetujui
- [ ] **AC 17** — teks kosong pada medan uang tersimpan tak-bernilai
- [ ] **AC 18** — teks kosong pada medan tanggal tersimpan tak-bernilai
- [ ] **AC 19** — nilai uang berdesimal sembilan **dibulatkan pada desimal kedelapan**, bukan dipotong ke dua
- [ ] **AC 20** — kolom uang berskala **delapan desimal**, ⭐ **tiga puluh digit di depan koma** *(`NUMBER(38,8)`; semula ~~dua belas~~ — dinaikkan 23-09-2026 sore)*
- [ ] **AC 21** — tanggal delapan digit terurai benar
- [ ] **AC 22** — cap waktu bersufiks zona terurai benar
- [ ] **AC 23** — pengurutan menurut tanggal menghasilkan urutan kronologis, bukan leksikal
- [ ] **AC 24** — nol kolom uang bertipe mengambang
- [ ] **AC 25** — pembandingan uang memakai toleransi, bukan kesamaan persis

#### ⛔ Kenapa tiket ini `blocked`

~~`[data DBA]` **Dua belas digit di depan koma belum diuji terhadap nilai terbesar.**~~ ⛔ **BUTIR GUGUR 23-09-2026 sore.** `[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."* Presisi dinaikkan ke **`NUMBER(38,8)`** — tiga puluh digit di depan koma. Dasarnya sapuan korpus: ambang dagang nyata **tepat di batas** dua belas digit *(`181500000000.00`)*, dan ada sentinel **tujuh belas digit** *(`99999999999999999.99`)*. ⭐ `NUMBER` di Oracle berpanjang berubah-ubah, jadi pelebaran ini **tidak memakan ruang tambahan**.

⚠️ **Dan satu pertentangan di dalam spec, dicatat di sini karena spec tidak boleh disunting:**
~~`[terverifikasi]` spec EDM **AC 49** menuntut sembilan desimal, bertentangan dengan AC 58.~~ ✅ **DITUTUP 23-09-2026 sore.** AC 49 diselaraskan ke **delapan**; bunyi lamanya dikutip di spec EDM. ⭐ Pertentangan ini ditemukan **pelaksana ronde tiket**, dan laporannya terbukti tepat..

## NB Treaty In - 19 - Pemecah dokumen menjadi baris

> ## ⭐ PENAHAN GUGUR — 23 September 2026 sore
>
> `[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`, ikuti dari data yang digunakan di Activity dan Section."* ⭐ Daftar medan disusun ulang dari **302 berkas aturan** — **394 medan unik**, 232 di antaranya tampil di layar. Hasilnya di **`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`**. Data guide turun derajat jadi **penambal**, sah hanya untuk panjang maksimum per medan.
>
> ⭐ **`blocked` → `ready-for-agent`.** ⛔ **Nol butir `[data DBA]` tersisa di tiket ini.**
>
> Rinciannya: `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026 sore)*
~~**Blocked by:** **16** · **17** · **18** · ⛔ `[data DBA]` **daftar kolom lengkap** — panduan bentuk dokumen dari DBA terbukti **basi**~~ ⛔ **gugur 23-09-2026 sore**
**Menutup:** NB AC **26–44** *(19 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-21..ID-31

#### Hasil & nilai pengguna

Dokumen polis lama dipecah menjadi baris tabel: data umum, kuotasi, ceding, angsuran, penyebaran,
lapisan, dan riwayat usulan. Sesudah tiket ini, bentuk data polis **dapat dinyatakan** — bukan
potret halaman kerja.

#### Yang dibangun

Pemecah dokumen dan penyusun baris, beserta aturan isi tiap tabel:

- **data umum** — medan skalar tingkat atas ditambah kolom datar dari tabel dokumen lama;
  ⛔ penanda lapisan **tidak** disimpan di sini, sebab di sistem lama ia pantulan baris pertama saja
- **ceding** — satu baris per ceding, id dan nama terpisah; ⭐ bentuk gabungan di data umum
  **disalin apa adanya**, tidak pernah dirangkai ulang
- **angsuran** — ⭐ **satu tingkat** pada bentuk proporsional
- **penyebaran** — pembagian rata menurut cacah baris, presisi sepuluh
- **lapisan** — pemegang penanda lapisan; ⚠️ potongan di sini adalah **uang**, bukan persen
- **riwayat usulan** — tabel yang sudah datar, ⛔ **tidak dibuat ulang**

⭐ **Penampung medan tak dikenal** disediakan — medan dokumen yang tidak dikenal **disimpan**,
bukan dibuang.

#### Batas — yang TIDAK termasuk

⛔ Uji pulang-pergi dan dua bentuk dokumen — tiket **21**.
⛔ Pemuatan dokumen lama secara massal — tiket **22**.
⛔ Pohon objek pertanggungan — milik modul fakultatif, **di luar lingkup**.

#### Cara mengujinya

Lewat seam `repository`. ⚠️ **Penampung medan tak dikenal wajib kosong** sebelum pekerjaan
dinyatakan selesai — isinya yang tidak kosong berarti sensus medan belum lengkap.

#### Acceptance criteria

- [ ] **AC 26** — penanda lapisan **tidak** ada di tabel data umum
- [ ] **AC 27** — tabel kuotasi menyimpan sepuluh medan, termasuk jenis kontrak dan kode bisnis
- [ ] **AC 28** — satu baris per ceding, id dan nama terpisah
- [ ] **AC 29** — bentuk gabungan di data umum tersimpan **persis** seperti di dokumen
- [ ] **AC 30** — menghapus satu ceding menghapus **barisnya**, bukan menyunting teks gabungan
- [ ] **AC 31** — polis proporsional tersimpan **tanpa** baris rincian angsuran bertingkat
- [ ] **AC 32** — polis proporsional tersimpan **tanpa** baris lapisan
- [ ] **AC 33** — menyimpan baris lapisan pada polis proporsional **ditolak**
- [ ] **AC 34** — penanda lapisan tersimpan di tabel lapisan, bukan di data umum
- [ ] **AC 35** — potongan pada lapisan bertipe **uang**
- [ ] **AC 36** — penyebaran pada polis baru memakai presisi **sepuluh**
- [ ] **AC 37** — tiga medan bagian tersimpan sebagai **persentase**
- [ ] **AC 38** — medan potongan dan total bagian tersimpan sebagai **persentase**
- [ ] **AC 39** — tabel riwayat usulan **tidak dibuat ulang**
- [ ] **AC 40** — keterangan usulan dipotong pada batas panjangnya
- [ ] **AC 41** — keputusan usulan diturunkan dari penandanya, per baris
- [ ] **AC 42** — penanda per baris usulan tersimpan **terpisah** dari penanda tingkat polis
- [ ] **AC 43** — kolom pelaku terisi dari identitas login
- [ ] **AC 44** — kolom penanggung jawab terisi dari nama tampilan

#### ⛔ Kenapa tiket ini `blocked`

~~`[data DBA]` **Panduan bentuk dokumen dari DBA terbukti basi.**~~ ⛔ **BUTIR GUGUR 23-09-2026 sore.** `[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`, ikuti dari data yang digunakan di Activity dan Section."*

⭐ Daftar medan disusun ulang dari **302 berkas aturan** — **394 medan unik**, **232** di antaranya tampil di layar. Sumbernya: **`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`**. Sapuan itu menutup dua titik buta yang sebelumnya meleset empat kali: **dua konvensi tag yang berlawanan**, dan **rujukan relatif di dalam loop**.

⚠️ Data guide turun derajat jadi **penambal** — sah untuk panjang maksimum per medan, dan untuk 24 nama yang ditulis sistem sehingga tidak tersapu aturan *(`OldData` · `ProdKe` · `EDMNo` · `OldPolicyNo` · `BranchCode`)*.

⛔ **394 tetap batas bawah, bukan total.** Penampung medan tak dikenal tetap wajib, dan wajib **kosong** sebelum pekerjaan dinyatakan selesai.

⭐ **Yang dapat dikerjakan lebih dulu tanpa menunggu:** kerangka pemecah, penampung medan tak
dikenal, dan seluruh AC yang tidak bergantung pada daftar kolom — ⚠️ tetapi **pekerjaan tidak
dapat dinyatakan selesai** sampai penampung itu kosong.

## NB Treaty In - 20 - Transaksi tunggal dan skema eksplisit

**Status:** ready-for-agent
**Blocked by:** **16** · **19**
**Menutup:** NB AC **45–48** *(4 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-32..ID-35

#### Hasil & nilai pengguna

Satu polis tersimpan **seluruhnya atau tidak sama sekali**. Sistem lama meng-`COMMIT` di dalam
empat program basis data yang berbeda, sehingga kegagalan di tengah meninggalkan data separuh yang
permanen.

⚠️ `[penyimpangan sadar]` **Ini lebih ketat daripada sistem lama, dan disengaja.**

#### Yang dibangun

Pembungkusan seluruh urutan penyimpanan satu polis — induk dan seluruh anaknya — dalam **satu
transaksi**. Dan dua aturan penulisan query:

| Aturan | Sebab |
| --- | --- |
| skema basis data ditulis **eksplisit** pada setiap pernyataan | sistem lama tidak konsisten; satu skema kedua terbukti disentuh |
| ⛔ **nol pemanggilan program tersimpan** dari aplikasi | mandat proyek; logikanya ditulis ulang |

#### Batas — yang TIDAK termasuk

⛔ Transaksi yang mencakup tabel proyeksi selisih — tiket **26**.
⛔ Pemuatan massal — tiket **22**.

#### Cara mengujinya

Lewat seam `repository`. Uji utama: kegagalan disuntikkan pada penulisan tabel anak terakhir —
tidak boleh ada satu baris pun tertinggal, termasuk baris induknya.

#### Acceptance criteria

- [ ] **AC 45** — kegagalan menulis tabel anak mana pun **membatalkan seluruh** penyimpanan
- [ ] **AC 46** — seluruh penyimpanan satu polis terjadi dalam **satu transaksi**
- [ ] **AC 47** — setiap pernyataan menyebut skema **eksplisit**
- [ ] **AC 48** — penyimpanan **tidak memanggil** satu pun program tersimpan

## NB Treaty In - 21 - Pulang-pergi dan dua bentuk dokumen

**Status:** ready-for-agent
**Blocked by:** **18** · **19**
**Menutup:** NB AC **49–54** *(6 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` Bab 7

#### Hasil & nilai pengguna

Dokumen yang disimpan lalu dibaca kembali menghasilkan nilai yang sama — pada **kedua** bentuk
dokumen, bukan hanya yang kebetulan diuji lebih dulu.

⛔⛔ **Bentuk dokumen tidak seragam.** Daftar angsuran **bersarang** pada satu bentuk dan **datar**
pada bentuk lain; keduanya nyata di sistem lama. Pengurai yang menganggapnya seragam akan patah.

#### Yang dibangun

Rangkaian uji pulang-pergi yang mencakup kedua bentuk, dan **aturan kecukupan cakupan**:
⛔ test suite yang hanya mencakup satu bentuk **dinyatakan tidak memadai** dan tidak boleh dihitung
sebagai cakupan.

⚠️ **Pulang-pergi saja tidak cukup.** Bila tulis dan baca sama-sama salah secara simetris, ujinya
tetap hijau — sebagian test karena itu memeriksa **nilai kolom langsung**.

#### Batas — yang TIDAK termasuk

⛔ Bentuk dokumen pada endorsemen — tiket **27**.
⛔ Contoh dokumen produksi **tidak disimpan** di dalam repositori; yang disimpan cetakan bentuknya.

#### Cara mengujinya

Lewat seam `repository`, dengan data buatan — ⛔ **bukan cuplikan produksi**.

#### Acceptance criteria

- [ ] **AC 49** — polis proporsional pulang-pergi menghasilkan nilai yang sama
- [ ] **AC 50** — polis non-proporsional pulang-pergi menghasilkan nilai yang sama
- [ ] **AC 51** — urutan baris anak saat dibaca sama dengan urutan nomor urutnya
- [ ] **AC 52** — bentuk daftar angsuran **datar** terurai benar
- [ ] **AC 53** — bentuk daftar angsuran **bersarang** terurai benar
- [ ] **AC 54** — cakupan yang hanya satu bentuk dinyatakan **tidak memadai**

## NB Treaty In - 22 - Pemuat dokumen lama

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Seluruh dokumen polis dipindahkan — setiap polis, setiap generasinya. Tidak ada penyaringan.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 5.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** **19** · **20** · ⛔ `[work owner]` **dokumen lama dipindahkan seluruhnya atau sebagian** — belum diputuskan~~ ⛔ **penahan gugur 23-09-2026**
**Menutup:** NB AC **55–59** *(5 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-3

#### Hasil & nilai pengguna

Dokumen polis lama dimuat ke dalam tabel baru **lewat antarmuka yang sama** dengan jalur biasa,
sehingga data lama melewati pemeriksaan yang sama dan tidak ada pintu belakang.

⭐⭐ **Nilai lama tidak pernah dihitung ulang** — disalin apa adanya, lengkap dengan galat
presisinya, supaya laporan lama masih dapat direkonsiliasi.

#### Yang dibangun

Pemuat massal yang: membaca dokumen lama · memecahnya dengan pemecah yang sama · menulis lewat
antarmuka penyimpanan yang sama · dan **melaporkan** dokumen yang gagal diurai beserta sebabnya.

⛔ **Dokumen yang gagal tidak dilewati diam-diam.**

#### Batas — yang TIDAK termasuk

⛔ Penanda migrasi khas endorsemen — tiket **28**.
⛔ Keputusan lingkup pemindahan — `[work owner]`, lihat di bawah.

#### Cara mengujinya

Lewat seam `repository` yang sama. ⭐ Uji utama: dokumen bergalat presisi dimuat, lalu nilainya
**dibaca dari kolomnya** — pembulatan ke presisi mata uang berarti gagal.

#### Acceptance criteria

- [ ] **AC 55** — dokumen lama dimuat **tanpa pembulatan ke presisi mata uang**
- [ ] **AC 56** — pemuat menulis lewat antarmuka penyimpanan yang **sama**
- [ ] **AC 57** — medan tak dikenal **tersimpan di penampung**, bukan dibuang
- [ ] **AC 58** — dokumen yang gagal diurai **dilaporkan beserta sebabnya**
- [ ] **AC 59** — penampung medan tak dikenal **wajib kosong** sebelum pekerjaan dinyatakan selesai

#### ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum diputuskan apakah dokumen lama dipindahkan seluruhnya, sebagian, atau tetap
dibaca lewat jalur lama.** Lingkup pemuat bergantung langsung padanya — memuat seluruh riwayat dan
memuat dua tahun terakhir adalah pekerjaan yang berbeda besarnya.

⚠️ **AC 59 tidak dapat dipenuhi sebelum tiket 19 lepas** — penampung medan tak dikenal tidak akan
kosong selama daftar kolom lengkap belum ada.

## NB Treaty In - 23 - Disiplin berkas, kerahasiaan, dan arah ketergantungan

**Status:** ready-for-agent
**Blocked by:** — *(dapat mulai segera)*
**Menutup:** NB AC **60–66** *(7 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-1 · ID-21 · ID-25 · Bab 9

#### Hasil & nilai pengguna

Pekerjaan penyimpanan tidak membocorkan data nasabah ke dalam repositori, dan susunan lapisannya
tidak terbalik.

⭐ **Tiket ini tanpa penahan dan dapat dikerjakan sejajar dengan tiket mana pun** — sebagian besar
isinya penjaga otomatis, bukan fitur.

#### Yang dibangun

Penjaga otomatis yang gagal bila dilanggar:

| Penjaga | Yang dicegahnya |
| --- | --- |
| nol nama orang di berkas proyek | nilai produksi tersalin ke spec, test, atau rancangan |
| nol nomor polis harfiah | pengenal nasabah masuk repositori |
| test berjalan di atas **data buatan** | cuplikan produksi menjadi fixture |
| perubahan skema menuntut persetujuan manusia | perubahan bentuk tabel berjalan otomatis |
| arah ketergantungan tidak terbalik | lapisan penyimpanan memanggil lapisan layanan |

Ditambah dua pernyataan bentuk yang diuji langsung: cacah medan tabel data umum, dan penerimaan
bahwa bentuk gabungan ceding **boleh tidak sinkron** dengan barisnya.

#### Batas — yang TIDAK termasuk

⛔ Isi tabel — tiket **19**.
⚠️ Cacah medan yang diuji AC 64 adalah **batas bawah**, bukan total.

#### Cara mengujinya

Sebagian lewat seam `repository`, sebagian lewat pemeriksaan repositori otomatis.

#### Acceptance criteria

- [ ] **AC 60** — nol nama orang tersalin ke berkas rancangan, spec, maupun test
- [ ] **AC 61** — nol nomor polis ditulis apa adanya
- [ ] **AC 62** — test berjalan di atas **data buatan**
- [ ] **AC 63** — perubahan skema memerlukan **persetujuan manusia**
- [ ] **AC 64** — tabel data umum menampung cacah medan yang ditetapkan spec
- [ ] **AC 65** — bentuk gabungan ceding **boleh tidak sinkron** dengan barisnya, tanpa penjaga
- [ ] **AC 66** — arah ketergantungan **tidak pernah dibalik**

## NB Treaty In - 24 - Rantai generasi endorsemen dan pembatalan [DIGANTIKAN]

> ## ⛔⛔ DIGANTIKAN — 23 September 2026
>
> **Tiket ini tidak dikerjakan.** Ia digantikan oleh
> **`.scratch\edm-treaty-in\issues\01 · 02 · 03 · 05`**.
>
> ### Kenapa
>
> Brief ronde tiket sebelumnya berlingkup **dua spec sekaligus**, sehingga pekerjaan **endorsemen**
> jatuh ke folder tiket **polis baru**. Ronde berikutnya membaginya ulang ke folder yang benar,
> lebih halus, dan — yang menentukan — **setiap tiket menyebut tiket NB mana yang membuat
> tabelnya**. Tanpa gate itu, tiket endorsemen bisa dikerjakan sebelum tabelnya ada.
>
> ⚠️ **Kekeliruan lingkup ini milik penyusun brief, bukan pelaksana ronde ini.**
>
> ⭐ Berkas ini **tidak dihapus** supaya jejaknya tidak hilang — tanpa ini, nomor tiket NB akan
> tampak melompat dari 23 ke 29 tanpa sebab.
>
> ### Yang berlaku
>
> | | |
> | --- | --- |
> | Tiket penyimpanan **NB** | `nb-treaty-in\issues\` **16–23** |
> | Tiket penyimpanan **EDM** | `edm-treaty-in\issues\` **01–11** |
>
> ⛔ Isi di bawah dibiarkan apa adanya sebagai catatan, **bukan sebagai pekerjaan.**

---


**Status:** ~~ready-for-agent~~ — **DIGANTIKAN**
**Blocked by:** **16** · **17**
**Menutup:** EDM AC **1–15** *(15 AC)*
**Sumber:** `edm-treaty-in\spec-penyimpanan-relasional.md` ID-8..ID-19 · ID-41

#### Hasil & nilai pengguna

Endorsemen tersimpan sebagai **generasi berikutnya** dari polis yang sama — bukan sebagai jenis
berkas tersendiri. Rantainya lurus: satu generasi, satu penerus.

⛔⛔ **Ini memperbaiki bahaya nyata.** `[terverifikasi]` Sistem lama membaca nomor generasi terakhir
**tanpa penguncian**, lalu menambah satu. Dua endorsemen serentak atas polis yang sama membaca angka
yang sama, dan **keduanya tersimpan**. ⭐ Di sistem baru yang kedua **gagal**.

#### Yang dibangun

Baris generasi dengan nomor ≥ 1 dan penunjuk ke generasi sebelumnya, beserta empat aturan:

| Aturan | Yang dicegahnya |
| --- | --- |
| penunjuk generasi **unik** | rantai menjadi pohon; *"selisih mana yang benar"* tidak terjawab |
| kunci alami **unik** | dua endorsemen serentak mendapat nomor sama |
| generasi lampau **tidak dapat disunting** | dasar selisih berubah sesudah disetujui |
| ⭐ **aturan keutuhan** | generasi baru **wajib memuat setiap nomor urut** milik generasi sebelumnya |

⭐ **Pada endorsemen baris rincian TIDAK dapat dihapus**, dan nomor urut lama terbawa apa adanya —
⚠️ **berbeda dari polis baru**, tempat penomoran ulang dibolehkan.

⭐ **Pembatalan adalah jenis endorsemen**, dipilih di awal — ia menghasilkan generasi bernilai nol,
⛔ **bukan penghapusan baris**.

⭐ **Nilai lama tidak menjadi tabel.** Ia baris yang ditunjuk penunjuk generasi.

#### Batas — yang TIDAK termasuk

⛔ Perhitungan selisih — tiket **25**.
⛔ Tabel proyeksi — tiket **26**.
⛔ Kolom khas endorsemen — tiket **27**.

#### Cara mengujinya

Lewat seam `repository`. ⭐ **Tiga uji yang wajib berdiri sendiri:** rantai **tiga generasi** ·
penolakan penerus kedua · penolakan generasi yang kehilangan satu nomor urut.
⚠️ Rantai dua generasi **tidak memadai** — selisih terhadap generasi tepat sebelumnya baru terbukti
mulai generasi ketiga.

#### Acceptance criteria

- [ ] **AC 1** — baris endorsemen tersimpan dengan nomor generasi ≥ 1
- [ ] **AC 2** — penunjuk generasi **terisi** dan menunjuk generasi sebelumnya
- [ ] **AC 3** — dua baris tidak boleh berbagi penunjuk generasi yang sama
- [ ] **AC 4** — dua endorsemen serentak: yang kedua **ditolak**
- [ ] **AC 5** — nomor endorsemen berbentuk nomor polis + pemisah + dua digit
- [ ] **AC 6** — menyunting generasi yang sudah punya penerus **ditolak**
- [ ] **AC 7** — selisih dihitung terhadap generasi **tepat sebelumnya**
- [ ] **AC 8** — generasi yang kehilangan satu nomor urut **ditolak**
- [ ] **AC 9** — nol tabel salinan nilai lama; ia dibaca lewat penunjuk generasi
- [ ] **AC 10** — pada endorsemen, baris rincian **tidak dapat dihapus**
- [ ] **AC 11** — nomor urut lama **terbawa apa adanya**
- [ ] **AC 12** — baris baru mendapat nomor urut **maksimum + 1**
- [ ] **AC 13** — pasangan yang bergeser diperlakukan sebagai **anomali**
- [ ] **AC 14** — pembatalan menghasilkan generasi bernilai **nol**, bukan penghapusan
- [ ] **AC 15** — jenis endorsemen tersimpan di kolomnya sendiri, dipilih di awal

## NB Treaty In - 25 - Perhitungan selisih di lapisan layanan [DIGANTIKAN]

> ## ⛔⛔ DIGANTIKAN — 23 September 2026
>
> **Tiket ini tidak dikerjakan.** Ia digantikan oleh
> **`.scratch\edm-treaty-in\issues\06`**.
>
> ### Kenapa
>
> Brief ronde tiket sebelumnya berlingkup **dua spec sekaligus**, sehingga pekerjaan **endorsemen**
> jatuh ke folder tiket **polis baru**. Ronde berikutnya membaginya ulang ke folder yang benar,
> lebih halus, dan — yang menentukan — **setiap tiket menyebut tiket NB mana yang membuat
> tabelnya**. Tanpa gate itu, tiket endorsemen bisa dikerjakan sebelum tabelnya ada.
>
> ⚠️ **Kekeliruan lingkup ini milik penyusun brief, bukan pelaksana ronde ini.**
>
> ⭐ Berkas ini **tidak dihapus** supaya jejaknya tidak hilang — tanpa ini, nomor tiket NB akan
> tampak melompat dari 23 ke 29 tanpa sebab.
>
> ### Yang berlaku
>
> | | |
> | --- | --- |
> | Tiket penyimpanan **NB** | `nb-treaty-in\issues\` **16–23** |
> | Tiket penyimpanan **EDM** | `edm-treaty-in\issues\` **01–11** |
>
> ⛔ Isi di bawah dibiarkan apa adanya sebagai catatan, **bukan sebagai pekerjaan.**

---


**Status:** ~~ready-for-agent~~ — **DIGANTIKAN**
**Blocked by:** **24**
**Menutup:** EDM AC **16–22** *(7 AC)*
**Sumber:** `edm-treaty-in\spec-penyimpanan-relasional.md` ID-4 · ID-28..ID-32

#### Hasil & nilai pengguna

Selisih antara generasi dihitung **di lapisan layanan**, dengan **satu rumus** untuk semua kasus —
endorsemen pertama maupun berlapis, proporsional maupun non-proporsional.

```
selisih.X = baris_ini.X − baris(generasi sebelumnya).X
```

⚠️ `[penyimpangan sadar]` **Sistem lama punya DUA rumus di sisi proporsional**, dan yang kedua
mengurangi terhadap **selisih** generasi lampau, bukan terhadap **nilainya**. Dengan angka contoh:
nilai baru 180, nilai lama 150, selisih lama 50 ⇒ varian kedua menghasilkan **130**, padahal yang
benar **30**. ⛔ **Varian kedua tidak ditiru** — `[keputusan work owner]`.

#### Yang dibangun

Perhitungan selisih di lapisan layanan, dengan empat aturan turunan:

| Golongan medan | Perlakuan |
| --- | --- |
| **uang** | dikurangi |
| **persentase** | ⭐ **disalin**, tidak dikurangi |
| **kunci** | disalin |
| **arah utang-piutang** | diturunkan dari **tanda** selisih |

⭐ Pada selisih lapisan, arah utang-piutang diturunkan dari **jumlah sepanjang daftar lapisan**,
bukan dari satu lapisan.

⭐ Lapisan penyimpanan hanya mengambil **dua baris** — baris ini dan baris yang ditunjuknya.
⛔ Tidak ada agregasi di basis data.

#### Batas — yang TIDAK termasuk

⛔ Penulisan hasilnya ke tabel proyeksi — tiket **26**.
⛔ Penanda migrasi untuk baris hasil rumus lama — tiket **28**.

#### Cara mengujinya

Lewat seam `repository` untuk operannya, dan lewat lapisan layanan untuk rumusnya.
⭐ **Uji regresi wajib:** sisi non-proporsional harus menghasilkan angka yang **sama** sebelum dan
sesudah penyeragaman — ia memang sudah seragam, dan penyeragaman tidak boleh mengubahnya.

#### Acceptance criteria

- [ ] **AC 16** — satu rumus untuk **seluruh** generasi
- [ ] **AC 17** — medan **uang** berisi hasil pengurangan
- [ ] **AC 18** — medan **persentase** disalin, tidak dikurangi
- [ ] **AC 19** — medan **kunci** disalin apa adanya
- [ ] **AC 20** — arah utang-piutang diturunkan dari **tanda** selisih
- [ ] **AC 21** — pada selisih lapisan, arah itu diturunkan dari **jumlah** sepanjang daftar
- [ ] **AC 22** — sisi non-proporsional menghasilkan angka yang **sama** sebelum dan sesudah

#### ⚠️ Satu hal yang tidak terbaca dari sistem lama

`[terbuka]` Kedua blok rumus di sistem lama **sama-sama bergerbang aktif**, dan pemilihnya hanya
tertulis di **keterangan langkah**, bukan di gerbang yang dijalankan. ⇒ Mana yang benar-benar
berjalan **tidak dapat dinyatakan dari ekspor**.

⭐ **Tidak menahan tiket ini** — varian kedua tidak ditiru apa pun jawabannya.

## NB Treaty In - 26 - Proyeksi selisih yang dapat dibaca langsung [DIGANTIKAN]

> ## ⛔⛔ DIGANTIKAN — 23 September 2026
>
> **Tiket ini tidak dikerjakan.** Ia digantikan oleh
> **`.scratch\edm-treaty-in\issues\07`**.
>
> ### Kenapa
>
> Brief ronde tiket sebelumnya berlingkup **dua spec sekaligus**, sehingga pekerjaan **endorsemen**
> jatuh ke folder tiket **polis baru**. Ronde berikutnya membaginya ulang ke folder yang benar,
> lebih halus, dan — yang menentukan — **setiap tiket menyebut tiket NB mana yang membuat
> tabelnya**. Tanpa gate itu, tiket endorsemen bisa dikerjakan sebelum tabelnya ada.
>
> ⚠️ **Kekeliruan lingkup ini milik penyusun brief, bukan pelaksana ronde ini.**
>
> ⭐ Berkas ini **tidak dihapus** supaya jejaknya tidak hilang — tanpa ini, nomor tiket NB akan
> tampak melompat dari 23 ke 29 tanpa sebab.
>
> ### Yang berlaku
>
> | | |
> | --- | --- |
> | Tiket penyimpanan **NB** | `nb-treaty-in\issues\` **16–23** |
> | Tiket penyimpanan **EDM** | `edm-treaty-in\issues\` **01–11** |
>
> ⛔ Isi di bawah dibiarkan apa adanya sebagai catatan, **bukan sebagai pekerjaan.**

---


**Status:** ~~blocked~~ — **DIGANTIKAN**
**Blocked by:** **20** · **25** · ⛔ `[work owner]` **anak proyeksi dibutuhkan pembaca SQL atau tidak** — belum dijelaskan
**Menutup:** EDM AC **23–34** *(12 AC)*
**Sumber:** `edm-treaty-in\spec-penyimpanan-relasional.md` ID-6 · ID-7 · ID-23..ID-27 · ID-20..ID-22

#### Hasil & nilai pengguna

Angka selisih disimpan sebagai **proyeksi** yang dapat dibaca langsung dengan SQL, di luar layar
aplikasi — `[keputusan work owner]` pembaca seperti itu memang ada.

⛔⛔ **Rumusnya tetap hanya di satu tempat.** View yang menghitung selisih **dibatalkan**: rumus di
dua tempat — di aplikasi dan di basis data — cepat atau lambat bercabang.

#### Yang dibangun

Tabel proyeksi selisih beserta **tiga aturan yang membuatnya bukan tabel sumber**:

1. ia **hanya ditulis** lapisan aplikasi, dalam transaksi yang sama dengan generasinya;
2. ia **boleh dihapus total dan dibangun ulang** — ⛔ **hanya baris yang dihitung sistem baru**;
3. bila isinya berbeda dari hasil hitung ulang, ⭐ **tabelnya yang salah**, bukan operannya.

⭐ Kolom asal membedakan baris **hasil migrasi** (beku) dari baris **hasil hitung** (boleh dibangun
ulang). ⛔ **Tanpa penjaga ini, satu perintah bangun ulang menimpa seluruh angka historis dan tidak
dapat dikembalikan.**

Kunci penyaring disertakan supaya pembaca SQL tidak perlu join balik.

⭐ **Dua proyeksi sengaja TIDAK dibuat:** untuk induk lapisan *(hanya salinan penanda yang sudah ada
di anaknya)* dan untuk rincian angsuran *(selisihnya hanya satu tingkat)*.

Ditambah tiga pernyataan kolom khas: kolom endorsemen kosong pada polis baru, kolom polis baru
kosong pada endorsemen, dan keadaan layar **tidak disimpan**.

#### Batas — yang TIDAK termasuk

⛔ Rumus selisihnya sendiri — tiket **25**.
⛔ Penanda migrasi pada baris proyeksi — tiket **28**.

#### Cara mengujinya

Lewat seam `repository`. ⭐ **Uji yang paling penting:** bangun ulang proyeksi dijalankan, lalu
baris hasil migrasi diperiksa **satu per satu** — satu baris yang berubah berarti gagal.

#### Acceptance criteria

- [ ] **AC 23** — rumus selisih **tidak pernah** ditulis di dalam basis data
- [ ] **AC 24** — tabel selisih **hanya ditulis** lapisan aplikasi
- [ ] **AC 25** — tabel selisih ditulis **dalam transaksi yang sama** dengan generasinya
- [ ] **AC 26** — bangun ulang **hanya menyentuh baris hasil hitung**
- [ ] **AC 27** — bila isinya beda dari hitung ulang, **tabelnya** yang dinyatakan salah
- [ ] **AC 28** — tabel selisih memuat kunci penyaring, sehingga pembaca tidak perlu join balik
- [ ] **AC 29** — **nol** proyeksi untuk induk lapisan
- [ ] **AC 30** — **nol** proyeksi untuk rincian angsuran
- [ ] **AC 31** — lapisan penyimpanan mengambil **dua baris**, bukan menjalankan agregasi
- [ ] **AC 32** — kolom khas endorsemen **kosong** pada baris polis baru
- [ ] **AC 33** — kolom khas polis baru **kosong** pada endorsemen, dan **tidak** dihapus dari skema
- [ ] **AC 34** — keadaan layar **tidak tersimpan** di tabel mana pun

#### ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum dijelaskan apakah anak tabel proyeksi — untuk penyebaran dan untuk angsuran —
memang dibutuhkan pembaca SQL, atau cukup induknya.** Membuat keduanya berarti memelihara tabel yang
mungkin tidak pernah dibaca; tidak membuatnya berarti pembaca harus join balik ke tabel sumber.

⭐ **Induknya dapat dikerjakan lebih dulu** — AC 23–28 dan 31–34 tidak bergantung pada keputusan itu.
⛔ AC 29 dan 30 menyatakan proyeksi mana yang **tidak** dibuat, dan itu sudah pasti.

## NB Treaty In - 27 - Perbedaan endorsemen pada tabel yang sama [DIGANTIKAN]

> ## ⛔⛔ DIGANTIKAN — 23 September 2026
>
> **Tiket ini tidak dikerjakan.** Ia digantikan oleh
> **`.scratch\edm-treaty-in\issues\08`**.
>
> ### Kenapa
>
> Brief ronde tiket sebelumnya berlingkup **dua spec sekaligus**, sehingga pekerjaan **endorsemen**
> jatuh ke folder tiket **polis baru**. Ronde berikutnya membaginya ulang ke folder yang benar,
> lebih halus, dan — yang menentukan — **setiap tiket menyebut tiket NB mana yang membuat
> tabelnya**. Tanpa gate itu, tiket endorsemen bisa dikerjakan sebelum tabelnya ada.
>
> ⚠️ **Kekeliruan lingkup ini milik penyusun brief, bukan pelaksana ronde ini.**
>
> ⭐ Berkas ini **tidak dihapus** supaya jejaknya tidak hilang — tanpa ini, nomor tiket NB akan
> tampak melompat dari 23 ke 29 tanpa sebab.
>
> ### Yang berlaku
>
> | | |
> | --- | --- |
> | Tiket penyimpanan **NB** | `nb-treaty-in\issues\` **16–23** |
> | Tiket penyimpanan **EDM** | `edm-treaty-in\issues\` **01–11** |
>
> ⛔ Isi di bawah dibiarkan apa adanya sebagai catatan, **bukan sebagai pekerjaan.**

---


**Status:** ~~blocked~~ — **DIGANTIKAN**
**Blocked by:** **19** · **24** · ⛔ `[work owner]` **beda dagang tabel pecahan penyebaran** · ⛔ `[data DBA]` **presisi fisik**
**Menutup:** EDM AC **35–38** · EDM AC **54–58** *(9 AC)*
**Sumber:** `edm-treaty-in\spec-penyimpanan-relasional.md` ID-5..ID-7 · ID-39..ID-42 · ID-46

#### Hasil & nilai pengguna

Tabel dasarnya sama dengan polis baru; yang berbeda **isinya**. Tiket ini menyatakan perbedaan itu
dengan tepat, dan membuktikan bahwa bentuk tabel dasar **tidak berubah**.

⭐⭐ **Nol tabel dasar dirancang ulang.** Bila pekerjaan ini menuntutnya, itu **temuan** yang wajib
dilaporkan, bukan perubahan yang dijalankan diam-diam.

#### Yang dibangun

| Perbedaan | Endorsemen | Polis baru |
| --- | --- | --- |
| presisi pembagian penyebaran | **20** | **10** |
| baris pertama penyebaran | nilai bawaan **100** bila persentasenya kosong | — |
| rincian angsuran bertingkat | **hidup** pada bentuk non-proporsional | tidak ada |
| asal rincian angsuran | ⭐ dapat berasal dari **salinan master kontrak** | diketik |

⚠️ **Perbedaan presisi ditiru apa adanya**, dan **wajib punya test tersendiri** — ia perilaku lama
yang disalin sadar, bukan cacat.

Ditambah pernyataan bentuk: cacah tabel pada kedua bentuk endorsemen, panjang minimum kolom
keterangan, dan skala kolom uang.

#### Batas — yang TIDAK termasuk

⛔ Rantai generasi — tiket **24**.
⛔ Proyeksi selisih — tiket **26**.
⛔ `CREATE TABLE` — presisi fisik dicocokkan DBA **di dalam** tiket ini.

#### Cara mengujinya

Lewat seam `repository`. ⭐ Uji presisi dijalankan **berdampingan** — kasus yang sama disimpan
sebagai polis baru dan sebagai endorsemen, dan hasilnya **harus berbeda** sesuai ketetapan.

#### Acceptance criteria

- [ ] **AC 35** — penyebaran endorsemen memakai presisi **20**; polis baru **10**
- [ ] **AC 36** — baris pertama penyebaran menerima nilai bawaan **100** ketika persentasenya kosong
- [ ] **AC 37** — rincian angsuran bertingkat **tersimpan** pada endorsemen non-proporsional
- [ ] **AC 38** — rincian dari salinan master kontrak tersimpan **sama** seperti yang diketik
- [ ] **AC 54** — bentuk tabel dasar **tidak berubah**
- [ ] **AC 55** — cacah tabel pada kedua bentuk endorsemen sesuai ketetapan
- [ ] **AC 56** — tabel pecahan penyebaran menyimpan empat medan yang ditetapkan
- [ ] **AC 57** — kolom keterangan berpanjang sekurangnya batas yang ditetapkan
- [ ] **AC 58** — nilai berdesimal lebih dari skala kolom **dibulatkan, bukan ditolak**

#### ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Beda dagang tabel pecahan penyebaran dari tabel penyebaran biasa belum
dijelaskan.** AC 56 menyatakan isinya, tetapi **kapan baris masuk ke yang satu dan bukan ke yang
lain** belum dapat diuji.

`[data DBA]` **Presisi fisik kolom uang** — sama dengan penahan tiket **18**.

⚠️ **Dan satu pertentangan di dalam spec, dicatat karena spec tidak boleh disunting:**
`[terverifikasi]` AC **49** pada spec yang sama masih menuntut **sembilan** desimal, sementara
AC **58** menetapkan **delapan**. ⛔ **Tidak diputuskan di tiket ini.**

## NB Treaty In - 28 - Penanda migrasi endorsemen dan kepatuhan lapisan [DIGANTIKAN]

> ## ⛔⛔ DIGANTIKAN — 23 September 2026
>
> **Tiket ini tidak dikerjakan.** Ia digantikan oleh
> **`.scratch\edm-treaty-in\issues\09 · 10 · 11`**.
>
> ### Kenapa
>
> Brief ronde tiket sebelumnya berlingkup **dua spec sekaligus**, sehingga pekerjaan **endorsemen**
> jatuh ke folder tiket **polis baru**. Ronde berikutnya membaginya ulang ke folder yang benar,
> lebih halus, dan — yang menentukan — **setiap tiket menyebut tiket NB mana yang membuat
> tabelnya**. Tanpa gate itu, tiket endorsemen bisa dikerjakan sebelum tabelnya ada.
>
> ⚠️ **Kekeliruan lingkup ini milik penyusun brief, bukan pelaksana ronde ini.**
>
> ⭐ Berkas ini **tidak dihapus** supaya jejaknya tidak hilang — tanpa ini, nomor tiket NB akan
> tampak melompat dari 23 ke 29 tanpa sebab.
>
> ### Yang berlaku
>
> | | |
> | --- | --- |
> | Tiket penyimpanan **NB** | `nb-treaty-in\issues\` **16–23** |
> | Tiket penyimpanan **EDM** | `edm-treaty-in\issues\` **01–11** |
>
> ⛔ Isi di bawah dibiarkan apa adanya sebagai catatan, **bukan sebagai pekerjaan.**

---


**Status:** ~~blocked~~ — **DIGANTIKAN**
**Blocked by:** **22** · **26** · ⛔ `[work owner]` **lingkup pemindahan dokumen lama** · ⛔ `[work owner]` **pembatalan dan batas endorse**
**Menutup:** EDM AC **39–53** *(15 AC)*
**Sumber:** `edm-treaty-in\spec-penyimpanan-relasional.md` ID-1..ID-3 · ID-33..ID-38 · ID-43..ID-47

#### Hasil & nilai pengguna

Baris hasil migrasi dibedakan dari baris hasil hitung, dan dua keadaan yang perlu diketahui
ditandai **tanpa mengubah satu angka pun**.

⭐ **Nilai lama tidak pernah dihitung ulang** — termasuk bila rumus lamanya diketahui keliru.
Yang dilakukan hanya **menandainya**, supaya kelak dapat diperlakukan berbeda bila diputuskan
demikian.

#### Yang dibangun

Dua penanda, keduanya **hanya berarti untuk baris hasil migrasi**:

| Penanda | Artinya |
| --- | --- |
| **pasangan bergeser** | kunci dagang pasangan menurut nomor urut ternyata berbeda ⇒ ⚠️ **anomali sungguhan**, sebab di endorsemen baris tidak dapat dihapus |
| **rumus berlapis** | nilainya lahir dari rumus lama yang mengurangi terhadap selisih ⇒ dapat ditentukan pasti dari generasi sebelumnya |

Ditambah sembilan pernyataan kepatuhan yang berlaku sama seperti pada polis baru: satu transaksi ·
skema eksplisit · kolom pelaku dari identitas login · nol tipe mengambang · skala kolom uang · kode
dan penanda tetap teks · teks kosong menjadi tak-bernilai · arah ketergantungan · satu antarmuka
penyimpanan.

#### Batas — yang TIDAK termasuk

⛔ Pemuat massalnya sendiri — tiket **22**.
⛔ Tabel proyeksi tempat penanda ini tinggal — tiket **26**.

#### Cara mengujinya

Lewat seam `repository`. ⭐ **Uji utama:** dokumen lama dimuat, penandaan dijalankan, lalu seluruh
nilai uang dibandingkan dengan sebelum penandaan — **satu angka yang berubah berarti gagal**.

#### Acceptance criteria

- [ ] **AC 39** — nilai hasil migrasi **tidak dihitung ulang**, termasuk digit galatnya
- [ ] **AC 40** — pemasangan antar generasi memakai **nomor urut**
- [ ] **AC 41** — penanda pasangan bergeser dihitung **tanpa mengubah satu angka pun**
- [ ] **AC 42** — penanda rumus berlapis terisi pada setiap baris yang memenuhi syaratnya
- [ ] **AC 43** — kedua penanda **hanya** terisi pada baris hasil migrasi
- [ ] **AC 44** — pemuat migrasi menulis lewat antarmuka penyimpanan yang **sama**
- [ ] **AC 45** — seluruh penyimpanan satu generasi berada dalam **satu transaksi**
- [ ] **AC 46** — setiap query menulis skema **eksplisit**
- [ ] **AC 47** — kolom pelaku diisi dari identitas login
- [ ] **AC 48** — nol kolom uang bertipe mengambang
- [ ] **AC 49** — kolom uang menerima skala desimal yang ditetapkan
- [ ] **AC 50** — kode dan penanda tersimpan sebagai **teks**, nol di depan utuh
- [ ] **AC 51** — teks kosong pada kolom angka atau tanggal tersimpan **tak-bernilai**
- [ ] **AC 52** — arah ketergantungan **tidak pernah dibalik**
- [ ] **AC 53** — hanya ada **satu** antarmuka penyimpanan polis treaty

#### ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Lingkup pemindahan dokumen lama belum diputuskan** — sama dengan penahan tiket
**22**. Penandaan hanya berarti atas baris yang benar-benar dipindahkan.

`[work owner]` **Dua butir endorsemen belum dijawab:** apakah polis yang sudah dibatalkan masih
boleh di-endorse lagi, dan berapa kali satu polis boleh di-endorse *(batas teknis 99)*. Keduanya
menentukan baris mana yang sah ada di dalam rantai, dan karena itu baris mana yang ditandai.

⚠️ **AC 49 membawa pertentangan spec** yang dicatat di tiket **18** dan **27** — skala delapan
lawan sembilan. ⛔ Tidak diputuskan di sini.

# EDM Treaty In

Jumlah tiket: **12**

## EDM Treaty In - 01 - Kolom khas endorsemen pada tabel yang sudah ada

**Status:** ready-for-agent
**Blocked by:** — *(dapat mulai sesudah tiket NB)*
**Bergantung pada tiket NB:** **16** *(kerangka penyimpanan)* · **19** *(pemecah dokumen)*
**Menutup:** AC **32–34** · AC **54–55** *(5 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-5 · ID-6 · ID-20 · ID-21 · ID-22

#### Hasil & nilai pengguna

Endorsemen menempati **tabel yang sama** dengan polis baru. Yang membedakannya hanya **kolom mana
yang terisi** — bukan bentuk tabelnya.

⭐⭐ **Ini tiket pertama rantai endorsemen, dan sekaligus penjaganya:** ia membuktikan bahwa bentuk
tabel dasar **tidak berubah** oleh seluruh pekerjaan endorsemen.

#### Yang dibangun

Pengisian dan pengosongan kolom menurut jenis berkas:

| Golongan | Perlakuan |
| --- | --- |
| **21 kolom khas endorsemen** | terisi pada endorsemen, ⛔ **kosong** pada polis baru |
| **4 kolom khas polis baru** | kosong pada endorsemen, ⛔ **tidak dihapus** dari skema — bukan kolom mati |
| **keadaan layar** | ⛔ **tidak disimpan** di tabel mana pun |

Ditambah dua pernyataan bentuk yang diuji langsung: bentuk kesembilan tabel dasar **tidak berubah**,
dan cacah tabel pada kedua bentuk endorsemen sesuai ketetapan.

⛔ **Nol tabel dasar dibuat tiket ini.** Sepuluh tabel dasar dibuat tiket NB **16** dan **19**.

#### Batas — yang TIDAK termasuk

⛔ Pembuatan tabel dasar — tiket NB **16** · **19**.
⛔ Tabel proyeksi selisih — tiket **07**.
⛔ Tipe dan presisi kolom — tiket **11**.

#### Cara mengujinya

Lewat seam `repository`. ⭐ **Uji berdampingan:** kasus yang sama disimpan sebagai polis baru dan
sebagai endorsemen — kolom yang terisi harus **berbeda**, dan kolom yang kosong **tetap ada**.

⚠️ Uji bentuk tabel dijalankan sebagai pemeriksaan skema, bukan lewat data.

#### Acceptance criteria

- [ ] **AC 32** — dua puluh satu kolom khas endorsemen **kosong** pada baris polis baru
- [ ] **AC 33** — empat kolom khas polis baru **kosong** pada endorsemen, dan **tidak** dihapus dari skema
- [ ] **AC 34** — keadaan layar **tidak tersimpan** di tabel mana pun
- [ ] **AC 54** — bentuk kesembilan tabel dasar **tidak berubah**
- [ ] **AC 55** — cacah tabel pada kedua bentuk endorsemen sesuai ketetapan

## EDM Treaty In - 02 - Rantai generasi dan larangan percabangan

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Satu polis boleh di-endorse **tanpa batas**. Batas teknis 99 tidak dipertahankan; kolom nomor generasi dilebarkan. Diuji ke korpus: kenaikannya aritmetika *(`Prodke+1`)*, nol topeng lebar tetap.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 1.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** ⛔ `[work owner]` **batas berapa kali satu polis boleh di-endorse** — batas teknis 99, batas dagang belum ditetapkan~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **16** *(kerangka penyimpanan dan kunci generasi)*
**Menutup:** AC **1–6** *(6 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-9 · ID-10 · ID-11 · ID-12 · ID-13

#### Hasil & nilai pengguna

Endorsemen tersimpan sebagai **generasi berikutnya** dari polis yang sama, dan rantainya **lurus**:
satu generasi, satu penerus.

⛔⛔ **Ini memperbaiki bahaya nyata, bukan kehati-hatian teoretis.** `[terverifikasi]` Sistem lama
membaca nomor generasi terakhir **tanpa penguncian**, lalu menambah satu. Dua endorsemen serentak
atas polis yang sama membaca angka yang sama, dan **keduanya tersimpan**.
⭐ Di sistem baru yang kedua **gagal terang-terangan**.

#### Yang dibangun

Baris generasi bernomor ≥ 1 dengan penunjuk ke generasi sebelumnya, beserta tiga aturan:

| Aturan | Yang dicegahnya |
| --- | --- |
| penunjuk generasi **unik** | rantai menjadi pohon — *"selisih mana yang benar"* tidak terjawab |
| kunci alami **unik** | dua endorsemen serentak mendapat nomor generasi yang sama |
| generasi lampau **tidak dapat disunting** | dasar selisih berubah sesudah disetujui |

Ditambah pembentukan nomor endorsemen: nomor polis induk, pemisah, lalu **dua digit** dengan nol di
depan bila kurang dari sepuluh.

#### Batas — yang TIDAK termasuk

⛔ Aturan keutuhan nomor urut antar generasi — tiket **03**.
⛔ Perhitungan selisih — tiket **06**.
⛔ Jalur penomoran kedua yang ada di sistem lama — **tidak dimigrasi**.

#### Cara mengujinya

Lewat seam `repository`. ⭐ **Tiga uji yang wajib berdiri sendiri:**

1. **rantai tiga generasi** — ⚠️ rantai dua generasi **tidak memadai**, sebab selisih terhadap
   generasi tepat sebelumnya baru terbukti mulai generasi ketiga;
2. **penolakan penerus kedua** atas generasi yang sama;
3. **dua endorsemen serentak** — yang kedua harus ditolak **oleh basis data**, bukan oleh
   pemeriksaan di aplikasi.

#### Acceptance criteria

- [ ] **AC 1** — baris endorsemen tersimpan dengan nomor generasi ≥ 1
- [ ] **AC 2** — penunjuk generasi **terisi** dan menunjuk generasi sebelumnya
- [ ] **AC 3** — dua baris tidak boleh berbagi penunjuk generasi yang sama
- [ ] **AC 4** — dua endorsemen serentak: yang kedua **ditolak**
- [ ] **AC 5** — nomor endorsemen berbentuk nomor polis + pemisah + **dua digit**
- [ ] **AC 6** — menyunting generasi yang sudah punya penerus **ditolak**

#### ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum ditetapkan berapa kali satu polis boleh di-endorse.** Batas teknisnya **99**
— nomor generasi dibentuk dua digit. Bila ada batas dagang yang lebih rendah, itu **validasi yang
perlu dibangun** di tiket ini; bila tidak ada, angka 99 tetap batas yang wajib diketahui sebelum
sistem baru memakainya.

⭐ **Lima dari enam AC tidak bergantung pada jawabannya** — hanya batas cacah generasi yang
menunggu.

## EDM Treaty In - 03 - Keutuhan nomor urut antar generasi

**Status:** ready-for-agent
**Blocked by:** **02**
**Bergantung pada tiket NB:** **17** *(nomor urut baris anak)*
**Menutup:** AC **7–9** *(3 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-8 · ID-14 · ID-15

#### Hasil & nilai pengguna

Generasi baru **wajib memuat setiap baris rincian** yang ada di generasi sebelumnya. Yang hilang
berarti endorsemen itu diam-diam menghapus komitmen yang masih berlaku — dan itu **ditolak**.

⭐ **Nilai lama tidak menjadi tabel.** Ia baris yang ditunjuk penunjuk generasi — dibaca, bukan
disalin ke tempat kedua.

#### Yang dibangun

Aturan keutuhan di lapisan layanan: sebelum generasi baru disimpan, setiap nomor urut milik
generasi sebelumnya diperiksa ada. ⛔ Yang hilang **membatalkan seluruh penyimpanan**, bukan
menyimpan sebagian.

Ditambah penegasan bahwa selisih dihitung terhadap **generasi tepat sebelumnya**, bukan terhadap
polis asli — dan bahwa tidak ada tabel salinan nilai lama.

⚠️ **Susunan bersarang yang dibuat sistem lama — nilai lama di dalam nilai lama — tidak ditiru.**
`[terverifikasi]` Korpus **tidak pernah membacanya**.

#### Batas — yang TIDAK termasuk

⛔ Perilaku nomor urut saat baris baru ditambahkan — tiket **04**.
⛔ Penanda pasangan bergeser — tiket **09**.

#### Cara mengujinya

Lewat seam `repository`. ⭐ **Uji utama:** generasi kedua dibuat dengan satu baris rincian
dihilangkan — penyimpanan harus **ditolak seluruhnya**, dan tidak boleh ada satu baris pun
tertinggal.

#### Acceptance criteria

- [ ] **AC 7** — selisih dihitung terhadap generasi **tepat sebelumnya**
- [ ] **AC 8** — generasi yang kehilangan salah satu nomor urut **ditolak**
- [ ] **AC 9** — **nol** tabel salinan nilai lama; ia dibaca lewat penunjuk generasi

## EDM Treaty In - 04 - Nomor urut di endorsemen — tanpa penghapusan

**Status:** ready-for-agent
**Blocked by:** **03**
**Bergantung pada tiket NB:** **17** *(nomor urut baris anak)*
**Menutup:** AC **10–13** *(4 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-16 · ID-17 · ID-18 · ID-34

#### Hasil & nilai pengguna

Di endorsemen, baris rincian **tidak dapat dihapus**, dan nomor urut lama **terbawa apa adanya**.

⚠️ **Ini berbeda dari polis baru**, tempat baris boleh dihapus dan nomornya dinomori ulang rapat —
aman di sana karena belum ada generasi untuk dipasangkan. ⭐ Begitu generasi pertama ditutup, nomor
urut **beku selamanya**.

#### Yang dibangun

Aturan penomoran khas endorsemen:

| Aturan | Sebab |
| --- | --- |
| baris rincian **tidak dapat dihapus** | pemasangan antar generasi lewat nomor urut menuntutnya |
| nomor urut lama **terbawa apa adanya** | pasangan lama tidak boleh bergeser |
| baris baru mendapat **maksimum + 1** | penyisipan di tengah menggeser seluruh pasangan sesudahnya |

⭐ **Akibat yang menguntungkan:** karena tidak ada penghapusan, pasangan yang bergeser berarti
**anomali sungguhan** — bukan derau yang wajar. Itu membuat penandanya berguna.

#### Batas — yang TIDAK termasuk

⛔ Penandaan pasangan bergeser saat migrasi — tiket **09**.
⛔ Pembatalan — tiket **05**; ia **bukan** penghapusan baris.

#### Cara mengujinya

Lewat seam `repository`. ⭐ Uji utama: permintaan menghapus baris rincian pada endorsemen harus
**ditolak**; dan baris baru yang ditambahkan harus mendapat nomor sesudah yang terbesar, bukan
mengisi lubang.

#### Acceptance criteria

- [ ] **AC 10** — pada endorsemen, baris rincian **tidak dapat dihapus**
- [ ] **AC 11** — nomor urut lama **terbawa apa adanya** ke generasi berikutnya
- [ ] **AC 12** — baris baru mendapat nomor urut **maksimum + 1**
- [ ] **AC 13** — pasangan yang bergeser diperlakukan sebagai **anomali**, bukan keadaan wajar

## EDM Treaty In - 05 - Pembatalan sebagai generasi bernilai nol

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Polis yang sudah dibatalkan **masih boleh di-endorse lagi**. Nol gerbang peran — layar cukup memberi **peringatan**, bukan penghalang.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 4.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** ⛔ `[work owner]` **sesudah dibatalkan, polis masih boleh di-endorse lagi atau tidak**~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **19** *(pemecah dokumen — kolom jenis berkas)*
**Menutup:** AC **14–15** *(2 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-19 · ID-41

#### Hasil & nilai pengguna

Pembatalan polis adalah **jenis endorsemen**, dipilih di awal saat berkas dibuat — bukan alur
tersendiri dan bukan tombol terpisah.

⛔⛔ **Pembatalan BUKAN penghapusan.** Ia menghasilkan **generasi baru berisi nol**, sehingga
jejaknya tetap ada dan dapat dipertanggungjawabkan.
⛔ **Tiket yang membuat jalur hapus untuk pembatalan salah.**

#### Yang dibangun

Jenis endorsemen tersimpan di kolomnya sendiri, dan salah satu nilainya berarti pembatalan.
Memilihnya menghasilkan generasi baru yang seluruh kolom uangnya **bernilai nol**, lewat jalur
penyimpanan yang **sama** dengan endorsemen lain.

⛔ **Nol jalur hapus dibangun untuk pembatalan.**

#### Batas — yang TIDAK termasuk

⛔ Aturan nomor urut — tiket **04**; pembatalan **tidak** menghapus baris, jadi aturannya berlaku apa adanya.
⛔ Perhitungan selisih terhadap generasi yang dibatalkan — tiket **06**.

#### Cara mengujinya

Lewat seam `repository`. ⭐ **Uji utama:** polis dibatalkan, lalu **cacah baris rincian dihitung**
— jumlahnya harus **sama** dengan generasi sebelumnya, dan seluruh kolom uangnya nol.
⛔ Baris yang berkurang berarti gagal.

#### Acceptance criteria

- [ ] **AC 14** — pembatalan menghasilkan generasi baru dengan kolom uang **bernilai nol**, bukan penghapusan
- [ ] **AC 15** — jenis endorsemen tersimpan di kolomnya sendiri dan dipilih di awal

#### ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum dijawab apakah polis yang sudah dibatalkan masih boleh di-endorse lagi.**
Jawabannya menentukan apakah generasi pembatalan menjadi **ujung rantai** — sehingga penunjuk
generasi berikutnya harus ditolak — atau sekadar generasi biasa yang nilainya nol.

⭐ **Kedua AC dapat dibangun tanpa menunggu**; yang menunggu hanya **penjaga rantai sesudahnya**.

## EDM Treaty In - 06 - Perhitungan selisih di lapisan layanan

**Status:** ready-for-agent
**Blocked by:** **02** · **03**
**Bergantung pada tiket NB:** **16** *(kerangka penyimpanan)*
**Menutup:** AC **16–22** *(7 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-4 · ID-28 · ID-29 · ID-30 · ID-31 · ID-32

#### Hasil & nilai pengguna

Selisih antar generasi dihitung **di lapisan layanan**, dengan **satu rumus** untuk semua kasus —
endorsemen pertama maupun berlapis, proporsional maupun non-proporsional:

```
selisih.X = baris_ini.X − baris(generasi sebelumnya).X
```

⛔⛔ **Rumus ini TIDAK pernah ditulis di dalam basis data.** View yang menghitung selisih
**dibatalkan** — rumus di dua tempat cepat atau lambat bercabang.
⛔ **Tiket ini tidak boleh menghidupkannya kembali.**

#### Yang dibangun

Perhitungan selisih dengan empat aturan turunan:

| Golongan medan | Perlakuan |
| --- | --- |
| **uang** | dikurangi |
| **persentase** | ⭐ **disalin**, tidak dikurangi |
| **kunci** | disalin apa adanya |
| **arah utang-piutang** | diturunkan dari **tanda** selisih |

⭐ Pada selisih lapisan, arah utang-piutang diturunkan dari **jumlah sepanjang daftar lapisan**,
bukan dari satu lapisan saja.

⭐ Lapisan penyimpanan hanya mengambil **dua baris** — baris ini dan baris yang ditunjuknya.
⛔ Nol agregasi di basis data.

⚠️ `[penyimpangan sadar]` **Sistem lama punya DUA rumus di sisi proporsional**, dan yang kedua
mengurangi terhadap **selisih** generasi lampau, bukan terhadap **nilainya**. Dengan angka contoh:
nilai baru 180, nilai lama 150, selisih lama 50 ⇒ varian kedua menghasilkan **130**, padahal yang
benar **30**. ⛔ **Varian kedua tidak ditiru** — `[keputusan work owner]`.

#### Batas — yang TIDAK termasuk

⛔ Penulisan hasilnya ke tabel proyeksi — tiket **07**.
⛔ Penandaan baris yang lahir dari rumus lama — tiket **09**.

#### Cara mengujinya

Operannya lewat seam `repository`; rumusnya lewat lapisan layanan.

⭐ **Uji regresi wajib:** sisi non-proporsional harus menghasilkan angka yang **sama** sebelum dan
sesudah penyeragaman — `[terverifikasi]` ia memang sudah seragam, dan penyeragaman tidak boleh
mengubahnya. Uji yang menemukan perubahan berarti rumusnya salah diterjemahkan.

#### Acceptance criteria

- [ ] **AC 16** — satu rumus untuk **seluruh** generasi
- [ ] **AC 17** — medan **uang** berisi hasil pengurangan
- [ ] **AC 18** — medan **persentase** disalin, tidak dikurangi
- [ ] **AC 19** — medan **kunci** disalin apa adanya
- [ ] **AC 20** — arah utang-piutang diturunkan dari **tanda** selisih
- [ ] **AC 21** — pada selisih lapisan, arah itu diturunkan dari **jumlah** sepanjang daftar
- [ ] **AC 22** — sisi non-proporsional menghasilkan angka **sama** sebelum dan sesudah

#### ⚠️ Satu hal yang tidak terbaca dari sistem lama

`[terbuka]` Kedua blok rumus di sistem lama **sama-sama bergerbang aktif**, dan pemilihnya hanya
tertulis di **keterangan langkah**, bukan di gerbang yang dijalankan. ⇒ Mana yang benar-benar
berjalan **tidak dapat dinyatakan dari ekspor**.

⭐ **Tidak menahan tiket ini** — varian kedua tidak ditiru apa pun jawabannya.

## EDM Treaty In - 07 - Proyeksi selisih yang dapat dibaca langsung

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Kedua anak tabel proyeksi **dibuat** — pembaca SQL langsung membutuhkan selisih sampai rincian per baris. Sesuai artefak rancangan Excel, sheet EDM Treaty In Prop.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 2.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** **06** · ⛔ `[work owner]` **anak proyeksi dibutuhkan pembaca SQL atau tidak**~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **20** *(transaksi tunggal dan skema eksplisit)*
**Menutup:** AC **23–31** *(9 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-6 · ID-7 · ID-23..ID-27 · ID-38 · ID-42

#### Hasil & nilai pengguna

Angka selisih disimpan sebagai **proyeksi** yang dapat dibaca langsung dengan SQL, di luar layar
aplikasi — `[keputusan work owner]` pembaca seperti itu memang ada.

⛔⛔ **Tabel proyeksi BUKAN tabel sumber**, dan ketiga aturan yang membuatnya begitu wajib dibangun
utuh — bukan diringkas.

#### Yang dibangun

##### Tiga aturan yang membuatnya bukan tabel sumber

**①** Ia **hanya ditulis lapisan aplikasi**, dalam **transaksi yang sama** dengan generasinya.

**②** Ia **boleh dihapus total dan dibangun ulang** — ⛔⛔ **tetapi hanya baris yang dihitung sistem
baru**. Baris hasil migrasi **beku**.
⛔ **Tanpa aturan ini, satu perintah bangun ulang menimpa seluruh angka historis dan tidak dapat
dikembalikan.**

**③** Bila isinya berbeda dari hasil hitung ulang, ⭐ **tabelnya yang salah**, bukan operannya.

##### Sisanya

Kolom asal membedakan baris **hasil migrasi** dari baris **hasil hitung**. Kunci penyaring
disertakan supaya pembaca SQL tidak perlu join balik.

⭐ **Dua proyeksi sengaja TIDAK dibuat:** untuk induk lapisan *(hanya salinan penanda yang sudah ada
di anaknya)* dan untuk rincian angsuran *(selisihnya hanya satu tingkat)*.

#### Batas — yang TIDAK termasuk

⛔ Rumus selisihnya sendiri — tiket **06**.
⛔ Penanda migrasi pada baris proyeksi — tiket **09**.
⛔ **View yang menghitung selisih** — dibatalkan; jangan dihidupkan.

#### Cara mengujinya

Lewat seam `repository`. ⭐ **Uji yang paling penting:** bangun ulang proyeksi dijalankan, lalu
**setiap baris hasil migrasi diperiksa satu per satu** — satu baris yang berubah berarti gagal.

⚠️ Uji itu wajib berdiri sendiri, tidak digabung dengan uji bangun ulang yang lain.

#### Acceptance criteria

- [ ] **AC 23** — rumus selisih **tidak pernah** ditulis di dalam basis data
- [ ] **AC 24** — tabel selisih **hanya ditulis** lapisan aplikasi
- [ ] **AC 25** — tabel selisih ditulis **dalam transaksi yang sama** dengan generasinya
- [ ] **AC 26** — bangun ulang **hanya menyentuh baris hasil hitung**
- [ ] **AC 27** — bila isinya beda dari hitung ulang, **tabelnya** yang dinyatakan salah
- [ ] **AC 28** — tabel selisih memuat kunci penyaring, sehingga pembaca tidak perlu join balik
- [ ] **AC 29** — **nol** proyeksi untuk induk lapisan
- [ ] **AC 30** — **nol** proyeksi untuk rincian angsuran
- [ ] **AC 31** — lapisan penyimpanan mengambil **dua baris**, bukan menjalankan agregasi

#### ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum dijelaskan apakah anak tabel proyeksi — untuk penyebaran dan untuk angsuran —
memang dibutuhkan pembaca SQL, atau cukup induknya.** Membuat keduanya berarti memelihara tabel yang
mungkin tidak pernah dibaca; tidak membuatnya berarti pembaca harus join balik.

⭐ **Induknya dapat dikerjakan lebih dulu** — AC 23–28 dan 31 tidak bergantung pada keputusan itu,
dan AC 29–30 justru menyatakan proyeksi mana yang **tidak** dibuat, yang sudah pasti.

## EDM Treaty In - 08 - Perbedaan perhitungan sebaran dan rincian angsuran

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` ⛔ **Tabel sebaran tambahan DIBATALKAN** — keempat medannya turunan, dua dari master `POOLDATA.PROPORTIONALARRG`, dua dihitung dari total yang sudah tersimpan. ⚠️ **Lingkup tiket ini menyusut**: bagian sebaran tambahan keluar, bagian rincian angsuran tetap berlaku.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 3.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** ⛔ `[work owner]` **beda dagang tabel sebaran tambahan** dari sebaran risiko belum dijelaskan~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **19** *(pemecah dokumen menjadi baris)*
**Menutup:** AC **35–38** · AC **56** *(5 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-39 · ID-40

#### Hasil & nilai pengguna

Tabelnya sama dengan polis baru; **perhitungannya berbeda**. Tiket ini menyatakan perbedaan itu
dengan tepat, dan membuatnya **diuji terpisah** — bukan diseragamkan.

⚠️ **Perbedaan ini ditiru apa adanya.** Ia perilaku lama yang disalin sadar, bukan cacat yang
diperbaiki.

#### Yang dibangun

| Perbedaan | Endorsemen | Polis baru |
| --- | --- | --- |
| presisi pembagian sebaran | **20** | **10** |
| baris pertama sebaran | menerima nilai bawaan **100** bila persentasenya kosong | — |
| rincian angsuran bertingkat | **hidup** pada bentuk non-proporsional | tidak ada |
| asal rincian angsuran | ⭐ dapat berasal dari **salinan master kontrak** | diketik pengguna |

⭐ **Sebab endorsemen memerlukan nilai bawaan dan polis baru tidak:** `[terverifikasi]` aturan yang
**mengisi** medan pembagi itu ada di modul polis baru dan **tidak ada** di modul endorsemen —
sehingga di endorsemen medannya dapat tiba kosong, sementara rumusnya **membaginya**.

~~Ditambah satu tabel sebaran tambahan dengan empat medannya.~~ ⛔ **DICABUT 23-09-2026 sore** — sebaran tambahan **tidak dimigrasi sama sekali**, tabel maupun perhitungannya.

#### Batas — yang TIDAK termasuk

⛔ Pemecahan dokumen menjadi baris sebaran — tiket NB **19**.
⛔ Selisih sebaran antar generasi — tiket **07**.

#### Cara mengujinya

Lewat seam `repository`. ⭐ **Uji berdampingan wajib:** kasus yang sama disimpan sebagai polis baru
dan sebagai endorsemen — hasilnya **harus berbeda** sesuai ketetapan presisi.
⛔ Uji yang menemukan keduanya sama berarti perbedaannya hilang.

#### Acceptance criteria

- [ ] **AC 35** — sebaran endorsemen memakai presisi **20**; polis baru **10**
- [ ] **AC 36** — baris pertama sebaran menerima nilai bawaan **100** ketika persentasenya kosong
- [ ] **AC 37** — rincian angsuran bertingkat **tersimpan** pada endorsemen non-proporsional
- [ ] **AC 38** — rincian dari salinan master kontrak tersimpan **sama** seperti yang diketik
- [x] ~~**AC 56** — tabel sebaran tambahan menyimpan empat medan~~ ⛔ **AC DITARIK DARI LINGKUP.** Cacah AC spec EDM **58 → 57**.

#### ⛔ Kenapa tiket ini `blocked`

~~`[work owner]` **Beda dagang tabel sebaran tambahan dari tabel sebaran risiko belum dijelaskan.**~~ ⛔ **BUTIR GUGUR 23-09-2026 sore.** `[keputusan work owner]` *"Itu tidak ada"* lalu *"itu tidak usah di migrasi perhitungannya"*.

⭐ Terbukti dari korpus: keempat medannya **turunan**. `TreatyType` dan `SharePercentage` diambil dari master `POOLDATA.PROPORTIONALARRG`; `PremiumSpreaded` dan `ClaimSpreaded` dihitung `TotalPremium × SharePercentage/100` dan `TotalClaim × SharePercentage/100`.

⇒ **Lingkup tiket ini menyusut.** Yang tetap berlaku: AC 35 · 36 *(penyebaran risiko)* dan AC 37 · 38 *(rincian angsuran)* — seluruhnya tidak terdampak.

⭐ **AC 35–38 dapat dikerjakan tanpa menunggu** — keempatnya tidak menyentuh tabel itu.

⚠️ **Dan satu ketergantungan yang wajib disebut:** rincian angsuran dapat berasal dari **salinan
master kontrak**, yang berarti ketergantungan pada dokumen kontrak **tetap ada** walau dokumen itu
di luar lingkup penyimpanan.

## EDM Treaty In - 09 - Dua penanda migrasi

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Seluruh dokumen polis dipindahkan, setiap generasinya.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 5.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** **04** · **07** · ⛔ `[work owner]` **lingkup pemindahan dokumen lama**~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **22** *(pemuat dokumen lama)*
**Menutup:** AC **39–43** *(5 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-33 · ID-34 · ID-35 · ID-36 · ID-37

#### Hasil & nilai pengguna

Baris hasil migrasi dibedakan dari baris hasil hitung, dan dua keadaan yang perlu diketahui
**ditandai tanpa mengubah satu angka pun**.

⭐⭐ **Nilai lama tidak pernah dihitung ulang** — termasuk bila rumus lamanya diketahui keliru.
Yang dilakukan hanya **menandainya**, supaya kelak dapat diperlakukan berbeda bila diputuskan
demikian. ⛔ Menghitung ulang berarti laporan lama tidak lagi dapat direkonsiliasi.

#### Yang dibangun

Dua penanda, keduanya **hanya berarti untuk baris hasil migrasi**:

| Penanda | Artinya | Kenapa berguna |
| --- | --- | --- |
| **pasangan bergeser** | kunci dagang pasangan menurut nomor urut ternyata berbeda | ⭐ di endorsemen baris tidak dapat dihapus ⇒ ini **anomali sungguhan** |
| **rumus berlapis** | nilainya lahir dari rumus lama yang mengurangi terhadap selisih | dapat ditentukan **pasti** dari generasi sebelumnya |

⭐ Pemasangan baris antar generasi memakai **nomor urut**; kunci dagang turun pangkat menjadi
**pemeriksa**, bukan pemasang.

#### Batas — yang TIDAK termasuk

⛔ Pemuat massalnya sendiri — tiket **10**.
⛔ Tabel proyeksi tempat penanda ini tinggal — tiket **07**.
⛔ Perhitungan ulang nilai lama — ⛔ **tidak dilakukan sama sekali**.

#### Cara mengujinya

Lewat seam `repository`. ⭐ **Uji utama:** dokumen lama dimuat, penandaan dijalankan, lalu
**seluruh nilai uang dibandingkan dengan sebelum penandaan** — satu angka yang berubah berarti
gagal.

#### Acceptance criteria

- [ ] **AC 39** — nilai hasil migrasi **tidak dihitung ulang**, termasuk digit galatnya
- [ ] **AC 40** — pemasangan antar generasi memakai **nomor urut**
- [ ] **AC 41** — penanda pasangan bergeser dihitung **tanpa mengubah satu angka pun**
- [ ] **AC 42** — penanda rumus berlapis terisi pada setiap baris yang memenuhi syaratnya
- [ ] **AC 43** — kedua penanda **hanya** terisi pada baris hasil migrasi

#### ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Lingkup pemindahan dokumen lama belum diputuskan** — seluruhnya, sebagian, atau
tetap dibaca lewat jalur lama. Penandaan hanya berarti atas baris yang benar-benar dipindahkan.

## EDM Treaty In - 10 - Pemuat migrasi endorsemen

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Seluruh dokumen polis dipindahkan, setiap generasinya.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 5.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** **09** · ⛔ `[work owner]` **lingkup pemindahan dokumen lama**~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **22** *(pemuat dokumen lama)*
**Menutup:** AC **44** *(1 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-3

#### Hasil & nilai pengguna

Dokumen endorsemen lama dimuat **lewat antarmuka yang sama** dengan jalur biasa, sehingga data lama
melewati pemeriksaan yang sama — termasuk aturan keutuhan nomor urut dan larangan percabangan.

⛔ **Nol pintu belakang.** Pemuat yang menulis langsung ke tabel akan melewati seluruh penjaga yang
dibangun tiket **02**, **03**, dan **04**.

#### Yang dibangun

Pemuat massal endorsemen yang memakai **antarmuka penyimpanan yang sama**, memulihkan rantai
generasi dari dokumen lama secara berurutan, dan menyalakan kedua penanda migrasi.

⚠️ **Urutan pemuatan penting:** generasi harus dimuat menurut nomornya, sebab penunjuk generasi
menuntut generasi sebelumnya sudah ada.

#### Batas — yang TIDAK termasuk

⛔ Kedua penanda migrasi itu sendiri — tiket **09**.
⛔ Pemuat dokumen polis baru — tiket NB **22**.
⛔ Keputusan lingkup pemindahan — `[work owner]`.

#### Cara mengujinya

Lewat seam `repository` yang sama. ⭐ **Uji utama:** rantai tiga generasi dimuat dari dokumen lama,
lalu diperiksa bahwa penjaga percabangan dan penjaga keutuhan **benar-benar berjalan** atasnya —
bukan dilewati karena ini jalur migrasi.

#### Acceptance criteria

- [ ] **AC 44** — pemuat migrasi menulis lewat antarmuka penyimpanan yang **sama**

#### ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Lingkup pemindahan belum diputuskan** — sama dengan penahan tiket **09**. Memuat
seluruh riwayat dan memuat dua tahun terakhir adalah pekerjaan yang berbeda besarnya.

⚠️ **Tiket ini sengaja kecil** — satu AC — sebab isinya hampir seluruhnya dipakai bersama tiket NB
**22**. Yang khas endorsemen hanya **urutan pemuatan generasi**.

## EDM Treaty In - 11 - Kepatuhan lapisan dan tipe kolom

> ## ⭐ PENAHAN GUGUR — 23 September 2026 sore
>
> `[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."* ⭐ Presisi dinaikkan ke **`NUMBER(38,8)`** — **30 digit di depan koma**, delapan di belakang, batas tertinggi Oracle. Dasarnya sapuan korpus: ambang dagang nyata sudah **tepat di batas** 12 digit *(`181500000000.00` · `150000000000.00`)*, dan ada sentinel **17 digit** *(`99999999999999999.99`)*. Penjumlahan lintas mata uang dapat melewati keduanya. ⭐ `NUMBER` di Oracle berpanjang **berubah-ubah** — hanya digit bermakna yang tersimpan, sehingga pelebaran ini **tidak memakan ruang tambahan**. Butir ditutup oleh bukti, bukan oleh DBA.
>
> ⭐ **`blocked` → `ready-for-agent`.** ⛔ **Nol butir `[data DBA]` tersisa di tiket ini.**
>
> Rinciannya: `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026 sore)*
~~**Blocked by:** ⛔ `[data DBA]` **presisi fisik kolom uang** — dua belas digit di depan koma belum diuji terhadap nilai terbesar~~ ⛔ **gugur 23-09-2026 sore**
**Bergantung pada tiket NB:** **18** *(tipe kolom dan presisi uang)* · **20** *(transaksi tunggal dan skema eksplisit)*
**Menutup:** AC **45–53** · AC **57–58** *(11 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-1 · ID-2 · ID-43..ID-47

#### Hasil & nilai pengguna

Sembilan aturan kepatuhan yang berlaku sama seperti pada polis baru, ditegakkan juga pada jalur
endorsemen — supaya jalur kedua tidak menjadi celah.

⭐ **Yang paling mudah terlewat:** satu antarmuka penyimpanan, bukan dua. Endorsemen yang mendapat
antarmukanya sendiri akan melewati penjaga yang dibangun untuk polis baru.

#### Yang dibangun

| Aturan | Yang dicegahnya |
| --- | --- |
| **satu transaksi** per generasi | generasi separuh tersimpan |
| skema basis data **eksplisit** | tabel terbaca dari skema yang salah |
| kolom pelaku dari **identitas login** | kolom pelaku kosong seperti di sistem lama |
| ⛔ nol tipe mengambang pada uang | sistem baru **menambah** galat baru |
| skala kolom uang | pemotongan digit yang dianggap presisi |
| kode dan penanda **tetap teks** | nol di depan hilang ⇒ penggolong jenis usaha gagal |
| teks kosong → **tak-bernilai** | kosong menjadi nol, dan nol punya arti lain |
| arah ketergantungan tidak terbalik | lapisan penyimpanan memanggil lapisan layanan |
| **satu** antarmuka penyimpanan | jalur endorsemen melewati penjaga jalur polis baru |

Ditambah panjang minimum kolom keterangan, dan aturan bahwa nilai berdesimal lebih dari skala kolom
**dibulatkan, bukan ditolak**.

#### Batas — yang TIDAK termasuk

⛔ Penetapan tipe kolom itu sendiri — tiket NB **18**; tiket ini **menegakkannya** pada jalur endorsemen.
⛔ `CREATE TABLE` — presisi fisik dicocokkan DBA **di dalam** tiket ini.

#### Cara mengujinya

Lewat seam `repository`, dan ⚠️ **sebagian test memeriksa nilai kolom langsung** — pulang-pergi
saja tidak cukup, sebab tulis dan baca yang sama-sama salah simetris tetap hijau.

⭐ Uji antarmuka tunggal dijalankan sebagai pemeriksaan susunan kode, bukan lewat data.

#### Acceptance criteria

- [ ] **AC 45** — seluruh penyimpanan satu generasi berada dalam **satu transaksi**
- [ ] **AC 46** — setiap query menulis skema **eksplisit**
- [ ] **AC 47** — kolom pelaku diisi dari identitas login
- [ ] **AC 48** — nol kolom uang bertipe mengambang
- [ ] **AC 49** — kolom uang menerima skala desimal yang ditetapkan
- [ ] **AC 50** — kode dan penanda tersimpan sebagai **teks**, nol di depan utuh
- [ ] **AC 51** — teks kosong pada kolom angka atau tanggal tersimpan **tak-bernilai**
- [ ] **AC 52** — arah ketergantungan **tidak pernah dibalik**
- [ ] **AC 53** — hanya ada **satu** antarmuka penyimpanan polis treaty
- [ ] **AC 57** — kolom keterangan berpanjang sekurangnya batas yang ditetapkan
- [ ] **AC 58** — nilai berdesimal lebih dari skala kolom **dibulatkan, bukan ditolak**

#### ⛔ Kenapa tiket ini `blocked`

~~`[data DBA]` **Dua belas digit di depan koma belum diuji.**~~ ⛔ **BUTIR GUGUR 23-09-2026 sore.** Presisi dinaikkan ke **`NUMBER(38,8)`** — tiga puluh digit di depan koma, delapan di belakang. Ditutup oleh bukti korpus, bukan oleh jawaban DBA.

⚠️ **Dan satu pertentangan di dalam spec, dicatat di sini karena spec tidak boleh disunting:**
`[terverifikasi]` **AC 49** berbunyi *"kolom uang menerima **sekurangnya sembilan** angka di belakang
koma"*, sementara **AC 58** pada spec yang sama menetapkan skala **delapan** dan menyatakan
kelebihannya **dibulatkan**.

⛔ Nilai produksi berdesimal sembilan **berubah** pada skala delapan. Spec penyimpanan polis baru
sudah diselaraskan ke delapan dan **mengutip bunyi lamanya**; spec ini **tertinggal**.
⛔ **Tidak diputuskan di tiket ini.** Pemiliknya `[work owner]` dan `[data DBA]`.

## EDM Treaty In - 12 - Medan catatan dan penguncian daftar kolom

> ## ⭐ PENAHAN GUGUR — 23 September 2026 sore
>
> `[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`, ikuti dari data yang digunakan di Activity dan Section."* ⭐ Daftar medan disusun ulang dari **302 berkas aturan** — **394 medan unik**, 232 di antaranya tampil di layar. Hasilnya di **`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`**. Data guide turun derajat jadi **penambal**, sah hanya untuk panjang maksimum per medan.
>
> ⭐ **`blocked` → `ready-for-agent`.** ⛔ **Nol butir `[data DBA]` tersisa di tiket ini.**
>
> Rinciannya: `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026 sore)*
~~**Blocked by:** ⛔ `[data DBA]` **panduan bentuk dokumen terbukti basi** — satu dokumen memuat~~ ⛔ **gugur 23-09-2026 sore**
95 jalur yang tidak ada di dalamnya
**Bergantung pada tiket NB:** **19** *(pemecah dokumen menjadi baris)*
**Rujukan:** `ID-27b` · `ID-27c`
**Menutup:** dua implementation decision yang tidak tercakup ronde tiket sebelumnya

> ⚠️ **Tiket ini lahir dari audit sesudah ronde, bukan dari ronde tiket.** Kedua butirnya
> disisipkan ke spec sesudah spec pertama selesai, sehingga tidak terbaca saat tiket 01–11 ditulis.
> Kekeliruan lingkup itu milik penyusun brief.

#### Hasil dan nilai pengguna

Dua hal kecil yang, bila terlewat, menghasilkan kerugian yang tidak kelihatan sampai terlambat.

Yang pertama satu medan catatan bebas yang selama ini ikut tersimpan pada dokumen polis dan tidak
pernah masuk rancangan tabel. Bila ia hilang saat pemindahan, isinya tidak dapat dipulihkan —
tidak ada tempat lain yang menyimpannya.

Yang kedua bukan medan melainkan **disiplin**. Daftar kolom disusun dari medan yang benar-benar dipakai aturan — Activity dan Section — bukan dari panduan bentuk dokumen, yang terbukti tidak lengkap. Menyusunnya dari panduan berarti membangun pemecah dokumen yang diam-diam membuang medan yang tidak dikenalinya.

#### Lingkup

##### Bagian 1 · medan catatan *(`ID-27b`)*

- Medan catatan bebas tingkat atas, panjang **128**, **ikut dipindahkan**.
- ⛔ Dua medan tampilan yang bertetangga dengannya **tidak** dipindahkan — keduanya keadaan layar,
  bukan data dagang. Keputusan lama, tetap berlaku.
- Uji: dokumen yang membawa medan catatan terisi, dipindahkan lalu dibaca kembali, isinya sama
  persis termasuk spasi di ujung.

##### Bagian 2 · penguncian daftar kolom *(`ID-27c`)*

- ⭐ **Dasar daftar kolom adalah `DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`** — 394 medan yang disapu dari 302 berkas aturan, bukan panduan bentuk dokumen.
- Panduan bentuk dokumen berkedudukan **penambal**: sah untuk panjang maksimum per medan, dan untuk 24 nama yang ditulis sistem sehingga tidak tersapu aturan.
- ⛔ Ia **tidak sah** dipakai untuk membuktikan bahwa sebuah medan tidak ada.
- Pemecah dokumen wajib punya **penampung medan tak dikenal**, dan penampung itu wajib **kosong**
  sebelum pekerjaan dinyatakan selesai.
- Uji: dokumen yang membawa medan di luar daftar tidak boleh diam-diam kehilangan medan itu —
  ia masuk penampung, dan keberadaan isi di penampung menggagalkan test.

#### Batas — yang TIDAK termasuk

- ⛔ Bukan tempat menetapkan daftar kolom akhir. Itu tiket NB 19.
- ⛔ Bukan tempat memutuskan presisi fisik kolom uang. Itu tiket EDM 11.
- ⛔ Tidak membuat tabel apa pun.

#### Uji

Lewat seam `repository`, seperti tiket lain di berkas ini.

⚠️ Uji pulang-pergi saja **tidak cukup** untuk bagian 2 — test wajib memeriksa **isi penampung
medan tak dikenal** secara langsung, karena dokumen yang kehilangan medan tetap lolos pulang-pergi.

#### Butir tertahan

| Butir | Pemilik |
| --- | --- |
| ~~⛔ panduan bentuk dokumen perlu disegarkan~~ ✅ **gugur 23-09 sore** — panduan diabaikan, daftar medan disusun dari korpus | ~~`[data DBA]`~~ |

⛔ Butir ini **tidak ditutup di tiket ini**. Surat permintaannya sudah disusun.

#### Area codebase

- Lapisan `repository`: pemecah dokumen dan penampung medan tak dikenal
- Pemuat pemindahan dokumen lama

#### Disiplin berkas

⛔ Nol `CREATE TABLE`. ⛔ Nol nama orang. ⛔ Nol nomor polis harfiah. ⛔ Nol cuplikan data produksi.

---

*Ditulis 23 September 2026, sesudah audit pasca-ronde menemukan dua implementation decision yang
tidak tercakup tiket mana pun.*
