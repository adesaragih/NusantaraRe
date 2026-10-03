// Beranda — layar awal sesudah identitas, butir **bg**.
//
// ⚠️ `[kerangka aplikasi, bukan menu Pega]`. Beranda adalah pengganti layar
// awal portal: PremiumList Life punya `PremiumLife_harness` *(kelas
// `Data-Portal`)* sebagai layar awalnya, sedangkan **Claim Life tidak punya
// harness portal yang terekspor**. Bentuk Beranda ini karena itu keputusan
// kami, dan ditandai begitu — bukan tiruan sebuah rule.
//
// ⛔ ANGKANYA DARI ENDPOINT YANG SUDAH ADA, bukan dari rute baru. Kotak masuk
// Claim Life mengembalikan `total` per tahap (`HalamanInbox.total`), jadi
// keempat kartunya dihitung dari sana. PremiumList dan Komite BELUM punya
// kotak masuk — kartunya menyatakan itu, bukan menampilkan nol.
//
// ⚠️ Nol adalah angka; "belum ada kotak masuk" adalah keadaan. Menampilkan
// nol untuk modul yang endpointnya belum ada berarti berbohong dengan angka
// yang terlihat benar.
//
// ⚠️ TATA LETAK mengikuti `workpage-template.html` (29-09-2026): kepala
// halaman, empat kartu ringkasan, lalu kisi dua kolom — tabel modul dan
// kartu peran. Kartu ringkasan berdiri SEJAK AWAL dengan tanda "—": kartu
// yang baru muncul sesudah data tiba menggeser seluruh layar ke bawah.

import { useEffect, useState, type ReactNode } from 'react'

import { BERANDA, KETERANGAN_BELUM_DIMIGRASI, PERAN_ID } from '../inti/frontend/labels'
import { FOLDER_KORPUS, LABEL_MENU } from './katalogKorpus'
import { TAHAP } from '../modul/claimlife/frontend/labels'
import {
  IkonBerkasCari,
  IkonJamPasir,
  IkonKotakMasuk,
  IkonPerisai,
  IkonStetoskop,
} from '../inti/frontend/components/ui/dasar'
import { ambilKotakMasuk, TAHAP_NOMOR, type NomorTahap } from '../modul/claimlife/frontend/api'
import { pesanGalat } from '../inti/frontend/klien'
import { type Sesi } from '../inti/frontend/store/sesi'
import { modulDipasang } from '../inti/frontend/lib/daftarMenu'
import { NAMA_CLAIMLIFE } from '../modul/claimlife/frontend/menu'
import { ENTRI_MENU, halamanAktif, type Halaman } from './daftar'

/** Satu tahap Claim Life beserta cacah antreannya. */
export interface AntreanTahap {
  nomor: NomorTahap
  nama: string
  total: number
}

/** Keadaan sebuah kartu modul. */
export interface KartuModul {
  nama: string
  /** Halaman awal modul itu (`HALAMAN_AWAL_<X>`), atau null bila belum dimigrasi. */
  tujuan: Halaman | null
  label: string | null
}

/**
 * Menyusun kartu untuk ketujuh belas modul.
 *
 * ⛔ Diturunkan dari `ENTRI_MENU`, sumber yang SAMA dengan sidebar dan palet.
 * Daftar keempat yang menyebut modul yang sama adalah daftar keempat yang
 * akan menyimpang.
 *
 * `aktif` - modul dari `GET /api/modul-aktif` (MODUL_AKTIF, refactor bentuk
 * B). Tombol kartu MEMBUKA modul, jadi ia menu juga: kartu modul yang
 * NONAKTIF tidak tampil, seperti kelompoknya di sidebar. `null` = semua.
 */
