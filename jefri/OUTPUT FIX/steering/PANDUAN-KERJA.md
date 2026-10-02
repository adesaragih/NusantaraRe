# Panduan Kerja — Migrasi Facultative Inward

Aturan operasional untuk siapa pun yang melanjutkan pekerjaan ini. Kosakata ada di `GLOSARIUM.md`;
keputusan rancangan di `../adr/`; keputusan lingkup di `../00-KEPUTUSAN-WORK-OWNER.md`.

---

## 1. Urutan baca sebelum menyentuh apa pun

| Butuh | Baca |
| --- | --- |
| Keputusan yang sudah ditutup — **selalu cek dulu** | `../00-KEPUTUSAN-WORK-OWNER.md` (K-001…) |
| Keputusan rancangan + alasannya | `../adr/` |
| Gambaran siklus NB | `../00-RINGKASAN-NB.md` |
| Yang harus diambil dari Pega selagi hidup | `../_EKSTRAKSI-PEGA-SELAGI-HIDUP.md` |
| Temuan lintas-siklus yang tetap berlaku | `../_ARSIP-lintas-siklus/_BACA-INI.md` |

Keputusan work owner **mengalahkan** label `[pertanyaan terbuka]` di dokumen mana pun. Bila sebuah
dokumen masih menyebut sesuatu sebagai terbuka padahal sudah ada keputusan, dokumen itulah yang
basi — bukan keputusannya.

---

## 2. Empat jebakan membaca korpus Pega

Keempatnya **terjadi** selama discovery ini dan menghasilkan angka salah sebelum dikoreksi.

| # | Jebakan | Koreksi |
| --- | --- | --- |
| 1 | Membaca kondisi `When` hanya dari `<pyConditionString>` | Baca **kedua** tag — 170 berkas (28,3 %) menyembunyikan kondisi di `<pyConditionValue1String>` |
| 2 | `Get-FileHash` mentah untuk membandingkan antar folder | Buang baris `px*`/`pz*` — metadata ekspor berbeda di setiap berkas |
| 3 | Membandingkan tanpa mengurutkan, atau dengan `Sort-Object` bawaan | Urutkan `[StringComparer]::Ordinal` — urutan tag ekspor tidak deterministik, dan ada nama rule yang hanya beda kapitalisasi |
| 4 | Membaca hanya `pySteps` tingkat atas | Pakai parser rekursif — **63 %** logika ada di langkah bersarang, sampai kedalaman 9 |

Tambahan: `pySectionReference` yang berisi nama section pada FlowAction, **bukan** `pyStreamName`.

---

## 3. Nama bukan bukti — enam bukti dari korpus ini

`CLAUDE.md` §3 butir 3, terbukti berulang kali dengan bentuk berbeda:

1. `When/IsFacout` namanya *fac out*, isinya menguji **banding**.
2. `When/IsOfferFacIn` memuat **dua kondisi berbeda** — teks tampilan vs ekspresi tersimpan.
3. Dua rule bernama `Insert…` sesungguhnya **menghapus lalu menyisipkan ulang**; deskripsi langkah
   `RDB-Delete`-nya sendiri berbunyi *"insert to db"*.
4. Dua shape flow berlabel sama *"Err Konversi?"*; yang satu menggerbangi hal yang sama sekali lain.
5. `Section/FormulaTreatyCapacityDesc` bernama *Treaty*, berkelas Fac In, dipanggil dari layar Fac In.
6. Seluruh tulisan Oracle memakai `RDB-List` — **metode baca** — dan tersimpan di tag `pyBrowseSQL`.

**Konsekuensi kerja:** arah operasi, arti nilai, dan lingkup tidak pernah disimpulkan dari nama. Hanya
isi tag yang membuktikan.

---

## 4. Aturan menulis kode

Dari `CLAUDE.md` dan ADR:

- **Uang tidak pernah `float`.** `Money{Amount decimal, Currency}` dan `Ratio{Value decimal, Scale}`
  adalah tipe berbeda yang tidak dapat dijumlahkan. → ADR-0004
- **Skala rasio melekat pada nilai**, diisi resolver COB terpusat yang mengikuti **rumus**, bukan
  label layar. COB tak dikenal → `panic`. → ADR-0004
- **Presisi pembulatan literal di tempatnya**, dengan komentar `// Asal: <rule>, langkah <n>,
  presisi <p>`. Jangan disentralkan, jangan diseragamkan. → ADR-0005
- **Urutan operasi dipertahankan persis.** `round(a,4) * b`, bukan `round(a*b,4)`.
- **Mata uang boleh `Unknown`** — `panic` hanya pada aritmetika lintas mata uang dan tulis Oracle.
  → ADR-0006
- `handlers → services → repository`, searah, tanpa memotong lapisan.
- Setiap query, predikat, dan transisi menyebut rule Pega asalnya dalam komentar (§4.6).
- Endpoint dan secret dari konfigurasi, **tidak pernah** literal.

---

## 5. Kapan gagal keras

| Situasi | Sikap |
| --- | --- |
| Arti enumerasi belum dijawab bisnis | `panic` |
| Isi stored procedure / tabel lookup tidak diketahui | `panic` |
| Rule dirujuk tetapi berkasnya tidak ada **di folder mana pun** | `panic` bila tercapai — **cabangnya tidak dihapus** |
| COB tidak ada di peta resolver skala | `panic` |
| Mata uang `Unknown` dipakai lintas mata uang / ditulis ke Oracle | `panic` |
| Status dikenali, baris keputusannya belum terverifikasi | `panic` saat paralel run · `decline`+log saat produksi → ADR-0002 |
| Kondisi tidak dikenali sama sekali | `panic` di **kedua** fase → ADR-0002 |

⚠️ **Hilang dari ekspor ≠ usang.** Ekspor korpus diambil dari titik waktu berbeda antar folder
(95 rule berbeda versi, arah tidak konsisten). Ketiadaan sebuah rule **bukan bukti** ia tidak dipakai.
Cabang yang memanggilnya **ditangguhkan, bukan dihapus**.

---

## 6. Perbaikan dipisahkan dari migrasi

`CLAUDE.md` §1. Sistem lama memuat banyak hal yang tampak keliru — perbandingan uang sebagai string,
tautologi pembanding limit, `ProRateType != 3` yang menihilkan premi, label rate yang bertentangan
dengan rumusnya, guard berbasis identitas orang.

**Semuanya direproduksi apa adanya**, diberi komentar sebagai kandidat perbaikan, dan didaftarkan ke
kuesioner. Memperbaikinya adalah keputusan bisnis yang terpisah.

Alasannya bukan kehati-hatian abstrak: rekonsiliasi menuntut **nol selisih sampai digit terakhir**
(ADR-0001). Satu "perbaikan" diam-diam menghasilkan selisih yang tidak dapat dijelaskan, dan
menyembunyikan salah-port yang sesungguhnya.

---

## 7. Aturan penulisan artefak

- **Nama orang tidak pernah disalin** ke artefak mana pun. Guard berbasis identitas dicatat sebagai
  **jumlah dan mekanismenya saja**. Metadata Pega (`pxCreateOpName`, `pxMoveImportOperName`, …)
  memuat nama orang — jangan ikut tersalin.
- Setiap klaim perilaku menyebut **path berkas + nama rule**, dan diberi label `[terverifikasi]` /
  `[dugaan]` / `[pertanyaan terbuka]`.
- Setiap angka disertai **perintah audit** yang menghasilkannya.
- Keputusan yang dibalik **tidak dihapus** — ditandai dibatalkan dan menunjuk penggantinya.
