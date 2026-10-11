// Fungsi murni tab Coverage keluarga Aneka (tiket 73 P2) - operasi pohon Lokasi → Okupasi → Objek → Coverage yang di Pega
// dikerjakan activity / data transform di layar (hitung premi tetap di backend `POST …/hitung-coverage-aneka`).
//
// Sumber `[terverifikasi]`:
// - Add: `DDL\AddCoverageAneka_act.xml` L294-295 (metode dari `LocationList(1)…OccupationList(1).AnekaList(1).CoverageList(1)`,
//   indeks keras) dan L315-379 (Currency + TSI dari objek); nilai awal sel `NB FacIn\Section\InputCoverageAneka_FacIn.xml`
//   `pyDefaultValue` (Calculate Method 1, % Short Period 0, Day 365, % Loading 0, % Scale 100, Indemnity 0, % Indemnity 100,
//   % Discount 0).
// - `DDL\SetTSI_DT.xml` L139-141: TSI coverage → TSI objek; `DDL\SetCurrency_DT.xml` L144-206: Currency coverage → objek;
//   `DDL\setTSIAneka_preDT.xml`: TSI tiap objek → semua coverage objek itu (dijalankan saat Add, tanpa hitung ulang).
// - Copy Coverage From This Section: `DDL\CopyCoverageFormSection_Act.xml` - daftar coverage objek sumber menggantikan daftar
//   coverage SETIAP objek (semua lokasi, semua okupasi) ber-`.Section` sama (L2370, termasuk sumber), TSI = TSI objek sasaran
//   (L1587 / L2082), lalu `CountPremiCoverageAneka` (Status Coverage → pemicu 'tsi'). Objek Golf tidak ikut (loop OccupationList).
// - Copy Deductible: `DDL\CopyDeductibleToAllLocationAneka_Act.xml` - DeductibleList setiap coverage (semua lokasi / okupasi /
//   objek) yang kode Coverage-nya sama dengan coverage sumber (L2221) diganti deductible sumber. Sumber = SEMUA coverage objek
//   sumber: langkah 2 meloop `…AnekaList(Param.idxAneka).CoverageList` (L505, repeat EMBEDDED L834) dan Page-Copy
//   `CoverageList(<CURRENT>)` (L721) = coverage iterasi berjalan (struktur terverifikasi sesi c3; arti <CURRENT> = [dugaan kuat]).
// - `addDeductibleAneka_ACT`: setiap Add deductible menyetel Currency deductible TERAKHIR setiap coverage = Currency coverage itu
//   (loop seluruh pohon; ditiru apa adanya).

import type { CoverageAneka, Deductible, LokasiAneka, ObjekAneka } from '../api'

/** Posisi objek: lokasi, okupasi (null = objek Golf langsung di lokasi), objek. Indeks 0-based. */
export interface PosisiObjekAneka {
  l: number
  k: number | null
  j: number
}

export function coverageAnekaBaru(o: Pick<ObjekAneka, 'currency' | 'tsi'>, metode: string): CoverageAneka {
  return {
    coverage: '', oldId: '', coverageNote: '', calculateMethodFacIn: metode === '' ? '1' : metode, calculateMethod: '',
    proRatePercent: '', pctShortPeriod: '0', day: '365', currency: o.currency, tsi: o.tsi, loading: '0', lostLimit: '100', rate: '',
    indemnity: '0', unit: '', indemnityPercentage: '100', premium: '', discountPercentage: '0', discount: '', limitOfLiability: '',
    pctLol: '', subLimitNote: '', deductibles: [],
  }
}

/** Coverage dari server dilengkapi medan yang absen (data lama / P1). */
export function rapikanCoverageAneka(c: Partial<CoverageAneka>): CoverageAneka {
  return { ...coverageAnekaBaru({ currency: '', tsi: '' }, '1'), ...c, deductibles: c.deductibles ?? [] }
}

export function objekDi(lokasi: LokasiAneka[], p: PosisiObjekAneka): ObjekAneka | undefined {
  const l = lokasi[p.l]
  if (l === undefined) return undefined
  return p.k === null ? l.objekGolf[p.j] : l.okupasi[p.k]?.objek[p.j]
}

/** Ganti satu objek - salinan baru seluruh jalur. */
export function gantiObjek(lokasi: LokasiAneka[], p: PosisiObjekAneka, f: (o: ObjekAneka) => ObjekAneka): LokasiAneka[] {
  return lokasi.map((l, a) => {
    if (a !== p.l) return l
    if (p.k === null) return { ...l, objekGolf: l.objekGolf.map((o, b) => (b === p.j ? f(o) : o)) }
    return { ...l, okupasi: l.okupasi.map((k, b) => (b !== p.k ? k : { ...k, objek: k.objek.map((o, c) => (c === p.j ? f(o) : o)) })) }
  })
}

/** Setiap objek di pohon (okupasi, lalu Golf) beserta posisinya. */
export function semuaObjek(lokasi: LokasiAneka[]): { p: PosisiObjekAneka; o: ObjekAneka }[] {
  const hasil: { p: PosisiObjekAneka; o: ObjekAneka }[] = []
  lokasi.forEach((l, a) => {
    l.okupasi.forEach((k, b) => k.objek.forEach((o, c) => hasil.push({ p: { l: a, k: b, j: c }, o })))
    l.objekGolf.forEach((o, c) => hasil.push({ p: { l: a, k: null, j: c }, o }))
  })
  return hasil
}

