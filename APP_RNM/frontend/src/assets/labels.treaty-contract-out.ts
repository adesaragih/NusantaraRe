// Label modul Treaty Contract Out — tiket 02 →.
//
// ⛔ VERBATIM dari korpus `D:\XML\RNM_BRD\Treaty Contract Out\`, dengan path +
// baris di sebelahnya (berkas korpus satu tag per baris; nomor baris mentah =
// nomor `sed -e 's/></>\n</g'`). Label yang diketik dari ingatan adalah label
// yang bergeser, dan yang bergeser tidak berbunyi.
//
// ⛔ Nol kata "Old" dan "testing" di nama apa pun (penyimpangan sadar 8);
// dijaga `components/treaty-contract-out/namaJujur.test.ts`.

/**
 * Menu kelompok **Treaty Contract Out** — tiga harness portal (kelas
 * `Data-Portal`), butirnya VERBATIM `<pyLabel>` tiap harness.
 *
 * ⚠️ Nama kelompok = nama FOLDER korpus, pola `MODUL` di `labels.ts`
 * (`MODUL.treatyContractOut`, ditambahkan ADITIF di tiket 03).
 */
export const MENU_TCO = {
  /** Nama folder korpus `D:\XML\RNM_BRD\Treaty Contract Out`. */
  kelompok: 'Treaty Contract Out',
  /** `Harness/InboxTreatyContract.xml` b151 `<pyLabel>`. */
  inboxTreatyContract: 'InboxTreatyContract',
  /** `Harness/InboxTreatyContractReinsType.xml` b151 `<pyLabel>`. */
  inboxTreatyContractReinsType: 'InboxTreatyContractReinsType',
  /** `Harness/InboxTreatyContractDescription.xml` b359 `<pyLabel>`. */
  inboxTreatyContractDescription: 'InboxTreatyContractDescription',
} as const

/**
 * Pemilih jenis reasuransi — tiket 02.
 *
 * Dua label korpus untuk medan yang SAMA (`ReinsTypeID` dari master
 * `REINSURANCETYPE`): form kontrak memakai `ReinsType`, form tahun memakai
 * `Reinsurance Type`. Keduanya dibawa apa adanya; pemanggil memilih.
 */
export const JENIS_REASURANSI_TCO = {
  /** `Section/InputTreatyContractReinsType.xml` b2652 `<pyLabelFieldValue>`. */
  reinsType: 'ReinsType',
  /** `Section/InputTreatyContract.xml` b7925 `<pyLabelFieldValue>` (dan b7948 `<pyLabelPreview>`). */
  reinsuranceType: 'Reinsurance Type',
  /** `[tidak ada di korpus]` — kosakata kami; keadaan yang ADR-0015 tuntut terlihat. */
  masterKosong: 'Master jenis reasuransi kosong atau tidak terbaca.',
  /** `[tidak ada di korpus]` — teks pilihan kosong `Pilih`. */
  belumDipilih: '-- pilih jenis reasuransi --',
} as const

/**
 * Layar tahun treaty — tiket 03: `Harness/InboxTreatyContract.xml` →
 * `Section/GridTreatyContract.xml` → `Section/InputTreatyContract.xml`
 * (grid `BrowseTreatyYear_RD` b17355) + `Section/InputDtlTreatyContact.xml`
 * (form lengkap "Input New Data", disertakan b1206).
 *
 * ⚠️ OQ-TCO-05 — LABEL BERSILANG DI KORPUS, dibawa apa adanya: di GRID
 * kolom `TREATYYEAR` berjudul `Transaction Year` (b17531 → `.TreatyYear`
 * b19125) dan `UNDERWRITINGYEAR` berjudul `Underwriting Year` (b17378 →
 * `.UnderwritingYear` b18964); di FORM `TREATYYEAR` berlabel
 * `Underwriting Year` (b8004 → `InputTreatyYear.TreatyYear` b8035) dan
 * `UNDERWRITINGYEAR` berlabel `Transaction Year` (b8284 → `.UnderwritingYear`
 * b8315). Pasangan judul↔sel grid dibaca menurut URUTAN di dalam satu grid
 * (`[dugaan kuat]`). Merapikannya adalah keputusan Product + UW, bukan kami.
 */
