// Isi tab polis NB Fac In hanya-baca Scoring Risk (FIRE; `PolisNB.tsx`, perintah work owner 11-10-2026 "tampilan view
// polis yang belum tergenerate tolong diperbaiki. cek tampilan nya di nb facin!"). Susunan disalin dari
// `modul/nbfacin/frontend/components/TabScoring.tsx` dalam keadaan baca-saja (pola, bukan impor): kepala dua kolom, satu
// tabel faktor menerus (pita faktor + subtotal, judul sub-faktor, pilihan dengan centang terkunci dan skor), ringkasan
// total, kartu Risk Classifications / Recomendation / Required approval. Isian % Loss Ratio, pilih Dominant Risk, pesan,
// dan Save tidak dibawa. Nilai = TERSIMPAN (keputusan work owner, `polisNBScoring.ts`).

import type { Halaman } from '../api'
import { angkaNB } from '../polisNB'
import {
  FAKTOR_SCORING_NB,
  HASIL_SCORING_NB as H,
  NADA_REKOMENDASI_NB,
  SCORING_NB as S,
  TEKS_SCORING_NB as T,
  nadaKelasNB,
  skorPolisNB,
  subtotalFaktorNB,
  type KunciHasilSkorNB,
  type SubFaktorScoringNB,
} from '../polisNBScoring'

/** Lencana kelas skor (General / kepala / ringkasan). */
export function LencanaSkorNB({ catatan, umum }: { catatan: string; umum?: boolean }) {
  if (catatan === '') return null
  return (
    <span
      className={'claimfacin__nb-skor-pil' + (umum ? ' claimfacin__nb-skor-pil--umum' : '')}
      data-nada={nadaKelasNB(catatan)}
    >
      {catatan}
    </span>
  )
}

