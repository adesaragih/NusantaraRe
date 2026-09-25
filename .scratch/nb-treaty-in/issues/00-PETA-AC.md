# 00-PETA-AC — tabel silang acceptance criteria ke tiket

⛔ **Ini bukan tiket.** Berkas ini hanya memetakan **96 acceptance criteria** `spec.md` ke
tiket yang menutupnya, dan mencatat telemetri eksekusinya. ⛔ Nol keputusan baru.

## 1 · Cakupan — dicacah DUA CARA

| Yang dicacah | Cara A | Cara B | |
| --- | ---: | ---: | :---: |
| nomor AC tertutup | **96** | **96** | ✅ |
| AC **yatim** | **0** | **0** | ✅ |
| AC **ganda** *(>1 tiket)* | **0** | — | ✅ |
| gabungan = **1–96 tanpa lompat** | ✅ | | |

⭐ **Cara A** — kumpulkan nomor dari baris `**Menutup:**` seluruh tiket.
⭐ **Cara B** — baca bab `Acceptance criteria` tiap tiket dan kumpulkan nomor yang disebut di sana.
✅ **Kedua himpunan sama.**

## 2 · Enam belas tiket

| # | Judul | Status | AC |
| ---: | --- | --- | ---: |
| **00** | Skema penyimpanan dan migrasi | ⚠️ `needs-info` | 0 |
| **01** | Sumber data realisasi treaty | `ready-for-agent` | 9 |
| **02** | Daur hidup berkas realisasi dan nomor polis | `ready-for-agent` | 4 |
| **03** | Tangga tiga jenjang dan empat cabang putusan | `ready-for-agent` | 11 |
| **04** | Antrean bersama dan pemeriksaan keanggotaan | `ready-for-agent` | 3 |
| **05** | Peran menggantikan nama orang | ⚠️ `needs-info` | 5 |
| **06** | Penggolongan jenis usaha | `ready-for-agent` | 8 |
| **07** | Uang | `ready-for-agent` | 9 |
| **08** | Keutuhan penyimpanan | `ready-for-agent` | 5 |
| **09** | Tanggal | `ready-for-agent` | 5 |
| **10** | Jejak audit dan kronologi | `ready-for-agent` | 8 |
| **11** | Layar realisasi | `ready-for-agent` | 10 |
| **12** | Layar jenjang ketiga | ⛔ `blocked` | 5 |
| **13** | Rantai perhitungan uang | ⛔ `blocked` | 2 |
| **14** | Lingkup | `ready-for-agent` | 6 |
| **15** | Migrasi, paritas, dan jejak keputusan | `ready-for-agent` | 6 |

⭐ **`ready-for-agent` 12** · ⚠️ **`needs-info` 2** · ⛔ **`blocked` 2**

## 3 · Tabel silang — AC → tiket

| AC | Tiket | AC | Tiket | AC | Tiket | AC | Tiket |
| ---: | :---: | ---: | :---: | ---: | :---: | ---: | :---: |
| 1 | **03** | 25 | **07** | 49 | **11** | 73 | **02** |
| 2 | **03** | 26 | **07** | 50 | **11** | 74 | **02** |
| 3 | **03** | 27 | **07** | 51 | **11** | 75 | **06** |
| 4 | **03** | 28 | **07** | 52 | **12** | 76 | **06** |
| 5 | **03** | 29 | **08** | 53 | **11** | 77 | **11** |
| 6 | **03** | 30 | **08** | 54 | **11** | 78 | **12** |
| 7 | **03** | 31 | **02** | 55 | **11** | 79 | **13** |
| 8 | **03** | 32 | **09** | 56 | **11** | 80 | **12** |
| 9 | **03** | 33 | **09** | 57 | **01** | 81 | **05** |
| 10 | **03** | 34 | **09** | 58 | **01** | 82 | **05** |
| 11 | **04** | 35 | **09** | 59 | **02** | 83 | **08** |
| 12 | **05** | 36 | **01** | 60 | **08** | 84 | **03** |
| 13 | **05** | 37 | **01** | 61 | **14** | 85 | **07** |
| 14 | **04** | 38 | **01** | 62 | **14** | 86 | **07** |
| 15 | **01** | 39 | **10** | 63 | **14** | 87 | **13** |
| 16 | **01** | 40 | **10** | 64 | **14** | 88 | **14** |
| 17 | **01** | 41 | **10** | 65 | **14** | 89 | **01** |
| 18 | **07** | 42 | **10** | 66 | **06** | 90 | **08** |
| 19 | **06** | 43 | **10** | 67 | **06** | 91 | **05** |
| 20 | **06** | 44 | **10** | 68 | **15** | 92 | **04** |
| 21 | **06** | 45 | **11** | 69 | **09** | 93 | **15** |
| 22 | **06** | 46 | **11** | 70 | **15** | 94 | **15** |
| 23 | **07** | 47 | **12** | 71 | **10** | 95 | **15** |
| 24 | **07** | 48 | **12** | 72 | **10** | 96 | **15** |

