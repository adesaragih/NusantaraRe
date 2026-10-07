// Tab Spreading kasus FIRE (tiket 48) - port `NB FacIn\Section\InputInwardFacultativeDtl.xml` tab "Spreading" dengan
// gambar layar Pega DEV (work owner 05-10-2026: "sama seperti dipega. awalnya pembuatan otomatis saat di isi % Share RNM
// (seprti dipega skrng), namun bisa juga dibuat manual, dengan klik add"):
// - % Share RNM (sel 80, `.OfferFacIn.PercentShare`): diubah -> `POST …/spreading/hitung-share` (CountPremiAndTSINusantaraRe_ACT
//   + GetKapasitasTreaty: TSI / premi Nusantara Re per coverage dan baris spreading otomatis QS / SPL) - TIDAK menyimpan;
// - Copy Spreading: template Type Treaty / % Share (Add / Delete, Σ % Share <= 100 seperti `cekSpreadingFactIn_Act`),
//   Copy To All Spreading -> `POST …/spreading/salin` (CopyToAllSpreadingFire_ACT);
// - lokasi -> item -> coverage -> detail (TSI 100% … Premium Nusantara Re) + baris spreading (baca-saja, seperti pengguna
//   biasa di Pega); total per lokasi; ringkasan per treaty dan per mata uang (`SummarySpreading_Section`);
// - Save (`IsThereAnyObjectLocation_Act`) -> `PUT …/spreading`: hitung ulang lalu simpan.
// Seluruh hitungan di backend (desimal eksak, ADR-0034); layar hanya menampilkan.
//
// Keputusan agent (tiket 48): S-1 Copy Spreading tampil untuk semua pengguna (Pega: hanya IsGroup / user tertentu) - kata
// work owner "bisa juga dibuat manual, dengan klik add". S-2 baris spreading per coverage baca-saja (grid editable Pega hanya
// untuk user ID tertentu). S-3 Total Accumulation = tahap berikut. S-4 hitung-share berjeda 500 ms; jawaban lama dibuang.
//
// Proteksi template Copy Spreading (work owner 05-10-2026: "kalo spreadingnya langsung spl tidak boleh, harus ada QS nya")
// = `NB FacIn\Activity\CalcultePersentageSpeading_Act.xml`, dipanggil saat Type Treaty grid Copy Spreading berubah
// (InputInwardFacultativeDtl, `SpreadingList.pxResults`). "Can not proceed spreading without QS" diperiksa BACKEND saja:
// "QS" dicari di REINSURANCETYPE.NOTE jenis treaty baris pertama (RDBList GetTreatyName), bukan di nama dropdown - layar
// tidak memegang NOTE, jadi galat 400 `…/salin` ditampilkan apa adanya. Type Treaty kembar ("Treaty Type can't be same")
// diperiksa juga di layar (ID, sama dengan backend) dan menahan Copy To All (keputusan agent S-5).

import { useEffect, useRef, useState } from 'react'

