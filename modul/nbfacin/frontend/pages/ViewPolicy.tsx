// Menu View Policy (tiket 71) - port tab Facultative "List Polis" section Pega `DATA-PORTAL!VIEWOPENPOLICYMARINEINBOX`
// (`D:\migrasi\RNM\DDL\ViewOpenPolicyMarineInbox.xml`): View Option + Filter + Refresh, grid `GridOpenPolicyMarine` (20 per
// halaman, link View Details), mode detail (Back, Refresh, "View Data" versi bila > 1).
//
// Sistem baru (work owner 09-10-2026 K71-1 / K71-2): cari = FACINPRODUCTION (`GetListViewPolicy_act`); versi = kolom relasional
// JSON_POLIS; detail versi yang case-nya ada di sistem baru = layar case baca-saja penuh (A71-4), polis lama Pega = ringkasan
// FACINPRODUCTION / CEDING_FACINPRODUCTION (A71-3). Tanpa DATA_JSON. Pencarian berjalan saat mengetik / ganti View Option /
// Enter (pola PopupFollowing).

import { useEffect, useRef, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { cariPolis, type BarisFollowing, type HalamanFollowing, type OpsiFollowing } from '../api'
import { DetailPolis } from '../components/DetailPolis'
import { JEDA_CARI_MS } from '../components/PopupPilihAgent'
import { tanggalFollowing } from '../components/PopupFollowing'
import { FOLLOWING as F, TEKS_FOLLOWING, VIEW_POLICY as V } from '../labels'
import InwardFacultative from './InwardFacultative'

// Detail polis (versi, ringkasan polis lama) dipindah ke components/DetailPolis.tsx - dipakai juga popup Following.
export { labelPolicyType, RingkasanPolisLama } from '../components/DetailPolis'

export default function ViewPolicy() {
  const [opsi, setOpsi] = useState<OpsiFollowing | ''>('')
  const [filter, setFilter] = useState('')
  const [hasil, setHasil] = useState<HalamanFollowing | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(false)
  const [detail, setDetail] = useState<BarisFollowing | null>(null)
  const nomorPermintaan = useRef(0)

  async function muat(o: OpsiFollowing | '', f: string, halaman: number) {
    const nomor = ++nomorPermintaan.current
    // GetListViewPolicy_act langkah 7: filter kosong -> keluar tanpa mencari.
    if (o === '' || f.trim() === '') {
      setHasil(null)
      setMemuat(false)
      return
    }
    setMemuat(true)
    setGalat(null)
    try {
      const h = await cariPolis(o, f, halaman)
      if (nomor === nomorPermintaan.current) setHasil(h)
    } catch (err) {
      if (nomor === nomorPermintaan.current) {
        setHasil(null)
        setGalat(err)
      }
    } finally {
      if (nomor === nomorPermintaan.current) setMemuat(false)
    }
  }

  useEffect(() => {
    const jadwal = window.setTimeout(() => void muat(opsi, filter, 1), JEDA_CARI_MS)
    return () => window.clearTimeout(jadwal)
  }, [opsi, filter])

  // Templat layar = menu Kelola User (work owner 09-10-2026 "template nya disamain dong kaya menu ini, biar ga beda2 sama
  // semua"; pola portal NB FacIn): kepala `inbox__kepala` + judul besar, bilah `toolbar`, tabel `inbox__tabel` berbaris klik.
  return (
    <div className="nbfacin">
      <section className="inbox nbf-tema">
        {detail !== null ? (
          <DetailPolis
            baris={detail}
            judul={`${V.judul} · ${detail.policyNo}`}
            onKembali={() => setDetail(null)}
            tampilKasus={(caseId) => <InwardFacultative kasus={{ caseId }} onBatal={() => setDetail(null)} hanyaLihat />}
          />
        ) : (
          <>
            <header className="inbox__kepala">
              <h2 className="inbox__judul nbf-tema__judul">{V.judul}</h2>
            </header>
            <form
              className="toolbar"
              onSubmit={(e) => {
                e.preventDefault()
                void muat(opsi, filter, 1)
              }}
            >
              <select
                className="field__input nbf-tema__isian nbf-vp__opsi"
                aria-label={F.viewOption}
                value={opsi}
                onChange={(e) => setOpsi(e.target.value as OpsiFollowing | '')}
              >
                <option value="">{`${F.viewOption}: ${F.pilihKosong}`}</option>
                {F.opsi.map((o) => (
                  <option key={o.value} value={o.value}>
                    {o.label}
                  </option>
                ))}
              </select>
              <input
                className="field__input nbf-tema__isian nbf-saring__kotak"
                type="search"
                aria-label={F.filter}
                placeholder={F.placeholder}
                value={filter}
                disabled={opsi === ''}
                onChange={(e) => setFilter(e.target.value)}
              />
              <button type="button" className="btn btn--sm nbf-tema__tombol" onClick={() => setFilter('')}>
                {F.refresh}
              </button>
            </form>
            {hasil && hasil.baris.length > 0 && (
              <div className="nbf-tema__pager">
                <Halaman halaman={hasil.halaman} ukuran={hasil.ukuran} total={hasil.total} onPindah={(h) => void muat(opsi, filter, h)} />
              </div>
            )}
            <table className="inbox__tabel nbf-tema__tabel">
              <thead>
                <tr>
                  {F.kolom.map((k) => (
                    <th key={k} scope="col">
                      {k}
                    </th>
                  ))}
                  <th scope="col" />
                </tr>
              </thead>
              {hasil && hasil.baris.length > 0 && (
                <tbody>
                  {hasil.baris.map((b) => (
                    <tr
                      key={`${b.caseId}-${b.policyNo}`}
                      className="inbox__baris"
                      onClick={() => {
                        if (b.policyNo !== '') setDetail(b)
                      }}
                    >
                      <td>{b.caseId}</td>
                      <td>
                        <button
                          type="button"
                          className="nbf-tautan"
                          disabled={b.policyNo === ''}
                          onClick={(e) => {
                            // Baris sudah membuka detail; cegah panggilan ganda.
                            e.stopPropagation()
                            setDetail(b)
                          }}
                        >
                          {b.policyNo}
                        </button>
                      </td>
                      <td>{b.reffNumber}</td>
                      <td>{b.insuredName}</td>
                      <td>{b.sourceOfBusiness}</td>
                      <td>{b.qqName}</td>
                      <td>{tanggalFollowing(b.startPeriod)}</td>
                      <td>{tanggalFollowing(b.endPeriod)}</td>
                      <td>
                        <button
                          type="button"
                          className="btn btn--sm nbf-tema__tombol"
                          disabled={b.policyNo === ''}
                          onClick={(e) => {
                            e.stopPropagation()
                            setDetail(b)
                          }}
                        >
                          {F.lihat}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              )}
            </table>
            {memuat && <Memuat />}
            {!memuat && hasil === null && <Kosong pesan={TEKS_FOLLOWING.isiFilter} />}
            {!memuat && hasil !== null && hasil.baris.length === 0 && <Kosong pesan={TEKS_FOLLOWING.kosong} />}
            <Gagal galat={galat} />
          </>
        )}
      </section>
    </div>
  )
}
