// Label modul NB FacIn - tiket 21.
//
// ⛔ VERBATIM dari `D:\migrasi\RNM\NB FacIn\Section\InputCoverageCargo_FacIn.xml`
// (ASM-FW-GISFW-DATA-COVERAGE!INPUTCOVERAGECARGO_FACIN): teks `pyLabelFor` sel yang
// disebut di sebelahnya. `labels.test.ts` membuka korpus dan memeriksa tiap pasangan.

/** Blok pertama section, urut layout Pega. `sel` = `pyCellId`, `properti` = `pyValue`. */
export const MEDAN_COVERAGE_CARGO = {
  keteranganJaminan: { sel: '3', properti: '.CoverageNote', label: 'Keterangan Jaminan' },
  coverageInitial: { sel: '4', properti: '.CoverageInitial', label: 'Coverage Initial' },
  mataUang: { sel: '10', properti: '.Currency.Name', label: 'Name' },
  rate: { sel: '11', properti: '.Rate', label: 'Rate (%)' },
  limitOfLiability: { sel: '12', properti: '.LimitofLiability', label: 'LimitofLiability' },
  currencyMaster: { sel: '15', properti: '.CurrencyMaster', label: '.CurrencyMaster' },
  diskonPersen: { sel: '18', properti: '.DiscountPercentage', label: 'Diskon (%)' },
  tsi: { sel: '21', properti: '.TSI', label: 'TSI' },
  premi: { sel: '22', properti: '.Premium', label: 'Premi' },
  minPremi: { sel: '23', properti: '.MinPremium', label: 'Min Premi' },
  diskon: { sel: '24', properti: '.Discount', label: 'Diskon' },
} as const

/** Tombol Pega di blok ini (`pyLabel`), belum diport - tampil nonaktif. */
export const TOMBOL_COVERAGE_CARGO = {
  /** Sel 5 - runActivity `SetIndexMarineCargo_Act` + showHarness popup. */
  pilihCoverage: { sel: '5', label: 'Select Coverage' },
  /** Sel 17 - refresh thisSection, `GetCurrMasterCargo`. */
  ambilCurrencyMaster: { sel: '17', label: 'Get Currency Master' },
} as const

/**
 * Teks sistem baru - BUKAN dari korpus (karena itu tidak diuji terhadap korpus).
 * Penyimpangan butir 61 (keputusan work owner 02-10-2026) dinyatakan terang di layar.
 */
export const TEKS_COVERAGE_CARGO = {
  judul: 'Coverage MARINE CARGO',
  alatPeriksa:
    'Alat periksa (butir 61): Name, Rate (%), dan TSI dapat diisi untuk menghitung premi; di Pega ketiganya hanya ' +
    'ditampilkan, dan Name berupa dropdown (di sini kotak teks). Data kasus NB belum tersambung (tiket 17) - medan ' +
    'lain kosong.',
  hitung: 'Hitung premi (alat periksa)',
  menghitung: 'Menghitung…',
  belumDiport: 'Belum diport',
  asalRumus: 'Asal rumus',
} as const

// ---------------------------------------------------------------------------
// Halaman depan - portal Opportunity (tiket 25)
// ---------------------------------------------------------------------------
//
// ⛔ VERBATIM dari korpus `D:\migrasi\RNM\NB FacIn\Section\`: `sel` = `pyCellId`, `tag` = elemen XML
// tempat teksnya tinggal. `labels.test.ts` membuka section dan memeriksa tiap pasangan.

/** `SFAPortalOpportunitiesHeader.xml` (PEGACRM-PORTAL!SFAPORTALOPPORTUNITIESHEADER). */
export const KEPALA_PORTAL = {
  /** Sel 57, LABEL - L516. */
  judul: { sel: '57', tag: 'pyValue', label: 'Opportunity' },
  /** Sel 72, pxButton - L2486. Di Pega `createWork` (`D_crmAppExtPage`, tidak ada di korpus). */
  buat: { sel: '72', tag: 'pyLabel', label: 'Create opportunity' },
} as const

