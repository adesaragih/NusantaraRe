// Popup isi Product Name - tombol "View" di samping Product Name layar Premium
// List Detail (permintaan work owner 05-10-2026: "tombol atau apapun yang bisa
// untuk melihat apa isi dari product name ini").
//
// BACA-SAJA. Sumbernya tabel Master Product Name Life (keputusan work owner
// 05-10-2026); judul bagian, label medan, dan kepala kolom datang dari server
// sehingga sama dengan layar Master Product Name Life.

import { useEffect, useState } from 'react'

import { Field, Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilRateProduk,
  ambilRincianProduk,
  ambilRiskProduk,
  type MedanRincianProduk,
  type RateProduk,
  type RincianProduk,
  type RiskProduk,
} from '../api'
import { pemisahRibuan } from '../angka'
import { RINCIAN_PRODUK } from '../labels'
import { tanggalTampil } from '../tanggal'
import AreaTeks from './AreaTeks'

/** Teks server → tampilan: angka berpemisah ribuan, tanggal dd/mm/yyyy. */
export function nilaiTampil(nilai: string, jenis: MedanRincianProduk['jenis']): string {
  if (jenis === 'angka') return pemisahRibuan(nilai)
  if (jenis === 'tanggal') return tanggalTampil(nilai)
  return nilai
}

/** Teks sepanjang ini atau berbaris banyak tampil sebagai kotak teks, bukan satu baris. */
const BATAS_SEBARIS = 60

/** R/I Rate yang sedang dilihat: nama (sel R/I Rate) dan ID-nya. */
interface RateTerbuka {
  id: string
  nama: string
}

