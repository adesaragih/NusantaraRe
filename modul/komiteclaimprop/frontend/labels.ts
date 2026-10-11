// Label modul Komite Claim Prop. Label LAYAR KASUS (judul, medan, judul grid, kolom, isian, tombol) datang dari server
// VERBATIM Section ShowTransfer (`models.Layar`). Berkas ini hanya memuat teks halaman awal dan perangkat yang tidak
// punya label di XML (ditandai).

/** Nama modul - nama folder korpus VERBATIM (`LayarModul.kelompok`, judul jendela berkas). Modul ini TANPA menu. */
export const KORPUS_KCP = { kelompok: 'Komite Claim Prop' } as const

/** Teks perangkat layar kasus komite yang tidak punya label di XML (ditandai). */
export const KCP = {
  memuat: 'Loading…',
  /** `[tidak ada di korpus]` - layar dibuka bukan oleh pemegang (mis. kasus sudah berpindah tingkat). */
  hanyaLihat: 'Read only - this case is waiting for another committee member.',
  memproses: 'Submitting…',
  /** `.AcceptStatus` `pyNoSelectionText` (ShowTransfer). */
  pilih: 'Choose',
  /** `[tidak ada di korpus]` - View more details: menu Claim Prop tidak dipegang akun ini (teks sama dengan Claim Prop). */
  berkasTakTerpasang: 'This file cannot be opened here: the menu of its module is not assigned to your account.',
  /** `[tidak ada di korpus]` - kembali ke inbox Claim Prop sesudah pesan Submit (kasus dibuka dari tabel komite). */
  kembali: 'Back',
} as const

/**
 * Tata letak layar komite (permintaan work owner 09-10-2026 "layout komite diperbaiki, lebih enak dilihat dan
 * userfriendly") - judul kelompok kartu, ringkasan kepala, kartu keputusan, dan langkah tangga `[tidak ada di korpus]`.
 * Label medan, judul grid, dan label isian tetap VERBATIM dari server.
 */
export const TATA_KCP = {
  kelompokKlaim: 'Claim Analysis',
  kelompokKerugian: 'Loss Details',
  kelompokEstimasi: 'Estimation',
  kelompokAdjustment: 'Adjustment',
  kelompokSpreading: 'Spreading',
  kelompokBayar: 'Payment & Bank',
  kelompokCatatan: 'Committee Notes',
  ringkasNoKlaim: 'Claim No',
  ringkasPolis: 'Policy No',
  ringkasTertanggung: 'Insured',
  ringkasAdjustment: 'Adjustment RNM',
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

/**
 * Panel lampiran klaim induk (perintah work owner 11-10-2026 "TAMBAHKAN ATTACHMENT ... KLAIM PROP KE KOMITE KLAIM
 * PROP"): teks sama dengan tab lampiran Claim Prop (screenshot Pega 09-10-2026), judul `[tidak ada di korpus]`. Hanya
 * kategori yang berisi; tombol attach nonaktif ("TOMBOL ATTACHNYA DI KOMITE DI DISABLE AJA").
 */
export const LAMPIRAN_KCP = {
  judul: 'Attachment',
  tambah: 'Add attachment',
  muatUlang: 'Refresh',
  menampilkan: 'Showing',
  dari: 'of',
  semua: 'Show All',
  kategori: 'Category',
  cacah: 'Count Attach',
  unggah: 'Upload File',
  lihat: 'View File',
  kosong: 'No attachment.',
  nama: 'File Name',
  tanggal: 'Create Date',
  oleh: 'Attached By',
  view: 'View',
  viewOffice: 'View Office Online',
  aksi: 'Action',
  unduhPilih: 'Download Selected',
  tutup: 'Close',
  /** `[tidak ada di korpus]` - keterangan tombol attach yang nonaktif. */
  hanyaDiKlaim: 'Attachments are added from the claim screen.',
} as const