## 4 · Urutan pengerjaan

```
01 ──▶ 02 ──▶ 03 ──▶ 04
  │       │      └──▶ 10 ──┐
  ├──▶ 06 │                ├──▶ 15
  ├──▶ 07 └──▶ 11 ──▶ ⛔12 │
  └──▶ 08 ──────────────────┘
         09 ──────────────────┘
```

⭐ **Dapat mulai segera:** **01** · **09** · **12** · **14** — keempatnya tanpa penahan.
⭐⭐ **Tertahan P18: NIHIL.**

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — baris ini semula berbunyi:
> > *"⭐ **Dapat mulai segera:** **01** · **09** · **14** — ketiganya tanpa penahan.
> > ⛔ **Tertahan P18:** **12** *(layar jenjang ketiga)* · **13** *(rantai perhitungan uang)*."*
>
> **P18 ditarik oleh tim migrasi** 2026-09-22 — lihat `VERIFIKASI-P18.md`. Tiket **12** lepas
> sepenuhnya. ⚠️ Tiket **13** **tetap `blocked`**, tetapi oleh **P30** *(aturan penjumlah yang
> hilang)* dan **P8** *(muatan efek keluar)* — **bukan** oleh P18.

⚠️ **Tertahan sebab lain:** **13** *(P30, P8)*.
⚠️ **Tertahan pemetaan peran:** **05**.
⚠️ **Tertahan bahan DBA:** **00** — ⭐ dan **tiket lain tidak bergantung padanya**.

## 5 · Pemeriksaan disiplin

| Yang diperiksa | Hasil |
| --- | ---: |
| ⛔ nama tabel / kolom fisik **di luar tiket 00** | ✅ **0** |
| ⛔ nilai berupa **nama orang** | ✅ **0** |
| ⛔ **nomor polis** apa adanya | ✅ **0** |
| ⛔ **DDL** / `CREATE TABLE` sebagai perintah | ✅ **0** |
| ⛔ **keputusan baru** | ✅ **0** — tiket hanya menerjemahkan `spec.md` |
| ⛔ **ADR baru** | ✅ **0** — hanya merujuk 15 yang ada |
| ⛔ butir terbuka **ditutup** | ✅ **0** |

⚠️ **Satu kemunculan `CREATE TABLE` ada di tiket 00**, di dalam kalimat yang justru
menyatakannya **nol** — ⭐ disebut, bukan dipakai.

## 6 · TELEMETRI EKSEKUSI

⛔⛔ **Pengukuran dari luar TIDAK dilakukan.** ⭐ Ronde ini berjalan **di dalam sesi
interaktif**, sehingga **biaya dan durasi tidak tersedia**. ⛔ Tidak ditaksir.

⚠️⚠️ **Dan angka di bawah diperoleh dengan MENGURANGKAN, bukan diukur langsung** —
baseline terakhir diambil sebelum ronde **verifikasi**, bukan sebelum ronde tiket ini.

| Yang dicatat | Nilai | Cara |
| --- | ---: | --- |
| ⭐ **token keluaran** | ⭐ **111.654** | ⚠️ selisih dikurangkan |
| ⭐ **token cache-read** | ⭐ **25.471.255** | sama |
| ⭐ **jumlah panggilan alat** | ⭐ **31** | sama |
| ⛔ durasi | ⛔ **TIDAK DAPAT DIUKUR** | butuh `--print` dari luar sesi |
| ⛔ biaya | ⛔ **TIDAK DAPAT DIUKUR** | butuh `--print` dari luar sesi |

⭐ **Cara pengurangannya, terbuka:** terukur sejak baseline verifikasi = keluaran **320.145** · cache-read **61.853.003** · panggilan **87**. ⭐ Dikurangi ronde **verifikasi** *(92.996 · 20.800.918 · 35)* dan ronde **spec** *(115.495 · 15.580.830 · 21)* yang keduanya sudah tercatat.

⚠️ **Pengurangan itu sah hanya bila tidak ada pekerjaan lain di antaranya** — dan memang
tidak ada, ⛔ tetapi **itu diingat, bukan diukur**. Dinyatakan apa adanya.

---

*Disusun 22 September 2026 dari `spec.md`, sesudah empat ronde grilling, satu ronde
verifikasi, dan 47 dari 49 pertanyaan terjawab.*
