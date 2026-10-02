# F06: Perhitungan premi retrosesi

**What to build:** Premi yang diteruskan ke reasuradur retro dihitung per coverage, lengkap dengan
komisi RI dan potongan diskon, sehingga nilainya siap ditulis ke produksi Fac Out.

## ⛔ DUA rule, bukan satu — dipilih menurut COB

**K-060 (amandemen):** `CountRateRetroCov` adalah **dua rule berbeda di kelas berbeda**, dipilih
menurut **lini bisnis**. **Keduanya berlaku dan keduanya harus diport** — bukan dipilih salah satu.

⛔ **Identitas rule = `pzInsKey` + `pyClassName`, BUKAN nama berkas.**

`[terverifikasi]`

| | Varian **FIRE** | Varian **ANEKA** |
| --- | --- | --- |
| `pyClassName` | `ASM-FW-GISFW-Data-PropertyItem` | `ASM-FW-GISFW-Data-Aneka` |
| Basis `pzInsKey` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-PROPERTYITEM COUNTRATERETROCOV` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-ANEKA COUNTRATERETROCOV` |
| **Berkas berlaku** | **`DDL\CountRateRetroCov.xml`** | **`DDL\CountRateRetroCov(ANEKA).xml`** |
| Salinan korpus | NB · RNW — ⛔ **BASI** | Endorsment — ✅ **mutakhir** (nol beda) |
| Langkah | **14** (10 puncak + 4 bersarang) | **12** (9 puncak + 3 bersarang) |

`[terverifikasi]` `pzOriginalInstanceKey` keduanya **sama** — varian Aneka lahir dari menyalin varian
PropertyItem lalu menyimpang.

### Perbedaan yang mengubah angka

| Aspek | **FIRE** | **ANEKA** |
| --- | --- | --- |
| **Pembagi premi** | **100.000** | ⛔ **10.000** |
| **Uji mata uang** | `.Currency=="IDR"` · `=="USD"` | ⛔ `.Currency.Name=="IDR"` · `=="USD"` |
| `Local.RIComIN` diakumulasi | ✅ ada | ⛔ **tidak ada** |
| `Local.TotalPremiCov` | ada | ⛔ tidak ada |
| Rumus koreksi langkah 10 | ✅ ada | ⛔ **tidak ada** |
| Gerbang `pySteps(2)`/`(3)`/`(4)` | `IsEDM` / `IsEdmExtendPeriod` / `IsEDM` | **identik** |

⛔ **Pembagi berbeda satu orde.** Memakai satu pembagi untuk kedua jalur menghasilkan premi retro
**10× salah** pada salah satunya.

📌 **Tabel langkah lengkap kedua varian ada di
`..\..\10-audit\01-pohon-langkah-countrateretrocov.md` §5 (FIRE) dan §7A (ANEKA); banding
berdampingan di §7B.**

## Rumus inti — varian FIRE

```
.PremiumRetro     = @Math.divide((ShareOffered * .Rate * Local.prorate), 100000, 20)
.PremiumRetro     = .PremiumRetro - @Math.divide((.PremiumRetro * .DiscountPercentage), 100, 20)
.RICommPercentage = @Math.divide(Local.RIComIN, Local.LengCov, 20)
.RIComm           = @Math.divide((.PremiumRetro * .RICommPercentage), 100, 20)
.Rate             = @Math.divide(Local.RateOut, Local.LengCov, 20)
```

`[terverifikasi]` **Kedua rata-rata itu sejajar.** `Local.RateIN` **dan** `Local.RIComIN` sama-sama
**diakumulasi per coverage** di `RH_2.pySteps(8)`, lalu sama-sama dibagi `Local.LengCov`:

```
RH_2.pySteps(8)  (Property-Set, iterasi .CoverageList, tanpa gerbang)
  [1] Local.RateIN  := Local.RateIN + .Rate
  [2] Local.LengCov := .pxListSubscript
  [3] Local.RIComIN := Local.RIComIN + .RIComm
```

⛔ **URUTAN OPERASI DIKUNCI: diskon dikurangkan LEBIH DULU, `.RIComm` dihitung SESUDAHNYA.**
`[terverifikasi]` Keduanya berada di langkah **berbeda** — diskon di `RH_2.pySteps(9).pySteps(2)`,
komisi di `RH_2.pySteps(9).pySteps(3)` — dan **keduanya tanpa gerbang**, sehingga tidak ada percabangan
yang dapat membalik urutannya. Menghitung `.RIComm` dari premi **sebelum** diskon menghasilkan angka
berbeda.

