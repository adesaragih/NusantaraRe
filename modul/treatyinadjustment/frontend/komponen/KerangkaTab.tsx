// Perender kerangka bangkitan — isi tab Old dan New, apa adanya dari ekspor.
//
// ⛔ Komponen ini TIDAK memutuskan tata letak; `kerangka.gen.ts` yang
// memutuskannya (urutan butir, judul blok, kolom grid, desimal per sel).
// Yang diputuskan di sini hanya cara MENAMPILKAN satu jenis butir.
//
// ⛔ Seluruh medan di dalam tab BACA-SAJA di kedua panel. Medan isian tab
// sisi New memang dapat disunting di sistem lama (mode Edit), tetapi tanpa
// jalur simpan suntingan itu tidak ke mana-mana; yang dapat disunting di
// layar ini hanya kepala form New.

import { Field, Kosong } from '../../../../inti/frontend/components/ui/dasar'
import type { SisiPenyesuaian } from '../api'
import { golongan } from '../ekspor/golongan'
import type { ButirKerangka, GridKerangka, MedanKerangka } from '../ekspor/jenis'
import { KERANGKA_INCLUDE } from '../ekspor/kerangka.gen'
import { syaratTerpenuhi } from '../ekspor/syarat'
import { PENYESUAIAN } from '../labelsPenyesuaian'
import { AreaBacaSaja, BelumDibangun, Centang, MedanTakAda, TanggalBacaSaja, selNilai } from './medan'

/**
 * Include yang ISINYA tidak dibangun — Retro, keputusan §17
 * (`treatyin/docs/KEPUTUSAN-PENYELARASAN-REPO.md`): jarang dipakai. Wadah
 * tab dan syaratnya tetap dibangun; isinya kotak "jarang".
 */
export const INCLUDE_RETRO: ReadonlySet<string> = new Set([
  'TreatyInFacultativeShareCalculation',
  'TreatyInFacultativeShareCalculationOldData',
  'TreatyInFacultativeRetro',
])

export interface KonteksKerangka {
  /** Halaman panel ini: `OLDDATA` (Old) atau akar (New). */
  sisi: SisiPenyesuaian
  /** Halaman akar `TreatyIn` — sebagian ikatan Old menunjuknya. */
  akar: SisiPenyesuaian
  /** Halaman akar + `ViewState` sesi — tempat SEMUA syarat dinilai. */
  halaman: Readonly<Record<string, string>>
}

const tanpaAksi = () => undefined

function halamanDari(dari: GridKerangka['dari'], k: KonteksKerangka): SisiPenyesuaian {
  return dari === 'akar' ? k.akar : k.sisi
}

