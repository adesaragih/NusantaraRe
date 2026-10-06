// Subsection `DetailPoliciesNonProportional` -> `DetailPolicyTreatyInNonProportional`
// (polis NB NonProporsional / XOL, `[keputusan work owner]` K8): grid master
// Limits / Share / Total Spreaded / Share Facultative (baca-saja), section
// `SpreadingRiskList`, dan jadwal angsuran per mata uang beserta rinciannya
// (`InstallmentList`, hanya-baca). Label, syarat tampil, dan aturan sunting: `../nonprop.ts`.
//
// ⛔ Nol perhitungan di sini - semua nilai datang dari backend (pilih bisnis:
// InputPolicyTreatyInDetail_NonProp; refresh `CountSpreading`).
//
// LABEL "NON EDM" (sebelum subsection) dan "EDM" (sesudahnya, sebelum subsection
// EDM yang tak terjangkau) tampil hanya bagi `OperatorID.pxInsName = <ID-operator-2>`
// di XML - di sini DIGERBANG tempat berperan tiket 05 (`tempat.ts`), tertunda
// selama pemetaan IAM kosong; begitu pemetaannya diisi, label ikut tampil.
//
// Judul wadah hanya yang ber-`pyIncludeHeader=true` (S1 Limits, S19 Share, S41 Total
// Spreaded, S51 Share Facultative); S73 (SpreadingRiskList) dan S74 (angsuran) NOHEADER -
// `Wadah` tanpa judul (W6 audit silang P3).

import { Panel } from '../../../../inti/frontend/components/ui/dasar'
import { MASTER, POLIS, daftar, nilai, type Baris, type Halaman, type Pilihan } from '../api'
import { KOLOM_ANGSURAN, KOLOM_SPREADING, TOMBOL } from '../labels'
import { labelTreatyType } from '../medan'
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
  kolomTotal,
  spreadingTerbuka,
  tampilFakultatif,
  type Kolom,
  type Total,
} from '../nonprop'
import { POLA_INTI, sajikan, type Sajian } from '../sajian'
import { TEMPAT_LABEL_NON_EDM, tempatTerbuka, type Tempat } from '../tempat'
import Wadah from './Wadah'

const SPREADING = POLIS + 'SpreadingRiskList'
const ANGSURAN = POLIS + 'ListInstallment'

/** pxNumber section NonProp dan SpreadingRiskList ber-`pyDecimalPlaces` 2 (K14). */
const DUA: Sajian = { desimal: 2 }
const angka2 = (v: string | undefined) => sajikan(v ?? '', DUA)
const TEKS = new Set(['Note', 'Currency'])

/** Teks satu sel grid: kolom tanpa properti kosong; Note/Currency apa adanya; format
 *  kolom (`Kolom.format`) bila ada; selain itu sajian grid. */
function teksSel(b: Baris, c: Kolom, sajian: Sajian): string {
  if (c.m === '') return ''
  const v = b[c.m] ?? ''
  if (TEKS.has(c.m) || c.format === 'mentah') return v
  return sajikan(v, c.m === 'DueDate' ? 'tanggal' : (c.format ?? sajian))
}

function Grid({ baris, kolom, sajian = DUA }: { baris: Baris[]; kolom: Kolom[]; sajian?: Sajian }) {
  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>
            {kolom.map((c, j) => (
              <th key={`${c.m}-${j}`} scope="col">
                {c.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.map((b, i) => (
            <tr key={i}>
              {kolom.map((c, j) => (
                <td key={`${c.m}-${j}`}>{teksSel(b, c, sajian)}</td>
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
          kolom={kolomTotal(t)}
        />
      ))}
    </div>
  )
}

export interface PropsDetailNonProp {
  halaman: Halaman
  /** Pelaku boleh bekerja di layar ini - admin MAUPUN atasan (W2 audit silang P3:
   *  `DetailDeptHeadTreatyIn_UW` S88 menyertakan subsection ini ber-`pyEditOptions=Auto`);
   *  grid spreading lalu terbuka menurut `spreadingTerbuka` (FacultativeShare 0/''). */
  sunting: boolean
  /** Tempat berperan pelaku (`Layar.tempat`, tiket 05). */
  tempat: Tempat
  opsiSpreading: Pilihan[]
  onUbahBaris: (jalur: string, i: number, kunci: string, v: string) => void
  onSetelDaftar: (jalur: string, b: Baris[]) => void
  onRefresh: (aksi: string, indeks: number) => void
}

export default function DetailNonProp({
  halaman: h,
  sunting,
  tempat,
  opsiSpreading,
  onUbahBaris,
  onSetelDaftar,
  onRefresh,
}: PropsDetailNonProp) {
  const spreading = daftar(h, SPREADING)
  const terbuka = sunting && spreadingTerbuka(h)
  return (
    <>
      {tempatTerbuka(tempat, TEMPAT_LABEL_NON_EDM) && <p className="nbti__label-np">{JUDUL_NONPROP.nonEdm}</p>}
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
            <span className="nbti__nilai">{nilai(h, MASTER + 'FacultativeShare')}</span>
          </div>
          <Grid baris={daftarMaster(h, 'LimitFacShareSummaryList')} kolom={KOLOM_FAKULTATIF} />
          <Totals h={h} daftarTotal={TOTAL_FAKULTATIF} />
        </Panel>
      )}

      {/* S73 NOHEADER: SUB_SECTION `SpreadingRiskList` tanpa judul wadah */}
      <Wadah>
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
                        {/* pyNoSelectionText */}
                        <option value="">{TOMBOL.pilihKosong}</option>
                        {opsiSpreading.map((o) => (
                          <option key={o.nilai} value={o.nilai}>
                            {o.label}
                          </option>
                        ))}
                      </select>
                    ) : (
                      // dropdown ro ber-RD BrowseReinsuranceType_RD: tampil .Note dari .ID
                      labelTreatyType(b, opsiSpreading)
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
                <td>{angka2(nilai(h, POLIS + 'TotalSharePercentagePremium'))}</td>
                <td>{angka2(nilai(h, POLIS + 'TotalPremium'))}</td>
                <td>{angka2(nilai(h, POLIS + 'TotalSharePercentageClaim'))}</td>
                <td>{angka2(nilai(h, POLIS + 'TotalClaim'))}</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </Wadah>

      {/* S74 NOHEADER (pyIncludeHeader=false): tanpa judul wadah; "Installment" = label sel */}
      <Wadah>
        <div className="field nbti__medan">
          <span className="field__label">{KOLOM_ANGSURAN.installment}</span>
          <span className="nbti__nilai">{nilai(h, POLIS + 'Installment')}</span>
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
      </Wadah>
      {/* LABEL "EDM": wadah S2 `TreatyMasterInEDM && IsEDMInputOnNB != true` - mustahil sesudah
          InputPolicyTreatyInDetail_NonProp langkah 10 (audit silang P3), tidak dirender. */}
    </>
  )
}
