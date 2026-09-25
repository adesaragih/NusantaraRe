# PROMPT — tiket penyimpanan Treaty In

> Ketikkan `/to-tickets` sebagai manusia, lalu tempel berkas ini sebagai briefnya.

---

## 0. LINGKUP

Dua spec penyimpanan, **keduanya sekaligus**:

```
.scratch\nb-treaty-in\spec-penyimpanan-relasional.md     67 AC · 36 ID
.scratch\edm-treaty-in\spec-penyimpanan-relasional.md    58 AC · 50 ID
```

Keluaran: **`.scratch\nb-treaty-in\issues\`** — lanjutkan penomoran dari tiket **16**.

`D:\XML\RNM_BRD\` READ-ONLY. `D:\XML\nusantara-re\` terlarang.
⛔ Kedua spec **jangan disunting**.

---

## 1. ATURAN TIKET

⛔ **JANGAN menulis nama tabel di judul tiket.** Aturan rumah dari ronde tiket sebelumnya.
Isi tiket boleh menyebutnya; judul tidak.

Satu tiket = satu potongan kerja yang bisa diselesaikan dan diuji sendiri. Tiap tiket membawa:

| | |
| --- | --- |
| Judul | tanpa nama tabel |
| Rujukan | nomor AC dan `ID-n` dari spec mana |
| Batas | apa yang **tidak** termasuk |
| Uji | lewat seam `repository` |
| Label | `ready-for-agent` atau `blocked` |

⛔ **Tiket yang bergantung pada butir terbuka diberi label `blocked`**, dengan butirnya disebut.

---

## 2. PEMBAGIAN YANG DISARANKAN

| # | Potongan | Sumber |
| ---: | --- | --- |
| 1 | Skema dasar + kunci + keunikan | NB ID-5..ID-10 · AC 1–7 |
| 2 | `NOURUT` dan penomoran ulang | NB ID-11..ID-13 · AC 8–11 |
| 3 | Tipe kolom dan konversi masuk | NB ID-14..ID-20 · AC 12–25 |
| 4 | Pemecah dokumen → baris | NB ID-21..ID-31 · AC 26–44 |
| 5 | Transaksi tunggal dan skema eksplisit | NB ID-32..ID-35 · AC 45–48 |
| 6 | Pulang-pergi dan dua bentuk dokumen | NB AC 49–54 |
| 7 | Pemuat migrasi | NB ID-3 · AC 55–59 |
| 8 | Rantai generasi `OLD_POLIS_ID` | EDM ID-7..ID-14 |
| 9 | Tabel proyeksi selisih + kolom `SUMBER` | EDM ID-20..ID-27 |
| 10 | Rumus selisih di `services` | EDM ID-28..ID-33 |
| 11 | Penanda migrasi `PASANGAN_BERGESER` · `RUMUS_BERLAPIS` | EDM ID-34..ID-39 |

Ubah pembagian ini bila spec menuntut lain — tetapi **sebutkan alasannya**.

---

## 3. YANG WAJIB DILABELI `blocked`

| Tiket menyentuh | Tertahan butir |
| --- | --- |
| presisi fisik kolom uang | ⛔ **12 digit di depan koma belum diuji** terhadap nilai terbesar `[data DBA]` |
| daftar kolom lengkap | ⛔ data guide **basi** — satu dokumen memuat 95 jalur yang tidak ada di dalamnya `[data DBA]` |
| `T_POLIS_BREAKDOWN_SPREAD` | ⛔ beda dagangnya dari `T_POLIS_SPREADING` belum dijelaskan `[work owner]` |
| pemuat migrasi | ⛔ dokumen lama dipindahkan seluruhnya atau sebagian `[work owner]` |
| anak tabel proyeksi | ⛔ `_SPREADING` dan `_INSTALMENT` dibutuhkan pembaca SQL atau tidak `[work owner]` |
| pembatalan dan batas endorse | ⛔ dua butir EDM `[work owner]` |

⇒ Perkiraan: **sebagian besar tiket `ready-for-agent`**, dan yang menyentuh DDL atau migrasi
`blocked`.

---

## 4. DISIPLIN

⛔ Nol `CREATE TABLE` di tiket — presisi fisik dicocokkan DBA **di dalam** tiket, bukan prasyaratnya.
⛔ Nol nama orang · nol nomor polis harfiah · nol cuplikan data produksi.
⛔ Jangan menutup butir `[terbuka]`.

⚠️ **Setiap cacah medan di spec adalah batas bawah, bukan total.** Tiket pemecah dokumen wajib
menyediakan penampung medan tak dikenal *(NB ID-27)*, dan penampung itu wajib kosong sebelum
pekerjaan dinyatakan selesai.

⚠️ Uji **pulang-pergi saja tidak cukup** — sebagian test memeriksa **nilai kolom langsung**.
Contohnya nyata: `GroupPanel "006"` disimpan sebagai `6` lalu dibaca balik jadi `"6"` akan lolos
pulang-pergi tetapi memecahkan penggolong `BusinessType_DeT`.

---

## 5. TANDA BERHASIL

- Setiap AC dari kedua spec tercakup **sekurangnya satu** tiket, dan pemetaannya ditulis.
- Nol judul tiket memuat nama tabel.
- Tiket yang tertahan berlabel `blocked` dengan butirnya disebut.
- Peta AC → tiket disimpan sebagai `issues\00-PETA-AC-PENYIMPANAN.md`.
- Bab telemetri memisahkan yang terukur dari yang ditaksir.

---

*Disusun 23 September 2026, sesudah kedua spec penyimpanan selesai dan koreksi 4quinque diterapkan.*
