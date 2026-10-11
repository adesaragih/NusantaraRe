// Tab Scoring Risk kasus FIRE (tiket 54) - port `NB FacIn\Section\CallDualScoringRisk.xml` -> `ScoringRisk.xml` +
// `ScoringRiskForm1..4` dengan gambar layar Pega (work owner 06-10-2026):
// - kepala: The Insured, Top Risk, Dominant Risk (pilih lokasi lain, `SetDataOccupationScoringRisk_Act`), Period, Risk
//   Location 1 / 2, Status, Score;
// - 8 faktor berbobot; tiap sub-faktor = satu centang per kolom (Top Risk / Dominant Risk), saling-eksklusif seperti Pega;
//   skor sub-faktor tampil di samping centang terpilih; isian sistem (Occupation, band Loss Ratio, Total Sum Insured)
//   dikunci; Loss Ratio % diketik -> band dicentang backend (`SetLossRatio_Act`);
// - Total per kolom, Scoring Result + Risk Classifications, Recomendation, Required approval (From / Due to);
// - Save -> `PUT …/scoring`. Seluruh hitungan (poin `ValueScoringRisk_act`, `SetScore_Act`, `ScoringResult`, isi otomatis)
//   di backend; layar hanya menampilkan.
//
// Keputusan work owner 06-10-2026: "hanya marketing" = Pega "Write" (posisi kasus tahap marketing & bukan on-going) -
// backend mengirim `bolehUbah`; selain itu baca-saja ("Read"). Cacat Pega dibetulkan (backend). Kurs TSI = Cedant.
// Keputusan agent S54-1: centang terpilih diklik lagi = dikosongkan (Pega mengunci centang lain sampai yang terpilih
// dilepas; di sini memilih centang lain langsung memindahkan pilihan). S54-2 kunci silang Growing Trees ditiru persis.

import { useEffect, useRef, useState } from 'react'

