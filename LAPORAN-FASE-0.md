# Laporan sesi — Fase 0, scaffold

Disusun 25 September 2026. Brief: `PROMPT-IMPLEMENTASI-GO-REACT.md` §1–§2.
Tiket: **tidak ada** — Fase 0 adalah pekerjaan pendahuluan.

---

## 1. ⛔ Fase 0 BELUM selesai menurut definisinya sendiri

Brief §2 menetapkan empat tanda selesai. **Tidak satu pun dapat dijalankan**, sebab
rantai alatnya tidak ada di mesin ini.

| Alat | Keadaan |
| --- | --- |
| `git` | ✅ 2.55.0 |
| `go` · `node` · `npm` · `docker` · `make` | ⛔ **tidak terpasang** |

| Tanda selesai §2 | Keadaan |
| --- | --- |
| `make build` lulus | ⛔ tak dapat dijalankan; **dan akan gagal** — `go.sum` dan `frontend/package-lock.json` belum ada |
| `make test` lulus dengan nol test | ⛔ tak dapat dijalankan; akan gagal karena sebab yang sama |
| `make run-api` menjawab `GET /healthz` | ⛔ tak dapat dijalankan |
| `make run-web` menampilkan satu halaman kosong | ⛔ tak dapat dijalankan |
| Nol aturan dagang | ✅ terpenuhi |

⛔ **`go.sum` dan `package-lock.json` tidak dapat saya buat.** Keduanya lahir dari
`go mod download` dan `npm install`, dan keduanya menuntut alat yang tidak ada. Hash
`go.sum` tidak dapat dikarang: ia dihitung Go sendiri. GitHub terjangkau dari mesin ini,
jadi begitu Go dan npm terpasang keduanya jadi dalam satu perintah.

⭐ **Yang tersisa untuk menutup Fase 0** — satu sesi pendek sesudah alat terpasang:
`go mod tidy` · `cd frontend && npm install` · lalu empat tanda selesai di atas.

---

## 2. Yang jadi

