// Tampilan View organisasi Company Detail (perintah work owner 04-10-2026: "tambah view sama edit"): isi yang sama
// dengan form - Company Detail, PIC, Address - HANYA DIBACA, dengan tombol Edit ke form.

import { useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambil, type Detail, type PilihanForm } from '../api'
import { labelKode } from '../aturan'
import { CD } from '../labels'

function Ruas({ label, nilai }: { label: string; nilai: string }) {
  return (
    <div className="companydetail__ruas">
      <span className="field__label">{label}</span>
      <span className="companydetail__nilai">{nilai === '' ? CD.kosongNilai : nilai}</span>
    </div>
  )
}

export default function LihatCompany({
  id,
  pilihan,
  onKembali,
  onUbah,
}: {
  id: string
  pilihan: PilihanForm
  onKembali: () => void
  /** Tidak ada = menu View only (M_LOGIN_GO_MENU.HAK, 04-10-2026): tanpa tombol Edit. */
  onUbah?: (id: string) => void
}) {
  const [detail, setDetail] = useState<Detail | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    ambil(id).then(
      (d) => {
        if (hidup) setDetail(d)
      },
      (g: unknown) => {
        if (hidup) setGalat(g)
      },
    )
    return () => {
      hidup = false
    }
  }, [id])

  return (
    <>
      <header className="inbox__kepala">
        <button type="button" className="btn btn--ghost companydetail__kembali" onClick={onKembali}>
          {CD.kembali}
        </button>
        <h2 className="inbox__judul">{CD.judulLihat(detail?.idView ?? '')}</h2>
        <span className="toolbar__spacer" />
        {detail !== null && onUbah !== undefined && (
          <button
            type="button"
            className="btn btn--primary"
            onClick={() => {
              onUbah(id)
            }}
          >
            {CD.ubah}
          </button>
        )}
      </header>

      {galat !== null && <Gagal galat={galat} />}
      {detail === null && galat === null && <Memuat pesan={CD.memuatDetail} />}

      {detail !== null && (
        <div className="companydetail__form">
          <section className="companydetail__bagian">
            <h3 className="companydetail__judul-bagian">{CD.companyDetail}</h3>
            <div className="companydetail__lihat">
              <Ruas label={CD.kolomId} nilai={detail.idView} />
              <Ruas label={CD.npwp} nilai={detail.npwp} />
              <Ruas label={CD.parent} nilai={detail.parentName} />
              <Ruas label={CD.country} nilai={detail.countryName === '' ? detail.country : detail.countryName} />
              <Ruas label={CD.title} nilai={detail.title} />
              <Ruas label={CD.organizationName} nilai={detail.nama} />
              <Ruas label={CD.businessField} nilai={labelKode(pilihan.businessField, detail.businessField)} />
              <Ruas label={CD.note} nilai={detail.note} />
            </div>
            {detail.createdBy !== '' && (
              <p className="muted companydetail__catatan">{CD.jejakDibuat(detail.createdBy, detail.createdAt)}</p>
            )}
            {detail.updatedBy !== '' && (
              <p className="muted companydetail__catatan">{CD.jejakDiubah(detail.updatedBy, detail.updatedAt)}</p>
            )}
          </section>

          <section className="companydetail__bagian">
            <h3 className="companydetail__judul-bagian">{CD.pic}</h3>
            {detail.pic.length === 0 && <p className="muted">{CD.kosongPIC}</p>}
            {detail.pic.length > 0 && (
              <div className="companydetail__gulir">
                <table className="inbox__tabel companydetail__tabel">
                  <thead>
                    <tr>
                      <th>{CD.name}</th>
                      <th>{CD.position}</th>
                      <th>{CD.gender}</th>
                      <th>{CD.email}</th>
                      <th>{CD.dateOfBirth}</th>
                      <th>{CD.phoneNumber}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {detail.pic.map((p) => (
                      <tr key={p.userIdentifier}>
                        <td>{p.nama}</td>
                        <td>{p.position}</td>
                        <td>{labelKode(pilihan.gender, p.gender)}</td>
                        <td>{p.email}</td>
                        <td>{p.dateOfBirth}</td>
                        <td>{p.phone}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </section>

          <section className="companydetail__bagian">
            <h3 className="companydetail__judul-bagian">{CD.address}</h3>
            {detail.alamat.length === 0 && <p className="muted">{CD.kosongAlamat}</p>}
            {detail.alamat.length > 0 && (
              <div className="companydetail__gulir">
                <table className="inbox__tabel companydetail__tabel">
                  <thead>
                    <tr>
                      <th>{CD.type}</th>
                      <th>{CD.address}</th>
                      <th>{CD.phoneAndFax}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {detail.alamat.map((a) => (
                      <tr key={a.asal}>
                        <td>{labelKode(pilihan.addressType, a.type)}</td>
                        <td className="companydetail__teks-panjang">{a.address}</td>
                        <td>
                          {a.telfax.length === 0 && <span className="muted">{CD.kosongNilai}</span>}
                          {a.telfax.map((t, j) => (
                            <span key={j} className="companydetail__kecil">
                              {`${labelKode(pilihan.telfax, t.type)}: ${t.code === '' ? '' : `${t.code} `}${t.no}`}
                            </span>
                          ))}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </section>
        </div>
      )}
    </>
  )
}
