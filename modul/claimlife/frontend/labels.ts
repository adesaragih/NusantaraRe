
// Label modul Claim Life - butir bi, dipisah 28-09-2026.
//
// ⛔ VERBATIM dari korpus `D:\XML\RNM_BRD\Claim Life\`, dengan nomor
// barisnya di sebelah masing-masing. Label yang diketik dari ingatan adalah
// label yang bergeser, dan yang bergeser tidak berbunyi.
//
// ⚠️ DIPISAH DARI `labels.ts` tanpa mengubah satu huruf pun. Yang tinggal di
// sana hanya yang dipakai LEBIH DARI SATU modul - menu, daftar modul, dan
// beranda. Label satu modul yang tinggal di berkas bersama membuat modul
// berikutnya menambahkan labelnya ke sana juga, dan berkas itu tumbuh
// menjadi tempat segala sesuatu.

/**
 * Judul tahap — VERBATIM `Claim Life/Flow/Register_Flow.xml`, tag
 * `<pyTaskName>` (masing-masing muncul 7 kali, juga sebagai `<pyMOName>`
 * dan `<pxLinkedRefTo>`).
 *
 * ⚠️ Keempatnya adalah nama ASSIGNMENT di alur Pega, bukan nama status.
 * Statusnya kode terpisah (`0`/`1`/`2`), dan menyamakannya akan salah.
 */
export const TAHAP = {
  inputRegister: 'Input Register',
  outstandingClaim: 'Outstanding Claim',
  medicalCheck: 'Medical Check',
  claimAnalis: 'Claim Analis',
} as const

/** Terjemahan pendamping tahap — DI SAMPING, bukan menggantikan. */
export const TAHAP_ID: Record<keyof typeof TAHAP, string> = {
  inputRegister: 'Input Register',
  outstandingClaim: 'Klaim Outstanding',
  medicalCheck: 'Pemeriksaan Medis',
  claimAnalis: 'Analis Klaim',
}

/**
 * Label tombol — VERBATIM section, dengan path dan barisnya.
 *
 * Korpus menuliskannya dalam dua bentuk huruf (`Save Adjustment` dan
 * `SAVE ADJUSTMENT`); yang dipakai bentuk berkapital judul, sebab itulah yang
 * `<pyLabelPreview>`/`<pyLabel>` bawa — bentuk HURUF BESAR datang dari
 * `pyButtonLabel …` yaitu nama RULE-nya, bukan teks yang tampil.
 */
export const TOMBOL = {
  /** `Section/ClaimLifeDetailGCNM.xml:22590` `<pyLabelPreview>` */
  simpanAdjustment: 'Save Adjustment',
  /** `Section/AdjustmentDetail_Section.xml:16169` `<pyLabelPreview>` */
  simpanKeOutstanding: 'Save to Outstanding',
  /** `Section/AdjustmentDetail_Section.xml:15101` `<pyLabelPreview>` */
  tolakOutstanding: 'Reject Outstanding',
  /** `Section/InputAkseptasiClaimLife.xml:20467` `<pyLabel>` */
  kirimBalikKeMedis: 'Send Back to Medical',
  /** `Section/InputAkseptasiClaimLife.xml:20221` `<pyLabel>` */
  kirimBalikKeAdmin: 'Send Back to Admin',
  /** `Section/CloseClaim_Section.xml:1028` `<pyLabelPreview>` */
  tutupKlaim: 'Close Claim',
  /** `Section/InputOSClaimLife.xml:24489` `pyButtonLabel Select All` */
  pilihSemua: 'Select All',
  /** `Section/InputOSClaimLife.xml:17865` `<pyLabel>` — cabut peserta (OQ-M6). */
  cabutPeserta: 'DELETE',
} as const

/**
 * Judul layar per flow action — VERBATIM nama rule `FlowAction\`, kecuali
 * yang punya label sendiri.
 */
