import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { KEPUTUSAN_POLIS } from '../labels'
import { bolehRejectDiTahap, TAHAP_POLIS } from '../api'
import { ringkasanAkibat } from './InputOffer'

// Uji layar keputusan penawaran — tiket 01 bagian 2; GILIRAN-14 butir bq.

const BERKAS = readFileSync(join(__dirname, 'InputOffer.tsx'), 'utf8')
/** Sumber tanpa komentar — prosa yang menerangkan tidak dituduh kode. */
const SUMBER = BERKAS.replace(/\{\/\*[\s\S]*?\*\/\}/g, '')
  .split('\n')
  .filter((b) => {
    const t = b.trimStart()
    return !t.startsWith('//') && !t.startsWith('*')
  })
  .join('\n')

describe('Reject hanya di tahap yang punya konektornya', () => {
  it('tahap penawaran TIDAK menawarkan Reject', () => {
    // ⛔ `Reject` muncul TEPAT SEKALI di seluruh flow — Transition9 b2306
    // pada Decision2, sesudah Input Premium Detail. Tombol yang pasti
    // dijawab 409 adalah tombol yang mengajari orang mengabaikan galat.
    expect(bolehRejectDiTahap(TAHAP_POLIS.penawaran)).toBe(false)
  })

  it('tahap detail menawarkannya', () => {
    expect(bolehRejectDiTahap(TAHAP_POLIS.detail)).toBe(true)
  })

  it('tahap summary tidak — ia nol konektor keputusan', () => {
    expect(bolehRejectDiTahap(TAHAP_POLIS.summary)).toBe(false)
  })

  it('tombol Reject tidak ditampilkan di tahap mana pun (keputusan work owner 03-10-2026)', () => {
    expect(SUMBER).not.toContain('KEPUTUSAN_POLIS.reject')
    expect(SUMBER).not.toContain('bolehRejectDiTahap(tahap) && (')
  })
})

describe('Decision3 tidak ditanyakan (butir bq)', () => {
  it('nol tombol Offer/Premium dan nol rute penggolong di layar', () => {
    // ⛔ `Offer`/`Premium` hasil decision table `IsFlagOnGoingPolicy` atas
    // bendera kasus — backend menerapkannya di dalam `Confirm`.
    expect(SUMBER).not.toContain('golongkanPenawaran')
    expect(SUMBER).not.toContain('PENGGOLONG_POLIS')
    expect(SUMBER).not.toContain('menungguPenggolong')
  })

  it('setiap keputusan yang berhasil melepas kasus dari layar ini', () => {
    // Sesudah bq, `Confirm` selalu menutup (Offer) atau memindahkan
    // (Premium) — tidak ada keadaan "menunggu" yang menahan layar. Sejak
    // 03-10-2026 keputusan yang menerbitkan PL Number melepasnya SESUDAH
    // popup nomornya ditutup; selebihnya langsung.
    expect(SUMBER).toMatch(/setAkibat\(hasil\)[\s\S]{0,400}\n\s*onSelesai\(\)\n\s*\} catch/)
    expect(SUMBER).toMatch(/setNomorTerbit\(null\)\n\s*onSelesai\(\)/)
  })
})

describe('ringkasanAkibat', () => {
  it('Confirm + bendera "0" (Offer) — tertutup Resolved-Completed', () => {
    const s = ringkasanAkibat({ tahapTujuan: '', statusWork: 'Resolved-Completed' })
    expect(s).toContain('Resolved-Completed')
  })

  it('Confirm + bendera "1" (Premium) — berpindah ke Input Premium Detail', () => {
    const s = ringkasanAkibat({ tahapTujuan: TAHAP_POLIS.detail, statusWork: '' })
    expect(s).toContain(TAHAP_POLIS.detail)
  })

  it('berpindah menyebut tahap tujuannya', () => {
    const s = ringkasanAkibat({ tahapTujuan: TAHAP_POLIS.penawaran, statusWork: '' })
    expect(s).toContain(TAHAP_POLIS.penawaran)
  })
})

describe('label VERBATIM', () => {
  it('ketiga keputusan', () => {
    expect(KEPUTUSAN_POLIS.confirm).toBe('Confirm')
    expect(KEPUTUSAN_POLIS.reject).toBe('Reject')
    expect(KEPUTUSAN_POLIS.decline).toBe('Decline')
  })
})

describe('periode produksi mm/yyyy', () => {
  it('YYYY-MM dari server tampil sebagai MM/YYYY', async () => {
    const { periodeTampil } = await import('./InputOffer')
    expect(periodeTampil('2026-10')).toBe('10/2026')
    expect(periodeTampil('2026-3')).toBe('03/2026')
  })

  it('bentuk asing ditampilkan apa adanya', async () => {
    const { periodeTampil } = await import('./InputOffer')
    expect(periodeTampil('Q4 2026')).toBe('Q4 2026')
    expect(periodeTampil('')).toBe('')
  })
})