/** `SFAPortal_OpportunitiesList.xml` (DATA-PORTAL!SFAPORTAL_OPPORTUNITIESLIST) - kotak saring. */
export const SARING_PORTAL = {
  /** Sel 9, pxTextInput `.FilterTermForOpportunity` - label TIDAK dirender (`pyIncludeLabel=false`). */
  label: { sel: '9', tag: 'pyLabelFor', label: 'Filter Term for Opportunity' },
  /** Sel 9 - L1655. */
  placeholder: { sel: '9', tag: 'pyPlaceholder', label: 'NB-1234 or Name' },
  /** Sel 13, pxButton - L2902. */
  tombol: { sel: '13', tag: 'pyLabel', label: 'Filter' },
} as const

/**
 * Grid `GetListOpportunityF` (keputusan agent B-1, tiket 25): judul kolom berurutan, sel 89-96 -
 * L14643-L15727. Kolom 3 dan 7 memang berjudul KOSONG di Pega (tanpa `pyValue`).
 */
export const KOLOM_PORTAL = [
  { sel: '89', label: 'Offer No', properti: '.TextNoQuotation' },
  { sel: '90', label: 'Name', properti: '.Name' },
  { sel: '91', label: '', properti: 'View' },
  { sel: '92', label: 'Group Business', properti: 'A.Quotation.BusinessName' },
  { sel: '93', label: 'Insured Name', properti: 'A.Quotation.InsuredName' },
  { sel: '94', label: 'Marketing', properti: 'A.Quotation.MarketingName' },
  { sel: '95', label: '', properti: 'A.NBStatus' },
  { sel: '96', label: 'Status', properti: 'A.NBStatusNew' },
] as const

/** Teks sistem baru - BUKAN dari korpus. */
export const TEKS_PORTAL = {
  /** Daftar berhasil dimuat tetapi kosong / tidak ada yang cocok dengan saringan. */
  tanpaCase: 'Tidak ada case yang cocok.',
  /** Nama aksesibel ikon hapus isian (sel 10, `pxIcon` tanpa teks). */
  hapusIsian: 'Hapus isian saring',
} as const

// ---------------------------------------------------------------------------
// Form "Opportunity" (tiket 26)
// ---------------------------------------------------------------------------
//
// ⛔ VERBATIM dari TANGKAPAN LAYAR Pega kiriman work owner 02-10-2026
// (`docs/02-layar/tangkapan/opportunity-pega-02-10-2026.png`). Form ini TIDAK ada di korpus
// (tiket 26 "Batas bukti"), jadi teks ini tidak dapat diuji terhadap XML.

export const FORM_OPPORTUNITY = {
  judul: 'Opportunity',
  tanggalTutup: 'Estimated Closing Date',
  /**
   * Di kanan Estimated Closing Date, berisi nama pengguna yang login - dari gambar keadaan awal
   * `D:\migrasi\RNM\DDL\HALAMAN DEPAN NB.JPG` (md5 9f4813a2c9c56ba12306232436453ca7; nama di gambar TIDAK disalin).
   */
  owner: 'Owner',
  namaProspek: 'Business Prospect Name',
  grupBisnis: 'Group Business',
  classOfBusiness: 'Class Of Business',
  typeOfInward: 'Type Of Inward',
  typeOfFacultative: 'Type Of Facultative',
  deskripsi: 'Description',
  phase: 'Phase',
  stage: 'Stage',
  sumber: 'Opportunity Source',
  statusBisnis: 'Business Status',
} as const

/** Tiga tombol Group Business di gambar - aksinya tidak diketahui (keputusan agent C-2: nonaktif). */
export const TOMBOL_FORM_OPPORTUNITY = {
  cariGrup: 'Search Group Business',
  perusahaanBaru: 'New Company Detail',
  grupBaru: 'New Group Business',
} as const

/**
 * Nilai yang TERLIHAT di gambar (keputusan agent C-1). Bukan daftar pilihan: daftar lengkapnya belum
 * terverifikasi, dan daftar karangan di React dilarang pola `Pilih` inti.
 */