export default function ModalRincianProduk({ produkID, onTutup }: { produkID: string; onTutup: () => void }) {
  const [isi, setIsi] = useState<RincianProduk | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  // View Rate baris PLAN LIST (05-10-2026) - tabel rate tampil di bawah grid plan.
  const [rateBuka, setRateBuka] = useState<RateTerbuka | null>(null)
  // View R/I Risk (05-10-2026) - tabel risk tampil di bawah bagian TREATY NAME.
  const [riskBuka, setRiskBuka] = useState(false)

  useEffect(() => {
    let hidup = true
    ambilRincianProduk(produkID).then(
      (r) => {
        if (hidup) setIsi(r)
      },
      (e: unknown) => {
        if (hidup) setGalat(e)
      },
    )
    return () => {
      hidup = false
    }
  }, [produkID])

  return (
    <Modal judul={RINCIAN_PRODUK.judul} onTutup={onTutup} labelBatal={RINCIAN_PRODUK.tutup} lebar>
      {galat !== null && <Gagal galat={galat} />}
      {galat === null && isi === null && <Memuat />}
      {isi !== null && (
        <div className="pl-rincian">
          {isi.bagian.map((b) => (
            <section key={b.judul} className="pl-rincian__bagian">
              <h4 className="pl-offer__subjudul">{b.judul}</h4>
              <div className="form-grid">
                {b.medan.map((m) => {
                  const teks = nilaiTampil(m.nilai, m.jenis)
                  return m.jenis === 'teks' && (teks.length > BATAS_SEBARIS || teks.includes('\n')) ? (
                    <div key={m.label} className="pl-rincian__lebar">
                      <AreaTeks label={m.label} value={teks} readOnly />
                    </div>
                  ) : m.label === RINCIAN_PRODUK.medanRisk && (isi.riRiskId ?? '') !== '' ? (
                    // Tombol View MENEMPEL di kanan kotak R/I Risk Name, seperti tombol
                    // Choose di form (permintaan work owner 05-10-2026).
                    <div key={m.label} className="pl-dp-pilih">
                      <Field label={m.label} value={teks} onChange={() => {}} readOnly />
                      <button
                        type="button"
                        className="btn btn--ghost"
                        aria-label={riskBuka ? RINCIAN_PRODUK.tutupRisk : RINCIAN_PRODUK.viewRisk}
                        title={riskBuka ? RINCIAN_PRODUK.tutupRisk : RINCIAN_PRODUK.viewRisk}
                        onClick={() => {
                          setRiskBuka((v) => !v)
                        }}
                      >
                        {riskBuka ? RINCIAN_PRODUK.tutupSingkat : RINCIAN_PRODUK.tombol}
                      </button>
                    </div>
                  ) : (
                    <Field key={m.label} label={m.label} value={teks} onChange={() => {}} readOnly />
                  )
                })}
              </div>
              {riskBuka && b.medan.some((m) => m.label === RINCIAN_PRODUK.medanRisk) && (
                <TabelRisk
                  riRiskId={isi.riRiskId ?? ''}
                  nama={b.medan.find((m) => m.label === RINCIAN_PRODUK.medanRisk)?.nilai ?? ''}
                />
              )}
            </section>
          ))}
          {isi.grid.map((g) => (
            <section key={g.judul} className="pl-rincian__bagian">
              <h4 className="pl-offer__subjudul">{g.judul}</h4>
              {g.baris.length === 0 ? (
                <Kosong pesan={RINCIAN_PRODUK.kosong} />
              ) : (
                <div className="pl-rincian__tabel">
                  <table className="inbox__tabel">
                    <thead>
                      <tr>
                        {g.kolom.map((k) => (
                          <th key={k.label}>{k.label}</th>
                        ))}
                        {g.kunci !== undefined && <th className="table__actions" />}
                      </tr>
                    </thead>
                    <tbody>
                      {g.baris.map((baris, i) => {
                        const rateId = g.kunci?.[i] ?? ''
                        const terbuka = rateBuka !== null && rateBuka.id === rateId
                        return (
                          <tr key={i} className="inbox__baris">
                            {baris.map((v, j) => (
                              <td key={j}>{nilaiTampil(v, g.kolom[j]?.jenis ?? 'teks')}</td>
                            ))}
                            {g.kunci !== undefined && (
                              <td className="table__actions">
                                {rateId !== '' && (
                                  <button
                                    type="button"
                                    className="btn btn--ghost btn--sm"
                                    onClick={() => {
                                      setRateBuka(terbuka ? null : { id: rateId, nama: baris[baris.length - 1] ?? '' })
                                    }}
                                  >
                                    {terbuka ? RINCIAN_PRODUK.tutupRate : RINCIAN_PRODUK.viewRate}
                                  </button>
                                )}
                              </td>
                            )}
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
              )}
              {g.kunci !== undefined && rateBuka !== null && (
                <TabelRate key={rateBuka.id} riRateId={rateBuka.id} nama={rateBuka.nama} />
              )}
            </section>
          ))}
        </div>
      )}
    </Modal>
  )
}

/** Isi satu R/I Rate - grid `ViewRate` Master Product Name Life (05-10-2026). */
function TabelRate({ riRateId, nama }: { riRateId: string; nama: string }) {
  const [rate, setRate] = useState<RateProduk | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    ambilRateProduk(riRateId).then(
      (r) => {
        if (hidup) setRate(r)
      },
      (e: unknown) => {
        if (hidup) setGalat(e)
      },
    )
    return () => {
      hidup = false
    }
  }, [riRateId])

  return (
    <div className="pl-rincian__rate">
      <h5 className="pl-offer__subjudul">
        {RINCIAN_PRODUK.judulRate}: {nama}
      </h5>
      {galat !== null && <Gagal galat={galat} />}
      {galat === null && rate === null && <Memuat />}
      {rate !== null && rate.baris.length === 0 && <Kosong pesan={RINCIAN_PRODUK.kosong} />}
      {rate !== null && rate.baris.length > 0 && (
        <div className="pl-rincian__tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {RINCIAN_PRODUK.kolomRate.map((k) => (
                  <th key={k}>{k}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rate.baris.map((b, i) => (
                <tr key={i} className="inbox__baris">
                  <td>{b.id}</td>
                  <td>{b.usedBy}</td>
                  <td>{b.gender}</td>
                  <td>{b.contract}</td>
                  <td>{b.age}</td>
                  <td>{b.rate}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {rate !== null && rate.terpotong && <p className="pl-offer__kurang">{RINCIAN_PRODUK.rateTerpotong}</p>}
    </div>
  )
}

/** Isi R/I Risk Name produk - tabel `RIRISK_LIFE` (05-10-2026; dulu view, tabel sejak migrasi inti 938-940). */
function TabelRisk({ riRiskId, nama }: { riRiskId: string; nama: string }) {
  const [risk, setRisk] = useState<RiskProduk | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    ambilRiskProduk(riRiskId).then(
      (r) => {
        if (hidup) setRisk(r)
      },
      (e: unknown) => {
        if (hidup) setGalat(e)
      },
    )
    return () => {
      hidup = false
    }
  }, [riRiskId])

  return (
    <div className="pl-rincian__rate">
      <h5 className="pl-offer__subjudul">
        {RINCIAN_PRODUK.judulRisk}: {nama}
      </h5>
      {galat !== null && <Gagal galat={galat} />}
      {galat === null && risk === null && <Memuat />}
      {risk !== null && risk.baris.length === 0 && <Kosong pesan={RINCIAN_PRODUK.kosong} />}
      {risk !== null && risk.baris.length > 0 && (
        <div className="pl-rincian__tabel">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {RINCIAN_PRODUK.kolomRisk.map((k) => (
                  <th key={k}>{k}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {risk.baris.map((b, i) => (
                <tr key={i} className="inbox__baris">
                  <td>{b.id}</td>
                  <td>{b.usedBy}</td>
                  <td>{b.age}</td>
                  <td>{b.year}</td>
                  <td>{b.month}</td>
                  <td>{b.risk}</td>
                  <td>{b.contract}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {risk !== null && risk.terpotong && <p className="pl-offer__kurang">{RINCIAN_PRODUK.rateTerpotong}</p>}
    </div>
  )
}
