// Label modul Cover Life - bahasa Inggris, VERBATIM dari section Pega `InboxCoverLife`
// (`D:\NUSARE DEV\Menu Cover\InboxCoverLife.xml`, nomor = baris XML): judul "COVER" b368, form "Cover" b848 (wajib b843)
// dan "Note" b1026, tombol Save b5407 dan Cancel b5671, grid "ID" b3023 / "Cover" b3161 / kolom tombol tanpa judul b3277
// dengan Edit b3784. Form TANPA medan ID. Teks lain (halaman, pesan) = teks antarmuka aplikasi, bukan label Pega.

/** Nama menu - `M_NAV_MENU.LABEL` baris modul ini (slot menu modul 957, keputusan work owner 08-10-2026 C4). */
export const MENU_CVL = { kelompok: 'Cover Life' } as const

export const CVL = {
  judul: 'COVER',
  id: 'ID',
  cover: 'Cover',
  note: 'Note',
  save: 'Save',
  menyimpan: 'Saving…',
  batal: 'Cancel',
  edit: 'Edit',

  modeUbah: (id: string) => `Editing ID ${id}`,
  galatWajib: 'Cover is required',
  galatPanjang: (medan: string, batas: number) => `${medan} is longer than ${batas} characters`,
  tersimpan: (id: string) => `Cover ${id} saved.`,
  memuat: 'Loading covers…',
  kosong: 'No cover yet.',
  jumlah: (n: number) => `${n} cover${n === 1 ? '' : 's'}`,
  halaman: (h: number, dari: number) => `Page ${h} of ${dari}`,
  sebelum: 'Previous',
  sesudah: 'Next',
} as const
