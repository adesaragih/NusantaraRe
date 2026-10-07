// Cabang NonProp baru `DetailPolicyTreatyInAddendum` S11 (`.IsNewPolicyNonProp = 1`):
//
//   INCLUDE `DetailPolicyTreatyInAddPremi`  S2 "Previous Premium" / S7 "Current Premium" / S12 "Total Difference"
//                                          (grid XOL per mata uang, expand pane `DetailPolicyAddPremiDetail` = rincian
//                                          per layer) + S17 grid `.SpreadingRiskList` (RO, 2 desimal)
//   S12 `.Installment` (pyMaxLength 2)     change -> FillPaymentInstallmentEDMT
//   S13 "Installment Data Information"     grid `.ListInstallment` Currency / Total, expand pane `InstallmentList`
//
// Label dan format: `../xol.ts`. ⛔ Grid XOL dan rinciannya seluruhnya hanya-baca; nol perhitungan di sini.

import { Panel } from '../../../../inti/frontend/components/ui/dasar'
import { POLIS, daftar, jalurAnak, nilai, type Halaman, type Pilihan } from '../api'
import { BAGIAN, KOLOM_ANGSURAN, KOLOM_SPREADING } from '../labels'
import { MEDAN_ANGSURAN_NP, SAJIAN_SPREADING_NP, labelTreatyType } from '../medan'
import { sajikan } from '../sajian'
import {
  ANAK_ANGSURAN,
  ANAK_XOL,
  DAFTAR_ANGSURAN,
  GRID_XOL,
  KOLOM_ANGSURAN_NP,
  KOLOM_RINCI_ANGSURAN,
  KOLOM_XOL,
  KOLOM_XOL_RINCI,
} from '../xol'
import { GridBaca, GridLipat } from './Grid'
import Wadah from './Wadah'

const SPREADING = POLIS + 'SpreadingRiskList'
const dua = (v: string | undefined) => sajikan(v ?? '', SAJIAN_SPREADING_NP)

export interface PropsAddPremi {
  halaman: Halaman
  /** Sel `.Installment` S12 dapat diisi (admin yang boleh bekerja - backend menerimanya hanya dari admin). */
  suntingAngsuran: boolean
  /** Teks `.TreatyType` grid spreading (RD `BrowseReinsuranceType_RD`). */
  opsiJenisReas: Pilihan[]
  onUbah: (jalur: string, v: string) => void
  onRefreshAngsuran: () => void
}

export default function AddPremi({
  halaman: h,
  suntingAngsuran,
  opsiJenisReas,
  onUbah,
  onRefreshAngsuran,
}: PropsAddPremi) {
  const angsuran = MEDAN_ANGSURAN_NP.jalur
  return (
    <>
      {GRID_XOL.map((g) => (
        // S2 / S7 / S12: wadah berjudul (pyIncludeHeader=true)
        <Panel key={g.daftar} judul={g.judul}>
          <GridLipat
            baris={daftar(h, g.daftar)}
            kolom={KOLOM_XOL}
            rinci={(_, i) => (
              <GridBaca baris={daftar(h, jalurAnak(g.daftar, i + 1, ANAK_XOL))} kolom={KOLOM_XOL_RINCI} />
            )}
          />
        </Panel>
      ))}

      {/* S17 NOHEADER: grid spreading hanya-baca, 2 desimal, tanpa Add / Delete */}
      <Wadah>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_SPREADING.treatyType}</th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_SPREADING.share}
                </th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_SPREADING.premium}
                </th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_SPREADING.claimPct}
                </th>
                <th scope="col" className="edmt__angka">
                  {KOLOM_SPREADING.claim}
                </th>
              </tr>
            </thead>
            <tbody>
              {daftar(h, SPREADING).map((b, i) => (
                <tr key={i}>
                  <td>{labelTreatyType(b, opsiJenisReas)}</td>
                  <td className="edmt__angka">{dua(b.SharePercentage)}</td>
                  <td className="edmt__angka">{dua(b.PremiumSpreaded)}</td>
                  <td className="edmt__angka">{dua(b.ClaimPercentage)}</td>
                  <td className="edmt__angka">{dua(b.ClaimSpreaded)}</td>
                </tr>
              ))}
            </tbody>
            <tfoot>
              <tr>
                <td>{KOLOM_SPREADING.total}</td>
                <td className="edmt__angka">{dua(nilai(h, POLIS + 'TotalSharePercentagePremium'))}</td>
                <td className="edmt__angka">{dua(nilai(h, POLIS + 'TotalPremium'))}</td>
                <td className="edmt__angka">{dua(nilai(h, POLIS + 'TotalSharePercentageClaim'))}</td>
                <td className="edmt__angka">{dua(nilai(h, POLIS + 'TotalClaim'))}</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </Wadah>

      {/* S12 (.Installment) + S13 "Installment Data Information" (pyIncludeHeader=true) */}
      <Panel judul={BAGIAN.angsuran}>
        <div className="edmt__aksi">
          <label className="field__label">{KOLOM_ANGSURAN.installment}</label>
          {suntingAngsuran ? (
            <input
              className="field__input edmt__pendek"
              aria-label={KOLOM_ANGSURAN.installment}
              maxLength={MEDAN_ANGSURAN_NP.panjangMaks}
              value={nilai(h, angsuran)}
              onChange={(e) => onUbah(angsuran, e.target.value)}
              onBlur={onRefreshAngsuran}
            />
          ) : (
            <span>{nilai(h, angsuran)}</span>
          )}
        </div>
        <GridLipat
          baris={daftar(h, DAFTAR_ANGSURAN)}
          kolom={KOLOM_ANGSURAN_NP}
          rinci={(_, i) => (
            <GridBaca
              baris={daftar(h, jalurAnak(DAFTAR_ANGSURAN, i + 1, ANAK_ANGSURAN))}
              kolom={KOLOM_RINCI_ANGSURAN}
              nomor
            />
          )}
        />
      </Panel>
    </>
  )
}
