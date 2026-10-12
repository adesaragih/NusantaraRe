import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { BahasaUI } from '../../../../inti/frontend/components/ui/bahasaUI'
import { kataStatusBaris, type Klaim, type PratinjauKomite } from '../api'
import { bolehKirimKomite, hitungRemarks, KonfirmasiKomite, POPUP_KOMITE, tanggalWaktuKomite } from './KonfirmasiKomite'
import { PESAN_KOMITE_TERKIRIM, statusBarisKomite } from './RincianPeserta'

// Pop-up Send Claim to Committee (keputusan work owner 12-10-2026 opsi A).

const en = (n: React.ReactNode) => renderToStaticMarkup(<BahasaUI.Provider value="en">{n}</BahasaUI.Provider>)
const PRATINJAU: PratinjauKomite = {
  waktu: new Date(2026, 9, 12, 9, 5).toISOString(),
  pic: 'UJI-SPV',
  jumlahKlaim: '1500000',
  mataUang: 'IDR',
  anggota: [
    { namaKomite: 'UJI Associate Director', urut: 1 },
    { namaKomite: 'UJI Director', urut: 2 },
  ],
}
const klaim = (isKpr: string) => ({ id: 'K1', nomorKlaim: 'UJI-CLM-1', isKpr, peserta: [] }) as unknown as Klaim
const modal = (isKpr = 'KPR', awal: PratinjauKomite = PRATINJAU) =>
  en(<KonfirmasiKomite klaim={klaim(isKpr)} pesertaID="P1" adjID="A1" onTutup={() => {}} onTerkirim={() => {}} awal={awal} />)

describe('pop-up Send Claim to Committee', () => {
  it('Date (waktu server, baca saja) · PIC · Adjustment · Remarks* · grid pratinjau "Waiting"', () => {
    const html = modal()
    expect(html).toContain('<dt>Date</dt><dd>12/10/2026 09:05</dd>')
    expect(html).toContain('<dt>PIC</dt><dd>UJI-SPV</dd>')
    expect(html).toMatch(/<dt>Adjustment<\/dt><dd>1\.500\.000 IDR<\/dd>/)
    expect(html).toMatch(/Remarks <span aria-hidden="true">\*<\/span><\/span><textarea[^>]*required=""[^>]*maxLength="4000"/)
    expect(html).not.toMatch(/<input[^>]*type="(date|datetime-local)"/)
    for (const k of POPUP_KOMITE.kolom) expect(html).toContain(`<th>${k}</th>`)
    expect(html).toContain('<td>UJI Associate Director</td><td><span class="komite-popup__chip">Waiting</span></td><td>—</td><td>—</td>')
    expect(html).toContain('<td>UJI Director</td><td><span class="komite-popup__chip">Waiting</span></td>')
  })

  it('tombol kirim nonaktif saat Remarks kosong (keadaan awal); tanpa konfirmasi kedua', () => {
    const html = modal()
    expect(html).toMatch(/<button[^>]*disabled=""[^>]*>Send Claim to Committee<\/button>/)
    expect(html).toContain('>Cancel<')
    expect(html).not.toContain('Are you sure')
  })

  it('aturan XML: roster ada + Remarks terisi + Is KPR terisi', () => {
    expect(bolehKirimKomite(PRATINJAU, 'UJI alasan', 'KPR')).toBe(true)
    expect(bolehKirimKomite(PRATINJAU, '   ', 'KPR')).toBe(false)
    expect(bolehKirimKomite(PRATINJAU, 'UJI alasan', '')).toBe(false)
    expect(bolehKirimKomite({ ...PRATINJAU, anggota: [] }, 'UJI alasan', 'KPR')).toBe(false)
    expect(bolehKirimKomite(null, 'UJI alasan', 'KPR')).toBe(false)
  })

  it('Is KPR kosong: keterangan mengapa tombol nonaktif', () => {
    expect(modal('')).toContain(POPUP_KOMITE.tanpaKpr)
    expect(modal('KPR')).not.toContain(POPUP_KOMITE.tanpaKpr)
  })

  it('tanggal dd/mm/yyyy HH:mm; tak terbaca = —', () => {
    expect(tanggalWaktuKomite(new Date(2026, 0, 3, 7, 4).toISOString())).toBe('03/01/2026 07:04')
    expect(tanggalWaktuKomite('bukan')).toBe('—')
  })

  it('satu pemanggil serahkanKeKomite, klik ganda ditahan, galat di PitaGalat modal; sukses → pita + muat ulang', () => {
    const src = readFileSync(join(__dirname, 'KonfirmasiKomite.tsx'), 'utf8')
    expect(src.split('await serahkanKeKomite(').length - 1).toBe(1)
    expect(src).toContain('if (sibuk || !bolehKirimKomite(pratinjau, remarks, klaim.isKpr)) return')
    expect(src).toContain('tanpaTutup={sibuk}')
    expect(src).toContain('<PitaGalat galat={galat} onTutup={() => setGalat(null)} />')
    const rincian = readFileSync(join(__dirname, 'RincianPeserta.tsx'), 'utf8')
    expect(rincian).toContain('setPesanKomite(PESAN_KOMITE_TERKIRIM)')
    expect(rincian).toContain('await onMuatUlang()')
    expect(rincian).not.toContain('serahkanKeKomite(')
    expect(PESAN_KOMITE_TERKIRIM).toBe('Claim sent to committee.')
  })

  it('baris adjustment menampilkan status komite sesudah diserahkan', () => {
    expect(statusBarisKomite({ kodeStatus: '0', komiteId: '' })).toBe(kataStatusBaris('0'))
    expect(statusBarisKomite({ kodeStatus: '0', komiteId: 'KMTLF-1' })).toBe(`${kataStatusBaris('0')} · In Committee (KMTLF-1)`)
  })
})

