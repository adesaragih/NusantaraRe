// Form tambah/ubah user — Kelola User (keputusan work owner 01-10-2026).
//
// Dua tab:
//   - Profil: username (tidak dapat diubah sesudah dibuat), nama, Organisasi →
//     Divisi → Unit berjenjang dari master aktif, kotak centang workbasket dan
//     menu.
//   - Security (permintaan work owner 01-10-2026): password - WAJIB saat
//     tambah, kosong = tidak diganti saat ubah - dan centang "Change Password
//     Next Login" (dicentang = user wajib ganti password saat login
//     berikutnya). Password akun SENDIRI boleh diganti; sesi ini tetap jalan.
//
// Satu tombol Simpan untuk kedua tab: saat ubah, profil disimpan dulu, lalu
// Security bila ada yang berubah.

import { useEffect, useState } from 'react'

import { Field, Gagal, Memuat, Modal, Pilih, StripTab, type Opsi } from '../components/ui/dasar'
import { KELOLA_USER, KETERANGAN_BELUM_DIMIGRASI } from '../labels'
import { KODE_MENU_KELOLA_USER } from '../lib/daftarMenu'
import {
  ambilPengguna,
  ambilPilihanPengguna,
  aturSandiPengguna,
  buatPengguna,
  ubahPengguna,
  type OpsiMaster,
  type PilihanKelola,
  type RingkasAkun,
  type RinciAkun,
} from './api'
import {
  alihkan,
  badanBaru,
  badanSandi,
  badanUbah,
  divisiUntuk,
  isianDari,
  isianKosong,
  kelompokMenu,
  periksaGanda,
  periksaIsian,
  setelDivisi,
  setelOrganisasi,
  tabGalat,
  unitUntuk,
  type IsianForm,
  type TabForm,
} from './aturan'

const keOpsi = (daftar: readonly OpsiMaster[]): Opsi[] =>
  daftar.map((o) => ({ value: o.kode, label: o.nama === '' ? o.kode : `${o.nama} (${o.kode})` }))

const TAB: readonly TabForm[] = ['profil', 'security']
const labelTab = (t: TabForm) => (t === 'profil' ? KELOLA_USER.tabProfil : KELOLA_USER.tabSecurity)

/** Kotak sandi tanpa isian otomatis peramban: sandi admin yang tersimpan tidak boleh masuk ke akun lain. */
function IsianSandiBaru({
  id,
  label,
  value,
  onChange,
  wajib,
}: {
  id: string
  label: string
  value: string
  onChange: (v: string) => void
  wajib: boolean
}) {
  return (
    <div className="field">
      <label className="field__label" htmlFor={id}>
        {label}
        {wajib && <span className="field__req">*</span>}
      </label>
      <input
        className="field__input"
        id={id}
        type="password"
        autoComplete="new-password"
        value={value}
        onChange={(e) => {
          onChange(e.target.value)
        }}
      />
    </div>
  )
}