export const NILAI_AWAL_OPPORTUNITY = {
  /**
   * Satu-satunya pilihan Type Of Inward yang terlihat (tangkapan layar pertama). Keadaan AWAL-nya kosong
   * berteks `Choose Type Of Inward` (`HALAMAN DEPAN NB.JPG`), dan Type Of Facultative baru tampil sesudah
   * Facultative dipilih - `[dugaan]` dari dua gambar itu.
   */
  typeOfInward: 'Facultative',
  inwardKosong: 'Choose Type Of Inward',
  typeOfFacultative: 'Facultative In',
  phase: 'Proposal',
  stage: 'Opportunity',
  /** Teks pilihan kosong dropdown Opportunity Source. */
  sumberKosong: 'Select...',
  statusBisnis: 'New Business',
} as const

/**
 * Pilihan dropdown `Opportunity Source` - VERBATIM dari tangkapan layar Pega dropdown terbuka, kiriman
 * work owner 02-10-2026 (`docs/02-layar/tangkapan/opportunity-source-pega-02-10-2026.png`), urutan
 * sama dengan gambar. `Select...` (pilihan kosong) tidak termasuk - ia `sumberKosong` di atas.
 *
 * ⚠️ `belum terverifikasi`: NILAI yang disimpan Pega untuk tiap pilihan (teks itu sendiri atau kode)
 * tidak terlihat di gambar - di sini nilai = teks tampilnya.
 */
export const OPSI_OPPORTUNITY_SOURCE = [
  'Iklan',
  'Analisa Referral',
  'Rujukan Pelanggan',
  'Direct Mail',
  'Email',
  'Rujukan Karyawan',
  'Telepon Masuk',
  'Partner',
  'Seminar',
  'Sosial Media',
  'Pameran',
  'Web',
  'Whatsapp',
] as const

/**
 * Popup tombol `Search Group Business` - VERBATIM dari tangkapan layar Pega kiriman work owner 02-10-2026
 * (md5 e4944778178e583eeac1ce1b7d87b0f3; gambar TIDAK disalin ke repo karena memuat nama pelanggan).
 * Rule-nya (`ChooseAccount`) tidak ada di korpus. Sumber data menurut work owner: tabel `T_M_ACCOUNT`
 * (DDL `D:\migrasi\RNM\DDL\T_M_ACCOUNT.txt`), lewat `GET /api/nbfacin/account` (tiket 27).
 */
export const POPUP_CHOOSE_ACCOUNT = {
  judul: 'ChooseAccount',
  cari: 'Search',
  tombolCari: 'Search',
  kolom: ['Insured ID', 'Insured Name', 'Group Business'],
  pilih: 'Choose',
} as const

/** Teks sistem baru untuk Group Business sesudah Choose - BUKAN dari Pega. */
export const TEKS_GRUP_BISNIS = {
  /** Nama aksesibel ikon roda gigi di samping Group Business terpilih (gambar 02-10-2026). */
  ganti: 'Ganti Group Business',
} as const

/** Teks sistem baru - BUKAN dari Pega. */
export const TEKS_FORM_OPPORTUNITY = {
  /** Nama aksesibel tombol kalender Estimated Closing Date. */
  kalender: 'Pilih tanggal',
  /** Isian tanggal 10 karakter tetapi bukan tanggal yang ada. */
  formatTanggal: 'Tanggal harus berformat dd/mm/yyyy.',
  /** Awal daftar medan wajib yang belum diisi saat Create opportunity. */
  wajibKosong: 'Lengkapi dulu:',
  /** Sesudah case berhasil dibuat; `{caseId}` diganti nomor case. */
  caseDibuat: 'Case {caseId} berhasil dibuat.',
  menyimpan: 'Menyimpan…',
  /** Pencarian berhasil tetapi tidak ada baris yang cocok. */
  tanpaAccount: 'Tidak ada account yang cocok.',
} as const

// ---------------------------------------------------------------------------
// Layar Inward Facultative - assignment pertama case NB (tiket 30, tahap 1)
// ---------------------------------------------------------------------------
//
// ⛔ VERBATIM dari korpus `D:\migrasi\RNM\NB FacIn\`; `labels.test.ts` membuka berkasnya. Susunan yang tampil
// untuk kasus FIRE dicocokkan dengan tangkapan layar Pega work owner 03-10-2026 (gambar TIDAK disalin ke repo:
// memuat nama pelanggan dan nama orang).

