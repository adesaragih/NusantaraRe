// Popup Master ID - padanan `ChooseMasterID`: kotak cari ceding (Enter = cari, `GetMasterIDAgg_Act`), grid treaty
// dengan kotak centang, Submit = Master ID terpilih DIGANTI baris yang dicentang (`SetMasterID`: hapus lalu tambah
// yang dicentang, urutan grid).

import { useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { cariMasterTreaty, type MasterTreaty } from '../api'
import { formatAngka } from '../aturan'
import { AG } from '../labels'

export default function DialogMasterID({
  terpilih,
  onPilih,
  onTutup,
}: {
  terpilih: MasterTreaty[]
  onPilih: (m: MasterTreaty[]) => void
  onTutup: () => void
}) {
  const [kata, setKata] = useState('')
  const [hasil, setHasil] = useState<MasterTreaty[]>(terpilih)
  const [centang, setCentang] = useState<ReadonlySet<string>>(new Set(terpilih.map((t) => t.id)))
  const [mencari, setMencari] = useState(false)
  const [sudahCari, setSudahCari] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)

  const cari = () => {
    if (mencari) return
    setMencari(true)
    setGalat(null)
    cariMasterTreaty(kata.trim()).then(
      (r) => {
        setMencari(false)
        setSudahCari(true)
        setHasil(r.daftar)
      },
      (g: unknown) => {
        setMencari(false)
        setGalat(g)
      },
    )
  }

  const ubahCentang = (id: string, ya: boolean) => {
    setCentang((c) => {
      const b = new Set(c)
      if (ya) b.add(id)
      else b.delete(id)
      return b
    })
  }

  return (
    <Modal
      judul={AG.judulMaster}
      onTutup={onTutup}
      onKirim={() => {
        onPilih(hasil.filter((t) => centang.has(t.id)))
      }}
      labelBatal={AG.batal}
      penuh
      aksi={
        <button type="submit" className="btn btn--primary">
          {AG.submit}
        </button>
      }
    >
      <div className="aggregate__master">
        <div className="aggregate__master-alat">
          <input
            className="field__input aggregate__cari"
            type="search"
            aria-label={AG.cariMaster}
            placeholder={AG.cariMaster}
            value={kata}
            onChange={(e) => {
              setKata(e.target.value)
            }}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                e.preventDefault()
                cari()
              }
            }}
          />
          <button type="button" className="btn btn--ghost" onClick={cari} disabled={mencari}>
            {AG.tombolCari}
          </button>
          <span className="muted">{AG.dipilih(hasil.filter((t) => centang.has(t.id)).length)}</span>
        </div>
        {galat !== null && <Gagal galat={galat} />}
        {mencari && <Memuat pesan={AG.mencari} />}
        {!mencari && sudahCari && hasil.length === 0 && <p className="muted">{AG.masterTidakAda}</p>}
        {hasil.length > 0 && (
          <div className="aggregate__gulir aggregate__master-tabel">
            <table className="inbox__tabel aggregate__grid">
              <thead>
                <tr>
                  <th>{AG.pilih}</th>
                  <th>{AG.masterId}</th>
                  <th>{AG.treatyContractName}</th>
                  <th>{AG.reinsuranceType}</th>
                  <th>{AG.ceding}</th>
                  <th>{AG.sob}</th>
                  <th>{AG.treatyGroup}</th>
                  <th className="aggregate__angka">{AG.rnmShare}</th>
                  <th>{AG.treatyYear}</th>
                </tr>
              </thead>
              <tbody>
                {hasil.map((t) => (
                  <tr key={t.id} className="inbox__baris">
                    <td>
                      <input
                        type="checkbox"
                        aria-label={`${AG.pilih} ${t.treatyId}`}
                        checked={centang.has(t.id)}
                        onChange={(e) => {
                          ubahCentang(t.id, e.target.checked)
                        }}
                      />
                    </td>
                    <td>{t.treatyId}</td>
                    <td>{t.treatyContractName}</td>
                    <td>{t.proportionType}</td>
                    <td>{t.ceding}</td>
                    <td>{t.sob}</td>
                    <td>{t.treatyGroup}</td>
                    <td className="aggregate__angka">{formatAngka(t.rnmShare)}</td>
                    <td>{t.treatyYear}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </Modal>
  )
}
