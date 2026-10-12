import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { BarisPerbandinganEDM } from '../api'
import { MATA_UANG } from '../kolomKorpus'
import { PERBANDINGAN_EDM } from '../labels'
import { arahSelisih, kelasNegatif, selBertanda, teksPeserta } from '../tampilan'
import { kolomPerbandingan, TabelPerbandingan } from './PerbandinganRekap'
import { TabelKorpus } from './TabelKorpus'

// Summary Comparison - permintaan work owner 08-10-2026 (bukan korpus); tampilan dirapikan 08-10-2026.

const sel = (cacah: number, premi: string, balance: string) => ({ cacah, nilai: { PREMIUM: premi, BALANCE: balance } })
const IDR: BarisPerbandinganEDM = {
  currency: 'IDR',
  before: sel(10, '358585.9984', '179292.9992'),
  deleted: sel(-1, '-22815.1511', '-11407.5755'),
  added: sel(3, '2047513.5625', '150'),
  after: sel(12, '2383284.4098', '167885.4237'),
  difference: sel(2, '2024698.4114', '-11257.5755'),
}
const USD: BarisPerbandinganEDM = { ...IDR, currency: 'USD' }
const KOLOM = kolomPerbandingan('QR').filter((k) => k.kolom === 'PREMIUM' || k.kolom === 'BALANCE')
const S = readFileSync(join(__dirname, '..', 'pages', 'InputEDMLife.tsx'), 'utf8').replace(/\r\n/g, '\n')

function baris(html: string, label: string): string {
  const i = html.indexOf(`<td>${label}`)
  return html.slice(html.lastIndexOf('<tr', i), html.indexOf('</tr>', i))
}

