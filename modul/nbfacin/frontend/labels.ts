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

/**
 * Popup Choose Risk Address (tiket 36) - ObjectDetails sel 16 showHarness `ChooseRiskAddress` (WindowName
 * "Choose Risk Location", Target popup). Saringan = `NB FacIn\Section\ChooseRiskAddress.xml` sel 78-84
 * (`pyLabelFieldValue`), tombol sel 73 `Search` / 74 `Add` (`pyLabel`); grid = `ChooseRiskAddress_ResultList.xml`
 * sel 19-26 (`pyValue`), 10 per halaman bernomor. `Pilih` = tangkapan layar work owner 03-10-2026 (di XML ikon
 * `IconChoose.png` tanpa teks).
 */
export const POPUP_RISK = {
  judul: 'Choose Risk Location',
  saring: {
    address: { sel: '78', tag: 'pyLabelFieldValue', label: 'Address' },
    zipCode: { sel: '79', tag: 'pyLabelFieldValue', label: 'Zip Code' },
    country: { sel: '80', tag: 'pyLabelFieldValue', label: 'Country' },
    province: { sel: '81', tag: 'pyLabelFieldValue', label: 'Province' },
    city: { sel: '82', tag: 'pyLabelFieldValue', label: 'City' },
    district: { sel: '83', tag: 'pyLabelFieldValue', label: 'District' },
    territory: { sel: '84', tag: 'pyLabelFieldValue', label: 'Territory' },
  },
  cari: { sel: '73', tag: 'pyLabel', label: 'Search' },
  tambah: { sel: '74', tag: 'pyLabel', label: 'Add' },
  kolom: [
    { sel: '19', tag: 'pyValue', label: 'Type' },
    { sel: '20', tag: 'pyValue', label: 'Address' },
    { sel: '21', tag: 'pyValue', label: 'Country' },
    { sel: '22', tag: 'pyValue', label: 'Province' },
    { sel: '23', tag: 'pyValue', label: 'City' },
    { sel: '24', tag: 'pyValue', label: 'District' },
    { sel: '25', tag: 'pyValue', label: 'Territory' },
    { sel: '26', tag: 'pyValue', label: 'Zip Code' },
  ],
  pilih: 'Pilih',
  ukuran: 10,
} as const

/** Teks sistem baru popup Risk Address - BUKAN dari Pega. */
export const TEKS_RISK = {
  /** Search tanpa satu pun saringan (Pega `SearchRiskAddressAct` tidak mencari bila semua kosong). */
  isiSaring: 'Isi minimal satu saringan lalu tekan Search.',
  tanpaHasil: 'Tidak ada alamat risiko yang cocok.',
} as const

/**
 * Popup Add alamat risiko (tiket 37) - tombol `Add` popup Choose Risk Location (ChooseRiskAddress sel 74 ->
 * harness `ChooseRiskLocation` -> section `NB FacIn\Section\InputRiskAddress.xml`). Label = `pyLabelFieldValue`
 * sel 3-21, tombol sel 28 `Save` / 29 `Close` (`pyLabel`). Medan tersembunyi (`never`): Code, Type, Status Trans
 * Pusat - tidak diport. Judul modal = label tombol pembukanya (harness tanpa WindowName).
 */
export const POPUP_TAMBAH_RISK = {
  judul: 'Add',
  country: { sel: '3', tag: 'pyLabelFieldValue', label: 'Country' },
  province: { sel: '5', tag: 'pyLabelFieldValue', label: 'Province' },
  city: { sel: '7', tag: 'pyLabelFieldValue', label: 'City' },
  district: { sel: '9', tag: 'pyLabelFieldValue', label: 'District' },
  territory: { sel: '11', tag: 'pyLabelFieldValue', label: 'Territory' },
  zipCode: { sel: '13', tag: 'pyLabelFieldValue', label: 'Zip Code' },
  title: { sel: '20', tag: 'pyLabelFieldValue', label: 'Title' },
  address: { sel: '21', tag: 'pyLabelFieldValue', label: 'Address' },
  simpan: { sel: '28', tag: 'pyLabel', label: 'Save' },
  tutup: { sel: '29', tag: 'pyLabel', label: 'Close' },
} as const

/**
 * Pilihan Title - tangkapan layar work owner 03-10-2026 (sumber Pega `associated`, aturan properti tidak ada di
 * korpus). Urutan sama dengan REGEXP `(DESA|DUSUN|GANG|GEDUNG|JL\.|KOMPLEK|PERUMAHAN|OTHERS)` di
 * `NB FacIn\RDBList\SearchAccumulationbypersetase_SQL.xml` (diuji). Tanpa pilihan kosong (`pyHasNoSelection=false`)
 * -> nilai awal = pilihan pertama (DESA). Nilai = teks (dirangkai ke Risk Location).
 */
export const OPSI_TITLE_RISK = ['DESA', 'DUSUN', 'GANG', 'GEDUNG', 'JL.', 'KOMPLEK', 'PERUMAHAN', 'OTHERS'] as const

/** Teks sistem baru popup Add - BUKAN dari Pega. */
export const TEKS_TAMBAH_RISK = {
  /** Keputusan agent I-4: Pega tidak memvalidasi; baris master kosong dicegah. */
  wajib: 'Zip Code dan Address wajib diisi.',
  menyimpan: 'Menyimpan…',
  tanpaZip: 'Zip Code tidak ditemukan di tabel RW.',
} as const

/**
 * Sub-tab Surrounding Risk (tiket 38) - `NB FacIn\Section\RiskAround.xml`. Empat blok sisi (pyTitle Front / Left /
 * Back / Right) masing-masing: Occupation (autocomplete `D_BrowseOccupationFacInFIRE`), Construction (dropdown),
 * Distance (meter), Note (baca-saja, diisi nama Occupation). Sel per sisi diuji `labels.test.ts`.
 */
export const SISI_SEKITAR = [
  { kunci: 'front', judul: 'Front', occupation: '9', construction: '12', distance: '13', note: '14' },
  { kunci: 'left', judul: 'Left', occupation: '23', construction: '26', distance: '27', note: '28' },
  { kunci: 'back', judul: 'Back', occupation: '37', construction: '40', distance: '41', note: '42' },
  { kunci: 'right', judul: 'Right', occupation: '51', construction: '54', distance: '55', note: '56' },
] as const

/** Label medan per sisi (`pyLabelFieldValue`, sama di keempat sisi). */
export const MEDAN_SISI = {
  occupation: 'Occupation',
  construction: 'Construction',
  distance: 'Distance (meter)',
  note: 'Note',
} as const

