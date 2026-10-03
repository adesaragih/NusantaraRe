// Subsection `DetailPoliciesNonProportional` -> `DetailPolicyTreatyInNonProportional`
// (polis NB NonProporsional / XOL, `[keputusan work owner]` K8): grid master
// Limits / Share / Total Spreaded / Share Facultative (baca-saja), section
// `SpreadingRiskList`, dan jadwal angsuran per mata uang beserta rinciannya
// (`InstallmentList`, hanya-baca). Label, syarat tampil, dan aturan sunting: `../nonprop.ts`.
//
// ⛔ Nol perhitungan di sini - semua nilai datang dari backend (pilih bisnis:
// InputPolicyTreatyInDetail_NonProp; refresh `CountSpreading`).

import { Panel } from '../../../../inti/frontend/components/ui/dasar'
import { daftar, nilai, type Baris, type Halaman, type Pilihan } from '../api'
import { BAGIAN, KOLOM_ANGSURAN, KOLOM_SPREADING, TOMBOL } from '../labels'
import {
  JUDUL_NONPROP,
  KOLOM_FAKULTATIF,
  KOLOM_LIMIT,
  KOLOM_RINCI,
  KOLOM_SHARE,
  TOTAL_FAKULTATIF,
  TOTAL_LIMIT,
  TOTAL_SHARE,
  TOTAL_SPREADED,
  daftarMaster,
  jalurRnmShare,
  spreadingTerbuka,
  tampilFakultatif,
  type Kolom,
  type Total,
} from '../nonprop'
import { sajikan, type Sajian } from '../sajian'

const P = 'PolicyTreatyIn.'
const SPREADING = P + 'SpreadingRiskList'
const ANGSURAN = P + 'ListInstallment'

/** pxNumber section NonProp dan SpreadingRiskList ber-`pyDecimalPlaces` 2 (K14). */
const DUA: Sajian = { desimal: 2 }
/** Section `InstallmentList`: `pyFormatType number` tanpa `pyDecimalPlaces` - pola inti. */
const POLA_INTI: Sajian = {}
const angka2 = (v: string | undefined) => sajikan(v ?? '', DUA)
const TEKS = new Set(['Note', 'Currency'])

