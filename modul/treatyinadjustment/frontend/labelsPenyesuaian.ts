// Label dan spesifikasi LAYAR ADJUSTMENT — `Section/InputTreatyInAdjustment.xml`
// beserta anak-anaknya.
//
// ⛔ SETIAP label, kolom, lebar, dan syarat di bawah DISALIN dari ekspor
// `D:\XML_NURE\Treaty In Adjustment` (2026-09). Offset `@n` adalah posisi bita
// di berkas Section-nya SESUDAH seluruh `<pyIncludedRuleXML>` dibuang dengan
// menghitung sarangnya — lihat `docs/LAYAR-ADJUSTMENT.md` §1. Yang tidak ada di
// ekspor TIDAK ada di sini.

/** Golongan nilai satu sel. `tanggal` diformat `formatDate` milik inti. */
export type JenisAngka = 'uang' | 'persen' | 'persenShare' | 'tanggal' | 'teks'

/** Desimal MAKSIMUM — nol di ekor dibuang `formatNumber`. */
export const DESIMAL_UANG = 4
export const DESIMAL_PERSEN = 2
/** Persen SHARE — 8, sama dengan batas simpan `NUMBER(38,8)`. */
export const DESIMAL_PERSEN_SHARE = 8

export const PENYESUAIAN = {
  judulMenu: 'Treaty In Adjustment — Penyesuaian',
  /** pyValue @26602 — judul layar. */
  judul: 'Adjustment Treaty In',
  // Kepala mode detail — pyLabelFieldValue tiap sel.
  idAsal: 'ID Original', // .OLDID @57593
  id: 'ID Revision', // .ID @63903
  jenisReasuransi: 'Reinsurance Type', // .ProportionType @69930 (pxRadioButtons)
  jenisPenyesuaian: 'Adjustment Type', // .EDMState @76962 (pxDropdown, tampil bila EDMState != 3 @79936)
  /** pyValue @82379 — LABEL pengganti dropdown bila EDMState = 3 (@84311). */
  premiPenyesuaian: 'Adjustment Premium',
  jenisMaterial: 'Material Type', // .EDMMaterialType @87067 (pxDropdown, ALWAYS)
  proporsional: 'Proportional',
  nonProporsional: 'NonProportional',

  // Panel mode detail.
  panelLama: 'Old Data', // pyTitle @192816 → TreatyInNONProportionalOldData
  panelBaru: 'New Data', // pyTitle @233311 → TreatyInNONProportional

  /**
   * ⚠️ Teks pilihan dropdown `EDMState`/`EDMMaterialType` TIDAK ada di
   * ekspor: daftarnya `pyListSource=associated`, milik rule properti yang
   * tidak ikut diekspor. Yang tampil KODENYA, ditandai.
   */
  kodeBelumBerteks: '',

  // Mode daftar (`DATASHOW != 1` @481804).
  tambahRevisi: 'Add Revision', // pyLabel @500554 → PickerTreatyInMasterRevisi
  tambahPremi: 'Add Adjustment Premium', // pyLabel @515429 → PickerTreatyInMaster

  // Picker tombol Add — `Section/PickerTreatyInMasterRevisi.xml` /
  // `Section/PickerTreatyInMaster.xml`.
  pickerRevisiInternal: 'Choose Master to Create Internal Revision', // @39970, EDMState = '1'
  pickerRevisiEksternal: 'Choose Master to Create External Revision', // @44930, EDMState = '2'
  pickerPremi: 'Choose Master to create Premium Adjustment', // @39472
  internalEksternal: 'Internal / External', // .EDMState @86844 (pxRadioButtons)
  materialRadio: 'Material Type', // .EDMMaterialType @108821 (pxRadioButtons)
  pilihMaster: 'Choose', // .pyTemplateInputBox @186747 / @122878
  petunjukMasterKosong: '',
  /**
   * ⛔ Di sistem lama `Choose` langsung MENYIMPAN (`TreatyInEDMSetValue` [9]
   * `SaveTreatyIn_EDM_Act`). Keputusan pemilik proses 7 Oktober 2026: draf
   * di layar, simpanannya menunggu tombol Save.
   */
  drafBelumTersimpan: 'Draf — belum tersimpan. Penyesuaian ini baru ada di layar; ia masuk basis data lewat tombol Save.',

  ubah: 'Edit', // pyLabel @782051
  lihat: 'View', // pyLabel @798870
  kembali: 'Kembali ke daftar',
  tanpaBaris: 'No items',
  petunjukDaftarKosong: '',
  // Panel Old / mode View: tombol ekspor yang tidak berlaku di sini.
  tombolTulisMati: 'Tombol ini tidak berlaku di panel atau mode ini — dimatikan, bukan dihilangkan.',
  rumusBelum: 'Rumus tombol ini belum disambungkan — tombolnya dimatikan, bukan dihilangkan.',
  /**
   * Rantai menyentuh langkah BERSYARAT yang belum dibangun (mis.
   * `TreatyInXOLAddSpreadingDetailActual` bila `EDMState = 3`, yang menulis
   * halaman `TreatyIn.ActualValue`). Seluruh rantai batal — nol perubahan.
   */
  langkahBelum: '',
  /**
   * Submit Achievement (`InsertToLogAchievement` → `POOLDATA.LOG_ACHIEVEMENT`)
   * — keputusan pemilik proses 7 Oktober 2026: tetap mati sampai tabel log
   * diputuskan bersama jalur Save. Angka Achievement sendiri tersimpan lewat Save.
   */
  achievementMenunggu:
    'Submit Achievement menulis log tersendiri; tabelnya menunggu keputusan bersama jalur Save. Angka Achievement tersimpan lewat Save.',
  /** `GenerateCSVTreaty` — nama berkas dari ekspor (`FSFileName`). */
  berkasAchievement: 'CSVAchievementTreatyIn.xlsx',
  /** Daftar master (mata uang, Treaty Group, Treaty Type) belum dimuat. */
  masterBelum: 'Daftar master belum dimuat; nama pasangannya tidak dapat diisi.',
  /** Nilai tersimpan yang tidak ada di daftar — ditawarkan, bukan dijatuhkan. */
  diLuarDaftar: '(di luar daftar)',

  // Medan bertanda dan keadaan.
  takAdaDiWarisan: '',
  tanpaTeks: '',
  belumDibangun: 'Tab ini belum dibangun.',
  belumDibangunPetunjuk: '',
  sebagian: '',
  petunjukGridKosong: '',
  petunjukHistory: '',
  // Grid `expandPane` yang flow action rinciannya TIDAK ada di ekspor —
  // hari ini nol; rincian lain dibangkitkan (`KERANGKA_RINCIAN`).
  rincianHilang: '',
  /** Tombol ▸/▾ rincian baris (`expandPane`). */
  rincianBuka: 'Rincian baris',
  includeTakAda: '',
  /** §17 — isi Retro tidak dibangun karena jarang. */
  retroJarang: 'Isi Retro tidak dibangun — jarang dipakai.',
  retroJarangPetunjuk: '',
} as const