/** Blok Other Description (sel 59 / 70, pyTitle) - label `pyLabelFieldValue` / `pyCheckboxCaption`. */
export const LAIN_SEKITAR = {
  judul: { sel: '', tag: 'pyTitle', label: 'Other Description' },
  ownership: { sel: '65', tag: 'pyLabelFieldValue', label: 'Ownership' },
  housekeepingStatus: { sel: '66', tag: 'pyLabelFieldValue', label: 'House keeping Status' },
  floodAreaStatus: { sel: '67', tag: 'pyLabelFieldValue', label: 'Flood Area Status' },
  floodArea: { sel: '68', tag: 'pyLabelFieldValue', label: 'Flood Area' },
  housekeepingRemark: { sel: '69', tag: 'pyLabelFieldValue', label: 'Housekeeping Remark' },
  productionProcess: { sel: '77', tag: 'pyCheckboxCaption', label: 'Production Process' },
  hotWork: { sel: '78', tag: 'pyCheckboxCaption', label: 'Job With a Chance of Fire' },
  flammable: { sel: '79', tag: 'pyCheckboxCaption', label: 'Flammable Item' },
} as const

/** Flood Area tampil hanya bila Flood Area Status = "0" (sel 68 `pyVisible OTHER`). */
export const FLOOD_STATUS_TAMPIL = '0'

/**
 * Daftar dropdown Surrounding Risk - `[terverifikasi]` aturan properti Pega (PromptList) yang ditambahkan work owner
 * 03-10-2026 di `D:\migrasi\RNM\DDL\`: value = `pyStandardValue`, label = `pyLocalizedValue`, pasangan PER
 * rowdata (diuji `labels.test.ts`).
 * - `FrontConstruction.xml` (`SURROUNDINGRISK!FRONTCONSTRUCTION`): baris pertama tanpa nilai = "Silahkan pilih".
 *   Left / Back / Right memakai daftar yang sama `[dugaan]` - hanya aturan Front yang dikirim; tangkapan layar
 *   menampilkan "Silahkan pilih" di keempat sisi. Nilai tersimpan = teks panjang (hingga ~210 karakter).
 * - `Ownership.xml`, `HousekeepingStatus.xml`, `FloodAreaStatus.xml`: TANPA baris kosong -> nilai awal = pilihan
 *   pertama "Not Informed" (cocok dengan data contoh: Ownership "2", HousekeepingStatus "0", FloodAreaStatus "2").
 * - `FloodArea.xml` (`SURROUNDINGRISK!FLOODAREA`): baris pertama tanpa nilai = "Silahkan pilih"; tampil hanya bila
 *   Flood Area Status = Yes ("0"). Nilainya dipakai `"0"+FloodArea` di `SetFloodParamFacIn_ACT`.
 */
export const OPSI_CONSTRUCTION = [
  { value: 'Reinforce concrete', label: 'I' },
  {
    value:
      'Building with all structural members and at least 80% of the walls constructed with non - combustible materials, roof covering to be non - combustible, but roof supports may be of wood or other combustible materials',
    label: 'II',
  },
  { value: 'All other buildings', label: 'III' },
]
export const CONSTRUCTION_KOSONG = 'Silahkan pilih'
export const OPSI_OWNERSHIP = [
  { value: '2', label: 'Not Informed' },
  { value: '0', label: 'Own' },
  { value: '1', label: 'Rent' },
]
export const OPSI_HOUSEKEEPING = [
  { value: '0', label: 'Not Informed' },
  { value: '1', label: 'Good' },
  { value: '2', label: 'Fair' },
  { value: '3', label: 'Poor' },
]
export const OPSI_FLOOD_STATUS = [
  { value: '2', label: 'Not Informed' },
  { value: '0', label: 'Yes' },
  { value: '1', label: 'No' },
]
export const OPSI_FLOOD_AREA = [
  { value: '1', label: 'Low' },
  { value: '2', label: 'Medium' },
  { value: '3', label: 'High' },
  { value: '4', label: 'Very High' },
]

/** Teks Surrounding Risk. */
export const TEKS_SEKITAR = {
  /**
   * Paragraf `InputFireObject_Q1_RiskFactorOption` (sel 76; isi aturan tidak ada di korpus) - teks dari tangkapan
   * layar work owner 03-10-2026.
   */
  faktorRisiko: 'Choose the factors below and give a description for any of it. Please Skip this step if the factor you need is not available.',
  /** `Activity\NegativeIsNotAllowed.xml` (`Local.notminus`). */
  jarakMinus: 'Jarak Resiko Sekitar tidak boleh Minus',
  /** Sistem baru: saran Occupation kosong. */
  tanpaOccupation: 'Tidak ada occupation yang cocok.',
} as const

/**
 * Sub-tab Object Item (tiket 39) - `NB FacIn\Section\PropertyItemList.xml` (grid `.Property.PropertyItemList`,
 * layout NB `!IsEDM`, sel kepala 19-24; grid Total `.Property.TotalTSIList` sel 59-60) dan form detail
 * `Section\PropertyItemFacIn_Section.xml` (flow action `PropertyItemFacIn_FlowAction`). Diuji `labels.test.ts`.
 */
export const GRID_ITEM = {
  kolom: [
    { sel: '19', label: 'Object Item Type' },
    { sel: '20', label: 'Condition' },
    { sel: '21', label: 'Year' },
    { sel: '22', label: 'Unit(s)' },
    { sel: '23', label: 'Currency' },
    { sel: '24', label: 'TSI Object Item' },
  ],
  total: [
    { sel: '59', label: 'Currency' },
    { sel: '60', label: 'Total TSI' },
  ],
} as const

/** Form detail baris item (`pyLabelFieldValue`, kecuali centang `pyCheckboxCaption`). */
export const FORM_ITEM = {
  itemType: { sel: '5', tag: 'pyLabelFieldValue', label: 'Object Item Type' },
  note: { sel: '6', tag: 'pyLabelFieldValue', label: 'Object Item Note' },
  year: { sel: '11', tag: 'pyLabelFieldValue', label: 'Year' },
  unit: { sel: '12', tag: 'pyLabelFieldValue', label: 'Unit(s)' },
  condition: { sel: '13', tag: 'pyLabelFieldValue', label: 'Condition' },
  currency: { sel: '16', tag: 'pyLabelFieldValue', label: 'Currency' },
  tsi: { sel: '17', tag: 'pyLabelFieldValue', label: 'TSI Object Item (All Unit)' },
  yearOfPlanting: { sel: '20', tag: 'pyLabelFieldValue', label: 'Year of Planting' },
  noOfTree: { sel: '21', tag: 'pyLabelFieldValue', label: 'No of Trees' },
  areaHectar: { sel: '22', tag: 'pyLabelFieldValue', label: 'Area ( Hectar )' },
  remark: { sel: '27', tag: 'pyLabelFieldValue', label: 'Remark of' },
  adjustable: { sel: '31', tag: 'pyCheckboxCaption', label: 'Adjustable' },
  pctAdjust: { sel: '36', tag: 'pyLabelFieldValue', label: 'Adjustment Pct. %' },
} as const

