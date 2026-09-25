# Papan tiket — Komite Claim Non Prop

**Tiket: aktif 1 · tertahan 11 · menunggu instance 0 · selesai 0 · mati 0**

> Pencacah ini menghitung **tiket**. Persyaratan `S-xxx` punya pencacahnya sendiri di `SPEC-KOMITE-01.md` bagian *Cakupan dan cacah*; keduanya diberi label bendanya supaya tidak pernah tertukar.

## Yang sebenarnya dapat dimulai

**1 dari 12 tiket di papan.** Sebelas sisanya tertahan **tiket lain**, bukan sesuatu di luar papan — rantai dependensinya nyata dan seluruhnya ada di sini. Itu keadaan yang berbeda dari papan lapisan data sisi Klaim, tempat seluruh sisa tertahan keputusan di luar papan.

Tidak ada penahan dari luar papan. Tiga tiket **menunggui** jawaban dari luar tanpa tertahan olehnya; daftarnya di bawah.

Diterbitkan 21 September 2026 dari `SPEC-KOMITE-01.md` bagian *Persyaratan* — 70 persyaratan `S-xxx`, **seluruhnya** diketiketkan sejak `PERISTIWA_ROSTER` ditambahkan sebagai objek ketujuh.
Tracker belum terkonfigurasi (`/setup-matt-pocock-skills` belum dijalankan), jadi papannya berkas lokal.

## Aturan papan

- **Status adalah field di dalam berkas tiket** (`status: aktif | tertahan | selesai | selesai-sebagian`), bukan lokasi foldernya. Pembangkit membaca field, dan **tidak pernah menghapus berkas tiket**.
- **Penahan dari luar papan ditulis dengan nama aslinya** — `C-01`, `PAGAR-01` — tidak diterjemahkan jadi nomor tiket. Yang menahan dari luar harus terlihat berasal dari luar.
- **Penahan yang sudah selesai dicoret, tidak dihapus**, supaya rantainya tetap terbaca: `~~05~~ *(selesai)*`.
- **Nomor yang lompat bukan kekeliruan.** Penomoran tidak disusun ulang, karena menyusunnya ulang memutus setiap rujukan yang sudah ada.
- **Uji hidup di dalam tiketnya**, bukan sebagai tiket tersendiri. Ini berbeda dari papan lapisan data sisi Klaim, dan sebabnya disebut di bagian terakhir.

## Papan

| # | Tiket | Ditahan oleh |
|---|---|---|
| `01` | Roster dan tabel seleksi berdiri sebagai data yang dikelola | — |
| `02` | Sirkulasi usulan pembayaran lahir lengkap dengan jenjangnya | `01` |
| `03` | Pembentukan yang gagal terlihat, dan klaim tidak membawa jejaknya | `02` |
| `04` | Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak | `02` |
| `05` | Akibat pada klaim dihitung di satu tempat dan ditulis dalam transaksi yang sama | `04` |
| `06` | Versi usulan naik saat penanda berubah, dan penjaga menolak versi yang berselisih | `04` |
| `07` | Sirkulasi menutup dan menolak klaim, pemegangnya diselesaikan dari peran | `05` |
| `08` | Nomor akseptasi terbit sekali, di dalam transaksi keputusan | `05` |
| `09` | Layar persetujuan komite dengan tiga masukan dan kunci yang ditegakkan dua kali | `04` |
| `10` | Giliran saya hari ini, dan umur yang tidak pernah disimpan | `04` |
| `11` | Lima port hilir berdiri kosong dengan niat tercatat lebih dulu | `05` · `08` |
| `12` | Pergantian pemegang jenjang tercatat sebagai peristiwa | `01` |

## Menunggui, tidak menahan

Tiket ini **boleh jalan**. Jawaban yang ditunggu mengubah satu hal yang tiketnya sudah menyebut — sebuah nama peran, sebuah tempat simpan, sebuah pembungkus — bukan rancangannya. Dipisahkan dari penahan sungguhan supaya papan tidak berteriak serigala.

| # | Menunggu | Yang berubah bila jawabannya lain |
|---|---|---|
| `01` | Nama peran administratif — keputusan pemilik proses | **Nama peran yang dicocokkan**, bukan adanya gerbang |
| `03` | Persetujuan pemilik tabel sisi Klaim atas kolom maksud pengajuan | **Tempat maksud disimpan**, bukan bahwa maksud dan akibat terpisah |
| `08` | `C-01` — isi badan prosedur penerbit pada basis data berjalan | **Pembungkusnya**, bukan bahwa nomor terbit di dalam transaksi |

## Dapat dimulai hari pertama

