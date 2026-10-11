// Tampilan polis NB Fac In HANYA-BACA - tab Scoring Risk FIRE (perintah work owner 11-10-2026 "tampilan view polis yang
// belum tergenerate tolong diperbaiki. cek tampilan nya di nb facin!"). Rupa = `modul/nbfacin/frontend/components/
// TabScoring.tsx` dalam keadaan baca-saja (pola, bukan impor); label VERBATIM `modul/nbfacin/frontend/labels.ts`.
//
// Keputusan work owner 11-10-2026: NILAI TERSIMPAN - pilihan, ScorePerFactor, Score kolom, FinalScore / NoteFinalScore,
// Status, dan ScoringResult dibaca apa adanya dari tabel rancangan nbfacin 205 (T_SCORINGRISK -> T_SCORINGRESULT /
// T_DATASCORINGRISKLIST -> T_SCORING_FACTOR -> T_SCORING_OPTION), bukan dihitung ulang (`ValueScoringRisk_act` /
// `SetScore_Act` nbfacin). Dua bentuk halaman: tabel datar (`.ScoringFactor(m)` ber-FactorGroup / FactorName, opsi
// `.ScoringOption(k)` ber-CheckBox / OptionNo - aturan loader V-30) dan halaman Pega (`.<Grup>.<Nama>.ChechBoxN`).

import type { Baris, Halaman } from './api'
import { AKAR_NB, NOL, anakNB, daftarNB, eksak, nilaiNB, tambahE, tanggalNB, teksE, umumNB } from './polisNB'

// ---------------------------------------------------------------- label (VERBATIM nbfacin/labels.ts)

/** `SCORING`. */
export const SCORING_NB = {
  judul: 'Risk Scoring - Property',
  insured: 'The Insured',
  topRisk: 'Top Risk',
  dominantRisk: 'Dominant Risk',
  periode: 'Period',
  lokasi1: 'Risk Location 1 (Top Risk)',
  lokasi2: 'Risk Location 2 (Dominant Risk)',
  status: 'Status',
  skor: 'Score',
  kolom: ['No', 'Scoring Factor', 'Percentage', 'Risk Detail', 'Top Risk', 'Dominant Risk'],
  total: 'Total',
  hasil: 'Scoring Result',
  klasifikasi: 'Risk Classifications',
  kelas: [
    { rentang: '86 - 100', nama: 'Good' },
    { rentang: '76 - 85', nama: 'Above Average' },
    { rentang: '60 - 75', nama: 'Average' },
    { rentang: '30 - 59', nama: 'Below Average' },
    { rentang: '0 - 29', nama: 'Poor' },
  ],
  rekomendasi: 'Recomendation',
  persetujuan: 'Required approval :',
  dari: 'From :',
  karena: 'Due to :',
  lossRatioInput: 'Loss Ratio 1 up to 5 Years',
} as const

/** Satu sub-faktor: kunci kontrak, label Risk Detail, bobot, pilihan ChechBox1..n; `terkunci` = nomor pilihan isian sistem. */
export interface SubFaktorScoringNB {
  kunci: string
  label: string
  bobot: string
  opsi: readonly string[]
  terkunci?: readonly number[]
}

export interface FaktorScoringNB {
  no: string
  faktor: string
  bobot: string
  sub: readonly SubFaktorScoringNB[]
}

