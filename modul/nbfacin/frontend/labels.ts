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
  // Sel 95 berjudul kosong di Pega; diberi judul "NB Status" atas permintaan work owner 07-10-2026 ("disini ditambahin NB
  // status dong, biar tau nb nya lgi dimana"). Isi = `.NBStatus` (T_GENERAL_POLIS.NB_STATUS; kosong = "NEW NB", awal flow).
  { sel: '95', label: 'NB Status', properti: 'A.NBStatus' },
  { sel: '96', label: 'Status', properti: 'A.NBStatusNew' },
] as const

/**
 * Kolom daftar kotak masuk Beranda NB FacIn = gambar layar inbox "New Business" Pega dari work owner 07-10-2026 (report-nya
 * tidak ada di korpus). Judul kosong = kolom tanpa judul di gambar yang berisi "Accept by <PIC>".
 */
export const KOLOM_INBOX_NBFACIN = [
  { kunci: 'id', label: 'ID' },
  { kunci: 'tertanggung', label: 'Insured Name' },
  { kunci: 'bisnis', label: 'Business name' },
  { kunci: 'sob', label: 'Source Of Business' },
  { kunci: 'ceding', label: 'Ceding Company' },
  { kunci: 'mulai', label: 'Inception Date' },
  { kunci: 'share', label: '%Share RNM' },
  { kunci: 'marketing', label: 'Marketing Name' },
  { kunci: 'status', label: 'Status' },
  { kunci: 'dibuat', label: 'Created date' },
  { kunci: 'penawaran', label: 'Offering Date' },
  { kunci: 'validitas', label: 'Validity Date' },
  { kunci: 'statusInbox', label: 'Status Inbox' },
  { kunci: 'aksiTerakhir', label: '' },
  { kunci: 'sejak', label: 'Time Since Last Update' },
] as const

/**
 * Kolom tambahan kotak masuk workbasket Underwriting (gambar layar inbox "New Business" Pega di underwriter, work owner
 * 07-10-2026): "Time in Underwriting" sesudah Validity Date; tampil untuk workbasket `WORKBASKET_UW_NBFACIN`.
 */
export const KOLOM_INBOX_UW_NBFACIN = { kunci: 'diUw', label: 'Time in Underwriting' } as const

