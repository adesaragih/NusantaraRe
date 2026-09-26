# 10: Kontrak Komite — penyerahan kasus

**Status:** claimed

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
- [x] Nilai uang yang menyeberang memakai representasi yang sama dengan di dalam sistem — tidak
      dikonversi menjadi *floating point* di batas. *(AC 22 spec)*
- [x] Penyerahan hanya mungkin pada baris yang **masih Outstanding**.
- [x] **Gerbang rekening pembayaran** `[terverifikasi]`: penyerahan **ditolak** bila salah satu dari
      **nama bank**, **id bank**, atau **nomor rekening** pada baris yang diserahkan **kosong**,
      dengan pesan yang setara `"Name of bank cannot be empty"`. Ini **paritas perilaku existing**,
      bukan penyimpangan — sumbernya `Claim Life/Activity/GetListKomiteLife.xml`
      (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `GETLISTKOMITELIFE` / `RULE-OBJ-ACTIVITY`), prasyarat
      `.NameOfBank=="" || .NoAccount=="" || .IDOfBank==""`. *(AC 57 spec)*
- [x] Penolakan gerbang rekening terjadi di **lapisan layanan**, dan tetap terjadi meskipun kontrol
      UI-nya ditampilkan. *(sejalan AC 12 spec)*
- [x] Penyerahan berlaku untuk **semua** adjustment, bukan hanya yang di atas ambang nilai tertentu.
      *(`[keputusan work owner]`, langkah 4 mesin status)*
- [x] Wewenang penyerahan mengikuti aturan tiket 07: `QP`/`QR` hanya SPV; `TP`/`TR` bebas peran.
- [x] Perubahan pada bentuk muatan diperlakukan sebagai **perubahan kontrak lintas konteks**, dan
      ditandai demikian di kode.
- [ ] **Jumlah tingkat komite = COUNT baris roster `EMAILKOMITE` yang aktif (`STS_AKTIF = "1"`) dan
      ber-`LIMIT_BOTTOM <= CLAIM_AMOUNT`** — dihitung saat penyerahan, **tidak** dibaca dari
      konstanta mana pun.
- [x] Ambang yang dipakai mencari roster adalah **nilai mutlak** klaim: klaim bernilai negatif
      dicari dengan tandanya dihilangkan — `@if(IsADj<0, IsADj*-1, IsADj)`.
- [x] ⚠️ `[keputusan work owner]` Bila tidak ada baris roster yang cocok, penyerahan **gagal
      terang-terangan** — bukan diam-diam membuat tangga nol tingkat. **Penjaga defensif**, bukan
      alur normal: roster dijamin **≥ 1** secara bisnis (limit berjenjang selalu menutup nilai
      klaim). `[terverifikasi]` Pega **tidak** punya gerbang ini — `CreateKMTLife_Act` men-set
      `KomiteCount = 1` meski `KomiteLoop = 0`.
- [x] Semua baris adjustment pada satu klaim bermata uang sama sebelum penyerahan (invariant
      **OQ-060**).

### Rujukan Komite ⚠️ BARU 2026-09-16

- [ ] ⚠️ Penyerahan ke Komite menyimpan **`KOMITE_ID`** = **identitas kasus komite**, yaitu
      `T_WORK_CLAIM.ID` baris komite yang baru lahir; baris yang belum pernah dikirim ber-`KOMITE_ID`
      **`NULL`**. *(AC 61 spec; tiket 14; penyimpangan sadar — rujukan, bukan salinan)* — ⚠️
      `[keputusan work owner]` REVISI 2026-09-17; bentuk penautan ini **tidak ada di korpus Pega**,
      lihat catatan korpus di tiket 14 §`KOMITE_ID`. `[terbuka]` tipenya mengikuti tipe
      `T_WORK_CLAIM.ID` yang belum ditetapkan.
- [x] ⚠️ **Roster dan keputusan komite per anggota TIDAK disimpan di konteks ini.** Test yang
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
- [x] ⚠️ Rujukan memakai **ID stabil**, bukan indeks posisi. Test yang menemukan padanan
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

---

## Pembacaan ulang XML — 27 September 2026

Empat rule dibaca ulang; jendela `D:\XML\RNM_BRD\Claim Life\**\*.xml`.