Struktur §1 lengkap, **tanpa menyimpang**, dan seluruh folder yang sudah ada
(`.scratch\`, `discovery\`, `docs\`, `dastin\`, `jefri\`, `scripts\`) tidak disentuh.

| Sisi | Cacah |
| --- | --- |
| Go | 8 berkas, 728 baris berisi, 216 baris komentar (29%) |
| Web | 8 berkas, 135 baris berisi, 26 baris komentar (19%) |
| Berkas scaffold seluruhnya | 24 |

Aturan §4 yang sudah berbentuk kode, masing-masing dengan rujukan ADR berawalan seri:

| Aturan | Tempatnya |
| --- | --- |
| Presisi desimal 38 dinyatakan **di satu tempat** | `pkg/utils/decimal.go` |
| Satu pasang konversi teks↔desimal (ADR-U-0034) | `pkg/utils/decimal.go` |
| Dua bentuk tanggal (ADR-U-0022) | `pkg/utils/date.go` |
| `Money` dan `Ratio` tipe berbeda, tak dapat dijumlahkan (ADR-F-0004) | `internal/models/money.go` |
| Uang keluar-masuk JSON sebagai **teks**, angka JSON ditolak (ADR-U-0003, ADR-U-0016) | `internal/models/money.go` |
| Nama skema eksplisit di setiap query (ADR-U-0033) | `repository.Qualify` |
| Nol `COMMIT` di teks SQL (ADR-U-0029) | `repository.PeriksaSQL` |
| Transaksi dibuka-ditutup aplikasi (ADR-U-0029) | `services.DalamTransaksi` |
| Wewenang di `services`, nol nama orang (ADR-U-0030) | `services.WajibPeran` |
| Hanya env var, nol literal host (ADR-U-0004) | `internal/config` |

---

## 3. Commit

Repositori git dibuat di `OUTPUT_HASIL_RNM\`. Korpus induk `D:\XML\RNM_BRD\` tidak ikut.
**1.334 berkas terlacak.**

| SHA | Pesan |
| --- | --- |
| `7263765` | fase-0: baseline dokumen |
| `0232686` | fase-0: scaffold |
| `beae599` | fase-0: perbaikan hasil code-review |
| `9f2ce68` | fase-0: seragamkan akhiran baris repository.go menjadi LF |

⚠️ **Penyimpangan sadar dari §2.** Brief menulis *"Commit pertama = scaffold kosong, pesan
`fase-0: scaffold`"*. Saya memecahnya menjadi dua: baseline dokumen lebih dulu, baru
scaffold. Sebabnya §6 butir 8 — *"Nol commit yang memuat berkas di luar lingkup"* — dan
menaruh 1.310 berkas dokumen ke dalam commit bernama "scaffold" melanggarnya. Pemecahan
ini juga membuat baris `Status:` pada tiket dapat ditelusuri sejak titik nol.
**Bila work owner lebih memilih satu commit, keempatnya dapat disatukan.**

---

## 4. Code review — dua poros, enam cacat nyata diperbaiki

`/code-review` dijalankan dengan titik tetap `7263765`, dua sub-agen paralel.

⚠️ **`CLAUDE.md` §8 keliru dan perlu diralat.** Ia menyatakan skill Matt Pocock
ber-`disable-model-invocation: true` sehingga *"agent tidak dapat memanggilnya"*, dan
memerintahkan berhenti bila diminta. Kenyataannya `/code-review` **termuat normal** dan
tidak menolak. Tidak ada pesan penolakan untuk dilampirkan sebagai bukti.

### Diperbaiki

| | Cacat | Akibat bila dibiarkan |
| ---: | --- | --- |
| 1 | `Router` menerima `*services.Service` lalu **tidak memakainya**; pemeriksaan kesehatan disuntik sebagai `*repository.DB` | Grafik impor bersih, tetapi **saat jalan** handlers menyentuh Oracle tanpa lewat services — persis yang dilarang §1 |
| 2 | `Money.UnmarshalJSON` menolak medan `amount` yang **absen** dengan galat parser, dan menuduh `{"amount":null}` sebagai angka JSON | Kolom kosong diperlakukan sebagai galat, melanggar ADR-U-0027 |
| 3 | `Ratio` tidak punya `UnmarshalJSON` | Penjaga "bukan angka JSON" tidak berlaku bagi rasio, padahal rasio ikut jalur uang |
| 4 | `PeriksaSQL` menolak kata `ROLLBACK` | Menolak pemanggilan **yang sah**: ADR-U-0029 justru mengizinkan procedure ber-`ROLLBACK` dipanggil terakhir |
| 5 | `DecimalContext()` membagikan pointer ke satu `apd.Context` | Satu pemanggil dapat mengubah `Precision` untuk **seluruh proses** — presisi 38 tidak lagi "di satu tempat" |
| 6 | Kutipan bukti di `go.mod` menunjuk bab **`T-keputusan` yang tidak ada** | Melanggar `CLAUDE.md` §4 butir 1. Bab yang nyata: **T-38**, baris 654, tabel putusan baris 661–666 |

Juga dibuang sebagai kode mati: medan `isPegaProd` di `repository` (penandanya sudah ada di
`config`) dan antarmuka `Aggregate` tanpa pengimplementasi. `App.jsx` dikembalikan menjadi
halaman benar-benar kosong sesuai tanda selesai §2.

### Disengketakan — tidak diubah, dan alasannya

- **`Pelaku`/`WajibPeran` di `services` disebut di luar lingkup.** §1 justru menamai isi
  paket itu *"aturan dagang, transaksi, **wewenang**"*, dan yang ada hanya seam kosong.
- **Flag `-migrate` disebut tugas kelima `main`.** Target `make migrate` diwajibkan §2 dan
  harus memanggil sesuatu; menaruhnya di `cmd/api` menghindari folder di luar pohon §1.
- **`.env.example` di luar pohon §1.** ⚠️ Penyimpangan kecil yang saya pertahankan: ia
  mendokumentasikan env var **tanpa satu pun nilai sungguhan**. Silakan dicabut bila §1
  dibaca harfiah.
- **`Add`/`Sub` kembar.** Mengekstrak helper untuk dua operasi menambah lapisan tanpa
  mengurangi risiko.

---

## 5. Pemeriksaan disiplin — dan instrumennya diuji lebih dulu

Tanpa kompilator, yang dapat diperiksa adalah aturan yang berupa teks. Instrumennya
diuji atas **7 contoh yang seharusnya kena** dan **5 yang seharusnya bersih**; ia lulus
keduabelasnya sebelum hasilnya dipakai.

⛔ Dua **tuduhan palsu** yang instrumennya sendiri hasilkan, dan sudah diperbaiki:
`float64` di dalam komentar yang justru melarangnya, dan `ORACLE_PASSWORD="$$ORACLE_DEV_PASSWORD"`
di `Makefile` yang dibaca sebagai sandi tertanam padahal ia rujukan env.

| Pemeriksaan | Hasil |
| --- | --- |
| Arah ketergantungan `handlers → services → repository` | ✅ 0 pelanggaran |
| `pkg/` mengimpor `internal/` | ✅ 0 |
| `float32`/`float64` di jalur uang Go | ✅ 0 |
| `Number(`/`parseFloat(` di jalur uang React | ✅ 0 |
| Rujukan ADR **tanpa** awalan seri | ✅ 0 (36 rujukan, seluruhnya berawalan) |
| Alamat host / IP harfiah | ✅ 0 |
| `COMMIT` di teks SQL | ✅ 0 |
| Rahasia atau sandi tertanam | ✅ 0 |
| Nama orang di kode | ✅ 0 |

⚠️ **Yang TIDAK diperiksa alat ini**: apakah kodenya benar-benar kompilasi. Diagnostik Go
di IDE hanya melaporkan `undefined: apd` — dan itu karena **modulnya belum pernah
diunduh** (nol cache modul Go), bukan cacat kode. Permukaan API `apd/v3` yang dipakai
(`BaseContext.WithPrecision`, `NewFromString`, `Context.Add`/`Sub`, `Text('f')`) **belum
diverifikasi kompilator**.

---

## 6. Butir terbuka — tidak saya tutup sendiri

1. ⛔ **Rantai alat.** Go, Node, npm, Docker, `make` perlu dipasang sebelum Fase 0 dapat
   ditutup. → pemilik pekerjaan / IT.
2. ⚠️ `[terbuka]` **Skala kolom uang berselisih antar-dokumen.** Brief §4 menulis
   `NUMBER(38,8)` (ADR-U-0003 · ADR-U-0016); korpus menulis **`NUMBER(38,20)` 171 kali**
   lawan `NUMBER(38,8)` 20 kali dan `NUMBER(38,2)` 7 kali, dan `TICKETS.md` T-38 memutuskan
   berdasarkan `NUMBER(38,20)`. **Presisinya sama (38), skalanya tidak.** Scaffold hanya
   menyatakan presisi, jadi kodenya tidak terpengaruh — tetapi migrasi pertama akan
   terpengaruh. → DBA / pemilik pekerjaan.
3. ⚠️ `[usulan]` **Pustaka yang belum pernah diputuskan lewat ADR**: `sijms/go-ora/v2`,
   `zustand`, kontainer `gvenzl/oracle-free`, router `net/http` bawaan. Seluruhnya
   ditandai `[usulan]` di `go.mod` dan di komentar.
4. ⚠️ **Gerbang Fase 2** belum dibuka: `KEPUTUSAN-TABEL-AKAR-BERSAMA.md` belum ada, dan
   memang bukan saya yang menulisnya. Fase 1 tidak menunggu gerbang ini.
5. ⚠️ **`CLAUDE.md` §8 perlu diralat** — lihat §4 laporan ini.

---

## 7. TELEMETRI EKSEKUSI

### Terukur

| Besaran | Nilai | Cara ukur |
| --- | ---: | --- |
| Commit | 4 | `git log` |
| Berkas terlacak | 1.334 | `git ls-files` |
| Berkas scaffold | 24 | `git diff --cached --name-only` |
| Baris Go berisi / komentar | 728 / 216 | pembacaan baris demi baris |
| Baris web berisi / komentar | 135 / 26 | idem |
| Diff scaffold | +1.009 | `git diff 7263765...HEAD --stat` |
| Diff perbaikan review | +122 / −79, 8 berkas | `git diff 0232686...HEAD --stat` |
| Rujukan ADR berawalan seri | 36 | pemeriksa disiplin |
| Swauji instrumen | 7 kena + 5 bersih, lulus 12/12 | `periksa_disiplin.py` |
| `NUMBER(38,20)` : `(38,8)` : `(38,2)` | 171 : 20 : 7 | `grep -o` lintas empat pohon dokumen |
| Sub-agen standar | 95.629 token, 17 panggilan alat, 247 detik | laporan harness |
| Sub-agen spec | 76.440 token, 9 panggilan alat, 165 detik | laporan harness |

### ⛔ Tidak diukur

| Besaran | Sebab |
| --- | --- |
| Token sesi utama | Angka token sejati tidak terlihat dari dalam sesi |
| Biaya | Turunan token; tidak diukur |
| Lama sesi, jam dinding | Tidak dicatat |
| Panggilan alat sesi utama | Tidak dicatat |
| Apakah kode kompilasi | ⛔ Tidak ada kompilator di mesin ini |

**Pengukuran luar (`claude --print --output-format json`) tidak dijalankan**, sebab ia
memulai sesi terpisah yang ongkosnya bukan ongkos ronde ini.