export const LAYAR = {
  /** `FlowAction/InputRegisterClaimLife.xml` */
  register: 'Input Register Claim Life',
  /** `FlowAction/OSClaimLife.xml` */
  outstanding: 'OS Claim Life',
  /** `FlowAction/Adjustment_Detail.xml` */
  adjustmentDetail: 'Adjustment_Detail',
  /** `FlowAction/RetroClaimLife.xml` */
  retro: 'Retro Claim Life',
  /** `FlowAction/AttachDocumentLife.xml:30` `<pyLabelOld>` */
  lampiran: 'Attach Document Life',
  /** `FlowAction/UploadCSV_ClaimLife.xml` */
  unggahCSV: 'Claim Life - Upload CSV',
  /** `FlowAction/ConfirmDeleteAttachment.xml` */
  hapusLampiran: 'Confirm Delete Attachment',
  /** `FlowAction/MedicalCheck.xml` */
  medis: 'Medical Check',
  /** `FlowAction/SendtoMedical.xml` */
  kirimKeMedis: 'Send to Medical',
  /** `FlowAction/SendtoAdmin.xml` */
  kirimKeAdmin: 'Send to Admin',
  /** `FlowAction/AkseptasiClaimLife.xml` */
  akseptasi: 'Akseptasi Claim Life',
  /** `FlowAction/RejectOSClaimLife.xml` */
  tolakOS: 'Reject OS Claim Life',
  /** `FlowAction/CloseClaim.xml` */
  tutupKlaim: 'Close Claim',
  /** `FlowAction/ShowEditClaimLife.xml` */
  ubahTanggal: 'Show Edit Claim Life',
  /** `FlowAction/ViewClaimDetailLifeGCNM.xml` */
  detailKlaim: 'View Claim Detail Life',
  /** `FlowAction/PL_DetailAction_ViewPolis.xml` */
  detailPolis: 'Detail Polis Life',
} as const

/**
 * Label layar Register — VERBATIM `Section/InputRegisterClaimLife.xml`,
 * dengan nomor barisnya.
 *
 * ⛔ Ketiga belas medan datanya terikat ke `.PolicyDataLife.*` — halaman
 * POLIS, bukan isian bebas. Artinya ia DIISI oleh pemilihan polis, dan layar
 * kita harus memperlakukannya begitu pula. Mengetiknya sebagai isian kosong
 * akan membuat orang mengisi ulang apa yang sudah ada di sistem polis.
 */
export const REGISTER = {
  /** b3776 `pxButton` — memuat data polis. */
  pilihPolis: 'Choose Policy No',
  /** b7057 `pxButton`. */
  cariTertanggung: 'Find Insured',
  /** b16277 `pyLabelFieldValue`, `pyLabelFor` CARI2 -> `SearchPolicyHolder.CARI2`. */
  sertifikat: 'Certificate No',
  /** b16553 `pxButton` -> `LoadDataPesertaSpesifik_Act` b16576. */
  cari: 'Search',
  /** b20008 `pxButton`. */
  pilihTertanggung: 'Select Insured',
  /** b9104 `.PolicyDataLife.Type` `pxDropdown`. */
  type: 'Type',
  /** b9301 `.PolicyDataLife.TypeCeding` (av-2). */
  sistemReasuransi: 'System Reinsurance',
  /** b9694 `.PolicyDataLife.ProRateType` (av-2). */
  metodePremi: 'Premium Method',
  /** b9890 `.PolicyDataLife.MarketingName` `pxAutoComplete`. */
  marketing: 'Marketing Officer',
  /** b10149 `.PolicyDataLife.WPC` (av-2). */
  wpc: 'WPC',
  /** b10335 `.PolicyDataLife.RetroName` (av-2) — tampil bila Type TP/TR (b10541). */
  retroName: 'Retro Name',
  /** b10617 `.PolicyDataLife.SecurityReinsurer` (av-2) — tampil bila Type TP/TR (b10824). */
  securityReinsurer: 'Security Reinsurer',
  /** b11541 `.PolicyDataLife.CedingCoName` `pxTextInput`. */
  ceding: 'Ceding',
  /** b11736 `.PolicyDataLife.PolicyHolderName` `pxTextInput`. */
  pemegangPolis: 'Policy Holder',
  /** b11930 `.PolicyDataLife.SobName` (av-2). */
  sob: 'SOB',
  /** b12171 `.PolicyDataLife.BusinessName` `pxAutoComplete`. */
  kelasBisnis: 'Class of Business',
  /** b12450 `.PolicyDataLife.ProductNameID` (av-2). */
  idProduk: 'Product Name ID',
  /** b12648 `.PolicyDataLife.ProductName` (av-2). */
  namaProduk: 'Product Name',
  /** b13384 `.PolicyDataLife.DateReceived` `pxDateTime`. */
  tanggalEmail: 'Date Received Email',
  /** b13590 `.PolicyDataLife.TanggalRespon` `pxDateTime`. */
  tanggalRespon: 'Response Date',
  /** b13795 `.PolicyDataLife.TanggalKonfirmasi` `pxDateTime`. */
  tanggalKonfirmasi: 'Confirmation Date',
  /** b14002 `.PolicyDataLife.Status` `pxTextInput`. */
  status: 'Status',
  /** b14198 `.PolicyDataLife.TanggalKonfirmasiBalik` (av-2) — tanpa kolom. */
  konfirmasiBalik: 'Confirmation Reserved',
  /** b14406 `.PolicyDataLife.StatusUpdate` `pxTextInput`. */
  statusDiperbarui: 'Updated Status',
  /** b14601 `.PolicyDataLife.TanggalRealisasi` `pxDateTime`. */
  tanggalRealisasi: 'Realization Date',
  /** b14806 `.PolicyDataLife.KetentuanUnderwriting` (av-2) — tanpa kolom. */
  catatanUnderwriter: 'Underwriter Note',
  /** b16064 `pxTextInput`. */
  namaTertanggung: 'Name of Insured',
} as const

