# ✅ SUDAH DIJAWAB — rumah bagi tab berisi TEKS PANJANG

> **Dijawab pemilik proses 4 Oktober 2026: JALAN B — biarkan di `JSONDATA`.**
> Keputusannya, aturan pemilihan ejaan menurut cabang, dan syarat pembalikannya di
> [`KEPUTUSAN-PENYELARASAN-REPO.md` §15](KEPUTUSAN-PENYELARASAN-REPO.md).
>
> Tab **Exclusions** dan **Special Conditions** dibangun ronde yang sama.
> **`ValueDifference` tetap dikeluarkan** — ia objek, bukan teks; tabnya tetap `.trin__belum`.
>
> ⛔ **Berkas ini TIDAK dihapus.** Ukuran di bawah — terutama "303 dokumen, nol yang isinya
> identik" — adalah dasar aturan "jangan jatuh ke ejaan lain", dan aturan itu tidak dapat
> diperiksa ulang tanpanya.

---

## Pertanyaan asli — rumah bagi tab berisi TEKS PANJANG

**Diajukan kepada pemilik proses 3 Oktober 2026.**

⛔ **Berkas ini TIDAK memuat keputusan, dan itu disengaja.** Yang ditanyakan menyentuh bentuk
layar dan kewajiban penyuntingan, dan keduanya bukan milik saya. Yang saya bawa **ukuran**, supaya
pilihannya dapat diambil dengan angka di depan mata.

---

## Pertanyaannya

`ValueDifference`, `Exclusions`, dan `SpecialConditions` adalah tab yang isinya **bukan larik**.
Kedelapan tabel pendaratan (migrasi `430`) **tidak** menampungnya. Dua jalan ke depan:

| | Jalan A — tabel kunci-nilai `M_TREATYIN_TEKS` | Jalan B — biarkan di dalam `JSONDATA` |
| --- | --- | --- |
| Bentuk | `MASTERID`, `KUNCI`, `ISI CLOB` | tidak ada tabel baru |
| Baca | `SELECT ... WHERE MASTERID=? AND KUNCI=?` | urai CLOB tiap kali tab dibuka |
| Tulis kelak | satu baris per medan | menulis ulang seluruh dokumen |
| Ongkos sekarang | satu migrasi, satu pemuat, satu rekonsiliasi | nol |

---

## Ukuran yang menentukan jawabannya

Diukur atas **seluruh 1.854 dokumen**, diurai utuh sebagai JSON.

### 1 · Panjangnya melampaui `VARCHAR2`, jadi Jalan A berarti `CLOB`

| Kunci | Dokumen yang punya | Panjang maksimum |
| --- | ---: | ---: |
| `Exclusions` | 896 | **23.453 aksara** |
| `ExclusionsP` | 1.155 | **20.763** |
| `SpecialConditions` | 686 | **21.805** |
| `SpecialConditionsP` | 1.022 | **21.805** |
| `SpecialConditionsp` | 292 | 4.154 |

⛔ Batas `VARCHAR2` Oracle **4.000 bita**. Keempat yang pertama melampauinya berkali lipat.
Jalan A **wajib `CLOB`** — dan dengan itu ia kehilangan sebagian keunggulannya atas Jalan B:
keduanya sama-sama menyimpan teks besar yang tidak dapat dicari dengan `=`.

### 2 · ⛔ Tiga ejaan `SpecialConditions` membawa isi yang BERBEDA

| | Cacah |
| --- | ---: |
| dokumen dengan lebih dari satu ejaan | **303** |
| …isinya **sama** | **0** |
| …isinya **berbeda** | **303** |
| dokumen dengan tepat satu ejaan | 1.155 |
| dokumen **tanpa satu pun** | 396 |

Kontrak `1000001`:

| Ejaan | 45 aksara pertama |
| --- | --- |
| `SpecialConditions` | `1. The reinsured to be the sole judge as to w` |
| `SpecialConditionsP` | `As attached oh schedule on page no. 9 until p` |
| `SpecialConditionsp` | `Section 1 - Fire\nContingent Business INterrup` |

⚠️ **Ini mengubah pertanyaannya.** Ketiganya bukan tiga ejaan satu medan — ketiganya **medan yang
berbeda**. Aturan "baca keduanya, mundur ke yang lain" yang berlaku bagi lima pasangan cabang
lainnya, bila diterapkan di sini, **menampilkan teks yang salah**, bukan salinan yang basi.
Jadi sebelum Jalan A atau B dipilih, satu hal harus diputuskan lebih dulu:

> **Apakah layar menampilkan ketiganya, atau salah satunya?**
> Bila salah satunya — yang mana, dan atas dasar apa?

Pertanyaan itu berdiri sendiri dan berlaku di kedua jalan.

### 3 · `ValueDifference` BUKAN medan teks, dan tidak termasuk pertanyaan ini

Ia **objek**, ada di 1.612 dokumen, berisi `EGNPI`, `Limits`, `Share`, `LimitSummaryList`, dan
sembilan medan `Total*` — **potret nilai sebelum perubahan**, bukan teks yang pemakai ketik.

⛔ Menaruhnya di tabel kunci-nilai teks salah dua kali: bentuknya bukan teks, dan artinya bukan
medan layar. **Saya keluarkan dari pertanyaan ini**; ia menuntut pertanyaannya sendiri, yaitu
apakah riwayat nilai ikut dipindahkan sama sekali.

---

## Yang saya sarankan dipertimbangkan, bukan yang saya putuskan

Bila layar tab itu **hanya menampilkan** teksnya hari ini — tanpa penyuntingan — **Jalan B** tidak
berutang apa pun: dokumennya sudah dibaca untuk medan lain, dan satu medan `CLOB` lagi dari
dokumen yang sama tidak menambah perjalanan ke basis data.

Bila layar itu **akan dapat disunting** sebelum tiket `44` selesai, **Jalan A** terbayar: menulis
satu baris jauh lebih murah dan jauh lebih aman daripada menulis ulang dokumen 158 KB yang 174
medan lain ikut menumpanginya.

**Keputusannya milik Anda.** Apa pun pilihannya, ia dicatat di
`KEPUTUSAN-PENYELARASAN-REPO.md` beserta syarat pembalikannya, seperti §12.
