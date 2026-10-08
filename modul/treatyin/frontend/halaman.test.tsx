// Uji PENAMPUNG HALAMAN `TreatyIn` (`halaman.tsx`) — permintaan pemakai
// 7 Oktober 2026: isian tab tidak hilang saat pindah tab, dan tab yang
// saling bergantung membaca nilai yang sama.
//
// ⚠️ Render statis (`react-dom/server`) — nol DOM, nol klik. "Pindah tab"
// diuji sebagai yang sebenarnya terjadi di React: tab DILEPAS lalu DIRENDER
// ULANG; yang membuktikan isiannya bertahan adalah tab yang baru dirender
// menampilkan isi penampung, bukan nilai kontrak yang dimuat (`baris`/`pohon`).
// Render statis BUKAN pengganti melihat layar di peramban.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import type { ReactNode } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { SimpulLimit } from './api'
import TabAkumulasi, { barisAkumulasiBaru } from './components/TabAkumulasi'
import TabAngsuran, { angsuranDariWarisan } from './components/TabAngsuran'
import TabEgnpi from './components/TabEgnpi'
import TabEventLimits from './components/TabEventLimits'
import TabRetensi from './components/TabRetensi'
import TabCoInsScale from './components/TabCoInsScale'
import TabInfoSubmit from './components/TabInfoSubmit'
import TabLimitsProp from './components/TabLimitsProp'
import TabPortofolio from './components/TabPortofolio'
import TabReportingPeriod from './components/TabReportingPeriod'
import TabShareProp from './components/TabShareProp'
import TabTeksPanjang from './components/TabTeksPanjang'
import { bacaProperti, PenyediaHalaman, PROPERTI_TAB_PROP, setelProperti, type Halaman } from './halaman'

const AKAR = __dirname
const sumber = (f: string) => readFileSync(join(AKAR, f), 'utf8')

/** Render satu tab di atas halaman yang sudah berisi — seperti sesudah kembali dari tab lain. */
function diAtas(halaman: Halaman, isi: ReactNode): string {
  return renderToStaticMarkup(
    <PenyediaHalaman penampung={{ halaman, ubah: () => undefined }}>{isi}</PenyediaHalaman>,
  )
}

const LIMITS: SimpulLimit[] = [
  {
    TreatyType: 'SPECIAL SURPLUS',
    Detail: [{ TreatyGroup: 'ENGINEERING', RNMShare: '8.42', ShareNote: ' of 100%', CessionList: [{ Currency: 'IDR', Value: '100' }] }],
  },
]

describe('fungsi murni penampung', () => {
  it('setelProperti mengganti SATU properti; yang lain utuh; masukan tidak berubah', () => {
    const h: Halaman = { Portfolio: [], ReportingStart: '20250101' }
    const h2 = setelProperti(h, 'ReportingStart', '20250201')
    expect(h2).toEqual({ Portfolio: [], ReportingStart: '20250201' })
    expect(h.ReportingStart).toBe('20250101')
    expect(bacaProperti(h2, 'ReportingEnd')).toBeUndefined()
  })

  it('⭐ properti yang dipakai DUA tab adalah jalur ketergantungannya (dari ekspor)', () => {
    const pemakai = (p: string) => Object.entries(PROPERTI_TAB_PROP).filter(([, ps]) => ps.includes(p)).map(([t]) => t)
    expect(pemakai('Limits')).toEqual(['Limits', 'Share', 'Achievement In IDR'])
    expect(pemakai('ReportingStart')).toEqual(['Reporting Period', 'Accumulation'])
    expect(pemakai('ReportingEnd')).toEqual(['Reporting Period', 'Accumulation'])
  })
})