/**
 * Teks layar Register untuk pendaftaran BANYAK premium list sekaligus - keputusan work owner 10-10-2026 (satu case
 * per PL number + mata uang). ⛔ BUKAN korpus: Pega mendaftarkan satu polis per register, jadi teks ini tidak punya
 * baris XML; sengaja di luar REGISTER yang cacahnya dikunci ke korpus.
 */
export const REGISTER_BANYAK = {
  /** Tampilan Register (keputusan work owner 10-10-2026 "rapihkan tampilan"): label layar seragam bahasa Inggris. */
  judulHalaman: 'Register Claim Life',
  labelPL: 'Premium List No',
  kolomPilih: 'Select',
  kolomSertifikat: 'Certificate No',
  kolomNama: 'Name of Insured',
  /** Kolom DOB menggantikan "Policy No" ceding (keputusan work owner 10-10-2026). */
  kolomTanggalLahir: 'Date of Birth',
  kolomNo: 'No',
  daftarkan: 'Register Claim',
  kosongTerpilih: 'No insured selected yet.',
  /** Modal konfirmasi Register Claim (10-10-2026). */
  judulKonfirmasi: 'Register Claim',
  kalimatKonfirmasi: '{n} claim case(s) will be created:',
  lanjutkan: 'Continue?',
  batal: 'Cancel',
  mendaftarkan: 'Registering…',
  /** Pemberitahuan inbox sesudah Register berhasil (10-10-2026): {n} jumlah case, {ids} Case ID dipisah koma. */
  terdaftar: '{n} claim case(s) registered: {ids}',
  tutup: 'Close',
  terpilih: 'Selected insured',
  /** {n} = jumlah case (PL + mata uang). */
  jumlahCase: '{n} claim case(s) will be created',
  /** Judul kelompok: "PL {pl}" + chip mata uang + "{n} insured". */
  judulKelompokPL: 'PL {pl}',
  jumlahPesertaKelompok: '{n} insured',
  hapus: 'Remove',
  kosongkan: 'Clear selection',
  hasil: 'Claim cases created',
  /** Hasil Register: Case ID (CLMLF-) - nomor klaim lahir di Save to RNM (10-10-2026). */
  kolomNomor: 'Case ID',
  kolomPL: 'PL Number',
  kolomUang: 'Currency',
  kolomJumlah: 'Insured',
} as const

/**
 * Label layar Outstanding — VERBATIM `Section/InputOSClaimLife.xml`.
 *
 * ⚠️ Kesebelas medan datanya SAMA dengan layar Register dan terikat
 * `.PolicyDataLife.*` pula: keduanya menampilkan polis yang sama pada tahap
 * yang berbeda. Karena itu keduanya memakai `PanelDataPolis` yang sama, yang
 * sejak butir av membacanya dari modul PremiumList Life (`/api/polis-life/ringkas`).
 */
export const OUTSTANDING = {
  /** b951 */
  nomorKlaim: 'Claim No',
  /** b6898 */
  type: 'Type',
  /** b7686 */
  marketing: 'Marketing Officer',
  /** b9407 */
  ceding: 'Ceding',
  /** b9602 */
  pemegangPolis: 'Policy Holder',
  /** b10036 */
  kelasBisnis: 'Class of Business',
  /** b11250 */
  tanggalEmail: 'Date Received Email',
  /** b11456 */
  tanggalRespon: 'Response Date',
  /** b11662 */
  tanggalKonfirmasi: 'Confirmation Date',
  /** b11868 */
  status: 'Status',
  /** b12271 */
  statusDiperbarui: 'Updated Status',
  /** b12468 */
  tanggalRealisasi: 'Realization Date',
} as const

/**
 * Layar Outstanding mengikuti Pega (10-10-2026) — `Section/InputOSClaimLife.xml` korpus
 * `kelvin/InboxClaimLife (Done)`, nomor baris berkas itu.
 */
