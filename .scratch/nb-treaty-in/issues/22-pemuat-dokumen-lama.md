# 22: Pemuat dokumen lama

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

## Hasil & nilai pengguna

Dokumen polis lama dimuat ke dalam tabel baru **lewat antarmuka yang sama** dengan jalur biasa,
sehingga data lama melewati pemeriksaan yang sama dan tidak ada pintu belakang.

⭐⭐ **Nilai lama tidak pernah dihitung ulang** — disalin apa adanya, lengkap dengan galat
presisinya, supaya laporan lama masih dapat direkonsiliasi.

## Yang dibangun

Pemuat massal yang: membaca dokumen lama · memecahnya dengan pemecah yang sama · menulis lewat
antarmuka penyimpanan yang sama · dan **melaporkan** dokumen yang gagal diurai beserta sebabnya.

⛔ **Dokumen yang gagal tidak dilewati diam-diam.**

## Batas — yang TIDAK termasuk

⛔ Penanda migrasi khas endorsemen — tiket **28**.
⛔ Keputusan lingkup pemindahan — `[work owner]`, lihat di bawah.

## Cara mengujinya

Lewat seam `repository` yang sama. ⭐ Uji utama: dokumen bergalat presisi dimuat, lalu nilainya
**dibaca dari kolomnya** — pembulatan ke presisi mata uang berarti gagal.

## Acceptance criteria

- [ ] **AC 55** — dokumen lama dimuat **tanpa pembulatan ke presisi mata uang**
- [ ] **AC 56** — pemuat menulis lewat antarmuka penyimpanan yang **sama**
- [ ] **AC 57** — medan tak dikenal **tersimpan di penampung**, bukan dibuang
- [ ] **AC 58** — dokumen yang gagal diurai **dilaporkan beserta sebabnya**
- [ ] **AC 59** — penampung medan tak dikenal **wajib kosong** sebelum pekerjaan dinyatakan selesai

## ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum diputuskan apakah dokumen lama dipindahkan seluruhnya, sebagian, atau tetap
dibaca lewat jalur lama.** Lingkup pemuat bergantung langsung padanya — memuat seluruh riwayat dan
memuat dua tahun terakhir adalah pekerjaan yang berbeda besarnya.

⚠️ **AC 59 tidak dapat dipenuhi sebelum tiket 19 lepas** — penampung medan tak dikenal tidak akan
kosong selama daftar kolom lengkap belum ada.