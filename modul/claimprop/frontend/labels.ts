// Label modul Claim Prop. Label MEDAN, TOMBOL, JUDUL GRID, dan KOLOM datang dari server (VERBATIM section XML, pohon
// tata). Berkas ini hanya memuat teks perangkat halaman awal dan pop-up yang tidak punya label di XML (ditandai).

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 900). */
export const MENU_CP = { kelompok: 'Claim Prop' } as const

export const CP = {
  judul: 'Claim Prop',
  /** Tab dan switch halaman awal - keputusan work owner 07-10 dan 08-10-2026 (`[tidak ada di korpus]`). */
  tabProses: 'Process',
  tabResolve: 'Resolve',
  /** Switch Teknik = Assignment1 "Input Acceptation" (workbasket ReasKlaimTeknik); mati = worklist sendiri. */
  wbTeknik: 'Teknik',
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
  hanyaLihat: 'Read only - you do not hold this assignment.',
  pilih: 'Choose',
  tutup: 'Close',
  cariPopup: 'Search',
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
  popAdjustment: 'AdjustmentDetail',
  // CatastrofeList_Sec
  addNew: 'Add New',
  note: 'Note',
  userInput: 'User Input',
  save: 'Save',
  cancel: 'Cancel',
  totalOutstanding: 'Total Outstanding',
} as const