export const OS_PEGA = {
  /** Judul kasus Pega "Claim Life (<Case ID>)". */
  judulKasus: 'Claim Life',
  /** Tab b1956 / b5909 / b14576 `pyTitle`. */
  tabRegister: 'Register Claim',
  tabPolis: 'Policy Details',
  tabKlaim: 'Claim Details',
  /** b15524 `pyTitle`. */
  pesertaDetail: 'Participant Details',
  /** b3054 */
  plNumber: 'PL Number',
  /** b3475 */
  katastrofe: 'Catastrophe',
  /** b3781 */
  isKpr: 'Is KPR',
  /** b4726 */
  stnc: 'STNC Treaty',
  /** b4905 */
  catatanKatastrofe: 'Catastrophe Note',
  /** Kepala kolom grid b16101 / b16256 / b16406. */
  kolomPL: 'PL NUMBER',
  kolomNama: 'NAME OF INSURED',
  kolomLahir: 'DATE OF BIRTH',
  /** b17597 `pyCheckboxCaption`. */
  accept: 'Accept',
  /** b16609 `pyLabel` - tombol kepala kolom Accept. */
  pilihSemua: 'SELECT ALL',
  /** b16839 / b17865 */
  hapus: 'DELETE',
  /** b19413 */
  simpan: 'Save',
  /** DetailPolisLife b631 `pyLabel` (11-10-2026). */
  viewDetail: 'View Detail',
  /** Kepala kolom grid DetailPolisLife b4700–b6082. */
  kolomPolis: ['No.', 'POLICY NO', 'CERTIFICATE NO', 'NAME OF INSURED', 'SEX', 'DATE OF BIRTH', 'AGE', 'PLAN', 'BEGIN DATE', 'EFFECTIVE DATE', 'EXPIRED DATE'],
  /** Rincian peserta `ClaimLifeDetailGCNM.xml` (11-10-2026), label VERBATIM dengan barisnya. */
  rincian: {
    sertifikat: 'CERTIFICATE NO', // b1395
    polis: 'POLICY NO', // b1574
    plan: 'PLAN', // b1753
    umur: 'AGE', // b1932
    sex: 'SEX', // b2111
    // Kelompok rincian peserta (11-10-2026, "tidak perlu ikuti pega 100%, hanya jadikan referensi").
    kelompokTertanggung: 'Insured',
    kelompokPertanggungan: 'Coverage',
    kelompokTanggalKlaim: 'Claim Dates',
    lokasiKematian: 'LOCATION OF DEATH', // screenshot 11-10-2026 - TIDAK ada di korpus
    tempatKematian: 'Place of Death',
    negara: 'COUNTRY',
    provinsi: 'Province',
    kota: 'CITY',
    catatanAdmin: 'ADMIN NOTES', // b2588
    dob: 'DOB', // b8691
    mulai: 'BEGIN DATE', // b8877
    expired: 'EXPIRED DATE', // b9063
    efektif: 'EFFECTIVE DATE', // b9249
    valuasiMulai: 'VALUATION BEGIN DATE', // b9435
    valuasiSelesai: 'VALUATION EXPIRED DATE', // b9623
    terima: 'RECEIVED DATE', // b10187
    stncTreaty: 'STNC TREATY', // b10373
    dol: 'DATE OF LOSS', // b11099
    terimaKlaim: 'CLAIM RECEIVED DATE', // b11383
    maksKadaluarsa: 'MAX EXPIRED CLAIM', // b11921
    maksTerima: 'MAX CLAIM RECEIVED', // b12102
    dokumenLengkap: 'DOCUMENT COMPLETE DATE', // b12663
    stncData: 'STNC DATA', // b12902
    maksData: 'MAX DATA RECEIVED', // b13381
    konfirmasi: 'CONFIRMATION DATE', // b13803
    editDate: 'Edit Date', // b14115
    analis: 'CLAIM ANALYST', // b15194
    medis: 'MEDICAL ADVISOR', // kotak diagnosa + Recommendation (b4908/b6401 pyPosition MedicalAdvisor; judul 11-10-2026)
    rekomendasi: 'RECOMMENDATION', // b7316 `.RECOMMENDATION` b7323
    tabAdjustment: 'Adjustment', // b16543
    tabDokumen: 'Document', // b23293
    kolomAdjustment: ['Status', 'OutStanding Date', 'Acceptation No', 'Acceptation Date'], // grid .AdjustmentList
  },
  /** Dialog AdjustmentDetail_Section (11-10-2026), label VERBATIM. */
  adjustment: {
    mataUang: 'Currency', // b1484
    persen: 'Percent Claim (%)', // b1649
    dibayar: 'Claim Paid', // b1972
    shareRnm: 'Share Nusantara Re', // b2679
    retensi: 'Ceding Retention', // b2843
    klaimRnm: 'Claim RNM', // b3005
    sumReasured: 'Sum Reasured', // b4133
    sumInsured: 'Sum Insured', // b4296
    shareRetro: 'Share Retro', // b4459
    retroShare: 'Retroceded Share', // b4649
    simpan: 'Save to Outstanding', // b16169
    payableTo: 'Payable To', // b5922
    namaBank: 'Name of Bank', // b6116
    cabangBank: 'Branch of Bank', // b6863
    nomorRekening: 'Account No', // b7042
    // Bagian Spreading (11-10-2026, bentuk tab Spreading PremiumList Life — PANEL_SPREADING premiumlistlife/labels.ts).
    spreading: 'Spreading',
    kolomSpreading: ['Reins Type', 'Retroceded Share', 'Reinsurer', '% Share', 'Amount', 'Rate', 'Commision %', 'Ovr Comm %', 'Gross', 'Net'],
    kolomSpreadingAngka: [false, true, false, true, true, true, true, true, true, true],
    kolomRekap: ['Reinsurer', 'Currency', 'Premium Spreaded Gross', 'Premium Spreaded Net'],
    // Subjudul (keputusan work owner 11-10-2026 "tambahkan yang di bawahnya sebagai summary spreading").
    spreadingDetail: 'Spreading Detail',
    spreadingSummary: 'Spreading Summary',
    total: 'TOTAL',
  },
  /** Konfirmasi Delete — TAMBAHAN aplikasi (Pega tanpa konfirmasi, b18021; permintaan 10-10-2026). */
  konfirmasiHapus: 'Remove this insured from the claim?',
} as const

