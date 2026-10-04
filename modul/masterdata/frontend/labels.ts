// Teks layar modul Master Data. Judul tiap master datang dari `GET /api/masterdata`; label kolom di bawah (kunci JSON
// -> teks tampil), dengan nama kolom tabel sebagai cadangan.

export const TEKS = {
  sub: 'Data rujukan yang dipakai modul lain. Data tidak dihapus - nonaktifkan bila tidak dipakai lagi.',
  cari: 'Cari ID / nama',
  semuaStatus: 'Semua status',
  aktif: 'Aktif',
  nonaktif: 'Nonaktif',
  kolomStatus: 'Status',
  kolomAksi: 'Aksi',
  tambah: 'Tambah',
  ubah: 'Ubah',
  aktifkan: 'Aktifkan',
  nonaktifkan: 'Nonaktifkan',
  simpan: 'Simpan',
  menyimpan: 'Menyimpan…',
  memuat: 'Memuat data…',
  kosong: 'Belum ada data.',
  tidakCocok: 'Tidak ada data yang cocok.',
  sebelumnya: 'Sebelumnya',
  berikutnya: 'Berikutnya',
  halaman: (h: number, dari: number, total: number) => `Halaman ${h} dari ${dari} · ${total} data`,
  judulTambah: (judul: string) => `Tambah ${judul}`,
  judulUbah: (judul: string) => `Ubah ${judul}`,
  idOtomatis: 'Dibuat otomatis saat disimpan',
  turunan: 'Diisi otomatis dari data rujukan',
  wajib: 'Wajib diisi.',
  terlaluPanjang: (lebar: number) => `Paling panjang ${lebar} karakter.`,
  tersimpan: (id: string) => `Data ${id} tersimpan.`,
  statusBerubah: (id: string, aktif: boolean) => `Data ${id} ${aktif ? 'diaktifkan' : 'dinonaktifkan'}.`,
  konfirmasiStatus: (id: string, aktif: boolean) =>
    aktif ? `Aktifkan kembali data ${id}?` : `Nonaktifkan data ${id}? Data nonaktif tidak muncul di pilihan modul lain.`,
  judulKonfirmasi: (aktif: boolean) => (aktif ? 'Aktifkan data' : 'Nonaktifkan data'),
  cariRujukan: 'Ketik untuk mencari…',
} as const

/** Label kolom per kunci JSON (sama di semua master). */
export const LABEL_KOLOM: Readonly<Record<string, string>> = {
  id: 'ID',
  oldId: 'Old ID',
  note: 'Note',
  nationInitial: 'Nation Initial',
  nationId: 'Nation',
  nationName: 'Nation Name',
  provinceId: 'Province',
  province: 'Province Name',
  branchId: 'Branch',
  email: 'Email',
  moId: 'MO ID',
  jabodetabekStatus: 'Jabodetabek Status',
  cityId: 'City',
  districtName: 'District Name',
  code: 'Code',
  description: 'Description',
  groupOf: 'Group Of',
  groupOfName: 'Group Of Name',
  accumulationType: 'Accumulation Type',
  keyword: 'Keyword',
  type: 'Type',
  accumulation: 'Accumulated Type',
  accumulationName: 'Accumulation Name',
  scopeArea: 'Scope Area',
  cZone: 'CZone',
  cZoneId: 'CZone ID',
  zipCode: 'Zip Code',
  syariahStatus: 'Syariah Status',
  objectItemType: 'Object Item Type',
  pctAdjustable1: 'Pct Adjustable 1',
  pctAdjustable2: 'Pct Adjustable 2',
  objectItemTypeIna: 'Object Item Type (INA)',
  group: 'Group',
  // Jejak ubah (MD-7, migrasi 882) - turunan, baca-saja.
  createOp: 'Dibuat oleh',
  tglCreate: 'Tanggal dibuat',
  updateOp: 'Diubah oleh',
  tglUpdate: 'Tanggal diubah',
  negara: 'Negara (awalan ID)',
}

/** Label satu kolom: kamus di atas, cadangannya nama kolom tabel. */
export const labelKolom = (kunci: string, kolom: string) => LABEL_KOLOM[kunci] ?? kolom

/**
 * Kolom rujukan -> master yang dirujuk (kontrak 02: "rujukan …"). Nilai yang disimpan = `nilai` baris rujukan (ID,
 * kecuali Group Of CZone = `code`). Branch (tabel BRANCH) bukan master modul ini -> isian teks biasa.
 */
export const RUJUKAN: Readonly<Record<string, Readonly<Record<string, { master: string; nilai: string }>>>> = {
  province: { nationId: { master: 'nation', nilai: 'id' } },
  city: { provinceId: { master: 'province', nilai: 'id' } },
  district: { cityId: { master: 'city', nilai: 'id' } },
  czone: { groupOf: { master: 'czone', nilai: 'code' } },
  accumulation: {
    accumulation: { master: 'accumulatedtype', nilai: 'id' },
    provinceId: { master: 'province', nilai: 'id' },
  },
}

/** Kolom nama yang menyertai ID saat menampilkan saran rujukan. */
export const KOLOM_NAMA: Readonly<Record<string, string>> = {
  nation: 'note',
  province: 'note',
  city: 'note',
  district: 'districtName',
  czone: 'description',
  accumulatedtype: 'accumulationType',
  accumulation: 'note',
  objectitemtype: 'objectItemType',
}