export const TAHUN_TCO = {
  /** `Section/GridTreatyContract.xml` b1057 `<pyValue>` — judul layar. */
  judul: 'TREATY CONTRACT OUT',
  /** `Section/InputDtlTreatyContact.xml` b5437 `<pyValue>` — judul form. */
  inputNewData: 'Input New Data',

  // --- kepala kolom grid, `Section/InputTreatyContract.xml` ---
  /** b17378 → sel `.UnderwritingYear` b18964. */
  kolomUnderwritingYear: 'Underwriting Year',
  /** b17531 → sel `.TreatyYear` b19125. */
  kolomTransactionYear: 'Transaction Year',
  /** b17684 → sel `.StartDate` b19286. */
  kolomStartDate: 'StartDate',
  /** b17837 → sel `.EndDate` b19446. */
  kolomEndDate: 'EndDate',
  /** b17990 → sel `.TreatyGroupName` b19605. */
  kolomTreatyGroup: 'Treaty Group',
  /** b18136 → sel `.Proportion` b19724 — kolom `PROPORTION` menyimpan pilihan jenis reasuransi. */
  kolomReinsuranceType: 'Reinsurance Type',

  // --- tombol grid, `Section/InputTreatyContract.xml` ---
  /** b16387 `<pyLabel>` → `NewInputTreatyYear_Act`. */
  add: 'Add',
  /** b19939 `<pyLabel>` → `SetTreatyYear_Act` b20030. */
  edit: 'Edit',
  /** b20778 `<pyLabel>` — membuka konteks kontrak (tiket 04). */
  reinsType: 'ReinsType',
  /** b22196 `<pyLabel>` — membuka konteks klausul (tiket 08). */
  listDescription: 'List Description',

  // --- form, `Section/InputDtlTreatyContact.xml` ---
  /** b6379 `<pyLabelFieldValue>` — `InputTreatyYear.ID`, hanya dibaca. */
  formId: 'ID',
  /** b6560 — `InputTreatyYear.TreatyGroupName` (pemilih `BrowseTreatyGroup_RD` b6624). */
  formTreatyGroup: 'Treaty Group',
  /** b6800 — `InputTreatyYear.Proportion` b6829 (`.ID` pilihan). */
  formReinsuranceType: 'Reinsurance Type',
  /** b7532 — `InputTreatyYear.StartDate`. */
  formStartDate: 'Start Date',
  /** b7816 — `InputTreatyYear.EndDate`. */
  formEndDate: 'End Date',
  /** b8004 — `InputTreatyYear.TreatyYear` b8035 (lihat OQ-TCO-05). */
  formUnderwritingYear: 'Underwriting Year',
  /** b8284 — `InputTreatyYear.UnderwritingYear` b8315 (lihat OQ-TCO-05). */
  formTransactionYear: 'Transaction Year',
  /** b9097 — `InputTreatyYear.TglUpdate`, hanya dibaca. */
  formModifiedDate: 'Modified Date',
  /** b9282 — `InputTreatyYear.UserID`, hanya dibaca. */
  formUsername: 'Username',
  /** b10332 `<pyLabel>` → `SaveTreatyYear_Act` b10356. */
  save: 'Save',
  /** b10622 `<pyLabel>` → `CancelActivityTreatyContract` b10640. */
  cancel: 'Cancel',

  /** `[tidak ada di korpus]` — kosakata kami. */
  kosong: 'Belum ada tahun treaty.',
  /** `[tidak ada di korpus]` — tombol yang menunggu tiketnya. */
  menungguTiket: 'menunggu tiket',
  /** `[tidak ada di korpus]` — catatan OQ-TCO-05 di layar. */
  catatanLabelBersilang:
    'Label tahun di grid dan di form bersilang di sistem lama (OQ-TCO-05); keduanya dibawa apa adanya.',
} as const
