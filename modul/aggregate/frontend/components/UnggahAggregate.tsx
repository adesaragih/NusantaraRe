// Layar unggah - padanan `ShowAggregateList`: grid Master ID terpilih (popup `ChooseMasterID`), Template (unduh
// kepala CSV), Upload CSV (versi Go `UploadCSVAggregate_Act`), grid pratinjau TempCSV HANYA DIBACA (work owner
// 04-10-2026), Save (`SaveAggregate_Act`), dan Close. Sesudah Save berhasil grid dikosongkan seperti Pega
// (Page-New TempCSV).

import { useRef, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { kodeStatusGalat, pesanGalat } from '../../../../inti/frontend/klien'
import { pratinjau, simpan, type Baris, type MasterTreaty } from '../api'
import { unduhTemplat } from '../../../../inti/frontend/templat/api'
import { barisPesan, formatAngka, KODE_TEMPLAT, NAMA_TEMPLATE } from '../aturan'
import { AG } from '../labels'
import DialogMasterID from './DialogMasterID'
import GridAggregate from './GridAggregate'

export default function UnggahAggregate({ onTutup, onTersimpan }: { onTutup: () => void; onTersimpan: () => void }) {
  const [master, setMaster] = useState<MasterTreaty[]>([])
  const [baris, setBaris] = useState<Baris[] | null>(null)
  const [dialog, setDialog] = useState(false)
  const [sibuk, setSibuk] = useState<'baca' | 'simpan' | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [tolak, setTolak] = useState<string[]>([])
  const [pesan, setPesan] = useState<string | null>(null)
  const berkas = useRef<HTMLInputElement>(null)

  const gagal = (g: unknown) => {
    const teks = kodeStatusGalat(g) === 422 ? pesanGalat(g) : undefined
    if (teks !== undefined) setTolak(barisPesan(teks))
    else setGalat(g)
  }

  const bersihkanPesan = () => {
    setGalat(null)
    setTolak([])
    setPesan(null)
  }

  const unggah = (f: File) => {
    bersihkanPesan()
    setSibuk('baca')
    f.text()
      .then((isi) => pratinjau(isi, master.map((m) => m.id)))
      .then(
        (p) => {
          setSibuk(null)
          setBaris(p.baris)
          setMaster(p.masterTreaty)
        },
        (g: unknown) => {
          setSibuk(null)
          gagal(g)
        },
      )
  }

  const simpanSemua = () => {
    if (baris === null || sibuk !== null) return
    bersihkanPesan()
    setSibuk('simpan')
    simpan(baris).then(
      (h) => {
        setSibuk(null)
        setPesan(h.pesan)
        setBaris(null)
        onTersimpan()
      },
      (g: unknown) => {
        setSibuk(null)
        gagal(g)
      },
    )
  }

  return (
    <>
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{AG.judul}</h2>
      </header>

      <section className="aggregate__kartu">
        <div className="aggregate__aksi-unggah">
          <button
            type="button"
            className="btn btn--ghost"
            onClick={() => {
              setDialog(true)
            }}
          >
            {AG.masterId}
          </button>
          <button
            type="button"
            className="btn btn--ghost"
            onClick={() => {
              unduhTemplat(KODE_TEMPLAT, NAMA_TEMPLATE).catch((g: unknown) => {
                setGalat(g)
              })
            }}
          >
            {AG.template}
          </button>
          <button
            type="button"
            className="btn btn--ghost"
            disabled={sibuk !== null}
            onClick={() => {
              berkas.current?.click()
            }}
          >
            {AG.uploadCsv}
          </button>
          <input
            ref={berkas}
            className="aggregate__berkas"
            type="file"
            accept=".csv,text/csv"
            aria-label={AG.uploadCsv}
            onChange={(e) => {
              const f = e.target.files?.[0]
              e.target.value = ''
              if (f !== undefined) unggah(f)
            }}
          />
          <span className="toolbar__spacer" />
          <button type="button" className="btn btn--primary" disabled={baris === null || sibuk !== null} onClick={simpanSemua}>
            {sibuk === 'simpan' ? AG.menyimpan : AG.save}
          </button>
          <button type="button" className="btn btn--ghost" onClick={onTutup} disabled={sibuk !== null}>
            {AG.close}
          </button>
        </div>

        {master.length === 0 && <p className="muted">{AG.masterKosong}</p>}
        {master.length > 0 && (
          <div className="aggregate__gulir">
            <table className="inbox__tabel aggregate__grid">
              <thead>
                <tr>
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
                {master.map((t) => (
                  <tr key={t.id} className="inbox__baris">
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
      </section>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {tolak.length > 0 && (
        <div className="alert alert--error" role="alert">
          <ul className="aggregate__tolak">
            {tolak.map((t) => (
              <li key={t}>{t}</li>
            ))}
          </ul>
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {sibuk === 'baca' && <Memuat pesan={AG.membaca} />}

      <section className="aggregate__kartu">
        {baris === null && sibuk !== 'baca' && <p className="muted">{AG.pratinjauKosong}</p>}
        {baris !== null && <GridAggregate baris={baris} />}
      </section>

      {dialog && (
        <DialogMasterID
          terpilih={master}
          onPilih={(m) => {
            setMaster(m)
            setDialog(false)
          }}
          onTutup={() => {
            setDialog(false)
          }}
        />
      )}
    </>
  )
}
