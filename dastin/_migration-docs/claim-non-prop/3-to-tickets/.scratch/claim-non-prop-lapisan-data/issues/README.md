# Papan tiket — lapisan data Claim Non Prop

**Tiket: aktif 0 · menunggu instance 8 · tertahan 1 · selesai 29 · mati 1**

> Pencacah ini menghitung **tiket**. Register REQ punya pencacahnya sendiri dengan format yang mirip; keduanya diberi label bendanya supaya tidak pernah tertukar.

## Yang sebenarnya dapat dimulai

**0 dari 9 tiket di papan.** Sisanya tertahan sesuatu di **luar papan** — keputusan yang belum dikonfirmasi atau REQ yang belum kembali — bukan tertahan tiket lain. Angka ini, bukan "aktif 0", yang menggambarkan keadaan proyek.

| Penahan dari luar papan | Tiket yang dibekukannya | Ditahan langsung |
|---|---|---|
| **lingkup-batch** | 1 | langsung: `11` |

**Tidak ada yang dapat dimulai sekarang.** Empat sapuan sumber yang tidak menyentuh gerbang sudah dikerjakan sampai habis; sisanya menunggu konfirmasi ADR-0028 dan REQ-032.

Diterbitkan 18 September 2026 dari `_migration-docs/claim-non-prop/TICKETS.md`, dibangun ulang dari field `status:` tiap berkas.
Tracker belum terkonfigurasi (`/setup-matt-pocock-skills` belum dijalankan), jadi papannya berkas lokal.

## Aturan papan

- **Status adalah field di dalam berkas tiket** (`status: aktif | tertahan | selesai | selesai-sebagian`), bukan lokasi foldernya. Folder `_selesai/` dan `_tertahan/` hanya kerapian; pembangkit membaca field, dan **tidak pernah menghapus berkas tiket**.
- **Penahan dari luar papan ditulis dengan nama aslinya** — `ADR-0028`, `REQ-032` — tidak diterjemahkan jadi nomor tiket. Yang menahan dari luar harus terlihat berasal dari luar.
- **Penahan yang sudah selesai dicoret, tidak dihapus**, supaya rantainya tetap terbaca: `~~05~~ *(selesai)*`.
- **Nomor yang lompat bukan kekeliruan.** Penomoran tidak disusun ulang, karena menyusunnya ulang memutus setiap rujukan yang sudah ada.

## Papan

| # | Asal | Tiket | Ditahan oleh |
|---|---|---|---|
| `11` | T-38 | Tipe desimal di sisi Golang untuk nilai uang | **lingkup-batch** |
| `24` | T-26 | Uji: batas tanggal inklusif | ~~`12`~~ |
| `32` | T-24 | Uji: kunci alami akseptasi | ~~`18`~~ |
| `33` | T-27 | Uji: NULL bukan nol | ~~`16`~~ |
| `34` | T-31 | Uji: paritas | ~~`23`~~ ~~`15`~~ |
| `36` | T-25 | Uji: pasangan uang dan mata uang | ~~`16`~~ ~~`28`~~ |
| `37` | T-28 | Uji: penularan keadaan menunggu kurs | ~~`17`~~ ~~`18`~~ ~~`26`~~ |
| `38` | T-30 | Uji: bentuk view kompatibilitas | ~~`29`~~ ~~`05`~~ |
| `39` | T-29 | Uji: satu pintu tulis | `35` |

## Menunggui, tidak menahan

Tiket ini **boleh jalan**. Jawaban yang ditunggu mengubah satu hal yang tiketnya sudah menyebut — sebuah constraint, sebuah aturan penolakan, sebuah jalur migrasi — bukan rancangannya. Dipisahkan dari penahan sungguhan supaya papan tidak berteriak serigala.

| # | Menunggu | Yang berubah bila jawabannya lain |
|---|---|---|

## Dapat dimulai hari pertama

**Tidak ada.** Setiap tiket yang tersisa tertahan sesuatu di luar papan. Itu keadaan yang sah dan bukan kebuntuan kerja: yang dapat dikerjakan tanpa menyentuh gerbang sudah dikerjakan sampai habis.