describe('⭐ isian bertahan saat pindah tab — tab dirender ulang dari penampung', () => {
  it('Portfolio: baris yang diketik, bukan baris kontrak', () => {
    const html = diAtas(
      { Portfolio: [{ TypePortfolio: 'Withdrawal', Type: 'Premium', Description: 'diketik sebelum pindah tab' }] },
      <TabPortofolio baris={[]} mode="lihat" petunjukKosong="" />,
    )
    expect(html).toContain('diketik sebelum pindah tab')
  })

  it('Co-Ins Scale: grid dan kedua Max Co-Insurance', () => {
    const html = diAtas(
      { CoInScale: [{ CoInShare: 'A', PctLimit: '12.5' }], MaxCoNonGroup: '7', MaxCoGroup: '9' },
      <TabCoInsScale baris={[]} mode="ubah" />,
    )
    expect(html).toContain('value="A"')
    expect(html).toContain('value="12.5"')
    expect(html).toContain('value="7"')
    expect(html).toContain('value="9"')
  })

  it('Exclusions / Special Conditions: properti per cabang (ExclusionsP / SpecialConditionsP)', () => {
    const h = { ExclusionsP: 'teks pengecualian', SpecialConditionsP: 'teks syarat' }
    expect(diAtas(h, <TabTeksPanjang judul="Exclusions" properti="ExclusionsP" tab={undefined} petunjukKosong="" mode="ubah" />)).toContain(
      'teks pengecualian',
    )
    expect(
      diAtas(h, <TabTeksPanjang judul="Special Conditions" properti="SpecialConditionsP" tab={undefined} petunjukKosong="" mode="ubah" />),
    ).toContain('teks syarat')
  })

  it('Information & Submit: Information dan Comment', () => {
    const html = diAtas({ Information: 'info tambahan', Comment: 'komentar' }, <TabInfoSubmit mode="ubah" statusAkseptasi="" />)
    expect(html).toContain('info tambahan')
    expect(html).toContain('komentar')
  })

  it('Reporting Period: medan kepala dan ReportingPeriodList (bentuk simpan → tampil)', () => {
    const h = {
      ReportingSubmission: '30',
      ReportingPeriodList: [{ Period: 'Q 1', AutoCalculate: '', InitialDate: '20250101', SubmissionDue: '20250131', ConfirmationDue: '', SettlementDue: '' }],
    }
    const ubah = diAtas(h, <TabReportingPeriod baris={[]} mode="ubah" />)
    expect(ubah).toContain('value="30"')
    expect(ubah).toContain('value="2025-01-31"')
    const lihat = diAtas(h, <TabReportingPeriod baris={[]} mode="lihat" />)
    expect(lihat).toContain('Q 1')
    expect(lihat).toContain('31/01/25')
  })

  it('Accumulation: AccumulationPeriod dan AccumulationList', () => {
    const html = diAtas(
      { AccumulationPeriod: 'quarter', AccumulationList: [{ Period: 'Q 1', ReportDate: '20250101', SubDays: '45', SubDueDate: '20250215' }] },
      <TabAkumulasi baris={[]} mode="lihat" />,
    )
    expect(html).toContain('Q 1')
    expect(html).toContain('01/01/25')
    expect(html).toContain('15/02/25')
    expect(html).toContain('<option value="quarter" selected="">')
  })

  it('Limits: layer yang ditambah di tab Limits, bukan pohon kontrak', () => {
    const html = diAtas({ Limits: LIMITS }, <TabLimitsProp pohon={[]} mode="lihat" opsi={{ jenisTreaty: [], kelompokTreaty: [], mataUang: [] }} />)
    expect(html).toContain('SPECIAL SURPLUS')
  })
})

describe('⭐ ketergantungan antartab — properti yang SAMA', () => {
  it('Share: grid Kind of Treaty = Limits yang diisi di tab Limits (pohon kontrak kosong)', () => {
    const html = diAtas({ Limits: LIMITS }, <TabShareProp pohon={[]} petunjukKosong="" mode="lihat" />)
    expect(html).toContain('SPECIAL SURPLUS')
    expect(html).not.toContain('(Kind of Treaty belum dipilih)')
  })

  it('Share: total hasil Refresh tinggal di halaman — tampil saat tab dibuka lagi', () => {
    const html = diAtas(
      {
        Limits: LIMITS,
        TotalShareRnmProp: [{ Currency: 'IDR', CurrencyID: '', Value: '3000000000' }],
        TotalSpreadedRnmProp: [{ Currency: 'IDR', CurrencyID: '', Value: '1200000000' }],
        TotalSpreadedRnmRIProp: [{ Currency: 'IDR', CurrencyID: '', Value: '1800000000' }],
      },
      <TabShareProp pohon={[]} petunjukKosong="" mode="lihat" />,
    )
    // Gambar Pega 17.
    for (const v of ['3.000.000.000,00', '1.200.000.000,00', '1.800.000.000,00']) expect(html).toContain(v)
  })

  it('Achievement In IDR dan Accumulation membaca properti tab lain dari penampung', () => {
    const ach = sumber('components/TabAchievement.tsx')
    expect(ach).toContain("useProperti<SimpulLimit[]>('Limits'")
    expect(ach).toContain('limits: limits.map(')
    const ak = sumber('components/TabAkumulasi.tsx')
    expect(ak).toContain("useProperti('ReportingStart', '')")
    expect(ak).toContain("useProperti('ReportingEnd', '')")
    expect(ak).toContain('ReportingStart: mulai, ReportingEnd: akhir')
    // Limits ditulis BALIK oleh Share (Refresh) ke properti yang sama.
    const sp = sumber('components/TabShareProp.tsx')
    expect(sp).toContain("useProperti<SimpulLimit[]>('Limits'")
    expect(sp).toContain('setLimits(h.Limits)')
  })

  it('⛔ penampung hidup di FORM, di luar tab yang dilepas — dan dikosongkan saat kontrak lain dimuat', () => {
    const form = sumber('pages/FormKontrakTreatyIn.tsx')
    const i = form.indexOf('<PenyediaHalaman penampung={penampung}>')
    expect(i).toBeGreaterThan(-1)
    expect(form.indexOf('<StripTab tab={tab}')).toBeLessThan(i)
    expect(form.match(/penampung\.kosongkan\(\)/g)?.length).toBe(2)
  })

  it('Accumulation Add: DataTransform TreatyInAddAccumulation', () => {
    expect(barisAkumulasiBaru('quarter', 2).Period).toBe('Q 3')
    expect(barisAkumulasiBaru('none', 0).Period).toBe('T 1')
    expect(barisAkumulasiBaru('other', 0).Period).toBe('1')
  })
})