export function kartuModul(aktif: readonly string[] | null = null): KartuModul[] {
  // Nama menu (nama tampilan `LABEL_TAMPIL` bila diputuskan) - sama dengan `kelompok` modul dan LABEL sidebar.
  return Object.values(LABEL_MENU).flatMap((nama) => {
    const milik = ENTRI_MENU.filter((e) => e.kelompok === nama)
    const pertama = milik.find((e) => halamanAktif(e.modul, aktif))
    if (milik.length > 0 && pertama === undefined) return []
    return [
      {
        nama,
        tujuan: pertama?.modul ?? null,
        label: pertama?.label ?? null,
      },
    ]
  })
}

/** Menyusun kalimat cacah antrean satu modul. */
export function ringkasanAntrean(antrean: AntreanTahap[] | null): string {
  if (antrean === null) return BERANDA.tanpaAntrean
  const jumlah = antrean.reduce((n, a) => n + a.total, 0)
  return `${jumlah} ${BERANDA.antrean}`
}

/** Keempat tahap Claim Life — urutan kartu SAMA dengan urutan permintaan. */
const TAHAP_BERANDA: readonly { nomor: NomorTahap; nama: string; ikon: ReactNode }[] = [
  { nomor: TAHAP_NOMOR.inputRegister, nama: TAHAP.inputRegister, ikon: <IkonKotakMasuk /> },
  { nomor: TAHAP_NOMOR.outstanding, nama: TAHAP.outstandingClaim, ikon: <IkonJamPasir /> },
  { nomor: TAHAP_NOMOR.medicalCheck, nama: TAHAP.medicalCheck, ikon: <IkonStetoskop /> },
  { nomor: TAHAP_NOMOR.claimAnalis, nama: TAHAP.claimAnalis, ikon: <IkonBerkasCari /> },
]

