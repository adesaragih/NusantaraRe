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
//
// Tata letak mengikuti screenshot layar Pega (work owner 06-10-2026: "tampilan NB Treaty Non Prop belum sesuai"):
// grid selebar isinya tanpa patah baris, grid total sebaris berdampingan, spreading berkolom `% RNM Share`,
// "Installment" sebaris dengan nilainya, angsuran per mata uang dapat dilipat dengan rincian bernomor.

import { Panel } from '../../../../inti/frontend/components/ui/dasar'
import { MASTER, POLIS, daftar, nilai, type Baris, type Halaman, type Pilihan } from '../api'
import { KOLOM_ANGSURAN } from '../labels'
import { labelTreatyType } from '../medan'
import {
  JUDUL_NONPROP,
  KOLOM_FAKULTATIF,
  KOLOM_LIMIT,
  KOLOM_RINCI,
  KOLOM_SHARE,
  KOLOM_SPREADING_NP,
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
import { kelasAngka, sajikan, type Sajian } from '../sajian'
import { TEMPAT_LABEL_NON_EDM, tempatTerbuka, type Tempat } from '../tempat'
import InputAngka from './InputAngka'
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

interface PropsGrid {
  baris: Baris[]
  kolom: Kolom[]
  sajian?: Sajian
  /** kolom nomor baris di depan (grid rincian angsuran Pega) */
  nomor?: boolean
}

function Grid({ baris, kolom, sajian = DUA, nomor = false }: PropsGrid) {
  return (
    <div className="table-wrap nbti__grid-np">
      <table>
        <thead>
          <tr>
            {nomor && <th scope="col" aria-label="No" />}
            {kolom.map((c, j) => (
              <th
                key={`${c.m}-${j}`}
                scope="col"
                // kepala ikut rata kanan bila kolomnya berisi angka, sejajar dengan selnya
                className={baris.some((b) => kelasAngka(b[c.m]) !== undefined) ? 'nbti__angka' : undefined}
              >
                {c.label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.map((b, i) => (
            <tr key={i}>
              {nomor && <td className="nbti__urut-np">{i + 1}</td>}
              {kolom.map((c, j) => (
                <td key={`${c.m}-${j}`} className={kelasAngka(b[c.m])}>
                  {teksSel(b, c, sajian)}
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
          kolom={kolomTotal(t)}
        />
      ))}
    </div>
  )
}

export interface PropsDetailNonProp {
  halaman: Halaman
  /** Admin pemegang berkas - HANYA admin (keputusan work owner 07-10-2026: "HANYA ADMIN YANG BISA EDIT";
   *  `DetailDeptHeadTreatyIn_UW` S88 ber-`pyEditOptions=Auto` disimpangkan sadar); grid spreading lalu
   *  terbuka menurut `spreadingTerbuka` (FacultativeShare 0/''). */
  sunting: boolean
  /** Tempat berperan pelaku (`Layar.tempat`, tiket 05). */
  tempat: Tempat
  opsiSpreading: Pilihan[]
  onUbahBaris: (jalur: string, i: number, kunci: string, v: string) => void
  onRefresh: (aksi: string, indeks: number) => void
}

export default function DetailNonProp({
  halaman: h,
  sunting,
  tempat,
  opsiSpreading,
  onUbahBaris,
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
        <div className="field nbti__medan nbti__medan-np">
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
          <div className="field nbti__medan nbti__medan-np">
            <span className="field__label">{JUDUL_NONPROP.fakultatifPersen}</span>
            <span className="nbti__nilai">{nilai(h, MASTER + 'FacultativeShare')}</span>
          </div>
          <Grid baris={daftarMaster(h, 'LimitFacShareSummaryList')} kolom={KOLOM_FAKULTATIF} />
          <Totals h={h} daftarTotal={TOTAL_FAKULTATIF} />
        </Panel>
      )}

      {/* S73 NOHEADER: SUB_SECTION `SpreadingRiskList` tanpa judul wadah */}
      <Wadah>
        {/* perintah work owner 06-10-2026: Add / Delete dibuang, Treaty Type hanya-baca ("ga boleh di ubah lagi") */}
        <div className="table-wrap nbti__grid-np">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_SPREADING_NP.treatyType}</th>
                <th scope="col" className="nbti__angka">{KOLOM_SPREADING_NP.rnmShare}</th>
                <th scope="col" className="nbti__angka">{KOLOM_SPREADING_NP.share}</th>
                <th scope="col" className="nbti__angka">{KOLOM_SPREADING_NP.premium}</th>
                <th scope="col" className="nbti__angka">{KOLOM_SPREADING_NP.claimPct}</th>
                <th scope="col" className="nbti__angka">{KOLOM_SPREADING_NP.claim}</th>
              </tr>
            </thead>
            <tbody>
              {spreading.map((b, i) => (
                <tr key={i}>
                  <td>
                    {/* hanya-baca: dropdown ro ber-RD BrowseReinsuranceType_RD - tampil .Note dari .ID */}
                    {labelTreatyType(b, opsiSpreading)}
                  </td>
                  {/* porsi RNM per treaty (.Pct master), hanya-baca */}
                  <td className="nbti__angka">{angka2(b.SplitRNMSharePct)}</td>
                  {(['SharePercentage', 'PremiumSpreaded', 'ClaimPercentage', 'ClaimSpreaded'] as const).map((k) => (
                    <td key={k} className="nbti__angka">
                      {terbuka && (k === 'SharePercentage' || k === 'ClaimPercentage') ? (
                        <InputAngka
                          label={k === 'SharePercentage' ? KOLOM_SPREADING_NP.share : KOLOM_SPREADING_NP.claimPct}
                          value={b[k] ?? ''}
                          sajian={DUA}
                          onChange={(v) => onUbahBaris(SPREADING, i, k, v)}
                          onBlur={() => onRefresh('CountSpreading', i + 1)}
                        />
                      ) : (
                        angka2(b[k])
                      )}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
            <tfoot>
              <tr>
                <td>{KOLOM_SPREADING_NP.total}</td>
                <td />
                <td className="nbti__angka">{angka2(nilai(h, POLIS + 'TotalSharePercentagePremium'))}</td>
                <td className="nbti__angka">{angka2(nilai(h, POLIS + 'TotalPremium'))}</td>
                <td className="nbti__angka">{angka2(nilai(h, POLIS + 'TotalSharePercentageClaim'))}</td>
                <td className="nbti__angka">{angka2(nilai(h, POLIS + 'TotalClaim'))}</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </Wadah>

      {/* S74 NOHEADER (pyIncludeHeader=false): tanpa judul wadah; "Installment" = label sel, sebaris dengan nilainya */}
      <Wadah>
        <div className="nbti__kolom nbti__kolom-np">
          <div className="field nbti__medan">
            <span className="field__label">{KOLOM_ANGSURAN.installment}</span>
            <span className="nbti__nilai">{nilai(h, POLIS + 'Installment')}</span>
          </div>
        </div>
        {/* grid "Currency" masterDetail: tiap mata uang dapat dilipat, rinciannya flow action `InstallmentList` */}
        <div className="table-wrap nbti__grid-np nbti__angsuran-np">
          <table>
            <thead>
              <tr>
                <th scope="col">{JUDUL_NONPROP.currency}</th>
              </tr>
            </thead>
            <tbody>
              {daftar(h, ANGSURAN).map((a, i) => (
                <tr key={i}>
                  <td>
                    <details className="nbti__rinci-np" open>
                      <summary>{a.Currency ?? ''}</summary>
                      <Grid baris={daftar(h, `${ANGSURAN}(${i + 1}).InstallmentList`)} kolom={KOLOM_RINCI} nomor />
                    </details>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Wadah>
      {/* LABEL "EDM": wadah S2 `TreatyMasterInEDM && IsEDMInputOnNB != true` - mustahil sesudah
          InputPolicyTreatyInDetail_NonProp langkah 10 (audit silang P3), tidak dirender. */}
    </>
  )
}