describe('Summary Comparison dirapikan', () => {
  it('kolom uang = kolom rekap per Type tanpa CURRENCY', () => {
    expect(kolomPerbandingan('QR').map((k) => k.kolom)).toEqual(['PREMIUM', 'DEDUCTION', 'BROKERAGE_FEE', 'BALANCE'])
    expect(kolomPerbandingan('TP').map((k) => k.kolom)).toEqual(['PREMIUM', 'DISCOUNT_PREMIUM_RETRO', 'SHARE_RETRO', 'BALANCE'])
  })

  it('satu blok bersub-judul per mata uang, tanpa kolom Currency; kepala angka rata kanan', () => {
    const html = renderToStaticMarkup(<TabelPerbandingan baris={[IDR, USD]} kolom={KOLOM} />)
    expect(html.match(/<h4 class="edm-banding__uang">/g)?.length).toBe(2)
    expect(html.indexOf('>IDR</h4>')).toBeLessThan(html.indexOf('>USD</h4>'))
    expect(html).not.toContain('rowspan')
    expect(html).not.toContain('>Currency<')
    const kepala = html.slice(html.indexOf('<thead>'), html.indexOf('</thead>'))
    expect(kepala).toContain(`<th>${PERBANDINGAN_EDM.kolomBaris}</th>`)
    expect(kepala.match(/<th class="edm-angka">/g)?.length).toBe(1 + KOLOM.length)
  })

  it('urutan baris; After tebal; Difference disorot, bertanda, berwarna', () => {
    const html = renderToStaticMarkup(<TabelPerbandingan baris={[IDR]} kolom={KOLOM} />)
    const urut = [PERBANDINGAN_EDM.before, PERBANDINGAN_EDM.deleted, PERBANDINGAN_EDM.added, PERBANDINGAN_EDM.after, PERBANDINGAN_EDM.difference]
    const posisi = urut.map((u) => html.indexOf(`<td>${u}`))
    expect(posisi.every((p, i) => p > 0 && (i === 0 || p > (posisi[i - 1] ?? 0)))).toBe(true)
    expect(PERBANDINGAN_EDM.difference).toBe('Difference (endorsement premium)')
    expect(baris(html, PERBANDINGAN_EDM.after)).toContain('class="edm-banding__total"')
    const selisih = baris(html, PERBANDINGAN_EDM.difference)
    expect(selisih).toContain('class="edm-banding__selisih"')
    expect(selisih).toContain('<span class="edm-banding__naik">+2.024.698,4114</span>')
    expect(selisih).toContain('<td class="edm-angka edm-negatif">-11.257,5755</td>')
    expect(selisih).toContain('<td class="edm-angka"><span class="edm-banding__naik">+2.024.698,4114</span></td>')
    // Nilai uang Deleted tetap negatif apa adanya; format id.
    expect(baris(html, PERBANDINGAN_EDM.deleted)).toContain('>-22.815,1511<')
    expect(html).toContain('>358.585,9984<')
  })

  it('Participants: jumlah orang tanpa tanda dan tanpa warna di semua baris (09-10-2026)', () => {
    const html = renderToStaticMarkup(<TabelPerbandingan baris={[IDR]} kolom={KOLOM} />)
    const mau: Record<string, string> = {
      [PERBANDINGAN_EDM.before]: '10',
      [PERBANDINGAN_EDM.deleted]: '1',
      [PERBANDINGAN_EDM.added]: '3',
      [PERBANDINGAN_EDM.after]: '12',
      [PERBANDINGAN_EDM.difference]: '2',
    }
    for (const [label, n] of Object.entries(mau)) {
      const sel = [...baris(html, label).matchAll(/<td[^>]*>([^<]*)/g)]
      // sel[0] = label, sel[1] = Participants
      expect(sel[1]?.[0], label).toBe(`<td class="edm-angka">${n}`)
    }
    expect(teksPeserta(-1)).toBe('1')
    expect(teksPeserta(0)).toBe('0')
    expect(teksPeserta(9)).toBe('9')
    // Kolom uang tidak berubah: Deleted tetap negatif merah.
    expect(baris(html, PERBANDINGAN_EDM.deleted)).toContain('<td class="edm-angka edm-negatif">-22.815,1511</td>')
  })

  it('tanpa catatan rumus (dihapus 08-10-2026); tanpa cabang "belum diproses" (10-10-2026)', () => {
    const html = renderToStaticMarkup(<TabelPerbandingan baris={[IDR]} kolom={KOLOM} />)
    expect(html).not.toContain('edm-banding__rumus')
    expect(html).not.toContain('Before = existing policy')
    expect(html).not.toContain('edm-banding__catatan')
    expect(html.trimEnd().endsWith('</div>')).toBe(true)
    // Deleted selalu berangka (tabel hanya tampil sesudah proses).
    expect(baris(html, PERBANDINGAN_EDM.deleted)).not.toContain('>—<')
  })

  it('selBertanda / arahSelisih', () => {
    expect(selBertanda('150')).toBe('+150')
    expect(selBertanda('-50.5')).toBe('-50,5')
    expect(selBertanda('0')).toBe('0')
    expect([arahSelisih('5'), arahSelisih('-5'), arahSelisih('0'), arahSelisih('0.000')]).toEqual(['naik', 'turun', 'nol', 'nol'])
  })
})

describe('bagian bawah layar kasus', () => {
  it('urutan: Process Endorsement → Summary Comparison → rekap berjudul → DECISION berjudul', () => {
    const iProses = S.indexOf('<div className="edm-kasus__proses">')
    const iBanding = S.indexOf('<PerbandinganRekap kasusId={kasusId}')
    const iRekap = S.indexOf('<section className="edm-rekap">')
    const iKeputusan = S.indexOf('<section className="edm-keputusan">')
    expect(iProses > 0 && iProses < iBanding && iBanding < iRekap && iRekap < iKeputusan).toBe(true)
    // Hanya sesudah Process Endorsement (keputusan work owner 10-10-2026).
    expect(S).toContain("{kasus.sudahSimpan && <PerbandinganRekap kasusId={kasusId} tipe={k.TYPE ?? ''} versi={versiBanding} />}")
    expect(S.split('<PerbandinganRekap').length - 1).toBe(1)
    expect(S.split('setVersiBanding((v) => v + 1)').length - 1).toBe(2)
  })

  it('tabel rekap berjudul tanpa keterangan (dihapus 08-10-2026); DECISION berjudul edm-subjudul', () => {
    const rekap = S.slice(S.indexOf('<section className="edm-rekap">'), S.indexOf('</section>', S.indexOf('<section className="edm-rekap">')))
    expect(rekap).toContain('<h3 className="edm-subjudul">{PERBANDINGAN_EDM.judulRekap}</h3>')
    expect(rekap).not.toContain('keterangan')
    expect(rekap).not.toContain('edm-catatan')
    expect(PERBANDINGAN_EDM.judul).toBe('Summary Comparison')
    expect(PERBANDINGAN_EDM.judulRekap).toBe('Summary to be posted on Confirm')
    const putusan = S.slice(S.indexOf('<section className="edm-keputusan">'))
    expect(putusan.indexOf('<h3 className="edm-subjudul">{PERBANDINGAN_EDM.judulKeputusan}</h3>')).toBeLessThan(putusan.indexOf('<KeputusanEDM'))
    expect(PERBANDINGAN_EDM.judulKeputusan).toBe('Decision')
  })

  it('kepala angka tabel rekap rata kanan (TabelKorpus jenis n)', () => {
    const html = renderToStaticMarkup(<TabelKorpus kolom={MATA_UANG.QR} baris={[{ CURRENCY: 'IDR', PREMIUM: '1' }]} kunci={(_, i) => String(i)} />)
    const kepala = html.slice(html.indexOf('<thead>'), html.indexOf('</thead>'))
    expect(kepala).toContain('<th>CURRENCY</th>')
    expect(kepala.match(/<th class="edm-angka">/g)?.length).toBe(MATA_UANG.QR.length - 1)
  })
})

