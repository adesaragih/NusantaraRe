# 09: Dua penanda migrasi

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

## Hasil & nilai pengguna

Baris hasil migrasi dibedakan dari baris hasil hitung, dan dua keadaan yang perlu diketahui
**ditandai tanpa mengubah satu angka pun**.

⭐⭐ **Nilai lama tidak pernah dihitung ulang** — termasuk bila rumus lamanya diketahui keliru.
Yang dilakukan hanya **menandainya**, supaya kelak dapat diperlakukan berbeda bila diputuskan
demikian. ⛔ Menghitung ulang berarti laporan lama tidak lagi dapat direkonsiliasi.

## Yang dibangun

Dua penanda, keduanya **hanya berarti untuk baris hasil migrasi**:

| Penanda | Artinya | Kenapa berguna |
| --- | --- | --- |
| **pasangan bergeser** | kunci dagang pasangan menurut nomor urut ternyata berbeda | ⭐ di endorsemen baris tidak dapat dihapus ⇒ ini **anomali sungguhan** |
| **rumus berlapis** | nilainya lahir dari rumus lama yang mengurangi terhadap selisih | dapat ditentukan **pasti** dari generasi sebelumnya |

⭐ Pemasangan baris antar generasi memakai **nomor urut**; kunci dagang turun pangkat menjadi
**pemeriksa**, bukan pemasang.

## Batas — yang TIDAK termasuk

⛔ Pemuat massalnya sendiri — tiket **10**.
⛔ Tabel proyeksi tempat penanda ini tinggal — tiket **07**.
⛔ Perhitungan ulang nilai lama — ⛔ **tidak dilakukan sama sekali**.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Uji utama:** dokumen lama dimuat, penandaan dijalankan, lalu
**seluruh nilai uang dibandingkan dengan sebelum penandaan** — satu angka yang berubah berarti
gagal.

## Acceptance criteria

- [ ] **AC 39** — nilai hasil migrasi **tidak dihitung ulang**, termasuk digit galatnya
- [ ] **AC 40** — pemasangan antar generasi memakai **nomor urut**
- [ ] **AC 41** — penanda pasangan bergeser dihitung **tanpa mengubah satu angka pun**
- [ ] **AC 42** — penanda rumus berlapis terisi pada setiap baris yang memenuhi syaratnya
- [ ] **AC 43** — kedua penanda **hanya** terisi pada baris hasil migrasi

## ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Lingkup pemindahan dokumen lama belum diputuskan** — seluruhnya, sebagian, atau
tetap dibaca lewat jalur lama. Penandaan hanya berarti atas baris yang benar-benar dipindahkan.