/** Workbasket berkolom "Time in Underwriting" (daftar work owner 07-10-2026). */
export const WORKBASKET_UW_NBFACIN = [
  'ReasFacInJuniorUnderwritingB',
  'ReasFacInJuniorUnderwritingA',
  'ReasFacInUnderwriting',
  'ReasFacInSeniorUnderwriting',
  'ReasFacInDepHeadUnderwriting',
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

/**
 * Tiga tombol Group Business di gambar. Dulu dua tombol terakhir nonaktif karena aksinya tidak diketahui (keputusan agent
 * C-2). Sejak 10-10-2026 work owner memutuskan "BUTTON NEW COMPANY DETAIL = ISI MENU COMPANY DETAIL, BUTTON NEW GROUP
 * BUSINESS = ISI MENU ACCOUNT", lihat `MENU_TOMBOL_OPPORTUNITY`.
 */
export const TOMBOL_FORM_OPPORTUNITY = {
  cariGrup: 'Search Group Business',
  perusahaanBaru: 'New Company Detail',
  grupBaru: 'New Group Business',
} as const

/** Modul yang menunya dibuka tombol New Company Detail / New Group Business (jendela `PropsRute.onBukaMenu`). */
export const MENU_TOMBOL_OPPORTUNITY = {
  perusahaanBaru: { modul: 'companydetail', menu: 'Company Detail' },
  grupBaru: { modul: 'accounts', menu: 'Accounts' },
} as const

/** Teks sistem baru - BUKAN dari Pega. */
export const TEKS_MENU_OPPORTUNITY = {
  /** Akun tidak memegang menu modul itu (atau modulnya tidak aktif). */
  tidakAda: (menu: string) => `Menu ${menu} tidak tersedia untuk akun Anda.`,
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
/**
 * Bilah CSV objek (tiket 66). `pleaseUpload` verbatim `NB FacIn\Activity\InsertUploadFire_act.xml` b.492; "Success Insert Data"
 * (b.7577) datang dari backend. Selebihnya teks sistem baru.
 */
export const TEKS_CSV_OBJEK = {
  pleaseUpload: 'Please Upload CSV',
  terbaca: (n: number) => `${n} baris CSV terbaca - periksa lewat View Upload, lalu tekan Save Data untuk menyimpan.`,
  konfirmasiGanti: 'Case ini sudah memiliki objek. Save Data akan MENGGANTI semua objek / lokasi dengan isi CSV. Lanjutkan?',
  menyimpan: 'Menyimpan…',
  ya: 'Ya',
  batal: 'Batal',
  tutup: 'Tutup',
  // Tiket 75 (unggah tanpa batas, proses latar) - teks sistem baru.
  mengunggah: (nama: string, mb: string) => `Mengunggah ${nama} (${mb} MB)… jangan tutup halaman sampai unggahan selesai.`,
  memvalidasi: (n: number) => `Memeriksa isi CSV di server… ${n.toLocaleString('id-ID')} baris diperiksa.`,
  valid: (n: number) => `${n.toLocaleString('id-ID')} baris CSV valid - periksa lewat View Upload, lalu tekan Save Data untuk menyimpan.`,
  gagalValidasi: (n: number) => `${n.toLocaleString('id-ID')} baris CSV tidak valid - perbaiki berkas lalu unggah ulang:`,
  menyimpanKemajuan: (a: number, b: number) =>
    `Menyimpan objek di server… ${a.toLocaleString('id-ID')} / ${b.toLocaleString('id-ID')} baris (${b > 0 ? Math.floor((a / b) * 100) : 0}%). Halaman boleh ditinggal; kemajuan tetap berjalan.`,
  gagal: (a: number, b: number) =>
    `Penyimpanan berhenti sesudah ${a.toLocaleString('id-ID')} / ${b.toLocaleString('id-ID')} baris. Tekan Ulangi untuk melanjutkan dari bagian yang gagal.`,
  ulangi: 'Ulangi',
  galatLain: (n: number) => `… dan ${n.toLocaleString('id-ID')} galat lain.`,
  /** Bug DEV 11-10-2026: template CSV FIRE hanya untuk case FIRE (Pega ChooseActivityMappingCSV_Act L439/L499). */
  belumUntukCob: (cob: string) => `Upload CSV untuk Class of Business ${cob} belum tersedia di sistem baru (template CSV FIRE hanya untuk Class of Business FIRE).`,
} as const

export const PILIHAN_PERIODE = {
  /**
   * Policy Type: disimpan STANDARD VALUE (`D:\migrasi\RNM\DDL\PolicyType.xml` `ASM-FW-GISFW-DATA-QUOTATION!POLICYTYPE`: 0 Individual
   * Policy, 1 Master Policy, 2 Declaration Policy) - work owner 09-10-2026 "YANG DISIMPAN ITU STANDARD VALUENYA. ( 0,1,2)",
   * ralat butir 78.2 ("Label text for now"). Layar tetap dua pilihan (gambar layar Pega).
   */
  policyType: [
    { value: '0', label: 'Individual Policy' },
    { value: '1', label: 'Master Policy' },
  ],
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
  /**
   * Isi tab untuk Class of Business selain FIRE (work owner 09-10-2026 "ini kenapa? saya baru buat cob liability"): layar Pega
   * keluarga Aneka (Liability dkk., 36 section) belum diport - bukan galat pemuatan. Teks sistem baru.
   */
  cobBelum: (tab: string, cob: string) => `Tab ${tab} untuk Class of Business ${cob} belum tersedia di sistem baru (saat ini baru FIRE).`,
  /** Sesudah Save for later berhasil. */
  tersimpan: 'Tersimpan.',
  /** Case baca-saja: sudah keluar dari tahap marketing (work owner 07-10-2026). */
  bacaSaja: 'Case ini baca-saja: sudah tidak berada di tahap marketing, jadi isinya tidak dapat diubah.',
  /** Tahap Underwriting (When IsUWSpread, tiket 59): hanya tab Spreading yang dapat diubah. Teks sistem baru. */
  bacaSajaKecualiSpreading: 'Case berada di tahap Underwriting: hanya tab Spreading (% Share RNM dan spreading) yang dapat diubah.',
  /** Tahap Underwriting dengan Accumulation terbuka (tiket 60). Teks sistem baru. */
  bacaSajaKecualiSpreadingAkumulasi:
    'Case berada di tahap Underwriting: hanya tab Spreading (% Share RNM dan spreading) dan Accumulation di tab Coverage yang dapat diubah.',
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

/**
 * Grid total per mata uang (`.Property.TotalTSIPremiGrossList`). "Total Gross Rate" = Premium / TSI seperti layar Pega DEV
 * (work owner 05-10-2026; korpus menulis x 1000 - beda versi server).
 */
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
/** Kepala kolom grid "Time Excess (Days)" (CoverageItem.xml); kolom lain berkepala kosong seperti Pega (work owner 08-10-2026). */
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
  /** Pesan Pega verbatim - `Activity\cekSpreadingFactIn.xml` (Σ % Share baris spreading satu coverage > 100; tiket 59). */
  shareBarisLebih: "Share Percentage can't be more than 100",
  /** Grid spreading tahap UW: baris tanpa Type Treaty / % Share (teks sistem baru, bukan korpus). */
  barisBelumLengkap: 'Lengkapi Type Treaty dan % Share setiap baris spreading.',
} as const

/**
 * Tab Inw Fac Cedant Panels FIRE (tiket 49) - `NB FacIn\Section\InputInwardFacultativeDtl.xml` baris 60510-67588
 * (pyLabelFieldValue Source of Business / Share Cedant Type / % Share RNM / Total TSI RNM / Total Premi RNM; kepala grid
 * "Ceding" / "% Share"; tombol "Add" / "Delete").
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
  /** Detail baris = `Section\CedingCedant.xml` (gambar layar work owner 05-10-2026). XML "OF", layar Pega "of". */
  dari: 'of',
  kolomMataUang: ['Currency', 'TSI Cedant', 'Premium Cedant'],
  bukaBaris: 'Buka / tutup detail cedant',
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
  typeWajib: 'Share Cedant Type wajib diisi.',
  /** Pesan Pega verbatim - `Activity\ProtectShareCedant_Act.xml`. */
  shareTidakSah: '% Share cedant more than 100 or less than 0!',
  cedingKosong: "Ceding on list inward facultative cedant can't be null!",
  listKosong: "List Inward facultative cedant can't be null!",
  /** Pesan Pega verbatim - `Activity\SetTSIPremiCedant_Act.xml` (`.ShareCeding>100`); berlaku TANPA gerbang. */
  nilaiLebih: "Value can't be more than 100 or less than 0!",
  sudahSemua: 'Semua Ceding Co tab General sudah ada di daftar.',
  /** Total RNM dalam IDR (C-8). */
  mataUangTotal: 'IDR',
  kursTidakAda: 'Kurs ke IDR belum ada untuk mata uang: ',
  /** Pesan agent (C-7) - bukan teks Pega. */
  totalLebih: 'Total % Share cedant tidak boleh lebih dari 100.',
  generalKosong: 'Ceding Co di tab General belum diisi.',
} as const

/**
 * Tab Premium Deduction / Brokerage Fee FIRE (tiket 50) - `NB FacIn\Section\InputInwardFacultativeDtl.xml` baris
 * 74726-90306 (pyLabelFieldValue / pyValue; "Exclude PPN ?" = pyCheckboxCaption); daftar komisi =
 * `CoverageCommisionList` / `PropertyItemListCoverageCommision` / `InputCoverageCommisionFire` / `CommisionItem`.
 * Kepala PPh / PPN = `SetTAXNote` ("PPH " + tarif + " %"). Diuji `TabDeduction.test.tsx`.
 */
export const DEDUCTION = {
  pctPremiumDeduction: '% Premium Deduction',
  pctBrokerageFee: '% Brokerage Fee',
  excludePpn: 'Exclude PPN ?',
  kolomMataUang: 'Currency',
  totalPremiumDeduction: 'Total Premium Deduction',
  totalPremiumDeduction2: 'Total Premium Deduction 2',
  totalBrokerageFee: 'Total Brokerage Fee',
  pph: 'PPH',
  ppn: 'PPN',
  kolomLokasi: ['No.', 'Object Name', 'Location'],
  kolomItem: ['Object Item Type', 'Currency', 'TSI Object Item', 'Total Gross Premium', 'Total Premium RNM'],
  kolomCoverage: ['Coverage', '‰ Standard Rate', 'Gross Premium', 'Premium RNM'],
  detailCoverage: {
    tsiLiability: 'Gross TSI Liability',
    premium: 'Gross Premium',
    tsiNusantaraRe: 'TSI Nusantara Re',
    premiNusantaraRe: 'Premium Nusantara Re',
  },
  /** Grid spreading di detail coverage = `InputDtlSpreadingCommisionCoverage_FacIn` (.TreatyType, .SharePercentage,
   *  .TSISpreaded, .ClaimEstimation, .PremiumSpreaded) - gambar layar Pega work owner 06-10-2026. */
  kolomSpreading: ['Type Treaty', '% Share', 'TSI Spreaded', 'Limit Of Liability', 'Premium Spreaded'],
} as const

export const TEKS_DEDUCTION = {
  kosong: 'No items',
  simpan: 'Save',
  menyimpan: 'Menyimpan…',
  menghitung: 'Menghitung…',
  tersimpan: 'Premium deduction tersimpan.',
  bukaBaris: 'Buka / tutup detail',
  /** Pesan Pega verbatim - `Activity\ProtekDeduction2_Act.xml` langkah 2. */
  deduction2Lebih: 'Deduction 2 cannot be greater than the Premium',
} as const

/**
 * Tab Payment FIRE (tiket 51) - `NB FacIn\Section\InputInwardFacultativeDtl.xml` (label "Payment Information"),
 * `PaymentCurrencyList` (Installment; Currency / Premium RNM / Total Payment) dan `InputDtlPayment_FacIn` (grid
 * "Installment Data Information"). Kepala PPH / PPN = `PPH_Note` / `PPN_Note`. Diuji `TabPayment.test.tsx`.
 */
export const PAYMENT = {
  judul: 'Payment Information',
  installment: 'Installment',
  kolomMataUang: ['Currency', 'Premium RNM', 'Total Payment'],
  judulCicilan: 'Installment Data Information',
  no: 'No',
  dueDate: 'Due Date',
  persen: '% Installment',
  premium: 'Premium RNM',
  deduction: 'Deduction',
  deduction2: 'Deduction 2',
  brokerage: 'Total Brokerage Fee',
  pph: 'PPH',
  ppn: 'PPN',
  discount: 'Discount',
  totalPayment: 'Total Payment',
} as const

export const TEKS_PAYMENT = {
  kosong: 'No items',
  simpan: 'Save',
  menyimpan: 'Menyimpan…',
  menghitung: 'Menghitung…',
  tersimpan: 'Payment tersimpan.',
  bukaBaris: 'Buka / tutup cicilan',
  installmentTidakSah: 'Installment diisi bilangan bulat 1–99.',
  /** Pesan Pega verbatim - `Endorsment Fac In\Activity\SetValidateInstallment_Act.xml` (versi kelas Currency). */
  persenLebih: "Total % Installment Can't Be More Than 100%",
  persenKurang: 'Total % Installment Must Be 100%',
} as const

/**
 * Zone / 4.2 Construction di form coverage (tiket 52) - `NB FacIn\Section\CoverageItem.xml` sel 6/11 ("Zone"), 12 (ikon,
 * pyTooltip "Switch to edit mode"), 13 ("4.2 Construction"). Kode coverage = syarat tampil sel (L4252, L6983).
 */
export const ZONA_COVERAGE = {
  zone: 'Zone',
  konstruksi: '4.2 Construction',
  ubahZona: 'Switch to edit mode',
} as const

/** Coverage yang menampilkan Zone: `.Coverage` 100835 / 100821 / 100829 / 100816 / 100817 / 100828 (sel 6). */
export const COVERAGE_ZONA: readonly string[] = ['100835', '100821', '100829', '100816', '100817', '100828']

/** Coverage yang menampilkan 4.2 Construction: `.Coverage` 100821 / 100829 (sel 13). */
export const COVERAGE_KONSTRUKSI: readonly string[] = ['100821', '100829']

/**
 * Tombol salin di bawah grid coverage item (tiket 53) - pyLabel `NB FacIn\Section\InputCoverageFire.xml`; popup =
 * `Section\CopyCoverageFrom.xml` (pyLabelFieldValue "Location Address" / "Object Item Name", tombol "Copy to Another
 * Object Item"). Diuji `salinCoverage.test.ts`.
 */
export const SALIN_COVERAGE = {
  coverage: 'Copy Coverage',
  semua: 'Copy Coverage To All',
  dari: 'Copy Coverage From',
  deductibleSemua: 'Copy Deductible To All',
  deductibleBI: 'Copy Deductible For BI',
  keItemLain: 'Copy to Another Object Item',
  lokasi: 'Location Address',
  item: 'Object Item Name',
  objek: 'Objek',
  menyalin: 'Menyalin dan menghitung ulang premi…',
  // Umpan balik Copy Deductible (work owner 08-10-2026 "dicoba tidak ada berfungsi"; Pega tanpa pesan - teks keputusan agent).
  deductibleKosong: 'Coverage item ini belum punya deductible - deductible coverage berkode sama di item lain dikosongkan.',
  deductibleTanpaPadanan: 'Tidak ada coverage di item lain yang kodenya sama dengan coverage item ini, jadi tidak ada deductible yang disalin.',
  deductibleTersalin: (coverage: number, item: number) =>
    `Deductible disalin ke ${coverage} coverage di ${item} item lain. Klik Save untuk menyimpan.`,
  // CopyDeductibleToAllLocation_Act "jangan copy ke Business Interuption 10000992"; tombol For BI hanya di item BI lokasi 1.
  deductibleBiDilewati: (n: number) =>
    `${n} item Business Interruption tidak ikut Copy Deductible To All (aturan Pega). Isi deductible di item Business Interruption, lalu pakai Copy Deductible For BI di item itu untuk menyalinnya ke item Business Interruption lain.`,
} as const

/**
 * Tab Scoring Risk FIRE (tiket 54) - teks verbatim `NB FacIn\Section\ScoringRisk.xml` (kepala, hasil, persetujuan) dan
 * `ScoringRiskForm1` / `4` / `2` / `3` (faktor, bobot, sub-faktor, pilihan ChechBox1..n berurutan; spasi ganda Pega
 * dipertahankan). Kunci sub-faktor = kontrak `…/scoring`. Diuji `TabScoring.test.tsx`.
 */
export const SCORING = {
  judul: 'Risk Scoring - Property',
  insured: 'The Insured',
  topRisk: 'Top Risk',
  dominantRisk: 'Dominant Risk',
  periode: 'Period',
  lokasi1: 'Risk Location 1 (Top Risk)',
  lokasi2: 'Risk Location 2 (Dominant Risk)',
  status: 'Status',
  skor: 'Score',
  kolom: ['No', 'Scoring Factor', 'Percentage', 'Risk Detail', 'Top Risk', 'Dominant Risk'],
  total: 'Total',
  hasil: 'Scoring Result',
  klasifikasi: 'Risk Classifications',
  kelas: [
    { rentang: '86 - 100', nama: 'Good' },
    { rentang: '76 - 85', nama: 'Above Average' },
    { rentang: '60 - 75', nama: 'Average' },
    { rentang: '30 - 59', nama: 'Below Average' },
    { rentang: '0 - 29', nama: 'Poor' },
  ],
  rekomendasi: 'Recomendation',
  persetujuan: 'Required approval :',
  dari: 'From :',
  karena: 'Due to :',
  lossRatioInput: 'Loss Ratio 1 up to 5 Years',
} as const

/** Satu sub-faktor: kunci kontrak, label Risk Detail, bobot, pilihan ChechBox1..n; `terkunci` = nomor pilihan isian sistem. */
export interface SubFaktorScoring {
  kunci: string
  label: string
  bobot: string
  opsi: readonly string[]
  terkunci?: readonly number[]
}

export interface FaktorScoring {
  no: string
  faktor: string
  bobot: string
  sub: readonly SubFaktorScoring[]
}

export const FAKTOR_SCORING: readonly FaktorScoring[] = [
  {
    no: '1', faktor: 'Occupation', bobot: '15,00',
    sub: [{ kunci: 'occupation', label: '', bobot: '', opsi: ['Low Risks', 'Medium Risks', 'High Risks', 'Non-Preferred Risk', 'Treaty Exclusion'], terkunci: [1, 2, 3, 4, 5] }],
  },
  {
    no: '2', faktor: 'Risk Location', bobot: '10,00',
    sub: [
      { kunci: 'riot', label: 'Riot', bobot: '1', opsi: ['Cover exclude Riot', 'Non - Riot prone area', 'Riot prone area', 'No Information'] },
      { kunci: 'flood', label: 'Flood', bobot: '5', opsi: ['Cover exclude Flood', 'Non - Flood prone area', 'Flood prone area', 'No Information'] },
      {
        kunci: 'earthquake', label: 'Earthquake', bobot: '2,5',
        opsi: ['Cover exclude Earthquake', 'Zone 1', 'Zone 2', 'Zone 3', 'Zone 4', 'Zone 5', 'Other than zona 1 up to zona 5'],
      },
      { kunci: 'tsunami', label: 'Tsunami', bobot: '1,5', opsi: ['Cover exclude Tsunami', 'Non - Tsunami prone area', 'Tsunami prone area', 'No Information'] },
    ],
  },
  {
    no: '3', faktor: 'Object Conditions', bobot: '12,50',
    sub: [
      { kunci: 'classConstruction', label: 'Class Construction', bobot: '2,5', opsi: ['I', 'II', 'III', 'Growing Trees', 'Open storage', 'No Information'] },
      {
        kunci: 'building', label: 'Building', bobot: '2,5',
        opsi: ['Good condition / housekeeping', 'Growing Trees / Open storage', 'Bad Condition /  housekeeping', 'Silent Risk', 'No Information'],
      },
      {
        kunci: 'machinery', label: 'Machinery', bobot: '2,5',
        opsi: ['Good condition / maintenance', 'Growing Trees / Open storage', 'Bad condition / maintenance', 'Silent Risk', 'No Information'],
      },
      {
        kunci: 'stock', label: 'Stock in the policy cover', bobot: '2,5',
        opsi: ['Policy does not cover stock', 'Sum Insured < 60%', 'Sum Insured Stock >= 60%', 'No Information'],
      },
      {
        kunci: 'typeOfStock', label: 'Type of Stock', bobot: '2,5',
        opsi: ['No Stock in the premises', 'Non Combustible or Non Perishable Goods', 'Combustible or Perishable Goods or obsolete', 'No Information'],
      },
    ],
  },
  {
    no: '4', faktor: 'FEA (Fire Extinguishing Appliance)', bobot: '10,00',
    sub: [
      {
        kunci: 'portableFireExtinguisher', label: 'Portable Fire Extinguisher', bobot: '3',
        opsi: ['Available', 'Available and  adequate', 'Not Available (Growing Trees)', 'Not available or No Information'],
      },
      { kunci: 'hydrant', label: 'Hydrant System', bobot: '1,5', opsi: ['Available', 'Not Available (Growing Trees)', 'Not available or No Information'] },
      {
        kunci: 'sprinkler', label: 'Sprinkler System', bobot: '1',
        opsi: ['Available', 'Not Available (Growing Trees, Open Storage)', 'Not available or No Information'],
      },
      { kunci: 'fireAlarm', label: 'Fire Alarm System', bobot: '1,5', opsi: ['Available', 'Not available or No Information'] },
      {
        kunci: 'securityGuard', label: 'Security Guard', bobot: '1,5',
        opsi: ['Available in the premises', 'Available on outside premises', 'Not available or No Information'],
      },
      {
        kunci: 'fireBrigade', label: 'Fire brigade', bobot: '1,5',
        opsi: [
          'Available in the premises',
          'Available on outside within < 30 minutes',
          'Available on outsidewithin 30-60 minutes',
          'Available on outsidewithin  > 60 minutes',
          'Not available or No Information',
        ],
      },
    ],
  },
  {
    no: '5', faktor: 'Net Loss Ratio', bobot: '22,50',
    sub: [
      {
        kunci: 'lossRatio', label: 'Loss Ratio (Claim/net premium)', bobot: '',
        opsi: [
          'Nil 5 years and more',
          'Nil 3 years',
          'Nil 1 years',
          '-- Loss Ratio > 0 up to 10%',
          '-- Loss Ratio > 10 up to 25%',
          '-- Loss Ratio > 25 up to 50%',
          '-- Loss Ratio > 50 up to 75%',
          '-- Loss Ratio> 75 up to 100%',
          '-- Loss Ratio > 100%',
          'No Information',
        ],
        terkunci: [4, 5, 6, 7, 8, 9],
      },
    ],
  },
  {
    no: '6', faktor: 'Risk Improvement', bobot: '10,00',
    sub: [
      {
        kunci: 'riskImprovement', label: '', bobot: '',
        opsi: [
          'Risk Improvement has been done or on progress',
          'Risk Improvement is not required due to nil loss record',
          'No risk improvement for incurred loss or No Information',
        ],
      },
    ],
  },
  {
    no: '7', faktor: 'Survey Report', bobot: '10,00',
    sub: [
      {
        kunci: 'surveyReport', label: '', bobot: '',
        opsi: [
          'Available - within 1 year',
          'Available - 2 years ago',
          'Available - 3 years ago',
          'Available - more than 3 years ago',
          'No Survey Report, but provided photographs or summary of the insured activities',
          'Not Available or No Information',
        ],
      },
    ],
  },
  {
    no: '8', faktor: 'Others', bobot: '10,00',
    sub: [
      { kunci: 'rate', label: 'Rate', bobot: '2', opsi: ['On Gross Rate (OGR)  standard', 'Market Rate', 'No Information'] },
      { kunci: 'deductible', label: 'Deductible', bobot: '2,5', opsi: ['Standard or higher than standard', 'Non Standard', 'No Information'] },
      {
        kunci: 'totalSumInsured', label: 'Total Sum Insured', bobot: '1',
        opsi: ['Less than IDR 30,000,000,000', 'IDR 30,000,000,000 or more'],
        terkunci: [1, 2],
      },
      {
        kunci: 'otherInformation', label: 'Other Information', bobot: '2,5',
        opsi: [
          'Available and not potentially to increase risk exposure',
          'Available and  potentially to increase risk exposure',
          'Not Available or No Information',
        ],
      },
      {
        kunci: 'clauses', label: 'Clauses', bobot: '2',
        opsi: ['Non Standard (un usual) Clauses', 'Standard (usual) Clauses', 'Not Available or No Information'],
      },
    ],
  },
]

/** Rekomendasi dan persetujuan (`ScoringRisk.xml`, `.ScoringRisk.ScoringResult.*`). */
export const HASIL_SCORING = {
  rekomendasi: [
    { kunci: 'acceptableNotApproval', label: 'Acceptable, not required approval' },
    { kunci: 'acceptableWithApproval', label: 'Acceptable, required approval' },
    { kunci: 'notAcceptable', label: 'Not Acceptable' },
  ],
  dari: [
    { kunci: 'opDivHead', label: 'Operational Div. Head' },
    { kunci: 'opDirector', label: 'Operational Director' },
    { kunci: 'techDirector', label: 'Technical Director' },
    { kunci: 'presDirector', label: 'President Director' },
  ],
  karena: [
    { kunci: 'dueOccupation', label: 'Occupation' },
    { kunci: 'dueObjectConditions', label: 'Object Conditions' },
    { kunci: 'dueLossRatio', label: 'Loss Ratio' },
    { kunci: 'dueOthers', label: 'Others' },
  ],
} as const

export const TEKS_SCORING = {
  bacaSaja: 'Scoring Risk hanya dapat diisi pada tahap marketing (case belum diajukan).',
  menghitung: 'Menghitung…',
  simpan: 'Save',
  menyimpan: 'Menyimpan…',
  tersimpan: 'Scoring Risk tersimpan.',
  kosong: '—',
  diisiSistem: 'Diisi sistem',
  terkunciGrowingTrees: 'Terkunci: Class Construction = Growing Trees',
  tanpaDominan: 'Pilih lokasi Dominant Risk di atas untuk mengisi kolom Dominant Risk.',
  belumDipilih: 'Belum dipilih',
} as const

/**
 * Tab utama Loss Record FIRE (tiket 55) - teks verbatim `NB FacIn\Section\ObjectList_LossRecord.xml` (bar unggah, premi,
 * grid), `SummaryLossRecord_Section.xml` (total, loss ratio) dan `InputLossRecord_Sec.xml` (ubah baris). Diuji
 * `TabLossRecord.test.tsx`.
 */
export const LOSS_RECORD = {
  judulUnggah: 'Please upload file with .csv format',
  unduhTemplate: 'Download Template CSV',
  unggah: 'Upload CSV',
  simpan: 'Save Data',
  /** Tombol Save tab (`InputInwardFacultativeDtl.xml` kontainer Loss Record, layout S89 sel 525). */
  simpanTab: 'Save',
  premiPar: 'Premi PAR',
  premiMb: 'Premi MB',
  premiPl: 'Premi PL',
  premiOther: 'Premi Other COB',
  uwYear: 'Uw Year',
  classOfBusiness: 'Class Of Business',
  dateOfLoss: 'Date of Loss',
  causeOfLoss: 'Cause Of Loss',
  causeOfLossUbah: 'Cause of Loss',
  currency: 'Currency',
  lossAmount: 'Loss Amount',
  lossAmountIdr: 'Loss Amount ( IDR )',
  insuredName: 'Insured Name',
  totalLoss: 'Total Loss Amount',
  totalLossIdr: 'Total Loss Amount ( IDR )',
  rasio: ['Loss Ratio 1 Year', 'Loss Ratio 3 Year', 'Loss Ratio 5 Year'],
  gross: 'Gross',
  net: 'Net',
} as const

/** Teks layar sistem baru (bukan Pega) tab Loss Record. */
export const TEKS_LOSS_RECORD = {
  tambah: 'Tambah',
  hapus: 'Hapus',
  bukaBaris: 'Ubah baris',
  kosong: 'No items',
  menghitung: 'Menghitung…',
  menyimpan: 'Menyimpan…',
  tersimpan: 'Loss Record tersimpan.',
  membaca: 'Membaca berkas CSV…',
  csvKosong: 'Berkas CSV tidak berisi baris kerugian.',
  csvKolom: 'Kolom CSV tidak lengkap - wajib ada: ',
  csvBaris: 'Baris CSV ',
  csvTanggal: ': Date of Loss bukan tanggal DD/MM/YYYY.',
  csvJumlah: ': Loss Amount bukan angka.',
  csvKelasKosong: ': Class Of Business kosong.',
  csvKelasTidakAda: ': Class Of Business tidak ada di daftar - ',
  memeriksaKelas: 'Memeriksa Class Of Business…',
  csvTerbaca: ' baris dibaca dari CSV - periksa lalu klik Save Data.',
  namaBerkas: 'templatecsvlossrecord.csv',
} as const

/**
 * Isi template CSV Loss Record - persis berkas `D:\migrasi\RNM\DDL\templatecsvlossrecord.csv` dari work owner 06-10-2026
 * (padanan Rule-File-Binary `TemplateCSVLossRecord` Pega, `Activity\DownloadMasterLossRecord.xml`; isinya tidak ada di
 * korpus). CRLF, pemisah `;`. Diuji sama dengan berkas itu.
 */
export const TEMPLATE_CSV_LOSS_RECORD =
  [
    'CLASS_OF_BUSINESS;DATE_OF_LOSS;CAUSE_OF_LOSS;CURRENCY;LOSS_AMOUNT;;;;;',
    'PROPERTY ALL RISK;;Windstorm;IDR;100.000.000;;;;;',
    'PROPERTY ALL RISK;01/05/2024;Fire;USD;10.000;;;;;',
    'PROPERTY ALL RISK;04/06/2025;Accident;IDR;100.000.000;;;;;',
    'MACHINERY BREAK DOWN;03/01/2025;Windstorm;IDR;100.000.000;;;;;Note:',
    'MACHINERY BREAK DOWN;01/05/2024;Fire;USD;1.000;;;;;- Date of loss diisi dengan tanggal bulan tahun (DD/MM/YYYY)',
    'MACHINERY BREAK DOWN;04/06/2025;Accident;IDR;100.000.000;;;;;- Cause of Loss diisi dengan xxxxx',
    'PUBLIC LIABILITY;03/01/2025;Windstorm;IDR;100.000.000;;;;;',
    'PUBLIC LIABILITY;01/05/2024;Fire;IDR;100.000.000;;;;;',
    'PUBLIC LIABILITY;04/06/2025;Accident;IDR;100.000.000;;;;;',
  ].join('\r\n') + '\r\n'

/**
 * Tab Correspondence (tiket 56) - teks verbatim `NB FacIn\Section\Correspondence.xml` (grid) dan
 * `CorrespondenceContent.xml` (editor). Status: "Not Sent" dari gambar layar Pega work owner 06-10-2026 (`.IsEditable`
 * = 1; daftar nilai properti tidak ada di korpus); "Sent" `[dugaan]` padanan `.IsEditable` = 0. Diuji
 * `TabCorrespondence.test.tsx`.
 */
export const KORESPONDENSI = {
  kolom: ['Date', 'Sender', 'Subject', 'Recipient', 'Status'],
  tambah: 'Add',
  hapus: 'Remove',
  sender: 'Sender',
  senderMail: 'Sender Mail Address',
  penerima: 'Recipients:',
  subject: 'Subject',
  message: 'Message',
  simpan: 'Save',
  kirim: 'Send Email',
  belumTerkirim: 'Not Sent',
  terkirim: 'Sent',
} as const

/** Teks layar sistem baru (bukan Pega) tab Correspondence. */
export const TEKS_KORESPONDENSI = {
  kosong: 'No items',
  bukaBaris: 'Buka email',
  namaPenerima: 'Recipient',
  emailPenerima: 'Email',
  menyimpan: 'Menyimpan…',
  mengirim: 'Mengirim…',
  tersimpan: 'Email tersimpan.',
  terkirimInfo: 'Email terkirim.',
  terkunci: 'Email sudah terkirim - tidak dapat diubah.',
  penerimaKosong: 'Tambahkan minimal satu penerima dengan alamat email.',
  emailSalah: 'Alamat email penerima tidak sah: ',
  /** Penghitung ukuran Message (backend: MESSAGE VARCHAR2(4000), HTML bersih dalam bita UTF-8). */
  ukuran: ' / 4000 bita',
  terlaluPanjang: 'Message terlalu panjang - kurangi isi atau format sampai paling banyak 4000 bita.',
  editor: {
    tebal: 'Bold',
    miring: 'Italic',
    garisBawah: 'Underline',
    daftarTitik: 'Bulleted list',
    daftarNomor: 'Numbered list',
    hapusFormat: 'Remove format',
    urung: 'Undo',
    ulang: 'Redo',
  },
} as const

/**
 * Blok Email Type tahap marketing (tiket 57) - teks verbatim `NB FacIn\Section\EmailSection.xml` (label, grid riwayat) dan
 * `RejectNotificationSection.xml` (konfirmasi). Opsi = kode `.EmailType` (`Activity\SaveViewSuggest.xml`: 1 Accept, 2 Reject,
 * 4 Banding, 7 Decline); empat opsi sesuai gambar layar work owner 06-10-2026 (Ask 3 / Revise 9 tidak ditampilkan). Diuji
 * `BlokEmailType.test.tsx`.
 */
export const EMAIL_TYPE = {
  label: 'Email Type',
  opsi: [
    { kode: '1', label: 'Accept' },
    { kode: '2', label: 'Reject' },
    { kode: '4', label: 'Banding' },
    { kode: '7', label: 'Decline' },
  ],
  comment: 'Comment',
  bandingTo: 'Banding to',
  bandingKosong: 'Choose',
  pilihSebab: 'Choose Cause of decline',
  kolomRiwayat: ['Date', 'PIC', 'Approval', 'Status', 'Comment'],
  yakinReject: 'Are you sure want to reject?',
  yakinDecline: 'Are you sure want to decline?',
} as const

/**
 * Email Type tahap Team Leader (tiket 58) - properti `.EmailTypeTL` (`D:\migrasi\RNM\DDL\EmailTypeTL.xml`,
 * `ASM-FW-GISFW-WORK!EMAILTYPETL`, PromptList: 1 Accept, 2 Reject, 3 Ask, 7 Decline) = gambar layar work owner 07-10-2026.
 * Isian khusus TL `Section\Periode.xml` (`ASM-FW-GISFW-DATA-OFFERFACIN!PERIODE`) sel 81 / 82.
 */
export const EMAIL_TYPE_TL = {
  opsi: [
    { kode: '1', label: 'Accept' },
    { kode: '2', label: 'Reject' },
    { kode: '3', label: 'Ask' },
    { kode: '7', label: 'Decline' },
  ],
  deductibleAcceptance: 'Deductible Acceptance Case',
  commentDeductible: 'Comment Deductible',
} as const

/**
 * Email Type tahap Underwriting (tiket 61) - properti `.EmailTypeUW` (`D:\migrasi\RNM\DDL\EmailTypeUW.xml`,
 * `ASM-FW-GISFW-WORK!EMAILTYPEUW`, PromptList: 1 Accept, 2 Reject, 3 Ask, 7 Decline). Isian UW = `NB FacIn\Section\EmailSection_IsUW.xml`
 * (pyLabelFieldValue): Validity Of Acceptance (`.OfferFacIn.DaysValidity`), "days, or up to" (`.OfferFacIn.DateValidity`,
 * disabled), Waiting Bind (`.OfferFacIn.WaitingBindDate`, disabled). Tanpa "Choose Cause of decline" (hanya di EmailSection).
 */
export const EMAIL_TYPE_UW = {
  opsi: [
    { kode: '1', label: 'Accept' },
    { kode: '2', label: 'Reject' },
    { kode: '3', label: 'Ask' },
    { kode: '7', label: 'Decline' },
  ],
  validitas: 'Validity Of Acceptance',
  hariAtau: 'days, or up to',
  waitingBind: 'Waiting Bind',
} as const

/**
 * Email Type UW tahap Policy (tiket 69) - properti `.EmailTypeUWPolicy` (`Section\EmailSection_IsUW.xml` L2026, tampil
 * `FlagOnGoingPolicy == 1`; `D:\migrasi\RNM\DDL\EmailTypeUWPolicy.xml` `ASM-FW-GISFW-WORK!EMAILTYPEUWPOLICY`: 1 Accept, 4 Banding,
 * 7 Decline). Banding TIDAK ditampilkan (work owner 08-10-2026 K69-6: "banding hanya ada di marketing, dan banding tidak boleh
 * di tahap policy").
 */
export const EMAIL_TYPE_UW_POLICY = {
  opsi: [
    { kode: '1', label: 'Accept' },
    { kode: '7', label: 'Decline' },
  ],
} as const

/**
 * Email Type tahap MARKETING (BINDING) (tiket 64) - properti `.EmailTypeBinding` (`D:\migrasi\RNM\DDL\EmailTypeBinding.xml`,
 * `ASM-FW-GISFW-WORK!EMAILTYPEBINDING`: 1 Accept, 4 Banding, 7 Decline, 9 Revise). Section `EmailSectionCeding` (tombol "Binding",
 * centang "Received RI Slip" `.OfferFacIn.ReceivedRiSlip`); konfirmasi Revise = `Section\ReviseNotification.xml`.
 */
export const EMAIL_TYPE_BINDING = {
  opsi: [
    { kode: '1', label: 'Accept' },
    { kode: '4', label: 'Banding' },
    { kode: '7', label: 'Decline' },
    { kode: '9', label: 'Revise' },
  ],
  tombolBinding: 'Binding',
  receivedRiSlip: 'Received RI Slip',
  yakinRevise: 'Are you sure you want to revise?',
} as const

/**
 * Blok Attachment File (tiket 67) - `NB FacIn\Section\AttachmentGridReas.xml` (pyTitle "Attachment File"; kolom Category · Count ·
 * Upload File · View File) dan harness `ViewPictureList` (File · Note · Upload Date · Delete).
 */
export const LAMPIRAN = {
  judul: 'Attachment File',
  kategori: 'Category',
  jumlah: 'Count',
  unggah: 'Upload File',
  lihat: 'View File',
  berkas: 'File',
  catatan: 'Note',
  tanggal: 'Upload Date',
  hapus: 'Delete',
  /** Work owner 11-10-2026: penampil tanpa unduh, label panel lampiran bersama (`ReasViewAttachment` b2670 / Product Name Life). */
  view: 'View',
  viewOffice: 'View Office Online',
} as const

/** Teks sistem baru (bukan Pega) blok Attachment File. Konfirmasi hapus = keputusan agent (Pega tanpa konfirmasi). */
export const TEKS_LAMPIRAN = {
  kosong: 'No items',
  memuat: 'Memuat…',
  mengunggah: 'Mengunggah…',
  terunggah: (n: number, kategori: string) => `${n} berkas terunggah ke ${kategori}.`,
  /** Status unggah asinkron (ADR-0008). */
  menunggu: 'Menunggu unggah…',
  gagal: 'Unggah gagal',
  ulangi: 'Ulangi',
  yakinHapus: (nama: string) => `Hapus berkas ${nama}?`,
  ya: 'Ya',
  batal: 'Batal',
  tutup: 'Tutup',
} as const

/** Teks layar sistem baru (bukan Pega) blok Email Type / Submit. */
export const TEKS_EMAIL_TYPE = {
  memeriksa: 'Memeriksa kelengkapan data…',
  mengajukan: 'Mengajukan…',
  ya: 'Ya',
  batal: 'Batal',
  wajibType: 'Pilih Email Type.',
  wajibComment: 'Comment wajib diisi.',
  wajibBanding: 'Pilih Banding to.',
  wajibSebab: 'Pilih Cause of decline lebih dulu.',
  /** Validity Of Acceptance tahap UW (tiket 61): bila diisi, angka bulat hari. */
  hariAngka: 'Validity Of Acceptance: isi jumlah hari (angka bulat).',
  /** Tahap Binding (tiket 64): Accept butuh tombol Binding lebih dulu (Pega: Send Email Accept butuh ConfirmBinding = 1). */
  wajibBinding: 'Tekan tombol Binding lebih dulu sebelum Submit Accept.',
  membinding: 'Binding…',
  bandingTanpaOpsi: 'Case ini belum pernah diproses approver di luar marketing - Banding belum dapat dipilih.',
  terkirim: 'Case berhasil di-submit.',
  /** Tiket 69a-2: Submit menyelesaikan Policy - nomor polis yang terbit. */
  nomorPolis: (no: string) => `Policy Number: ${no}`,
  /** Jendela konfirmasi sesudah Submit dari kotak masuk Beranda (work owner 08-10-2026: NB Status tidak terlihat). */
  judulSelesai: 'Submit',
  ok: 'OK',
  riwayatKosong: 'No items',
} as const

/**
 * Popup "Choose Cause of decline" (tiket 57) - teks verbatim `NB FacIn\Section\InputCauseOfDecline.xml` (judul, kolom,
 * tombol) dan pesan `Activity\CekCauseOfDecline_ACT.xml` (APA ADANYA, termasuk spasi depan dan salah ketik "Remaks").
 * Nama T/C = nilai `POOLDATA.CAUSEOFDECLINE.TERMCONDITION` yang disebut syarat aktif / wajib di section dan activity itu.
 */
export const SEBAB_DECLINE = {
  judul: 'Cause of Decline',
  kolom: ['Term and Conditions (T/C)', 'Option', 'Remarks', '% RNM', '% SOB'],
  nonTc: 'Non - T/C',
  prioritas: ['I', 'II'],
  cek: 'Check',
  submit: 'Submit',
  pesan: {
    pctRnmSob: " PctRNM or PctSOB Can't be NULL",
    pctRnm: " PctRNM Can't be NULL",
    remaks: " Remaks Can't be NULL",
    option: "Option Can't be NULL",
    remarks: "Remarks Can't be NULL",
  },
  /** Syarat kolom (`InputCauseOfDecline` pyDisabledWhen / `CekCauseOfDecline_ACT`). */
  term: {
    option: ['Period Of Insurance (POI)', 'Extension of coverage', 'Risk Exposure'],
    remarks: ['Occupancy', 'TSI (in IDR)', 'Clauses', 'Lack of data for acceptance', 'Others T/C'],
    remarksWajib: ['Occupancy', 'TSI (in IDR)', 'Lack of data for acceptance', 'Others T/C'],
    pctRnm: ['Rate', 'R/I Comm, Deduction', 'OJK Discount', 'Loss Ratio (net)'],
    pctSob: ['Rate', 'R/I Comm, Deduction', 'OJK Discount'],
    rate: 'Rate',
    riComm: 'R/I Comm, Deduction',
    lossRatio: 'Loss Ratio (net)',
    lainnya: 'Others',
  },
} as const

export const TEKS_SEBAB_DECLINE = {
  memuat: 'Memuat daftar Cause of Decline…',
  pilihPrioritas: 'Prioritas',
  lolos: 'Cause of Decline lengkap - klik Submit.',
  kosong: 'Pilih minimal satu Cause of Decline.',
} as const

/**
 * Popup Following Policy (tiket 70) - verbatim `NB FacIn\Section\ViewOpenFollowingPolicy.xml` ("Policy List", "View Option"
 * + "Choose", "Filter" + placeholder "Policy No / ID", "Refresh", "Back"), grid `Section\GridViewFollowingNB.xml` (kolom
 * dan link), View Option `D:\migrasi\RNM\DDL\FilterForViewRetro.xml` (3 No Slip ditunda K70-3). Judul popup = label sel
 * Periode "Following Policy Number". Teks konfirmasi / kosong = sistem baru (keputusan agent).
 */
export const FOLLOWING = {
  judul: 'Following Policy Number',
  daftar: 'Policy List',
  viewOption: 'View Option',
  pilihKosong: 'Choose',
  filter: 'Filter',
  placeholder: 'Policy No / ID',
  refresh: 'Refresh',
  back: 'Back',
  opsi: [
    { value: '1', label: 'ID Pega' },
    { value: '2', label: 'Policy Number' },
    { value: '4', label: 'Reff. Number' },
    { value: '5', label: 'Insured Name' },
    { value: '6', label: 'Source Of Business' },
    { value: '7', label: 'QQ Name' },
  ],
  kolom: ['ID', 'Policy Number', 'Reff. Number', 'Insured Name', 'Source of Business', 'QQ Name', 'Start Period', 'End Period'],
  lihat: 'View Details',
  salinSemua: 'Copy All',
  salinNopolis: 'Copy Nopolis',
} as const

export const TEKS_FOLLOWING = {
  isiFilter: 'Pilih View Option lalu isi Filter untuk mencari polis.',
  kosong: 'Tidak ada polis yang cocok.',
  yakinSalinSemua: (no: string) =>
    `Copy All akan MENGGANTI seluruh lokasi / objek case ini dengan milik polis ${no}, serta mengisi Old Policy Number, % Share dan RI Commission. Lanjutkan?`,
  ya: 'Ya, Copy All',
  menyalin: 'Menyalin…',
  judulDetail: (no: string) => `Policy ${no}`,
} as const

/**
 * Menu View Policy (tiket 71) - tab Facultative "List Polis" `D:\migrasi\RNM\DDL\ViewOpenPolicyMarineInbox.xml` (label "List
 * Polis" L3255; View Option / Filter / grid = `FOLLOWING` karena sel & grid `GridOpenPolicyMarine` sama dengan
 * `GridViewFollowingNB`); "View Data" (`TempPolis.CARI12`). Judul halaman / tombol portal = nama menu Pega "View Policy"
 * (work owner 09-10-2026). Teks ringkasan polis lama = sistem baru (A71-3).
 */
export const VIEW_POLICY = {
  judul: 'View Policy',
  daftar: 'List Polis',
  viewData: 'View Data',
  kolomRingkasan: [
    'Object No', 'Risk Location', 'Object Item', 'Coverage', 'Accumulation Code', 'Currency', 'TSI 100%', 'TSI RNM', 'Rate',
    'Premium RNM', 'Treaty', '% Share', 'RI Commission',
  ],
  kepala: {
    policyNo: 'Policy Number',
    noEndors: 'Endorsement Number',
    caseId: 'ID',
    insuredName: 'Insured Name',
    qqName: 'QQ Name',
    sourceOfBusiness: 'Source of Business',
    cedingCo: 'Ceding Co',
    periode: 'Period',
    policyType: 'Policy Type',
    reffNumber: 'Reff. Number',
    groupName: 'Group Name',
    statusBusiness: 'Business Status',
    noMasterPolis: 'Master Policy Number',
  },
  ceding: ['Ceding', '% Share'],
} as const

export const TEKS_VIEW_POLICY = {
  polisLama: 'Polis ini dibuat di sistem lama (Pega) - ditampilkan ringkasan dari data produksi.',
  memuatDetail: 'Memuat detail polis…',
} as const

// ---------------------------------------------------------------------------
// Tab Object keluarga Aneka (tiket 73 P1). Label = label tampil sel Pega (preview, atau pyLabelFormatValue bila preview
// kosong); sumber per blok di komentarnya.
// ---------------------------------------------------------------------------

/** Grid lokasi Aneka - `Section\InputDtlObject_FacIn.xml` L37067-37366 (kontainer L36463 / Bonding L22724). */
export const GRID_LOKASI_ANEKA = ['No', 'Top Risk', 'Location'] as const

/** Grid Occupation di form lokasi - `Section\InputDtlObjectLocation_FacIn.xml` L15910 / L16080. */
export const GRID_OKUPASI_ANEKA = ['Occupation Code', 'Occupation Name'] as const

/**
 * Grid Object List Golf - `InputDtlObjectLocation_FacIn.xml` L12144-12602. Kolom TSI Old tidak ditampilkan di NB / RNW (work owner 09-10-2026: "tsi old hanya untuk EDM, di nb,rnw tidak perlu") -
 * menyimpang dari XML (kolomnya ALWAYS; TSI Old hanya terisi endorsement).
 */
export const GRID_GOLF_ANEKA = ['Object Name', 'Currency', 'TSI'] as const

/**
 * Kolom grid objek di detail okupasi - `Section\InputOkupasiAneka_FacIn.xml`: K2 L3852-4956, K3 (HE) L9821-10609, K4 (Growing
 * Trees / Yield Shortfall) L14377-14672.
 */
export const GRID_OBJEK_ANEKA = {
  objectName: 'Object Name',
  section: 'Section',
  limitOfLiability: 'Limit of Liability',
  year: 'Year',
  quantity: 'Quantity',
  unitQuantity: 'Unit Quantity',
  chassisNumber: 'Chassis Number',
  engineNumber: 'Engine Number',
  currency: 'Currency',
  tsi: 'TSI',
} as const

/** Popup Choose Occupation Aneka - HARNESS `Harness\ChooseOccupation.xml` (L289 judul, L1955 Search, kolom ID / Occupation Name, L3260 Choose; 10 per halaman). */
export const POPUP_OKUPASI_ANEKA = {
  judul: 'Choose Occupation',
  cari: 'Search',
  kolom: ['ID', 'Occupation Name'],
  pilih: 'Choose',
  ukuran: 10,
} as const

/** Opsi `.Section` - `DDL\Section.xml` L189-208. */
export const OPSI_SECTION_ANEKA = [
  { value: '1', label: 'Section I' },
  { value: '2', label: 'Section II' },
  { value: '3', label: 'Warranty' },
  { value: '4', label: 'No Section' },
] as const

/** Label medan form objek - `Section\InputDtlObjectAneka_FacIn.xml` (nomor baris sel per kontainer di FormObjekAneka.tsx). */
export const FORM_OBJEK_ANEKA = {
  objectDescription: 'Object Description',
  objectName: 'Object Name',
  vesselName: 'Vessel Name',
  modelTipe: 'Model dan Tipe',
  section: 'Section',
  noKontrak: 'Contract No',
  startDate: 'Start Date',
  endDate: 'End Date',
  noKtp: 'KTP No',
  noNpwp: 'NPWP No',
  currency: 'Currency',
  tsi: 'TSI',
  vesselPrice: 'Vessel Price (TSI)',
  draftWordingId: 'Draft Wording ID',
  maintenancePeriod: 'Maintenance Period',
  others: 'Others (Note)',
  productionYear: 'Production Year',
  builtIn: 'Place of Manufacture',
  flag: 'Flag',
  quantity: 'Quantity',
  invoiceNo: 'Invoice No',
  chassisNumber: 'Chassis Number',
  engineNumber: 'Engine Number',
  serialNo: 'Serial No',
  year: 'Year',
  unit: 'Unit',
  objectNameHe: 'Object Name',
  merk: 'Merk',
  type: 'Type',
  tahun: 'Tahun',
  batasGeografis: 'Batas Geografis',
  yearOfPlanting: 'Year of Planting',
  noOfTree: 'No of Trees',
  areaHectar: 'Area ( Hectar )',
  limitIdemnityDay: 'Limit Idemnity Day',
  limitIdemnityMonth: 'Limit Idemnity Month',
  briguna: {
    loanAccountNumber: 'Loan Account No',
    loanName: 'Loan Name',
    job: 'Job',
    depositAccountNo: 'Deposit Account No',
    idCardNo: 'ID Card Number',
    institution: 'Institution',
    loanType: 'Loan Type',
    dateOfBirth: 'Date Of Birth',
    salary: 'Salary',
  },
  marineHull: {
    typeOfVessel: 'Type of Vessel',
    nrtrt: 'NRT(RT)',
    lengthm: 'Length (M)',
    classification: 'Classification',
    grtrt: 'GRT(RT)',
    breadthm: 'Breadth (M)',
    construction: 'Construction',
    dwtrt: 'DWT(RT)',
    depthm: 'Depth (M)',
  },
  aviationHull: {
    registrationMarks: 'Kode Registrasi',
    maxPassengers: 'Maksimal Penumpang',
    pilots: 'Pilot',
    standardUses: 'Penggunaan Standar',
    specialUses: 'Penggunaan Khusus',
    specialRentalUses: 'Penggunaan Khusus Rental',
  },
} as const

/** Sub-form `Section\InputDtlMemorandumAccept.xml` (L470 judul; L1068-3633). */
export const FORM_MEMORANDUM = {
  judul: 'Akseptasi Memorandum',
  applicationNumber: 'Permohanan Ke',
  accumulation: 'Akumulasi',
  collateral: 'Colateral',
  insuredSiup: 'S.I.U.P',
  keyPerson: 'Key Person ( PIC )',
  originalRate: 'Rate Baku ( % )',
  principalsOutGo: 'Out Go ( % )',
  typeOfWording: 'Jenis Wording',
  buildingStatus: 'Status Bangunan',
  supportingDocument: 'Dok Pendukung',
  principalsData: 'Data Prinsipal',
  conclusionsProposals: 'Kesimpulan Dan Usulan',
  principalsBackground: 'Latar Belakang',
  summaryFinancial: 'Summary Laporan Keuangan',
  principalsGuarantee: 'Penjaminan',
} as const

/** Teks layar tab Object Aneka (bukan dari Pega). */
export const TEKS_ANEKA = {
  kosongLokasi: 'Belum ada lokasi.',
  kosongOkupasi: 'Belum ada okupasi.',
  kosongObjek: 'Belum ada objek.',
  bukaBaris: 'Buka / tutup detail',
  ketikCari: 'Ketik untuk mencari…',
  tanpaKontainer: 'Form objek untuk Class of Business ini tidak tersedia di Pega (tanpa kontainer).',
  tersimpan: 'Data objek tersimpan.',
  ditolak: 'Data objek belum dapat disimpan:',
  belumDiisi: '(belum diisi)',
  /**
   * Judul kartu pengelompokan form objek (perapian work owner 09-10-2026 "RAPIKAN BIAR USER FRIENDLY") - BUKAN dari Pega;
   * Pega menaruh semua sel dalam satu layout per kontainer.
   */
  grup: {
    objek: 'Object',
    kapal: 'Vessel Details',
    kendaraan: 'Vehicle / Equipment Details',
    pesawat: 'Aircraft Details',
    periode: 'Period',
    nilai: 'Sum Insured',
    identitas: 'Identity',
    briguna: 'Loan Details',
    catatan: 'Note',
  },
} as const

// ---------------------------------------------------------------------------
// Tab Coverage keluarga Aneka (tiket 73 P2).
// ---------------------------------------------------------------------------

/** Grid lokasi - `DDL\InputDtlCoverage_FacIn.xml` kontainer 7 L25661 / L25804. */
export const GRID_COV_LOKASI_ANEKA = ['No.', 'Location'] as const
/** Grid okupasi - `DDL\ObjectOccupation_FacIn.xml` L1277 (kontainer !IsEDM). */
export const GRID_COV_OKUPASI_ANEKA = ['Occupation Name'] as const
/** Grid objek - `NB FacIn\Section\ObjectDtlAneka_FacIn.xml` L2893-3179; grid Golf `ObjectOccupation_FacIn.xml` L3641-3927. */
export const GRID_COV_OBJEK_ANEKA = ['Object Name', 'Currency', 'TSI'] as const
/** Total per mata uang - `ObjectOccupation_FacIn.xml` L6151-6577 (`.Property.TotalTSIPremiGrossList`, `!IsMultiCOB`). */
export const GRID_COV_TOTAL_ANEKA = ['Currency', 'Grand Total Gross TSI', 'Grand Total Gross Premium', 'Grand Total Gross Rate'] as const
/** Grid `.CoverageList` - `DDL\InputCoverageAneka_FacIn.xml` (Data-Aneka) L1417 / L1561. */
export const GRID_COVERAGE_ANEKA = ['Coverage', 'Gross Premium'] as const

/** Tombol salin kontainer `IsAneka` - `DDL\InputCoverageAneka_FacIn.xml` L7966 / L9356 / L9619. */
export const SALIN_ANEKA = {
  dariSeksi: 'Copy Coverage From This Section',
  deductibleSemua: 'Copy Deductible To All',
  deductibleSeksi: 'Copy Deductible (This Section)',
} as const

/** Label sel form coverage - `NB FacIn\Section\InputCoverageAneka_FacIn.xml` (nomor baris di FormCoverageAneka.tsx). */
export const FORM_COVERAGE_ANEKA = {
  coverageNote: 'Coverage Note',
  calculateMethod: 'Calculate Method',
  proRate: '% Pro Rate',
  shortPeriod: '% Short Period',
  day: 'Day',
  currency: 'Currency',
  tsi: 'TSI',
  loading: '% Loading',
  scale: '% Scale',
  rate: '% Rate',
  indemnity: 'Indemnity',
  unit: 'Indemnity Unit',
  pctIndemnity: '% Indemnity',
  premium: 'Premium',
  pctDiscount: '% Discount',
  discount: 'Discount',
  lol: 'Limit of Liability',
  pctLol: '% Limit of Liability / % Sub Limit / % First Loss Limit',
  subLimit: 'Sub Limit',
} as const

/** `DDL\CalculateMethod_FacIn.xml` L198-212. */
export const OPSI_METODE_ANEKA = [
  { value: '1', label: 'PRO RATA' },
  { value: '2', label: 'SHORT PERIOD' },
  { value: '3', label: 'FIX RATE' },
] as const
/** `DDL\Day.xml` (radio). */
export const OPSI_DAY_ANEKA = ['365', '366', '360'] as const
/** `DDL\Unit.xml` L218-229. */
export const OPSI_UNIT_INDEMNITY = [
  { value: '0', label: 'Day (s)' },
  { value: '1', label: 'Month (s)' },
  { value: '2', label: 'Year (s)' },
] as const

/** Teks layar tab Coverage Aneka (bukan dari Pega). */
export const TEKS_COV_ANEKA = {
  menghitung: 'Menghitung premi…',
  kosongCoverage: 'Belum ada coverage.',
  tersimpan: 'Data coverage tersimpan.',
  deductibleTersalin: (n: number) => `Deductible disalin ke ${n} coverage berkode sama.`,
  deductibleKosong: 'Objek ini belum punya coverage berkode untuk disalin.',
  seksiTersalin: (n: number) => `Coverage disalin ke ${n} objek dengan Section yang sama.`,
} as const

/**
 * Tab Spreading keluarga Aneka (tiket 73 P3) - label rantai Pega: lokasi `NB FacIn\Section\CoverageSpreadingList.xml` L3591 / L3733;
 * okupasi `PropertyItemListCoverageSpreading.xml` L18903; objek `SpreadingObjectAneka.xml` L1229-1653; coverage
 * `SpreadingCoverageAneka.xml` L1218 / L1354; Golf objek `PropertyItemListCoverageSpreading.xml` L17071, coverage
 * `InputCoverageSpreadingAneka.xml` L992-1568; detail `SpreadingItemAneka.xml` L794-1285. Header % Share RNM, Copy Spreading,
 * baris spreading, total dan ringkasan = `SPREADING` (FIRE).
 */
export const SPREADING_ANEKA = {
  kolomLokasi: ['No', 'Location'],
  kolomOkupasi: ['Occupation Name'],
  /** Kolom TSI Old tidak ditampilkan di NB / RNW (work owner 09-10-2026: "tsi old hanya untuk EDM, di nb,rnw tidak perlu") - menyimpang dari XML (kolomnya ALWAYS; TSI Old hanya terisi endorsement). */
  kolomObjek: ['Object Name', 'Currency', 'TSI'],
  kolomCoverage: ['Coverage', 'Gross Premium'],
  kolomObjekGolf: ['Object Name'],
  kolomCoverageGolf: ['Coverage Note', '% Loading', '% Rate', 'TSI', 'Premi'],
  detailCoverage: {
    tsi: 'Gross TSI',
    tsiNusantaraRe: 'TSI Nusantara Re',
    premium: 'Gross Premium',
    premiNusantaraRe: 'Premium Nusantara Re',
  },
} as const

/**
 * Daftar komisi tab Premium Deduction untuk keluarga Aneka (tiket 73 P3) - rantai Pega `NB FacIn\Section\CoverageCommisionList.xml`
 * kontainer L6319 (No L6900 · Location L7040) → `PropertyItemListCoverageCommision.xml` L10644 (Occupation Name L11273) / Golf L8795
 * (Object Name L9424) → `CommisionObjectAneka.xml` (Object Name · Currency · TSI Old · TSI, L1239-1663) → `CommisionCoverageAneka.xml`
 * (Coverage · Gross Premium, L1228 / L1364) → `CommisionItemAneka.xml` (TSI Liability `.TSI` L795 · TSI Nusantara Re L980 · Gross
 * Premium L1127 · Premium Nusantara Re L1285) + grid `InputDtlSpreadingCommisionCoverage_FacIn` (= `DEDUCTION.kolomSpreading`).
 */
export const DEDUCTION_ANEKA = {
  kolomLokasi: ['No', 'Location'],
  kolomOkupasi: ['Occupation Name'],
  /** Kolom TSI Old tidak ditampilkan di NB / RNW (work owner 09-10-2026: "tsi old hanya untuk EDM, di nb,rnw tidak perlu") - menyimpang dari XML (kolomnya ALWAYS; TSI Old hanya terisi endorsement). */
  kolomObjek: ['Object Name', 'Currency', 'TSI'],
  kolomObjekGolf: ['Object Name'],
  kolomCoverage: ['Coverage', 'Gross Premium'],
  detailCoverage: {
    tsi: 'TSI Liability',
    tsiNusantaraRe: 'TSI Nusantara Re',
    premium: 'Gross Premium',
    premiNusantaraRe: 'Premium Nusantara Re',
  },
} as const

/**
 * Copy Spreading per lokasi (tiket 73 P5) - `NB FacIn\Section\PropertyItemListCoverageSpreading_IsUW.xml` kepala kolom L1571 /
 * L1735 / L1876, tombol Add L2025 · Delete L3290 · Copy To All Spreading L4444. Judul kartu = teks layar (bukan dari Pega).
 */
export const SPREADING_LOKASI = {
  judul: 'Copy Spreading Lokasi Ini',
  kolom: ['Type Treaty', '% Share', 'Type Facultative'],
  tambah: 'Add',
  hapus: 'Delete',
  salin: 'Copy To All Spreading',
  fakultatifWajib: 'Type Facultative wajib untuk treaty 10015.',
} as const

/** Opsi `.FlagSpreading` - `DDL\FlagSpreading.xml` (ASM-FW-GISFW-DATA-SPREADINGRISK!FLAGSPREADING). */
export const OPSI_FLAG_SPREADING = ['SF-MR 2023', 'Open Fac Market'] as const

// Tiket 74 Marine Cargo - label VERBATIM korpus (`pyLabelFieldValue` / `pyLabelPreview`; `pyLabelFor` hanya sasaran label).

/** Kolom grid CargoList tab Object - `NB FacIn\Section\InputDtlObject_FacIn.xml` L3014-3693 (non-UW). */
export const GRID_CARGO = ['No.', 'Top Risk', 'SailDate', 'Goods', 'Ship Name', 'Conveyance'] as const

/** Kolom grid CargoList tahap UW - `InputDtlObject_FacIn_IsUW.xml` L2989-3605. */
export const GRID_CARGO_UW = ['No.', 'Top Risk', 'Trading', 'Goods', 'Packing', 'Conveyance'] as const

/** Form cargo - `InputDtlCargo_FacIn.xml` (kepala) + `MarineCargoDtl.xml` (blok). */
export const FORM_CARGO = {
  topRisk: 'Top Risk',
  trading: 'Trading',
  goods: 'Goods',
  conveyance: 'Conveyance',
  packing: 'Packing',
  /** Teks tanpa pilihan dropdown (`pyNoSelectionText` L1692). */
  pilih: 'Choose',
  grupAgent: 'Survey Agent',
  selectAgent: 'Select Agent',
  fullName: 'Full Name',
  address: 'Address',
  city: 'City/Town',
  country: 'Country Code',
  grupLc: 'LC',
  lcStatus: 'LC Status',
  lcBank: 'LC Bank',
  lcNumber: 'LC Number',
  lcCondition: 'LC Condition',
  lcRemark: 'LC Remark',
  lcEndDate: 'LC Expired Date',
  grupKapal: 'Ship Data',
  selectShip: 'Select Ship',
  shipName: 'Ship Name',
  shipType: 'Ship Type',
  grt: 'GRT',
  nrt: 'NRT',
  dwt: 'DWT',
  age: 'Age',
  flagName: 'Flag Name',
  construction: 'Construction',
  yMake1: 'Year Make 1',
  remark: 'Remark',
  sailDate: 'Sail Date',
  fromRute: 'From Rute',
  toRute: 'To Rute',
  blNumber: 'No BL',
  invoiceNumber: 'Invoice Number',
  invoiceDate: 'Invoice Date',
  maxAge: 'Max Age ID',
  maxAgeDetail: 'Max Age Detail',
  grupKapalTambahan: 'Additional Ship',
  kolomKapalTambahan: ['Register ID', 'Ship Name', 'Voyage No'],
  tambah: 'Add',
  hapus: 'Delete',
} as const

/** Popup Select Ship - `Section\SelectShip.xml` (kolom grid `BrowseShip_RD`, tombol Select). */
export const POPUP_KAPAL_CARGO = {
  judul: 'Select Ship',
  kolom: ['Register', 'Imo', 'ID', 'Ship Name', 'Former Ship Name', 'Ship Type', 'Flag Name', 'Construction', 'Classification',
    'GRT', 'NRT', 'DWT', 'Year Make 1', 'Year Make 2', 'Rebuilt', 'Remark'],
  pilih: 'Select',
} as const

/** Popup Select Agent - `Section\SelectAgent.xml` (`BrowseLLOYDAGENT_RD`). */
export const POPUP_AGENT_CARGO = {
  judul: 'Select Agent',
  kolom: ['ID', 'Name', 'Address', 'City', 'Country'],
  pilih: 'Select',
} as const

/** Teks sistem baru - BUKAN dari Pega. */
export const TEKS_CARGO = {
  kosong: 'Belum ada objek cargo.',
  /** Pega memakai filter kolom grid; sistem baru satu kotak cari. */
  cari: 'Cari nama…',
  kosongCari: 'Tidak ada data.',
} as const

// Tiket 74 P2 - tab Coverage Marine Cargo. Label = `pyLabelPreview` (bila ada) / `pyLabelFieldValue`, sama dengan form Aneka;
// `pyLabelFor` ("Keterangan Jaminan", "Premi", "Diskon (%)") hanya sasaran label (diurai per sel `InputCoverageCargo_FacIn.xml`).

/** Grid cargo tab Coverage - `WORK!INPUTDTLCOVERAGE_FACIN` S131 L970-1718. */
export const GRID_CARGO_COVERAGE = ['No.', 'Trading', 'Goods', 'Packing', 'Conveyance'] as const

/** Grid coverage per cargo + grid total - `DATA-CARGO!INPUTDTLCOVERAGE_FACIN` L1378 / L1519, L3674-4110. */
export const GRID_COVERAGE_CARGO = {
  kolom: ['Coverage Note', 'Coverage Initial'],
  total: ['Currency', 'Grand Total Gross TSI', 'Grand Total Gross Premium', 'Grand Total Gross Rate'],
} as const

/** Form `InputCoverageCargo_FacIn.xml`. */
export const FORM_COVERAGE_CARGO = {
  coverageNote: 'Coverage Note',
  coverageInitial: 'Coverage Initial',
  selectCoverage: 'Select Coverage',
  currency: 'Currency',
  rate: 'Rate (%)',
  limitOfLiability: 'Limit of Liability',
  discountPercentage: 'Discount (%)',
  tsi: 'TSI',
  premium: 'Premium',
  minPremium: 'Min Premium',
  discount: 'Discount',
  grupWarranty: 'Warranty Extension',
  kolomWarranty: ['Coverage', 'Deductible'],
  tambah: 'Add',
  hapus: 'Delete',
  /** `.Descriptions` deductible - label "Deskripsi" `Section\InputDeductible_FacIn.xml` kontainer S3 (sel grid cargo berlabel generik "Text Area"). */
  deskripsi: 'Deskripsi',
} as const

/** Popup Select Coverage - `Section\SelectCoverage.xml` S1 (Coverage ID / Coverage Note, tombol Select). */
export const POPUP_COVERAGE_CARGO = {
  judul: 'Select Coverage',
  kolom: ['Coverage ID', 'Coverage Note'],
  pilih: 'Select',
} as const

/** Tab Spreading Marine Cargo (tiket 74 P3) - grid cargo `CoverageSpreadingList.xml` L15422-16018, grid coverage
 *  `InputCoverageSpreadingMarineCargo.xml` (pyValue kepala kolom; tanpa "Old" - K73-39 TSI / Premium Old hanya EDM). Detail
 *  coverage + baris spreading = `SPREADING` FIRE (`SpreadingItem` / `InputDtlSpreadingCoverage_FacIn` sama persis). */
export const SPREADING_CARGO = {
  kolomCargo: ['No.', 'Trading', 'Goods', 'Packing', 'Conveyance'],
  kolomCoverage: ['Coverage', 'Standard Rate (%)', 'Gross Premium', 'Premium Nusantara Re'],
} as const

/** Premium Deduction Marine Cargo (tiket 74 P4) - `CoverageCommisionList.xml` cabang Marine (cargo) → `InputCommissionMC_FacIn.xml`
 *  (Coverage / TSI / Gross Premium; TSI Old / Gross Premium Old hanya EDM - K73-39) → `CommisionItemAneka` (= `DEDUCTION_ANEKA.detailCoverage`). */
export const DEDUCTION_CARGO = {
  kolomCargo: ['No.', 'Trading', 'Goods', 'Packing', 'Conveyance'],
  kolomCoverage: ['Coverage', 'TSI', 'Gross Premium'],
} as const

/** Tab Object per halaman untuk case besar (tiket 75 F3, K75-3) - teks sistem baru, bukan Pega. */
export const TEKS_OBJEK_HALAMAN = {
  cari: 'Cari Object No / Object Name / alamat risiko / kode pos…',
  jumlah: (n: number) => `${n.toLocaleString('id-ID')} lokasi - ditampilkan per halaman dari server.`,
  belumDisimpan: (n: number) => `${n.toLocaleString('id-ID')} perubahan belum disimpan.`,
  baru: 'Lokasi baru (belum disimpan)',
  dihapus: 'akan dihapus',
  batalHapus: 'Batal hapus',
  terlaluBanyak: 'Paling banyak 500 lokasi diubah / dihapus dan 500 lokasi ditambah per Save - simpan dulu, lalu lanjutkan.',
  tersimpan: (n: number) => `Tersimpan. Case kini berisi ${n.toLocaleString('id-ID')} lokasi.`,
} as const