/** Pilihan kosong dropdown Object Item Type / Currency (form detail: "Choose"). */
export const ITEM_KOSONG = 'Choose'

/**
 * Daftar Condition Object Item - `[terverifikasi]` aturan properti `ASM-FW-GISFW-DATA-PROPERTYITEM!CONDITION`
 * (PromptList, `DDL\ConditionObjectItem.xml`, dikirim ulang work owner 03-10-2026 setelah `DDL\Condition.xml` diisi
 * aturan Deductible bernama sama - uji mencari menurut pxInsName). Baris pertama tanpa nilai = "Please Select".
 */
export const OPSI_CONDITION = [
  { value: '1', label: 'Good' },
  { value: '2', label: 'Fair' },
  { value: '3', label: 'Poor' },
]
export const CONDITION_KOSONG = 'Please Select'

/**
 * Daftar Adjustment Pct. (PctAdjust2) - `[terverifikasi]` aturan properti `ASM-FW-GISFW-DATA-PROPERTYITEM!PCTADJUST2`
 * (`D:\migrasi\RNM\DDL\PctAdjust2.xml`, pyTableOption LocalList): satu nilai "100". Data contoh: PctAdjust2 = 100
 * di semua item -> nilai awal item baru "100" (keputusan agent K-7).
 */
export const OPSI_PCT_ADJUST = [{ value: '100', label: '100' }]
export const PCT_ADJUST_AWAL = '100'

/** Teks Object Item. */
export const TEKS_ITEM = {
  /** `Activity\ValidateAdjustPct.xml` (PctAdjustOther < 60 atau > 100). */
  pctAdjust: "%Adjustment can't be less than 60% or more than 100%",
  /** Sistem baru: field value Pega `ErrorMessageUnit` tidak ada di korpus (`SetErrorMessageUnit_Act`, Unit <= 0). */
  unit: 'Unit(s) harus lebih dari 0.',
  /** Sistem baru: field value Pega `TSIObjectItemErrorMessage` tidak ada di korpus (TSI < 0). */
  tsiMinus: 'TSI Object Item tidak boleh minus.',
  /** Sistem baru: TSI bukan angka desimal bertitik. */
  tsiBukanAngka: 'TSI Object Item harus angka (pemisah desimal titik).',
  /**
   * Sistem baru (keputusan agent A133, tiket 39): Currency wajib - kolom rancangan CURRENCY NOT NULL DEFAULT 'UNKNOWN'
   * dan aplikasi dilarang menulis 'UNKNOWN' (K-012). Pega sendiri tidak mewajibkannya.
   */
  currencyWajib: 'Currency wajib diisi.',
} as const

/**
 * Sub-tab Occupation (tiket 40) - `NB FacIn\Section\OccupationList.xml` (grid `.Property.OccupationList`, kepala sel
 * 14-16), form `OccupationItemFacIn_Section.xml`, popup `ChooseOccupation.xml` dan `ChooseClassofContraction.xml`.
 * Diuji `labels.test.ts`.
 */
export const GRID_OKUPASI = [
  { sel: '14', label: 'Occupation ID' },
  { sel: '15', label: 'Occupation Name' },
  { sel: '16', label: 'Class Of Construction' },
] as const

export const FORM_OKUPASI = {
  pilihOkupasi: { sel: '3', tag: 'pyLabel', label: 'Choose Occupation' },
  occupationId: { sel: '4', tag: 'pyLabelFieldValue', label: 'Occupation ID' },
  occupationName: { sel: '5', tag: 'pyLabelFieldValue', label: 'Occupation Name' },
  pilihKonstruksi: { sel: '8', tag: 'pyLabel', label: 'Choose Class of Construction' },
  konstruksi: { sel: '9', tag: 'pyLabelFieldValue', label: 'Class of Construction' },
} as const

/** Popup Choose Occupation (`ChooseOccupation.xml`): kotak sel 1, kolom sel 16-17, tombol sel 22. */
export const POPUP_OKUPASI = {
  cari: { sel: '1', tag: 'pyLabelFieldValue', label: 'Search Name/ID' },
  kolom: [
    { sel: '16', label: 'ID' },
    { sel: '17', label: 'Name' },
  ],
  pilih: { sel: '22', tag: 'pyLabel', label: 'Choose' },
  ukuran: 20,
} as const

/** Popup Choose Class of Construction (`ChooseClassofContraction.xml`): kolom sel 15 (Limit sel 16 tersembunyi). */
export const POPUP_KONSTRUKSI = {
  kolom: { sel: '15', label: 'Description' },
  pilih: { sel: '21', tag: 'pyLabel', label: 'Choose' },
} as const

/** Teks sistem baru sub-tab Occupation. */
export const TEKS_OKUPASI = {
  tanpaHasil: 'Tidak ada data yang cocok.',
  /** Choose Class of Construction sebelum Occupation dipilih (Category kosong). */
  pilihOkupasiDulu: 'Pilih Occupation lebih dulu.',
} as const

/**
 * Sub-tab FEA = Fire Extinguisher Availability (flow action `InputFEA`, pyLabel "Input Fire Extinguisher Availability")
 * - tiket 41. Grid `NB FacIn\Section\FEAList.xml` (`.FEAList` baris objek, kepala sel 15-20). Form isian: section
 * `OfferFEAList!InputFEA` TIDAK ada di korpus; dipakai padanannya `Section\InputFEA_IsUW.xml` (kelas sama, baca-saja)
 * sel 12-21. Diuji `labels.test.ts`.
 */
export const GRID_FEA = [
  { sel: '15', label: 'APAR' },
  { sel: '16', label: 'Sprinkler' },
  { sel: '17', label: 'Smoke Detector & Alarm' },
  { sel: '18', label: 'Hydrant' },
  { sel: '19', label: 'Private Truck Brigade' },
  { sel: '20', label: 'Others Info' },
] as const

export const FORM_FEA = {
  apar: { sel: '12', label: 'APAR (Unit)' },
  sprinkler: { sel: '13', label: 'Sprinkler (Unit)' },
  smokeDetector: { sel: '14', label: 'Smoke Detector & Alarm (Unit)' },
  hydrant: { sel: '15', label: 'Hydrant (Unit)' },
  privateTruckBrigade: { sel: '16', label: 'Private Truck Brigade (Unit)' },
  privateFireBrigade: { sel: '17', label: 'Private Team Fire Brigade' },
  teamSopSafety: { sel: '18', label: 'Team & SOP Safety' },
  teamSopRiskManagement: { sel: '19', label: 'Team & SOP Risk Management' },
  info: { sel: '21', label: 'Others Info' },
} as const

/**
 * Daftar Private Team Fire Brigade / Team & SOP Safety / Team & SOP Risk Management - `[terverifikasi]` aturan properti
 * `ASM-FW-GISFW-DATA-FEA!PRIVATEFIREBRIGADE` / `!TEAMSOPSAFETY` / `!TEAMSOPRISKMANAGEMENT` (PromptList,
 * `D:\migrasi\RNM\DDL\PrivateFireBrigade.xml` dst.; tanpa baris kosong). Diuji `labels.test.ts`.
 */
