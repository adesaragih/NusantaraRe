# 10: Kontrak Komite — penyerahan kasus

**Status:** ready-for-agent

**Blocked by:** 07 (penegakan peran + wewenang per `Type`)

## Hasil & nilai pengguna

Sebagai **ReasLifeSPV**, saya dapat menyerahkan baris adjustment ke **Komite Life** untuk diputuskan
— dan Komite menerima nilai klaim beserta mata uangnya, sehingga mereka tidak perlu membaca balik ke
sistem ini. *(User story 17, 18, 19 di spec)*

Komite adalah **sistem luar**; tiket ini membangun **batasnya**, bukan isinya.

## Area codebase

`internal/models` (muatan penyerahan), `internal/repository` (pembuatan rekam kasus Komite),
`internal/services` (aturan penyerahan + pemilihan roster), `internal/handlers` (endpoint serahkan),
`frontend/` (kontrol "Send ke Komite" pada baris adjustment).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/CreateKMTLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`, 121.652 byte | `[terverifikasi]` sepuluh langkah: `Property-Set` → `Call pxRetrieveReportData` → `Property-Set` ×3 → `Call pxAddChildWork` → `Obj-Refresh-And-Lock` → `Property-Set` → `Obj-Save` → `Call SendEmailKlaimLF` |
| `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml` | `ASM-FW-GCNMFW-INT-EMAILKOMITE` / `FILTEREMAILKOMITEWITHLIMIT` / `RULE-OBJ-REPORT-DEFINITION`, 67.827 byte | `[terverifikasi]` pemilihan roster; menerima `Param.LIMIT_BOTTOM` dan `Param.STS_KLAIM` |
| `Claim Life/Section/ClaimComite.xml`, `Claim Life/Harness/Committe_Life.xml` | — | `[terverifikasi]` pemicu dari UI (`<pyActivity>CreateKMTLife_Act</pyActivity>`, 2× masing-masing) — **bukan** shape flow |

`[terverifikasi]` Kelas kasus anak: `ASM-FW-GCNMFW-Work-KomiteLife`. Muatan yang menyeberang
sekarang:

```
childPageKomite.CLMNO            childPageKomite.KomiteCount
childPageKomite.KomiteLoop       childPageKomite.IndexAdjustment
childPageKomite.IndexPremiumList childPageKomite.KomiteList(<APPEND>).KomiteID
childPageKomite.KomiteList(<LAST>).IDKomite / .KomiteAproval / .KomiteEmail
sisi induk: .IsKomite  .KomiteNo  .TotalKomite
```

## ADR terkait

**ADR-0001** (tiga kontrak batas; muatan **diperluas** atas keputusan work owner),
**ADR-0003** (nilai uang yang menyeberang memakai representasi yang sama di kedua sisi),
**ADR-0012** (siapa yang boleh menyerahkan, bergantung `Type`).

## Acceptance criteria

- [ ] Penyerahan membuat kasus anak berkelas Komite Life, membawa penunjuk **baris** yang diserahkan.
- [ ] Muatan penyerahan memuat **nilai klaim**, **`CURRENCY`**, dan **status baris saat penyerahan**
      — tiga hal yang **tidak** ada di sistem lama. *(AC 24 spec; `[keputusan work owner]`)*
- [ ] Nilai uang yang menyeberang memakai representasi yang sama dengan di dalam sistem — tidak
      dikonversi menjadi *floating point* di batas. *(AC 22 spec)*
- [ ] Penyerahan hanya mungkin pada baris yang **masih Outstanding**.
- [ ] **Gerbang rekening pembayaran** `[terverifikasi]`: penyerahan **ditolak** bila salah satu dari
      **nama bank**, **id bank**, atau **nomor rekening** pada baris yang diserahkan **kosong**,
      dengan pesan yang setara `"Name of bank cannot be empty"`. Ini **paritas perilaku existing**,
      bukan penyimpangan — sumbernya `Claim Life/Activity/GetListKomiteLife.xml`
      (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`), prasyarat
      `.NameOfBank=="" || .NoAccount=="" || .IDOfBank==""`. *(AC 57 spec)*
- [ ] Penolakan gerbang rekening terjadi di **lapisan layanan**, dan tetap terjadi meskipun kontrol
      UI-nya ditampilkan. *(sejalan AC 12 spec)*
- [ ] Penyerahan berlaku untuk **semua** adjustment, bukan hanya yang di atas ambang nilai tertentu.
      *(`[keputusan work owner]`, langkah 4 mesin status)*
- [ ] Wewenang penyerahan mengikuti aturan tiket 07: `QP`/`QR` hanya SPV; `TP`/`TR` bebas peran.
- [ ] Perubahan pada bentuk muatan diperlakukan sebagai **perubahan kontrak lintas konteks**, dan
      ditandai demikian di kode.
- [ ] **Jumlah tingkat komite = COUNT baris roster `EMAILKOMITE` yang aktif (`STS_AKTIF = "1"`) dan
      ber-`LIMIT_BOTTOM <= CLAIM_AMOUNT`** — dihitung saat penyerahan, **tidak** dibaca dari
      konstanta mana pun.
- [ ] Ambang yang dipakai mencari roster adalah **nilai mutlak** klaim: klaim bernilai negatif
      dicari dengan tandanya dihilangkan — `@if(IsADj<0, IsADj*-1, IsADj)`.
