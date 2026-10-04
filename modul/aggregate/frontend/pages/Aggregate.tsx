// Halaman Aggregate (perintah work owner 04-10-2026): chart RNM Value (USD) Ceding > Treaty Type > Coverage, daftar unggahan
// (`GridDasbordAgg`: Add Data, klik ganda = rincian, Delete dengan konfirmasi), dan layar unggah (`ShowAggregateList`)
// yang menggantikan daftar saat Add Data ditekan.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDaftar, hapus, type Halaman, type Kelompok, type Kunci } from '../api'
import { formatBulat, jumlahHalaman, kunciDari, kunciTeks } from '../aturan'
import DialogRincian from '../components/DialogRincian'
import ChartAggregate from '../components/ChartAggregate'
import UnggahAggregate from '../components/UnggahAggregate'
import { AG } from '../labels'

export default function Aggregate() {
  const [unggah, setUnggah] = useState(false)
  const [ketik, setKetik] = useState('')
  const [kueri, setKueri] = useState('')
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [rincian, setRincian] = useState<Kunci | null>(null)
  const [akanHapus, setAkanHapus] = useState<Kelompok | null>(null)
  const [menghapus, setMenghapus] = useState(false)
  const [galatHapus, setGalatHapus] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)

  const muat = useCallback((q: string, h: number) => {
    ambilDaftar(q, h).then(
      (d) => {
        setData(d)
        setGalat(null)
      },
      (g: unknown) => {
        setGalat(g)
      },
    )
  }, [])

  useEffect(() => {
    muat(kueri, halaman)
  }, [muat, kueri, halaman, segar])

  const jalankanHapus = () => {
    if (akanHapus === null || menghapus) return
    setMenghapus(true)
    setGalatHapus(null)
    hapus(kunciDari(akanHapus)).then(
      (r) => {
        setMenghapus(false)
        setAkanHapus(null)
        setPesan(AG.terhapus(r.dihapus))
        setSegar((n) => n + 1)
      },
      (g: unknown) => {
        setMenghapus(false)
        setGalatHapus(g)
      },
    )
  }

  if (unggah) {
    return (
      <section className="inbox aggregate__akar">
        <UnggahAggregate
          onTutup={() => {
            setUnggah(false)
            setSegar((n) => n + 1)
          }}
          onTersimpan={() => {
            setSegar((n) => n + 1)
          }}
        />
      </section>
    )
  }

  const total = data?.total ?? 0
  const jumlah = jumlahHalaman(total, data?.ukuran ?? 1)

  return (
    <section className="inbox aggregate__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{AG.judul}</h2>
      </header>

      <ChartAggregate segar={segar} />

      <form
        className="toolbar"
        onSubmit={(e) => {
          e.preventDefault()
          setPesan(null)
          setHalaman(1)
          setKueri(ketik.trim())
        }}
      >
        <input
          className="field__input aggregate__cari"
          type="search"
          aria-label={AG.cari}
          placeholder={AG.cari}
          value={ketik}
          onChange={(e) => {
            setKetik(e.target.value)
          }}
        />
        <button type="submit" className="btn btn--ghost">
          {AG.tombolCari}
        </button>
        <span className="toolbar__spacer" />
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setPesan(null)
            setUnggah(true)
          }}
        >
          {AG.addData}
        </button>
      </form>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {data === null && galat === null && <Memuat pesan={AG.memuat} />}
      {galat !== null && <Gagal galat={galat} />}

      {data !== null && (
        <>
          {data.daftar.length === 0 && <Kosong pesan={kueri === '' ? AG.kosong : AG.tidakCocok} />}
          {data.daftar.length > 0 && (
            <table className="inbox__tabel aggregate__tabel">
              <thead>
                <tr>
                  <th>{AG.tanggalInput}</th>
                  <th>{AG.cedingCode}</th>
                  <th>{AG.cedingName}</th>
                  <th>{AG.treatyType}</th>
                  <th>{AG.asAt}</th>
                  <th>{AG.uwYear}</th>
                  <th className="aggregate__angka">{AG.jumlahBaris}</th>
                  <th className="table__actions">{AG.aksi}</th>
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((k) => (
                  <tr
                    key={kunciTeks(k)}
                    className="inbox__baris aggregate__baris"
                    title={AG.petunjukBaris}
                    tabIndex={0}
                    onDoubleClick={() => {
                      setRincian(kunciDari(k))
                    }}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') setRincian(kunciDari(k))
                    }}
                  >
                    <td>{k.tanggalInput}</td>
                    <td>{k.cedingCode}</td>
                    <td>{k.cedingName}</td>
                    <td>{k.treatyType}</td>
                    <td>{k.asAt}</td>
                    <td>{k.uwYear}</td>
                    <td className="aggregate__angka">{formatBulat(k.jumlahBaris)}</td>
                    <td className="table__actions">
                      <span className="aggregate__aksi">
                        <button
                          type="button"
                          className="btn btn--ghost"
                          onClick={() => {
                            setPesan(null)
                            setGalatHapus(null)
                            setAkanHapus(k)
                          }}
                        >
                          {AG.hapus}
                        </button>
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          <div className="aggregate__halaman">
            <span className="muted">{AG.jumlah(total)}</span>
            <span className="toolbar__spacer" />
            <button
              type="button"
              className="btn btn--ghost"
              disabled={halaman <= 1}
              onClick={() => {
                setHalaman((h) => h - 1)
              }}
            >
              {AG.sebelumnya}
            </button>
            <span>{AG.halaman(halaman, jumlah)}</span>
            <button
              type="button"
              className="btn btn--ghost"
              disabled={halaman >= jumlah}
              onClick={() => {
                setHalaman((h) => h + 1)
              }}
            >
              {AG.berikutnya}
            </button>
          </div>
        </>
      )}

      {rincian !== null && (
        <DialogRincian
          kunci={rincian}
          onTutup={() => {
            setRincian(null)
          }}
        />
      )}

      {akanHapus !== null && (
        <Modal
          judul={AG.judulHapus}
          onTutup={() => {
            setAkanHapus(null)
          }}
          onKirim={jalankanHapus}
          labelBatal={AG.batal}
          aksi={
            <button type="submit" className="btn btn--primary" disabled={menghapus}>
              {menghapus ? AG.menghapus : AG.hapus}
            </button>
          }
        >
          {galatHapus !== null && <Gagal galat={galatHapus} />}
          <p>{AG.kalimatHapus(akanHapus.jumlahBaris)}</p>
          <p className="muted">
            {[akanHapus.tanggalInput, akanHapus.cedingName, akanHapus.treatyType, akanHapus.asAt, akanHapus.uwYear]
              .filter((s) => s !== '')
              .join(' · ')}
          </p>
        </Modal>
      )}
    </section>
  )
}