describe('nilai negatif merah (08-10-2026)', () => {
  it('kelasNegatif dari tanda teks (tanpa parseFloat); -0 / nol / positif tidak', () => {
    for (const v of ['-1', '-22815.1511', ' -0.0001 ', -1]) expect(kelasNegatif(v), String(v)).toBe('edm-negatif')
    for (const v of ['0', '-0', '-0.000', '1', '+5', '', undefined, 0, 3]) expect(kelasNegatif(v), String(v)).toBeUndefined()
  })

  it('Summary Comparison: SEMUA sel uang negatif merah di baris mana pun; positif/nol tidak; Participants tak bertanda', () => {
    const html = renderToStaticMarkup(<TabelPerbandingan baris={[IDR]} kolom={KOLOM} />)
    const sel = [...html.matchAll(/<td class="([^"]*)">(?:<span[^>]*>)?([^<]*)/g)].map((m) => ({ kelas: m[1] ?? '', isi: m[2] ?? '' }))
    expect(sel.length).toBeGreaterThan(10)
    for (const { kelas, isi } of sel) {
      expect(kelas.includes('edm-negatif'), `${isi} ${kelas}`).toBe(isi.startsWith('-'))
    }
    // Before negatif pun merah (baris mana pun).
    const neg: BarisPerbandinganEDM = { ...IDR, before: sel2(-1, '-5', '0') }
    const h2 = renderToStaticMarkup(<TabelPerbandingan baris={[neg]} kolom={KOLOM} />)
    expect(baris(h2, PERBANDINGAN_EDM.before)).toContain('<td class="edm-angka edm-negatif">-5</td>')
    expect(baris(h2, PERBANDINGAN_EDM.before)).toContain('<td class="edm-angka">0</td>')
  })

  it('tabel rekap: angka negatif merah (tandaiNegatif); tabel korpus lain tidak', () => {
    const b = [{ CURRENCY: 'IDR', PREMIUM: '-312955.6962', DEDUCTION: '0', BROKERAGE_FEE: '5', BALANCE: '-1' }]
    const html = renderToStaticMarkup(<TabelKorpus kolom={MATA_UANG.QR} baris={b} kunci={(_, i) => String(i)} tandaiNegatif />)
    expect(html).toContain('<td class="edm-angka edm-negatif">-312.955,6962</td>')
    expect(html).toContain('<td class="edm-angka edm-negatif">-1</td>')
    expect(html).toContain('<td class="edm-angka">0</td>')
    expect(html).toContain('<td class="edm-angka">5</td>')
    expect(renderToStaticMarkup(<TabelKorpus kolom={MATA_UANG.QR} baris={b} kunci={(_, i) => String(i)} />)).not.toContain('edm-negatif')
    expect(S).toContain('tandaiNegatif />')
  })
})

function sel2(cacah: number, premi: string, balance: string) {
  return { cacah, nilai: { PREMIUM: premi, BALANCE: balance } }
}