/** `Section\Periode.xml` (ASM-FW-GISFW-DATA-OFFERFACIN!PERIODE): `sel` = pyCellId, `tag` = elemen teksnya. */
export const PERIODE = {
  judul: { sel: '', tag: 'pyTitle', label: 'General' },
  reffNumber: { sel: '9', tag: 'pyLabelFieldValue', label: 'Reff. number' },
  businessStatus: { sel: '10', tag: 'pyLabelFieldValue', label: 'Business status' },
  insuredName: { sel: '16', tag: 'pyLabelFieldValue', label: 'Insured name' },
  qqName: { sel: '20', tag: 'pyLabelFieldValue', label: 'QQ name' },
  beginDate: { sel: '21', tag: 'pyLabelFieldValue', label: 'Begin date' },
  offeringDate: { sel: '22', tag: 'pyLabelFieldValue', label: 'Offering date' },
  policyType: { sel: '26', tag: 'pyLabelFieldValue', label: 'Policy Type' },
  riskScoring: { sel: '', tag: 'pyLabelFieldValue', label: 'Risk Scoring' },
  uploadQuotation: { sel: '38', tag: 'pyLabel', label: 'Upload QUOTATION/PLACING SLIP AI' },
  uploadRISlip: { sel: '39', tag: 'pyLabel', label: 'Upload R/I SLIP AI' },
  classOfBusiness: { sel: '42', tag: 'pyLabelFieldValue', label: 'Class of business' },
  typeFacultative: { sel: '43', tag: 'pyLabelFieldValue', label: 'Type facultative' },
  sourceOfBusiness: { sel: '48', tag: 'pyLabelFieldValue', label: 'Source of business' },
  cedingCoName: { sel: '49', tag: 'pyLabelFieldValue', label: 'Ceding co name' },
  changeSob: { sel: '52', tag: 'pyLabel', label: 'Change SOB' },
  changeCedingCo: { sel: '53', tag: 'pyLabel', label: 'Change Ceding Co' },
  groupName: { sel: '56', tag: 'pyLabelFieldValue', label: 'Group Name' },
  endDate: { sel: '60', tag: 'pyLabelFieldValue', label: 'End date' },
  followingPolicyNumber: { sel: '66', tag: 'pyValue', label: 'Following Policy Number' },
  search: { sel: '69', tag: 'pyLabel', label: 'Search' },
  oldPolicyNumber: { sel: '72', tag: 'pyLabelFieldValue', label: 'Old Policy Number' },
  marketingName: { sel: '75', tag: 'pyLabelFieldValue', label: 'Marketing Name' },
  /** Teks pilihan kosong dropdown Marketing Name (dekat sel 75). */
  marketingKosong: { sel: '', tag: 'pyNoSelectionText', label: 'Choose' },
  day: { sel: '78', tag: 'pyLabelFieldValue', label: 'Day' },
  judulCsv: { sel: '', tag: 'pyTitle', label: 'Please upload file with .csv format' },
  downloadTemplateCsv: { sel: '93', tag: 'pyLabel', label: 'Download Template CSV' },
  uploadCsv: { sel: '94', tag: 'pyLabel', label: 'Upload CSV' },
  viewUpload: { sel: '95', tag: 'pyLabel', label: 'View Upload' },
  saveData: { sel: '96', tag: 'pyLabel', label: 'Save Data' },
  insertAccumulation: { sel: '97', tag: 'pyLabel', label: 'Insert Accumulation' },
} as const

/**
 * Pilihan radio Policy Type (sel 26) dan Day (sel 78) - daftarnya ikut definisi properti (`associated`, tipe
 * rule Property tidak ada di korpus), jadi teksnya dari TANGKAPAN LAYAR work owner 03-10-2026.
 */
export const PILIHAN_PERIODE = {
  policyType: ['Individual Policy', 'Master Policy'],
  day: ['365', '366'],
} as const


/** `Section\InputInwardFacultative.xml` - checkbox `.IsShowDetail` (`pyCheckboxCaption` L3288). */
export const SHOW_DETAIL = 'Show Detail'

