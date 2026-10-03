# 00: Skema penyimpanan komite — tiga tabel + dua kolom usul — **PREFACTOR**

**Status:** ready-for-agent

**Blocked by:** **skema klaim Claim Prop** — `T_WORK_CLAIM` harus sudah ada, karena baris komite
tinggal di tabel yang sama. Bukan penyelesaian modul Claim Prop, hanya tabelnya.

## Hasil & nilai pengguna

Sebuah kasus komite dapat **dibuat, disimpan, dan dibaca kembali** — beserta daftar penyetujunya dan
penunjuk ke baris penyesuaian induknya. Tanpa ini tidak ada tiket lain yang bisa menyimpan apa pun.

Ini **prefactor**: tidak ada perilaku pengguna yang terlihat, tetapi seluruh tiket berikutnya
bersandar padanya.

## Bentuk yang dibangun

⛔ **Bentuk ditulis dalam kalimat.** Acuan tunggalnya `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md`.

**`T_WORK_CLAIM` — baris komite.** Tabel **lintas-lini** yang sudah ada; modul ini hanya menambah
**baris**, bukan kolom. Baris komite ber-`ID` berawalan **`TKMT-`**, ber-`COVER_KEY` berisi `ID`
baris klaim induk yang berawalan **`CLMP-`**, dan ber-`LINI` **`PROP`**.

**`T_GENERAL_KOMITE` — header kasus komite.** Satu baris per penyerahan satu baris penyesuaian ke
komite. `ID`-nya **sama persis** dengan baris komite di `T_WORK_CLAIM` — **shared primary key**,
tanpa kolom penyambung. Kolomnya: penunjuk ke baris penyesuaian induk (**wajib isi**, index
**unik**), pencacah **berapa penyetuju yang dibutuhkan**, pencacah **penyetuju ke berapa yang
sedang berjalan**, **hasil akseptasi terakhir**, dan **dua penanda usul** dari tiket 05.

⭐ **Kedua kolom penanda usul — tipe dan nilainya sudah dikunci.** `[keputusan work owner]`
2026-09-19 — keduanya **teks satu huruf**: **`CHAR(1)`**, **wajib isi** *(NOT NULL)*, dan
**nilai sah hanya `'1'` dan `'0'`** — `'1'` bila diusulkan, `'0'` bila tidak. **Kotak centang yang
tidak pernah disentuh tersimpan `'0'`, bukan kosong.** Nilai di luar keduanya **ditolak**.
⚠️ Migrasi data lama: benar/salah → `'1'`/`'0'`, kosong → `'0'` *(tiket 13)*.

**`T_KOMITE_KOMITELIST` — satu baris per penyetuju.** Kolomnya: kunci baris, penunjuk ke header
kasus komite, **urutan penyetuju**, **akun operator**, **jabatan**, **email**, **keputusan**,
**catatan**, dan **tanggal diputuskan**.

⚠️ **Ketiga pencacah — jumlah penyetuju, penyetuju keberapa, urutan penyetuju — adalah bilangan
bulat biasa**, tidak ikut aturan ketelitian angka.

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml` · `ASM-FW-GCNMFW-DATA-ADJUSTMENT!ADDKOMITETREATYCHILD_ACT` | `[terverifikasi]` **langkah 13** mengisi penunjuk posisi baris penyesuaian · **langkah 14** menyalin isi baris penyesuaian ke halaman komite · **langkah 22.1** membangun daftar penyetuju dengan keputusan awal kosong |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` · `ASM-FW-GCNMFW-WORK-KOMITETREATY!KOMITEPOSTADJUSTMENT` | `[terverifikasi]` **langkah 3** menyalin penunjuk posisi ke variabel kerja · **langkah 4** membuka kasus klaim induk dengan kunci induk |

⭐ `[terverifikasi]` **Halaman penyesuaian pada kasus komite adalah SATU, bukan daftar** — dasar
keputusan menaruh dua kolom usul di header kasus.

## Jalur baca 17 properti penyesuaian

`[keputusan work owner]` 2026-09-18 — **17 properti yang hanya dibaca TIDAK disalin** ke tabel
komite; keduanya dibaca dari baris penyesuaian milik Claim Prop.

Di Pega kuncinya **dua bagian dan bersifat posisional**: kunci kasus klaim induk, ditambah **nomor
urut baris** di dalam daftar penyesuaian klaim itu.

⭐ **Di sistem baru, penunjuk baris penyesuaian di header kasus komite menggantikan pasangan kunci
itu** — itulah sebabnya ia **wajib isi** dan ber-index **unik**.

## Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | Nama kolom **snake_case**; berkas struktur tabel adalah **acuan tunggal** nama kolom |
| 2026-09-18 | Pencacah jumlah penyetuju adalah **kolom yang disimpan**, bukan dihitung ulang tiap dibaca |
| 2026-09-18 | **Satu kolom tanggal persetujuan saja** — properti tanggal kedua di Pega sengaja tidak dijadikan kolom |
| 2026-09-18 | 17 properti penyesuaian **dibaca**, 2 penanda usul **disimpan** |
| **2026-09-19** | **Dua kolom usul disimpan di header kasus komite** — *"komite menyimpan catatan usulnya sendiri, terpisah dari nilai akhir yang tercatat di kasus klaim"* |
| **2026-09-19** | **Batas digit angka mengikuti bentuk kolom yang sudah ada apa adanya** |

## ⚠️ TITIK YANG SENGAJA DIUBAH — ketelitian angka

`[keputusan work owner]` 2026-09-18, diralat 2026-09-19. **Ini perubahan sadar, bukan peniruan.**

| Di mana | Ketetapan |
| --- | --- |
| Kolom uang, persen, kurs | **mengikuti bentuk yang sudah ada — 20 digit seluruhnya, 8 di belakang koma** |
| Pencacah | **bilangan bulat biasa** |
| Hitungan | sampai **20 angka di belakang koma**, nol pembulatan di tengah jalan |
| Tampilan | **4 angka di belakang koma** |

⚠️ **Pembulatan terjadi di batas penyimpanan, diterima sadar.** Pega memakai **empat rupa
ketelitian** tanpa pembulatan sama sekali; hasil sistem baru **akan berbeda di angka belakang
koma**, dan itu disengaja.

⚠️ `[keputusan work owner]` 2026-09-19 — **bila sebuah nilai melewati batas digit, basis data
menolak menyimpan, bukan membulatkan.** Kegagalan itu **harus terlihat**, tidak boleh ditelan.

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 1 · 40 · 42 · 54 · 55 · 60 · 61 · 71 · 72 · **74**

- [ ] Satu kasus komite dapat dibuat dan dibaca kembali lengkap dengan daftar penyetujunya.
- [ ] Baris komite dan header kasus komite memakai **kunci yang sama persis**; test yang menemukan
      kolom penyambung terpisah **gagal**.
- [ ] Baris komite ber-`ID` berawalan `TKMT-`, ber-`COVER_KEY` berawalan `CLMP-`, ber-`LINI` `PROP`.
- [ ] Penunjuk ke baris penyesuaian **wajib isi** dan **unik**; dua kasus komite tidak bisa menunjuk
      baris penyesuaian yang sama.
- [ ] Ketiga pencacah tersimpan sebagai **bilangan bulat**, bukan angka desimal.
- [ ] Kedua kolom penanda usul menerima **hanya** `'1'` dan `'0'`; nilai lain **ditolak**, dan
      baris tanpa nilai **ditolak** — kotak yang tidak disentuh tersimpan `'0'`, bukan kosong.
- [ ] Nilai uang melewati aplikasi **tanpa melewati bilangan pecahan biner**.
- [ ] Nilai yang melewati batas digit **ditolak dengan galat yang terlihat**, bukan diam-diam
      dibulatkan atau ditelan.

## Butir `[terbuka]` yang menyentuh tiket ini

- ✅ ~~**Nilai apa yang tersimpan** untuk kedua penanda usul~~ — **TERJAWAB**
  `[keputusan work owner]` 2026-09-19: **`CHAR(1)`, wajib isi, nilai sah hanya `'1'` dan `'0'`;
  kotak yang tidak disentuh tersimpan `'0'`.** Lihat blok bentuk tabel di atas.
  > ⛔ **Kalimat lamanya dikutip, tidak dihapus:** *"**Nilai apa yang tersimpan** untuk kedua
  > penanda usul — di Pega ia kotak centang, tetapi nilai tersimpannya **tidak terbaca dari
  > korpus**."*
  >
  > ⚠️ Butir ini **tidak pernah masuk register 13 butir** di `spec.md`, jadi menjawabnya
  > **tidak mengubah jumlah register**. ⛔ **Nol butir register ditutup.**

## Seam & verifikasi

`[keputusan work owner]` 2026-09-18 — **memakai ulang seam Claim Prop, tidak menambah seam baru.**
⚠️ Tiket ini **belum bisa diverifikasi ujung-ke-ujung** sampai seam Claim Prop berdiri.