describe('⭐ alur XML yang kini tersambung — 7 Oktober 2026', () => {
  // ⚠️ Pemicu `change` butuh DOM; berkas ini merender statis. Yang dijaga di
  // sini adalah SAMBUNGANNYA di sumber.
  it('Accumulation Period: DT pra-refresh TreatyInDeleteAccumulationLists — daftar dikirim KOSONG', () => {
    expect(sumber('components/TabAkumulasi.tsx')).toContain("jalankan('periode', [], v)")
  })

  it('Reporting Period: sel Initial Date memicu TreatyInSetReport HANYA bila Auto Calculate barisnya dicentang', () => {
    expect(sumber('components/TabReportingPeriod.tsx')).toContain(
      "if (kunci === 'InitialDate' && b.AutoCalculate === 'true') terapkan(keSimpan(v))",
    )
  })

  it('Exclusions/Special Conditions Prop: DT TreatyInCopyConditions saat isian ditinggalkan', () => {
    const s = sumber('components/TabTeksPanjang.tsx')
    expect(s).toContain("ExclusionsP: 'Exclusions'")
    expect(s).toContain("SpecialConditionsP: 'SpecialConditions'")
    expect(s).toContain('onBlur={salinSyarat}')
  })
})

describe('tanpa penyedia — perilaku tab sebelum penampung ada', () => {
  it('nilai awal dari kontrak yang dimuat', () => {
    const html = renderToStaticMarkup(
      <TabPortofolio baris={[{ jenis: 'Loss', jenisPortfolio: 'Assumption', keterangan: 'dari kontrak' }]} mode="lihat" petunjukKosong="" />,
    )
    expect(html).toContain('dari kontrak')
  })
})

// ---------------------------------------------------------------------
// ⭐ 7 Oktober 2026 (lanjutan) — "input tiba-tiba hilang" dan tab Non-Prop.
// ---------------------------------------------------------------------

