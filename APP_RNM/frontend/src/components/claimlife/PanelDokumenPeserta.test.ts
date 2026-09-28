import { describe, expect, it } from 'vitest'

import { DOKUMEN } from '../../assets/labels.claimlife'
import { dokumenTerunggah, tautanDokumen, type Dokumen } from '../../services/api'
import { barisDokumen, PRANALA_MENUNGGU } from './PanelDokumenPeserta'

// Uji daftar dokumen pendukung — A3 kelompok 1.

function dok(lebih: Partial<Dokumen> = {}): Dokumen {
  return {
    id: 1,
    pesertaId: 'UJI-PES-1',
    namaFile: 'UJI-berkas.pdf',
    mime: 'application/pdf',
    kategori1: 'UJI-K1',
    kategori2: 'Surat Keterangan',
    tStorageId: 'UJI-STORAGE-1',
    tanggal: '2026-01-15',
    paymentDate: '',
    ...lebih,
  } as unknown as Dokumen
}

describe('barisDokumen', () => {
  it('kosong tetap kosong, bukan baris kosong', () => {
    expect(barisDokumen([])).toEqual([])
  })

  it('menampilkan KATEGORI_2, bukan KATEGORI_1', () => {
    // ⛔ `DocumentLife.xml` menampilkan KATEGORI_2 kepada manusia; KATEGORI_1
    // adalah penggolong yang dipakai saringan `Obj-Browse`
    // (`Field .KATEGORI_1`, `Value .DOCUMENT`). Menukarnya membuat layar
    // menampilkan kode penggolong alih-alih kata yang dibaca orang.
    const [b] = barisDokumen([dok({ kategori1: 'SALAH', kategori2: 'BENAR' })])
    expect(b?.kategori).toBe('BENAR')
  })

  it('membedakan "belum terunggah" dari "pranala belum tersambung"', () => {
    // ⚠️ Dua keadaan berbeda, dan hanya yang kedua menunggu kita.
    // `T_STORAGE_ID` kosong berarti berkasnya memang belum naik;
    // terisi berarti berkasnya ada dan yang belum ada adalah URL-nya
    // (`GetUrlGoogleStorage_Act` belum disambungkan).
    const [naik] = barisDokumen([dok({ tStorageId: 'UJI-STORAGE-1' } as Partial<Dokumen>)])
    const [belum] = barisDokumen([dok({ tStorageId: '' } as Partial<Dokumen>)])
    expect(naik?.terunggah).toBe(true)
    expect(belum?.terunggah).toBe(false)
  })

  it('tanggal KOSONG tetap kosong, bukan tanggal nol', () => {
    // ADR-U-0027 dan ADR-U-0022 Akibat 2: kolom tanggal yang NULL menjadi
    // KOSONG, bukan 0001-01-01. Tanggal nol akan tampil sebagai tanggal
    // sungguhan dan tidak ada yang tahu ia karangan.
    const [b] = barisDokumen([dok({ tanggal: '' })])
    expect(b?.tanggal).toBe('')
  })

  it('tanggal NULL diperlakukan sama dengan kosong', () => {
    // ⛔ CACAT LINTAS-LAPIS KELIMA, dan ia lahir di giliran yang sama
    // ketika keempat pendahulunya didaftarkan. Go mengirim `null` untuk
    // kolom DATE yang NULL (`*time.Time`), sedangkan `api.ts` sempat
    // menyatakan `tanggal: string`. TypeScript tidak dapat menangkapnya -
    // tipe yang berbohong tentang data dari jaringan tetap dikompilasi.
    // Akibatnya `=== ''` tidak pernah menyala, dan tanggal yang memang
    // kosong tampil sebagai sel kosong alih-alih penanda.
    const [b] = barisDokumen([dok({ tanggal: null } as Partial<Dokumen>)])
    expect(b?.tanggal).toBe('')
  })

  it('membawa nama BERKAS, dan tidak ada medan lain yang dapat memuat nama orang', () => {
    // ⛔ Penjaga arah yang sama dengan marshaller peserta: medan pada baris
    // ini mudah bertambah satu per satu, dan satu di antaranya suatu hari
    // dapat berisi nama tertanggung. Himpunan kuncinya dikunci UTUH.
    const [b] = barisDokumen([dok()])
    expect(Object.keys(b ?? {}).sort()).toEqual([
      'id',
      'kategori',
      'namaFile',
      'tanggal',
      'terunggah',
    ])
  })

  it('kalimat penanda pranala disebut apa adanya', () => {
    expect(PRANALA_MENUNGGU).toContain('menunggu')
  })
})

describe('label DocumentLife', () => {
  it('keempat tombol VERBATIM', () => {
    // Dibaca dari `Section/DocumentLife.xml`: b611, b1245, b3502, b4288.
    expect(DOKUMEN.muatUlang).toBe('Refresh')
    expect(DOKUMEN.tambahLampiran).toBe('Add attachment')
    expect(DOKUMEN.lihatOfficeOnline).toBe('View Office Online')
    expect(DOKUMEN.hapus).toBe('Delete')
  })

  it('cacahnya dikunci — tombol yang hilang dari layar tidak berbunyi sendiri', () => {
    expect(Object.keys(DOKUMEN)).toHaveLength(4)
  })
})

// Butir be — unggah/unduh/hapus lewat outbox, 27-09-2026.
describe('dokumen lewat outbox', () => {
  it('tautan unduh memakai rute kita, bukan URL layanan luar', () => {
    // ⛔ Butir be: selama penyambungan nyata belum disetujui, tautannya
    // rute KITA — dan rute itu berbatas identitas.
    // ⛔ TEKS, dan uji inilah yang menemukan sebabnya. Ditulis sebagai
    // angka, `20260927103000123` menjadi `...124` di JavaScript — 17 angka
    // berada di luar `Number.MAX_SAFE_INTEGER`. Pengenal dokumen karena itu
    // menyeberang sebagai teks di kedua sisi.
    const t = tautanDokumen('20260927103000123')
    expect(t).toContain('/api/dokumen/20260927103000123/isi')
    expect(t).not.toContain('googleapis')
  })

  it('tStorageId kosong berarti SEDANG diproses, bukan hilang', () => {
    // ⛔ Dibedakan supaya pemakai tahu menunggu, bukan mengunggah ulang.
    // Baris ganda adalah akibat langsung dari layar yang tidak membedakan.
    expect(dokumenTerunggah(dokUji(''))).toBe(false)
    expect(dokumenTerunggah(dokUji('   '))).toBe(false)
    expect(dokumenTerunggah(dokUji('20260927103000123'))).toBe(true)
  })
})

/** Dokumen seadanya; hanya `tStorageId` yang diuji di sini. */
function dokUji(tStorageId: string): Dokumen {
  return {
    id: '20260927103000123',
    pesertaId: 'P-1',
    namaFile: 'UJI-berkas.pdf',
    mime: 'application/pdf',
    kategori1: 'DL-1',
    kategori2: 'KTP',
    tStorageId,
    tanggal: null,
    paymentDate: null,
  }
}