/** `Section\InputInwardFacultativeDtl.xml` - judul tab (`pyTitle`), urut tangkapan layar kasus FIRE. */
export const TAB_DETAIL = [
  'Object',
  'Coverage',
  'Clauses',
  'Spreading',
  'Inw Fac Cedant Panels',
  'Premium Deduction / Brokerage Fee',
  'Loss Record',
  'Payment',
  'Scoring Risk',
  'Correspondence',
] as const

/** `FlowAction\InwardFacultative.xml` - tombol kaki (pySubmitLabel / pySaveLabel / pyCancelLabel). */
export const TOMBOL_KAKI_INWARD = {
  submit: { tag: 'pySubmitLabel', label: 'Submit' },
  simpan: { tag: 'pySaveLabel', label: 'Save for later' },
  batal: { tag: 'pyCancelLabel', label: 'Cancel' },
} as const

/** Teks sistem baru - BUKAN dari Pega. */
export const TEKS_INWARD = {
  /**
   * End date lebih awal dari Begin date. Pega (`SetValidateDate_Act` langkah 1, Property-Set-Messages di
   * `.PolicyData.EndDateTime`) memakai Rule-Message `ErrorMSG` yang TIDAK ada di korpus - teks ini sistem baru.
   */
  endSebelumBegin: 'End date tidak boleh lebih awal dari Begin date.',
  /** Ringkasan objek kosong (Pega "No items"). */
  kosong: 'No items',
  /** Untuk `BelumTersedia` isi tab detail (tahap 3). */
  isiTab: 'Isi tab',
  /** Sesudah Save for later berhasil. */
  tersimpan: 'Tersimpan.',
  /** Popup Change SOB: pencarian berhasil tetapi kosong. */
  tanpaSob: 'Tidak ada SOB yang cocok.',
} as const

/**
 * Popup tombol `Change SOB` (Periode sel 52 → harness `SOB`, `D:\migrasi\RNM\NB FacIn\Harness\SOB.xml`,
 * WindowName "Change SOB"). `Search` (`pyLabelFieldValue`) dan `Choose` (`pyLabel`) diuji terhadap SOB.xml; judul
 * kolom dari TANGKAPAN LAYAR work owner 03-10-2026 (properti grid `.ID`, `.ClientName`; Client ID tak terbaca di XML).
 */
export const POPUP_SOB = {
  judul: 'Change SOB',
  cari: { tag: 'pyLabelFieldValue', label: 'Search' },
  pilih: { tag: 'pyLabel', label: 'Choose' },
  kolom: ['ID', 'Client ID', 'Name'],
} as const

/**
 * Popup tombol `Change Ceding Co` (Periode sel 53 → showHarness `ShowCedingCoList`, Target popup) - tiket 34.
 * `NB FacIn\Section\ShowCedingCoList.xml`: grid atas `.Quotation.CedingCoList`; kepala sel 14 `Ceding Co` (wajib) dan
 * sel 15 tombol `Add / Select Ceding` (AddCedingList_act → showHarness `CedingCompany`, WindowName "Ceding Company");
 * baris sel 17 `.CedingCoName` (baca-saja) dan sel 18 `Delete` (DeleteCeding_Act). Judul jendela = pyLabel harness
 * `ShowCedingCoList`; `Submit` = tombol harness (SetCedingCo_Act).
 */
export const POPUP_CEDING = {
  judul: 'Ceding Co List',
  kolom: { sel: '14', tag: 'pyValue', label: 'Ceding Co' },
  tambah: { sel: '15', tag: 'pyLabel', label: 'Add / Select Ceding' },
  hapus: { sel: '18', tag: 'pyLabel', label: 'Delete' },
  submit: { tag: 'pyLabel', label: 'Submit' },
  /** WindowName showHarness `CedingCompany` (sel 15). */
  judulCari: 'Ceding Company',
} as const

/**
 * Tab Object (FIRE) - tiket 35. Dtl tab "Object" sel 21 meng-include section `ObjectList` (visible `IsFire`);
 * grid `.LocationList` (master-detail, 10 per halaman). Judul kolom = `pyValue` sel 13-16
 * `NB FacIn\Section\ObjectList.xml`. `Tambah` / `Hapus` = tangkapan layar work owner 03-10-2026 (di XML tombol
 * ikon `IconAdd.png` / `IconTrash.png` tanpa teks).
 */