const OPSI_ADA_TIDAK = [
  { value: 'Have', label: 'Have' },
  { value: 'Not Have', label: 'Not Have' },
  { value: 'No Info', label: 'No Info' },
]
export const OPSI_FIRE_BRIGADE = OPSI_ADA_TIDAK
export const OPSI_SOP_SAFETY = OPSI_ADA_TIDAK
export const OPSI_SOP_RISIKO = OPSI_ADA_TIDAK

/** Teks sistem baru sub-tab FEA. */
export const TEKS_FEA = {
  /** Jumlah unit (pxNumber) - keputusan agent M-2: bilangan bulat >= 0. */
  unit: 'Isi jumlah unit (bilangan bulat, 0 atau lebih).',
} as const

/**
 * Sub-tab Loss Record (tiket 42) - `NB FacIn\Section\CauseOfLoss_FacIn.xml` (grid `.Property.ListCauseOfLoss`, kepala
 * sel 17-22), form `InputCauseOfLoss_FacIn.xml` (sel 3-11), dan `InputOfferFacInLossRatio.xml` (sel 3-6). Diuji
 * `labels.test.ts`.
 */
export const GRID_KERUGIAN = [
  { sel: '17', label: 'Date of Loss' },
  { sel: '18', label: 'Insured Name' },
  { sel: '19', label: 'Loss Object' },
  { sel: '20', label: 'Currency' },
  { sel: '21', label: 'Total of Loss' },
  { sel: '22', label: 'Total Claim' },
] as const

export const FORM_KERUGIAN = {
  insuredName: { sel: '3', label: 'Insured Name' },
  dateOfLoss: { sel: '4', label: 'Date of Loss' },
  lossObject: { sel: '5', label: 'Loss Object' },
  currency: { sel: '6', label: 'Currency' },
  claim: { sel: '7', label: 'Total Claim (100%)' },
  preventionOfLoss: { sel: '8', label: 'Prevention Of Loss' },
  causeOfLoss: { sel: '9', label: 'Cause of Loss' },
  detail: { sel: '11', label: 'Loss Detail' },
} as const

/**
 * Label medan `.Remarks` (sel 10, label dari aturan properti) - `[terverifikasi]` pyLabel "Remarks" aturan
 * `ASM-FW-GISFW-DATA-CAUSEOFLOSS!REMARKS` (`D:\migrasi\RNM\DDL\Remarks.xml`, ditambahkan work owner 03-10-2026).
 */
export const LABEL_REMARKS = 'Remarks'

/** Daftar `.Remarks` - PromptList aturan yang sama (tanpa baris kosong; nilai = label). Diuji `labels.test.ts`. */
export const OPSI_REMARKS = ['Settled', 'Ex Gratia Payment', 'Withdraw', 'Others', '--'].map((v) => ({ value: v, label: v }))

/** Loss ratio objek (`InputOfferFacInLossRatio`, baca-saja). Desimal tampilan = pxNumber Pega. */
export const LOSS_RATIO = [
  { sel: '3', label: 'LR 1 Year', kunci: 'oneYearAmount', desimal: 3 },
  { sel: '4', label: '%LR 1 Years', kunci: 'oneYearPercent', desimal: 2 },
  { sel: '5', label: 'LR 3 - 5 Years', kunci: 'threeFiveYearAmount', desimal: 3 },
  { sel: '6', label: '%LR 3 - 5 Years', kunci: 'threeFiveYearPercent', desimal: 2 },
] as const

/** Sub-tab Loss Record Internal (tiket 42) - grid `CauseOfLossClaim_FacIn.xml` (`.Property.ListCauseOfLossClaim`). */
export const GRID_KLAIM_INTERNAL = [
  { sel: '24', label: 'Year' },
  { sel: '25', label: 'Date of Loss' },
  { sel: '26', label: 'Location No' },
  { sel: '27', label: 'Location' },
  { sel: '28', label: 'Currency' },
  { sel: '29', label: 'Premium' },
  { sel: '30', label: 'O/S Claim' },
  { sel: '31', label: 'Acccepted Claim' },
  { sel: '32', label: 'Incurred Claim' },
  { sel: '33', label: 'Loss Ratio' },
  { sel: '34', label: 'Remark' },
] as const

/** Teks sistem baru Loss Record. */
export const TEKS_KERUGIAN = {
  uang: 'Isi angka (pemisah desimal titik).',
} as const

/**
 * Tab Coverage FIRE (tiket 43, tahap C1) - `NB FacIn\Section\InputInwardFacultativeDtl.xml` tab "Coverage" (visible
 * IsFire) -> `CoverageList.xml` (grid objek, sel 37-39) -> flow action `CoverageListFire_FlowAction` ->
 * `PropertyItemListCoverage.xml` (grid item sel 17-21, grid total per mata uang sel 51-54) -> flow action
 * `PropertyItemCoverageFacIn_FlowAction` -> `InputCoverageFire.xml` (grid coverage sel 16-18) -> flow action
 * `CoverageItem` (form, `CoverageItem.xml`). Tombol Save = Dtl `SaveFacIn_Act`. Diuji `labels.test.ts`.
 */
export const GRID_COV_OBJEK = [
  { sel: '37', label: 'No.' },
  { sel: '38', label: 'Object Name' },
  { sel: '39', label: 'Location' },
] as const

export const GRID_COV_ITEM = [
  { sel: '17', label: 'Object Item Type' },
  { sel: '18', label: 'Currency' },
  { sel: '19', label: 'TSI Object Item' },
  { sel: '20', label: 'Total Gross Premium' },
  { sel: '21', label: '‰ Total Net Rate' },
] as const

/** Grid total per mata uang (`.Property.TotalTSIPremiGrossList`). "Total Gross Rate" = permil (Pega salah memberi akhiran %). */
export const GRID_COV_TOTAL = [
  { sel: '51', label: 'Currency' },
  { sel: '52', label: 'Total TSI' },
  { sel: '53', label: 'Total Premi' },
  { sel: '54', label: 'Total Gross Rate' },
] as const

/**
 * Ringkasan seluruh coverage kasus (permintaan work owner 05-10-2026, gambar layar Pega) - `NB FacIn\Section\
 * SummaryCoverage_Section.xml`: grid per Object Item Type + Currency (`TempTotalItem.pxResults`) dan grid per Currency
 * (`TempTotal.pxResults`, Rate = `.TotalNetRate`). Diuji `TabCoverage.test.tsx`.
 */
export const RINGKASAN_COV_ITEM = ['Object Item Type', 'Currency', 'Total TSI', 'Total Premium'] as const
export const RINGKASAN_COV_MATA_UANG = ['Currency', 'Total TSI', 'Total Premium', 'Rate'] as const