/** Metode coverage pertama objek pertama okupasi pertama lokasi pertama (indeks keras AddCoverageAneka_act L294-295). */
export function metodeAwal(lokasi: LokasiAneka[]): string {
  return lokasi[0]?.okupasi[0]?.objek[0]?.coverages?.[0]?.calculateMethodFacIn ?? ''
}

/** setTSIAneka_preDT: TSI tiap objek didorong ke semua coverage objek itu (tanpa hitung ulang premi). */
export function dorongTsiObjek(lokasi: LokasiAneka[]): LokasiAneka[] {
  const satu = (o: ObjekAneka): ObjekAneka => (o.coverages ? { ...o, coverages: o.coverages.map((c) => ({ ...c, tsi: o.tsi })) } : o)
  return lokasi.map((l) => ({
    ...l,
    okupasi: l.okupasi.map((k) => ({ ...k, objek: k.objek.map(satu) })),
    objekGolf: l.objekGolf.map(satu),
  }))
}

/** Add coverage (AddCoverageAneka_act lalu refresh ber-pre-DT setTSIAneka_preDT). */
export function tambahCoverage(lokasi: LokasiAneka[], p: PosisiObjekAneka): LokasiAneka[] {
  const metode = metodeAwal(lokasi)
  const baru = gantiObjek(lokasi, p, (o) => ({ ...o, coverages: [...(o.coverages ?? []), coverageAnekaBaru(o, metode)] }))
  return dorongTsiObjek(baru)
}

/** Ubah satu coverage; TSI / Currency coverage ikut ke objek (SetTSI_DT / SetCurrency_DT). */
export function ubahCoverageDi(lokasi: LokasiAneka[], p: PosisiObjekAneka, n: number, c: CoverageAneka): LokasiAneka[] {
  return gantiObjek(lokasi, p, (o) => {
    const lama = o.coverages?.[n]
    const coverages = (o.coverages ?? []).map((x, i) => (i === n ? c : x))
    return {
      ...o,
      coverages,
      tsi: lama !== undefined && lama.tsi !== c.tsi ? c.tsi : o.tsi,
      currency: lama !== undefined && lama.currency !== c.currency ? c.currency : o.currency,
    }
  })
}

/**
 * addDeductibleAneka_ACT: Currency deductible TERAKHIR setiap coverage di seluruh pohon = Currency coverage itu. Dipanggil
 * sesudah satu deductible ditambahkan.
 */
export function setelCurrencyDeductibleTerakhir(lokasi: LokasiAneka[]): LokasiAneka[] {
  const cov = (c: CoverageAneka): CoverageAneka => {
    if (c.deductibles.length === 0) return c
    const d = c.deductibles.map((x, i) => (i === c.deductibles.length - 1 ? { ...x, currency: c.currency } : x))
    return { ...c, deductibles: d }
  }
  const obj = (o: ObjekAneka): ObjekAneka => (o.coverages ? { ...o, coverages: o.coverages.map(cov) } : o)
  return lokasi.map((l) => ({
    ...l,
    okupasi: l.okupasi.map((k) => ({ ...k, objek: k.objek.map(obj) })),
    objekGolf: l.objekGolf.map(obj),
  }))
}

const salinDeductibles = (d: Deductible[]): Deductible[] => d.map((x) => ({ ...x }))

/**
 * Copy Coverage From This Section: hasil pohon + posisi sasaran (yang perlu dihitung ulang pemicu 'tsi'). Hanya objek okupasi.
 */
export function salinCoverageSeksi(lokasi: LokasiAneka[], sumber: PosisiObjekAneka): { lokasi: LokasiAneka[]; sasaran: PosisiObjekAneka[] } {
  const s = objekDi(lokasi, sumber)
  if (s === undefined) return { lokasi, sasaran: [] }
  const daftar = s.coverages ?? []
  const sasaran = semuaObjek(lokasi).filter((x) => x.p.k !== null && x.o.section === s.section).map((x) => x.p)
  let hasil = lokasi
  for (const p of sasaran) {
    hasil = gantiObjek(hasil, p, (o) => ({
      ...o,
      coverages: daftar.map((c) => ({ ...c, tsi: o.tsi, deductibles: salinDeductibles(c.deductibles) })),
    }))
  }
  return { lokasi: hasil, sasaran }
}

/** Copy Deductible To All / (This Section): deductible coverage berkode sama diganti deductible sumber; jumlah coverage tersalin. */
export function salinDeductibleAneka(lokasi: LokasiAneka[], sumber: PosisiObjekAneka): { lokasi: LokasiAneka[]; jumlah: number } {
  const s = objekDi(lokasi, sumber)
  const peta = new Map((s?.coverages ?? []).filter((c) => c.coverage !== '').map((c) => [c.coverage, c.deductibles]))
  let jumlah = 0
  const cov = (c: CoverageAneka): CoverageAneka => {
    const d = peta.get(c.coverage)
    if (d === undefined) return c
    jumlah++
    return { ...c, deductibles: salinDeductibles(d) }
  }
  const hasil = lokasi.map((l) => ({
    ...l,
    okupasi: l.okupasi.map((k) => ({ ...k, objek: k.objek.map((o) => (o.coverages ? { ...o, coverages: o.coverages.map(cov) } : o)) })),
  }))
  return { lokasi: hasil, jumlah }
}
