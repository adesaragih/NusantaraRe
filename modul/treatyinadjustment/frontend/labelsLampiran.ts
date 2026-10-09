// Label panel Attachment layar Adjustment — tombol kelola lampiran
// (8 Oktober 2026).
//
// ⛔ TEKSNYA SAMA PERSIS dengan `LAMPIRAN` modul Treaty In
// (`modul/treatyin/frontend/labels.ts`): harness Adjustment memakai Section
// `WorkAttachments` dan `ShowAttachmentTreaty` yang sama. DISALIN, tidak
// diimpor — modul hanya boleh mengimpor `inti/**` dan dirinya sendiri
// (`inti/frontend/lapisan.guard.test.ts`).

/** Kolom modal `View File` — DUA, `File Name` dan `Type`. */
export const KOLOM_LIHAT_BERKAS = ['File Name', 'Type'] as const

export const LAMPIRAN_KELOLA = {
  judul: 'Attachment',
  /** `pyCaption Recommended safe substitute should be . or _` — disalin apa adanya. */
  spanduk: 'Recommended safe substitute should be . or _',
  unduhSemua: 'Download All', // pyButtonLabel DOWNLOAD ALL
  segarkan: 'Refresh', // pyButtonLabel REFRESH
  unggah: 'Upload', // kolom `Upload file`
  lihatBerkas: 'View', // kolom `View File`
  tutup: 'Close',
  /** Judul modal `View File` — nama panelnya, bukan judul teknis jendela Pega. */
  judulLihatBerkas: 'View File',
  tanpaBaris: 'No items',
  tanpaLampiran: 'No items',
  petunjukLampiran: '',
  /** FlowAction `TreatyAttachContent` — `pyCaption ASM Attach Content`. */
  judulUnggah: 'ASM Attach Content',
  /** Tombol FlowAction — `pyButtonLabel Attach` / `Cancel`. */
  lampirkan: 'Attach',
  batal: 'Cancel',
  /** `pyAttachmentScreen` — pemilih berkas. */
  pilihBerkas: 'Select file(s)',
  /** Kotak unggah — bentuk Master Product Name Life, sama dengan Treaty In. */
  seretBerkas: 'Drag and drop files here, or click to choose files',
  mengunggah: 'Uploading',
  buangPilihan: 'Remove',
  gagal: 'Failed',
  /** Draf penyesuaian belum ber-ID — lampiran menempel pada ID tersimpan. */
  simpanDulu: 'Save the contract first to upload attachments.',
  /** `ShowAttachmentTreaty` — `pyLabel` apa adanya. */
  viewOffice: 'View Office Online',
  /** Popup penampil kantor sebelum URL-nya datang. */
  memuatPenampil: 'Loading…',
  hapus: 'Delete',
  gantiKategori: 'Change Category',
  simpanKategori: 'Save',
} as const
