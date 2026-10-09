// Label panel lampiran "Reas" (`inti/backend/dokumenpolis`) - bahasa Inggris. Caption = XML folder korpus NB FacIn,
// tempat grid `AttachmentGridReas` dan popup `ReasViewAttachment` berada (dipinjam NB / EDM Treaty In lewat
// `SetCategoryAttach`). Nomor baris = berkas yang disebut.

export const LAMPIRAN_REAS = {
  /** `Section/AttachmentGridReas.xml` b702 `<pyTitle>`. */
  judul: 'Attachment File',
  category: 'Category', // AttachmentGridReas b1340
  count: 'Count', // AttachmentGridReas b1493
  uploadFile: 'Upload File', // AttachmentGridReas b1643
  viewFile: 'View File', // AttachmentGridReas b1791
  file: 'File', // ReasViewAttachment b1363
  note: 'Note', // ReasViewAttachment b1618
  uploadDate: 'Upload Date', // ReasViewAttachment b1762
  viewOffice: 'View Office Online', // ReasViewAttachment b2670
  delete: 'Delete', // ReasViewAttachment b3434
  attach: 'Attach', // FlowAction/AttachContentGIS b24 `<pySubmitLabel>`
  cancel: 'Cancel', // FlowAction/AttachContentGIS b22 `<pyCancelLabel>`
  // `[tidak ada di korpus]` - kalimat layar ini (pola lampiran Bordereaux / Product Name Life).
  view: 'View',
  close: 'Close',
  seretBerkas: 'Drag and drop files here, or click to choose files',
  pilihBerkas: 'Choose files',
  buangPilihan: 'Remove',
  tanpaBerkas: 'No file attached',
  memuat: 'Loading attachments…',
  kosong: 'No items',
  berkasTerpilih: (n: number) => `${n} file${n === 1 ? '' : 's'} selected`,
  mengunggah: (ke: number, total: number, nama: string) => `Uploading ${ke}/${total}: ${nama}`,
  gagalUnggah: 'These files were not uploaded:',
  judulUnggah: (kategori: string) => `Upload File · ${kategori}`,
  judulLihat: (kategori: string) => `View File · ${kategori}`,
} as const
