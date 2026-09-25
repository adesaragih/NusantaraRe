# Jenis aturan yang NOL kemunculan di ekspor — dan lima yang mengubah klaim kita

**Tanggal:** 24 September 2026
**Sebab:** *"tidak ada di ekspor" bukan "tidak ada di sistem"* — `CONTEXT.md` §2.0-h

> Sapuan ini **tidak menunggu siapa pun** dan dijalankan hari itu juga. Jenis aturan diambil dari
> `pxObjClass` pertama tiap berkas — blok `Rule-*` di kepalanya, **yaitu dari ISI berkas, bukan dari
> nama foldernya.**

> ### KOREKSI 24 September 2026 — laporan pertama menjelaskannya lewat FOLDER
>
> Sapuannya memang sudah berbasis isi. **Narasinya yang keliru:** saya menulis *"ekspornya berfolder
> menurut jenis, jadi jenis tanpa folder memang nol berkas"* — dan itu menjadikan **susunan folder**
> sebagai dasar, padahal buktinya lebih kuat daripada itu.
>
> Bedanya bukan gaya. Bila foldernya dibuat **sesudah** ekspor oleh skrip penyortir, *"tidak ada
> folder"* punya sebab kedua: **skripnya tidak mengenali jenis itu**, sehingga berkasnya tidak
> dipindahkan ke mana pun. Ini `CONTEXT.md` §2.0-h satu tingkat lebih dalam — *"tidak ada foldernya"*
> bukan *"tidak ada di ekspor"*.
>
> **Dua pemeriksaan dijalankan, dan keduanya bersih:**
>
> | Yang diperiksa | Hasil |
> |---|---|
> | berkas XML **di luar** folder jenis mana pun — tertinggal di akar | **NOL**, di kedua modul. Tidak ada berkas lain berjenis apa pun di akar |
> | kecocokan **folder ↔ `pxObjClass` isinya** | **satu-ke-satu sempurna**, dua belas folder, dua belas jenis, **nol** berkas yang jenisnya tidak cocok dengan foldernya |
>
> **Batasnya tetap disebut:** bila skrip penyortir **membuang** yang tidak dikenalinya alih-alih
> meninggalkannya, hasil bersih ini akan terlihat sama. Yang terbantah hanya sebab *"tertinggal di
> folder utama"*; sebab *"dibuang penyortir"* hanya dapat ditutup dengan bertanya. **Daftar 21
> berdiri, dan sekarang ia berdiri di atas isi berkas.**

---

## 1. Yang MUNCUL

| Jenis | Berkas |
|---|---:|
| `Rule-Obj-Activity` | 304 |
| `Rule-HTML-Section` | 118 |
| `Rule-Connect-SQL` | 84 |
| `Rule-Obj-FlowAction` | 75 |
| `Rule-Obj-Model` *(data transform)* | 73 |
| `Rule-Obj-Report-Definition` | 34 |
| `Rule-HTML-Harness` | 9 |
| `Rule-Obj-When` | 4 |
| `Rule-Connect-REST` · `Rule-Declare-DecisionTable` · `Rule-Admin-System-Settings` | 2 each |
| `Rule-Navigation` | 1 |

**Dua belas jenis**, dihitung dari **`pxObjClass` isi berkas**. Ekspornya kebetulan juga berfolder
menurut jenis, dan **kecocokan folder ↔ isi diperiksa satu-ke-satu** — tetapi folder bukan dasar
angka ini; isi berkaslah dasarnya.

---

## 2. Yang NOL — dua puluh satu jenis

| Jenis | Apa yang hilang bersamanya | Berat |
|---|---|---|
| **`Rule-Declare-Expressions`** | **menulis nilai tanpa dipanggil siapa pun** | **berat** |
| **`Rule-Declare-Trigger`** | idem, dipicu perubahan data | **berat** |
| **`Rule-Obj-Flow`** | orkestrasi — urutan tugas dan siapa pemegangnya | **berat** |
| **`Rule-Obj-CaseType`** | daur hidup kasus, tahapan | **berat** |
| **`Rule-Access-Role-Obj`** · `Rule-Access-Deny-Obj` | **wewenang** — siapa boleh apa | **berat** |
| `Rule-Obj-ServiceLevel` | tenggat dan pengingat | sedang |
| `Rule-Obj-Validate` · `Rule-Declare-Constraints` | validasi deklaratif | sedang |
| `Rule-Declare-OnChange` · `Rule-Declare-Index` | reaksi perubahan, indeks terdeklarasi | sedang |
| `Rule-Agent-Queue` | pekerjaan latar | sedang |
| `Rule-Obj-Property` | **definisi properti** — sudah diketahui tidak ada | ringan |
| `Rule-Obj-MapValue` · `Rule-Obj-ListView` · `Rule-Message` · `Rule-Obj-Corr` · `Rule-HTML-Property` · `Rule-File-Binary` · `Rule-Application` · `Rule-Obj-Class` | pemetaan, pesan, korespondensi, kelas | ringan |

