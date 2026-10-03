# 28: Penanda migrasi endorsemen dan kepatuhan lapisan

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

## Hasil & nilai pengguna

Baris hasil migrasi dibedakan dari baris hasil hitung, dan dua keadaan yang perlu diketahui
ditandai **tanpa mengubah satu angka pun**.

⭐ **Nilai lama tidak pernah dihitung ulang** — termasuk bila rumus lamanya diketahui keliru.
Yang dilakukan hanya **menandainya**, supaya kelak dapat diperlakukan berbeda bila diputuskan
demikian.

## Yang dibangun

Dua penanda, keduanya **hanya berarti untuk baris hasil migrasi**:

| Penanda | Artinya |
| --- | --- |
| **pasangan bergeser** | kunci dagang pasangan menurut nomor urut ternyata berbeda ⇒ ⚠️ **anomali sungguhan**, sebab di endorsemen baris tidak dapat dihapus |
| **rumus berlapis** | nilainya lahir dari rumus lama yang mengurangi terhadap selisih ⇒ dapat ditentukan pasti dari generasi sebelumnya |

Ditambah sembilan pernyataan kepatuhan yang berlaku sama seperti pada polis baru: satu transaksi ·
skema eksplisit · kolom pelaku dari identitas login · nol tipe mengambang · skala kolom uang · kode
dan penanda tetap teks · teks kosong menjadi tak-bernilai · arah ketergantungan · satu antarmuka
penyimpanan.

## Batas — yang TIDAK termasuk

⛔ Pemuat massalnya sendiri — tiket **22**.
⛔ Tabel proyeksi tempat penanda ini tinggal — tiket **26**.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Uji utama:** dokumen lama dimuat, penandaan dijalankan, lalu seluruh
nilai uang dibandingkan dengan sebelum penandaan — **satu angka yang berubah berarti gagal**.

## Acceptance criteria

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

## ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Lingkup pemindahan dokumen lama belum diputuskan** — sama dengan penahan tiket
**22**. Penandaan hanya berarti atas baris yang benar-benar dipindahkan.

`[work owner]` **Dua butir endorsemen belum dijawab:** apakah polis yang sudah dibatalkan masih
boleh di-endorse lagi, dan berapa kali satu polis boleh di-endorse *(batas teknis 99)*. Keduanya
menentukan baris mana yang sah ada di dalam rantai, dan karena itu baris mana yang ditandai.

⚠️ **AC 49 membawa pertentangan spec** yang dicatat di tiket **18** dan **27** — skala delapan
lawan sembilan. ⛔ Tidak diputuskan di sini.