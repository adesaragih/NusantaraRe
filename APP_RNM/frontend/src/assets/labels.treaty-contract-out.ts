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
 * Menu kelompok **Treaty Contract Out** — tco5 `[DIPUTUSKAN work owner
 * 29-09-2026]`: *"untuk menu hanya Treaty Contract Out; Treaty Contract
 * ReinsType dan Treaty Contract Description dihapus"*. SATU butir.
 *
 * XML membenarkannya: `InboxTreatyContractReinsType` dan
 * `InboxTreatyContractDescription` BUKAN menu portal melainkan popup dari form
 * kontrak — `Section/InputTreatyContract.xml` tombol `ReinsType` b20778 →
 * harness b20947/b21627, tombol `List Description` b22196 → harness
 * b22323/b23088. Keduanya dibuka tombol itu di layar tahun treaty.
 *
 * ⚠️ Nama kelompok = nama FOLDER korpus, pola `MODUL` di `labels.ts`
 * (`MODUL.treatyContractOut`, ditambahkan ADITIF di tiket 03).
 */
export const MENU_TCO = {
  /** Nama folder korpus `D:\XML\RNM_BRD\Treaty Contract Out`. */
  kelompok: 'Treaty Contract Out',
  /** tco5: label SATU-SATUNYA butir menu `[keputusan work owner]`; asalnya harness portal `InboxTreatyContract`. */
  treatyContractOut: 'Treaty Contract Out',
  /** `Harness/InboxTreatyContract.xml` b151 `<pyLabel>` — nama harness asal butir itu (bukan label menu sejak tco5). */
  inboxTreatyContract: 'InboxTreatyContract',
  /** `Harness/InboxTreatyContractReinsType.xml` b151 `<pyLabel>` — popup tombol `ReinsType` b20778, bukan menu. */
  inboxTreatyContractReinsType: 'InboxTreatyContractReinsType',
  /** `Harness/InboxTreatyContractDescription.xml` b359 `<pyLabel>` — popup tombol `List Description` b22196, bukan menu. */
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

/**
 * Tiket 12 — panel lampiran tahun treaty (FITUR BARU, penyimpangan sadar 9).
 *
 * Nomor baris = `Section/GridTreatyArrangementAttachment.xml` kecuali disebut
 * lain; panel itu disertakan `Section/InputTreatyContract.xml` b13074.
 *
 * ⛔ `Download` b2391 (→ `DownloadAll_Act` generik) TIDAK dibawa sebagai
 * tombol kedua: `Download All` b2659 (→ `TreatyOutDownloadAll_Act`) adalah
 * aksi yang sama untuk modul ini, dan dua tombol unduh-semua berdampingan
 * hanya membingungkan.
 */
export const LAMPIRAN_TCO = {
  /** `Section/InputTreatyContract.xml` b11721 `<pyValue>`. */
  attachmentFor: 'Attachment for',
  /** b1785 `<pyValue>`. */
  forTreatyContractOut: 'For Treaty Contract Out',
  /** b578 `<pyLabel>` → `SetCategory_act` b596 → flow action `TreatyOutAttachContent` b642. */
  addAttachment: 'Add attachment',
  /** b1023 `<pyLabel>` → `LoadAttachmentTreatyOut` b1041. */
  refresh: 'Refresh',
  /** b2659 `<pyLabel>` → `TreatyOutDownloadAll_Act` b2677. */
  downloadAll: 'Download All',
  /** b3032 `<pyValue>` — judul kolom; sel `.pyFileName` b3428 → `TreatyOutDownloadOne` b3488. */
  kolomFileName: 'File Name',
  /** b3170 `<pyValue>` — judul kolom; sel `.pyCategory` b3705. Juga label pemilih kategori. */
  kolomType: 'Type',
  /** b3897 `<pyLabel>` → `DeleteAttachmentTreaty` b3915. */
  delete: 'Delete',
  /** `Activity/TreatyOutSaveAttachment.xml` b376 `Local.Err` — VERBATIM. */
  tanpaBerkas: 'Tidak ada file yg diattach',

  /** `[tidak ada di korpus]` — kosakata kami (fiturnya tidak ada di Pega). */
  kolomStatus: 'Status',
  statusTerkirim: 'terkirim',
  statusTertunda: 'tertunda',
  statusGagal: 'gagal',
  ulangi: 'Ulangi',
  periksaSelaras: 'Periksa keselarasan',
  selarasBersih: 'Rekam lampiran dan berkas di penyimpanan sejalan.',
  kosong: 'Belum ada lampiran pada tahun treaty ini.',
  simpanDulu: 'Simpan tahun treaty lebih dulu; lampiran melekat pada tahun treaty yang sudah ber-ID.',
  pilihKategori: '— pilih —',
} as const

/**
 * Tiket 04 — editor kontrak treaty di dalam tahun treaty.
 *
 * Nomor baris = `Section/InputTreatyContractReinsType.xml` (disertakan
 * `PanggilReinsType.xml` b1247, disertakan harness
 * `InboxTreatyContractReinsType.xml` b1799) kecuali disebut lain.
 */
export const KONTRAK_TCO = {
  /** `Harness/InboxTreatyContractReinsType.xml` b1670 `<pyValue>` — judul. */
  judul: 'ReinsType',
  /** b1145 — `InputTreatyContractReinsType.UnderwritingYear` b1176, hanya dibaca. */
  headerUnderwritingYear: 'Underwriting Year',
  /**
   * b1358 — label `ReinsType` pada medan `InputTreatyContractReinsType.TreatyGroupName`
   * b1386. ⚠️ Label dan nilainya bersilang di korpus (nama GRUP berlabel ReinsType);
   * dibawa apa adanya, catatannya tampil di layar.
   */
  headerReinsType: 'ReinsType',
  /** b2905 — `InputData.CARIDATETIME`; perubahan → `SetTanggalTreatyContract` b3007. */
  formStartDate: 'Start Date',
  /** b3244 — `InputData.CARIENDDATE`. */
  formEndDate: 'End Date',
  /** b3618 `<pyLabel>` → `SaveTreatyContract_Act` b3642. */
  save: 'Save',
  /** b4363 — `InputTreatyContract.TglUpdate`, hanya dibaca. */
  formModifiedDate: 'Modified Date',
  /** b4547 — hanya dibaca. */
  formUsername: 'Username',
  /** b5343 `<pyLabel>` → `UndoOperation` b5366. */
  undo: 'Undo',
  /** b6400 — `OutputData.HASIL1`. */
  information: 'Information',
  /** b8528 `<pyLabel>` → `NewInputTreatyContract_Act` b8552. */
  add: 'Add',
  /** b9164 `<pyValue>` — sel `.ReinsTypeName` b9992. */
  kolomReinsType: 'Reins Type',
  /** b9304 `<pyValue>` — sel `.TreatyStartDate` b10100. */
  kolomTreatyStart: 'Treaty Start',
  /** b9444 `<pyValue>` — sel `.TreatyEndDate` b10285. */
  kolomTreatyEnd: 'Treaty End',
  /** b10519 `<pyLabel>` → `SetUbahTreatyContract` b10543. */
  edit: 'Edit',
  /** b10842 `<pyLabel>` → `BrowseTreatyBusinessList_Act` — tiket 07. */
  businessList: 'Business List',
  /** b11308 `<pyLabel>` → `BrowseTreatyReinsurerList_Act` — tiket 05. */
  reinsurerList: 'Reinsurer List',
  /** b11809 `<pyLabel>` → `BrowseDeleteRowTreatyInContract` — tiket 10. */
  delete: 'Delete',

  /** `[tidak ada di korpus]` — label medan ID; korpus hanya punya label bawaan kontrol `Formatted Text` b2478. */
  formId: 'ID',
  /** `[tidak ada di korpus]` — pemilih tahun saat layar dibuka dari menu (Pega membukanya sebagai popup berkonteks). */
  pilihTahun: 'Tahun treaty',
  /** `[tidak ada di korpus]` */
  pilihTahunDulu: 'Pilih tahun treaty lebih dulu, atau buka lewat tombol ReinsType di layar InboxTreatyContract.',
  /** `[tidak ada di korpus]` */
  kosong: 'Belum ada kontrak pada tahun treaty ini.',
  /** `[tidak ada di korpus]` */
  tutup: 'Tutup',
  /** `[tidak ada di korpus]` */
  tersimpan: 'Kontrak tersimpan.',
  /** `[tidak ada di korpus]` — catatan label bersilang b1358/b1386. */
  catatanLabelBersilang: 'Label ReinsType di kepala layar memuat nama grup treaty di sistem lama; dibawa apa adanya.',
} as const

/**
 * Tiket 05 — panel reinsurer pada kombinasi (tahun, grup, jenis).
 *
 * Nomor baris = `Section/ViewDetailTreatyReinsurerGrid1.xml` (disertakan
 * `InputTreatyContractReinsType.xml` b13311, dibuka tombol baris kontrak
 * `Reinsurer List` b11308 → `BrowseTreatyReinsurerList_Act` b11325).
 *
 * ⛔ `Tambah` b1996 (tombol kedua ke `NewTreatyReinsurerDetail_Act` yang sama
 * dengan `Add` b1730) TIDAK dibawa sebagai tombol kedua.
 */
export const REINSURER_TCO = {
  /** b1730 `<pyLabel>` → `NewTreatyReinsurerDetail_Act` b1754. */
  add: 'Add',
  /** b2418 `<pyValue>` — sel `.ReinsurerID`. */
  kolomReinsId: 'ReinsID',
  /** b2560 `<pyValue>` — sel `.NAME`. */
  kolomReinsurer: 'Reinsurer',
  /** b2702 `<pyValue>` — sel `.PctShare`. */
  kolomShare: '%Share',
  /** b2844 `<pyValue>` — sel `.Ricomm`. */
  kolomComm: '%Comm',
  /** b2990 `<pyValue>` — sel `.StdRating`. */
  kolomRating: 'Rating',
  /** b3138 `<pyValue>` — sel `.OperatorName`. */
  kolomOperatorName: 'Operator Name',
  /** b4491 `<pyLabel>` → `SetUbahTreatyReinsurerList_Act` b4543. */
  edit: 'Edit',
  /** b4936 `<pyLabel>` → `DeleteTreatyReins_Act` b4953 — tiket 10. */
  delete: 'Delete',
  /** b5277 `<pyLabel>` — tiket 06. */
  securityReinsurer: 'Security Reinsurer',
  /** b6186 `<pyValue>` (`--&gt;&gt;` di XML) — nilai `InputTreatyReinsurer.TotalShare` b6330. */
  totalShare: 'Total Share -->>',
  /** b7842. */
  formId: 'ID',
  /** b8042 — `InputTreatyReinsurer.ReinsurerID`, diisi pemilih. */
  formReinsId: 'Reins.ID',
  /** b8226 — pemilih `BrowseAgentReinsSOA_RD`, tampil `.ClientName`. */
  formReinsurer: 'Reinsurer',
  /** b8522 — `InputTreatyReinsurer.PctShare` (koma atau titik desimal). */
  formShare: '%Share',
  /** b8800 — `InputTreatyReinsurer.Ricomm`. */
  formComm: '%Comm',
  /** b9076 — `InputTreatyReinsurer.StdRating`. */
  formRating: 'Rating',
  /** b11100 — hanya dibaca. */
  formOperatorName: 'Operator Name',
  /** b11405 `<pyLabel>` → `SaveTreatyReinsurerDetail1_Act` b11429. */
  save: 'Save',
  /** b12131 — `OutputParam.ERRMSG6`. */
  error: 'Error',
  /** b12868 — `OutputParam.ERRMSG`. */
  informasi: 'Informasi',

  /** `[tidak ada di korpus]` */
  cariReinsurer: 'Cari nama reinsurer',
  /** `[tidak ada di korpus]` */
  kosong: 'Belum ada reinsurer pada kombinasi ini.',
  /** `[tidak ada di korpus]` */
  tersimpan: 'Reinsurer tersimpan.',
  /** `[tidak ada di korpus]` */
  tutup: 'Tutup',
  /** `[tidak ada di korpus]` — kepala kombinasi (`OutputData.HASIL1/3/2` b2284–b2296). */
  kombinasi: 'Kombinasi',
} as const

/**
 * Tiket 07 — panel business pada kombinasi (tahun, grup, jenis).
 *
 * Nomor baris = `Section/ViewDetailTreatyBusinessGrid.xml` (disertakan
 * `InputTreatyContractReinsType.xml` b14064, dibuka tombol baris kontrak
 * `Business List` b10842 → `BrowseTreatyBusinessList_Act` b10859).
 */
export const BUSINESS_TCO = {
  /** b1988 — kepala panel; nilai `InputTreatyBusiness.ReinsTypeName` b2019. */
  judul: 'Business List',
  /** b2785 `<pyLabel>` → `NewTreatyBusinessDetail_Act` b2809. */
  add: 'Add',
  /** b3422 `<pyValue>` — sel `.TreatyGroupName` b4070. */
  kolomTreatyGroup: 'Treaty Group',
  /** b3566 `<pyValue>` — sel `.ID` b4228 (ID BARIS, bukan kode bisnis). */
  kolomBusinessId: 'Business ID',
  /** b3710 `<pyValue>` — sel `.BIZNAME` b4380. */
  kolomBusinessName: 'Business Name',
  /** b4547 `<pyLabel>` → `SetUbahTreatyBusinessList_Act` b4571. */
  edit: 'Edit',
  /** b4826 `<pyLabel>` → `DeleteRowBusiness` b4850. */
  delete: 'Delete',
  /** b6241 — pemilih `BrowseFilterBusiness_RD`, tampil `.Note`. */
  formBusinessName: 'Business Name',
  /** b6499 — radio `InputTreatyBusiness.IsActive`, wajib. Juga kolom status grid (AC 22). */
  formActive: 'Active',
  /** b6680 — `InputTreatyBusiness.BizCode` (tersembunyi di Pega; tampil baca-saja di sini). */
  formBusinessCode: 'Business Code',
  /** b6966 `<pyLabel>` → `SaveTreatyBusinessDetail_Act` b6993. */
  save: 'Save',
  /** b9319 / b10064 — `OutputData.HASIL5` / `OutputParam.ERRMSG4` b10095. */
  information: 'Information',
  /** b10889 `<pyLabel>` → `CancelActivity` b10914. */
  closeList: 'Close List',

  /** `[tidak ada di korpus]` — nilai radio Active; `0` = nonaktif [keputusan work owner 29-09-2026] (OQ-TCO-13). */
  aktif: 'Aktif',
  nonaktif: 'Nonaktif',
  /** `[tidak ada di korpus]` */
  kosong: 'Belum ada bisnis pada kombinasi ini.',
  /** `[tidak ada di korpus]` */
  tersimpan: 'Bisnis tersimpan.',
} as const

/**
 * Tiket 08 — layar klausul (harness `InboxTreatyContractDescription` b359).
 *
 * Nomor baris = `Harness/InboxTreatyContractDescription.xml` kecuali disebut lain.
 */
export const KLAUSUL_TCO = {
  /** b2003 — `InputTreatyArrangementDesc.TreatyGroupID`. */
  headerTreatyGroupId: 'TreatyGroupID',
  /** b2489 — `.TreatyYear` (label bersilang, lihat OQ-TCO-05). */
  headerUnderwritingYear: 'Underwriting Year',
  /** b2665 — `.UnderwritingYear`. */
  headerTransactionYear: 'Transaction Year',
  /** b2839. */
  headerStartDate: 'Start Date',
  /** b3025. */
  headerEndDate: 'End Date',
  /** b3209 — `.TreatyGroupName`. */
  headerTreatyDescription: 'Treaty Description',
  /** b3382 — `.Proportion`. */
  headerProportionType: 'Proportion Type',
  /** b4880 `<pyValue>` — grid `BrowseTreatyDesc_RD` `IsXOL = 0` (b5335). */
  gridNonXol: 'For Non XOL',
  /** b7971 `<pyValue>` — grid `IsXOL = 1` (b8426). */
  gridXol: 'For XOL',
  /** b5440 `<pyValue>` — sel `.ID`. */
  kolomId: 'ID',
  /** b5549 `<pyValue>` — sel `.DescName`. */
  kolomDescriptionName: 'Description Name',
  /** b6059 `<pyLabel>` → `BrowseDescriptionLimit` + `testingKurs` + `SetKirimIDDesc` + `PanggilID` + `RefreshErrorProportionalarrg`. */
  show: 'Show',
  /** `Section/GridTreatyArrangementEpi.xml` b4487. */
  formModifiedDate: 'Modified Date',
  /** `GridTreatyArrangementEpi.xml` b5501. */
  save: 'Save',
  /** `GridTreatyArrangementEpi.xml` b8980. */
  add: 'Add',
  /** `GridTreatyArrangementEpi.xml` b10917. */
  edit: 'Edit',
  /** `GridTreatyArrangementEpi.xml` b11200 → `BrowseTreatyArrEpiParentList`. */
  showChild: 'Show Child',
  /** `Section/GridTreatyArrTreatyEpiList.xml` b8657 → `CancelActivityTreatyLimitChild`. */
  closeChild: 'Close Child',
  /** `Section/GridTreatyArrangementExclutionTreaty.xml` b955. */
  exclusionTreaty: 'Exclusion Treaty',

  /** `[tidak ada di korpus]` */
  cancel: 'Cancel',
  /** `[tidak ada di korpus]` */
  tutup: 'Tutup',
  /** `[tidak ada di korpus]` */
  kosong: 'Belum ada baris klausul.',
  /** `[tidak ada di korpus]` */
  totalPct: 'Total Pct',
  /** `[tidak ada di korpus]` */
  tersimpan: 'Klausul tersimpan.',
  /** `[tidak ada di korpus]` */
  cariPilihan: 'Cari',
  /** `[tidak ada di korpus]` */
  pilihTahun: 'Tahun treaty',
  /** `[tidak ada di korpus]` */
  pilihTahunDulu: 'Pilih tahun treaty lebih dulu, atau buka lewat tombol List Description di layar InboxTreatyContract.',
  /** `[tidak ada di korpus]` — Rp/Usd anak dihitung server (HitungRpUsd). */
  turunanServer: 'Rp dan Usd baris anak dihitung dari induknya (Pct × nilai induk ÷ 100).',
} as const

/**
 * Label medan klausul — VERBATIM dari form tiap jenis (tiket 08).
 *
 * Bawaan per medan; `LABEL_MEDAN_KHUSUS` menimpa per jenis/subjenis. ⚠️ Form
 * exclusion memakai label rujukan properti `.Occupation` b2296 / `.Clause`
 * b2320 — yang tampil adalah nama propertinya.
 */
export const LABEL_MEDAN_KLAUSUL = {
  /** `GridTreatyArrangementEpi.xml` b2777. */
  ReinsTypeID: 'ReinsType',
  /** b3215. */
  Line: 'Line',
  /** b3372. */
  Rp: 'Rp',
  /** b3798. */
  Usd: 'Usd',
  /** `GridTreatyArrTreatyEpiList.xml` b3029. */
  Pct: 'Pct',
  /** `GridTreatyArrangementProfitCommision.xml` b3223. */
  PctMe: 'PctMe',
  /** b3481. */
  Ydcf: 'Ydcf',
  /** `GridTreatyArrangementBordereAux.xml` b2702. */
  Method: 'Method',
  /** `GridTreatyArrangementTerrLimit.xml` b2742. */
  TerritorialLimit: 'Territorial Limit',
  /** `GridTreatyArrangementCoins.xml` b3394. */
  CoIns_Min: 'From',
  /** b4193. */
  CoIns_Max: 'To',
  /** b4880. */
  TreatyLimit: 'Treaty Limit',
  /** `GridTreatyArrangementExclutionTreatyOccupation.xml` b2007. */
  ID_Occupation: 'ID Occupation',
  /** b2296 (`.Occupation`). */
  Occupation: 'Occupation',
  /** `GridTreatyArrangementExclutionTreatyClausule.xml` b2030. */
  ID_Clause: 'ID Clause',
  /** b2320 (`.Clause`). */
  Clause: 'Clause',
  /** `GridTreatyArrangementExclutionTreatyPeriode.xml` b500. */
  Layer: 'Max Periode (Month)',
} as const satisfies Readonly<Record<string, string>>

/** Penimpaan label per `jenis` atau `jenis/subjenis`. */
export const LABEL_MEDAN_KHUSUS = {
  /** `GridTreatyArrangementMinLOL.xml` b1540. */
  MinLOL: { Pct: 'Minimum LOL (%)' },
  /** `GridTreatyArrangementMInLOLMB.xml` b1548. */
  MinLOLMB: { Pct: 'Minimum LOL MB (%)' },
  /** `GridTreatyArrangementMaxCoinsPanel.xml` b1525. */
  MaxCoinsPanel: { CoIns_Max: 'Max Coins Panel' },
  /** `GridTreatyArrangementExclutionTreatyOccupation.xml` b2745 / b2935 / b3222. */
  'ExclutionTreaty/Occupation': { Line: 'Class of Contruction', Usd: 'TSI Less Than (USD)', Rp: 'TSI Less Than (IDR)' },
  /** `GridTreatyArrangementExclutionTreatyObject.xml` b566 (`TSI BI &gt;`). */
  'ExclutionTreaty/Object': { Pct: 'TSI BI >' },
} as const satisfies Readonly<Record<string, Readonly<Record<string, string>>>>

/**
 * Tiket 06 — grid security di bawah reinsurer (`InputTreatyContractReinsType.xml`
 * b14885, tampil bila `HASILD21 == 1`; dibuka tombol baris reinsurer
 * `Security Reinsurer`, `ViewDetailTreatyReinsurerGrid1.xml` b5277).
 *
 * Nomor baris = `Section/InputTreatyContractReinsType.xml`.
 */
export const SECURITY_TCO = {
  /** b15459 `<pyLabel>` → `InputNewSecurityReinsurer` (b15487). */
  add: 'Add',
  /** b16088 `<pyValue>` — sel `.REAS_SECURITY` b16724. */
  kolomReasSecurity: 'Reas Security',
  /** b16228 `<pyValue>` — sel `.CLIENTNAME` b16878. */
  kolomSecurityName: 'Security Name',
  /** b16368 `<pyValue>` — sel `.PCT_SHARE` b17013. */
  kolomPercentShare: 'Percent Share',
  /** b17252 `<pyLabel>` → `ShowEditSecurityReinsurer` b17276. */
  edit: 'Edit',
  /** b17559 `<pyLabel>` → `DeleteSecurityReinsurer` b17583 (tanpa konfirmasi, aksi `refresh`). */
  delete: 'Delete',
  /** b19468 — `InputTreatySecurity.REAS_SECURITY` b19499, selalu nonaktif b19515. */
  formSecurityId: 'Security ID',
  /** b19648 — wajib b19642; pemilih `BrowseAgentReinsSOA_RD` b19711. */
  formSecurityName: 'Security Name',
  /** b19888 — `InputTreatySecurity.PCT_SHARE` b19919. */
  formShare: '%Share',
  /** b20246 `<pyLabel>` → `SaveSecurityReinsurer_Act` b20270. */
  save: 'Save',
  /** b20980. */
  error: 'Error',
  /** b21717. */
  informasi: 'Informasi',

  /** `[tidak ada di korpus]` */
  judul: 'Security',
  /** `[tidak ada di korpus]` */
  cariSecurity: 'Cari security',
  /** `[tidak ada di korpus]` */
  kosong: 'Belum ada security pada reinsurer ini.',
  /** `[tidak ada di korpus]` */
  tersimpan: 'Security tersimpan.',
  /** `[tidak ada di korpus]` */
  cancel: 'Cancel',
  /** `[tidak ada di korpus]` */
  tutup: 'Tutup',
} as const

/**
 * Tiket 11 — kurs USD → IDR di layar klausul.
 *
 * ⚠️ Korpus tidak punya label tampil untuk kurs: section `NitipKurs`
 * (`Harness/InboxTreatyContractDescription.xml` b3882) hanya menitip
 * `InputTreatyArrangement.Kurs` (`Section/NitipKurs.xml` b512). Pesan kurs
 * kosong datang VERBATIM dari server (`NewTreatyArrEpi.xml` b870).
 */
export const KURS_TCO = {
  /** `[tidak ada di korpus]` */
  kurs: 'Kurs USD → IDR',
  /** `[tidak ada di korpus]` */
  berlaku: 'berlaku',
  /** `[tidak ada di korpus]` */
  sampai: 's.d.',
  /** `[tidak ada di korpus]` — `HitungRpUsd_depan`: Usd = Rp ÷ Kurs, dihitung server. */
  catatanRpKeUsd: 'Usd dihitung server dari Rp ÷ Kurs.',
  /** `[tidak ada di korpus]` — `CalculateTSIExcludeTreaty`. */
  catatanDuaArah: 'Mengisi IDR menghitung USD (Rp ÷ Kurs), mengisi USD menghitung IDR (Usd × Kurs) — di server.',
  /** `[tidak ada di korpus]` — keputusan work owner 29-09-2026: baris kembar = satu kurs. */
  catatanKembar: 'baris kembar (TOIDR sama) di master kurs — dipakai sebagai satu kurs',
  /** `[tidak ada di korpus]` — lanjutan 6: baris yang tanggalnya ditolak Oracle dilewati. */
  catatanDitolak: 'baris master kurs dilewati karena tanggalnya ditolak Oracle',
} as const

/**
 * Tiket 10 — popup konfirmasi hapus (penyimpangan sadar 4).
 *
 * ⚠️ Pega menghapus TANPA konfirmasi (`Delete` kontrak `InputTreatyContractReinsType.xml`
 * b11809, `Delete` reinsurer `ViewDetailTreatyReinsurerGrid1.xml` b4936). Seluruh
 * teks popup karena itu `[tidak ada di korpus]`, kecuali pesan sukses kontrak yang
 * datang VERBATIM dari server (`BrowseDeleteRowTreatyInContract.xml` b762).
 */
export const HAPUS_TCO = {
  /** `[tidak ada di korpus]` */
  judulKontrak: 'Hapus kontrak?',
  /** `[tidak ada di korpus]` */
  judulReinsurer: 'Hapus reinsurer?',
  /** `[tidak ada di korpus]` */
  ya: 'Ya',
  /** `[tidak ada di korpus]` */
  batal: 'Batal',
  /** `[tidak ada di korpus]` */
  ikutTerhapus: 'Ikut terhapus:',
  /** `[tidak ada di korpus]` */
  reinsurer: 'reinsurer',
  /** `[tidak ada di korpus]` */
  security: 'security',
  /** `[tidak ada di korpus]` */
  business: 'business',
  /** `[tidak ada di korpus]` — AC 44: klausul milik tahun/grup/jenis, bukan milik satu kontrak. */
  klausulTetap: 'baris klausul TIDAK ikut terhapus — klausul milik tahun/grup/jenis reasuransi, bukan milik satu kontrak.',
  /** `[tidak ada di korpus]` */
  memuatDampak: 'Menghitung baris yang akan ikut terhapus…',
  /** `[tidak ada di korpus]` — OQ-TCO-21 [keputusan work owner 29-09-2026]: hapus seperti Pega, tetapi tidak diam. */
  bersama: 'kontrak lain memakai kombinasi yang sama — reinsurer dan security-nya IKUT terhapus bersama kontrak ini (seperti Pega); business milik tahun treaty lain tidak ikut.',
} as const