import { Field, Gagal, Kosong, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { desimalSah, jumlahDesimal } from '../../../../inti/frontend/lib/desimal'
import { formatNumber } from '../../../../inti/frontend/lib/format'
import {
  ambilSpreading,
  hitungShareSpreading,
  salinSpreading,
  simpanSpreading,
  type TampilanSpreading,
  type TemplateSpreading,
  type TotalSpreading,
} from '../api'
import { SPREADING as S, TEKS_OBJEK, TEKS_SPREADING as T } from '../labels'

const DESIMAL = 4
const JEDA_MS = 500

/** Angka sah: kosong atau desimal bertitik. */
const angkaSah = (v: string) => v.trim() === '' || desimalSah(v.trim())

/** Total desimal (teks, eksak) melebihi 100 - tanpa konversi ke float. */
export function lebihDari100(total: string): boolean {
  const [bulat = '0', pecahan = ''] = total.trim().replace(/^\+/, '').split('.')
  if (bulat.startsWith('-')) return false
  const n = BigInt(bulat === '' ? '0' : bulat)
  return n > 100n || (n === 100n && /[1-9]/.test(pecahan))
}

/** Σ % Share template Copy Spreading melebihi 100 (`cekSpreadingFactIn_Act`). */
export const templateLebih = (t: TemplateSpreading[]) => lebihDari100(jumlahDesimal(t.map((x) => x.sharePercentage)).total)

/** Ada Type Treaty yang dipilih lebih dari sekali (`Local.Counter>=2`); ID dipangkas spasi seperti backend. */
export function treatyKembar(t: TemplateSpreading[]): boolean {
  const dipilih = t.map((x) => x.treatyType.trim()).filter((x) => x !== '')
  return new Set(dipilih).size !== dipilih.length
}

const angka = (v: string | undefined) => formatNumber(v ?? '', DESIMAL)

function Buka({ buka, onKlik }: { buka: boolean; onKlik: () => void }) {
  return (
    <button
      type="button"
      className={'btn btn--ghost btn--sm nbf-buka' + (buka ? ' nbf-buka--terbuka' : '')}
      aria-label={TEKS_OBJEK.bukaBaris}
      aria-expanded={buka}
      onClick={onKlik}
    >
      ▸
    </button>
  )
}

/** Tabel total (per lokasi, atau ringkasan per treaty dengan Total First Loss). */
function TabelTotal({ kolom, baris, firstLoss }: { kolom: readonly string[]; baris: TotalSpreading[]; firstLoss?: boolean }) {
  return (
    <div className="nbf-cov-wrap">
      <table className="nbf-tabel">
        <thead>
          <tr>
            {kolom.map((k) => (
              <th key={k} scope="col">
                {k}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {baris.length === 0 && (
            <tr>
              <td colSpan={kolom.length}>{T.kosong}</td>
            </tr>
          )}
          {baris.map((t, i) => (
            <tr key={`${t.currency}-${t.treatyType}-${i}`}>
              <td>{t.currency}</td>
              <td>{t.treatyName}</td>
              <td className="nbf-angka">{angka(t.sharePercentage)}</td>
              <td className="nbf-angka">{angka(t.claimSpreaded)}</td>
              <td className="nbf-angka">{angka(t.tsiSpreaded)}</td>
              <td className="nbf-angka">{angka(t.claimEstimation)}</td>
              {firstLoss && <td className="nbf-angka">{angka(t.claimAmountIdr)}</td>}
              <td className="nbf-angka">{angka(t.premiumSpreaded)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function TabSpreading({ caseId }: { caseId: string }) {
  const [v, setV] = useState<TampilanSpreading | null>(null)
  const [share, setShare] = useState('')
  const [galat, setGalat] = useState<unknown>(null)
  const [menghitung, setMenghitung] = useState(false)
  const [menyimpan, setMenyimpan] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  const [terbuka, setTerbuka] = useState<string[]>([])
  const nomor = useRef(0)
  const jadwal = useRef<number | undefined>(undefined)

  useEffect(() => {
    let batal = false
    ambilSpreading(caseId).then(
      (h) => {
        if (batal) return
        setV(h)
        setShare(h.percentShare)
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

  const balik = (k: string) => setTerbuka((t) => (t.includes(k) ? t.filter((x) => x !== k) : [...t, k]))

  /** Jalankan satu panggilan hitung; jawaban lama dibuang. */
  function hitung(janji: () => Promise<TampilanSpreading>) {
    const n = ++nomor.current
    setMenghitung(true)
    setGalat(null)
    janji().then(
      (h) => {
        if (n !== nomor.current) return
        setV(h)
        setMenghitung(false)
      },
      (err: unknown) => {
        if (n !== nomor.current) return
        setGalat(err)
        setMenghitung(false)
      },
    )
  }

  function ubahShare(nilai: string) {
    setShare(nilai)
    setTersimpan(false)
    window.clearTimeout(jadwal.current)
    // Backend menolak di luar 0..100 (kontrak spreading c3); jangan memanggilnya untuk nilai yang pasti ditolak.
    if (!angkaSah(nilai) || nilai.trim() === '' || lebihDari100(nilai)) return
    jadwal.current = window.setTimeout(() => hitung(() => hitungShareSpreading(caseId, nilai.trim())), JEDA_MS)
  }

  function ubahTemplate(template: TemplateSpreading[]) {
    setV((x) => (x ? { ...x, template } : x))
    setTersimpan(false)
  }

  async function simpan() {
    if (v === null) return
    setMenyimpan(true)
    setGalat(null)
    try {
      const h = await simpanSpreading(caseId, { ...v, percentShare: share.trim() })
      setV(h)
      setShare(h.percentShare)
      setTersimpan(true)
    } catch (err) {
      setGalat(err)
    } finally {
      setMenyimpan(false)
    }
  }

  if (v === null) return galat ? <Gagal galat={galat} /> : <Memuat />

  const opsiTreaty = v.treaty.map((t) => ({ value: t.id, label: t.name }))
  const lebih = templateLebih(v.template)
  const kembar = treatyKembar(v.template)

  return (
    <div className="nbf-objek nbf-cov-tab nbf-spread">
      <Gagal galat={galat} />
      {tersimpan && <div className="alert alert--ok">{T.tersimpan}</div>}
      {v.pesan.map((p) => (
        <div key={p} className="alert alert--warn">
          {p}
        </div>
      ))}

      {/* % Share RNM + Copy Spreading */}
      <section className="nbf-lapis nbf-spread__atas">
        <div className="nbf-spread__share">
          <Field
            label={S.percentShare}
            value={share}
            onChange={ubahShare}
            required
            error={!angkaSah(share) ? T.angka : lebihDari100(share) ? T.percentShareLebih : undefined}
          />
          {menghitung && <span className="muted">{T.menghitung}</span>}
        </div>
        <div className="nbf-spread__copy">
          <h5 className="nbf-akum__kartu-judul">{S.judulCopy}</h5>
          <div className="nbf-cov-wrap">
            <table className="nbf-tabel">
              <thead>
                <tr>
                  {S.kolomTemplate.map((k) => (
                    <th key={k} scope="col">
                      {k}
                    </th>
                  ))}
                  <th scope="col" className="table__actions">
                    <button
                      type="button"
                      className="btn btn--sm"
                      onClick={() => ubahTemplate([...v.template, { treatyType: '', treatyName: '', sharePercentage: '' }])}
                    >
                      {S.tambah}
                    </button>
                  </th>
                </tr>
              </thead>
              <tbody>
                {v.template.length === 0 && (
                  <tr>
                    <td colSpan={3}>{T.kosong}</td>
                  </tr>
                )}
                {v.template.map((t, i) => (
                  <tr key={i}>
                    <td>
                      <Pilih
                        label={S.kolomTemplate[0]}
                        value={t.treatyType}
                        kosong={T.pilihTreaty}
                        opsi={opsiTreaty}
                        onChange={(id) =>
                          ubahTemplate(
                            v.template.map((x, k) =>
                              k === i ? { ...x, treatyType: id, treatyName: v.treaty.find((y) => y.id === id)?.name ?? '' } : x,
                            ),
                          )
                        }
                      />
                    </td>
                    <td>
                      <Field
                        label={S.kolomTemplate[1]}
                        value={t.sharePercentage}
                        onChange={(s) => ubahTemplate(v.template.map((x, k) => (k === i ? { ...x, sharePercentage: s } : x)))}
                        error={angkaSah(t.sharePercentage) ? undefined : T.angka}
                      />
                    </td>
                    <td className="table__actions">
                      <button
                        type="button"
                        className="btn btn--sm"
                        onClick={() => ubahTemplate(v.template.filter((_, k) => k !== i))}
                      >
                        {S.hapus}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {lebih && <div className="field__error">{T.shareLebih}</div>}
          {kembar && <div className="field__error">{T.treatySama}</div>}
          <div>
            <button
              type="button"
              className="btn btn--primary btn--sm"
              disabled={lebih || kembar || v.template.length === 0 || menghitung}
              onClick={() => hitung(() => salinSpreading(caseId, share.trim(), v.template))}
            >
              {S.salinSemua}
            </button>
          </div>
        </div>
      </section>

      {/* Lokasi -> item -> coverage -> spreading */}
      <div className="nbf-cov-wrap">
        <table className="nbf-tabel">
          <thead>
            <tr>
              <th scope="col" />
              {S.kolomLokasi.map((k) => (
                <th key={k} scope="col">
                  {k}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {v.lokasi.length === 0 && (
              <tr>
                <td colSpan={4}>
                  <Kosong pesan={T.kosong} />
                </td>
              </tr>
            )}
            {v.lokasi.map((l, o) => [
              <tr key={`o${o}`} className={terbuka.includes(`o${o}`) ? 'nbf-baris--terbuka' : undefined}>
                <td>
                  <Buka buka={terbuka.includes(`o${o}`)} onKlik={() => balik(`o${o}`)} />
                </td>
                <td>{l.objectNo}</td>
                <td>{l.objectName}</td>
                <td>{l.location}</td>
              </tr>,
              terbuka.includes(`o${o}`) && (
                <tr key={`od${o}`} className="nbf-objek__detail">
                  <td colSpan={4}>
                    <div className="nbf-lapis nbf-lapis--objek">
                      <div className="nbf-cov-wrap">
                        <table className="nbf-tabel">
                          <thead>
                            <tr>
                              <th scope="col" />
                              {S.kolomItem.map((k) => (
                                <th key={k} scope="col">
                                  {k}
                                </th>
                              ))}
                            </tr>
                          </thead>
                          <tbody>
                            {l.items.map((it, i) => [
                              <tr key={`i${i}`} className={terbuka.includes(`i${o}-${i}`) ? 'nbf-baris--terbuka' : undefined}>
                                <td>
                                  <Buka buka={terbuka.includes(`i${o}-${i}`)} onKlik={() => balik(`i${o}-${i}`)} />
                                </td>
                                <td>{it.itemType}</td>
                                <td>{it.currency}</td>
                                <td className="nbf-angka">{angka(it.tsi)}</td>
                                <td className="nbf-angka">{angka(it.totalGrossPremi)}</td>
                                <td className="nbf-angka">{angka(it.totalPremiumRnm)}</td>
                              </tr>,
                              terbuka.includes(`i${o}-${i}`) && (
                                <tr key={`id${i}`} className="nbf-objek__detail">
                                  <td colSpan={6}>
                                    <div className="nbf-lapis nbf-lapis--item">
                                      <div className="nbf-cov-wrap">
                                        <table className="nbf-tabel">
                                          <thead>
                                            <tr>
                                              <th scope="col" />
                                              {S.kolomCoverage.map((k) => (
                                                <th key={k} scope="col">
                                                  {k}
                                                </th>
                                              ))}
                                            </tr>
                                          </thead>
                                          <tbody>
                                            {it.coverages.map((c, n) => {
                                              const kunci = `c${o}-${i}-${n}`
                                              return [
                                                <tr key={kunci} className={terbuka.includes(kunci) ? 'nbf-baris--terbuka' : undefined}>
                                                  <td>
                                                    <Buka buka={terbuka.includes(kunci)} onKlik={() => balik(kunci)} />
                                                  </td>
                                                  <td>
                                                    <span className="nbf-cov-kode">{c.oldId}</span>
                                                  </td>
                                                  <td className="nbf-angka">{angka(c.rate)}</td>
                                                  <td className="nbf-angka">{angka(c.premium)}</td>
                                                  <td className="nbf-angka">{angka(c.premiNusantaraRe)}</td>
                                                </tr>,
                                                terbuka.includes(kunci) && (
                                                  <tr key={`${kunci}d`} className="nbf-objek__detail">
                                                    <td colSpan={5}>
                                                      <div className="nbf-spread__coverage">
                                                        <dl className="nbf-spread__angka">
                                                          {(
                                                            [
                                                              [S.detailCoverage.tsi, c.tsi],
                                                              [S.detailCoverage.tsiLiability, c.tsiLiability],
                                                              [S.detailCoverage.tsiNusantaraRe, c.tsiNusantaraRe],
                                                              [S.detailCoverage.premium, c.premium],
                                                              [S.detailCoverage.premiNusantaraRe, c.premiNusantaraRe],
                                                            ] as const
                                                          ).map(([lbl, nilai]) => (
                                                            <div key={lbl}>
                                                              <dt>{lbl}</dt>
                                                              <dd className="nbf-angka">{angka(nilai)}</dd>
                                                            </div>
                                                          ))}
                                                        </dl>
                                                        <div className="nbf-cov-wrap">
                                                          <table className="nbf-tabel">
                                                            <thead>
                                                              <tr>
                                                                {S.kolomSpreading.map((k) => (
                                                                  <th key={k} scope="col">
                                                                    {k}
                                                                  </th>
                                                                ))}
                                                              </tr>
                                                            </thead>
                                                            <tbody>
                                                              {c.spreading.length === 0 && (
                                                                <tr>
                                                                  <td colSpan={6}>{T.kosong}</td>
                                                                </tr>
                                                              )}
                                                              {c.spreading.map((b, k) => (
                                                                <tr key={`${b.treatyType}-${k}`}>
                                                                  <td>{b.treatyName || b.treatyType}</td>
                                                                  <td className="nbf-angka">{angka(b.sharePercentage)}</td>
                                                                  <td className="nbf-angka">{angka(b.tsiGrossSpreaded)}</td>
                                                                  <td className="nbf-angka">{angka(b.tsiSpreaded)}</td>
                                                                  <td className="nbf-angka">{angka(b.claimEstimation)}</td>
                                                                  <td className="nbf-angka">{angka(b.premiumSpreaded)}</td>
                                                                </tr>
                                                              ))}
                                                            </tbody>
                                                          </table>
                                                        </div>
                                                      </div>
                                                    </td>
                                                  </tr>
                                                ),
                                              ]
                                            })}
                                          </tbody>
                                        </table>
                                      </div>
                                    </div>
                                  </td>
                                </tr>
                              ),
                            ])}
                          </tbody>
                        </table>
                      </div>
                      <div className="nbf-spread__total">
                        <TabelTotal kolom={S.kolomTotal} baris={l.total} />
                      </div>
                    </div>
                  </td>
                </tr>
              ),
            ])}
          </tbody>
        </table>
      </div>

      {/* Ringkasan (SummarySpreading_Section) */}
      <section className="nbf-spread__ringkasan">
        <TabelTotal kolom={S.kolomRingkasanTreaty} baris={v.ringkasanTreaty} firstLoss />
        <div className="nbf-cov-wrap">
          <table className="nbf-tabel">
            <thead>
              <tr>
                {S.kolomRingkasanMataUang.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {v.ringkasanMataUang.length === 0 && (
                <tr>
                  <td colSpan={5}>{T.kosong}</td>
                </tr>
              )}
              {v.ringkasanMataUang.map((r) => (
                <tr key={r.currency}>
                  <td>{r.currency}</td>
                  <td className="nbf-angka">{angka(r.tsiTopRisk)}</td>
                  <td className="nbf-angka">{angka(r.tsi)}</td>
                  <td className="nbf-angka">{angka(r.lol)}</td>
                  <td className="nbf-angka">{angka(r.premium)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <div className="nbf-objek__kaki">
        <button type="button" className="btn btn--primary" onClick={() => void simpan()} disabled={menyimpan || !angkaSah(share) || lebihDari100(share)}>
          {menyimpan ? T.menghitung : S.simpan}
        </button>
      </div>
    </div>
  )
}