export const GRID_COVERAGE = [
  { sel: '16', label: 'Coverage' },
  { sel: '17', label: '‰ Standard Rate' },
  { sel: '18', label: 'Premi' },
] as const

/** Form coverage (`CoverageItem.xml`, `pyLabelFieldValue`; tombol `pyLabel`). */
export const FORM_COV = {
  pilihCoverage: { sel: '4', tag: 'pyLabel', label: 'Choose Coverage' },
  coverage: { sel: '5', tag: 'pyLabelFieldValue', label: 'Coverage' },
  accumulationCode: { sel: '14', tag: 'pyLabelFieldValue', label: 'Accumulation Code' },
  accumulationAddress: { sel: '15', tag: 'pyLabelFieldValue', label: 'Accumulation Address' },
  pilihAkumulasi: { sel: '18', tag: 'pyLabel', label: 'Choose Accumulation Code' },
  salinAkumulasi: { sel: '19', tag: 'pyLabel', label: 'Copy Accumulation' },
  conditions: { sel: '20', tag: 'pyLabelFieldValue', label: 'Conditions' },
  day: { sel: '26', tag: 'pyLabelFieldValue', label: 'Days' },
  tsi: { sel: '27', tag: 'pyLabelFieldValue', label: 'TSI' },
  indemnity: { sel: '28', tag: 'pyLabelFieldValue', label: 'Indemnity' },
  rate: { sel: '29', tag: 'pyLabelFieldValue', label: '‰ Gross Rate' },
  firstLoss: { sel: '30', tag: 'pyLabelFieldValue', label: '% First Loss' },
  discountPercentage: { sel: '31', tag: 'pyLabelFieldValue', label: '% Discount' },
  tsiLiability: { sel: '32', tag: 'pyLabelFieldValue', label: 'First Loss' },
  netRate: { sel: '33', tag: 'pyLabelFieldValue', label: '‰ Net Rate' },
  limitOfLiability: { sel: '34', tag: 'pyLabelFieldValue', label: 'Limit of Liability' },
  pctLol: { sel: '35', tag: 'pyLabelFieldValue', label: '% Limit of Liability' },
  proRate: { sel: '39', tag: 'pyLabelFieldValue', label: '% Pro Rate' },
  unit: { sel: '40', tag: 'pyLabelFieldValue', label: 'Indemnity Unit' },
  indemnityPercentage: { sel: '41', tag: 'pyLabelFieldValue', label: '% Indemnity' },
  firstScale: { sel: '42', tag: 'pyLabelFieldValue', label: '% First Scale' },
  sublimit: { sel: '43', tag: 'pyLabelFieldValue', label: '% Sub Limit' },
  lostLimit: { sel: '44', tag: 'pyLabelFieldValue', label: '% Loss Limit' },
  emlPml: { sel: '45', tag: 'pyLabelFieldValue', label: '% EML/PML' },
  discount: { sel: '46', tag: 'pyLabelFieldValue', label: 'Discount' },
  premium: { sel: '47', tag: 'pyLabelFieldValue', label: 'Gross Premium' },
} as const

/**
 * Days coverage (sel 26) = RADIO BUTTON mendatar (`pxRadioButtons`, `pyOrientation` horizontal), bukan dropdown. Daftar =
 * PromptList aturan `ASM-FW-GISFW-DATA-COVERAGE!DAY` (`DDL\Day.xml`, dikirim work owner 04-10-2026; sama dengan gambar
 * layar). Ubah = `CountPremi_ACT` (percent / percent). Diuji `labels.test.ts`.
 */
export const PILIHAN_DAY_COVERAGE = ['365', '366', '360'] as const

/**
 * Indemnity Unit (sel 40, dropdown `pyHasNoSelection` false) - PromptList aturan `ASM-FW-GISFW-DATA-COVERAGE!UNIT`
 * (`DDL\Unit.xml`, 04-10-2026). Tanpa pilihan kosong: browser menampilkan pilihan pertama, sehingga nilai kosong = "0"
 * (gambar layar "Day (s)"; data contoh mayoritas "0"). Diuji `labels.test.ts`.
 */
export const OPSI_INDEMNITY_UNIT = [
  { value: '0', label: 'Day (s)' },
  { value: '1', label: 'Month (s)' },
  { value: '2', label: 'Year (s)' },
]

/** Accumulation Code / Address kosong tampil "---" (gambar layar Pega 03-10-2026). */
export const AKUMULASI_KOSONG = '---'

/** Label Coverage Basis = pyLabel aturan `ASM-FW-GISFW-DATA-COVERAGE!COVERAGEBASIS` (`DDL\CoverageBasis.xml`). */
export const LABEL_COVERAGE_BASIS = 'Coverage Basis'

/**
 * Pilihan Coverage Basis - PromptList `DDL\CoverageBasis.xml` (diuji). C1 menghitung basis 1-4; Layering (5) = tahap
 * berikut (C2; grid LayerList).
 */
export const OPSI_COVERAGE_BASIS = [
  { value: '1', label: 'Sum Insured Basis' },
  { value: '2', label: 'First Loss Basis' },
  { value: '3', label: 'EML / PML Basis' },
  { value: '4', label: 'Sub Limit Basis' },
  { value: '5', label: 'Layering Basis' },
]

/** Tombol Save tab Coverage (Dtl, `SaveFacIn_Act`). */
export const SIMPAN_COVERAGE = 'Save'

/** Teks sistem baru tab Coverage. */
export const TEKS_COVERAGE = {
  rateWajib: '‰ Gross Rate wajib diisi.',
  angka: 'Isi angka (pemisah desimal titik).',
  menghitung: 'Menghitung…',
  layeringBelum: 'Layering Basis belum tersedia (tahap berikut).',
  tanpaCoverage: 'Tidak ada coverage yang cocok.',
  cariCoverage: 'Search',
} as const

/**
 * Net rate item (tiket 44, tahap C2) - `NB FacIn\Activity\CekNetRate_ACT.xml`: ‰ Total Net Rate dapat diisi
 * (`FlagNetRate` "true") HANYA bila CoverageList item memuat kelima `.OLDID` ini (cocok harfiah); selain itu
 * TotalNetRate = 0. `CalculateNetRate_ACT`: NetRate coverage = Rate / ΣRate × TotalNetRate, lalu `CountPremi_ACT`.
 */
export const OLDID_NET_RATE = ['FLEXAS', '4.1A CC', '4.3', '4.2 PRGBI', 'OTHERS'] as const

/**
 * Deductible (tahap C3) - `[terverifikasi]` aturan properti `ASM-FW-GISFW-DATA-DEDUCTIBLE!MINMAX` (`DDL\MinMax.xml`) dan
 * `!CONDITION` (`DDL\Condition.xml`, ditambahkan work owner 03-10-2026). PromptList tanpa baris kosong. Diuji menurut
 * pxInsName.
 */