/** `FAKTOR_SCORING` (urut layar; `kunci` = kontrak `DaftarFaktorScoring` nbfacin). */
export const FAKTOR_SCORING_NB: readonly FaktorScoringNB[] = [
  {
    no: '1',
    faktor: 'Occupation',
    bobot: '15,00',
    sub: [
      {
        kunci: 'occupation',
        label: '',
        bobot: '',
        opsi: ['Low Risks', 'Medium Risks', 'High Risks', 'Non-Preferred Risk', 'Treaty Exclusion'],
        terkunci: [1, 2, 3, 4, 5],
      },
    ],
  },
  {
    no: '2',
    faktor: 'Risk Location',
    bobot: '10,00',
    sub: [
      {
        kunci: 'riot',
        label: 'Riot',
        bobot: '1',
        opsi: ['Cover exclude Riot', 'Non - Riot prone area', 'Riot prone area', 'No Information'],
      },
      {
        kunci: 'flood',
        label: 'Flood',
        bobot: '5',
        opsi: ['Cover exclude Flood', 'Non - Flood prone area', 'Flood prone area', 'No Information'],
      },
      {
        kunci: 'earthquake',
        label: 'Earthquake',
        bobot: '2,5',
        opsi: [
          'Cover exclude Earthquake',
          'Zone 1',
          'Zone 2',
          'Zone 3',
          'Zone 4',
          'Zone 5',
          'Other than zona 1 up to zona 5',
        ],
      },
      {
        kunci: 'tsunami',
        label: 'Tsunami',
        bobot: '1,5',
        opsi: ['Cover exclude Tsunami', 'Non - Tsunami prone area', 'Tsunami prone area', 'No Information'],
      },
    ],
  },
  {
    no: '3',
    faktor: 'Object Conditions',
    bobot: '12,50',
    sub: [
      {
        kunci: 'classConstruction',
        label: 'Class Construction',
        bobot: '2,5',
        opsi: ['I', 'II', 'III', 'Growing Trees', 'Open storage', 'No Information'],
      },
      {
        kunci: 'building',
        label: 'Building',
        bobot: '2,5',
        opsi: [
          'Good condition / housekeeping',
          'Growing Trees / Open storage',
          'Bad Condition /  housekeeping',
          'Silent Risk',
          'No Information',
        ],
      },
      {
        kunci: 'machinery',
        label: 'Machinery',
        bobot: '2,5',
        opsi: [
          'Good condition / maintenance',
          'Growing Trees / Open storage',
          'Bad condition / maintenance',
          'Silent Risk',
          'No Information',
        ],
      },
      {
        kunci: 'stock',
        label: 'Stock in the policy cover',
        bobot: '2,5',
        opsi: ['Policy does not cover stock', 'Sum Insured < 60%', 'Sum Insured Stock >= 60%', 'No Information'],
      },
      {
        kunci: 'typeOfStock',
        label: 'Type of Stock',
        bobot: '2,5',
        opsi: [
          'No Stock in the premises',
          'Non Combustible or Non Perishable Goods',
          'Combustible or Perishable Goods or obsolete',
          'No Information',
        ],
      },
    ],
  },
  {
    no: '4',
    faktor: 'FEA (Fire Extinguishing Appliance)',
    bobot: '10,00',
    sub: [
      {
        kunci: 'portableFireExtinguisher',
        label: 'Portable Fire Extinguisher',
        bobot: '3',
        opsi: [
          'Available',
          'Available and  adequate',
          'Not Available (Growing Trees)',
          'Not available or No Information',
        ],
      },
      {
        kunci: 'hydrant',
        label: 'Hydrant System',
        bobot: '1,5',
        opsi: ['Available', 'Not Available (Growing Trees)', 'Not available or No Information'],
      },
      {
        kunci: 'sprinkler',
        label: 'Sprinkler System',
        bobot: '1',
        opsi: ['Available', 'Not Available (Growing Trees, Open Storage)', 'Not available or No Information'],
      },
      {
        kunci: 'fireAlarm',
        label: 'Fire Alarm System',
        bobot: '1,5',
        opsi: ['Available', 'Not available or No Information'],
      },
      {
        kunci: 'securityGuard',
        label: 'Security Guard',
        bobot: '1,5',
        opsi: ['Available in the premises', 'Available on outside premises', 'Not available or No Information'],
      },
      {
        kunci: 'fireBrigade',
        label: 'Fire brigade',
        bobot: '1,5',
        opsi: [
          'Available in the premises',
          'Available on outside within < 30 minutes',
          'Available on outsidewithin 30-60 minutes',
          'Available on outsidewithin  > 60 minutes',
          'Not available or No Information',
        ],
      },
    ],
  },
  {
    no: '5',
    faktor: 'Net Loss Ratio',
    bobot: '22,50',
    sub: [
      {
        kunci: 'lossRatio',
        label: 'Loss Ratio (Claim/net premium)',
        bobot: '',
        opsi: [
          'Nil 5 years and more',
          'Nil 3 years',
          'Nil 1 years',
          '-- Loss Ratio > 0 up to 10%',
          '-- Loss Ratio > 10 up to 25%',
          '-- Loss Ratio > 25 up to 50%',
          '-- Loss Ratio > 50 up to 75%',
          '-- Loss Ratio> 75 up to 100%',
          '-- Loss Ratio > 100%',
          'No Information',
        ],
        terkunci: [4, 5, 6, 7, 8, 9],
      },
    ],
  },
  {
    no: '6',
    faktor: 'Risk Improvement',
    bobot: '10,00',
    sub: [
      {
        kunci: 'riskImprovement',
        label: '',
        bobot: '',
        opsi: [
          'Risk Improvement has been done or on progress',
          'Risk Improvement is not required due to nil loss record',
          'No risk improvement for incurred loss or No Information',
        ],
      },
    ],
  },
  {
    no: '7',
    faktor: 'Survey Report',
    bobot: '10,00',
    sub: [
      {
        kunci: 'surveyReport',
        label: '',
        bobot: '',
        opsi: [
          'Available - within 1 year',
          'Available - 2 years ago',
          'Available - 3 years ago',
          'Available - more than 3 years ago',
          'No Survey Report, but provided photographs or summary of the insured activities',
          'Not Available or No Information',
        ],
      },
    ],
  },
  {
    no: '8',
    faktor: 'Others',
    bobot: '10,00',
    sub: [
      {
        kunci: 'rate',
        label: 'Rate',
        bobot: '2',
        opsi: ['On Gross Rate (OGR)  standard', 'Market Rate', 'No Information'],
      },
      {
        kunci: 'deductible',
        label: 'Deductible',
        bobot: '2,5',
        opsi: ['Standard or higher than standard', 'Non Standard', 'No Information'],
      },
      {
        kunci: 'totalSumInsured',
        label: 'Total Sum Insured',
        bobot: '1',
        opsi: ['Less than IDR 30,000,000,000', 'IDR 30,000,000,000 or more'],
        terkunci: [1, 2],
      },
      {
        kunci: 'otherInformation',
        label: 'Other Information',
        bobot: '2,5',
        opsi: [
          'Available and not potentially to increase risk exposure',
          'Available and  potentially to increase risk exposure',
          'Not Available or No Information',
        ],
      },
      {
        kunci: 'clauses',
        label: 'Clauses',
        bobot: '2',
        opsi: ['Non Standard (un usual) Clauses', 'Standard (usual) Clauses', 'Not Available or No Information'],
      },
    ],
  },
]