function Grid({ baris, kolom, sajian = DUA }: { baris: Baris[]; kolom: Kolom[]; sajian?: Sajian }) {
  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>
            {kolom.map((c) => (
              <th key={c.m} scope="col">
                {c.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.map((b, i) => (
            <tr key={i}>
              {kolom.map((c) => (
                <td key={c.m}>
                  {TEKS.has(c.m) ? (b[c.m] ?? '') : sajikan(b[c.m] ?? '', c.m === 'DueDate' ? 'tanggal' : sajian)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function Totals({ h, daftarTotal }: { h: Halaman; daftarTotal: Total[] }) {
  return (
    <div className="nbti__total-np">
      {daftarTotal.map((t) => (
        <Grid
          key={t.daftar}
          baris={daftarMaster(h, t.daftar)}
          kolom={[{ m: 'Currency', label: t.label }, { m: 'Value', label: JUDUL_NONPROP.value }, ...(t.tambahan ?? [])]}
        />
      ))}
    </div>
  )
}

export interface PropsDetailNonProp {
  halaman: Halaman
  /** Layar admin dan pelaku boleh bekerja - grid spreading dapat disunting. */
  sunting: boolean
  opsiSpreading: Pilihan[]
  onUbahBaris: (jalur: string, i: number, kunci: string, v: string) => void
  onSetelDaftar: (jalur: string, b: Baris[]) => void
  onRefresh: (aksi: string, indeks: number) => void
}

export default function DetailNonProp({ halaman: h, sunting, opsiSpreading, onUbahBaris, onSetelDaftar, onRefresh }: PropsDetailNonProp) {
  const spreading = daftar(h, SPREADING)
  const terbuka = sunting && spreadingTerbuka(h)
  return (
    <>
      <Panel judul={JUDUL_NONPROP.limits}>
        <Grid baris={daftarMaster(h, 'LimitSummaryList')} kolom={KOLOM_LIMIT} />
        <Totals h={h} daftarTotal={TOTAL_LIMIT} />
      </Panel>

      <Panel judul={JUDUL_NONPROP.share}>
        <div className="field nbti__medan">
          <span className="field__label">{JUDUL_NONPROP.rnmShare}</span>
          <span className="nbti__nilai">{nilai(h, jalurRnmShare(h))}</span>
        </div>
        <Grid baris={daftarMaster(h, 'LimitShareSummaryList')} kolom={KOLOM_SHARE} />
        <Totals h={h} daftarTotal={TOTAL_SHARE} />
        <h4>{JUDUL_NONPROP.totalSpreaded}</h4>
        <Totals h={h} daftarTotal={TOTAL_SPREADED} />
      </Panel>

      {tampilFakultatif(h) && (
        <Panel judul={JUDUL_NONPROP.fakultatif}>
          <div className="field nbti__medan">
            <span className="field__label">{JUDUL_NONPROP.fakultatifPersen}</span>
            <span className="nbti__nilai">{nilai(h, 'TreatyIn.FacultativeShare')}</span>
          </div>
          <Grid baris={daftarMaster(h, 'LimitFacShareSummaryList')} kolom={KOLOM_FAKULTATIF} />
          <Totals h={h} daftarTotal={TOTAL_FAKULTATIF} />
        </Panel>
      )}

      <Panel judul={BAGIAN.spreading}>
        {terbuka && (
          <button type="button" className="btn btn--sm" onClick={() => onSetelDaftar(SPREADING, [...spreading, {} as Baris])}>
            {TOMBOL.add}
          </button>
        )}
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_SPREADING.treatyType}</th>
                <th scope="col">{KOLOM_SPREADING.share}</th>
                <th scope="col">{KOLOM_SPREADING.premium}</th>
                <th scope="col">{KOLOM_SPREADING.claimPct}</th>
                <th scope="col">{KOLOM_SPREADING.claim}</th>
                {terbuka && <th scope="col" />}
              </tr>
            </thead>
            <tbody>
              {spreading.map((b, i) => (
                <tr key={i}>
                  <td>
                    {terbuka ? (
                      <select className="field__input" value={b.TreatyType ?? ''} onChange={(e) => onUbahBaris(SPREADING, i, 'TreatyType', e.target.value)}>
                        <option value="" />
                        {opsiSpreading.map((o) => (
                          <option key={o.nilai} value={o.nilai}>
                            {o.label}
                          </option>
                        ))}
                      </select>
                    ) : (
                      (b.TreatyName ?? b.TreatyType ?? '')
                    )}
                  </td>
                  {(['SharePercentage', 'PremiumSpreaded', 'ClaimPercentage', 'ClaimSpreaded'] as const).map((k) => (
                    <td key={k}>
                      {terbuka && (k === 'SharePercentage' || k === 'ClaimPercentage') ? (
                        <input
                          className="field__input"
                          value={b[k] ?? ''}
                          onChange={(e) => onUbahBaris(SPREADING, i, k, e.target.value)}
                          onBlur={() => onRefresh('CountSpreading', i + 1)}
                        />
                      ) : (
                        angka2(b[k])
                      )}
                    </td>
                  ))}
                  {terbuka && (
                    <td>
                      <button
                        type="button"
                        className="btn btn--sm btn--danger"
                        onClick={() => onSetelDaftar(SPREADING, spreading.filter((_, j) => j !== i))}
                      >
                        {TOMBOL.delete}
                      </button>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
            <tfoot>
              <tr>
                <td>{KOLOM_SPREADING.totalShare}</td>
                <td>{angka2(nilai(h, P + 'TotalSharePercentagePremium'))}</td>
                <td>{angka2(nilai(h, P + 'TotalPremium'))}</td>
                <td>{angka2(nilai(h, P + 'TotalSharePercentageClaim'))}</td>
                <td>{angka2(nilai(h, P + 'TotalClaim'))}</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </Panel>

      <Panel judul={JUDUL_NONPROP.installment}>
        <div className="field nbti__medan">
          <span className="field__label">{KOLOM_ANGSURAN.installment}</span>
          <span className="nbti__nilai">{nilai(h, P + 'Installment')}</span>
        </div>
        {daftar(h, ANGSURAN).map((a, i) => (
          <div key={i} className="nbti__angsuran-np">
            <div className="field nbti__medan">
              <span className="field__label">{JUDUL_NONPROP.currency}</span>
              <span className="nbti__nilai">{a.Currency ?? ''}</span>
            </div>
            {/* masterDetail -> flow action `InstallmentList` (Section ber-pyReadOnly; angka tanpa pyDecimalPlaces) */}
            <Grid baris={daftar(h, `${ANGSURAN}(${i + 1}).InstallmentList`)} kolom={KOLOM_RINCI} sajian={POLA_INTI} />
          </div>
        ))}
      </Panel>
    </>
  )
}
