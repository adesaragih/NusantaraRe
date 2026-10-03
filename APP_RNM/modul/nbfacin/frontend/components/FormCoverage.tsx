// Form satu coverage FIRE (tiket 43, tahap C1) - port `NB FacIn\Section\CoverageItem.xml` (flow action
// `CoverageItem`) dan popup Choose Coverage (local action `ChooseCoverage`, RD `BrowseCoverageFacIn_RD`).
//
// Premi TIDAK dihitung di layar: setiap perubahan medan hitung memanggil `POST …/hitung-coverage` (rumus
// `CountPremi_ACT`, desimal eksak di backend) - mode "percent" (premi dari rate) kecuali Gross Premium yang diubah
// (mode "amount", rate dihitung balik). Medan tampil per Coverage Basis mengikuti syarat sel `CoverageItem.xml`.
//
// Keputusan agent (tiket 43): P-1 Layering (basis 5), Zone / 4.2 Construction, Accumulation, Indemnity Unit, View
// Indemnity Table = tahap berikut. P-2 pilihan Days = 365 / 366 (aturan `.Day` tidak ada di korpus; data contoh 365)
// `[dugaan]`. P-3 tanda wajib ‰ Gross Rate tanpa menahan Save (pola L-2). P-4 jeda hitung 500 ms; jawaban lama dibuang.

import { useEffect, useRef, useState } from 'react'

import { Area, Field, Gagal, Kosong, Memuat, Modal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah } from '../../../../inti/frontend/lib/desimal'
import { formatNumber } from '../../../../inti/frontend/lib/format'
import { cariCoverage, daftarMataUang, hitungCoverage, type BarisCoverage, type CoverageObjek, type ModeHitung } from '../api'
import GridDeductible from './GridDeductible'
import IsianUang from './IsianUang'
import { FORM_COV as F, LABEL_COVERAGE_BASIS, OPSI_COVERAGE_BASIS, PILIHAN_PERIODE, POPUP_OKUPASI, TEKS_COVERAGE } from '../labels'

const JEDA_MS = 500
const DESIMAL = 4

/** Coverage kosong: basis Sum Insured, % Indemnity 100 (`AddCoverageAutoFire`). */
export function coverageBaru(b?: BarisCoverage): CoverageObjek {
  return {
    coverage: b?.id ?? '', oldId: b?.oldId ?? '', coverageNote: b?.nama ?? '', coverageBasis: '1', day: '', tsi: '',
    indemnity: '', rate: '', rateOjk: '', firstLoss: '', discountPercentage: '', tsiLiability: '', netRate: '',
    limitOfLiability: '', pctLol: '', proRatePercent: '', indemnityPercentage: '100', firstScale: '', sublimit: '',
    lostLimit: '', emlPml: '', discount: '', premium: '', conditions: '',
  }
}

/** Medan angka yang diisi pengguna. */
const MEDAN_ANGKA = [
  'indemnity', 'rate', 'firstLoss', 'discountPercentage', 'netRate', 'limitOfLiability', 'pctLol', 'indemnityPercentage',
  'firstScale', 'sublimit', 'lostLimit', 'emlPml', 'discount', 'premium',
] as const

/** Angka sah: kosong atau desimal bertitik. */
export const angkaSah = (v: string) => v.trim() === '' || desimalSah(v.trim())

/** Ada coverage dengan isian angka tidak sah. */
export function adaGalatCoverage(c: CoverageObjek[]): boolean {
  return c.some((x) => MEDAN_ANGKA.some((k) => !angkaSah(x[k])))
}

function Tampil({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="field">
      <span className="field__label">{label}</span>
      <div className="nbf-inward__teks">{nilai}</div>
    </div>
  )
}