/** `HASIL_SCORING`. */
export const HASIL_SCORING_NB = {
  rekomendasi: [
    { kunci: 'acceptableNotApproval', label: 'Acceptable, not required approval' },
    { kunci: 'acceptableWithApproval', label: 'Acceptable, required approval' },
    { kunci: 'notAcceptable', label: 'Not Acceptable' },
  ],
  dari: [
    { kunci: 'opDivHead', label: 'Operational Div. Head' },
    { kunci: 'opDirector', label: 'Operational Director' },
    { kunci: 'techDirector', label: 'Technical Director' },
    { kunci: 'presDirector', label: 'President Director' },
  ],
  karena: [
    { kunci: 'dueOccupation', label: 'Occupation' },
    { kunci: 'dueObjectConditions', label: 'Object Conditions' },
    { kunci: 'dueLossRatio', label: 'Loss Ratio' },
    { kunci: 'dueOthers', label: 'Others' },
  ],
} as const

/** `TEKS_SCORING` (yang dipakai tampilan hanya-baca). */
export const TEKS_SCORING_NB = {
  kosong: '—',
  diisiSistem: 'Diisi sistem',
} as const

// ---------------------------------------------------------------- nilai tersimpan

export const SKOR_NB = AKAR_NB + '.ScoringRisk'
const DATA_SKOR_NB = SKOR_NB + '.DataScoringRiskList'