import { Gagal, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { jumlahDesimal } from '../../../../inti/frontend/lib/desimal'
import { angka4 } from '../angka4'
import { ambilScoring, hitungScoring, simpanScoring, type HasilScoring, type KolomScoring, type TampilanScoring } from '../api'
import {
  FAKTOR_SCORING,
  HASIL_SCORING as H,
  SCORING as S,
  TEKS_SCORING as T,
  type FaktorScoring,
  type SubFaktorScoring,
} from '../labels'
import { dariKetikan } from './IsianAngka4'

const JEDA_MS = 300

/** Pilihan baru satu sub-faktor: klik pilihan terpilih = kosongkan (S54-1), selain itu pindah ke pilihan itu. */
export function pilihScoring(kolom: KolomScoring, kunci: string, n: number): KolomScoring {
  const kini = kolom.pilihan[kunci] ?? 0
  return { ...kolom, pilihan: { ...kolom.pilihan, [kunci]: kini === n ? 0 : n } }
}

/** Isian % Loss Ratio terkunci bila Nil 5 / 3 / 1 tahun atau No Information dipilih (Pega pyDisabledWhen CB1/2/3/10). */
export function persenLossRatioTerkunci(kolom: KolomScoring): boolean {
  return [1, 2, 3, 10].includes(kolom.pilihan.lossRatio ?? 0)
}

/**
 * Kunci silang Pega (S54-2, work owner 06-10-2026 "ditiru persis seperti Pega") - `ScoringRiskForm2.xml` pyDisabledWhen
 * `… || .DataScoringRiskList(k).ObjectConditions.ClassConstruction.ChechBox4 !=""`: pilihan "Not Available (Growing
 * Trees)" dikunci bila Class Construction = Growing Trees (pilihan 4) pada kolom yang sama. Portable Fire Extinguisher
 * pilihan 3 di kedua kolom; Hydrant / Sprinkler pilihan 2 hanya kolom Top Risk (kolom Dominant Risk tanpa syarat ini).
 */
const KUNCI_SILANG: readonly { kunci: string; n: number; kolom: readonly number[] }[] = [
  { kunci: 'portableFireExtinguisher', n: 3, kolom: [0, 1] },
  { kunci: 'hydrant', n: 2, kolom: [0] },
  { kunci: 'sprinkler', n: 2, kolom: [0] },
]

export function terkunciSilang(k: number, kolom: KolomScoring, kunci: string, n: number): boolean {
  return (kolom.pilihan.classConstruction ?? 0) === 4 && KUNCI_SILANG.some((x) => x.kunci === kunci && x.n === n && x.kolom.includes(k))
}

/** Tanggal kabel DD-MM-YYYY -> tampil DD/MM/YYYY. */
const tanggal = (v: string) => (v === '' ? T.kosong : v.replace(/-/g, '/'))

/** Nada warna kelas skor / rekomendasi: baik (hijau), sedang (kuning), buruk (merah), '' = netral. */
export function nadaKelas(nama: string): '' | 'baik' | 'sedang' | 'buruk' {
  if (nama === 'Good' || nama === 'Above Average') return 'baik'
  if (nama === 'Average') return 'sedang'
  if (nama === 'Below Average' || nama === 'Poor') return 'buruk'
  return ''
}

const NADA_REKOMENDASI: Record<string, '' | 'baik' | 'sedang' | 'buruk'> = {
  acceptableNotApproval: 'baik',
  acceptableWithApproval: 'sedang',
  notAcceptable: 'buruk',
}

/** Subtotal skor satu faktor pada satu kolom (eksak, `jumlahDesimal`); '' bila belum ada sub-faktor yang dicentang. */
export function subtotalFaktor(kolom: KolomScoring | undefined, f: FaktorScoring): string {
  if (!kolom || !f.sub.some((s) => (kolom.pilihan[s.kunci] ?? 0) > 0)) return ''
  const j = jumlahDesimal(f.sub.map((s) => kolom.skor[s.kunci] ?? ''))
  return j.takTerbaca.length > 0 ? '' : j.total
}

/**
 * `awal` hanya untuk uji render statis (tanpa DOM, efek tidak berjalan); layar memuat lewat `ambilScoring`.
 * `saatTersimpan` - skor tersimpan dikabarkan ke General (kolom Risk Scoring, `Section\Periode.xml`).
 */
export default function TabScoring({
  caseId,
  awal = null,
  saatTersimpan,
}: {
  caseId: string
  awal?: TampilanScoring | null
  saatTersimpan?: (skor: { finalScore: string; noteFinalScore: string }) => void
}) {
  const [v, setV] = useState<TampilanScoring | null>(awal)
  const [galat, setGalat] = useState<unknown>(null)
  const [menghitung, setMenghitung] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  const nomor = useRef(0)
  const jadwal = useRef<number | undefined>(undefined)

  useEffect(() => {
    let batal = false
    ambilScoring(caseId).then(
      (h) => {
        if (!batal) setV(h)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
      window.clearTimeout(jadwal.current)
    }
  }, [caseId])

  if (v === null) return galat ? <Gagal galat={galat} /> : <Memuat />

  /** Ubah lalu jadwalkan hitung ulang di backend; jawaban lama dibuang. */
  function ubahHitung(baru: TampilanScoring) {
    setV(baru)
    setTersimpan(false)
    window.clearTimeout(jadwal.current)
    jadwal.current = window.setTimeout(() => {
      const n = ++nomor.current
      setMenghitung(true)
      setGalat(null)
      hitungScoring(caseId, baru).then(
        (h) => {
          if (n === nomor.current) {
            setV(h)
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

  const ubahKolom = (k: number, kolom: KolomScoring) => ubahHitung({ ...v, kolom: v.kolom.map((x, i) => (i === k ? kolom : x)) })

  async function simpan() {
    if (v === null) return
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = await simpanScoring(caseId, v)
      setV(h)
      setTersimpan(true)
      saatTersimpan?.({ finalScore: h.finalScore, noteFinalScore: h.noteFinalScore })
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  const adaDominan = (v.kolom[1]?.occupationNote ?? '') !== ''
  const NAMA_KOLOM = [S.topRisk, S.dominantRisk]

  /** Sel satu kolom untuk satu pilihan: centang + skor (bila terpilih). */
  const sel = (k: number, s: SubFaktorScoring, n: number) => {
    const kolom = v.kolom[k]
    if (!kolom) return <td key={k} />
    const dipilih = (kolom.pilihan[s.kunci] ?? 0) === n
    const sistem = (s.terkunci ?? []).includes(n)
    const silang = terkunciSilang(k, kolom, s.kunci, n)
    const kunci = !v.bolehUbah || sistem || (k === 1 && !adaDominan) || silang
    return (
      <td key={k} className={'nbf-score__sel' + (dipilih ? ' nbf-score__sel--pilih' : '')}>
        <label title={silang ? T.terkunciGrowingTrees : sistem ? T.diisiSistem : undefined}>
          <input
            type="checkbox"
            aria-label={`${NAMA_KOLOM[k]} - ${s.label || s.kunci} - ${s.opsi[n - 1]}`}
            checked={dipilih}
            disabled={kunci}
            onChange={() => ubahKolom(k, pilihScoring(kolom, s.kunci, n))}
          />
          <span className="nbf-score__poin">{dipilih ? angka4(kolom.skor[s.kunci] ?? '') : ''}</span>
        </label>
      </td>
    )
  }

  /** Baris isian % Loss Ratio (sesudah pilihan ke-3 Net Loss Ratio). */
  const barisPersen = () => (
    <tr key="persen" className="nbf-score__persen">
      <td className="nbf-score__opsi">{S.lossRatioInput}</td>
      {v.kolom.slice(0, 2).map((kolom, k) => (
        <td key={k}>
          <input
            className="nbf-angka"
            inputMode="decimal"
            aria-label={`${S.lossRatioInput} - ${NAMA_KOLOM[k]}`}
            value={kolom.persenLossRatio}
            disabled={!v.bolehUbah || persenLossRatioTerkunci(kolom) || (k === 1 && !adaDominan)}
            onChange={(e) => setV({ ...v, kolom: v.kolom.map((x, i) => (i === k ? { ...x, persenLossRatio: dariKetikan(e.target.value) } : x)) })}
            onBlur={() => ubahHitung(v)}
          />{' '}
          %
        </td>
      ))}
    </tr>
  )

  const hasil = (kunci: keyof HasilScoring) => v.hasil[kunci]
  const skorAkhir = v.finalScore === '' ? T.kosong : angka4(v.finalScore)
  const totalKolom = (k: number) => {
    const t = v.kolom[k]?.total ?? ''
    return t === '' ? T.kosong : angka4(t)
  }

  return (
    <div className="nbf-objek nbf-score">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{T.tersimpan}</div>}
      {!v.bolehUbah && <div className="alert">{T.bacaSaja}</div>}
      {v.pesan.map((p) => (
        <div key={p} className="alert alert--warn">
          {p}
        </div>
      ))}

      <h4 className="nbf-score__judul">{S.judul}</h4>

      <section className="nbf-score__kepala">
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
          <dd className="nbf-score__skor" data-nada={nadaKelas(v.noteFinalScore)}>
            {skorAkhir}
            {v.noteFinalScore !== '' && <span className="nbf-score__pil">{v.noteFinalScore}</span>}
          </dd>
        </dl>
        <dl>
          <dt>{S.topRisk}</dt>
          <dd>{v.kolom[0]?.occupationNote || T.kosong}</dd>
          <dt>{S.lokasi1}</dt>
          <dd>{v.kolom[0]?.topRiskLocation || T.kosong}</dd>
          <dt>{S.dominantRisk}</dt>
          <dd className="nbf-score__dominan">
            {v.bolehUbah ? (
              <Pilih
                label={S.dominantRisk}
                value={v.kolom[1]?.indexLocation ?? ''}
                onChange={(x) => ubahKolom(1, { ...(v.kolom[1] as KolomScoring), indexLocation: x })}
                opsi={v.lokasiDominan.map((l) => ({ value: l.indexLocation, label: l.label }))}
              />
            ) : (
              v.kolom[1]?.occupationNote || T.kosong
            )}
            {v.bolehUbah && !adaDominan && <small className="nbf-score__catatan">{T.tanpaDominan}</small>}
          </dd>
          <dt>{S.lokasi2}</dt>
          <dd>{v.kolom[1]?.topRiskLocation || T.kosong}</dd>
        </dl>
      </section>

      <div className="nbf-score__bungkus">
        <table className="nbf-tabel nbf-score__tabel">
          <thead>
            <tr>
              <th scope="col">
                {S.kolom[1]} / {S.kolom[3]}
              </th>
              <th scope="col">{S.kolom[4]}</th>
              <th scope="col">{S.kolom[5]}</th>
            </tr>
          </thead>
          {FAKTOR_SCORING.map((f) => (
            <tbody key={f.no} className="nbf-score__faktor">
              <tr className="nbf-score__baris-faktor">
                <th scope="rowgroup">
                  <span className="nbf-score__no">{f.no}</span>
                  {f.faktor}
                  <span className="nbf-score__bobot">
                    {S.kolom[2]} {f.bobot}
                  </span>
                </th>
                {NAMA_KOLOM.map((nama, k) => {
                  const t = subtotalFaktor(v.kolom[k], f)
                  return (
                    <td key={nama} className="nbf-score__sel">
                      <span className="nbf-score__slot" />
                      <span className="nbf-score__poin nbf-score__poin--sub">{t === '' ? '' : angka4(t)}</span>
                    </td>
                  )
                })}
              </tr>
              {f.sub.flatMap((s) => [
                ...(s.label !== ''
                  ? [
                      <tr key={`${s.kunci}-judul`} className="nbf-score__judul-sub">
                        <td colSpan={3}>
                          {s.label}
                          {s.bobot !== '' && <span className="nbf-score__bobot">{s.bobot}</span>}
                        </td>
                      </tr>,
                    ]
                  : []),
                ...s.opsi.flatMap((o, oi) => {
                  const n = oi + 1
                  const terpilih = v.kolom.some((x) => (x.pilihan[s.kunci] ?? 0) === n)
                  const tr = (
                    <tr key={`${s.kunci}-${n}`} className={terpilih ? 'nbf-score__baris--pilih' : undefined}>
                      <td className="nbf-score__opsi">
                        {o}
                        {(s.terkunci ?? []).includes(n) && <span className="nbf-score__sistem">{T.diisiSistem}</span>}
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

      <dl className="nbf-score__ringkas nbf-score__ringkas--kaki">
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
        <div className="nbf-score__utama" data-nada={nadaKelas(v.noteFinalScore)}>
          <dt>{S.hasil}</dt>
          <dd>
            {skorAkhir}
            {v.noteFinalScore !== '' && <span className="nbf-score__pil">{v.noteFinalScore}</span>}
          </dd>
        </div>
      </dl>

      <section className="nbf-score__hasil">
        <div className="nbf-score__kartu">
          <h5>{S.klasifikasi}</h5>
          <ul>
            {S.kelas.map((k) => (
              <li key={k.nama} className={v.noteFinalScore === k.nama ? 'nbf-score__aktif' : undefined} data-nada={nadaKelas(k.nama)}>
                <span className="nbf-score__rentang">{k.rentang}</span>
                {k.nama}
              </li>
            ))}
          </ul>
        </div>
        <div className="nbf-score__kartu">
          <h5>{S.rekomendasi}</h5>
          <ul>
            {H.rekomendasi.map((r) => (
              <li key={r.kunci} className={hasil(r.kunci) ? 'nbf-score__aktif' : undefined} data-nada={NADA_REKOMENDASI[r.kunci]}>
                <span className="nbf-score__tanda">{hasil(r.kunci) ? '✓' : ''}</span>
                {r.label}
              </li>
            ))}
          </ul>
        </div>
        <div className="nbf-score__kartu">
          <h5>{S.persetujuan}</h5>
          <div className="nbf-score__setuju">
            <div>
              <h6>{S.dari}</h6>
              <ul>
                {H.dari.map((r) => (
                  <li key={r.kunci} className={hasil(r.kunci) ? 'nbf-score__aktif' : undefined}>
                    <span className="nbf-score__tanda">{hasil(r.kunci) ? '✓' : ''}</span>
                    {r.label}
                  </li>
                ))}
              </ul>
            </div>
            <div>
              <h6>{S.karena}</h6>
              <ul>
                {H.karena.map((r) => (
                  <li key={r.kunci} className={hasil(r.kunci) ? 'nbf-score__aktif' : undefined}>
                    <span className="nbf-score__tanda">{hasil(r.kunci) ? '✓' : ''}</span>
                    {r.label}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </section>

      {v.bolehUbah && (
        <div className="nbf-objek__kaki">
          {menghitung && <span className="muted">{T.menghitung}</span>}
          <button type="button" className="btn nbf-tema__tombol nbf-tema__tombol--utama btn--primary" onClick={() => void simpan()} disabled={menyimpan || menghitung}>
            {menyimpan ? T.menyimpan : T.simpan}
          </button>
        </div>
      )}
    </div>
  )
}