/**
 * Grid daftar — `Section/InputTreatyInAdjustment.xml` @655629.
 *
 * ⛔ Urutan dan lebar DISALIN (`pyWidth` sel kepala). Lebar dipakai sebagai
 * PERBANDINGAN, bukan piksel — tata letak kita responsif, Pega tidak.
 */
/**
 * Teks tombol IKON — `pyLabel` dan `pyCaption` KOSONG di XML, hanya
 * `pyImage`. Pega yang berjalan mengisinya sendiri: `Tambah` terlihat di
 * tangkapan layar Co-Ins Scale Treaty In (pemakai, 6 Oktober 2026), `Hapus`
 * dari pola Pega yang SAMA di NB FacIn (tangkapan layar 3 Oktober 2026).
 */
export const TOMBOL_IKON: Readonly<Record<string, string>> = {
  'IconAdd.png': 'Tambah',
  'IconTrash.png': 'Hapus',
}

export const KOLOM_DAFTAR = [
  'ID Revision', // @663170 W136 .ID
  'ID Original', // @667220 W134 .OLDID
  'Type', // @671196 W93 .EDMState
  'Material Type', // @674776 W93 .EDMMaterialType
  'Contract Name', // @677990 W141 .TreatyContractName
  'Reinsurance Type', // @682129 W106 .ProportionType
  'Source of Business', // @686185 W137 .LeadingReinsSource
  'Ceding', // @690161 W135 .Ceding
  'Commencement', // @694125 W160 .Commencement
  'Termination', // @698096 W149 .Termination
  'Position To', // @702199 W83 .PositionUsername
  'Status Accept', // @706364 W72 .StatusAkseptasi
] as const
export const LEBAR_DAFTAR = [136, 134, 93, 93, 141, 106, 137, 135, 160, 149, 83, 72] as const
export const JENIS_DAFTAR: readonly JenisAngka[] = [
  'teks', 'teks', 'teks', 'teks', 'teks', 'teks', 'teks', 'teks', 'tanggal', 'tanggal', 'teks', 'teks',
]
/** Kolom tombol baris (Edit/View) — W56 @782773, tanpa judul. */
export const LEBAR_TOMBOL_BARIS = 56