/** Kunci kontrak -> [grup, nama] halaman faktor Pega (`DaftarFaktorScoring` nbfacin; grup '' = langsung di bawah kolom). */
export const HALAMAN_FAKTOR_NB: Record<string, readonly [string, string]> = {
  occupation: ['', 'Occupation'],
  riot: ['RiskLocation', 'Riot'],
  flood: ['RiskLocation', 'Flood'],
  earthquake: ['RiskLocation', 'Earthquake'],
  tsunami: ['RiskLocation', 'Tsunami'],
  classConstruction: ['ObjectConditions', 'ClassConstruction'],
  building: ['ObjectConditions', 'Building'],
  machinery: ['ObjectConditions', 'Machinery'],
  stock: ['ObjectConditions', 'Stockinthepolicycover'],
  typeOfStock: ['ObjectConditions', 'TypeofStock'],
  portableFireExtinguisher: ['FEA', 'PortableFireExtinguisher'],
  hydrant: ['FEA', 'HydrantSystem'],
  sprinkler: ['FEA', 'SprinklerSystem'],
  fireAlarm: ['FEA', 'FireAlarmSystem'],
  securityGuard: ['FEA', 'SecurityGuard'],
  fireBrigade: ['FEA', 'Firebrigade'],
  lossRatio: ['LossRatio', 'LossRatio'],
  riskImprovement: ['', 'RiskImprovement'],
  surveyReport: ['', 'SurveyReport'],
  rate: ['Others', 'Rate'],
  deductible: ['Others', 'Deductible'],
  totalSumInsured: ['Others', 'TotalSumInsured'],
  otherInformation: ['Others', 'OtherInformation'],
  clauses: ['Others', 'Clauses'],
}

/** Kunci `HASIL_SCORING_NB` -> medan `.ScoringRisk.ScoringResult.*` (nbfacin handlers/scoring.go `keScoringKabel`). */
const MEDAN_HASIL_NB = {
  acceptableNotApproval: 'AcceptableNotApproval',
  acceptableWithApproval: 'AcceptableWithApproval',
  notAcceptable: 'NotAcceptable',
  opDivHead: 'OperationalDivHead',
  opDirector: 'OperationalDirector',
  techDirector: 'TechnicalDirector',
  presDirector: 'PresidentDirector',
  dueOccupation: 'Occupation',
  dueObjectConditions: 'ObjectConditions',
  dueLossRatio: 'LossRatio',
  dueOthers: 'Others',
} as const

export type KunciHasilSkorNB = keyof typeof MEDAN_HASIL_NB

export interface KolomSkorNB {
  occupationNote: string
  topRiskLocation: string
  /** ObjectNo lokasi Dominant Risk ('' kolom Top Risk). */
  indexLocation: string
  /** Kunci sub-faktor -> nomor ChechBox tercentang (1..n); absen = tidak dicentang. */
  pilihan: Record<string, number>
  /** Kunci sub-faktor -> ScorePerFactor tersimpan. */
  skor: Record<string, string>
  persenLossRatio: string
  /** `.Score` kolom. */
  total: string
}

export interface TampilanSkorNB {
  insuredName: string
  status: string
  mulai: string
  selesai: string
  /** [Top Risk, Dominant Risk] - selalu dua. */
  kolom: KolomSkorNB[]
  finalScore: string
  noteFinalScore: string
  hasil: Record<KunciHasilSkorNB, boolean>
}

const benar = (v: string | undefined) => (v ?? '').trim() === 'true'

/** Opsi tercentang ganda -> nomor terkecil (nbfacin `BacaScoring`); 0 = tidak ada. */
function nomorTerkecil(nomor: number[]): number {
  const sah = nomor.filter((n) => Number.isInteger(n) && n > 0)
  return sah.length === 0 ? 0 : Math.min(...sah)
}