export const GRID_OBJEK = {
  kolom: [
    { sel: '13', tag: 'pyValue', label: 'Top Risk' },
    { sel: '14', tag: 'pyValue', label: 'No.' },
    { sel: '15', tag: 'pyValue', label: 'Object Name' },
    { sel: '16', tag: 'pyValue', label: 'Location' },
  ],
  tambah: 'Tambah',
  hapus: 'Hapus',
  ukuran: 10,
} as const

/** Sub-tab baris objek - `pyTitle` `NB FacIn\Section\Property.xml` (flow action `Property_FlowAction`). */
export const SUBTAB_OBJEK = [
  'Object Address',
  'Surrounding Risk',
  'Object Item',
  'Occupation',
  'FEA',
  'Loss Record',
  'Loss Record Internal',
] as const

/** Sub-tab Object Address - `NB FacIn\Section\ObjectDetails.xml` (sel -> label, tag pembawa). */
export const OBJECT_ADDRESS = {
  objectNo: { sel: '5', tag: 'pyLabelFieldValue', label: 'Object No.' },
  objectType: { sel: '6', tag: 'pyLabelFieldValue', label: 'Object Type' },
  objectTypeKosong: { sel: '6', tag: 'pyNoSelectionText', label: 'Please Select' },
  materialDamage: { sel: '11', tag: 'pyCheckboxCaption', label: 'Material Damage' },
  topRisk: { sel: '12', tag: 'pyCheckboxCaption', label: 'Top Risk' },
  objectName: { sel: '13', tag: 'pyLabelFieldValue', label: 'Object Name' },
  chooseRisk: { sel: '16', tag: 'pyLabel', label: 'Choose Risk Address' },
  clearRisk: { sel: '17', tag: 'pyLabel', label: 'Clear Risk Address' },
  judulRisk: { sel: '', tag: 'pyTitle', label: 'Risk Address' },
  type: { sel: '28', tag: 'pyLabelFieldValue', label: 'Type' },
  address: { sel: '29', tag: 'pyLabelFieldValue', label: 'Address' },
  buildingNo: { sel: '30', tag: 'pyLabelFieldValue', label: 'Building No.' },
  zipCode: { sel: '31', tag: 'pyLabelFieldValue', label: 'Zip Code' },
  country: { sel: '32', tag: 'pyLabelFieldValue', label: 'Country' },
  riskLocation: { sel: '40', tag: 'pyLabelFieldValue', label: 'Risk Location' },
  territory: { sel: '35', tag: 'pyLabelFieldValue', label: 'Territory' },
  city: { sel: '36', tag: 'pyLabelFieldValue', label: 'City' },
  district: { sel: '37', tag: 'pyLabelFieldValue', label: 'District' },
  province: { sel: '38', tag: 'pyLabelFieldValue', label: 'Province' },
  riskAddressId: { sel: '39', tag: 'pyLabelFieldValue', label: 'Risk Address ID' },
  judulBangunan: { sel: '', tag: 'pyTitle', label: 'Building Construction' },
  numberOfFloor: { sel: '49', tag: 'pyLabelFieldValue', label: 'Number of Floor' },
  roofType: { sel: '50', tag: 'pyLabelFieldValue', label: 'Roof Type' },
  wallType: { sel: '51', tag: 'pyLabelFieldValue', label: 'Wall Type' },
  floorType: { sel: '52', tag: 'pyLabelFieldValue', label: 'Floor Type' },
  partitionType: { sel: '56', tag: 'pyLabelFieldValue', label: 'Partition Type' },
  supportWallType: { sel: '57', tag: 'pyLabelFieldValue', label: 'Support Wall Type' },
  otherType: { sel: '58', tag: 'pyLabelFieldValue', label: 'Other Type' },
} as const

/** Tombol Save tab Object - Dtl sel 25 (`pyLabel`, runActivity `SaveFacIn_Act`). */
export const SIMPAN_OBJEK = { sel: '25', tag: 'pyLabel', label: 'Save' } as const

