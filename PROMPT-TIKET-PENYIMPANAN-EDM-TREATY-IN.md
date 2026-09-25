# PROMPT — tiket penyimpanan EDM Treaty In

> Ketikkan `/to-tickets` sebagai manusia, lalu tempel berkas ini sebagai briefnya.

---

## 0. LINGKUP

Satu spec: **`.scratch\edm-treaty-in\spec-penyimpanan-relasional.md`** — 58 AC · 50 ID.

Keluaran: **`.scratch\edm-treaty-in\issues\`** — penomoran mulai dari **01**.

`D:\XML\RNM_BRD\` READ-ONLY. `D:\XML\nusantara-re\` terlarang.
⛔ Spec **jangan disunting**. ⛔ Tiket NB di `.scratch\nb-treaty-in\issues\` **jangan disentuh**.

### ⚠️ Prasyarat urutan kerja

**EDM tidak membuat satu pun tabel dasar.** Sepuluh tabel dasar dibuat oleh tiket NB. Setiap tiket
EDM yang menyentuhnya **wajib menyebut ketergantungan itu** — dikerjakan sesudah tiket NB yang
bersangkutan, bukan bersamaan.

---

## 1. ATURAN TIKET

⛔ **JANGAN menulis nama tabel di judul.** Isi tiket boleh; judul tidak.

Satu tiket = satu potongan kerja yang bisa diselesaikan dan diuji sendiri. Tiap tiket membawa:

| | |
| --- | --- |
| Judul | tanpa nama tabel |
| Rujukan | nomor AC dan `ID-n` |
| Bergantung pada | tiket NB mana, bila ada |
| Batas | apa yang **tidak** termasuk |
| Uji | lewat seam `repository` |
| Label | `ready-for-agent` atau `blocked` |

---

## 2. PEMBAGIAN YANG DISARANKAN

| # | Potongan | Sumber |
| ---: | --- | --- |
| 1 | Dua puluh satu kolom khas endorsemen pada tabel yang sudah ada | ID-15..ID-19 |
| 2 | Rantai generasi dan larangan percabangan | ID-7..ID-11 |
| 3 | Aturan keutuhan nomor urut antar generasi | ID-12..ID-14 |
| 4 | Nomor urut di endorsemen — tanpa penghapusan, tambah di belakang | ID-20..ID-23 |
| 5 | Tabel proyeksi selisih dan kolom asal barisnya | ID-24..ID-27 |
| 6 | Rumus selisih di lapisan `services` | ID-28..ID-33 |
| 7 | Dua penanda migrasi | ID-34..ID-39 |
| 8 | Pembatalan sebagai generasi bernilai nol | ID-23, P56 |
| 9 | Perbedaan perhitungan sebaran dari polis baru | ID-40, P60 |
| 10 | Pemuat migrasi endorsemen | ID-3, ID-34..ID-39 |

Ubah bila spec menuntut lain — **sebutkan alasannya**.

---

## 3. YANG WAJIB DILABELI `blocked`

| Tiket menyentuh | Tertahan butir |
| --- | --- |
| presisi fisik kolom uang | ⛔ 12 digit di depan koma **belum diuji** terhadap nilai terbesar `[data DBA]` |
| daftar kolom lengkap | ⛔ data guide **basi** — satu dokumen memuat 95 jalur yang tidak ada di dalamnya `[data DBA]` |
| anak tabel proyeksi | ⛔ `_SPREADING` dan `_INSTALMENT` dibutuhkan pembaca SQL atau tidak `[work owner]` |
| pembatalan | ⛔ sesudah dibatalkan, polis masih boleh di-endorse lagi atau tidak `[work owner]` |
| rantai endorsemen | ⛔ batas berapa kali satu polis boleh di-endorse — batas teknis 99 `[work owner]` |
| pemuat migrasi | ⛔ dokumen lama dipindahkan seluruhnya atau sebagian `[work owner]` |
| tabel sebaran tambahan | ⛔ beda dagangnya dari sebaran risiko belum dijelaskan `[work owner]` |

---

## 4. TIGA HAL YANG PALING MUDAH SALAH DITERJEMAHKAN

**1 · Tabel proyeksi bukan tabel sumber.** Tiga aturannya wajib masuk tiket, bukan diringkas:
① hanya ditulis Go, dalam transaksi yang sama dengan generasinya · ② boleh dihapus total dan
dibangun ulang, **tetapi hanya baris ber-`SUMBER='GO'`** · ③ bila isinya beda dari hasil hitung
ulang, **tabelnya yang salah**, bukan operannya.

⛔ Tanpa aturan ②, satu perintah bangun ulang menimpa seluruh angka historis dan **tidak bisa
dikembalikan**.

**2 · Rumus selisih ada di `services`, bukan di Oracle.** `V_POLIS_DIFFERENCE` **dibatalkan** —
view berarti rumus ditulis dua kali dan pasti bercabang. Tiket tidak boleh menghidupkannya.

**3 · Pembatalan bukan penghapusan.** Ia generasi baru berisi nol. Tiket yang membuat jalur
`DELETE` untuk pembatalan **salah**.

---

## 5. DISIPLIN

⛔ Nol `CREATE TABLE` di tiket. ⛔ Nol nama orang · nol nomor polis harfiah · nol cuplikan data
produksi. ⛔ Jangan menutup butir `[terbuka]`.

⚠️ Uji **pulang-pergi saja tidak cukup** — sebagian test memeriksa **nilai kolom langsung**.

⚠️ Khas endorsemen, wajib punya test tersendiri: rantai tiga generasi · penolakan percabangan ·
penolakan generasi yang kehilangan satu nomor urut · pembatalan menghasilkan generasi nol ·
bangun ulang proyeksi **tidak menyentuh** baris hasil migrasi.

---

## 6. TANDA BERHASIL

- Setiap AC tercakup sekurangnya satu tiket, pemetaannya ditulis.
- Nol judul memuat nama tabel.
- Setiap tiket yang menyentuh tabel dasar menyebut **ketergantungan pada tiket NB**.
- Tiket tertahan berlabel `blocked` dengan butirnya disebut.
- Peta AC → tiket disimpan sebagai `issues\00-PETA-AC.md`.
- Bab telemetri memisahkan yang terukur dari yang ditaksir.

---

*Disusun 23 September 2026.*