| Rule (path) | Baris pecahan | Yang dipastikan |
| --- | --- | --- |
| `Activity/GetListKomiteLife.xml` | 402-403 | `Local.IsADj = .CLAIM_AMOUNT` |
| idem | **449-450** | `local.errmsg = "Name of bank cannot be empty"` — teks gerbangnya, `[terverifikasi]` |
| idem | 587-589 | pesan dipasang ke `.NameOfBank` sebagai `pyMessageLabel` |
| idem | **655** | `.NameOfBank==""\|\|.NoAccount==""\|\|.IDOfBank==""`, `WhenTrue=2` (lanjut), `WhenFalse=3` (lewati) |
| idem | **790** | `@if(Local.IsADj<0,Local.IsADj* -1,Local.IsADj)` — ambang MUTLAK |
| idem | 810-811 | `Param.STS_KLAIM = "LIFE"` |
| `ReportDefinition/FilterEmailKomiteWithLimit.xml` | **662** | `pyFilterLogic = A AND C AND B` |
| idem | 670-680 / 683-697 / 701-715 | A `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM` · **C** `.STS_KLAIM = Param.STS_KLAIM` · **B** `.STS_AKTIF = "1"` — ⚠️ ralat: ronde pertama menukar B dan C |

**Pertanyaan yang XML tidak jawab:**

- Mata uang pita roster. `POOLDATA.EMAILKOMITE` `[data DBA]` tanpa kolom mata uang, sehingga
  `LIMIT_BOTTOM` berlaku atas satu mata uang implisit. `[terbuka — work owner]`
- Apa yang terjadi bila baris diserahkan dua kali. Korpus tidak punya gerbangnya; di sini
  `KOMITE_ID IS NULL` di klausa `WHERE` yang memutuskannya. **Penyimpangan sadar.**
- Nasib penyerahan bila klaim bermata uang campur. Korpus diam; OQ-060 yang menjawabnya.

## Implementasi — 27 September 2026 (tiket 10)

**Status: `claimed`** — **12 dari 17 AC milik tiket ini tertutup.** Titik tetap `2598778`.

⚠️ **Jendela penyebutnya dinamai.** Berkas ini memuat **18** kotak centang, tetapi satu di
antaranya — *"✅ TERTUTUP 2026-09-17 … `T_GENERAL_KOMITE.ADJUSTMENT_ID`"* — ditutup **tiket 00
Komite Claim Life**, bukan tiket ini. Penyebut pekerjaan tiket 10 karena itu **17**, bukan 18.
Dihitung dua cara: `grep -c '^- \[x\]'` → 12 dan `grep -cE '^- \[[x ]\]'` → 18; lalu pemindaian
per bab → Acceptance criteria 10/13, Rujukan Komite 2/5 dengan 1 milik tiket lain. `[terverifikasi]`

| Angka | Nilai | Cara 1 | Cara 2 | Label |
| --- | --- | --- | --- | --- |
| test Go | **200 PASS · 0 FAIL · 31 SKIP** *(dari 189/28)* | `go test -tags=db ./internal/... -v`, cacah awalan `--- PASS`/`--- FAIL`/`--- SKIP` → 231 | `grep -rhoE '^func Test[A-Za-z0-9_]+' internal/ --include=*_test.go \| wc -l` → **231** (tanpa `sort -u`; lihat ralat di bawah) | `[terverifikasi]` |
| vet | bersih, termasuk `-tags=db` | `go vet ./...` | `go vet -tags=db ./...` | `[terverifikasi]` |
| gofmt | nol berkas | `gofmt -l internal/ cmd/ pkg/` | `gofmt -l` atas tiap berkas baru satu per satu | `[terverifikasi]` |
| test JS | **11 PASS** *(dari 7)* | `npm test` → `Tests 11 passed` | `grep -c "  it(" src/services/api.test.ts` → 11 | `[terverifikasi]` |
| modul frontend | **88** | `npm run build` → `88 modules transformed` | 8 berkas sumber `.ts`/`.tsx` + dependensi node | `[terverifikasi]` |

