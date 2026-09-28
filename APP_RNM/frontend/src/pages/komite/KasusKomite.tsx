// Kasus Komite — tiket 01 Komite Claim Life (baca saja).
//
// Tangga persetujuan kasus dan baris yang diputuskan. Layar keputusan
// `ShowTransfer` (`FlowAction/ViewTransferDtl.xml`) dibangun di tiket 02.

import { useCallback, useEffect, useState } from 'react'

import { KASUS_KOMITE, KOLOM_INBOX_KOMITE } from '../../assets/labels.komite'
import { Gagal, Memuat } from '../../components/ui/dasar'
import { ambilKasusKomite, type KasusKomite as Kasus } from '../../services/api'
import { selKomite, tingkatKomite } from './InboxKomite'

export default function KasusKomite({
  kasusID,
  onKembali,
}: {
  kasusID: string
  onKembali: () => void
}) {
  const [k, setK] = useState<Kasus | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(true)

  const muat = useCallback(async () => {
    setSibuk(true)
    setGalat(null)
    try {
      setK(await ambilKasusKomite(kasusID))
    } catch (e) {
      // ⚠️ 403 bila pelaku bukan anggota tangga — pesan server tampil apa adanya.
      setGalat(e)
      setK(null)
    } finally {
      setSibuk(false)
    }
  }, [kasusID])

  useEffect(() => {
    void muat()
  }, [muat])

  return (
    <section className="komite-kasus">
      <button type="button" onClick={onKembali}>
        {KASUS_KOMITE.kembali}
      </button>
      <h2>
        {KASUS_KOMITE.judul} {kasusID}
      </h2>
      {sibuk && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {k !== null && (
        <>
          <dl className="komite-kasus__kepala">
            <dt>{KOLOM_INBOX_KOMITE.nomorKlaim}</dt>
            <dd>{selKomite(k.kasus.nomorKlaim)}</dd>
            <dt>{KOLOM_INBOX_KOMITE.tingkat}</dt>
            <dd>{tingkatKomite(k.kasus)}</dd>
            <dt>{KOLOM_INBOX_KOMITE.nilaiKlaim}</dt>
            <dd>
              {selKomite(k.kasus.nilaiKlaim)} {k.kasus.mataUang}
            </dd>
            <dt>{KOLOM_INBOX_KOMITE.statusBaris}</dt>
            <dd>{k.kasus.statusBaris}</dd>
          </dl>
          <p role="status">
            {k.giliranSaya ? KASUS_KOMITE.giliranAnda : KASUS_KOMITE.bukanGiliran}
          </p>
          <h3>{KASUS_KOMITE.tangga}</h3>
          <table className="komite-kasus__tangga">
            <thead>
              <tr>
                <th>{KASUS_KOMITE.urut}</th>
                <th>{KASUS_KOMITE.jabatan}</th>
                <th>{KASUS_KOMITE.approval}</th>
                <th>{KASUS_KOMITE.komentar}</th>
                <th>{KASUS_KOMITE.tanggal}</th>
              </tr>
            </thead>
            <tbody>
              {k.tangga.map((a) => (
                <tr key={a.urut} className={a.saya ? 'komite-kasus__saya' : undefined}>
                  <td>{a.urut}</td>
                  <td>{selKomite(a.jabatan)}</td>
                  <td>{a.kataApproval}</td>
                  <td>{selKomite(a.komentar)}</td>
                  <td>{selKomite(a.tglApprove)}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <p>{KASUS_KOMITE.keputusanMenyusul}</p>
        </>
      )}
    </section>
  )
}