describe('pop-up saat pratinjau ditolak (12-10-2026)', () => {
  it('konteks tetap tampil di bawah pita galat: Date, PIC —, Adjustment dari baris, Remarks nonaktif, grid kosong', () => {
    const html = en(
      <KonfirmasiKomite
        klaim={klaim('KPR')}
        pesertaID="P1"
        adjID="A1"
        jumlahKlaim="2500000"
        mataUang="IDR"
        galatAwal="No committee level covers this claim amount."
        onTutup={() => {}}
        onTerkirim={() => {}}
      />,
    )
    const iGalat = html.indexOf('No committee level covers this claim amount.')
    expect(iGalat).toBeGreaterThan(0)
    expect(html.indexOf('<dt>Date</dt>')).toBeGreaterThan(iGalat)
    expect(html).toMatch(/<dt>Date<\/dt><dd>\d{2}\/\d{2}\/\d{4} \d{2}:\d{2}<\/dd>/)
    expect(html).toContain('<dt>PIC</dt><dd>—</dd>')
    expect(html).toMatch(/<dt>Adjustment<\/dt><dd>2\.500\.000 IDR<\/dd>/)
    expect(html).toMatch(/<textarea[^>]*disabled=""/)
    expect(html).toContain('<td colSpan="4">No committee members found.</td>')
    expect(html).toMatch(/<button[^>]*disabled=""[^>]*>Send Claim to Committee<\/button>/)
  })
})

describe('tampilan pop-up dirapikan (12-10-2026)', () => {
  it('info tiga kolom, Remarks 4 baris + penghitung "0 / 4000", judul Committee di atas tabel penuh lebar', () => {
    const html = modal()
    expect(html).toContain('<dl class="komite-popup__info">')
    expect(html).toMatch(/<textarea class="field__input" rows="4"/)
    expect(html).toContain('<span class="komite-popup__hitung" aria-live="polite">0 / 4000</span>')
    expect(html).toContain('<h4 class="komite-popup__judul">Committee</h4><div class="komite-popup__kartu"><table class="inbox__tabel komite-popup__tabel">')
    expect(html).not.toContain('inbox__gulir')
    // "Remarks is required." baru tampil sesudah disentuh (blur), bukan saat modal dibuka.
    expect(html).not.toContain('Remarks is required.')
  })
  it('kaki: teks bantuan "Fill in Remarks to continue." di sebelah tombol nonaktif; Cancel tetap kiri', () => {
    const html = modal()
    expect(html).toMatch(/>Cancel<\/button><span class="komite-popup__bantu">Fill in Remarks to continue\.<\/span><button[^>]*disabled=""[^>]*>Send Claim to Committee<\/button>/)
  })
  it('penghitung karakter', () => {
    expect(hitungRemarks('')).toBe('0 / 4000')
    expect(hitungRemarks('UJI alasan')).toBe('10 / 4000')
  })
  it('CSS: modal ≈720px, jarak 16px, tabel tanpa bayangan terpotong, chip, satu kolom di layar sempit', () => {
    const css = readFileSync(join(__dirname, '..', 'claimlife.css'), 'utf8')
    for (const p of ['.claimlife .modal:has(.komite-popup) {\n  width: min(720px, 92vw);', '.claimlife .komite-popup {\n  display: grid;\n  gap: 16px;',
      '.claimlife .komite-popup .komite-popup__tabel {\n  width: 100%;\n  box-sizing: border-box;\n  margin: 0;\n  box-shadow: none;',
      '.claimlife .komite-popup__chip {', '.claimlife .komite-popup__info {\n    grid-template-columns: minmax(0, 1fr);'])
      expect(css, p).toContain(p)
  })
})

describe('Adjustment tanpa kata Claim RNM; kartu tabel tidak terpotong (12-10-2026)', () => {
  it('Adjustment = nominal + mata uang saja', () => {
    expect(modal()).not.toContain('Claim RNM')
  })
  it('kartu tabel = pembungkus sendiri (tepi, sudut, overflow hidden); tabel polos; ruang bawah 4px', () => {
    const css = readFileSync(join(__dirname, '..', 'claimlife.css'), 'utf8')
    expect(css).toContain('.claimlife .komite-popup__kartu {\n  width: 100%;\n  box-sizing: border-box;\n  overflow: hidden;')
    expect(css).toContain('.claimlife .komite-popup__kartu .komite-popup__tabel {\n  border: 0;\n}')
    expect(css).toContain('  padding-bottom: 4px;')
  })
})
