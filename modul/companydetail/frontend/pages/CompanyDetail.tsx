// Halaman Company Detail (perintah work owner 03/04-10-2026): daftar organisasi `CLIENT` (`FLAG` Org) bercari dan
// berhalaman, tombol Create, dan form organisasi (`FormCompany`) yang menggantikan daftar saat dibuka. Tampilan
// mengikuti tema modul (`companydetail.css`, akar `.companydetail__akar`). Nol hapus organisasi.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, ambilHak, ambilPilihan, type Detail, type Halaman, type Hak, type PilihanForm } from '../api'
import { jumlahHalaman, labelKode } from '../aturan'
import DialogCopyOld from '../components/DialogCopyOld'
import FormCompany from '../components/FormCompany'
import LihatCompany from '../components/LihatCompany'
import { CD } from '../labels'
import { NAMA_CD } from '../menu'

/** Tampilan halaman: daftar, View satu organisasi, atau form (`id` null = Create). */
type Tampilan = { jenis: 'daftar' } | { jenis: 'lihat'; id: string } | { jenis: 'form'; id: string | null }

export default function CompanyDetail() {
  // Menu View only (M_LOGIN_GO_MENU.HAK, 04-10-2026): tanpa Create dan Edit; View tetap.
  const bolehUbah = useBolehUbah(NAMA_CD)
  const [tampilan, setTampilan] = useState<Tampilan>({ jenis: 'daftar' })
  const [ketik, setKetik] = useState('')
  const [kueri, setKueri] = useState('')
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [pilihan, setPilihan] = useState<PilihanForm | null>(null)
  const [galatPilihan, setGalatPilihan] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [hak, setHak] = useState<Hak>({ copyOld: false })
  const [copyOld, setCopyOld] = useState(false)

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
  }, [muat, kueri, halaman])

  // Tombol Copy Old hanya bagi superadmin; backend menegakkannya juga (403).
  useEffect(() => {
    let hidup = true
    ambilHak().then(
      (h) => {
        if (hidup) setHak(h)
      },
      () => undefined,
    )
    return () => {
      hidup = false
    }
  }, [])

  useEffect(() => {
    let hidup = true
    ambilPilihan().then(
      (p) => {
        if (hidup) setPilihan(p)
      },
      (g: unknown) => {
        if (hidup) setGalatPilihan(g)
      },
    )
    return () => {
      hidup = false
    }
  }, [])

  const cari = () => {
    setPesan(null)
    setHalaman(1)
    setKueri(ketik.trim())
  }

  const buka = (id: string | null) => {
    setPesan(null)
    setTampilan({ jenis: 'form', id })
  }

  const lihat = (id: string) => {
    setPesan(null)
    setTampilan({ jenis: 'lihat', id })
  }

  const tersimpan = (d: Detail) => {
    setTampilan({ jenis: 'daftar' })
    setPesan(CD.tersimpan(d.idView))
    muat(kueri, halaman)
  }

  if (tampilan.jenis === 'lihat') {
    return (
      <section className="inbox companydetail__akar">
        {galatPilihan !== null && <Gagal galat={galatPilihan} />}
        {pilihan === null && galatPilihan === null && <Memuat pesan={CD.memuatPilihan} />}
        {pilihan !== null && (
          <LihatCompany
            key={tampilan.id}
            id={tampilan.id}
            pilihan={pilihan}
            onKembali={() => {
              setTampilan({ jenis: 'daftar' })
            }}
            onUbah={bolehUbah ? buka : undefined}
          />
        )}
      </section>
    )
  }

  if (tampilan.jenis === 'form') {
    return (
      <section className="inbox companydetail__akar">
        {galatPilihan !== null && <Gagal galat={galatPilihan} />}
        {pilihan === null && galatPilihan === null && <Memuat pesan={CD.memuatPilihan} />}
        {pilihan !== null && (
          <FormCompany
            key={tampilan.id ?? 'baru'}
            id={tampilan.id}
            pilihan={pilihan}
            onKembali={() => {
              setTampilan({ jenis: 'daftar' })
            }}
            onTersimpan={tersimpan}
          />
        )}
      </section>
    )
  }

  const total = data?.total ?? 0
  const jumlah = jumlahHalaman(total, data?.ukuran ?? 1)

  return (
    <section className="inbox companydetail__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{CD.judul}</h2>
      </header>

      <form
        className="toolbar"
        onSubmit={(e) => {
          e.preventDefault()
          cari()
        }}
      >
        <input
          className="field__input companydetail__cari"
          type="search"
          aria-label={CD.cari}
          placeholder={CD.cari}
          value={ketik}
          onChange={(e) => {
            setKetik(e.target.value)
          }}
        />
        <button type="submit" className="btn btn--ghost">
          {CD.tombolCari}
        </button>
        <span className="toolbar__spacer" />
        {hak.copyOld && bolehUbah && (
          <button
            type="button"
            className="btn btn--ghost"
            onClick={() => {
              setPesan(null)
              setCopyOld(true)
            }}
          >
            {CD.copyOld}
          </button>
        )}
        {bolehUbah && (
          <button
            type="button"
            className="btn btn--primary"
            onClick={() => {
              buka(null)
            }}
          >
            {CD.create}
          </button>
        )}
      </form>

      {copyOld && (
        <DialogCopyOld
          onTutup={(adaYangDisalin) => {
            setCopyOld(false)
            if (adaYangDisalin) muat(kueri, halaman)
          }}
        />
      )}

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {data === null && galat === null && <Memuat pesan={CD.memuat} />}
      {galat !== null && <Gagal galat={galat} />}

      {data !== null && (
        <>
          {data.daftar.length === 0 && <Kosong pesan={kueri === '' ? CD.kosong : CD.tidakCocok} />}
          {data.daftar.length > 0 && (
            <table className="inbox__tabel companydetail__tabel">
              <thead>
                <tr>
                  <th>{CD.kolomId}</th>
                  <th>{CD.organizationName}</th>
                  <th>{CD.title}</th>
                  <th>{CD.npwp}</th>
                  <th>{CD.country}</th>
                  <th>{CD.businessField}</th>
                  <th>{CD.parent}</th>
                  <th className="table__actions">{CD.kolomAksi}</th>
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((b) => (
                  <tr key={b.id} className="inbox__baris">
                    <td>{b.idView}</td>
                    <td>{b.nama}</td>
                    <td>{b.title}</td>
                    <td>{b.npwp}</td>
                    <td>{b.countryName}</td>
                    <td>{pilihan === null ? b.businessField : labelKode(pilihan.businessField, b.businessField)}</td>
                    <td>{b.parentName}</td>
                    <td className="table__actions">
                      <span className="companydetail__aksi">
                        <button
                          type="button"
                          className="btn btn--ghost"
                          onClick={() => {
                            lihat(b.id)
                          }}
                        >
                          {CD.lihat}
                        </button>
                        {bolehUbah && (
                          <button
                            type="button"
                            className="btn btn--ghost"
                            onClick={() => {
                              buka(b.id)
                            }}
                          >
                            {CD.ubah}
                          </button>
                        )}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          <div className="companydetail__halaman">
            <span className="muted">{CD.jumlah(total)}</span>
            <span className="toolbar__spacer" />
            <button
              type="button"
              className="btn btn--ghost"
              disabled={halaman <= 1}
              onClick={() => {
                setHalaman((h) => h - 1)
              }}
            >
              {CD.sebelumnya}
            </button>
            <span>{CD.halaman(halaman, jumlah)}</span>
            <button
              type="button"
              className="btn btn--ghost"
              disabled={halaman >= jumlah}
              onClick={() => {
                setHalaman((h) => h + 1)
              }}
            >
              {CD.berikutnya}
            </button>
          </div>
        </>
      )}
    </section>
  )
}