- [ ] ⚠️ `[keputusan work owner]` Bila tidak ada baris roster yang cocok, penyerahan **gagal
      terang-terangan** — bukan diam-diam membuat tangga nol tingkat. **Penjaga defensif**, bukan
      alur normal: roster dijamin **≥ 1** secara bisnis (limit berjenjang selalu menutup nilai
      klaim). `[terverifikasi]` Pega **tidak** punya gerbang ini — `CreateKMTLife_Act` men-set
      `KomiteCount = 1` meski `KomiteLoop = 0`.
- [ ] Semua baris adjustment pada satu klaim bermata uang sama sebelum penyerahan (invariant
      **OQ-060**).

### Rujukan Komite ⚠️ BARU 2026-09-16

- [ ] ⚠️ Penyerahan ke Komite menyimpan **`KOMITE_ID`** = **identitas kasus komite**, yaitu
      `T_WORK_CLAIM.ID` baris komite yang baru lahir; baris yang belum pernah dikirim ber-`KOMITE_ID`
      **`NULL`**. *(AC 61 spec; tiket 14; penyimpangan sadar — rujukan, bukan salinan)* — ⚠️
      `[keputusan work owner]` REVISI 2026-09-17; bentuk penautan ini **tidak ada di korpus Pega**,
      lihat catatan korpus di tiket 14 §`KOMITE_ID`. `[terbuka]` tipenya mengikuti tipe
      `T_WORK_CLAIM.ID` yang belum ditetapkan.
- [ ] ⚠️ **Roster dan keputusan komite per anggota TIDAK disimpan di konteks ini.** Test yang
      menemukan tabel/kolom penyimpan `KomiteAproval`, `KomiteComment`, atau `DateApprove` di Claim
      Life **gagal**. Keduanya milik **Komite Claim Life**. *(AC 61 spec; tiket 14; penyimpangan sadar)*
- [ ] ⚠️ Keputusan komite **dibaca lewat join**, bukan disalin ke adjustment. Rantainya **tiga
      lompatan**, bukan satu: `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` → `T_WORK_CLAIM` (`ID` = `KOMITE_ID`,
      `COVER_KEY` = `ID` baris klaim) → `T_GENERAL_KOMITE` (**`ID` = `T_WORK_CLAIM.ID`, shared
      primary key — ⛔ tidak ada kolom `WORK_CLAIM_ID`, REVISI 2026-09-18**) →
      `T_KOMITE_KOMITELIST` (keputusan per anggota, diurut `KOMITE_URUT`). *(tiket 14 §`KOMITE_ID`;
      tiket 00 Komite Claim Life)* — ⚠️ `[keputusan work owner]` REVISI 2026-09-17; rantai ini
      **bentuk baru, tidak ada di korpus Pega**. Di Pega penautannya lewat `pxAddChildWork` +
      `CLMNO` + indeks posisi; lihat catatan korpus di tiket 14 §`KOMITE_ID`.
- [ ] ⚠️ Rujukan memakai **ID stabil**, bukan indeks posisi. Test yang menemukan padanan
      `IndexPremiumList` / `IndexAdjustment` sebagai kunci rujukan **gagal**. *(`[terverifikasi]`
      `Claim Life/Activity/CreateKMTLife_Act.xml` memakai `.pxListSubscript`; penyimpangan sadar; **AC 62 spec**)*
- [ ] ✅ **TERTUTUP 2026-09-17** — tabel komite **sudah** menampung keputusan per baris adjustment
      lewat **`T_GENERAL_KOMITE.ADJUSTMENT_ID`** (FK → `T_CLAIMLF_ADJUSTMENT.ID`), ditetapkan di
      **tiket 00 Komite Claim Life** (§Tabel, plus AC "`ADJUSTMENT_ID` berada di `T_GENERAL_KOMITE`,
      **BUKAN** di `T_WORK_CLAIM`"). Penunjuk dua arah — `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` dan
      `T_GENERAL_KOMITE.ADJUSTMENT_ID` — diisi dalam **satu transaksi** saat kirim komite.

## Catatan penutupan (2026-09-14)

**OQ-032 TERTUTUP** `[terverifikasi]` — tangga komite **sepenuhnya data-driven**, tanpa konstanta:

1. `Claim Life/Activity/GetListKomiteLife.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
   `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`) set `Local.IsADj = .CLAIM_AMOUNT`, panggil report
   `FilterEmailKomiteWithLimit`, hasil ke `.KomiteList`.
2. `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml` (`ASM-FW-GCNMFW-INT-EMAILKOMITE` /
   `FILTEREMAILKOMITEWITHLIMIT` / `RULE-OBJ-REPORT-DEFINITION`) filter
   `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM && .STS_KLAIM = Param.STS_KLAIM && .STS_AKTIF = "1"`.
3. `Claim Life/Activity/CreateKMTLife_Act.xml` set
   `childPageKomite.KomiteLoop = @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)`.

**OQ-037 TERTUTUP** — ambang adalah **data di tabel** `POOLDATA.EMAILKOMITE`
(`LIMIT_BOTTOM INTEGER`, `LIMIT_TOP INTEGER`, `STS_AKTIF`, `STS_KLAIM`), **bukan** hardcode seperti
di Komite Claim FacIn. ⚠️ `[data DBA]` roster **tanpa kolom mata uang** → pita berlaku atas satu
mata uang implisit.

**OQ-060 TERTUTUP** — bentuk uang `(amount, currency)` per baris, **invariant: satu klaim satu mata
uang**.

`[terbuka]` **OQ-035** — kepemilikan `serviceInsertArasapasClaimLife_act` (satu salinan, dua
konteks). **Tidak memblokir** tiket ini; menyentuh tiket 12.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
