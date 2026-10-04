// Sensus tombol Endorsement Life - PARITAS bab 10: 22 sel `pxButton` di 14 Section.
//
// ⛔ Setiap tombol korpus berstatus DIBANGUN (kunci label yang dirujuk sebuah komponen `.tsx`) atau
// TIDAK DIBANGUN beserta alasannya. Tombol yang hilang diam-diam, label yang dikarang, atau kunci yang
// tidak pernah dirender menggagalkan uji. Korpus READ-ONLY; bila tak terjangkau, bagian korpus DILEWATI.

import { existsSync, readFileSync, readdirSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import * as LABEL from './labels'

const KORPUS = 'D:\\XML\\RNM_BRD\\Endorsement Life'
const adaKorpus = existsSync(KORPUS)

interface Tombol {
  berkas: string
  label: string
  baris: number
  /** Kunci `OBJEK.medan` di `labels.ts`, atau alasan tidak dibangun. */
  kunci?: string
  alasan?: string
}

export const SENSUS: readonly Tombol[] = [
  { berkas: 'ConfirmSubmitEDM', label: 'Close', baris: 1366, kunci: 'TERIMA_EDM.close' },
  { berkas: 'EndorsmentLife_Section', label: 'Process Policy No', baris: 2414, alasan: 'VIS=never' },
  { berkas: 'EndorsmentLife_Section', label: 'Submit', baris: 4226, kunci: 'BUAT_EDM.submit' },
  { berkas: 'InboxEndorsementLife', label: 'Create', baris: 3780, alasan: 'R11 panel tak terjangkau (InData.CARI1==1)' },
  { berkas: 'InboxEndorsementLife', label: 'Cancel', baris: 4228, alasan: 'R11 panel tak terjangkau' },
  { berkas: 'InboxEndorsementLife', label: 'Create Case Endorsement', baris: 5879, alasan: 'VIS=never' },
  { berkas: 'InboxEndorsementLife', label: 'Create Addendum', baris: 6620, kunci: 'INBOX_EDM.createAddendum' },
  { berkas: 'InputEDMLife', label: 'Process Policy No', baris: 5233, alasan: 'VIS=never' },
  { berkas: 'InputEDMLife', label: 'Upload CSV', baris: 8973, kunci: 'UNGGAH_EDM.uploadCsv' },
  { berkas: 'InputEDMLife', label: 'View Upload', baris: 9340, kunci: 'UNGGAH_EDM.viewUpload' },
  { berkas: 'InputEDMLife', label: 'Add CSV Data', baris: 10405, kunci: 'UNGGAH_EDM.addCsvData' },
  { berkas: 'InputEDMLife', label: 'DELETE ALL', baris: 13607, kunci: 'SIMPAN_EDM.deleteAll' },
  { berkas: 'InputEDMLife', label: 'View Premium', baris: 36494, alasan: 'VIS 1=2' },
  { berkas: 'InputEDMLife', label: 'Save', baris: 37202, kunci: 'SIMPAN_EDM.save' },
  { berkas: 'InputEDMLife', label: 'Submit', baris: 37494, kunci: 'PUTUSAN_EDM.submit' },
  { berkas: 'InputEDMLife', label: 'Submit', baris: 38109, kunci: 'PUTUSAN_EDM.submit' },
  { berkas: 'ShowLifePremiumSummary_EDM', label: 'View Old Policy', baris: 64965, kunci: 'POLIS_LAMA_EDM.viewOldPolicy' },
  { berkas: 'ShowLifePremiumSummary_EDM', label: 'View Old Policy', baris: 65522, kunci: 'POLIS_LAMA_EDM.viewOldPolicy' },
  { berkas: 'ShowLifePremiumSummary_EDM', label: 'View Old Policy', baris: 66083, kunci: 'POLIS_LAMA_EDM.viewOldPolicy' },
  { berkas: 'ShowLifePremiumSummary_EDM', label: 'View Old Policy', baris: 66640, kunci: 'POLIS_LAMA_EDM.viewOldPolicy' },
  { berkas: 'ShowLifePremiumSummary_EDM', label: 'Submit', baris: 67733, alasan: 'R02 layar yatim; efeknya = jalur Confirm' },
  { berkas: 'ViewCSVResult_LifeEDM', label: 'Generate Data Detail', baris: 14322, kunci: 'UNGGAH_EDM.generateDataDetail' },
]

function nilaiLabel(kunci: string): string | undefined {
  const [objek, medan] = kunci.split('.') as [string, string]
  return (LABEL as unknown as Record<string, Record<string, string> | undefined>)[objek]?.[medan]
}

/** Seluruh sumber `.tsx` modul ini. */
function sumberTsx(): string {
  const akar = join(__dirname)
  const hasil: string[] = []
  const jelajah = (dir: string) => {
    for (const n of readdirSync(dir, { withFileTypes: true })) {
      const j = join(dir, n.name)
      if (n.isDirectory()) jelajah(j)
      else if (n.name.endsWith('.tsx')) hasil.push(readFileSync(j, 'utf8'))
    }
  }
  jelajah(akar)
  return hasil.join('\n')
}

describe('sensus tombol Endorsement Life', () => {
  it('22 tombol; setiap tombol dibangun ATAU beralasan, tidak keduanya', () => {
    expect(SENSUS).toHaveLength(22)
    for (const t of SENSUS) expect(t.kunci === undefined, `${t.berkas} b${t.baris}`).not.toBe(t.alasan === undefined)
  })

  it('tombol yang dibangun: label VERBATIM dan dirujuk komponen', () => {
    const tsx = sumberTsx()
    for (const t of SENSUS.filter((x) => x.kunci !== undefined)) {
      const kunci = t.kunci as string
      expect(nilaiLabel(kunci), kunci).toBe(t.label)
      expect(tsx.includes(kunci), `${kunci} tidak dirender di komponen mana pun`).toBe(true)
    }
  })
})

describe.skipIf(!adaKorpus)('sensus tombol terhadap korpus', () => {
  const baris = new Map<string, string[]>()
  const baca = (berkas: string): string[] => {
    if (!baris.has(berkas)) baris.set(berkas, readFileSync(join(KORPUS, 'Section', `${berkas}.xml`), 'utf8').split('\n'))
    return baris.get(berkas) ?? []
  }

  it('cacah pxButton per Section sama dengan sensus (cara 1 bab 0)', () => {
    const korpus = new Map<string, number>()
    for (const f of readdirSync(join(KORPUS, 'Section')).filter((n) => n.endsWith('.xml'))) {
      const c = readFileSync(join(KORPUS, 'Section', f), 'utf8')
        .split('\n')
        .filter((l) => l.trim() === '<pyFormat>pxButton</pyFormat>').length
      if (c > 0) korpus.set(f.replace(/\.xml$/, ''), c)
    }
    const sensus = new Map<string, number>()
    for (const t of SENSUS) sensus.set(t.berkas, (sensus.get(t.berkas) ?? 0) + 1)
    expect(Object.fromEntries(sensus)).toEqual(Object.fromEntries(korpus))
  })

  it.each(SENSUS.map((t) => [t.berkas, t.baris, t.label]))('%s b%i = %s', (berkas, nomor, label) => {
    expect(baca(berkas as string)[(nomor as number) - 1]?.trim()).toBe(`<pyLabel>${label as string}</pyLabel>`)
  })
})
