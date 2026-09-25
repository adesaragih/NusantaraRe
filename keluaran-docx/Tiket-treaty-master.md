# Tiket - Kontrak Treaty dan Adjustment

> Dokumen ini memuat **badan tiket lengkap**, disusun per modul lalu per nomor.
> Disusun 25 September 2026 dari berkas tiket proyek migrasi Nusantara Re.

## Matriks status

| Modul | Tiket | Siap | Tertahan | needs-info | wontfix | Lain |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Treaty In | **52** | 0 | 0 | 0 | 0 | 52 |
| Treaty In Adjustment | **14** | 0 | 0 | 0 | 0 | 14 |
| **Jumlah** | **66** | **0** | **0** | **0** | **0** | **66** |

---

# Treaty In

Jumlah tiket: **52**

## Treaty In - 14 - Kontrak dan versi pertamanya berdiri, dan dapat ditemukan kembali dengan pengenalnya

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-01` (**pecahan identitas** — lihat `Tidak termasuk`) · `SPEC-MODEL-DATA.md` §10.1, §10.2 · `ADR-D-TI-0040`.*

**What to build:** **PK** membuat sebuah kontrak baru, mengisi kepala versinya, menyimpannya, dan
**menemukannya kembali dengan pengenal yang sistem berikan**. Identitasnya terpecah dua sebagaimana
`ADR-D-TI-0040` menetapkan: `KONTRAK` memegang lapisan beku — cedant, asal bisnis, sifat proporsi,
periode — dan `VERSI_KONTRAK` memegang segala yang dapat berbeda antar versi.

Artefak: tabel `KONTRAK` dan `VERSI_KONTRAK` seluruhnya, kunci utamanya, kunci asing di antara
keduanya, dan sequence pengenalnya.

**PEMBUAT PERTAMA** untuk `KONTRAK` dan `VERSI_KONTRAK`.

**Kenapa begini:** **Tidak ada satu pun tiket di papan yang membuat `VERSI_KONTRAK`.** Tiket `01`,
`02`, dan `05` menambahkan **kolom** padanya; tiket `04` menambahkan tabel yang merujuknya. Keempatnya
mengandaikan tabelnya sudah ada, dan pembuatnya tidak pernah ditiketkan — sebab ia milik `P-01`,
kemampuan Treaty In, yang ronde tiket sebelumnya memang di luar lingkup. Irisan ini menutup lubang
itu, dan karena itu **seluruh papan menggantung padanya**.

**Persyaratan:** `INV-01` (setiap tabel berkunci utama) · `INV-02` (pengenal dari `SEQUENCE`, tidak
pernah dari cap waktu maupun teks) · `INV-03` (pengenal tidak dipakai ulang, sequence tanpa `CYCLE`) ·
`INV-04` (`NOMOR_URUT_VERSI` unik di dalam satu `KONTRAK`) · `INV-18` (perilaku hapus tiap kunci asing
ditetapkan sadar) · `INV-29` (`SIFAT_PROPORSI` dua nilai) · `INV-53` (`TANGGAL_MULAI` ≤
`TANGGAL_BERAKHIR`, keduanya inklusif) · `ADR-D-TI-0040` · nama dan tipe kolom **mengikat** pada
`2-to-spec/KAMUS-KOLOM.md` (urutan wewenang butir 5)

**Tidak termasuk:** **Daftar nilai sah `KEADAAN_SIKLUS_HIDUP` dan mesin perpindahannya.** Kolomnya
berdiri di sini — tabelnya tidak dapat berdiri tanpanya — tetapi `INV-20`, `INV-22`, `INV-23`,
`INV-24`, dan `INV-25` **tegas bukan bagian irisan ini**. Ia batch 2, sebab `REV-3` merevisi ADR-D-TI-0055
§4 yang menjadi sumber daftar keadaan itu. **Ini pecahan yang disengaja**, dan bunyi lengkap `P-01`
memuat keduanya.
**Peringatan kunci alami ganda** — irisan `16`. **Pencarian lewat nomor warisan** — irisan `17`.
**Pembekuan kunci alami** — irisan `18`.
**Mesin pro rata** — `GRL-15` memutuskan ia **sengaja tidak dibangun**; atribut *"berlaku sejak"*
dibawa, mesinnya tidak. Itu **pernyataan keputusan**, bukan pekerjaan tertunda.

**Jalur gagal:** Menyimpan versi tanpa kontrak induk -> ditolak kunci asing · Dua versi bernomor urut
sama pada satu kontrak -> **ditolak** `INV-04`, pesannya menyebut nomor yang bentrok · Kontrak
bertanggal berakhir lebih awal daripada tanggal mulai -> ditolak `INV-53` · Pengenal yang diminta
dari cap waktu atau dari teks -> **tidak ada jalurnya**; satu-satunya sumber pengenal adalah sequence.

**Uji:** **Negatif:** simpan versi yatim; simpan dua versi bernomor urut sama; simpan kontrak
bertanggal terbalik; minta pengenal berulang sesudah sequence berputar penuh (`INV-03`).
**Positif — dan ia yang menangkap pemisahan identitas yang terlalu ketat:** satu kontrak dengan
**tiga** versi berturut-turut **diterima**, dan ketiganya berbagi lapisan beku yang sama tanpa
menyalinnya. Constraint yang menolak versi kedua lulus setiap uji negatif yang pernah ditulis untuk
irisan ini.

**Menggantikan:** `P-01` bunyi lama tidak bergeser; yang bergeser **cakupan irisannya**, bukan
kemampuannya. Sistem lama menyimpan kontrak dan addendum di **satu baris `M_TREATY_IN`** dengan
seluruh halaman clipboard sebagai satu kolom `JSONDATA`, dan pengenalnya `ID` bertipe teks yang
dibentuk dari pola bernomor revisi (`TDA-11`). Identitas yang terpecah dua **tidak pernah ada** di
sana.

**Blocked by:** None (can start immediately)

**Dasar:**
```
EVIDENCED(M_TREATY_IN@Table/, TREATYINDETAIL@Table/, SaveTreatyInDetail_Act@ekspor-2026-09)
        DECIDED(ADR-D-TI-0040, ADR-D-TI-0042, KTV-A, KTV-C)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(T-6)
```

- [ ] `KONTRAK` dan `VERSI_KONTRAK` berdiri dengan **seluruh** kolom yang `2-to-spec/KAMUS-KOLOM.md` sebutkan, bernama dan bertipe persis
- [ ] pengenal keduanya datang dari `SEQUENCE` tanpa `CYCLE`; tidak ada jalur lain yang dapat memberi pengenal
- [ ] `INV-04` menolak nomor urut versi ganda di dalam satu kontrak, pesannya menyebut nomornya
- [ ] perilaku hapus kunci asing `VERSI_KONTRAK` → `KONTRAK` **ditetapkan sadar dan tertulis**, bukan dibiarkan bawaan (`INV-18`)
- [ ] uji positif lulus: satu kontrak dengan tiga versi diterima, lapisan bekunya tidak disalin
- [ ] kolom `KEADAAN_SIKLUS_HIDUP` **ada**, dan ketiadaan constraint daftar nilainya **tertulis sebagai pernyataan keputusan** di dalam berkas DDL-nya — bukan dibiarkan terbaca sebagai kelalaian
- [ ] `KTV-A` dan `T-6` tercatat di `ASUMSI-CLEAR.md` sebagai asumsi yang tiket ini bersandar padanya

## Treaty In - 15 - Himpunan acuan bertambah tanpa mengubah arti apa pun yang sudah tercatat

---
status: aktif
golongan: perubahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-48` · `SPEC-MODEL-DATA.md` §10.22 · `ADR-D-TI-0038`.*

**What to build:** Keenam himpunan yang dapat bertambah — **mata uang, jenis potongan, jenis
reasuransi, bahaya, kelompok treaty, kelas bisnis** — berdiri sebagai **tabel acuan**, dan menambah
satu baris baru **tidak menyentuh satu pun kontrak yang sudah tercatat** dan **tidak menuntut
perubahan skema**.

Artefak: keenam tabel acuan, kunci alaminya, dan penanda aktifnya.

**PEMBUAT PERTAMA** untuk `MATA_UANG`, `JENIS_POTONGAN`, `JENIS_REASURANSI`, `BAHAYA`,
`KELOMPOK_TREATY`, `KELAS_BISNIS`.

**Kenapa begini:** Sistem lama menaruh **empat bahaya sebagai empat pasang kolom di kepala kontrak** —
`Earthquake`, `FloodJab`, `FloodNation`, `RSMDLimit`, masing-masing berpasangan dengan kolom mata
uangnya. Bentuk itu menuntut **kolom kesembilan** begitu ada bahaya baru, dan bahaya baru adalah hal
yang pasti terjadi di reasuransi. `ADR-D-TI-0038` membalikkannya: **yang bertambah tanpa mengubah arti
disimpan sebagai data, bukan sebagai nama kolom dan bukan sebagai `CHECK`.**

**Persyaratan:** `INV-62` (himpunan yang dapat bertambah disimpan sebagai tabel acuan, bukan `CHECK`) ·
`INV-68` (`KODE` unik di seluruh tabel, pada keenamnya) · `INV-44` (kode mata uang merujuk tabel acuan
mata uang) · `INV-01` · `INV-02` · `ADR-D-TI-0038`

**Tidak termasuk:** **Pemindahan isinya dari sistem lama** — irisan `44`. Irisan ini membuat
tabelnya dan membuktikan ia dapat bertambah; mengisinya dari sumber lama pekerjaan lain.
**`JENIS_REASURANSI` bersusun** — kolom `ID_INDUK` yang menunjuk baris lain di tabel yang sama
**ikut di sini**, tetapi arti susunannya dipakai pertama kali oleh irisan `38`.
**Pemilik tabel acuan pembagian kapasitas** — `P-49`, dan `G2` menyatakannya **tegas tidak masuk**
penyerahan pertama: ia mewujudkan butir eskalasi 7 yang belum diputuskan.

**Jalur gagal:** Dua baris berkode sama di satu tabel acuan -> **ditolak** `INV-68`, pesannya
menyebut kodenya · Kode mata uang pada kontrak yang tidak ada di tabel acuan -> ditolak `INV-44` ·
Sebuah `CHECK` berisi daftar nilai bahaya -> **tidak ada**, dan ketiadaannya diperiksa, bukan
diandaikan.

**Uji:** **Negatif:** sisipkan kode ganda di tiap tabel; rujuk mata uang yang tidak terdaftar.
**Positif — dan ia pokok irisan ini:** tambahkan **bahaya kesembilan** ke tabel acuan, lalu catat
batas tanggungan untuknya pada sebuah kontrak. Keduanya **berhasil tanpa satu pun perubahan skema**.
Uji negatif saja tidak dapat membuktikannya: sebuah tabel acuan yang benar dan sebuah tabel acuan
yang `CHECK`-nya masih tertinggal di tempat lain **menolak hal yang sama**.

**Menggantikan:** `P-48` tidak menggantikan aturan bernomor cacat, melainkan **bentuk** yang §10.0b
uraikan: delapan kolom bernama bahaya di kepala kontrak. Golongannya **PERUBAHAN** — sistem lama
melakukannya, dan kita sengaja mengubahnya — sehingga hasilnya **tidak dapat diuji terhadap data
lama**: data lama justru menunjukkan bentuk yang dibuang.

**Blocked by:** None (can start immediately)

**Dasar:**
```
EVIDENCED(SPEC-MODEL-DATA.md §10.0b - delapan kolom bernama bahaya, dibaca dari kepala TreatyIn)
        DECIDED(ADR-D-TI-0038, INV-62, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] keenam tabel acuan berdiri dengan bentuk seragam — kode, nama, penanda aktif — sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-68` menolak kode ganda di **keenamnya**, diperiksa satu per satu bukan disimpulkan dari satu
- [ ] **nol `CHECK` berisi daftar nilai** untuk keenam himpunan ini; diperiksa dengan sapuan, hasilnya dicetak
- [ ] `JENIS_REASURANSI` membawa `ID_INDUK` yang menunjuk barisnya sendiri, dengan perilaku hapusnya ditetapkan sadar
- [ ] uji positif lulus: bahaya kesembilan ditambahkan dan langsung dapat dipakai, **tanpa perubahan skema**
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 16 - Pengisi kontrak diperingatkan bahwa kunci alaminya sama dengan kontrak lain, tanpa dihalangi menyimpan

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-02` · `SPEC-MODEL-DATA.md` §10.1 · `ADR-D-TI-0040` §2.*

**What to build:** **PK** yang menyimpan kontrak berkunci alami sama dengan kontrak lain — cedant +
asal bisnis + tanggal mulai + sifat proporsi — **tetap dapat menyimpannya**, dan melihat peringatan
yang **menyebut kontrak pembandingnya**.

**Kenapa begini:** `ADR-D-TI-0040` §2 memutuskan kunci alami kontrak **memperingatkan, tidak melarang** —
dan bedanya bukan selera. Kontrak yang benar-benar kembar memang terjadi: pembaruan tahunan yang
dicatat ulang, atau dua perjanjian terpisah dengan cedant yang sama pada tanggal yang sama.
Melarangnya berarti **menolak data yang sah**, yaitu kekeliruan lingkup yang sudah tiga kali terjadi
di modul ini. Dan peringatan tanpa menyebut pembandingnya adalah peringatan yang tidak dapat
ditindaklanjuti: yang membacanya harus dapat membuka kontrak yang dimaksud.

**Persyaratan:** `ADR-D-TI-0040` §2 · `SPEC-MODEL-DATA.md` §10.1 *("kunci alami dan lingkupnya —
`ID_CEDANT` + `ID_ASAL_BISNIS` + `TANGGAL_MULAI` + sifat proporsi")* · **ketiadaan nomor `INV`**-nya
**disengaja** dan dinyatakan di §10.1 — ia bukan constraint

**Tidak termasuk:** **`UNIQUE` atas kunci alami kontrak** — tegas tidak dipasang, dan ketiadaannya
adalah **pernyataan keputusan** yang sudah berdiri di `ddl-usulan/Z00_KUNCI_ALAMI.sql`. Siapa pun yang
menambahkannya membalikkan `ADR-D-TI-0040` §2 tanpa membukanya.
**Pembekuan kunci alami sesudah kontrak lahir** — irisan `18`. Yang di sini soal **lahirnya**, yang di
sana soal **mengubahnya**.

**Jalur gagal:** Menyimpan kontrak kembar -> **BERHASIL**, dan peringatannya muncul menyebut pengenal
kontrak pembandingnya · Peringatan muncul tanpa menyebut pembandingnya -> **kriteria selesai tidak
terpenuhi** · Penyimpanan ditolak -> **kriteria selesai tidak terpenuhi**, sebab itu perilaku yang
`ADR-D-TI-0040` §2 larang.

**Uji:** **Negatif — dan di irisan ini "negatif" berarti sesuatu yang lain:** yang harus gagal bukan
penyimpanannya melainkan **kebisuan**. Simpan kontrak kembar dan periksa bahwa **peringatannya ada**;
ketiadaan peringatan adalah kegagalan.
**Positif:** dua kontrak yang **berbeda pada satu ruas kunci saja** — misalnya sifat proporsi —
disimpan **tanpa** peringatan. Tanpa uji ini, peringatan yang terlalu lebar akan berbunyi untuk
hampir setiap kontrak dan orang berhenti membacanya.

**Menggantikan:** `P-02` tidak menggantikan aturan sistem lama, sebab **sistem lama tidak punya
pemeriksaan kembar sama sekali** pada kunci ini. Yang ada hanya `CheckDuplicateOffer`, yang bekerja
atas **penawaran** dan berbunyi *"Have similarities in SoB, Ceding, and Period, please check for
duplicates"* — sumbu yang **mirip tetapi tidak sama**: ia tidak memuat sifat proporsi. Golongan
**PELESTARIAN** dipertahankan karena perilakunya — memperingatkan, bukan menolak — memang
dilestarikan; yang berubah **sumbunya**.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(CheckDuplicateOffer@ekspor-2026-09 - "Have similarities in SoB, Ceding, and Period")
        DECIDED(ADR-D-TI-0040)
```

- [ ] menyimpan kontrak berkunci alami kembar **berhasil**, dan tidak ada constraint yang menolaknya
- [ ] peringatannya menyebut **pengenal kontrak pembandingnya**, bukan hanya menyatakan ada kembar
- [ ] uji positif lulus: beda pada satu ruas kunci saja -> **tidak ada peringatan**
- [ ] ketiadaan `UNIQUE` atas kunci alami kontrak tertulis sebagai **pernyataan keputusan**, dan sapuan atas `ddl-usulan/` membuktikannya masih berupa pernyataan — bukan sudah dipasang diam-diam
- [ ] perbedaan sumbu terhadap `CheckDuplicateOffer` **tercatat**: sistem lama tidak memuat sifat proporsi

## Treaty In - 17 - Kontrak ditemukan kembali lewat nomor warisan sistem lama, berdampingan dengan pengenal baru

---
status: aktif
golongan: baru
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-03` · `SPEC-MODEL-DATA.md` §10.1 · `ADR-D-TI-0042`.*

**What to build:** **PK** yang memegang **nomor kontrak sistem lama** — angka yang tertulis di kertas,
di surel, dan di sistem tetangga — dapat menemukan kontraknya di sistem baru dengan mengetik nomor
itu, **tanpa perlu tahu pengenal barunya**.

**Kenapa begini:** Pada hari peralihan, **setiap rujukan yang beredar di luar sistem memakai nomor
lama**. Bordereau yang sudah dikirim, surat yang sudah ditandatangani, dan sistem hilir yang belum
dipindahkan seluruhnya menyebut nomor itu. Tanpa jalan masuk lewat nomor lama, orang yang memegang
selembar kertas **tidak punya cara menemukan barisnya** — dan yang terjadi adalah ia membuat kontrak
baru.

**Persyaratan:** `ADR-D-TI-0042` (sejarah pindah apa adanya) · `SPEC-MODEL-DATA.md` §10.1
`NOMOR_KONTRAK_WARISAN` *("hanya terisi pada baris hasil migrasi")* · `INV-02` (pengenal baru tetap
dari sequence; nomor lama **tidak pernah** menjadi pengenal)

**Tidak termasuk:** **Pengisian `NOMOR_KONTRAK_WARISAN`** — itu jalur migrasi, dan migrasinya batch 2.
Irisan ini membangun **jalan masuknya** dan mengujinya dengan baris uji.
**`NOMOR_PENAWARAN_WARISAN`** — kolom kedua yang `KEPUTUSAN-TANPA-VERIFIKASI.md` §6 putuskan,
bergantung pada **Uji AL**. Ia dicari lewat jalur yang sama begitu terisi, tetapi keputusan
mencabutnya **belum jatuh**.
**Pencarian menurut cedant, periode, atau keadaan** — irisan batch 2 (`P-56`), sebab ia mencari
**menurut keadaan siklus hidup**.

**Jalur gagal:** Mencari nomor warisan yang tidak ada -> hasil **kosong yang dinyatakan**, bukan galat
dan bukan daftar penuh · Baris non-warisan bernomor warisan kosong -> **tidak muncul** pada pencarian
nomor warisan mana pun, dan **tidak menghalangi** pencarian yang lain · Nomor warisan dipakai sebagai
kunci gabung ke tabel lain -> **tidak ada jalurnya**; ia kolom pencarian, bukan pengenal.

**Uji:** **Negatif:** cari nomor yang tidak ada; pastikan baris berkolom kosong tidak ikut terambil
oleh pencarian bernilai kosong.
**Positif — dan ia yang menangkap jalur yang terlalu sempit:** dua kontrak warisan bernomor lama
**berbeda** ditemukan masing-masing, dan sebuah kontrak yang **lahir di sistem baru** — nomor
warisannya kosong — tetap ditemukan lewat pengenal barunya. Jalur pencarian yang diam-diam menuntut
nomor warisan terisi lulus setiap uji negatif yang pernah ditulis untuknya.

**Menggantikan:** tidak ada. Golongannya **BARU**: sistem lama tidak punya "nomor warisan" sebab ia
**adalah** sistem lamanya. `ID`-nya bertipe teks dan menjadi pengenal sekaligus pembawa nomor revisi
lewat pola `@substring(ID,10,12)` — dua arti di satu kolom, dan itu persis yang `ADR-D-TI-0040` larang
dibawa.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **PK yang menangani pertanyaan masuk dari cedant**, pada setiap pertanyaan yang menyebut nomor lama. Bukan pemantau berkala: kegagalannya muncul saat seseorang mencari dan tidak menemukan |
| 3 | **ambang berangka** — nyala penuh ketika **seluruh** kontrak warisan sudah punya `NOMOR_KONTRAK_WARISAN` terisi, dan pencarian atas contoh acak **50 nomor lama** menemukan **50**. Di bawah itu jalurnya tetap hidup tetapi kekosongan dilaporkan, bukan didiamkan |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(SaveTreatyInDetail_Act@ekspor-2026-09 - SaveData.ID; @substring(ID,10,12) sebagai nomor revisi)
        DECIDED(ADR-D-TI-0040, ADR-D-TI-0042)
```

- [ ] pencarian lewat `NOMOR_KONTRAK_WARISAN` menemukan barisnya, dan hasil kosong **dinyatakan kosong**
- [ ] baris bernomor warisan kosong tidak ikut terambil oleh pencarian bernilai kosong
- [ ] nomor warisan **tidak** dipakai sebagai kunci gabung di mana pun — diperiksa dengan sapuan
- [ ] uji positif lulus: kontrak yang lahir di sistem baru tetap ditemukan lewat pengenal barunya
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk **ambang berangka**
- [ ] `NOMOR_PENAWARAN_WARISAN` dinyatakan **di luar irisan ini**, dengan Uji AL sebagai penagihnya

## Treaty In - 18 - Kunci alami kontrak tidak dapat diubah sesudah kontraknya lahir, dan penolakannya menyebut apa yang diubah

---
status: aktif
golongan: baru
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-05` · `SPEC-INVARIAN.md` `INV-19` · `ADR-D-TI-0040`.*

**What to build:** Sesudah sebuah kontrak tersimpan, **cedant, asal bisnis, sifat proporsi, tanggal
mulai, dan tanggal berakhirnya tidak dapat diubah oleh siapa pun**. Percobaannya ditolak, dan
pesannya **menyebut ruas mana** yang ia coba ubah.

**Kenapa begini:** Kelima ruas itu **adalah identitas kontraknya**. Mengubah salah satunya bukan
menyunting kontrak melainkan **menjadikannya kontrak yang lain**, sementara seluruh versi, layer,
penyebaran, dan catatan persetujuan di bawahnya tetap menggantung padanya. Sistem lama membiarkannya:
`SaveTreatyInDetail_Act` menulis ulang seluruh kepala pada setiap penyimpanan, tanpa satu pun
pemeriksaan bahwa yang beku tetap beku. Dan penolakan yang tidak menyebut ruasnya memaksa orang
menebak — pada borang berkolom puluhan, menebak berarti mencoba satu per satu.

**Persyaratan:** **`INV-19`** — *"kunci alami kontrak tidak berubah sepanjang hidup kontrak"*,
bergolongan **TRIGGER** · `ADR-D-TI-0040` · `SPEC-MODEL-DATA.md` §10.1 · `INV-29` · `INV-53`

**Tidak termasuk:** **Peringatan kunci alami ganda saat lahir** — irisan `16`. Yang di sini soal
**mengubahnya**, yang di sana soal **kembarannya**.
**Lapisan beku yang tidak berubah antar versi** — irisan `19`. Bedanya nyata: irisan ini menjaga
**satu baris `KONTRAK`** terhadap `UPDATE`; irisan `19` menjaga **versi baru** agar tidak membawa
nilai yang berbeda.
**Sakelar penegakan trigger selama migrasi** — irisan `43`. Trigger ini **harus dapat dimatikan** saat
pemindahan, dan mekanismenya milik irisan itu.

**Jalur gagal:** `UPDATE` atas `ID_CEDANT` pada kontrak yang sudah ada -> **ditolak** `INV-19`,
pesannya menyebut `ID_CEDANT` · `UPDATE` atas dua ruas beku sekaligus -> ditolak, pesannya menyebut
**keduanya**, bukan yang pertama saja · `UPDATE` atas ruas yang **bukan** kunci alami -> **diterima**.

**Uji:** **Negatif:** ubah masing-masing dari kelima ruas, satu per satu — **lima uji, bukan satu**.
Satu uji atas satu ruas tidak membuktikan keempat lainnya dijaga, dan trigger yang hanya memeriksa
ruas pertama lulus uji tunggal mana pun.
**Positif — dan ia yang menangkap trigger yang terlalu lebar:** ubah `NAMA_KONTRAK` dan
`LINGKUP_WILAYAH` pada versi yang sama -> **diterima**. Trigger yang menolak setiap `UPDATE` atas
`KONTRAK` lulus kelima uji negatif di atas dan **tetap salah**.

**Menggantikan:** tidak ada aturan lama yang digantikan — **golongannya BARU**, dan itu yang membuat
medan `CARA MENYALAKANNYA` wajib. Sistem lama **tidak pernah memeriksanya**; kontrak yang cedantnya
diganti diam-diam akan lolos, dan **tidak ada cara mengetahui berapa kali itu terjadi**, sebab
penyimpanan tidak meninggalkan jejak siapa maupun kapan (`ADR-D-TI-0045` konteks).

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemilik proses Treaty In, mingguan**, atas catatan pelanggaran mode peringatan |
| 3 | **ambang berangka** — kuasa memblokir dinyalakan ketika **nol pelanggaran tercatat selama 4 minggu berturut-turut**. Bila masih ada, tiap pelanggaran diperiksa satu per satu lebih dulu: ia dapat berarti data lama yang memang perlu dibetulkan, bukan orang yang salah |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(SaveTreatyInDetail_Act@ekspor-2026-09 - kepala ditulis ulang tiap simpan, nol pemeriksaan)
        DECIDED(ADR-D-TI-0040, INV-19)
```

- [ ] `INV-19` terpasang sebagai **trigger**, bukan sebagai pemeriksaan aplikasi
- [ ] kelima ruas diuji **satu per satu**; lima uji negatif lulus
- [ ] pesan penolakan menyebut **seluruh** ruas yang dicoba, bukan yang pertama saja
- [ ] uji positif lulus: ruas non-kunci tetap dapat diubah
- [ ] trigger ini **terdaftar** di daftar yang irisan `43` matikan saat pemindahan — bukan ditemukan belakangan saat migrasi gagal
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk ambang berangka

## Treaty In - 19 - Lapisan beku tidak berubah antar versi

---
status: aktif
golongan: baru
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-43` · `SPEC-MODEL-DATA.md` §10.1 · `ADR-D-TI-0040`.*

**What to build:** Versi baru pada kontrak yang sama **tidak dapat membawa** cedant, asal bisnis,
sifat proporsi, atau periode yang berbeda dari kontraknya. Percobaannya ditolak, dan pesannya
menyebut ruas yang menyimpang beserta nilai yang berlaku.

**Kenapa begini:** Pemisahan `KONTRAK` / `VERSI_KONTRAK` hanya berarti bila **lapisan bekunya
benar-benar beku**. Bila sebuah versi dapat menyatakan cedant yang berbeda, maka kontraknya bukan
lagi satu benda, dan pertanyaan *"siapa cedant kontrak ini"* punya jawaban yang berbeda-beda
tergantung versi mana yang dibaca. Di sistem lama pertanyaan itu memang punya jawaban berbeda-beda —
setiap addendum menyimpan salinan penuh kepala kontraknya di `JSONDATA`-nya sendiri, dan tidak ada
yang mencocokkannya.

**Persyaratan:** `ADR-D-TI-0040` · `SPEC-MODEL-DATA.md` §10.1 · `INV-29` (sifat proporsi dua nilai) ·
`INV-53` (tanggal mulai ≤ tanggal berakhir) · `INV-59` (tidak ada nilai yang disalin dari entitas
lain; hilir diberi rujukan)

**Tidak termasuk:** **Pembekuan satu baris `KONTRAK` terhadap `UPDATE`** — irisan `18`.
**Mesin pro rata.** `GRL-15` memutuskan **ia sengaja tidak dibangun**: atribut *"berlaku sejak"*
dibawa, mesinnya tidak. Ini **pernyataan keputusan**, bukan `TODO` — siapa pun yang membangunnya
membalikkan `GRL-15` tanpa membukanya. Akibat yang disengaja: perubahan yang berlaku di tengah
periode **dicatat tanggalnya** dan **tidak dihitung prorata oleh sistem**.
**Keadaan versi baru saat lahir** — batch 2.

**Jalur gagal:** Menyimpan versi kedua yang membawa `ID_CEDANT` berbeda -> **ditolak**, pesannya
menyebut ruas dan nilai yang berlaku · Menyalin lapisan beku ke kolom pada `VERSI_KONTRAK` ->
**tidak ada kolomnya**; `INV-59` melarangnya dan ketiadaannya diperiksa dengan sapuan atas
`KAMUS-KOLOM.md`.

**Uji:** **Negatif:** untuk masing-masing dari kelima ruas beku, buat versi kedua yang menyimpang.
**Lima uji.**
**Positif — dan ia yang menangkap pemeriksaan yang terlalu lebar:** versi kedua yang mengubah
`NAMA_KONTRAK`, `LINGKUP_WILAYAH`, dan seluruh kepala yang **memang boleh berbeda** -> **diterima**.
Sebuah pemeriksaan yang membandingkan seluruh kepala versi terhadap versi sebelumnya lulus kelima uji
negatif dan **membekukan hal yang tidak pernah dimaksudkan beku**.

**Menggantikan:** tidak ada. **Golongannya BARU.** Sistem lama menyimpan salinan penuh kepala di
setiap addendum, sehingga *"tidak berubah antar versi"* bukan aturan melainkan **kebetulan** — dan
kebetulan yang tidak pernah diperiksa. Uji apa pun terhadap data lama yang menemukan penyimpangan
**bukan bukti sistem baru salah**; ia bukti bahwa aturan ini memang belum pernah ada.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemilik proses Treaty In, mingguan**. Pelanggaran di sini menunjuk **data warisan yang memang menyimpang**, bukan orang yang salah mengisi |
| 3 | **ambang berangka** — kuasa memblokir dinyalakan ketika penyimpangan pada baris **yang lahir di sistem baru** nol selama 4 minggu. Baris warisan **dihitung terpisah** dan tidak ikut menahan penyalaan |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(M_TREATY_IN_EDM@Table/ - tiap addendum menyimpan salinan penuh kepala di JSONDATA-nya sendiri)
        DECIDED(ADR-D-TI-0040, GRL-15, INV-59)
```

- [ ] kelima ruas beku diuji **satu per satu**; lima uji negatif lulus
- [ ] pesan penolakan menyebut ruas **dan** nilai yang berlaku pada kontraknya
- [ ] uji positif lulus: kepala versi yang memang boleh berbeda tetap dapat berbeda
- [ ] sapuan atas `KAMUS-KOLOM.md` membuktikan **tidak ada** kolom salinan lapisan beku di `VERSI_KONTRAK`
- [ ] pemisahan hitungan baris warisan versus baris baru **tertulis**, bukan digabung
- [ ] ketiadaan mesin pro rata tertulis sebagai **pernyataan keputusan** yang menyebut `GRL-15`
- [ ] keempat butir **CARA MENYALAKANNYA** terisi

## Treaty In - 20 - Lebih dari satu mata uang berlaku pada satu kontrak, masing-masing dengan kurs dan periodenya

---
status: aktif
golongan: perubahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-06` · `SPEC-MODEL-DATA.md` §10.10 · `ADR-D-TI-0053`.*

**What to build:** **PK** mencatat **beberapa** mata uang yang berlaku pada satu versi kontrak,
masing-masing dengan **kursnya dan periode berlakunya**. Mata uang yang sama dua kali pada satu versi
ditolak.

Artefak: entitas `MATA_UANG_KONTRAK`, kunci asingnya ke `VERSI_KONTRAK` dan ke tabel acuan
`MATA_UANG`, dan `INV-07`.

**PEMBUAT PERTAMA** untuk `MATA_UANG_KONTRAK`.

**Kenapa begini:** Sistem lama menyimpan mata uang kontrak sebagai **satu skalar** `Currency` di
kepala. Kontrak yang benar-benar bermata uang lebih dari satu **tidak dapat dinyatakan**, dan yang
terjadi adalah orang memilih satu lalu mengonversi sisanya dengan tangan — konversi yang **kursnya
tidak tercatat di mana pun**. `ADR-D-TI-0053` menjadikan mata uang sebuah **daftar**, dan itu perubahan
bentuk yang tidak dapat dipasang belakangan: skalar yang sudah terisi **tidak menyimpan baris kedua
yang dulu ada**.

**Persyaratan:** `ADR-D-TI-0053` · `INV-07` (kode mata uang unik di dalam satu versi) · `INV-44` (kode
mata uang merujuk tabel acuan) · `INV-38` (kurs selalu lebih besar dari nol; tidak pernah nol, tidak
pernah dipaksa satu) · `INV-57` (tanggal kurs tidak lebih akhir dari tanggal transaksi yang
memakainya) · `INV-43` (kurs dicatat **sebagai nilai pada versi**, tidak dibaca ulang saat dibutuhkan)

**Tidak termasuk:** **Pembekuan kurs pada versi yang sudah disetujui** — `P-27`, batch 2: titik
bekunya **adalah sebuah keadaan**, dan `REV-3` dapat mengubah daftar keadaan. Irisan ini memasang
`INV-43` sebagai bentuk kolom — kursnya **tersimpan**, bukan dibaca ulang — tanpa menyatakan **kapan**
ia membeku.
**Keterangan kegagalan konversi** — irisan `21`.
**`CurrencyRelation`** — `ADR-D-TI-0053` sudah memutuskan ia **pasif** dan tidak disimpan.

**Jalur gagal:** Dua baris bermata uang sama pada satu versi -> **ditolak** `INV-07`, pesannya
menyebut kodenya · Kurs bernilai nol atau kosong -> **ditolak** `INV-38` · Kode mata uang yang tidak
ada di tabel acuan -> ditolak `INV-44` · Tanggal kurs lebih akhir daripada tanggal transaksi ->
ditolak `INV-57`.

**Uji:** **Negatif:** sisipkan mata uang kembar; kurs nol; kurs satu yang dipaksakan; kode tak
terdaftar; tanggal kurs di masa depan.
**Positif — dan ia pokok irisan ini:** satu versi dengan **tiga** mata uang, masing-masing berkurs
dan berperiode berbeda, **diterima seluruhnya**. Sebuah rancangan yang diam-diam masih menyimpan mata
uang sebagai skalar lulus setiap uji negatif di atas — sebab satu baris tidak pernah bentrok dengan
dirinya sendiri.

**Menggantikan:** skalar `TreatyIn.Currency` di kepala kontrak, dan pasangan kolom kembar `Rp` / `Usd`
yang §12.4 nyatakan **dibuang** (`INV-46`). Golongannya **PERUBAHAN** — hasilnya **tidak dapat diuji
terhadap data lama**, sebab data lama hanya pernah memuat satu mata uang per kontrak dan itu justru
bentuk yang dibuang.

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(TreatyIn.Currency@ekspor-2026-09 skalar tunggal; kolom kembar Rp/Usd, SPEC-MODEL-DATA §12.4)
        DECIDED(ADR-D-TI-0053, INV-46, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `MATA_UANG_KONTRAK` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`, dengan kunci asing ke `VERSI_KONTRAK` dan `MATA_UANG`
