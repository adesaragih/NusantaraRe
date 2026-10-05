// Form organisasi Company Detail - padanan layar Pega SFAGIS Company Detail (screenshot work owner 03-10-2026):
// bagian Company Detail (NPWP, Parent organization, COUNTRY*, Title, Organization Name*, Business Field*, Note),
// grid PIC (Name* - dropdown akun login aktif, Position*, Gender, Email, Date of birth, Phone number), dan grid
// Address (Type, Address, Phone and Fax). ORG ID hanya ditampilkan: backend yang membentuknya saat Create.

import { useCallback, useEffect, useState } from 'react'

import { Area, Field, Gagal, Memuat, Modal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { PilihSaring, type OpsiSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import { dariInputTanggal, keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'
import {
  ambil,
  cariInduk,
  periksaNama,
  tambah,
  ubah,
  type Detail,
  type Isian,
  type PeriksaNama,
  type Pilihan,
  type PilihanForm,
} from '../api'
import {
  alamatKosong,
  buangDi,
  gantiDi,
  isianDari,
  isianKosong,
  jabatanAkun,
  KODE_LAIN,
  kodeAreaLain,
  opsiAkun,
  opsiDari,
  opsiKodeArea,
  opsiNegara,
  opsiTitle,
  periksa,
  picKosong,
  telfaxKosong,
  galatTitle,
} from '../aturan'
import { CD } from '../labels'

const tetap = () => undefined

function PilihKecil({
  label,
  value,
  opsi,
  onChange,
}: {
  label: string
  value: string
  opsi: { value: string; label: string }[]
  onChange: (v: string) => void
}) {
  return (
    <select
      className="field__input"
      aria-label={label}
      value={value}
      onChange={(e) => {
        onChange(e.target.value)
      }}
    >
      <option value="">{CD.pilih}</option>
      {opsi.map((o) => (
        <option key={o.value} value={o.value}>
          {o.label}
        </option>
      ))}
    </select>
  )
}

/**
 * Kode area satu nomor: dropdown kodehp + "Others" yang membuka isian kode sendiri (perintah work owner 05-10-2026).
 * Kode terpasang yang tidak ada di daftar langsung tampil sebagai Others.
 */
function KodeArea({ kode, value, onChange }: { kode: Pilihan[]; value: string; onChange: (v: string) => void }) {
  const [lain, setLain] = useState(() => kodeAreaLain(kode, value))
  const modeLain = lain || kodeAreaLain(kode, value)
  return (
    <div className="companydetail__kode">
      <PilihKecil
        label={CD.kodeArea}
        value={modeLain ? KODE_LAIN : value}
        opsi={opsiKodeArea(kode, modeLain ? '' : value)}
        onChange={(v) => {
          setLain(v === KODE_LAIN)
          onChange(v === KODE_LAIN ? '' : v)
        }}
      />
      {modeLain && (
        <input
          className="field__input"
          aria-label={CD.isiKodeArea}
          placeholder={CD.isiKodeArea}
          inputMode="tel"
          maxLength={10}
          value={value}
          onChange={(e) => {
            onChange(e.target.value)
          }}
        />
      )}
    </div>
  )
}

export default function FormCompany({
  id,
  pilihan,
  onKembali,
  onTersimpan,
}: {
  /** `null` = Create. */
  id: string | null
  pilihan: PilihanForm
  onKembali: () => void
  onTersimpan: (d: Detail) => void
}) {
  const baru = id === null
  const [detail, setDetail] = useState<Detail | null>(null)
  const [galatMuat, setGalatMuat] = useState<unknown>(null)
  const [isi, setIsi] = useState<Isian>(isianKosong)
  const [namaInduk, setNamaInduk] = useState('')
  const [opsiInduk, setOpsiInduk] = useState<OpsiSaring[]>([])
  const [memuatInduk, setMemuatInduk] = useState(false)
  const [galatLokal, setGalatLokal] = useState<string | null>(null)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  /** Temuan pemeriksaan nama sebelum Create; terisi = popup Continue save / No terbuka. */
  const [konfirmasi, setKonfirmasi] = useState<PeriksaNama | null>(null)
  const [memeriksa, setMemeriksa] = useState(false)

  useEffect(() => {
    if (id === null) return
    let hidup = true
    ambil(id).then(
      (d) => {
        if (!hidup) return
        setDetail(d)
        setIsi(isianDari(d))
        setNamaInduk(d.parentName)
      },
      (g: unknown) => {
        if (hidup) setGalatMuat(g)
      },
    )
    return () => {
      hidup = false
    }
  }, [id])

  // Parent organization: dropdown bersaring (`PilihSaring`) - 18 ribuan organisasi dicari di server, paling banyak
  // 20 per jawaban; kata kosong = 20 nama pertama.
  const cariOpsiInduk = useCallback(
    (kata: string) => {
      setMemuatInduk(true)
      cariInduk(kata, id ?? '').then(
        (h) => {
          setOpsiInduk(h.daftar.map((b) => ({ value: b.id, label: b.nama, keterangan: b.idView })))
          setMemuatInduk(false)
        },
        () => {
          setOpsiInduk([])
          setMemuatInduk(false)
        },
      )
    },
    [id],
  )

  const ubahIsi = (sebagian: Partial<Isian>) => {
    setIsi((x) => ({ ...x, ...sebagian }))
  }

  // Create: title di nama DITOLAK (periksa); lalu organisasi bernama sama/mirip diperiksa - ada temuan = popup
  // Continue save / No (perintah work owner 04-10-2026). Ubah langsung disimpan.
  const kirim = () => {
    if (sibuk || memeriksa) return
    const awal = periksa(isi, pilihan.title, detail?.nama)
    setGalatLokal(awal)
    setGalatSimpan(null)
    if (awal !== null) return
    if (!baru) {
      simpan()
      return
    }
    setMemeriksa(true)
    periksaNama(isi.nama, '').then(
      (h) => {
        setMemeriksa(false)
        if (h.serupa.length > 0) {
          setKonfirmasi(h)
          return
        }
        simpan()
      },
      (g: unknown) => {
        setMemeriksa(false)
        setGalatSimpan(g)
      },
    )
  }

  const simpan = () => {
    setKonfirmasi(null)
    setSibuk(true)
    const janji = id === null ? tambah(isi) : ubah(id, isi)
    janji.then(
      (d) => {
        setSibuk(false)
        onTersimpan(d)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalatSimpan(g)
      },
    )
  }

  const galatNamaTitle = galatTitle(isi.nama, pilihan.title, detail?.nama)

  return (
    <>
      <header className="inbox__kepala">
        <button type="button" className="btn btn--ghost companydetail__kembali" onClick={onKembali}>
          {CD.kembali}
        </button>
        <h2 className="inbox__judul">{baru ? CD.judulBaru : CD.judulUbah(detail?.idView ?? '')}</h2>
      </header>

      {galatMuat !== null && <Gagal galat={galatMuat} />}
      {!baru && detail === null && galatMuat === null && <Memuat pesan={CD.memuatDetail} />}

      {(baru || detail !== null) && (
        <form
          className="companydetail__form"
          onSubmit={(e) => {
            e.preventDefault()
            kirim()
          }}
        >
          {galatLokal !== null && (
            <div className="alert alert--error" role="alert">
              {galatLokal}
            </div>
          )}
          {galatSimpan !== null && <Gagal galat={galatSimpan} />}

          <fieldset className="companydetail__bagian">
            <legend className="companydetail__judul-bagian">{CD.companyDetail}</legend>
            <div className="form-grid">
              <Field label={CD.kolomId} value={detail?.idView ?? ''} placeholder={CD.idOtomatis} onChange={tetap} readOnly />
              <Field
                label={CD.npwp}
                value={isi.npwp}
                onChange={(v) => {
                  ubahIsi({ npwp: v })
                }}
              />
              <div className="companydetail__induk">
                <PilihSaring
                  label={CD.parent}
                  value={isi.parentId}
                  teksTerpilih={namaInduk}
                  opsi={opsiInduk}
                  memuat={memuatInduk}
                  onCari={cariOpsiInduk}
                  onPilih={(o) => {
                    ubahIsi({ parentId: o.value })
                    setNamaInduk(o.label)
                  }}
                />
                {(isi.parentId !== '' || namaInduk !== '') && (
                  <button
                    type="button"
                    className="btn btn--ghost companydetail__induk-hapus"
                    onClick={() => {
                      ubahIsi({ parentId: '' })
                      setNamaInduk('')
                    }}
                  >
                    {CD.hapusParent}
                  </button>
                )}
              </div>
              <Pilih
                label={CD.country}
                value={isi.country}
                onChange={(v) => {
                  ubahIsi({ country: v })
                }}
                opsi={opsiNegara(pilihan.negara, isi.country, detail?.countryName ?? '')}
                kosong={CD.pilih}
                required
              />
              <Pilih
                label={CD.title}
                value={isi.title}
                onChange={(v) => {
                  ubahIsi({ title: v })
                }}
                opsi={opsiTitle(pilihan.title, isi.title)}
                kosong={CD.pilih}
              />
              <div>
                {/* Edit: Organization Name tidak boleh diubah (perintah work owner 05-10-2026); backend menolaknya juga. */}
                <Field
                  label={CD.organizationName}
                  value={isi.nama}
                  onChange={(v) => {
                    // Huruf besar saat diketik (perintah work owner 04-10-2026); backend menegakkannya juga.
                    ubahIsi({ nama: v.toUpperCase() })
                  }}
                  readOnly={!baru}
                  required
                />
                {galatNamaTitle !== null && (
                  <p className="field__error" role="alert">
                    {galatNamaTitle}
                  </p>
                )}
              </div>
              <Pilih
                label={CD.businessField}
                value={isi.businessField}
                onChange={(v) => {
                  ubahIsi({ businessField: v })
                }}
                opsi={opsiDari(pilihan.businessField, isi.businessField)}
                kosong={CD.pilih}
                required
              />
              <Area
                label={CD.note}
                value={isi.note}
                onChange={(v) => {
                  ubahIsi({ note: v })
                }}
                baris={3}
              />
            </div>
            {detail !== null && detail.createdBy !== '' && (
              <p className="muted companydetail__catatan">{CD.jejakDibuat(detail.createdBy, detail.createdAt)}</p>
            )}
            {detail !== null && detail.updatedBy !== '' && (
              <p className="muted companydetail__catatan">{CD.jejakDiubah(detail.updatedBy, detail.updatedAt)}</p>
            )}
          </fieldset>

          <fieldset className="companydetail__bagian">
            <legend className="companydetail__judul-bagian">{CD.pic}</legend>
            {isi.pic.length === 0 && <p className="muted">{CD.kosongPIC}</p>}
            {isi.pic.length > 0 && (
              <div className="companydetail__gulir">
                <table className="companydetail__grid">
                  <thead>
                    <tr>
                      <th>
                        {CD.name}
                        <span className="field__req">*</span>
                      </th>
                      <th>
                        {CD.position}
                        <span className="field__req">*</span>
                      </th>
                      <th>{CD.gender}</th>
                      <th>{CD.email}</th>
                      <th>{CD.dateOfBirth}</th>
                      <th>{CD.phoneNumber}</th>
                      <th />
                    </tr>
                  </thead>
                  <tbody>
                    {isi.pic.map((p, i) => {
                      const ganti = (sebagian: Partial<typeof p>) => {
                        ubahIsi({ pic: gantiDi(isi.pic, i, { ...p, ...sebagian }) })
                      }
                      return (
                        <tr key={p.userIdentifier === '' ? `baru-${i}` : p.userIdentifier}>
                          <td>
                            {/* Name dari akun login aktif M_LOGIN_GO (perintah work owner 05-10-2026). */}
                            <PilihKecil
                              label={CD.name}
                              value={p.nama}
                              opsi={opsiAkun(pilihan.akun, p.nama)}
                              onChange={(v) => {
                                // Position = JOB_POSITION akun yang dipilih; nama PIC lama (bukan akun) membawa
                                // position-nya sendiri.
                                const jabatan = jabatanAkun(pilihan.akun, v)
                                ganti(jabatan === undefined ? { nama: v } : { nama: v, position: jabatan })
                              }}
                            />
                          </td>
                          <td>
                            {/* Position dari akun login yang dipilih (perintah work owner 05-10-2026) - tidak diketik. */}
                            <input
                              className="field__input field__input--readonly"
                              aria-label={CD.position}
                              value={p.position}
                              readOnly
                            />
                          </td>
                          <td>
                            <PilihKecil
                              label={CD.gender}
                              value={p.gender}
                              opsi={opsiDari(pilihan.gender, p.gender)}
                              onChange={(v) => {
                                ganti({ gender: v })
                              }}
                            />
                          </td>
                          <td>
                            <input
                              className="field__input"
                              type="email"
                              aria-label={CD.email}
                              value={p.email}
                              onChange={(e) => {
                                ganti({ email: e.target.value })
                              }}
                            />
                          </td>
                          <td>
                            <input
                              className="field__input"
                              type="date"
                              aria-label={CD.dateOfBirth}
                              value={keInputTanggal(p.dateOfBirth)}
                              onChange={(e) => {
                                ganti({ dateOfBirth: dariInputTanggal(e.target.value) })
                              }}
                            />
                          </td>
                          <td>
                            <input
                              className="field__input"
                              aria-label={CD.phoneNumber}
                              value={p.phone}
                              onChange={(e) => {
                                ganti({ phone: e.target.value })
                              }}
                            />
                          </td>
                          <td>
                            <button
                              type="button"
                              className="btn btn--ghost"
                              onClick={() => {
                                ubahIsi({ pic: buangDi(isi.pic, i) })
                              }}
                            >
                              {CD.hapus}
                            </button>
                          </td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              </div>
            )}
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                ubahIsi({ pic: [...isi.pic, picKosong()] })
              }}
            >
              {CD.tambahPIC}
            </button>
          </fieldset>

          <fieldset className="companydetail__bagian">
            <legend className="companydetail__judul-bagian">{CD.address}</legend>
            {isi.alamat.length === 0 && <p className="muted">{CD.kosongAlamat}</p>}
            {isi.alamat.map((a, i) => {
              const ganti = (sebagian: Partial<typeof a>) => {
                ubahIsi({ alamat: gantiDi(isi.alamat, i, { ...a, ...sebagian }) })
              }
              return (
                <div key={a.asal === '' ? `baru-${i}` : a.asal} className="companydetail__alamat">
                  <div className="companydetail__alamat-kepala">
                    <label className="companydetail__isian">
                      <span className="field__label">
                        {CD.type}
                        <span className="field__req">*</span>
                      </span>
                      <PilihKecil
                        label={CD.type}
                        value={a.type}
                        opsi={opsiDari(pilihan.addressType, a.type)}
                        onChange={(v) => {
                          ganti({ type: v })
                        }}
                      />
                    </label>
                    <label className="companydetail__isian companydetail__isian--lebar">
                      <span className="field__label">
                        {CD.address}
                        <span className="field__req">*</span>
                      </span>
                      <textarea
                        className="field__input"
                        rows={2}
                        value={a.address}
                        onChange={(e) => {
                          ganti({ address: e.target.value })
                        }}
                      />
                    </label>
                    <button
                      type="button"
                      className="btn btn--ghost"
                      onClick={() => {
                        ubahIsi({ alamat: buangDi(isi.alamat, i) })
                      }}
                    >
                      {CD.hapus}
                    </button>
                  </div>
                  <div className="companydetail__telfax">
                    <span className="field__label">{CD.phoneAndFax}</span>
                    {a.telfax.length === 0 && <p className="muted">{CD.kosongNomor}</p>}
                    {a.telfax.map((t, j) => {
                      const gantiNomor = (sebagian: Partial<typeof t>) => {
                        ganti({ telfax: gantiDi(a.telfax, j, { ...t, ...sebagian }) })
                      }
                      return (
                        <div key={j} className="companydetail__nomor">
                          <PilihKecil
                            label={CD.type}
                            value={t.type}
                            opsi={opsiDari(pilihan.telfax, t.type)}
                            onChange={(v) => {
                              gantiNomor({ type: v })
                            }}
                          />
                          <KodeArea
                            kode={pilihan.kodeArea}
                            value={t.code}
                            onChange={(v) => {
                              gantiNomor({ code: v })
                            }}
                          />
                          <input
                            className="field__input"
                            aria-label={CD.nomor}
                            placeholder={CD.nomor}
                            value={t.no}
                            onChange={(e) => {
                              gantiNomor({ no: e.target.value })
                            }}
                          />
                          <button
                            type="button"
                            className="btn btn--ghost"
                            onClick={() => {
                              ganti({ telfax: buangDi(a.telfax, j) })
                            }}
                          >
                            {CD.hapus}
                          </button>
                        </div>
                      )
                    })}
                    <button
                      type="button"
                      className="btn btn--ghost"
                      onClick={() => {
                        ganti({ telfax: [...a.telfax, telfaxKosong()] })
                      }}
                    >
                      {CD.tambahNomor}
                    </button>
                  </div>
                </div>
              )
            })}
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                ubahIsi({ alamat: [...isi.alamat, alamatKosong()] })
              }}
            >
              {CD.tambahAlamat}
            </button>
          </fieldset>

          <div className="companydetail__aksi-form">
            <button type="button" className="btn btn--ghost" onClick={onKembali} disabled={sibuk}>
              {CD.batal}
            </button>
            <button type="submit" className="btn btn--primary" disabled={sibuk || memeriksa}>
              {sibuk ? CD.menyimpan : memeriksa ? CD.memeriksa : baru ? CD.create : CD.simpan}
            </button>
          </div>
        </form>
      )}

      {konfirmasi !== null && (
        <Modal
          judul={CD.judulCek}
          onTutup={() => {
            setKonfirmasi(null)
          }}
          onKirim={simpan}
          labelBatal={CD.tidak}
          lebar
          aksi={
            <button type="submit" className="btn btn--primary">
              {CD.lanjutSimpan}
            </button>
          }
        >
          {konfirmasi.serupa.length > 0 && (
            <>
              <p>{CD.kalimatSerupa(isi.nama.trim())}</p>
              <div className="companydetail__gulir">
                <table className="inbox__tabel companydetail__tabel">
                  <thead>
                    <tr>
                      <th>{CD.kolomId}</th>
                      <th>{CD.organizationName}</th>
                      <th>{CD.title}</th>
                      <th>{CD.country}</th>
                      <th>{CD.kolomKemiripan}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {konfirmasi.serupa.map((s) => (
                      <tr key={s.id}>
                        <td>{s.idView}</td>
                        <td>{s.nama}</td>
                        <td>{s.title}</td>
                        <td>{s.countryName}</td>
                        <td>
                          <span className={s.jenis === 'sama' ? 'companydetail__tanda companydetail__tanda--sama' : 'companydetail__tanda'}>
                            {s.jenis === 'sama' ? CD.jenisSama : CD.jenisMirip}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </>
          )}
        </Modal>
      )}
    </>
  )
}