/**
 * Tombol layar Outstanding — VERBATIM, dengan barisnya.
 *
 * ⛔ Butir **aw**: kedua tombol perpindahan memanggil
 * `POST /api/klaim-life/{id}/tahap/{tujuan}`. Di Pega keduanya TIDAK menulis
 * apa pun pada posisi Admin karena prasyarat yang tampaknya salah tempel —
 * cacat rule warisan, dilaporkan `OQ-untuk-tim.md` OQ-C. Yang ditiru
 * MAKSUDnya: label tombol, penyambung `Decision3 → Assignment2`, dan
 * ADR-U-0002 ketiganya menyebut jalur balik ini sebagai fitur.
 */
export const TOMBOL_OS = {
  /**
   * b21102 `pyLabel` → `refresh` b21112 → `SaveOutStandingLife_Act` b21126.
   * Mati bila `pyWorkPage.Save = 1` (b21095) — bendera tanpa kolom (OQ-N1).
   */
  simpanRNM: 'Save to RNM',
  /** b21404 → `pyLocalAction SendtoAdmin` 21433 → tahap `input-register`. */
  kembaliKeRegister: 'Send Back to Register',
  /** b21349 / b21839 → `SendtoAdmin_Act1` 21863 → tahap `medical-check`. */
  kirimKeMedis: 'Send to Medical Check',
  /** b22750 */
  tutupKlaim: 'Close Claim',
  /** b24489 `pyButtonLabel Select All` */
  pilihSemua: 'Select All',
} as const

/**
 * Layar detail klaim — `Section/ClaimLifeDetailGCNM.xml`.
 *
 * ⛔ RALAT 27-09-2026 ATAS KESIMPULAN KAMI SENDIRI. Blok sebelumnya di sini
 * menyatakan “kelima `total…` dihitung oleh rule yang tidak ada di ekspor”,
 * lalu menyimpulkan bahwa angkanya tidak dapat ditiru. DUA hal keliru:
 *
 * 1. Totalnya ENAM, bukan lima. `Total Ceding Retention` b20629 luput karena
 *    pencacahannya memakai rujukan `CheckTotalAdjustmentClaim`, dan ia
 *    satu-satunya total yang TIDAK punya aksi refresh - jadi ia tidak ikut
 *    tercacah. Mencacah lewat pemanggil, bukan lewat label.
 * 2. Nilainya BUKAN misteri. `CheckTotalAdjustmentClaim` memang nol berkas
 *    rule-nya (rujukan menggantung itu nyata, dan tetap OQ-H), tetapi yang
 *    hilang hanya pemanggil REFRESH di layar. Yang MENGHITUNG keenam angka
 *    itu ada di korpus dan berprasyarat kosong:
 *
 *       `Activity/SavePesertaClaim.xml`        langkah 8 b4002 / 8.1 b4221 / 8.2 b4592
 *       `Activity/SaveOutStandingLife_Act.xml` langkah 23 b10638 / 23.1 b10841 / 23.2 b11067
 *
 *    Rumusnya: jumlah SELURUH baris `.AdjustmentList` peserta itu, tanpa
 *    memandang `STS_REJECT` dan tanpa memandang `IsCheck`.
 *
 * Keenamnya kini dihitung `models.HitungTotalPeserta` dan menyeberang di
 * JSON peserta. Sebab salah bacanya: berhenti pada rule yang NAMANYA
 * tertulis di section, tanpa menanyakan siapa lagi yang menulis medan itu.
 */