- [ ] `INV-07` menolak mata uang kembar di dalam satu versi
- [ ] `INV-38` menolak kurs nol **dan** menolak kurs yang dipaksakan menjadi satu — dua uji, bukan satu
- [ ] `INV-57` terpasang: tanggal kurs tidak dapat lebih akhir daripada transaksi yang memakainya
- [ ] uji positif lulus: tiga mata uang pada satu versi diterima seluruhnya
- [ ] sapuan membuktikan **tidak ada** kolom mata uang skalar yang tertinggal di `VERSI_KONTRAK` selain `KODE_MATA_UANG_KONTRAK` yang §10.2 sebutkan
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 21 - Perhitungan yang gagal menghasilkan keterangan apa yang gagal, bukan angka nol dan bukan kurs satu

---
status: aktif
golongan: perubahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-28` · `ADR-D-TI-0035` · `SPEC-INVARIAN.md` `INV-38`, `INV-42`.*

**What to build:** **PK** yang perhitungannya tidak dapat diselesaikan — kurs tidak ada, pembagi nol,
nilai sumber kosong — melihat **keterangan apa yang gagal**. Ia **tidak** melihat angka nol, dan
**tidak** melihat hasil yang diam-diam memakai kurs satu.

**Kenapa begini:** Ini instans `ADR-D-TI-0035` yang paling sering terjadi, dan sistem lama melakukannya di
dua tempat sekaligus. `SaveTreatyInDetail_Act` menulis **`DEDUCTION1 = 0`, `DEDUCTION2 = 0`,
`SPREAD_RNM_SHARE_VALUE = 0`** ke tabel datar pada setiap penyimpanan, dan `DetailCalculationROL`
menulis **`9989998`** sebagai persentase ketika limitnya nol. Keduanya **kegagalan yang menyamar
sebagai nilai**, dan bedanya cuma satu: yang kedua begitu ganjil sampai orang curiga, yang pertama
**sama sekali tidak dapat dibedakan dari nol yang benar**.

> **Nol yang sebenarnya kegagalan tidak dapat dibedakan lagi setelah tersimpan.** Itu sebab `G2`
> menempatkan `ADR-D-TI-0035` di **butir (b)** — perubahan yang menentukan bentuk dan tidak dapat dipasang
> belakangan.

**Persyaratan:** `ADR-D-TI-0035` · `INV-42` (nilai uang tidak pernah memuat angka penanda kegagalan;
ketidakmampuan menghitung menghasilkan **kosong**, bukan angka) · `INV-38` (kurs tidak pernah nol,
tidak pernah dipaksa satu) · `INV-37` (`NILAI_IDR` terisi ⟹ kurs, tanggal kurs, sumber kurs terisi)

**Tidak termasuk:** **Penulisan nol oleh jalur migrasi** — baris warisan yang **sudah** memuat nol
atau `9989998` dibawa **apa adanya** (`ADR-D-TI-0042`), dan pemisahannya dari nol yang benar adalah
pekerjaan irisan migrasi batch 2. Irisan ini melarang sistem baru **membuat** yang baru.
**Berapa baris warisan yang terkena** — **`Uji AQ`**, belum dijalankan.
**Nasib kolom `PERSEN_ROL` itu sendiri** — `F-15`, menunggu satu kalimat pemilik proses.

**Jalur gagal:** Konversi tanpa kurs -> hasilnya **kosong beserta keterangannya**, bukan nol ·
Pembagian dengan nol -> keterangan, bukan angka penanda · `NILAI_IDR` terisi sementara kursnya kosong
-> **ditolak** `INV-37` · Sebuah kolom uang memuat `9989998` -> **ditolak** `INV-42`.

**Uji:** **Negatif:** simpan nilai IDR tanpa kurs; paksakan kurs satu; sisipkan `9989998` ke kolom
persentase; bagi dengan nol lewat jalur perhitungan.
**Positif — dan ia yang menangkap larangan yang terlalu lebar:** sebuah besaran yang **memang bernilai
nol** — potongan nol persen yang disepakati — **tersimpan sebagai nol dan diterima**. Sebuah
constraint yang menolak setiap nol lulus keempat uji negatif di atas dan **menolak data yang sah**.

**Menggantikan:** `SaveTreatyInDetail_Act` yang menulis nol tetap ke tiga kolom datar, dan
`DetailCalculationROL` yang menulis `9989998` — keduanya tercatat di `TETAPAN-DI-KODE.md` §2b dan §3.
Golongannya **PERUBAHAN**: perilakunya sengaja dibalik, sehingga **cocok-tidaknya terhadap data lama
berarti kebalikan** dari yang biasa.

**Blocked by:** 20

**Dasar:**
```
EVIDENCED(SaveTreatyInDetail_Act@ekspor-2026-09 - DEDUCTION1=0, DEDUCTION2=0, SPREAD_RNM_SHARE_VALUE=0)
        EVIDENCED(DetailCalculationROL@ekspor-2026-09 - @if(Local.TotalLimit==0, 9989998, ...))
        DECIDED(ADR-D-TI-0035, INV-42)
```

- [ ] `INV-42` terpasang dan menolak angka penanda kegagalan pada kolom uang dan persentase
- [ ] `INV-38` menolak kurs nol **dan** kurs yang dipaksakan satu
- [ ] keterangan kegagalan **tersimpan dan terbaca**, bukan hanya muncul sesaat di layar
- [ ] uji positif lulus: nol yang memang disepakati tetap tersimpan sebagai nol
- [ ] daftar angka penanda yang dilarang **menyebut `9989998` dengan namanya**, bukan hanya "angka besar"
- [ ] pembatasan lingkup tertulis: baris warisan dibawa apa adanya, dan `Uji AQ` yang menghitungnya

## Treaty In - 22 - Retensi cedant dicatat per kelompok treaty per mata uang

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-07` · `SPEC-MODEL-DATA.md` §10.11 · `SPEC-INVARIAN.md` `INV-08`.*

**What to build:** **PK** mencatat berapa yang **ditahan cedant sendiri**, dirinci **per kelompok
treaty** dan **per mata uang**. Baris kembar pada sumbu itu ditolak.

Artefak: entitas `RETENSI_CEDANT`, kunci asingnya ke `VERSI_KONTRAK` dan `KELOMPOK_TREATY`, kolom
`KODE_MATA_UANG`, dan `INV-08`.

**PEMBUAT PERTAMA** untuk `RETENSI_CEDANT`.

**Kenapa begini:** Kunci alaminya **menyebut mata uang**, dan itu bukan kerapian. Sebuah kontrak dapat
menahan jumlah berbeda dalam dua mata uang untuk kelompok treaty yang sama; menulis `UNIQUE` tanpa
mata uang **mengubah artinya** menjadi *"dilarang dua baris bermata uang berbeda"* dan **menolak data
yang sah**. `INV-08` sempat berstatus **klaim, bukan penegakan**, justru karena kolom mata uangnya
belum ada — `P-8` golongan B yang membebaskannya dengan menaruh mata uang **di baris yang sama**,
sebagaimana bentuknya memang di sistem lama.

**Persyaratan:** `INV-08` (kelompok treaty + mata uang unik di dalam satu versi) · `INV-36` (nilai uang
terisi ⟹ mata uangnya terisi) · `INV-44` · `INV-47` — retensi masuk rekonsiliasi bersama penyerahan
(`INV-51`), dan **`SUMBU_REKONSILIASI`**-nya dinyatakan §10.11

**Tidak termasuk:** **`INV-51`** — *"retensi + penyerahan sama dengan 100 persen"* — bergolongan
**CONSTRAINT BELUM DIBUKTIKAN** (`SPEC-INVARIAN.md` §5.1) dan ditegakkan lewat *materialized view*
bersama `INV-47` dan `INV-50`. Ia berdiri di irisan `38`, **beserta pemantau kebasiannya**.
**Tingkat pencatatan** besaran retensi — `T-2`…`T-4` belum kembali dari teknik treaty.

**Jalur gagal:** Dua baris berkelompok treaty **dan** mata uang sama pada satu versi -> **ditolak**
`INV-08` · Nilai retensi terisi dengan mata uang kosong -> ditolak `INV-36` · Kelompok treaty yang
tidak ada di tabel acuan -> ditolak.

**Uji:** **Negatif:** sisipkan baris kembar pada sumbu penuh; isi nilai tanpa mata uang.
**Positif — dan ia yang menangkap lingkup kunci yang terlalu sempit:** satu versi dengan **dua baris
berkelompok treaty sama tetapi mata uang berbeda** -> **diterima**. Inilah uji yang memisahkan
`UNIQUE` yang benar dari `UNIQUE` tanpa mata uang; keduanya lulus uji negatif di atas.