export default function FormPengguna({
  akunId,
  akunSaya,
  daftar,
  onTutup,
  onTersimpan,
}: {
  /** `null` = tambah user. */
  akunId: string | null
  /** Akun yang sedang login - Kelola User miliknya tidak dapat dicabut. */
  akunSaya: string
  /** Daftar user yang sudah termuat - cek username/email sudah terdaftar sebelum kirim; `null` = belum termuat. */
  daftar: readonly RingkasAkun[] | null
  onTutup: () => void
  onTersimpan: (r: RinciAkun) => void
}) {
  const baru = akunId === null
  const [pilihan, setPilihan] = useState<PilihanKelola | null>(null)
  const [isi, setIsi] = useState<IsianForm>(isianKosong)
  // Contact ID hanya ditampilkan (migrasi 905): diberi backend saat user dibuat, tidak dapat diubah.
  const [idKontak, setIdKontak] = useState('')
  // Centang wajib ganti saat dibuka - Security hanya dikirim bila berubah.
  const [wajibAwal, setWajibAwal] = useState(true)
  const [tab, setTab] = useState<TabForm>('profil')
  const [galatMuat, setGalatMuat] = useState<unknown>(null)
  const [galatLokal, setGalatLokal] = useState<string | null>(null)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    let hidup = true
    Promise.all([ambilPilihanPengguna(), akunId === null ? Promise.resolve(null) : ambilPengguna(akunId)]).then(
      ([p, r]) => {
        if (!hidup) return
        setPilihan(p)
        if (r !== null) {
          setIsi(isianDari(r))
          setIdKontak(r.contactId)
          setWajibAwal(r.wajibGantiSandi)
        }
      },
      (g: unknown) => {
        if (hidup) setGalatMuat(g)
      },
    )
    return () => {
      hidup = false
    }
  }, [akunId])

  const diriSendiri = !baru && akunId === akunSaya
  const ubah = (sebagian: Partial<IsianForm>) => {
    setIsi((x) => ({ ...x, ...sebagian }))
  }

  const kirim = () => {
    if (sibuk || pilihan === null) return
    const awal = periksaIsian(isi, baru, akunSaya) ?? periksaGanda(isi, akunId, daftar)
    setGalatLokal(awal)
    setGalatSimpan(null)
    if (awal !== null) {
      setTab(tabGalat(awal))
      return
    }
    setSibuk(true)
    const janji =
      akunId === null
        ? buatPengguna(badanBaru(isi))
        : ubahPengguna(akunId, badanUbah(isi)).then((r) => {
            const sandi = badanSandi(isi, wajibAwal)
            return sandi === null ? r : aturSandiPengguna(akunId, sandi)
          })
    janji.then(
      (r) => {
        setSibuk(false)
        onTersimpan(r)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalatSimpan(g)
      },
    )
  }

  return (
    <Modal
      judul={baru ? KELOLA_USER.judulBaru : `${KELOLA_USER.judulUbah} — ${akunId}`}
      onTutup={onTutup}
      onKirim={kirim}
      labelBatal={KELOLA_USER.batal}
      lebar
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk || pilihan === null}>
          {sibuk ? KELOLA_USER.menyimpan : KELOLA_USER.simpan}
        </button>
      }
    >
      {galatMuat !== null && <Gagal galat={galatMuat} />}
      {pilihan === null && galatMuat === null && <Memuat pesan={KELOLA_USER.memuatPilihan} />}
      {pilihan !== null && (
        <div className="kelola-user__form">
          <StripTab tab={TAB} aktif={tab} onPilih={setTab} label={labelTab} />
          {galatLokal !== null && (
            <div className="alert alert--error" role="alert">
              {galatLokal}
            </div>
          )}
          {galatSimpan !== null && <Gagal galat={galatSimpan} />}

          {tab === 'profil' && (
            <>
              <div className="form-grid">
                <div>
                  <Field
                    label={KELOLA_USER.akun}
                    value={isi.akunId}
                    onChange={(v) => {
                      ubah({ akunId: v })
                    }}
                    required={baru}
                    readOnly={!baru}
                    autoFocus={baru}
                  />
                  {!baru && <p className="muted kelola-user__catatan">{KELOLA_USER.akunTetap}</p>}
                </div>
                {!baru && (
                  <div>
                    <Field label={KELOLA_USER.contactId} value={idKontak} onChange={() => undefined} readOnly />
                    <p className="muted kelola-user__catatan">{KELOLA_USER.catatanContactId}</p>
                  </div>
                )}
                <Field
                  label={KELOLA_USER.nama}
                  value={isi.nama}
                  onChange={(v) => {
                    ubah({ nama: v })
                  }}
                  required
                  autoFocus={!baru}
                />
                {/* Kontak akun (Kelola User 03-10-2026) - opsional, label berbahasa Inggris. */}
                <Field
                  label={KELOLA_USER.email}
                  type="email"
                  value={isi.email}
                  placeholder={KELOLA_USER.contohEmail}
                  onChange={(v) => {
                    ubah({ email: v })
                  }}
                />
                <Field
                  label={KELOLA_USER.telepon}
                  type="tel"
                  value={isi.telepon}
                  placeholder={KELOLA_USER.contohTelepon}
                  onChange={(v) => {
                    ubah({ telepon: v })
                  }}
                />
                <Field
                  label={KELOLA_USER.nik}
                  value={isi.nik}
                  onChange={(v) => {
                    ubah({ nik: v })
                  }}
                />
                <Field
                  label={KELOLA_USER.jabatan}
                  value={isi.jabatan}
                  onChange={(v) => {
                    ubah({ jabatan: v })
                  }}
                />
                <Pilih
                  label={KELOLA_USER.organisasi}
                  value={isi.organisasi}
                  opsi={keOpsi(pilihan.organisasi)}
                  kosong={KELOLA_USER.tidakDiisi}
                  onChange={(v) => {
                    setIsi((x) => setelOrganisasi(x, v, pilihan))
                  }}
                />
                <Pilih
                  label={KELOLA_USER.divisi}
                  value={isi.divisi}
                  opsi={keOpsi(divisiUntuk(pilihan, isi.organisasi))}
                  kosong={KELOLA_USER.tidakDiisi}
                  onChange={(v) => {
                    setIsi((x) => setelDivisi(x, v, pilihan))
                  }}
                />
                <Pilih
                  label={KELOLA_USER.unit}
                  value={isi.unit}
                  opsi={keOpsi(unitUntuk(pilihan, isi.divisi))}
                  kosong={KELOLA_USER.tidakDiisi}
                  onChange={(v) => {
                    ubah({ unit: v })
                  }}
                />
              </div>

              <fieldset className="kelola-user__centang">
                <legend>
                  {KELOLA_USER.workbasket}{' '}
                  <span className="muted">{KELOLA_USER.dipilih(isi.workbasket.length, pilihan.workbasket.length)}</span>
                </legend>
                <div className="kelola-user__aksi-centang">
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      ubah({ workbasket: pilihan.workbasket.map((w) => w.kode).sort() })
                    }}
                  >
                    {KELOLA_USER.pilihSemua}
                  </button>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      ubah({ workbasket: [] })
                    }}
                  >
                    {KELOLA_USER.kosongkan}
                  </button>
                </div>
                <div className="kelola-user__kisi">
                  {pilihan.workbasket.map((w) => {
                    const dipilih = isi.workbasket.includes(w.kode)
                    // Nama di atas, kode workbasket kecil di bawahnya (tampilan sama dengan kartu menu, 05-10-2026).
                    const nama = w.nama !== '' && w.nama !== w.kode ? w.nama : ''
                    return (
                      <div
                        key={w.kode}
                        className={dipilih ? 'kelola-user__butir-menu kelola-user__butir-menu--dipilih' : 'kelola-user__butir-menu'}
                      >
                        <label className="kelola-user__butir">
                          <input
                            type="checkbox"
                            checked={dipilih}
                            onChange={(e) => {
                              ubah({ workbasket: alihkan(isi.workbasket, w.kode, e.target.checked) })
                            }}
                          />
                          <span className="kelola-user__nama-wb">
                            <span>{nama === '' ? w.kode : nama}</span>
                            {nama !== '' && <span className="kelola-user__kode-wb">{w.kode}</span>}
                          </span>
                        </label>
                      </div>
                    )
                  })}
                </div>
              </fieldset>

              <fieldset className="kelola-user__centang">
                <legend>
                  {KELOLA_USER.menu}{' '}
                  <span className="muted">{KELOLA_USER.dipilih(isi.menu.length, pilihan.menu.length)}</span>
                </legend>
                <div className="kelola-user__aksi-centang">
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      ubah({ menu: pilihan.menu.map((m) => m.kode).sort() })
                    }}
                  >
                    {KELOLA_USER.pilihSemua}
                  </button>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      ubah({ menu: diriSendiri ? [KODE_MENU_KELOLA_USER] : [] })
                    }}
                  >
                    {KELOLA_USER.kosongkan}
                  </button>
                </div>
                {kelompokMenu(pilihan.menu).map((g) => (
                  <div key={g.golongan} className="kelola-user__golongan">
                    <p className="kelola-user__golongan-judul">{g.golongan}</p>
                    <div className="kelola-user__kisi">
                      {g.menu.map((m) => {
                        const terkunci = diriSendiri && m.kode === KODE_MENU_KELOLA_USER
                        const dipilih = isi.menu.includes(m.kode)
                        const lihat = isi.menuLihat.includes(m.kode)
                        return (
                          <div
                            key={m.kode}
                            className={dipilih ? 'kelola-user__butir-menu kelola-user__butir-menu--dipilih' : 'kelola-user__butir-menu'}
                          >
                            <label className="kelola-user__butir" title={terkunci ? KELOLA_USER.menuDiriSendiri : undefined}>
                              <input
                                type="checkbox"
                                checked={dipilih}
                                disabled={terkunci}
                                onChange={(e) => {
                                  ubah({ menu: alihkan(isi.menu, m.kode, e.target.checked) })
                                }}
                              />
                              <span className="kelola-user__nama-menu">
                                {m.label}
                                {!m.dimigrasi && <span className="kelola-user__tanda">{KETERANGAN_BELUM_DIMIGRASI}</span>}
                              </span>
                            </label>
                            {/* Full / View only (migrasi 914): hanya menu modul yang mendaftar; aktif bila menunya dipilih. */}
                            {m.bisaLihat === true && (
                              <div
                                className="kelola-user__hak"
                                role="radiogroup"
                                aria-label={KELOLA_USER.hakMenu(m.label)}
                                title={dipilih ? KELOLA_USER.catatanHak : KELOLA_USER.hakPilihDulu}
                              >
                                <button
                                  type="button"
                                  role="radio"
                                  aria-checked={!lihat}
                                  disabled={!dipilih}
                                  className={!lihat ? 'kelola-user__hak-pilih kelola-user__hak-pilih--aktif' : 'kelola-user__hak-pilih'}
                                  onClick={() => {
                                    ubah({ menuLihat: alihkan(isi.menuLihat, m.kode, false) })
                                  }}
                                >
                                  {KELOLA_USER.hakPenuh}
                                </button>
                                <button
                                  type="button"
                                  role="radio"
                                  aria-checked={lihat}
                                  disabled={!dipilih}
                                  className={
                                    lihat
                                      ? 'kelola-user__hak-pilih kelola-user__hak-pilih--lihat kelola-user__hak-pilih--aktif'
                                      : 'kelola-user__hak-pilih kelola-user__hak-pilih--lihat'
                                  }
                                  onClick={() => {
                                    ubah({ menuLihat: alihkan(isi.menuLihat, m.kode, true) })
                                  }}
                                >
                                  {KELOLA_USER.hakLihat}
                                </button>
                              </div>
                            )}
                          </div>
                        )
                      })}
                    </div>
                  </div>
                ))}
                {diriSendiri && <p className="muted kelola-user__catatan">{KELOLA_USER.menuDiriSendiri}</p>}
                {pilihan.menu.some((m) => m.bisaLihat === true) && (
                  <p className="muted kelola-user__catatan">{KELOLA_USER.catatanHak}</p>
                )}
              </fieldset>
            </>
          )}

          {tab === 'security' && (
            <div className="form-grid">
              <IsianSandiBaru
                id="kelola-user-sandi"
                label={baru ? KELOLA_USER.sandi : KELOLA_USER.sandiBaru}
                value={isi.sandi}
                wajib={baru}
                onChange={(v) => {
                  ubah({ sandi: v })
                }}
              />
              <IsianSandiBaru
                id="kelola-user-ulangi"
                label={baru ? KELOLA_USER.ulangiSandi : KELOLA_USER.ulangiSandiBaru}
                value={isi.ulangiSandi}
                wajib={baru}
                onChange={(v) => {
                  ubah({ ulangiSandi: v })
                }}
              />
              <p className="muted kelola-user__catatan field--lebar">
                {baru ? KELOLA_USER.catatanSandi : KELOLA_USER.catatanSandiUbah}
                {diriSendiri && ` ${KELOLA_USER.catatanSandiSendiri}`}
              </p>
              <label className="kelola-user__butir kelola-user__wajib field--lebar">
                <input
                  type="checkbox"
                  checked={isi.wajibGanti}
                  onChange={(e) => {
                    ubah({ wajibGanti: e.target.checked })
                  }}
                />
                <span>
                  <strong>{KELOLA_USER.wajibGantiCentang}</strong>
                  <span className="muted"> — {KELOLA_USER.catatanWajibGanti}</span>
                </span>
              </label>
            </div>
          )}
        </div>
      )}
    </Modal>
  )
}
