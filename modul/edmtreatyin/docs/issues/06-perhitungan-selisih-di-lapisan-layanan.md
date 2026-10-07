# 06: Perhitungan selisih di lapisan layanan

**Status:** ready-for-agent
**Blocked by:** **02** · **03**
**Bergantung pada tiket NB:** **16** *(kerangka penyimpanan)*
**Menutup:** AC **16–22** *(7 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-4 · ID-28 · ID-29 · ID-30 · ID-31 · ID-32

> ## ⛔ KOREKSI 06-10-2026 — berbukti XML (log: `../KOREKSI-DOKUMEN-2026-10-06.md`)
>
> | Bunyi lama (dikutip) | Bunyi baru | Bukti |
> | --- | --- | --- |
> | *"**arah utang-piutang** — diturunkan dari **tanda** selisih"* (AC 20, untuk seluruh selisih) | ⛔ **Tidak berlaku untuk jalur proporsional.** `EDMTCalculateTreatyDifference` tidak menulis `DueTo` (juga tidak `Currency`/`CurrencyID`); nol aturan korpus menulis `TreatyDifference.DueTo`. Arah dari tanda selisih **hanya pada lapisan XOL**. ⇒ `T_POLIS_DIFFERENCE` tanpa `DUE_TO` (migrasi 360) benar | `Activity/EDMTCalculateTreatyDifference.xml` langkah 1–6 (medan tertulis: 26 `TreatyDifference.*`, tanpa `DueTo`); `Activity/CalculateDifferenceEDM_act.xml` langkah 1.2.1 / 1.2.3 / 1.2.5; pembaca satu-satunya `Activity/InsetTreatyInProdAddendum_Act.xml` 9.4 / 9.5 (`.DueTo==1` / `==0`) |
> | *"Pada selisih lapisan, arah utang-piutang diturunkan dari **jumlah sepanjang daftar lapisan**"* (AC 21) | ⛔ **Keliru.** `DueTo` **lapisan** = tanda selisih `DueToValue` **lapisan itu** (`@If(local.duetovalue>0,"DUE TO US","DUE TO YOU")`, sesudah batas bawah 0). `DueTo` **induk** (langkah 1.1) diset **sebelum** kalang lapisan, memakai `local.duetovalue` sisa **iterasi sebelumnya**; jumlah lapisan (langkah 2.2.1) hanya masuk `DueToValue` induk (2.3). **Nol pembaca** `DueTo` induk maupun lapisan: produksi NonProp menulis `"DUE TO US"` tetap | `CalculateDifferenceEDM_act.xml` langkah 1.1, 1.2.1, 2.1–2.3; `InsetTreatyInProdAddendum_Act.xml` langkah 10.1; `Section/DetailPolicyTreatyInAddPremi.xml`, `DetailPolicyAddPremiDetail.xml` (kolom grid tanpa `DueTo`) |
> | *"`[terbuka]` Kedua blok rumus … pemilihnya hanya tertulis di keterangan langkah"* | ✅ **Gugur oleh bukti korpus.** Pemilihnya **dijalankan**: langkah 1 PRE `.OldData.EDMNo==""` → salah ⇒ **lompat** ke label `HasEDMNo` (langkah 4); langkah 4 PRE sama → benar ⇒ **keluar** activity. Sejalan `../spec.md` bab 5.3 dan P57 | `EDMTCalculateTreatyDifference.xml` langkah 1 (`pyStepsPreCondParamsWhenFalse=1`, `…FalsePrms=HasEDMNo`), langkah 4 (`…WhenTrue=6`) |
> | *"sisi non-proporsional … memang sudah seragam"* (AC 22, ID-31) | ⚠️ Pengurangan selalu terhadap **nilai** generasi lampau — benar; tetapi **bukan pengurangan polos**: 1.2.1 batas bawah **0** (`@if(baru<lama,0,baru−lama)`), 1.2.3 **prorata** bila master `IsProRate`, 1.2.4 pajak dihitung ulang (`FlagPPH` ∧ `EDMType=="3"`), 1.2.5 **mentah** untuk `EDMType` 4/2. ⇒ *"satu rumus untuk semua"* (AC 16) bila dipaksakan ke NonProp **melanggar AC 22**. Migrasi 363 mengikuti XML — **butir WO** | `CalculateDifferenceEDM_act.xml` langkah 1.2.1–1.2.5; `../../backend/migrations/363_t_polis_xol_layer_difference.sql` baris 2–4 |

## Hasil & nilai pengguna

Selisih antar generasi dihitung **di lapisan layanan**, dengan **satu rumus** untuk semua kasus —
endorsemen pertama maupun berlapis, proporsional maupun non-proporsional:

```
selisih.X = baris_ini.X − baris(generasi sebelumnya).X
```

⛔⛔ **Rumus ini TIDAK pernah ditulis di dalam basis data.** View yang menghitung selisih
**dibatalkan** — rumus di dua tempat cepat atau lambat bercabang.
⛔ **Tiket ini tidak boleh menghidupkannya kembali.**

## Yang dibangun

Perhitungan selisih dengan empat aturan turunan:

| Golongan medan | Perlakuan |
| --- | --- |
| **uang** | dikurangi |
| **persentase** | ⭐ **disalin**, tidak dikurangi |
| **kunci** | disalin apa adanya |
| **arah utang-piutang** | diturunkan dari **tanda** selisih — ⚠️ *hanya lapisan XOL (koreksi 06-10)* |

~~⭐ Pada selisih lapisan, arah utang-piutang diturunkan dari **jumlah sepanjang daftar lapisan**,
bukan dari satu lapisan saja.~~ ⭐ *(koreksi 06-10)* Pada selisih lapisan, arah utang-piutang diturunkan
dari tanda selisih `DueToValue` **lapisan itu**; induk lapisan tidak bertabel (ID-7).

⭐ Lapisan penyimpanan hanya mengambil **dua baris** — baris ini dan baris yang ditunjuknya.
⛔ Nol agregasi di basis data.

⚠️ `[penyimpangan sadar]` **Sistem lama punya DUA rumus di sisi proporsional**, dan yang kedua
mengurangi terhadap **selisih** generasi lampau, bukan terhadap **nilainya**. Dengan angka contoh:
nilai baru 180, nilai lama 150, selisih lama 50 ⇒ varian kedua menghasilkan **130**, padahal yang
benar **30**. ⛔ **Varian kedua tidak ditiru** — `[keputusan work owner]`.

## Batas — yang TIDAK termasuk

⛔ Penulisan hasilnya ke tabel proyeksi — tiket **07**.
⛔ Penandaan baris yang lahir dari rumus lama — tiket **09**.

## Cara mengujinya

Operannya lewat seam `repository`; rumusnya lewat lapisan layanan.

⭐ **Uji regresi wajib:** sisi non-proporsional harus menghasilkan angka yang **sama** sebelum dan
sesudah penyeragaman — `[terverifikasi]` ia memang sudah seragam, dan penyeragaman tidak boleh
mengubahnya. Uji yang menemukan perubahan berarti rumusnya salah diterjemahkan.

## Acceptance criteria

- [ ] **AC 16** — satu rumus untuk **seluruh** generasi
- [ ] **AC 17** — medan **uang** berisi hasil pengurangan
- [ ] **AC 18** — medan **persentase** disalin, tidak dikurangi
- [ ] **AC 19** — medan **kunci** disalin apa adanya
- [ ] **AC 20** — arah utang-piutang diturunkan dari **tanda** selisih *(koreksi 06-10: hanya lapisan XOL; proporsional tanpa `DUE_TO`)*
- [ ] **AC 21** — pada selisih lapisan, arah itu diturunkan dari ~~**jumlah** sepanjang daftar~~ **tanda selisih `DueToValue` lapisan itu** *(koreksi 06-10)*
- [ ] **AC 22** — sisi non-proporsional menghasilkan angka **sama** sebelum dan sesudah *(koreksi 06-10: rumus XML 1.2.1–1.2.5 dipertahankan; butir WO atas cakupan AC 16)*

## ⚠️ ~~Satu hal yang tidak terbaca dari sistem lama~~ — ✅ terbaca *(koreksi 06-10)*

~~`[terbuka]` Kedua blok rumus di sistem lama **sama-sama bergerbang aktif**, dan pemilihnya hanya
tertulis di **keterangan langkah**, bukan di gerbang yang dijalankan. ⇒ Mana yang benar-benar
berjalan **tidak dapat dinyatakan dari ekspor**.~~ ⭐ Pemilihnya **`.OldData.EDMNo==""`** pada prasyarat
langkah 1 (salah ⇒ lompat `HasEDMNo`) dan langkah 4 (benar ⇒ keluar) — lihat blok KOREKSI.

⭐ **Tidak menahan tiket ini** — varian kedua tidak ditiru apa pun jawabannya.