export const OPSI_MINMAX = [
  { value: '1', label: 'Min' },
  { value: '2', label: 'Max' },
  { value: '3', label: 'Or' },
]
export const OPSI_KONDISI_DEDUCTIBLE = [
  { value: '1', label: 'Any One Occurrence' },
  { value: '2', label: 'Any One Accident' },
  { value: '3', label: 'Each and Every Loss' },
  { value: '4', label: 'Each and Every Claim' },
  { value: '5', label: 'Other' },
]

/**
 * Deductible coverage FIRE (tiket 45, tahap C3) - grid `.DeductibleList` di `NB FacIn\Section\CoverageItem.xml`
 * (judul "Deductible"; kepala kolom Pega kosong kecuali "Time Excess (Days)") dan form flow action
 * `InputDtlDeductibleFire_FacIn` -> section `addDeductible.xml` (sel 5-14). Diuji `labels.test.ts`.
 */
export const JUDUL_DEDUCTIBLE = 'Deductible'
export const FORM_DEDUCTIBLE = {
  typeDeductible: { sel: '5', label: 'Type Deductible' },
  pctDeductible: { sel: '6', label: 'Pct Deductible' },
  minMax: { sel: '7', label: 'MinMax' },
  currency: { sel: '8', label: 'Currency' },
  typeDeductible2: { sel: '9', label: 'Type Deductible' },
  pctDeductible2: { sel: '10', label: 'Pct Deductible' },
  condition: { sel: '11', label: 'Condition' },
  amount: { sel: '12', label: 'Amount' },
  inputCondition: { sel: '13', label: 'Condition' },
  timeExcess: { sel: '14', label: 'Time Excess (In Days)' },
} as const
/** Kepala kolom grid "Time Excess (Days)" (CoverageItem.xml); kolom lain memakai label form (keputusan agent R-1). */
export const KOLOM_TIME_EXCESS = 'Time Excess (Days)'
/** Pilihan kosong dropdown deductible (`pyNoSelectionText` addDeductible). */
export const DEDUCTIBLE_KOSONG = 'Choose'

/** Type Deductible / Type Deductible 2 - aturan `DATA-DEDUCTIBLE!TYPEDEDUCTIBLE(2)` (`DDL\TypeDeductible*.xml`). */
export const OPSI_TYPE_DEDUCTIBLE = [
  { value: '1', label: '% Of Claim' },
  { value: '2', label: '% Of Loss' },
  { value: '3', label: '% Of TSI' },
  { value: '4', label: '% Of Approved Loss Value' },
  { value: '5', label: '% Of Recoverable Claim Amount' },
  { value: '6', label: '% Of Recoverable Amount' },
  { value: '7', label: 'In Amount' },
  { value: '0', label: 'NIL' },
]
export const OPSI_TYPE_DEDUCTIBLE2 = [
  { value: '1', label: '% Of Claim' },
  { value: '2', label: '% Of Loss' },
  { value: '3', label: '% Of TSI' },
  { value: '4', label: '% Of TSI Whichever Is Higher' },
  { value: '0', label: 'NIL' },
]

/**
 * Popup Choose Accumulation (tiket 46) - harness `ChooseAccumulation_FacIn` (judul jendela "Choose Accumulation") ->
 * section `SearchRiskAccumCov.xml`. Tahap 1 = jalur Report Definition `SearchRiskAccumulation_RD` (grid sel 90:
 * `.ID` / `.AccumulationName` / `.Note`). Diuji `labels.test.ts`.
 */
export const POPUP_AKUMULASI = {
  judul: 'Choose Accumulation',
  cari: { sel: '1', tag: 'pyValue', label: 'Search Risk Accumulation' },
  /** Judul bar tata letak saringan (`pyTitle`). */
  kondisi: 'Conditions',
  accumulationCode: { sel: '17', tag: 'pyLabelFieldValue', label: 'Accumulation Code' },
  policyNo: { sel: '20', tag: 'pyLabelFieldValue', label: 'Policy No' },
  road: { sel: '21', tag: 'pyLabelFieldValue', label: 'Road' },
  zipCode: { sel: '26', tag: 'pyLabelFieldValue', label: 'Zip Code' },
  country: { sel: '27', tag: 'pyLabelFieldValue', label: 'Country' },
  province: { sel: '28', tag: 'pyLabelFieldValue', label: 'Province' },
  accumType: { sel: '29', tag: 'pyLabelFieldValue', label: 'Accum. Type' },
  city: { sel: '32', tag: 'pyLabelFieldValue', label: 'City' },
  district: { sel: '33', tag: 'pyLabelFieldValue', label: 'District' },
  area: { sel: '34', tag: 'pyLabelFieldValue', label: 'Area' },
  czone: { sel: '35', tag: 'pyLabelFieldValue', label: 'CZone' },
  keyword: { sel: '36', tag: 'pyLabelFieldValue', label: 'Key Word' },
  filter: { sel: '40', tag: 'pyLabel', label: 'Filter' },
  bersih: { sel: '41', tag: 'pyLabel', label: 'Clear Column' },
  kolom: [
    { sel: '107', tag: 'pyValue', label: 'Accumulation Code' },
    { sel: '108', tag: 'pyValue', label: 'Type' },
    { sel: '109', tag: 'pyValue', label: 'Description' },
  ],
  pilih: { sel: '119', tag: 'pyLabel', label: 'Choose' },
} as const

/**
 * Key Word popup accumulation (sel 36, dropdown dengan pilihan kosong bertulisan kosong) - PromptList aturan
 * `ASM-FW-GISFW-INT-ACCUMULATION!KEYWORD` (`DDL\Keyword.xml`, 04-10-2026). Diuji `labels.test.ts`.
 */
export const OPSI_KEYWORD = [
  'APARTEMEN/RUMAH SUSUN',
  'BUSINESS CENTRE / OFFICE BUILDING',
  'DESA/DUSUN',
  'HOTEL/RESORT',
  'KAWASAN BERIKAT/INDUSTRI',
  'KAWASAN PERGUDANGAN',
  'KAWASAN WISATA',
  'KOMPLEKS/KAV.PERUMAHAN/ VILLA/PEMUKIMAN',
  'MAL/SHOPPING CENTRE/TRADE CENTRE',
  'PASAR TRADISIONAL',
  'RESTORAN',
  'RUKO/RUKAN',
  'RUMAH IBADAH',
  'RUMAH SAKIT',
  'STASIUN',
  'SUPERMARKET/HYPERMARKET',
  'TEMPAT KURSUS/SEKOLAH/AKADEMI/UNIVERSITAS',
  'TERMINAL',
  'ZONA LAIN-LAIN',
].map((k) => ({ value: k, label: k }))