describe('judul halaman', () => {
  it('tahap penawaran berjudul "Input Offer Life", tahap detail "Input Premium Detail"', async () => {
    const { judulKeputusan } = await import('./InputOffer')
    expect(judulKeputusan(TAHAP_POLIS.penawaran)).toBe('Input Offer Life')
    expect(judulKeputusan(TAHAP_POLIS.detail)).toBe('Input Premium Detail')
  })
})

describe('period di tahap Input Premium Detail (02-10-2026)', () => {
  it('panel Decision tidak lagi memuat period; ia tampil di kepala Premium List Detail', () => {
    const io = readFileSync(join(__dirname, 'InputOffer.tsx'), 'utf8')
    expect(io).not.toContain('pl-keputusan__meta')
    expect(io).toContain('{!diDetail && galatPeriode !== null')
  })
})

describe('konfirmasi tombol Decision (03-10-2026)', () => {
  it('tombol hanya membuka popup — tidak langsung memutuskan', () => {
    const io = readFileSync(join(__dirname, 'InputOffer.tsx'), 'utf8')
    expect(io).toContain('setTanya(KEPUTUSAN_POLIS.confirm)')
    expect(io).toContain('setTanya(KEPUTUSAN_POLIS.decline)')
    // Satu-satunya pemanggil putuskanPenawaran adalah tombol "Yes" di popup.
    expect(io.match(/putuskanPenawaran\(/g)?.length).toBe(1)
    expect(io).toContain('<Modal')
  })

  it('kalimatnya menyebut akibat per tahap', async () => {
    const { kalimatKonfirmasi } = await import('./InputOffer')
    const { KONFIRMASI_KEPUTUSAN } = await import('../labels')
    expect(kalimatKonfirmasi(KEPUTUSAN_POLIS.confirm, TAHAP_POLIS.detail)).toBe(KONFIRMASI_KEPUTUSAN.confirmDetail)
    expect(kalimatKonfirmasi(KEPUTUSAN_POLIS.confirm, TAHAP_POLIS.penawaran)).toBe(KONFIRMASI_KEPUTUSAN.confirmPenawaran)
    expect(kalimatKonfirmasi(KEPUTUSAN_POLIS.decline, TAHAP_POLIS.detail)).toBe(KONFIRMASI_KEPUTUSAN.decline)
    expect(KONFIRMASI_KEPUTUSAN.confirmDetail).toContain('PL Number')
  })
})

describe('PL Number ditampilkan sebelum kembali ke kotak masuk (03-10-2026)', () => {
  it('nomor yang terbit membuka popup; kembali ke kotak masuk hanya saat popup ditutup', () => {
    const io = readFileSync(join(__dirname, 'InputOffer.tsx'), 'utf8')
    expect(io).toContain('setNomorTerbit(hasil.plNumber')
    expect(io).toMatch(/if \(\(hasil\.plNumber \?\? ''\)\.trim\(\) !== ''\) \{[\s\S]{0,160}return\n/)
    expect(io).toMatch(/onTutup=\{\(\) => \{\n\s*setNomorTerbit\(null\)\n\s*onSelesai\(\)/)
    expect(io).toContain('HASIL_NOMOR_PL.judul')
  })
})

describe('Confirm terkunci selama ada perubahan belum disimpan (03-10-2026)', () => {
  it('form melapor, rute meneruskan, Confirm dikunci dan alasannya dikatakan', () => {
    const io = readFileSync(join(__dirname, 'InputOffer.tsx'), 'utf8')
    const form = readFileSync(join(__dirname, 'FormDataPolis.tsx'), 'utf8')
    const detail = readFileSync(join(__dirname, 'PremiumListDetail.tsx'), 'utf8')
    const rute = readFileSync(join(__dirname, '..', 'rute.tsx'), 'utf8')
    expect(form).toContain('JSON.stringify(isi) !== dasar')
    expect(form).toContain('onBelumTersimpan?.(belumTersimpan)')
    expect(detail).toContain('onBelumTersimpan={onBelumTersimpan}')
    expect(rute).toContain('onBelumTersimpan={setDataBelumTersimpan}')
    expect(rute).toContain('confirmTerkunci={polis.tahap === TAHAP_POLIS.detail && dataBelumTersimpan}')
    expect(io).toContain('disabled={sibuk || confirmTerkunci}')
    expect(io).toContain('KONFIRMASI_KEPUTUSAN.belumTersimpan')
  })
})

describe('WPC ditampilkan bersama PL Number (03-10-2026)', () => {
  it('popup hasil memuat PL Number dan WPC (dd/mm/yyyy)', () => {
    const io = readFileSync(join(__dirname, 'InputOffer.tsx'), 'utf8')
    expect(io).toContain("setWpcTerbit(hasil.wpc ?? '')")
    expect(io).toContain('{tanggalTampil(wpcTerbit)}')
    expect(io).toContain('HASIL_NOMOR_PL.labelWpc')
  })
})