export function TabScoringNB({ h }: { h: Halaman }) {
  const v = skorPolisNB(h)
  const adaDominan = (v.kolom[1]?.occupationNote ?? '') !== ''
  const NAMA_KOLOM = [S.topRisk, S.dominantRisk]
  const skorAkhir = v.finalScore === '' ? T.kosong : angkaNB(v.finalScore)
  const totalKolom = (k: number) => {
    const t = v.kolom[k]?.total ?? ''
    return t === '' ? T.kosong : angkaNB(t)
  }
  const tanggal = (t: string) => (t === '' ? T.kosong : t)
  const hasil = (k: KunciHasilSkorNB) => v.hasil[k]

  /** Sel satu kolom untuk satu pilihan: centang terkunci + skor (bila terpilih). */
  const sel = (k: number, s: SubFaktorScoringNB, n: number) => {
    const kolom = v.kolom[k]
    const dipilih = (kolom?.pilihan[s.kunci] ?? 0) === n
    return (
      <td key={k} className="claimfacin__nb-skor-sel">
        <label>
          <input
            type="checkbox"
            aria-label={`${NAMA_KOLOM[k]} - ${s.label || s.kunci} - ${s.opsi[n - 1]}`}
            checked={dipilih}
            disabled
            readOnly
          />
          <span className="claimfacin__nb-skor-poin">{dipilih ? angkaNB(kolom?.skor[s.kunci] ?? '') : ''}</span>
        </label>
      </td>
    )
  }

  /** Baris % Loss Ratio (sesudah pilihan ke-3 Net Loss Ratio) - nilai tersimpan. */
  const barisPersen = () => (
    <tr key="persen" className="claimfacin__nb-skor-persen">
      <td className="claimfacin__nb-skor-opsi">{S.lossRatioInput}</td>
      {v.kolom.slice(0, 2).map((kolom, k) => (
        <td key={k} className="claimfacin__nb-skor-sel">
          {kolom.persenLossRatio === '' ? '' : `${angkaNB(kolom.persenLossRatio)} %`}
        </td>
      ))}
    </tr>
  )

  return (
    <div className="claimfacin__nb-skor">
      <h4 className="claimfacin__nb-skor-judul">{S.judul}</h4>

      <section className="claimfacin__nb-skor-kepala">
        <dl>
          <dt>{S.insured}</dt>
          <dd>{v.insuredName || T.kosong}</dd>
          <dt>{S.periode}</dt>
          <dd>
            {tanggal(v.mulai)} - {tanggal(v.selesai)}
          </dd>
          <dt>{S.status}</dt>
          <dd>{v.status || T.kosong}</dd>
          <dt>{S.skor}</dt>
          <dd className="claimfacin__nb-skor-nilai" data-nada={nadaKelasNB(v.noteFinalScore)}>
            {skorAkhir}
            <LencanaSkorNB catatan={v.noteFinalScore} />
          </dd>
        </dl>
        <dl>
          <dt>{S.topRisk}</dt>
          <dd>{v.kolom[0]?.occupationNote || T.kosong}</dd>
          <dt>{S.lokasi1}</dt>
          <dd>{v.kolom[0]?.topRiskLocation || T.kosong}</dd>
          <dt>{S.dominantRisk}</dt>
          <dd>{v.kolom[1]?.occupationNote || T.kosong}</dd>
          <dt>{S.lokasi2}</dt>
          <dd>{v.kolom[1]?.topRiskLocation || T.kosong}</dd>
        </dl>
      </section>

      <div className="claimfacin__nb-skor-bungkus">
        <table className="claimfacin__nb-skor-tabel">
          <thead>
            <tr>
              <th scope="col">
                {S.kolom[1]} / {S.kolom[3]}
              </th>
              <th scope="col">{S.kolom[4]}</th>
              <th scope="col">{S.kolom[5]}</th>
            </tr>
          </thead>
          {FAKTOR_SCORING_NB.map((f) => (
            <tbody key={f.no}>
              <tr className="claimfacin__nb-skor-faktor">
                <th scope="rowgroup">
                  <span className="claimfacin__nb-skor-no">{f.no}</span>
                  {f.faktor}
                  <span className="claimfacin__nb-skor-bobot">
                    {S.kolom[2]} {f.bobot}
                  </span>
                </th>
                {NAMA_KOLOM.map((nama, k) => {
                  const t = subtotalFaktorNB(v.kolom[k], f)
                  return (
                    <td key={nama} className="claimfacin__nb-skor-sel">
                      <span className="claimfacin__nb-skor-slot" />
                      <span className="claimfacin__nb-skor-poin claimfacin__nb-skor-poin--sub">
                        {t === '' ? '' : angkaNB(t)}
                      </span>
                    </td>
                  )
                })}
              </tr>
              {f.sub.flatMap((s) => [
                ...(s.label !== ''
                  ? [
                      <tr key={`${s.kunci}-judul`} className="claimfacin__nb-skor-judul-sub">
                        <td colSpan={3}>
                          {s.label}
                          {s.bobot !== '' && <span className="claimfacin__nb-skor-bobot">{s.bobot}</span>}
                        </td>
                      </tr>,
                    ]
                  : []),
                ...s.opsi.flatMap((o, oi) => {
                  const n = oi + 1
                  const terpilih = v.kolom.some((x) => (x.pilihan[s.kunci] ?? 0) === n)
                  const tr = (
                    <tr key={`${s.kunci}-${n}`} className={terpilih ? 'claimfacin__nb-skor-baris--pilih' : undefined}>
                      <td className="claimfacin__nb-skor-opsi">
                        {o}
                        {(s.terkunci ?? []).includes(n) && (
                          <span className="claimfacin__nb-skor-sistem">{T.diisiSistem}</span>
                        )}
                      </td>
                      {sel(0, s, n)}
                      {sel(1, s, n)}
                    </tr>
                  )
                  return s.kunci === 'lossRatio' && n === 3 ? [tr, barisPersen()] : [tr]
                }),
              ])}
            </tbody>
          ))}
        </table>
      </div>

      <dl className="claimfacin__nb-skor-ringkas">
        <div>
          <dt>
            {S.total} {S.topRisk}
          </dt>
          <dd>{totalKolom(0)}</dd>
        </div>
        <div>
          <dt>
            {S.total} {S.dominantRisk}
          </dt>
          <dd>{adaDominan ? totalKolom(1) : T.kosong}</dd>
        </div>
        <div data-nada={nadaKelasNB(v.noteFinalScore)}>
          <dt>{S.hasil}</dt>
          <dd>
            {skorAkhir}
            <LencanaSkorNB catatan={v.noteFinalScore} />
          </dd>
        </div>
      </dl>

      <section className="claimfacin__nb-skor-hasil">
        <div className="claimfacin__nb-skor-kartu">
          <h5>{S.klasifikasi}</h5>
          <ul>
            {S.kelas.map((k) => (
              <li
                key={k.nama}
                className={v.noteFinalScore === k.nama ? 'claimfacin__nb-skor-aktif' : undefined}
                data-nada={nadaKelasNB(k.nama)}
              >
                <span className="claimfacin__nb-skor-rentang">{k.rentang}</span>
                {k.nama}
              </li>
            ))}
          </ul>
        </div>
        <div className="claimfacin__nb-skor-kartu">
          <h5>{S.rekomendasi}</h5>
          <ul>
            {H.rekomendasi.map((r) => (
              <li
                key={r.kunci}
                className={hasil(r.kunci) ? 'claimfacin__nb-skor-aktif' : undefined}
                data-nada={NADA_REKOMENDASI_NB[r.kunci]}
              >
                <span className="claimfacin__nb-skor-tanda">{hasil(r.kunci) ? '✓' : ''}</span>
                {r.label}
              </li>
            ))}
          </ul>
        </div>
        <div className="claimfacin__nb-skor-kartu">
          <h5>{S.persetujuan}</h5>
          <div className="claimfacin__nb-skor-setuju">
            <div>
              <h6>{S.dari}</h6>
              <ul>
                {H.dari.map((r) => (
                  <li key={r.kunci} className={hasil(r.kunci) ? 'claimfacin__nb-skor-aktif' : undefined}>
                    <span className="claimfacin__nb-skor-tanda">{hasil(r.kunci) ? '✓' : ''}</span>
                    {r.label}
                  </li>
                ))}
              </ul>
            </div>
            <div>
              <h6>{S.karena}</h6>
              <ul>
                {H.karena.map((r) => (
                  <li key={r.kunci} className={hasil(r.kunci) ? 'claimfacin__nb-skor-aktif' : undefined}>
                    <span className="claimfacin__nb-skor-tanda">{hasil(r.kunci) ? '✓' : ''}</span>
                    {r.label}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