---

## Menunggu instance — bukan aktif

Rancangannya lengkap; yang kurang mesinnya. Uji ini baru dapat hijau setelah DDL dijalankan di instance nyata, dan itu di luar batch lapisan data. Ditandai begini supaya papan tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun.

| # | Tiket |
|---|---|
| `24` | Uji: batas tanggal inklusif |
| `32` | Uji: kunci alami akseptasi |
| `33` | Uji: NULL bukan nol |
| `34` | Uji: paritas |
| `36` | Uji: pasangan uang dan mata uang |
| `37` | Uji: penularan keadaan menunggu kurs |
| `38` | Uji: bentuk view kompatibilitas |
| `39` | Uji: satu pintu tulis |

## Tertahan di luar papan

| # | Tiket | Ditahan |
|---|---|---|
| `11` | Tipe desimal di sisi Golang untuk nilai uang | **lingkup-batch** |

## Mati — dibatalkan, bukan ditunda

Tiket yang premisnya gugur. Berkasnya tidak dihapus: tiket yang pernah ada adalah bukti bahwa premisnya diperiksa, bukan dilewatkan. Nomornya tidak dipakai ulang.

| # | Tiket | Sebab |
|---|---|---|
| `02` | Tabel singkatan tertutup ditulis dan dibekukan | premis gugur — lihat kepala berkasnya di `_mati/` |

## Selesai

| # | Tiket | Status |
|---|---|---|
| `01` | Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis | selesai |
| `03` | Penamaan yang tidak nyaris kembar | selesai |
| `04` | Index penopang setiap foreign key | selesai |
| `05` | Sapuan S1: kolom yang sesungguhnya diterima Arasapas dan kasir | selesai |
| `06` | Keadaan kasus: nilai mana yang ada, dan siapa yang membacanya | selesai |
| `07` | Sapuan S3: perilaku per nilai `PaymentType` | selesai |
| `08` | Sapuan S4: adakah padanan IDR untuk biaya penilaian, salvage, dan biaya lain | selesai |
| `09` | Sapuan S9: apa persisnya yang ditulis `EditXOLAlokasi` | selesai |
| `10` | Sapuan S5, S6, S7, S8, dan dua sumber yang belum pernah dibaca | selesai |
| `12` | Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap | selesai |
| `13` | Bentuk lama boleh masuk utuh, di satu tempat saja | selesai |
| `14` | Nilai yang tidak dapat diurai tercatat, tidak dibulatkan, tidak dibuang | selesai |
| `15` | Jembatan ke sistem lama berdiri sebagai tabel terpisah | selesai |
| `16` | Nilai kerugian per mata uang, dengan nilai IDR yang tidak dapat lahir tanpa kurs | selesai |
| `17` | Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang | selesai |
| `18` | Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran | selesai |
| `19` | Masukan mesin alokasi dan penyebaran tersimpan, bukan hanya keluarannya | selesai |
| `20` | Objek, kronologi, dan dokumen klaim | selesai |
| `21` | Tanggal tutup buku dan tarif menjadi data bertanggal berlaku | selesai |
| `22` | Koreksi bernilai tercatat menggantikan tambalan di dalam kode | selesai |
| `23` | Setiap baris baru dapat ditunjuk balik ke barisnya di sistem lama | selesai |
| `25` | Arsip muatan keluar: apa yang benar-benar dikirim, tersimpan sebagai catatan | selesai |
| `26` | Premi pemulihan dapat dihitung ulang dari barisnya sendiri | selesai |
| `27` | Domain kolom keadaan | selesai |
| `28` | Pasangan IDR–kurs | selesai |
| `29` | Arasapas dan kasir menerima bentuk lama, dan yang tidak dapat diringkas tidak hilang diam-diam | selesai |
| `30` | Rekap per klaim per mata uang, dengan Retensi Cedant yang tidak pernah tercampur | selesai |
| `31` | Lima PageList lama tersedia sebagai view, tanpa menyimpan hasil yang dapat menyimpang | selesai |
| `35` | Satu pintu tulis ditegakkan hak akses, bukan kesepakatan | selesai |
