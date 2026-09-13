# Contributing

Terima kasih sudah tertarik berkontribusi! Proyek ini adalah compiler
pseudocode (gaya notasi algoritmik kuliah Indonesia) ke Go, lengkap
dengan backend API dan frontend Angular. Dokumen ini menjelaskan cara
mulai berkontribusi.

## Struktur repo

```
.
├── LICENSE
├── docs/
│   └── ROADMAP.md      # rencana pengembangan, termasuk task OCR yang belum digarap
├── backend/            # Go: lexer, parser, codegen, HTTP API
└── frontend/           # Angular
```

## Mulai dari mana?

- Cek [docs/ROADMAP.md](docs/ROADMAP.md) — ada breakdown task konkret
  di sana, termasuk beberapa yang cocok untuk kontribusi pertama
  (ditandai dengan checklist, akan dipetakan ke GitHub Issues berlabel
  `good first issue`).
- Kalau menemukan pseudocode valid yang gagal di-compile, atau
  pseudocode salah yang malah lolos, itu juga kontribusi berharga —
  laporkan lewat Issue dengan contoh `.pseudo`-nya.

## Setup development

### Backend

```bash
cd backend
go build ./...
go test ./...
```

Semua perubahan di `internal/lexer`, `internal/parser`, atau
`internal/codegen` **wajib disertai test** — proyek ini sengaja disiplin
soal ini karena grammar-nya berkembang dari kasus nyata (termasuk
tulisan tangan mahasiswa asli), bukan spekulasi. Kalau menambah
dukungan sintaks baru, tambahkan juga file `.pseudo` di `testdata/`
yang merepresentasikan kasus tersebut, lalu pastikan `go test ./...`
tetap hijau untuk seluruh testdata yang sudah ada.

### Frontend

```bash
cd frontend
npm install
ng serve 
ng test    
```

## Menambah dukungan sintaks pseudocode baru

Proyek ini mendukung beberapa variasi penulisan yang secara historis
muncul dari cara mahasiswa sungguhan menulis pseudocode (misalnya
`endif` maupun `end if`, `=` maupun `==`). Kalau menambah toleransi
serupa:

1. Sertakan **sumber/alasan konkret** di deskripsi PR — idealnya
   contoh nyata yang memicunya, bukan asumsi.
2. Pastikan perubahan tidak memecah gaya yang sudah didukung
   sebelumnya (jalankan regresi penuh: `go test ./...`).
3. Kalau perubahan mengubah AST atau semantik codegen (bukan cuma
   lexer), jelaskan di PR bagaimana itu memengaruhi kode Go yang
   dihasilkan.

## Alur kontribusi

1. Fork & buat branch dari `main`: `feature/nama-singkat` atau
   `fix/nama-singkat`.
2. Commit dengan pesan yang jelas menyatakan **apa** dan **kenapa**,
   bukan cuma "update code".
3. Buka Pull Request, tautkan ke Issue terkait kalau ada.
4. Pastikan `go build ./...` dan `go test ./...` lolos sebelum minta
   review.

## Bahasa

Pesan error compiler dan sebagian komentar kode ditulis dalam Bahasa
Indonesia (menyesuaikan audiens target: mahasiswa). Konsistensi ini
tolong dijaga di kontribusi baru — kalau menambah pesan error baru,
tulis dalam Bahasa Indonesia mengikuti gaya yang sudah ada.

## Lisensi

Dengan berkontribusi, Anda menyetujui kontribusi tersebut dilisensikan
di bawah [MIT License](LICENSE) yang sama dengan proyek ini.

## Etika berkontribusi

Bersikap baik dan konstruktif terhadap sesama kontributor. Perbedaan
pendapat teknis itu wajar — sampaikan dengan hormat dan fokus pada
masalahnya, bukan orangnya.