---

## 3. LIMA yang mengubah klaim yang sudah kita buat

### 3.1 `Declare Expression` dan `Declare Trigger` — setiap "tidak pernah ditulis" punya lubang

Keduanya **menulis nilai tanpa dipanggil siapa pun**. Seluruh sapuan penulis kita mencari
`<PropertiesName>` di aktivitas dan *data transform* — **bentuk itu tidak akan pernah melihat
keduanya.**

**Klaim yang terdampak, dan masing-masing perlu dibaca ulang dengan lubang ini di kepalanya:**

| Klaim | Di mana |
|---|---|
| `ReisuredParticipant` **"nol aturan menulisnya"** | `SPEC-INVARIAN.md` §4.4a — dipakai untuk membuang sumbu pihak |
| `ReasTreatyInGroupLeader` **"disetel nol kali"** | butir 7.10 |
| pembagian **masukan versus turunan** — 187 masukan, 119 daun turunan | `SPEC-MODEL-DATA.md` §3, §4 |
| ~~`CedingStatusActive` dan `SourceStatusActive` "tidak pernah ditulis satu aturan pun"~~ | §14.5 — **TERBANTAH 24 Sep 2026, lihat §3.1b** |
| `TREATY_IN` **"nol `INSERT`/`UPDATE` langsung dari aturan"** — INV-60 aman | `PRA-PEMISAHAN…` §1c |

> **Tidak satu pun dari kelimanya otomatis salah.** Yang berubah: masing-masing berdiri di atas
> ketiadaan pada semesta yang **kurang satu jenis penulis**.

### 3.1b Yang paling dicurigai ternyata SALAH — dan sebabnya bukan L-10

§14.5 ditandai paling patut dicurigai karena di situ ketiadaan penulis dipakai sebagai **alasan
membuang**. Ia diperiksa, dan **klaimnya memang salah** — tetapi bukan karena jenis aturan yang
tidak terekspor. **Karena bentuk penulisan yang tidak pernah kami cari.**

`Harness/InputTreatyInOffer.xml`:

```xml
<pyDisplayProperty>.StatusActive</pyDisplayProperty>
<pyPropertyTarget>TreatyIn.CedingStatusActive</pyPropertyTarget>
<pySetValueOnSelect>true</pySetValueOnSelect>
```

Itu **pemetaan hasil pencarian di layar**: ketika orang memilih satu cedant dari master, kolom
`StatusActive` baris itu **disalin** ke kontrak. Tidak ada `Property-Set` di mana pun — dan seluruh
sapuan penulis kami hanya mengenali `<PropertiesName>`.

| | |
|---|---|
| Sasaran `pyPropertyTarget` berawalan `TreatyIn.` di seluruh korpus | **empat** — `CedingID`, `CedingStatusActive`, `LeadingReinsSourceID`, `SourceStatusActive` |
| Berkas yang memakai bentuk ini | 35 |
| Lubangnya sistemik? | **tidak** — empat sasaran, dan dua di antaranya (`CedingID`, `LeadingReinsSourceID`) memang sudah dipertahankan |

> ### Dan hasilnya membalik arah: §14.5 sekarang LEBIH KUAT
>
> Alasan **pertama** — *"tidak pernah ditulis"* — **dicabut, dan ia memang salah.**
>
> Alasan **kedua** — *"salinan status data master"* — ditanyakan golongannya sesuai
> `BENTUK-TIKET.md` §7.2, dan jawabannya **`EVIDENCED`, bukan `DECIDED`**: pemetaan di atas
> **memperlihatkan penyalinannya terjadi**, dari baris master, pada saat dipilih, dengan berkas dan
> baris yang dapat ditunjuk.
>
> **Putusan membuang kedua atribut BERTAHAN, di atas dasar yang terbukti alih-alih dua dasar yang
> salah satunya runtuh.** §14.5 **keluar** dari daftar yang harus dilihat ulang saat L-10 dijawab.