export const TEKS_AKUMULASI = {
  tanpaHasil: 'Tidak ada accumulation yang cocok. Coba longgarkan kondisi pencarian.',
  /** Keterangan di bawah judul popup (tampilan dirapikan 04-10-2026, permintaan work owner). */
  sub: 'Isi satu atau beberapa kondisi, lalu tekan Filter. Accumulation yang dipilih harus ber-Zip sama dengan lokasi risiko.',
  zipLokasi: 'Zip lokasi risiko',
  zipKosong: 'belum diisi',
  kondisiAktif: (n: number) => `${n} kondisi aktif`,
  hasil: (n: number) => `${n} accumulation ditemukan`,
  petunjukAwal: 'Hasil pencarian tampil di sini.',
  pertama: 'Pertama',
  sebelumnya: 'Sebelumnya',
  berikutnya: 'Berikutnya',
  terakhir: 'Terakhir',
  navigasiHalaman: 'Navigasi halaman hasil',
  /** "Menampilkan 11–20 dari 35 data" (10 per halaman). */
  menampilkan: (halaman: number, ukuran: number, total: number) =>
    `Menampilkan ${(halaman - 1) * ukuran + 1}–${Math.min(halaman * ukuran, total)} dari ${total} data`,
  zipSesuai: 'Zip sesuai',
  zipBerbeda: 'Zip berbeda',
  kolomZip: 'Zip',
  /** Pesan `SetDataAccum_Act` langkah 5 (verbatim, termasuk dua spasi dan ejaan "ZIpCode"). */
  zipBeda: (risiko: string, dariId: string) => `ZIpCode Harus Sama dengan ZIpCode  Yang di Object Item >>>> ${risiko} != ${dariId}`,
} as const

/**
 * Add New accumulation (tiket 46, permintaan work owner 04-10-2026 dengan gambar layar) - tombol sel 39
 * `SearchRiskAccumCov` (pre-DT `InputAccumulationAdd_PreDT`, `GetCzone_Act`) membuka sub-section
 * `NB FacIn\Section\InputAccumulationCov.xml`; Choose Zip Code = flow action `ChooseZipCode` -> section
 * `ChooseZipCodeDtl` (RD `BrowseRiskAddressZipCode_RD`). Diuji `labels.test.ts`.
 */
export const TAMBAH_AKUMULASI = {
  tambahBaru: { sel: '39', tag: 'pyLabel', label: 'Add New' },
  accumulationType: { sel: '10', tag: 'pyLabelFieldValue', label: 'Accumulation Type' },
  province: { sel: '13', tag: 'pyLabelFieldValue', label: 'Province' },
  scopeArea: { sel: '16', tag: 'pyLabelFieldValue', label: 'Scope Area' },
  crestaZone: { sel: '17', tag: 'pyLabelFieldValue', label: 'Cresta Zone' },
  primaryZip: { sel: '20', tag: 'pyLabelFieldValue', label: 'Primary Zip Code' },
  pilihZip: { sel: '21', tag: 'pyLabel', label: 'Choose Zip Code' },
  keyword: { sel: '24', tag: 'pyLabelFieldValue', label: 'Key Word' },
  note: { sel: '25', tag: 'pyLabelFieldValue', label: 'Accumulation Description' },
  simpan: { sel: '28', tag: 'pyLabel', label: 'Save' },
  tutup: { sel: '29', tag: 'pyLabel', label: 'Close' },
  /** Kolom grid `ChooseZipCodeDtl` + tombol Choose sel 23. */
  kolomZip: ['Zip Code', 'City', 'Province', 'Nation'],
  pilih: 'Choose',
} as const

/** Pilihan kosong Key Word / Scope Area di form Add New (`pyNoSelectionText` "--"). */
export const PILIHAN_KOSONG_TAMBAH = '--'

/**
 * Scope Area (sel 16, dropdown wajib) - PromptList aturan `ASM-FW-GISFW-INT-ACCUMULATION!SCOPEAREA` (`DDL\ScopeArea.xml`,
 * dikirim work owner 04-10-2026). Diuji `labels.test.ts`.
 */
export const OPSI_SCOPE_AREA = ['AREA', 'DISTRICT', 'CITY', 'PROVINCE', 'COUNTRY'].map((v) => ({ value: v, label: v }))

export const TEKS_TAMBAH_AKUMULASI = {
  judul: 'Add New Accumulation',
  sub: 'Isi data accumulation baru. Primary Zip Code dipilih lewat tombol Choose Zip Code; Cresta Zone mengikuti zip.',
  wajib: 'Wajib diisi.',
  pilihDariSaran: 'Pilih dari daftar saran.',
  /** Pesan `SaveAccumulation_Act` langkah 10 (verbatim). */
  tidakLengkap: "Postal code, Nation, CZone and Accumulation Description can't be null!",
  /** Pesan prosedur `RDBMASTERACCUMULATION` bila Note + zip sudah ada (verbatim, tanpa tag HTML). */
  ganda: 'Error master item Akumulasi Sudah Ada',
  menyimpan: 'Menyimpan…',
  cariZip: 'Cari zip code',
  /** Kepala panel Choose Zip Code. */
  judulZip: 'Choose Zip Code',
  provinsiZip: 'Province',
  semuaProvinsi: 'Semua province',
  tutupZip: 'Tutup daftar zip code',
  petunjukZip: 'Daftar mengikuti Province yang dipilih di form. Pilih satu zip untuk mengisi Primary Zip Code dan Cresta Zone.',
  tanpaZip: 'Tidak ada zip code yang cocok.',
  zipBelum: '---',
} as const

/**
 * Tab Clauses kasus FIRE (tiket 47) - `NB FacIn\Section\InputInwardFacultativeDtl.xml` tab "Clauses" (`IsFire`) ->
 * `InputDtlClause_FacIn` (sel 209 Choose Clause; grid `.ClauseList` -> `InputClauseFire_FacIn` per baris); popup harness
 * `ChooseClauseFire` -> section `ChooseClauseFire`; local action `InputClauseFire_ViewDtl` ("Isi Klasula"). Diuji
 * `labels.test.ts`.
 */
export const KLAUSA = {
  pilihKlausa: { sel: '209', tag: 'pyLabel', label: 'Choose Clause' },
  clauseCode: { sel: '3', tag: 'pyLabelFieldValue', label: 'Clause Code' },
  totalArgumen: { sel: '5', tag: 'pyLabelFieldValue', label: 'Total Argument' },
  title: { sel: '6', tag: 'pyLabelFieldValue', label: 'Title' },
  description: { sel: '8', tag: 'pyLabelFieldValue', label: 'Description' },
  lihatArgumen: { sel: '9', tag: 'pyLabel', label: 'See Clause and Argument' },
  /** Popup ChooseClauseFire. */
  language: { sel: '3', tag: 'pyLabelFieldValue', label: 'Language' },
  keyword: { sel: '4', tag: 'pyLabelFieldValue', label: 'Key word' },
  cari: { sel: '7', tag: 'pyLabel', label: 'Search' },
  submit: { sel: '8', tag: 'pyLabel', label: 'Submit' },
  kolomPilih: ['Clause ID', 'Title', 'Description'],
  /** Local action InputClauseFire_ViewDtl. */
  judulArgumen: 'Isi Klasula',
  kolomArgumen: ['No. Argumen', 'Deskripsi Argumen', 'Isi Argumen'],
  simpan: 'Save',
} as const