export default function Beranda({
  masuk,
  onBuka,
  modulAktif = null,
}: {
  masuk: Sesi
  onBuka: (modul: Halaman) => void
  /** Modul aktif dari backend; `null` = semua (lihat `halamanAktif`). */
  modulAktif?: readonly string[] | null
}) {
  const [antrean, setAntrean] = useState<AntreanTahap[] | null>(null)
  const [galat, setGalat] = useState<string | null>(null)
  // Cacah antrean milik Claim Life: modul itu NONAKTIF = kartunya tidak
  // tampil dan kotak masuknya tidak diminta (rutenya memang tidak ada).
  const claimLifeAktif = modulDipasang(NAMA_CLAIMLIFE, modulAktif)

  useEffect(() => {
    if (!claimLifeAktif) {
      setAntrean(null)
      setGalat(null)
      return
    }
    let hidup = true
    void (async () => {
      try {
        // Keempat tahap Claim Life. Dimuat berbarengan: empat permintaan
        // berurutan membuat Beranda terasa lambat tanpa sebab.
        const hasil = await Promise.all(
          TAHAP_BERANDA.map(async ({ nomor, nama }) => {
            const h = await ambilKotakMasuk(nomor, 1, 1)
            return { nomor, nama: h.namaTahap || nama, total: h.total }
          }),
        )
        if (hidup) setAntrean(hasil)
      } catch (e) {
        // ⛔ Galat DINYATAKAN, bukan menjadi daftar kosong. Beranda yang
        // diam sesudah gagal memuat terbaca "tidak ada pekerjaan".
        if (hidup) setGalat(pesanGalat(e) ?? 'Cacah antrean gagal dimuat.')
      }
    })()
    return () => {
      hidup = false
    }
  }, [claimLifeAktif])

  const kartu = kartuModul(modulAktif)
  const jumlahAktif = kartu.filter((k) => k.tujuan !== null).length
  const memuat = antrean === null && galat === null

  /** Isi kolom Antrean satu modul. "—" = tidak ada angka untuk ditampilkan. */
  function antreanModul(k: KartuModul): string {
    if (k.tujuan === null) return '—'
    if (k.nama !== FOLDER_KORPUS.claimLife) return BERANDA.tanpaAntrean
    return antrean === null ? '—' : ringkasanAntrean(antrean)
  }

  return (
    <div className="beranda">
      <section className="kepala-halaman">
        <div>
          <h2 className="beranda__judul">
            {BERANDA.salam}, {masuk.akunID}
          </h2>
          <p className="beranda__peran">{BERANDA.subjudul}</p>
        </div>
      </section>

      {galat !== null && (
        <p className="alert alert--error" role="alert">
          {galat}
        </p>
      )}

      {claimLifeAktif && (
        <section className="beranda__antrean" aria-label={BERANDA.ringkasan} aria-busy={memuat}>
          {TAHAP_BERANDA.map((t) => {
            const a = antrean?.find((x) => x.nomor === t.nomor)
            return (
              <div key={t.nomor} className="kartu beranda__antrean-butir">
                <p className="beranda__antrean-nama">
                  {t.ikon}
                  {a?.nama ?? t.nama}
                </p>
                <p className="beranda__antrean-angka">{a === undefined ? '—' : a.total}</p>
                <p className="beranda__antrean-catatan">{BERANDA.catatanTahap}</p>
              </div>
            )
          })}
        </section>
      )}

      <div className="beranda__kisi">
        <section className="kartu" aria-labelledby="beranda-modul">
          <div className="kartu__kepala">
            <div>
              <h3 id="beranda-modul">{BERANDA.judulModul}</h3>
              <p>
                {jumlahAktif} {BERANDA.aktif} · {kartu.length - jumlahAktif}{' '}
                {KETERANGAN_BELUM_DIMIGRASI}
              </p>
            </div>
          </div>
          <div className="kartu__tabel">
            <table className="beranda__kartu">
              <thead>
                <tr>
                  <th scope="col">{BERANDA.kolomModul}</th>
                  <th scope="col">{BERANDA.kolomStatus}</th>
                  <th scope="col" className="sembunyi-ponsel">
                    {BERANDA.kolomAntrean}
                  </th>
                  <th scope="col">
                    <span className="sr-only">{BERANDA.kolomAksi}</span>
                  </th>
                </tr>
              </thead>
              <tbody>
                {kartu.map((k) => (
                  <tr
                    key={k.nama}
                    className={`beranda__kartu-butir${
                      k.tujuan === null ? ' beranda__kartu-butir--pasif' : ''
                    }`}
                  >
                    <td>
                      <span className="beranda__kartu-nama">{k.nama}</span>
                      {/* Di ponsel kolom Antrean disembunyikan — angkanya
                          pindah ke bawah nama, seperti tenggat di template. */}
                      {k.tujuan !== null && (
                        <span className="beranda__kartu-keadaan tampil-ponsel">
                          {antreanModul(k)}
                        </span>
                      )}
                    </td>
                    <td>
                      {/* ⛔ Menyebut keadaannya, bukan menyembunyikan barisnya. */}
                      <span
                        className={`status ${k.tujuan === null ? 'status--pasif' : 'status--aktif'}`}
                      >
                        {k.tujuan === null ? KETERANGAN_BELUM_DIMIGRASI : BERANDA.aktif}
                      </span>
                    </td>
                    <td className="sembunyi-ponsel">{antreanModul(k)}</td>
                    <td className="beranda__aksi">
                      {k.tujuan !== null && k.label !== null && (
                        <button
                          type="button"
                          className="beranda__kartu-tautan"
                          onClick={() => {
                            onBuka(k.tujuan as Halaman)
                          }}
                        >
                          {k.label}
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        <section className="kartu" aria-labelledby="beranda-peran">
          <div className="kartu__kepala">
            <h3 id="beranda-peran">{BERANDA.judulPeran}</h3>
          </div>
          <ul className="beranda__peran-daftar">
            {masuk.peran.map((p) => (
              <li key={p}>
                <span className="beranda__peran-ikon" aria-hidden="true">
                  <IkonPerisai />
                </span>
                <span>
                  <strong>{PERAN_ID[p]}</strong>
                  <span>{p}</span>
                </span>
              </li>
            ))}
          </ul>
        </section>
      </div>
    </div>
  )
}