⚠️ **Gerbang langkah 3 berbeda antara kedua salinan.** `[terverifikasi]` Salinan **DDL yang berlaku**
menggerbangi `Local.prorate := 100` **hanya** dengan `IsEdmExtendPeriod`; salinan korpus yang basi
menuntut `IsEDM` **dan** `IsEdmExtendPeriod`. **Ikuti DDL.** Perincian di audit 01 §4.2.

✅ Pembagi **100.000** = `1.000 × 100` adalah konsekuensi **aturan pembagi komposit K-018**, bukan
rumus lain.

## `Local.prorate` — enam penugasan, tetapi **nilai akhirnya sudah pasti**

> ⛔ **Keputusan work owner 21 September 2026.** Pada `CountRateRetroCov`, `Local.prorate` **bernilai
> `1`** keluar dari langkah 1. **Dua penugasan sebelumnya di langkah yang sama TIDAK berlaku.**
>
> Rumusan lama di tiket ini menyajikan keenam penugasan sebagai sama-sama mungkin dan menggantungkan
> jawabannya pada `REPEATINGINDEX`. **Digantikan bagian ini**; tidak dihapus dari riwayat
> (`PANDUAN-KERJA` §7).

`[terverifikasi]` Keempat langkah prorata **identik di kedua varian** — gerbang, T/F dan penugasannya
sama persis di `DDL\CountRateRetroCov.xml` dan `DDL\CountRateRetroCov(ANEKA).xml`, diukur dengan
parser yang membedakan **anak langsung** dari keturunan:

| Langkah | Gerbang | T | F | Penugasan |
| --- | --- | ---: | ---: | --- |
| `pySteps(1)` | — | 2 | 2 | `[1]` `EndPeriod − StartPeriod` ⛔ **mati**<br>`[2]` `@divide(prorate, @if(EDMDay=="",365,@toDecimal(EDMDay)), 20) * 100` ⛔ **mati**<br>`[3]` **`1`** ✅ berlaku |
| `pySteps(2)` | `IsEDM` | 2 | 3 | `(ProrateEDMEnd + ProrateStartEDM) * 100` |
| `pySteps(3)` | **`IsEdmExtendPeriod` saja** — ⛔ salinan korpus yang basi menuntut `IsEDM` **dan** `IsEdmExtendPeriod` | 2 | 3 | `100` |
| `pySteps(4)` | `IsEDM` | **3** | **2** | `100` — ⛔ **terbalik**: jalan saat `IsEDM` **SALAH** |

### Nilai akhir `Local.prorate` per jalur

`[terverifikasi]` dari tabel di atas, dengan penambatan kode transisi (`2` = jalankan lalu lanjut,
`3` = lewati langkah — **`[dugaan kuat]`**, `..\..\10-audit\02-peta-kode-transisi.md` §2.2):

| Jalur | Langkah yang berjalan | **Nilai akhir `Local.prorate`** |
| --- | --- | --- |
| **NB / RNW** (non-EDM) | 1 → 4 | **`100`** |
| **EDM biasa** | 1 → 2 | **`(ProrateEDMEnd + ProrateStartEDM) × 100`** |
| **EDM extend period** | 1 → 2 → 3 | **`100`** |

⛔ **Pada jalur non-EDM, `Local.prorate` BUKAN prorata.** Nilainya selalu `100` dan berfungsi sebagai
**konstanta skala** yang menyatu dengan pembagi — **bukan** pembagian waktu. Kosakatanya dikunci di
`..\..\steering\GLOSARIUM.md` bagian *Kosakata yang dihindari*.

📌 `[dugaan]` `<pyStepsDescription>` langkah 4 berbunyi **`IsNB`** di **kedua** varian — sejalan dengan
pembacaan "jalan saat bukan EDM". **Deskripsi adalah label, bukan bukti** (`PANDUAN-KERJA` §3);
dicatat sebagai penguat saja, tidak dijadikan dasar.

### ⛔ Dua rumus mati di langkah 1 — **tetap diport**

`[terverifikasi]` Ketiga penugasan berada di `pyParamArray` yang sama, ber-atribut `REPEATINGINDEX`
**1**, **2**, **3**; yang ber-indeks **3** adalah `Local.prorate := 1`.