export const DETAIL = {
  /** b5061 `pyLabel` `pxButton` -> `showHarness` b5071 `Diagnose_Harness`. */
  cariPenyakit: 'Find Disease',
  /** b14115 `pyLabel` -> `pyLocalAction ShowEditClaimLife` b14144. */
  ubahTanggal: 'Edit Date',
  /** b12124 `pyLabelPreview` -> `.MAXCLAIM_RECEIVED` b12131 (butir bk), read-only. */
  maxClaimReceived: 'MAX CLAIM RECEIVED',
  /** b22641 `pyLabel` -> `pyActivity SaveAdjustment_Act` b22665. */
  simpanAdjustment: 'Save Adjustment',
  /**
   * Butir bo — grid `.AdjustmentList` b17126: b17937 `pyLabel` -> `addRow`
   * b17947 + `SetIndexAdjustmentList` b17991; tampil hanya bila
   * `pyWorkPage.pyPosition =='ReasLifeSPV'` b18160.
   */
  tambahAdjustment: 'Add',

  /**
   * `CloseClaim_Section.xml` b1081 `pyLabel`.
   *
   * ⚠️ Tombolnya menjalankan DUA aksi pada satu klik: `refresh` ->
   * `ProtectCloseClaim_act` b1101, DAN `closeContainer` b1129. Laporan ronde
   * pertama menulis "tidak punya aksi lain" dan itu keliru.
   */
  tutupKlaim: 'Close Claim',
  /** `CloseClaim_Section.xml` b499 `pyValue` -> `pyCaption` b1499. */
  konfirmasiTutup: 'Are you sure want to Close Claim?',

  /**
   * b20629 `pyLabelPreview` -> `.TotalCedingRetention` b20636.
   *
   * ⛔ Total PERTAMA di layar, dan satu-satunya yang tidak punya aksi
   * refresh - sebab itu ia sempat hilang dari sini selama enam hari.
   */
  totalCedingRetention: 'Total Ceding Retention',
  /** b20914 `pyLabelPreview` -> `.TotalShareRNM` b20921. */
  totalShareNusantaraRe: 'Total Share Nusantara Re',
  /** b21201 `pyLabelPreview`. */
  totalSumInsured: 'Total Sum Insured',
  /** b21488 `pyLabelPreview`. */
  totalSumReasured: 'Total Sum Reasured',
  /** b21775 `pyLabelPreview`. */
  totalShareRetro: 'Total Share Retro',
  /** b22063 `pyLabelPreview`. */
  totalClaimAmount: 'Total Claim Amount',
} as const

/**
 * Dialog `Edit Date` — `Section/EditDateClaimLife_Section.xml`, dibuka tombol
 * `DETAIL.ubahTanggal` b14115. Label VERBATIM `pyLabelPreview`, kecuali Save.
 *
 * ⚠️ DOL tidak ada di sini: kotaknya sudah berdiri sendiri sejak tiket 06
 * dengan label `DETAIL.ubahTanggal`, dan menyimpannya lewat rute sendiri.
 */
export const EDIT_DATE = {
  /** b1069 → `.CLAIM_RECEIVED_DATE` b1076. */
  terimaKlaim: 'CLAIM RECEIVED DATE',
  /** b1381 → `.COMPLETE_DATE` b1387. */
  dokumenLengkap: 'DOCUMENT COMPLETE DATE',
  /** b1619 → `.CONFIRMATION_DATE` b1626. */
  konfirmasi: 'CONFIRMATION DATE',
  /** b1910 `pyLabel` → `UpdateDateClaimLife_Act` b1929. */
  simpan: 'Save',
} as const

/**
 * Dialog Reject Outstanding — `Section/RejectOSClaimLife_Sec.xml`, dibuka flow
 * action `FlowAction/RejectOSClaimLife.xml`. OQ-M5 ditutup 29-09-2026 (GILIRAN-17):
 * alasannya disimpan di `T_CLAIMLF_JEJAK.KOMENTAR` (migrasi 021).
 */
export const REJECT_OS = {
  /** b783 → `.DataCommitteeTreaty.TanggalComitee` b790; tidak pernah ditulis activity. */
  tanggal: 'Date',
  /** b969 → `OperatorID.pyUserName` b975, baca-saja b933. */
  pic: 'PIC',
  /** b1680 (`pyUseLabelDesc` b1674) → `.DataCommitteeTreaty.Remarks` b1687, wajib b1695. */
  alasan: 'Remarks',
  /** b3098 `pyLabel` → `RejectOSClaimLife_Act` b3117. */
  kirim: 'Submit',
  /** `FlowAction/RejectOSClaimLife.xml` b19 `pyCancelLabel`. */
  batal: 'Cancel',
} as const