/**
 * Grid picker — `TempMasterList.pxResults` (`CARI1`…`CARI7`). Kolom PERTAMA
 * tombol `Choose` tanpa judul (W109), lalu ketujuh kolom ini.
 */
export const KOLOM_PICKER = [
  'ID', // .CARI1
  'Contract Name', // .CARI2
  'Reinsurance Type', // .CARI3
  'Source of Business', // .CARI4
  'Ceding', // .CARI5
  'Commencement', // .CARI6 pxDateTime
  'Termination', // .CARI7 pxDateTime
] as const
export const JENIS_PICKER: readonly JenisAngka[] = ['teks', 'teks', 'teks', 'teks', 'teks', 'tanggal', 'tanggal']
/** `pyWidth` — tombol `Choose` lebih dulu. Kedua Section berbeda beberapa piksel. */
export const LEBAR_PICKER_REVISI = [109, 213, 215, 134, 215, 215, 218, 225] as const
export const LEBAR_PICKER_PREMI = [109, 208, 209, 132, 209, 212, 213, 217] as const

/**
 * Pilihan radio `Internal / External` (`TreatyIn.EDMState`).
 *
 * ⚠️ Daftarnya `associated` — milik rule properti yang TIDAK diekspor. Teks
 * `Internal`/`External` diambil dari Section yang SAMA: judul picker
 * `EDMState = '1'` berbunyi "…Internal Revision", `= '2'` "…External Revision".
 */
export const OPSI_JENIS_REVISI = [
  { nilai: '1', label: 'Internal' },
  { nilai: '2', label: 'External' },
] as const
/** Radio `Material Type` — kodenya saja; teks pilihannya tidak ada di ekspor (`kodeBelumBerteks`). */
export const OPSI_MATERIAL = ['1', '2'] as const

/**
 * Satu medan kepala form — `Section/TreatyInNONProportional[OldData].xml`.
 *
 * `dari`   halaman yang diikat. ⛔ Pada panel Old, `Contract Ref No`
 *          mengikat halaman AKAR (`TreatyIn.ContractRefNo` @40302), bukan
 *          `OLDDATA` — itu yang layar lama tampilkan, dan disalin.
 * `bentuk` kontrol di ekspor (`pyFormat`).
 * `syarat` cabang tempat medan tampil (`pyCondition` ProportionType).
 */