function PopupCoverage({ onTutup, onPilih }: { onTutup: () => void; onPilih: (b: BarisCoverage) => void }) {
  const [kotak, setKotak] = useState('')
  const [hasil, setHasil] = useState<BarisCoverage[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const nomor = useRef(0)

  useEffect(() => {
    const n = ++nomor.current
    const jadwal = window.setTimeout(
      () => {
        cariCoverage(kotak).then(
          (h) => {
            if (n === nomor.current) setHasil(h.baris)
          },
          (err: unknown) => {
            if (n === nomor.current) setGalat(err)
          },
        )
      },
      kotak === '' ? 0 : 400,
    )
    return () => window.clearTimeout(jadwal)
  }, [kotak])

  return (
    <Modal judul={F.pilihCoverage.label} onTutup={onTutup} lebar>
      <div className="nbf-popup__kotak">
        <Field label={TEKS_COVERAGE.cariCoverage} value={kotak} onChange={setKotak} />
      </div>
      <div className="table-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col">{F.coverage.label}</th>
              <th scope="col">{POPUP_OKUPASI.kolom[1].label}</th>
              <th scope="col" />
            </tr>
          </thead>
          {hasil && hasil.length > 0 && (
            <tbody>
              {hasil.map((b) => (
                <tr key={b.id}>
                  <td>{b.oldId}</td>
                  <td>{b.nama}</td>
                  <td className="table__actions">
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => onPilih(b)}>
                      {POPUP_OKUPASI.pilih.label}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          )}
        </table>
      </div>
      {hasil === null && galat === null && <Memuat />}
      {hasil !== null && hasil.length === 0 && <Kosong pesan={TEKS_COVERAGE.tanpaCoverage} />}
      <Gagal galat={galat} />
    </Modal>
  )
}

export default function FormCoverage({
  caseId,
  c,
  tsiItem,
  item,
  ubah,
}: {
  caseId: string
  c: CoverageObjek
  /** TSI Object Item (teks desimal). */
  tsiItem: string
  /** Item pemilik (PctAdjustment `CountPremi_ACT`; mata uang = awal deductible baru). */
  item?: { isAdjustable: boolean; pctAdjustOther: string; currency?: string }
  ubah: (c: CoverageObjek) => void
}) {
  const [pilih, setPilih] = useState(false)
  const [menghitung, setMenghitung] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const nomor = useRef(0)
  const jadwal = useRef<number | undefined>(undefined)
  // `Param.DiscountStatus`: medan diskon terakhir yang diubah pengguna.
  const modeDiskon = useRef<ModeHitung>('percent')

  useEffect(() => () => window.clearTimeout(jadwal.current), [])

  // Daftar mata uang untuk Currency deductible (RD `BrowseCurrency_RD`).
  const [mataUang, setMataUang] = useState<string[]>([])
  useEffect(() => {
    let batal = false
    daftarMataUang().then(
      (h) => {
        if (!batal) setMataUang(h.baris)
      },
      () => {},
    )
    return () => {
      batal = true
    }
  }, [])

  /** Ubah lalu jadwalkan hitung ulang di backend. */
  function ubahHitung(baru: CoverageObjek, mode: ModeHitung) {
    ubah(baru)
    window.clearTimeout(jadwal.current)
    if (MEDAN_ANGKA.some((k) => !angkaSah(baru[k])) || baru.coverageBasis === '5') return
    jadwal.current = window.setTimeout(() => {
      const n = ++nomor.current
      setMenghitung(true)
      setGalat(null)
      hitungCoverage(caseId, {
        coverage: baru,
        tsi: tsiItem,
        mode,
        modeDiskon: modeDiskon.current,
        isAdjustable: item?.isAdjustable,
        pctAdjustOther: item?.pctAdjustOther,
      }).then(
        (h) => {
          if (n === nomor.current) {
            ubah(h)
            setMenghitung(false)
          }
        },
        (err: unknown) => {
          if (n === nomor.current) {
            setGalat(err)
            setMenghitung(false)
          }
        },
      )
    }, JEDA_MS)
  }

  const set = (k: keyof CoverageObjek, mode: ModeHitung = 'percent') => (v: string) => ubahHitung({ ...c, [k]: v }, mode)
  const angka = (k: (typeof MEDAN_ANGKA)[number], label: string, mode: ModeHitung = 'percent', required = false) => (
    <Field label={label} value={c[k]} onChange={set(k, mode)} required={required} error={angkaSah(c[k]) ? undefined : TEKS_COVERAGE.angka} />
  )
  /** Medan UANG: berformat ribuan saat diketik (IsianUang, permintaan work owner 03-10-2026). */
  const uang = (k: 'limitOfLiability' | 'premium', label: string, mode: ModeHitung = 'percent') => (
    <IsianUang label={label} value={c[k]} onChange={set(k, mode)} error={angkaSah(c[k]) ? undefined : TEKS_COVERAGE.angka} />
  )
  /** Medan diskon: menandai mode diskon (`DiscountStatus`) sebelum hitung. */
  const angkaDiskon = (k: 'discountPercentage' | 'discount', label: string, md: ModeHitung) => {
    const props = {
      label,
      value: c[k],
      onChange: (v: string) => {
        modeDiskon.current = md
        ubahHitung({ ...c, [k]: v }, 'percent')
      },
      error: angkaSah(c[k]) ? undefined : TEKS_COVERAGE.angka,
    }
    // Discount = uang (berformat ribuan); % Discount = persen (isian biasa).
    return k === 'discount' ? <IsianUang {...props} /> : <Field {...props} />
  }
  const b = c.coverageBasis

  return (
    <div className="nbf-objek__isi">
      <div className="nbf-opp__kolom">
        <div className="nbf-opp__tumpuk">
          <Pilih label={LABEL_COVERAGE_BASIS} value={b} onChange={set('coverageBasis')} opsi={OPSI_COVERAGE_BASIS} />
          {b === '5' && <div className="alert alert--warn">{TEKS_COVERAGE.layeringBelum}</div>}
          <div className="nbf-opp__tombol">
            <button type="button" className="btn btn--sm" onClick={() => setPilih(true)}>
              {F.pilihCoverage.label}
            </button>
          </div>
          <Tampil label={F.coverage.label} nilai={[c.oldId, c.coverageNote].filter(Boolean).join(' - ')} />
          <Area label={F.conditions.label} value={c.conditions} onChange={(v) => ubah({ ...c, conditions: v })} baris={2} />
        </div>
        <div className="nbf-opp__tumpuk">
          <Pilih
            label={F.day.label}
            value={c.day}
            onChange={set('day')}
            opsi={PILIHAN_PERIODE.day.map((d) => ({ value: d, label: d }))}
          />
          <Tampil label={F.tsi.label} nilai={formatNumber(c.tsi || tsiItem, DESIMAL)} />
          {b === '2' && angka('indemnity', F.indemnity.label)}
          {angka('rate', F.rate.label, 'percent', true)}
          {b === '2' && angka('firstLoss', F.firstLoss.label)}
          {angkaDiskon('discountPercentage', F.discountPercentage.label, 'percent')}
          {b === '2' && <Tampil label={F.tsiLiability.label} nilai={formatNumber(c.tsiLiability, DESIMAL)} />}
          {angka('netRate', F.netRate.label)}
          {uang('limitOfLiability', F.limitOfLiability.label)}
          {angka('pctLol', F.pctLol.label)}
          <Tampil label={F.proRate.label} nilai={formatNumber(c.proRatePercent, DESIMAL)} />
          {angka('indemnityPercentage', F.indemnityPercentage.label)}
          {b === '2' && angka('firstScale', F.firstScale.label)}
          {(b === '4' || b === '5') && angka('sublimit', F.sublimit.label)}
          {angka('lostLimit', F.lostLimit.label)}
          {b === '3' && angka('emlPml', F.emlPml.label)}
          {angkaDiskon('discount', F.discount.label, 'amount')}
          {uang('premium', F.premium.label, 'amount')}
          {menghitung && <span className="muted">{TEKS_COVERAGE.menghitung}</span>}
          <Gagal galat={galat} />
        </div>
      </div>
      {/* Deductible (tahap C3, tiket 45): tidak memicu hitung premi. */}
      <GridDeductible
        daftar={c.deductibles ?? []}
        currency={item?.currency ?? ''}
        mataUang={mataUang}
        ubah={(deductibles) => ubah({ ...c, deductibles })}
      />
      {pilih && (
        <PopupCoverage
          onTutup={() => setPilih(false)}
          onPilih={(row) => {
            setPilih(false)
            ubahHitung({ ...c, coverage: row.id, oldId: row.oldId, coverageNote: row.nama }, 'percent')
          }}
        />
      )}
    </div>
  )
}