/** Satu kolom dari baris `DataScoringRiskList(n)`: faktor tabel datar dulu, bila tidak ada bentuk halaman Pega. */
function kolomSkor(h: Halaman, b: Baris | undefined, n: number): KolomSkorNB {
  const k: KolomSkorNB = {
    occupationNote: b?.['OccupationNote'] ?? '',
    topRiskLocation: b?.['TopRiskLocation'] ?? '',
    indexLocation: b?.['IndexLocation'] ?? '',
    pilihan: {},
    skor: {},
    persenLossRatio: '',
    total: b?.['Score'] ?? '',
  }
  if (!b) return k
  const jalurFaktor = anakNB(DATA_SKOR_NB, n, 'ScoringFactor')
  const faktor = daftarNB(h, jalurFaktor)
  for (const [kunci, [grup, nama]] of Object.entries(HALAMAN_FAKTOR_NB)) {
    let pilih = 0
    let skor = ''
    let persen = ''
    const m = faktor.findIndex((f) => (f['FactorGroup'] ?? '') === grup && (f['FactorName'] ?? '') === nama)
    if (m >= 0) {
      const f = faktor[m]!
      const opsi = daftarNB(h, anakNB(jalurFaktor, m + 1, 'ScoringOption'))
      pilih = nomorTerkecil(opsi.filter((o) => benar(o['CheckBox'])).map((o) => Number(o['OptionNo'])))
      skor = f['ScorePerFactor'] ?? ''
      if (skor === '' && pilih > 0) skor = opsi.find((o) => Number(o['OptionNo']) === pilih)?.['Score'] ?? ''
      persen = f['PercentLossRatio'] ?? ''
    } else {
      const awal = (grup === '' ? '' : grup + '.') + nama + '.'
      const nomor: number[] = []
      for (const [kb, v] of Object.entries(b)) {
        if (kb.startsWith(awal + 'ChechBox') && benar(v)) nomor.push(Number(kb.slice(awal.length + 'ChechBox'.length)))
      }
      pilih = nomorTerkecil(nomor)
      skor = b[awal + 'ScorePerFactor'] ?? ''
      if (skor === '' && pilih > 0) skor = b[awal + 'Score' + pilih] ?? ''
      persen = b[awal + 'PercentLossRatio'] ?? ''
    }
    if (pilih > 0) k.pilihan[kunci] = pilih
    k.skor[kunci] = skor
    if (kunci === 'lossRatio') k.persenLossRatio = persen
  }
  return k
}

/** Tab Scoring Risk dari halaman polis (nilai tersimpan). */
export function skorPolisNB(h: Halaman): TampilanSkorNB {
  const data = daftarNB(h, DATA_SKOR_NB)
  const hasil = {} as Record<KunciHasilSkorNB, boolean>
  for (const [k, medan] of Object.entries(MEDAN_HASIL_NB)) {
    hasil[k as KunciHasilSkorNB] = benar(nilaiNB(h, 'ScoringRisk.ScoringResult.' + medan))
  }
  return {
    insuredName: umumNB(h).insured,
    status: nilaiNB(h, 'ScoringRisk.Status'),
    mulai: tanggalNB(nilaiNB(h, 'PolicyData.StartDateTime')),
    selesai: tanggalNB(nilaiNB(h, 'PolicyData.EndDateTime')),
    kolom: [kolomSkor(h, data[0], 1), kolomSkor(h, data[1], 2)],
    finalScore: nilaiNB(h, 'ScoringRisk.FinalScore'),
    noteFinalScore: nilaiNB(h, 'ScoringRisk.NoteFinalScore'),
    hasil,
  }
}

const DESIMAL_SAH = /^[+-]?(\d+(\.\d*)?|\.\d+)$/

/** Subtotal skor satu faktor pada satu kolom (`subtotalFaktor` nbfacin): '' bila belum ada sub dicentang atau skor tak
 * terbaca; skor kosong dilewati; jumlah eksak. */
export function subtotalFaktorNB(kolom: Pick<KolomSkorNB, 'pilihan' | 'skor'> | undefined, f: FaktorScoringNB): string {
  if (!kolom || !f.sub.some((s) => (kolom.pilihan[s.kunci] ?? 0) > 0)) return ''
  let total = NOL
  for (const s of f.sub) {
    const t = (kolom.skor[s.kunci] ?? '').trim()
    if (t === '') continue
    if (!DESIMAL_SAH.test(t)) return ''
    total = tambahE(total, eksak(t))
  }
  return teksE(total)
}

/** Nada warna kelas skor (`nadaKelas` nbfacin): baik / sedang / buruk / '' netral. */
export function nadaKelasNB(nama: string): '' | 'baik' | 'sedang' | 'buruk' {
  if (nama === 'Good' || nama === 'Above Average') return 'baik'
  if (nama === 'Average') return 'sedang'
  if (nama === 'Below Average' || nama === 'Poor') return 'buruk'
  return ''
}

/** Nada rekomendasi (`NADA_REKOMENDASI` nbfacin). */
export const NADA_REKOMENDASI_NB: Partial<Record<KunciHasilSkorNB, 'baik' | 'sedang' | 'buruk'>> = {
  acceptableNotApproval: 'baik',
  acceptableWithApproval: 'sedang',
  notAcceptable: 'buruk',
}
