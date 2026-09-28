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
 * Peran — `[terverifikasi]` `Register_Flow.xml` `<pyPosition>` tiap
 * Assignment, dan ADR-U-0002.
 *
 * ⛔ Nilainya adalah PENGENAL yang dikirim ke backend sebagai `X-Peran`;
 * mengubahnya mengubah wewenang, bukan tampilan.
 */
export const PERAN = {
  admin: 'ReasLifeAdmin',
  medis: 'ReasLifeMedicalAdvisor',
  spv: 'ReasLifeSPV',
} as const

export type KodePeran = (typeof PERAN)[keyof typeof PERAN]

/** Sebutan peran di layar — pendamping, bukan pengganti. */
export const PERAN_ID: Record<KodePeran, string> = {
  [PERAN.admin]: 'Admin Klaim Jiwa',
  [PERAN.medis]: 'Penasihat Medis',
  [PERAN.spv]: 'Supervisor Klaim',
}

/**
 * Nama produk di topbar dan judul dokumen.
 *
 * ⛔ `[tidak ada di korpus]` — aset dan judul **e-Treaty** milik produk lain
 * dan TIDAK dipakai (brief §2 aturan 1). Teks netral dipakai sampai aset
 * resmi diberikan.
 */
export const PRODUK = {
  nama: 'Nusantara Re',
  sub: 'Reasuransi',
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
  /** b9890 `.PolicyDataLife.MarketingName` `pxAutoComplete`. */
  marketing: 'Marketing Officer',
  /** b11541 `.PolicyDataLife.CedingCoName` `pxTextInput`. */
  ceding: 'Ceding',
  /** b11736 `.PolicyDataLife.PolicyHolderName` `pxTextInput`. */
  pemegangPolis: 'Policy Holder',
  /** b12171 `.PolicyDataLife.BusinessName` `pxAutoComplete`. */
  kelasBisnis: 'Class of Business',
  /** b13384 `.PolicyDataLife.DateReceived` `pxDateTime`. */
  tanggalEmail: 'Date Received Email',
  /** b13590 `.PolicyDataLife.TanggalRespon` `pxDateTime`. */
  tanggalRespon: 'Response Date',
  /** b13795 `.PolicyDataLife.TanggalKonfirmasi` `pxDateTime`. */
  tanggalKonfirmasi: 'Confirmation Date',
  /** b14002 `.PolicyDataLife.Status` `pxTextInput`. */
  status: 'Status',
  /** b14406 `.PolicyDataLife.StatusUpdate` `pxTextInput`. */
  statusDiperbarui: 'Updated Status',
  /** b14601 `.PolicyDataLife.TanggalRealisasi` `pxDateTime`. */
  tanggalRealisasi: 'Realization Date',
  /** b16064 `pxTextInput`. */
  namaTertanggung: 'Name of Insured',
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
  /** b22641 `pyLabel` -> `pyActivity SaveAdjustment_Act` b22665. */
  simpanAdjustment: 'Save Adjustment',

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