/** Grid bangkitan; kolom bersyarat sel (`Auto Calculate` @313578) dinilai per kolom. */
export function GridEkspor({ g, k }: { g: GridKerangka; k: KonteksKerangka }) {
  const hal = halamanDari(g.dari, k)
  const ada = Object.prototype.hasOwnProperty.call(hal.larik, g.larik)
  const baris = hal.larik[g.larik] ?? []
  const tampil = g.kunci.map((_, i) => i).filter((i) => {
    const s = g.syaratSel[i]
    return s === null || s === undefined || syaratTerpenuhi([s], k.halaman)
  })
  const total = tampil.reduce((a, i) => a + (g.lebar[i] ?? 0), 0) || 1
  return (
    <div className="tria__grid">
      <div className="table-wrap">
        <table className="tria__tabel">
          {/* Lebar DARI ekspor sebagai perbandingan — tata letak responsif. */}
          <colgroup>
            {tampil.map((i) => (
              <col key={i} style={{ width: `${(((g.lebar[i] ?? 0) / total) * 100).toFixed(2)}%` }} />
            ))}
          </colgroup>
          <thead>
            <tr>
              {tampil.map((i) => (
                <th key={i} scope="col">
                  {g.kolom[i]}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={tampil.length || 1}>
                  <Kosong
                    pesan={PENYESUAIAN.tanpaBaris}
                    petunjuk={ada ? PENYESUAIAN.petunjukGridKosong : PENYESUAIAN.takAdaDiWarisan}
                  />
                </td>
              </tr>
            )}
            {baris.map((b, r) => (
              <tr key={r}>
                {tampil.map((i) => {
                  const kunci = g.kunci[i] ?? ''
                  const fmt = g.format[i] ?? ''
                  const v = b[kunci] ?? ''
                  if (fmt === 'pxCheckbox') {
                    return (
                      <td key={i}>
                        <input type="checkbox" checked={v === 'true'} disabled readOnly aria-label={g.kolom[i]} />
                      </td>
                    )
                  }
                  const jenis = golongan(kunci, fmt)
                  return (
                    <td key={i} className={jenis === 'teks' || jenis === 'tanggal' ? undefined : 'tria__angka'}>
                      {selNilai(jenis, v, g.desimal[i] ?? null)}
                    </td>
                  )
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {g.templatBaris?.endsWith('!pyGridRowDetails') === true && (
        <span className="tria__redup">
          {PENYESUAIAN.rincianBaris} ({g.templatBaris})
        </span>
      )}
    </div>
  )
}

/** Medan bangkitan — baca-saja, kontrolnya menurut `pyFormat`. */
export function MedanEkspor({ m, k }: { m: MedanKerangka; k: KonteksKerangka }) {
  const hal = halamanDari(m.dari, k)
  const label = m.label !== '' ? m.label : (m.caption ?? '')
  if (!Object.prototype.hasOwnProperty.call(hal.medan, m.kunci)) return <MedanTakAda label={label} />
  const v = hal.medan[m.kunci] ?? ''
  switch (m.format) {
    case 'pxTextArea':
      return <AreaBacaSaja label={label} nilai={v} />
    case 'pxDateTime':
      // ⛔ Nilai TERSIMPAN ke kotak tanggal — bukan terjemahan tampil (§5).
      return <TanggalBacaSaja label={label} nilai={v} />
    case 'pxCheckbox':
      return <Centang label={m.caption ?? label} nilai={v} bacaSaja onUbah={tanpaAksi} />
    default:
      return <Field label={label} value={selNilai(golongan(m.kunci, m.format), v, m.desimal)} onChange={tanpaAksi} readOnly />
  }
}

/** Satu daftar butir kerangka. */
export function RenderKerangka({ isi, k }: { isi: readonly ButirKerangka[]; k: KonteksKerangka }) {
  return (
    <>
      {isi.map((b) => {
        if (!syaratTerpenuhi(b.syarat, k.halaman)) return null
        switch (b.t) {
          case 'blok':
            return (
              <div key={b.at} className="tria__blok">
                {b.judul !== '' && <h5 className="tria__subjudul">{b.judul}</h5>}
                <RenderKerangka isi={b.anak} k={k} />
              </div>
            )
          case 'grid':
            return <GridEkspor key={b.at} g={b} k={k} />
          case 'medan':
            return <MedanEkspor key={b.at} m={b} k={k} />
          case 'teks':
            return (
              <span key={b.at} className="tria__teks-sel">
                {b.teks}
              </span>
            )
          case 'tombol':
            // ⛔ Tombol tab menjalankan Activity (Update Total, Refresh,
            // Submit, Decline offer …) — jalur yang belum ada: MATI.
            return (
              <button key={b.at} type="button" className="btn btn--ghost btn--sm" disabled title={PENYESUAIAN.tombolTulisMati}>
                {b.label}
              </button>
            )
          case 'include': {
            if (INCLUDE_RETRO.has(b.nama)) {
              return <BelumDibangun key={b.at} judul={PENYESUAIAN.retroJarang} petunjuk={PENYESUAIAN.retroJarangPetunjuk} />
            }
            const anak = KERANGKA_INCLUDE[b.nama]
            if (anak === undefined) {
              return <BelumDibangun key={b.at} judul={PENYESUAIAN.belumDibangun} petunjuk={`${b.nama}: ${PENYESUAIAN.includeTakAda}`} />
            }
            return <RenderKerangka key={b.at} isi={anak} k={k} />
          }
          default:
            return null
        }
      })}
    </>
  )
}