/**
 * Layar dokumen pendukung — `Section/DocumentLife.xml`.
 *
 * Nomor baris dibaca 27-09-2026 dari korpus apa adanya.
 *
 * ⛔ Ketiga tombol selain `Refresh` BELUM terpasang: unggah, unduh, dan
 * hapus dokumen adalah kelompok Dokumen, yang menyentuh penyimpanan luar
 * (Google Storage, ADR-U-0010) dan outbox `T_LOG_SERVICE_RNM`. Labelnya ada
 * di sini supaya layar dapat MENYEBUT apa yang belum ada, bukan diam —
 * layar yang tampak lengkap padahal tidak adalah layar yang tidak akan
 * dicari lagi (pelajaran butir av).
 */
export const DOKUMEN = {
  /** b611 `pyLabel` -> `pyAction refresh` b619. */
  muatUlang: 'Refresh',
  /** b1245 `pyLabel` -> `localAction` b1254 `AttachDocumentLife` b1273. */
  tambahLampiran: 'Add attachment',
  /** b3502 `pyLabel` -> `runActivity` b3511 `DownloadDocumentClaim` b3519. */
  lihatOfficeOnline: 'View Office Online',
  /**
   * b4288 `pyLabel` -> `localAction` b4297 `ConfirmDeleteAttachment` b4317.
   *
   * ⚠️ Tombolnya menjalankan DUA aksi pada satu klik: local action DAN
   * `closeContainer` b4441 — pola yang sama dengan `Close Claim`.
   */
  hapus: 'Delete',
} as const

/**
 * Tombol layar **Medical Check** — `Section/MedicalCheckClaimLife.xml`.
 *
 * ⛔ Dua tombol perpindahan, tidak lebih. `Send to Medical Check` TIDAK ada
 * di sini — ia milik layar Outstanding *(`TOMBOL_OS`)*. Tiap layar hanya
 * menawarkan perpindahan yang section-nya sendiri punya; menyalin tombol
 * antarlayar berarti membuka jalur yang di sistem lama tidak ada.
 */
export const TOMBOL_MEDIS = {
  /** b18572 `pyLabel`. Bukan perpindahan — ia menyimpan. */
  simpan: 'Save',
  /** b20256 → `pyLocalAction SendtoAdmin` b20285 → tahap `outstanding`. */
  kembaliKeAdmin: 'Send Back to Admin',
  /**
   * b21151 → `SendtoAdmin_Act1` b21174 **dan** `finishAssignment` b21202.
   *
   * ⚠️ DUA aksi pada satu klik, pola yang sama dengan `Close Claim` dan
   * `Delete`. Yang membawa akibat tahap adalah `finishAssignment` →
   * Transition1 → Decision1 → `Else` → Claim Analis.
   */
  kirimKeAnalis: 'Send to Claim Analyst',
} as const

/**
 * Tombol layar **Claim Analis** — `Section/InputAkseptasiClaimLife.xml`.
 *
 * `Close Claim` b21428 ada di layar ini pula, tetapi labelnya sudah hidup di
 * `DETAIL.tutupKlaim` — satu label, satu tempat.
 */
export const TOMBOL_AKSEPTASI = {
  /** b18543 `pyLabel`. */
  simpan: 'Save',
  /** b20221 → `SendtoAdmin` b20250 → tahap `outstanding`. */
  kembaliKeAdmin: 'Send Back to Admin',
  /** b20467 → `SendtoMedical` b20496 → tahap `medical-check`. */
  kembaliKeMedis: 'Send Back to Medical',
} as const

/**
 * Konfirmasi jalur balik — DUA local action, bukan per tombol.
 *
 * `[terverifikasi]` tombol yang membuka local action (bukan activity langsung):
 *
 *   `SendtoAdmin`   InputOSClaimLife b21433 · MedicalCheckClaimLife b20285 ·
 *                   InputAkseptasiClaimLife b20250 → `SendtoAdmin_Section`
 *   `SendtoMedical` InputAkseptasiClaimLife b20496 → `SendtoMedical_Section`
 *
 * ⚠️ `Send Back to Register` (Outstanding) pun bertanya "Send Back to Admin?" —
 * ia membuka local action yang SAMA. Ditiru apa adanya.
 *
 * `Send to Medical Check` dan `Send to Claim Analyst` TIDAK bertanya: keduanya
 * `SendtoAdmin_Act1` + `finishAssignment` langsung (b21863/b21891, b21174/b21202).
 */
export const KONFIRMASI_BALIK = {
  /** `SendtoAdmin_Section.xml` b566 `pyValue`. */
  keAdmin: 'Send Back to Admin?',
  /** `SendtoMedical_Section.xml` b577 `pyValue`. */
  keMedis: 'Send Back to Medical?',
  /** b1229 (Admin) / b1267 (Medical) `pyLabel` → activity lalu `finishAssignment`. */
  kirim: 'Submit',
  /**
   * `FlowAction/SendtoAdmin.xml` b19 / `SendtoMedical.xml` b21 `pyCancelLabel`
   * — VERBATIM, menggantikan teks karangan "Batal" (GILIRAN-11 paket 4).
   */
  batal: 'Cancel',
} as const