**Pelajaran yang lebih luas daripada §14.5 sendiri:** sebelum menyapu *"siapa menulis X"*,
**daftarkan dulu BENTUK-BENTUK penulisan yang mungkin** — penugasan langsung, pemetaan hasil
pencarian, pemetaan koneksi basis data, ekspresi terdeklarasi, pemicu terdeklarasi. Sapuan satu
bentuk menjawab pertanyaan yang lebih sempit daripada yang ditanyakan, dan **diamnya terbaca sebagai
ketiadaan**. `CONTEXT.md` §2.0-i.

### 3.2 `Flow` dan `CaseType` — mesin keadaan mungkin kurang

ADR-0055 disusun dari aktivitas dan *data transform*. Bila orkestrasi yang sebenarnya hidup di
`Rule-Obj-Flow`, maka **tiga belas perpindahan itu adalah perpindahan yang TERLIHAT DARI SISI
AKTIVITAS**, bukan tentu seluruhnya.

Dan ia menyambung langsung ke **7.10**: `Flow` adalah tempat yang **wajar** untuk menyetel jabatan
pemegang tugas. Bila `Flow` ada dan tidak terekspor, **bacaan A menjadi jauh lebih mungkin** — dan
ia tidak lagi bersandar pada *"sesuatu dinonaktifkan karena pernah menyala"*.

### 3.3 `Access-Role` — **P-41 TIDAK berubah.** Koreksi 24 September 2026

Saya menulis bahwa P-41 berubah dari *"belum ditemukan"* menjadi *"dapat diminta"*. **Itu terlalu
cerah, dan penghalangnya tidak boleh dipindahkan.**

`Rule-Access-Role-Obj` mengatur **siapa boleh melakukan apa terhadap sebuah kelas** — membuka,
mengubah, menghapus. Ia **tidak memuat ambang nilai**. *"Direktur boleh menyetujui kontrak sampai
sekian miliar"* bukan bentuk yang muat di sana; ia aturan bisnis, dan tempatnya di **kebijakan
perusahaan** atau di tabel acuan.

> **P-41 TETAP TERTAHAN pada dokumen batas wewenang, dan penghalangnya tetap MANAJEMEN, bukan
> ekspor.** Penghalang yang dipindah ke tempat yang salah akan **ditunggu di tempat yang salah**.

**Permintaan aturan wewenang tetap berjalan** — ia murah dan ikut daftar L-10 — tetapi untuk hal
lain, dan ketiganya memang menyangkut siapa boleh apa:

| | |
|---|---|
| **P-36** | pemberi keputusan bukan pengisi |
| **P-40** | tidak ada jalan pintas |
| **eskalasi butir 1** | persetujuan tanpa penyetuju |

Ketiganya **lebih kuat dengan daftar peran di tangan**, dan tidak satu pun adalah P-41.

---

## 4. Yang diminta, dan ia MURAH

Ditambahkan ke permintaan ekspor kedua yang sudah berjalan (`LUBANG-SPESIFIKASI.md` §8):

> **ADAKAH ATURAN BERJENIS `Rule-Obj-Flow`, `Rule-Obj-CaseType`, `Rule-Declare-Expressions`,
> `Rule-Declare-Trigger`, dan `Rule-Access-Role-Obj` UNTUK APLIKASI INI? BILA ADA, BERAPA BANYAK
> DAN NAMANYA APA?**
>
> **Tidak perlu isinya. Cukup DAFTARNYA.**

| Jawaban | Artinya |
|---|---|
| **nol** | kesimpulan lama berdiri — dan sekarang ia berdiri **di atas jawaban, bukan di atas ketiadaan** |
| **bukan nol** | temuan terbesar sejak L-8, dan **cakupannya langsung terukur dari daftarnya** |

**Penagihnya:** L-10 di `LUBANG-SPESIFIKASI.md`, dan kelima klaim di §3.1 yang masing-masing
menyebut lubang ini di tempatnya sendiri.