export interface MedanForm {
  label: string
  kunci: string
  dari: 'sisi' | 'akar'
  bentuk: 'teks' | 'area' | 'pilih' | 'tanggal' | 'tampil' | 'centang'
  syarat?: 'Proportional' | 'NonProportional'
  /**
   * Sisi New saja — kapan medan ini BACA-SAJA (`pyReadOnlyCondition`).
   * Dievaluasi atas halaman akar dengan `ViewState` sesi (Edit `0`, View
   * `1`). Tanpa ini dan tanpa `selaluBacaSaja`, medan dapat disunting.
   */
  bacaSajaBila?: (m: Readonly<Record<string, string>>) => boolean
  /**
   * Medan diisi lewat jendela pencarian agen (`localAction`), bukan diketik:
   * `reinsured` → Ceding + CedingID, `source` → LeadingReinsSource +
   * LeadingReinsSourceID (`DataTransform/TreatyInSetReinsured`). Tombolnya
   * tampil bila `TreatyIn.ViewState != 1`.
   */
  pilihAgen?: 'reinsured' | 'source'
  /** Sel `pyReadOnly=true` TANPA syarat — baca-saja di mode apa pun. */
  selaluBacaSaja?: boolean
  jejak: string
}

/** `TreatyIn.ViewState = 1` — mode View. */
const modeLihat = (m: Readonly<Record<string, string>>) => m.ViewState === '1'
/** `TreatyIn.IsEditData= 1`. */
const dataDiubah = (m: Readonly<Record<string, string>>) => m.IsEditData === '1'
/** `TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 1`. */
const dataDiubahAtauMaterial = (m: Readonly<Record<string, string>>) =>
  m.IsEditData === '1' || m.EDMMaterialType === '1'

/**
 * ⛔ Panel Old: NOL medan dapat disunting — keputusan, bukan salinan.
 * Ekspor membuat dua textarea Old dapat disunting (`pyReadOnly=false`
 * @43977 TeritorialScope, @56021 BordereauxNote); Old Data yang bisa diketik
 * adalah cara tercepat merusak jejak perubahan, jadi itu TIDAK ditiru.
 */
