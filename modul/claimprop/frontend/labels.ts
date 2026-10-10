// Label modul Claim Prop. Label MEDAN, TOMBOL, JUDUL GRID, dan KOLOM datang dari server (VERBATIM section XML, pohon
// tata). Berkas ini hanya memuat teks perangkat halaman awal dan pop-up yang tidak punya label di XML (ditandai).

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 900). */
export const MENU_CP = { kelompok: 'Claim Prop' } as const

export const CP = {
  judul: 'Claim Prop',
  /** Tab dan switch halaman awal - keputusan work owner 07-10 dan 08-10-2026 (`[tidak ada di korpus]`). */
  tabProses: 'Process',
  tabResolve: 'Resolve',
  /** Switch "Inbox ( ) Technical" (08-10-2026): mati = worklist sendiri; nyala = Assignment1 "Input Acceptation"
   *  (workbasket ReasKlaimTeknik). */
  inbox: 'Inbox',
  wbTeknik: 'Technical',
  teknikTanpaHak: 'Requires workbasket ReasKlaimTeknik',
  /** Start1 -> Assignment2 (`Flow_TreatyIn`); harness New tidak diekspor (OQ-CP-13). */
  tambahKlaim: 'Add Claim',
  cari: 'Search case ID, claim no, policy no, insured name',
  memuat: 'Loading…',
  kosong: 'No case.',
  kembali: 'Back to list',
  kolomID: 'Case ID',
  kolomTahap: 'Assignment',
  kolomNoClaim: 'Claim No',
  kolomPolis: 'Policy No',
  kolomTreaty: 'Treaty Name',
  kolomPembuat: 'Create Operator',
  kolomTanggal: 'Create Date/Time',
  /**
   * Tabel komite di bawah inbox (menu Komite Claim Prop dibuang, keputusan work owner 09-10-2026) - kolom sama dengan
   * daftar kerja Komite Claim Prop (`[tidak ada di korpus]`).
   */
  judulKomite: 'Committee',
  kosongKomite: 'No committee case is waiting for your decision.',
  komiteKolomTgl: 'Update Date/Time',
  komiteKolomKlaim: 'Claim',
  komiteKolomTingkat: 'Level',
  komiteKolomJabatan: 'Committe Name',
  komiteKolomNilai: 'Adjustment RNM',
  komiteKolomMataUang: 'Currency',
  komiteTakTerpasang: 'This committee case cannot be opened here: the Komite Claim Prop module is not active.',
  hanyaLihat: 'Read only - you do not hold this assignment.',
  pilih: 'Choose',
  tutup: 'Close',
  cariPopup: 'Search',
  /** Popup master: filter per kolom, 50 baris per halaman, batas 500 (keputusan work owner 08-10-2026). */
  saring: 'Filter',
  /** Popup tombol "+" Consultant / Adjuster (harness MstAdjusterConsultant, WindowName AdjusterConsultant). */
  popTambahAdjuster: 'Adjuster Consultant',
  /** Tombol View polis (08-10-2026): polis Pega lama yang belum disalin lewat Copy Old tidak punya berkas. */
  polisTanpaBerkas: 'This policy has no NB / EDM Treaty In file yet (old Pega policy not copied).',
  view: 'View',
  berkasTakTerpasang: 'This file cannot be opened here: the menu of its module is not assigned to your account.',
  menampilkan: 'Showing',
  dari: 'of',
  batasMaster: 'First 500 rows - narrow the filter',
  // judul pop-up = WindowName / judul layout section
  popMaster: 'Data Master TreatyIn',
  popPolis: 'Data Polis',
  popSebab: 'List Cause of Loss',
  popKatastrofe: 'Catastrophe',
  popRingkasan: 'Summary Outstanding Claim',
  popKomite: 'Komite klaim Treaty',
  popTutup: 'PreventRejectClaimProp',
  popPLA: 'Generate File PLA',
  popDLA: 'GenerateDLATreaty',
  // CatastrofeList_Sec
  addNew: 'Add New',
  note: 'Note',
  userInput: 'User Input',
  save: 'Save',
  cancel: 'Cancel',
  /** Tab "Lampiran" - VERBATIM screenshot layar Pega (work owner 09-10-2026); kolom jendela View File = Section
   *  `ViewAttachment`. Isi modal unggah (pilih berkas, Attach) `[tidak ada di korpus]` (widget bawaan Pega). */
  lampiran: 'Lampiran',
  lampiranTambah: 'Add attachment',
  lampiranMuatUlang: 'Refresh',
  lampiranMenampilkan: 'Showing',
  lampiranSemua: 'Show All',
  lampiranKategori: 'Category',
  lampiranCacah: 'Count Attach',
  lampiranUnggah: 'Upload File',
  lampiranLihat: 'View File',
  lampiranPilih: 'Select file(s)',
  lampiranSeret: 'Drag and drop files here or',
  lampiranFile: 'File',
  lampiranBuang: 'Remove',
  lampiranKirim: 'Attach',
  lampiranMengirim: 'Attaching…',
  lampiranKosong: 'No attachment.',
  /** Jendela View File - ikut NB Treaty In (`ReasViewAttachment` korpus NB FacIn: File b1363, Upload Date b1762,
   *  View Office Online b2670, Delete b3434); kolom No Acceptation / No Prekas dibuang, diganti pengunggah dan tanggal
   *  unggah (work owner 09-10-2026). `View`, `Attached By` `[tidak ada di korpus]`. */
  lampiranNama: 'File Name',
  lampiranTanggal: 'Create Date',
  lampiranOleh: 'Attached By',
  lampiranView: 'View',
  lampiranViewOffice: 'View Office Online',
  lampiranHapus: 'Delete',
  /** Jendela View File - VERBATIM screenshot layar Pega (work owner 09-10-2026): pilih kategori, centang Action,
   *  Download Selected / Delete Selected / Change Category, Create Date. "Move to" `[tidak ada di korpus]`. */
  lampiranAksi: 'Action',
  lampiranUnduhPilih: 'Download Selected',
  lampiranHapusPilih: 'Delete Selected',
  lampiranPindah: 'Change Category',
  lampiranPindahKe: 'Move to',
  totalOutstanding: 'Total Outstanding',
} as const
