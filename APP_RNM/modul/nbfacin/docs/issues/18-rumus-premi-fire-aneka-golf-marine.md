# 18: Rumus premi FIRE, ANEKA (+Bonding), GOLF, MARINE CARGO

> ⚠️ **Disusun agent dari XML atas perintah work owner — bukan hasil `/to-tickets`.** Perintah:
> `PROMPT-LANJUT-IMPLEMENT-TANPA-TANYA-NB-RNW.md` (root repo, 01-10-2026) aturan 5 dan urutan 4.

**What to build:** `premium.Calculate` menghitung premi satu coverage untuk keempat lini yang belum punya
rumus, persis seperti sistem lama — termasuk presisi pembulatan per langkah — sehingga kelima kasus nyata P-5 di
kerangka rekonsiliasi (tiket 16) berubah dari "belum tercakup" menjadi "cocok" sampai digit terakhir.

**Asal** `[terverifikasi]` — semua `D:\migrasi\RNM\NB FacIn\Activity\`. Pemilih rumus `CountGrossPremi_Act.xml`
(ASM-FW-GISFW-WORK!COUNTGROSSPREMI_ACT) langkah 5.1 `IsFire`, 5.2 `IsAneka`, 5.3 `isGolfInsurance`, 5.4 →
`CountGPWMarinePAMbu_Act.xml` langkah 1.1 `IsMarineCargo`:

| Lini | Rule rumus | Langkah |
| --- | --- | --- |
| FIRE | `CountPremi_ACT.xml` (ASM-FW-GISFW-DATA-COVERAGE!COUNTPREMI_ACT) | 17–18 (LossLimit/SubLimit bawaan 100), 19–50 jalur `Param.PremiStatus == "percent"` (basis 1–4; 5 = Layering, tidak dipakai, butir 30), 51 diskon |
| ANEKA (termasuk Bonding, A11) | `CountPremiCoverageAneka.xml` (…!COUNTPREMICOVERAGEANEKA) | 11–12 (Losslimit), 14–16 (metode 1/2/3), 17 (MBD, indemnity) |
| GOLF | `FillPremiGolf.xml` (…!FILLPREMIGOLF) | 7–8, 9–11, 12 — rumusnya identik ANEKA |
| MARINE CARGO | `CountGPWMarinePAMbu_Act.xml` | 1.1.1.1.1 `round4(TSI × Rate / 100)`; 1.1.1.1.2 master policy (`PolicyType == 1 && IsMOP == "MOP"`) → 0 |

`[terverifikasi]` rumus diuji dulu dengan Python `decimal` atas fixture P-5 (01-10-2026), dengan pro-rata
tersimpan (`ProRatePercent`) dan `IndemnityPercentage` tersimpan: FIRE 54/54 coverage, ANEKA 1/1, MARINE CARGO
104/104 cocok. GOLF tanpa fixture.

**Di luar tiket ini** (masukan disediakan pemanggil, dicatat sebagai pertanyaan/tiket lanjutan):
- pro-rata dari aritmetika tanggal (`CountPremi_ACT` langkah 3–9, `CountPremiCoverageAneka` langkah 4/14.1:
  `@DateTime.DateTimeDifference`, zona Asia/Jakarta, `T050000.000 GMT`) — masukan `ProRatePercent`;
- `IndemnityPercentage` (lookup `SearchIndemnityRate_SQL`) dan `FirstScale` (`SearchFirstLossScale1_SQL`) —
  masukan apa adanya;
- `PctAdjustment` dari flag item properti (langkah 11–16) — masukan apa adanya;
- jalur `Param.PremiStatus == "amount"` (rate diturunkan dari premi) dan TSI/premi Nusantara Re, spreading, syariah.

**Blocked by:** 03

**Status:** ready-for-human — diport 01-10-2026, 159/159 coverage P-5 cocok; A37 dikonfirmasi untuk premi (butir 56); A36, A38, A39, A42, A43 ditahan work owner

- [x] Satu fungsi rumus per lini, presisi `@Math.divide` per langkah apa adanya; ~~seri setengah → `ErrPembulatanSeri`~~
  → **diganti A37**: data membedakan mode, setengah-ke-atas (lihat Comments)
- [x] FIRE: urutan langkah yang saling menimpa direproduksi (adjustable, NetRate); basis 5 → galat (Layering)
- [x] ANEKA dan GOLF berbagi satu rumus; cabang MBD hanya bila `IndemnityPercentage` terisi
- [x] MARINE CARGO master policy = 0
- [x] Pembagi komposit diturunkan dari satuan resolver K-018, bukan angka ajaib
- [x] Kerangka rekonsiliasi: kelima kasus P-5 cocok sampai digit terakhir per coverage; yang selisih dilaporkan

## Comments

### 2026-10-01 — diport (agent)

**Kode:** `backend/services/premium/lini_lain.go` (`calculateFire`, `calculateAnekaGolf`, `calculateMarine`,
`premiKomposit`, `kaliLossLimitPersen`, `lossLimitFire`, `lossLimitBawaan`); `premium.go` (medan `Input` baru, dispatch,
`AsalRumus`); `resolver.go` (`rasioPersen`); `rekonsiliasi/rekonsiliasi.go` (`barisPremi` per coverage,
`jalurCoverage`, `masukanCoverage`); kontrak `inti/backend/kontrak/facin.go` `MasukanPremiFacIn` + 8 medan (aditif) dan
adapter `kontrakfacin.masukan` (dijaga `TestMasukanSemuaMedan`).

**Patokan:** `TestBerkasKasusP5` — **159/159 coverage cocok sebagai teks persis**: edm-fire-1 FIRE 45, nb-fire-1 FIRE 5,
rnw-fire-1 FIRE 4, nb-kredit-1 ANEKA 1, nb-marinecargo-1 MARINE CARGO 104. Jumlah per kasus dihitung dua cara (tes + urai
Python; perintah dan jendela di `../../backend/services/premium/testdata/kasus/README.md`). FIRE juga dihitung ulang dengan
Python `decimal` independen dari Go: 54/54.

**Dua perubahan yang lahir dari data:**
- **A37 — pembulatan setengah-ke-atas** (mengganti butir 14). Dua coverage FIRE EDM jatuh **tepat seri** di desimal ke-21
  dengan digit sebelumnya genap, dan tersimpan dibulatkan ke atas; setengah-genap akan meleset. *(Ralat 01-10 malam: seri
  tepat ada **3**; yang ketiga, `nb-fire-1`, ber-digit ganjil dan tidak membedakan mode. Skrip pertama hanya menghitung seri
  yang membedakan mode — register bab Ralat.)* `ErrPembulatanSeri`
  dihapus di `premium` dan `acceptance` (`bulatNol` — berdasarkan analogi, `[dugaan]`).
- **A38 — skala** `.Premium*Local.Losslimit/100`: skala pembilang, bukan 38 digit (kasus Kredit).

**Cakupan data nyata terbatas** (temuan review spec): seluruh 54 FIRE basis 1, NetRate `"0"`, LostLimit 100; satu ANEKA
metode 3, Loading 0, bukan MBD; tanpa GOLF, tanpa master policy. Basis 2/3/4, NetRate, MBD, metode 1/2, LossLimit ≠ 100,
master policy diuji dengan masukan sintetis saja.

**Penjaga dari review:** NetRate/LossLimit FIRE `"0.0"` → `ErrTeksNolAmbigu` (A43); CoverageBasis kosong/tak dikenal →
`ErrCoverageBasisTakDikenal`; ANEKA metode 1 pro-rata 0/kosong → `ErrBentukBelumDiport` (langkah 14.1 dari tanggal).
Langkah 51 (diskon FIRE) hanya menulis `.Discount`, tidak mengubah `.Premium` — tidak diport. `Param.PremiStatus` tidak
dimodelkan: coverage `amount` akan muncul sebagai selisih (dicatat di komentar paket `rekonsiliasi`).

**Uji mutasi:** 26/28 lalu 6/6 penjaga baru tertangkap. Dua lolos = **ekuivalen terhadap data**: `IsMBD` dipaksa false
(coverage ANEKA nyata tanpa `IndemnityPercentage`, urai + grep) dan jalur GOLF ditukar (tanpa fixture GOLF). Dua lain
yang semula lolos (master policy, ejaan `LostLimit`) dikunci `TestMasukanCoverage`.

**Code review dua sumbu:** standar — komentar basi diperbaiki (`premium.go`, `resolver.go`, `rekonsiliasi.go`, doc
`TestPremiFireDitolak`), bukti `path + rule` ditambahkan, Loading kosong ANEKA/GOLF diberi label A42; *smell* (switch
lini berulang di `Calculate`/`AsalRumus`, `Input` makin lebar) dicatat, tidak direfaktor di tiket ini. Spec — penjaga di
atas; `bulatNol` dicatat sebagai analogi.