Kedua rumus pertama — selisih periode, lalu pembagian terhadap `EDMDay`/365 — **tidak berpengaruh
pada jalur mana pun**. Keduanya **diport apa adanya** (`CLAUDE.md` §1), ditandai **KODE MATI**, dan
didaftarkan sebagai **kandidat perbaikan (K-046)** bersama keanehan `Local.RIComIN` varian Aneka di
bawah — `K046_Prorate_DuaRumusMati_Langkah1`.

⚠️ **Yang dijawab work owner adalah NILAI `Local.prorate`, bukan aturan umum** bahwa urutan
`REPEATINGINDEX` = urutan eksekusi. Aturan umum itu **tetap `[pertanyaan terbuka]`** di tempat lain
yang memakainya, dan **tidak ikut ditutup** oleh keputusan ini.

### Pemeriksaan silang K-018 — kedua pembagi benar

`[terverifikasi]` Pada jalur non-EDM, dengan `Local.prorate = 100`:

| Varian | Rumus efektif | Skala `.Rate` | Sesuai K-018? |
| --- | --- | :-: | :-: |
| **FIRE** (`Data-PropertyItem`) | `Share × Rate × 100 / 100000` = `Share × Rate / 1000` | **‰** | ✅ |
| **ANEKA** (`Data-Aneka`) | `Share × Rate × 100 / 10000` = `Share × Rate / 100` | **%** | ✅ |

⛔ Keduanya konsisten dengan **skala rasio per lini bisnis** yang dikunci **K-018** — FIRE per mille,
ANEKA persen. Pembagi yang berbeda **bukan kekeliruan ekspor**; ia konsekuensi langsung **aturan
pembagi komposit**. Menyeragamkannya menghasilkan premi **10× salah** pada salah satu jalur.

## Penskalaan mata uang sebelum dipakai

`[terverifikasi]` `Local.ShareOffered` diskala ulang **sebelum** masuk rumus premi:

| Gerbang | Penskalaan |
| --- | --- |
| `.Currency=="IDR" && …PctPremiAllObj != ""` | `ShareOffered * PctPremiAllObj / 100` |
| `.Currency=="USD" && …PctPremiAllObjUSD != ""` | `ShareOffered * PctPremiAllObjUSD / 100` |

⛔ Hanya **dua** mata uang yang punya cabang. Mata uang lain melewati keduanya tanpa penskalaan —
**diport apa adanya**, bukan digeneralisasi.

## Rumus kedua adalah KOREKSI, bukan alternatif

`[terverifikasi]` Langkah `RH_1.pySteps(10)` bergerbang kesetaraan **4 desimal**:
`@Math.divide(Local.TotalPremiCov,1,4) == @Math.divide(.CoverageList(1).FacOutObjectList(1).ObjectPremi,1,4)`
dengan T=3 / F=2 → langkah bersarang `10.1` berjalan **ketika keduanya TIDAK sama**:

```
.PremiumRetro = @Math.divide((@toDecimal(Primary.CoverageList(1).FacOutObjectList(1).ShareOffered)
                              * .Rate * Local.prorate), 100000, 20)
```

⛔ Rumus koreksi ini memakai `ShareOffered` **langsung dari halaman Primary**, sehingga **MELEWATI
penskalaan mata uang** di atas. Itu perbedaan nyata, bukan penyederhanaan.

⚠️ Ada **dua gerbang kesetaraan pada presisi berbeda**: rate pada **10 desimal**
(`RH_1.pySteps(9).pySteps(1)`), premi pada **4 desimal** (`RH_1.pySteps(10)`). Presisinya **literal di
tempatnya** (ADR-0005), tidak diseragamkan.

## ⛔ Keanehan `Local.RIComIN` — **hidup di varian ANEKA saja**

> Sesi sebelumnya **mencabut keanehan ini seluruhnya**. ⛔ **Itu terlalu jauh, dan dibatasi ulang.**
> Pencabutan **benar untuk varian FIRE**, tetapi **salah untuk varian ANEKA**.

`[terverifikasi]` Pengukuran per varian:

| Salinan | Kelas | `<PropertiesName>` (menyetel) | `<PropertiesValue>` (membaca) | Keanehan |
| --- | --- | ---: | ---: | :-: |
| `DDL\CountRateRetroCov.xml` (FIRE, berlaku) | `Data-PropertyItem` | **1** | 2 | ✅ tidak ada |
| `NB FacIn\Activity\CountRateRetroCov.xml` (FIRE, basi) | `Data-PropertyItem` | 0 | 1 | *(artefak salinan basi)* |
| **`DDL\CountRateRetroCov(ANEKA).xml` (ANEKA, berlaku)** | `Data-Aneka` | **0** | **1** | ⛔ **ADA** |

⛔ **Di jalur ANEKA, `.RICommPercentage = @Math.divide(Local.RIComIN, Local.LengCov, 20)` membagi dari
variabel yang TIDAK PERNAH DIISI.** Langkah 8 varian Aneka hanya mengakumulasi `Local.RateIN` dan
`Local.LengCov`; **tidak ada** `Local.RIComIN := Local.RIComIN + .RIComm`.

**Diport apa adanya; kandidat perbaikan (K-046) — terbatas varian ANEKA.**

```powershell
# FIRE(DDL) setel=1 baca=2 ; ANEKA(DDL) setel=0 baca=1
foreach($p in @('D:\migrasi\RNM\DDL\CountRateRetroCov.xml',
                'D:\migrasi\RNM\DDL\CountRateRetroCov(ANEKA).xml')){
  $t=[IO.File]::ReadAllText($p); $s=0;$b=0
  foreach($m in [regex]::Matches($t,'<PropertiesName>([^<]*)</PropertiesName>')){ if(([regex]'RIComIN').IsMatch($m.Groups[1].Value)){$s++} }
  foreach($m in [regex]::Matches($t,'<PropertiesValue>([^<]*)</PropertiesValue>')){ if(([regex]'RIComIN').IsMatch($m.Groups[1].Value)){$b++} }
  "$([IO.Path]::GetFileName($p)) : setel=$s baca=$b" }
```

⚠️ `[pertanyaan terbuka]` Bahkan di varian FIRE, `Local.RIComIN` **tidak pernah di-nol-kan** —
`pySteps(5)` menginisialisasi `Local.RateIN := 0` dan `Local.TotalPremiCov := 0`, tetapi tidak ada
`Local.RIComIN := 0`. Apakah variabel lokal Pega otomatis nol adalah **`[di luar korpus]`**.

## ⚠️ Jalur properti mata uang berbeda antar varian

`[terverifikasi]` FIRE menguji **`.Currency`**; ANEKA menguji **`.Currency.Name`**.

⛔ **Satu implementasi yang membaca mata uang lewat satu jalur akan SALAH pada salah satu varian** —
gerbangnya tidak pernah benar, sehingga **penskalaan mata uang terlewat diam-diam** dan premi keluar
tanpa diskala. Kegagalannya **senyap**, bukan galat.

## `[pertanyaan terbuka]` Kemungkinan varian KETIGA — kelas `Data-Cargo`

`[terverifikasi]` Indeks `Embed-Reference-Rule` pada **7 berkas** mencatat `CountRateRetroCov` pada
**tiga** kelas: `ASM-FW-GISFW-Data-PropertyItem`, `ASM-FW-GISFW-Data-Aneka`, dan
**`ASM-FW-GISFW-Data-Cargo`**. Tiap pemanggil produksi memuat **tiga** langkah `Call
CountRateRetroCov`, seluruhnya bergerbang `.CoverageList(1).PremiumRetro==""`.