**Menggantikan:** `P-07` melestarikan daftar `Retention` di sistem lama. Yang bergeser hanya
**kunci alaminya menjadi dapat dikompilasi** — di sistem lama tidak ada constraint apa pun atasnya.

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(TreatyIn.Retention@ekspor-2026-09 - daftar per kelompok treaty; Currency pada baris yang sama)
        DECIDED(INV-08, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `RETENSI_CEDANT` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`, dengan `KODE_MATA_UANG` pada barisnya
- [ ] `INV-08` terpasang sebagai `UNIQUE` **bertiga kolom**, dan mata uangnya ada di dalamnya
- [ ] uji positif lulus: kelompok treaty sama, mata uang berbeda, **diterima**
- [ ] `INV-51` dinyatakan **di luar irisan ini** dengan irisan `38` sebagai pemiliknya
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 23 - EGNPI dicatat per kelompok treaty per mata uang

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-08` · `SPEC-MODEL-DATA.md` §10.12 · `SPEC-INVARIAN.md` `INV-09`.*

**What to build:** **PK** mencatat **perkiraan premi bruto bersih tahunan** — EGNPI — per kelompok
treaty per mata uang. Baris kembar pada sumbu itu ditolak.

Artefak: entitas `EGNPI`, kunci asingnya, kolom `KODE_MATA_UANG`, dan `INV-09`.

**PEMBUAT PERTAMA** untuk `EGNPI`.

**Kenapa begini:** EGNPI adalah **penyebut hampir setiap perhitungan premi** di modul ini —
`DetailCalculation` menggelung `.EgnpiTotalList` dan mengalikannya dengan tarif penyesuaian untuk
menghasilkan premi diperoleh. Besaran yang menjadi penyebut orang lain harus punya **kunci yang tidak
dapat bertabrakan**, dan kuncinya menyebut mata uang dengan alasan yang sama seperti `RETENSI_CEDANT`.

**Persyaratan:** `INV-09` (kelompok treaty + mata uang unik di dalam satu versi) · `INV-36` · `INV-44` ·
`INV-39` dan `INV-40` **tidak dipasang di sini** — lihat `Tidak termasuk`

**Tidak termasuk:** **Tingkat pencatatan `NILAI_EGNPI`.** `INV-39` dan `INV-40` menuntut setiap paket
uang menyatakan apakah ia `TREATY_100_PERSEN` atau `BAGIAN_NURE`; untuk EGNPI **tingkatnya belum
ditentukan**, dan yang menjawabnya **`T-4`** dari teknik treaty. Constraint yang dipasang sekarang
**menuntut nilai yang migrasi tidak tahu cara mengisinya** — itu sebab ketiadaannya berdiri sebagai
**pernyataan keputusan** di `ddl-usulan/`, bukan sebagai pekerjaan tertunda.
**Premi diperoleh yang diturunkan dari EGNPI** — `F-15`, menunggu pemilik proses.

**Jalur gagal:** Dua baris berkelompok treaty **dan** mata uang sama -> **ditolak** `INV-09` · Nilai
EGNPI terisi tanpa mata uang -> ditolak `INV-36`.

**Uji:** **Negatif:** baris kembar pada sumbu penuh; nilai tanpa mata uang.
**Positif:** satu versi dengan dua baris berkelompok treaty sama tetapi **mata uang berbeda** ->
diterima.
**Positif kedua — dan ia khas entitas ini:** sebuah kontrak **tanpa satu pun baris EGNPI** tetap dapat
disimpan dan disetujui. EGNPI adalah perkiraan, bukan syarat; constraint yang menuntutnya ada akan
menolak kontrak yang sah pada hari pertama.

**Menggantikan:** `P-08` melestarikan daftar `EGNPI` sistem lama. Dua kolom yang **tidak dibawa**:
`CurrencyEgnpiAmount` dan `TotalEgnpiAmount` — keduanya agregat turunan yang tidak disimpan
(`SPEC-MODEL-DATA.md` §10.0c, `ADR-D-TI-0037`).

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(TreatyIn.EGNPI@ekspor-2026-09; DetailCalculation@ekspor-2026-09 - gelung .EgnpiTotalList)
        DECIDED(INV-09, ADR-D-TI-0037, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(T-4)
```

- [ ] `EGNPI` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md` dengan `KODE_MATA_UANG` pada barisnya
- [ ] `INV-09` terpasang sebagai `UNIQUE` bertiga kolom
- [ ] kedua uji positif lulus — mata uang berbeda diterima, **dan** kontrak tanpa EGNPI diterima
- [ ] ketiadaan `INV-39`/`INV-40` tertulis sebagai **pernyataan keputusan** yang menyebut `T-4` sebagai penagihnya
- [ ] `TotalEgnpiAmount` dan `CurrencyEgnpiAmount` **tidak ada** sebagai kolom — diperiksa dengan sapuan
- [ ] `KTV-A` dan `T-4` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 24 - Portofolio masuk dan keluar yang menyertai kontrak dicatat

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-09` · `SPEC-MODEL-DATA.md` §10.13 · `SPEC-INVARIAN.md` `INV-66`.*

**What to build:** **PK** mencatat **portofolio masuk** dan **portofolio keluar** yang menyertai
sebuah versi kontrak — premi maupun klaim — dengan arah dan jenisnya terbedakan.

Artefak: entitas `PORTOFOLIO`, kunci asingnya ke `VERSI_KONTRAK`, dan `INV-66`.

**PEMBUAT PERTAMA** untuk `PORTOFOLIO`.

**Kenapa begini:** **Arah adalah pembedanya, bukan nama kolom.** Sistem lama menyimpan portofolio
sebagai daftar `Portfolio` dengan arah sebagai nilai; membawanya sebagai dua kolom bernama — "masuk"
dan "keluar" — akan menuntut kolom ketiga begitu ada arah lain, dan itu bentuk yang sama dengan
delapan kolom bahaya yang `ADR-D-TI-0038` bongkar. Kunci alaminya **arah + jenis** di dalam satu versi,
dan itu baru bernomor pada langkah 3 to-spec (`INV-66`) — sebelumnya entitas ini berdiri **tanpa kunci
alami sama sekali**.

**Persyaratan:** `INV-66` (arah + jenis portofolio unik di dalam satu versi) · `INV-36` · `INV-01`

**Tidak termasuk:** **Perhitungan nilai portofolio** — ia turunan dari premi dan klaim periode
berjalan, dan `ADR-D-TI-0037` melarang turunan disimpan. Yang dicatat **kesepakatannya**, bukan hasilnya.
**Portofolio pada cabang retro keluar** — `RETRO_KELUAR` adalah **GEL-2**, tegas di luar penyerahan
pertama.

**Jalur gagal:** Dua baris berarah **dan** berjenis sama pada satu versi -> **ditolak** `INV-66` ·
Nilai portofolio terisi tanpa mata uang -> ditolak `INV-36`.

**Uji:** **Negatif:** baris kembar pada sumbu arah + jenis.
**Positif — dan ia yang menangkap kunci yang terlalu sempit:** satu versi dengan portofolio **masuk**
dan portofolio **keluar** berjenis sama -> **diterima**. Sebuah `UNIQUE` yang hanya menyebut jenis
lulus uji negatif di atas dan menolak pasangan masuk-keluar yang justru **bentuk paling lazim**.

**Menggantikan:** `P-09` melestarikan daftar `Portfolio`. Yang bergeser: ia memperoleh **kunci alami
bernomor**, yang sebelumnya tidak ada — `SPEC-MODEL-DATA.md` §5.0a menandainya *"kunci alaminya belum
bernomor"* sampai langkah 3 to-spec menutupnya.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(TreatyIn.Portfolio@ekspor-2026-09 - daftar berarah sebagai nilai)
        DECIDED(INV-66, ADR-D-TI-0037, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `PORTOFOLIO` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-66` terpasang atas **arah + jenis** di dalam satu versi
- [ ] uji positif lulus: masuk dan keluar berjenis sama **diterima**
- [ ] arah tersimpan sebagai **nilai**, bukan sebagai dua kolom bernama — diperiksa terhadap `KAMUS-KOLOM.md`
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 25 - Periode pelaporan dicatat beserta tanggal jatuh temponya, dan ditolak bila di luar periode kontrak

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-10` · `SPEC-MODEL-DATA.md` §10.14 · `SPEC-INVARIAN.md` `INV-10`, `INV-55`.*

**What to build:** **PK** mencatat periode pelaporan sebuah versi beserta **batas penyerahan, batas
konfirmasi, batas pelunasan, dan hari pengingatnya**. Periode yang jatuh **di luar periode
kontraknya ditolak**.

Artefak: entitas `PERIODE_PELAPORAN`, kunci asingnya, `INV-10`, dan `INV-55`.

**PEMBUAT PERTAMA** untuk `PERIODE_PELAPORAN`.

**Kenapa begini:** Periode pelaporan yang jatuh di luar periode kontrak **tidak dapat dipenuhi
siapa pun** — tidak ada premi yang dilaporkan untuk waktu yang kontraknya belum atau sudah tidak
berlaku. Sistem lama menyimpan keempat batas hari itu sebagai skalar di kepala (`ReportingSubmission`,
`ReportingConfirmation`, `ReportingSettlement`, `ReminderDays`) **tanpa satu pun pemeriksaan
rentang**, sehingga baris semacam itu memang dapat masuk dan **tidak ada yang tahu berapa banyak**.

**Persyaratan:** `INV-10` (periode unik di dalam satu versi) · **`INV-55`** (periode pelaporan berada
di dalam periode kontraknya) — bergolongan **APLIKASI**, bukan constraint, dan alasannya ada di
`SPEC-INVARIAN.md` §3.5 · `INV-53`

**Tidak termasuk:** **`CARA_PEMBUKUAN` dan `PERIODE_PELAPORAN_KONTRAK` di kepala versi** — keduanya
kolom `VERSI_KONTRAK` yang lahir bersama irisan `14`. Yang di sini **barisnya**, bukan penandanya.
**`CARA_PEMBUKUAN_XOL`** — `P-20`, dan `G2` menyatakannya **tegas tidak masuk** penyerahan pertama:
ia menunggu **Uji X-2**.

**Jalur gagal:** Dua baris berperiode sama pada satu versi -> **ditolak** `INV-10` · Periode yang
mulai sebelum `TANGGAL_MULAI` kontrak atau berakhir sesudah `TANGGAL_BERAKHIR` -> **ditolak**
`INV-55`, pesannya menyebut periode kontraknya · Batas pelunasan lebih awal daripada batas penyerahan
-> ditolak.

**Uji:** **Negatif:** periode kembar; periode yang mulai sehari sebelum kontrak; periode yang berakhir
sehari sesudah kontrak; urutan batas hari yang terbalik.
**Positif — dan ia menguji batas inklusif:** periode yang **mulai tepat pada** `TANGGAL_MULAI` dan
**berakhir tepat pada** `TANGGAL_BERAKHIR` -> **diterima**. `ADR-D-CNP-0022` menetapkan kedua batas
**inklusif**; pemeriksaan yang memakai perbandingan tegas menolak periode yang justru paling lazim.

**Menggantikan:** `P-10` melestarikan daftar `ReportingPeriodList`. Yang bergeser: rentangnya
**diperiksa**, dan di sistem lama tidak.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(TreatyIn.ReportingPeriodList@ekspor-2026-09; ReportingSubmission/Confirmation/Settlement skalar di kepala)
        DECIDED(INV-10, INV-55, ADR-D-CNP-0022, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `PERIODE_PELAPORAN` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-10` terpasang; `INV-55` terpasang **di lapisan aplikasi**, dan letaknya itu tertulis sebagai keputusan
- [ ] pesan penolakan `INV-55` menyebut **periode kontraknya**, bukan hanya "di luar rentang"
- [ ] uji positif lulus: periode yang berimpit tepat pada kedua batas **diterima** — batas inklusif
- [ ] `CARA_PEMBUKUAN_XOL` dinyatakan di luar irisan ini, dengan `P-20` dan Uji X-2 sebagai penagihnya
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 26 - Periode akumulasi dicatat, dan ditolak bila di luar periode kontrak

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-11` · `SPEC-MODEL-DATA.md` §10.15 · `SPEC-INVARIAN.md` `INV-11`, `INV-56`.*

**What to build:** **PK** mencatat periode akumulasi sebuah versi — jendela waktu yang beberapa
kejadian di dalamnya dihitung sebagai **satu kejadian**. Periode yang jatuh di luar periode
kontraknya ditolak.

Artefak: entitas `PERIODE_AKUMULASI`, kunci asingnya, `INV-11`, dan `INV-56`.

**PEMBUAT PERTAMA** untuk `PERIODE_AKUMULASI`.

**Kenapa begini:** Periode akumulasi **menentukan berapa besar klaim yang dibayar**, sebab ia
menentukan apakah dua kejadian dijumlahkan ke satu deductible atau ke dua. Sebuah periode akumulasi
yang jatuh di luar periode kontrak **tidak dapat memuat satu kejadian pun yang tertanggung**, dan
karena itu ia selalu kekeliruan pengisian — bukan pilihan yang sah.

**Persyaratan:** `INV-11` (periode unik di dalam satu versi) · **`INV-56`** (periode akumulasi berada
di dalam periode kontraknya), bergolongan **APLIKASI** · `INV-53` · `ADR-D-CNP-0022` (kedua batas inklusif)

**Tidak termasuk:** **Perhitungan akumulasi kejadian itu sendiri** — ia perilaku klaim, dan modul
klaim tidak di dalam penyerahan ini. Yang dicatat **jendelanya**, bukan hasilnya.
**`PERIODE_AKUMULASI_KONTRAK` di kepala versi** — kolom `VERSI_KONTRAK`, lahir bersama irisan `14`.

**Jalur gagal:** Dua baris berperiode sama pada satu versi -> **ditolak** `INV-11` · Periode di luar
periode kontrak -> **ditolak** `INV-56`, pesannya menyebut periode kontraknya · Periode bertanggal
terbalik -> ditolak `INV-53`.

**Uji:** **Negatif:** periode kembar; periode yang mulai sehari sebelum kontrak; periode bertanggal
terbalik.
**Positif — dan ia yang menangkap kekeliruan lingkup:** satu versi dengan **dua periode akumulasi
yang saling bertumpang tindih** tetapi berperiode berbeda -> **diterima**. Tumpang tindih **tidak
dilarang** oleh satu pun invarian, dan pemeriksaan yang melarangnya akan menolak susunan yang sah.
Ini uji yang memisahkan `INV-11` — keunikan periode — dari larangan tumpang tindih yang **tidak pernah
diputuskan siapa pun**.

**Menggantikan:** `P-11` melestarikan daftar `AccumulationList`. Yang bergeser: rentangnya diperiksa.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(TreatyIn.AccumulationList@ekspor-2026-09; AccumulationPeriod skalar di kepala)
        DECIDED(INV-11, INV-56, ADR-D-CNP-0022, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `PERIODE_AKUMULASI` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-11` terpasang; `INV-56` terpasang di lapisan aplikasi, dan letaknya tertulis sebagai keputusan
- [ ] uji positif lulus: dua periode bertumpang tindih **diterima**
- [ ] ketiadaan larangan tumpang tindih tertulis sebagai **pernyataan keputusan** — supaya orang berikutnya tidak menambahkannya sebagai "kerapian"
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 27 - Termin pembayaran premi dicatat bernomor urut, per mata uang

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-12` · `SPEC-MODEL-DATA.md` §10.16, §10.16a · `SPEC-INVARIAN.md` `INV-12`.*

**What to build:** **PK** mencatat termin pembayaran premi — angsuran bernomor urut beserta tanggal
jatuh temponya dan jumlahnya — **per mata uang**.

Artefak: entitas `TERMIN`, kunci asingnya, kolom `KODE_MATA_UANG`, dan `INV-12`.

**PEMBUAT PERTAMA** untuk `TERMIN`.

**Kenapa begini:** **Kunci alaminya sempat ditulis salah, dan koreksinya adalah alasan irisan ini ada
dalam bentuk ini.** `INV-12` semula berbunyi *"nomor termin unik di dalam satu versi"*. Aktivitas yang
**menambah barisnya** menunjukkan ada satu tingkat pengelompokan di antaranya: daftar yang tampak
"termin" ternyata **baris per mata uang**, dan termin yang sebenarnya hidup satu tingkat di bawahnya —
langkah penambahnya bahkan berlabel jujur, *"Set currency list"*. `UNIQUE` tanpa mata uang **menolak
data yang sah**: kontrak yang menagih empat angsuran dalam dua mata uang punya dua baris bernomor
satu, dan keduanya benar.

> Ini **instans ketiga** dari pola *"lingkup ditulis dari bentuk yang terlihat, bukan dari
> penulisnya"*, dan yang paling bersih. Siapa pun yang menyederhanakan kunci ini mengulanginya.

**Persyaratan:** **`INV-12`** — *"kode mata uang + nomor termin unik di dalam satu versi"*,
**sebagaimana diperbaiki 24 September 2026** · `INV-36` · `INV-44` · `INV-41` (persentase di rentang
0–100, bila terminnya dinyatakan dalam persen)

**Tidak termasuk:** **`JUMLAH_TERMIN` dan `MEMAKAI_PRORATA` di kepala versi** — keduanya kolom
`VERSI_KONTRAK`, lahir bersama irisan `14`.
**Mesin pro rata** — `GRL-15`, **sengaja tidak dibangun**. Pernyataan keputusan, bukan `TODO`.
**Rekonsiliasi jumlah termin terhadap premi** — `INV-47` lewat *materialized view*, irisan `38`.

**Jalur gagal:** Dua baris bernomor termin **dan** mata uang sama -> **ditolak** `INV-12` · Nilai
termin terisi tanpa mata uang -> ditolak `INV-36` · Nomor termin nol atau negatif -> ditolak.

**Uji:** **Negatif:** baris kembar pada sumbu penuh; nilai tanpa mata uang; nomor termin nol.
**Positif — dan ia yang menangkap kunci yang terlalu sempit, sekaligus membuktikan koreksi `INV-12`
memang diterapkan:** satu versi dengan **termin nomor 1 dalam IDR dan termin nomor 1 dalam USD** ->
**diterima**. Uji inilah yang gagal pada rumusan `INV-12` yang lama, dan ia harus dijalankan — bukan
diargumentasikan.

**Menggantikan:** `P-12` melestarikan daftar termin. Yang bergeser **kunci alaminya**, dari *"nomor
termin"* menjadi *"mata uang + nomor termin"* — dan bunyi lamanya disebut di sini justru supaya yang
sudah membacanya tahu apa yang bergeser.

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(aktivitas penambah baris termin@ekspor-2026-09 - langkah "Set currency list"; SPEC-MODEL-DATA §10.16a)
        DECIDED(INV-12, GRL-15, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `TERMIN` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md` dengan `KODE_MATA_UANG` pada barisnya
- [ ] `INV-12` terpasang sebagai `UNIQUE` **bertiga kolom**, dan mata uangnya ada di dalamnya
- [ ] uji positif lulus: termin nomor 1 dalam dua mata uang **diterima** — dijalankan, bukan diargumentasikan
- [ ] ketiadaan mesin pro rata tertulis sebagai **pernyataan keputusan** yang menyebut `GRL-15`
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 28 - Skala ko-asuransi dicatat sebagai beberapa baris, bukan satu nilai

---
status: aktif
golongan: perubahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-13` · `SPEC-MODEL-DATA.md` §10.17 · `SPEC-INVARIAN.md` `INV-13`.*

**What to build:** **PK** mencatat skala ko-asuransi sebagai **beberapa baris** — tiap baris satu
persen limit dengan persen ko-asuransinya — bukan sebagai satu angka di kepala kontrak.

Artefak: entitas `SKALA_KOASURANSI`, kunci asingnya ke `VERSI_KONTRAK`, dan `INV-13`.

**PEMBUAT PERTAMA** untuk `SKALA_KOASURANSI`.

**Kenapa begini:** Skala ko-asuransi **adalah sebuah tangga**: persentase yang berbeda berlaku pada
lapis limit yang berbeda. Sistem lama punya daftarnya — `CoInScale[]` — **dan juga** sebuah skalar
bernama sama di kepala, `TreatyIn.CoInScale`, yang **nol penulis** dan hanya muncul di empat
harness/section serta **tidak terdeklarasi sebagai properti**. Skalar itu **penampung layar**, bukan
atribut, dan `COCOK-SILANG-CACAH-ATRIBUT.md` §4 mengadilinya demikian. Membawanya sebagai satu nilai
akan meratakan tangga menjadi satu anak tangga.

**Persyaratan:** `INV-13` (persen limit unik di dalam satu versi) · `INV-41` (persentase di rentang
0–100) · `INV-01`

**Tidak termasuk:** **Skalar `TreatyIn.CoInScale`** — **tidak dibawa**, dan itu putusan bernomor:
nol penulis atas kelima bentuk penulis, tak terdeklarasi, hanya di layar. Ini **pernyataan
keputusan**; siapa pun yang menambahkannya kembali membalikkan adjudikasi `F-2`.
**Perhitungan bagian ko-asuransi** — turunan, `ADR-D-TI-0037`.

**Jalur gagal:** Dua baris berpersen limit sama pada satu versi -> **ditolak** `INV-13` · Persen di
luar 0–100 -> ditolak `INV-41` · Sebuah kolom skalar ko-asuransi di `VERSI_KONTRAK` -> **tidak ada**,
dan ketiadaannya diperiksa dengan sapuan.

**Uji:** **Negatif:** persen limit kembar; persen 101; persen negatif.
**Positif — dan ia pokok irisan ini:** satu versi dengan **empat baris skala** berpersen limit menaik
-> **diterima seluruhnya**, dan keempatnya terbaca kembali **dalam urutan yang sama**. Sebuah
rancangan yang diam-diam masih meratakannya menjadi satu nilai lulus ketiga uji negatif di atas —
satu baris tidak pernah bentrok dengan dirinya sendiri.

**Menggantikan:** skalar `TreatyIn.CoInScale` di kepala. Golongannya **PERUBAHAN**: bentuknya sengaja
diubah dari satu nilai menjadi daftar, sehingga **cocok-tidaknya terhadap data lama berarti
kebalikan** dari yang biasa.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(datar-treatyin-lama.csv baris 172 - TreatyIn.CoInScale Skalar, "(tidak dideklarasikan)", 0 penulis)
        EVIDENCED(COCOK-SILANG-CACAH-ATRIBUT.md §4 - adjudikasi kelima calon F-2)
        DECIDED(INV-13, ADR-D-TI-0037, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `SKALA_KOASURANSI` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-13` terpasang atas persen limit di dalam satu versi
- [ ] `INV-41` terpasang dan diuji pada **kedua** ujung rentang
- [ ] uji positif lulus: empat baris skala diterima dan terbaca kembali berurutan
- [ ] sapuan membuktikan **tidak ada** kolom skalar ko-asuransi di `VERSI_KONTRAK`
- [ ] ketidakbawaan skalar `CoInScale` tertulis sebagai **pernyataan keputusan** yang menyebut buktinya
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 29 - Batas tanggungan dicatat untuk bahaya apa pun yang ada di daftar bahaya, termasuk yang ditambahkan kemudian

---
status: aktif
golongan: perubahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-14` · `SPEC-MODEL-DATA.md` §10.0b, §10.18 · `ADR-D-TI-0038` · `SPEC-INVARIAN.md` `INV-14`.*

**What to build:** **PK** mencatat batas tanggungan **per bahaya**, untuk **bahaya apa pun yang ada di
tabel acuan** — termasuk bahaya yang ditambahkan sesudah sistem berjalan, **tanpa perubahan skema**.

Artefak: entitas `BATAS_PER_BAHAYA`, kunci asingnya ke `VERSI_KONTRAK` dan `BAHAYA`, kolom
`KODE_MATA_UANG`, dan `INV-14`.

**PEMBUAT PERTAMA** untuk `BATAS_PER_BAHAYA`.

**Kenapa begini:** Sistem lama memakai **delapan kolom untuk empat bahaya** yang **namanya menjadi
nama kolom** — `Earthquake`, `FloodJab`, `FloodNation`, `RSMDLimit`, masing-masing berpasangan dengan
kolom mata uangnya. Bahaya kelima menuntut kolom kesembilan dan kesepuluh, dan bahaya baru adalah hal
yang pasti terjadi di reasuransi. **Bahayanya adalah nilai, bukan nama kolom** — dan itu `ADR-D-TI-0038`.

**Persyaratan:** `INV-14` (bahaya unik di dalam satu versi) · `INV-36` · `INV-44` · `ADR-D-TI-0038` ·
`INV-46` (tidak ada pasangan kolom kembar untuk dua mata uang)

**Tidak termasuk:** **Tingkat pencatatan `NILAI_BATAS`.** `INV-39` dan `INV-40` **tidak dipasang**;
yang menjawabnya **`T-2`** dari teknik treaty. Ketiadaannya berdiri sebagai **pernyataan keputusan**
di `ddl-usulan/`, sebab constraint yang dipasang sekarang menuntut nilai yang migrasi tidak tahu cara
mengisinya.
**Glosarium `RSMD`** — istilah pasar yang **dipertahankan** dan **wajib punya entri glosarium**
(`CONTEXT.md` §3.3b). Kepanjangannya **belum tercatat di mana pun**; korpusnya sudah disapu habis dan
hasilnya nol. Itu butir wawancara, bukan pekerjaan irisan ini.

**Jalur gagal:** Dua baris berbahaya sama pada satu versi -> **ditolak** `INV-14` · Nilai batas terisi
tanpa mata uang -> ditolak `INV-36` · Bahaya yang tidak ada di tabel acuan -> ditolak · Sebuah kolom
bernama bahaya di `VERSI_KONTRAK` -> **tidak ada**, diperiksa dengan sapuan.

**Uji:** **Negatif:** baris berbahaya kembar; nilai tanpa mata uang; bahaya tak terdaftar.
**Positif — dan ia pokok irisan ini:** tambahkan **bahaya kesembilan** ke tabel acuan lewat irisan
`15`, lalu catat batas tanggungan untuknya pada sebuah kontrak. **Berhasil tanpa satu pun perubahan
skema.** Inilah satu-satunya uji yang memisahkan tabel acuan yang sungguhan dari daftar nilai yang
diam-diam masih tertanam di tempat lain.

**Menggantikan:** delapan kolom bernama bahaya di kepala `TreatyIn`, beserta keempat kolom
`Currency*`-nya yang **melebur ke paket uangnya** (`COCOK-SILANG-CACAH-ATRIBUT.md` §3.3).
Golongannya **PERUBAHAN** — bentuknya sengaja diubah, jadi data lama menunjukkan bentuk yang dibuang.

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(SPEC-MODEL-DATA.md §10.0b - Earthquake/FloodJab/FloodNation/RSMDLimit + empat kolom Currency*)
        DECIDED(ADR-D-TI-0038, INV-14, INV-46, KTV-A, KTV-C)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(T-2)
```

- [ ] `BATAS_PER_BAHAYA` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`, dengan `KODE_MATA_UANG` **boleh kosong** — ia bukan bagian kunci alami (`KTV-C` koreksi 1)
- [ ] `INV-14` terpasang atas bahaya di dalam satu versi
- [ ] uji positif lulus: **bahaya kesembilan** dapat dipakai tanpa perubahan skema
- [ ] sapuan membuktikan **nol** kolom bernama bahaya di `VERSI_KONTRAK`
- [ ] ketiadaan `INV-39`/`INV-40` tertulis sebagai **pernyataan keputusan** yang menyebut `T-2`
- [ ] `RSMD` **terdaftar sebagai butir glosarium yang belum terisi**, bukan didiamkan
- [ ] `KTV-A` dan `T-2` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 30 - Dokumen dilampirkan pada kontrak, dan dokumennya dirujuk bukan disalin

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-15` · `SPEC-MODEL-DATA.md` §10.19 · `ADR-D-CNP-0027` · `SPEC-INVARIAN.md` `INV-59`, `INV-67`.*

**What to build:** **PK** melampirkan dokumen pada sebuah versi kontrak — slip, wording, surat — dan
yang tersimpan adalah **rujukan ke dokumennya**, bukan salinan isinya.

Artefak: entitas `DOKUMEN_KONTRAK`, kunci asingnya ke `VERSI_KONTRAK`, dan `INV-67`.

**PEMBUAT PERTAMA** untuk `DOKUMEN_KONTRAK`.

**Kenapa begini:** Menyalin isi berkas ke dalam basis data melahirkan **dua sumber kebenaran untuk
satu dokumen**, dan yang kedua tidak pernah ikut berubah ketika yang pertama diganti — persis yang
`INV-59` larang. `ADR-D-CNP-0027` sudah memutuskan lampiran **dirujuk, tidak dimiliki**; sistem lama pun
demikian, lewat `M_ATTACHMENTTREATY_2`. Yang berubah hanya bahwa aturannya kini **tertulis dan
diperiksa**.

**Persyaratan:** `ADR-D-CNP-0027` · `INV-59` (tidak ada nilai yang disalin dari entitas lain; hilir diberi
rujukan) · `INV-67` (dokumen unik di dalam satu versi) · `INV-61` **tidak berlaku di sini** — ia
tentang arsip JSON, irisan `42`

**Tidak termasuk:** **`DOKUMEN_ADDENDUM`** — entitas yang **berbeda**, milik irisan `04`, dan
namanya sengaja dibedakan: `DOKUMEN_KONTRAK` berstatus *"dirujuk, tidak dimiliki"* untuk lampiran,
sementara `DOKUMEN_ADDENDUM` adalah **benda yang berdiri sendiri di atas versi**, dengan nomor yang
beredar di luar sistem. Menyatukannya membuat dua arti pada satu nama.
**Penyimpanan berkasnya** — di luar skema; yang di sini rujukannya.
**Arsip JSON sistem lama** — irisan `42`, dan ia benda lain lagi.

**Jalur gagal:** Dua baris merujuk dokumen yang sama pada satu versi -> **ditolak** `INV-67` ·
Sebuah kolom yang memuat isi berkas -> **tidak ada**, dan ketiadaannya diperiksa dengan sapuan atas
`KAMUS-KOLOM.md` · Nama dokumen disalin dari sumbernya sebagai teks yang kemudian basi -> dilarang
`INV-59`.

**Uji:** **Negatif:** lampirkan dokumen yang sama dua kali pada satu versi.
**Positif — dan ia yang menangkap kunci yang terlalu ketat:** **dokumen yang sama dilampirkan pada
dua versi berbeda** dari kontrak yang sama -> **diterima**. Wording yang tidak berubah antar versi
adalah keadaan yang lazim; `UNIQUE` yang ditulis atas dokumen saja — tanpa versinya — menolaknya, dan
lulus uji negatif di atas.

**Menggantikan:** `P-15` melestarikan `M_ATTACHMENTTREATY_2`. Yang bergeser: keunikannya
**diperiksa**, dan rujukan-bukan-salinan **dinyatakan sebagai invarian** alih-alih dibiarkan sebagai
kebiasaan.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(M_ATTACHMENTTREATY_2@Table/ - lampiran sebagai rujukan)
        DECIDED(ADR-D-CNP-0027, INV-59, INV-67, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `DOKUMEN_KONTRAK` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-67` terpasang atas **dokumen di dalam satu versi** — bukan atas dokumen saja
- [ ] uji positif lulus: satu dokumen melekat pada dua versi berbeda **diterima**
- [ ] sapuan membuktikan **tidak ada** kolom pembawa isi berkas
- [ ] perbedaannya dari `DOKUMEN_ADDENDUM` **tertulis di dalam tiket ini**, supaya dua nama tidak berarti satu benda
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 31 - Layer dicatat beserta limit, deductible, MDP, dan ketentuan pemulihan limitnya

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-16` · `SPEC-MODEL-DATA.md` §10.3, §10.3a, §10.3b · `SPEC-INVARIAN.md` `INV-05`.*

**What to build:** **PK** mencatat layer sebuah kontrak non-proporsional — nomor layer dan bagiannya,
limit, deductible, tarif penyesuaian, minimum deposit premium, dan **ketentuan pemulihan limitnya** —
lalu membacanya kembali utuh.

Artefak: entitas `LAYER`, `PEMULIHAN_LIMIT`, dan kedua tabel anak paket uang `NILAI_MDP` dan
`NILAI_MDP_MINIMUM`; kunci asingnya; `INV-05`.

**PEMBUAT PERTAMA** untuk `LAYER`, `PEMULIHAN_LIMIT`, `NILAI_MDP`, `NILAI_MDP_MINIMUM`.

**Kenapa begini:** `LAYER` **diselesaikan atas sumber yang kurang**, dan itu tercatat: pohon 985
simpul tidak pernah melihat **enam belas properti** kelas `Data-TreatyInLimits` (`L-8`). Dari yang
kemudian diadili lahir `LIMIT_AGREGAT` — batas total sepanjang periode, terpisah dari limit per
kejadian — yang **hilang tanpa disengaja** karena kembaran mata uang keduanya dibuang lebih dulu
sementara kembaran pertamanya tidak pernah ikut masuk. **Menghilangkannya mengubah arti kontrak.**
Dan `MDP` bukan satu angka: ia **daftar per mata uang** di sistem lama, sehingga ia tabel anak —
golongan A menurut `P-8`.

**Persyaratan:** `INV-05` (`NOMOR_LAYER` + `BAGIAN_LAYER` unik di dalam satu versi) · `INV-36` ·
`INV-41` · `INV-46` (tidak ada pasangan kolom kembar untuk dua mata uang) · `INV-49` (pemulihan
adalah **besaran berulang**, dikecualikan dari partisi `INV-47`) · `ADR-D-TI-0035` — `DEDUCTIBLE_KEDUA`
tersimpan **sebagai teks** di sistem lama, dan baris yang gagal dikonversi **tidak disamarkan menjadi
nol**

**Tidak termasuk:** **Syarat yang berbeda antar pemulihan** — irisan `32`. Di sini pemulihan
tersimpan sebagai daftar; **nilainya boleh berbeda** adalah kemampuan tersendiri dan bergolongan
**BARU**.
**Bagian NuRe atas layer** — irisan `33`.
**Nasib `PERSEN_ROL` dan premi diperoleh** — **`F-15`**. `PERSEN_ROL` berdiri di §10.3 atas putusan
`TDA-17`, dan `F-15` membantahnya dengan rumus `DetailCalculationROL`: `premi ÷ limit × 100`. Tabel
anak `NILAI_PREMI_DIPEROLEH` **tidak dibuat**. Keputusannya pemilik proses; irisan ini **tidak
mendahuluinya ke arah mana pun**.
**Tingkat pencatatan `LIMIT_AGREGAT`** — `T-3`.

**Jalur gagal:** Dua layer bernomor **dan** berbagian sama pada satu versi -> **ditolak** `INV-05` ·
`DEDUCTIBLE_KEDUA` warisan yang tidak dapat dikonversi dari teks -> **kegagalan bernama**, bukan nol
(`ADR-D-TI-0035`) · Baris MDP tanpa mata uang -> ditolak · `PERSEN_ROL` bernilai `9989998` -> ditolak
`INV-42`.

**Uji:** **Negatif:** layer kembar pada sumbu penuh; MDP tanpa mata uang; sisipkan `9989998`;
konversi `DEDUCTIBLE_KEDUA` dari teks yang bukan angka.
**Positif — dan ia yang menangkap kunci yang terlalu sempit:** satu versi dengan **layer 1 bagian A
dan layer 1 bagian B** -> **diterima**. `UNIQUE` atas nomor layer saja lulus uji negatif di atas dan
menolak susunan berlapis yang justru pokok kontrak non-proporsional.
**Positif kedua:** satu layer dengan **dua baris MDP bermata uang berbeda** -> diterima.
**Positif ketiga:** satu layer dengan **dua baris pemulihan limit** -> diterima, dan `INV-47`
**tidak** menuntut keduanya berjumlah utuh — itu yang `INV-49` kecualikan.

**Menggantikan:** `P-16` melestarikan daftar `Limits[]`. Yang bergeser: `LIMIT_AGREGAT`,
`MDP_MINIMUM`, `PERSEN_MDP_MINIMUM`, `MDP_DIGABUNG`, dan `TANPA_HITUNG_PREMI_PEMULIHAN`
**ditambahkan** — kelimanya ada di sistem lama dan **tidak pernah terlihat** oleh pohon (`L-8`,
§10.3a). Dan `PEMULIHAN_LIMIT` berdiri sebagai entitas tersendiri (§10.23c butir 1).

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(SPEC-MODEL-DATA.md §10.3a - 16 properti Data-TreatyInLimits, diadili satu per satu)
        EVIDENCED(SetReinstatementPct@ekspor-2026-09 - ReinstatementPct="100", AdditionalPct="100" sebagai tetapan)
        DECIDED(INV-05, INV-49, ADR-D-TI-0035, KTV-A, KTV-C)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(T-6)
        DIASUMSIKAN-CLEAR(T-3)
        DIASUMSIKAN-CLEAR(F-15)
```

- [ ] `LAYER`, `PEMULIHAN_LIMIT`, `NILAI_MDP`, `NILAI_MDP_MINIMUM` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-05` terpasang atas `NOMOR_LAYER` + `BAGIAN_LAYER` di dalam satu versi
- [ ] ketiga uji positif lulus — layer berbagian, MDP dua mata uang, dua baris pemulihan
- [ ] konversi `DEDUCTIBLE_KEDUA` dari teks **melaporkan kegagalannya**; tidak ada baris yang diam-diam menjadi nol
- [ ] `INV-49` tertulis sebagai pengecualian bernama terhadap `INV-47` — bukan dibiarkan terbaca sebagai pelanggaran
- [ ] penanda `F-15` pada `PERSEN_ROL` **masih berbunyi** di `SPEC-MODEL-DATA.md` §10.3 saat tiket ini dinyatakan selesai, atau keputusannya sudah turun dan tiket ini disesuaikan
- [ ] `KTV-A`, `T-6`, `T-3`, `F-15` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 32 - Syarat berbeda ditetapkan untuk tiap pemulihan limit pada sebuah layer

---
status: aktif
golongan: baru
---



*Asal: `DAFTAR-PEKERJAAN.md` §2.2b `P-59` · `SPEC-MODEL-DATA.md` §10.3b · `TETAPAN-DI-KODE.md` §3.*

**What to build:** **PK** menetapkan **syarat yang berbeda untuk tiap pemulihan limit** pada sebuah
layer — pemulihan pertama, kedua, dan seterusnya boleh berpersentase berbeda, baik porsi limit yang
dipulihkan maupun tarif premi pemulihannya.

**Kenapa begini:** Sistem lama **tidak dapat menyatakannya.** `SetReinstatementPct.xml` menulis
`ReinstatementPct = "100"` dan `AdditionalPct = "100"` sebagai **tetapan di dalam kode**, bukan
sebagai masukan — sehingga setiap pemulihan selalu seragam, dan ketentuan pasar yang biasa
*("pemulihan pertama gratis, kedua berbayar 100%")* **tidak punya tempat**. Yang **PELESTARIAN**
hanyalah menyimpan pemulihan sebagai daftar; yang **BARU** adalah **nilainya boleh berbeda**.

> Ini **instans keempat** pola *"cacat bersembunyi di balik nilai bawaan"*: selama semuanya 100, satu
> baris dan sepuluh baris menghasilkan angka yang sama — dan tidak ada yang pernah melihat batasnya.

**Persyaratan:** `INV-49` (pemulihan adalah **besaran berulang**, dikecualikan dari partisi `INV-47`) ·
`INV-41` (persentase di rentang 0–100) · `SPEC-MODEL-DATA.md` §10.3b

**Tidak termasuk:** **Entitas `PEMULIHAN_LIMIT` itu sendiri** — irisan `31` yang membuatnya.
**Premi pemulihan yang dihitung darinya** — turunan, `ADR-D-TI-0037`; dan `TANPA_HITUNG_PREMI_PEMULIHAN`
sudah menjadi kolom `LAYER` di irisan `31`.
**Baris selisih untuk perubahan persentase pemulihan** — tiket `13`, jalur addendum.

**Jalur gagal:** Persentase pemulihan di luar 0–100 -> **ditolak** `INV-41` · Dua baris pemulihan
bernomor urut sama pada satu layer -> ditolak · `INV-47` menuntut baris pemulihan berjumlah utuh ->
**tidak terjadi**; `INV-49` mengecualikannya, dan bila ia tetap menolak maka pengecualiannya belum
terpasang.

**Uji:** **Negatif:** persentase 101; persentase negatif; nomor urut pemulihan kembar.
**Positif — dan ia pokok irisan ini:** satu layer dengan **pemulihan pertama berporsi 100% bertarif
0%** dan **pemulihan kedua berporsi 100% bertarif 100%** -> **diterima**, dan keduanya terbaca kembali
berbeda. Data uji **wajib memuat nilai selain 100**; uji yang seluruhnya bernilai 100 **tidak
memisahkan apa pun** — ia menghasilkan hasil yang sama pada rancangan lama maupun baru.
**Positif kedua:** `INV-47` **tidak** menolak layer berpemulihan dua baris.

**Menggantikan:** tetapan `ReinstatementPct = "100"` dan `AdditionalPct = "100"` di
`SetReinstatementPct.xml`. **Golongannya BARU**, dan karena itu `CARA MENYALAKANNYA` wajib.

**Satu pertanyaan wawancara menyertainya, dan ia tidak menahan pembangunannya:** apakah keseragaman
100/100 di sistem lama **kehendak bisnis** atau **akibat cara sistem dibangun**? **Datanya tidak
dapat menjawab** — seluruhnya disemai 100, sehingga *"selalu seragam"* dan *"tidak pernah bisa
berbeda"* menghasilkan data yang persis sama.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **PK yang menyusun layer non-proporsional**, pada setiap kontrak berpemulihan. **Bukan pemantau berkala**: kegagalannya muncul saat seseorang mencoba mengisi nilai berbeda |
| 3 | **ambang berangka** — nyala penuh ketika **satu kontrak** tercatat dengan dua pemulihan berpersentase berbeda **dan lolos persetujuan**. Sebelum itu ia tersedia tetapi bernilai bawaan seragam, sehingga **tidak mengubah apa pun** |
| 4 | **siapa boleh menyalakan** — pemilik proses, setelah bisnis mengonfirmasi bahwa pemulihan bertingkat memang dipakai NuRe |

**Blocked by:** 31

**Dasar:**
```
EVIDENCED(SetReinstatementPct@ekspor-2026-09 - ReinstatementPct="100" dan AdditionalPct="100" sebagai tetapan)
        DECIDED(INV-49, SPEC-MODEL-DATA §10.3b)
```

- [ ] porsi dan tarif tiap baris pemulihan dapat **berbeda satu sama lain**, dan terbaca kembali berbeda
- [ ] `INV-41` terpasang dan diuji pada **kedua** ujung rentang
- [ ] uji positif dijalankan dengan **nilai selain 100**; data uji seluruhnya-100 **tidak diterima sebagai bukti**
- [ ] `INV-47` tidak menolak layer berpemulihan dua baris — `INV-49` terpasang sebagai pengecualian bernama
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk ambang berangka
- [ ] pertanyaan wawancara *"kehendak bisnis atau akibat cara sistem dibangun"* **terdaftar**, beserta catatan bahwa datanya tidak dapat menjawabnya

## Treaty In - 33 - Bagian NuRe atas sebuah layer dicatat, dan bagian hanya dapat dibuat pada kontrak non-proporsional

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-17` · `SPEC-MODEL-DATA.md` §10.5, §12.3 · `SPEC-INVARIAN.md` `INV-33`, `INV-64`.*

**What to build:** **PK** mencatat **bagian NuRe** atas sebuah layer — premi bruto dan minimumnya per
mata uang — dan sistem **menolak** pembuatan bagian pada kontrak yang sifat proporsinya
`PROPORSIONAL`.

Artefak: entitas `BAGIAN` dan kedua tabel anak paket uang `NILAI_PREMI_BRUTO` dan
`NILAI_PREMI_BRUTO_MINIMUM`; `INV-64`; `INV-33`.

**PEMBUAT PERTAMA** untuk `BAGIAN`, `NILAI_PREMI_BRUTO`, `NILAI_PREMI_BRUTO_MINIMUM`.

**Kenapa begini:** §12.3 **mengoreksi** §8: `BAGIAN` **bukan** satu entitas untuk kedua cabang. Di
cabang non-proporsional ia kelas tersendiri dengan 49 properti dan baris `Share[]`-nya sendiri; di
cabang proporsional **tidak ada kelas tersendiri sama sekali** — bagian NuRe di sana adalah **atribut
pada baris `DETAIL_PROPORSIONAL`**. Uji yang dipakai semula — *"berbeda arti versus tidak dipakai"* —
dijalankan pada sumbu yang salah; sumbu yang menentukan **masukan versus turunan**, dan pada sumbu itu
kedua sisi tidak sebanding.

**Persyaratan:** `INV-64` (satu `BAGIAN` per `LAYER` — `ID_LAYER` unik di `BAGIAN`) · **`INV-33`**
(`BAGIAN` hanya ada di bawah kontrak ber-`SIFAT_PROPORSI` `NON_PROPORSIONAL`), bergolongan **INDEKS
UNIK** · `INV-36` · `INV-52` (premi bruto dikurangi seluruh potongan sama dengan premi bersih) —
**dibawa sebagai invarian**, penegakannya di irisan `37`

**Tidak termasuk:** **Potongan atas premi** — irisan `37`; `POTONGAN` menggantung pada `BAGIAN`
**dan** `DETAIL_PROPORSIONAL`, sehingga ia menunggu keduanya.
**Penyebaran bagian NuRe ke susunan retro** — irisan `38`.
**Sepuluh daftar `RNMSpreadedList…XOL`** — **turunan, tidak disimpan** (`ADR-D-TI-0037`, §12.3). Itu
pernyataan keputusan: siapa pun yang menyimpannya menanam kolom turunan yang `INV-58` larang.

**Jalur gagal:** Dua bagian pada satu layer -> **ditolak** `INV-64` · Bagian pada kontrak
`PROPORSIONAL` -> **ditolak** `INV-33`, dan pesannya menyebut sifat proporsi kontraknya · Premi bruto
terisi tanpa mata uang -> ditolak `INV-36`.

**Uji:** **Negatif:** dua bagian pada satu layer; bagian di bawah kontrak proporsional; premi bruto
tanpa mata uang.
**Positif:** satu bagian dengan **dua baris premi bruto bermata uang berbeda** -> diterima.
**Positif kedua — dan ia yang menangkap `INV-33` yang salah arah:** sebuah `DETAIL_PROPORSIONAL` pada
kontrak **proporsional** -> **diterima**. Indeks unik yang ditulis terbalik menolak cabang yang justru
sah, dan ia lulus setiap uji negatif di atas.

**Menggantikan:** `P-17` melestarikan baris `Share[]` cabang non-proporsional. Yang bergeser:
`BAGIAN` **tidak lagi mengaku melayani kedua cabang** — §12.3 mengoreksi §8 — dan premi brutonya
menjadi **tabel anak per mata uang** (`P-8` golongan A) alih-alih satu angka.

**Blocked by:** 31

**Dasar:**
```
EVIDENCED(Data-TreatyInShare@ekspor-2026-09 - 49 properti, baris Share[]; cabang proporsional tanpa kelas tersendiri)
        DECIDED(INV-33, INV-64, ADR-D-TI-0037, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `BAGIAN`, `NILAI_PREMI_BRUTO`, `NILAI_PREMI_BRUTO_MINIMUM` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-64` terpasang; `INV-33` terpasang **sebagai indeks unik**, dan bentuk itu tertulis sebagai keputusan
- [ ] pesan penolakan `INV-33` menyebut **sifat proporsi kontraknya**
- [ ] kedua uji positif lulus — dua mata uang diterima, **dan** cabang proporsional tidak ikut tertolak
- [ ] sapuan membuktikan **nol** kolom `RNMSpreadedList…` tersimpan
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 34 - Ketentuan proporsional dicatat per kelompok treaty di dalam sebuah layer

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-21` · `SPEC-MODEL-DATA.md` §10.4, §12.3 · `SPEC-INVARIAN.md` `INV-06`, `INV-32`.*

**What to build:** **PK** mencatat ketentuan cabang proporsional — bagian NuRe, jenis treaty, dan
cadangan premi per mata uang — **per kelompok treaty di dalam sebuah layer**. Baris kembar pada sumbu
itu ditolak, dan cabang ini **hanya ada di bawah kontrak proporsional**.

Artefak: entitas `DETAIL_PROPORSIONAL` dan tabel anak `NILAI_CADANGAN_PREMI`; `INV-06`; `INV-32`.

**PEMBUAT PERTAMA** untuk `DETAIL_PROPORSIONAL`, `NILAI_CADANGAN_PREMI`.

**Kenapa begini:** **Bagian NuRe pada cabang proporsional tidak punya baris sendiri** — ia atribut
`RNMShare` **pada** baris ini, terisi bila `RNMShareAcrossTheBoard` mati. §12.3 menetapkannya sesudah
uji semula dijalankan pada sumbu yang salah. Dan lingkupnya **di dalam layer**, bukan di dalam versi:
`INV-06` berbunyi *"`ID_KELOMPOK_TREATY` unik di dalam satu `LAYER`"* — menulisnya *"di dalam versi"*
akan menolak dua layer yang sama-sama memuat kelompok treaty yang sama, dan itu susunan yang lazim.

**Persyaratan:** `INV-06` (kelompok treaty unik di dalam satu `LAYER`) · **`INV-32`**
(`DETAIL_PROPORSIONAL` hanya di bawah kontrak ber-`SIFAT_PROPORSI` `PROPORSIONAL`), bergolongan
**INDEKS UNIK** · `INV-30` (`JENIS_TREATY` dua nilai) · `INV-36` · `INV-41`

**Tidak termasuk:** **Aturan tepat-satu antara persen quota share dan jumlah lines surplus** — irisan
`35`. Kolomnya berdiri di sini; **aturannya** irisan berikutnya.
**Baris surplus yang menuntut adanya baris quota share** — irisan `36`.
**`KAPASITAS_SURPLUS` sebagai kolom tersimpan** — **tidak ada**: ia **kelipatan** retensi × jumlah
lines, dan `INV-48` menamainya pengecualian bernama terhadap `INV-47`. Turunan tidak disimpan
(`ADR-D-TI-0037`).
**Potongan dan penyebaran** — irisan `37` dan `38`.

**Jalur gagal:** Dua baris berkelompok treaty sama pada satu layer -> **ditolak** `INV-06` ·
`DETAIL_PROPORSIONAL` di bawah kontrak non-proporsional -> **ditolak** `INV-32` · Cadangan premi
terisi tanpa mata uang -> ditolak `INV-36`.

**Uji:** **Negatif:** baris kembar di dalam satu layer; baris di bawah kontrak non-proporsional;
cadangan premi tanpa mata uang.
**Positif — dan ia yang menangkap lingkup yang terlalu lebar:** **dua layer pada versi yang sama,
masing-masing memuat kelompok treaty yang sama** -> **diterima**. `UNIQUE` yang ditulis di dalam
versi alih-alih di dalam layer menolaknya, dan lulus uji negatif di atas. Ini bentuk yang sama dengan
ketiga kekeliruan lingkup yang sudah tercatat di modul ini.
**Positif kedua:** `BAGIAN` pada kontrak **non-proporsional** tetap diterima — `INV-32` dan `INV-33`
tidak saling menelan.

**Menggantikan:** `P-21` melestarikan `Data-TreatyInLimitsDetail`. Yang bergeser: **cadangan premi
menjadi tabel anak per mata uang** (`P-8` golongan A), dan bagian NuRe proporsional **dinyatakan
sebagai atribut**, bukan entitas — koreksi §12.3 terhadap §8.

**Blocked by:** 31 · 15

**Dasar:**
```
EVIDENCED(Data-TreatyInLimitsDetail@ekspor-2026-09 - RNMShare per baris bila RNMShareAcrossTheBoard mati)
        DECIDED(INV-06, INV-32, INV-48, ADR-D-TI-0037, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `DETAIL_PROPORSIONAL` dan `NILAI_CADANGAN_PREMI` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-06` terpasang **di dalam `LAYER`**, bukan di dalam versi — diperiksa dengan uji positif, bukan dibaca dari rumusannya
- [ ] `INV-32` terpasang sebagai indeks unik, dan tidak menelan `INV-33`
- [ ] kedua uji positif lulus
- [ ] sapuan membuktikan **tidak ada** kolom `KAPASITAS_SURPLUS` tersimpan; `INV-48` tertulis sebagai pengecualian bernama
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 35 - Tepat satu dari persen quota share atau jumlah lines surplus terisi, ditentukan jenis treaty-nya

---
status: aktif
golongan: baru
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-22` · `SPEC-MODEL-DATA.md` §10.4 · `SPEC-INVARIAN.md` `INV-30`, `INV-31`.*

**What to build:** Pada tiap baris ketentuan proporsional, **tepat satu** dari `PERSEN_QUOTA_SHARE`
dan `JUMLAH_LINES_SURPLUS` terisi, dan **yang mana** ditentukan `JENIS_TREATY` baris itu. Mengisi
keduanya ditolak; mengisi yang tidak sesuai jenisnya ditolak; mengosongkan keduanya ditolak.

**Kenapa begini:** Kedua kolom itu **mengukur hal yang berbeda dengan satuan yang berbeda** —
persentase terhadap jumlah lines — dan keduanya terisi berarti barisnya **tidak dapat ditafsirkan**:
tidak ada yang tahu mana yang dipakai menghitung. Sistem lama membiarkan keduanya terisi, sebab
pemeriksaannya ada di layar saja, dan layar dapat dilewati lewat jalur simpan yang lain. Bentuk
constraint-nya **lintas baris pada satu baris** — `INV-31` menyebutnya demikian — sehingga ia dapat
ditegakkan basis data, bukan aplikasi.

**Persyaratan:** **`INV-31`** (*tepat satu* dari keduanya terisi, ditentukan `JENIS_TREATY`) ·
`INV-30` (`JENIS_TREATY` hanya `QUOTA_SHARE` atau `SURPLUS`) · `INV-41` (persentase di rentang 0–100)

**Tidak termasuk:** **Entitas `DETAIL_PROPORSIONAL`** — irisan `34` yang membuatnya beserta kedua
kolomnya.
**Baris surplus yang menuntut baris quota share** — irisan `36`. Bedanya nyata: yang di sini **satu
baris terhadap dirinya sendiri**, yang di sana **satu baris terhadap saudaranya di versi yang sama**.
**`KAPASITAS_SURPLUS`** — turunan, `INV-48`.

**Jalur gagal:** Keduanya terisi -> **ditolak** `INV-31`, pesannya menyebut jenis treaty barisnya ·
`JENIS_TREATY` `QUOTA_SHARE` dengan `JUMLAH_LINES_SURPLUS` terisi -> ditolak · Keduanya kosong ->
ditolak · `JENIS_TREATY` bernilai selain kedua nilai sah -> ditolak `INV-30`.

**Uji:** **Negatif — empat, bukan satu:** keduanya terisi; keduanya kosong; jenis `QUOTA_SHARE` dengan
kolom surplus; jenis `SURPLUS` dengan kolom quota share.
**Positif:** baris `QUOTA_SHARE` berpersen 30 -> diterima; baris `SURPLUS` berjumlah 9 lines ->
diterima.
**Positif kedua — dan ia yang menangkap constraint yang terlalu ketat:** **satu layer yang memuat
baris `QUOTA_SHARE` dan baris `SURPLUS` sekaligus** -> **diterima**. Susunan quota share bersurplus
adalah bentuk yang lazim; constraint yang ditulis per layer alih-alih per baris menolaknya dan lulus
keempat uji negatif di atas.

**Menggantikan:** tidak ada aturan lama yang digantikan — **golongannya BARU**, dan karena itu
`CARA MENYALAKANNYA` wajib. Sistem lama memeriksanya **hanya di layar**; berapa banyak baris warisan
yang memuat keduanya **tidak diketahui**, dan itu yang mode peringatan hitung.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemilik proses Treaty In, mingguan**, atas catatan pelanggaran. Angka pertamanya menjawab pertanyaan yang belum pernah dijawab: berapa baris warisan memuat keduanya |
| 3 | **ambang berangka** — kuasa memblokir dinyalakan ketika pelanggaran pada baris **yang lahir di sistem baru** nol selama **4 minggu berturut-turut**. Baris warisan dihitung terpisah dan tidak menahan penyalaan |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 34

**Dasar:**
```
EVIDENCED(SPEC-INVARIAN.md INV-31 - "bentuk lintas baris"; pemeriksaan sistem lama hanya di badan seksi)
        DECIDED(INV-30, INV-31)
```

- [ ] `INV-31` terpasang **di basis data**, bukan hanya di aplikasi — dan letaknya tertulis sebagai keputusan
- [ ] keempat uji negatif lulus, masing-masing dijalankan terpisah
- [ ] kedua uji positif lulus — termasuk satu layer bercampur `QUOTA_SHARE` dan `SURPLUS`
- [ ] pesan penolakan menyebut **jenis treaty barisnya**, bukan hanya "salah satu harus kosong"
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk ambang berangka
- [ ] hitungan baris warisan yang melanggar **dilaporkan terpisah** dari baris baru

## Treaty In - 36 - Baris surplus tanpa baris quota share pada versi yang sama menghasilkan kegagalan yang menyebutkan apa yang kurang

---
status: aktif
golongan: baru
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-23` · `SPEC-INVARIAN.md` `INV-35` · `ADR-D-TI-0035`.*

**What to build:** **PK** yang mencatat baris `SURPLUS` sementara versi itu **tidak punya** baris
`QUOTA_SHARE` menerima **kegagalan yang menyebutkan apa yang kurang** — bukan galat tanpa isi, dan
bukan hasil perhitungan bernilai nol.

**Kenapa begini:** Kapasitas surplus dihitung **retensi × jumlah lines**, dan retensinya datang dari
susunan quota share. Tanpa baris quota share, **penyebutnya tidak ada** — dan sistem lama menyelesaikan
keadaan itu dengan menghasilkan **nol**. Nol yang sebenarnya *"tidak dapat dihitung"* **tidak dapat
dibedakan lagi setelah tersimpan**, dan itu persis larangan `ADR-D-TI-0035`. `INV-35` menyebutnya dengan
kata-kata itu: *"ketiadaannya kegagalan yang dilaporkan, bukan nol"*.

**Persyaratan:** **`INV-35`** (baris `SURPLUS` menuntut adanya baris `QUOTA_SHARE` pada **versi yang
sama**), bergolongan **APLIKASI** · `ADR-D-TI-0035` · `INV-42` (nilai uang tidak pernah memuat angka
penanda kegagalan) · `INV-48` (kapasitas surplus adalah **kelipatan**, pengecualian bernama terhadap
`INV-47`)

**Tidak termasuk:** **Aturan tepat-satu pada satu baris** — irisan `35`. Yang di sini **satu baris
terhadap saudaranya di versi yang sama**.
**`KAPASITAS_SURPLUS` sebagai kolom tersimpan** — tidak ada; turunan.
**Kenapa `INV-35` di lapisan aplikasi dan bukan constraint** — ia melintasi baris di dalam versi lewat
`LAYER`, dan bentuk itu tidak dapat ditulis sebagai `CHECK`. Letaknya **keputusan**, bukan kelalaian.

**Jalur gagal:** Simpan baris `SURPLUS` pada versi tanpa baris `QUOTA_SHARE` -> **kegagalan bernama**
yang menyebut *"tidak ada baris quota share pada versi ini"* · Kapasitas surplus dihitung dan
menghasilkan **nol** -> **kriteria selesai tidak terpenuhi**; itu perilaku yang irisan ini hapus ·
Galat tanpa isi -> idem.

**Uji:** **Negatif:** simpan surplus tanpa quota share pada versi kosong; dan pada versi yang punya
layer lain tetapi tanpa quota share sama sekali.
**Positif — dan ia yang menangkap pemeriksaan yang dijalankan pada lingkup yang salah:** baris
`SURPLUS` pada versi yang punya baris `QUOTA_SHARE` **di layer yang berbeda** -> **diterima**.
`INV-35` berbunyi *"pada versi yang sama"*, bukan *"pada layer yang sama"* — pemeriksaan yang
dijalankan per layer menolak susunan yang sah dan lulus kedua uji negatif di atas.

**Menggantikan:** perilaku sistem lama yang menghasilkan **nol** ketika penyebutnya tidak ada.
**Golongannya BARU** — pemeriksaannya memang belum pernah ada — sehingga `CARA MENYALAKANNYA` wajib,
dan hasilnya **tidak dapat diuji terhadap data lama**: data lama memuat nol yang artinya bukan nol.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemilik proses Treaty In, mingguan**. Angka pertamanya menjawab berapa versi warisan memuat surplus tanpa quota share |
| 3 | **ambang berangka** — kuasa memblokir dinyalakan ketika pelanggaran pada versi **yang lahir di sistem baru** nol selama **4 minggu berturut-turut**; versi warisan dihitung terpisah |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 34

**Dasar:**
```
EVIDENCED(SPEC-INVARIAN.md INV-35 - "ketiadaannya kegagalan yang dilaporkan, bukan nol")
        DECIDED(ADR-D-TI-0035, INV-35, INV-48)
```

- [ ] `INV-35` terpasang **di lapisan aplikasi**, dan letaknya tertulis sebagai keputusan yang menyebut sebabnya
- [ ] kegagalannya **menyebut apa yang kurang**, dan keterangannya **tersimpan**, bukan hanya muncul di layar
- [ ] tidak ada jalur yang menghasilkan kapasitas surplus bernilai nol saat penyebutnya tidak ada — diperiksa, bukan diandaikan
- [ ] uji positif lulus: quota share di layer berbeda pada versi yang sama **diterima**
- [ ] keempat butir **CARA MENYALAKANNYA** terisi
- [ ] hitungan versi warisan yang melanggar dilaporkan terpisah

## Treaty In - 37 - Potongan atas premi dicatat, dengan aturan yang sama di kedua tempat ia melekat

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-19` · `SPEC-MODEL-DATA.md` §14.1, §14.2 · `SPEC-INVARIAN.md` `INV-15`, `INV-17`, `INV-63`.*

**What to build:** **PK** mencatat potongan atas premi — brokerase, komisi, dan sejenisnya — pada
**bagian non-proporsional maupun ketentuan proporsional**, dan **rumus perhitungannya tertulis
sekali** serta berlaku di kedua pelekatan.

Artefak: entitas `POTONGAN` dengan **dua kolom induk bernama**, `CHECK` tepat satu terisi, dan
`INV-15` sebagai **dua** `UNIQUE`.

**PEMBUAT PERTAMA** untuk `POTONGAN`.

**Kenapa begini:** `POTONGAN` **satu konsep dengan dua pelekatan**, dan buktinya satu aktivitas
melayani keduanya dengan satu rumus: `CalculateDeduction` menghitung
`Deduction = @divide(DeductionPct, 100, 4) * GrossPremiumList(1).Value`, dan pemanggilnya mencakup
kedua sisi — `Share.xml`, `ShareRetro.xml`, `TreatyInXOLAddSpreadingDetail` di non-proporsional,
`DetailLimits.xml` di proporsional. **Rumusnya tidak bercabang.** Menyalinnya menjadi dua akan
membuat keduanya berpisah dalam dua tahun.

Bentuk fisiknya **`KTV-B`**: dua kolom bernama, bukan diskriminator. `INV-17` melarang rujukan yang
sasarannya bergantung nilai kolom lain — diskriminator melanggarnya secara harfiah, dua kolom bernama
tidak. Dan keuntungannya bukan sampingan: pada bentuk lama, **basis data tidak menolak baris yatim**;
dengan dua kolom bernama, **dua kunci asing berdiri** dan baris yatim menjadi mustahil.

**Persyaratan:** **`INV-15`** (jenis potongan unik di dalam satu induk) — terpasang sebagai **dua**
`UNIQUE`, satu per pelekatan · **`INV-17`** · **`INV-63`** (aturan potongan tertulis **sekali**,
berlaku pada kedua pelekatannya), bergolongan **APLIKASI, tinjauan skema** · `INV-52` (premi bruto
dikurangi seluruh potongan sama dengan premi bersih) · `INV-41`

**Tidak termasuk:** **Nama potongan sebagai teks.** Sistem lama menuliskannya tetap —
`DeductionList(1).Comment = "Comm to NuRe"`, `"Brokerage fee"`, `"Facultative Brokerage fee"` — dan di
model baru ia **rujukan ke tabel acuan `JENIS_POTONGAN`** (§14.2). Tetapan itu **hilang dengan
sendirinya**; `INV-59` melarang menyalin namanya.
**Penyebaran** — irisan `38`, bentuk induk gandanya sama tetapi entitasnya lain.
**Kolom `ID_INDUK_POTONGAN`** — **tidak ada lagi**. Ia diganti `ID_BAGIAN` + `ID_DETAIL_PROPORSIONAL`
oleh `KTV-B`; siapa pun yang memulihkannya membalikkan keputusan itu tanpa membukanya.

**Jalur gagal:** Kedua kolom induk terisi -> **ditolak** `CHECK` · Kedua kolom induk kosong ->
ditolak · Dua potongan berjenis sama pada induk yang sama -> **ditolak** `INV-15` · Baris potongan
menunjuk induk yang tidak ada -> ditolak kunci asing.

**Uji:** **Negatif:** kedua induk terisi; kedua induk kosong; jenis potongan kembar pada satu
`BAGIAN`; jenis potongan kembar pada satu `DETAIL_PROPORSIONAL`; induk yatim.
**Positif — dan ia yang menangkap `UNIQUE` yang mencampur dua ruang pengenal:** sebuah `BAGIAN`
bernomor 7 dan sebuah `DETAIL_PROPORSIONAL` bernomor 7, **masing-masing dengan potongan berjenis
sama** -> **keduanya diterima**. Pada bentuk lama — satu `UNIQUE` atas `ID_INDUK_POTONGAN` — keduanya
bertabrakan meskipun keduanya sah, dan itu **menolak data yang sah**.
**Positif kedua:** rumus potongan yang sama menghasilkan angka yang sama di kedua pelekatan,
diperiksa dengan **satu** himpunan data uji yang dijalankan dua kali — bukan dua himpunan.

**Menggantikan:** `P-19` melestarikan `DeductionList` dan rumus `CalculateDeduction`. Yang bergeser:
namanya menjadi **rujukan tabel acuan**, dan induk gandanya memperoleh **kunci asing sungguhan** yang
sebelumnya tidak ada sama sekali.

**Blocked by:** 33 · 34 · 15

**Dasar:**
```
EVIDENCED(CalculateDeduction@ekspor-2026-09 - satu rumus, pemanggil di kedua cabang: Share.xml, ShareRetro.xml, TreatyInXOLAddSpreadingDetail, DetailLimits.xml)
        EVIDENCED(DeductionList Comment="Comm to NuRe"/"Brokerage fee" - nama sebagai tetapan)
        DECIDED(KTV-B, INV-15, INV-17, INV-63, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `POTONGAN` berdiri dengan `ID_BAGIAN` dan `ID_DETAIL_PROPORSIONAL` **keduanya nullable**, masing-masing berkunci asing
- [ ] `CHECK` tepat-satu-terisi terpasang, dan **kedua arah** pelanggarannya diuji
- [ ] `INV-15` terpasang sebagai **dua** `UNIQUE`; uji positif membuktikan keduanya tidak saling mencampur
- [ ] `INV-63` diperiksa lewat **tinjauan skema tertulis**: rumus potongan ditemukan **tepat satu kali** di seluruh basis kode
- [ ] uji positif kedua lulus: satu himpunan data uji, dijalankan di kedua pelekatan, hasilnya sama
- [ ] sapuan membuktikan **tidak ada** kolom `ID_INDUK_POTONGAN` yang tertinggal
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 38 - Bagian NuRe disebarkan ke susunan retro internal per jenis reasuransi per mata uang

---
status: tertahan
golongan: belum pasti — L atau B, lihat PENGHALANG
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-18` dan §2.3 · `SPEC-MODEL-DATA.md` §10.6, §10.7, §10.8, §12.3, §12.4 · `SPEC-INVARIAN.md` `INV-16`, `INV-47`, `INV-50`, `INV-51`, `INV-65`.*

**What to build:** **PK** menyebarkan bagian NuRe ke **susunan retro internal**, dirinci **per jenis
reasuransi** dan **per mata uang**, dan jumlah persen penyebarannya **per induk, dikelompokkan menurut
jenis reasuransi, sama dengan 100**.

Artefak: `PENYEBARAN` dengan **dua kolom induk bernama** (`KTV-B`), `RINCIAN_PENYEBARAN`,
`NILAI_PENYEBARAN`; `INV-16` sebagai **dua** `UNIQUE`; `INV-65`; dan **`INV-47`, `INV-50`, `INV-51`
beserta pemantau kebasiannya**.

**PEMBUAT PERTAMA** untuk `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN`.

**Kenapa begini:** Tiga invarian terpenting modul ini **melintasi entitas** dan ditegakkan lewat
*materialized view* — dan **tidak satu pun punya pemantau kebasian hari ini**.

> **MV yang gagal me-refresh berhenti menegakkan tanpa satu galat pun.** Irisan yang hanya memasang
> constraint-nya akan **dinyatakan selesai padahal tidak menegakkan apa pun**, dan tidak ada yang akan
> tahu sampai data yang melanggar sudah masuk.

Dan sumbunya **sempat ditulis salah**: kedua baris rincian semula berbunyi *"per pihak"*.
`SetSpreadName.xml` mengisi `BreakDownSprdList(…).ReinsID` dari `pxResults(…).ReinsTypeID` — **nama
ruasnya menyebut pihak; yang ditulis ke dalamnya jenis**. Seluruh 866 pasangan kelas–properti disapu:
identitas pihak **tidak ada sama sekali** di keluarga penyebaran.

**Persyaratan:** **`INV-16`** (jenis reasuransi unik di dalam satu induk) — **dua** `UNIQUE`, satu per
pelekatan · **`INV-65`** (jenis reasuransi unik di dalam satu `PENYEBARAN`) · **`INV-47`** (jumlah
baris turunan menurut `SUMBU_REKONSILIASI`) · **`INV-50`** (jumlah persen penyebaran per induk,
**dikelompokkan menurut jenis reasuransi**, sama dengan 100) · **`INV-51`** (retensi + penyerahan =
100) — bergolongan **CONSTRAINT BELUM DIBUKTIKAN** · `INV-17` · `KTV-B` · `INV-58` (penyebaran adalah
**fakta terbukukan yang membawa penunjuk asalnya**, dan `INV-58` menyediakan pengecualiannya)

**Tidak termasuk:** **Cabang retro keluar** — `RETRO_KELUAR` adalah **GEL-2**, tegas di luar
penyerahan pertama. Yang di sini **susunan retro internal**, dan pembedanya **arah**.
**Tetapan 15% ke `ORS`** — `FetchQSfromMasterXOL` menulis `Pct = "15.00"` dan `ReinsTypeID = "10007"`
sebagai tetapan di dalam kode. Di sistem baru ia **baris data di susunan retro**, bukan baris kode;
irisan ini **tidak menyalinnya sebagai tetapan**. Pertanyaan apakah 15% ke `ORS` ketentuan umum atau
tambalan **terdaftar sebagai butir wawancara** dan tidak menahan irisan ini.
**Pemantau kebasian untuk `INV-69`/`INV-70`** — tiket `11`, yang **menetapkan polanya**. Irisan ini
**mengikuti** pola itu, bukan sebaliknya.

**Jalur gagal:** Kedua kolom induk terisi atau kedua-duanya kosong -> **ditolak** `CHECK` · Dua
penyebaran berjenis reasuransi sama pada induk yang sama -> **ditolak** `INV-16` · Jumlah persen
penyebaran per kelompok jenis ≠ 100 -> **ditolak** `INV-50` · MV yang **basi** -> **kegagalan
penegakan bernama** yang menyebut **sejak kapan**, bukan diam.

**Uji:** **Negatif:** kedua induk terisi; kedua induk kosong; jenis reasuransi kembar di kedua
pelekatan; jumlah persen 99; jumlah persen 101; **dan uji negatif MV yang dijalankan, bukan
diargumentasikan**.
**Positif:** `BAGIAN` bernomor 7 dan `DETAIL_PROPORSIONAL` bernomor 7, masing-masing dengan penyebaran
berjenis sama -> **keduanya diterima**.
**Positif kedua — dan ia menguji pengelompokan `INV-50` yang dikoreksi:** satu induk dengan **dua
kelompok jenis reasuransi**, masing-masing berjumlah 100 -> **diterima**. Rumusan lama —
*"per susunan"* tanpa pengelompokan — menolaknya, sebab jumlah seluruhnya 200.
**Positif ketiga:** pemantau kebasian **berbunyi kepada seseorang** ketika MV sengaja dibuat basi.

**Menggantikan:** `P-18` melestarikan `SpreadingList` dan `BreakDownSprdList`. Yang bergeser: sumbunya
**per jenis reasuransi**, bukan per pihak — dan itu koreksi yang sudah masuk `INV-50` dan
`SPEC-INVARIAN.md` §4.4a.

**PENGHALANG:**

| | |
|---|---|
| **Apa yang ditunggu (1)** | **`Uji AD`** — apakah persentase rincian penyebaran **pernah menyimpang** dari tabel master. Bila **pernah**, `P-18` bergolongan **L**; bila **tidak pernah**, ia **B** dan irisan ini **menuntut medan `CARA MENYALAKANNYA`** yang belum ada |
| **Siapa dapat menjawabnya (1)** | **kantor** — izin kueri baca-saja ke produksi |
| **Apa yang berubah (1)** | golongan tiket, dan ada-tidaknya medan `CARA MENYALAKANNYA`. **Bentuk tabelnya tidak berubah** |
| **Apa yang ditunggu (2)** | **`L-3`** — tidak ada instans Oracle yang terjangkau, sehingga uji negatif MV untuk `INV-50` dan `INV-51` **belum pernah dijalankan di mana pun** |
| **Siapa dapat menjawabnya (2)** | **kantor** — sediakan instans yang terjangkau |
| **Apa yang berubah (2)** | dari **klaim** menjadi **penegakan**. Tanpa itu, ketiga MV berstatus klaim dan tiket ini **tidak dapat dinyatakan selesai** |

**Taksiran: kosong**, dan kekosongan itu **disengaja** — `G3` melarang taksiran pada tiket yang
penghalangnya belum terjawab.

**Blocked by:** 33 · 34 · 15 · **`Uji AD`** · **`L-3`**

**Dasar:**
```
EVIDENCED(SetSpreadName@ekspor-2026-09 - BreakDownSprdList.ReinsID diisi dari pxResults.ReinsTypeID)
        EVIDENCED(sapuan 866 pasangan kelas-properti - identitas pihak NOL di keluarga penyebaran)
        EVIDENCED(FetchQSfromMasterXOL@ekspor-2026-09 - Pct="15.00", ReinsTypeID="10007" sebagai tetapan)
        DECIDED(KTV-B, INV-16, INV-50, INV-58, INV-65, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(Uji AD)
        DIASUMSIKAN-CLEAR(L-3)
```

- [ ] `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `PENYEBARAN` memakai **dua kolom induk bernama** + `CHECK`, sama bentuknya dengan `POTONGAN` di irisan `37`
- [ ] `INV-16` terpasang sebagai **dua** `UNIQUE`; uji positif membuktikan keduanya tidak mencampur ruang pengenal
- [ ] `INV-50` terpasang **dengan pengelompokan menurut jenis reasuransi**, dan uji positif kedua lulus
- [ ] **uji negatif MV dijalankan** untuk `INV-47`, `INV-50`, `INV-51` — bukan diargumentasikan
- [ ] **pemantau kebasian berdiri untuk ketiganya**: `REFRESH_MODE`, `STALENESS`, terjadwal, dan **berbunyi kepada seseorang yang bernama**
- [ ] pemantau memakai **pola yang tiket `11` tetapkan**, bukan pola kedua
- [ ] kegagalan pemantau tercatat sebagai **kegagalan penegakan bernama** yang menyebut **sejak kapan**
- [ ] tetapan 15% ke `ORS` **tidak disalin sebagai tetapan**; pertanyaan wawancaranya terdaftar
- [ ] golongan tiket ditetapkan sesudah `Uji AD` kembali, dan bila **B**, medan `CARA MENYALAKANNYA` ditambahkan dengan keempat butirnya
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`

## Treaty In - 39 - Pemeriksa jejak melihat siapa mengubah fakta apa dan kapan, untuk kontrak mana pun

---
status: aktif
golongan: baru
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-47` · `SPEC-MODEL-DATA.md` §10.21 · `ADR-D-TI-0045` · `ADR-D-TI-0044` · `KTV-D`.*

**What to build:** **PJ** membuka kontrak mana pun dan melihat **siapa mengubah ruas apa, dari nilai
apa ke nilai apa, kapan, dan di bawah peran apa**. Jejaknya terisi pada **setiap penulisan**, bukan
hanya pada perpindahan persetujuan.

Artefak: entitas `JEJAK_PERUBAHAN` dengan kedelapan atributnya, dan mekanisme pengisiannya.

**PEMBUAT PERTAMA** untuk `JEJAK_PERUBAHAN`.

**Kenapa begini:** **Jejak perubahan adalah bukti bahwa pembekuan benar-benar terjadi.** `ADR-D-TI-0036`
menetapkan angka yang disetujui dibekukan, dan modul Adjustment mengambil data lama lewat `SELECT`
alih-alih menyalinnya. Keduanya hanya dapat dipercaya bila ada cara **membuktikan data lama tidak
berubah di antara dua pembacaan**. Tanpa jejak, pembekuan adalah **janji**, bukan fakta — dan selisih
yang dibukukan Adjustment berdiri di atas janji.

Di sistem lama satu-satunya jejak adalah daftar komentar, ditulis **hanya pada perpindahan
persetujuan**. Penyimpanan biasa — **termasuk yang mengubah angka** — tidak meninggalkan jejak siapa
maupun kapan, dan **tidak ada kolom waktu sama sekali** di tabel kontrak maupun tabel addendum.

**Persyaratan:** **`ADR-D-TI-0045`** — isi minimalnya *"siapa, kapan, apa yang berubah dari nilai apa ke
nilai apa, dan **di bawah peran apa**"* · **`KTV-D`** — `PERAN_PELAKU` adalah **POTRET**: teks, final,
**tidak dinormalisasi ulang** ketika entitas peran kelak lahir · `ADR-D-TI-0041` (satu fakta, satu penulis) ·
`INV-01`

**Tidak termasuk:** **Catatan persetujuan** — benda yang **berbeda**, dan `ADR-D-TI-0045` memisahkannya
tegas: catatan adalah narasi manusia yang **dapat disunting**, jejak adalah fakta mesin yang **tidak
dapat**. Menggabungkannya menghasilkan yang terburuk dari keduanya — catatan yang bisa disunting dan
karena itu tidak membuktikan apa-apa. `CATATAN_PERSETUJUAN` milik batch 2.
**Entitas peran dan penugasan bertanggal** — **`F-16`**, lubang terbuka: `ADR-D-TI-0044` menempatkan
keduanya di gelombang 1 dan §10 tidak memuat satu pun. **Ia tidak menahan irisan ini**, sebab `KTV-D`
sudah memutuskan bentuk `PERAN_PELAKU`, dan bentuk itu **tidak berubah** apa pun jadinya entitas peran
nanti.
**Layar pemeriksa jejak** — `L-4`: tidak ada spesifikasi layar di mana pun, dan lubang itu milik batch
lapisan aplikasi. Yang dibangun di sini **jalur bacanya**, bukan tampilannya.

**Jalur gagal:** Sebuah penulisan yang **tidak** meninggalkan baris jejak -> **kriteria selesai tidak
terpenuhi** · Baris jejak yang disunting sesudah ditulis -> **ditolak**; entitas ini **tambah-saja** ·
`PERAN_PELAKU` kosong -> ditolak · Baris jejak dihapus -> ditolak.

**Uji:** **Negatif:** `UPDATE` atas baris jejak; `DELETE` atas baris jejak; simpan jejak tanpa peran;
lakukan penulisan lewat jalur yang **tidak** menuliskan jejak dan pastikan jalurnya **tidak ada**.
**Positif — dan ia yang menangkap jejak yang terlalu sempit:** ubah **satu ruas yang bukan bagian
persetujuan** — misalnya `LINGKUP_WILAYAH` — dan pastikan jejaknya **tetap tercatat**. Sistem lama
mencatat hanya pada perpindahan persetujuan; mekanisme yang meniru bentuk lama lulus ketiga uji
negatif di atas dan **melewatkan justru perubahan yang paling sering terjadi**.
**Positif kedua:** dua perubahan pada **detik yang sama** menghasilkan **dua baris**. `WAKTU_PERUBAHAN`
bertipe `DATE` beresolusi detik, dan entitas ini **sengaja tanpa kunci alami** justru supaya keadaan
itu sah.

**Menggantikan:** daftar komentar sistem lama sebagai satu-satunya jejak. **Golongannya BARU** —
jejak perubahan sebagai fakta mesin memang belum pernah ada — sehingga `CARA MENYALAKANNYA` wajib, dan
**hasilnya tidak punya pembanding sama sekali** di data lama.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemeriksa jejak (PJ), bulanan**, atas cacah baris jejak per kontrak. Angka nol pada kontrak yang jelas berubah adalah tanda mekanismenya tidak terpasang di salah satu jalur |
| 3 | **ambang berangka** — jejak dinyatakan lengkap ketika **setiap** penulisan atas `VERSI_KONTRAK` dan anaknya menghasilkan sedikitnya satu baris jejak, diperiksa atas contoh **100 penulisan berturut-turut** dengan hasil **100** |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In, bersama `PJ` |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(ADR-D-TI-0045 konteks - daftar komentar sebagai satu-satunya jejak; nol kolom waktu di tabel kontrak dan addendum)
        DECIDED(ADR-D-TI-0045, ADR-D-TI-0044, KTV-D, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(F-16)
```

- [ ] `JEJAK_PERUBAHAN` berdiri dengan **kedelapan** atributnya sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] entitas ini **tambah-saja**: `UPDATE` dan `DELETE` atasnya ditolak, keduanya diuji terpisah
- [ ] `PERAN_PELAKU` tersimpan sebagai **teks beku**, dan `KTV-D` dirujuk di dalam berkas DDL-nya sebagai pernyataan keputusan
- [ ] kedua uji positif lulus — ruas non-persetujuan tetap berjejak, **dan** dua perubahan sedetik menghasilkan dua baris
- [ ] tidak ada jalur penulisan yang melewati jejak — diperiksa dengan sapuan, hasilnya dicetak
- [ ] `F-16` tercatat di `ASUMSI-CLEAR.md` sebagai lubang yang **tidak menahan** tiket ini, beserta sebabnya
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk ambang berangka

## Treaty In - 40 - Pengisi kontrak melihat nilai versi sebelumnya berdampingan dengan versi yang sedang disusun

---
status: aktif
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-44` · `ADR-D-TI-0048` butir 2 · `GRL-03`, `GRL-10` · `SPEC-INVARIAN.md` `INV-58`, `INV-59`.*

**What to build:** **PK** yang menyusun versi penyesuaian melihat **nilai versi dasarnya berdampingan**
dengan nilai yang sedang ia isi — dan nilai lama itu **di-SELECT lewat `ID_VERSI_KONTRAK_DASAR`,
bukan disalin** ke kolom mana pun.

**Kenapa begini:** Sistem lama menyalin **seluruh halaman clipboard** versi dasar ke dalam sebuah
pohon bernama `OLDDATA` di dalam `JSONDATA` versi baru. Salinan itu melahirkan **dua sumber kebenaran
untuk satu fakta**: ketika versi dasarnya dibetulkan, salinannya **tidak ikut berubah**, dan selisih
yang dihitung terhadap salinan basi **tidak menghasilkan galat** — ia menghasilkan angka yang salah
dengan tenang. `ADR-D-TI-0048` butir 2 membalikkannya: **tidak ada tabel `OLDDATA`**; sisi lama di-`SELECT`
ke versi dasar lewat rujukannya.

**Persyaratan:** `ADR-D-TI-0048` butir 2 · `INV-59` (tidak ada nilai yang disalin dari entitas lain; hilir
diberi rujukan) · `INV-58` (tidak ada kolom turunan yang disimpan, kecuali ia fakta terbukukan yang
membawa penunjuk asalnya) · `GRL-11` (turunan *"versi berlaku"* **dihitung, bukan disimpan**)

**Tidak termasuk:** **Kolom `ID_VERSI_KONTRAK_DASAR` itu sendiri** — **tiket `01`** yang
membangunnya, beserta constraint keterisiannya. Irisan ini **membacanya**.
**Baris selisih** — tiket `06`. Yang di sini **penyandingan untuk dibaca orang**; yang di sana
**besaran berselisih yang dibukukan**.
**Pohon `ActualValue`** — **dibuang seluruhnya** (`GRL-14`); premi aktual menjadi nilai versinya
sendiri. Pernyataan keputusan, bukan pekerjaan tertunda.
**Layar penyandingnya** — `L-4`; yang dibangun **jalur bacanya**.

**Jalur gagal:** Sebuah kolom pada `VERSI_KONTRAK` yang memuat salinan nilai versi dasar ->
**tidak ada**, dan ketiadaannya diperiksa dengan sapuan atas `KAMUS-KOLOM.md` · Versi dasar dibetulkan
sesudah versi baru dibuat, lalu nilai lama dibaca -> **menampilkan nilai yang sudah dibetulkan**,
bukan yang basi · Versi pertama sebuah kontrak — dasarnya kosong -> **jalur bacanya menyatakan tidak
ada versi sebelumnya**, bukan galat dan bukan nilai kosong yang menyerupai nol.

**Uji:** **Negatif:** minta nilai versi sebelumnya untuk versi pertama; pastikan hasilnya **pernyataan
ketiadaan**, bukan nol. Sapu `KAMUS-KOLOM.md` untuk kolom salinan.
**Positif — dan ia yang membuktikan SELECT benar-benar SELECT:** buat versi penyesuaian, lalu
**betulkan sebuah nilai pada versi dasarnya**, lalu baca nilai lamanya. Hasilnya **harus nilai yang
sudah dibetulkan**. Sebuah rancangan yang diam-diam menyalin lulus setiap uji negatif di atas — sebab
salinan dan rujukan **memberi jawaban yang sama sampai salah satunya berubah**.

**Menggantikan:** pohon `OLDDATA` di dalam `JSONDATA`. Golongannya **PELESTARIAN** — kemampuan
melihat nilai lama berdampingan memang sudah ada — tetapi **caranya berubah**, dan bunyi lamanya
disebut di sini supaya yang sudah membacanya tahu apa yang bergeser.

**Blocked by:** 14 · **01**

**Dasar:**
```
EVIDENCED(TreatyIn.OLDDATA@ekspor-2026-09 - salinan penuh halaman clipboard versi dasar di dalam JSONDATA)
        DECIDED(ADR-D-TI-0048, GRL-03, GRL-10, GRL-11, GRL-14, INV-58, INV-59)
```

- [ ] jalur baca nilai versi dasar berdiri dan memakai `ID_VERSI_KONTRAK_DASAR` — diperiksa di rencana kueri, bukan diandaikan
- [ ] sapuan atas `KAMUS-KOLOM.md` membuktikan **nol** kolom salinan nilai versi dasar
- [ ] uji positif lulus: pembetulan pada versi dasar **terlihat** pada pembacaan berikutnya
- [ ] versi pertama menghasilkan **pernyataan ketiadaan**, bukan nol
- [ ] ketiadaan tabel `OLDDATA` dan pembuangan `ActualValue` tertulis sebagai **pernyataan keputusan** yang menyebut `ADR-D-TI-0048` dan `GRL-14`

## Treaty In - 41 - Konsumen hilir yang tidak dikenal tetap dapat membaca identitas kontrak dalam bentuk lama

---
status: aktif
golongan: baru
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-54` · `ADR-D-TI-0051` · `SPEC-MODEL-DATA.md` §10.1, §10.2.*

**What to build:** Sebuah **bentuk baca** menyajikan identitas kontrak **dalam bentuk sistem lama** —
pengenal bertipe teks berpola lama, cedant, asal bisnis, periode — sehingga konsumen hilir yang
**belum diketahui siapa saja** tetap dapat membacanya pada hari peralihan.

**Kenapa begini:** `ADR-D-TI-0051` memerintahkan **anggap ada konsumen hilir**, dan alasannya terbaca dari
ekspor: tabel datar `TREATYINDETAIL` dan `PROPORTIONALARRG` ditulis pada **setiap** penyimpanan,
dengan 52 kolom bertipe `VARCHAR2(1000)`, dan **tidak satu pun aturan di dalam ekspor membacanya
kembali**. Tabel yang ditulis tetapi tidak dibaca oleh sistemnya sendiri **dibaca oleh sesuatu yang di
luar ekspor** — dan kita tidak tahu apa.

> **Diam bukan bukti.** Ketiadaan pembaca di dalam ekspor tidak membuktikan ketiadaan pembaca; ia
> hanya membuktikan pembacanya **tidak ada di sini**.

**Persyaratan:** `ADR-D-TI-0051` · `INV-59` (hilir diberi **rujukan**, bukan salinan) · `INV-58` (bentuk
baca adalah **turunan yang tidak disimpan**) · `ADR-D-TI-0042`

**Tidak termasuk:** **Menulis ulang tabel datar lama.** Bentuk baca ini **diturunkan**, bukan
dituliskan aplikasi — `SaveTreatyInDetail_Act` yang menulis nol tetap ke tiga kolomnya **tidak
ditiru**, dan `TETAPAN-DI-KODE.md` §3 sudah menyatakan tabel datar di model baru **diturunkan**.
**Jalur penerbitan ke luar** — `G2` menyatakannya **tegas tidak masuk**: `TREATYINOFFER` tidak punya
satu pun penulis yang terjangkau dan penggantinya belum diketahui (`DAFTAR-ESKALASI-MANAJEMEN.md`
butir 3).
**Isi kontrak selengkapnya dalam bentuk lama** — hanya **identitas** yang dijanjikan di sini. Menjanji
lebih berarti menjanji bentuk yang belum ada yang memintanya.

**Jalur gagal:** Bentuk baca yang **menyimpan** hasilnya sebagai tabel -> melanggar `INV-58`;
ia diturunkan saat dibaca · Pengenal bentuk lama dibangkitkan dari cap waktu -> ditolak `INV-02`; ia
**dibaca dari `NOMOR_KONTRAK_WARISAN`**, dan kontrak yang lahir di sistem baru **tidak punya** — itu
keadaan yang dinyatakan, bukan ditambal.

**Uji:** **Negatif:** minta bentuk lama untuk kontrak yang lahir di sistem baru -> **pernyataan bahwa
nomor lamanya tidak ada**, bukan nomor karangan.
**Positif — dan ia yang membuktikan bentuknya benar-benar sama:** ambil **sepuluh kontrak warisan**,
bandingkan keluaran bentuk baca ini dengan baris `TREATYINDETAIL` lama untuk kontrak yang sama.
**Kesepuluhnya harus identik pada kolom identitas.** Perbandingan terhadap satu kontrak tidak cukup:
sebuah bentuk baca yang benar untuk satu baris dan salah untuk sisanya lulus uji tunggal mana pun.

**Menggantikan:** tidak ada aturan yang digantikan. **Golongannya BARU** — jaminan kompatibilitas
baca memang belum pernah dinyatakan — sehingga `CARA MENYALAKANNYA` wajib.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemilik proses Treaty In, mingguan** selama delapan minggu pertama sesudah peralihan, atas **catatan akses** ke bentuk baca ini. Catatan itulah satu-satunya cara mengetahui **siapa** konsumen hilirnya |
| 3 | **ambang berangka** — bentuk baca ini **dapat dicabut** ketika catatan akses menunjukkan **nol akses selama 8 minggu berturut-turut**. Sampai itu terjadi, ia dipertahankan — dan bila ada akses, pemiliknya **dicari dan dicatat namanya** |
| 4 | **siapa boleh menyalakan atau mencabutnya** — pemilik proses Treaty In |

> Butir 3 sengaja berbentuk **syarat pencabutan**, bukan syarat penyalaan: fitur ini menyala sejak
> hari pertama, dan yang belum diketahui adalah **kapan ia boleh mati**.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(TREATYINDETAIL@Table/ dan PROPORTIONALARRG@Table/ - ditulis tiap simpan, NOL pembaca di dalam ekspor)
        DECIDED(ADR-D-TI-0051, ADR-D-TI-0042, INV-58, INV-59)
```

- [ ] bentuk baca identitas dalam bentuk lama berdiri, dan ia **diturunkan** — bukan tabel yang ditulis aplikasi
- [ ] kontrak yang lahir di sistem baru menghasilkan **pernyataan ketiadaan nomor lama**, bukan nomor karangan
- [ ] uji positif lulus atas **sepuluh** kontrak warisan, bukan satu
- [ ] **catatan akses** berdiri, sebab tanpanya butir 3 `CARA MENYALAKANNYA` tidak dapat diukur
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk **syarat pencabutan berangka**
- [ ] jalur penerbitan ke luar dinyatakan **di luar irisan ini**, dengan eskalasi butir 3 sebagai penagihnya

## Treaty In - 42 - Pelaksana migrasi menyimpan dokumen JSON sistem lama sebagai arsip, dan dapat menunjukkan arsipnya ada

---
status: aktif
golongan: baru
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-55` · `ADR-D-TI-0034` · `SPEC-INVARIAN.md` `INV-61`.*

**What to build:** **PM** menyimpan dokumen `JSONDATA` sistem lama **sebagaimana adanya** pada saat
pemindahan, dan dapat **menunjukkan bahwa arsip itu ada** untuk kontrak mana pun — tanpa aplikasi
pernah membacanya sebagai sumber.

**Kenapa begini:** `JSONDATA` adalah **tempat kebenaran sistem lama** — seluruh halaman clipboard
sebagai satu kolom — dan model baru menggantikannya dengan tabel relasional. Selama beberapa tahun
pertama, pertanyaan *"apakah migrasi kehilangan sesuatu"* akan muncul, dan **satu-satunya cara
menjawabnya adalah membuka dokumen aslinya**. Tetapi arsip yang punya jalur baca akan **dipakai
sebagai sumber** oleh orang yang sedang buru-buru, dan sejak saat itu ada dua kebenaran. `ADR-D-TI-0034`
menyelesaikannya dengan tegas: **arsipnya disimpan, jalur bacanya tidak ada.**

> Bentuk kemampuannya sempat ditulis sebagai *"arsip tidak punya jalur baca"* dan **gugur uji U-6** —
> itu **larangan**, bukan kemampuan, dan tidak ada pelaku yang "melakukannya". Yang nyata:
> **menyimpannya**, dipegang `PM`, dan dapat dinyatakan selesai. `INV-61` menjadi invarian yang dibawa,
> bukan tiket.

**Persyaratan:** `ADR-D-TI-0034` · **`INV-61`** (dokumen JSON arsip **tidak punya jalur baca**; ia bukan
sumber kanonik) · `ADR-D-TI-0042`

**Tidak termasuk:** **Jalur baca isi arsip** — **tegas tidak dibangun**, dan itu **pernyataan
keputusan**. Siapa pun yang menambahkannya membalikkan `ADR-D-TI-0034` tanpa membukanya.
**Pembandingan isi arsip terhadap hasil migrasi** — itu **uji paritas migrasi**, milik irisan migrasi
batch 2.
**Kalimat untuk pembaca arsip** — **`D-7`**, diff modul Adjustment yang **masih diparkir**; dasarnya
sudah lewat dan ia menunggu satu kalimat izin pemilik proses. Irisan ini **tidak mendahuluinya**.
**Penyimpanan berkasnya di luar basis data** — bentuk fisik penyimpanan bukan keputusan irisan ini.

**Jalur gagal:** Sebuah kueri aplikasi yang **membaca isi** arsip -> **tidak ada jalurnya**, dan
ketiadaannya diperiksa dengan sapuan · Arsip tersimpan tetapi **tidak dapat ditunjukkan ada** untuk
sebuah kontrak -> kriteria selesai tidak terpenuhi · Arsip disimpan sesudah migrasi, bukan **pada
saat** migrasi -> ditolak; isinya sudah tidak asli.

**Uji:** **Negatif:** sapu seluruh basis kode untuk kueri yang menyentuh isi arsip; hasilnya **harus
nol**, dan angkanya dicetak.
**Positif — dan ia yang menangkap arsip yang tidak lengkap:** untuk **setiap** kontrak warisan yang
dipindahkan, keberadaan arsipnya dapat ditunjukkan. **Cacah arsip = cacah kontrak warisan**, dan
selisihnya **dilaporkan per kontrak**, bukan sebagai satu angka. Sebuah mekanisme yang menyimpan
arsip untuk sebagian kontrak lulus uji negatif di atas dan **kehilangan justru kontrak yang paling
aneh** — yang biasanya yang paling ingin diperiksa orang.

**Menggantikan:** tidak ada. **Golongannya BARU**: sistem lama **adalah** `JSONDATA`-nya; gagasan
"arsip" baru ada karena ada penggantinya.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **PM, satu kali pada tiap gelombang pemindahan**, atas selisih cacah arsip terhadap cacah kontrak warisan. Sesudah gelombang terakhir, **pemilik proses Treaty In, sekali** |
| 3 | **ambang berangka** — pemindahan sebuah gelombang dinyatakan lengkap ketika **selisih cacahnya nol**. Selisih bukan nol **menahan gelombang itu**, dan tiap kontrak yang hilang arsipnya disebut namanya |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(M_TREATY_IN.JSONDATA@Table/ - seluruh halaman clipboard sebagai satu kolom)
        DECIDED(ADR-D-TI-0034, INV-61)
        DIASUMSIKAN-CLEAR(D-7)
```

- [ ] arsip tersimpan **pada saat** pemindahan, dan keberadaannya dapat ditunjukkan per kontrak
- [ ] sapuan basis kode untuk pembaca isi arsip kembali **nol**, dan angkanya **dicetak** — bukan disimpulkan
- [ ] uji positif lulus: cacah arsip = cacah kontrak warisan, selisihnya dilaporkan **per kontrak**
- [ ] ketiadaan jalur baca tertulis sebagai **pernyataan keputusan** yang menyebut `ADR-D-TI-0034` dan `INV-61`
- [ ] `D-7` tercatat di `ASUMSI-CLEAR.md`; kalimat untuk pembaca arsip **bukan** bagian tiket ini
- [ ] keempat butir **CARA MENYALAKANNYA** terisi

## Treaty In - 43 - Pelaksana migrasi mematikan penegakan trigger selama pemindahan dan menyalakannya kembali, dan keadaan sakelarnya terlihat

---
status: aktif
golongan: baru
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-53` · `SPEC-INVARIAN.md` §4.3b, CO-1…CO-8 · `ADR-D-TI-0042`.*

**What to build:** **PM** mematikan penegakan trigger selama pemindahan data, menyalakannya kembali
sesudahnya, dan **keadaan sakelar itu terlihat tanpa membuka basis data**. Menyalakan kembali
**memeriksa ulang** baris yang masuk selama ia mati.

**Kenapa begini:** `ADR-D-TI-0042` memerintahkan sejarah **pindah apa adanya**, dan sebagian data lama
**memang melanggar** aturan yang sistem baru tegakkan — kunci alami yang berubah di tengah jalan,
keadaan yang tidak punya padanan, nilai nol yang sebenarnya kegagalan. Trigger yang hidup akan
menolaknya, dan pemindahan **berhenti pada baris pertama yang menyimpang**. Mematikannya adalah jalan
yang benar; yang berbahaya adalah **lupa menyalakannya kembali**.

> Sakelar yang keadaannya hanya terbaca dengan membuka basis data akan dibiarkan mati, dan **tidak
> ada yang tahu sampai data yang melanggar sudah masuk lewat jalur biasa**. Itu bentuk kegagalan yang
> sama dengan *materialized view* yang berhenti me-refresh — penegakan yang berhenti **tanpa satu
> galat pun**.

**Persyaratan:** `SPEC-INVARIAN.md` §4.3b dan **CO-1…CO-8** · `ADR-D-TI-0042` · `INV-26`
(`WARISAN_TAK_TERPETAKAN` hanya dapat dimasuki lewat migrasi, **tidak pernah oleh sistem berjalan**) —
sakelar ini **yang membuat larangan itu dapat ditegakkan**

**Tidak termasuk:** **Pemindahan datanya sendiri** — irisan `44` untuk tabel acuan, dan sisanya batch
2. Irisan ini membangun **sakelarnya**, dan mengujinya dengan data uji.
**Daftar trigger yang dimatikan** — daftarnya **dikumpulkan dari tiket yang memasangnya**, dan tiap
tiket bertrigger wajib mendaftarkan diri: irisan `18` (`INV-19`) dan irisan-irisan batch 2 (`INV-28`,
`INV-54`). Yang dibangun di sini **mekanismenya**, bukan daftar isinya.
**Mematikan constraint yang bukan trigger** — `CHECK` dan `UNIQUE` **tidak dimatikan**. Membiarkan
kunci ganda masuk merusak hal yang tidak dapat diperbaiki tanpa memilih baris mana yang dibuang.

**Jalur gagal:** Menyalakan kembali tanpa memeriksa ulang baris yang masuk -> **kriteria selesai tidak
terpenuhi** · Keadaan sakelar tidak terlihat di luar basis data -> idem · Pemeriksaan ulang menemukan
pelanggaran -> **dilaporkan per baris**, dan pemindahannya **tidak dinyatakan selesai** · Sistem
berjalan menulis baris saat sakelar mati -> **dilarang**; sakelar mati hanya sah selama jendela
pemindahan, dan jendelanya tercatat.

**Uji:** **Negatif:** nyalakan kembali dengan sengaja menyisipkan baris melanggar saat mati —
pemeriksaan ulang **harus menemukannya**; matikan sakelar lalu tulis lewat jalur aplikasi biasa —
**harus ditolak**.
**Positif — dan ia yang menangkap sakelar yang terlalu lebar:** saat sakelar mati, `UNIQUE` dan
`CHECK` **tetap menolak** baris yang melanggar. Sebuah sakelar yang mematikan **segalanya** lulus uji
negatif pertama dan **membiarkan kerusakan yang tidak dapat diperbaiki**.
**Positif kedua:** keadaan sakelar terbaca **benar** pada kedua posisinya, diperiksa dari luar basis
data.

**Menggantikan:** tidak ada. **Golongannya BARU** — sistem lama tidak punya trigger yang perlu
dimatikan, sebab ia hampir tidak punya penegakan di basis data sama sekali.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-1…CO-8**, dan §4.3b |
| 2 | **siapa membaca, seberapa sering** — **PM, pada setiap pembukaan dan penutupan jendela pemindahan**; dan **pemilik proses Treaty In, harian** selama gelombang pemindahan berjalan. Keadaan sakelar adalah hal yang **harus dilihat tiap hari**, bukan tiap minggu |
| 3 | **ambang berangka** — jendela pemindahan dinyatakan tertutup ketika sakelar **menyala** dan pemeriksaan ulang mengembalikan **nol pelanggaran yang belum diadili**. Pelanggaran yang sudah diadili dan diputuskan tetap dibawa **dihitung terpisah** dan tidak menahan penutupan |
| 4 | **siapa boleh menyalakan atau mematikan** — **PM**, dan **hanya** di dalam jendela pemindahan yang pemilik proses buka |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(SPEC-INVARIAN.md §4.3b - urutan penyalaan penegakan; CO-1..CO-8)
        DECIDED(ADR-D-TI-0042, INV-26)
        DIASUMSIKAN-CLEAR(L-3)
```

- [ ] sakelar berdiri, dan keadaannya **terbaca dari luar basis data** pada kedua posisinya
- [ ] menyalakan kembali **memeriksa ulang** baris yang masuk selama mati, dan melaporkan pelanggaran **per baris**
- [ ] uji positif lulus: `UNIQUE` dan `CHECK` **tetap menolak** saat sakelar mati
- [ ] jalur aplikasi biasa **ditolak** selama sakelar mati — diuji, bukan diandaikan
- [ ] **daftar trigger** yang tunduk pada sakelar ini dibangkitkan dari tiket yang memasangnya, bukan diketik ulang
- [ ] jendela pemindahan **tercatat** — kapan dibuka, kapan ditutup, oleh siapa
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk ambang berangka
- [ ] `L-3` tercatat di `ASUMSI-CLEAR.md`: ujinya menuntut instans yang terjangkau

## Treaty In - 44 - Isi keenam tabel acuan dipindahkan dari sistem lama

---
status: tertahan
golongan: pelestarian
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-50` (**pecahan tabel acuan** — lihat `Tidak termasuk`) · `ADR-D-TI-0042` · `ADR-D-TI-0043` · `SPEC-MODEL-DATA.md` §10.22.*

**What to build:** **PM** memindahkan isi keenam himpunan acuan — **mata uang, jenis potongan, jenis
reasuransi, bahaya, kelompok treaty, kelas bisnis** — dari sumber lamanya **apa adanya, tanpa
menghitung ulang apa pun**. Baris yang tidak dapat dipetakan menjadi **pengecualian migrasi
bernomor**, bukan baris yang didiamkan.

**Kenapa begini:** `P-50` adalah **payung**, dan `DAFTAR-PEKERJAAN.md` memerintahkan ia **dipecah per
kelompok entitas sebelum penaksiran** — kemampuan payung yang ditaksir sebagai satu **selalu**
ditaksir terlalu rendah, sebab yang ditaksir **kalimatnya**, bukan isinya. Pecahan tabel acuan
**sendirian masuk batch 1**, dan sebabnya tajam: ia **tidak menulis satu pun nilai keadaan**. Kelima
pecahan lain — kepala kontrak, daftar anak versi, cabang non-proporsional, cabang proporsional,
potongan dan penyebaran — menggantung pada pemindahan kepala, yang **menulis `KEADAAN_SIKLUS_HIDUP`**
termasuk `WARISAN_TAK_TERPETAKAN`, dan karena itu **menunggu `REV-3`**.

Dan ia **didahulukan** bukan karena kecil: **setiap kunci asing di seluruh skema menunggunya**.
Memindahkan kontrak sebelum tabel acuannya terisi berarti setiap baris ditolak kunci asing.

**Persyaratan:** `ADR-D-TI-0042` (sejarah pindah apa adanya) · **`ADR-D-TI-0043`** (menghitung ulang adalah
**peristiwa bisnis tersendiri**, bukan bagian migrasi) · `INV-68` (kode unik) · `INV-44` · `INV-62`

**Tidak termasuk:** **Kelima pecahan `P-50` yang lain** — **batch 2**, dan sebabnya tertulis di
`USULAN-IRISAN-TREATY-IN-BATCH-1.md` §0.2.
**Perbaikan data lama dan hitung ulang** — `G2` menyatakannya **tegas tidak masuk** penyerahan
pertama: `ADR-D-TI-0043` menjadikannya peristiwa bisnis tersendiri. Baris acuan yang **salah eja** di
sistem lama dipindahkan **salah eja**, dan pembetulannya pekerjaan lain yang punya pelakunya sendiri.
**Penambahan baris acuan baru** — irisan `15` sudah membuktikan ia mungkin; yang di sini
**pemindahan isi yang sudah ada**.

**Jalur gagal:** Baris acuan lama berkode ganda -> **ditolak** `INV-68`, dan **kedua** baris dilaporkan
sebagai pengecualian bernomor — bukan salah satunya dibuang diam-diam · Baris acuan lama tanpa kode ->
pengecualian bernomor · Migrasi membetulkan ejaan -> **dilarang** `ADR-D-TI-0042`; ia hitung ulang yang
menyamar sebagai kerapian · Sakelar trigger menyala saat pemindahan berjalan -> pemindahan **tidak
dimulai**.

**Uji:** **Negatif:** sisipkan sumber berkode ganda; sumber tanpa kode; jalankan pemindahan dengan
sakelar menyala.
**Positif — dan ia pokok irisan ini:** sesudah pemindahan, **cacah baris acuan baru = cacah baris
sumber dikurangi cacah pengecualian bernomor**, dan persamaan itu **menutup untuk keenam tabel**.
Sebuah pemindahan yang diam-diam membuang baris yang tidak dikenalinya lulus setiap uji negatif di
atas — sebab yang dibuangnya **tidak pernah muncul di mana pun**.
**Positif kedua:** sebuah baris acuan yang **salah eja di sumber** mendarat **salah eja** — dan itu
**keberhasilan**, bukan kegagalan. Uji paritas yang menandainya sebagai cacat **salah membaca
`ADR-D-TI-0042`**.

**Menggantikan:** `P-50` bunyi lama: *"**PM** dapat memindahkan kontrak lama apa adanya, **tanpa
menghitung ulang apa pun**"* — utuh, sebagai payung. Yang bergeser: **lingkupnya dipecah**, dan
pecahan ini membawa **hanya tabel acuan**. Bunyi lengkapnya tetap berlaku untuk kelima pecahan lain.

**PENGHALANG:**

| | |
|---|---|
| **Apa yang ditunggu** | **`KTV-A`**, dan ia **bertenggat**. Presisi kolom boleh dipersempit **hanya sebelum data dimuat** — dan irisan inilah **yang memuat data pertama**. Sesudah ia berjalan, mempersempit kolom menuntut memeriksa setiap baris, dan baris yang tidak muat **tidak punya tempat pergi** |
| **Siapa dapat menjawabnya** | **gerbang pembuka sesi DDL**, bersama pemilik proses |
| **Apa yang berubah** | tidak ada isi tiket yang berubah. Yang berubah **kapan ia boleh dijalankan**: sesudah sesi DDL menutup presisinya, atau sesudah pemilik proses menyatakan presisi `KTV-A` diterima apa adanya |

**Taksiran: kosong**, dan kekosongan itu **disengaja**.

**Blocked by:** 15 · 43 · **`KTV-A`**

**Dasar:**
```
EVIDENCED(SPEC-MODEL-DATA.md §10.22 - enam tabel acuan, bentuk seragam)
        DECIDED(ADR-D-TI-0042, ADR-D-TI-0043, INV-68)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] keenam tabel acuan terisi dari sumber lama, dan **persamaan cacah menutup** untuk keenamnya
- [ ] pengecualian migrasi **bernomor** dan dilaporkan **per baris**, bukan sebagai satu angka
- [ ] baris berkode ganda melaporkan **keduanya**, bukan membuang salah satunya
- [ ] uji positif kedua lulus: salah eja di sumber mendarat salah eja, dan **tidak** ditandai cacat
- [ ] pemindahan berjalan **hanya** di dalam jendela sakelar irisan `43`, dan jendelanya tercatat
- [ ] **sebelum dijalankan**, `ASUMSI-CLEAR.md` diperiksa — `KTV-A` bertenggat pada momen ini, dan `KTV-2` menuntut hal yang sama untuk tiket `08`
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md` sebagai asumsi **bertenggat** yang tiket ini menutup tenggatnya

## Treaty In - 45 - Daftar keadaan dan perpindahan sah berdiri, dan perpindahan di luar daftar ditolak

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-01` (pecahan keadaan) dan `P-37` · `ADR-D-TI-0055` · `ADR-D-TI-0046` · `INV-20`…`INV-23`.*

**What to build:** Kolom `KEADAAN_SIKLUS_HIDUP` hanya menerima **delapan** nilai sah, dan hanya **tiga belas**
perpindahan yang diterima. Setiap perpindahan lain ditolak — termasuk **setiap** perpindahan keluar
dari `DISETUJUI`, `DITOLAK`, dan `DIBATALKAN`.

Ini **mesinnya**. Tiket `14` sudah membuat kolomnya; tiket ini memberi kolom itu **arti**.

**Persyaratan:** `ADR-D-TI-0055` (delapan keadaan, tiga belas perpindahan) · `ADR-D-TI-0046` (satu keadaan, himpunan tertutup) · `INV-20`…`INV-23` · `ADR-D-TI-0038` (himpunan tertutup **boleh** `CHECK`, sebab ia tidak bertambah tanpa mengubah arti)

**Tidak termasuk:** **Siapa boleh memicu perpindahan mana** — itu wewenang, bukan daftar (`ADR-D-TI-0044`), dan ia
milik tiket `48`…`52`.
**Catatan persetujuan** yang menyertai tiap perpindahan — tiket `54`.
**Pembekuan terminal** — tiket `46`; daftar putih menolak *perpindahan*, tidak menjaga *nilai*.

**Jalur gagal:** Menyetel keadaan ke nilai di luar delapan -> **ditolak** · `DISETUJUI` -> `DRAFT` ->
**ditolak**, dan pesannya menyebut bahwa keadaan itu terminal · Melompati `DIAJUKAN` langsung ke
`DISETUJUI` -> ditolak.

**Uji:** **Negatif:** kedelapan keadaan diuji terhadap setiap perpindahan yang **tidak** ada di daftar.
**POSITIF — dan ia yang membuat tiket ini tidak boleh dipecah:** **ketiga belas** perpindahan sah
diuji satu per satu dan **seluruhnya diterima**.

> Daftar putih yang kurang **lolos setiap uji negatif** — uji negatif hanya memeriksa bahwa yang
> terlarang ditolak. Daftar yang baru memuat lima perpindahan menolak delapan yang sah, dan
> **satu-satunya yang menangkapnya adalah uji positif yang lengkap.**

**Menggantikan:** `TDA-08` — *rantai memendek lewat `RevisionState`, dan jejaknya dihapus.* Dan lebih
luas: sistem lama **tidak punya daftar perpindahan sama sekali**. Keadaan disetel lewat
`StatusAkseptasi` dari enam layar berbeda, termasuk satu yang menyetelnya ke nilai `"test"` —
nilai yang tidak punya padanan di antara keadaan sah mana pun (`ADR-D-TI-0054`).

**Blocked by:** `14`

**Dasar:**
```
EVIDENCED(TreatyInSetValue@ekspor-2026-09 - StatusAkseptasi="test" di enam layar)
        EVIDENCED(Akseptasi_DT@ekspor-2026-09 - cabang 1.4 pyDisabled, jalur kelompok tak terjangkau)
        DECIDED(ADR-D-TI-0055, ADR-D-TI-0046, ADR-D-TI-0038)
```

- [ ] kolom keadaan menerima **tepat delapan** nilai, dan menolak yang lain
- [ ] **tiga belas** perpindahan diterima; **uji positif ketiga belasnya lulus satu per satu**
- [ ] setiap perpindahan keluar dari `DISETUJUI`, `DITOLAK`, `DIBATALKAN` ditolak
- [ ] pesan penolakan menyebut **keadaan asal dan tujuan**, bukan hanya "tidak sah"

## Treaty In - 46 - Versi berkeadaan terminal beku, dan pembekuannya menjangkau seluruh entitas anaknya

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-38` · `INV-24` · `ADR-D-TI-0036`.*

**What to build:** Versi berkeadaan `DISETUJUI`, `DITOLAK`, atau `DIBATALKAN` **tidak dapat diubah nilainya
oleh siapa pun** — dan larangan itu menjangkau **seluruh entitas anaknya**, bukan hanya barisnya
sendiri.

Trigger `BEFORE INSERT OR UPDATE OR DELETE` pada tiap tabel anak.

**Persyaratan:** `INV-24` · `ADR-D-TI-0036` (angka dasar persetujuan dibekukan saat disetujui)

**Tidak termasuk:** **Daftar perpindahan** — tiket `45`. Tiket ini menjaga **nilai**, bukan perpindahan.
**Pembekuan jenis dan materialitas sejak `AJUKAN`** — tiket `07`, dan ia **lebih awal** daripada
pembekuan ini.

**Jalur gagal:** `UPDATE` pada versi `DISETUJUI` -> ditolak · **`INSERT` baris layer baru** pada versi
`DISETUJUI` -> **ditolak** · `DELETE` baris penyebaran pada versi `DITOLAK` -> ditolak.

**Uji:** **Negatif:** `UPDATE`, `INSERT`, dan `DELETE` diuji **pada tiap tabel anak**, bukan hanya
pada `VERSI_KONTRAK`.
**Positif:** versi `DRAFT` dan versi yang sedang menunggu persetujuan **tetap dapat diubah** —
pembekuan tidak boleh bocor ke keadaan non-terminal.

> **`UPDATE`/`DELETE` saja masih mengizinkan versi disetujui BERTAMBAH baris.** Itu sebab
> triggernya harus `BEFORE INSERT OR UPDATE OR DELETE`, dan sebab uji negatifnya wajib menyentuh
> `INSERT`.

**Menggantikan:** `TDA-07` penggantinya — *dua kontrol layar mengosongkan status akseptasi kontrak yang
sudah disetujui.* Dan `Force Edit (dev)` ditambah `Save EDM(dev)`: operator divisi IT dapat membuka
kunci addendum yang **sudah disetujui**, menyuntingnya, lalu menyimpannya di tempat — tanpa mesin
selisih dihitung ulang.

**Blocked by:** `45`

**Dasar:**
```
EVIDENCED(TreatyInForceEdit@ekspor-2026-09 - pyVisible=ALWAYS di atas gerbang divisi IT)
        EVIDENCED(SaveTreatyIn_EDM_Act@ekspor-2026-09 - tidak menghitung ulang selisih)
        DECIDED(INV-24, ADR-D-TI-0036)
```

- [ ] trigger `BEFORE INSERT OR UPDATE OR DELETE` terpasang pada **setiap** tabel anak versi
- [ ] uji negatif mencakup `INSERT`, bukan hanya `UPDATE`
- [ ] uji positif lulus: versi non-terminal tetap dapat diubah
- [ ] daftar tabel anak yang dijangkau **tertulis**, dan dicocokkan terhadap `KAMUS-KOLOM.md` — bukan disusun dari ingatan

## Treaty In - 47 - Paling banyak satu versi tak-terminal per kontrak

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-39` · `INV-25`.*

**What to build:** Sebuah kontrak hanya dapat punya **satu** versi yang belum selesai. Membuat versi kedua
yang belum selesai **ditolak**, dan pesannya menyebut versi mana yang masih terbuka.

**Persyaratan:** `INV-25` · `ADR-D-TI-0055`

**Tidak termasuk:** **Pembatalan draf** yang membebaskan slotnya — tiket `56`.
**Penomoran versi** — sudah dipegang `INV-04` di tiket `14`.

**Jalur gagal:** Membuat versi penyesuaian sementara versi sebelumnya masih `DRAFT` atau menunggu persetujuan -> **ditolak**, pesannya menyebut nomor versi yang masih terbuka.

**Uji:** **Negatif:** buat versi kedua saat yang pertama `DRAFT`; saat menunggu SH; saat menunggu DR.
**Positif:** sesudah versi pertama menjadi `DISETUJUI`, `DITOLAK`, **atau `DIBATALKAN`**, versi
berikutnya **diterima** — ketiga keadaan terminal membebaskan slotnya, dan menguji hanya
`DISETUJUI` akan melewatkan dua.

**Menggantikan:** `TDA-01` — *penjaga duplikat `TreatyInEdmCheckDuplicate` **mati**; tabrakan masuk cabang
`UPDATE`, menimpa addendum lain, dan melapor berhasil.*

**Blocked by:** `45`

**Dasar:**
```
EVIDENCED(TreatyInEdmCheckDuplicate@ekspor-2026-09 - MATI)
        DECIDED(INV-25, ADR-D-TI-0055)
```

- [ ] constraint menolak versi tak-terminal kedua, pesannya menyebut versi yang terbuka
- [ ] uji positif mencakup **ketiga** keadaan terminal, bukan hanya `DISETUJUI`

## Treaty In - 48 - Pengajuan menolak dua syarat dan memperingatkan enam, dan peringatannya tercatat

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-29` dan `P-30` · `INV-27` · `K1-1`…`K1-8` · `CO-6`…`CO-8`.*

**What to build:** **PK** mengajukan versi `DRAFT` untuk persetujuan. `K1-1` dan `K1-2` **menolak**
pengajuan, dengan pesan yang menyebut apa yang kurang. `K1-3` sampai `K1-8` menghasilkan
**peringatan yang tercatat** — belum penolakan.

**Persyaratan:** `INV-27` · `K1-1`…`K1-8` · `CO-6`…`CO-8` · `ADR-D-TI-0055` (`AJUKAN`)

**Tidak termasuk:** **Menaikkan `K1-3`…`K1-8` menjadi penolakan** — itu kemampuan **BARU** tersendiri, dan
`CARA MENYALAKANNYA` miliknya belum ditulis. Tiket ini hanya memasang mode peringatan.

**Jalur gagal:** Pengajuan tanpa syarat `K1-1` -> **ditolak**, pesannya menyebut apa yang kurang ·
Pengajuan yang melanggar `K1-5` -> **diterima**, dan peringatannya **tercatat dan dapat dibaca**.

**Uji:** **Negatif:** kedelapan syarat diuji satu per satu; dua menolak, enam meloloskan.
**Positif:** pengajuan yang memenuhi seluruh delapan **diterima tanpa satu peringatan pun** —
peringatan palsu sama merusaknya dengan penolakan palsu.

> **Mode peringatan tanpa pembaca adalah fitur yang dimatikan, ditambah biaya log.** Kriteria
> selesai menuntut peringatannya **dapat dibaca seseorang**, bukan hanya tertulis.

**Menggantikan:** **`TreatyInSubmitEDM`** — pada jalur addendum **kedelapan syarat `K1` MATI**: `CheckID`,
`CheckError`, dan keluar-bergalat seluruhnya dinonaktifkan. Addendum diajukan **tanpa satu pun
pemeriksaan**.

**Blocked by:** `45`

**Dasar:**
```
EVIDENCED(TreatyInSubmitEDM@ekspor-2026-09 - CheckID/CheckError/exit MATI)
        DECIDED(INV-27, ADR-D-TI-0055)
```

- [ ] `K1-1` dan `K1-2` menolak, pesannya menyebut apa yang kurang
- [ ] `K1-3`…`K1-8` menghasilkan peringatan **tersimpan**
- [ ] peringatan **dapat dibaca seseorang** — pembacanya disebut namanya, bukan "tersedia di log"
- [ ] uji positif lulus: pengajuan yang memenuhi delapan syarat **nol peringatan**

## Treaty In - 49 - Section head menyetujui versi yang menunggunya, dan versinya berpindah ke antrian dept head

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-31` · `ADR-D-TI-0055` (`SETUJUI`) · `ADR-D-TI-0044`.*

**What to build:** **SH** menyetujui versi yang berada di antriannya; versinya berpindah ke antrian **DH**,
dan keputusannya meninggalkan satu catatan persetujuan.

Pecahan **pertama dari tiga**. `U-2` melarang satu tiket melintasi lebih dari satu perpindahan, dan
`GRL-08` menulis `SETUJUI×3` — **tiga perpindahan berbeda**, bukan satu perpindahan berparameter.

**Persyaratan:** `ADR-D-TI-0055` (`SETUJUI` tingkat 1) · `ADR-D-TI-0044` (keadaan menjawab *"mungkinkah"*, peran menjawab *"boleh oleh siapa"*, dan **setiap larangan punya tepat satu sebab**)

**Tidak termasuk:** **Tingkat DH dan DR** — tiket `50` dan `51`.
**Larangan pengaju menyetujui sendiri** — tiket `52`; itu sebab yang **berbeda** dari wewenang
tingkat, dan `ADR-D-TI-0044` menuntut tiap penolakan punya tepat satu sebab.

**Jalur gagal:** Orang tanpa peran SH menyetujui -> ditolak **karena peran** · SH menyetujui versi yang
**tidak** di antriannya -> ditolak **karena keadaan**. Kedua pesan **berbeda**, dan itu uji
`ADR-D-TI-0044`.

**Uji:** **Negatif:** peran salah; keadaan salah; keduanya salah sekaligus — dan pesannya tetap
menyebut **satu** sebab.
**Positif:** SH yang sah menyetujui versi yang sah -> berpindah ke antrian DH, **dan catatannya
ada**.

**Menggantikan:** `TDA-09` — *peran dibaca dari `pyWorkBasketList(2)`, `pyTelephone`, dan nama orang yang
ditanam langsung di aturan.* Empat nama orang tersemat di penyaluran persetujuan yang **hidup sejak
2019**; orang baru tidak dapat menerima tugas sampai ruas teleponnya disetel, dan langkah itu tidak
tertulis di mana pun.

**Blocked by:** `48` · `54`

**Dasar:**
```
EVIDENCED(Akseptasi_DT@ekspor-2026-09 - nama orang sebagai tetapan, TreatyInSetValue 2.6 PositionUsername="BERNARD")
        DECIDED(ADR-D-TI-0055, ADR-D-TI-0044)
```

- [ ] SH menyetujui, versi berpindah ke antrian DH
- [ ] penolakan karena **peran** dan penolakan karena **keadaan** menghasilkan pesan yang **berbeda**
- [ ] satu catatan persetujuan tertulis, berisi pelaku, tingkat, dan waktunya
- [ ] **nol nama orang** di dalam aturan mana pun — peran, bukan nama

## Treaty In - 50 - Dept head menyetujui versi yang menunggunya, dan versinya berpindah ke antrian direktur

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-32` · `ADR-D-TI-0055` (`SETUJUI` tingkat 2).*

**What to build:** **DH** menyetujui versi di antriannya; versinya berpindah ke antrian **DR**, dan
keputusannya meninggalkan catatan.

Pecahan **kedua dari tiga**. Begitu `49` mendarat, tiket ini hampir gratis — dan itu alasan ia
**murah**, bukan alasan ia digabung.

**Persyaratan:** `ADR-D-TI-0055` (`SETUJUI` tingkat 2) · `ADR-D-TI-0044`

**Tidak termasuk:** **Tingkat SH dan DR** — tiket `49` dan `51`. **Wewenang bernilai** — tiket `58`.

**Jalur gagal:** DH menyetujui versi yang masih di antrian SH -> ditolak karena **keadaan** · SH mencoba menyetujui di tingkat DH -> ditolak karena **peran**.

**Uji:** **Negatif:** peran salah; keadaan salah.
**Positif:** DH yang sah -> antrian DR, catatan ada, dan **catatan tingkat SH tetap utuh** — tiket
ini tidak boleh menimpa jejak tingkat sebelumnya.

**Menggantikan:** `TDA-09`, dan `TDA-08` — *rantai dipendekkan `RevisionState = 1` menjadi dua tingkat, tanpa ambang nilai apa pun.*

**Blocked by:** `49`

**Dasar:**
```
EVIDENCED(Akseptasi_DT@ekspor-2026-09 - cabang 2 Admin->SecHead->selesai)
        DECIDED(ADR-D-TI-0055, ADR-D-TI-0044)
```

- [ ] DH menyetujui, versi berpindah ke antrian DR
- [ ] catatan tingkat sebelumnya **tetap utuh**
- [ ] penolakan peran dan penolakan keadaan berpesan berbeda

## Treaty In - 51 - Direktur menyetujui versi yang menunggunya, dan versinya menjadi disetujui bila seluruh K2 terpenuhi

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-33` · `ADR-D-TI-0055` (`SETUJUI` tingkat 3) · `K2-1`…`K2-6`.*

**What to build:** **DR** menyetujui versi di antriannya; versinya menjadi **`DISETUJUI`** — tetapi **hanya
bila seluruh `K2` terpenuhi**. Perpindahan ini membuat versinya **terminal**, dan karena itu ia
memicu pembekuan tiket `46`.

Pecahan **ketiga dari tiga**.

**Persyaratan:** `ADR-D-TI-0055` (`SETUJUI` tingkat 3) · `K2-1`…`K2-6` · `INV-24` (versinya menjadi terminal)

**Tidak termasuk:** **Pembekuan nilainya** — tiket `46` yang memasangnya; tiket ini hanya memicu keadaannya.

**Jalur gagal:** DR menyetujui sementara satu syarat `K2` tidak terpenuhi -> **ditolak**, pesannya menyebut
syarat mana · Sesudah `DISETUJUI`, mengubah nilai versinya -> ditolak oleh `46`.

**Uji:** **Negatif:** keenam syarat `K2` diuji satu per satu.
**Positif:** DR menyetujui versi yang memenuhi seluruh `K2` -> `DISETUJUI`, **dan pembekuan `46`
langsung berlaku** — diuji dengan mencoba mengubah satu nilai sesudahnya.

**Menggantikan:** **`TreatyInForceResolveComplete`** — *persetujuan dapat terjadi tanpa penyetuju.* Tombol
itu menyetel selesai-disetujui **tanpa satu pun penyetuju**, dan `TreatyInForceEdit` yang membukanya
ber-`pyVisible = ALWAYS`. Dan `TreatyInSetToDirector` **seluruh langkahnya mati** — tingkat direktur
di sistem lama tidak pernah benar-benar dijalankan lewat jalur itu.

**Blocked by:** `50`

**Dasar:**
```
EVIDENCED(TreatyInForceResolveComplete@ekspor-2026-09 - menyetel selesai tanpa penyetuju)
        EVIDENCED(TreatyInSetToDirector@ekspor-2026-09 - langkah 1-4 blok //, 4.1-4.7 ikut mati)
        DECIDED(ADR-D-TI-0055, INV-24)
```

- [ ] DR menyetujui -> `DISETUJUI`, hanya bila seluruh `K2` terpenuhi
- [ ] keenam `K2` diuji negatif satu per satu
- [ ] uji positif membuktikan **pembekuan `46` langsung berlaku** sesudahnya
- [ ] **tidak ada jalan** menyetel `DISETUJUI` selain lewat perpindahan ini

## Treaty In - 52 - Pengaju tidak menyetujui, tiap tingkat orang berbeda, dan tidak ada jalan pintas

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-36` dan `P-40` · `ADR-D-TI-0052` · `ADR-D-TI-0044`.*

**What to build:** Pemberi keputusan yang **sama dengan pengisi kontraknya**, atau **sama dengan pemberi
keputusan tingkat sebelumnya**, **ditolak**. Dan tidak ada jalur yang melompati tingkat.

**Persyaratan:** `ADR-D-TI-0052` (satu rantai, tanpa pengecualian) · `ADR-D-TI-0044` (tepat satu sebab per larangan)

**Tidak termasuk:** **Wewenang menurut nilai kontrak** — tiket `58`. Itu sebab ketiga yang berbeda lagi.

**Jalur gagal:** Pengaju menyetujui di tingkat mana pun -> ditolak, pesannya menyebut **ia pengajunya** ·
Orang yang sama menyetujui di dua tingkat berurutan -> ditolak, pesannya menyebut **tingkat mana**
ia sudah memutuskan.

**Uji:** **Negatif:** pengaju = SH; SH = DH; DH = DR; dan pengaju = DR.
**Positif:** **empat orang berbeda** melewati rantai penuh **diterima** — dan rantai tiga orang
dengan pengaju keempat juga diterima.

**Menggantikan:** `TDA-08` dan **`RevisionState == 1`** — *jalur revisi memendekkan persetujuan menjadi satu
tingkat; kepala departemen dan direktur dilewati, **tanpa ambang nilai apa pun**.*

**Blocked by:** `51`

**Dasar:**
```
EVIDENCED(Akseptasi_DT@ekspor-2026-09 - RevisionState=1 memotong rantai jadi dua tingkat)
        DECIDED(ADR-D-TI-0052, ADR-D-TI-0044)
```

- [ ] pengaju ditolak di setiap tingkat, pesannya menyebut sebab itu saja
- [ ] orang yang sama di dua tingkat berurutan ditolak
- [ ] uji positif rantai empat orang lulus
- [ ] **tidak ada** jalur yang melompati tingkat — diperiksa terhadap ketiga belas perpindahan

## Treaty In - 53 - Section head, dept head, atau direktur mengembalikan versi ke draft, dan pengembalian tanpa alasan ditolak

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-34` · `ADR-D-TI-0055` (`KEMBALIKAN`).*

**What to build:** **SH / DH / DR** mengembalikan versi yang menunggunya ke `DRAFT`. Pengembalian **tanpa
alasan ditolak**, dan alasannya tersimpan di catatan persetujuan.

Perpindahan ini yang membuat pembekuan jenis dan materialitas (tiket `07`) **tidak menjadi
penguncian permanen**.

**Persyaratan:** `ADR-D-TI-0055` (`KEMBALIKAN` ×3) · `INV-27`

**Tidak termasuk:** **Penolakan** — tiket `55`. Mengembalikan dan menolak adalah dua perpindahan berbeda dengan akibat berbeda.

**Jalur gagal:** Pengembalian tanpa alasan -> ditolak · Pengembalian oleh orang yang bukan pemegang antriannya -> ditolak karena peran.

**Uji:** **Negatif:** alasan kosong; alasan hanya spasi; pelaku salah.
**Positif:** sesudah dikembalikan, versinya **dapat disunting lagi** termasuk jenis dan
materialitasnya — itu yang membuktikan tiket `07` tidak mengunci permanen.

**Menggantikan:** `TDA-02` dan bentuk umumnya: sistem lama **tidak punya perpindahan mengembalikan** untuk
addendum sama sekali. Satu-satunya jalan keluar dari pengajuan adalah **ditolak**, dan penolakan
**menghapus barisnya**.

**Blocked by:** `49`

**Dasar:**
```
EVIDENCED(TreatyInDeclineConfirmation_postactEDM@ekspor-2026-09 - langkah 1-4 MATI, dua RDB remove hidup)
        DECIDED(ADR-D-TI-0055)
```

- [ ] ketiga tingkat dapat mengembalikan; alasan **wajib**
- [ ] alasannya tersimpan di catatan persetujuan
- [ ] uji positif: sesudah kembali ke `DRAFT`, jenis dan materialitas dapat diubah lagi

## Treaty In - 54 - Tiap perpindahan keadaan meninggalkan satu catatan persetujuan, dan tidak ada jalur yang melewatinya

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-46` · `ADR-D-TI-0045` · `SPEC-MODEL-DATA.md` §10.20.*

**What to build:** Setiap perpindahan keadaan menulis **satu** baris `CATATAN_PERSETUJUAN` berisi pelaku,
tingkat, waktu, dan alasannya bila ada. **Tidak ada jalur yang melewatinya** — ia fakta mesin.

**PEMBUAT PERTAMA** untuk `CATATAN_PERSETUJUAN`.

**Persyaratan:** `ADR-D-TI-0045` (jejak sebagai fakta mesin, tidak dapat dilewati) · `SPEC-MODEL-DATA.md` §10.20 · `INV-26`

**Tidak termasuk:** **`JEJAK_PERUBAHAN`** — entitas berbeda, sudah dibangun tiket `39`. Yang satu mencatat
**perpindahan keadaan**, yang lain mencatat **perubahan fakta**.
**`PERISTIWA_KONTRAK`** — lihat `PENGHALANG` di bawah; ia **tidak punya kemampuan** dan tidak
dibangun di sini.

**Jalur gagal:** Perpindahan yang berhasil tanpa catatan -> **mustahil**; bila mungkin, tiket ini belum
selesai · Catatan tanpa pelaku -> ditolak.

**Uji:** **Negatif:** coba lakukan tiap perpindahan lewat jalur yang melewati penulisan catatan.
**Positif:** **ketiga belas** perpindahan diuji, dan **ketiga belasnya** menghasilkan tepat satu
catatan — bukan nol, bukan dua.

**Menggantikan:** **`CommentList` sistem lama** hanya punya `OperatorName`, `IsApproved`, `Suggest`, dan
`Date` — **tidak merekam tingkat penyetuju sama sekali**. Dan zona waktunya campuran: komentar
persetujuan memakai jam Pega (GMT), *"Create Revision"* memakai jam Oracle, dan langkah jam Oracle
di `AddCommentList_Act` langkah 1 **mati**.

**Blocked by:** `45`

**Dasar:**
```
EVIDENCED(AddCommentList_Act@ekspor-2026-09 - langkah 1 MATI, CommentList tanpa ruas tingkat)
        DECIDED(ADR-D-TI-0045, INV-26)
```

- [ ] `CATATAN_PERSETUJUAN` berdiri sesuai `KAMUS-KOLOM.md`, **termasuk ruas tingkat** yang sistem lama tidak punya
- [ ] ketiga belas perpindahan menghasilkan tepat **satu** catatan masing-masing
- [ ] **tidak ada jalur** yang dapat berpindah tanpa menulis catatan — diuji, bukan diargumentasikan
- [ ] waktu disimpan dalam **satu** zona yang dinyatakan, bukan campuran

## Treaty In - 55 - Penolakan versi — kontrak maupun addendum — meninggalkan catatan, dan barisnya tidak dihapus

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-35` dan `P-58` · `ADR-D-TI-0055` (`TOLAK`).*

**What to build:** **SH / DH / DR** menolak versi. Penolakan tanpa alasan **ditolak**. Versinya menjadi
`DITOLAK` — **terminal** — dan **barisnya tetap tersimpan**, nomornya **tidak dipakai ulang**.

**Persyaratan:** `ADR-D-TI-0055` (`TOLAK`, `DITOLAK` terminal) · `INV-04` (nomor tidak dipakai ulang) · `INV-24`

**Tidak termasuk:** **Pengembalian ke `DRAFT`** — tiket `53`. Menolak dan mengembalikan berakibat berbeda: yang satu terminal, yang lain tidak.

**Jalur gagal:** Penolakan tanpa alasan -> ditolak · Sesudah `DITOLAK`, mengubah versinya -> ditolak oleh
`46` · Nomor versi yang ditolak dipakai ulang -> ditolak oleh `INV-04`.

**Uji:** **Negatif:** alasan kosong; mencoba memakai ulang nomor versi yang ditolak.
**Positif:** versi yang ditolak **masih dapat dibaca** — beserta alasannya, pelakunya, dan
waktunya. Itu yang membuktikan barisnya tidak dihapus.

**Menggantikan:** `TDA-02` — ***menolak addendum MENGHAPUS barisnya.*** Dan kode pencatatannya **ada dan
dimatikan**: `TreatyInDeclineConfirmation_postactEDM` langkah 1–4 mati, dua langkah `RDB remove`
hidup. Akibatnya nomor revisi mungkin dipakai ulang, dan **tidak ada yang pernah dapat menjawab
berapa kali sebuah kontrak gagal diubah**.

**Blocked by:** `49` · `54`

**Dasar:**
```
EVIDENCED(TreatyInDeclineConfirmation_postactEDM@ekspor-2026-09 - RDB remove hidup, pencatatan MATI)
        DECIDED(ADR-D-TI-0055, INV-04, INV-24)
```

- [ ] penolakan tanpa alasan ditolak
- [ ] versi `DITOLAK` **tetap tersimpan** dan dapat dibaca beserta alasannya
- [ ] nomornya **tidak dapat** dipakai ulang
- [ ] berlaku sama untuk versi kontrak **dan** versi addendum — diuji keduanya

## Treaty In - 56 - Pengisi kontrak membatalkan draf yang ia buat sendiri, dan kontraknya langsung terbuka untuk versi berikutnya

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-57` · `ADR-D-TI-0055` perubahan 24 Sep 2026 · `INV-04`.*

**What to build:** **PK** membatalkan versi `DRAFT` yang **ia buat sendiri**. Versinya menjadi `DIBATALKAN` —
terminal — **barisnya tetap tersimpan**, nomornya tidak dipakai ulang, dan kontraknya **langsung
terbuka** untuk versi berikutnya.

**Persyaratan:** `ADR-D-TI-0055` (`DRAFT → DIBATALKAN`) · `INV-04` · `INV-25`

**Tidak termasuk:** **Pembatalan oleh orang lain** — tidak ada, dan itu keputusan: hanya pembuatnya.

**Jalur gagal:** Orang lain membatalkan draf -> ditolak · Membatalkan versi yang sudah diajukan -> ditolak;
jalurnya `KEMBALIKAN` lalu batal · Nomor versi yang dibatalkan dipakai ulang -> ditolak.

**Uji:** **Negatif:** pelaku bukan pembuat; versi bukan `DRAFT`; pakai ulang nomornya.
**Positif:** sesudah pembatalan, versi berikutnya **langsung diterima** — itu yang membuktikan
`INV-25` membebaskan slotnya.

**Menggantikan:** **Tidak ada padanannya — ini kemampuan BARU.** Sistem lama **tidak punya cara membuang
draf** yang tidak merusak apa pun: satu-satunya jalan adalah mengajukannya supaya ditolak, dan
penolakan **menghapus barisnya** (`TDA-02`).

**Blocked by:** `45` · `54`

**Dasar:**
```
EVIDENCED(ADR-D-TI-0055 perubahan 24 Sep 2026 - L-7 ditutup)
        DECIDED(ADR-D-TI-0055, INV-04, INV-25)
```

- [ ] hanya pembuatnya dapat membatalkan, dan hanya saat `DRAFT`
- [ ] baris tetap tersimpan, nomor tidak dipakai ulang
- [ ] uji positif: versi berikutnya langsung diterima sesudahnya
- [ ] **CARA MENYALAKANNYA** ditulis — ini golongan BARU

## Treaty In - 57 - Angka rupiah pada versi yang sudah disetujui tidak bergeser, hari ini maupun tahun depan

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-27` · `ADR-D-TI-0036` · `INV-43`.*

**What to build:** Kurs **dibekukan pada versinya** saat disetujui, dan **tidak dibaca ulang** saat
ditampilkan. Angka rupiah sebuah versi `DISETUJUI` sama hari ini dan tahun depan.

**Persyaratan:** `ADR-D-TI-0036` (angka dasar persetujuan beku; hasil beku membawa penunjuk ke masukan yang dipakai) · `INV-43`

**Tidak termasuk:** **Kurs pada versi `DRAFT`** — ia memang boleh mengikuti kurs berjalan sampai disetujui.
**Perhitungan ulang** — `ADR-D-TI-0043`: itu peristiwa bisnis tersendiri, bukan bagian tiket ini.

**Jalur gagal:** Mengubah baris kurs tahunan lalu membuka versi `DISETUJUI` -> angkanya **tidak berubah** ·
Bila berubah, tiket ini belum selesai.

**Uji:** **Negatif:** ubah kurs master, buka versi disetujui, bandingkan — **harus identik**.
**Positif:** versi `DRAFT` **mengikuti** kurs terbaru; pembekuan tidak boleh bocor ke draf.

**Menggantikan:** **Angka yang sudah disetujui dapat bergeser sendiri** — salah satu dari empat sebabnya
**kurs tahun berjalan**. Sistem lama membaca kurs saat menampilkan, sehingga nilai rupiah sebuah
kontrak yang sudah disetujui **berubah setiap ganti tahun** tanpa ada yang menyentuhnya.

**Blocked by:** `46` · `20`

**Dasar:**
```
EVIDENCED(TREATYEXCHANGEYEARLY@ekspor-2026-09 - kurs per tahun dibaca saat tampil)
        DECIDED(ADR-D-TI-0036, INV-43)
```

- [ ] kurs tersimpan pada versi saat disetujui, beserta tanggal dan sumbernya
- [ ] uji negatif kurs-master-berubah lulus: angka versi disetujui **identik**
- [ ] uji positif lulus: versi `DRAFT` tetap mengikuti kurs terbaru

## Treaty In - 58 - Wewenang persetujuan dibatasi nilai kontrak menurut batas yang berlaku

---
status: tertahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-41` · `ADR-D-TI-0044` · eskalasi manajemen butir 5.*

**What to build:** Penyetuju yang wewenangnya **di bawah nilai kontrak** ditolak, dan pesannya menyebut
batasnya.

**Persyaratan:** `ADR-D-TI-0044` (peran menjawab *"boleh oleh siapa"*) · dokumen batas wewenang yang berlaku

**Tidak termasuk:** **Pemisahan pelaku** — tiket `52`; itu sebab yang berbeda, dan `ADR-D-TI-0044` menuntut tiap penolakan punya tepat satu sebab.

**Jalur gagal:** Penyetuju berwenang di bawah nilai kontrak -> ditolak, pesannya menyebut **batasnya**, bukan hanya "tidak berwenang".

**Uji:** **Negatif:** nilai tepat di atas batas tiap tingkat. **Positif:** nilai **tepat pada** batas diterima — batas inklusif atau eksklusif adalah bagian dari jawaban yang ditunggu.

**Menggantikan:** Sistem lama **tidak punya batas wewenang sama sekali** — sapuan menemukan **nol nama hak
akses terisi**; mekanisme peran bawaan Pega hadir 2.202 kali sebagai tempat kosong dan tidak pernah
dipakai sekali pun (`ADR-D-TI-0044` Konteks).

**Blocked by:** `51` · **dokumen batas wewenang**

**Dasar:**
```
EVIDENCED(sapuan hak akses@ekspor-2026-09 - 2.202 tempat kosong, nol terisi)
        DECIDED(ADR-D-TI-0044)
```

- [ ] batas wewenang dibaca dari **tabel acuan**, bukan ditanam di kode
- [ ] penolakan menyebut batasnya
- [ ] **PENGHALANG:** dokumen batas wewenang yang berlaku · **siapa menjawab:** manajemen, eskalasi butir 5 · **yang berubah:** angka ambangnya, dan apakah batasnya inklusif

## Treaty In - 59 - Kepala kontrak warisan dan seluruh daftar anaknya dipindahkan apa adanya, beserta keadaannya

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-50` (pecahan selain tabel acuan) · `ADR-D-TI-0042` · `ADR-D-TI-0043`.*

**What to build:** **PM** memindahkan kepala kontrak lama dan seluruh daftar anaknya **apa adanya, tanpa
menghitung ulang apa pun**, dan menulis keadaannya ke kolom yang sudah punya himpunan nilai sah.

**Persyaratan:** `ADR-D-TI-0042` (pindah apa adanya, penanda warisan, sentuh-perbaiki) · `ADR-D-TI-0043` (setiap angka identik) · `KTV-A`

**Tidak termasuk:** **Keadaan yang tidak berpadanan** — tiket `60`.
**Isi tabel acuan** — sudah dipindahkan tiket `44`.
**Menghitung ulang** apa pun — `ADR-D-TI-0043` melarangnya; itu peristiwa bisnis tersendiri.

**Jalur gagal:** Satu angka berbeda antara sumber dan hasil -> **migrasi gagal dan diulang**; bukan
"dalam toleransi" · Jumlah baris anak per kontrak tidak cocok -> gagal, dan **ini kegagalan yang
paling sering lolos ketiga ukuran pertama**.

**Uji:** **Negatif:** kontrak yang kehilangan satu baris layer **harus** terdeteksi.
**Positif:** kontrak yang pindah utuh lolos **keempat** ukuran `ADR-D-TI-0043` — cacah kontrak, pasangan
ID dua arah, nilai kunci identik, **dan jumlah baris anak per kontrak**.

**Menggantikan:** Tidak ada perilaku lama yang digantikan — ini **pemindahan**. Yang digantikan adalah
ketiadaannya: data lama tinggal di `M_TREATY_IN.JSONDATA`, **dua kolom saja**, tanpa kolom waktu.

**Blocked by:** `45` · `44`

**Dasar:**
```
EVIDENCED(M_TREATY_IN DDL@ekspor-2026-09 - dua kolom, ID dan JSONDATA, tanpa kolom waktu)
        DECIDED(ADR-D-TI-0042, ADR-D-TI-0043)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] kepala dan seluruh daftar anak pindah; **jumlah baris anak per kontrak dicocokkan**
- [ ] baris hasil migrasi membawa **penanda warisan** yang terbaca siapa pun
- [ ] nol angka dihitung ulang — diperiksa, bukan diandaikan
- [ ] **bertenggat `KTV-A`:** presisi kolom **tidak dapat dipersempit** sesudah tiket ini berjalan

## Treaty In - 60 - Keadaan warisan yang tidak punya padanan mendarat di warisan-tak-terpetakan, dengan nilai aslinya tersimpan

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-51` · `ADR-D-TI-0054` · `INV-21`, `INV-26`.*

**What to build:** Nilai keadaan lama yang tidak punya padanan — seperti `"test"` — mendarat di
**`WARISAN_TAK_TERPETAKAN`**, dan **nilai aslinya tersimpan** di `KEADAAN_WARISAN_ASLI`.

**Persyaratan:** `ADR-D-TI-0054` · `INV-21` · `INV-26`

**Tidak termasuk:** **Perbaikannya** — tiket `61`. Tiket ini hanya menempatkan; memperbaiki adalah perbuatan orang.

**Jalur gagal:** Nilai liar dipetakan diam-diam ke keadaan sah -> **cacat**; itu membersihkan sejarah ·
`KEADAAN_WARISAN_ASLI` terisi pada baris yang **bukan** `WARISAN_TAK_TERPETAKAN` -> ditolak.

**Uji:** **Negatif:** nilai `"test"` dipetakan ke `DRAFT` -> harus tidak terjadi.
**Positif:** baris ber-`WARISAN_TAK_TERPETAKAN` **tidak dapat** memasuki alur kerja biasa — tidak
ada perpindahan keluar selain `PERBAIKAN_WARISAN`.

**Menggantikan:** **`TreatyInSetValue` menyetel `StatusAkseptasi = "test"`**, dan aturan itu terpasang di
**enam layar**. Nilai itu tidak punya padanan di antara keadaan sah mana pun.

**Blocked by:** `59`

**Dasar:**
```
EVIDENCED(TreatyInSetValue@ekspor-2026-09 - StatusAkseptasi="test" di enam layar)
        DECIDED(ADR-D-TI-0054, INV-21, INV-26)
```

- [ ] nilai tanpa padanan mendarat di `WARISAN_TAK_TERPETAKAN`
- [ ] `KEADAAN_WARISAN_ASLI` terisi **hanya** di baris itu
- [ ] nilai aslinya **tidak pernah dibaca** perhitungan mana pun
- [ ] uji positif: tidak ada perpindahan keluar selain `PERBAIKAN_WARISAN`

## Treaty In - 61 - Baris warisan tak terpetakan diperbaiki ke keadaan sah yang dipilih secara eksplisit

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-52` · `ADR-D-TI-0054` · `ADR-D-TI-0045`.*

**What to build:** **PK** memindahkan baris ber-`WARISAN_TAK_TERPETAKAN` ke keadaan sah **yang dipilihnya
secara eksplisit** — bukan ditebak sistem — dan pilihannya tercatat beserta siapa dan kapan.
Sesudahnya baris itu berjalan seperti baris lain.

**Persyaratan:** `ADR-D-TI-0054` (`PERBAIKAN_WARISAN`) · `ADR-D-TI-0045` (tercatat di jejak) · `ADR-D-TI-0042` (sentuh-perbaiki: penanda warisannya dicabut)

**Tidak termasuk:** **Penebakan otomatis** — dilarang `ADR-D-TI-0054`. Tidak ada jalur yang memilihkan keadaannya.

**Jalur gagal:** Perbaikan tanpa memilih keadaan -> ditolak · Perbaikan yang tidak meninggalkan jejak ->
mustahil · Baris yang diperbaiki tetap membawa penanda warisan -> cacat; `ADR-D-TI-0042` mencabutnya.

**Uji:** **Negatif:** perbaikan tanpa pilihan; perbaikan ke keadaan di luar delapan.
**Positif — dan ia yang membuktikan tiket ini punya kriteria sendiri:** sesudah diperbaiki, baris
itu **dapat menjalani perpindahan biasa** dan **memenuhi invarian yang berlaku sekarang**. Itu
perilaku yang tidak dimiliki tiket `60`, dan sebab keduanya tidak digabung.

**Menggantikan:** Tidak ada — sistem lama tidak punya keadaan tak-terpetakan, sebab ia tidak punya himpunan tertutup.

**Blocked by:** `60`

**Dasar:**
```
DECIDED(ADR-D-TI-0054, ADR-D-TI-0045, ADR-D-TI-0042)
```

- [ ] keadaan dipilih **orang**, eksplisit, dari delapan yang sah
- [ ] pilihannya tercatat di jejak beserta pelaku dan waktu
- [ ] penanda warisan **dicabut** sesudah perbaikan
- [ ] uji positif: baris yang diperbaiki menjalani perpindahan biasa dan memenuhi invarian sekarang

## Treaty In - 62 - Kontrak dicari menurut cedant, asal bisnis, periode, sifat proporsi, dan keadaan siklus hidupnya

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-56` · `SPEC-MODEL-DATA.md` §10.1 dan §10.2 · `ADR-D-TI-0055`.*

**What to build:** **PK** mencari kontrak yang ada menurut lima sumbu — termasuk **keadaan siklus hidupnya** —
dan melihat hasilnya sebagai daftar.

**Persyaratan:** `SPEC-MODEL-DATA.md` §10.1, §10.2 · `ADR-D-TI-0055` (himpunan keadaan yang dapat dicari)

**Tidak termasuk:** **Bentuk layarnya** — `L-4`: tidak ada spesifikasi layar di mana pun, dan ia milik batch
lapisan aplikasi. Yang dibangun di sini **perilaku pencariannya**, bukan tampilannya.

**Jalur gagal:** Mencari menurut keadaan yang tidak ada di daftar -> mengembalikan kosong, bukan galat · Kontrak tanpa versi berlaku -> tetap ditemukan, dengan keadaannya apa adanya.

**Uji:** **Negatif:** keadaan di luar delapan.
**Positif:** kontrak ber-`WARISAN_TAK_TERPETAKAN` **ditemukan** — bila ia tidak muncul di pencarian,
tidak akan ada yang tahu ia perlu diperbaiki.

**Menggantikan:** `TDA-11` — *daftar pilihan menyatukan kontrak dan addendum tanpa pembeda jenis dan
**tanpa saringan keadaan sama sekali***: `TreatyLoadMasterJoinEdm` meng-UNION `treaty_in` dengan
`treaty_in_edm` **tanpa `WHERE`**.

**Blocked by:** `45`

**Dasar:**
```
EVIDENCED(TreatyLoadMasterJoinEdm@ekspor-2026-09 - UNION tanpa WHERE, tanpa pembeda jenis)
        DECIDED(ADR-D-TI-0055)
```

- [ ] kelima sumbu pencarian bekerja, termasuk keadaan
- [ ] uji positif: baris `WARISAN_TAK_TERPETAKAN` ikut ditemukan
- [ ] hasil pencarian **tidak** mencampur kontrak dengan versi tanpa pembeda

## Treaty In - 63 - Cara pembukuan XOL dicatat, dan hanya pada kontrak non-proporsional

---
status: tertahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-20` · `INV-34` · `SPEC-MODEL-DATA.md` §14.4.*

**What to build:** **PK** mencatat cara pembukuan XOL pada versi, dan `INV-34` menolaknya pada kontrak
**proporsional**.

**Persyaratan:** `INV-34` · `SPEC-MODEL-DATA.md` §14.4

**Tidak termasuk:** **Perilaku pembukuannya sendiri** — di luar gelombang ini; yang dibangun hanya pencatatan dan penolakannya.

**Jalur gagal:** Mencatat cara pembukuan XOL pada kontrak proporsional -> **ditolak**, pesannya menyebut sifat proporsi kontraknya.

**Uji:** **Negatif:** kontrak proporsional. **Positif:** kontrak non-proporsional menerima seluruh nilai sahnya — dan daftar nilai sahnya **belum diketahui** sampai `Uji X-2` kembali.

**Menggantikan:** Tidak ada cacat yang digantikan — `AccountingMode` tersimpan di sistem lama dan **belum diketahui apakah dipakai**.

**Blocked by:** `14` · **`Uji X-2`**

**Dasar:**
```
DECIDED(INV-34)
        DIASUMSIKAN-CLEAR(Uji X-2)
```

- [ ] `INV-34` menolak cara pembukuan XOL pada kontrak proporsional
- [ ] **PENGHALANG:** hasil `Uji X-2` atas `AccountingMode` · **siapa menjawab:** DBA, sesudah izin kueri baca-saja turun · **yang berubah:** apakah kolomnya **dipakai sama sekali**, dan himpunan nilai sahnya

## Treaty In - 64 - Tabel acuan pembagian kapasitas punya pemilik yang bertanggung jawab atas isinya

---
status: tertahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-49` · `ADR-D-TI-0050` · eskalasi manajemen butir 7.*

**What to build:** Tabel acuan pembagian kapasitas membawa **pemiliknya**, dan setiap perubahan isinya
berjejak beserta siapa dan kapan.

**Persyaratan:** `ADR-D-TI-0050` (sumber luar yang tidak dipercaya) · `ADR-D-TI-0045`

**Tidak termasuk:** **Isi tabelnya** — ia data, bukan bentuk; pemiliknya yang mengisinya.

**Jalur gagal:** Perubahan isi tabel acuan tanpa jejak -> mustahil · Tabel tanpa pemilik -> tiket ini tidak dapat dinyatakan selesai.

**Uji:** **Negatif:** ubah isi tanpa pelaku. **Positif:** perubahan oleh pemilik yang sah tercatat lengkap.

**Menggantikan:** **`PROPORTIONALARRG`** — master susunan retro sistem lama: **seluruh kolom
`VARCHAR2(1000)`, tanpa kunci utama maupun indeks**, dan tanpa pemilik yang tercatat di mana pun.

**Blocked by:** `15` · **penetapan pemilik**

**Dasar:**
```
EVIDENCED(PROPORTIONALARRG DDL@ekspor-2026-09 - seluruh kolom VARCHAR2(1000), nol PK, nol indeks)
        DECIDED(ADR-D-TI-0050, ADR-D-TI-0045)
```

- [ ] tabel acuan kapasitas membawa pemiliknya
- [ ] perubahan isinya berjejak
- [ ] **PENGHALANG:** penetapan pemilik · **siapa menjawab:** manajemen, eskalasi butir 7 · **yang berubah:** siapa yang bertanggung jawab, bukan bentuk tabelnya

## Treaty In - REA - Papan tiket — Treaty In

**Tiket: aktif 46 · tertahan 5 · selesai 0 · mati 0** — **51 tiket**, bernomor `14`…`64`

**Yang sebenarnya dapat dimulai: 2** — `14`, `15`. *Dihitung ulang atas **seluruh papan**, bukan per batch.*

> ## ANGKA LAMA *"dapat dimulai: 4"* SALAH, DAN SEBABNYA LAYAK DIBACA
>
> Papan ini pernah mencetak **"dapat dimulai: 4 — `01`, `02`, `03`, `05`"**. **Keempatnya tidak dapat
> dimulai.** Ketiganya menambahkan **kolom** pada `VERSI_KONTRAK` dan satu menambahkan tabel yang
> merujuknya — dan **tidak ada satu pun tiket di papan yang membuat tabelnya**.
>
> **Sebab kekeliruannya, dan inilah yang akan terulang:**
>
> > Angka itu dihitung atas **penghalang di dalam papan**, sementara penghalangnya ada **di luar** —
> > di lingkup yang ronde itu belum digarap.
>
> Selama sebuah batch hanya memuat sebagian kemampuan sebuah modul, **"tidak ada penghalang di papan"
> tidak berarti "tidak ada penghalang"**. Papan yang tidak lengkap **selalu** melaporkan angka yang
> terlalu besar, dan ia melaporkannya dengan tenang.
>
> **Cara memeriksanya, sebelum mencetak angka itu di batch berikutnya:** untuk tiap tiket yang mengaku
> dapat dimulai, sebutkan **entitas yang disentuhnya** dan tunjukkan **tiket mana yang membuatnya**.
> Bila jawabannya "tidak ada tiket, tabelnya sudah ada", itu **andaian** — dan andaian itulah yang
> keliru di sini.

> **Lingkup papan ini: kemampuan Treaty In saja** — batch 1 (`14`…`44`) dan batch 2 (`45`…`64`).
>
> ### PEMISAHAN PAPAN — 25 September 2026
>
> Tiket **`01`…`13`** (kemampuan Adjustment) **dipindahkan** ke
> [`../../../treaty-in-adjustment/5-tiket/issues/`](../../../treaty-in-adjustment/5-tiket/issues/)
> atas perintah pemilik proses. **Ini membalik `GRL-01`** — *"satu model, satu spesifikasi, satu
> papan"* — dan pembalikannya dicatat, bukan disamarkan.
>
> **Yang dibalik hanya papannya.** Satu model data, satu daftar invarian, satu
> `DAFTAR-PEKERJAAN.md`, dan satu `ASUMSI-CLEAR.md` **tetap**. `P-60`…`P-66` **tidak ikut pindah**;
> ia tetap di `DAFTAR-PEKERJAAN.md` papan ini.
>
> **Penomoran tidak disusun ulang.** `14`…`64` tetap, `01`…`13` tetap — tidak ada tabrakan nomor
> antar-papan, dan setiap rujukan yang sudah ada tetap sah.
> Dari **66 kemampuan** di [`../DAFTAR-PEKERJAAN.md`](../DAFTAR-PEKERJAAN.md). **Batch 2 — 22 butir
> yang belum ditiketkan** (sapuan mekanis atas medan `*Asal:*`, bukan hitungan tangan).
>
> ### RALAT 25 September 2026 — penahan batch 2 BUKAN `REV-3`
>
> Baris ini pernah berbunyi: *"20 butir yang menyentuh daftar keadaan — belum ditiketkan;
> **penahannya tanggapan `REV-3`**."* **Itu keliru, dan kekeliruannya menyebar ke tiga berkas.**
>
> `REV-3` menyatakan sebaliknya di barisnya sendiri: keputusan ADR-D-TI-0055 — *"tidak ada satu pun jalan
> menyetel keadaan selain melalui perpindahan di daftar"* — **tetap, seluruhnya**. Yang diusulkan
> berubah hanya **tabel §4**, dan keempat butirnya **koreksi atas pemerian sistem lama**. Tidak satu
> pun mengubah daftar keadaan, kolom keadaan, maupun perpindahannya. `KTV-4` sudah mengadilinya:
> ***"nol dari enam menentukan letak kolom."***
>
> | Salah | Benar |
> |---|---|
> | *"batch 2 ditunda, menunggu `REV-3`"* | *"batch 2 dikerjakan sesi berikutnya, sebab 51 kemampuan tidak muat satu sesi"* |
>
> Yang pertama membuat orang **menunggu jawaban yang tidak akan mengubah apa pun**; yang kedua
> membuat mereka **menjadwalkan sesi**. **Pemisahan batchnya tidak dicabut** — ia benar sebagai
> ukuran sesi. `REV` tetap diserahkan: ia **utang dokumentasi ADR induk**, bukan utang tiket.

#### Papan

##### Tiket `01`…`13` TIDAK lagi di papan ini

Ketiga belasnya pindah ke papan Adjustment. **Satu tepi tetap melintas ke sini dan tidak boleh
hilang:**

| Tiket papan ini | Diblokir | Yang diantarkan tiket itu |
|---|---|---|
| **`40`** — nilai versi sebelumnya berdampingan, di-SELECT bukan disalin | **`01`** ⟦papan Adjustment⟧ | rujukan versi dasar yang disimpan eksplisit |

Dan **lima tepi ke arah sebaliknya**: `01`, `02`, `03`, `04`, `05` di papan Adjustment seluruhnya
diblokir **`14`** di papan ini. Itu sebab papan Adjustment mencetak **"dapat dimulai: 0"**.

##### Batch 1 Treaty In — `14`…`44`

| # | Tiket | Ditahan oleh |
|---|---|---|
| `14` | **Kontrak dan versi pertamanya berdiri**, dan dapat ditemukan kembali dengan pengenalnya | — |
| `15` | **Himpunan acuan bertambah** tanpa mengubah arti apa pun yang sudah tercatat | — |
| `16` | Peringatan kunci alami ganda tanpa menghalangi simpan | `14` |
| `17` | Kontrak ditemukan kembali lewat nomor warisan | `14` |
| `18` | Kunci alami kontrak beku sesudah kontraknya lahir | `14` |
| `19` | Lapisan beku tidak berubah antar versi | `14` |
| `20` | Lebih dari satu mata uang berlaku, masing-masing berkurs dan berperiode | `14` · `15` |
| `21` | Kegagalan hitung menghasilkan keterangan, bukan nol | `20` |
| `22` | Retensi cedant per kelompok treaty per mata uang | `14` · `15` |
| `23` | EGNPI per kelompok treaty per mata uang | `14` · `15` |
| `24` | Portofolio masuk dan keluar yang menyertai kontrak | `14` |
| `25` | Periode pelaporan beserta jatuh temponya | `14` |
| `26` | Periode akumulasi, ditolak bila di luar periode kontrak | `14` |
| `27` | Termin pembayaran premi bernomor urut, per mata uang | `14` · `15` |
| `28` | Skala ko-asuransi sebagai beberapa baris | `14` |
| `29` | Batas tanggungan untuk bahaya apa pun di daftar | `14` · `15` |
| `30` | Dokumen dilampirkan, dirujuk bukan disalin | `14` |
| `31` | Layer beserta limit, deductible, MDP, dan pemulihan limitnya | `14` · `15` |
| `32` | Syarat berbeda tiap pemulihan limit | `31` |
| `33` | Bagian NuRe atas layer, hanya pada non-proporsional | `31` |
| `34` | Ketentuan proporsional per kelompok treaty di dalam layer | `31` · `15` |
| `35` | Tepat satu dari persen quota share atau jumlah lines surplus | `34` |
| `36` | Baris surplus tanpa quota share menerima kegagalan yang menyebut apa | `34` |
| `37` | Potongan atas premi, aturan sama di kedua pelekatannya | `33` · `34` · `15` |
| `38` | Penyebaran bagian NuRe per jenis reasuransi per mata uang | `33` · `34` · `15` · **`Uji AD`** · **`L-3`** |
| `39` | Jejak perubahan: siapa mengubah fakta apa dan kapan | `14` |
| `40` | Nilai versi sebelumnya berdampingan, di-SELECT bukan disalin | `14` · **`01`** ⟦papan Adjustment⟧ |
| `41` | Konsumen hilir membaca identitas kontrak dalam bentuk lama | `14` |
| `42` | Arsip JSON sistem lama disimpan dan dapat ditunjukkan | `14` |
| `43` | Sakelar penegakan trigger selama pemindahan | `14` |
| `44` | Isi tabel acuan dipindahkan dari sistem lama | `15` · `43` · **`KTV-A`** |

##### Dapat dimulai hari pertama

**`14` dan `15` — dua, dan keduanya tidak bergantung pada tiket mana pun.**

Keduanya **PEMBUAT PERTAMA**, dan itu bukan kebetulan: `14` membuat `KONTRAK` dan `VERSI_KONTRAK`
yang **seluruh papan** menggantung padanya, `15` membuat keenam tabel acuan yang **setiap kunci
asing** menunggunya. Begitu keduanya mendarat, **yang dapat dimulai melompat ke tiga belas**.

> **Taksiran keduanya tidak dibandingkan dengan tiket mana pun**, sesuai aturan `PEMBUAT PERTAMA`.

##### Menunggui, tidak menahan

Tiket ini **boleh jalan**. Jawaban yang ditunggu mengubah satu constraint, satu lebar kolom, atau satu
jalur migrasi — bukan rancangannya.

| # | Menunggu | Yang berubah bila jawabannya lain |
|---|---|---|
| `01` | `REV-3` | **keadaan mana yang dihitung "berlaku"** — bukan bahwa ada penunjuk dasar |
| `02` | `DB-20` | **titik beku** materialitas, dan titik beku milik `07` |
| `09` | `DB-16a` | **bentuk kolom** dokumennya, yang ikut `04` |
| `14`, `31` | `T-6` | tiga kolom mata uang kosong pada `VERSI_KONTRAK`, dua pada `LAYER`. **Dicabut**, bukan diubah |
| `23` | `T-4` · `29` `T-2` · `31` `T-3` | **tingkat pencatatan** paket uangnya — satu constraint, bukan bentuk tabel |
| `31` | **`F-15`** | `PERSEN_ROL` dapat **dicabut**, dan tabel anak premi diperoleh dapat **ditambahkan**. Bentuk `LAYER` tidak berubah |
| `39` | **`F-16`** | **dari mana** nilai `PERAN_PELAKU` diambil. **Bentuknya sudah diputuskan** `KTV-D` — potret, bukan rujukan — dan bentuk itu tidak berubah |
| `43` | `L-3` | **ujinya**, bukan bentuk sakelarnya |
| enam belas **PEMBUAT PERTAMA** | `KTV-A` | **lebar kolom**. Dapat dipersempit **hanya sebelum data dimuat** — dan yang memuat data pertama adalah `44` |

##### Menahan sungguhan — lima

| # | Penahan | Kenapa ia menahan, bukan menunggui |
|---|---|---|
| `04` | `DB-16a` | bila dokumen ternyata dikirim ke luar, **jenis bernilai tiga** dan penanda kembali. Itu kolom, bukan penyetelan |
| `07` | `DB-20` | titik beku adalah **pokok** tiketnya |
| `08` | `DB-16b` | tanggalnya **adalah** pokok tiketnya, dan `KTV-2` menuntut kolomnya **dicabut sebelum data masuk** bila dibantah |
| `38` | `Uji AD` · `L-3` | `Uji AD` menentukan **golongannya**, dan bila **B** tiket ini menuntut medan `CARA MENYALAKANNYA` yang belum ada. `L-3` berarti uji negatif MV-nya **belum pernah dijalankan di mana pun** |
| `44` | `KTV-A` | **bertenggat**: ia yang **memuat data pertama**, dan sesudahnya presisi tidak dapat dipersempit lagi |

**Empat pertanyaan bisnis di antaranya — `DB-16a`, `DB-16b`, `DB-20` — belum dikirim.**

---

#### UKURAN SELESAI DI PAPAN INI — dibaca sebelum menilai satu tiket pun

Skill `to-tickets` menuntut tiap tiket memotong tegak **skema · API · UI · uji**. **Di papan ini, dua
dari empat lapisan itu belum ada.**

`L-4` menyatakan: *"tidak ada spesifikasi layar di mana pun."* Lubang itu sudah punya pemilik —
**batch lapisan aplikasi** — dan **bukan batch ini**.

> ### JANGAN MENGARANG LAYAR UNTUK MEMENUHI BENTUK IRISAN TEGAK.
>
> Layar yang dikarang akan dipakai sebagai kriteria selesai oleh orang yang **tidak tahu ia
> karangan** — dan ia akan membangun yang salah dengan keyakinan penuh.

Yang berlaku di sini gerbang **`G1`** modul ini: **pembagian tiket memakai kemampuan, bukan layar.**
Irisan tegak di papan ini berarti tegak melalui **lapisan yang ada**:

```
skema  →  constraint dan invarian  →  jalur migrasi  →  uji negatif DAN positif
```

| | |
|---|---|
| **Sebuah tiket selesai bila** | perilakunya **dapat gagal** dan **dapat diperiksa** |
| **Bukan bila** | ia dapat diperagakan di layar |

##### Dan satu kewajiban yang lahir dari ukuran ini

Karena "dapat diperagakan" tidak tersedia sebagai bukti, **uji menjadi satu-satunya bukti**. Maka tiap
tiket membawa **uji negatif dan uji positif**, bukan hanya yang pertama.

> Alasannya bukan kesetaraan: ketiga kekeliruan lingkup di modul ini **lolos uji negatif**. Yang
> menangkapnya hanya uji positif — uji yang menuntut **data sah diterima**. Constraint yang menolak
> terlalu banyak lulus setiap uji negatif yang pernah ditulis untuknya.

---

#### Aturan papan

- **Status adalah medan di dalam berkas tiket** (`status: aktif | tertahan | selesai | mati`), bukan
  lokasi foldernya. Papan **dibangkitkan dari berkas**, dan **tidak pernah menghapus tiket**.
- **Empat golongan, tidak ada yang kelima.**
- **Penahan dari luar papan ditulis dengan nama aslinya** — `REV-3`, `DB-16b`, `Uji AD`, `L-3`,
  `KTV-A`, `F-15`, `F-16`, `T-2`…`T-6` — tidak diterjemahkan menjadi nomor tiket.
- **Penahan yang selesai dicoret, tidak dihapus**, supaya rantainya tetap terbaca:
  `~~05~~ *(selesai)*`.
- **Nomor yang lompat bukan kekeliruan.** Penomoran tidak disusun ulang.
- **"Menunggui, tidak menahan" dipisahkan dari penahan sungguhan**, supaya papan tidak berteriak
  serigala.
- **Pencacah diberi label bendanya.** Papan ini mencacah **tiket**; kemampuan `P-nn` punya
  pencacahnya sendiri di `DAFTAR-PEKERJAAN.md`, dan keduanya **tidak pernah dijumlahkan**.
- **Tiket kemampuan yang DIUBAH menyebut bunyi lamanya** di medan `Menggantikan:` — bukan hanya kode
  cacatnya.
- **Satu kode per `DIASUMSIKAN-CLEAR(...)`.** Penanda berisi dua kode **tidak dapat disapu**; tiket
  `07` sempat memuat `(DB-20, REV-3)` dan **sudah dipecah** menjadi dua baris, 24 September 2026.

#### Satu angka yang menggambarkan keadaan, dan ia bukan "aktif"

Papan ini mencetak **"Yang sebenarnya dapat dimulai"** — tiket yang **tidak tertahan** dan **seluruh
penghalangnya, di dalam maupun di luar papan, sudah selesai**.

Angka itu, bukan cacah "aktif", yang menjawab *"apakah batch ini layak dimulai sekarang"*. Hari ini
ia **2 dari 44** — dan itu keadaan yang sehat, bukan buruk: keduanya fondasi yang selebihnya
menunggu.

#### Asumsi yang dianggap clear

Sebagian tiket bersandar pada penahan yang **dianggap clear atas perintah pemilik proses**, bukan yang
benar-benar tertutup. Masing-masing ditulis di medan `Dasar:` sebagai `DIASUMSIKAN-CLEAR(<kode>)`, dan
daftarnya beserta **apa yang ditinjau ulang bila ia patah** ada di
[`../ASUMSI-CLEAR.md`](../ASUMSI-CLEAR.md).

Satu sapuan atas `DIASUMSIKAN-CLEAR(` di folder ini mengeluarkan daftarnya. **Tidak ada yang perlu
diingat.**

#### Satu ketaksamaan bentuk yang tercatat, bukan ditambal

Ketiga belas tiket `01`…`13` **tidak punya medan `golongan:` di frontmatter**, sementara
`BENTUK-TIKET.md` §6.1 mewajibkannya dan tiket `14`…`44` memuatnya. Golongan sebuah tiket —
pelestarian, perubahan, atau baru — **menentukan cara mengujinya**, dan tanpa medan itu pembacanya
harus menyimpulkannya dari isi.

**Tidak ditambal dari sini:** mengisinya berarti **memutuskan golongan tiket yang sudah disetujui**,
dan itu wewenang pemilik proses. **Ditagih** pada putaran berikutnya yang menyentuh ketiga belas tiket
itu.


---

#### Batch 2 Treaty In — `45`…`64` · daur hidup dan warisan

| # | Tiket | Ditahan oleh |
|---|---|---|
| `45` | **Daftar keadaan dan perpindahan sah berdiri**; perpindahan di luar daftar ditolak | `14` |
| `46` | Versi terminal beku, menjangkau **seluruh entitas anaknya** | `45` |
| `47` | Paling banyak satu versi tak-terminal per kontrak | `45` |
| `48` | Pengajuan menolak dua syarat, memperingatkan enam | `45` |
| `49` | **SH** menyetujui, versi pindah ke antrian DH | `48` · `54` |
| `50` | **DH** menyetujui, versi pindah ke antrian DR | `49` |
| `51` | **DR** menyetujui, versi menjadi `DISETUJUI` bila seluruh `K2` | `50` |
| `52` | Pengaju tidak menyetujui; tiap tingkat orang berbeda; tanpa jalan pintas | `51` |
| `53` | Pengembalian ke `DRAFT` dengan alasan wajib | `49` |
| `54` | **Tiap perpindahan meninggalkan catatan persetujuan** | `45` |
| `55` | Penolakan versi — kontrak maupun addendum — meninggalkan catatan | `49` · `54` |
| `56` | Pembatalan draf oleh pembuatnya sendiri | `45` · `54` |
| `57` | Angka rupiah pada versi disetujui tidak bergeser | `46` · `20` |
| `58` | Wewenang persetujuan dibatasi nilai kontrak | `51` · **dokumen batas wewenang** |
| `59` | Kepala kontrak warisan dipindahkan beserta keadaannya | `45` · `44` |
| `60` | Keadaan warisan tanpa padanan → `WARISAN_TAK_TERPETAKAN` | `59` |
| `61` | Baris warisan diperbaiki ke keadaan sah yang dipilih eksplisit | `60` |
| `62` | Kontrak dicari, termasuk menurut keadaan siklus hidupnya | `45` |
| `63` | Cara pembukuan XOL, hanya pada non-proporsional | `14` · **`Uji X-2`** |
| `64` | Tabel acuan pembagian kapasitas punya pemiliknya | `15` · **penetapan pemilik** |

##### Tiga PEMBUAT PERTAMA baru di batch ini

| Tiket | Entitas |
|---|---|
| `54` | **`CATATAN_PERSETUJUAN`** — ditemukan lewat pemeriksaan entitas-per-irisan; tanpa pemeriksaan itu, `49` akan menulis catatan ke tabel yang belum ada |
| `45` | himpunan nilai `KEADAAN_SIKLUS_HIDUP` — **kolomnya** dibuat `14`, **artinya** dibuat di sini |
| `59` | isi kepala kontrak warisan — data, bukan tabel |

##### Kenapa `49`, `50`, `51` tiga tiket dan bukan satu

`U-2` ditegakkan harfiah: satu tiket tidak melintasi lebih dari **satu** perpindahan, dan `GRL-08`
menulis `SETUJUI×3` — **tiga perpindahan berbeda**, bukan satu perpindahan berparameter.

> Satu tiket yang melintasi tiga tingkat **gagal sebagai satu blok**. Bila tingkat kedua yang salah,
> tiketnya *"belum selesai"* — dan **tidak ada yang tahu tingkat mana yang patah**.

Dan wewenangnya berbeda per tingkat: `ADR-D-TI-0044` menuntut **setiap larangan punya tepat satu sebab**,
sehingga penolakan karena **peran** dan penolakan karena **keadaan** harus berpesan berbeda — tiga
kali, dengan jawaban yang berbeda tiap tingkat.

---

#### Satu entitas yang TIDAK punya tiket, dan tidak punya kemampuan

Pemeriksaan entitas-per-irisan atas **seluruh 35 entitas** `KAMUS-KOLOM.md` terhadap **64 tiket**
menemukan satu yatim:

| Entitas | Keadaan |
|---|---|
| **`PERISTIWA_KONTRAK`** | **nol tiket menyebutnya, dan nol kemampuan `P-nn` menghasilkannya** |

Ia berdiri di `SPEC-MODEL-DATA.md` §10.21a dan di `ddl-usulan/`, lahir dari pemisahan `CommentList`
sistem lama menjadi **dua** entitas: persetujuan (`IsApproved` terisi) dan **peristiwa**
(`IsApproved` kosong — *"Copied from ID …"*, *"Had Created Internal Edit"*).

> **Tabelnya akan berdiri kosong selamanya**, dan tidak ada pencacah yang memperlihatkannya. Ia
> bentuk `K-3` satu tingkat lebih dalam: di sana **kemampuan** jatuh di antara dua batch; di sini
> **entitas** tidak punya kemampuan sama sekali.

| | |
|---|---|
| **Siapa menutup** | pemilik proses — satu kalimat: kemampuan `P-67` tersendiri, atau ia ikut tiket `59` sebagai bagian pemindahan warisan |
| **Yang menagih** | berkas ini, dan `ddl-usulan/` yang memuat tabelnya |

# Treaty In Adjustment

Jumlah tiket: **14**

## Treaty In Adjustment - 01 - Versi baru lahir dari versi berlaku terakhir, dan rujukan dasarnya disimpan eksplisit

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-42` (diubah 24 Sep 2026) · `SPEC-MODEL-DATA.md` §10.2 · `GRL-10`, bahan to-spec `B-3` dan `B-4`.*

**What to build:** **PK** membuat versi penyesuaian atas sebuah kontrak, dan versi itu lahir `DRAFT` sambil
menyimpan **rujukan eksplisit** ke versi yang menjadi dasarnya. Dasarnya **bukan baris yang dipilih
di layar** dan **bukan versi `DISETUJUI` mana pun** — ia **versi berlaku terakhir** pada saat versi
baru dibuat.

Artefak: kolom `ID_VERSI_KONTRAK_DASAR` pada `VERSI_KONTRAK`, kunci asingnya, dan constraint yang
menjaga keterisiannya.

**Persyaratan:** `INV-04` · bahan to-spec `B-3` (wajib terisi pada versi penyesuaian, wajib kosong pada versi pertama) · `B-4` (yang ditunjuk harus versi berlaku terakhir, bukan `DITOLAK` maupun `DIBATALKAN`) · `ADR-D-TI-0040`

**Tidak termasuk:** **Penomoran ulang baris warisan** — itu irisan 10. Baris warisan boleh punya dasar kosong
sampai irisan itu selesai.
**Turunan "versi berlaku"** — `GRL-11` menetapkan ia dihitung, bukan disimpan; tidak ada kolom yang
dibangun untuknya di sini.

**Jalur gagal:** Menunjuk versi berkeadaan `DITOLAK` sebagai dasar -> **ditolak saat simpan**, dan pesannya
menyebut keadaan versi yang ditunjuk · Versi pertama sebuah kontrak dengan dasar terisi -> ditolak ·
Versi penyesuaian dengan dasar kosong -> ditolak.

**Uji:** **Negatif:** simpan versi penyesuaian tanpa dasar; simpan versi pertama dengan dasar; tunjuk
versi `DITOLAK`; tunjuk versi `DIBATALKAN`.
**Positif — dan ia yang menangkap lingkup yang terlalu sempit:** kontrak yang punya **dua** versi
`DISETUJUI` dalam sejarahnya menerima versi baru yang menunjuk **yang terakhir**, bukan ditolak
karena ada dua.

**Menggantikan:** `P-42` bunyi lama: *"**PK** dapat membuat versi baru **dari versi yang sudah `DISETUJUI`**,
dan versi baru itu mulai dari `DRAFT` serta melewati keempat tingkat."*

Yang bergeser: dasarnya menjadi **versi berlaku terakhir**, dan rujukannya **disimpan eksplisit**.
Bunyi lama mengizinkan memilih versi `DISETUJUI` yang bukan terakhir — dan selisih yang dihitung
terhadap dasar yang salah **tidak menghasilkan galat**. Sistem lama: `OLDID` = ID baris yang
**dipilih di picker**, tanpa pembeda jenis dan tanpa saringan keadaan (`TDA-11`).

**Blocked by:** **14** — *ditambahkan 24 September 2026.* Tiket ini menambahkan **kolom** atau **tabel anak** pada `VERSI_KONTRAK`, dan tidak satu pun tiket di papan membuat tabelnya. Irisan `14` adalah **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`.

**Dasar:**
```
EVIDENCED(TreatyInRevisi_post@ekspor-2026-09, TreatyLoadMasterJoinEdm@ekspor-2026-09)
        DECIDED(GRL-10, ADR-D-TI-0040)
        DIASUMSIKAN-CLEAR(REV-3)
```

- [ ] `ID_VERSI_KONTRAK_DASAR` berdiri di `VERSI_KONTRAK` dengan kunci asingnya
- [ ] versi penyesuaian tanpa dasar **ditolak**, versi pertama dengan dasar **ditolak**
- [ ] menunjuk versi `DITOLAK` atau `DIBATALKAN` **ditolak**, pesannya menyebut keadaannya
- [ ] uji positif lulus: kontrak berdua versi `DISETUJUI` menerima versi baru
- [ ] `REV-3` tercatat di `ASUMSI-CLEAR.md` sebagai asumsi yang tiket ini bersandar padanya

## Treaty In Adjustment - 02 - Materialitas mengunci ruas mana yang boleh disunting, dan penolakannya ditegakkan di sisi simpan

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-60` · `GRL-20` · `SPEC-MODEL-DATA.md` §10.2a.*

**What to build:** **PK** menyatakan materialitas sebuah versi — `MATERIAL` atau `TIDAK_MATERIAL` — dan
pernyataan itu **menentukan ruas mana yang boleh disunting**. Penolakannya berdiri **di sisi
simpan**, bukan hanya di layar.

Artefak: kolom `SIFAT_MATERIAL_ADDENDUM`, dan penegakan per-ruas yang menolak penyimpanan.

**PEMBUAT PERTAMA** untuk `SIFAT_MATERIAL_ADDENDUM`.

**Persyaratan:** `GRL-20` · `SPEC-MODEL-DATA.md` §10.2a · `ADR-D-TI-0037` butir 3 (besaran yang boleh disepakati dinaikkan menjadi masukan)

**Tidak termasuk:** **`INV-69` dan `INV-70`** — penegakan atas baris `NILAI_SELISIH` — itu irisan 11, dan ia
menunggu `NILAI_SELISIH` lahir di irisan 06.
**Titik beku materialitas** — itu irisan 07.
**Nilai warisan** `EDMMaterialType` — dibawa apa adanya oleh jalur migrasi, bukan oleh irisan ini.

**Jalur gagal:** Versi `TIDAK_MATERIAL` yang penyimpanannya mengubah ruas terkunci -> **ditolak**, dan
pesannya menyebut **ruas mana** yang terkunci · Versi `MATERIAL` yang mengubah ruas yang sama ->
diterima.

**Uji:** **Negatif:** simpan versi `TIDAK_MATERIAL` yang mengubah ruas uang terkunci.
**Positif — dan ia yang menangkap penguncian yang terlalu lebar:** versi `TIDAK_MATERIAL` yang
mengubah ruas yang **memang boleh** berubah **diterima**. Ketiga kekeliruan lingkup di modul ini
lolos uji negatif; yang menangkapnya hanya uji positif.

**Menggantikan:** `TDA-10` — *materialitas hanya ditegakkan di layar; nol penegakan di sisi simpan.*
Sistem lama memasang kondisi penguncian di badan seksi dan **tidak memeriksa apa pun saat
menyimpan**, sehingga nilai yang seharusnya terkunci tetap dapat masuk lewat jalur lain.

**Blocked by:** **14** — *ditambahkan 24 September 2026.* Tiket ini menambahkan **kolom** atau **tabel anak** pada `VERSI_KONTRAK`, dan tidak satu pun tiket di papan membuat tabelnya. Irisan `14` adalah **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`.

**Dasar:**
```
EVIDENCED(TreatyInEDMSetValue@ekspor-2026-09, ekspor-tambahan/EDMMaterialType.xml)
        DECIDED(GRL-20, ADR-D-TI-0037)
        DIASUMSIKAN-CLEAR(DB-20)
```

- [ ] `SIFAT_MATERIAL_ADDENDUM` berdiri dengan dua nilainya, boleh kosong pada versi pertama
- [ ] daftar ruas yang dikunci tiap nilai **tertulis**, dan diturunkan dari kondisi layar lama — bukan disusun ulang
- [ ] simpan yang melanggar **ditolak**, pesannya menyebut ruas yang terkunci
- [ ] uji positif lulus: ruas yang boleh berubah tetap dapat berubah
- [ ] `DB-20` tercatat di `ASUMSI-CLEAR.md`

## Treaty In Adjustment - 03 - Tanggal berlaku pada versi ditolak bila jatuh di luar periode kontrak

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-45` (diubah 24 Sep 2026) · `INV-54` · `GRL-15` · `KTV-2`.*

**What to build:** **PK** menyatakan sejak kapan sebuah versi berlaku, dengan bawaan **tanggal mulai
kontrak**. Tanggal yang jatuh di luar periode kontraknya **ditolak**.

Artefak: atribut tanggal berlaku pada `VERSI_KONTRAK`, bawaannya, dan trigger `INV-54`.

**Persyaratan:** `INV-54` (trigger) · `GRL-15`

**Tidak termasuk:** **Mesin pro rata — sengaja TIDAK dibangun.** Ini pernyataan keputusan, bukan `TODO`:

| | |
|---|---|
| apa yang sengaja tidak dilakukan | rumus yang membagi besaran versi menurut porsi periode tersisa |
| kenapa | sistem lama **tidak pernah menjalankannya sungguhan** — `EDMEffective` punya satu penulis yang menyamakannya dengan tanggal mulai, dan kontrol layarnya hanya-baca, sehingga `ProRatePercent` **tidak pernah dapat selain 100** |
| akibatnya pada angka | **nol** — setiap versi berlaku penuh, sama seperti sistem lama |
| di mana selisihnya terlihat | `UA-19` — sebaran waktu baris ber-`ProRatePercent` bukan 100 |
| kapan ditagih | saat `DB-5` dijawab |

**Tanggal berlaku pada DOKUMEN** adalah ruas yang berbeda — irisan 08.

**Jalur gagal:** Tanggal berlaku versi sebelum tanggal mulai kontrak atau sesudah tanggal berakhir -> **ditolak** trigger, pesannya menyebut periode kontraknya.

**Uji:** **Negatif:** tanggal berlaku sebelum mulai; sesudah berakhir; tepat di batas luar.
**Positif:** tanggal berlaku **tepat pada** tanggal mulai dan tepat pada tanggal berakhir
**diterima** — batas inklusif, dan uji negatif saja tidak pernah membuktikannya.

**Menggantikan:** `P-45` bunyi lama: *"**Tanggal berlaku addendum** yang jatuh di luar periode kontraknya
**ditolak**."*

Yang bergeser: **pembawanya disebut** — *"tanggal berlaku **pada versi**"*. `KTV-2` memberi dokumen
addendum tanggal berlakunya sendiri, sehingga tanpa menyebut pembawanya irisan ini dan irisan 08
**sama-sama mengklaim ruas bernama sama**, dan yang menang adalah yang ditulis lebih dulu.

**Blocked by:** **14** — *ditambahkan 24 September 2026.* Tiket ini menambahkan **kolom** atau **tabel anak** pada `VERSI_KONTRAK`, dan tidak satu pun tiket di papan membuat tabelnya. Irisan `14` adalah **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`.

**Dasar:**
```
EVIDENCED(TreatyInSetEditPre@ekspor-2026-09, TreatyCalculateProratePct@ekspor-2026-09)
        DECIDED(GRL-15, KTV-2, ADR-D-TI-0037)
```

- [ ] atribut tanggal berlaku berdiri pada `VERSI_KONTRAK` dengan bawaan tanggal mulai kontrak
- [ ] trigger `INV-54` menolak tanggal di luar periode, pesannya menyebut periodenya
- [ ] uji positif batas inklusif lulus
- [ ] pernyataan keputusan **mesin pro rata tidak dibangun** tertulis di tiket ini dan di §10.2, bukan sebagai `TODO`

## Treaty In Adjustment - 04 - Dokumen addendum berdiri bernomor sendiri, dan satu dokumen dapat menyentuh beberapa kontrak

---
status: tertahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-61` · `GRL-19` · `INV-71` · `STRUKTUR-ADDENDUM.md` §3.*

**What to build:** **PK** mencatat sebuah **dokumen addendum** dengan **nomornya sendiri**, dan dokumen itu
dapat menunjuk versi di **beberapa kontrak sekaligus**. Persetujuan **tetap per kontrak** — dokumen
tidak membawa keadaan persetujuan apa pun.

Artefak: entitas `DOKUMEN_ADDENDUM`, hubungannya ke `VERSI_KONTRAK`, dan keunikan nomornya.

**PEMBUAT PERTAMA** untuk `DOKUMEN_ADDENDUM`.

**Persyaratan:** `GRL-19` · `INV-71` (keunikan nomor dokumen) · `KTV-3` (tidak ada penanda "dikirim ke luar"; jenis tetap bernilai dua)

**Tidak termasuk:** **Tanggal berlaku dokumen** — irisan 08.
**Pengisian nomor untuk baris warisan** — irisan 09; migrasi **tidak punya sumber** untuk mengisinya.
**Keadaan persetujuan pada dokumen** — tidak ada, dan itu keputusan: persetujuan melekat pada versi
per kontrak (*"1 kontrak di aksep 1 per 1"*).

**Jalur gagal:** Dua dokumen bernomor sama -> **ditolak** `INV-71`, pesannya menyebut nomor yang bentrok ·
Dokumen tanpa satu pun versi yang ditunjuknya -> ditolak.

**Uji:** **Negatif:** simpan dua dokumen bernomor sama; simpan dokumen tanpa versi.
**Positif — dan ia pokok irisan ini:** satu dokumen menunjuk versi di **tiga kontrak berbeda**
**diterima**, dan ketiga versi tetap berjalan di rantai persetujuannya masing-masing.

**Menggantikan:** `TDA-11` — *picker menyatukan kontrak dan addendum tanpa pembeda dan tanpa saringan
keadaan.* Sistem lama **tidak punya benda bernama dokumen sama sekali**: sapuan dua inventaris
properti atas dua belas akar kata mengembalikan **nol** properti penyimpan nomor dokumen, dan kelas
integrasi addendum memuat tepat 19 properti tanpa satu pun di antaranya.

**Blocked by:** **14** — *ditambahkan 24 September 2026.* Tiket ini menambahkan **kolom** atau **tabel anak** pada `VERSI_KONTRAK`, dan tidak satu pun tiket di papan membuat tabelnya. Irisan `14` adalah **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`.

**Dasar:**
```
EVIDENCED(sapuan-properti-dokumen@ekspor-2026-09 - NOL hasil, dikalibrasi atas EDMState/EDMMaterialType/OLDID)
        DECIDED(GRL-19, KTV-3, INV-71)
        DIASUMSIKAN-CLEAR(DB-16a)
```

- [ ] `DOKUMEN_ADDENDUM` berdiri dengan nomornya sebagai kunci alami
- [ ] `INV-71` menolak nomor ganda, pesannya menyebut nomornya
- [ ] satu dokumen menunjuk versi di beberapa kontrak; persetujuan tiap versi **tidak saling menyentuh**
- [ ] dokumen **tidak** membawa kolom keadaan persetujuan — diperiksa, bukan diandaikan
- [ ] `DB-16a` tercatat di `ASUMSI-CLEAR.md`; bila dibantah, jenis bernilai **tiga** dan penanda "dikirim ke luar" kembali

## Treaty In Adjustment - 05 - Perluas — kolom nomor urut versi berdiri berdampingan dengan pengenal warisan

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-65` bagian 1 · `GRL-17` · `ADR-D-TI-0042`, `ADR-D-TI-0043`.*

**What to build:** Kolom `NOMOR_URUT_VERSI` berdiri pada `VERSI_KONTRAK` dan **kosong**; pengenal warisan
`‹kontrak›/Rnn` tetap apa adanya. **Tidak ada pembaca yang berubah.**

Ini bagian **perluas** dari perubahan lebar: bentuk baru berdiri di samping yang lama sehingga
tidak ada yang patah.

**Persyaratan:** `GRL-17` · `ADR-D-TI-0042` (sejarah pindah apa adanya) · `GRL-09` (pengenal warisan dilestarikan)

**Tidak termasuk:** **Pengisian nilainya** — irisan 10.
**Pencabutan pembacaan dari pengenal** — irisan 12.
**Pengenal `/Rnn`** tidak disentuh sama sekali; ia dilestarikan apa adanya.

**Jalur gagal:** Kolom berdiri `NOT NULL` -> setiap baris warisan gagal dimuat. Ia **wajib boleh kosong** pada tahap ini, dan itu pokok bentuk perluas.

**Uji:** **Negatif:** memuat baris warisan tanpa nomor urut **tidak boleh gagal**.
**Positif:** seluruh pembaca yang ada — pencarian, pembaca hilir bentuk lama — **mengembalikan hasil
yang sama persis** sebelum dan sesudah kolomnya berdiri.

**Menggantikan:** tidak ada — bentuk lama tidak dicabut di irisan ini; ia berdiri berdampingan.

**Blocked by:** **14** — *ditambahkan 24 September 2026.* Tiket ini menambahkan **kolom** atau **tabel anak** pada `VERSI_KONTRAK`, dan tidak satu pun tiket di papan membuat tabelnya. Irisan `14` adalah **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`.

**Dasar:**
```
EVIDENCED(TreatyInRevisi_post@ekspor-2026-09 - offset @substring(ID,10,12) meleset satu)
        DECIDED(GRL-17, ADR-D-TI-0042)
```

- [ ] `NOMOR_URUT_VERSI` berdiri, **boleh kosong**
- [ ] pengenal `/Rnn` tidak berubah pada satu baris pun
- [ ] seluruh pembaca yang ada mengembalikan hasil identik sebelum dan sesudah

## Treaty In Adjustment - 06 - Baris selisih berkunci bisnis, dan mata uang ada di dalam kuncinya

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-63` · `ADR-D-TI-0048` butir 3 · `GRL-14` · `GRL-16` · bahan to-spec `E-1a`.*

**What to build:** Menyimpan versi penyesuaian menghasilkan **baris selisih per besaran yang berubah**,
dipadankan terhadap versi dasarnya lewat **kunci bisnis** — bukan menurut posisi baris. Mata uang
**ada di dalam kunci**, sehingga mengubah mata uang tampil sebagai baris dihapus ditambah baris
baru, bukan sebagai selisih angka.

Artefak: entitas `NILAI_SELISIH`, kunci padanannya, dan jalur yang mengisinya saat simpan.

**PEMBUAT PERTAMA** untuk `NILAI_SELISIH`.

**Persyaratan:** `ADR-D-TI-0048` butir 3 · `ADR-D-TI-0053` (mata uang adalah daftar) · bahan to-spec `E-1a` · `GRL-14` · `GRL-16`

**Tidak termasuk:** **`INV-69` dan `INV-70`** — irisan 11.
**Pohon `ActualValue`** — `GRL-14` menghapusnya; tidak ada pohon ketiga yang dibangun di sini.
**Perhitungan ulang selisih versi yang sudah disetujui** — `ADR-D-TI-0036`: angka dasar persetujuan beku.

**Jalur gagal:** Baris disisipkan di **tengah** daftar pada versi baru -> baris sesudahnya **tetap dipadankan
dengan pasangan bisnisnya**, bukan bergeser · Daftar baru lebih panjang daripada daftar lama ->
barisnya muncul sebagai **baris baru**, bukan dikurangi baris yang tidak ada · Mata uang berubah ->
**dua baris**, bukan satu selisih.

**Uji:** **Negatif:** sisipkan baris di tengah dan pastikan padanan **tidak** bergeser; ubah mata uang
dan pastikan hasilnya bukan pengurangan lintas mata uang.
**Positif — dan ia WAJIB, sebab `GRL-16` bersandar padanya:** mengubah **share fakultatif**
menghasilkan sedikitnya **satu** baris selisih. Uji yang tidak pernah gagal atas besaran itu berarti
besaran itu **tidak ada di dalam himpunan yang dipadankan** — dan itu persis cacat yang `GRL-16`
temukan di sistem lama.

**Menggantikan:** `TDA-04` — *selisih dipadankan menurut posisi baris.* Dan `TDA-05` — *selisih dikurangkan
tanpa memeriksa mata uang.* Sistem lama memadankan lewat `pxListSubscript` dan `<CURRENT>`; bila
sebuah baris disisipkan di tengah, seluruh baris sesudahnya dipadankan dengan baris lama yang salah
**tanpa satu galat pun**. Ditambah `TDA-06` — sebagian selisih ditulis ke pohon yang salah lalu
ditimpa.

**Blocked by:** 01

**Dasar:**
```
EVIDENCED(TreatyEDMCalculateDifference@ekspor-2026-09, TreatyEDMDifferenceShare@ekspor-2026-09, TreatyEDMDifferenceDeduction@ekspor-2026-09)
        DECIDED(ADR-D-TI-0048, ADR-D-TI-0053, GRL-14, GRL-16)
```

- [ ] `NILAI_SELISIH` berdiri dengan kunci padanannya, dan **mata uang ada di dalam kunci**
- [ ] penyisipan baris di tengah daftar **tidak** menggeser padanan
- [ ] perubahan mata uang menghasilkan dua baris, bukan satu selisih
- [ ] **uji positif share fakultatif lulus** — perubahannya menghasilkan baris selisih
- [ ] selisih versi yang sudah disetujui **tidak berubah** setelah rumusnya diperbaiki

## Treaty In Adjustment - 07 - Jenis dan materialitas beku sejak diajukan, dan tiap perubahannya meninggalkan jejak

---
status: tertahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-04` (diubah 24 Sep 2026) · `GRL-18` · `KTV-1` · `ADR-D-TI-0045`.*

**What to build:** Selama versinya `DRAFT`, **PK** dapat membetulkan jenis dan materialitas, dan **tiap
perubahan meninggalkan baris jejak** berisi nilai lama, nilai baru, siapa, dan kapan. Begitu versinya
**diajukan**, keduanya **beku** — mengubahnya menuntut versinya dikembalikan ke `DRAFT` lebih dulu.

Artefak: penegakan titik beku pada perpindahan `AJUKAN`, dan baris `JEJAK_PERUBAHAN` untuk kedua
sumbu.

**Persyaratan:** `GRL-18` · `KTV-1` · `ADR-D-TI-0045` (jejak perubahan sebagai fakta mesin, tidak dapat dilewati) · `ADR-D-TI-0055` (perpindahan `KEMBALIKAN`)

**Tidak termasuk:** **Pembekuan seluruh kepala versi** — `P-38` sudah memegangnya pada `DISETUJUI`; irisan ini
hanya memajukan titik beku **untuk dua sumbu**.
**Penguncian ruas oleh materialitas** — irisan 02.

**Jalur gagal:** Mengubah jenis pada versi yang sudah diajukan -> **ditolak**, dan pesannya menyebut bahwa
versinya harus dikembalikan ke `DRAFT` lebih dulu · Mengubah materialitas pada versi yang menunggu
persetujuan -> ditolak · Perubahan pada `DRAFT` tanpa baris jejak -> **tidak mungkin**; jejaknya
fakta mesin.

**Uji:** **Negatif:** ubah jenis sesudah `AJUKAN`; ubah materialitas sesudah `AJUKAN`; ubah dua kali
lalu kembali ke nilai semula dan pastikan jejaknya **tiga baris**, bukan nol.
**Positif:** mengubah jenis pada versi `DRAFT` **berhasil**, dan sesudah `KEMBALIKAN` ia **berhasil
lagi** — pembekuan tidak boleh menjadi penguncian permanen.

**Menggantikan:** `P-04` bunyi lama: *"**PK** dapat mengisi dan mengubah seluruh kepala versi **selama
versinya `DRAFT`**."*

Yang bergeser: pembekuan sesudah `DRAFT` **dinyatakan**, bukan disiratkan; dan perubahan jenis
maupun materialitas **meninggalkan jejak**. Bunyi lama diam tentang apa yang terjadi sesudah `DRAFT`
dan diam sama sekali tentang jejak. Di sistem lama, radio picker **menimpa jenis tanpa syarat apa
pun**, kapan saja, dengan **nol pemeriksaan di sisi simpan** dan **tanpa mencatat nilai
sebelumnya** — sehingga pertanyaan *"jenis apa versi ini ketika kepala seksi menyetujuinya"* **tidak
dapat dijawab untuk satu pun baris warisan**.

**Blocked by:** 02

**Dasar:**
```
EVIDENCED(PickerTreatyInMasterRevisi@ekspor-2026-09, TreatyInEDMSetValue@ekspor-2026-09)
        DECIDED(GRL-18, KTV-1, ADR-D-TI-0045, ADR-D-TI-0055)
        DIASUMSIKAN-CLEAR(DB-20)
        DIASUMSIKAN-CLEAR(REV-3)
```

- [ ] perubahan jenis dan materialitas pada `DRAFT` menghasilkan baris `JEJAK_PERUBAHAN` bernilai lama dan baru
- [ ] perubahan sesudah `AJUKAN` **ditolak**, pesannya menyebut jalan keluarnya
- [ ] uji positif lulus: sesudah `KEMBALIKAN` keduanya dapat diubah lagi
- [ ] `DB-20` tercatat di `ASUMSI-CLEAR.md`; bila ia menjawab **beku sejak lahir**, kriteria selesai tiket ini **berbalik**

## Treaty In Adjustment - 08 - Dokumen addendum membawa tanggal berlakunya sendiri, dan boleh kosong

---
status: tertahan
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-62` · `KTV-2` · `STRUKTUR-ADDENDUM.md` §3.*

**What to build:** Sebuah dokumen addendum dapat menyimpan **tanggal berlakunya sendiri**, terpisah dari
tanggal berlaku versi. Kosong **diterima** — tidak setiap dokumen membawa tanggal.

**Persyaratan:** `KTV-2` · `ADR-D-TI-0035` (kegagalan tidak disamarkan menjadi nilai — kosong berarti kosong, bukan tanggal mulai kontrak)

**Tidak termasuk:** **Tanggal berlaku pada versi** — irisan 03. Keduanya ruas berbeda pada pembawa berbeda, dan
irisan ini **tidak boleh** menuliskan nilainya ke versi.
**Aturan bisnis yang memakai tanggal dokumen** — belum ada; `DB-16b` yang menentukannya.

**Jalur gagal:** Tanggal dokumen diisi lalu dibaca sebagai tanggal berlaku versi -> **cacat**; keduanya tidak boleh saling menimpa · Kosong ditolak -> salah; kosong adalah nilai sah.

**Uji:** **Negatif:** pastikan mengisi tanggal dokumen **tidak** mengubah tanggal berlaku versi mana pun.
**Positif:** dokumen tanpa tanggal **diterima** dan tetap dapat menunjuk versi.

**Menggantikan:** tidak ada — sistem lama tidak punya entitas dokumen, sehingga tidak punya tanggal untuknya.

**Blocked by:** 04

**Dasar:**
```
DECIDED(KTV-2, ADR-D-TI-0035)
        DIASUMSIKAN-CLEAR(DB-16b)
```

- [ ] kolom tanggal berlaku berdiri pada `DOKUMEN_ADDENDUM`, **boleh kosong**
- [ ] uji negatif lulus: tanggal dokumen tidak menyentuh tanggal versi
- [ ] `DB-16b` tercatat di `ASUMSI-CLEAR.md` **beserta tenggatnya**: bila dibantah, kolomnya **dicabut sebelum data dimuat** — sesudah migrasi berjalan, mencabutnya berongkos

## Treaty In Adjustment - 09 - Nomor dokumen baris warisan diisi dari arsip kertas, dan yang belum terisi terlihat sebagai daftar

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-66` · `GRL-19` §3.1 · `STRUKTUR-ADDENDUM.md` §3.1.*

**What to build:** **PK** dapat mengisi nomor dokumen pada baris addendum warisan, satu per satu, dari arsip
kertas. Baris yang **belum** terisi muncul sebagai **daftar yang dapat dibaca** — bukan sebagai diam.

Ini **pekerjaan orang**, bukan kemampuan sistem yang berjalan sendiri. Migrasi **tidak mengisinya**:
ia tidak punya sumber.

**Persyaratan:** `GRL-19` §3.1 · `ADR-D-TI-0042` (baris warisan ditandai, diterima apa adanya)

**Tidak termasuk:** **Pengisian otomatis oleh migrasi** — mustahil, dan itu temuan bukan kekurangan: sapuan
ekspor mengembalikan **nol** properti penyimpan nomor dokumen.
**Menolak baris warisan tanpa nomor** — dilarang; `ADR-D-TI-0042` menerima warisan apa adanya.

**Jalur gagal:** Baris warisan tanpa nomor dokumen ditolak oleh invarian mana pun -> **cacat**; ia harus diterima dan terdaftar · Daftar yang belum terisi tidak dapat dibaca siapa pun -> pekerjaan ini tidak akan pernah dijadwalkan.

**Uji:** **Negatif:** invarian keunikan nomor dokumen **tidak boleh** menolak beberapa baris warisan
yang sama-sama kosong.
**Positif:** daftar baris tanpa nomor dapat dibaca, dan berkurang satu setiap kali sebuah nomor
diisi.

**Menggantikan:** tidak ada — sistem lama tidak pernah merekam nomor dokumen, sehingga tidak ada perilaku yang digantikan. **Ini kemampuan BARU yang menutup lubang yang dibuat oleh ketiadaan sumber.**

**Blocked by:** 04

**Dasar:**
```
EVIDENCED(sapuan-properti-dokumen@ekspor-2026-09 - NOL hasil, dikalibrasi)
        DECIDED(GRL-19, ADR-D-TI-0042)
        DIASUMSIKAN-CLEAR(DB-16a)
```

- [ ] baris warisan dapat menyimpan nomor dokumen yang diketik orang
- [ ] beberapa baris warisan **tanpa** nomor tidak saling bertabrakan pada invarian keunikan
- [ ] daftar baris yang belum terisi **dapat dibaca**, dan berkurang saat diisi
- [ ] ongkos pengisiannya **diukur** dan dilaporkan — berapa baris, dan apakah sepadan

## Treaty In Adjustment - 10 - Pindahkan — nomor urut warisan diisi menurut kronologi, per batch kontrak

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-65` bagian 2 · `GRL-17` · `ADR-D-TI-0043`.*

**What to build:** Batch kontrak warisan memperoleh `NOMOR_URUT_VERSI` yang **diberikan ulang menurut
kronologi** — sumber utama tanggal komentar pembuatan, cadangan `EDMDATE`. Pengenal `/Rnn`
**dilestarikan apa adanya**.

Batch demi batch, dan **hijau di antara batch**: bentuk lama masih berdiri.

**Persyaratan:** `GRL-17` · `ADR-D-TI-0042` · `ADR-D-TI-0043` · `GRL-11` (turunan versi berlaku bersandar pada nomor urut)

**Tidak termasuk:** **Pencabutan pembacaan dari pengenal** — irisan 12.
**Pembersihan nomor `/Rnn` yang bentrok** — dilarang; `ADR-D-TI-0042` melarang membersihkan sejarah.

**Jalur gagal:** Dua versi satu kontrak memperoleh nomor urut **sama** -> kontraknya masuk **pengecualian
migrasi bernomor**, dilaporkan per kontrak, dan **tidak** memperoleh turunan "versi berlaku" sampai
seseorang memutuskannya · Kronologi tidak dapat ditentukan dari sumber mana pun -> idem, bukan
ditebak.

**Uji:** **Negatif:** batch yang menghasilkan nomor urut kembar **harus** memunculkan pengecualian,
bukan lolos diam-diam.
**Positif:** kontrak yang kronologinya jelas memperoleh urutan yang **sama dengan urutan
`/Rnn`-nya** — bila keduanya berbeda tanpa sebab, yang salah adalah pengisian ini, bukan datanya.

**Menggantikan:** `TDA-01` dan `TDA-12` — *penjaga duplikat mati, tabrakan berakhir sebagai `UPDATE` yang
menimpa dan melapor berhasil*, dan *offset pengurai nomor revisi meleset satu*. Sistem lama
menghitung nomor dari `@substring(ID,10,12)` — memilih `/R01` ketika `R02` sudah ada menghasilkan
`R02` lagi, dan barisnya **tertimpa tanpa galat**. **Nomor yang dapat dipakai ulang bukan
urutan.**

**Blocked by:** 05

**Dasar:**
```
EVIDENCED(TreatyInRevisi_post@ekspor-2026-09, TreatyInEdmCheckDuplicate@ekspor-2026-09 - MATI)
        DECIDED(GRL-17, ADR-D-TI-0042, ADR-D-TI-0043)
```

- [ ] batch mengisi nomor urut dari tanggal komentar pembuatan, cadangan `EDMDATE`
- [ ] baris yang tidak dapat diurutkan menjadi **pengecualian migrasi bernomor**, dilaporkan **per kontrak**
- [ ] pengenal `/Rnn` tidak berubah pada satu baris pun
- [ ] **`NOMOR_URUT_VERSI` dikecualikan dari kriteria identik `ADR-D-TI-0043`**, dan sebabnya tertulis di daftar harapan uji paritas: ia mekanisme yang **diberikan**, bukan angka yang **dipindahkan**
- [ ] `UA-21` dijalankan atau tercatat belum dapat dijalankan: berapa kontrak bernomor `/Rnn` ganda, berapa addendum tertimpa

## Treaty In Adjustment - 11 - INV-69 dan INV-70 ditegakkan atas baris selisih, beserta pemantau kebasiannya

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-60` · `INV-69`, `INV-70` · `UJI-NEGATIF-INVARIAN.md` §2 dan §3.*

**What to build:** Versi yang **dinyatakan `TIDAK_MATERIAL`** tetapi memiliki baris `NILAI_SELISIH` bertipe
**uang** atau **porsi** **ditolak saat simpan**. Penegakannya lewat *materialized view*
ber-`REFRESH ON COMMIT`.

Dan karena mekanisme itu **dapat berhenti tanpa memberi tahu siapa pun**, irisan ini berdiri hanya
bila **pemantau kebasiannya** juga berdiri.

**Persyaratan:** `INV-69` · `INV-70` · `UJI-NEGATIF-INVARIAN.md` §2 (uji negatif dijalankan, bukan diargumentasikan) dan §3 (pemantau kebasian)

**Tidak termasuk:** **Penguncian ruas** — irisan 02; ia penegakan di lapisan yang berbeda.
**Pemantau kebasian untuk `INV-47`, `INV-50`, `INV-51`** — itu **lubang induk**, `F-13`. Irisan ini
menutup dua invarian dan **membiarkan tiga**; itu dinyatakan, bukan disamarkan.

**Jalur gagal:** Versi `TIDAK_MATERIAL` dengan baris selisih bertipe porsi -> **ditolak**, pesannya menyebut
besaran mana yang melanggarnya · **MV gagal me-refresh** -> penegakan berhenti **tanpa galat**, dan
**hanya pemantau kebasian yang memperlihatkannya.**

**Uji:** **Negatif — dijalankan, bukan diargumentasikan:** simpan versi `TIDAK_MATERIAL` dengan baris
selisih bertipe uang; lalu bertipe porsi.
**Positif:** versi `TIDAK_MATERIAL` yang baris selisihnya **hanya** bertipe lain **diterima**.
**Uji pemantau:** matikan refresh MV-nya dengan sengaja, dan pastikan pemantau **berbunyi** — sebuah
pemantau yang belum pernah dilihat berbunyi belum terbukti ada.

**Menggantikan:** `TDA-10` bagian kedua — *nol penegakan di sisi simpan.* Dan di sistem lama versi yang
`TIDAK_MATERIAL` tetap dapat mengubah angka uang lewat jalur yang tidak melewati layar.

**Blocked by:** 02, 06

**Dasar:**
```
EVIDENCED(sensus-kondisi-penguncian@ekspor-2026-09 - lihat GRILL-D/01-TEMUAN TD-02)
        DECIDED(GRL-20, INV-69, INV-70)
        DIASUMSIKAN-CLEAR(DB-20)
```

- [ ] constraint `INV-69` dan `INV-70` terpasang di atas *materialized view*-nya
- [ ] **uji negatif dijalankan** dan gagal sebagaimana seharusnya — hasilnya dilampirkan, bukan dinyatakan
- [ ] **pemantau kebasian berdiri**: `REFRESH_MODE`, `STALENESS`, terjadwal, dan **berbunyi kepada seseorang yang bernama**
- [ ] pemantau diuji dengan mematikan refresh secara sengaja, dan ia **berbunyi**
- [ ] `F-13` dirujuk: tiga MV induk masih **tanpa** pemantau, dan itu bukan bagian tiket ini

## Treaty In Adjustment - 12 - Kerutkan — tidak ada lagi pembaca yang menurunkan urutan dari pengenal warisan

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-65` bagian 3 · `GRL-17` · `GRL-11`.*

**What to build:** Turunan *"versi berlaku"* dan setiap pengurutan versi dibaca **hanya** dari
`NOMOR_URUT_VERSI`. Tidak ada lagi kode di mana pun yang mengurai angka dari pengenal `/Rnn`.

Pengenal `/Rnn` **tetap disimpan dan tetap ditampilkan** — yang dicabut adalah **pembacaannya
sebagai urutan**, bukan keberadaannya.

**Persyaratan:** `GRL-17` · `GRL-11` · `GRL-09` (pengenal warisan dilestarikan apa adanya)

**Tidak termasuk:** **Penghapusan pengenal `/Rnn`** — dilarang. `GRL-09` melestarikannya; nomor itu sudah
beredar di luar sistem.

**Jalur gagal:** Sebuah pembaca yang masih mengurai `/Rnn` tertinggal -> ia akan **benar untuk hampir semua kontrak** dan salah untuk yang nomornya pernah dipakai ulang — kegagalan yang paling sulit ditemukan.

**Uji:** **Negatif:** sapuan seluruh kode untuk pengurai `/Rnn` mengembalikan **nol**, dan sapuannya
**dikalibrasi** terhadap satu pemakaian yang diketahui ada sebelum irisan ini.
**Positif:** turunan "versi berlaku" atas kontrak yang nomor `/Rnn`-nya **pernah bentrok**
mengembalikan versi yang benar.

**Menggantikan:** `TDA-12` — *offset pengurai nomor revisi meleset satu; patah pada revisi kesepuluh.*
Deret yang dihasilkannya `R09 -> R010 -> R11 -> R02`, bergantung pada perilaku `@substring` di luar
panjang teks — yang **tidak terbaca dari ekspor**.

**Blocked by:** 10

**Dasar:**
```
EVIDENCED(TreatyInRevisi_post@ekspor-2026-09 - @substring(ID,10,12))
        DECIDED(GRL-17, GRL-11, GRL-09)
```

- [ ] sapuan pengurai `/Rnn` mengembalikan **nol**, dan sapuannya **dikalibrasi**
- [ ] pengenal `/Rnn` tetap tersimpan dan tetap ditampilkan
- [ ] uji positif lulus atas kontrak yang nomornya pernah bentrok

## Treaty In Adjustment - 13 - Perubahan persentase pemulihan limit menghasilkan baris selisih, dan versi tidak-material yang mengubahnya ditolak

---
status: aktif
---



*Asal: `DAFTAR-PEKERJAAN.md` `P-64` · `TDA-18` · `INV-69` · `STRUKTUR-ADDENDUM.md` §2.2.*

**What to build:** Mengubah persentase pemulihan limit pada sebuah versi penyesuaian **menghasilkan baris
`NILAI_SELISIH`** — sesuatu yang sistem lama tidak pernah lakukan.

Akibat langsungnya, dan ia yang harus diberitahukan sebelum peralihan:

> Versi yang mengubah persentase pemulihan limit **dan dinyatakan `TIDAK_MATERIAL`** **ditolak saat
> simpan** oleh `INV-69` — sebab perubahan itu kini menghasilkan baris selisih bertipe **porsi**.
> Di sistem lama hal yang sama **lolos tanpa sepatah galat**.

**Persyaratan:** `TDA-18` · `INV-69` · bahan to-spec `E-1a`

**Tidak termasuk:** **Perhitungan ulang atas versi warisan** — `ADR-D-TI-0043`: menghitung ulang adalah peristiwa
bisnis tersendiri, bukan bagian migrasi. Baris warisan **tidak** diklasifikasi ulang.

**Jalur gagal:** Versi `TIDAK_MATERIAL` yang mengubah persentase pemulihan -> **ditolak saat simpan**, dan
pesannya menyebut besaran yang melanggarnya · Perubahan yang sama pada versi `MATERIAL` -> diterima,
dan baris selisihnya muncul.

**Uji:** **Negatif:** simpan versi `TIDAK_MATERIAL` yang mengubah persentase pemulihan limit.
**Positif:** versi `MATERIAL` yang mengubahnya **diterima**, dan baris selisihnya **ada** — sebab
ketiadaan baris di sini adalah persis cacat lama, dan nol terbaca seperti "cocok".

**Menggantikan:** `TDA-18` — ***`ReinstatementPct` disalin, tidak dikurangi.*** Di
`TreatyEDMDifferenceLimits` bentuknya `ValueDifference…ReinstatementPct = .ReinstatementPct`, bukan
pengurangan terhadap `OLDDATA`. Sapuan seluruh korpus 708 berkas, dua ekspor, dua bentuk penulis:
**10 penugasan, NOL pengurangan**. Kalibrasi **lulus** — `Limit` ditemukan 6 pengurangan nyata.
Akibatnya di sistem lama: **perubahan persentase pemulihan tidak pernah menghasilkan baris
selisih**.

**Blocked by:** 11

**Dasar:**
```
EVIDENCED(TreatyEDMDifferenceLimits@ekspor-2026-09 - 10 penugasan, nol pengurangan, kalibrasi lulus atas Limit)
        DECIDED(TDA-18, INV-69)
        DIASUMSIKAN-CLEAR(DB-20)
```

- [ ] perubahan persentase pemulihan limit menghasilkan baris `NILAI_SELISIH` bertipe porsi
- [ ] versi `TIDAK_MATERIAL` yang mengubahnya **ditolak saat simpan**, pesannya menyebut besarannya
- [ ] uji positif lulus: versi `MATERIAL` diterima **dan barisnya ada**
- [ ] baris warisan **tidak** diklasifikasi ulang (`ADR-D-TI-0043`)
- [ ] **pemberitahuan pra-peralihan tertulis dan terkirim**: pekerjaan yang hari ini lolos tanpa galat akan ditolak sesudah peralihan

## Treaty In Adjustment - REA - Papan tiket — Treaty In Adjustment

**Tiket: aktif 10 · tertahan 3 · selesai 0 · mati 0** — **13 tiket**

**Yang sebenarnya dapat dimulai: 0.** Sebabnya di §3, dan ia **seluruhnya ada di papan sebelah**.

> ## PAPAN INI DIPISAHKAN DARI PAPAN INDUK — 25 September 2026
>
> Tiket `01`…`13` sebelumnya berdiri di `treaty-in/5-tiket/issues/` bersama tiket `14`…`64`.
> **Dipisahkan atas perintah pemilik proses.**
>
> ### Ini membalik `GRL-01`, dan itu dicatat, bukan disamarkan
>
> `GRL-01` berbunyi: ***"satu model, satu spesifikasi, satu papan."*** Perintah ronde-ronde
> sebelumnya menegaskannya: *"Jangan membuat papan kedua."*
>
> **Yang dibalik hanya PAPANNYA.** Yang **tidak** berubah:
>
> | Tetap | Di mana |
> |---|---|
> | satu model data | `treaty-in/SPEC-MODEL-DATA.md` — modul ini **tidak punya §10 sendiri** |
> | satu daftar invarian | `treaty-in/SPEC-INVARIAN.md` |
> | satu daftar kemampuan | `treaty-in/5-tiket/DAFTAR-PEKERJAAN.md` — `P-60`…`P-66` **tetap di sana**, tidak ikut pindah |
> | satu daftar asumsi | `treaty-in/5-tiket/ASUMSI-CLEAR.md` |
> | penomoran tiket | `01`…`13` di sini, `14`…`64` di sana — **tidak ada tabrakan nomor**, dan tidak satu pun disusun ulang |
>
> **Ongkos pemisahan ini nyata dan tertulis di §3:** papan ini **tidak dapat menghitung sendiri**
> tiket mana yang dapat dimulai, sebab lima dari tiga belas tiketnya diblokir tiket yang **tidak ada
> di sini**.

---

#### 1. Ukuran selesai di papan ini

Sama dengan papan induk, dan **bukan** ukuran batch lapisan aplikasi.

`L-4` menyatakan *"tidak ada spesifikasi layar di mana pun"*, dan lubang itu milik **batch lapisan
aplikasi**. Irisan di sini tegak melalui **lapisan yang ada**:

```
skema  →  constraint dan invarian  →  jalur migrasi  →  uji negatif DAN positif
```

| | |
|---|---|
| **Sebuah tiket selesai bila** | perilakunya **dapat gagal** dan **dapat diperiksa** |
| **Bukan bila** | ia dapat diperagakan di layar |

**Jangan mengarang layar untuk memenuhi bentuk irisan tegak.** Layar yang dikarang akan dipakai
sebagai kriteria selesai oleh orang yang **tidak tahu ia karangan**.

Dan karena "dapat diperagakan" tidak tersedia sebagai bukti, **uji menjadi satu-satunya bukti**:
tiap tiket membawa **uji negatif dan uji positif**. Ketiga kekeliruan lingkup di modul ini **lolos
uji negatif**; yang menangkapnya hanya uji positif — uji yang menuntut **data sah diterima**.

---

#### 2. Papan

| # | Tiket | Ditahan oleh |
|---|---|---|
| `01` | Versi baru lahir dari versi berlaku terakhir, rujukan dasarnya disimpan eksplisit | **`14`** ⟦luar⟧ |
| `02` | Materialitas mengunci ruas yang boleh disunting, ditegakkan di sisi simpan | **`14`** ⟦luar⟧ |
| `03` | Tanggal berlaku pada versi ditolak bila di luar periode kontrak | **`14`** ⟦luar⟧ |
| `04` | Dokumen addendum bernomor sendiri, menyentuh beberapa kontrak | **`14`** ⟦luar⟧ · **`DB-16a`** |
| `05` | *Perluas* — kolom nomor urut versi berdiri berdampingan | **`14`** ⟦luar⟧ |
| `06` | Baris selisih berkunci bisnis, mata uang di dalam kuncinya | `01` |
| `07` | Jenis dan materialitas beku sejak diajukan, dan berjejak | `02` · **`DB-20`** |
| `08` | Dokumen addendum membawa tanggal berlakunya sendiri | `04` · **`DB-16b`** |
| `09` | Nomor dokumen baris warisan diisi dari arsip kertas | `04` |
| `10` | *Pindahkan* — nomor urut warisan diisi menurut kronologi | `05` |
| `11` | `INV-69` dan `INV-70` ditegakkan, beserta pemantau kebasiannya | `02` · `06` |
| `12` | *Kerutkan* — urutan tidak lagi diturunkan dari pengenal | `10` |
| `13` | Perubahan reinstatement menghasilkan baris selisih | `11` |

**⟦luar⟧** berarti penahannya **tidak ada di papan ini**. Ia ditulis dengan **nomor aslinya**, tidak
diterjemahkan — yang menahan dari luar harus terlihat berasal dari luar.

---

#### 3. Yang sebenarnya dapat dimulai: NOL — dan sebabnya seluruhnya di papan sebelah

Pemeriksaan `K-2` dijalankan: *untuk tiap tiket yang mengaku dapat dimulai, sebutkan entitas yang
disentuhnya dan tunjukkan tiket mana yang membuatnya.*

| Tiket | Entitas yang disentuhnya | Tiket yang membuatnya |
|---|---|---|
| `01`, `02`, `03`, `05` | `VERSI_KONTRAK` | **`14`** — di papan `treaty-in` |
| `04` | `DOKUMEN_ADDENDUM`, merujuk `VERSI_KONTRAK` | **`14`** — idem |

> **Papan ini tidak dapat menghitung angkanya sendiri.** Lima dari tiga belas tiketnya diblokir satu
> tiket yang tidak ada di sini, dan **delapan sisanya menggantung pada kelima itu**.
>
> Itu ongkos pemisahan papan, dan ia dicatat di sini supaya siapa pun yang membaca *"aktif 10"*
> tidak menyimpulkan sepuluh tiket dapat dikerjakan hari ini. **Nol dapat.**

**Begitu `14` mendarat di papan induk**, yang dapat dimulai di sini melompat ke **lima** — `01`,
`02`, `03`, `05`, dan `04` bila `DB-16a` sudah dijawab.

---

#### 4. Tepi yang melintasi kedua papan

Dipetakan sebelum pemisahan, bukan sesudah. **Enam tepi, dua arah:**

| Arah | Tepi |
|---|---|
| papan ini **diblokir** papan induk | `01` ← `14` · `02` ← `14` · `03` ← `14` · `04` ← `14` · `05` ← `14` |
| papan induk **diblokir** papan ini | **`40` ← `01`** — *nilai versi sebelumnya berdampingan, di-SELECT bukan disalin* |

**Tepi terakhir yang paling mudah hilang:** ia satu-satunya arah sebaliknya, dan ia ada di papan
yang **tidak memuat tiket `01`**. Ia tertulis di `treaty-in/5-tiket/issues/README.md` §penahan luar.

##### Rujukan lain antar-papan, di luar `Blocked by`

Tujuh, dan seluruhnya tetap sah — ia rujukan, bukan penghalang:

```
14 → 04     32 → 13     38 → 11     40 → 01, 06     44 → 08     46 → 07     53 → 07
```

---

#### 5. Aturan papan

Sama dengan papan induk, dan keduanya **tidak boleh berbeda**:

- **Status adalah medan di dalam berkas tiket**, bukan lokasi foldernya. Papan **dibangkitkan dari
  berkas** dan **tidak pernah menghapus tiket**.
- **Empat golongan, tidak ada yang kelima.**
- **Penahan dari luar papan ditulis dengan nama aslinya** — `DB-16a`, `DB-20`, dan sekarang juga
  **nomor tiket papan induk**.
- **Penahan yang selesai dicoret, tidak dihapus.**
- **Nomor yang lompat bukan kekeliruan**, dan **penomoran tidak disusun ulang** — `01`…`13` tetap
  `01`…`13` meski papannya pindah.
- **Pencacah diberi label bendanya.** Papan ini mencacah **tiket Adjustment**; jangan dijumlahkan
  dengan pencacah papan induk tanpa menyebut keduanya.

#### 6. Menunggui, tidak menahan

| # | Menunggu | Yang berubah bila jawabannya lain |
|---|---|---|
| `01` | `REV-3` | **keadaan mana yang dihitung "berlaku"** — bukan bahwa ada penunjuk dasar |
| `02` | `DB-20` | **titik beku** materialitas, dan titik beku milik `07` |
| `09` | `DB-16a` | **bentuk kolom** dokumennya, yang ikut `04` |

#### 7. Menahan sungguhan — tiga, dan ketiga pertanyaannya BELUM DIKIRIM

| # | Penahan | Kenapa ia menahan |
|---|---|---|
| `04` | `DB-16a` | bila dokumen ternyata dikirim ke luar, **jenis bernilai tiga** dan penanda kembali. Itu kolom, bukan penyetelan |
| `07` | `DB-20` | titik beku adalah **pokok** tiketnya |
| `08` | `DB-16b` | tanggalnya **adalah** pokok tiketnya, dan `KTV-2` menuntut kolomnya **dicabut sebelum data masuk** bila dibantah |

#### 8. Asumsi yang dianggap clear

Daftarnya **tetap satu**, di `treaty-in/5-tiket/ASUMSI-CLEAR.md` — ia tidak ikut dipisah, sebab
sebagian kodenya dipakai kedua papan (`KTV-A` dipakai 19 tiket papan induk, `DB-20` dipakai empat
tiket di sini).

Sapuan yang mengeluarkan daftarnya kini menuntut **dua folder**:

```
grep -r "DIASUMSIKAN-CLEAR(" treaty-in/5-tiket/issues/ treaty-in-adjustment/5-tiket/issues/
```

> **Itu ongkos kedua dari pemisahan papan**, dan ia lebih halus daripada yang pertama: sapuan yang
> hanya menyisir satu folder akan mengembalikan daftar yang **terlihat lengkap**.