/**
 * Tombol penyerahan ke Komite — `Section/ClaimComite.xml`.
 *
 * ⛔ RALAT 27-09-2026. Layar Detail sebelumnya memakai teks karangan
 * **`Send ke Komite`** — campuran Indonesia-Inggris yang tidak ada di korpus
 * mana pun. Label tombol diambil dari XML, bukan dikarang: pemakai sistem
 * lama mencari kalimat yang sama, dan pengujinya membandingkan layar lama
 * dengan layar baru kata demi kata.
 *
 * ⚠️ Section-nya bernama `ClaimComite` tetapi `pyRuleName`-nya
 * `ClaimComiteeLife` (b139) — dua ejaan, keduanya salah eja bahasa Inggris,
 * dan keduanya warisan. Disebut apa adanya supaya dapat dicari.
 */
export const TOMBOL_KOMITE = {
  /** b7033 `pyLabel`. */
  serahkan: 'Send Claim to Committee',
  /** b7715 `pyLabel`. */
  batal: 'Cancel',
} as const

/**
 * Grid diagnosa — `Section/ClaimLifeDetailGCNM.xml` b3923 `.DiagnoseList`.
 *
 * ⛔ Label VERBATIM, termasuk HURUF BESARNYA. Kepala kolom di section itu
 * ditulis `DIAGNOSE`, bukan `Diagnose` — dan merapikannya berarti layar ini
 * menyebut kolom dengan kata yang tidak ada di sistem lama. Orang yang
 * beralih membaca layar, bukan kode.
 *
 * ⛔ `GROUP DIAGNOSE` punya kepala kolom tetapi **tidak punya daftar
 * pilihan**: dropdown b5863 ber-`pyListSource associated`, artinya daftarnya
 * hidup pada rule properti `.GROUPDIAGNOSE` kelas `Data-DiagnoseLife` — yang
 * tidak ada di ekspor dan tidak ada di katalog DEV. `pyLabelPreview` sel itu
 * (b5854) pun KOSONG, jadi kepala kolom b4490 adalah satu-satunya kata yang
 * sah untuknya. Butir **bf**, **OQ-L**.
 */
export const DIAGNOSA = {
  /** b4188 `pyValue` — kepala kolom pertama. */
  kolomNama: 'DIAGNOSE',
  /** b4337 `pyValue` — kepala kolom kedua. */
  kolomIcd: 'ICD CODE',
  /** b4490 `pyValue` — kepala kolom ketiga. */
  kolomKelompok: 'GROUP DIAGNOSE',
  /** b4690 `pyLabel` -> `addRow` b4700 lalu `refresh` b4730. */
  tambah: 'Add',
  /**
   * b6160 `pyLabel` -> `deleteRow` b6170 **lalu `save` b6191**.
   *
   * ⚠️ DUA aksi pada satu klik, dan yang kedua itulah bedanya dengan `Add`:
   * penghapusan MENETAP seketika, penambahan tidak. Pola yang sama dengan
   * `Close Claim` dan `Delete` dokumen.
   */
  hapus: 'Delete',
  /** `Diagnose_Section.xml` b2509 `pyLabel` -> `SetDisease` b2528. */
  pilih: 'Choose',
} as const

/**
 * Tabel komite di bawah inbox (perintah work owner 09-10-2026: menu Komite Claim Life dihapus) - judul dan kolom
 * disalin dari Inbox Komite modul Komite Claim Life. ⚠️ Tiga kolom pertama dipinjam dari kolom Pega-standar
 * `InboxPremiumList` (`pyFieldLabel` b721/b735/b764): worklist `KomiteRouter` tidak punya section inbox sendiri.
 */
export const TABEL_KOMITE = {
  judul: 'Inbox Komite',
  kosong: 'Tidak ada kasus komite yang menunggu keputusan Anda.',
  kasusId: 'Case ID',
  tglUpdate: 'Update Date/Time',
  statusWork: 'Work Status',
  nomorKlaim: 'Klaim induk',
  tingkat: 'Tingkat',
  nilaiKlaim: 'Nilai klaim',
  mataUang: 'CURRENCY',
  statusBaris: 'Status baris',
  /** Laporan "perlu intervensi" — tiket 08 Komite Claim Life, tidak ada di korpus (ADR-0015). */
  laporan: 'Perlu intervensi hari ini',
  laporanKosong: 'Laporan hari ini: tidak ada efek keluar yang perlu intervensi.',
  /** `[tidak ada di korpus]` - modul Komite Claim Life tidak dipasang bagi akun ini. */
  takTerpasang: 'Kasus komite ini tidak dapat dibuka di sini: modul Komite Claim Life tidak aktif.',
} as const
