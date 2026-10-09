// Label modul Komite Claim Non Prop. Label LAYAR KASUS (judul, medan, judul grid, kolom, isian, tombol) datang dari
// server VERBATIM Section ShowTransfer (`models.Layar`). Berkas ini hanya memuat teks halaman awal dan perangkat yang
// tidak punya label di XML (ditandai).

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (migrasi inti 900) = nama folder korpus. */
export const MENU_KCNP = { kelompok: 'Komite Claim Non Prop' } as const

/**
 * Daftar kerja penyetuju (worklist assignment "KomiteRouter"). Worklist itu tidak punya section inbox di korpus; kolom
 * mengikuti daftar kerja Komite Claim Prop tanpa nilai / mata uang (akseptasi Non Prop bernilai per layer) -
 * `[tidak ada di korpus]`.
 */
export const KCNP = {
  judul: 'Komite Claim Non Prop',
  memuat: 'Loading…',
  kosong: 'No committee case is waiting for your decision.',
  kolomKasus: 'Case ID',
  kolomTgl: 'Update Date/Time',
  kolomStatus: 'Work Status',
  kolomKlaim: 'Claim',
  kolomTingkat: 'Level',
  kolomJabatan: 'Committe Name',
  kolomStatusBaris: 'Status',
  /** `[tidak ada di korpus]` - layar dibuka bukan oleh pemegang (mis. kasus sudah berpindah tingkat). */
  hanyaLihat: 'Read only - this case is waiting for another committee member.',
  memproses: 'Submitting…',
  /** `.AcceptStatus` `pyNoSelectionText` (ShowTransfer). */
  pilih: 'Choose',
  /** `[tidak ada di korpus]` - View Claim: menu Claim Non Prop tidak dipegang akun ini (teks sama dengan Claim Non Prop). */
  berkasTakTerpasang: 'This file cannot be opened here: the menu of its module is not assigned to your account.',
} as const

/**
 * Tata letak layar komite (pola Komite Claim Prop, permintaan work owner 09-10-2026 "ikuti tampilan klaim prop") - judul
 * kelompok kartu, ringkasan kepala, kartu keputusan, dan langkah tangga `[tidak ada di korpus]`. Label medan, judul
 * grid, dan label isian tetap VERBATIM dari server.
 */
export const TATA_KCNP = {
  kelompokKlaim: 'Claim Analysis',
  kelompokKerugian: 'Loss Details',
  kelompokAkseptasi: 'Claim Acceptation',
  kelompokAlokasi: 'Allocation',
  kelompokSpreading: 'Spreading',
  kelompokBayar: 'Payment & Bank',
  kelompokCatatan: 'Committee Notes',
  ringkasNoKlaim: 'Claim No',
  ringkasPolis: 'Policy No',
  ringkasTertanggung: 'Insured',
  ringkasTingkat: 'Committee Level',
  statusGiliranAnda: 'Waiting for your decision',
  statusMenunggu: 'Waiting for',
  statusSelesai: 'Completed',
  keputusanAnda: 'Your Decision',
  langkahSetuju: 'Approved',
  langkahTolak: 'Rejected',
  langkahBerjalan: 'In review',
  langkahMenunggu: 'Waiting',
  tanpaData: 'No data',
} as const