describe('⛔ penyebab isian hilang — dicabut', () => {
  it('Retention dan EGNPI tidak lagi mereset baris setiap form dirender', () => {
    // `baris` dibuat ulang (`.map`) tiap render form → efek ini dahulu
    // mengembalikan baris ke data kontrak setiap kali apa pun berubah.
    const efekReset = /useEffect\(\(\) => \{\s*set(Rows|Total)\(\[\.\.\.(baris|totalAwal)\]\)/
    expect(sumber('components/TabRetensi.tsx')).not.toMatch(efekReset)
    expect(sumber('components/TabEgnpi.tsx')).not.toMatch(efekReset)
  })

  it('Limits Prop: hasil rumus asinkron diterapkan ke simpul TERKINI, bukan salinan saat dipicu', () => {
    const lp = sumber('components/TabLimitsProp.tsx')
    // LimitCalculation, CalculateDeduction, PremiumReserveCalculate.
    expect(lp.match(/terapkanKini\(\(kini\) =>/g)?.length).toBe(3)
    expect(lp).not.toContain('onUbah({ ...dBaru, DeductionList: h.DeductionList')
    expect(lp).toContain('setLimits((ls) => ganti(ls, i, f(ls[i] ?? {})))')
    expect(lp).toContain('setLimits((ls) => ganti(ls, i, baru))')
  })
})

describe('⭐ tab Non-Prop di penampung — bertahan saat pindah tab', () => {
  const opsi = { jenisTreaty: [], kelompokTreaty: [], mataUang: [] }

  it('Maximum Retention', () => {
    const html = diAtas(
      { Retention: [{ ID: '', TreatyGroup: 'FIRE', TreatyGroupID: '', Currency: 'IDR', CurrencyID: '', Amount: '3500000000', ClassOfBusiness: '', ClassOfBusinessID: '', Note: 'diketik' }] },
      <TabRetensi baris={[]} opsiAwal={opsi} mode="lihat" />,
    )
    expect(html).toContain('FIRE')
  })

  it('EGNPI', () => {
    const html = diAtas(
      {
        EGNPI: [{ ID: '', TreatyGroup: 'MARINE', TreatyGroupID: '', AsDate: '', Proportion: '', Currency: 'USD', CurrencyID: '', Amount: '2', AmountIDR: '', ClassOfBusiness: '', ClassOfBusinessID: '', Note: '' }],
        TotalEgnpiAmount: '200',
      },
      <TabEgnpi baris={[]} kurs={[]} retensi={[]} opsiAwal={opsi} mode="lihat" />,
    )
    expect(html).toContain('MARINE')
  })

  it('Event Limits — properti akar per kunci', () => {
    const html = diAtas({ RSMDLimit: '500000000000', CurrencyRSMD: 'IDR' }, <TabEventLimits mode="ubah" />)
    expect(html).toContain('value="500000000000"')
  })

  it('Installment — InstallmentNo, halaman per mata uang, dan total', () => {
    const html = diAtas(
      {
        InstallmentNo: '4',
        Installment: [{ Currency: 'IDR', AmountTotal: '644674819.59', PctTotal: '100', InstallmentList: [] }],
        TotalInstallmentNP: [{ Currency: 'IDR', CurrencyID: '', Value: '644674819.59' }],
      },
      <TabAngsuran baris={[]} netPremium={[]} mode="ubah" />,
    )
    expect(html).toContain('value="4"')
    expect(html).toContain('<td>IDR</td>')
    expect(html).toContain('644.674.819,59')
  })
})

describe('⭐ ketergantungan Non-Prop — tab tujuan membaca isian TERKINI tab sumber', () => {
  const form = sumber('pages/FormKontrakTreatyIn.tsx')

  it('Retention → EGNPI (mata uang baris baru dari Retention(1))', () => {
    expect(form).toContain("bacaProperti(penampung.halaman, 'Retention')")
    expect(form).toContain('retensi={retensiKini.map(')
  })

  it('EGNPI → Limits Non-Prop (`TotalEgnpi`)', () => {
    expect(form).toContain("bacaProperti(penampung.halaman, 'EGNPI')")
    expect(form).toContain('egnpi={egnpiKini.map(')
  })

  it('Share Non-Prop → Installment (`TotalShareNetNP`)', () => {
    expect(form).toContain('Total.TotalShareNetNP')
    expect(form).toContain('netPremium={netPremiumKini}')
  })

  it('Installment: jadwal kontrak dikelompokkan per mata uang (halaman `TreatyIn.Installment`)', () => {
    const a = angsuranDariWarisan([
      { angsuran: '1', mataUang: 'IDR', jumlah: '10', persen: '50', jatuhTempo: '', tanggalBayar: '', wpc: '', jatuhTempoAsli: '20250101', tanggalBayarAsli: '' },
      { angsuran: '1', mataUang: 'USD', jumlah: '1', persen: '100', jatuhTempo: '', tanggalBayar: '', wpc: '', jatuhTempoAsli: '', tanggalBayarAsli: '' },
      { angsuran: '2', mataUang: 'IDR', jumlah: '10', persen: '50', jatuhTempo: '', tanggalBayar: '', wpc: '', jatuhTempoAsli: '', tanggalBayarAsli: '' },
    ])
    expect(a.map((x) => `${x.Currency}:${String(x.InstallmentList.length)}`)).toEqual(['IDR:2', 'USD:1'])
    expect(a[0]?.InstallmentList[0]?.DueDate).toBe('20250101')
  })
})
