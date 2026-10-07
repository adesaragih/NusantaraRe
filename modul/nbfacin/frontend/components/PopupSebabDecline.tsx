// Popup "Choose Cause of decline" (tiket 57) - port harness / section `NB FacIn\InputCauseOfDecline` (pre-activity
// `GetCauseOfDecline`): grid T/C (centang prioritas I / II, Option, Remarks, % RNM, % SOB) dan grid Non - T/C; Check
// menjalankan validasi `CekCauseOfDecline_ACT`; Submit (muncul sesudah Check lolos, Pega `Protect.CARI1 = 1`) menyusun
// Comment (`SetCauseOfDecline_ACT`). Kolom isian baris yang tidak dicentang dinonaktifkan (keputusan agent K57-6).

import { useState } from 'react'

import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import type { BarisSebabTc, PilihanSebabTc } from '../api'
import { SEBAB_DECLINE as S, TEKS_SEBAB_DECLINE as T } from '../labels'
import { dariKetikan } from './IsianAngka4'
import {
  centangPrioritas,
  komentarDecline,
  optionAktif,
  pctRnmAktif,
  pctSobAktif,
  pilihanKosong,
  remarksAktif,
  validasiSebab,
} from './sebabDecline'

export interface HasilSebabDecline {
  pilihan: PilihanSebabTc[]
  nonTc: string
  komentar: string
}

export default function PopupSebabDecline({
  tc,
  nonTc,
  pctOtomatis,
  awal,
  onTutup,
  onSubmit,
}: {
  tc: BarisSebabTc[]
  nonTc: string[]
  pctOtomatis: { rate: string; riComm: string }
  awal: { pilihan: PilihanSebabTc[]; nonTc: string } | null
  onTutup: () => void
  onSubmit: (h: HasilSebabDecline) => void
}) {
  const [pilihan, setPilihan] = useState<PilihanSebabTc[]>(() =>
    tc.map((b) => awal?.pilihan.find((p) => p.termCondition === b.termCondition) ?? pilihanKosong(b.termCondition)),
  )
  const [pilihNonTc, setPilihNonTc] = useState(awal?.nonTc ?? '')
  const [pesan, setPesan] = useState<string[] | null>(null)

  const ubah = (baru: PilihanSebabTc[]) => {
    setPilihan(baru)
    setPesan(null)
  }
  const ubahBaris = (i: number, isi: Partial<PilihanSebabTc>) => ubah(pilihan.map((x, j) => (j === i ? { ...x, ...isi } : x)))
  const otomatis = (term: string) => (term === S.term.rate ? pctOtomatis.rate : term === S.term.riComm ? pctOtomatis.riComm : '')
  const tcMati = pilihNonTc !== ''
  const lolos = pesan !== null && pesan.length === 0

  return (
    <Modal
      judul={S.judul}
      onTutup={onTutup}
      lebar
      aksi={
        <>
          <button type="button" className="btn" onClick={() => setPesan(validasiSebab(pilihan, pilihNonTc))}>
            {S.cek}
          </button>
          {lolos && (
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => onSubmit({ pilihan, nonTc: pilihNonTc, komentar: komentarDecline(pilihan, pilihNonTc) })}
            >
              {S.submit}
            </button>
          )}
        </>
      }
    >
      <div className="nbf-sebab">
        {pesan !== null &&
          (pesan.length > 0 ? (
            <ul className="nbf-sebab__galat">
              {pesan.map((p) => (
                <li key={p}>{p}</li>
              ))}
            </ul>
          ) : (
            <div className="alert alert--ok">{T.lolos}</div>
          ))}
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel nbf-sebab__tc">
            <thead>
              <tr>
                {S.prioritas.map((p) => (
                  <th key={p} scope="col" title={T.pilihPrioritas}>
                    {p}
                  </th>
                ))}
                {S.kolom.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {tc.map((b, i) => {
                const p = pilihan[i] ?? pilihanKosong(b.termCondition)
                const terpilih = p.prioritas !== ''
                const mati = (aktif: boolean) => tcMati || !terpilih || !aktif
                return (
                  <tr key={b.termCondition} className={terpilih ? 'nbf-sebab__pilih' : undefined}>
                    {(['I', 'II'] as const).map((pr) => (
                      <td key={pr} className="nbf-sebab__centang">
                        <input
                          type="checkbox"
                          aria-label={`${pr} - ${b.termCondition}`}
                          checked={p.prioritas === pr}
                          disabled={tcMati}
                          onChange={() => ubah(centangPrioritas(pilihan, i, pr, pr === 'I' ? otomatis(b.termCondition) : ''))}
                        />
                      </td>
                    ))}
                    <td>{b.termCondition}</td>
                    <td>
                      {b.opsi.length > 0 ? (
                        <select
                          aria-label={`${S.kolom[1]} - ${b.termCondition}`}
                          value={p.description1}
                          disabled={mati(optionAktif(b.termCondition))}
                          onChange={(e) => ubahBaris(i, { description1: e.target.value })}
                        >
                          <option value="" />
                          {b.opsi.map((o) => (
                            <option key={o} value={o}>
                              {o}
                            </option>
                          ))}
                        </select>
                      ) : (
                        <input
                          aria-label={`${S.kolom[1]} - ${b.termCondition}`}
                          value={p.description1}
                          disabled={mati(optionAktif(b.termCondition))}
                          onChange={(e) => ubahBaris(i, { description1: e.target.value })}
                        />
                      )}
                    </td>
                    <td>
                      <input
                        aria-label={`${S.kolom[2]} - ${b.termCondition}`}
                        value={p.description2}
                        disabled={mati(remarksAktif(b.termCondition, p.description1))}
                        onChange={(e) => ubahBaris(i, { description2: e.target.value })}
                      />
                    </td>
                    <td>
                      <input
                        className="nbf-angka"
                        inputMode="decimal"
                        aria-label={`${S.kolom[3]} - ${b.termCondition}`}
                        value={p.pctRnm}
                        disabled={mati(pctRnmAktif(b.termCondition))}
                        onChange={(e) => ubahBaris(i, { pctRnm: dariKetikan(e.target.value) })}
                      />
                    </td>
                    <td>
                      <input
                        className="nbf-angka"
                        inputMode="decimal"
                        aria-label={`${S.kolom[4]} - ${b.termCondition}`}
                        value={p.pctSob}
                        disabled={mati(pctSobAktif(b.termCondition))}
                        onChange={(e) => ubahBaris(i, { pctSob: dariKetikan(e.target.value) })}
                      />
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>

        <div className="nbf-cov-wrap">
          <table className="nbf-tabel nbf-sebab__nontc">
            <thead>
              <tr>
                <th scope="col" />
                <th scope="col">{S.nonTc}</th>
              </tr>
            </thead>
            <tbody>
              {nonTc.map((n) => (
                <tr key={n} className={pilihNonTc === n ? 'nbf-sebab__pilih' : undefined}>
                  <td className="nbf-sebab__centang">
                    <input
                      type="checkbox"
                      aria-label={`${S.nonTc} - ${n}`}
                      checked={pilihNonTc === n}
                      onChange={() => {
                        const baru = pilihNonTc === n ? '' : n
                        setPilihNonTc(baru)
                        setPesan(null)
                        if (baru !== '') setPilihan(tc.map((b) => pilihanKosong(b.termCondition)))
                      }}
                    />
                  </td>
                  <td>{n}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </Modal>
  )
}