⛔ **Varian berkelas `Data-Cargo` TIDAK ADA** di ketiga folder korpus maupun di `DDL\`.

⚠️ Work owner menyatakan hanya ada **dua** rule di Pega. Indeks rujukan merekam resolusi **saat berkas
terakhir disimpan**, jadi ia **bukan bukti** rule itu masih ada sekarang. **Jangan disimpulkan salah
satu.** Rincian: `..\..\10-audit\09-pemanggil-countrateretrocov.md`.

**Asal (Pega).** **`DDL\CountRateRetroCov.xml`** (varian FIRE) · **`DDL\CountRateRetroCov(ANEKA).xml`**
(varian ANEKA) — keduanya berlaku (K-060). Pembanding basi:
`NB FacIn\Activity\CountRateRetroCov.xml`.

**Keputusan.** **K-060** · **K-058** · **K-048** §8.1 · K-018 · K-010/K-012 ·
K-027 · **K-046** · ADR-0004/0005 · `CLAUDE.md` §1, §4.1, §4.6

**Blocked by:** F05 · `..\01-tracer-money-ratio-premi-pa.md` · `..\03-resolver-cob-skala.md`
⛔ Kedua tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] ⛔ **Diskon sebelum komisi RI** — urutan terkunci dan teruji secara terpisah
- [ ] Pembagi **100.000** diturunkan **resolver K-018**, bukan angka hafalan
- [ ] Enam penugasan `prorate` **ditimpa berurutan**, tidak dipecah jadi rasio bernama berbeda (K-048 §8.1)
- [ ] ⛔ `Local.prorate` keluar langkah 1 bernilai **`1`** (keputusan work owner 21 Sept) — dua rumus pertama **KODE MATI**, **tetap diport**, `K046_Prorate_DuaRumusMati_Langkah1`
- [ ] Nilai akhir prorata per jalur **diuji terpisah**: non-EDM `100` · EDM `(ProrateEDMEnd+ProrateStartEDM)×100` · EDM extend `100`
- [ ] Pada jalur non-EDM, `100` diperlakukan sebagai **konstanta skala**, bukan prorata — tidak dinamai `prorate` di kode baru
- [ ] Gerbang langkah 4 **terbalik** (jalan saat bukan EDM) — arah dibaca dari kode transisi, bukan dari label
- [ ] Penskalaan mata uang hanya untuk **dua** mata uang; sisanya lewat tanpa penskalaan
- [ ] Rumus koreksi berjalan **saat total TIDAK sama**, dan **melewati** penskalaan mata uang
- [ ] Dua presisi pembanding (**10** dan **4** desimal) **literal di tempatnya**, tidak diseragamkan (ADR-0005)
- [ ] ⛔ **DUA varian diport, bukan satu** — dipilih menurut **COB** (K-060 amandemen)
- [ ] Pembagi **100.000** (FIRE) dan **10.000** (ANEKA) **tidak diseragamkan**
- [ ] Mata uang dibaca lewat **jalur yang benar per varian** — `.Currency` (FIRE) vs `.Currency.Name` (ANEKA)
- [ ] Varian ANEKA **tanpa** rumus koreksi dan **tanpa** `Local.TotalPremiCov`
- [ ] `Local.RIComIN` diakumulasi **hanya** di varian FIRE; **K-046** `K046_RIComIN_TakDiisi_HanyaAneka`
- [ ] Gerbang langkah 3 mengikuti **salinan DDL**: `IsEdmExtendPeriod` **saja**
- [ ] **K-046** `K046_Prorate_EnamPenugasanSalingMenimpa`
- [ ] Nilai uang bertipe `Money`; rasio bertipe `Ratio`; `Money + Ratio` **gagal saat kompilasi**
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6), **beserta kelasnya** — nama saja tidak cukup

**`[pertanyaan terbuka]` yang tersisa:**
1. ✅ **DIJAWAB WORK OWNER 21 September 2026 — ditutup untuk tiket ini.** ~~Apakah urutan
   `REPEATINGINDEX` dalam satu `Property-Set` = urutan eksekusi; bila ya, dua penugasan `prorate`
   pertama di langkah 1 mati.~~ Work owner menetapkan `Local.prorate` **bernilai `1`** keluar dari
   langkah 1, sehingga kedua rumus itu **memang mati** — lihat bagian *Dua rumus mati di langkah 1*.
   ⚠️ **Yang ditutup adalah NILAI `Local.prorate` pada rule ini**, **bukan** aturan umum bahwa
   `REPEATINGINDEX` = urutan eksekusi. Aturan umum itu tetap **`[di luar korpus]`** dan tetap terbuka
   di tiket/dokumen lain yang menyandarinya.
2. Mengapa gerbang `IsEDM` hilang di langkah 3 salinan DDL varian FIRE — disengaja atau kekeliruan
   ekspor. Korpus tidak dapat menjawab.
3. Apakah `Local.RIComIN` perlu di-nol-kan eksplisit di awal.
4. **Apakah ada varian ketiga berkelas `Data-Cargo`** — lihat bagian di atas.
5. **Apa yang memilih varian saat runtime.** Korpus tidak memuat gerbang pemilih; pemilihan terjadi
   lewat **kelas halaman** tempat rule dipanggil. Pemetaan COB → kelas belum terbaca dari korpus.