⛔ **Ralat 27-09-2026 — selisih itu cacat CARA HITUNG saya, bukan selisih sungguhan.** Cara 2
memakai `sort -u`, yang membuang `TestPagarSkemaUjiMenuntutDuaSyarat` — nama yang memang ada di
**dua** paket dan memang dijalankan dua kali. Tanpa `sort -u` kedua cara sepakat **231**.
CLAUDE.md §4 bab 4a melarang memilih salah satu bila keduanya berselisih; yang benar di sini
bukan memilih 231, melainkan **membetulkan cara keduanya** sehingga tidak ada selisih untuk
dipilih.

| Berkas | Isi |
| --- | --- |
| `services/komite.go` *(baru)* | gerbang, ambang, muatan, tiga antarmuka yang gagal terang, `Serahkan` |
| `services/komite_test.go` *(baru)* | 9 kasus tanpa Oracle |
| `services/komite_db_test.go` *(baru)* | 3 kasus terhadap skema uji |
| `services/komite_statik_test.go` *(baru)* | 2 penjaga batas konteks |
| `repository/klaimlife.go` | `PerbaruiKomiteID` — `WHERE … AND KOMITE_ID IS NULL` |
| `handlers/komite.go` *(baru)* + `handlers.go` | satu rute, sepuluh pemetaan galat |
| `frontend/` `api.ts` · `KlaimLife.tsx` · `api.test.ts` | `serahkanKeKomite`, `pesanGalat`, kolom Komite |

⭐ **Pemanggil produksi pertama `WajibWewenangKomite`.** Tiket 07 melahirkannya tanpa pemanggil dan
AC-nya sengaja dibiarkan terbuka; `Serahkan` menutupnya di sini.

**Gerbang dagang, tiap-tiap dibuktikan gagal lalu dipulihkan:**

| Gerbang | Mutasi | Hasil |
| --- | --- | --- |
| rekening | `TrimSpace(...)==""` → `false` | merah |
| ambang mutlak | `mutlak.Negative = false` dihapus | merah |
| tangga nol tingkat | `tingkat <= 0` → `tingkat < 0` | merah |
| invariant mata uang | pembandingan dilucuti | merah |
| serah ganda | penjagaan `KomiteID` dilucuti | merah |
| batas konteks *(statik)* | `KomiteAproval` + `IndexAdjustment` disisipkan | merah, keduanya disebut |
| tak memutus status *(statik)* | `.KodeStatus =` disisipkan | merah |

⛔ **Penjaga statik kedua semula salah tuduh.** Ia melarang `KodeStatus:` dan langsung menangkap
`MuatanKomite{KodeStatus: …}` — yang justru AC 24: status baris saat penyerahan memang ikut
menyeberang. Membaca status untuk dilaporkan bukan memutuskannya; polanya dipersempit ke penulisan.

### ⛔ AC yang belum tertutup — 6

| AC | Sebab |
| --- | --- |
| kasus anak Komite dibuat | menuntut tabelnya — butir **af**, `[USULAN]`. Yang ada `KasusKomiteBelumDiputuskan`, gagal terang → HTTP 501 |
| muatan memuat nilai klaim + `CURRENCY` + status baris | bentuknya ada dan diisi; **pembuktiannya** hanya di test db yang MELEWATI tanpa Oracle. Melewati bukan lulus |
| jumlah tingkat = COUNT roster | nol konstanta di kode — tetapi query-nya menuntut `POOLDATA.EMAILKOMITE`, tabel **produksi**. Butir **af** + persetujuan manusia |
| `KOMITE_ID` = identitas kasus komite | penulisnya ada dan ber-`WHERE … IS NULL`; pembuktiannya di test db yang melewati |
| keputusan komite dibaca lewat join tiga lompatan | menuntut `T_GENERAL_KOMITE` — konteks Komite Claim Life, butir **af** |
| penyerahan membawa penunjuk baris ke kasus anak | sama: menunggu **af** |

⚠️ **Yang TIDAK saya putuskan sendiri:** nilai `KMT-xxxxxx` tidak dikarang — `PembuatKasusKomite`
mengembalikan ID dan `Serahkan` memakainya apa adanya. Penomorannya milik butir **o1–o3**.

### Hasil /code-review — titik tetap `2598778`

**Standards — 3 pelanggaran keras + 1 kelemahan terbukti + 5 bau. Semua yang keras diperbaiki.**

