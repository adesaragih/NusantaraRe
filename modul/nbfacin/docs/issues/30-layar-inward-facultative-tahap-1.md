# 30: Layar Inward Facultative — tahap 1 (General, ringkasan, Show Detail, tombol kaki)

> Disusun agent (sesi `nusantarare-0f`) dari XML korpus + tangkapan layar Pega work owner 03-10-2026 atas
> permintaan work owner ("kenapa blm lari ke flow InputInwardFacultative", lalu "lanjut") — **bukan hasil
> `/to-tickets`**. Tangkapan layar kasus FIRE tidak disalin ke repo (memuat nama pelanggan dan nama orang).

**What to build:** sesudah **Create opportunity** berhasil, layar pindah ke assignment pertama case NB —
section `InputInwardFacultative` — dengan blok **General** (section `Periode`), ringkasan objek, checkbox
**Show Detail** + strip tab detail, dan tombol kaki flow action.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- Flow `InputInwardFacultativeOffer`: Start → konektor (mengisi PositionNote, IsCedingConfirm=Offer,
  FlagOnGoingPolicy=0, NBStatus "NEW NB") → assignment MARKETING → flow action `InwardFacultative` → section
  `InputInwardFacultative`. Pra-tampil `InwardFacultative_PreDT` (OfferingDate kosong → hari ini).
- Susunan tingkat atas `InputInwardFacultative`: Periode → `AllSummarySection` (FIRE: `FireSummarySection`) →
  checkbox Show Detail (`pyCheckboxCaption`) → `InputInwardFacultativeDtl` / `_IsUW` (bila Show Detail atau Admin;
  `ConfirmBinding`) → `ComfirmPolis_NonLife` / `_Life`.
- Label: `labels.ts` `PERIODE` (sel + tag, diuji `labels.test.ts` terhadap `Periode.xml`), `TAB_DETAIL` (`pyTitle`
  `InputInwardFacultativeDtl.xml`), `KOLOM_RINGKASAN` (`FireSummarySection.xml`), `TOMBOL_KAKI_INWARD` (flow action:
  Submit / Save for later / Cancel).
- Sel 42 "Class of business" = `.QuotationData.BusinessName` — di tangkapan layar berisi Group Business (FIRE),
  bukan Class Of Business form Opportunity.

## Keputusan agent (menunggu konfirmasi work owner)

- **D-1** nilai awal dari isian Create opportunity yang dikirim (belum ada endpoint baca case): Business status,
  Insured name (akun ChooseAccount), Class of business = Group Business, Type facultative; Offering date = hari ini.
- **D-2** medan yang sumbernya belum tersambung (Source of business, Ceding co name, Group Name, Old Policy Number,
  Risk Scoring, pilihan Marketing Name — RD `BrowseMarketingOfficer_RD` / tabel `MARKETINGOFFICER`) tampil kosong.
- **D-3** pilihan radio Policy Type (Individual Policy / Master Policy) dan Day (365 / 366) dari tangkapan layar
  (definisi properti tidak ada di korpus); belum terpilih pada case baru `[dugaan]`.
- **D-4** tombol fitur lain (upload, Change SOB/Ceding Co, Search, CSV) serta Submit / Save for later nonaktif
  sampai backend-nya ada; Cancel kembali ke portal. Isi tab detail = tahap 3 (`BelumTersedia`).
- **D-5** judul "SUMMARY" dan "No items" hanya dari tangkapan layar (tidak ditemukan di XML).

## Tahap berikutnya

- Tahap 2: endpoint baca/simpan case (backend sesi c3) → Save for later / Submit, pilihan Marketing Name, Source
  of business / Ceding co.
- Tahap 3+: isi tab detail per lini, mulai FIRE (Object → Coverage → Clauses → Spreading → …).

**Status:** ready-for-human — tahap 1 dibangun 03-10-2026 (uji frontend nbfacin 139/139); menunggu tinjauan D-1…D-5

- [x] Sesudah Create opportunity berhasil, pindah ke layar Inward Facultative case `NB-<n>`
- [x] Blok General: label verbatim korpus, susunan = tangkapan layar FIRE
- [x] Ringkasan + Show Detail + strip tab; tombol kaki Submit / Save for later / Cancel
- [ ] Baca/simpan case — tahap 2
- [ ] Isi tab detail — tahap 3+

## Comments