| # | Tiket |
|---|---|
| `01` | Roster dan tabel seleksi berdiri sebagai data yang dikelola |

Satu tiket, dan itu disengaja: tanpa roster dan tabel seleksi sebagai data, tidak ada sirkulasi yang dapat dibentuk sama sekali. Tiket `02` adalah peluru penjejak intinya dan terbuka begitu `01` hijau, bersama tiket `12`.

---

## Tidak diketiketkan

Dua daftar, dan pembedaannya penting. Yang pertama memuat **persyaratan** — baris `S-xxx` di
spesifikasi yang belum punya tiket. Yang kedua memuat **keputusan yang belum diambil**, yang
bukan persyaratan dan tidak akan pernah ditemukan dengan mencari `S-xxx` di spesifikasi.
Dicampur, pembaca berikutnya akan mengira ada persyaratan hilang dan mencarinya di tempat
yang tidak memuatnya.

### Persyaratan ditahan

**Kosong.** Ketujuh puluh persyaratan `S-xxx` seluruhnya punya tiket.

`S-052` pernah berada di sini, karena ia menuntut pergantian pemegang tercatat sebagai
peristiwa sementara satu-satunya tabel peristiwa mewajibkan pengenal sirkulasi — dan
pergantian roster tidak punya sirkulasi. Objek ketujuh `PERISTIWA_ROSTER` ditambahkan
21 September 2026, `S-052` lolos `T-1`, dan ia kini menjadi tiket `12`. Judul daftar ini
tidak dihapus: kosong adalah keadaan yang perlu terbaca, bukan bagian yang perlu hilang.

### Keputusan belum diambil

Bukan persyaratan, dan **bukan lubang bukti**. Ketiganya menunggu orang, bukan berkas.

| Yang ditunggu | Keadaannya |
|---|---|
| **Nama peran administratif** pengelola roster dan tabel seleksi | Keputusan pemilik proses. `H-4` menyebut "peran administratif tersendiri di luar keanggotaan komite" sebagai **default**, bukan sebagai ketetapan. Gerbangnya sudah dispesifikasikan dan dibangun di tiket `01`; yang menunggu hanya namanya, dan itu **menunggui, tidak menahan** |
| **Bentuk migrasi sirkulasi berjalan** | Satu cacah yang belum diambil: berapa klaim di produksi membawa pengenal sirkulasi yang bukan miliknya — `INVENTARIS-BUKTI.md` §2.5 baris 6. Yang **sudah** diputuskan ada di tiket `02` sebagai `S-070` |
| **Port hilir bagi `KonversiKlaim_Act` dan `InsertJsonClaimTreatyNonProp_act`** | Keputusan pemilik proses. Berkas keempat rule-nya **ada di repo** (`INVENTARIS-BUKTI.md` §2.4), sehingga keduanya **lubang rancangan, bukan lubang bukti**, dan berkeadaan `BELUM DIPUTUSKAN` — bukan `DIPAGARI` |

## Delapan pagar, nol tiket

Aliran A-4, A-5, dan **isi** A-6 tidak punya tiket sama sekali, dan itu benar. Yang dibangun hanyalah **pintunya** — tiket `11`, lima port kosong. PAGAR-01 menahan satu nilai di dalam tiket `08`, bukan tiketnya.

Pagar berlaku bagi **alasan**, bukan hanya bagi kesimpulan: tiket mana pun yang argumennya bersandar pada isi aliran beku tidak sah, meski kesimpulannya berada di aliran terbuka.

## Dua penyimpangan dari papan lapisan data sisi Klaim

Konvensi berkasnya ditiru apa adanya. Dua hal berbeda, dan sebabnya dicatat supaya tidak terbaca sebagai kelalaian.

**Uji tidak menjadi tiket tersendiri.** Papan sisi Klaim memberi delapan uji tiketnya masing-masing, bertanda `menunggu-instance`, karena batch-nya lapisan data saja — ujinya memang tidak dapat hijau sampai DDL berjalan di instance nyata. Batch ini memotong seluruh lapisan sekaligus, sehingga ujinya hijau di dalam tiketnya sendiri. Memisahkannya justru melanggar aturan irisan tegak.

**Tiap tiket membawa bagian `Persyaratan`, `Jalur gagal`, `Uji`, dan `Menggantikan`.** Keempatnya tidak punya tempat di konvensi sisi Klaim, dan ditambahkan sebagai bagian baru — bukan dibuang. Ketiganya yang pertama datang dari lapisan `S-xxx`; yang terakhir datang dari bagian *Ketertelusuran*, dan tanpanya shadow-run kehilangan dasar pembandingnya.
