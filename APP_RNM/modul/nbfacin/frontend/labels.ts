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
  /** Untuk `BelumTersedia` (keputusan agent B-2): daftar belum punya sumber data. */
  daftar: 'Daftar opportunity',
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
  /** Pencarian berhasil tetapi tidak ada baris yang cocok. */
  tanpaAccount: 'Tidak ada account yang cocok.',
} as const
