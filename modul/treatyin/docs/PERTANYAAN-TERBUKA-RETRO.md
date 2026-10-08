# ✅ SUDAH DIJAWAB — apakah tab Retro layak dibangun sekarang

> **Dijawab pemilik proses 4 Oktober 2026: DITUNDA, dan kedua kontraknya DITANDAI.**
> Keputusannya, penandaannya yang tiga lapis, kewajiban mengadili `ShareSumary` bila kelak
> dibangun, dan syarat pembalikannya di
> [`KEPUTUSAN-PENYELARASAN-REPO.md` §14](KEPUTUSAN-PENYELARASAN-REPO.md).
>
> Kedua kontraknya: **`1000493`** dan **`1000755`**, keduanya `NonProportional`, 2 elemen
> masing-masing. Ditandai di `repository/warisan_retro.go`, diukur ulang dari Oracle oleh
> `TestKontrakRetroTertundaMasihDuaItu`, dan dicetak pada tiap `-cocokkan`.
>
> ⛔ **Berkas ini TIDAK dihapus.** Ukuran di bawah yang membuat jawabannya dapat diperiksa
> ulang — dan ia pula yang menjadi dasar pembalikan bila kontrak ketiga muncul.

---

## Pertanyaan asli — apakah tab Retro layak dibangun sekarang

**Diajukan kepada pemilik proses 3 Oktober 2026.**

⛔ **Tidak diputuskan di sini, dan tidak dibangun "supaya lengkap".** Yang dibawa ukurannya.

---

## Ukuran

Diurai utuh atas **seluruh 1.854 dokumen** `POOLDATA.M_TREATY_IN.JSONDATA`:

| | |
| --- | ---: |
| kontrak dengan `RetroList` berisi | **2** |
| total elemen | **4** |
| terbanyak per kontrak | 2 |

Dua kontrak. Dari 1.854.

### Bentuk satu elemen — 6 skalar, 12 larik bersarang

| Medan skalar | ada | maks aksara | contoh |
| --- | ---: | ---: | --- |
| `RetroType` | 4 | 16 | `2020 XOL ODYSSEY` |
| `RetroTypeID` | 4 | 5 | `10218` |
| `RetroPct` | 4 | 2 | `10` |
| `RetroTotalPct` | 4 | 3 | `100` |
| `RetroBrokerage` | **3** | 4 | `1.15` |
| `pxObjClass` | 4 | 36 | `ASM-FW-GISFW-Data-TreatyInRetroShare` |

| Larik bersarang | ada |
| --- | ---: |
| `RetroMemberList` | 4 |
| `RetroSpreadingList` | 4 |
| `Share` | 4 |
| `ShareSumary` *(ejaan ekspor)* | 4 |
| `TotalDeduction` · `TotalNet` · `TotalNetOR` · `TotalNetRI` · `TotalPremium` · `TotalShare` · `TotalShareOR` · `TotalShareRI` | 4 masing-masing |

---

## Kenapa ini belum dapat dirancang

**Delapan dari dua belas larik bersarang adalah `Total*` — nilai TURUNAN.** `INV-58` melarang
menyimpan larik turunan: angka yang disimpan dan angka yang dihitung akan berselisih, dan yang
menang selalu yang salah. Jadi yang tersisa untuk dimodelkan empat: `RetroMemberList`,
`RetroSpreadingList`, `Share`, `ShareSumary` — **tiga tingkat bersarang**, dan sebuah tabel induk.

**Empat elemen tidak cukup untuk merancang empat tabel.** Bentuk yang diturunkan dari 4 baris
akan mengunci rancangan pada dua kontrak yang kebetulan ada, dan tiap kolom yang ternyata
opsional baru terlihat ketika kontrak ketiga masuk — sesudah data pertama dimuat, ketika
menambah kolom jauh lebih mahal.

⚠️ Pembandingnya: tab yang sudah dibangun punya 60 sampai 11.365 baris. `AccumulationList`, yang
**paling jarang** di antara kesembilan tabel pendaratan, punya **60 elemen di 18 kontrak** — lima
belas kali lebih banyak bukti daripada Retro.

---

## Yang ditanyakan

1. **Apakah Retro dibangun sekarang**, atas dasar dua kontrak — atau ditunda sampai sistem lama
   memberi lebih banyak contoh?
2. Bila ditunda: apakah kedua kontrak itu (`RetroList` berisi) perlu ditandai supaya tidak ikut
   dianggap "selesai dipindahkan" pada rekonsiliasi tiket `44`?
3. Bila dibangun: ketiga tingkat bersarang itu **tiga tabel**, dan rancangannya menuntut jawaban
   lebih dulu atas apakah `ShareSumary` turunan dari `Share` — bila ya, `INV-58` melarangnya juga,
   dan yang tersisa dua tabel.

**Keadaan hari ini:** tab Retro tampil `.trin__belum` (*"belum ada KODE"*), bukan grid kosong —
yang benar, sebab yang kurang memang layarnya, bukan datanya.