export const MEDAN_KIRI_LAMA: readonly MedanForm[] = [
  { label: 'Treaty Contract Name', kunci: 'TreatyContractName', dari: 'sisi', bentuk: 'teks', jejak: 'OldData @34625' },
  { label: 'Contract Ref No', kunci: 'ContractRefNo', dari: 'akar', bentuk: 'teks', jejak: 'OldData @40302 TreatyIn.ContractRefNo' },
  { label: 'Teritorial Scope', kunci: 'TeritorialScope', dari: 'sisi', bentuk: 'area', jejak: 'OldData @45884' },
  { label: 'Bordereaux', kunci: 'Bordeaux', dari: 'sisi', bentuk: 'pilih', syarat: 'Proportional', jejak: 'OldData @51628' },
  { label: 'Bordereaux Note', kunci: 'BordereauxNote', dari: 'sisi', bentuk: 'area', jejak: 'OldData @57926' },
]
export const MEDAN_KANAN_LAMA: readonly MedanForm[] = [
  { label: 'Commencement', kunci: 'Commencement', dari: 'sisi', bentuk: 'tanggal', jejak: 'OldData @78744' },
  { label: 'Termination', kunci: 'Termination', dari: 'sisi', bentuk: 'tanggal', jejak: 'OldData @88918' },
  { label: 'Treaty Year', kunci: 'TreatyYear', dari: 'sisi', bentuk: 'teks', jejak: 'OldData @95117' },
  { label: 'Accounting Mode', kunci: 'AccountingMode', dari: 'sisi', bentuk: 'pilih', syarat: 'Proportional', jejak: 'OldData @100988' },
  { label: 'Accounting Mode', kunci: 'AccountingModeNonProp', dari: 'sisi', bentuk: 'pilih', syarat: 'NonProportional', jejak: 'OldData @107318' },
  { label: 'Ceding', kunci: 'Ceding', dari: 'sisi', bentuk: 'tampil', jejak: 'OldData @113604' },
  { label: 'Source of Business', kunci: 'LeadingReinsSource', dari: 'sisi', bentuk: 'tampil', jejak: 'OldData @118843' },
]
export const MEDAN_KIRI_BARU: readonly MedanForm[] = [
  { label: 'Treaty Contract Name', kunci: 'TreatyContractName', dari: 'sisi', bentuk: 'teks', bacaSajaBila: dataDiubah, jejak: 'NONProportional @35853' },
  { label: 'Contract Ref No', kunci: 'ContractRefNo', dari: 'sisi', bentuk: 'teks', bacaSajaBila: dataDiubah, jejak: 'NONProportional @42427' },
  { label: 'Teritorial Scope', kunci: 'TeritorialScope', dari: 'sisi', bentuk: 'area', bacaSajaBila: dataDiubahAtauMaterial, jejak: 'NONProportional @48962' },
  { label: 'Bordereaux', kunci: 'Bordeaux', dari: 'sisi', bentuk: 'pilih', syarat: 'Proportional', bacaSajaBila: modeLihat, jejak: 'NONProportional @55076' },
  { label: 'Bordereaux Note', kunci: 'BordereauxNote', dari: 'sisi', bentuk: 'area', bacaSajaBila: dataDiubahAtauMaterial, jejak: 'NONProportional @61498' },
]
/**
 * ⛔ Kolom kanan New LEBIH PANJANG: `RNM as Treaty Leader` @189183,
 * `Effective Date` @238886, `Is Pro Rate` @249049 hanya ada di sisi ini.
 * Tiga sel `1=2` (Ceding / Business Source / Leading Reinsurer pxAutoComplete
 * @119100 @131001 @142956) TIDAK dibangun — salinan hidupnya `pxDisplayText`
 * @168493 dan @221597.
 */
export const MEDAN_KANAN_BARU: readonly MedanForm[] = [
  { label: 'Commencement', kunci: 'Commencement', dari: 'sisi', bentuk: 'tanggal', bacaSajaBila: modeLihat, jejak: 'NONProportional @83060' },
  { label: 'Termination', kunci: 'Termination', dari: 'sisi', bentuk: 'tanggal', bacaSajaBila: modeLihat, jejak: 'NONProportional @93504' },
  { label: 'Treaty Year', kunci: 'TreatyYear', dari: 'sisi', bentuk: 'teks', selaluBacaSaja: true, jejak: 'NONProportional @99854' },
  { label: 'Accounting Mode', kunci: 'AccountingMode', dari: 'sisi', bentuk: 'pilih', syarat: 'Proportional', bacaSajaBila: modeLihat, jejak: 'NONProportional @106135' },
  { label: 'Accounting Mode', kunci: 'AccountingModeNonProp', dari: 'sisi', bentuk: 'pilih', syarat: 'NonProportional', bacaSajaBila: modeLihat, jejak: 'NONProportional @112579' },
  { label: 'Ceding', kunci: 'Ceding', dari: 'sisi', bentuk: 'tampil', pilihAgen: 'reinsured', jejak: 'NONProportional @168493' },
  // Kotak centang menampilkan `pyCheckboxCaption`; label sel ketiga sel ini
  // (`Effective Date`) adalah sisa salin Pega.
  { label: 'RNM as Treaty Leader', kunci: 'TreatyLeader', dari: 'sisi', bentuk: 'centang', bacaSajaBila: modeLihat, jejak: 'NONProportional @189183' },
  { label: 'Source of Business', kunci: 'LeadingReinsSource', dari: 'sisi', bentuk: 'tampil', pilihAgen: 'source', jejak: 'NONProportional @221597' },
  { label: 'Effective Date', kunci: 'EDMEffective', dari: 'sisi', bentuk: 'tanggal', bacaSajaBila: modeLihat, jejak: 'NONProportional @238886 (TreatyMasterInEDM)' },
  { label: 'Is Pro Rate', kunci: 'IsProRate', dari: 'sisi', bentuk: 'centang', bacaSajaBila: modeLihat, jejak: 'NONProportional @249049 (TreatyMasterInEDM)' },
]