/**
 * Language popup Choose Clause (`Clause.ClauseLanguageID`) - PromptList aturan `ASM-FW-GISFW-DATA-CLAUSE!CLAUSELANGUAGEID`
 * (`DDL\ClauseLanguageID.xml`, dikirim work owner 05-10-2026): kode 0 / 1 / 2. "Indonesia" terpilih di gambar layar
 * Pega DEV -> awal "0". Diuji `labels.test.ts`.
 */
export const OPSI_BAHASA_KLAUSA = [
  { value: '0', label: 'Indonesia' },
  { value: '1', label: 'Inggris' },
  { value: '2', label: 'Dual Bahasa' },
]
export const BAHASA_KLAUSA_AWAL = '0'

export const TEKS_KLAUSA = {
  judulPilih: 'Pilih Klausula Fire',
  kosong: 'Belum ada klausa. Tekan Choose Clause untuk menambah.',
  tanpaHasil: 'Tidak ada klausa yang cocok.',
  dipilih: (n: number) => `${n} klausa dipilih`,
  hapus: 'Hapus klausa',
  memuat: 'Memuat…',
  tersimpan: 'Klausa tersimpan.',
  sudahAda: 'Sudah ada di daftar klausa',
} as const

/**
 * Tab Spreading kasus FIRE (tiket 48) - `NB FacIn\Section\InputInwardFacultativeDtl.xml` tab "Spreading" (`!IsLife`):
 * sel 80 % Share RNM, blok "Copy Spreading", `CoverageSpreadingList` (lokasi) -> `PropertyItemListCoverageSpreading` (item +
 * total per lokasi) -> `InputCoverageSpreadingFire` (coverage) -> `SpreadingItem` + `InputDtlSpreadingCoverage_FacIn`
 * (baris spreading); `SummarySpreading_Section` (ringkasan). Gambar layar Pega DEV work owner 05-10-2026.
 */
export const SPREADING = {
  percentShare: '% Share RNM',
  judulCopy: 'Copy Spreading',
  kolomTemplate: ['Type Treaty', '% Share'],
  tambah: 'Add',
  hapus: 'Delete',
  salinSemua: 'Copy To All Spreading',
  kolomLokasi: ['No.', 'Object Name', 'Location'],
  kolomItem: ['Object Item Type', 'Currency', 'TSI Object Item', 'Total Gross Premium', 'Total Premium RNM'],
  kolomCoverage: ['Coverage', '‰ Standard Rate', 'Gross Premium', 'Premium RNM'],
  detailCoverage: {
    tsi: 'TSI 100%',
    tsiLiability: 'TSI Liability',
    tsiNusantaraRe: 'TSI Nusantara Re',
    premium: 'Gross Premium',
    premiNusantaraRe: 'Premium Nusantara Re',
  },
  kolomSpreading: ['Type Treaty', '% Share', 'TSI 100% Spreaded', 'TSI Spreaded (RNM)', 'Limit Of Liability', 'Premium Spreaded'],
  kolomTotal: [
    'Currency',
    'Treaty Name',
    '(%) Share',
    'Total TSI Top Risk Spreaded',
    'Total TSI Spreaded',
    'Total Limit of Liability',
    'Total Premium Spreaded',
  ],
  kolomRingkasanTreaty: [
    'Currency',
    'Treaty Name',
    '(%) Share',
    'Total TSI Top Risk Spreaded',
    'Total TSI Spreaded',
    'Total Limit of Liability',
    'Total First Loss',
    'Total Premium Spreaded',
  ],
  kolomRingkasanMataUang: ['Currency', 'Total TSI Top Risk Spreaded', 'Total TSI Spreaded', 'Total Limit of Liability', 'Total Premium Spreaded'],
  simpan: 'Save',
} as const

export const TEKS_SPREADING = {
  menghitung: 'Menghitung…',
  shareLebih: 'Total % Share Copy Spreading melebihi 100.',
  percentShareLebih: '% Share RNM paling besar 100.',
  pilihTreaty: '--',
  kosong: 'Belum ada spreading.',
  tersimpan: 'Spreading tersimpan.',
  angka: 'Isi angka (pemisah desimal titik).',
  /** Pesan Pega verbatim - `Activity\CalcultePersentageSpeading_Act.xml` langkah "check same treaty type". */
  treatySama: "Treaty Type can't be same",
} as const

/**
 * Tab Inw Fac Cedant Panels FIRE (tiket 49) - `NB FacIn\Section\InputInwardFacultativeDtl.xml` baris 60510-67588
 * (pyLabelFieldValue Source of Business / Share Cedant Type / % Share RNM / Total TSI RNM / Total Premi RNM; kepala grid
 * "Ceding" / "% Share"; tombol "Add" / "Delete"); "Choose Ceding" = pyLabel + pyWindowName `Section\CedingCedant.xml`.
 * Diuji `TabCedant.test.tsx`.
 */
export const CEDANT = {
  sob: 'Source of Business',
  shareCedantType: 'Share Cedant Type',
  percentShare: '% Share RNM',
  totalTsi: 'Total TSI RNM',
  totalPremi: 'Total Premi RNM',
  kolomCeding: 'Ceding',
  kolomShare: '% Share',
  tambah: 'Add',
  hapus: 'Delete',
  pilihCeding: 'Choose Ceding',
  judulPilih: 'Choose Ceding',
} as const

/** `DDL\ShareCedantType.xml` (aturan properti `ASM-FW-GISFW-DATA-OFFERFACIN!SHARECEDANTTYPE`, pyPromptTableList). */
export const OPSI_SHARE_CEDANT_TYPE = [
  { value: '0', label: 'Gross' },
  { value: '1', label: 'Share RNM' },
] as const

export const TEKS_CEDANT = {
  kosong: 'No items',
  simpan: 'Save',
  menyimpan: 'Menyimpan…',
  tersimpan: 'Cedant tersimpan.',
  typeWajib: 'Share Cedant Type wajib dipilih.',
  /** Pesan Pega verbatim - `Activity\ProtectShareCedant_Act.xml`. */
  shareTidakSah: '% Share cedant more than 100 or less than 0!',
  cedingKosong: "Ceding on list inward facultative cedant can't be null!",
  listKosong: "List Inward facultative cedant can't be null!",
  /** Pesan Pega verbatim - `Activity\SetTSIPremiCedant_Act.xml` (`.ShareCeding>100`); berlaku TANPA gerbang. */
  nilaiLebih: "Value can't be more than 100 or less than 0!",
} as const
