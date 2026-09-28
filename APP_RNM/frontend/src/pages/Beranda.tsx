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

import { useEffect, useState } from 'react'

import { BERANDA, KETERANGAN_BELUM_DIMIGRASI, MENU, MODUL } from '../assets/labels'
import { ENTRI_MENU, type ModulTetap } from '../lib/daftarMenu'
import {
  ambilKotakMasuk,
  pesanGalat,
  TAHAP_NOMOR,
  type NomorTahap,
} from '../services/api'
import { type Sesi } from '../store/sesi'

/** Satu tahap Claim Life beserta cacah antreannya. */
export interface AntreanTahap {
  nomor: NomorTahap
  nama: string
  total: number
}

/** Keadaan sebuah kartu modul. */
export interface KartuModul {
  nama: string
  /** Butir pertama modul itu, atau null bila belum dimigrasi. */
  tujuan: ModulTetap | null
  label: string | null
}

/**
 * Menyusun kartu untuk ketujuh belas modul.
 *
 * ⛔ Diturunkan dari `ENTRI_MENU`, sumber yang SAMA dengan sidebar dan palet.
 * Daftar keempat yang menyebut modul yang sama adalah daftar keempat yang
 * akan menyimpang.
 */
export function kartuModul(): KartuModul[] {
  return Object.values(MODUL).map((nama) => {
    const pertama = ENTRI_MENU.find((e) => e.kelompok === nama)
    return {
      nama,
      tujuan: pertama?.modul ?? null,
      label: pertama?.label ?? null,
    }
  })
}

/** Menyusun kalimat cacah antrean satu modul. */
export function ringkasanAntrean(antrean: AntreanTahap[] | null): string {
  if (antrean === null) return BERANDA.tanpaAntrean
  const jumlah = antrean.reduce((n, a) => n + a.total, 0)
  return `${jumlah} ${BERANDA.antrean}`
}

export default function Beranda({
  masuk,
  onBuka,
}: {
  masuk: Sesi
  onBuka: (modul: ModulTetap) => void
}) {
  const [antrean, setAntrean] = useState<AntreanTahap[] | null>(null)
  const [galat, setGalat] = useState<string | null>(null)

  useEffect(() => {
    let hidup = true
    void (async () => {
      try {
        // Keempat tahap Claim Life. Dimuat berbarengan: empat permintaan
        // berurutan membuat Beranda terasa lambat tanpa sebab.
        const hasil = await Promise.all(
          (
            [
              [TAHAP_NOMOR.inputRegister, 'Input Register'],
              [TAHAP_NOMOR.outstanding, 'Outstanding Claim'],
              [TAHAP_NOMOR.medicalCheck, 'Medical Check'],
              [TAHAP_NOMOR.claimAnalis, 'Claim Analis'],
            ] as const
          ).map(async ([nomor, nama]) => {
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
  }, [])

  const kartu = kartuModul()

  return (
    <section className="beranda">
      <h2 className="beranda__judul">
        {BERANDA.salam}, {masuk.akunID}
      </h2>
      <p className="beranda__peran">{masuk.peran.join(', ')}</p>

      {galat !== null && <p role="alert">{galat}</p>}

      {antrean !== null && (
        <ul className="beranda__antrean">
          {antrean.map((a) => (
            <li key={a.nomor} className="beranda__antrean-butir">
              <span className="beranda__antrean-nama">{a.nama}</span>
              <span className="beranda__antrean-angka">{a.total}</span>
            </li>
          ))}
        </ul>
      )}

      <ul className="beranda__kartu">
        {kartu.map((k) => (
          <li
            key={k.nama}
            className={`beranda__kartu-butir${
              k.tujuan === null ? ' beranda__kartu-butir--pasif' : ''
            }`}
          >
            <h3 className="beranda__kartu-nama">{k.nama}</h3>
            {k.tujuan === null || k.label === null ? (
              /* ⛔ Menyebut keadaannya, bukan menyembunyikan kartunya. */
              <p className="beranda__kartu-keadaan">{KETERANGAN_BELUM_DIMIGRASI}</p>
            ) : (
              <>
                <p className="beranda__kartu-keadaan">
                  {BERANDA.aktif}
                  {k.nama === MODUL.claimLife && ` — ${ringkasanAntrean(antrean)}`}
                  {k.nama !== MODUL.claimLife && ` — ${BERANDA.tanpaAntrean}`}
                </p>
                <button
                  type="button"
                  className="beranda__kartu-tautan"
                  onClick={() => {
                    onBuka(k.tujuan as ModulTetap)
                  }}
                >
                  {k.label === MENU.inbox ? MENU.inbox : k.label}
                </button>
              </>
            )}
          </li>
        ))}
      </ul>
    </section>
  )
}