/**
 * Pilihan dropdown — NILAI TERSIMPAN, diukur atas 280 dokumen × 2 sisi
 * (5 Oktober 2026). Teks tampilnya tidak ada di ekspor; yang tersimpan
 * itulah yang ditawarkan. Diukur ulang `penyesuaian_db_test.go`.
 */
/**
 * Jendela pencarian agen — `Section/TreatyInSearchReinsured.xml` dan
 * `TreatyInSearchSoB.xml` (isi keduanya sama; yang berbeda `type` DataTransform
 * `TreatyInSetReinsured`). Judul jendela = label flow action-nya: ekspor tidak
 * memberi judul lain.
 */
/** Panel `Existing Policy for Master ID` — `Section/InputTreatyInAdjustment.xml` @111283. */
export const POLIS_MASTER = {
  judul: 'Existing Policy for Master ID', // pyTitle
  kolom: ['Policy No', 'Pega ID'] as const, // sel 31 · sel 32
  tanpaBaris: 'No items',
} as const

export const PILIH_AGEN = {
  kataKunci: 'Search', // sel 5 InputData.CARI2 — tanpa label; nama akses dari tombolnya
  cari: 'Search', // sel 8 pyLabel
  kolomId: 'ID', // sel 21
  kolomNama: 'Name', // sel 22
  pilih: 'Choose', // sel 27 pyLabel
  tanpaBaris: 'No items',
  judulCeding: 'TreatyInSearchReinsured', // FlowAction pyLabel
  judulSob: 'TreatyInSearchSoB', // FlowAction pyLabel
  tombolCeding: 'Choose Ceding', // TreatyInNONProportional sel 27
  tombolSob: 'Choose Source of Business', // TreatyInNONProportional sel 29
} as const

export const PILIHAN_MEDAN: Readonly<Record<string, readonly string[]>> = {
  Bordeaux: ['reporting', 'nonreporting'],
  AccountingMode: ['accounting', 'underwriting'],
  AccountingModeNonProp: ['loss', 'risk'],
}

/** Blok Pro Rate — tampil bila `TreatyIn.IsProRate == true` @264580. */
export const PRO_RATA = {
  judul: 'Pro Rate:', // pyValue @278917
  pemisah: '/', // pyValue @297724 — ProRateDays / ProRateTotalDays
  persen: 'Pro Rate Percentage', // pyLabelFieldValue @327367
  tandaPersen: '%', // pyValue @337597
  /** `pyDecimalPlaces` sel `ProRatePercent` @327220 — desimal per kolom dari ekspor. */
  desimal: 4,
} as const

/**
 * Grid Rate of Exchange — kolomnya dari kerangka bangkitan (`GRID_KURS`,
 * `TreatyInNONProportional[OldData].xml`).
 *
 * ⚠️ Sel `Currency` sisi New mengikat `.CurrencyID` (dropdown); sisi Old
 * `.Currency`. Diukur: `CurrencyID` adalah KODE (`10026`↔IDR 265 kali,
 * `10001`↔USD 240 kali, nol pasangan ganda — dijaga
 * `TestCurrencyIDBerpasanganTetapDenganCurrency`). Dropdown menampilkan NAMA
 * mata uangnya, jadi yang tampil di kedua sisi `.Currency`.
 */
export const KURS = {
  judul: 'Rate of Exchange', // pyTitle @157541 (Old) / @390565 (New)
  gantiKunciTampil: { CurrencyID: 'Currency' } as Readonly<Record<string, string>>,
} as const

