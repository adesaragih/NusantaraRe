// PUSAT LABEL UI — F0.1, brief lanjutan 6 §2 aturan 4.
//
// # Aturan yang berkas ini tegakkan
//
// SEGALA yang soal TAMPILAN diubah di FRONTEND, dan sebisanya di BERKAS INI:
// "mau ganti sebutan" tidak pernah berarti membuka backend.
//
// ⛔ SETIAP TEKS MEMBAWA BUKTI XML-nya — path dan nomor baris. Teks tanpa
// bukti adalah teks yang dikarang, dan sistem lama akan berbeda sebutan dari
// sistem baru tanpa satu pun orang menyadarinya. Bila sebuah teks memang
// tidak ada di korpus, ia ditandai `[tidak ada di korpus]` dengan terang,
// bukan diam-diam dikarang.
//
// ⚠️ Bahasa Indonesia boleh berdiri DI SAMPING, bukan menggantikan. Judul
// tahap dan label tombol adalah kosakata yang dipakai pengguna sejak sistem
// lama; menerjemahkannya diam-diam membuat pelatihan, dokumen, dan percakapan
// sehari-hari tidak lagi cocok dengan layar.

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
 * Menu sidebar — HANYA yang berbukti korpus (brief lanjutan 7 §1.2).
 *
 * ⛔ Aturan 7 brief 6 DICABUT. Sidebar tidak memuat sebelas modul lain, tidak
 * memuat butir `BelumTersedia`, dan tidak memuat *Dokumen / Komite / Detail &
 * Tutup / Cari Polis / Medical Check / Claim Analis* sebagai menu. Semua itu
 * dibuka DARI DALAM kasus lewat flow action dan popup — begitu Pega
 * melakukannya, dan menu yang tidak ada di sistem lama adalah menu yang
 * dikarang.
 *
 * ⛔ `Claim Life` sebagai nama kelompok: `[terverifikasi]`
 * `Flow/Register_Flow.xml:270` `<pyWorkTypeName>ClaimLife</pyWorkTypeName>`.
 *
 * ⚠️ `inbox` `[tidak ada di korpus]` sebagai label. Yang ada hanyalah
 * KOSAKATA-nya: berkas struktur pengekspor bernama `Struktur_InboxClaimLife`
 * dan report `InboxPremiumList`. Tidak ada harness portal Claim Life di
 * ekspor — modul yang punya menu portal diekspor bersama harness portalnya
 * (NB Treaty In `SFAPortalOpportunities`), dan Claim Life tidak.
 * `[terbuka — pemilik ekspor Pega]` apakah portalnya ada; bila jawabannya
 * datang, menu MENGIKUTI XML-nya.
 */
export const MENU = {
  /** Nama kelompok — `Register_Flow.xml:270` `pyWorkTypeName`. */
  kelompokClaimLife: 'Claim Life',
  /** `[tidak ada di korpus]` — kosakata `InboxPremiumList` / `Struktur_InboxClaimLife`. */
  inbox: 'Inbox Claim Life',
  /** VERBATIM `Flow/Register_Flow.xml:155` `<pyLabel>Register</pyLabel>`. */
  register: 'Register',
} as const

/**
 * Nama modul lain yang TIDAK BOLEH muncul di sidebar.
 *
 * ⛔ Daftar ini dipakai penjaga, bukan tampilan. Ia ada supaya "menambah satu
 * menu saja" untuk modul yang belum berbukti menjadi kegagalan uji, bukan
 * keputusan sepi yang tidak ada yang tinjau.
 */
export const MODUL_LAIN_TERLARANG = [
  'Treaty',
  'Endorsement',
  'Fac',
  'Prop',
  'Master',
  'Premium',
] as const
