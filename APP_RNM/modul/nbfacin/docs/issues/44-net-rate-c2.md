# 44: Tab Coverage FIRE — tahap C2 Net Rate

> ⚠️ **Disusun agent atas perintah sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Frontend (`TabCoverage.tsx`
> `bolehNetRate` / `IsianNetRate`, `api.ts` `hitungNetRate`) dibangun sesi 0f; berkas ini mencatat backend (sesi c3).

**What to build:** ‰ Total Net Rate item dibagi ke setiap coverage secara proporsional terhadap Rate, lalu premi dihitung
ulang — hanya bila item memuat kelima coverage net rate. Tanpa migrasi (kolom `TOTAL_NET_RATE` / `NET_RATE` dari 193).

**Status:** backend selesai 03-10-2026 sesi c3 (uji hijau). ⛔ Butuh migrasi **193** sudah dijalankan (tiket 43).

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Activity\CekNetRate_ACT.xml` (kelas `ASM-FW-GISFW-Data-Coverage`; dipanggil `CountPremi_ACT` langkah 5 bila
  `Param.CallAct=="rate"`): langkah 1 `Local.FlagNetRate = "false"`; langkah 2 menelusuri
  `.Property.PropertyItemList(Param.IdxPropertyItem).CoverageList`, syarat `.OLDID=="FLEXAS"` / `"4.1A CC"` / `"4.3"` /
  `"4.2 PRGBI"` / `"OTHERS"` menandai `Local.FLEXAS` / `RSMD` / `FLOOD` / `PRGBI` / `OTHERS`, kelimanya `"1"` →
  `FlagNetRate = "true"`; langkah 3 menyalin flag ke `PropertyItemList(…).FlagNetRate`; langkah 4 `FlagNetRate=="false"` →
  `PropertyItemList(…).TotalNetRate = 0`.
- `Activity\CalculateNetRate_ACT.xml` (kelas `ASM-FW-GISFW-Data-PropertyItem`; dipanggil dari
  `Section\PropertyItemListCoverage.xml`): langkah 1 `Local.TotalNetrate = @divide(Local.TotalNetrate,1,20)`; langkah 2
  (flag `"true"`, per coverage) `Local.TotalRate = Local.TotalRate+.Rate`; langkah 3 (flag `"true"`, per coverage)
  `.NetRate = @divide(.Rate,Local.TotalRate,20)*Local.TotalNetrate`, lalu langkah 3.2 `Call CountPremi_ACT` dengan
  `PremiStatus = "percent"` (DiscountStatus dan CallAct tidak diisi).
- `CountPremi_ACT` cabang NetRate (tiket 43): NetRate terisi dan ≠ 0 → Premium dari NetRate.

## Backend (sesi c3)

- `POST /api/nbfacin/kasus/{caseId}/hitung-net-rate` badan `{ tsi, totalNetRate, coverages, isAdjustable?, pctAdjustOther? }`
  → `{ coverages, totalNetRate, flagNetRate }` (`handlers/coverage.go`, `services/netrate.go`). Tanpa menyimpan, tanpa
  identitas. Case selalu dibaca (404 / 503 seperti hitung-coverage). Flag false → `totalNetRate "0"`, coverages dikembalikan
apa adanya, periode tidak diperiksa. Flag true → NetRate
  dibagi, premi dihitung ulang (mode percent). 400: desimal tak sah (`coverages[k].<medan>`, `tsi`, `totalNetRate`,
  `pctAdjustOther`), basis tidak 1–4, ΣRate = 0; 404; **409** Begin / End case kosong (hanya bila flag true); 503.
- `PUT …/objek`: per item ber-coverage — flag true dan `totalNetRate` terisi → dibagi seperti POST; flag true tanpa
  `totalNetRate` → hitung biasa (A158 / A162), Total Net Rate tetap kosong; flag false → `TOTAL_NET_RATE = 0`.
  `TOTAL_NET_RATE` dan `NET_RATE` disimpan (NUMBER(38,8)). FlagNetRate **dihitung, tidak disimpan** (A167).
- Uji: `services/netrate_test.go` (lima OLDID lengkap / kurang satu / huruf kecil / berspasi / ganda; pembagian eksak
  1-2-1 dan sepertiga; Total Net Rate kosong; ΣRate 0; diskon tidak dihitung ulang; PUT tiga jalan),
  `handlers/coverage_test.go` `TestHitungNetRateHandler`. Mutasi tertangkap: pencocokan OLDID dilonggarkan
  (huruf/spasi), NetRate dibulatkan sebelum menghitung premi.

## Keputusan agent — ✅ A165–A168 DISETUJUI work owner 03-10-2026 (butir 100: *"setuju keputusan Net Rate A165–A168"*)

| # | Keputusan | Dasar |
| --- | --- | --- |
| A165 | ΣRate coverage = 0 pada item ber-flag dengan Total Net Rate → **400** (`totalNetRate: Σ Rate coverage = 0, Total Net Rate tidak dapat dibagi`), POST dan PUT | Pega membagi nol (`@divide(.Rate, 0, 20)`), hasilnya tidak tertulis di korpus `[pertanyaan terbuka]`. Alternatif yang ditolak: NetRate dibiarkan diam-diam (pengguna tidak tahu Total Net Rate tidak diterapkan) |
| A166 | Saat PUT, pembagian hanya bila flag true DAN `totalNetRate` terisi (termasuk "0" → NetRate 0); kosong → NetRate coverage dibiarkan (ketikan manual tidak terhapus). Flag false → TotalNetRate 0 hanya pada item ber-coverage; item tanpa coverage → apa adanya. Menggantikan bagian TotalNetRate A159 | Pega menjalankan CalculateNetRate hanya saat Total Net Rate diubah di layar dan CekNetRate hanya dalam konteks coverage; PUT tidak membawa peristiwa itu |
| A167 | `FLAG_NET_RATE` (kolom rancangan `T_PROPERTYITEMLIST`) tidak ditambahkan; flag dihitung dari OLDID coverage tiap kali | Flag sepenuhnya turunan coverage item; tanpa migrasi baru. Bila laporan / modul lain membaca kolom itu → migrasi ALTER kemudian |
| A168 | Hitung ulang premi sesudah pembagian memakai mode percent **tanpa** menghitung ulang diskon (Discount / % Discount tetap) | `[terverifikasi]` CalculateNetRate langkah 3.2 tidak mengisi `DiscountStatus`; CountPremi langkah 51–52 hanya berjalan bila DiscountStatus "percent" / "amount" |

### Catatan review dua sumbu (03-10-2026)

- ✅ **NetRate basi — keputusan work owner 03-10-2026 (butir 101): *"pertahankan sesuai pega"*.** Flag false hanya
  membuat `TotalNetRate = 0` — `NET_RATE` coverage TIDAK dikosongkan (= CekNetRate langkah 4, `[terverifikasi]`). Bila
  sesudah pembagian satu dari lima coverage dihapus, NetRate lama tetap ada dan premi tetap dihitung dari NetRate itu
  (cabang NetRate CountPremi). Kode tidak diubah.
- A165 vs A166: POST dengan `totalNetRate` kosong dihitung sebagai 0 dan tetap dibagi (ΣRate 0 → 400), sedangkan PUT
  dengan Total Net Rate kosong melewati pembagian (A166) — perbedaan disengaja (POST = peristiwa ubah Total Net Rate di
  layar, PUT tidak).
- A166 menerapkan "flag false → 0" pada setiap PUT; Pega hanya menjalankan CekNetRate saat rate diubah
  (`CallAct=="rate"`) — pendekatan, `[dugaan]` setara pada data yang konsisten.

## Acceptance criteria

- [x] POST hitung-net-rate; flag lima OLDID harfiah; flag false → 0 tanpa mengubah coverage.
- [x] Pembagian proporsional eksak (20 desimal), premi dihitung ulang dari NetRate presisi penuh.
- [x] ΣRate 0 tidak membagi nol (A165).
- [x] PUT objek menerapkan aturan yang sama; TOTAL_NET_RATE / NET_RATE tersimpan.
