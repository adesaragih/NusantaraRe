// Label modul Claim Non Prop. Label MEDAN, TOMBOL, JUDUL GRID, dan KOLOM datang dari server (VERBATIM section XML, pohon
// tata). Berkas ini hanya memuat judul pop-up (pyLabel / pyWindowName harness dan flow action), judul kolom grid pop-up
// (VERBATIM section harness), dan teks perangkat halaman awal pola Claim Prop yang tidak punya label di XML (ditandai).

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 900). */
export const MENU_CNP = { kelompok: 'Claim Non Prop' } as const

export const CNP = {
  judul: 'Claim Non Prop',
  /** Tab dan switch halaman awal - pola Claim Prop (prompt Claim Non Prop §7; `[tidak ada di korpus]`). */
  tabProses: 'Process',
  tabResolve: 'Resolve',
  inbox: 'Inbox',
  wbTeknik: 'Technical',
  teknikTanpaHak: 'Requires workbasket ReasKlaimTeknik',
  /** Start1 -> Assignment2 (`Flow_TreatyIn`); harness New tidak diekspor (OQ-CNP-23). */
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
  /** Popup master: filter per kolom, 50 baris per halaman, batas 500 (pola Claim Prop, prompt §8 P2). */
  saring: 'Filter',
  polisTanpaBerkas: 'This policy has no NB / EDM Treaty In file yet (old Pega policy not copied).',
  view: 'View',
  berkasTakTerpasang: 'This file cannot be opened here: the menu of its module is not assigned to your account.',
  menampilkan: 'Showing',
  dari: 'of',
  batasMaster: 'First 500 rows - narrow the filter',
  // judul pop-up = pyLabel / pyWindowName harness dan flow action
  popMaster: 'Choose Master Treaty XOL',
  popPolis: 'ViewListPolicyCNP',
  popSebab: 'Detail Cause of Loss',
  popKatastrofe: 'CatastropheList',
  popSelisih: 'Hitung_Test',
  popAlokasiLama: 'View Old Allocation',
  popLampiranBayar: 'ViewAttachmentNP',
  popRiwayatMaster: 'ViewHistoryMasterID_NP',
  popKomite: 'Send To Commitee',
  popTutup: 'CloseClaimMD',
  popCWP: 'CloseClaimNP',
  popPLA: 'GeneratePLACNP',
  // ChooseMasterTNonProp S1
  catatanMaster:
    'Note : Jika saat di filter COB nya kosong, silahkan search ulang dengan menginput treaty id pada kolom dibawah ini !',
  // CatastrofeList_Sec
  addNew: 'Add New',
  note: 'Note',
  userInput: 'User Input',
  save: 'Save',
  cancel: 'Cancel',
  // CauseofLoss_Section
  judulSebab: 'List Cause of Loss',
  // Hitung_Test
  nilaiEstimasi: 'Nilai Estimasi',
  nilaiAkseptasi: 'Nilai Akseptasi',
  nilaiSelisih: 'NIlai Selisih Actual Premi',
  cadangan100: '100% Reserve Updated',
  totalClaim: 'Total Claim',
  claimSpreadedKecil: 'Claim spreaded',
  claimSpreaded: 'Claim Spreaded',
  // ViewOldAllocation
  alokasiLama: 'Old XOL Allocation',
} as const