| Temuan | Tindakan |
| --- | --- |
| ⛔ **bypass terbukti**: peta pengecualian penjaga statik berkunci `filepath.Base`, sehingga berkas bernama sama **di mana pun** ikut dikecualikan; akar telusurnya `internal/` saja; `diperiksa++` berada **sesudah** pengecualian sehingga swa-periksa mustahil menangkapnya. Reviewer membangun bypass-nya dan menjalankannya — hijau | **diperbaiki bertiga**: kunci jalur penuh, akar modul, cacah sebelum pengecualian, plus penguncian jumlah pengecualian |
| §4 rule 1: `CreateKMTLife_Act` dikutip 5× tanpa nomor baris | **diperbaiki** — 1322-1323, 1377-1378, 1398-1399, 1440-1441, 1543 |
| §4 rule 9 + §4a: "12 dari 18" tanpa perintah audit, tanpa label, satu arah, jendela tak bernama | **diperbaiki** — dua arah, jendela dinamai, penyebut dibetulkan ke **17** |
| §4a: telemetri satu arah | **diperbaiki** — dua cara; selisih 1 dijelaskan |
| bau: tiga `Dengan*` kembar · `GalatRekening` nol konsumen · `stsKlaim` menyesatkan · `pesertaDariPeta` Feature Envy | `GalatRekening` **dibuang**; `stsKlaim` → `lini`; `PeriksaSatuMataUang` menerima `[]Peserta`. Tiga `Dengan*` **dibiarkan** — menyatukannya menuntut generik yang menyamarkan tiga ketergantungan berbeda jadi satu |
| salah kutip "CLAUDE.md bab 4 butir 10" | **diterima, tidak diperbaiki** — salah kutip yang sama sudah ada di `register.go` dan `tolak.go`; membetulkannya di satu tempat saja membuat ketiganya tampak berbeda padahal sama. Milik tiket pembersihan |

**Spec — 1 cacat nyata + 2 ralat kutipan + 3 scope. Cacatnya yang termahal.**

| Temuan | Tindakan |
| --- | --- |
| ⛔ **`PeriksaSatuMataUang` menjaga kolom yang SALAH.** Ia membandingkan `CurrencyID` (`CURRENCY_ID`), padahal yang menyeberang ke Komite lewat `AmbangRoster` adalah `JumlahKlaim.Currency` (`CURRENCY`) — kolom lain. Klaim ber-`CURRENCY` campur di balik `CURRENCY_ID` seragam **lolos utuh**, persis kegagalan yang komentar fungsi itu sendiri peringatkan. Nol test menangkapnya: keduanya hanya mengisi `CurrencyID` | **diperbaiki** — kedua kolom dijaga, masing-masing dibuktikan merah sendiri-sendiri; test db kini mengisi keduanya |
| AC gerbang rekening "di lapisan layanan" dicentang padahal buktinya hanya test db yang MELEWATI — standar yang sama yang saya pakai untuk **mencabut** centang AC muatan | **diperbaiki dengan menutupnya sungguhan**: predikat tampil layar diangkat jadi `bolehSerahkanDiLayar` dan diuji — ia mengembalikan `true` untuk baris tanpa bank sama sekali. Layar memang menampilkan tombolnya, jadi bila gerbangnya tidak di services ia tidak ada di mana pun. Dua test yang **berjalan** |
| label penyaring B/C tertukar di kutipan | **diralat** di kode dan di tiket |
| tafsir "tangga satu tingkat" bukan bunyi XML | **diturunkan ke `[dugaan]`**; yang `[terverifikasi]` hanya bahwa `KomiteCount` konstan `1` tanpa prasyarat |
| `GalatRekening` nol konsumen | **dibuang** |
| komentar `pastikanSatuBaris` teryatim oleh sisipan | **dipulihkan** |

⭐ Dua review menemukan dua hal berbeda, dan keduanya hal yang saya tidak akan temukan sendiri:
Standards menemukan penjaga yang **lebih lemah dari klaimnya**, Spec menemukan penjaga yang
**menjaga kolom yang salah**. Cacat kedua akan lolos ke produksi tanpa satu pun test memerah.

**lanjut dari sini:** tiket 10 selesai. Lima AC sisanya menunggu butir **af** disahkan work owner.
Berikutnya tiket 12.