/** Strip tab tiap (sisi × cabang) — urutan `pyTitle` tingkat tab. */
export const TAB_LAMA_NP = [
  'Maximum Retention', 'EGNPI', 'Limits', 'Share', 'Retro', 'Installment', 'Exclusions', 'Special Conditions',
] as const // TreatyInTabsNonProportionalOldData @23400 … @1798936
export const TAB_BARU_NP = [
  'Maximum Retention', 'Event Limits', 'EGNPI', 'Limits', 'Share', 'Retro', 'Installment',
  'Value Difference', 'Exclusions', 'Special Conditions', 'Information & Submit',
] as const // TreatyInTabsNonProportional @15962 … @3977799
/**
 * Cabang ADJUST PREMIUM panel New — `TreatyInNONProportional` @566980
 * menampilkan `TreatyInTabsNonProportionalAdjustPremi` bila
 * `ProportionType='NonProportional' && EDMState=3` (dan
 * `TreatyInTabsNonProportional` bila `EDMState!=3`). `Actual Retro`
 * dibangkitkan tetapi disembunyikan bersama Retro lain.
 */
export const TAB_BARU_NP_ADJ = [
  'Actual GNPI', 'Actual Limits', 'Actual Share', 'Actual Retro', 'Premium Adjustment', 'Information & Submit',
] as const // TreatyInTabsNonProportionalAdjustPremi @20240 … @596060
export const TAB_LAMA_P = [
  'Reporting Period', 'Portfolio', 'Limits', 'Share', 'Accumulation', 'Exclusions', 'Special Conditions',
] as const // TreatyInTabsProportionalOldData @23096 … @1167630
export const TAB_BARU_P = [
  'Reporting Period', 'Portfolio', 'Limits', 'Share', 'Retro', 'Co-Ins Scale', 'Accumulation',
  'Exclusions', 'Special Conditions', 'Information & Submit', 'Achievement In IDR',
] as const // TreatyInTabsProportional @26514 … @1581141

/** Deret tombol — `Section/TreatyInActionButtons.xml`, blok `TreatyMasterInEDM` @103200. */
export const TOMBOL = {
  simpan: 'Save', // pyLabel @113510 — syarat @119719
  tutup: 'Close', // pyLabel @130374 — ALWAYS @137426
  aksi: 'Actions', // pyLabel @148044 → TreatyInActionEDM — syarat peran @155994
  menyimpan: 'Menyimpan…',
  /** Modal `TreatyInActionEDM` — teks disalin dari Section-nya. */
  judulAksi: 'Action',
  dariID: 'Of ID :', // pyLabelFieldValue, `TreatyIn.ID`
  pilihan: 'Status', // `TreatyIn.ChooseStatusAkseptasi` (tanpa label di ekspor)
  komentar: 'Comment', // `TreatyIn.Comment`
  batal: 'Cancel',
  kirim: 'Submit', // → `TreatyInAkseptasiEDM_Act`
  /** Modal `TreatyInDeclineConfirmationEDM`. */
  judulTolak: 'Decline offer',
  tanyaTolak: 'Are you sure you want to DECLINE this offer?',
  tolak: 'Decline', // → `TreatyInDeclineConfirmation_postactEDM`
  /** Properti terkirim yang tidak tersimpan — dilaporkan, tidak ditelan. */
  takTersimpan: 'TIDAK tersimpan (belum punya kolom/tabel):',
} as const

/**
 * Pilihan `ChooseStatusAkseptasi` — nilai yang `Akseptasi_DT` periksa.
 * ⚠️ Daftarnya `associated` (rule properti tidak diekspor); teksnya nilai itu
 * sendiri.
 */
export const PILIHAN_AKSEPTASI = ['Accept', 'Reject', 'Decline'] as const

/** Kolom History mode detail — `TreatyIn.CommentList` @407001, sel @428448…@446102. */
export const KUNCI_HISTORY = ['Date', 'OperatorName', 'IsApproved', 'Suggest'] as const