/**
 * Pilihan Object Type - urutan = tangkapan layar dropdown work owner 03-10-2026. Sumber daftar Pega
 * (`pyListSource=associated`) tidak ada di korpus; HIMPUNANNYA cocok dengan ekspresi
 * `@if(.OBJECT_TYPE!="Dwelling House"&&...,"Others",...)` `Activity\InsertUploadFire_act.xml` (diuji).
 */
export const OPSI_OBJECT_TYPE = [
  'Dwelling House',
  'Shop Houses',
  'Office',
  'Shop',
  'Apartment',
  'Restaurant',
  'Private Warehouse',
  'Public Warehouse',
  'Factory',
  'Others',
] as const

/**
 * Nilai Object Type yang membuka medan Object Name (sel 13 visible `.Property.ObjectType = 'Others'`).
 * ⚠️ `SetValueOnObjectName_Act` memeriksa 'Lainnya' - nilai yang tidak ada di daftar; diikuti 'Others' (G-2).
 */
export const OBJECT_TYPE_LAINNYA = 'Others'

/** Teks tab Object. */
export const TEKS_OBJEK = {
  /** `Activity\SetErrorMessageFloorNumber_Act.xml` (Property-Set-Messages, `NumberOfFloor < 0`). */
  lantaiMinus: "Floor number can't be minus",
  /** Sistem baru - BUKAN dari Pega: Object Type wajib (sel 6 `pyRequired`) sebelum Save. */
  typeWajib: 'Object Type wajib diisi.',
  /** Sistem baru: nama aksesibel tombol buka/tutup baris. */
  bukaBaris: 'Buka/tutup detail objek',
} as const

/**
 * Pilihan Roof / Wall / Floor Type (Building Construction, sel 50-52, `pyListSource=associated`) - `[terverifikasi]`
 * aturan properti `ASM-FW-GISFW-DATA-BUILDINGCONSTRUCTION!ROOFTYPE` / `!WALLTYPE` / `!FLOORTYPE` (PromptList,
 * `D:\migrasi\RNM\DDL\RoofType.xml`, `WallType.xml`, `FloorType.xml`, ditambah work owner 03-10-2026):
 * value = `pyStandardValue`, label = `pyLocalizedValue`; baris tanpa nilai = "Silahkan Pilih". Diuji `labels.test.ts`.
 */
export const OPSI_ROOF_TYPE = [
  'Dak Beton',
  'Genteng Beton',
  'Bilik',
  'Genteng Tanah Liat',
  'Sirap',
  'Seng Gelombang',
  'Seng Lembaran',
  'Aluminium Gelombang',
  'Aluminium Lembaran',
  'Kaca',
  'Plastik / Policarbon Lembaran',
  'Plastik / Policarbon Gelombang',
  'Daun',
  'Lain-lain',
].map((label, i) => ({ value: String(i + 1), label }))

export const OPSI_WALL_TYPE = [
  'Batu bata',
  'Kayu / Papan',
  'Semi Permanen',
  'Bilik',
  'Batako',
  'Panel Beton',
  'Seng / Plat',
  'Kaca',
  'Lain-lain',
].map((label, i) => ({ value: String(i + 1), label }))

export const OPSI_FLOOR_TYPE = [
  { value: 'KELAS III', label: 'Keramik' },
  { value: 'KELAS II', label: 'Kayu' },
  { value: 'KELAS I', label: 'Lain-lain' },
]

/**
 * Nilai awal Building Construction objek baru = "Lain-lain" di ketiganya (Roof 14, Wall 9, Floor "KELAS I"):
 * seluruh 332 entri data contoh `DDL\CONTOH` dan tangkapan layar objek baru. Aturan properti tidak memuat
 * `pyDefaultValue` - sumber nilai awalnya di Pega `belum terverifikasi`.
 */
export const AWAL_BANGUNAN = { roofType: '14', wallType: '9', floorType: 'KELAS I' } as const

/** Teks pilihan kosong Roof / Wall / Floor Type - baris pertama PromptList ketiga aturan properti (sel 50: "Please Select"). */
export const BANGUNAN_KOSONG = 'Silahkan Pilih'
